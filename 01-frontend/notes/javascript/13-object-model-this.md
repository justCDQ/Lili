# Prototype、原型链、Class、this 与对象模型

## 是什么

对象通过内部原型委托属性查找；class 是基于原型的声明语法；this 由普通函数调用方式决定，箭头函数捕获外层 this。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

实例方法位于 prototype；new 创建对象、连接原型并绑定 this；Object.create 显式设原型；组合常比深继承稳定。

## 实际使用

```js
class User{ constructor(name){this.name=name;} greet(){return `Hi ${this.name}`;} }
const u=new User('Li');
const bound=u.greet.bind(u);
```

## 常见错误与边界

解构方法后 this 可能丢失；箭头函数不能作构造器；修改内建原型会污染全局。

## 相关补充知识

属性查找沿原型链进行，class 是建立原型方法的语法机制。普通函数的 `this` 由调用方式决定，箭头函数捕获外层 `this`；提取方法时需显式绑定或改造 API。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Inheritance_and_the_prototype_chain)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Classes)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/this)

访问日期：2026-07-16。
