# Contributing to portly

Thanks for taking the time to contribute!

## Getting started

**Requirements:** Go 1.22+, `lsof` (macOS/Linux)

```bash
git clone https://github.com/mohokh67/portly.git
cd portly
go mod download
go build ./...
go test ./...
```

## Before you open a PR

- **Format:** `gofmt -w .` — CI will reject unformatted code
- **Lint:** `golangci-lint run` (install: `brew install golangci-lint`)
- **Tests:** `go test ./... -race` — all must pass on macOS and Linux
- **New code needs tests** — add them in the same package as what you changed

## Project layout

```
cmd/                    # Cobra commands (root, check, kill)
internal/
  scanner/              # Port scanning (lsof + /proc fallback)
  killer/               # SIGTERM → SIGKILL logic
  icons/                # Process icon lookup + terminal detection
  tui/                  # Bubble Tea TUI (model, update, view)
main.go
```

## Making changes

### Adding a process icon

Edit `internal/icons/icons.go` — add an entry to `table`:

```go
{keys: []string{"myapp"}, nerdFont: "\uXXXX", emoji: "🔧"},
```

Keys are matched case-insensitively as substrings of the process name.

### Adding a CLI command

Create `cmd/<name>.go`, register it with `rootCmd.AddCommand(...)` in an `init()` func. Follow the pattern in `cmd/check.go`.

### Changing TUI behaviour

- Key handling → `internal/tui/update.go`
- Rendering → `internal/tui/view.go`
- State/model fields → `internal/tui/model.go`

## Pull request checklist

- [ ] `go test ./... -race` passes
- [ ] `gofmt -l .` outputs nothing
- [ ] Description explains *why*, not just *what*
- [ ] Linked to a relevant issue if one exists

## Reporting bugs

Open an issue and include:
- OS + architecture (`uname -srm`)
- portly version (`portly --version`)
- Steps to reproduce
- Expected vs actual behaviour

## Feature requests

Open an issue first — discuss before building. Large changes without prior discussion may not be merged.
