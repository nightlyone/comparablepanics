---
description: "Use when editing Go analyzer tests, analysistest fixtures, or files under testdata/src. Covers fixture layout, // want diagnostics, and keeping analyzer_test.go in sync."
name: "Analysistest Fixtures"
applyTo: "analyzer_test.go,testdata/src/**/*.go"
---

# Analysistest Fixtures

- Follow the existing `golang.org/x/tools/go/analysis/analysistest` pattern in [analyzer_test.go](../../analyzer_test.go).
- Keep reproductions in `testdata/src/...` and prefer a focused fixture case before broadening analyzer logic.
- Express expected diagnostics with exact `// want "..."` comments in fixture files.
- If you add a new top-level fixture package under `testdata/src/`, also add that package name to `analysistest.Run(...)` in [analyzer_test.go](../../analyzer_test.go).
- Do not edit files under [testdata/pkg/](../../testdata/pkg/); they are supporting module fixtures, not the primary test cases.
- After changing analyzer behavior or fixtures, validate with `go test ./...`.