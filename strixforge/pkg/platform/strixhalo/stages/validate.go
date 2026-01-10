package stages

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/daveweinstein1/strixforge/pkg/core"
	"github.com/daveweinstein1/strixforge/pkg/system"
)

// ValidateStage verifies the installation
type ValidateStage struct{}

func NewValidateStage() *ValidateStage { return &ValidateStage{} }

func (s *ValidateStage) ID() string   { return "validate" }
func (s *ValidateStage) Name() string { return "Validation" }
func (s *ValidateStage) Description() string {
	return "Verify kernel, GPU, IOMMU, and LXD configuration"
}
func (s *ValidateStage) Optional() bool { return false }

func (s *ValidateStage) Run(ctx context.Context, ui core.UI) error {
	systemd := system.NewSystemd()
	failures := 0

	// Check 1: Kernel version
	ui.Progress(10, "Checking kernel version...")
	version, err := getKernelVersion()
	if err != nil {
		ui.Log(core.LogError, fmt.Sprintf("✗ Could not get kernel version: %v", err))
		failures++
	} else {
		major, minor := parseVersion(version)
		if major < 6 || (major == 6 && minor < 18) {
			ui.Log(core.LogError, fmt.Sprintf("✗ Kernel %s does not meet 6.18+ requirement", version))
			failures++
		} else {
			ui.Log(core.LogInfo, fmt.Sprintf("✓ Kernel %s meets requirements", version))
		}
	}

	// Check 2: Kernel parameters
	ui.Progress(25, "Checking kernel parameters...")

	// We check /proc/cmdline directly (current running kernel)
	cmdlineBytes, err := os.ReadFile("/proc/cmdline")
	hasIOMMU := false
	if err == nil {
		if strings.Contains(string(cmdlineBytes), "iommu=pt") {
			hasIOMMU = true
		}
	}

	if hasIOMMU {
		ui.Log(core.LogInfo, "✓ IOMMU passthrough enabled")
	} else {
		// Log a warning: verify if it's configured in bootloader but we just haven't rebooted?
		// Checking all bootloaders is complex, so for validation we focus on "active" state
		ui.Log(core.LogWarn, "✗ iommu=pt not active (reboot needed?)")
	}

	// Check 3: GPU rendering
	ui.Progress(40, "Checking GPU...")
	glxOut, err := exec.CommandContext(ctx, "glxinfo").Output()
	if err != nil {
		ui.Log(core.LogWarn, "✗ glxinfo not available")
	} else {
		glxStr := string(glxOut)
		if strings.Contains(glxStr, "AMD") || strings.Contains(glxStr, "Radeon") {
			ui.Log(core.LogInfo, "✓ AMD GPU detected in OpenGL renderer")
		} else {
			ui.Log(core.LogWarn, "✗ AMD GPU not detected in OpenGL renderer")
			failures++
		}
	}

	// Check 4: Vulkan
	ui.Progress(55, "Checking Vulkan...")
	result, err := system.Exec(ctx, "vulkaninfo", "--summary")
	if err != nil || result.ExitCode != 0 {
		ui.Log(core.LogWarn, "✗ Vulkan check failed")
		failures++
	} else {
		ui.Log(core.LogInfo, "✓ Vulkan is functional")
	}

	// Check 5: LXD service
	ui.Progress(70, "Checking LXD service...")
	if systemd.IsActive(ctx, "lxd.socket") {
		ui.Log(core.LogInfo, "✓ LXD socket is active")
	} else {
		ui.Log(core.LogError, "✗ LXD socket is not running")
		failures++
	}

	// Check 6: LXD can run
	ui.Progress(85, "Testing LXD access...")
	result, err = system.Exec(ctx, "lxc", "list")
	if err != nil || result.ExitCode != 0 {
		ui.Log(core.LogWarn, "✗ Cannot run 'lxc list' - may need to log out/in for group changes")
	} else {
		ui.Log(core.LogInfo, "✓ LXD access working")
	}

	ui.Progress(100, "Validation complete")

	if failures > 0 {
		return fmt.Errorf("%d validation checks failed", failures)
	}

	ui.Log(core.LogInfo, "All validation checks passed!")
	return nil
}

func (s *ValidateStage) Rollback(ctx context.Context) error {
	return nil
}
