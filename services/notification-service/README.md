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
- [golangci-lint](https://golangci-lint.run/) v2 — Windows: `winget install GolangCI.golangci-lint`
- [Task](https://taskfile.dev/) v3 — Windows: `winget install Task.Task`
- Local RabbitMQ running with the topology created — see [docs/local-setup.md](../../docs/local-setup.md)

## Tasks

Common commands are defined in [`Taskfile.yml`](./Taskfile.yml). Run `task --list` to see them all.

| Command | Description |
| --- | --- |
| `task lint` | Run golangci-lint (linters + formatting check) |
| `task lint:fix` | Run golangci-lint and apply automatic fixes |
| `task fmt` | Format code with goimports |

Linter configuration lives in [`.golangci.yml`](./.golangci.yml).

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

## Quality gates

Rules every change to this service must meet. Implementation is tracked in the lint, tests, hooks and CI issues of the Foundation milestone.

| Gate | Rule |
| --- | --- |
| Formatting | `goimports` (run through golangci-lint) |
| Linting | golangci-lint with **zero issues** |
| Tests | `go test` with the standard library only (no assertion libraries) |
| Test style | Table-driven tests with `t.Run` subtests when there is more than one scenario |
| Race detector | `-race` enabled in CI (Linux); optional locally |
| Coverage | **≥ 80%** over `internal/` (`cmd/` excluded), counting unit + integration tests, enforced in CI |

### Linters

Default set: `govet`, `errcheck`, `staticcheck`, `unused`, `ineffassign`.

Extra:

| Linter | Why |
| --- | --- |
| `revive` (`exported` rule) | Requires doc comments on exported identifiers |
| `gosec` | Security issues |
| `errorlint` | Correct error wrapping (`%w`, `errors.Is` / `errors.As`) |
| `misspell` | English typos |
| `gocritic` | Style and performance suggestions |

### Test types

| Type | Scope | Command | Needs Docker |
| --- | --- | --- | --- |
| Unit | Pure logic (e.g. event handler, config) | `go test ./...` | No |
| Integration | RabbitMQ consumer against a real broker via [testcontainers-go](https://golang.testcontainers.org/), files tagged `//go:build integration` | `go test -tags integration ./...` | Yes |

To keep the RabbitMQ code testable, the **consumer** (connection, delivery, ack/nack) is separated from the **handler** (what to do with an event), which is an interface tested with plain unit tests.
