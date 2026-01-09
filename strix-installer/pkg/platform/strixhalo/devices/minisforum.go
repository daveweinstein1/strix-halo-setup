package devices

import (
	"context"

	"github.com/daveweinstein1/strix-installer/pkg/core"
)

// MinisforumS1Max represents the Minisforum MS-S1 Max with Strix Halo
type MinisforumS1Max struct {
	BaseDevice
}

// NewMinisforumS1Max creates a new Minisforum S1 Max device
func NewMinisforumS1Max(manufacturer, product string) *MinisforumS1Max {
	device := &MinisforumS1Max{
		BaseDevice: BaseDevice{
			Manufacturer_: manufacturer,
			Product_:      product,
		},
	}
	device.Quirks_ = device.buildQuirks()
	return device
}

// Name returns the device display name
func (d *MinisforumS1Max) Name() string {
	return "Minisforum MS-S1 Max"
}

// buildQuirks returns Minisforum-specific quirks
func (d *MinisforumS1Max) buildQuirks() []core.Quirk {
	return []core.Quirk{
		{
			ID:          "ethernet-unreliable",
			Description: "Onboard Ethernet may be unreliable. Consider using a USB Ethernet adapter.",
			Type:        core.QuirkAdvisory,
			Apply: func(ctx context.Context) error {
				// Advisory only
				return nil
			},
		},
		{
			ID:          "usb4-display-issues",
			Description: "USB4 display output may not work. Use HDMI instead.",
			Type:        core.QuirkAdvisory,
			Apply: func(ctx context.Context) error {
				// Advisory only
				return nil
			},
		},
		{
			ID:          "fan-sleep-issues",
			Description: "Fan may run loud during sleep. Check BIOS for sleep settings.",
			Type:        core.QuirkAdvisory,
			Apply: func(ctx context.Context) error {
				// Advisory only
				return nil
			},
		},
	}
}
