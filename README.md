# Software Product Engineer Roadmap

从高级前端工程师到软件产品工程师的长期成长路线。

五个方向也分别提供从零基础开始的独立路径，任何人都可以按自己的目标选择一个入口，在理论、日常练习、项目和笔记中持续成长。

## 为什么叫「狸力」

「狸力」（lí lì）出自《山海经·南山经》。书中记载它“其状如豚，有距，其音如狗吠”，出现时当地会兴起许多土木工程，即“见则其县多土功”。

这里借用的是“土功”背后的建造者意象：向下深挖问题，打牢技术基础，再把前端、产品、交互、AI、后端与数据组合成真正可用、能够持续演进的软件产品。

狸力代表的不是一条速成路线，而是一种长期建设的工程精神——既能看见需求，也愿意亲手把系统一层层建起来。

## 最终目标

成为能够发现问题、定义产品、设计体验、构建系统、应用 AI，并持续推动产品演进的软件产品工程师。

## 五个核心方向

| 方向 | 目标 | 路线图 | 每日记录 | 模块笔记 | 实践项目 |
| --- | --- | --- | --- | --- | --- |
| 01 前端深化 | 从 Web 入门到架构、性能与平台能力 | [路线图](01-frontend/roadmap.md) | [每日前端项目](01-frontend/daily/README.md) | [笔记](01-frontend/notes/README.md) | [项目](01-frontend/projects.md) |
| 02 产品能力 | 从产品观察到定义、数据与商业能力 | [路线图](02-product/roadmap.md) | [每日产品/功能拆解](02-product/daily/README.md) | [笔记](02-product/notes/README.md) | [项目](02-product/projects.md) |
| 03 交互设计 | 从交互原则到复杂流程与 AI UX | [路线图](03-interaction-design/roadmap.md) | [每日交互拆解](03-interaction-design/daily/README.md) | [笔记](03-interaction-design/notes/README.md) | [项目](03-interaction-design/projects.md) |
| 04 AI 工程 | 从模型 API 到 RAG、Agent 与评估体系 | [路线图](04-ai/roadmap.md) | [每日 AI 实验](04-ai/daily/README.md) | [笔记](04-ai/notes/README.md) | [项目](04-ai/projects.md) |
| 05 后端与数据 | 从编程入门到分布式与云原生能力 | [路线图](05-backend-data/roadmap.md) | [每日算法题](05-backend-data/daily/README.md) | [笔记](05-backend-data/notes/README.md) | [项目](05-backend-data/projects.md) |

五个方向当前已经形成 445 篇入门、初级与中级笔记，其中 229 篇构成入门与初级基础。统一覆盖范围、学习方式和维护约定见 [学习笔记知识库](learning-notes.md)。

## 一份持续演进的路线图

这里没有固定的三年期限，也不以追赶进度为目标。技术、产品和我们的理解都会持续变化，因此路线图本身也是产品：学习、实践、记录、验证，再根据新证据迭代。

- 五个方向都是可以从零开始、独立完成的成长路线。
- 理论建立心智模型，日常练习形成手感，项目验证综合能力，笔记沉淀自己的知识体系。
- 具体框架和工具会过时，稳定原理、问题解决过程和可复现证据应长期保留。
- 新内容优先依据官方文档、标准、论文、源码和可复现实验。
- 路线图的更新记录在 [CHANGELOG](CHANGELOG.md)，综合项目维护在 [项目索引](projects/README.md)，通用模板位于 [templates](templates/README.md)。

## 学习闭环

```text
理解概念 → 完成最小实验 → 加入真实项目 → 主动制造边界与故障
→ 使用工具排障 → 总结方案取舍 → 沉淀文档、模式或工具
```

“学会”意味着能解释、能实现、能测试、能排障、能做取舍，并能沉淀复用，而不只是看完课程或跑通 Demo。

## 如何维护

1. 用 Obsidian 将仓库根目录作为 Vault 打开，或直接用 VS Code 编辑 Markdown。
2. 任选一个方向，从该方向的 `README.md` 和 `roadmap.md` 开始，不要求五条路线同步推进。
3. 每完成一个模块，把笔记放进该方向的 `notes/`；每天的最小练习放进 `daily/`。
4. 项目结束后使用 [项目复盘模板](templates/project-retrospective.md)，把代码、数据、文档或体验证据链接回来。
5. 通过 Git 提交和 GitHub 同步；路线图更新、学习笔记和日常记录尽量分别提交。

## 目录结构

```text
.
├── README.md
├── learning-notes.md
├── 01-frontend/
├── 02-product/
├── 03-interaction-design/
├── 04-ai/
├── 05-backend-data/
├── templates/
├── projects/
└── CHANGELOG.md
```
