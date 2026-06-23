---
layout: default
title: catnet — Scriptable Network Scanner CLI
nav_order: 1
description: A fast, scriptable network scanner for the command line. Built in Go. Zero engine dependencies. Made for pipelines.
---

> **catnet** — A fast, scriptable network scanner for the command line. Built in Go. Zero engine dependencies. Made for pipelines.

## What catnet does

| 🔍 Discover              | 🔒 Enumerate                | 📤 Export                |
|---------------------------|-----------------------------|--------------------------|
| ICMP ping sweep           | TCP port scanning           | JSON · CSV · XML         |
| ARP resolution            | Reverse DNS lookup          | Pipeline-ready stdout    |

## Install in 30 seconds

**Linux / macOS:**
```bash
curl -sSL https://github.com/mendsec/catnet/releases/latest/download/catnet_Linux_x86_64.tar.gz | tar xz
sudo mv catnet /usr/local/bin/
```

**Windows:**
Download `catnet_Windows_x86_64.zip` from [Releases](https://github.com/mendsec/catnet/releases), extract, and add to PATH.

## Designed for pipelines

```bash
catnet scan 192.168.1.0/24 --format json | jq '.devices[] | select(.isAlive) | {ip, hostname, openPorts}'
```

## Part of the CatNet Ecosystem

| Repository      | Role                      |
|-----------------|---------------------------|
| catnet-core     | Shared scanning engine    |
| catnet          | **Scriptable CLI**        |
| catnet-tui      | Interactive TUI           |
| catnet-scanner  | Desktop GUI               |

## Latest Release

[![Release](https://img.shields.io/github/v/release/mendsec/catnet)](https://github.com/mendsec/catnet/releases)
[![CI](https://github.com/mendsec/catnet/actions/workflows/ci.yml/badge.svg)](https://github.com/mendsec/catnet/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**v0.2.0** — catnet-core v0.2.0 dependency, stable export flag, English standardisation, CI hardening.

- [Full documentation on the Wiki](https://github.com/mendsec/catnet/wiki)
- [GitHub Repository](https://github.com/mendsec/catnet)
- [Report an Issue](https://github.com/mendsec/catnet/issues/new)
