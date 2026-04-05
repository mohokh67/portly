# portly — Design Spec

**Date:** 2026-04-04
**Status:** Approved

## Overview

`portly` is a cross-platform CLI (Mac + Linux) for managing ports. It has two modes: an interactive TUI for browsing and killing ports, and direct commands for scripting and quick use.

## Goals

- See all listening ports with process info at a glance
- Kill processes by port interactively or directly
- Distribute via Homebrew and GitHub Releases
- Work on macOS (amd64 + arm64) and Linux (amd64 + arm64)

## Non-Goals

- Windows support (out of scope for v1)
- Network traffic inspection / packet capture
- Port forwarding or remapping

---

## Commands

```
portly                  # open interactive TUI
portly <port>           # show what's on <port>; if occupied, ask "Kill it? [y/N]"
portly check <port>     # show info only, no prompt
portly kill <port>      # kill immediately, no prompt, exit
```

**Cobra routing for `portly <port>`:** root command uses `cobra.ArbitraryArgs`. On run, if exactly one arg is provided and it is a valid integer (1–65535), treat as port lookup. Otherwise print usage and exit 1.

**`portly check <port>` output (non-interactive):**
```
port 3000: node (PID 98234, user: mehr, proto: TCP, addr: 0.0.0.0)
```
If port is free:
```
port 3000: free
```
Exit code 0 in both cases. Pipe-friendly plain text (not a table).

**Port not found / free:**
- `portly <port>` (interactive): prints `Port 3000 is free.` and exits 0
- `portly check <port>`: prints `port 3000: free` and exits 0
- `portly kill <port>`: prints `Nothing on port 3000.` and exits 2

**Flags (global):**

```
--icons=nerdfont|emoji|none    # icon style (default: nerdfont, falls back to emoji)
--no-color                     # disable color output
--version                      # print version and exit
```

**Exit codes:**

| Code | Meaning |
|------|---------|
| 0 | Success (port free, process killed, info shown) |
| 1 | General error (permission denied, lsof missing) |
| 2 | Port not found / nothing listening |

---

## Architecture

Three packages:

### 1. `scanner` — port data

**macOS:** calls `lsof -i -P -n`.

**Linux:** tries `lsof -i -P -n` first. If `lsof` is not installed, falls back to parsing `/proc/net/tcp` and `/proc/net/tcp6` directly (reads hex address/port pairs, resolves PID via `/proc/<pid>/net/tcp`). If both fail, exits with error: `portly requires lsof or /proc/net/tcp (Linux kernel 2.6+)`.

Parses output into:

```go
type Process struct {
    Port    int
    Proto   string // TCP | UDP
    PID     int
    User    string
    Address string // e.g. 0.0.0.0 or 127.0.0.1
    Name    string // process name
}
```

Two modes: `ListeningOnly` (default) and `AllConnections` (includes established).

### 2. `tui` — interactive UI

Built with **Bubble Tea** (model/update/view) + **Lipgloss** (styling).

**Columns (full layout):**

| PORT | PROTO | PROCESS | PID | USER | ADDRESS |

- Process column prefixed with icon (Nerd Font or emoji)
- Per-process color coding
- Selected rows highlighted
- Status bar at bottom: key hints + current mode

**Controls:**

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate |
| `Space` | Toggle row selection |
| `k` | Kill selected rows (or current if none selected) — shows confirmation prompt |
| `t` | Toggle listening only / all connections |
| `/` | Search / filter by port or process name |
| `q` / `Ctrl+C` | Quit |

**Kill flow in TUI:**
1. Press `k`
2. Single selected: `Kill node (PID 98234) on :3000? [y/N]`
   Multiple selected: `Kill 3 processes (node:3000, postgres:5432, redis:6379)? [y/N]`
3. On `y`: for each process, send SIGTERM → poll 2s → SIGKILL if still alive
4. Show inline result per process: `✓ killed node:3000` or `✗ failed (permission denied)`
5. Refresh list, deselect all

**Search/filter (`/`):**
- Opens an inline text input at the bottom of the TUI (replaces key hint bar)
- Filters in real-time as you type — matches port number or process name (case-insensitive, partial match)
- `Esc` or empty + `Enter` dismisses filter and restores full list
- If no results: show `No matching ports` in list area
- Filtered view preserves selection state

### 3. `cmd` — CLI entrypoint

Built with **Cobra**. Registers root command and subcommands (`check`, `kill`). Passes flags down to scanner and tui.

---

## Icon System

Lookup table mapping process names → icons. Ships with ~20 common entries:

| Process | Nerd Font | Emoji |
|---------|-----------|-------|
| node / npm | `` (green) | 🟢 |
| postgres | `` (blue) | 🐘 |
| docker-proxy | `` (blue) | 🐳 |
| redis-server | `` (red) | ⚡ |
| mongod | `` (green) | 🍃 |
| nginx | `` (green) | 🌐 |
| python | `` (yellow) | 🐍 |
| ruby | `` (red) | 💎 |
| java | `` (orange) | ☕ |
| go / air | `` (cyan) | 🔵 |
| unknown | `` | ⚙️ |

Detection: check `$TERM_PROGRAM` and `$TERM` against a known list of Nerd Font-friendly terminals (iTerm.app, WezTerm, Alacritty, kitty, ghostty). Also check `$NERD_FONTS=1` env var as an explicit opt-in. If none match, fall back to emoji. User override via `--icons` flag always wins.

Detection order:
1. `--icons` flag (highest priority)
2. `$NERD_FONTS=1` env var
3. `$TERM_PROGRAM` / `$TERM` known list
4. Emoji fallback

README includes a section on installing a Nerd Font (recommended: JetBrains Mono NF or FiraCode NF).

---

## Kill Behavior

- **TUI `k`**: confirmation prompt → SIGTERM → poll 2s → SIGKILL if still alive → show result inline
- **`portly kill <port>`**: SIGTERM → poll 2s → SIGKILL, no prompt → print `killed node (PID 98234) on :3000` or error
- **`portly <port>` interactive**: inline `[y/N]` prompt, same kill sequence
- **Already dead before 2s**: detected via polling PID existence; skip SIGKILL, report success
- **SIGKILL fails**: print `failed to kill PID 98234: <reason>`, exit 1
- **Permission denied**: print `permission denied — try running with sudo`, exit 1

---

## Distribution

### Homebrew (Mac + Linux)
- Separate tap repo: `github.com/<owner>/homebrew-portly`
- Formula lives in that repo (not inside the main `portly` repo)
- Formula points to GitHub Release binaries (generated by GoReleaser)
- Install: `brew tap <owner>/portly && brew install portly`

### GitHub Releases (via GoReleaser)
- Triggered on git tag (`v*`)
- Builds: `darwin/amd64`, `darwin/arm64`, `linux/amd64`, `linux/arm64`
- Archives: `.tar.gz` with binary + README
- Checksums file included

### Direct install script (Linux)
```bash
curl -sfL https://raw.githubusercontent.com/<owner>/portly/main/install.sh | sh
```
Detects arch, downloads correct binary from latest release, installs to `/usr/local/bin`.

---

## Project Structure

```
portly/                  # main repo
├── cmd/
│   ├── root.go         # portly (no args) → TUI; portly <port> → lookup
│   ├── check.go        # portly check <port>
│   └── kill.go         # portly kill <port>
├── internal/
│   ├── scanner/
│   │   ├── scanner.go  # dispatcher: lsof or /proc fallback
│   │   ├── lsof.go     # lsof-based implementation
│   │   └── proc.go     # /proc/net/tcp fallback (Linux)
│   ├── tui/
│   │   ├── model.go    # Bubble Tea model
│   │   ├── update.go   # key handling, state transitions
│   │   └── view.go     # rendering, Lipgloss styles
│   └── icons/
│       └── icons.go    # lookup table, detection logic
├── main.go
├── .goreleaser.yaml
└── install.sh          # Linux direct install script

homebrew-portly/         # separate tap repo
└── Formula/
    └── portly.rb
```

---

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/charmbracelet/bubbletea` | TUI framework |
| `github.com/charmbracelet/lipgloss` | TUI styling |
| `github.com/spf13/cobra` | CLI framework |
| `github.com/charmbracelet/bubbles` | Table, textinput components |

Go stdlib only for `lsof` invocation and process killing (`os/exec`, `syscall`).

---

## Open Questions (resolved)

- Language: **Go**
- Binary name: **portly**
- Default icon style: **Nerd Font, emoji fallback**
- Default port view: **listening only, `t` to toggle**
- Kill UX in TUI: **Space to select, `k` to kill with confirmation**
- Non-interactive default (`portly 3000`): **show info + ask [y/N]**
