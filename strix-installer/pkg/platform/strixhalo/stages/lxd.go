package stages

import (
	"context"
	"fmt"
	"os/user"

	"github.com/daveweinstein1/strix-installer/pkg/core"
	"github.com/daveweinstein1/strix-installer/pkg/system"
)

// LXDStage installs and configures LXD with GPU passthrough
type LXDStage struct{}

func NewLXDStage() *LXDStage { return &LXDStage{} }

func (s *LXDStage) ID() string   { return "lxd" }
func (s *LXDStage) Name() string { return "LXD Containerization" }
func (s *LXDStage) Description() string {
	return "Install LXD, configure GPU passthrough, enable nesting"
}
func (s *LXDStage) Optional() bool { return false }

func (s *LXDStage) Run(ctx context.Context, ui core.UI) error {
	pacman := system.NewPacman()
	systemd := system.NewSystemd()
	lxd := system.NewLXD()

	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("could not determine current user: %v", err)
	}
	username := currentUser.Username

	// Step 1: Install LXD
	ui.Progress(10, "Installing LXD...")
	if err := pacman.Install(ctx, "lxd"); err != nil {
		return fmt.Errorf("failed to install LXD: %v", err)
	}
	ui.Log(core.LogInfo, "✓ LXD installed")

	// Step 2: Enable and start LXD socket
	ui.Progress(25, "Enabling LXD service...")
	if err := systemd.EnableAndStart(ctx, "lxd.socket"); err != nil {
		return fmt.Errorf("failed to enable LXD: %v", err)
	}
	ui.Log(core.LogInfo, "✓ LXD service enabled")

	// Step 3: Add user to lxd group
	ui.Progress(40, "Configuring user permissions...")
	if !lxd.IsUserInGroup(ctx, username) {
		if err := lxd.AddUserToGroup(ctx, username); err != nil {
			return fmt.Errorf("failed to add user to lxd group: %v", err)
		}
		ui.Log(core.LogInfo, fmt.Sprintf("✓ Added %s to lxd group", username))
		ui.Log(core.LogWarn, "NOTE: Log out and back in for group changes to take effect")
	} else {
		ui.Log(core.LogInfo, fmt.Sprintf("✓ User %s already in lxd group", username))
	}

	// Step 4: Initialize LXD
	ui.Progress(55, "Initializing LXD...")
	if err := lxd.Init(ctx); err != nil {
		ui.Log(core.LogWarn, fmt.Sprintf("LXD init warning: %v", err))
		// Continue - may already be initialized
	}
	ui.Log(core.LogInfo, "✓ LXD initialized")

	// Step 5: Add GPU device to default profile
	ui.Progress(70, "Configuring GPU passthrough...")
	if err := lxd.AddGPUDevice(ctx); err != nil {
		ui.Log(core.LogWarn, fmt.Sprintf("GPU device config warning: %v", err))
	} else {
		ui.Log(core.LogInfo, "✓ GPU passthrough configured")
	}

	// Step 6: Enable nesting for Docker-in-LXD
	ui.Progress(85, "Enabling container nesting...")
	if err := lxd.EnableNesting(ctx); err != nil {
		ui.Log(core.LogWarn, fmt.Sprintf("Nesting config warning: %v", err))
	} else {
		ui.Log(core.LogInfo, "✓ Container nesting enabled")
	}

	ui.Progress(100, "LXD setup complete")
	return nil
}

func (s *LXDStage) Rollback(ctx context.Context) error {
	return nil
}
