# Contributing to keep-card

## Workflow

1. Every change starts from a **GitHub issue**.
2. Create a branch from `main` (see naming below).
3. Commit following **Conventional Commits**.
4. Open a **Pull Request** to `main` using the PR template and reference the issue with `Closes #<number>`.
5. The PR is merged with **squash merge**: the whole PR becomes a single commit on `main`, and the PR title becomes its message. The branch is deleted automatically after merge.

`main` is protected: no direct pushes, no force pushes, and every change goes through a PR (rules apply to admins too).

## Branch naming

```
<type>/<short-description>
```

Use the same types as commits, in kebab-case, without the issue number.

Examples: `docs/base-architecture`, `feat/purchase-registration`, `chore/general-quality-gates`.

## Commit messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<optional scope>): <description>
```

- **Language:** English, imperative mood (`add`, not `added` / `adds`).
- **Scope:** optional and free-form; usually the service name (e.g. `card-service`).
- **PR titles** follow the same format, since they become the squash commit message.

| Type | When to use |
| --- | --- |
| `feat` | A new feature |
| `fix` | A bug fix |
| `docs` | Documentation only |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `test` | Adding or fixing tests |
| `perf` | Performance improvement |
| `build` | Build system or dependencies |
| `ci` | CI configuration (GitHub Actions) |
| `chore` | Maintenance tasks that don't touch app code |

Examples:

```
feat(card-service): add purchase registration endpoint
fix(identity-service): reject expired tokens
docs: add base architecture diagram
```

## Language

- **Code, commits, branch names and checklists:** English.
- **Issue and PR descriptions:** English and Brazilian Portuguese (EN + PT-BR).

## Code style

Formatting is defined in [`.editorconfig`](./.editorconfig) and line endings in [`.gitattributes`](./.gitattributes):

- UTF-8, LF line endings, final newline, no trailing whitespace.
- 2 spaces by default (TypeScript, JSON, YAML), tabs for Go, 4 spaces for Python.

VS Code users: install the recommended extensions from [`.vscode/extensions.json`](./.vscode/extensions.json).

## Local setup

See [docs/architecture](./docs/architecture/README.md) for an overview of the system and [docs/local-setup.md](./docs/local-setup.md) to run the local infrastructure (Docker, RabbitMQ). More setup steps (git hooks, services) will be added here as the project grows.
