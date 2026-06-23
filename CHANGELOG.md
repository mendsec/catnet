# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-06-23

### Changed
- `export --format` flag is now local to the export subcommand and required. It no longer inherits the root `--format` value. Existing scripts that relied on `catnet --format csv export input.json` must be updated to `catnet export input.json --format csv`. (Resolves analysis finding C8.)
- Updated catnet-core dependency from development pseudo-version to stable v0.2.0.
- All repository comments, documentation strings, and user-facing messages standardized to English. No functional changes.

### Fixed
- Integration tests no longer share Cobra flag state between cases. `TestMain` now resets `rootCmd` before each test. (Resolves analysis finding C6.)
- `os.Stdout` pipe is now always restored via `defer` in integration tests, preventing file descriptor leaks on early test failure. (Resolves analysis finding C7.)
- `TestScanCancelledByContext` rewritten as a subprocess test. Signal is sent only to the child process, eliminating risk of terminating the test runner. (Resolves analysis finding C10.)

### Added
- Unit tests for `output/human.go` and `output/json.go` using an injected `io.Writer`. Coverage now includes TTY detection fallback and color flag propagation. (Resolves analysis finding C9.)

## [0.1.0] - 2026-06-06

### Added
- Initial scaffolding of the CLI repository (`github.com/mendsec/catnet`).
- Cobra CLI structure with `root`, `scan`, `export`, and `version` subcommands.
- Graceful cancellation handling via `context` and `os/signal` (Exit Code 130).
- Human-readable output formatting with progress bars and terminal color support.
- JSON output formatting for CI/CD and scriptability.
- Feature to re-export previous JSON scans into CSV, XML, and JSON.
- GitHub Actions CI pipelines for cross-compilation (Windows, macOS, Linux).
- GoReleaser configuration for automated binary publishing.
- Integration tests simulating End-to-End behavior via `127.0.0.1` loopback.

[unreleased]: https://github.com/mendsec/catnet/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/mendsec/catnet/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/mendsec/catnet/compare/2721a3346032d02831f4f0594ad6332a57c4f145...v0.1.0
