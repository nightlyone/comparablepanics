---
name: add-analyzer-regression-test
description: 'Add or update a comparablepanics regression test. Use for analyzer bugs, false positives, false negatives, fixture-based reproductions, or changes that need analysistest coverage.'
argument-hint: 'Describe the analyzer bug or behavior change to cover'
---

# Add Analyzer Regression Test

Use this skill when you need to reproduce and fix analyzer behavior with a focused `analysistest` case.

## Procedure

1. Identify the smallest failing scenario and decide whether it fits an existing fixture package in `testdata/src/` or needs a new one.
2. Add or update a fixture `.go` file under `testdata/src/` with the minimal code needed to exercise the analyzer.
3. Add exact `// want "..."` expectations for diagnostics that should be reported.
4. If you created a new top-level fixture package, add its package name to `analysistest.Run(...)` in [../../../analyzer_test.go](../../../analyzer_test.go).
5. Make the smallest analyzer change needed in [../../../analyzer.go](../../../analyzer.go) or, if necessary, the thin CLI wiring in [../../../cmd/comparablepanics/main.go](../../../cmd/comparablepanics/main.go).
6. Run `go test ./...`.
7. If the change touches build or entrypoint behavior, also run `go build ./cmd/comparablepanics`.

## Repo-Specific Checks

- Prefer extending existing packages `a`, `b`, or `c` unless a separate package boundary is needed.
- Keep production changes narrower than the fixture that motivated them.
- Never modify `testdata/pkg/`; treat it as external fixture data.
- If diagnostic wording changes intentionally, update the matching `// want` strings in fixtures.