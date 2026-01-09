package stages

import (
	"context"

	"github.com/daveweinstein1/strix-installer/pkg/core"
	"github.com/daveweinstein1/strix-installer/pkg/system"
)

// CleanupStage removes orphaned packages and cleans cache
type CleanupStage struct{}

func NewCleanupStage() *CleanupStage { return &CleanupStage{} }

func (s *CleanupStage) ID() string          { return "cleanup" }
func (s *CleanupStage) Name() string        { return "Cleanup" }
func (s *CleanupStage) Description() string { return "Remove orphaned packages, clean package cache" }
func (s *CleanupStage) Optional() bool      { return true }

func (s *CleanupStage) Run(ctx context.Context, ui core.UI) error {
	pacman := system.NewPacman()

	// Step 1: Remove orphaned packages
	ui.Progress(30, "Removing orphaned packages...")
	if err := pacman.CleanOrphans(ctx); err != nil {
		ui.Log(core.LogWarn, "No orphans to remove or cleanup failed")
	} else {
		ui.Log(core.LogInfo, "✓ Orphaned packages removed")
	}

	// Step 2: Clean package cache
	ui.Progress(70, "Cleaning package cache...")
	if err := pacman.CleanCache(ctx); err != nil {
		ui.Log(core.LogWarn, "Cache cleanup had issues")
	} else {
		ui.Log(core.LogInfo, "✓ Package cache cleaned")
	}

	ui.Progress(100, "Cleanup complete")
	return nil
}

func (s *CleanupStage) Rollback(ctx context.Context) error {
	return nil
}
