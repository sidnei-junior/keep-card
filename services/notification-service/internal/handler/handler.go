// Package handler defines how notification-service reacts to keep-card events.
package handler

import (
	"context"
	"log/slog"
)

// Event is a message received from the keep-card event bus.
type Event struct {
	RoutingKey string
	Body       []byte
}

// Handler processes a single event. Returning an error tells the caller that
// the event could not be processed.
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// LogHandler handles events by logging them. It stands in for real
// notification delivery until email and push providers are integrated.
type LogHandler struct {
	logger *slog.Logger
}

// NewLogHandler returns a LogHandler that writes to logger.
func NewLogHandler(logger *slog.Logger) *LogHandler {
	return &LogHandler{logger: logger}
}

// Handle logs the event and never fails.
func (h *LogHandler) Handle(ctx context.Context, event Event) error {
	h.logger.InfoContext(ctx, "event received",
		"routing_key", event.RoutingKey,
		"body", string(event.Body),
	)
	return nil
}
