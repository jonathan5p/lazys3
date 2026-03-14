# Agent Development Guide — lazys3 (learning branch)

Authoritative reference for all AI agents working in this repository.

This is the **learning branch**. The project is being rebuilt from scratch, one milestone at a time.
The reference implementation lives on `main`. Do not copy from it — guide the student to write every line.

---

## Project Goal

Rebuild `lazys3` — a terminal UI for browsing AWS S3 — from scratch as a Go learning project.
See `CURRICULUM.md` for the full milestone plan.

---

## Build & Run Commands

```sh
make build      # compile binary → ./lazys3
make run        # go run ./cmd/lazys3 (no binary)
make clean      # remove ./lazys3
make tidy       # go mod tidy
```

---

## Test Commands

```sh
make test                              # run all tests
make test-verbose                      # go test -v ./...
make test-race                         # with race detector
make cover                             # tests + per-function coverage summary
make cover-html                        # generates coverage.html

# Run a single test by name
go test ./... -run TestMyFunction -v

# Run all tests in one package
go test ./internal/s3/ -v
```

---

## Lint & Format Commands

```sh
make fmt        # gofmt -w . (in-place)
make vet        # go vet ./...
```

No `golangci-lint` or `staticcheck` installed. Use `go vet` as the linter.

---

## Pre-Commit Checklist (MANDATORY)

Run before **every** commit. Zero tolerance — fix all errors first.

```sh
make pre-commit     # gofmt check + go vet + go test ./...
```

Order: formatting → vet → tests. If any step fails, fix and re-run from the top.

---

## Code Style

### Formatting

All code must pass `gofmt`. Tabs for indentation — never spaces.

### Imports

Three groups separated by blank lines:

```go
import (
    "context"
    "fmt"

    "github.com/some/thirdparty"

    "github.com/jonathan5p/lazys3/internal/something"
)
```

1. Standard library
2. Third-party packages
3. Internal packages

### Naming

- Unexported identifiers: `camelCase`. Exported: `PascalCase`.
- Acronyms stay uppercase: `S3`, `TUI`, `JSON`, `URL`.
- Constructors: `New<Type>` (e.g. `NewApp`, `NewClient`).
- Receiver names: short, consistent per type (e.g. `a` for `AppModel`, `c` for `Client`).
- Boolean fields: positive framing (`loading`, `active`), not negations (`notLoaded`, `inactive`).

### Types & Receivers

- Use pointer receivers when a method mutates state.
- Prefer concrete structs over interfaces; introduce an interface only when needed for testing or multiple implementations.
- No `interface{}` / `any` unless unavoidable.

### Functions & Methods

- Keep every function under **20 lines**. Flag and refactor anything longer.
- One responsibility per function.

### Error Handling

- Always check errors immediately — never discard with `_`.
- Propagate with `return err` at the lowest level.
- Add context when wrapping: `fmt.Errorf("load config: %w", err)`.
- `log.Fatalf` only in `main()`. All other packages return errors to the caller.

### Comments

- No comments that restate what the code does.
- Doc comments (starting with the exported symbol name) only on exported identifiers.
- Unexported helpers must be self-explanatory through naming alone.

---

## TDD Rules

- Write the failing test first; confirm it is RED before implementing.
- Use only the standard `testing` package — no third-party assertion libraries.
- Prefer table-driven tests (`[]struct{ ... }`) for multiple input cases.
- Test case names must be descriptive: `"empty bucket list"`, not `"test1"`.
- After every code change, run `go test ./...` automatically without asking.

---

## Language

All code, comments, commit messages, test names, and documentation must be in **English**.
