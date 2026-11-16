// Package consumer connects to RabbitMQ and delivers keep-card events from a
// queue to a handler.
package consumer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/sidnei-junior/keep-card/services/notification-service/internal/handler"
)

const (
	consumerTag      = "notification-service"
	prefetchCount    = 1
	autoAck          = false
	exclusive        = false
	noLocal          = false
	noWait           = false
	multiple         = false
	requeueOnFailure = false
)

// Consumer reads messages from a single RabbitMQ queue and passes them to a
// handler.
type Consumer struct {
	url     string
	queue   string
	handler handler.Handler
	logger  *slog.Logger
}

// New returns a Consumer that delivers messages from queue to h.
func New(url, queue string, h handler.Handler, logger *slog.Logger) *Consumer {
	return &Consumer{url: url, queue: queue, handler: h, logger: logger}
}

// Run connects to RabbitMQ and consumes messages until ctx is canceled or
// the connection is lost.
//
// Messages are delivered one at a time and acknowledged manually only after
// the handler succeeds, so if the service stops mid-processing RabbitMQ
// requeues the message instead of losing it. Messages the handler fails to
// process are rejected without requeue, so a permanently invalid message
// cannot block the queue in an endless redelivery loop.
func (c *Consumer) Run(ctx context.Context) error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("connect to rabbitmq: %w", err)
	}
	defer c.closeResource("connection", conn)

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %w", err)
	}
	defer c.closeResource("channel", ch)

	if err := ch.Qos(prefetchCount, 0, false); err != nil {
		return fmt.Errorf("set qos: %w", err)
	}

	msgs, err := ch.Consume(c.queue, consumerTag, autoAck, exclusive, noLocal, noWait, nil)
	if err != nil {
		return fmt.Errorf("consume queue %q: %w", c.queue, err)
	}

	c.logger.Info("waiting for messages", "queue", c.queue)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("shutting down consumer")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return errors.New("delivery channel closed by rabbitmq")
			}
			c.process(ctx, msg)
		}
	}
}

func (c *Consumer) process(ctx context.Context, msg amqp.Delivery) {
	event := handler.Event{RoutingKey: msg.RoutingKey, Body: msg.Body}

	if err := c.handler.Handle(ctx, event); err != nil {
		c.logger.Error("handle event failed", "routing_key", msg.RoutingKey, "error", err)
		if err := msg.Nack(multiple, requeueOnFailure); err != nil {
			c.logger.Error("nack failed", "error", err)
		}
		return
	}

	if err := msg.Ack(multiple); err != nil {
		c.logger.Error("ack failed", "error", err)
	}
}

func (c *Consumer) closeResource(name string, closer io.Closer) {
	if err := closer.Close(); err != nil && !errors.Is(err, amqp.ErrClosed) {
		c.logger.Warn("failed to close "+name, "error", err)
	}
}
