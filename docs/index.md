---
layout: home

hero:
  name: ghset
  text: Declarative GitHub Repository Settings
  tagline: "Describe an existing repo into YAML. Spin up a new repo from that YAML. One command each way."
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
