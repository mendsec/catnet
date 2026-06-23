---
layout: default
title: Architecture
nav_order: 6
description: How catnet is structured internally and how it relates to catnet-core.
---

# Architecture

## Design Philosophy

1. **catnet is a consumer, not an engine.** All network logic lives in catnet-core. catnet MUST NOT implement scanning logic directly.
2. **stdout is for data, stderr is for humans.**
3. **Exit codes are a contract.**

## Repository Structure

```
cmd/catnet/main.go              — entry point, ExitError unwrap
internal/cli/
  root.go                       — Cobra root, persistent flags
  scan.go                       — scan subcommand
  export.go                     — export subcommand
  version.go                    — version subcommand
  signals.go                    — SIGINT/SIGTERM context wiring
  errors.go                     — ExitError type, exit codes
  output/
    human.go                    — tabwriter + ANSI colour adapter
    json.go                     — silent adapter (stderr only)
tests/
  integration_test.go           — subprocess-based E2E tests
testdata/
  expected_output.json          — canonical JSON fixture
```

## Event Flow

```
User → catnet scan args
        │
        ▼
  targets.ParseRange (catnet-core)
        │
        ▼
  engine.StartScan (catnet-core)
        │
        ▼
  EventCallback → output adapter → stdout/stderr
        │
        ▼
  exporter.ExportJSON (catnet-core) → file or stdout
```

## Dependency Graph

```
catnet
  ├── catnet-core v0.2.0 (zero external deps)
  ├── github.com/spf13/cobra v1.10.2
  └── github.com/spf13/pflag v1.0.9 (indirect)
```

## Signal Handling

`signals.go` uses `signal.NotifyContext` to wire SIGINT (all platforms) and SIGTERM (Unix) to a context cancellation. When the context is cancelled, the engine stops early and returns a partial report. The process exits with code 130.

## Adding a New Output Format

1. Implement the `OutputHandler` interface in `internal/cli/output/`
2. Register in `scan.go` format switch
3. Add integration test case
4. Document in the output reference
