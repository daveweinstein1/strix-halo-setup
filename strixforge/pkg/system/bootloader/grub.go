package bootloader

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/daveweinstein1/strixforge/pkg/system"
)

// Grub provides bootloader management for standard GRUB
type Grub struct {
	configPath string
}

// NewGrub creates a new Grub manager
func NewGrub() *Grub {
	return &Grub{
		configPath: "/etc/default/grub",
	}
}

func (g *Grub) Name() string { return "GRUB" }

// IsInstalled checks if GRUB seems to be the active/installed bootloader
func (g *Grub) IsInstalled() bool {
	// Check for config dir
	if _, err := os.Stat("/boot/grub"); os.IsNotExist(err) {
		return false
	}
	// Check for defaults file
	if _, err := os.Stat(g.configPath); os.IsNotExist(err) {
		return false
	}
	// Logic check: if /boot/loader/entries exists, it might be systemd-boot
	// But GRUB can coexist. We return true if GRUB is configured.
	return true
}

// Backup creates a timestamped backup of grub config
func (g *Grub) Backup(ctx context.Context) (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup-%s", g.configPath, timestamp)

	result, err := system.ExecSudo(ctx, "cp", g.configPath, backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to backup grub: %s\n%s", err, result.Stderr)
	}
	return backupPath, nil
}

// GetCmdlineParams returns current kernel command line parameters
func (g *Grub) GetCmdlineParams(ctx context.Context) (string, error) {
	data, err := os.ReadFile(g.configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read grub config: %v", err)
	}

	re := regexp.MustCompile(`GRUB_CMDLINE_LINUX_DEFAULT="([^"]*)"`)
	matches := re.FindSubmatch(data)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not find GRUB_CMDLINE_LINUX_DEFAULT")
	}

	return string(matches[1]), nil
}

// AddParam adds a parameter to kernel command line if not present
func (g *Grub) AddParam(ctx context.Context, param string) error {
	current, err := g.GetCmdlineParams(ctx)
	if err != nil {
		return err
	}

	// Check if already present
	if strings.Contains(current, param) {
		return nil // Already there
	}

	// Add parameter
	newParams := strings.TrimSpace(current + " " + param)
	return g.SetCmdlineParams(ctx, newParams)
}

// SetCmdlineParams sets the kernel command line parameters
func (g *Grub) SetCmdlineParams(ctx context.Context, params string) error {
	// Use sed to update the config
	// Note: We use a slightly more robust sed command here, but basically same logic
	sedCmd := fmt.Sprintf(`sed -i 's/GRUB_CMDLINE_LINUX_DEFAULT="[^"]*"/GRUB_CMDLINE_LINUX_DEFAULT="%s"/' %s`, params, g.configPath)
	result, err := system.ExecShellSudo(ctx, sedCmd)
	if err != nil {
		return fmt.Errorf("failed to update grub config: %s\n%s", err, result.Stderr)
	}

	// Update GRUB
	return g.update(ctx)
}

// update regenerates grub configuration
func (g *Grub) update(ctx context.Context) error {
	// Try grub-mkconfig first (Arch/CachyOS)
	result, err := system.ExecSudo(ctx, "grub-mkconfig", "-o", "/boot/grub/grub.cfg")
	if err != nil {
		return fmt.Errorf("failed to update grub: %s\n%s", err, result.Stderr)
	}
	return nil
}
