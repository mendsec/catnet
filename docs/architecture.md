# CatNet CLI Architecture

This document describes the design decisions and architecture of the `catnet` scriptable CLI.

## Objective
The `catnet` repository acts as a thin, scriptable wrapper around the `catnet-core` engine. It is designed to be composable, adhering to Unix philosophies by providing structured outputs (JSON) and clean human-readable text via `stdout`, while keeping operational progress and logs strictly in `stderr`.

## Key Components

### 1. Root & Commands (`internal/cli/root.go`)
- Uses `github.com/spf13/cobra` for robust CLI scaffolding.
- Exposes persistent flags such as `--format` and `--no-color`.
- Silences default Cobra error printing to allow manual handling and specific Unix exit codes.

### 2. Error Handling (`internal/cli/errors.go`)
- Custom `ExitError` type allows nested commands to bubble up deterministic exit codes.
- `main.go` unwraps these errors and calls `os.Exit(code)` appropriately.

### 3. Graceful Cancellation (`internal/cli/signals.go`)
- Uses `signal.NotifyContext` to listen for `SIGINT` and `SIGTERM`.
- Passes the cancelable context directly to `catnet-core`, ensuring that network operations halt cleanly without resource leakage.

### 4. Scanner Command (`internal/cli/scan.go`)
- Responsible for argument parsing (comma-separated IPs, CIDRs, and ranges).
- Subscribes to the `engine.StartScan` Event Stream.
- Dispatches UI updates through the `output` package adapters.

### 5. Output Adapters (`internal/cli/output/`)
- **Human Output:** Leverages `text/tabwriter` for column alignment. Emits dynamic progress bars using `\r` on `stderr` to avoid corrupting data streams. Detects TTY presence dynamically.
- **JSON Output:** Silences real-time results and outputs the entire scan artifact purely to stdout or a defined file using the core exporter.

## CI/CD and Release
- Github Actions workflows matrix test against Windows, macOS, and Linux.
- Integrates `GoReleaser` via `.goreleaser.yml` to automatically cross-compile and publish binaries to GitHub Releases upon tag creation.
