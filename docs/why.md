---
title: "Why ghset? — GitHub Repository Settings CLI"
description: "Every tool for declarative GitHub repo settings is dead, broken, or requires Terraform. ghset is the simple alternative: describe and init in one command."
head:
  - - meta
    - name: keywords
      content: github repository settings cli, declarative github config, copy github repo settings, github settings as code, alternative to terraform github, probot settings alternative, safe-settings alternative
---

# Why ghset?

GitHub template repos copy files. Not settings. Not labels. Not rulesets. Not security config.

Every developer who creates repos regularly hits the same wall: 20 minutes of clicking through Settings, Security, Labels, and Rulesets — for configuration they've already done on another repo.

## What exists today

We researched every tool in the ecosystem. The results:

| Tool | Status | Problem |
|---|---|---|
| **probot/settings** | DEAD | Framework deprecated, unmaintained, vulnerable |
| **github/safe-settings** | DEAD | Crashes on rulesets (#744), ignores deleted rules (#655) |
| **gitstrap** | DEAD | Installation broken since 2021 — brew and curl both fail |
| **dothub** | DEAD | Last commit 2017, leaks internal IDs into configs |
| **google/github-repo-automation** | DEAD | Officially DEPRECATED by Google |
| **Terraform GitHub provider** | ALIVE | Requires HCL, state management, import workflow — overkill |
| **xfg** | ALIVE | Multi-repo file sync tool, not a settings copier |

Full analysis with issue numbers and user quotes: [Devlog #1 — The Graveyard](/devlog/001-the-graveyard).

## What's missing everywhere

No living tool can do this:

```bash
# Snapshot an existing repo's settings
ghset describe thumbrise/resilience > template.yml

# Create a new repo with those exact settings
ghset init my-new-repo --from template.yml
```

The `describe` direction — snapshotting a repo into a portable config — doesn't exist in any maintained tool. Everyone assumes you'll write YAML by hand.

## What ghset does

| Feature | ghset | Terraform | Template repos |
|---|---|---|---|
| Copy settings | ✅ one command | ✅ HCL + state | ❌ |
| Copy security | ✅ | ✅ | ❌ |
| Copy labels | ✅ | ✅ | ❌ |
| Copy rulesets | ✅ | ✅ complex | ❌ |
| Snapshot existing repo | ✅ `describe` | ❌ | ❌ |
| Setup time | 0 (uses `gh`) | Hours (HCL + state) | 0 |
| Infrastructure | None | State file + backend | None |
| Learning curve | 2 commands | Terraform language | Click UI |

## Install

```bash
go install github.com/thumbrise/ghset@latest
```

Only prerequisite: [gh CLI](https://cli.github.com/) installed and authenticated.

## Usage

```bash
# Describe → YAML
ghset describe owner/repo > settings.yml

# Init from file
ghset init new-repo --from settings.yml

# Init directly from another repo
ghset init new-repo --from owner/repo

# Pipe
ghset describe owner/repo | ghset init new-repo
```

No tokens. No state files. No infrastructure. Just settings, copied.
