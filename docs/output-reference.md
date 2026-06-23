---
layout: default
title: Output Reference
nav_order: 4
description: Complete schema documentation for JSON, CSV, XML, and human output formats.
---

# Output Reference

## Human Output (default)

catnet separates progress and results across two output streams:

- **stderr**: progress bar, lifecycle messages, warnings
- **stdout**: results table (tabwriter-aligned)

**Table columns:** IP, HOSTNAME, MAC, STATUS, PORTS

Colours are disabled when `--no-color` is set or when stdout is not a TTY.

## JSON Output

**Schema version:** `"2.0.0"`

### Top-level fields

| Field          | Type    | Description                              |
|----------------|---------|------------------------------------------|
| `schemaVersion`| string  | Always `"2.0.0"`                         |
| `startTime`    | string  | ISO 8601 UTC                             |
| `endTime`      | string  | ISO 8601 UTC                             |
| `total`        | integer | IPs submitted for scanning               |
| `alive`        | integer | Hosts that responded to ICMP             |
| `devices`      | array   | One entry per scanned host               |

### Device object fields

| Field      | Type         | Description                              |
|------------|--------------|------------------------------------------|
| `ip`       | string       | IPv4 address                             |
| `isAlive`  | boolean      | `true` if host responded to ping         |
| `hostname` | string       | Reverse DNS result (`""` if none)        |
| `mac`      | string       | MAC address (`""` if not on local subnet)|
| `openPorts`| array[int]   | Sorted list of open TCP ports            |

### Full example

```json
{
  "schemaVersion": "2.0.0",
  "startTime": "2026-06-06T12:00:00Z",
  "endTime": "2026-06-06T12:00:05Z",
  "total": 1,
  "alive": 1,
  "devices": [
    {
      "ip": "127.0.0.1",
      "mac": "",
      "hostname": "",
      "isAlive": true,
      "openPorts": [22, 80]
    }
  ]
}
```

### Forward compatibility

Consumers MUST NOT fail on unknown fields. New fields may appear in future minor versions.

## CSV Output

**Header:** `IP,Hostname,MAC,Status,Open Ports`

**Status:** `Alive` | `Dead`

**Open Ports:** semicolon-separated integers (`"80;443"`)

> **Security:** Fields starting with formula-trigger characters are prefixed with a single quote to prevent CSV injection.

## XML Output

**Root:** `<results>` — **Child:** `<device>` (repeated) — **Sub-elements:** `<ip>`, `<hostname>`, `<mac>`, `<status>`

> `openPorts` are not included in XML output.

## Stderr vs Stdout Contract

| Stream | Content |
|--------|---------|
| stdout | Machine-readable data (JSON, CSV, XML, human table) |
| stderr | Human-readable progress, warnings, lifecycle events |

**Pipeline example:**
```bash
catnet scan 192.168.1.0/24 --format json 2>/dev/null | jq '.devices[]'
```
