package stages

import (
	"context"
	"fmt"

	"github.com/daveweinstein1/strix-installer/pkg/core"
	"github.com/daveweinstein1/strix-installer/pkg/system"
)

// SystemStage performs system update and installs essentials
type SystemStage struct{}

func NewSystemStage() *SystemStage { return &SystemStage{} }

func (s *SystemStage) ID() string   { return "system" }
func (s *SystemStage) Name() string { return "System Update" }
func (s *SystemStage) Description() string {
	return "Update mirrors, system packages, install essentials"
}
func (s *SystemStage) Optional() bool { return false }

func (s *SystemStage) Run(ctx context.Context, ui core.UI) error {
	pacman := system.NewPacman()

	// Step 1: Rate mirrors (optional, CachyOS specific)
	ui.Progress(10, "Checking for mirror optimization...")
	if system.CheckCommand("cachyos-rate-mirrors") {
		ui.Log(core.LogInfo, "Running CachyOS mirror ranking...")
		result, err := system.ExecSudo(ctx, "cachyos-rate-mirrors")
		if err != nil {
			ui.Log(core.LogWarn, fmt.Sprintf("Mirror ranking failed: %v", err))
		} else {
			ui.Log(core.LogInfo, "✓ Mirrors optimized")
			_ = result // suppress unused warning
		}
	} else if system.CheckCommand("rate-mirrors") {
		ui.Log(core.LogInfo, "Running rate-mirrors...")
		_, err := system.ExecShellSudo(ctx, "rate-mirrors --save /etc/pacman.d/mirrorlist arch")
		if err != nil {
			ui.Log(core.LogWarn, fmt.Sprintf("Mirror ranking failed: %v", err))
		}
	} else {
		ui.Log(core.LogInfo, "Mirror ranking tools not found, skipping")
	}

	// Step 2: Full system update
	ui.Progress(30, "Updating system packages...")
	if err := pacman.Update(ctx); err != nil {
		return fmt.Errorf("system update failed: %v", err)
	}
	ui.Log(core.LogInfo, "✓ System updated")

	// Step 3: Install essential packages
	ui.Progress(60, "Installing essential packages...")
	essentials := []string{
		"base-devel",
		"git",
		"wget",
		"curl",
		"vim",
		"neovim",
		"btop",
		"neofetch",
		"fastfetch",
	}

	if err := pacman.Install(ctx, essentials...); err != nil {
		return fmt.Errorf("failed to install essentials: %v", err)
	}
	ui.Log(core.LogInfo, "✓ Essential packages installed")

	ui.Progress(100, "System update complete")
	return nil
}

func (s *SystemStage) Rollback(ctx context.Context) error {
	return nil
}
