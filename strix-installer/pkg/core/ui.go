package core

// UI defines the interface for user interaction
// Both TUI and GUI implement this interface
type UI interface {
	// Stage lifecycle
	StageStart(stage Stage)
	StageComplete(result StageResult)

	// Progress reporting
	Progress(percent int, message string)

	// Logging
	Log(level LogLevel, message string)

	// User prompts
	Confirm(message string, defaultYes bool) bool
	Select(message string, options []string) int
	Input(message string, defaultVal string) string
}

// NullUI is a no-op implementation for testing
type NullUI struct{}

func (n *NullUI) StageStart(stage Stage)                         {}
func (n *NullUI) StageComplete(result StageResult)               {}
func (n *NullUI) Progress(percent int, message string)           {}
func (n *NullUI) Log(level LogLevel, message string)             {}
func (n *NullUI) Confirm(message string, defaultYes bool) bool   { return defaultYes }
func (n *NullUI) Select(message string, options []string) int    { return 0 }
func (n *NullUI) Input(message string, defaultVal string) string { return defaultVal }
