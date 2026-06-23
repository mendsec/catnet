---
layout: default
title: Use Cases
nav_order: 5
description: Real-world recipes for catnet — CI/CD pipelines, reconnaissance, jq integration, Elasticsearch, and cron scans.
---

# Use Cases

## Asset Discovery in CI/CD

Nightly GitHub Actions job that scans a staging subnet and fails the pipeline if unexpected hosts appear.

```yaml
name: nightly-asset-discovery
on:
  schedule:
    - cron: '0 2 * * *'

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install catnet
        run: |
          curl -sSL https://github.com/mendsec/catnet/releases/latest/download/catnet_Linux_x86_64.tar.gz | tar xz
          sudo mv catnet /usr/local/bin/
      - name: Scan staging subnet
        run: catnet scan 10.0.1.0/24 --quiet --format json -o scan.json
      - name: Check for unexpected hosts
        run: |
          EXPECTED="known-hosts.txt"
          cat scan.json | jq -r '.devices[] | select(.isAlive) | .ip' | sort > alive.txt
          comm -13 <(sort "$EXPECTED") alive.txt > unexpected.txt
          if [ -s unexpected.txt ]; then
            echo "UNEXPECTED HOSTS FOUND:"; cat unexpected.txt; exit 1
          fi
```

## Red Team Reconnaissance

```bash
catnet scan 10.0.0.0/8 --no-ports -t 256 --quiet --format json -o hosts.json
catnet export hosts.json --format csv -o hosts.csv
```

> Always obtain written authorisation before scanning.

## jq Pipelines

```bash
# List alive hosts
catnet scan 10.0.0.0/24 --format json | jq '.devices[] | select(.isAlive)'

# Filter by open port
catnet scan 10.0.0.0/24 --format json | jq '.devices[] | select(.openPorts | contains([22]))'

# Count alive vs total
catnet scan 10.0.0.0/24 --format json | jq '{total: .total, alive: .alive}'

# Extract IPs as newline-separated list
catnet scan 10.0.0.0/24 --format json | jq -r '.devices[] | select(.isAlive) | .ip'
```

## Feeding Results into Elasticsearch

```bash
catnet scan 10.0.0.0/24 --format json | \
  curl -X POST "https://elasticsearch:9200/catnet-scans/_doc/" \
    -H "Content-Type: application/json" \
    -d @-
```

## Scheduled Cron Scan with Change Detection

```bash
#!/bin/bash
TODAY="/tmp/scan-$(date +%Y%m%d).json"
YESTERDAY="/tmp/scan-$(date -d yesterday +%Y%m%d).json"

catnet scan 192.168.0.0/16 --quiet --format json -o "$TODAY"

if [ -f "$YESTERDAY" ]; then
  diff <(jq -r '.devices[] | select(.isAlive) | .ip' "$TODAY" | sort) \
       <(jq -r '.devices[] | select(.isAlive) | .ip' "$YESTERDAY" | sort)
fi
```

## Scripting with Exit Codes

```bash
catnet scan 192.168.1.1 --quiet --format json -o /tmp/out.json
case $? in
  0)   echo "Success" ;;
  130) echo "Interrupted" ;;
  *)   echo "Failed" >&2 ;;
esac
```
