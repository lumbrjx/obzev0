package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── Styles ─────────────────────────────────────────────────────────────────

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	selectedItem = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	dimItem      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	statusOK     = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	statusErr    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// ── State machine ──────────────────────────────────────────────────────────

type tuiState int

const (
	stateMenu      tuiState = iota // arrow-key chaos type menu
	stateConfigure                 // text inputs for the selected type
	stateRunning                   // spinner while experiment is active
	stateStopped                   // summary after stop
)

var chaosTypes = []string{
	"Latency Injection",
	"Packet Drop / Corrupt",
	"TCP RST Reset",
	"DNS Chaos",
	"Bandwidth Throttle",
	"HTTP Fault Injection",
}

// configFields returns the label+placeholder pairs for each chaos type.
func configFields(idx int) [][2]string {
	switch idx {
	case 0: // Latency
		return [][2]string{
			{"gRPC address", "127.0.0.1:50051"},
			{"Request delay (ms)", "200"},
			{"Response delay (ms)", "200"},
		}
	case 1: // Packet
		return [][2]string{
			{"gRPC address", "127.0.0.1:50051"},
			{"Drop rate (0.0-1.0)", "0.1"},
			{"Corrupt rate (0.0-1.0)", "0.0"},
			{"Duration (seconds)", "30"},
		}
	case 2: // TCP RST
		return [][2]string{
			{"gRPC address", "127.0.0.1:50051"},
			{"Listen address", ":9080"},
			{"Reset rate (0.0-1.0)", "0.5"},
		}
	case 3: // DNS Chaos
		return [][2]string{
			{"gRPC address", "127.0.0.1:50051"},
			{"Mode (nxdomain/servfail/delay)", "nxdomain"},
			{"DNS listen address", "127.0.0.1:5353"},
			{"Upstream DNS", "8.8.8.8:53"},
		}
	case 4: // Bandwidth
		return [][2]string{
			{"gRPC address", "127.0.0.1:50051"},
			{"Interface", "eth0"},
			{"Rate (kbps)", "1000"},
		}
	default: // HTTP Fault
		return [][2]string{
			{"gRPC address", "127.0.0.1:50051"},
			{"Listen address", ":8880"},
			{"Target URL", "http://localhost:9090"},
			{"Error rate (0.0-1.0)", "0.3"},
			{"Error code", "500"},
			{"Delay (ms)", "0"},
		}
	}
}

// ── Model ──────────────────────────────────────────────────────────────────

type tuiModel struct {
	state      tuiState
	cursor     int
	chaosIdx   int
	inputs     []textinput.Model
	focusIdx   int
	spinner    spinner.Model
	statusMsg  string
	errMsg     string
}

func newTUI() tuiModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	return tuiModel{spinner: sp}
}

func (m tuiModel) Init() tea.Cmd {
	return nil
}

// ── Update ─────────────────────────────────────────────────────────────────

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateMenu:
		return m.updateMenu(msg)
	case stateConfigure:
		return m.updateConfigure(msg)
	case stateRunning:
		return m.updateRunning(msg)
	case stateStopped:
		return m.updateStopped(msg)
	}
	return m, nil
}

func (m tuiModel) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(chaosTypes)-1 {
				m.cursor++
			}
		case "enter":
			m.chaosIdx = m.cursor
			m = m.buildInputs()
			m.state = stateConfigure
		}
	}
	return m, nil
}

func (m tuiModel) buildInputs() tuiModel {
	fields := configFields(m.chaosIdx)
	m.inputs = make([]textinput.Model, len(fields))
	for i, f := range fields {
		ti := textinput.New()
		ti.Placeholder = f[1]
		ti.CharLimit = 120
		ti.Width = 40
		// Use the label as prompt
		ti.Prompt = f[0] + ": "
		m.inputs[i] = ti
	}
	m.inputs[0].Focus()
	m.focusIdx = 0
	return m
}

func (m tuiModel) updateConfigure(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = stateMenu
			return m, nil
		case "tab", "down":
			m.inputs[m.focusIdx].Blur()
			m.focusIdx = (m.focusIdx + 1) % len(m.inputs)
			m.inputs[m.focusIdx].Focus()
		case "shift+tab", "up":
			m.inputs[m.focusIdx].Blur()
			m.focusIdx = (m.focusIdx - 1 + len(m.inputs)) % len(m.inputs)
			m.inputs[m.focusIdx].Focus()
		case "enter":
			if m.focusIdx == len(m.inputs)-1 {
				// Last field: launch experiment.
				m.statusMsg = fmt.Sprintf("Started %s", chaosTypes[m.chaosIdx])
				m.state = stateRunning
				return m, m.spinner.Tick
			}
			// Advance to next field.
			m.inputs[m.focusIdx].Blur()
			m.focusIdx++
			m.inputs[m.focusIdx].Focus()
		}
	}
	// Update focused input.
	var cmd tea.Cmd
	m.inputs[m.focusIdx], cmd = m.inputs[m.focusIdx].Update(msg)
	return m, cmd
}

func (m tuiModel) updateRunning(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "s":
			m.statusMsg = fmt.Sprintf("Stopped %s", chaosTypes[m.chaosIdx])
			m.state = stateStopped
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, m.spinner.Tick
}

func (m tuiModel) updateStopped(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			// Any key → back to menu.
			m.state = stateMenu
			m.statusMsg = ""
		}
	}
	return m, nil
}

// ── View ───────────────────────────────────────────────────────────────────

func (m tuiModel) View() string {
	switch m.state {
	case stateMenu:
		return m.viewMenu()
	case stateConfigure:
		return m.viewConfigure()
	case stateRunning:
		return m.viewRunning()
	case stateStopped:
		return m.viewStopped()
	}
	return ""
}

func (m tuiModel) viewMenu() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("obzevMini — chaos type") + "\n\n")
	for i, t := range chaosTypes {
		if i == m.cursor {
			b.WriteString(selectedItem.Render("▶ "+t) + "\n")
		} else {
			b.WriteString(dimItem.Render("  "+t) + "\n")
		}
	}
	b.WriteString("\n" + helpStyle.Render("↑/↓ navigate • enter select • q quit"))
	return b.String()
}

func (m tuiModel) viewConfigure() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(chaosTypes[m.chaosIdx]+" — configure") + "\n\n")
	for _, inp := range m.inputs {
		b.WriteString(inp.View() + "\n")
	}
	b.WriteString("\n" + helpStyle.Render("tab/↓ next • shift+tab/↑ prev • enter confirm • esc back"))
	return b.String()
}

func (m tuiModel) viewRunning() string {
	return fmt.Sprintf("\n  %s %s\n\n%s\n",
		m.spinner.View(),
		statusOK.Render(m.statusMsg),
		helpStyle.Render("s or q to stop"),
	)
}

func (m tuiModel) viewStopped() string {
	return fmt.Sprintf("\n  %s\n\n%s\n",
		statusOK.Render("✓ "+m.statusMsg),
		helpStyle.Render("any key to return to menu • q quit"),
	)
}

// ── Entry point ────────────────────────────────────────────────────────────

func runTUI() error {
	p := tea.NewProgram(newTUI(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
