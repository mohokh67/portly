# portly

A CLI for managing ports — list what's running, kill by port, or browse interactively.

## Install

### Homebrew (Mac + Linux)
```bash
brew tap mohokh67/portly
brew install portly
```

### Linux (direct)
```bash
curl -sfL https://raw.githubusercontent.com/mohokh67/portly/main/install.sh | sh
```

### From source
```bash
go install github.com/mohokh67/portly@latest
```

## Usage

```
portly                  # interactive TUI — browse and kill ports
portly 3000             # show what's on port 3000, ask to kill
portly check 3000       # show info only
portly kill 3000        # kill immediately, no prompt
```

### Flags

```
--icons=nerdfont|emoji|none|auto    # icon style (default: auto-detect)
--no-color                          # disable color output
--version                           # print version and exit
```

## TUI Controls

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate |
| `Space` | Select / deselect row |
| `k` | Kill selected (or current row) — shows confirm prompt |
| `t` | Toggle listening only / all connections |
| `/` | Search by port number or process name |
| `q` / `Ctrl+C` | Quit |

## Icons

portly shows icons next to process names. For the best experience install a [Nerd Font](https://www.nerdfonts.com/font-downloads) (recommended: **JetBrains Mono NF** or **FiraCode Nerd Font**) and set it as your terminal font.

| Process | Emoji | Nerd Font |
|---------|-------|-----------|
| node / npm | 🟢 | |
| postgres | 🐘 | |
| docker | 🐳 | |
| redis | ⚡ | |
| nginx | 🌐 | |
| python | 🐍 | |
| ruby / rails | 💎 | |
| java | ☕ | |
| go / air | 🔵 | |

Override icon style:
```bash
portly --icons=nerdfont   # force Nerd Font
portly --icons=emoji      # use emoji
portly --icons=none       # no icons
```

Set `NERD_FONTS=1` in your shell profile to always use Nerd Font icons regardless of terminal detection.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (permission denied, lsof missing) |
| 2 | Port not found / nothing listening |

## Platform Support

- macOS (amd64 + arm64)
- Linux (amd64 + arm64) — uses `lsof` or falls back to `/proc/net/tcp`
