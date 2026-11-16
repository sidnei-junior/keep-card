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
	"github.com/sidnei-junior/keep-card/services/notification-service/internal/handler"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("consumer stopped with error", "error", err)
		os.Exit(1)
	}

	logger.Info("notification-service stopped")
}

func run(logger *slog.Logger) error {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("starting notification-service", "queue", cfg.QueueName)

	h := handler.NewLogHandler(logger)

	return consumer.New(cfg.AMQPURL, cfg.QueueName, h, logger).Run(ctx)
}
