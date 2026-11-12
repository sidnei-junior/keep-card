// Package consumer connects to RabbitMQ and processes keep-card events from a queue.
package consumer

import (
	"context"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	consumerTag   = "notification-service"
	prefetchCount = 1
	autoAck       = false
	exclusive     = false
	noLocal       = false
	noWait        = false
)

// Consumer reads messages from a single RabbitMQ queue.
type Consumer struct {
	url    string
	queue  string
	logger *slog.Logger
}

// New returns a Consumer for the given RabbitMQ URL and queue.
func New(url, queue string, logger *slog.Logger) *Consumer {
	return &Consumer{url: url, queue: queue, logger: logger}
}

// Run connects to RabbitMQ and consumes messages until ctx is cancelled or
// the connection is lost.
//
// Messages are delivered one at a time and acknowledged manually only after
// being processed, so if the service stops mid-processing RabbitMQ requeues
// the message instead of losing it.
func (c *Consumer) Run(ctx context.Context) error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("connect to rabbitmq: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %w", err)
	}
	defer ch.Close()

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
				return fmt.Errorf("delivery channel closed by rabbitmq")
			}
			c.handle(msg)
		}
	}
}

func (c *Consumer) handle(msg amqp.Delivery) {
	c.logger.Info("event received",
		"routing_key", msg.RoutingKey,
		"body", string(msg.Body),
	)

	if err := msg.Ack(false); err != nil {
		c.logger.Error("ack failed", "error", err)
	}
}
