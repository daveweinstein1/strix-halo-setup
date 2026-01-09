package stages

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/daveweinstein1/strix-installer/pkg/core"
	"github.com/daveweinstein1/strix-installer/pkg/system"
)

// GraphicsStage installs and verifies graphics stack
type GraphicsStage struct{}

func NewGraphicsStage() *GraphicsStage { return &GraphicsStage{} }

func (s *GraphicsStage) ID() string   { return "graphics" }
func (s *GraphicsStage) Name() string { return "Graphics Stack" }
func (s *GraphicsStage) Description() string {
	return "Install Mesa 25.3+, Vulkan, LLVM 21.x, firmware"
}
func (s *GraphicsStage) Optional() bool { return false }

func (s *GraphicsStage) Run(ctx context.Context, ui core.UI) error {
	pacman := system.NewPacman()

	// Step 1: Install graphics packages
	ui.Progress(10, "Installing graphics packages...")
	packages := []string{
		"mesa", "lib32-mesa", "mesa-utils",
		"vulkan-radeon", "lib32-vulkan-radeon", "vulkan-tools",
		"linux-firmware",
		"llvm", "lib32-llvm",
	}

	if err := pacman.Install(ctx, packages...); err != nil {
		return fmt.Errorf("failed to install graphics packages: %v", err)
	}
	ui.Log(core.LogInfo, "✓ Graphics packages installed")

	// Step 2: Verify Mesa version
	ui.Progress(50, "Verifying Mesa version...")
	mesaVersion, err := pacman.GetVersion(ctx, "mesa")
	if err != nil {
		return fmt.Errorf("failed to get Mesa version: %v", err)
	}

	major, minor := parseMesaVersion(mesaVersion)
	ui.Log(core.LogInfo, fmt.Sprintf("Mesa version: %s (parsed: %d.%d)", mesaVersion, major, minor))

	if major < 25 || (major == 25 && minor < 3) {
		return fmt.Errorf("Mesa 25.3+ required, found %s", mesaVersion)
	}
	ui.Log(core.LogInfo, "✓ Mesa version meets requirements")

	// Step 3: Verify LLVM version
	ui.Progress(70, "Verifying LLVM version...")
	llvmVersion, err := pacman.GetVersion(ctx, "llvm")
	if err != nil {
		ui.Log(core.LogWarn, fmt.Sprintf("Could not verify LLVM version: %v", err))
	} else {
		llvmMajor, _ := parseMesaVersion(llvmVersion)
		ui.Log(core.LogInfo, fmt.Sprintf("LLVM version: %s", llvmVersion))

		if llvmMajor < 21 {
			return fmt.Errorf("LLVM 21.x required, found %s", llvmVersion)
		}
		ui.Log(core.LogInfo, "✓ LLVM version meets requirements")
	}

	// Step 4: Quick Vulkan check
	ui.Progress(90, "Checking Vulkan...")
	result, err := system.Exec(ctx, "vulkaninfo", "--summary")
	if err != nil {
		ui.Log(core.LogWarn, "vulkaninfo check failed - Vulkan may not be working")
	} else if result.ExitCode == 0 {
		ui.Log(core.LogInfo, "✓ Vulkan is functional")
	}

	ui.Progress(100, "Graphics stack complete")
	return nil
}

func (s *GraphicsStage) Rollback(ctx context.Context) error {
	return nil
}

// parseMesaVersion extracts major.minor from version string
func parseMesaVersion(version string) (int, int) {
	re := regexp.MustCompile(`^(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) < 3 {
		return 0, 0
	}
	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	return major, minor
}
