# JavaScript 值、变量、类型、表达式与类型转换

JavaScript 程序通过表达式产生值，通过绑定保存值，通过运算符组合值。理解类型转换的目的不是记忆几组特殊结果，而是能够沿着规范规则判断：输入是什么类型、运算符需要什么类型、转换在哪一步发生、最终得到什么值或异常。

## 1. 从源码到结果

下面是阅读一个表达式时最基本的检查顺序。

```mermaid
flowchart LR
  A[读取标识符或字面量] --> B[得到操作数的值]
  B --> C[按运算符规则求值]
  C --> D{是否需要转换}
  D -- 是 --> E[执行 ToPrimitive / ToNumeric / ToString 等抽象操作]
  D -- 否 --> F[直接计算]
  E --> F
  F --> G[产生值或抛出异常]
```

规范中的 `ToNumber`、`ToBoolean` 等是抽象操作，不是必须由开发者直接调用的函数。`Number(value)`、`Boolean(value)` 是可以调用的显式转换接口；运算符内部也可能执行相应抽象操作。

```js
const unitPrice = 25;        // Number 字面量产生值，const 创建绑定
const rawCount = '4';        // String 值
const count = Number(rawCount); // 显式转换为 Number
const total = unitPrice * count; // * 对两个数值求值

console.log(total); // 100
```

## 2. 值、绑定与变量

值是程序计算的数据。绑定是名称与值之间的关联。日常所说的“变量”通常指可通过名称访问的绑定，但绑定本身不是装值的对象，也不代表对象内容不可变化。

### 2.1 `const`、`let` 与 `var`

| 声明 | 绑定能否重新赋值 | 作用域 | 声明前访问 | 同一作用域重复声明 |
| --- | --- | --- | --- | --- |
| `const` | 否 | 块级 | `ReferenceError` | 不允许 |
| `let` | 是 | 块级 | `ReferenceError` | 不允许 |
| `var` | 是 | 函数级或脚本全局 | 得到 `undefined` | 允许 `var` 重复声明 |

默认使用 `const` 表达“该名称不会重新指向另一个值”；确实需要改变绑定时使用 `let`。新代码通常不需要 `var`，因为函数作用域和声明提升更容易让状态跨出预期边界。

```js
const settings = { theme: 'light' };
settings.theme = 'dark'; // 合法：改变对象属性

// settings = {};        // TypeError：试图重新赋值 const 绑定

let retries = 0;
retries += 1;            // 合法：let 绑定可以重新赋值
```

`const` 只限制绑定重新赋值，不会递归冻结对象。需要限制对象表层属性时可用 `Object.freeze()`，但它默认也不是深冻结。

### 2.2 声明、初始化与赋值

三个动作应区分：

- 声明：在作用域中创建名称，例如 `let score`。
- 初始化：第一次给绑定设置值；未带初始化器的 `let score` 初始化为 `undefined`。
- 赋值：之后用 `score = 10` 替换绑定当前的值。

`let` 和 `const` 绑定从作用域开始处已存在，但直到声明执行前都不可访问，这一段称为暂时性死区。`const` 必须在声明时提供初始化器。

```js
{
  // console.log(status); // ReferenceError，不是 undefined
  const status = 'ready';
  console.log(status); // ready
}
```

作用域、提升与闭包的完整机制将在后续专篇展开；此处需要记住：不能用“代码写在后面但会提升”作为提前读取 `let` 或 `const` 的理由。

## 3. ECMAScript 的八种语言类型

ECMAScript 语言值分为七种原始类型和 Object 类型。

| 类型 | 示例 | 关键语义 |
| --- | --- | --- |
| Undefined | `undefined` | 只有一个值，常表示尚未提供或尚未产生结果 |
| Null | `null` | 只有一个值，常用于显式表达“没有对象/值” |
| Boolean | `true`、`false` | 逻辑真值 |
| String | `'Lili'` | UTF-16 码元序列，不可变 |
| Symbol | `Symbol('id')` | 每次创建通常得到唯一值，可作属性键 |
| Number | `42`、`NaN`、`Infinity` | IEEE 754 双精度浮点数体系 |
| BigInt | `42n` | 任意精度整数，不能直接与 Number 混合算术 |
| Object | `{}`、`[]`、函数、日期等 | 可包含属性，按身份比较 |

原始值本身不可变。字符串方法返回新字符串，不会改变原字符串。

```js
const label = 'roadmap';
const upper = label.toUpperCase();

console.log(label); // roadmap
console.log(upper); // ROADMAP
```

对象是属性集合并具有身份。两个内容相同但分别创建的对象不是同一个值。

```js
const first = { id: 1 };
const second = { id: 1 };
const alias = first;

console.log(first === second); // false
console.log(first === alias);  // true
```

### 3.1 `undefined` 与 `null`

两者都可表达缺失，但来源和意图不同。一个可维护的约定是：`undefined` 表示未提供、未初始化或属性不存在；`null` 表示业务或 API 主动确认该位置当前没有值。具体项目应统一，而不是在同一数据模型中随机混用。

```js
const profile = {
  nickname: undefined, // 也可能直接省略该属性
  avatarUrl: null,     // 已确认没有头像
};

console.log(profile.missing); // undefined
console.log('nickname' in profile); // true
console.log('missing' in profile);  // false
```

`JSON.stringify()` 对对象中的 `undefined` 属性会省略，对数组中的 `undefined` 会输出 `null`；JSON 专篇会进一步说明序列化边界。

### 3.2 Number、`NaN` 与安全整数

Number 包含有限数、正负无穷和 `NaN`。浮点表示会产生十进制小数不能精确表示的问题。

```js
console.log(0.1 + 0.2); // 0.30000000000000004
console.log(Number.isSafeInteger(9_007_199_254_740_991)); // true
console.log(Number.isSafeInteger(9_007_199_254_740_992)); // false
```

`NaN` 表示某次数值运算没有得到有效数值。它的类型仍是 `number`，且不等于自身。检查转换结果应使用 `Number.isNaN(value)`，不要用全局 `isNaN()` 代替，因为全局版本会先转换参数。

```js
const invalid = Number('12px');

console.log(typeof invalid);          // number
console.log(invalid === NaN);         // false
console.log(Number.isNaN(invalid));   // true
console.log(Number.isNaN('invalid')); // false，不先转换
```

货币计算可在业务边界明确时转换成整数最小单位，例如分；金融、计费或跨币种领域仍应使用经过验证的十进制定点方案，不能仅靠 `toFixed()` 掩盖累计误差。

### 3.3 BigInt

BigInt 用于超出安全整数范围且只需要整数运算的场景。字面量后缀是 `n`。

```js
const nextId = 9_007_199_254_740_993n;
console.log(nextId + 1n); // 9007199254740994n

// nextId + 1; // TypeError：BigInt 与 Number 不能直接混合算术
console.log(Number(10n)); // 10；转换大值时可能丢失精度
```

BigInt 除法丢弃小数部分：`5n / 2n` 得到 `2n`。`JSON.stringify()` 默认不能序列化 BigInt，必须先定义明确的字符串或其他协议表示。

### 3.4 String 与 Unicode 边界

规范把 String 定义为 UTF-16 码元序列，因此 `length` 统计码元，不保证等于用户看到的字符数量。一些字符由两个码元组成。

```js
const icon = '😀';
console.log(icon.length);       // 2
console.log([...icon].length);  // 1，按 Unicode 码点迭代
```

码点也不总等于用户感知字符；组合附加符和家庭 emoji 可能包含多个码点。需要按可见字符截断时应使用适合本地化的分段能力并做真实语言测试。

### 3.5 Symbol

Symbol 常用于创建不会与普通字符串键冲突的属性键，或参与语言协议，例如迭代协议中的 `Symbol.iterator`。

```js
const internalId = Symbol('internalId');
const item = { title: 'CSS' };
item[internalId] = 42;

console.log(item[internalId]); // 42
console.log(Object.keys(item)); // ['title']，不包含 Symbol 键
```

`Symbol('id') !== Symbol('id')`，描述文本不决定身份。`Symbol.for('id')` 使用全局 Symbol 注册表，语义不同，使用前要确认是否确实需要跨模块共享键。

## 4. `typeof` 能回答什么

`typeof` 返回一个字符串，适合做有限的运行时分类，但它不是完整类型系统。

| 表达式 | 结果 |
| --- | --- |
| `typeof undefined` | `'undefined'` |
| `typeof null` | `'object'`，历史兼容结果 |
| `typeof true` | `'boolean'` |
| `typeof 'x'` | `'string'` |
| `typeof 1` / `typeof NaN` | `'number'` |
| `typeof 1n` | `'bigint'` |
| `typeof Symbol()` | `'symbol'` |
| `typeof {}` / `typeof []` | `'object'` |
| `typeof function () {}` | `'function'` |

区分数组用 `Array.isArray()`，区分 `null` 用严格相等，识别具体对象通常要依据数据结构或协议，不能只依赖 `typeof`。

```js
function classify(value) {
  if (value === null) return 'null';
  if (Array.isArray(value)) return 'array';
  return typeof value;
}

console.log(classify(null)); // null
console.log(classify([]));   // array
```

`typeof neverDeclared` 对从未声明的标识符返回 `'undefined'`，但对暂时性死区中的 `let`/`const` 绑定仍会抛出 `ReferenceError`。

## 5. 字面量、表达式与语句

字面量直接写出值，例如 `10`、`'ok'`、`true`、`[1, 2]`。表达式会求值得到一个值，例如 `price * count`、`user.name`、`createUser()`。语句控制程序执行，例如声明、`if`、`for`、`return`。

```js
const subtotal = 20 * 3; // 20、3 是字面量；20 * 3 是表达式；整行是声明语句
```

运算符优先级决定分组，结合性决定同级运算符如何归组。维护代码时，括号往往比依赖读者记忆完整优先级表更清楚。

```js
const total = price * count + shipping;
const explicitTotal = (price * count) + shipping;

const ratio = 100 / 10 / 2; // (100 / 10) / 2，结果 5
const power = 2 ** 3 ** 2;   // 2 ** (3 ** 2)，结果 512
```

不要把有副作用的复杂表达式压在一行。先求值、再验证、再更新状态，可以让错误定位和测试更直接。

## 6. 显式类型转换

外部数据经常以字符串进入程序：表单字段、URL 查询参数、存储值和部分响应字段都需要在边界转换。转换必须和验证相邻进行。

### 6.1 转成 Boolean

`Boolean(value)` 使用真值规则。以下值为假值：

- `false`
- `undefined`
- `null`
- `+0`、`-0`、`0n`
- `NaN`
- 空字符串 `''`

其余值为真值，包括字符串 `'false'`、字符串 `'0'`、空数组 `[]` 和空对象 `{}`。

```js
console.log(Boolean('false')); // true
console.log(Boolean('0'));     // true
console.log(Boolean([]));      // true
console.log(Boolean(0));       // false
```

因此不能用 `Boolean(raw)` 解析文本形式的布尔字段。应明确接受值集合。

```js
function parseBoolean(raw) {
  if (raw === 'true') return true;
  if (raw === 'false') return false;
  throw new TypeError('布尔值必须是 true 或 false');
}
```

### 6.2 转成 Number

`Number(value)` 要求整个字符串符合数值语法，忽略首尾空白；空字符串转换为 `0`。`parseInt()` 和 `parseFloat()` 从字符串开头解析，遇到不能继续的字符就停止，因此用途不同。

```js
console.log(Number(' 42 '));       // 42
console.log(Number(''));           // 0
console.log(Number('12px'));       // NaN
console.log(Number.parseInt('12px', 10)); // 12
console.log(Number.parseFloat('1.5rem')); // 1.5
```

如果字段不允许空字符串，必须先检查空值，不能让 `Number('') === 0` 将“未填写”误当成合法零。解析整数时明确传入十进制基数，并验证业务范围和整数性。

```js
function parsePositiveInteger(raw, fieldName) {
  const normalized = raw.trim();
  if (normalized === '') {
    throw new TypeError(`${fieldName} 不能为空`);
  }

  const value = Number(normalized);
  if (!Number.isInteger(value) || value <= 0) {
    throw new RangeError(`${fieldName} 必须是正整数`);
  }
  return value;
}
```

### 6.3 转成 String

`String(value)` 是通用显式转换。模板字符串插值和字符串连接也可能触发转换。

```js
console.log(String(null));      // 'null'
console.log(String(undefined)); // 'undefined'
console.log(String(42n));       // '42'

const count = 3;
const message = `共 ${count} 项`; // 共 3 项
```

直接执行 `String(Symbol('x'))` 可得到描述形式，但某些隐式字符串化路径对 Symbol 会抛错。日志和序列化协议不应假设所有值都能无损转成字符串再恢复。

### 6.4 转成 BigInt

`BigInt()` 可从符合整数语法的字符串或整数 Number 创建 BigInt。小数 Number、无效字符串会抛出异常，而不是返回 `NaN`。

```js
console.log(BigInt('9007199254740993')); // 9007199254740993n
console.log(BigInt(42));                 // 42n

// BigInt(1.5); // RangeError
// BigInt('1.5'); // SyntaxError
```

## 7. 隐式转换与常见运算符

隐式转换本身不是错误；问题在于输入类型不明确时，代码行为难以从接口契约判断。

### 7.1 `+` 的双重职责

二元 `+` 在操作数转成原始值后，只要一方是字符串，就执行字符串连接；否则执行数值加法。其他算术运算符通常按数值路径处理。

```js
console.log('5' + 2); // '52'
console.log('5' - 2); // 3
console.log(1 + 2 + '3'); // '33'
console.log('1' + 2 + 3); // '123'
```

边界输入先转换可消除歧义：`Number(rawA) + Number(rawB)` 明确要求数值加法。

### 7.2 `&&`、`||` 与 `??` 返回操作数

逻辑运算符不保证返回 Boolean：

- `a && b`：`a` 为假值时返回 `a`，否则返回 `b`。
- `a || b`：`a` 为真值时返回 `a`，否则返回 `b`。
- `a ?? b`：`a` 是 `null` 或 `undefined` 时返回 `b`，否则返回 `a`。

它们都会短路：确定结果后不再求值右侧表达式。

```js
const configuredRetries = 0;

console.log(configuredRetries || 3); // 3，错误地覆盖合法的 0
console.log(configuredRetries ?? 3); // 0，只对 null/undefined 回退

false && console.log('不会执行');
```

不能在不加括号时直接混用 `??` 与 `&&`/`||`，语法会拒绝这种歧义。即使加了括号，也应让回退策略一眼可见。

### 7.3 对象到原始值

运算需要原始值时，规范会执行 `ToPrimitive`，可能使用 `Symbol.toPrimitive`、`valueOf()` 或 `toString()`。这解释了部分内置对象参与运算的结果，但业务代码不应依赖难读的对象隐式转换。

```js
const amount = {
  cents: 500,
  [Symbol.toPrimitive](hint) {
    if (hint === 'number') return this.cents;
    return `${this.cents} cents`;
  },
};

console.log(Number(amount)); // 500
console.log(String(amount)); // '500 cents'
```

自定义转换协议会影响比较、模板和算术，必须记录契约并测试所有使用路径。多数数据对象直接提供 `toCents()`、`format()` 等具名方法更清晰。

## 8. 四种相等语义

JavaScript 不只有一套“相等”。

| 机制 | 类型转换 | `NaN` 与自身 | `+0` 与 `-0` | 常见位置 |
| --- | --- | --- | --- | --- |
| `==` | 会按宽松相等算法转换 | 不等 | 相等 | 遗留代码、明确的空值判断 |
| `===` | 不转换不同类型 | 不等 | 相等 | 默认业务比较 |
| `Object.is()` | 不转换 | 相等 | 不等 | 需要 SameValue 语义 |
| SameValueZero | 不转换 | 相等 | 相等 | `Set`、`Map` 键、`includes()` 等内置行为 |

```js
console.log(0 == false);          // true
console.log(0 === false);         // false
console.log(NaN === NaN);         // false
console.log(Object.is(NaN, NaN)); // true
console.log(Object.is(+0, -0));   // false
console.log([NaN].includes(NaN)); // true，SameValueZero
```

默认使用 `===` 和 `!==`，让类型契约显式。`value == null` 是一个有意同时匹配 `null` 和 `undefined` 的特例；如果团队接受这种写法，应注明意图，否则写成 `value === null || value === undefined`。

无论哪套相等语义，对象都按身份比较，不会递归比较属性内容。

## 9. 完整案例：解析订单行输入

输入来自表单或 URL，类型全部是字符串：

```js
const rawOrderLine = {
  sku: '  css-book ',
  unitPriceCents: '2590',
  quantity: '2',
  expedited: 'false',
};
```

目标输出必须满足：SKU 非空；单价是非负安全整数；数量是 1 到 99 的整数；加急字段只接受 `'true'` 或 `'false'`；总价继续用整数分表示。

### 9.1 分步处理

```js
function parseInteger(raw, { name, min, max }) {
  const text = raw.trim();
  if (text === '') {
    throw new TypeError(`${name} 不能为空`);
  }

  const value = Number(text);
  if (!Number.isSafeInteger(value)) {
    throw new TypeError(`${name} 必须是安全整数`);
  }
  if (value < min || value > max) {
    throw new RangeError(`${name} 必须在 ${min} 到 ${max} 之间`);
  }
  return value;
}

function parseOrderLine(raw) {
  const sku = raw.sku.trim();
  if (sku === '') {
    throw new TypeError('SKU 不能为空');
  }

  const unitPriceCents = parseInteger(raw.unitPriceCents, {
    name: '单价',
    min: 0,
    max: Number.MAX_SAFE_INTEGER,
  });
  const quantity = parseInteger(raw.quantity, {
    name: '数量',
    min: 1,
    max: 99,
  });
  const expedited = parseBoolean(raw.expedited);
  const totalCents = unitPriceCents * quantity;

  if (!Number.isSafeInteger(totalCents)) {
    throw new RangeError('总价超出安全整数范围');
  }

  return { sku, unitPriceCents, quantity, expedited, totalCents };
}

console.log(parseOrderLine(rawOrderLine));
```

### 9.2 可观察输出

```js
const expectedOrderLine = {
  sku: 'css-book',
  unitPriceCents: 2590,
  quantity: 2,
  expedited: false,
  totalCents: 5180,
};

console.log(expectedOrderLine);
```

处理顺序是：规范化文本、拒绝空值、显式转换、检查语言层数值边界、检查业务范围、计算派生值、再次检查计算结果。返回对象中的每个字段已经具有后续代码可依赖的类型。

### 9.3 失败分支

```js
const failures = [
  { ...rawOrderLine, quantity: '' },
  { ...rawOrderLine, quantity: '2.5' },
  { ...rawOrderLine, quantity: '100' },
  { ...rawOrderLine, unitPriceCents: '25元' },
  { ...rawOrderLine, expedited: 'yes' },
];

for (const input of failures) {
  try {
    parseOrderLine(input);
  } catch (error) {
    console.log(error.name, error.message);
  }
}
```

这些失败输入分别覆盖空字符串被 `Number()` 转为零、非整数、业务越界、部分数值字符串和模糊布尔文本。异常名称与消息是可观察证据，可在单元测试中断言错误类型。

## 10. 调试与测试清单

遇到类型转换问题时，按以下顺序检查：

1. 在输入边界记录原始值和 `typeof`，敏感信息只记录经过脱敏的结构。
2. 将“规范化、转换、验证、计算”拆开，定位最先偏离预期的步骤。
3. 对 Number 同时检查 `Number.isNaN()`、`Number.isFinite()`、整数性和业务范围。
4. 对对象确认是在比较身份，还是业务上确实需要比较字段。
5. 对默认值确认应使用假值回退 `||`，还是只对空值回退 `??`。
6. 对大整数确认序列化协议、数据库字段和其他语言端是否能无损承载。
7. 对字符串确认限制的是 UTF-16 长度、码点数、字节数，还是用户感知字符数。

最低测试集合应包括正常值、空字符串、首尾空白、零、负数、边界值、超大数、`NaN` 产生路径、`null`、`undefined` 和错误类型。来自 URL、表单或 JSON 的数据不能因为“看起来像数字”就跳过运行时验证。

## 11. 练习与完成标准

实现一个查询参数解析器，接受 `page`、`pageSize`、`archived`：

- `page` 是从 1 开始的整数。
- `pageSize` 是 10、20、50 之一。
- `archived` 只接受 `'true'`、`'false'`，缺失时为 `false`。
- 返回的数据中不保留未验证的原始字符串。
- 为 `''`、`'0'`、`'1.5'`、`'20px'`、`'TRUE'`、缺失值分别写测试。

完成标准是：能解释每个输入在哪一步转换；不会把空字符串误判为合法零；不会用字符串真值解析布尔值；错误输入具有稳定、可断言的失败结果。

## 来源

- [ECMAScript® Language Specification：ECMAScript Data Types and Values](https://tc39.es/ecma262/multipage/ecmascript-data-types-and-values.html)（访问日期：2026-07-17）
- [ECMAScript® Language Specification：Type Conversion 与 Testing and Comparison Operations](https://tc39.es/ecma262/multipage/abstract-operations.html)（访问日期：2026-07-17）
- [MDN：Grammar and types](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Grammar_and_types)（访问日期：2026-07-17）
- [MDN：Equality comparisons and sameness](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Equality_comparisons_and_sameness)（访问日期：2026-07-17）
