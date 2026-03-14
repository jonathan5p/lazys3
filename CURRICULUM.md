# lazys3 — Go Learning Curriculum

Rebuild this TUI from scratch, one milestone at a time.
Each milestone introduces specific Go concepts. By the end, the result matches the current `lazys3` codebase.

**Background:** Advanced Python developer, some Go exposure, familiar with goroutines at a surface level.
**Pace:** ~5 hours/week. Estimated total: ~11 weeks.
**Method:** TDD. I write the first test per new pattern; student writes subsequent tests in the same style.

---

## Progress Key

- `[ ]` Not started
- `[~]` In progress
- `[x]` Complete

---

## Milestone 1 — The Language Foundation

> Go syntax, types, functions, packages, modules

- [ ] Set up a new Go module (`go mod init`)
- [ ] Write and run a Hello World program using `fmt`
- [ ] Learn value types, pointers, and the difference between them
- [ ] Define structs and methods on them
- [ ] Understand `iota` for enumerations
- [ ] Write table-driven tests with the `testing` package
- [ ] **Deliverable:** A standalone `Object` struct (key, size, lastModified) with a `String()` method and passing tests

**Go concepts introduced:** `package`, `import`, `func`, `struct`, `iota`, `*T` vs `T`, `go test`, table-driven tests

---

## Milestone 2 — Interfaces and Dependency Injection

> Go interfaces, duck typing, mocking in tests

- [ ] Learn how Go interfaces work (implicit satisfaction — no `implements` keyword)
- [ ] Define a `Client` interface with `ListBuckets` and `ListObjects`
- [ ] Implement a `FakeClient` in tests (Go's answer to `unittest.mock`)
- [ ] Write tests that inject `FakeClient` to test logic without hitting AWS
- [ ] **Deliverable:** `internal/s3/client.go` with the `Client` interface, a real `AWSClient` stub, and a `FakeClient` used in tests

**Go concepts introduced:** `interface`, implicit satisfaction, interface-based mocking, `_test.go` files

---

## Milestone 3 — Error Handling

> Go error values, wrapping, propagation

- [ ] Learn why Go uses explicit `error` returns instead of exceptions
- [ ] Understand `fmt.Errorf("context: %w", err)` for wrapping
- [ ] Implement `ListBuckets` and `ListObjects` on `AWSClient` with proper error propagation
- [ ] Write tests that exercise error paths using `FakeClient`
- [ ] **Deliverable:** `internal/s3/client.go` fully implemented with error handling and passing tests

**Go concepts introduced:** `error` type, `fmt.Errorf + %w`, `errors.Is`, `errors.As`, nil checks

---

## Milestone 4 — Goroutines and Channels

> goroutines, channels, `context.Context`

- [ ] Learn what a goroutine is and why it differs from Python `asyncio`
- [ ] Understand `context.Context` as a cancellation/timeout mechanism
- [ ] Write a simple goroutine that fetches buckets and sends the result back via a channel
- [ ] **Deliverable:** A standalone program that fetches bucket names concurrently and prints them

**Go concepts introduced:** `go func()`, `chan`, `context.Context`, `context.Background()`, `context.WithTimeout`

---

## Milestone 5 — The TUI Shell

> Bubble Tea model/update/view pattern, messages, commands

- [ ] Understand the Elm Architecture: `Model`, `Update(msg) → (Model, Cmd)`, `View() → string`
- [ ] Create an empty `AppModel` with `Init`, `Update`, `View`
- [ ] Render a static string — run the TUI and see it in the terminal
- [ ] Handle `WindowSizeMsg` to track terminal dimensions
- [ ] Handle `q` / `ctrl+c` to quit
- [ ] Write tests for `Update` with synthetic messages
- [ ] **Deliverable:** A minimal `cmd/lazys3/main.go` + `internal/ui/app.go` that boots a TUI, shows dimensions, and quits cleanly

**Go concepts introduced:** Bubble Tea v2 API (`Init/Update/View`), `tea.Cmd`, `tea.Msg`, type switch, `tea.KeyPressMsg`

---

## Milestone 6 — The Status Bar

> Value receiver patterns, immutable updates, lipgloss styling

- [ ] Learn the difference between value and pointer receivers in Go
- [ ] Understand why Bubble Tea sub-models use value receivers and return copies
- [ ] Build `StatusBar` struct with `SetMessage`, `SetError`, `SetMode`, `View`
- [ ] Apply basic lipgloss styling (color, width)
- [ ] Write tests asserting `View()` output contains expected strings
- [ ] **Deliverable:** A TUI that shows a styled status bar at the bottom with a static message

**Go concepts introduced:** value vs pointer receivers, lipgloss `Style`, immutable update pattern

---

## Milestone 7 — The Bucket Pane

> Slices, cursor navigation, rendering lists, async message dispatch

- [ ] Learn Go slices (length vs capacity, `append`, `make`)
- [ ] Build `BucketPane` with `items`, `cursor`, `SetItems`, `Update` (j/k navigation), `View`
- [ ] Dispatch a `BucketsLoadedMsg` asynchronously from `Init` using a `tea.Cmd`
- [ ] Handle `BucketsLoadedMsg` in `AppModel.Update` to populate the pane
- [ ] Write tests for cursor clamping and navigation
- [ ] **Deliverable:** TUI that loads and displays real or fake bucket names; navigate up/down

**Go concepts introduced:** slices, `make`, `append`, `tea.Cmd` as a lazy message producer, async dispatch

---

## Milestone 8 — The Object Pane and Two-Stage Navigation

> Maps, structs with multiple fields, stage/mode state machine

- [ ] Build `ObjectPane` similarly to `BucketPane`
- [ ] Introduce `AppStage` (`StageBuckets` / `StageObjects`) using `iota`
- [ ] Wire `Enter` on a bucket to trigger `ListObjects` and transition to `StageObjects`
- [ ] Show a peek of objects in stage 1, full navigation in stage 2
- [ ] Write tests for stage transitions triggered by key events
- [ ] **Deliverable:** TUI that navigates from bucket list into object list and back

**Go concepts introduced:** `map`, multi-field struct rendering, state machine with `iota` modes

---

## Milestone 9 — The Filter Mode

> String manipulation, mode-based key routing, sub-state

- [ ] Introduce `AppMode` enum and mode-based key dispatch in `Update`
- [ ] Build `FilterState` with live filtering logic (`strings.Contains`)
- [ ] Connect filter to both `BucketPane` and `ObjectPane`
- [ ] Write tests that send key presses and assert filtered results
- [ ] **Deliverable:** Press `/` to enter filter mode; typing narrows the list in real time

**Go concepts introduced:** `strings` package, mode dispatch pattern, sub-state structs

---

## Milestone 10 — The Preview Pane

> `io.Reader`, file I/O, temp files, goroutines for S3 streaming

- [ ] Learn `io.Reader` / `io.ReadCloser` — Go's streaming abstraction
- [ ] Build `PreviewPane` with scrollable content and `+`/`-` zoom
- [ ] Implement `GetObject` on `AWSClient` and stream content to a temp file
- [ ] Dispatch `PreviewReadyMsg` asynchronously; handle it to populate the pane
- [ ] Manage temp file lifecycle (delete on replacement)
- [ ] Write tests for preview state transitions
- [ ] **Deliverable:** Select an object and see its raw content in the preview pane

**Go concepts introduced:** `io.Reader`, `os.CreateTemp`, `defer`, `io.Copy`, temp file management

---

## Milestone 11 — File Type Detection and Formatters

> Interfaces for polymorphism, build tags, JSON parsing

- [ ] Learn how to use an interface to select a formatter at runtime
- [ ] Implement a `Formatter` interface with `Format(path) (string, error)`
- [ ] Build `JSONFormatter`: pretty-print JSON using `encoding/json`
- [ ] Build file type detector using file extension and magic bytes
- [ ] Introduce Go build tags (`//go:build`) for optional DuckDB/Parquet support
- [ ] Write tests for each formatter with fixture files
- [ ] **Deliverable:** JSON files render as pretty-printed in the preview pane

**Go concepts introduced:** `encoding/json`, `bytes.HasPrefix` for magic bytes, build tags, `//go:build`

---

## Milestone 12 — Overlays: Metadata and Download Prompt

> Struct composition, overlay rendering, cursor-in-input pattern

- [ ] Build `MetadataOverlay` — fetch and display S3 object metadata
- [ ] Build `DownloadPrompt` — text input with cursor, backspace, left/right movement
- [ ] Implement Tab-based directory autocomplete using `os.ReadDir`
- [ ] Write tests for all input behaviors (insert, delete, cursor, tab-complete)
- [ ] **Deliverable:** Press `m` for metadata overlay; press `d` for download prompt with working input

**Go concepts introduced:** `os.ReadDir`, string slicing for cursor manipulation, overlay rendering with lipgloss

---

## Milestone 13 — Help Overlay and SQL Query Mode

> String building, multi-mode state, wrapping complex state in structs

- [ ] Build `HelpOverlay` that lists all keybindings in a formatted box
- [ ] Build `QueryState` and `ModeQuery` for SQL query input
- [ ] Wire `:` to enter query mode and execute a DuckDB query on the current temp file
- [ ] **Deliverable:** Press `?` for help; press `:` to run a SQL query on a Parquet/CSV file

**Go concepts introduced:** `strings.Builder`, multi-mode state management, CGO integration pattern

---

## Milestone 14 — The CLI Entry Point

> Cobra, `context.Context` lifecycle, wiring everything together

- [ ] Learn how Cobra structures CLI commands with flags
- [ ] Wire `--region`, `--endpoint`, `--profile` flags to `NewAWSClient`
- [ ] Write a smoke test for the root command
- [ ] **Deliverable:** `lazys3 --region us-east-1` boots the full app — feature-complete, matching current state

**Go concepts introduced:** Cobra `Command`, `PersistentFlags`, `RunE`, `os.Exit` conventions

---

## Reference: Current File Structure

```
cmd/lazys3/
  main.go           — Cobra entry point; wires S3 client + TUI
  main_test.go

internal/
  s3/
    client.go       — Client interface + AWSClient (aws-sdk-go-v2)
    client_test.go
    download.go
    types.go

  preview/
    detector.go     — File type detection (extension + magic bytes)
    formatter.go    — Formatter interface
    json.go         — JSON pretty-printer
    parquet.go      — Parquet formatter (pure Go)
    duckdb.go       — DuckDB formatter (CGO, build tag: cgo)
    duckdb_stub.go  — Stub when CGO disabled (build tag: !cgo)
    *_test.go

  ui/
    app.go          — AppModel: Init / Update / View, message handlers
    app_test.go
    buckets.go      — BucketPane
    objects.go      — ObjectPane
    preview_pane.go — PreviewPane
    statusbar.go    — StatusBar
    overlays.go     — MetadataOverlay, DownloadPrompt, HelpOverlay
    filter.go       — FilterState
    keys.go         — Key handlers per mode
    layout.go       — renderLayout()
    messages.go     — Message types
    styles.go       — Lipgloss styles
    query.go        — QueryState
    *_test.go
```

---

## Stack Reference

| Concern | Library |
|---|---|
| TUI framework | `charm.land/bubbletea/v2` |
| Styling | `github.com/charmbracelet/lipgloss` |
| AWS S3 | `github.com/aws/aws-sdk-go-v2` |
| CLI flags | `github.com/spf13/cobra` |
| Parquet | `github.com/parquet-go/parquet-go` |
| DuckDB (optional) | `github.com/marcboeker/go-duckdb` (CGO) |
