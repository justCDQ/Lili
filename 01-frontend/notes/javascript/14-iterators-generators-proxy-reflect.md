# JavaScript 迭代器、生成器、Proxy 与 Reflect

迭代协议把“数据如何存储”和“数据如何逐项消费”分开；生成器用可暂停函数实现迭代器；`Proxy` 拦截对象的基础操作；`Reflect` 以函数形式执行相应默认操作。这四部分共同提供了 JavaScript 的控制抽象与元编程基础。

## 1. 迭代器协议

迭代器是包含 `next()` 方法的对象。每次调用 `next()` 都必须返回一个对象：

- `done` 为 `false` 时，`value` 是本次产出的值；
- `done` 为 `true` 时，序列已经结束，`value` 可以表示最终返回值；
- 若省略 `done`，消费方通常把它视为假值。

```js
function createRangeIterator(start, end, step = 1) {
  if (step <= 0) {
    throw new RangeError("step 必须大于 0");
  }

  let current = start;

  return {
    next() {
      if (current >= end) {
        return { value: undefined, done: true };
      }

      const value = current;
      current += step;
      return { value, done: false };
    },
  };
}

const iterator = createRangeIterator(1, 6, 2);
console.log(iterator.next()); // { value: 1, done: false }
console.log(iterator.next()); // { value: 3, done: false }
console.log(iterator.next()); // { value: 5, done: false }
console.log(iterator.next()); // { value: undefined, done: true }
```

迭代器是有状态、一次性向前移动的对象。把同一个迭代器交给两个消费者时，它们共享当前位置。

### 1.1 可迭代协议

可迭代对象必须实现键为 `Symbol.iterator` 的无参数方法。该方法返回迭代器。

```js
const range = {
  start: 2,
  end: 5,

  [Symbol.iterator]() {
    let current = this.start;
    const end = this.end;

    return {
      next() {
        if (current < end) {
          return { value: current++, done: false };
        }
        return { value: undefined, done: true };
      },
    };
  },
};

console.log([...range]); // [2, 3, 4]
console.log([...range]); // [2, 3, 4]
```

这里每次调用 `range[Symbol.iterator]()` 都创建新迭代器，所以可重复遍历。迭代器和可迭代对象可以是同一个对象，但这样通常只能消费一次。

### 1.2 哪些语法消费 iterable

以下操作使用同步可迭代协议：

- `for...of`；
- 数组展开和函数调用展开；
- 数组解构；
- `Array.from()`；
- `new Map(iterable)`、`new Set(iterable)`；
- `Promise.all()` 等 Promise 组合方法接收的输入集合。

```js
const values = new Set(["html", "css", "javascript"]);

const [first, ...rest] = values;
console.log(first); // html
console.log(rest); // ["css", "javascript"]

for (const value of values) {
  console.log(value);
}
```

普通对象默认不可迭代。遍历其键值对时先选择明确视图：

```js
const scores = { html: 80, css: 90 };

for (const [name, score] of Object.entries(scores)) {
  console.log(name, score);
}
```

### 1.3 提前终止与 iterator.return

消费方提前退出 `for...of` 时，会在迭代器存在 `return()` 的情况下调用它，让生产方释放资源。

```js
function createLoggedIterator() {
  let current = 0;

  return {
    next() {
      return { value: current++, done: false };
    },
    return() {
      console.log("迭代已关闭");
      return { value: undefined, done: true };
    },
    [Symbol.iterator]() {
      return this;
    },
  };
}

for (const value of createLoggedIterator()) {
  console.log(value);
  if (value === 2) break;
}
```

`break`、`return` 或循环体抛错都可能触发迭代器关闭。直接手动调用 `iterator.next()` 后停止，不会自动调用 `return()`。

### 1.4 无限序列

迭代器可以按需计算无限序列，因为它不必预先保存所有值。

```js
function* naturalNumbers() {
  let value = 1;
  while (true) {
    yield value++;
  }
}

for (const value of naturalNumbers()) {
  if (value > 3) break;
  console.log(value);
}
```

不要对无限 iterable 使用数组展开、`Array.from()` 或没有退出条件的完整消费，这会持续运行并不断占用内存。

## 2. 生成器函数

`function*` 声明生成器函数。调用生成器函数不会立即执行函数体，而是返回生成器对象。生成器对象同时实现迭代器和可迭代协议。

```js
function* lessonIds() {
  console.log("开始执行");
  yield "js-01";
  yield "js-02";
  return "结束";
}

const lessons = lessonIds();
console.log(lessons.next());
console.log(lessons.next());
console.log(lessons.next());
console.log(lessons.next());
```

执行顺序如下：

1. `lessonIds()` 只创建生成器，尚未打印；
2. 第一次 `next()` 从函数开头运行到第一个 `yield`；
3. 第二次 `next()` 从暂停位置恢复到下一个 `yield`；
4. 第三次恢复并执行 `return`，得到 `{ value: "结束", done: true }`；
5. 已结束的生成器继续 `next()`，得到完成结果。

`for...of` 只消费 `done: false` 的值，不会产出生成器的最终 `return` 值。

### 2.1 yield 的输入与输出

`yield expression` 先向调用方产出 `expression` 的值。下一次调用 `next(input)` 时，`yield` 表达式本身求值为 `input`。

```js
function* questionnaire() {
  const language = yield "首选语言？";
  const level = yield "当前阶段？";
  return { language, level };
}

const questions = questionnaire();
console.log(questions.next());
console.log(questions.next("JavaScript"));
console.log(questions.next("初级"));
```

第一次 `next(argument)` 传入的参数会被忽略，因为此时生成器还没有暂停在某个 `yield` 表达式上。

### 2.2 throw 与 return

生成器对象的三个控制方法是：

- `next(value)`：正常恢复，并把值交给暂停的 `yield`；
- `throw(error)`：在暂停位置抛出错误；
- `return(value)`：请求结束生成器，并执行途经的 `finally` 清理。

```js
function* managedTask() {
  try {
    yield "running";
    yield "still-running";
  } finally {
    console.log("释放任务资源");
  }
}

const task = managedTask();
console.log(task.next());
console.log(task.return("cancelled"));
console.log(task.next());
```

如果生成器内部捕获了通过 `throw()` 注入的错误，可以继续产出；未捕获时错误传播给调用者，生成器结束。

### 2.3 yield* 委托

`yield* iterable` 逐项转发另一个 iterable 的值。被委托生成器完成后，其最终返回值会成为整个 `yield*` 表达式的结果。

```js
function* frontendBasics() {
  yield "HTML";
  yield "CSS";
  return 2;
}

function* roadmap() {
  const count = yield* frontendBasics();
  yield `基础模块数：${count}`;
  yield "JavaScript";
}

console.log([...roadmap()]);
```

`yield*` 适合递归遍历、组合多个序列和把控制方法转发给子生成器。

### 2.4 生成器不是并行线程

生成器只在调用 `next()` 等方法时同步执行，直到再次暂停或结束。它不会自动调度，也不会让 CPU 工作并行。涉及异步值时应使用异步生成器和 `for await...of`，详见异步迭代章节。

## 3. 用 iterable 设计 API

返回 iterable 的 API 可以延迟计算、提前停止，也能被多种语言语法消费。

```js
class NoteCollection {
  #notes;

  constructor(notes) {
    this.#notes = [...notes];
  }

  *[Symbol.iterator]() {
    yield* this.#notes;
  }

  *filterByLevel(level) {
    for (const note of this.#notes) {
      if (note.level === level) {
        yield note;
      }
    }
  }
}

const collection = new NoteCollection([
  { id: "js-01", level: "beginner" },
  { id: "js-14", level: "junior" },
]);

console.log([...collection]);
console.log([...collection.filterByLevel("junior")]);
```

设计时应明确：

- 是否允许重复遍历；
- 迭代期间底层集合变化时采用快照还是实时视图；
- 是否可能无限；
- 提前结束时是否需要释放文件、连接或锁；
- 迭代错误由生产方抛出还是转换为普通结果。

## 4. Proxy：拦截对象操作

`new Proxy(target, handler)` 返回代理对象。对代理执行属性读取、写入、删除、枚举、函数调用、构造等操作时，对应 trap 可以介入。未提供 trap 的操作会转发到目标对象。

```js
const target = { title: "迭代器" };

const proxy = new Proxy(target, {
  get(currentTarget, key, receiver) {
    console.log("读取", String(key));
    return Reflect.get(currentTarget, key, receiver);
  },
});

console.log(proxy.title);
```

常见 trap 与操作：

| trap | 常见触发操作 |
| --- | --- |
| `get` | `proxy.key` |
| `set` | `proxy.key = value` |
| `has` | `key in proxy` |
| `deleteProperty` | `delete proxy.key` |
| `ownKeys` | `Reflect.ownKeys(proxy)`、部分枚举操作 |
| `getOwnPropertyDescriptor` | 属性描述符查询 |
| `defineProperty` | 定义属性 |
| `getPrototypeOf` | 原型查询 |
| `setPrototypeOf` | 修改原型 |
| `isExtensible` | 可扩展性查询 |
| `preventExtensions` | 禁止扩展 |
| `apply` | 调用函数代理 |
| `construct` | 对构造器代理使用 `new` |

### 4.1 用 Reflect 正确转发

`Reflect` 是命名空间对象，方法名称与 Proxy trap 对应。它不是构造器。`Reflect.get()`、`Reflect.set()` 等方法直接执行对象基础操作，并以返回值表达部分成功或失败状态。

```js
const state = { progress: 0 };

const validated = new Proxy(state, {
  set(target, key, value, receiver) {
    if (key === "progress") {
      if (!Number.isFinite(value) || value < 0 || value > 100) {
        throw new RangeError("progress 必须在 0 到 100 之间");
      }
    }
    return Reflect.set(target, key, value, receiver);
  },
});

validated.progress = 60;
console.log(state.progress); // 60
```

`set` trap 必须返回布尔值。严格模式下返回假值会使赋值抛出 `TypeError`。业务校验失败时直接抛出具体错误，通常比静默返回 `false` 更清晰。

`receiver` 很重要：若目标或其原型上存在 getter/setter，`Reflect.get(target, key, receiver)` 和 `Reflect.set(...)` 能保持访问器中的 `this` 指向原始接收者。

### 4.2 Proxy 不变量

Proxy 不能随意伪造与目标对象完整性冲突的结果。引擎会检查关键不变量。例如：

- 不能把目标上不可配置、不可写的数据属性报告成另一个值；
- `ownKeys` 不能遗漏目标的不可配置自有键；
- 目标不可扩展时，`ownKeys` 不能增加或遗漏自有键；
- `getPrototypeOf` 对不可扩展目标必须返回真实原型；
- `construct` trap 必须返回对象。

```js
const target = {};
Object.defineProperty(target, "fixed", {
  value: 1,
  writable: false,
  configurable: false,
});

const invalidProxy = new Proxy(target, {
  get(currentTarget, key, receiver) {
    if (key === "fixed") return 2;
    return Reflect.get(currentTarget, key, receiver);
  },
});

try {
  console.log(invalidProxy.fixed);
} catch (error) {
  console.log(error.name); // TypeError
}
```

trap 返回了 `2`，但目标不可变属性真实值为 `1`，因此引擎拒绝该结果。Proxy 提供拦截能力，不提供绕过对象完整性规则的能力。

### 4.3 身份、内部槽与私有字段边界

代理是与目标不同的对象：

```js
const target = {};
const proxy = new Proxy(target, {});

console.log(proxy === target); // false
const set = new Set([target]);
console.log(set.has(proxy)); // false
```

某些内建对象的方法依赖接收者拥有特定内部槽。把 `Map`、`Set`、`Date` 等直接包在空代理中，再通过代理调用方法，可能因接收者缺少内部槽而抛出 `TypeError`。

```js
const map = new Map();
const proxy = new Proxy(map, {});

try {
  proxy.set("topic", "proxy");
} catch (error) {
  console.log(error.name); // TypeError
}
```

类似地，类方法访问 `#private` 字段时，空代理不自动获得目标的私有品牌。不能把 Proxy 当成所有对象的透明包装。

### 4.4 可撤销代理

`Proxy.revocable(target, handler)` 返回 `{ proxy, revoke }`。调用 `revoke()` 后，几乎所有代理操作都会抛出 `TypeError`，重复撤销无副作用。

```js
const session = Proxy.revocable(
  { scope: "lesson:read" },
  {
    get(target, key, receiver) {
      return Reflect.get(target, key, receiver);
    },
  },
);

console.log(session.proxy.scope);
session.revoke();

try {
  console.log(session.proxy.scope);
} catch (error) {
  console.log(error.name); // TypeError
}
```

它适合有明确生命周期的临时能力对象、插件接口或会话视图。撤销代理不是取消已经启动的异步副作用；异步任务仍需 `AbortSignal` 等取消协议。

## 5. Reflect 的独立用途

Reflect 不只用于 Proxy 转发。它为对象操作提供一致的函数 API：

```js
const note = {};

const defined = Reflect.defineProperty(note, "title", {
  value: "元编程",
  writable: true,
  enumerable: true,
  configurable: true,
});

console.log(defined); // true
console.log(Reflect.has(note, "title")); // true
console.log(Reflect.ownKeys(note)); // ["title"]
console.log(Reflect.deleteProperty(note, "title")); // true
```

与部分 `Object` API 相比，`Reflect.defineProperty()`、`Reflect.deleteProperty()`、`Reflect.preventExtensions()` 等返回布尔结果，便于在控制流中处理失败。`Reflect.apply()` 和 `Reflect.construct()` 可分别执行函数调用和构造调用。

```js
function format(prefix, suffix) {
  return `${prefix}${this.title}${suffix}`;
}

const result = Reflect.apply(
  format,
  { title: "Reflect" },
  ["主题：", "。"],
);
console.log(result);

class Lesson {
  constructor(id) {
    this.id = id;
  }
}

const lesson = Reflect.construct(Lesson, ["js-14"]);
console.log(lesson instanceof Lesson); // true
```

## 6. 完整案例：可迭代、可观察的学习计划

案例目标：计划对象可重复遍历；筛选使用生成器延迟产出；代理校验写入并记录变更；`Reflect` 保留默认语义；撤销后禁止继续修改。

```js
"use strict";

class StudyPlan {
  #items;

  constructor(items) {
    if (!Array.isArray(items)) {
      throw new TypeError("items 必须是数组");
    }
    this.#items = items.map((item) => StudyPlan.#normalize(item));
  }

  static #normalize(item) {
    if (item === null || typeof item !== "object") {
      throw new TypeError("计划项必须是对象");
    }
    if (typeof item.id !== "string" || item.id.trim() === "") {
      throw new TypeError("计划项 id 必须是非空字符串");
    }
    if (!["planned", "learning", "completed"].includes(item.status)) {
      throw new RangeError(`无效状态：${item.status}`);
    }
    return { id: item.id, status: item.status };
  }

  *[Symbol.iterator]() {
    for (const item of this.#items) {
      yield { ...item };
    }
  }

  *byStatus(status) {
    for (const item of this.#items) {
      if (item.status === status) {
        yield { ...item };
      }
    }
  }

  update(id, status) {
    if (!["planned", "learning", "completed"].includes(status)) {
      throw new RangeError(`无效状态：${status}`);
    }
    const item = this.#items.find((candidate) => candidate.id === id);
    if (!item) {
      throw new RangeError(`找不到计划项：${id}`);
    }
    item.status = status;
  }
}

function observePlan(plan, onChange) {
  if (typeof onChange !== "function") {
    throw new TypeError("onChange 必须是函数");
  }

  return Proxy.revocable(plan, {
    get(target, key, receiver) {
      const value = Reflect.get(target, key, receiver);

      if (key !== "update" || typeof value !== "function") {
        return value;
      }

      return function observedUpdate(id, status) {
        const before = [...target];
        const result = Reflect.apply(value, target, [id, status]);
        const after = [...target];
        onChange({ operation: "update", id, status, before, after });
        return result;
      };
    },
  });
}

const plan = new StudyPlan([
  { id: "js-13", status: "completed" },
  { id: "js-14", status: "learning" },
  { id: "js-15", status: "planned" },
]);

const changes = [];
const access = observePlan(plan, (change) => changes.push(change));

console.log([...access.proxy]);
console.log([...access.proxy.byStatus("planned")]);

access.proxy.update("js-14", "completed");
console.log(changes.length); // 1
console.log(changes[0].after);

try {
  access.proxy.update("missing", "completed");
} catch (error) {
  console.log(error.name, error.message);
}

try {
  access.proxy.update("js-15", "invalid");
} catch (error) {
  console.log(error.name, error.message);
}

console.log(changes.length); // 失败操作没有写入变更记录，仍为 1

access.revoke();
try {
  console.log([...access.proxy]);
} catch (error) {
  console.log(error.name); // TypeError
}
```

关键实现理由：

- 迭代时返回条目副本，避免消费者直接改动内部数组对象；
- `byStatus()` 是延迟生成器，消费者可提前结束；
- trap 只包装 `update`，其余读取通过 `Reflect.get()` 转发；
- 调用原始方法时显式把 `target` 作为接收者，使私有字段品牌检查通过；
- 只在更新成功后记录变更，失败分支不会生成虚假审计记录；
- 撤销代理控制访问生命周期，但原始 `plan` 若仍被持有，依然可以使用。

## 7. 常见错误与调试清单

### 7.1 常见错误

1. 让 `next()` 返回普通值，而不是迭代结果对象。
2. 把一次性迭代器当成可重复遍历集合。
3. 对无限序列使用展开或 `Array.from()`。
4. 期望 `for...of` 取得生成器的最终 `return` 值。
5. 忘记第一次 `next(value)` 的 `value` 不会进入生成器。
6. 把生成器当作并行或自动执行的任务。
7. `set`、`deleteProperty` 等 trap 忘记返回布尔值。
8. trap 直接使用 `target[key]`，破坏 getter 的接收者语义。
9. 返回违反不可配置属性约束的结果，引发 `TypeError`。
10. 假设代理与目标身份相同，或代理能透明转发所有内部槽。
11. 在热点循环中无评估地引入 Proxy，增加间接调用和调试成本。
12. 认为撤销代理会自动取消已经发起的网络请求。

### 7.2 调试清单

- 手动连续调用 `next()`，检查每一步的 `value` 和 `done`；
- 检查 `value[Symbol.iterator]()` 是否每次返回新迭代器；
- 用 `try/finally` 和提前 `break` 验证资源清理；
- 检查生成器暂停位置，再判断下次 `next(value)` 的接收点；
- 为无限序列设置消费上限或提前退出条件；
- 在 trap 中记录操作名、键和结果，不记录密码或令牌值；
- 对照 `Reflect` 的同名方法转发默认行为；
- 用 `Object.getOwnPropertyDescriptor()` 检查目标属性不变量；
- 检查方法是否依赖私有字段或内建内部槽；
- 单独测试代理撤销后的读取、写入和方法调用；
- 对性能敏感路径用实际工作负载测量，不依据抽象结论猜测。

## 8. 练习

### 练习一：树的深度优先遍历

实现 `function* walk(node)`，先产出当前节点，再通过 `yield*` 遍历子节点。为空节点和深层数据添加测试，并考虑循环引用的防护。

### 练习二：分页 iterable

实现同步的分页数据 iterable，每页只在需要时从本地数据源读取；提前 `break` 时记录关闭。说明它为何仍不适合真实异步网络请求。

### 练习三：双向生成器

实现一个生成器，依次请求题目答案，通过 `next(answer)` 接收结果，最终返回得分。分别测试正常完成、`throw()` 注入错误和 `return()` 取消。

### 练习四：配置代理

用 Proxy 限制允许写入的配置键，拒绝删除必需键，并通过 `Reflect` 转发。再把必需键定义为不可配置，验证 trap 必须遵守的不变量。

### 练习五：撤销访问

为插件对象创建可撤销代理，只暴露白名单方法。撤销后验证读取和调用失败，并说明为什么仍不能把它当成完整安全沙箱。

## 9. 补充知识

- 字符串、数组、TypedArray、Map、Set 都提供内建同步迭代器；普通对象没有。
- 数组的默认迭代器会读取迭代期间的索引状态；修改集合可能影响后续结果，业务 API 应明确修改规则。
- Proxy 只能拦截语言定义的特定基础操作，局部变量访问等操作不经过 Proxy。
- 代理不能改变目标是否可调用或可构造：只有目标本身可调用时才有调用能力，只有目标本身可构造时才有构造能力。
- `Reflect.ownKeys()` 同时返回字符串键与 Symbol 键，且包含不可枚举自有键。
- Proxy 适合观察、验证、适配和权限生命周期，但不等于安全边界；不可信代码仍需受控执行环境。

## 来源

- [ECMAScript 2026：Control Abstraction Objects](https://tc39.es/ecma262/2026/multipage/control-abstraction-objects.html)（访问日期：2026-07-17）
- [ECMAScript 2026：Ordinary and Exotic Object Behaviours](https://tc39.es/ecma262/2026/multipage/ordinary-and-exotic-objects-behaviours.html)（访问日期：2026-07-17）
- [ECMAScript 2026：Reflection](https://tc39.es/ecma262/2026/multipage/reflection.html)（访问日期：2026-07-17）
- [MDN：Iterators and generators](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Iterators_and_generators)（访问日期：2026-07-17）
- [MDN：Meta programming](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Meta_programming)（访问日期：2026-07-17）
