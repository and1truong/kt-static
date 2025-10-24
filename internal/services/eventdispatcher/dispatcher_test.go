package eventdispatcher

import (
	"context"
	"fmt"
	"htruong/kt-crawler/internal/services/logging/logging"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- Mocks for Testing ---

type mockEvent struct {
	name string
	data string
}

func (e mockEvent) Name() string {
	return e.name
}

type mockListener struct {
	handledEvents chan string
}

func newMockListener() *mockListener {
	return &mockListener{
		handledEvents: make(chan string, 50), // Buffered channel
	}
}

func (l *mockListener) Handle(ctx context.Context, event Event) error {
	if e, ok := event.(mockEvent); ok {
		msg := fmt.Sprintf("%s:%s", e.Name(), e.data)
		l.handledEvents <- msg
	}
	return nil
}

// --- New Mock Listener for Sequential Tests ---

type sequentialListener struct {
	id            int
	executionLog  *[]int
	shouldStop    bool
	handledEvents chan string
}

func newSequentialListener(id int, log *[]int, shouldStop bool) *sequentialListener {
	return &sequentialListener{
		id:            id,
		executionLog:  log,
		shouldStop:    shouldStop,
		handledEvents: make(chan string, 1),
	}
}

func (l *sequentialListener) Handle(ctx context.Context, event Event) error {
	// Record execution order
	*l.executionLog = append(*l.executionLog, l.id)

	// Check if it's a propagatable event and if we should stop
	if l.shouldStop {
		if pe, ok := event.(PropagatableEvent); ok {
			pe.StopPropagation()
		}
	}

	// Send a signal that we were called (for verification)
	l.handledEvents <- fmt.Sprintf("listener-%d-called", l.id)

	return nil
}

// --- Test Cases ---

func TestNewDispatcher(t *testing.T) {
	d := NewDispatcher(&logging.MockLogger{})
	assert.NotNil(t, d)
	assert.NotNil(t, d.listeners)
}

func TestRegisterAndDispatch(t *testing.T) {
	d := NewDispatcher(&logging.MockLogger{})
	listener := newMockListener()

	eventName := "test.event"
	d.Register(eventName, listener)

	event := mockEvent{name: eventName, data: "payload"}
	assert.NoError(t, d.Dispatch(context.Background(), event))

	select {
	case handled := <-listener.handledEvents:
		assert.Equal(t, "test.event:payload", handled)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Listener was not called within the expected time")
	}
}

func TestDispatch_MultipleListeners(t *testing.T) {
	d := NewDispatcher(&logging.MockLogger{})
	listener1 := newMockListener()
	listener2 := newMockListener()

	eventName := "test.event"
	d.Register(eventName, listener1)
	d.Register(eventName, listener2)

	event := mockEvent{name: eventName, data: "shared_payload"}
	assert.NoError(t, d.Dispatch(context.Background(), event))

	count := 0
	for i := 0; i < 2; i++ {
		select {
		case handled := <-listener1.handledEvents:
			assert.Equal(t, "test.event:shared_payload", handled)
			count++
		case handled := <-listener2.handledEvents:
			assert.Equal(t, "test.event:shared_payload", handled)
			count++
		case <-time.After(100 * time.Millisecond):
			t.Fatal("A listener was not called within the expected time")
		}
	}

	assert.Equal(t, 2, count, "Both listeners should have been called")
}

func TestDispatch_NoListeners(t *testing.T) {
	d := NewDispatcher(&logging.MockLogger{})

	event := mockEvent{name: "unheard.event", data: " lonely_payload"}

	// Dispatch should not panic or block
	assert.NotPanics(t, func() {
		// The error is ignored here as we are only checking for panics
		_ = d.Dispatch(context.Background(), event)
	})
}

func TestDispatch_CorrectListener(t *testing.T) {
	d := NewDispatcher(&logging.MockLogger{})
	listener1 := newMockListener()
	listener2 := newMockListener()

	d.Register("event1", listener1)
	d.Register("event2", listener2)

	event1 := mockEvent{name: "event1", data: "payload1"}
	assert.NoError(t, d.Dispatch(context.Background(), event1))

	// Check listener1 was called
	select {
	case handled := <-listener1.handledEvents:
		assert.Equal(t, "event1:payload1", handled)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Listener 1 was not called")
	}

	// Check listener2 was NOT called
	select {
	case <-listener2.handledEvents:
		t.Fatal("Listener 2 should not have been called")
	default:
		// All good
	}
}

func TestConcurrentDispatch(t *testing.T) {
	d := NewDispatcher(&logging.MockLogger{})
	listener := newMockListener()
	d.Register("concurrent.event", listener)

	var wg sync.WaitGroup
	numDispatches := 50

	for i := 0; i < numDispatches; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			event := mockEvent{name: "concurrent.event", data: fmt.Sprintf("payload-%d", i)}
			// We ignore the error here as the test is focused on concurrency and counting
			_ = d.Dispatch(context.Background(), event)
		}(i)
	}

	wg.Wait()

	close(listener.handledEvents)

	count := 0
	for range listener.handledEvents {
		count++
	}

	assert.Equal(t, numDispatches, count, "Listener should have been called for every dispatch")
}

func TestSequentialDispatchWithPropagationControl(t *testing.T) {
	d := NewDispatcher(&logging.MockLogger{})
	eventName := "sequential.event"
	var executionLog []int

	// Listener 1: Will run and NOT stop propagation
	listener1 := newSequentialListener(1, &executionLog, false)
	// Listener 2: Will run and STOP propagation
	listener2 := newSequentialListener(2, &executionLog, true)
	// Listener 3: Should NOT run
	listener3 := newSequentialListener(3, &executionLog, false)

	d.Register(eventName, listener1)
	d.Register(eventName, listener2)
	d.Register(eventName, listener3)

	// Use the BaseEvent implementation which is Propagatable
	event := NewBaseEvent(eventName)

	// Dispatch the event
	assert.NoError(t, d.Dispatch(context.Background(), event))

	// --- Verification ---

	// 1. Check execution order (should be sequential: 1, 2)
	expectedLog := []int{1, 2}
	assert.Equal(t, expectedLog, executionLog, "Listeners should be executed sequentially until propagation is stopped.")

	// 2. Check that listener 1 and 2 were called
	select {
	case <-listener1.handledEvents:
		// OK
	case <-time.After(10 * time.Millisecond):
		t.Fatal("Listener 1 was not called")
	}

	select {
	case <-listener2.handledEvents:
		// OK
	case <-time.After(10 * time.Millisecond):
		t.Fatal("Listener 2 was not called")
	}

	// 3. Check that listener 3 was NOT called
	select {
	case <-listener3.handledEvents:
		t.Fatal("Listener 3 should not have been called because propagation was stopped by Listener 2")
	default:
		// All good
	}

	// 4. Check that the event's propagation status is stopped
	assert.True(t, event.IsPropagationStopped(), "Event propagation status should be stopped")
}
