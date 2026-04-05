package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mohokh67/portly/internal/icons"
	"github.com/mohokh67/portly/internal/scanner"
)

const (
	colPort    = 7
	colProto   = 5
	colPID     = 7
	colUser    = 12
	colAddress = 15
	colSel     = 2
	// total fixed chars: sel + port + proto + pid + user + address + 5 spaces between 6 cols
	colsFixed = colSel + colPort + colProto + colPID + colUser + colAddress + 5
)

var (
	headerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("63")).
			Foreground(lipgloss.Color("255")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("63")).
			Foreground(lipgloss.Color("255")).
			Bold(true).
			Padding(0, 1)

	filterStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("255"))

	// cursor row: subtle dark highlight, white text
	cursorRowStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("237")).
			Foreground(lipgloss.Color("255"))

	// selected row: bold yellow
	selectedRowStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("220")).
				Foreground(lipgloss.Color("0"))

	portStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	userStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
	protoTCPStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
	protoUDPStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69"))
	versionStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func (m Model) termW() int {
	if m.termWidth == 0 {
		return 80
	}
	return m.termWidth
}

func (m Model) processColWidth() int {
	w := m.termW() - colsFixed - 1 // 1 extra sep before process col
	if w < 16 {
		w = 16
	}
	return w
}

// truncate clips s to at most maxW visible columns.
func truncate(s string, maxW int) string {
	for lipgloss.Width(s) > maxW {
		runes := []rune(s)
		s = string(runes[:len(runes)-1])
	}
	return s
}

// renderRow builds a full-width row string. All cells use lipgloss Width()
// which is ANSI-aware, so colored content never breaks column alignment.
func (m Model) renderRow(p scanner.Process, idx int, isCursor bool) string {
	isSelected := m.selected[idx]
	colProcess := m.processColWidth()

	icon := icons.Resolve(p.Name, m.iconStyle)
	label := p.Name
	if icon != "" {
		label = icon + " " + p.Name
	}
	label = truncate(label, colProcess)

	port := fmt.Sprintf("%d", p.Port)
	pid := fmt.Sprintf("%d", p.PID)
	user := truncate(p.User, colUser)
	addr := truncate(p.Address, colAddress)

	var s lipgloss.Style
	switch {
	case isSelected:
		s = selectedRowStyle
	case isCursor:
		s = cursorRowStyle
	default:
		var protoCell string
		if p.Proto == "UDP" {
			protoCell = protoUDPStyle.Width(colProto).Render(p.Proto)
		} else {
			protoCell = protoTCPStyle.Width(colProto).Render(p.Proto)
		}
		return "  " +
			portStyle.Width(colPort).Render(port) + " " +
			protoCell + " " +
			lipgloss.NewStyle().Width(colProcess).Render(label) + " " +
			lipgloss.NewStyle().Width(colPID).Render(pid) + " " +
			userStyle.Width(colUser).Render(user) + " " +
			dimStyle.Width(colAddress).Render(addr)
	}

	return s.Width(colSel).Render("") +
		s.Width(colPort).Render(port) + " " +
		s.Width(colProto).Render(p.Proto) + " " +
		s.Width(colProcess).Render(label) + " " +
		s.Width(colPID).Render(pid) + " " +
		s.Width(colUser).Render(user) + " " +
		s.Width(colAddress).Render(addr)
}

func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("error: %v\n\nPress q to quit.", m.err))
	}

	w := m.termW()
	var b strings.Builder

	// ── title bar ─────────────────────────────────────────────────
	modeLabel := "listening"
	if m.mode == modeAll {
		modeLabel = "all connections"
	}
	ver := m.version
	if ver == "" {
		ver = "dev"
	}
	left := titleStyle.Render("portly") + dimStyle.Render("  "+modeLabel)
	right := versionStyle.Render("v" + ver)
	gap := w - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	b.WriteString(left + strings.Repeat(" ", gap) + right + "\n")

	// ── header ────────────────────────────────────────────────────
	// Built identically to data rows: sel+port+space+proto+space+...
	// (no separator between sel and port, matching renderRow layout)
	colProcess := m.processColWidth()
	h := headerStyle
	sp := h.Render(" ")
	headerRow := h.Width(colSel).Render("") +
		h.Width(colPort).Render("PORT") + sp +
		h.Width(colProto).Render("PROTO") + sp +
		h.Width(colProcess).Render("PROCESS") + sp +
		h.Width(colPID).Render("PID") + sp +
		h.Width(colUser).Render("USER") + sp +
		h.Width(colAddress).Render("ADDRESS")
	b.WriteString(headerRow + "\n")

	// ── rows ──────────────────────────────────────────────────────
	visible := m.visibleIndices()
	if len(visible) == 0 {
		msg := "  No listening ports"
		if m.filterInput != "" {
			msg = "  No matching ports"
		}
		b.WriteString(dimStyle.Render(msg) + "\n")
	}
	for i, idx := range visible {
		b.WriteString(m.renderRow(m.processes[idx], idx, i == m.cursor) + "\n")
	}

	// ── padding ───────────────────────────────────────────────────
	used := 2 + len(visible) + 2 // title + header + rows + footer + spare
	if m.termHeight > 0 && m.termHeight > used {
		b.WriteString(strings.Repeat("\n", m.termHeight-used))
	} else {
		b.WriteString("\n")
	}

	// ── footer ────────────────────────────────────────────────────
	switch m.state {
	case stateFilter:
		b.WriteString(filterStyle.Width(w).Render(fmt.Sprintf(" / %s", m.filterInput)) + "\n")
	case stateConfirm:
		prompt := killPrompt(m.killTargets())
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true).Width(w).Render(prompt) + "\n")
	case stateResult:
		for _, r := range m.killResults {
			if strings.HasPrefix(r, "✓") {
				b.WriteString(successStyle.Render(r) + "\n")
			} else {
				b.WriteString(errorStyle.Render(r) + "\n")
			}
		}
		b.WriteString(dimStyle.Render("press any key to continue") + "\n")
	default:
		b.WriteString(statusStyle.Width(w).Render(" ↑↓/jk navigate   space select   x kill   t toggle   / search   q quit") + "\n")
	}

	return b.String()
}

func killPrompt(targets []scanner.Process) string {
	if len(targets) == 1 {
		p := targets[0]
		return fmt.Sprintf("Kill %s (PID %d) on :%d? [y/N] ", p.Name, p.PID, p.Port)
	}
	names := make([]string, len(targets))
	for i, p := range targets {
		names[i] = fmt.Sprintf("%s:%d", p.Name, p.Port)
	}
	return fmt.Sprintf("Kill %d processes (%s)? [y/N] ", len(targets), strings.Join(names, ", "))
}
