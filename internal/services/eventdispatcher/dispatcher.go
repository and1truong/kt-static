package eventdispatcher

import (
	"context"
	"fmt"
	"sync"

	"htruong/kt-crawler/internal/services/logging"
)

// Dispatcher manages event listeners and dispatches events.
type Dispatcher struct {
	listeners map[string][]Listener
	mu        sync.RWMutex
	logger    logging.Logger
}

// NewDispatcher creates a new Dispatcher.
func NewDispatcher(logger logging.Logger) *Dispatcher {
	return &Dispatcher{
		listeners: make(map[string][]Listener),
		logger:    logger,
	}
}

// Register adds a listener for a specific event name.
func (d *Dispatcher) Register(eventName string, listener Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.listeners[eventName] = append(d.listeners[eventName], listener)
}

// Dispatch sends an event to all registered listeners.
// If the event implements PropagatableEvent, listeners are executed sequentially
// with propagation control. Otherwise, they are executed concurrently.
func (d *Dispatcher) Dispatch(ctx context.Context, event Event) error {
	eventName := event.Name()
	d.logger.Info("Dispatching event", "event", eventName)

	d.mu.RLock()
	listeners, ok := d.listeners[eventName]
	d.mu.RUnlock()

	if !ok {
		d.logger.Warn("No listeners registered for event", "event", eventName)
		return nil
	}

	// Check for PropagatableEvent and use sequential dispatch if found
	if pe, ok := event.(PropagatableEvent); ok {
		d.logger.Debug("Dispatching event sequentially", "event", eventName, "listeners", len(listeners))
		for _, listener := range listeners {
			listenerType := fmt.Sprintf("%T", listener)
			if ctx.Err() != nil {
				d.logger.Warn("Context cancelled during sequential dispatch", "event", eventName, "error", ctx.Err())
				return ctx.Err()
			}
			if pe.IsPropagationStopped() {
				d.logger.Info("Event propagation stopped", "event", eventName)
				return nil
			}
			if err := listener.Handle(ctx, event); err != nil {
				d.logger.Error("Listener failed during sequential dispatch", "event", eventName, "listener", listenerType, "error", err)
				return err
			}
			d.logger.Debug("Listener handled event sequentially", "event", eventName, "listener", listenerType)
		}
		return nil
	}

	// Default: Concurrent dispatch
	d.logger.Info("Dispatching event concurrently", "event", eventName, "listeners", len(listeners))
	var wg sync.WaitGroup
	errCh := make(chan error, len(listeners))

	for _, listener := range listeners {
		// Check context before starting a new goroutine
		if ctx.Err() != nil {
			d.logger.Warn("Context cancelled before starting concurrent listener", "event", eventName, "error", ctx.Err())
			return ctx.Err()
		}

		wg.Add(1)
		go func(l Listener) {
			defer wg.Done()
			listenerType := fmt.Sprintf("%T", l)
			d.logger.Info("Handling event concurrently", "event", eventName, "listener", listenerType)
			if err := l.Handle(ctx, event); err != nil {
				d.logger.Error("Listener failed during concurrent dispatch", "event", eventName, "listener", listenerType, "error", err)
				errCh <- err
			}
		}(listener)
	}

	// Wait for all listeners to complete
	wg.Wait()
	close(errCh)

	// Collect and return the first error, if any
	if err, ok := <-errCh; ok {
		d.logger.Error("One or more listeners failed during concurrent dispatch", "event", eventName, "first_error", err)
		return err
	}

	// Check context one last time in case it was canceled while waiting
	if ctx.Err() != nil {
		d.logger.Warn("Context cancelled after concurrent dispatch completed", "event", eventName, "error", ctx.Err())
		return ctx.Err()
	}

	d.logger.Debug("Event dispatch completed successfully", "event", eventName)
	return nil
}
