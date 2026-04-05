package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mohokh67/portly/internal/killer"
	"github.com/mohokh67/portly/internal/scanner"
)

type killDoneMsg struct {
	results []string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case portsLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.processes = msg.processes
		m.cursor = 0
		m.selected = make(map[int]bool)
		m = m.withFilter()
		return m, nil

	case killDoneMsg:
		m.killResults = msg.results
		m.state = stateResult
		scanMode := scanner.ListeningOnly
		if m.mode == modeAll {
			scanMode = scanner.AllConnections
		}
		return m, loadPorts(scanMode)

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateFilter:
		return m.handleFilterKey(msg)
	case stateConfirm:
		return m.handleConfirmKey(msg)
	case stateResult:
		m.state = stateList
		m.killResults = nil
		return m, nil
	}
	return m.handleListKey(msg)
}

func (m Model) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visible := m.visibleIndices()
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(visible)-1 {
			m.cursor++
		}

	case "k":
		return m.startKillConfirm()

	case " ":
		if len(visible) > 0 {
			idx := visible[m.cursor]
			m.selected[idx] = !m.selected[idx]
		}

	case "t":
		if m.mode == modeListening {
			m.mode = modeAll
			return m, loadPorts(scanner.AllConnections)
		}
		m.mode = modeListening
		return m, loadPorts(scanner.ListeningOnly)

	case "/":
		m.state = stateFilter
		m.filterInput = ""
	}
	return m, nil
}

func (m Model) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateList
		m.filterInput = ""
		m = m.withFilter()
	case "enter":
		if m.filterInput == "" {
			m.state = stateList
		}
		m = m.withFilter()
	case "backspace":
		if len(m.filterInput) > 0 {
			m.filterInput = m.filterInput[:len(m.filterInput)-1]
		}
		m = m.withFilter()
	default:
		if len(msg.String()) == 1 {
			m.filterInput += msg.String()
			m = m.withFilter()
		}
	}
	return m, nil
}

func (m Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "y":
		return m, m.executeKill()
	default:
		m.state = stateList
	}
	return m, nil
}

func (m Model) startKillConfirm() (Model, tea.Cmd) {
	if len(m.killTargets()) == 0 {
		return m, nil
	}
	m.state = stateConfirm
	return m, nil
}

// killTargets returns the processes to kill: selected ones, or current cursor if none selected.
func (m Model) killTargets() []scanner.Process {
	visible := m.visibleIndices()
	var targets []scanner.Process
	for idx := range m.selected {
		if m.selected[idx] {
			targets = append(targets, m.processes[idx])
		}
	}
	if len(targets) == 0 && len(visible) > 0 {
		targets = []scanner.Process{m.processes[visible[m.cursor]]}
	}
	return targets
}

func (m Model) executeKill() tea.Cmd {
	targets := m.killTargets()
	return func() tea.Msg {
		var results []string
		for _, p := range targets {
			if err := killer.KillPID(p.PID); err != nil {
				results = append(results, fmt.Sprintf("✗ failed %s:%d — %v", p.Name, p.Port, err))
			} else {
				results = append(results, fmt.Sprintf("✓ killed %s (PID %d) on :%d", p.Name, p.PID, p.Port))
			}
		}
		return killDoneMsg{results: results}
	}
}

// withFilter returns a model with filtered updated based on filterInput.
func (m Model) withFilter() Model {
	m.filtered = nil
	q := strings.ToLower(m.filterInput)
	for i, p := range m.processes {
		if q == "" ||
			strings.Contains(strings.ToLower(p.Name), q) ||
			strings.Contains(fmt.Sprintf("%d", p.Port), q) {
			m.filtered = append(m.filtered, i)
		}
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
	return m
}

// visibleIndices returns the process indices currently shown.
func (m Model) visibleIndices() []int {
	if m.filterInput == "" && m.filtered == nil {
		idx := make([]int, len(m.processes))
		for i := range idx {
			idx[i] = i
		}
		return idx
	}
	return m.filtered
}
