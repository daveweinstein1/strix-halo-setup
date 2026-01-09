package core

import (
	"context"
	"time"
)

// Status represents the outcome of a stage execution
type Status int

const (
	StatusPending Status = iota
	StatusRunning
	StatusSuccess
	StatusFailed
	StatusSkipped
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusRunning:
		return "running"
	case StatusSuccess:
		return "success"
	case StatusFailed:
		return "failed"
	case StatusSkipped:
		return "skipped"
	default:
		return "unknown"
	}
}

// Stage defines the interface for an installation stage
type Stage interface {
	ID() string
	Name() string
	Description() string
	Optional() bool
	Run(ctx context.Context, ui UI) error
	Rollback(ctx context.Context) error
}

// StageResult captures the outcome of running a stage
type StageResult struct {
	StageID   string
	StageName string
	Status    Status
	Duration  time.Duration
	Error     error
	Logs      []LogEntry
}

// LogEntry represents a single log message
type LogEntry struct {
	Time    time.Time
	Level   LogLevel
	Message string
}

// LogLevel indicates severity of log messages
type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
)

func (l LogLevel) String() string {
	switch l {
	case LogDebug:
		return "DEBUG"
	case LogInfo:
		return "INFO"
	case LogWarn:
		return "WARN"
	case LogError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
