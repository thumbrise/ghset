# ghset

[![CI](https://github.com/thumbrise/ghset/actions/workflows/ci.yml/badge.svg)](https://github.com/thumbrise/ghset/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/thumbrise/ghset.svg)](https://pkg.go.dev/github.com/thumbrise/ghset)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](/LICENSE)

Declarative GitHub repository settings. Describe an existing repo into YAML, spin up a new repo from that YAML.

## Install

```bash
go install github.com/thumbrise/ghset@latest
```

Prerequisites: [gh CLI](https://cli.github.com/) installed and authenticated (`gh auth login`).

## Usage

```bash
# Snapshot an existing repo's settings
ghset describe thumbrise/ghset > opensource.yml

# Create a new repo with those settings
ghset init my-new-repo --from opensource.yml

# Or directly from another repo — no intermediate file
ghset init my-new-repo --from thumbrise/ghset

# Apply settings to an existing repo
ghset apply thumbrise/ghset --from opensource.yml

# Pipe — describe one, create another
ghset describe thumbrise/ghset | ghset init my-new-repo
```

## What it manages

- **Settings** — visibility, merge strategies, wiki, issues, projects, discussions
- **Security** — secret scanning, push protection, vulnerability alerts, automated fixes
- **Labels** — create custom, skip existing
- **Rulesets** — branch and tag protection rules

## Config format

Single YAML format — `describe` writes it, `init` reads it:

```yaml
settings:
  visibility: public
  allow_rebase_merge: true
  delete_branch_on_merge: true

security:
  secret_scanning: true
  vulnerability_alerts: true

labels:
  - name: "T: bug"
    color: "d73a4a"
    description: "Something isn't working"

rulesets:
  - name: "main"
    target: branch
    enforcement: active
    rules:
      - type: deletion
      - type: pull_request
        parameters:
          required_approving_review_count: 1
```

## How it works

All GitHub API calls go through `gh` CLI as a subprocess. No tokens in the tool, no OAuth — `gh` handles authentication entirely.

## Why not Terraform / Probot / safe-settings?

They're either dead, broken, or overkill. [Full analysis →](https://thumbrise.github.io/ghset/why)

## Documentation

**[thumbrise.github.io/ghset](https://thumbrise.github.io/ghset/)** — why ghset, devlog, full docs.

## License

Apache 2.0
