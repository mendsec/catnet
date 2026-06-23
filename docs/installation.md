---
layout: default
title: Installation
nav_order: 2
description: Install catnet — binary download, go install, or build from source.
---

# Installation

## Requirements

- **Go 1.26 or later** (for source builds only)
- **Linux, macOS, or Windows** (amd64 or arm64)
- **Root / administrator privileges** for ICMP scanning on some systems

## Method 1 — Pre-built Binary (Recommended)

**Linux / macOS:**
```bash
curl -sSL https://github.com/mendsec/catnet/releases/latest/download/catnet_Linux_x86_64.tar.gz | tar xz
sudo mv catnet /usr/local/bin/
catnet version
```

**Windows:**
1. Download `catnet_Windows_x86_64.zip` from the [Releases page](https://github.com/mendsec/catnet/releases).
2. Extract the archive.
3. Add the extracted folder to your `PATH`.
4. Verify: `catnet.exe version`

Always verify your download against the `checksums.txt` file published with each release.

## Method 2 — `go install`

```bash
go install github.com/mendsec/catnet/cmd/catnet@latest
```

For reproducible builds, pin to a specific tag:

```bash
go install github.com/mendsec/catnet/cmd/catnet@v0.2.0
```

## Method 3 — Build from Source

```bash
git clone https://github.com/mendsec/catnet.git
cd catnet
go build -o catnet ./cmd/catnet
./catnet version
```

## Advanced Options

For build-from-source with version injection, Homebrew setup, and shell completion installation, see the [Installation wiki page](https://github.com/mendsec/catnet/wiki/Installation).
