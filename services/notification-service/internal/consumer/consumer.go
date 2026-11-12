// Package consumer connects to RabbitMQ and processes events from a queue.
package consumer

import (
	"context"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer reads messages from a single RabbitMQ queue.
type Consumer struct {
	url    string
	queue  string
	logger *slog.Logger // "*" means a pointer: we share the same logger, not a copy
}

// New is the conventional Go "constructor": a plain function returning a value.
func New(url, queue string, logger *slog.Logger) *Consumer {
	return &Consumer{url: url, queue: queue, logger: logger}
}

// Run connects, consumes messages and blocks until ctx is cancelled
// (e.g. Ctrl+C) or the connection is lost.
//
// Go has no exceptions: functions return an "error" value that the caller
// must check. Returning nil means "no error".
func (c *Consumer) Run(ctx context.Context) error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		// %w wraps the original error so callers can still inspect it.
		return fmt.Errorf("connect to rabbitmq: %w", err)
	}
	// defer runs when Run returns, no matter how: perfect for cleanup.
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %w", err)
	}
	defer ch.Close()

	// Prefetch 1: RabbitMQ sends one message at a time and waits for its ack
	// before sending the next one.
	if err := ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("set qos: %w", err)
	}

	msgs, err := ch.Consume(
		c.queue,
		"notification-service", // consumer tag (shows up in the management UI)
		false,                  // autoAck = false: we ack manually after processing
		false,                  // exclusive
		false,                  // noLocal (not supported by RabbitMQ)
		false,                  // noWait
		nil,                    // extra arguments
	)
	if err != nil {
		return fmt.Errorf("consume queue %q: %w", c.queue, err)
	}

	c.logger.Info("waiting for messages", "queue", c.queue)

	// A for + select loop waits on several channels at once.
	// Go "channels" (chan) are pipes used to pass values between goroutines.
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

// handle processes a single message and acknowledges it.
func (c *Consumer) handle(msg amqp.Delivery) {
	// For now "sending a notification" is just logging the event.
	c.logger.Info("event received",
		"routing_key", msg.RoutingKey,
		"body", string(msg.Body),
	)

	// Ack only after processing: if the service crashes before this line,
	// RabbitMQ puts the message back in the queue.
	if err := msg.Ack(false); err != nil {
		c.logger.Error("ack failed", "error", err)
	}
}
