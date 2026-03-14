# lazys3

A fast, keyboard-driven terminal UI for browsing and querying AWS S3 — inspired by lazygit and Neovim.

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/github/license/jonathan5p/lazys3)
![CI](https://github.com/jonathan5p/lazys3/actions/workflows/ci.yml/badge.svg)

---

## Features

- **Two-stage navigation** — bucket list with live object peek → object browser with file preview
- **Fuzzy filter** — `/` to filter buckets or objects in real time; paste a full `s3://` URL to jump directly to any object
- **File preview** — JSON (syntax-highlighted), Parquet (ASCII table, column-width fitted to pane)
- **SQL queries** — press `:` in the preview pane to run DuckDB SQL against the previewed file (`SELECT * FROM data WHERE ...`)
- **Column resize** — `+`/`-` in preview to widen or narrow table columns
- **Download** — `d` opens a prompt for the destination path
- **Metadata** — `m` shows HeadObject metadata for any file
- **Neovim/lazygit keybindings** — `j`/`k`, `h`/`l`, `g`/`G`, `Ctrl-d`/`Ctrl-u`, `?` for help

---

## Installation

### Prerequisites

- **Go 1.24+** — [install Go](https://go.dev/dl/)
- **GCC or Clang** — required for DuckDB SQL support (`gcc --version` to check)
- **AWS credentials** — configured via `~/.aws/credentials`, environment variables, or IAM role

### Option 1 — `go install` (recommended)

```sh
go install github.com/jonathan5p/lazys3/cmd/lazys3@latest
```

The binary is placed in `$GOPATH/bin` (usually `~/go/bin`). Make sure that directory is in your `$PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Option 2 — Download a pre-built binary

Download the latest release for your platform from the [Releases page](https://github.com/jonathan5p/lazys3/releases):

```sh
# Linux (amd64)
curl -L https://github.com/jonathan5p/lazys3/releases/latest/download/lazys3-linux-amd64 \
  -o lazys3 && chmod +x lazys3 && sudo mv lazys3 /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/jonathan5p/lazys3/releases/latest/download/lazys3-darwin-arm64 \
  -o lazys3 && chmod +x lazys3 && sudo mv lazys3 /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/jonathan5p/lazys3/releases/latest/download/lazys3-darwin-amd64 \
  -o lazys3 && chmod +x lazys3 && sudo mv lazys3 /usr/local/bin/
```

### Option 3 — Build from source

```sh
git clone https://github.com/jonathan5p/lazys3.git
cd lazys3
make build-duckdb   # CGO_ENABLED=1, includes DuckDB SQL support
# or
make build          # CGO auto-detected; DuckDB included if GCC present
```

The binary is output to `./lazys3` in the project root.

---

## AWS Credentials

lazys3 uses the standard AWS credential chain — no configuration needed if you already use the AWS CLI:

1. Environment variables: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`
2. `~/.aws/credentials` and `~/.aws/config`
3. EC2/ECS instance role

---

## Usage

```sh
lazys3                          # use default AWS profile and region
lazys3 --profile staging        # use a named AWS profile
lazys3 --region eu-west-1       # override region
lazys3 --endpoint http://localhost:9000  # use custom endpoint (MinIO, LocalStack)
```

---

## Keybindings

| Key                 | Action                                                            |
| ------------------- | ----------------------------------------------------------------- |
| `j` / `k`           | Move cursor down / up                                             |
| `h` / `←`           | Focus left pane / go back                                         |
| `l` / `→` / `Enter` | Enter bucket or focus right pane                                  |
| `Esc`               | Go up prefix level / return to bucket list                        |
| `/`                 | Filter mode — type to filter, `Enter` to confirm, `Esc` to cancel |
| `p`                 | Preview selected file                                             |
| `d`                 | Download selected file (prompts for path)                         |
| `m`                 | Show object metadata                                              |
| `:`                 | SQL query mode (DuckDB) — requires preview to be open             |
| `+` / `=`           | Widen preview table columns                                       |
| `-`                 | Narrow preview table columns                                      |
| `r`                 | Refresh current pane                                              |
| `g` `g`             | Jump to top                                                       |
| `G`                 | Jump to bottom                                                    |
| `?`                 | Toggle keybinding help overlay                                    |
| `q`                 | Quit                                                              |

### Filter mode tips

- Type any substring to narrow the list in real time
- Paste a full `s3://bucket/path/to/file.parquet` URL and press `Enter` to jump directly to that object
- Type a prefix with `/` (e.g. `data/2024/`) and press `Enter` to navigate to that S3 prefix

### SQL query mode

Press `:` while viewing a Parquet or JSON file to open the SQL prompt. The table alias is always `data`:

```sql
SELECT col1, col2 FROM data WHERE col3 > 100 LIMIT 50
SELECT COUNT(*) FROM data
SELECT DISTINCT category FROM data ORDER BY 1
```

Press `Enter` to execute, `Esc` to cancel. Use `←`/`→` to move the cursor, `Ctrl+A`/`Ctrl+E` for line start/end.

---

## Development

See [AGENTS.md](./AGENTS.md) for the full agent dev guide including build commands, code style, and testing conventions.

```sh
make test           # run all tests
make pre-commit     # fmt + vet + tests (run before every commit)
make help           # list all make targets
```

---

## License

MIT
