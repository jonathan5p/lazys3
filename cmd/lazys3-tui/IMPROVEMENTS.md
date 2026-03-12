# TUI Improvement Backlog

A prioritized list of planned enhancements for the lazys3 TUI.

---

## 1. View Highlighting on Focus Switch

When the user switches between panels (bucket list, object list, preview, etc.),
the currently active view should be visually distinguished from inactive ones.
This mirrors the behavior of tools like `lazygit` and `ranger`, where a colored
border or highlighted title makes the active panel immediately obvious.

---

## 2. Pagination

Long bucket and object lists should be paginated rather than rendered in a single
unbounded scroll. Each panel should show a fixed page of results and expose
keybindings to move forward and backward through pages. This keeps memory usage
predictable and render time consistent regardless of bucket size.

---

## 3. In-Memory Caching

S3 API calls for bucket lists and object listings should be cached in memory for
the lifetime of the session. Revisiting a bucket or prefix that was already
fetched should return the cached result instantly without hitting the network.
A manual refresh keybinding (e.g. `r`) should allow the user to invalidate the
cache for the current view when fresh data is needed.

---

## 4. Parquet File Preview

When the cursor lands on a `.parquet` file, the preview panel should display a
small sample of rows (e.g. the first 10-20 rows) rendered as a table. The file
should be downloaded to a temporary location on demand, decoded using a Go
Parquet reader, and the result displayed inline. Only the row sample is held in
memory; the full file is not retained.

---

## 5. Search by Bucket Name

A search mode triggered by a keybinding (e.g. `/`) should filter the bucket list
in real time as the user types. The filter should be case-insensitive and match
any substring of the bucket name. Pressing `Escape` or `Enter` should exit search
mode and leave the list filtered or fully restored respectively.

---

## 6. Search by S3 Path

A separate search mode for the object list should allow the user to filter objects
by key prefix or substring within the current bucket. This is distinct from the
S3 `ListObjects` prefix parameter: the filter operates on already-fetched results
client-side for instant feedback, and a deeper S3-level prefix search is issued
only when the user explicitly submits the query.

---

## 7. Keybindings Inspired by lazygit and Neovim

The keybinding model should feel familiar to users of `lazygit` and Neovim:

| Key | Action |
|-----|--------|
| `j` / `k` | Move cursor down / up |
| `gg` | Jump to top of list |
| `G` | Jump to bottom of list |
| `Ctrl-d` / `Ctrl-u` | Scroll half-page down / up |
| `Enter` | Confirm / expand |
| `Backspace` | Go up one prefix level |
| `/` | Enter search mode |
| `Escape` | Cancel / exit current mode |
| `r` | Refresh current view (invalidates cache) |
| `?` | Toggle keybinding help overlay |
| `q` | Quit |

All keybindings should be documented in a help overlay accessible via `?`.

---

## 8. Download to Custom Path

When the user initiates a download (e.g. `d`), a small input prompt should appear
at the bottom of the screen asking for a destination path. The field should be
pre-filled with a sensible default (e.g. the current working directory). The user
can edit the path before confirming with `Enter` or cancel with `Escape`. Progress
should be reported inline in the status bar for large files.
