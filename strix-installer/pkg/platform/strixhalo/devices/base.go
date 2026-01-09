package devices

import (
	"github.com/daveweinstein1/strix-installer/pkg/core"
)

// BaseDevice provides common device functionality
type BaseDevice struct {
	Manufacturer_ string
	Product_      string
	Quirks_       []core.Quirk
}

func (d *BaseDevice) Manufacturer() string { return d.Manufacturer_ }
func (d *BaseDevice) Model() string        { return d.Product_ }
func (d *BaseDevice) Quirks() []core.Quirk { return d.Quirks_ }
