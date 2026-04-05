package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mohokh67/portly/internal/icons"
	"github.com/mohokh67/portly/internal/scanner"
)

// viewMode controls which connections are shown.
type viewMode int

const (
	modeListening viewMode = iota
	modeAll
)

// uiState tracks what the TUI is currently showing.
type uiState int

const (
	stateList    uiState = iota // normal list navigation
	stateConfirm                // kill confirmation prompt
	stateFilter                 // search input active
	stateResult                 // showing kill results
)

// Model is the Bubble Tea model for portly TUI.
type Model struct {
	processes   []scanner.Process
	cursor      int
	selected    map[int]bool // index → selected
	mode        viewMode
	state       uiState
	iconStyle   icons.IconStyle
	filterInput string
	filtered    []int // indices into processes that match filter
	killResults []string
	err         error
	width       int
	height      int
}

// Config holds TUI startup options.
type Config struct {
	IconStyle icons.IconStyle
}

func newModel(cfg Config) Model {
	return Model{
		selected:  make(map[int]bool),
		iconStyle: cfg.IconStyle,
		mode:      modeListening,
		state:     stateList,
	}
}

// Init loads the initial port list.
func (m Model) Init() tea.Cmd {
	return loadPorts(scanner.ListeningOnly)
}
