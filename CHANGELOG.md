# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Created `human_test.go` for isolated unit testing of standard and quiet human output formatting through dependency injection.

### Changed
- Updated `catnet-core` dependency to `v0.1.1`.
- Migrated integration tests to use `os/exec` subprocesses with a `TestMain` binary builder to achieve true state isolation and reliable signal cancellation tests.

### Fixed
- Fixed `--format` flag shadowing issue in `exportCmd` by making it locally scoped and mutually exclusive, avoiding inheritance from `rootCmd`.
- Fixed test state bleeding by ensuring `os.Args` and `os.Stdout` are restored properly using `defer` during tests that required temporary global mutation.
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
