// Package config loads the notification-service settings from environment variables.
package config

import "os"

// Config holds the settings needed to run the service.
type Config struct {
	AMQPURL   string
	QueueName string
}

// Load reads the configuration from the environment, falling back to
// defaults that match the local docker-compose setup.
func Load() Config {
	return Config{
		AMQPURL:   getEnv("AMQP_URL", "amqp://guest:guest@localhost:5672/"),
		QueueName: getEnv("QUEUE_NAME", "notification.all"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
