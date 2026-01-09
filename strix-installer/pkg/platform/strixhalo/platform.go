package strixhalo

import (
	"context"

	"github.com/daveweinstein1/strix-installer/pkg/core"
	"github.com/daveweinstein1/strix-installer/pkg/platform/strixhalo/stages"
)

// Platform implements the Strix Halo installation platform
type Platform struct {
	device core.Device
}

// New creates a new Strix Halo platform
func New() *Platform {
	return &Platform{}
}

// Name returns the platform display name
func (p *Platform) Name() string {
	return "AMD Strix Halo (gfx1151)"
}

// Detect identifies the specific hardware device
func (p *Platform) Detect() (core.Device, error) {
	device, err := Detect(context.Background())
	if err != nil {
		return nil, err
	}
	p.device = device
	return device, nil
}

// Stages returns all installation stages in order
func (p *Platform) Stages() []core.Stage {
	return []core.Stage{
		stages.NewKernelStage(p.device),
		stages.NewGraphicsStage(),
		stages.NewSystemStage(),
		stages.NewLXDStage(),
		stages.NewCleanupStage(),
		stages.NewValidateStage(),
		stages.NewAppsStage(),
		stages.NewWorkspaceStage(),
	}
}

// Validate checks prerequisites (called before running stages)
func (p *Platform) Validate() error {
	// Could check for CachyOS, Arch-based distro, etc.
	return nil
}
