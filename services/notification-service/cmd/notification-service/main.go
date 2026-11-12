// Command notification-service consumes keep-card events and sends notifications.
//
// "package main" + "func main()" is the entry point of a Go program.
// By convention, executables live in cmd/<name>/main.go.
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
	// Structured logs in JSON: easy to read by humans and by log tools.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config.Load()

	// The context is cancelled when we receive Ctrl+C (SIGINT) or SIGTERM
	// (sent by Docker on "stop"), allowing a graceful shutdown.
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
