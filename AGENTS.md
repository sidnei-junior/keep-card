# AGENTS.md

Guidelines for AI coding agents working on keep-card. Human contributors follow the same rules — see [CONTRIBUTING.md](./CONTRIBUTING.md) for the full workflow.

## Project

Event-driven monorepo for a loyalty card platform. Services live in `services/<name>` and communicate through RabbitMQ. Architecture: [docs/architecture](./docs/architecture/README.md). Local infrastructure: [docs/local-setup.md](./docs/local-setup.md).

## Code comments

### Go

Go code must contain **only doc comments** (godoc):

- Comment the `package` and exported identifiers (types, functions, methods, constants, variables).
- Start each doc comment with the identifier name, e.g. `// Load reads the configuration...`.
- Doc comments may explain *why* a design decision was made (e.g. manual ack after processing).

Do **not** write:

- Inline comments inside function bodies.
- Tutorial-style comments explaining language features.
- Comments that restate what the code does.

Clear names are the documentation: prefer named constants and well-named variables over a comment, e.g. `ch.Consume(queue, consumerTag, autoAck, exclusive, noLocal, noWait, nil)` instead of `false, // autoAck`.

## Go testing

- Use the standard library only (`testing`); no assertion libraries such as testify.
- Write table-driven tests with `t.Run` subtests whenever there is more than one scenario.
- Use `t.Setenv`, `t.TempDir` and similar helpers so tests clean up after themselves.
- Keep business logic out of infrastructure code so it can be unit tested (e.g. an event handler interface separated from the RabbitMQ consumer).
- Tests that need real infrastructure (RabbitMQ via testcontainers-go) go in files tagged `//go:build integration`.
- Each service's quality gates (linters, coverage minimum) are listed in its README.

## Conventions

- **Language:** code, commits, branch names and checklists in English; issue and PR descriptions in English and Brazilian Portuguese (EN + PT-BR).
- **Commits and PR titles:** [Conventional Commits](https://www.conventionalcommits.org/), imperative mood, optional free-form scope (usually the service name).
- **Branches:** `<type>/<short-description>`, without the issue number.
- **Merging:** `main` is protected; every change goes through a PR and is squash-merged.
- **Git hooks:** Lefthook runs lint, commit message checks and unit tests (see [CONTRIBUTING.md](./CONTRIBUTING.md#git-hooks)). Never bypass them with `--no-verify`; fix the reported problem instead.
- **Formatting:** follow [`.editorconfig`](./.editorconfig) (LF, UTF-8; tabs for Go, 4 spaces for Python, 2 spaces otherwise). Go code must be `gofmt`-clean.
- **Decisions:** record significant architecture or technology decisions as ADRs in `docs/adr/`.
