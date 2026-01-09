package stages

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/daveweinstein1/strix-installer/pkg/core"
	"github.com/daveweinstein1/strix-installer/pkg/system"
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
	grub := system.NewGrub()

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

	// Step 2: Backup GRUB
	ui.Progress(25, "Backing up GRUB configuration...")
	backupPath, err := grub.Backup(ctx)
	if err != nil {
		ui.Log(core.LogWarn, fmt.Sprintf("Could not backup GRUB: %v", err))
	} else {
		ui.Log(core.LogInfo, fmt.Sprintf("GRUB backup: %s", backupPath))
	}

	// Step 3: Add required kernel parameters
	ui.Progress(40, "Configuring kernel parameters...")

	// IOMMU for GPU passthrough
	if err := grub.AddCmdlineParam(ctx, "iommu=pt"); err != nil {
		ui.Log(core.LogWarn, fmt.Sprintf("Could not add iommu=pt: %v", err))
	}

	// AMD P-State driver
	if err := grub.AddCmdlineParam(ctx, "amd_pstate=active"); err != nil {
		ui.Log(core.LogWarn, fmt.Sprintf("Could not add amd_pstate=active: %v", err))
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

	// Step 6: Update GRUB
	ui.Progress(90, "Updating GRUB configuration...")
	if err := grub.UpdateGrub(ctx); err != nil {
		return fmt.Errorf("failed to update GRUB: %v", err)
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
