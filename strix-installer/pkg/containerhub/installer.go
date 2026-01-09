package containerhub

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/daveweinstein1/strix-installer/pkg/system"
)

type Installer struct {
	lxd *system.LXD
}

func NewInstaller() *Installer {
	return &Installer{
		lxd: system.NewLXD(),
	}
}

// InstallImage installs a selected toolbox image into a target LXD container
// cmd: lxc exec <targetContainer> -- toolbox create <toolboxName> --image <imageURL>
func (i *Installer) InstallImage(ctx context.Context, targetContainer, toolboxName, imageURL string) error {
	// 1. Ensure target container exists
	if !i.lxd.ContainerExists(ctx, targetContainer) {
		return fmt.Errorf("target container '%s' does not exist", targetContainer)
	}

	// 2. Run toolbox create command inside the container
	// Note: toolbox create might prompt or take time. We assume non-interactive here?
	// toolbox create -c <name> -i <image> -y (to auto-accept)

	cmd := exec.CommandContext(ctx, "lxc", "exec", targetContainer, "--",
		"toolbox", "create", "-c", toolboxName, "-i", imageURL, "-y")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("toolbox installation failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}
