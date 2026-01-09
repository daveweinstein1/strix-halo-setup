package system

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// Grub provides bootloader management
type Grub struct {
	configPath string
}

// NewGrub creates a new Grub instance
func NewGrub() *Grub {
	return &Grub{
		configPath: "/etc/default/grub",
	}
}

// Backup creates a timestamped backup of grub config
func (g *Grub) Backup(ctx context.Context) (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup-%s", g.configPath, timestamp)

	result, err := ExecSudo(ctx, "cp", g.configPath, backupPath)
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

// AddCmdlineParam adds a parameter to kernel command line if not present
func (g *Grub) AddCmdlineParam(ctx context.Context, param string) error {
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
	sedCmd := fmt.Sprintf(`sed -i 's/GRUB_CMDLINE_LINUX_DEFAULT="[^"]*"/GRUB_CMDLINE_LINUX_DEFAULT="%s"/' %s`, params, g.configPath)
	result, err := ExecShellSudo(ctx, sedCmd)
	if err != nil {
		return fmt.Errorf("failed to update grub config: %s\n%s", err, result.Stderr)
	}
	return nil
}

// UpdateGrub regenerates grub configuration
func (g *Grub) UpdateGrub(ctx context.Context) error {
	// Try grub-mkconfig first (Arch/CachyOS)
	result, err := ExecSudo(ctx, "grub-mkconfig", "-o", "/boot/grub/grub.cfg")
	if err != nil {
		return fmt.Errorf("failed to update grub: %s\n%s", err, result.Stderr)
	}
	return nil
}

// GetCurrentCmdline returns the currently running kernel's command line
func (g *Grub) GetCurrentCmdline(ctx context.Context) (string, error) {
	data, err := os.ReadFile("/proc/cmdline")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// HasParam checks if a parameter is in the current kernel command line
func (g *Grub) HasParam(ctx context.Context, param string) (bool, error) {
	cmdline, err := g.GetCurrentCmdline(ctx)
	if err != nil {
		return false, err
	}
	return strings.Contains(cmdline, param), nil
}
