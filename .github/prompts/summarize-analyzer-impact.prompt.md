---
description: "Summarize how a code change affects comparablepanics diagnostics, analysistest coverage, and analyzer pass wiring. Use for diffs, PRs, or local edits touching analyzer behavior."
name: "Summarize Analyzer Impact"
argument-hint: "Diff, PR summary, changed files, or describe the edit to inspect"
agent: "agent"
---

Review the provided change and summarize its impact on this repository's Go analyzer.

Use [AGENTS.md](../../AGENTS.md), [analyzer.go](../../analyzer.go), and [analyzer_test.go](../../analyzer_test.go) as the primary repo context.

In the response:

1. State whether the change likely affects reported diagnostics, generic-type handling, or analyzer traversal.
2. State whether existing `analysistest` fixtures under [testdata/src/](../../testdata/src/) are still sufficient.
3. Call out whether `analysistest.Run(...)` package coverage in [analyzer_test.go](../../analyzer_test.go) needs to change.
4. Mention whether analyzer dependency wiring such as `inspect.Analyzer`, `usesgenerics.Analyzer`, or CLI integration in [cmd/comparablepanics/main.go](../../cmd/comparablepanics/main.go) is affected.
5. Recommend the smallest validation command to run next, usually `go test ./...`.

Keep the summary concise and specific to the changed behavior.