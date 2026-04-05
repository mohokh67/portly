package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mohokh67/portly/internal/scanner"
)

// portsLoadedMsg carries freshly scanned processes.
type portsLoadedMsg struct {
	processes []scanner.Process
	err       error
}

// loadPorts is a Bubble Tea command that scans ports asynchronously.
func loadPorts(mode scanner.ScanMode) tea.Cmd {
	return func() tea.Msg {
		procs, err := scanner.Scan(mode)
		return portsLoadedMsg{processes: procs, err: err}
	}
}

// Run starts the TUI.
func Run(cfg Config) error {
	m := newModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("tui error: %w", err)
	}
	return nil
}
