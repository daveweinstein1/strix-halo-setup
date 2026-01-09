package stages

import (
	"context"
	"fmt"

	"github.com/daveweinstein1/strix-installer/pkg/core"
	"github.com/daveweinstein1/strix-installer/pkg/system"
)

// WorkspaceStage provisions development containers
type WorkspaceStage struct{}

func NewWorkspaceStage() *WorkspaceStage { return &WorkspaceStage{} }

func (s *WorkspaceStage) ID() string   { return "workspace" }
func (s *WorkspaceStage) Name() string { return "AI & Dev Workspaces" }
func (s *WorkspaceStage) Description() string {
	return "Create ai-lab (ROCm/PyTorch) and dev-lab (Rust/Go) containers"
}
func (s *WorkspaceStage) Optional() bool { return true }

func (s *WorkspaceStage) Run(ctx context.Context, ui core.UI) error {
	lxd := system.NewLXD()

	// Create ai-lab container
	ui.Progress(10, "Creating ai-lab container...")
	if err := s.createAILab(ctx, ui, lxd); err != nil {
		return err
	}

	// Create dev-lab container
	ui.Progress(55, "Creating dev-lab container...")
	if err := s.createDevLab(ctx, ui, lxd); err != nil {
		return err
	}

	ui.Progress(100, "Workspaces ready")
	return nil
}

func (s *WorkspaceStage) createAILab(ctx context.Context, ui core.UI, lxd *system.LXD) error {
	name := "ai-lab"

	// Check if already exists
	if lxd.ContainerExists(ctx, name) {
		ui.Log(core.LogInfo, fmt.Sprintf("Container %s already exists, skipping creation", name))
		return nil
	}

	// Create container
	ui.Log(core.LogInfo, "Launching ai-lab container from archlinux image...")
	if err := lxd.CreateContainer(ctx, name, "images:archlinux/current"); err != nil {
		return fmt.Errorf("failed to create %s: %v", name, err)
	}

	// Wait for network
	ui.Progress(20, "Waiting for container network...")
	if err := lxd.WaitForNetwork(ctx, name); err != nil {
		ui.Log(core.LogWarn, "Network wait timed out, continuing anyway")
	}

	// Install ROCm and AI packages
	ui.Progress(25, "Installing ROCm stack...")
	packages := []string{
		"rocm-hip-sdk",
		"python-pytorch-rocm",
		"python-numpy",
		"python-pip",
		"git",
		"base-devel",
		"fastfetch",
		"vim",
	}

	for i, pkg := range packages {
		ui.Progress(25+(i*3), fmt.Sprintf("Installing %s...", pkg))
		_, err := lxd.ExecInContainer(ctx, name, "pacman", "-S", "--needed", "--noconfirm", pkg)
		if err != nil {
			ui.Log(core.LogWarn, fmt.Sprintf("Failed to install %s: %v", pkg, err))
		}
	}

	ui.Log(core.LogInfo, "✓ ai-lab container ready")
	return nil
}

func (s *WorkspaceStage) createDevLab(ctx context.Context, ui core.UI, lxd *system.LXD) error {
	name := "dev-lab"

	// Check if already exists
	if lxd.ContainerExists(ctx, name) {
		ui.Log(core.LogInfo, fmt.Sprintf("Container %s already exists, skipping creation", name))
		return nil
	}

	// Create container
	ui.Log(core.LogInfo, "Launching dev-lab container from archlinux image...")
	if err := lxd.CreateContainer(ctx, name, "images:archlinux/current"); err != nil {
		return fmt.Errorf("failed to create %s: %v", name, err)
	}

	// Wait for network
	ui.Progress(65, "Waiting for container network...")
	if err := lxd.WaitForNetwork(ctx, name); err != nil {
		ui.Log(core.LogWarn, "Network wait timed out, continuing anyway")
	}

	// Install development packages
	ui.Progress(70, "Installing development tools...")
	packages := []string{
		"base-devel",
		"git",
		"rust",
		"go",
		"nodejs",
		"npm",
		"python",
		"python-pip",
		"vim",
		"neovim",
		"fastfetch",
	}

	for i, pkg := range packages {
		ui.Progress(70+(i*2), fmt.Sprintf("Installing %s...", pkg))
		_, err := lxd.ExecInContainer(ctx, name, "pacman", "-S", "--needed", "--noconfirm", pkg)
		if err != nil {
			ui.Log(core.LogWarn, fmt.Sprintf("Failed to install %s: %v", pkg, err))
		}
	}

	ui.Log(core.LogInfo, "✓ dev-lab container ready")
	return nil
}

func (s *WorkspaceStage) Rollback(ctx context.Context) error {
	lxd := system.NewLXD()
	lxd.DeleteContainer(ctx, "ai-lab", true)
	lxd.DeleteContainer(ctx, "dev-lab", true)
	return nil
}
