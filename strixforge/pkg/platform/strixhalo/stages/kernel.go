package stages

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/daveweinstein1/strixforge/pkg/core"

	"github.com/daveweinstein1/strixforge/pkg/system/bootloader"
)

// KernelStage configures kernel and bootloader
type KernelStage struct {
	device core.Device
}

// NewKernelStage creates a new kernel configuration stage
func NewKernelStage(device core.Device) *KernelStage {
	return &KernelStage{device: device}
}

func (s *KernelStage) ID() string   { return "kernel" }
func (s *KernelStage) Name() string { return "Kernel Configuration" }
func (s *KernelStage) Description() string {
	return "Verify kernel version, configure GRUB, apply device quirks"
}
func (s *KernelStage) Optional() bool { return false }

func (s *KernelStage) Run(ctx context.Context, ui core.UI) error {
	// Step 1: Check kernel version
	ui.Progress(10, "Checking kernel version...")
	version, err := getKernelVersion()
	if err != nil {
		return fmt.Errorf("failed to get kernel version: %v", err)
	}

	major, minor := parseVersion(version)
	ui.Log(core.LogInfo, fmt.Sprintf("Kernel version: %s (parsed: %d.%d)", version, major, minor))

	// Require 6.18+
	if major < 6 || (major == 6 && minor < 18) {
		return fmt.Errorf("kernel 6.18+ required, found %s. Please update your kernel", version)
	}
	ui.Log(core.LogInfo, "✓ Kernel version meets requirements")

	// Step 2: Configure Bootloader(s)
	// We now support GRUB, systemd-boot, Limine, and rEFInd
	loaders := bootloader.Detect()
	if len(loaders) > 0 {
		ui.Log(core.LogInfo, fmt.Sprintf("Detected %d active bootloader(s)", len(loaders)))

		for _, loader := range loaders {
			ui.Log(core.LogInfo, fmt.Sprintf("Configuring: %s", loader.Name()))

			// Backup
			ui.Progress(25, fmt.Sprintf("Backing up %s config...", loader.Name()))
			if backupPath, err := loader.Backup(ctx); err != nil {
				ui.Log(core.LogWarn, fmt.Sprintf("Could not backup %s: %v", loader.Name(), err))
			} else {
				ui.Log(core.LogInfo, fmt.Sprintf("Backup created: %s", backupPath))
			}

			// Add Parameters
			ui.Progress(40, fmt.Sprintf("Adding kernel parameters to %s...", loader.Name()))

			// IOMMU
			if err := loader.AddParam(ctx, "iommu=pt"); err != nil {
				ui.Log(core.LogWarn, fmt.Sprintf("Failed to add iommu=pt to %s: %v", loader.Name(), err))
			}
			// AMD P-State
			if err := loader.AddParam(ctx, "amd_pstate=active"); err != nil {
				ui.Log(core.LogWarn, fmt.Sprintf("Failed to add amd_pstate=active to %s: %v", loader.Name(), err))
			}
		}
	} else {
		// No bootloader detected (or custom setup)
		ui.Log(core.LogWarn, "⚠ No supported bootloader detected! (Checked: GRUB, systemd-boot, Limine, rEFInd)")
		ui.Log(core.LogWarn, "MANUAL ACTION REQUIRED: Add 'iommu=pt amd_pstate=active' to your kernel arguments.")
	}

	// Step 4: Apply device-specific quirks
	ui.Progress(60, "Applying device quirks...")
	for _, quirk := range s.device.Quirks() {
		if quirk.Type == core.QuirkAuto {
			ui.Log(core.LogInfo, fmt.Sprintf("Applying quirk: %s", quirk.Description))
			if err := quirk.Apply(ctx); err != nil {
				ui.Log(core.LogWarn, fmt.Sprintf("Quirk failed: %v", err))
			}
		} else {
			// Advisory quirk - just inform user
			ui.Log(core.LogWarn, fmt.Sprintf("ADVISORY: %s", quirk.Description))
		}
	}

	// Step 5: ZRAM Optimization
	ui.Progress(75, "Checking ZRAM config...")
	ramGB, err := getTotalRAM()
	if err == nil {
		if ramGB >= 64 {
			ui.Log(core.LogInfo, fmt.Sprintf("High memory system (%d GB) detected. Disabling ZRAM to prevent GTT conflicts.", ramGB))

			// Disable ZRAM generator service
			cmd := exec.CommandContext(ctx, "systemctl", "disable", "--now", "zram-generator@zram0.service")
			if out, err := cmd.CombinedOutput(); err != nil {
				// Don't fail if service doesn't exist, just log
				ui.Log(core.LogWarn, fmt.Sprintf("Failed to disable ZRAM (might not be active): %v %s", err, string(out)))
			} else {
				ui.Log(core.LogInfo, "✓ ZRAM disabled")
			}
		} else {
			ui.Log(core.LogInfo, fmt.Sprintf("System memory %d GB < 64GB. Keeping ZRAM enabled.", ramGB))
		}
	} else {
		ui.Log(core.LogWarn, fmt.Sprintf("Could not determine system RAM: %v. Skipping ZRAM optimization.", err))
	}

	ui.Progress(100, "Kernel configuration complete")
	ui.Log(core.LogInfo, "NOTE: Reboot may be required for kernel parameter changes")

	return nil
}

func (s *KernelStage) Rollback(ctx context.Context) error {
	// Could restore GRUB backup here
	return nil
}

// getKernelVersion returns the current kernel version
func getKernelVersion() (string, error) {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// parseVersion extracts major.minor from kernel version string
func parseVersion(version string) (int, int) {
	re := regexp.MustCompile(`^(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) < 3 {
		return 0, 0
	}
	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	return major, minor
}

// getTotalRAM returns system RAM in GB
func getTotalRAM() (int, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				kb, err := strconv.Atoi(fields[1])
				if err != nil {
					return 0, err
				}
				// Convert kB to GB check
				return kb / 1024 / 1024, nil
			}
		}
	}
	return 0, fmt.Errorf("MemTotal not found in /proc/meminfo")
}
