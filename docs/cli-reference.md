---
layout: default
title: CLI Reference
nav_order: 3
description: Complete flag reference for catnet subcommands — scan, export, version, completion.
---

# CLI Reference

This page is also available on the [CLI Reference wiki page](https://github.com/mendsec/catnet/wiki/CLI-Reference).

## Global Flags

| Flag         | Shorthand | Default | Description                           |
|--------------|-----------|---------|---------------------------------------|
| `--format`   |           | `human` | Output format: `human`, `json`        |
| `--no-color` |           | `false` | Disable ANSI colour in human output   |

---

## `catnet scan`

**Usage:** `catnet scan [targets] [flags]`

**Arguments:** `targets` — comma-separated list of IPs, CIDRs, or dash ranges. Multiple positional arguments are also accepted.

**Flags:**

| Flag             | Shorthand | Default                 | Description                                      |
|------------------|-----------|-------------------------|--------------------------------------------------|
| `--ports`        | `-p`      | `22,80,443,139,445,3389` | Ports to scan                                    |
| `--threads`      | `-t`      | `64`                    | Max concurrent goroutines                        |
| `--ping-timeout` |           | `1000`                  | ICMP timeout in milliseconds                     |
| `--port-timeout` |           | `500`                   | TCP connect timeout in milliseconds              |
| `--timeout`      |           | (none)                  | Hard wall-clock scan deadline (e.g. `30s`)       |
| `--no-ports`     |           | `false`                 | Skip port scanning entirely                      |
| `--output`       | `-o`      | (none)                  | Write JSON to file (implies JSON output)         |
| `--quiet`        | `-q`      | `false`                 | Suppress progress output to stderr               |
| `--format`       |           | `human`                 | Output format: `human`, `json`                   |

**`--timeout` behaviour:** Sets a `context.Deadline` on the entire scan. Partial results are returned if the timeout expires (exit code 2).

**Examples:**

1. **Scan /24 with ports 22 and 80 only:**
   ```bash
   catnet scan 192.168.1.0/24 -p 22,80 --quiet
   ```

2. **JSON output piped to jq:**
   ```bash
   catnet scan 10.0.0.0/24 --format json | jq '.devices[] | select(.isAlive)'
   ```

3. **Save results and convert to CSV:**
   ```bash
   catnet scan 10.0.0.1-254 -o scan.json
   catnet export scan.json --format csv -o scan.csv
   ```

4. **Hard 60-second timeout:**
   ```bash
   catnet scan 172.16.0.0/16 --timeout 60s
   ```

---

## `catnet export`

**Usage:** `catnet export [input.json] --format <fmt> [-o output]`

**Arguments:** `input.json` — path to a previously saved scan result.

**Flags:**

| Flag       | Shorthand | Default    | Description                              |
|------------|-----------|------------|------------------------------------------|
| `--format` | `-f`      | (required) | Output format: `json`, `csv`, `xml`      |
| `--output` | `-o`      | stdout     | Write to file instead of stdout          |

**Forward compatibility:** Unknown JSON fields in the input file are silently ignored. A WARN is printed to stderr if the `schemaVersion` major digit is not 1.

---

## `catnet version`

**Usage:** `catnet version [--short]`

**Flags:**

| Flag      | Shorthand | Default | Description                           |
|-----------|-----------|---------|---------------------------------------|
| `--short` |           | `false` | Print only the version number         |

**Output:**
```
catnet/v0.2.0 (linux/amd64) go1.26.3
catnet-core/v0.2.0
```

---

## `catnet completion`

**Usage:** `catnet completion [bash|zsh|fish|powershell]`

Pipe to a file or eval for persistent completion.
