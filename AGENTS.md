# Agent Development Guide — lazys3

This is the authoritative reference for all AI agents working in this repository.
It is Go-project-specific and supersedes any generic global agent instructions.

---

## Project Overview

`lazys3` is a terminal UI (TUI) for browsing and downloading AWS S3 objects.
Stack: Go 1.24, `jesseduffield/gocui` (TUI), `aws-sdk-go-v2` (S3).

**Package layout:**

```
main.go                            # root stub — not the real entry point
cmd/lazys3-tui/
  main.go                          # real entry point: wires logger, S3 client, App
  logger/logger.go                 # file-backed structured logger (Info/Debug/Error)
  s3/client.go                     # AWS S3 wrapper (ListBuckets, ListObjects, DownloadObject)
  ui/app.go                        # gocui TUI: layout, keybindings, render loop
  IMPROVEMENTS.md                  # backlog of planned TUI enhancements (see below)
```

**Improvement backlog:**

`cmd/lazys3-tui/IMPROVEMENTS.md` is the authoritative list of planned enhancements for the TUI.
Before starting any new feature work, consult this file to understand the intended direction of the
project and to avoid implementing something that conflicts with a planned improvement. When a backlog
item is fully implemented, remove or mark it as done in that file.

---

## Build & Run

There is no Makefile. Use the `go` toolchain directly.

```sh
# Build
go build -o lazys3-tui ./cmd/lazys3-tui

# Run without building
go run ./cmd/lazys3-tui
```

---

## Test Commands

```sh
# Run all tests
go test ./...

# Run all tests with verbose output
go test -v ./...

# Run tests for a single package
go test ./cmd/lazys3-tui/s3/

# Run a single named test
go test ./cmd/lazys3-tui/s3/ -run TestListBuckets

# Run with race detector (recommended for concurrent/TUI code)
go test -race ./...
```

---

## Lint & Format Commands

```sh
# Format all code in place
gofmt -w .

# Check formatting without rewriting (must produce no output before a commit)
gofmt -l .

# Static analysis (built-in, no external linter installed)
go vet ./...

# Tidy dependencies
go mod tidy
```

> No `golangci-lint` or `staticcheck` is installed. Use `go vet` as the linter.
> If a Makefile is added later, switch to `make <target>` instead of direct commands.

---

## Pre-Commit Checklist (MANDATORY)

Run all three in order and fix every error before committing:

```sh
gofmt -l .    # must produce no output
go vet ./...  # must produce no output
go test ./... # all tests must pass
```

Zero tolerance. Never commit with outstanding errors.

---

## Code Style

### Formatting

- All code must pass `gofmt`. Tabs for indentation — never spaces.

### Imports

Three groups, separated by blank lines:

1. Standard library
2. Third-party packages
3. Internal packages (`github.com/jonathan5p/lazys3/...`)

```go
import (
    "context"
    "io"

    "github.com/aws/aws-sdk-go-v2/aws"

    "github.com/jonathan5p/lazys3/cmd/lazys3-tui/logger"
)
```

### Naming

- Unexported: `camelCase`. Exported: `PascalCase`.
- Acronyms stay uppercase: `S3`, `TUI`, `HTTP`, `URL`.
- Receiver names: short, consistent per type (`a` for `*App`, `c` for `*Client`, `l` for `*Logger`).
- Constructors follow `New<Type>` (`NewApp`, `NewClient`, `New`).
- Boolean fields use positive framing (`loaded`, `IsDir`), not negations.

### Types & Receivers

- Use pointer receivers on all methods of a struct if any method mutates state.
- Prefer concrete structs over interfaces; introduce an interface only when needed for testing or multiple implementations.
- No `interface{}` / `any` unless unavoidable.

### Error Handling

- Always check errors immediately — never discard with `_`.
- Propagate with `return err` at the lowest level.
- Add context when wrapping: `fmt.Errorf("loadObjects: %w", err)`.
- `log.Fatalf` only in `main()`. All other packages return errors to the caller.
- The gocui view-init idiom is correct and must be followed exactly:

  ```go
  if err != nil && err != gocui.ErrUnknownView { return err }
  ```

### Functions & Methods

- Keep every function under 20 lines. Flag and refactor anything longer.
- Keybinding handlers must match the gocui signature:

  ```go
  func(g *gocui.Gui, v *gocui.View) error
  ```

- Every log call must include function context and key variable values:

  ```go
  a.log.Info("loadObjects: bucket=%q prefix=%q", bucket, prefix)
  ```

### Comments

- No comments that restate what the code does.
- Doc comments (starting with the symbol name) only on exported identifiers.
- Unexported helpers must be self-explanatory through naming alone.

---

## TUI-Specific Rules (gocui)

- Always check `gocui.ErrUnknownView` on every `g.SetView(...)` call.
- Never call `g.SetCurrentView` outside of layout or a keybinding handler.
- The `layout` function is called on every frame — keep it idempotent. Use the `!a.loaded` guard for one-time initialization.
- Colors are set via `g.FgColor`, `g.BgColor`, and per-view `v.FgColor` / `v.BgColor`. Use `gocui.ColorDefault` to inherit from the user's terminal theme.

---

## TDD Rules

- Follow TDD in spirit: tests are part of the development cycle, not an afterthought. Every non-trivial piece of logic must have a test before the feature is considered done.
- Red-green cycles are not required to be minimal. It is acceptable to write a fuller implementation before returning to green, as long as tests are written as part of the same cycle — not after the fact.
- Use only the standard `testing` package (no third-party assertion libraries).
- Prefer table-driven tests (`[]struct{ ... }`) for multiple input cases.
- Test case names must be descriptive — they appear in failure output (e.g. `"directory with nested prefix"`, not `"test"`).
- Mock AWS calls by extracting an interface over `*s3.Client` and providing a fake in tests.
- Do not test `main()` directly; test the packages it wires together.

---

## Language

- All code, comments, commit messages, test names, and documentation must be in English.
- Conversation with the developer may be in Spanish or English.
