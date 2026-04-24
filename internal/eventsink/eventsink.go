// Package eventsink provides a fan-out event sink that dispatches
// port-change diff events to one or more registered handlers.
package eventsink

import (
	"fmt"
	"log"
	"sync"

	"github.com/example/portwatch/internal/snapshot"
)

// Handler is a function that receives a diff event.
type Handler func(diff snapshot.Diff) error

// Sink fans out diff events to all registered handlers.
type Sink struct {
	mu       sync.RWMutex
	handlers []namedHandler
	logger   *log.Logger
}

type namedHandler struct {
	name    string
	handle  Handler
}

// New creates a new Sink. logger must not be nil.
func New(logger *log.Logger) *Sink {
	if logger == nil {
		panic("eventsink: logger must not be nil")
	}
	return &Sink{logger: logger}
}

// Register adds a named handler to the sink.
// Duplicate names are allowed; all registered handlers are called.
func (s *Sink) Register(name string, h Handler) {
	if h == nil {
		panic("eventsink: handler must not be nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers = append(s.handlers, namedHandler{name: name, handle: h})
}

// Dispatch sends the diff to every registered handler sequentially.
// Errors from individual handlers are logged but do not abort dispatch.
// Returns a combined error if any handler failed.
func (s *Sink) Dispatch(diff snapshot.Diff) error {
	s.mu.RLock()
	handlers := make([]namedHandler, len(s.handlers))
	copy(handlers, s.handlers)
	s.mu.RUnlock()

	var errs []error
	for _, nh := range handlers {
		if err := nh.handle(diff); err != nil {
			s.logger.Printf("eventsink: handler %q error: %v", nh.name, err)
			errs = append(errs, fmt.Errorf("%s: %w", nh.name, err))
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("eventsink: %d handler(s) failed: %v", len(errs), errs)
}

// Len returns the number of registered handlers.
func (s *Sink) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.handlers)
}
