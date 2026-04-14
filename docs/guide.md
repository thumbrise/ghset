---
title: "Guide — ghset"
description: "Install ghset and copy GitHub repository settings in one command."
head:
  - - meta
    - name: keywords
      content: github repository settings cli, ghset install, ghset guide, copy github repo settings, declarative github config
---

# Guide

## Install

Only prerequisite: [gh CLI](https://cli.github.com/) installed and authenticated (`gh auth login`).

**macOS / Linux**
```bash
curl -sfL https://raw.githubusercontent.com/thumbrise/ghset/main/install.sh | sh
```

**Go**
```bash
go install github.com/thumbrise/ghset@latest
```

**Manual** — grab a binary from [Releases](https://github.com/thumbrise/ghset/releases/latest).

## Quick start

```bash
ghset init my-new-repo --from thumbrise/ghset
```

That's it. The new repo gets everything copied:
- **Settings** — visibility, merge strategies, wiki, issues, projects, discussions
- **Security** — secret scanning, push protection, vulnerability alerts, automated fixes
- **Labels** — create custom, skip existing
- **Rulesets** — branch and tag protection rules

## Commands

```bash
# Snapshot repo settings → YAML
ghset describe owner/repo > config.yml

# Create new repo from config
ghset init my-new-repo --from config.yml

# Or directly from another repo — no intermediate file
ghset init my-new-repo --from owner/repo

# Apply config to an existing repo
ghset apply owner/repo --from config.yml

# Pipe — describe one, create another
ghset describe owner/repo | ghset init my-new-repo
```

| Command | What it does |
|---------|-------------|
| `describe` | Snapshot repo → YAML to stdout |
| `init` | Create new repo + apply settings |
| `apply` | Apply settings to existing repo |

## Config format

Single YAML format — `describe` writes it, `init` and `apply` read it:

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

All GitHub API calls go through `gh` CLI as a subprocess. No tokens in the tool, no OAuth — `gh` handles auth entirely.

No state files. No infrastructure. Just settings, copied.
