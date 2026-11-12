// Package config loads the service settings from environment variables.
//
// In Go, a "package" is a folder of .go files. Anything inside "internal/"
// can only be imported by code in this same module, which keeps it private.
package config

import "os"

// Config holds all settings the service needs.
//
// A struct is Go's way to group fields (similar to a class without methods).
// Fields starting with an uppercase letter are exported (public).
type Config struct {
	AMQPURL   string // RabbitMQ connection string
	QueueName string // queue this service consumes from
}

// Load reads the configuration from the environment, falling back to
// defaults that match the local docker-compose setup.
func Load() Config {
	return Config{
		AMQPURL:   getEnv("AMQP_URL", "amqp://guest:guest@localhost:5672/"),
		QueueName: getEnv("QUEUE_NAME", "notification.all"),
	}
}

// getEnv returns the value of an environment variable or a fallback.
// Lowercase name = unexported (only visible inside this package).
func getEnv(key, fallback string) string {
	// "if value, ok := ...; ok" declares variables scoped to the if block.
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
