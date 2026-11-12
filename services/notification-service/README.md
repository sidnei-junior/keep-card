# notification-service

Consumes keep-card events from RabbitMQ and sends notifications to customers and merchants.

> **Current stage:** skeleton. "Sending" a notification just logs the event. Real email/push providers come later.

Language choice: [ADR 0001](../../docs/adr/0001-notification-service-language.md) (Go).

## Events consumed

From the `notification.all` queue (bound to the `keepcard.events` topic exchange):

`user.created`, `purchase.registered`, `reward.redeemed`, `stock.low`

Messages are acknowledged **after** processing (manual ack, prefetch 1). If the service stops before acking, RabbitMQ redelivers the message.

## Requirements

- Go 1.27+
- Local RabbitMQ running with the topology created — see [docs/local-setup.md](../../docs/local-setup.md)

## Configuration

| Variable | Default | Description |
| --- | --- | --- |
| `AMQP_URL` | `amqp://guest:guest@localhost:5672/` | RabbitMQ connection string |
| `QUEUE_NAME` | `notification.all` | Queue to consume from |

## Running

From this folder (`services/notification-service`):

```bash
go run ./cmd/notification-service
```

Stop with `Ctrl+C` (graceful shutdown).

To build a binary:

```bash
go build -o bin/notification-service ./cmd/notification-service
```

## Trying it out

1. Start the service.
2. In the RabbitMQ Management UI, publish to the `keepcard.events` exchange with routing key `user.created` and any JSON payload.
3. The service logs something like:

   ```json
   {"level":"INFO","msg":"event received","routing_key":"user.created","body":"{\"userId\":\"u-1\"}"}
   ```

## Project layout

```
cmd/notification-service/   entry point (main.go)
internal/config/            settings from environment variables
internal/consumer/          RabbitMQ connection and message handling
```

Code under `internal/` can only be imported by this module.
