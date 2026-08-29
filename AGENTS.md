# comparablepanics Agent Notes

This repository is a small Go static analyzer built on `golang.org/x/tools/go/analysis`.

## Start Here

- Read [README.md](README.md) for project intent, background, and current TODOs.
- The core analyzer implementation is in [analyzer.go](analyzer.go).
- The standalone CLI entrypoint is in [cmd/comparablepanics/main.go](cmd/comparablepanics/main.go).
- Tests are driven by [analyzer_test.go](analyzer_test.go) and fixtures under [testdata/](testdata/).

## Build And Test

- Use `go test ./...` for normal validation.
- Use `go build ./cmd/comparablepanics` to verify the analyzer binary still builds.
- The module targets Go `1.24.0` and pins toolchain `go1.25.6` in [go.mod](go.mod). Avoid introducing older-version compatibility assumptions.

## Code Shape

- Keep analyzer logic in the root package `comparablepanics`; `cmd/comparablepanics` should stay as a thin `singlechecker` wrapper.
- `Analyzer` in [analyzer.go](analyzer.go) is the public integration surface. Changes to required analyzers, reported diagnostics, or traversal flow should remain intentional and minimal.
- The analysis relies on `types.Info`, AST nodes, and `inspect`/`usesgenerics` passes from `golang.org/x/tools` rather than ad hoc parsing.

## Test Conventions

- Follow `analysistest` conventions when adding coverage.
- Put analyzer fixtures under `testdata/src/...`.
- Express expected diagnostics with `// want "..."` comments in fixture files.
- Prefer adding focused fixture cases over expanding production code without a reproducer.

## Editing Guidance

- Preserve the existing small-scope style: targeted functions, early returns, and minimal abstraction.
- Do not vendor or modify contents under [testdata/pkg/](testdata/pkg/); those are module fixtures used by tests.
- If behavior changes, update or add the corresponding `analysistest` case in `testdata/src/`.