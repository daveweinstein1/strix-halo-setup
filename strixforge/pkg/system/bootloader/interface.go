package bootloader

import "context"

// Bootloader interface defines operations for boot configuration
type Bootloader interface {
	// Name returns the name of the bootloader (e.g., "GRUB", "systemd-boot")
	Name() string

	// IsInstalled checks if this bootloader is present and active
	IsInstalled() bool

	// AddParam adds a kernel parameter if it's not already present
	AddParam(ctx context.Context, param string) error

	// Backup creates a backup of the bootloader configuration
	// Returns the path to the backup file
	Backup(ctx context.Context) (string, error)
}
