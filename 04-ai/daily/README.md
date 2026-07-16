# 每日 AI 实验

AI 侧每天做一个“小型能力验证”，而不是只追新闻或收藏工具。固定任务和评测样例，每次只改变一个变量，才可能看见模型、Prompt、Context、RAG、Tool 或 Workflow 的真实差异。

## 实验题目

- Prompt：指令结构、示例、约束、输出格式和失败行为。
- Model/API：结构化输出、Streaming、多模态、Token 与上下文限制。
- Context：选择、排序、压缩、记忆、冲突和权限。
- RAG：解析、Chunk、查询、召回、重排、引用和无答案。
- Tool/MCP：Schema、工具选择、参数、确认、幂等和错误恢复。
- Workflow/Agent：路由、并行、计划、停止、恢复和任务完成率。
- Evaluation/Safety：Grader、回归样例、注入、越权、隐私和成本攻击。

## 建议节奏

- 周一：模型 API 或 Prompt。
- 周二：Context 或 RAG。
- 周三：Tool、MCP 或 Workflow。
- 周四：Evaluation、Observability 或 Safety。
- 周五：复测一个新模型/版本或拆解一个真实 AI 产品能力。

节奏只是选题辅助，不要求为了打卡每天切换主题。

## 记录方式

1. 复制 [_template.md](_template.md)。
2. 按 `YYYY/MM/YYYY-MM-DD-hypothesis.md` 保存。
3. 评测样例和脚本放在记录的同名目录；密钥、敏感输入和完整用户数据禁止提交。
4. 模型或框架升级时复跑旧实验，把结果链接到原记录，形成纵向比较。

## 完成标准

- [ ] 有明确假设、固定样例、变量和基线。
- [ ] 记录完整模型标识、版本、日期和关键参数。
- [ ] 同时观察质量、延迟、Token、成本和失败类型。
- [ ] 结论可复现、可证伪，并说明适用边界。

