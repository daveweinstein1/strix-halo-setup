package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/daveweinstein1/strix-installer/pkg/core"
	"github.com/daveweinstein1/strix-installer/pkg/platform/strixhalo"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF00")).
			MarginBottom(1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00AAFF"))

	warnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAA00"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))
)

// Model holds the TUI state
type Model struct {
	platform     core.Platform
	device       core.Device
	stages       []core.Stage
	currentStage int
	progress     progress.Model
	spinner      spinner.Model
	logs         []string
	running      bool
	done         bool
	err          error
	width        int
	height       int
}

// Messages
type stageStartMsg struct{ stage core.Stage }
type stageCompleteMsg struct{ result core.StageResult }
type progressMsg struct {
	percent int
	message string
}
type logMsg struct {
	level   core.LogLevel
	message string
}
type doneMsg struct{ err error }

func main() {
	fmt.Println(titleStyle.Render("Strix Halo Post-Installer"))
	fmt.Println()

	// Initialize platform
	platform := strixhalo.New()

	// Detect device
	fmt.Println("Detecting hardware...")
	device, err := platform.Detect()
	if err != nil {
		fmt.Printf("Warning: Could not detect device: %v\n", err)
	} else {
		fmt.Printf("Detected: %s\n", device.Name())

		// Show quirks
		quirks := device.Quirks()
		if len(quirks) > 0 {
			fmt.Println("\nDevice-specific notes:")
			for _, q := range quirks {
				fmt.Printf("  • %s\n", q.Description)
			}
		}
	}
	fmt.Println()

	// Initialize TUI
	p := progress.New(progress.WithDefaultGradient())
	s := spinner.New()
	s.Spinner = spinner.Dot

	m := Model{
		platform:     platform,
		device:       device,
		stages:       platform.Stages(),
		currentStage: -1,
		progress:     p,
		spinner:      s,
		logs:         make([]string, 0),
		width:        80,
		height:       24,
	}

	// Run TUI
	program := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// Init initializes the TUI
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.EnterAltScreen,
	)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if !m.running && !m.done {
				m.running = true
				return m, m.runInstall()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 10

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case stageStartMsg:
		m.currentStage++
		m.logs = append(m.logs, infoStyle.Render(fmt.Sprintf("→ Starting: %s", msg.stage.Name())))
		return m, m.spinner.Tick

	case stageCompleteMsg:
		if msg.result.Status == core.StatusSuccess {
			m.logs = append(m.logs, successStyle.Render(fmt.Sprintf("✓ Complete: %s", msg.result.StageName)))
		} else if msg.result.Status == core.StatusFailed {
			m.logs = append(m.logs, errorStyle.Render(fmt.Sprintf("✗ Failed: %s - %v", msg.result.StageName, msg.result.Error)))
		}
		return m, nil

	case progressMsg:
		m.logs = append(m.logs, fmt.Sprintf("  %s", msg.message))
		return m, m.progress.SetPercent(float64(msg.percent) / 100)

	case logMsg:
		var styled string
		switch msg.level {
		case core.LogError:
			styled = errorStyle.Render(msg.message)
		case core.LogWarn:
			styled = warnStyle.Render(msg.message)
		default:
			styled = msg.message
		}
		m.logs = append(m.logs, styled)
		return m, nil

	case doneMsg:
		m.running = false
		m.done = true
		m.err = msg.err
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("Strix Halo Post-Installer"))
	b.WriteString("\n\n")

	// Device info
	if m.device != nil {
		b.WriteString(infoStyle.Render(fmt.Sprintf("Device: %s", m.device.Name())))
		b.WriteString("\n\n")
	}

	// Stage list
	b.WriteString("Stages:\n")
	for i, stage := range m.stages {
		prefix := "  "
		if i == m.currentStage && m.running {
			prefix = m.spinner.View() + " "
		} else if i < m.currentStage {
			prefix = successStyle.Render("✓ ")
		}

		name := stage.Name()
		if stage.Optional() {
			name += " (optional)"
		}
		b.WriteString(fmt.Sprintf("%s%s\n", prefix, name))
	}
	b.WriteString("\n")

	// Progress bar
	if m.running {
		b.WriteString(m.progress.View())
		b.WriteString("\n\n")
	}

	// Logs (last 10)
	if len(m.logs) > 0 {
		b.WriteString("Log:\n")
		start := 0
		if len(m.logs) > 10 {
			start = len(m.logs) - 10
		}
		for _, log := range m.logs[start:] {
			b.WriteString(fmt.Sprintf("  %s\n", log))
		}
		b.WriteString("\n")
	}

	// Status / Instructions
	if m.done {
		if m.err != nil {
			b.WriteString(errorStyle.Render(fmt.Sprintf("Installation failed: %v\n", m.err)))
		} else {
			b.WriteString(successStyle.Render("Installation complete!\n"))
		}
		b.WriteString("\nPress 'q' to exit")
	} else if !m.running {
		b.WriteString("Press ENTER to start installation, or 'q' to quit")
	}

	return b.String()
}

// runInstall runs the installation in the background
func (m Model) runInstall() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		ui := &tuiAdapter{program: nil} // Would need program reference for real impl

		engine := core.NewEngine(m.platform, ui)
		err := engine.Run(ctx)

		return doneMsg{err: err}
	}
}

// tuiAdapter implements core.UI for the TUI
type tuiAdapter struct {
	program *tea.Program
}

func (t *tuiAdapter) StageStart(stage core.Stage) {
	// Would send message to program
}

func (t *tuiAdapter) StageComplete(result core.StageResult) {
	// Would send message to program
}

func (t *tuiAdapter) Progress(percent int, message string) {
	// Would send message to program
}

func (t *tuiAdapter) Log(level core.LogLevel, message string) {
	// Would send message to program
}

func (t *tuiAdapter) Confirm(message string, defaultYes bool) bool {
	return defaultYes // Simplified for now
}

func (t *tuiAdapter) Select(message string, options []string) int {
	return 0
}

func (t *tuiAdapter) Input(message string, defaultVal string) string {
	return defaultVal
}
