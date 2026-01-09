package devices

import (
	"github.com/daveweinstein1/strix-installer/pkg/core"
)

// GenericDevice represents an unknown Strix Halo device
type GenericDevice struct {
	BaseDevice
}

// NewGenericDevice creates a new generic device
func NewGenericDevice(manufacturer, product string) *GenericDevice {
	return &GenericDevice{
		BaseDevice: BaseDevice{
			Manufacturer_: manufacturer,
			Product_:      product,
			Quirks_:       []core.Quirk{}, // No specific quirks
		},
	}
}

// Name returns the device display name
func (d *GenericDevice) Name() string {
	if d.Product_ != "" {
		return d.Product_
	}
	return "Generic Strix Halo Device"
}
