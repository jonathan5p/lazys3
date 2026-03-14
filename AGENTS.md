# Agent Development Guide — lazys3

Authoritative reference for all AI agents working in this repository.

---

## Project Overview

`lazys3` is a terminal UI (TUI) for browsing, filtering, and previewing AWS S3 objects.

**Stack:** Go 1.24, Bubble Tea v2 (`charm.land/bubbletea/v2`), Lipgloss, aws-sdk-go-v2, DuckDB (CGO), parquet-go (pure Go), Cobra.

**Package layout:**
```
cmd/lazys3/         — cobra entry point; wires S3 client + TUI
internal/
  ui/               — Bubble Tea models: AppModel, panes, overlays, keybindings
  s3/               — S3 Client interface + AWSClient (aws-sdk-go-v2)
  preview/          — file type detection, JSON/Parquet/DuckDB formatters
```

---

## Build & Run Commands

```sh
make build          # compile binary → ./lazys3
make build-duckdb   # compile with CGO_ENABLED=1 (required for DuckDB)
make run            # go run ./cmd/lazys3 (no binary)
make clean          # remove ./lazys3
make tidy           # go mod tidy
```

---

## Test Commands

```sh
make test                                    # run all tests
make test-verbose                            # go test -v ./...
make test-race                               # with race detector (recommended for TUI code)

# Run a single test by name
go test ./internal/ui/ -run TestNewAppDefaultMode -v

# Run all tests in one package
go test ./internal/s3/ -v
go test ./internal/preview/ -v
go test ./internal/ui/ -v

# Run tests matching a pattern
go test ./... -run TestFilter -v
```

---

## Lint & Format Commands

```sh
make fmt            # gofmt -w . (in-place)
make vet            # go vet ./...

# Check formatting without rewriting (must produce no output before a commit)
gofmt -l .
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

## Bubble Tea v2 API (Critical — differs from v1)

- Module path: `charm.land/bubbletea/v2` (NOT `github.com/charmbracelet/bubbletea`)
- `Init() tea.Cmd` (no return model)
- `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
- `View() tea.View` — use `tea.NewView(string)`, NOT `string`
- Alt screen: `v.AltScreen = true` on the returned `tea.View` struct
- Key events: `tea.KeyPressMsg`, `msg.String()` e.g. `"k"`, `"j"`, `"enter"`, `"esc"`
- Paste events: `tea.PasteMsg{Content string}` — separate from KeyPressMsg
- `bubbles/list` and `bubbles/viewport` are NOT compatible with v2 — use custom pane models

---

## Application Architecture

### Modal state machine

```go
type AppMode int  // ModeNormal | ModeFilter | ModeDownload | ModeHelp | ModeMetadata | ModeQuery
type AppStage int // StageBuckets (buckets+objects peek) | StageObjects (objects+preview)
```

`Update()` routes key events by mode before falling through to `handleKey`.

### Two-stage layout

- **StageBuckets**: `Buckets (50%) | Objects peek (50%)` — j/k on buckets lazily loads object peek
- **StageObjects**: `Objects (50%) | Preview (50%)` — entered via Enter/l on a bucket

### Message types (async S3 ops)

`BucketsLoadedMsg`, `ObjectsLoadedMsg`, `PreviewReadyMsg` (carries `RawPath`+`FileKey` for rescale), `MetadataLoadedMsg`, `DownloadDoneMsg`

### Preview temp file lifecycle

`handlePreviewReady` stores `rawPath` on `PreviewPane`. Rescale (`+`/`-`) and query (`:`) re-use the same path. The file is deleted only when a genuinely new preview replaces it (different `RawPath` in the new `PreviewReadyMsg`).

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

    tea "charm.land/bubbletea/v2"
    "github.com/charmbracelet/lipgloss"

    s3pkg "github.com/jonathan5p/lazys3/internal/s3"
)
```

1. Standard library
2. Third-party packages
3. Internal packages (`github.com/jonathan5p/lazys3/...`)

### Naming

- Unexported: `camelCase`. Exported: `PascalCase`.
- Acronyms stay uppercase: `S3`, `TUI`, `JSON`, `URL`.
- Constructors: `New<Type>` (e.g. `NewApp`, `NewAWSClient`).
- Receiver names: short, consistent per type (`a` for `*AppModel`, `c` for `*Client`).
- Boolean fields: positive framing (`loading`, `active`), not negations.

### Types & Receivers

- Use pointer receivers on structs when any method mutates state.
- Prefer concrete structs over interfaces; introduce an interface only for testing or multiple implementations (e.g. `s3.Client`, `preview.Formatter`).
- No `interface{}` / `any` unless unavoidable.

### Functions & Methods

- Keep every function under **20 lines**. Flag and refactor anything longer.
- Sub-models (`BucketPane`, `ObjectPane`, `PreviewPane`, `StatusBar`) follow the pattern:
  ```go
  func (p PaneType) Update(msg tea.Msg) (PaneType, tea.Cmd)
  func (p PaneType) View(active bool) string
  ```
- Keybinding handlers return `(tea.Model, tea.Cmd)` — never modify app state in-place.

### Error Handling

- Always check errors immediately — never discard with `_`.
- Propagate with `return err` at the lowest level.
- Add context when wrapping: `fmt.Errorf("load aws config: %w", err)`.
- `log.Fatalf` only in `main()`. All other packages return errors to the caller.
- S3/async errors surface via the message types (e.g. `BucketsLoadedMsg{Err: err}`) and are displayed in the status bar via `statusbar.SetError(err)`.

### Comments

- No comments that restate what the code does.
- Doc comments (starting with the symbol name) only on exported identifiers.
- Unexported helpers must be self-explanatory through naming alone.

---

## TDD Rules

- Write the failing test first; confirm it is RED before implementing.
- Use only the standard `testing` package — no third-party assertion libraries.
- Prefer table-driven tests (`[]struct{ ... }`) for multiple input cases.
- Test case names must be descriptive: `"directory with nested prefix"`, not `"test"`.
- Mock S3 calls by implementing the `s3.Client` interface with a fake struct in tests.
- Test mode transitions and state changes — not gocui/Bubble Tea rendering internals.
- After every code change, run `go test ./...` automatically without asking.

---

## DuckDB / CGO

- DuckDB requires `CGO_ENABLED=1` and GCC/Clang at build time.
- Source file: `internal/preview/duckdb.go` — build tag `//go:build cgo`
- Fallback stub: `internal/preview/duckdb_stub.go` — build tag `//go:build !cgo`
- Use `make build-duckdb` to compile with DuckDB support.
- Default `make build` works without explicit CGO flag (CGO on by default when GCC present).

---

## Language

All code, comments, commit messages, test names, and documentation must be in **English**.
