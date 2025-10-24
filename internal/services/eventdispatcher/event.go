package eventdispatcher

import "sync"

// Event defines the interface for any event.
type Event interface {
	Name() string
}

// PropagatableEvent extends Event with propagation control.
type PropagatableEvent interface {
	Event
	StopPropagation()
	IsPropagationStopped() bool
}

// BaseEvent provides a basic implementation of the Event and PropagatableEvent interfaces.
type BaseEvent struct {
	name    string
	mu      sync.RWMutex
	stopped bool
}

// NewBaseEvent creates a new BaseEvent with the given name.
func NewBaseEvent(name string) *BaseEvent {
	return &BaseEvent{
		name: name,
	}
}

// Name returns the name of the event.
func (e *BaseEvent) Name() string {
	return e.name
}

// StopPropagation sets the event's propagation status to stopped.
func (e *BaseEvent) StopPropagation() {
	e.mu.Lock()
	e.stopped = true
	e.mu.Unlock()
}

// IsPropagationStopped returns true if a listener has called StopPropagation().
func (e *BaseEvent) IsPropagationStopped() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.stopped
}
