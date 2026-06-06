# catnet

[![CI](https://github.com/mendsec/catnet/actions/workflows/ci.yml/badge.svg)](https://github.com/mendsec/catnet/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mendsec/catnet)](https://golang.org/doc/go1.23)
[![Release](https://img.shields.io/github/v/release/mendsec/catnet)](https://github.com/mendsec/catnet/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Scriptable network scanner CLI. Part of the [CatNet ecosystem](https://github.com/mendsec/catnet-core).

## Installation

### From binary (recommended)
Download the latest release from the [GitHub Releases](https://github.com/mendsec/catnet/releases) page.

### From source
```bash
go install github.com/mendsec/catnet/cmd/catnet@latest
```

## Usage

### Basic scan
```bash
catnet scan 192.168.1.0/24
```

### Multiple targets
```bash
catnet scan 192.168.1.0/24,10.0.0.1-10
```

### JSON output (for scripting)
```bash
catnet scan 192.168.1.0/24 --format json | jq '.devices[] | select(.isAlive)'
```

### Save and re-export
```bash
catnet scan 192.168.1.0/24 --format json -o result.json
catnet export result.json --format csv -o result.csv
```

### Custom port scan
```bash
catnet scan 192.168.1.0/24 --ports 22,80,443,8080,3306
```

### Quiet mode (scripts/CI)
```bash
catnet scan 192.168.1.0/24 --format json --quiet
```

## Output Reference
The JSON output contains the following structure:
| Field | Type | Description |
|---|---|---|
| `schemaVersion` | string | The version of the JSON schema (e.g., "1.0.0") |
| `startTime` | string | ISO 8601 timestamp when the scan started |
| `endTime` | string | ISO 8601 timestamp when the scan completed |
| `total` | int | Total number of IP addresses scanned |
| `alive` | int | Total number of hosts found alive |
| `devices` | array | List of scanned devices with their details |

Device fields:
| Field | Type | Description |
|---|---|---|
| `ip` | string | IPv4 address |
| `isAlive` | boolean | True if the host responded to ping |
| `hostname` | string | Reverse DNS hostname |
| `mac` | string | MAC address (if on local subnet) |
| `openPorts` | array[int] | List of open ports discovered |

## Exit Codes
| Code | Description |
|---|---|
| `0` | Success |
| `1` | Input error (invalid targets or flags) |
| `2` | Runtime error (engine failure or export failure) |
| `130` | Interrupted (cancelled by signal, e.g., Ctrl+C) |

## Ecosystem
| Repository | Role |
|---|---|
| `catnet-core` | Shared scanning engine (no CLI/GUI) |
| `catnet` | **This repository — Scriptable CLI** |
| `catnet-tui` | Interactive TUI interface |
| `catnet-scanner` | Desktop GUI application |

## License
MIT
