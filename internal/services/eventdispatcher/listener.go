package eventdispatcher

import "context"

// Listener defines the interface for any event listener.
// The Handle method is called when an event it subscribes to is dispatched.
type Listener interface {
	Handle(ctx context.Context, event Event) error
}
