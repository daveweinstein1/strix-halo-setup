package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/daveweinstein1/strixforge/pkg/containerhub"
	"github.com/daveweinstein1/strixforge/pkg/core"
	"github.com/daveweinstein1/strixforge/pkg/platform/strixhalo"
	"github.com/daveweinstein1/strixforge/pkg/system"
)

//go:embed frontend/*
var frontendFS embed.FS

//go:embed all:frontend/dist
var guiAssets embed.FS

var (
	forceTUI        = flag.Bool("tui", false, "Force TUI mode")
	forceGUI        = flag.Bool("gui", false, "Force GUI mode (native window)")
	autoMode        = flag.Bool("auto", false, "Run all stages without prompts")
	manualMode      = flag.Bool("manual", false, "Manually select stages to run")
	marketplaceMode = flag.Bool("hub", false, "Browse Container Hub")
	checkVersions   = flag.Bool("check-versions", false, "Check package versions and exit")
	dryRun          = flag.Bool("dry-run", false, "Simulate installation without changes")
)

func main() {
	flag.Parse()

	// Version check mode
	if *checkVersions {
		runVersionCheck()
		return
	}

	// Auto mode: run without TUI/GUI
	if *autoMode {
		runAutoMode()
		return
	}

	// Marketplace mode
	if *marketplaceMode {
		runMarketplace()
		return
	}

	// Mode selection
	if *forceTUI || *manualMode {
		runTUI()
		return
	}

	if *forceGUI {
		runGUI()
		return
	}

	// Auto-detect: GUI if display available, otherwise TUI
	if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" {
		runGUI()
	} else {
		runTUI()
	}
}

// runVersionCheck displays version comparison and exits
func runVersionCheck() {
	fmt.Println(titleStyle.Render("Strix Halo Version Check"))
	fmt.Println()

	ctx := context.Background()
	checks, err := system.CheckAllVersions(ctx)
	if err != nil {
		fmt.Printf("Error checking versions: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(system.FormatVersionTable(checks))
	fmt.Println()
	fmt.Println(system.SummarizeVersionChecks(checks))

	if system.HasCriticalMismatches(checks) {
		fmt.Println(warnStyle.Render("⚠ Some packages have older versions than expected."))
		fmt.Println("Run the installer to update, or use --auto to proceed anyway.")
		os.Exit(1)
	}

	fmt.Println(successStyle.Render("✓ All versions look good!"))
}

// runAutoMode runs all stages without prompts
func runAutoMode() {
	fmt.Println(titleStyle.Render("Strix Halo Auto-Install"))
	fmt.Println()

	// Detect hardware
	fmt.Println("Detecting hardware...")
	device, err := platform.Detect()
	if err != nil {
		fmt.Printf("Warning: Could not detect device: %v\n", err)
	} else {
		fmt.Printf("Detected: %s\n", device.Name())
	}
	fmt.Println()

	// Create UI adapter that auto-accepts
	ui := &autoUIAdapter{dryRun: *dryRun}

	// Run engine
	ctx := context.Background()
	engine := core.NewEngine(platform, ui)
	if *dryRun {
		engine.SetDryRun(true)
		fmt.Println(warnStyle.Render("DRY RUN MODE - No changes will be made"))
		fmt.Println()
	}

	err = engine.Run(ctx)
	if err != nil {
		fmt.Printf(errorStyle.Render("Installation failed: %v\n"), err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(successStyle.Render("✓ Installation complete!"))
}

// runMarketplace launches the TUI marketplace
func runMarketplace() {
	// Initialize manager
	mgr := containerhub.NewManager()

	// Load config (or use defaults if file missing)
	if err := mgr.LoadConfigFromPath("configs/registries.yaml"); err != nil {
		fmt.Printf("Warning: Could not load registries.yaml: %v. Using defaults.\n", err)
		// TODO: Add default/hardcoded fallback if file missing
	}

	// Launch Bubble Tea program
	p := tea.NewProgram(NewMarketplaceModel(mgr, 80, 24, func() {
		// On Back, we exit for now
		os.Exit(0)
	}))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running marketplace: %v\n", err)
		os.Exit(1)
	}
}

// autoUIAdapter implements UI interface for auto mode
type autoUIAdapter struct {
	dryRun bool
}

func (a *autoUIAdapter) StageStart(stage core.Stage) {
	fmt.Printf("→ Starting: %s\n", stage.Name())
}

func (a *autoUIAdapter) StageComplete(result core.StageResult) {
	if result.Status == core.StatusSuccess {
		fmt.Printf(successStyle.Render("✓ Complete: %s\n"), result.StageName)
	} else if result.Status == core.StatusFailed {
		fmt.Printf(errorStyle.Render("✗ Failed: %s - %v\n"), result.StageName, result.Error)
	} else if result.Status == core.StatusSkipped {
		fmt.Printf("○ Skipped: %s\n", result.StageName)
	}
}

func (a *autoUIAdapter) Progress(percent int, message string) {
	fmt.Printf("  [%d%%] %s\n", percent, message)
}

func (a *autoUIAdapter) Log(level core.LogLevel, message string) {
	switch level {
	case core.LogError:
		fmt.Println(errorStyle.Render(message))
	case core.LogWarn:
		fmt.Println(warnStyle.Render(message))
	default:
		fmt.Println(message)
	}
}

func (a *autoUIAdapter) Confirm(message string, defaultYes bool) bool {
	// Auto mode: always use default
	return defaultYes
}

func (a *autoUIAdapter) Select(message string, options []string) int {
	return 0
}

func (a *autoUIAdapter) Input(message string, defaultVal string) string {
	return defaultVal
}

// runGUI starts the native Wails GUI
func runGUI() {
	// Create an instance of the app structure
	app := NewGUIApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Strixforge",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: guiAssets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		fmt.Println("Could not start GUI, falling back to TUI")
		runTUI()
	}
}

// GUIApp struct for Wails
type GUIApp struct {
	ctx            context.Context
	marketplaceMgr *containerhub.Manager
}

// NewGUIApp creates a new GUIApp
func NewGUIApp() *GUIApp {
	mgr := containerhub.NewManager()
	_ = mgr.LoadConfigFromPath("configs/registries.yaml")
	return &GUIApp{marketplaceMgr: mgr}
}

func (a *GUIApp) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetSystemStatus returns detected hardware info
func (a *GUIApp) GetSystemStatus() map[string]string {
	status := make(map[string]string)
	device, err := strixhalo.Detect(a.ctx)
	if err != nil {
		status["error"] = err.Error()
		return status
	}
	status["name"] = device.Name()
	status["manufacturer"] = device.Manufacturer()
	status["model"] = device.Model()
	status["quirks_count"] = fmt.Sprintf("%d", len(device.Quirks()))
	return status
}

// FetchHubImages returns all available images from Container Hub
func (a *GUIApp) FetchHubImages() []containerhub.Image {
	images, err := a.marketplaceMgr.FetchAllImages(a.ctx)
	if err != nil {
		return []containerhub.Image{}
	}
	return images
}

// API handlers kept for compatibility (unused now)
var platform = strixhalo.New()

// =============================================================================
// TUI Mode (Bubble Tea)
// =============================================================================

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

func runTUI() {
	fmt.Println(titleStyle.Render("Strix Halo Post-Installer"))
	fmt.Println()

	fmt.Println("Detecting hardware...")
	device, err := platform.Detect()
	if err != nil {
		fmt.Printf("Warning: Could not detect device: %v\n", err)
	} else {
		fmt.Printf("Detected: %s\n", device.Name())

		quirks := device.Quirks()
		if len(quirks) > 0 {
			fmt.Println("\nDevice-specific notes:")
			for _, q := range quirks {
				fmt.Printf("  • %s\n", q.Description)
			}
		}
	}
	fmt.Println()

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

	program := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.EnterAltScreen,
	)
}

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

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Strix Halo Post-Installer"))
	b.WriteString("\n\n")

	if m.device != nil {
		b.WriteString(infoStyle.Render(fmt.Sprintf("Device: %s", m.device.Name())))
		b.WriteString("\n\n")
	}

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

	if m.running {
		b.WriteString(m.progress.View())
		b.WriteString("\n\n")
	}

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

func (m Model) runInstall() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		ui := &tuiAdapter{}

		engine := core.NewEngine(m.platform, ui)
		err := engine.Run(ctx)

		return doneMsg{err: err}
	}
}

type tuiAdapter struct{}

func (t *tuiAdapter) StageStart(stage core.Stage)                    {}
func (t *tuiAdapter) StageComplete(result core.StageResult)          {}
func (t *tuiAdapter) Progress(percent int, message string)           {}
func (t *tuiAdapter) Log(level core.LogLevel, message string)        {}
func (t *tuiAdapter) Confirm(message string, defaultYes bool) bool   { return defaultYes }
func (t *tuiAdapter) Select(message string, options []string) int    { return 0 }
func (t *tuiAdapter) Input(message string, defaultVal string) string { return defaultVal }
