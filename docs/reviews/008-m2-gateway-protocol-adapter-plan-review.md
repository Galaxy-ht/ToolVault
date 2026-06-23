# M2 Gateway + Protocol Adapter Proposal Review Report

## Verdict

**PASS**

更新版已满足进入 M2 实现阶段的 proposal 条件。原先的关键 blocker 已闭环：M1 Registry deviations 已作为 M2 internal input contract 接受，M2 scope 明确限制为 read-only discovery + deterministic projection，并明确排除 Runtime、Policy、Credential、execution、SDK integration、full MCP server、streaming、admin/debug discovery 和 production ingress。

## Score

**96 / 100**

- 范围控制：20 / 20
- 架构一致性：19 / 20
- 任务可执行性：19 / 20
- 接口设计清晰度：14 / 15
- 测试与验证：15 / 15
- 风险与依赖管理：9 / 10

## Blockers

无 proposal 内容 blocker。

Process note: [docs/proposals/002-m2-gateway-protocol-adapter-plan.md](/Users/egon/project/ToolVault/docs/proposals/002-m2-gateway-protocol-adapter-plan.md) 当前仍是 untracked 文件。进入实现前应确保它进入正式变更集，否则 Builder/Reviewer 不能把它当作稳定批准基线。

## Required Changes

无内容层面的 required change。

进入实现前需要完成一项流程动作：将 proposal 纳入 git 变更集或通过项目认可的批准记录固化。这个不是 proposal 设计缺陷，但会影响 M2 task contract 的可追溯性。

## Optional Improvements

- 可以在 Gateway error model 中后续补充 `non_discoverable`，区分 “不存在” 与 “存在但 M2 不暴露”。当前用 `404` 统一处理是可接受的。
- `ToolDescription.Payload any` 后续实现时应测试 JSON marshal 稳定性，避免投影输出不可比较。
- M2.1 guardrail 实现时应用具体模式检查 execution endpoint drift，例如 `invoke`、`execute`、`run`、`call` 等命名。

## Boundary Check

- 是否只聚焦 Gateway / Protocol Projection？**是。** 第 10-13 行和第 26-50 行明确为 Registry-backed read-only discovery + deterministic projection。
- 是否越界到 Runtime？**否。** Non-goals 明确排除 runtime behavior 和 runtime selection。
- 是否越界到 Policy？**否。** 明确不做 auth、authorization、policy decision。
- 是否越界到 Credential？**否。** 明确不做 secret retrieval 或 injection。
- 是否越界到 Sandbox / Dashboard / Streaming / Composition？**否。** 均列入 non-goals，并在 guardrail 中拒绝。
- 是否引入不必要依赖？**否。** 第 826 行后明确禁止第三方依赖、OpenAI SDK、MCP SDK、router、DB driver、observability libraries、Web UI toolchains。

## Implementation Readiness

- Builder 是否知道要做什么？**是。** M2.1-M2.8 拆分清楚。
- Builder 是否知道不能做什么？**是。** 第 32-34 行和 non-goals/guardrail 明确拒绝越界实现。
- Builder 是否知道能改哪些目录？**是。** Allowed / forbidden directories 明确，`cmd/toolvault/` 默认排除。
- Builder 是否知道如何验证完成？**是。** 每个任务都有 verification commands，并要求新增 `make m2-gateway-check`。
- Reviewer 是否知道如何审查产出？**是。** Scope、目录、依赖、投影边界、HTTP surface、guardrail 检查点都足够具体。

## Final Recommendation

**批准进入 M2 实现阶段。**

建议下一步从 **M2.1 M2 Guardrail** 开始，而不是直接写 Gateway 业务代码。这样可以先把 forbidden modules、execution endpoints、SDK imports、third-party dependencies、debug/admin discovery 和 UI/persistence drift 卡住，再允许 `internal/gateway/` 实现进入仓库。
