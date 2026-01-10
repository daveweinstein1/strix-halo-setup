package bootloader

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/daveweinstein1/strixforge/pkg/system"
)

// Refind manages rEFInd bootloader configuration
type Refind struct {
	configPath string
}

// NewRefind creates a new Refind manager
func NewRefind() *Refind {
	return &Refind{
		configPath: "/boot/refind_linux.conf", // Primary linux options file
	}
}

func (r *Refind) Name() string { return "rEFInd" }

// IsInstalled checks if rEFInd is installed/configured for linux
func (r *Refind) IsInstalled() bool {
	// Check for config file
	if _, err := os.Stat(r.configPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// Backup creates a timestamped backup
func (r *Refind) Backup(ctx context.Context) (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup-%s", r.configPath, timestamp)

	result, err := system.ExecSudo(ctx, "cp", r.configPath, backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to backup refind config: %s\n%s", err, result.Stderr)
	}
	return backupPath, nil
}

// AddParam adds a parameter to the default boot options
func (r *Refind) AddParam(ctx context.Context, param string) error {
	data, err := os.ReadFile(r.configPath)
	if err != nil {
		return fmt.Errorf("failed to read refind config: %v", err)
	}
	content := string(data)

	if strings.Contains(content, param) {
		return nil
	}

	// Strategy: refind_linux.conf format:
	// "Description" "loader options..."
	// We want to target the first entry or entries containing "default" or "Default"?
	// Or just append to all lines? Usually safe to append to all lines in this file.

	// Better yet: "Boot using default options" is the standard generated line.
	// sed -i 's/"$/ param"/' to append BEFORE the closing quote of options?
	// The file format is:
	// "Title"  "options"

	// So we need to match the SECOND quoted string of each line.
	// sed 's/"$/ param"/' applies to end of line, which is usually the closing quote of the options.

	sedCmd := fmt.Sprintf(`sed -i 's/"$/ %s"/' %s`, param, r.configPath)

	result, err := system.ExecShellSudo(ctx, sedCmd)
	if err != nil {
		return fmt.Errorf("failed to add param to refind: %s\n%s", err, result.Stderr)
	}

	return nil // No update command needed for rEFInd, it reads config at boot
}
