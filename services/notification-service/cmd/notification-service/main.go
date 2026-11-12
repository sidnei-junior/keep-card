// Command notification-service consumes keep-card events from RabbitMQ and
// sends notifications. It shuts down gracefully on SIGINT or SIGTERM.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sidnei-junior/keep-card/services/notification-service/internal/config"
	"github.com/sidnei-junior/keep-card/services/notification-service/internal/consumer"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("starting notification-service", "queue", cfg.QueueName)

	c := consumer.New(cfg.AMQPURL, cfg.QueueName, logger)
	if err := c.Run(ctx); err != nil {
		logger.Error("consumer stopped with error", "error", err)
		os.Exit(1)
	}

	logger.Info("notification-service stopped")
}
