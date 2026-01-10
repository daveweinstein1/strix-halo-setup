package bootloader

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/daveweinstein1/strixforge/pkg/system"
)

// SystemdBoot manages systemd-boot via sdboot-manage
type SystemdBoot struct {
	configPath string
}

// NewSystemdBoot creates a new SystemdBoot manager
func NewSystemdBoot() *SystemdBoot {
	return &SystemdBoot{
		configPath: "/etc/sdboot-manage.conf",
	}
}

func (s *SystemdBoot) Name() string { return "systemd-boot" }

// IsInstalled checks if systemd-boot is active
func (s *SystemdBoot) IsInstalled() bool {
	// 1. Check if sdboot-manage exists
	if !system.CheckCommand("sdboot-manage") {
		return false
	}
	// 2. Check if config exists
	if _, err := os.Stat(s.configPath); os.IsNotExist(err) {
		return false
	}
	// 3. Confirm we are actually booted via EFI (usually systemd-boot manages this)
	if _, err := os.Stat("/sys/firmware/efi"); os.IsNotExist(err) {
		return false
	}
	// 4. Look for loader config
	if _, err := os.Stat("/boot/loader/loader.conf"); os.IsNotExist(err) {
		return false
	}
	return true
}

// Backup creates a timestamped backup of config
func (s *SystemdBoot) Backup(ctx context.Context) (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup-%s", s.configPath, timestamp)

	result, err := system.ExecSudo(ctx, "cp", s.configPath, backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to backup systemd-boot config: %s\n%s", err, result.Stderr)
	}
	return backupPath, nil
}

// AddParam adds a parameter to LINUX_OPTIONS
func (s *SystemdBoot) AddParam(ctx context.Context, param string) error {
	// Read file
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return fmt.Errorf("failed to read sdboot-manage config: %v", err)
	}
	content := string(data)

	// Check if param already exists
	// Note: simplistic check, but safe enough for "iommu=pt" etc.
	if strings.Contains(content, param) {
		return nil
	}

	// We look for LINUX_OPTIONS="..."
	// If it doesn't exist, we append it. If it does, we assume sed will handle it.
	// Actually for simplicity let's use sed to append inside the quotes
	// match: LINUX_OPTIONS="...
	// replace: LINUX_OPTIONS="... param "

	// Since we don't have robust config parsing, we'll try to use sed carefully
	// Pattern: Find LINUX_OPTIONS="... and make sure we don't break the closing quote
	// Actually, easier strategy:
	// 1. If LINUX_OPTIONS detected, use sed substitution to append inside last quote
	// 2. If not detected, append new line

	if strings.Contains(content, "LINUX_OPTIONS=") {
		// Append to existing options
		// sed -i 's/^LINUX_OPTIONS="/LINUX_OPTIONS="param /'
		// Note: sdboot-manage usually has LINUX_OPTIONS="" empty by default or commented?
		// Default in cachyos config checks... it's usually LINUX_OPTIONS=""

		// Let's use a safer approach: append the param to the beginning of the string inside quotes
		// This handles the empty case LINUX_OPTIONS="" -> LINUX_OPTIONS="param "
		sedCmd := fmt.Sprintf(`sed -i 's/^LINUX_OPTIONS="/LINUX_OPTIONS="%s /'`, param)

		// Wait, if line is commented out? We should uncomment it.
		// sed -i 's/^#LINUX_OPTIONS=/LINUX_OPTIONS=/'
		uncommentCmd := `sed -i 's/^#LINUX_OPTIONS=/LINUX_OPTIONS=/' ` + s.configPath
		system.ExecShellSudo(ctx, uncommentCmd)

		cmd := fmt.Sprintf("%s %s", sedCmd, s.configPath)
		if _, err := system.ExecShellSudo(ctx, cmd); err != nil {
			return fmt.Errorf("failed to add param to systemd-boot: %v", err)
		}
	} else {
		// perform append
		newLine := fmt.Sprintf(`echo 'LINUX_OPTIONS="%s"' | sudo tee -a %s`, param, s.configPath)
		if _, err := system.ExecShellSudo(ctx, newLine); err != nil {
			return fmt.Errorf("failed to append config: %v", err)
		}
	}

	// Regenerate entries
	return s.update(ctx)
}

func (s *SystemdBoot) update(ctx context.Context) error {
	result, err := system.ExecSudo(ctx, "sdboot-manage", "gen")
	if err != nil {
		return fmt.Errorf("failed to update systemd-boot: %s\n%s", err, result.Stderr)
	}
	return nil
}
