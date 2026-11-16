package consumer

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/sidnei-junior/keep-card/services/notification-service/internal/handler"
)

func TestRunFailsWhenBrokerIsUnreachable(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewLogHandler(logger)
	c := New("amqp://guest:guest@127.0.0.1:1/", "notification.all", h, logger)

	if err := c.Run(context.Background()); err == nil {
		t.Fatal("Run() error = nil, want connection error")
	}
}
