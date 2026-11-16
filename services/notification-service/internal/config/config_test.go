package config

import "testing"

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		setEnv   bool
		envValue string
		want     string
	}{
		{name: "uses env value when set", setEnv: true, envValue: "custom", want: "custom"},
		{name: "falls back when unset", setEnv: false, want: "fallback"},
		{name: "falls back when empty", setEnv: true, envValue: "", want: "fallback"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv("KEEPCARD_TEST_KEY", tt.envValue)
			}

			if got := getEnv("KEEPCARD_TEST_KEY", "fallback"); got != tt.want {
				t.Errorf("getEnv() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want Config
	}{
		{
			name: "defaults match local docker-compose",
			env:  map[string]string{"AMQP_URL": "", "QUEUE_NAME": ""},
			want: Config{
				AMQPURL:   "amqp://guest:guest@localhost:5672/",
				QueueName: "notification.all",
			},
		},
		{
			name: "environment overrides defaults",
			env:  map[string]string{"AMQP_URL": "amqp://user:pass@broker:5672/", "QUEUE_NAME": "custom.queue"},
			want: Config{
				AMQPURL:   "amqp://user:pass@broker:5672/",
				QueueName: "custom.queue",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			if got := Load(); got != tt.want {
				t.Errorf("Load() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
