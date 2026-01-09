package devices

import (
	"context"

	"github.com/daveweinstein1/strix-installer/pkg/core"
)

// FrameworkDesktop represents the Framework Desktop with Strix Halo
type FrameworkDesktop struct {
	BaseDevice
}

// NewFrameworkDesktop creates a new Framework Desktop device
func NewFrameworkDesktop(manufacturer, product string) *FrameworkDesktop {
	device := &FrameworkDesktop{
		BaseDevice: BaseDevice{
			Manufacturer_: manufacturer,
			Product_:      product,
		},
	}
	device.Quirks_ = device.buildQuirks()
	return device
}

// Name returns the device display name
func (d *FrameworkDesktop) Name() string {
	return "Framework Desktop"
}

// buildQuirks returns Framework-specific quirks
func (d *FrameworkDesktop) buildQuirks() []core.Quirk {
	return []core.Quirk{
		{
			ID:          "fan-noise-advisory",
			Description: "At 140W TDP, fan noise is louder. Consider setting TDP to 110-120W in BIOS.",
			Type:        core.QuirkAdvisory,
			Apply: func(ctx context.Context) error {
				// Advisory only - no automatic action
				return nil
			},
		},
	}
}
