# JavaScript 解构、展开、模板字符串、可选链与空值合并

这些语法用于读取结构化数据、组合新数据、生成文本，以及安全处理可能缺失的路径。它们减少样板代码，但不会自动完成深复制、数据验证或输出转义。使用时必须知道每段表达式读取了什么、复制到哪一层、何时短路，以及默认值到底针对哪些缺失状态。

## 1. 五种语法的责任边界

| 语法 | 核心作用 | 不会自动完成 |
| --- | --- | --- |
| 解构 | 从可迭代值或对象属性建立绑定/赋值 | 验证数据类型、深复制 |
| 展开 | 在调用、数组或对象字面量中展开值 | 深复制、保留全部对象元数据 |
| 模板字符串 | 插值、多行文本、标签处理 | HTML/SQL/URL 安全转义 |
| 可选链 `?.` | 左侧为空值时停止连续属性/调用链 | 吞掉其他异常、验证方法可调用 |
| 空值合并 `??` | 只为 `null`/`undefined` 提供回退 | 为 `0`、`false`、空字符串回退 |

```mermaid
flowchart LR
  A[外部原始数据] --> B[先验证结构]
  B --> C[解构读取所需字段]
  C --> D[可选链读取允许缺失的路径]
  D --> E[空值合并应用业务默认]
  E --> F[展开生成下一份浅层结构]
  F --> G[模板字符串生成文本]
  G --> H[按输出上下文编码或渲染]
```

语法简洁不等于输入可信。来自 API、存储或 URL 的值仍需在使用前验证。

## 2. 解构：用模式读取结构

解构模式出现在绑定位置或赋值左侧。数组解构通过迭代协议按顺序取值；对象解构按属性键读取。

### 2.1 数组解构

```js
const coordinates = [121.47, 31.23];
const [longitude, latitude] = coordinates;

console.log(longitude); // 121.47
console.log(latitude);  // 31.23
```

模式比来源长时，多出的绑定得到 `undefined`；逗号可以跳过位置；剩余元素收集为新数组并且必须位于末尾。

```js
const values = ['first', 'second', 'third', 'fourth'];
const [first, , third, ...rest] = values;

console.log(first); // first
console.log(third); // third
console.log(rest);  // ['fourth']
```

数组解构不要求右侧一定是 Array，只要求可迭代。普通对象默认不可迭代，会抛出 `TypeError`。

```js
const [firstCodePoint] = '😀A';
console.log(firstCodePoint); // 😀，字符串迭代按码点前进

// const [value] = { 0: 'x', length: 1 }; // TypeError：不可迭代
```

变量交换可以用一次解构赋值完成；右侧会先求值，再写入左侧目标。

```js
let left = 'L';
let right = 'R';
[left, right] = [right, left];

console.log(left, right); // R L
```

### 2.2 对象解构与重命名

对象解构中的左侧属性名是要读取的键，冒号右侧才是新绑定名。

```js
const lesson = {
  id: 'js-04',
  title: '现代表达式语法',
  status: 'learning',
};

const { id, title: lessonTitle } = lesson;

console.log(id);          // js-04
console.log(lessonTitle); // 现代表达式语法
// console.log(title);    // ReferenceError：没有创建 title 绑定
```

动态属性键写在方括号中。

```js
const fieldName = 'status';
const { [fieldName]: currentStatus } = lesson;
console.log(currentStatus); // learning
```

### 2.3 默认值只处理 `undefined`

解构默认值在属性缺失或值严格为 `undefined` 时求值；`null`、`false`、`0` 和空字符串都会被保留。

```js
const settings = {
  pageSize: 0,
  theme: null,
  compact: false,
};

const {
  pageSize = 20,
  theme = 'system',
  compact = true,
  locale = 'zh-CN',
} = settings;

console.log(pageSize); // 0
console.log(theme);    // null
console.log(compact);  // false
console.log(locale);   // zh-CN
```

默认表达式是惰性求值，只有确实需要回退时才执行。

```js
function createDefaultId() {
  console.log('创建默认 id');
  return 'generated';
}

const { id: existingId = createDefaultId() } = { id: 'saved' };
// 不输出日志
```

### 2.4 嵌套解构与缺失中间层

嵌套模式可直接读取深层字段，但中间值若是 `undefined` 或 `null`，继续解构会抛出 `TypeError`。需要给可能缺失的中间对象提供默认值，或先验证后分步读取。

```js
const payload = { user: { profile: { name: 'Lili' } } };

const {
  user: {
    profile: { name = '匿名' } = {},
  } = {},
} = payload;

console.log(name); // Lili
```

这里的默认对象仍只在值为 `undefined` 时生效。若 `user: null`，模式依然失败。外部数据存在 `null` 的可能时，先用 schema/条件验证，或用可选链读取更清楚。

过深解构还会隐藏字段来源，错误信息也难对应业务层。超过两三层时，分步绑定通常更可维护。

### 2.5 对象剩余属性

对象 rest 把尚未取出的自有可枚举属性浅复制到一个新对象；它必须是模式最后一项。

```js
const account = {
  id: 7,
  passwordHash: 'secret',
  displayName: 'Lili',
};

const { passwordHash, ...publicAccount } = account;
console.log(publicAccount); // { id: 7, displayName: 'Lili' }
```

这可用于明确排除字段，但不能作为安全边界：未来新增的敏感字段会自动进入 `publicAccount`。对外响应更适合显式白名单选择允许字段。

### 2.6 绑定模式与赋值模式

声明会创建新绑定：

```js
const { status } = lesson;
```

赋值模式写入已经存在的目标。对象解构赋值作为独立语句时需要括号，否则开头花括号会被解析为块。

```js
let status;
let title;

({ status, title } = lesson);
console.log(status, title);
```

无分号代码风格中，要防止前一行与括号表达式意外连接。格式化和 lint 规则应统一处理。

### 2.7 函数参数中的解构

参数解构适合少量稳定字段，并能配合整个参数的默认值。

```js
function formatLesson({ title, status = 'draft' } = {}) {
  if (typeof title !== 'string' || title.trim() === '') {
    throw new TypeError('title 必须是非空字符串');
  }
  return `${title}：${status}`;
}

console.log(formatLesson({ title: 'CSS' })); // CSS：draft
```

`= {}` 只处理调用时缺少参数或传 `undefined`；`formatLesson(null)` 仍会抛出 `TypeError`。参数字段多、验证复杂时先接收 `options`，在函数体中验证和解构更容易调试。

## 3. 展开语法：三个上下文、三套要求

`...value` 的行为取决于语法位置。

| 位置 | 来源要求 | 结果 |
| --- | --- | --- |
| 函数调用 `fn(...value)` | 可迭代 | 每个迭代值成为实参 |
| 数组字面量 `[...value]` | 可迭代 | 每个迭代值成为数组元素 |
| 对象字面量 `{...value}` | 可转对象 | 复制自有可枚举属性 |

### 3.1 调用展开

```js
const dimensions = [10, 20, 30];
console.log(Math.max(...dimensions)); // 30
```

这等价于把数组元素分别传入参数位置。引擎对一次调用的实参数量存在上限，大数组不能安全地全部展开调用；使用循环或 `reduce()`。

```js
const maximum = dimensions.reduce(
  (current, value) => Math.max(current, value),
  -Infinity,
);
```

### 3.2 数组展开

数组展开可以连接或浅复制可迭代值。

```js
const basics = ['HTML', 'CSS'];
const roadmap = [...basics, 'JavaScript'];
const copy = [...roadmap];

console.log(roadmap); // ['HTML', 'CSS', 'JavaScript']
console.log(copy === roadmap); // false
```

普通对象不可迭代，不能直接写 `[...object]`。Map 展开产生键值二元数组，Set 展开产生成员。

```js
const selected = new Set(['html', 'css']);
console.log([...selected]); // ['html', 'css']
```

### 3.3 对象展开与覆盖顺序

对象展开按出现顺序复制自有可枚举属性，后面的同名属性覆盖前面的值，但属性在枚举顺序中的原位置不会因为覆盖而简单移动。

```js
const defaults = { theme: 'system', pageSize: 20 };
const userInput = { theme: 'dark' };
const options = { ...defaults, ...userInput, pageSize: 50 };

console.log(options); // { theme: 'dark', pageSize: 50 }
```

展开 `null` 或 `undefined` 到对象字面量不会添加属性；其他原始值会按对象属性规则处理。字符串的索引字符是自有可枚举属性。

```js
console.log({ ...'ok' }); // { 0: 'o', 1: 'k' }
console.log({ ...null }); // {}
```

不要依赖这类原始值展开构造业务对象，它通常说明输入类型未验证。

### 3.4 条件属性

可以通过条件表达式展开对象或空对象。

```js
const includeDebug = false;
const request = {
  endpoint: '/lessons',
  ...(includeDebug ? { debug: true } : {}),
};
```

表达式 `...(condition && { debug: true })` 也常见，因为假值没有可枚举属性，但三元形式更明确地表达展开来源始终是对象。

### 3.5 展开是浅复制

展开只创建外层数组或对象，嵌套引用仍共享。

```js
const state = {
  profile: { name: 'Lili' },
  tags: ['js'],
};

const shallow = { ...state };
shallow.profile.name = 'Changed';

console.log(state.profile.name); // Changed
```

不可变更新要复制从根到目标字段的每一层。

```js
const nextState = {
  ...state,
  profile: {
    ...state.profile,
    name: 'Next',
  },
  tags: [...state.tags, 'syntax'],
};
```

对象展开不会保留原型、不可枚举属性和属性描述符；读取 getter 时会执行 getter，并把结果写成普通数据属性。类实例不能靠 `{...instance}` 正确克隆。

## 4. 模板字符串

模板字符串使用反引号，可包含换行，并在 `${expression}` 中求值后转为字符串。

```js
const title = 'JavaScript';
const completed = 4;
const total = 18;

const summary = `${title}：${completed}/${total}`;
console.log(summary); // JavaScript：4/18
```

插值中可以写任意表达式，但复杂计算应先绑定具名结果，避免文本模板同时承担业务逻辑。

```js
const percent = Math.round((completed / total) * 100);
const readable = `完成度：${percent}%`;
```

### 4.1 多行、转义与原始文本

```js
const message = `第一行
第二行`;

const escaped = `反引号：\`，插值标记：\${notEvaluated}`;
```

缩进会进入字符串内容。生成协议文本或快照时，应明确处理前导空格和换行，而不是假设编辑器缩进会被去掉。

`String.raw` 是内置标签函数，可让反斜杠序列按原始形式出现在结果中，适合需要保留转义符的文本；它不是通用安全编码器。

```js
const pathPattern = String.raw`C:\notes\javascript`;
console.log(pathPattern);
```

### 4.2 标签模板

标签模板把模板拆成静态字符串片段和插值值传给函数。标签的返回值可以是任意类型，不要求是字符串。

```js
function inspectTemplate(strings, ...values) {
  return { strings: [...strings], values };
}

const result = inspectTemplate`topic=${'CSS'} count=${3}`;
console.log(result);
// strings 含三个静态片段，values 为 ['CSS', 3]
```

模板对象还提供 `strings.raw` 查看未处理转义的片段。同一个标签模板位置在多次求值时会获得具有稳定身份的模板字符串数组，标签实现可以据此缓存解析结果。

### 4.3 模板字符串不提供安全转义

```js
const untrusted = '<img src=x onerror=alert(1)>';
const html = `<p>${untrusted}</p>`;
```

`html` 仍包含可被解释为标记的文本。若交给 `innerHTML`，可能形成注入风险。纯文本写入 DOM 使用 `textContent`；必须生成 HTML 时使用经过审计的模板系统、上下文相关编码或可信类型策略。SQL、URL、Shell 也各有不同的安全边界，不能靠模板字符串本身处理。

## 5. 可选链 `?.`

当链左侧是 `null` 或 `undefined` 时，可选链返回 `undefined` 并停止该条连续链的后续求值；其他假值不会触发短路。

### 5.1 属性、动态属性与可选调用

```js
const response = {
  user: {
    profile: { name: 'Lili' },
  },
};

console.log(response.user?.profile?.name); // Lili
console.log(response.account?.profile?.name); // undefined

const field = 'name';
console.log(response.user?.profile?.[field]); // Lili

const onComplete = undefined;
onComplete?.({ ok: true }); // 不调用，也不抛错
```

`object.method?.()` 只在 `method` 为空值时跳过；若它存在但不是函数，仍会抛出 `TypeError`。可选链不会验证接口类型。

```js
const plugin = { run: 'not a function' };
// plugin.run?.(); // TypeError
```

### 5.2 短路范围与副作用

只有连续的可选链受保护。索引或参数表达式在提前短路时不会求值。

```js
let index = 0;
const items = null;
const value = items?.[index++];

console.log(value); // undefined
console.log(index); // 0
```

括号会结束连续链，后续普通属性访问可能再次抛错。

```js
const data = null;
console.log(data?.user?.name); // undefined
// console.log((data?.user).name); // TypeError
```

每个允许缺失的位置都应显式写出 `?.`。`root?.child.name` 在 root 非空但 child 为 `undefined` 时仍会访问 `.name` 并失败；写成 `root?.child?.name` 才同时允许两层缺失。

### 5.3 不能使用的位置

可选链的结果不能作为赋值目标，也不能用于构造调用的构造器位置或标签模板标签。

```js
// response.user?.name = 'new'; // SyntaxError
// new Service?.();             // SyntaxError
// String.raw?.`text`;          // SyntaxError
```

删除可选属性是允许的；对象为空值时结果为 `true`，否则执行普通删除语义。

```js
const draft = null;
console.log(delete draft?.temporary); // true
```

### 5.4 不要掩盖必需数据

如果业务契约要求 `response.user.profile.name` 一定存在，全部写成可选链会把协议错误变成 `undefined` 并延后暴露。可选链应只放在允许缺失的边界；必需字段应验证并明确失败。

## 6. 空值合并 `??`

`left ?? right` 在 left 为 `null` 或 `undefined` 时返回 right，否则返回 left。右侧只在需要时求值。

```js
console.log(0 ?? 20);       // 0
console.log(false ?? true); // false
console.log('' ?? '匿名');  // ''
console.log(null ?? '匿名');// 匿名
```

`||` 根据真值回退，会覆盖合法的零、false 和空字符串。

```js
const configuredPageSize = 0;
console.log(configuredPageSize || 20); // 20
console.log(configuredPageSize ?? 20); // 0
```

选择规则取决于业务：若空字符串本来就应该视为缺失，先执行显式规范化，再使用 `??`，不要依靠 `||` 同时承担多种数据清洗语义。

```js
function emptyToUndefined(value) {
  if (typeof value === 'string' && value.trim() === '') return undefined;
  return value;
}

const label = emptyToUndefined('  ') ?? '未命名';
```

语法禁止在无括号时直接把 `??` 与 `&&` 或 `||` 混合，以避免优先级含糊。

```js
const result = (cached || local) ?? fallback;
```

### 6.1 空值合并赋值

`target ??= value` 只在 target 为空值时求值右侧并赋值。属性访问只求值一次，适合有 getter/setter 或复杂索引的目标。

```js
const preferences = { pageSize: 0 };
preferences.pageSize ??= 20;
preferences.locale ??= 'zh-CN';

console.log(preferences); // { pageSize: 0, locale: 'zh-CN' }
```

相应的 `||=` 和 `&&=` 根据真值决定是否赋值，语义不同。选择时写清哪些值属于“缺失”。

## 7. 组合读取、默认和更新

下面的组合是常见数据流：

```js
const locale = response.user?.preferences?.locale ?? 'zh-CN';

const nextUser = {
  ...response.user,
  preferences: {
    ...response.user?.preferences,
    locale,
  },
};
```

执行顺序是：

1. 可选链读取允许缺失的 preferences 和 locale。
2. `??` 只在 locale 为 `null`/`undefined` 时采用默认值。
3. 外层展开创建新 user 对象。
4. 内层展开创建新 preferences 对象；展开 `undefined` 不添加属性。
5. 最后写入经过解析的 locale。

这仍是浅层不可变更新。`response.user` 若本身是必需字段，应先验证，而不是让展开 `undefined` 静默创建一个只有 preferences 的新对象。

## 8. 完整案例：解析并更新学习面板配置

输入允许部分字段缺失，但字段一旦提供就必须满足类型和范围。合法的 `0`、`false` 和空标题要按业务分别处理。

```js
const defaults = {
  locale: 'zh-CN',
  pageSize: 20,
  compact: false,
  title: '学习面板',
  filters: {
    topic: undefined,
    completed: undefined,
  },
};

const rawConfig = {
  pageSize: 50,
  compact: false,
  title: '',
  filters: { topic: 'JavaScript' },
  callbacks: {
    onReady(config) {
      console.log(`ready:${config.pageSize}`);
    },
  },
};
```

### 8.1 验证与规范化

```js
const PAGE_SIZES = new Set([10, 20, 50]);

function resolveConfig(raw = {}) {
  if (raw === null || typeof raw !== 'object' || Array.isArray(raw)) {
    throw new TypeError('配置必须是对象');
  }

  const {
    locale = defaults.locale,
    pageSize = defaults.pageSize,
    compact = defaults.compact,
    title = defaults.title,
    filters = {},
    callbacks = {},
  } = raw;

  if (typeof locale !== 'string' || locale.trim() === '') {
    throw new TypeError('locale 必须是非空字符串');
  }
  if (!PAGE_SIZES.has(pageSize)) {
    throw new RangeError('pageSize 只能是 10、20 或 50');
  }
  if (typeof compact !== 'boolean') {
    throw new TypeError('compact 必须是 Boolean');
  }
  if (typeof title !== 'string') {
    throw new TypeError('title 必须是 String');
  }
  if (filters === null || typeof filters !== 'object' || Array.isArray(filters)) {
    throw new TypeError('filters 必须是对象');
  }
  if (callbacks === null || typeof callbacks !== 'object') {
    throw new TypeError('callbacks 必须是对象');
  }

  const topic = filters.topic ?? defaults.filters.topic;
  const completed = filters.completed ?? defaults.filters.completed;

  if (topic !== undefined && typeof topic !== 'string') {
    throw new TypeError('filters.topic 必须是 String');
  }
  if (completed !== undefined && typeof completed !== 'boolean') {
    throw new TypeError('filters.completed 必须是 Boolean');
  }
  if (callbacks.onReady !== undefined && typeof callbacks.onReady !== 'function') {
    throw new TypeError('callbacks.onReady 必须是函数');
  }

  return {
    config: {
      ...defaults,
      locale: locale.trim(),
      pageSize,
      compact,
      title,
      filters: { topic, completed },
    },
    callbacks,
  };
}
```

`title: ''` 被保留，因为需求允许用户主动清空标题；`compact: false` 不会被默认值覆盖；filters 只重建允许字段，避免未知字段直接扩散。

### 8.2 输出与可选调用

```js
function initializePanel(raw) {
  const { config, callbacks } = resolveConfig(raw);
  callbacks.onReady?.(config);

  const topicLabel = config.filters.topic ?? '全部主题';
  const titleLabel = config.title === '' ? '（无标题）' : config.title;

  return {
    config,
    summary: `${titleLabel}｜${topicLabel}｜每页 ${config.pageSize} 项`,
  };
}

console.log(initializePanel(rawConfig));
```

可观察结果包括控制台先输出 `ready:50`，返回配置保留 `compact: false` 和空标题，filters 的 completed 为 `undefined`，summary 使用显式业务规则把空标题显示为“（无标题）”。

### 8.3 更新嵌套状态

```js
function updateTopic(config, topic) {
  if (topic !== undefined && typeof topic !== 'string') {
    throw new TypeError('topic 必须是 String 或 undefined');
  }

  return {
    ...config,
    filters: {
      ...config.filters,
      topic,
    },
  };
}
```

调用后原 config 和原 filters 都保持不变，next config 的未修改字段保留。这里只复制两层，因为目标字段位于第二层。

### 8.4 失败注入

```js
const failures = [
  null,
  { pageSize: 0 },
  { compact: 'false' },
  { filters: null },
  { filters: { completed: 0 } },
  { callbacks: { onReady: true } },
];

for (const input of failures) {
  try {
    initializePanel(input);
  } catch (error) {
    console.log(error.name, error.message);
  }
}
```

失败路径验证整个参数默认值不会接管 `null`、合法集合检查不等于真值检查、字符串 `'false'` 不是布尔值、嵌套对象必须验证，以及可选调用不会把存在但不可调用的值当成缺失。

## 9. 调试与审查清单

1. 解构冒号右侧是绑定名，不是类型标注；确认最终创建了哪些名称。
2. 默认值只处理 `undefined`，检查 `null` 是否应先规范化。
3. 嵌套解构的每个中间层是否可能为 `null`；过深时改成分步读取。
4. rest 和 spread 都是浅复制，检查嵌套对象是否仍共享引用。
5. 调用或数组展开要求可迭代；对象展开遵循自有可枚举属性规则。
6. 展开类实例时检查原型、方法、getter 和不可枚举属性是否丢失。
7. 模板字符串进入 HTML、URL、SQL 或 Shell 前使用对应上下文的安全接口。
8. 可选链是否只用于允许缺失的路径；必需字段不能被静默转成 undefined。
9. 括号是否结束了连续可选链；方法存在但非函数仍会抛错。
10. `??` 与 `||` 的业务缺失定义是否一致；零、false 和空字符串是否有效。

## 10. 练习与完成标准

实现一个搜索请求构造器：

- 输入允许 `query`、`page`、`pageSize`、`filters.status` 和可选 `onBuilt`。
- `page: 0` 必须保留，不能被默认值覆盖。
- 空 query 按业务规范化为 `undefined`，不是依靠 `||` 隐式处理。
- 使用解构读取，但所有字段都要验证。
- 使用两层展开创建新请求，不能修改输入。
- 使用可选调用执行回调；回调存在但不是函数时明确报错。
- 使用模板字符串生成仅供日志的摘要，不把未转义文本写入 HTML。
- 测试 `undefined`、`null`、`0`、`false`、空字符串、错误嵌套结构和未知状态。

完成标准是：能逐步解释解构默认、空值合并和可选链的求值顺序；证明输入对象未被修改；所有共享嵌套引用都符合预期；异常路径具有稳定的类型和信息。

## 来源

- [MDN：Destructuring](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Destructuring)（访问日期：2026-07-17）
- [MDN：Spread syntax](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Spread_syntax)（访问日期：2026-07-17）
- [MDN：Template literals](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Template_literals)（访问日期：2026-07-17）
- [MDN：Optional chaining](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Optional_chaining)（访问日期：2026-07-17）
- [MDN：Nullish coalescing operator](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Nullish_coalescing)（访问日期：2026-07-17）
