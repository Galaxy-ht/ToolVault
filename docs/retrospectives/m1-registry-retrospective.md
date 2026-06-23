# M1 Registry Retrospective

Status: completed

Date: 2026-06-23

## 1. M1 Original Goal

M1 的原始目标是交付一个小而可测试的 Tool Registry 内部模块，作为后续 Gateway、Runtime、Policy、Credential 和 Observability 等模块的依赖基础。

Proposal 001 对 M1 的边界定义很明确：

- Registry 只拥有 tool metadata、lifecycle state、validation、lookup、listing、update 和 deletion semantics。
- 实现必须留在 `internal/registry/`。
- 实现必须使用 Go 标准库。
- 存储必须是 in-memory。
- 不引入数据库、网络协议、Web UI、第三方依赖、`pkg/` 公共包或其他 core module 逻辑。
- 先建立 M1 Registry verification gate，再允许 Registry `.go` 文件进入仓库。
- 每个实现阶段都必须有测试和 review gate。

原计划的最小 `ToolSpec` 只包含 `ID`、`Name`、`DisplayName`、`Description`、`Version`、`Status` 和 `Tags`。生命周期规则要求支持 `draft`、`active`、`deprecated`、`disabled`，并保持 hard delete 语义。

## 2. Actual Completed Work

M1 实际完成了一个可运行、可测试的内部 Registry baseline：

- `internal/registry/spec` 定义了 `ToolSpec`、`ToolID`、`ToolStatus`、`ToolAction` 和 validation。
- `internal/registry` 定义了 Registry interface、request/option 类型、lifecycle helper 和 error contract。
- `internal/registry/memory` 实现了标准库 in-memory Registry。
- 支持 `Register`、`Get`、`List`、`Update`、`Delete` 和 `SetStatus`。
- 支持 `draft`、`active`、`deprecated`、`disabled` 生命周期状态。
- 支持 hard delete、稳定 list 排序、status/tag filtering、duplicate ID/name 检查和 defensive copy。
- 支持 opaque version validation，并在 `Update` 中通过 `ExpectedVersion` 做 optimistic update check。
- 补齐了 M1 guardrail：`make m1-registry-check`。
- 补齐了验收文档、示例和 release notes。
- M1 相关阶段使用了 review 文档记录边界、测试、残余风险和合并建议。

M1 也产生了已记录的 proposal deviations：

- `ToolSpec` 增加了 `Actions` 和 `Metadata`。
- `Update` 增加了 `ExpectedVersion` 和 `VersionConflict`。
- 状态变更通过 `SetStatus` 完成，而不是通过 `Update`。
- in-memory Registry 除了拒绝重复 `ID`，也拒绝重复 `Name`。

这些偏差已在 `docs/acceptance/m1-registry.md` 和 M1 review 中记录，但 M2 使用 Registry 前仍需要人工确认接受或收窄。

## 3. What Worked

Proposal-first 流程有效。`docs/proposals/001-m1-registry-implementation-plan.md` 在实现前先固定了目标、非目标、目录边界、验证命令、任务拆分和 rollback plan，使 Builder 有明确输入，Reviewer 有明确审查标准。

M1 guardrail 有效。`make m1-registry-check` 把 M1 从 M0 的“禁止内部 Go 实现”推进到“只允许 Registry Go 实现”，同时继续拒绝 forbidden core modules、数据库、Web UI、dependency artifacts 和第三方依赖。

小步提交和 tag 有效。M1 使用了 `m1-registry-plan-approved`、`m1.1-toolspec-pass`、`m1.2-interface-pass`、`m1.3-memory-store-pass`、`m1.4-lifecycle-pass`、`m1.6-acceptance-pass`、`m1-registry-complete` 等节点，便于回看阶段边界和定位变更来源。

Reviewer gate 有效。每个阶段的 review 都检查了 scope drift、boundary violations、dependency creep、missing tests 和 Git diff 完整性。尤其是多次发现 untracked 文件没有进入 diff，避免了断链文档和不完整合并。

Acceptance checklist 有效。最终验收没有只看测试通过，而是把 proposal alignment、out-of-scope 内容和 residual risks 一起记录下来，为 M2 提供了明确的前置决策清单。

## 4. Problems Exposed

最大问题是 contract drift。`Actions`、`Metadata`、`ExpectedVersion`、`VersionConflict`、`SetStatus` 和 duplicate `Name` 都是有价值的实现选择，但它们超过了 Proposal 001 的最小模型。Review 发现了这些偏差，但部分偏差是在后续 acceptance 中被记录，而不是在引入前通过 proposal 或 task contract 明确批准。

第二个问题是 Git diff hygiene 不稳定。M1.3 和 M1.6 review 都发现关键文件仍处于 untracked 状态。测试能通过不代表变更集完整；如果按当时 diff 合并，会出现实现遗漏或 README 引用断链。

第三个问题是任务编号和阶段命名存在轻微漂移。早期 `docs/06-m1-registry-tasks.md` 和 Proposal 001 的子任务编号并不完全一致，例如 CLI probe 在一个文档中更靠前，而 Proposal 中是 optional M1.5。这不会影响最终实现，但会增加 Builder 和 Reviewer 对齐成本。

第四个问题是 verification chain 曾需要修复。`be99891 Fix M1 verify chain` 表明从 M0 bootstrap 到 M1 check 的验证组合需要在阶段切换时同步维护。后续 milestone 不应假设旧验证命令自然适配新范围。

第五个问题是条件通过容易积累残余风险。M1.1、M1.2、M1.3、M1.4、M1.6 都有 conditional pass。条件通过可以保持推进速度，但如果没有一个最终显式决策点，残余风险会被后续 milestone 误认为已经完全批准。

## 5. Git Management Lessons

每个阶段合并前必须同时看 `git status --short` 和 `git diff --stat`。只看 diff 会漏掉 untracked 文件；只看测试会漏掉文档引用和提交完整性问题。

阶段提交应保持 narrow diff。M1 中最容易 review 的提交是只触碰当前任务允许目录的提交，例如 guardrail、ToolSpec、interface、memory implementation 和 lifecycle hardening。混入 docs、examples 或 unrelated cleanup 会降低 reviewer 判断质量。

Tag 对 milestone 很有帮助。M1 的阶段 tag 让回溯变得简单，M2 应继续在 plan approved、interface pass、implementation pass、acceptance pass 等节点打 tag。

Merge commit 应保留阶段语义。`merge: complete M1 Registry` 和 `docs: add v0.1.0 registry MVP release notes` 清楚区分了实现闭环和发布说明。M2 也应避免把 implementation、acceptance 和 release notes 混在一个不可拆的提交里。

Git 管理不能替代治理记录。即使提交历史能显示偏差何时引入，M2 仍需要在 proposal、task contract 或 acceptance 中明确写出“接受 / 收窄 / 延后”的决策。

## 6. Codex Builder / Reviewer Collaboration Lessons

Builder 在明确 allowed directories、forbidden directories、acceptance criteria 和 verification commands 时表现最好。M1 的成功依赖于任务合同足够具体，Builder 不需要自行推断是否可以 touching Gateway、Runtime、Policy、Credential 或 `pkg/`。

Reviewer 的价值不只是找代码 bug。M1 review 最重要的贡献是发现范围漂移、proposal deviation、untracked 文件、dependency risk 和 verification gap。这些问题单靠 `go test` 不会暴露。

Reviewer 应在发现 contract drift 时要求明确决策，而不是只记录为 optional note。`Actions` / `Metadata` 从 M1.1 开始被反复指出，说明 Reviewer 的信号是正确的；M2 需要更早把这种信号升级为 human decision gate。

Builder 和 Reviewer 都需要把“测试通过”和“任务完成”分开。M1 多次出现测试通过但仍有治理条件的情况。后续阶段应明确：verification pass 是必要条件，不是 acceptance decision。

## 7. Rules M2 Should Reuse

M2 应继续沿用以下规则：

- Proposal-first：任何 Gateway 或 Registry consumption 工作先有 proposal，明确依赖方向和模块边界。
- Task contract-first：每个 Builder 任务必须包含 goal、allowed directories、forbidden directories、acceptance criteria、verification commands 和 risks。
- Guardrail-first：在实现 Gateway 前先建立 M2 verification gate，继续拒绝未批准的 Runtime、Policy、Credential、Protocol、Observability、Web UI、database 和 third-party dependency drift。
- Interface-first：Gateway 只能通过批准的 Registry contract 依赖 Registry，不应读取 Registry internals。
- Standard-library-first：除非 proposal 明确获批，否则 M2 继续避免新增 heavyweight dependency。
- Review-before-merge：每个阶段都保留 reviewer 文档，检查 scope、boundary、dependencies、tests、Git diff completeness 和 residual risks。
- Acceptance-before-next-milestone：进入下一个模块前必须有 acceptance checklist，列出已完成内容、未完成内容和必须人工确认的偏差。
- Tag milestone gates：继续为 proposal approval、implementation pass、acceptance pass 和 release point 打 tag。

## 8. Problems M2 Should Avoid

M2 不应在未确认 M1 deviations 前直接把 Registry 当成稳定依赖。Gateway planning 必须先决定是否接受 `Actions` / `Metadata`、`ExpectedVersion` / `VersionConflict`、`SetStatus` 和 duplicate `Name` 规则。

M2 不应把 Gateway 变成 Runtime、Policy、Credential 或 Protocol Adapter 的提前实现。Gateway 只能处理已批准的入口、路由和 Registry consumption 边界；执行、鉴权、密钥注入、观测和协议适配仍应等待对应 milestone 或 proposal。

M2 不应依赖 untracked 文件或隐含本地状态。合并前必须确认所有被 README、docs、tests 或 examples 引用的文件都在 Git 变更集中。

M2 不应让 conditional pass 累积到阶段末尾。每个 conditional pass 都应该产生一个明确后续动作：立即修正、创建 proposal、记录 human approval，或明确从当前 milestone 移除。

M2 不应扩大 public surface。除非有边界变更 proposal，Registry 和 Gateway 的 M2 work 都应保持 internal module，不新增 `pkg/` 公共包或外部服务承诺。

M2 不应把验证命令当成固定不变。每个 milestone 开始时都要重新检查 `make verify` 是否覆盖当前阶段真实风险，并在 proposal 中写明 narrow verification command。
