package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mohokh67/portly/internal/icons"
	"github.com/mohokh67/portly/internal/scanner"
)

var (
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("8")).Padding(0, 1)
	cursorStyle  = lipgloss.NewStyle().BorderLeft(true).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69")).PaddingLeft(1)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	portStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	userStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
	protoTCP     = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
	protoUDP     = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	statusStyle  = lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("8")).Padding(0, 1)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69")).Padding(0, 1)
)

func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("error: %v\n\nPress q to quit.", m.err))
	}

	var b strings.Builder

	modeLabel := "listening"
	if m.mode == modeAll {
		modeLabel = "all connections"
	}
	b.WriteString(titleStyle.Render(fmt.Sprintf("portly — %s", modeLabel)))
	b.WriteString("\n")

	b.WriteString(headerStyle.Render(
		fmt.Sprintf("%-8s %-6s %-24s %-8s %-12s %s",
			"PORT", "PROTO", "PROCESS", "PID", "USER", "ADDRESS"),
	))
	b.WriteString("\n")

	visible := m.visibleIndices()

	if len(visible) == 0 {
		if m.filterInput != "" {
			b.WriteString(dimStyle.Render("  No matching ports\n"))
		} else {
			b.WriteString(dimStyle.Render("  No listening ports\n"))
		}
	}

	for i, idx := range visible {
		p := m.processes[idx]
		icon := icons.Resolve(p.Name, m.iconStyle)
		if icon != "" {
			icon += " "
		}

		proto := protoTCP.Render(p.Proto)
		if p.Proto == "UDP" {
			proto = protoUDP.Render(p.Proto)
		}

		sel := "  "
		if m.selected[idx] {
			sel = selectedStyle.Render("● ")
		}

		line := fmt.Sprintf("%s%-8s %-6s %-24s %-8d %-12s %s",
			sel,
			portStyle.Render(fmt.Sprintf("%d", p.Port)),
			proto,
			icon+p.Name,
			p.PID,
			userStyle.Render(p.User),
			dimStyle.Render(p.Address),
		)

		if i == m.cursor {
			b.WriteString(cursorStyle.Render(line))
		} else {
			b.WriteString("  " + line)
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")

	switch m.state {
	case stateFilter:
		b.WriteString(statusStyle.Render(fmt.Sprintf("/ %s", m.filterInput)) + "  esc to clear\n")
	case stateConfirm:
		targets := m.killTargets()
		prompt := killPrompt(targets)
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true).Render(prompt))
		b.WriteString("\n")
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
		b.WriteString(statusStyle.Render("↑↓ navigate  space select  k kill  t toggle  / search  q quit"))
		b.WriteString("\n")
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
