# hello-cli

A simple CLI for testing GitHub releases and self-update.

## Install

**macOS / Linux:**

```bash
curl -sSL https://raw.githubusercontent.com/ludleth/hello-cli/main/install.sh | sh
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/ludleth/hello-cli/main/install.ps1 | iex
```

**From source:**

```bash
go install github.com/ludleth/hello-cli@latest
```

## Update

```bash
# Update to latest stable
hello-cli update

# Update to a specific version
hello-cli update v0.0.5

# Skip confirmation
hello-cli update -y
```
