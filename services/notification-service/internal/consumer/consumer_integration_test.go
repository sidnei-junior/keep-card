//go:build integration

package consumer

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"

	"github.com/sidnei-junior/keep-card/services/notification-service/internal/handler"
)

const (
	rabbitmqImage   = "rabbitmq:4-management"
	eventsExchange  = "test.keepcard.events"
	eventsQueue     = "test.events"
	deadLetterQueue = "test.events.dlq"
	waitTimeout     = 10 * time.Second
	pollInterval    = 100 * time.Millisecond
)

type recordingHandler struct {
	events chan handler.Event
	err    error
}

func newRecordingHandler(err error) *recordingHandler {
	return &recordingHandler{events: make(chan handler.Event, 10), err: err}
}

func (h *recordingHandler) Handle(_ context.Context, event handler.Event) error {
	h.events <- event
	return h.err
}

func TestConsumerIntegration(t *testing.T) {
	ctx := context.Background()

	container, err := rabbitmq.Run(ctx, rabbitmqImage)
	testcontainers.CleanupContainer(t, container)
	if err != nil {
		t.Fatalf("start rabbitmq container: %v", err)
	}

	url, err := container.AmqpURL(ctx)
	if err != nil {
		t.Fatalf("get amqp url: %v", err)
	}

	tests := []struct {
		name           string
		handlerErr     error
		wantDeadLetter int
	}{
		{name: "acks message when handler succeeds", handlerErr: nil, wantDeadLetter: 0},
		{name: "rejects message without requeue when handler fails", handlerErr: errors.New("boom"), wantDeadLetter: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := setupQueues(t, url)

			h := newRecordingHandler(tt.handlerErr)
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			c := New(url, eventsQueue, h, logger)

			runCtx, cancel := context.WithCancel(ctx)
			runErr := make(chan error, 1)
			go func() { runErr <- c.Run(runCtx) }()

			publish(t, ch, "user.created", `{"userId":"u-1"}`)

			select {
			case got := <-h.events:
				if got.RoutingKey != "user.created" || string(got.Body) != `{"userId":"u-1"}` {
					t.Errorf("handler got %+v, want routing key user.created with original body", got)
				}
			case <-time.After(waitTimeout):
				t.Fatal("handler was not called")
			}

			waitForMessageCount(t, ch, eventsQueue, 0)
			waitForMessageCount(t, ch, deadLetterQueue, tt.wantDeadLetter)

			select {
			case extra := <-h.events:
				t.Errorf("message was redelivered: %+v", extra)
			default:
			}

			cancel()
			if err := <-runErr; err != nil {
				t.Errorf("Run() error = %v, want nil after cancel", err)
			}
		})
	}
}

func setupQueues(t *testing.T, url string) *amqp.Channel {
	t.Helper()

	conn, err := amqp.Dial(url)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("open channel: %v", err)
	}

	for _, queue := range []string{eventsQueue, deadLetterQueue} {
		if _, err := ch.QueueDelete(queue, false, false, false); err != nil {
			t.Fatalf("delete queue %s: %v", queue, err)
		}
	}

	if _, err := ch.QueueDeclare(deadLetterQueue, true, false, false, false, nil); err != nil {
		t.Fatalf("declare dead-letter queue: %v", err)
	}

	deadLetterArgs := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": deadLetterQueue,
	}
	if _, err := ch.QueueDeclare(eventsQueue, true, false, false, false, deadLetterArgs); err != nil {
		t.Fatalf("declare events queue: %v", err)
	}

	if err := ch.ExchangeDeclare(eventsExchange, "topic", true, false, false, false, nil); err != nil {
		t.Fatalf("declare exchange: %v", err)
	}

	if err := ch.QueueBind(eventsQueue, "#", eventsExchange, false, nil); err != nil {
		t.Fatalf("bind events queue: %v", err)
	}

	return ch
}

func publish(t *testing.T, ch *amqp.Channel, routingKey, body string) {
	t.Helper()

	msg := amqp.Publishing{ContentType: "application/json", Body: []byte(body)}
	if err := ch.PublishWithContext(context.Background(), eventsExchange, routingKey, false, false, msg); err != nil {
		t.Fatalf("publish: %v", err)
	}
}

func waitForMessageCount(t *testing.T, ch *amqp.Channel, queue string, want int) {
	t.Helper()

	deadline := time.Now().Add(waitTimeout)
	for {
		q, err := ch.QueueDeclarePassive(queue, true, false, false, false, nil)
		if err != nil {
			t.Fatalf("inspect queue %s: %v", queue, err)
		}
		if q.Messages == want {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("queue %s has %d messages, want %d", queue, q.Messages, want)
		}
		time.Sleep(pollInterval)
	}
}
