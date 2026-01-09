package core

// Event types for the event system
// Both TUI and GUI subscribe to these events

type Event interface {
	eventMarker()
}

// StageStartedEvent is emitted when a stage begins
type StageStartedEvent struct {
	Stage Stage
}

func (e StageStartedEvent) eventMarker() {}

// StageCompletedEvent is emitted when a stage finishes
type StageCompletedEvent struct {
	Stage  Stage
	Result StageResult
}

func (e StageCompletedEvent) eventMarker() {}

// ProgressEvent is emitted for progress updates
type ProgressEvent struct {
	Percent int
	Message string
}

func (e ProgressEvent) eventMarker() {}

// LogEvent is emitted for log messages
type LogEvent struct {
	Level   LogLevel
	Message string
}

func (e LogEvent) eventMarker() {}

// PromptEvent is emitted when user input is needed
type PromptEvent struct {
	Type     PromptType
	Message  string
	Options  []string
	Default  interface{}
	Response chan interface{}
}

func (e PromptEvent) eventMarker() {}

// PromptType indicates the kind of prompt
type PromptType int

const (
	PromptConfirm PromptType = iota
	PromptSelect
	PromptInput
)

// EventBus distributes events to subscribers
type EventBus struct {
	subscribers []chan Event
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make([]chan Event, 0),
	}
}

// Subscribe returns a channel that receives events
func (b *EventBus) Subscribe() chan Event {
	ch := make(chan Event, 100)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

// Publish sends an event to all subscribers
func (b *EventBus) Publish(event Event) {
	for _, ch := range b.subscribers {
		select {
		case ch <- event:
		default:
			// Drop event if channel is full
		}
	}
}

// Close closes all subscriber channels
func (b *EventBus) Close() {
	for _, ch := range b.subscribers {
		close(ch)
	}
}
