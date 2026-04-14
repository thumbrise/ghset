# Contributing

## Requirements

- **Go 1.26.x** — exact major.minor match required.
- **[golangci-lint](https://golangci-lint.run/)** — install: `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`
- **[Task](https://taskfile.dev/)** — task runner. Install: `go install github.com/go-task/task/v3/cmd/task@latest`
- **[gh CLI](https://cli.github.com/)** — required at runtime. Install and authenticate: `gh auth login`.
- **Node.js** (optional) — only for commitlint and docs build.

## First time setup

```bash
git clone https://github.com/thumbrise/ghset.git
cd ghset
go build ./...
```

## Development workflow

```bash
# Run tests
go test ./... -v

# Run lint (Docker, no local golangci-lint needed)
task lint

# Fix license headers
task generate

# Build and run
go run . describe thumbrise/ghset
go run . init my-repo --from config.yml
```

## Commit messages

Conventional commits. English only.

```
feat(describe): add rulesets support
fix(init): handle private repo visibility
docs: update getting-started guide
```

Only `feat` and `fix` trigger releases. See [REVIEW.md](REVIEW.md) for full guidelines.

## Tests

- All tests use `package xxx_test` (blackbox only).
- Bug fix = test first: red test commit, then fix commit. Never combined.
- Run `task test` before pushing.

## Code review

See [REVIEW.md](REVIEW.md) — structural review rules, naming conventions, hard thresholds.
