---
layout: home

hero:
  name: ghset
  text: Copy GitHub Repository Settings
  tagline: "ghset init --from owner/repo — copy settings, security, labels, and rulesets. One command."
  actions:
    - theme: brand
      text: Why ghset?
      link: /why
    - theme: alt
      text: Devlog
      link: /devlog/
    - theme: alt
      text: GitHub
      link: https://github.com/thumbrise/ghset

features:
  - icon: 📸
    title: Describe
    details: "Snapshot any repo's settings, security, labels, and rulesets into a single YAML file. One command: ghset describe owner/repo."
  - icon: 🚀
    title: Init
    details: "Create a new repo and apply all settings from YAML — or directly from another repo. ghset init my-repo --from owner/repo."
  - icon: 🔧
    title: Everything That Matters
    details: "Settings, security (secret scanning, vulnerability alerts), labels, rulesets — all in one config. Not just files — actual repo configuration."
  - icon: 🪶
    title: Zero Infrastructure
    details: "No Terraform state. No GitHub App. No tokens. Just gh CLI you already have. go install and go."
  - icon: 🔗
    title: Composable
    details: "Pipe-friendly. ghset describe repo | ghset init new-repo. Works with flags, stdin, or interactive prompts."
  - icon: ⚰️
    title: Because Everything Else Is Dead
    details: "Probot Settings — dead. safe-settings — buggy. gitstrap — broken install. Google's tool — deprecated. We checked."
    link: /devlog/001-the-graveyard
    linkText: Read the graveyard report →
---

## What is ghset?

ghset is an open-source CLI tool that lets you **copy GitHub repository settings** from one repo to another — in one command:

```bash
ghset init my-new-repo --from owner/existing-repo
```

That single line creates `my-new-repo` and applies every setting, security toggle, label, and ruleset from the source. No clicking through the UI, no Terraform, no state files. Just `init --from`.

Need a portable snapshot first? `ghset describe owner/repo > settings.yml` exports everything to YAML. Edit it, commit it, apply it later with `ghset init`. **GitHub repository settings as code** — that's the idea.

## Why another GitHub repository settings tool?

Every existing tool for **declarative GitHub settings** is either dead, broken, or requires heavy infrastructure. GitHub template repos copy files — not settings, labels, rulesets, or security config. ghset fills the gap: a lightweight, pipe-friendly **GitHub repository settings tool** built around `init --from` — the fastest way to clone a repo's configuration.
