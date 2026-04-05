---
title: "ghset Devlog #1 — The Graveyard of GitHub Settings Tools"
description: "Every tool for declarative GitHub repo settings is either dead, broken, or requires a PhD in Terraform. We built ghset because nobody else did it right."
head:
  - - meta
    - name: keywords
      content: github repository settings cli, declarative github config, github repo settings yaml, copy github repo settings, github settings as code, ghset, probot settings deprecated, safe-settings broken, gitstrap broken
---

# #1 — The Graveyard of GitHub Settings Tools

> "I just want to copy settings from one repo to another. Why is this so hard?"
> — Every developer, eventually.

## The Trigger

I needed to create a new Go repository — `clienter-go`. Private, same settings as my open-source repos: rebase-only merges, delete branch on merge, secret scanning, custom labels, branch protection rulesets. The usual.

I opened `thumbrise/resilience` in the browser. Settings tab. Clicked through 6 pages of checkboxes. Then Security. Then Labels — 30+ labels to recreate. Then Rulesets — two rulesets with bypass actors, required reviews, status checks.

Twenty minutes of clicking. For settings I've configured before. On a repo template that GitHub promises will "help you get started quickly" — but copies files, not settings.

I thought: surely someone has built a tool for this.

## The Research

I spent an evening looking. What I found was a graveyard.

### probot/settings — Dead (deprecated framework)

The original idea was beautiful: store settings in `.github/settings.yml`, apply on push. A GitHub App powered by the Probot framework.

**Status: DEAD.** The Probot framework itself is deprecated, unmaintained, and potentially vulnerable. The settings app hasn't been updated in years. No rulesets support. No security settings. The ecosystem it was built on no longer exists.

### github/safe-settings — Dead (GitHub's own, abandoned)

GitHub's official replacement for probot/settings. Policy-as-code for organizations. Sounds perfect.

**Status: DEAD.** The repository is a graveyard of unfixed bugs:

- **Issue #744** — when no changes are needed for a ruleset, the app crashes with `Cannot read properties of undefined` and silently skips all remaining settings.
- **Issue #655** — manually deleted rules are not detected. The tool doesn't notice.
- **Issue #723** — config merging produces duplicates.
- **Issue #540** — check runs are created but never execute.

A tool from GitHub that can't reliably manage GitHub settings.

### gitstrap — Dead (broken installation)

A Go CLI for bootstrapping repos from YAML. Exactly what I needed.

**Status: BROKEN.** The installation instructions from the official README don't work. `brew install` fails with `formula requires at least a URL`. The curl install script fails with `tar: Error opening archive: Unrecognized archive format`. Last commit: 2021.

I couldn't even install it, let alone use it.

### dothub — Dead (2017)

A Python CLI for repo settings. 3 stars.

**Status: DEAD.** Last commit: 2017. Issues describe fundamental problems: the tool leaks internal IDs into configs that can't be used for creation (#47), and it "guesses" which Git repo you mean (#48). Dead for almost a decade.

### google/github-repo-automation — Dead (officially deprecated)

A Node.js tool from Google. The name alone inspires confidence.

**Status: DEPRECATED.** The repository page says it in bold: **THIS REPOSITORY IS DEPRECATED.** Google built it, used it internally, and abandoned it. No settings management, no labels, no rulesets.

### Terraform GitHub Provider — Alive, but overkill

The industrial-grade option. Full GitHub API coverage. Actively maintained.

**Status: ALIVE, but WRONG TOOL.** It requires:

- HCL configuration language
- State file management (where do you store it? S3? Local? Git?)
- `terraform init`, `terraform plan`, `terraform apply` workflow
- Understanding of Terraform lifecycle, imports, drift detection

To copy settings from one repo to another, you'd need to:

1. Write HCL for every setting
2. Import the existing repo into state
3. Modify the HCL for the new repo
4. Apply

That's not "copy settings." That's "become a Terraform engineer."

And it has its own bugs: import errors on branch protection rules (#2122), issues with empty check settings (#2261).

### xfg — Alive, but bloated

The only tool that's actually maintained (last update: February 2026). Supports settings, rulesets, file sync.

**Status: ALIVE, but BLOATED.** xfg is a multi-repo file synchronization tool that also happens to manage settings. It's designed for organizations managing hundreds of repos with shared configs. For "copy settings from A to B" it's like using a forklift to move a chair.

## What People Actually Want

The evidence is everywhere:

On **GitHub Discussions**, users complain:
> "I assumed branch protection rules and access settings would be copied to the new repository when creating from a template. They are not."

On **Stack Overflow**, the same question appears repeatedly:
> "Is there a way to copy GitHub repository settings, including webhooks, access, etc., to another repository in the same organization?"

Developers resort to **PowerShell scripts** and **curl one-liners** to solve this. Everyone reinvents the same wheel because no tool does it simply.

## The Gap

| Tool | Status | Settings | Security | Labels | Rulesets | Describe |
|---|---|---|---|---|---|---|
| probot/settings | DEAD | ✅ | ❌ | ✅ | ❌ | ❌ |
| safe-settings | DEAD | ✅ | ✅ | ✅ | buggy | ❌ |
| gitstrap | DEAD | ✅ | ❌ | ✅ | ✅ | ✅ |
| dothub | DEAD | ✅ | ❌ | ✅ | ❌ | ✅ |
| google/repo-auto | DEAD | ✅ | ❌ | ✅ | ❌ | ❌ |
| Terraform | ALIVE | ✅ | ✅ | ✅ | ✅ | ❌ |
| xfg | ALIVE | ✅ | ❌ | ✅ | ✅ | ❌ |
| **ghset** | **ALIVE** | **✅** | **✅** | **✅** | **✅** | **✅** |

The `Describe` column tells the story. Every tool assumes you'll write the config by hand. None of them can snapshot an existing repo. The dead ones (gitstrap, dothub) tried, but they're dead.

ghset is the only living tool that can describe an existing repo into YAML and create a new repo from that YAML. One command each way.

## What We Built

Two commands. One YAML format. Zero infrastructure.

```bash
# Snapshot everything
ghset describe thumbrise/resilience > template.yml

# Create a new repo with the same settings
ghset init clienter-go --from template.yml

# Or skip the file entirely
ghset init clienter-go --from thumbrise/resilience
```

Settings, security, labels, rulesets — all captured, all applied. Best-effort: if secret scanning isn't available on your plan, ghset warns and continues. If a label already exists, it skips. No state files, no HCL, no YAML schema to memorize.

The tool delegates authentication to `gh` CLI — the tool every developer already has installed. No tokens to manage, no OAuth flows, no GitHub App to register.

Twenty minutes of clicking became one command.

## The Lesson

The graveyard exists because everyone tried to build a platform. Probot built a framework. safe-settings built policy-as-code. Terraform built infrastructure management. Google built... something, then deprecated it.

Nobody built a tool that does one thing: **copy settings from repo A to repo B.**

ghset does one thing. It does it now. It works today.

---

*The tool that exists because every other tool doesn't.*
