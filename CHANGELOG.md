# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
