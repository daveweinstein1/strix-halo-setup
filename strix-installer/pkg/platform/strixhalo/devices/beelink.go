package devices

import (
	"context"

	"github.com/daveweinstein1/strix-installer/pkg/core"
	"github.com/daveweinstein1/strix-installer/pkg/system"
)

// BeelinkGTR9 represents the Beelink GTR9 Pro with Strix Halo
type BeelinkGTR9 struct {
	BaseDevice
}

// NewBeelinkGTR9 creates a new Beelink GTR9 device
func NewBeelinkGTR9(manufacturer, product string) *BeelinkGTR9 {
	device := &BeelinkGTR9{
		BaseDevice: BaseDevice{
			Manufacturer_: manufacturer,
			Product_:      product,
		},
	}
	device.Quirks_ = device.buildQuirks()
	return device
}

// Name returns the device display name
func (d *BeelinkGTR9) Name() string {
	return "Beelink GTR9 Pro"
}

// buildQuirks returns Beelink-specific quirks
func (d *BeelinkGTR9) buildQuirks() []core.Quirk {
	return []core.Quirk{
		{
			ID:          "e610-blacklist",
			Description: "Blacklist Intel E610 Ethernet driver (crashes under GPU load)",
			Type:        core.QuirkAuto,
			Apply: func(ctx context.Context) error {
				grub := system.NewGrub()
				return grub.AddCmdlineParam(ctx, "modprobe.blacklist=ice")
			},
		},
		{
			ID:          "tdp-tool",
			Description: "Install RyzenAdj for TDP control",
			Type:        core.QuirkAuto,
			Apply: func(ctx context.Context) error {
				pacman := system.NewPacman()
				// ryzenadj may be in AUR
				if !pacman.IsInstalled(ctx, "ryzenadj") {
					// Try AUR via yay
					yay := system.NewYay("") // Will use current user
					return yay.Install(ctx, "ryzenadj")
				}
				return nil
			},
		},
	}
}
