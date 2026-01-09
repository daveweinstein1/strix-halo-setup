package core

import (
	"context"
	"fmt"
	"time"
)

// Engine orchestrates the installation process
type Engine struct {
	platform Platform
	ui       UI
	bus      *EventBus
	results  []StageResult
	dryRun   bool
}

// Platform defines the interface for a target platform (e.g., Strix Halo)
type Platform interface {
	Name() string
	Detect() (Device, error)
	Stages() []Stage
	Validate() error
}

// Device represents detected hardware
type Device interface {
	Name() string
	Manufacturer() string
	Model() string
	Quirks() []Quirk
}

// Quirk represents a device-specific fix or advisory
type Quirk struct {
	ID          string
	Description string
	Type        QuirkType
	Apply       func(ctx context.Context) error
}

// QuirkType indicates if a quirk is automatic or advisory
type QuirkType int

const (
	QuirkAuto     QuirkType = iota // Applied automatically
	QuirkAdvisory                  // Shown to user as recommendation
)

// NewEngine creates a new installation engine
func NewEngine(platform Platform, ui UI) *Engine {
	return &Engine{
		platform: platform,
		ui:       ui,
		bus:      NewEventBus(),
		results:  make([]StageResult, 0),
		dryRun:   false,
	}
}

// SetDryRun enables dry-run mode (no actual changes)
func (e *Engine) SetDryRun(dryRun bool) {
	e.dryRun = dryRun
}

// EventBus returns the event bus for UI subscription
func (e *Engine) EventBus() *EventBus {
	return e.bus
}

// Run executes all stages in order
func (e *Engine) Run(ctx context.Context) error {
	stages := e.platform.Stages()

	for i, stage := range stages {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Skip optional stages if user declines
		if stage.Optional() {
			if !e.ui.Confirm(fmt.Sprintf("Run optional stage: %s?", stage.Name()), true) {
				result := StageResult{
					StageID:   stage.ID(),
					StageName: stage.Name(),
					Status:    StatusSkipped,
				}
				e.results = append(e.results, result)
				e.bus.Publish(StageCompletedEvent{Stage: stage, Result: result})
				continue
			}
		}

		// Run the stage
		result := e.runStage(ctx, stage, i+1, len(stages))
		e.results = append(e.results, result)

		// Stop on failure
		if result.Status == StatusFailed {
			return result.Error
		}
	}

	return nil
}

// runStage executes a single stage with timing and error handling
func (e *Engine) runStage(ctx context.Context, stage Stage, num, total int) StageResult {
	e.ui.Log(LogInfo, fmt.Sprintf("[%d/%d] Starting: %s", num, total, stage.Name()))
	e.ui.StageStart(stage)
	e.bus.Publish(StageStartedEvent{Stage: stage})

	start := time.Now()

	var err error
	if !e.dryRun {
		err = stage.Run(ctx, e.ui)
	} else {
		e.ui.Log(LogInfo, "[DRY RUN] Would execute stage")
	}

	duration := time.Since(start)

	result := StageResult{
		StageID:   stage.ID(),
		StageName: stage.Name(),
		Duration:  duration,
		Error:     err,
	}

	if err != nil {
		result.Status = StatusFailed
		e.ui.Log(LogError, fmt.Sprintf("Failed: %s - %v", stage.Name(), err))
	} else {
		result.Status = StatusSuccess
		e.ui.Log(LogInfo, fmt.Sprintf("Complete: %s (%v)", stage.Name(), duration.Round(time.Second)))
	}

	e.ui.StageComplete(result)
	e.bus.Publish(StageCompletedEvent{Stage: stage, Result: result})

	return result
}

// Results returns all stage results
func (e *Engine) Results() []StageResult {
	return e.results
}
