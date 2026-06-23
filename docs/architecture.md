# CatNet CLI Architecture

This document describes the design decisions and architecture of the `catnet` scriptable CLI.

## Objective
The `catnet` repository acts as a thin, scriptable wrapper around the `catnet-core` engine. It is designed to be composable, adhering to Unix philosophies by providing structured outputs (JSON) and clean human-readable text via `stdout`, while keeping operational progress and logs strictly in `stderr`.

## Key Components

### 1. Root & Commands (`internal/cli/root.go`)
- Uses `github.com/spf13/cobra` for robust CLI scaffolding.
- Exposes persistent flags such as `--no-color`.
- Silences default Cobra error printing to allow manual handling and specific Unix exit codes.
- Format (`--format`) is registered per-subcommand (not persistent) to avoid flag inheritance issues.

### 2. Error Handling (`internal/cli/errors.go`)
- Custom `ExitError` type allows nested commands to bubble up deterministic exit codes (0 = success, 1 = input error, 2 = runtime error, 130 = interrupted).
- `main.go` unwraps these errors and calls `os.Exit(code)` appropriately.
- JSON format validation uses a switch/case to reject invalid formats with a clear error message.

### 3. Graceful Cancellation (`internal/cli/signals.go`)
- Uses `signal.NotifyContext` to listen for `SIGINT` and `SIGTERM`.
- `SIGTERM` is only registered on Unix systems (excluded on Windows where it is not supported).
- Passes the cancelable context directly to `catnet-core`, ensuring that network operations halt cleanly without resource leakage.
- Cancellation produces exit code 130.

### 4. Scanner Command (`internal/cli/scan.go`)
- Responsible for argument parsing (comma-separated IPs, CIDRs, and ranges) via `catnet-core`'s `ParseRange`.
- Subscribes to the `engine.StartScan` Event Stream and dispatches events to the configured output adapter.
- Default ports scanned: `22, 80, 443, 139, 445, 3389` — configurable via `--ports / -p`.
- Supports `--no-ports` to skip port scanning entirely.
- JSON output is written to file with `0600` permissions (owner-only) for security.
- `--quiet` suppresses human-readable progress output on stderr.

### 5. Export Command (`internal/cli/export.go`)
- Re-exports a previously saved JSON scan result to CSV, XML, or JSON format.
- Reads the JSON file, validates it, and converts using `catnet-core`'s encoder.
- Schema version validation: warns if major version is outside supported range (1-2).
- Output file written with `0600` permissions.
- Errors handled: invalid input file, invalid JSON, unsupported format.

### 6. Version Command (`internal/cli/version.go`)
- Prints version, commit, and build date from build-time injected variables.
- Supports `--short` flag for concise output (just the version number).
- Also prints the linked `catnet-core` dependency version.

### 7. Output Adapters (`internal/cli/output/`)
- **Human Output:** Uses `text/tabwriter` for column alignment. Emits dynamic progress bars using `\r` on `stderr` to avoid corrupting data streams. Detects TTY presence dynamically.
  - Constructor uses dependency injection: `NewHumanOutput` accepts `io.Writer` for stdout and stderr (improves testability). Internal helper `newHumanOutputWithWriters` is used by tests.
- **JSON Output:** Silences real-time results and outputs the entire scan artifact to stdout or a defined file using the core exporter.
  - `JSONOutput` struct accepts `io.Writer` for dependency injection, enabling testability without file I/O.

### 8. Event Stream Architecture
- `scan.go` subscribes to `catnet-core`'s event stream via `HandleEvent`.
- Event types dispatched:
  - `Lifecycle` (scan start/end)
  - `Progress` (hosts scanned so far)
  - `Result` (device found with its details)
  - `Warning` (non-fatal errors during scan)
- Human output renders progress live on stderr; JSON output accumulates results and prints at the end.

### 9. Integration Tests (`tests/integration_test.go`)
- `TestMain` compiles the CLI binary via `go build` before running tests.
- All tests use `os/exec` subprocess isolation (no in-process `cli.Execute()` or global state mutation).
- Mock TCP server on `127.0.0.1:0` replaces dependency on real open ports for scan tests.
- Tests cover: JSON output, cancellation via context, invalid targets, XML/CSV/JSON export, I/O errors (invalid file, invalid JSON, unsupported format), schema version warnings, version output.

## CI/CD and Release
- GitHub Actions matrix test against Windows, macOS, and Linux (Go 1.26.x).
- Aikido Security SAST scan runs on pull requests with severity threshold HIGH.
- Integrates `GoReleaser` via `.goreleaser.yml` to automatically cross-compile and publish binaries to GitHub Releases upon tag creation.
- Auto-merge workflow: pushes to `develop` are signed via SSH `git filter-branch` on `develop-signed`, then a PR is opened to `main`.
- Sync workflow: after any PR merges to `main`, `develop` is fast-forwarded automatically using `--force-with-lease`.
