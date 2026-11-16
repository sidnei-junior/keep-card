package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestLogHandlerHandle(t *testing.T) {
	tests := []struct {
		name  string
		event Event
	}{
		{
			name:  "logs user.created",
			event: Event{RoutingKey: "user.created", Body: []byte(`{"userId":"u-1"}`)},
		},
		{
			name:  "logs purchase.registered",
			event: Event{RoutingKey: "purchase.registered", Body: []byte(`{"purchaseId":"p-1"}`)},
		},
		{
			name:  "logs empty body",
			event: Event{RoutingKey: "stock.low"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := NewLogHandler(slog.New(slog.NewJSONHandler(&buf, nil)))

			if err := h.Handle(context.Background(), tt.event); err != nil {
				t.Fatalf("Handle() error = %v, want nil", err)
			}

			var entry struct {
				Msg        string `json:"msg"`
				RoutingKey string `json:"routing_key"`
				Body       string `json:"body"`
			}
			if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
				t.Fatalf("log output is not valid JSON: %v (%q)", err, buf.String())
			}

			if entry.Msg != "event received" {
				t.Errorf("msg = %q, want %q", entry.Msg, "event received")
			}
			if entry.RoutingKey != tt.event.RoutingKey {
				t.Errorf("routing_key = %q, want %q", entry.RoutingKey, tt.event.RoutingKey)
			}
			if entry.Body != string(tt.event.Body) {
				t.Errorf("body = %q, want %q", entry.Body, string(tt.event.Body))
			}
		})
	}
}
