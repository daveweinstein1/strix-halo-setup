package main

import (
	"context"
	"fmt"

	"github.com/daveweinstein1/strix-installer/pkg/containerhub"
	"github.com/daveweinstein1/strix-installer/pkg/platform/strixhalo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	marketplaceMgr *containerhub.Manager
}

// NewApp creates a new App application struct
func NewApp() *App {
	mgr := containerhub.NewManager()
	// Ignore error for now, defaults used if file missing
	_ = mgr.LoadConfigFromPath("configs/registries.yaml")

	return &App{
		marketplaceMgr: mgr,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogInfo(ctx, "Strix Installer GUI Started")
}

// GetSystemStatus returns detected hardware info
func (a *App) GetSystemStatus() map[string]string {
	status := make(map[string]string)

	device, err := strixhalo.Detect(a.ctx)
	if err != nil {
		status["error"] = err.Error()
		return status
	}

	status["name"] = device.Name()
	status["manufacturer"] = device.Manufacturer()
	status["model"] = device.Model()
	status["quirks_count"] = fmt.Sprintf("%d", len(device.Quirks()))

	return status
}

// FetchMarketplaceImages returns all available images from configured registries
func (a *App) FetchMarketplaceImages() []containerhub.Image {
	images, err := a.marketplaceMgr.FetchAllImages(a.ctx)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to fetch images: %v", err))
		return []containerhub.Image{}
	}
	return images
}
