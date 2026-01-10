package bootloader

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/daveweinstein1/strixforge/pkg/system"
)

// Limine manages Limine bootloader configuration
type Limine struct {
	configPath string
}

// NewLimine creates a new Limine manager
func NewLimine() *Limine {
	return &Limine{
		configPath: "/etc/default/limine", // Primary config for auto-generation
	}
}

func (l *Limine) Name() string { return "Limine" }

// IsInstalled checks if Limine is active
func (l *Limine) IsInstalled() bool {
	// 1. Check if limine command/tools exist
	if !system.CheckCommand("limine-mkinitcpio") {
		return false
	}
	// 2. Check for defaults file
	if _, err := os.Stat(l.configPath); os.IsNotExist(err) {
		return false
	}
	// 3. Check for boot dir
	if _, err := os.Stat("/boot/limine"); os.IsNotExist(err) {
		return false
	}
	return true
}

// Backup creates a timestamped backup
func (l *Limine) Backup(ctx context.Context) (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup-%s", l.configPath, timestamp)

	result, err := system.ExecSudo(ctx, "cp", l.configPath, backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to backup limine config: %s\n%s", err, result.Stderr)
	}
	return backupPath, nil
}

// AddParam adds a parameter to KERNEL_CMDLINE
func (l *Limine) AddParam(ctx context.Context, param string) error {
	data, err := os.ReadFile(l.configPath)
	if err != nil {
		return fmt.Errorf("failed to read limine config: %v", err)
	}
	content := string(data)

	if strings.Contains(content, param) {
		return nil
	}

	// Strategy: Append to KERNEL_CMDLINE[default]="...
	// or KERNEL_CMDLINE="... if that's the syntax used (bash array style usually)

	// Default CachyOS config uses: KERNEL_CMDLINE[default]="..."
	// Or sometimes just KERNEL_CMDLINE="..."

	// We'll look for KERNEL_CMDLINE...="
	// Safer: just append a new line which bash will use to override/append?
	// No, that might overwrite.

	// Sed replace strategy:
	// Insert param inside the first KERNEL_CMDLINE.*=" quote
	// sed -i 's/KERNEL_CMDLINE.*="/&param /'

	sedCmd := fmt.Sprintf(`sed -i 's/KERNEL_CMDLINE.*="/&%s /'`, param)

	// NOTE: This assumes the lines exists. If commented out, we must uncomment.
	if strings.Contains(content, "#KERNEL_CMDLINE") || strings.Contains(content, "# KERNEL_CMDLINE") {
		// Try to uncomment
		uncommentCmd := `sed -i 's/^#\s*KERNEL_CMDLINE/KERNEL_CMDLINE/' ` + l.configPath
		system.ExecShellSudo(ctx, uncommentCmd)
	}

	cmd := fmt.Sprintf("%s %s", sedCmd, l.configPath)
	if _, err := system.ExecShellSudo(ctx, cmd); err != nil {
		return fmt.Errorf("failed to add param to limine: %v", err)
	}

	return l.update(ctx)
}

func (l *Limine) update(ctx context.Context) error {
	// Run limine-mkinitcpio which triggers the hook to update limine.conf
	result, err := system.ExecSudo(ctx, "limine-mkinitcpio")
	if err != nil {
		return fmt.Errorf("failed to update limine: %s\n%s", err, result.Stderr)
	}
	// Also ensure bootloader is installed/updated on disk if needed?
	// Usually mkinitcpio is enough for config updates.
	return nil
}
