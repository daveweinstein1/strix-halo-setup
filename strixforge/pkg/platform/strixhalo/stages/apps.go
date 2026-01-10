package stages

import (
	"context"
	"fmt"
	"os"
	"os/user"

	"github.com/daveweinstein1/strixforge/pkg/core"
	"github.com/daveweinstein1/strixforge/pkg/system"
)

// AppsStage installs desktop applications
type AppsStage struct{}

func NewAppsStage() *AppsStage { return &AppsStage{} }

func (s *AppsStage) ID() string          { return "apps" }
func (s *AppsStage) Name() string        { return "Desktop Software" }
func (s *AppsStage) Description() string { return "Install browsers, office suite, and utilities" }
func (s *AppsStage) Optional() bool      { return true }

func (s *AppsStage) Run(ctx context.Context, ui core.UI) error {
	pacman := system.NewPacman()

	// Get current user for AUR operations
	username := os.Getenv("SUDO_USER")
	if username == "" {
		if u, err := user.Current(); err == nil {
			username = u.Username
		}
	}
	yay := system.NewYay(username)

	// Step 1: Install yay (AUR helper)
	ui.Progress(5, "Setting up AUR helper...")
	if !pacman.IsInstalled(ctx, "yay") {
		if err := pacman.Install(ctx, "yay"); err != nil {
			ui.Log(core.LogWarn, fmt.Sprintf("Could not install yay: %v", err))
		} else {
			ui.Log(core.LogInfo, "✓ yay installed")
		}
	}

	// Step 2: Official repo packages
	ui.Progress(15, "Installing browsers and utilities...")
	officialPackages := []string{
		"firefox",
		"vlc",
		"signal-desktop",
	}
	if err := pacman.Install(ctx, officialPackages...); err != nil {
		ui.Log(core.LogWarn, fmt.Sprintf("Some official packages failed: %v", err))
	}
	ui.Log(core.LogInfo, "✓ Firefox, VLC, Signal installed")

	// Step 3: AUR packages (optional, ask user)
	ui.Progress(40, "AUR packages...")

	aurPackages := map[string]string{
		"google-chrome":          "Google Chrome",
		"ungoogled-chromium-bin": "Ungoogled Chromium",
		"helium":                 "Helium Browser",
		"onlyoffice-bin":         "OnlyOffice",
	}

	for pkg, name := range aurPackages {
		if ui.Confirm(fmt.Sprintf("Install %s?", name), false) {
			ui.Log(core.LogInfo, fmt.Sprintf("Installing %s...", name))
			if err := yay.Install(ctx, pkg); err != nil {
				ui.Log(core.LogWarn, fmt.Sprintf("Failed to install %s: %v", name, err))
			} else {
				ui.Log(core.LogInfo, fmt.Sprintf("✓ %s installed", name))
			}
		}
	}

	// Step 4: Special Suites (NordVPN)
	if ui.Confirm("Install NordVPN Suite (CLI + KDE Tray Icon)?", false) {
		ui.Progress(60, "Installing NordVPN Suite...")

		// 1. Install packages
		nordPackages := []string{"nordvpn-bin", "nordvpn-plasmoid"}
		if err := yay.Install(ctx, nordPackages...); err != nil {
			ui.Log(core.LogWarn, fmt.Sprintf("Failed to install NordVPN packages: %v", err))
		} else {
			ui.Log(core.LogInfo, "✓ NordVPN packages installed")

			// 2. Configure Group
			// usermod -aG nordvpn $USER
			if _, err := system.ExecSudo(ctx, "usermod", "-aG", "nordvpn", username); err != nil {
				ui.Log(core.LogWarn, fmt.Sprintf("Failed to add user to nordvpn group: %v", err))
			} else {
				ui.Log(core.LogInfo, "✓ User added to 'nordvpn' group")
			}

			// 3. Enable Service
			// systemctl enable --now nordvpnd
			sysd := system.NewSystemd()
			if err := sysd.EnableAndStart(ctx, "nordvpnd"); err != nil {
				ui.Log(core.LogWarn, fmt.Sprintf("Failed to enable nordvpnd: %v", err))
			} else {
				ui.Log(core.LogInfo, "✓ NordVPN service enabled and started")
			}
		}
	}

	ui.Progress(100, "Desktop software installation complete")
	return nil
}

func (s *AppsStage) Rollback(ctx context.Context) error {
	return nil
}
