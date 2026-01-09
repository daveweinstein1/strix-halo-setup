package strixhalo

import (
	"context"
	"strings"

	"github.com/daveweinstein1/strix-installer/pkg/core"
	"github.com/daveweinstein1/strix-installer/pkg/platform/strixhalo/devices"
	"github.com/daveweinstein1/strix-installer/pkg/system"
)

// Detect identifies the specific Strix Halo device
func Detect(ctx context.Context) (core.Device, error) {
	dmi := system.NewDMIDecode()

	manufacturer, _ := dmi.GetSystemManufacturer(ctx)
	product, _ := dmi.GetProductName(ctx)

	manufacturer = strings.ToLower(manufacturer)
	product = strings.ToLower(product)

	switch {
	case strings.Contains(manufacturer, "beelink"):
		return devices.NewBeelinkGTR9(manufacturer, product), nil
	case strings.Contains(manufacturer, "framework"):
		return devices.NewFrameworkDesktop(manufacturer, product), nil
	case strings.Contains(manufacturer, "minisforum"):
		return devices.NewMinisforumS1Max(manufacturer, product), nil
	default:
		return devices.NewGenericDevice(manufacturer, product), nil
	}
}
