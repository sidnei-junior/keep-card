# ADR 0001 — Language for notification-service

- **Status:** Accepted
- **Date:** 2026-11-12

## Context

`notification-service` is the first service of the keep-card monorepo. It only consumes events from RabbitMQ (`notification.all` queue) and sends notifications through an external provider. It has no HTTP API and no business rules beyond formatting messages, so its workload is almost entirely I/O.

The goal of the Foundation milestone is to set up the base tooling (lint, tests, git hooks, CI) on a small service. The project targets free hosting tiers, which usually offer around 256 MB of RAM per instance.

Each service may use a different language; this decision applies only to `notification-service`.

## Options considered

| Criteria | TypeScript (Node) | Go | Python |
| --- | --- | --- | --- |
| RabbitMQ client | `amqplib` (mature) | `amqp091-go` (official, maintained by the RabbitMQ team) | `aio-pika` / `pika` |
| I/O concurrency | Event loop | Goroutines | asyncio |
| Production RAM | ~60–100 MB | ~10–20 MB | ~40–70 MB |
| Deploy artifact | Node runtime + `node_modules` | Single static binary | Python runtime + dependencies |
| Lint / test / format | ESLint + Prettier + Vitest | `gofmt` + `go test` built in, plus golangci-lint | Ruff + pytest |
| Email / push ecosystem | Very strong (Nodemailer, SDKs) | Good, more manual | Good (official SDKs) |
| Typing | Static (TS) | Static, strong | Optional (type hints + mypy) |

## Decision

Use **Go** for `notification-service`.

## Consequences

- **Positive:** very low memory footprint fits free hosting tiers; single binary and small Docker image; formatting and testing come with the toolchain, simplifying the lint and test setup; good opportunity to learn Go.
- **Negative:** learning curve, since the team is new to Go; email/push integrations need a bit more manual code than in Node.
- **Follow-up:** other services will get their own ADR when their language is chosen.
