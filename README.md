# twx — Declarative tmux Workspace Manager

`twx` is a Go CLI, built with Cobra, for managing tmux workspaces declaratively on Ubuntu.

Define your workspaces once in YAML, then create, inspect, restart, and manage tmux sessions from a single command.

---

## Why

Development workspaces often start as useful tmux shell scripts and slowly become repeated, project-specific setup commands.

`twx` replaces those scripts with:

- A readable YAML config
- A focused CLI with safe mutations
- Tmux workspace lifecycle commands
- Optional TPM setup helpers
- Single-binary, Ubuntu-first design

---

## Target Ubuntu Versions

- Ubuntu 20.04 LTS
- Ubuntu 22.04 LTS
- Ubuntu 24.04 LTS

---

## Tech Stack

`twx` is built with:

| Component | Purpose |
|-----------|---------|
| [Go](https://go.dev/) | Single-binary CLI implementation |
| [Cobra](https://github.com/spf13/cobra) | Command structure and CLI framework |
| [yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) | YAML config parsing and writing |
| [tmux](https://github.com/tmux/tmux) | Workspace/session backend |
| [GoReleaser](https://goreleaser.com/) | Release automation for Linux binaries and `.deb` packages |

The CLI is intentionally dependency-light and Ubuntu-first.

---

## Features

### Config Management

```bash
twx config init      # Initialize config
twx config path      # Show config path
twx config validate  # Validate config syntax
twx list             # List all configured workspaces
```

**Runtime config location:** `~/.config/twx/config.yaml`

> **Note:** `examples/config.yaml` is sample/development data only. It is not the default runtime config.

### Workspace Management

```bash
twx workspace add <workspace>     # Add a new workspace
twx workspace show <workspace>    # Show workspace config
twx workspace remove <workspace>  # Remove workspace from config
```

Workspace config changes are written to YAML and backed up before mutation.

### Window Management

```bash
twx window add <workspace> <window>              # Add a window to workspace config
twx window remove <workspace> <window>           # Remove window from workspace config
twx window set-command <workspace> <window> CMD  # Update window command
```

Window mutation commands modify config only. They do not create, kill, or restart tmux sessions directly.

### Tmux Lifecycle

```bash
twx start <workspace>    # Create and attach tmux session from config
twx attach <workspace>   # Attach to an existing tmux session
twx restart <workspace>  # Kill and recreate tmux session from config
twx kill <workspace>     # Kill tmux session only; config is unchanged
```

**Behavior:**

- `twx start` creates tmux sessions and windows from config
- `twx kill` kills the tmux session only; it does not remove workspace config
- Config changes can be applied to a live session with `twx restart`

### Tmux Inspection

```bash
twx sessions         # List active tmux sessions
twx windows SESSION  # List windows in a tmux session
```

Read-only inspection commands are safe to run anytime.

### TPM Helpers

```bash
twx tpm status   # Check TPM status and config
twx tpm install  # Install TPM and plugin config
```

TPM support is optional. `twx tpm install` backs up `~/.tmux.conf` before modifying it.

---

## Installation

`twx` is Ubuntu-first. Multiple installation methods are available.

### .deb Package (Recommended)

Download the latest `.deb` package for your architecture from the [GitHub Releases](https://github.com/jrarkaan/tmux-workspace/releases) page:

```bash
sudo dpkg -i twx_<version>_linux_amd64.deb
```

**Supported architectures:**

- `amd64` — Intel/AMD 64-bit
- `arm64` — ARM 64-bit

### Binary Archive

Download the `.tar.gz` archive, extract it, and move the binary to your PATH:

```bash
tar -xzf twx_<version>_linux_amd64.tar.gz
sudo mv twx /usr/local/bin/twx
```

### Local Development Install

For contributors installing from source:

```bash
git clone https://github.com/jrarkaan/tmux-workspace.git
cd tmux-workspace
./scripts/install-local.sh
```

This builds from source and installs the binary locally to `~/bin/twx`.

> **For normal Ubuntu usage, prefer the `.deb` package from GitHub Releases.**

---

## Quickstart

### Step 1: Initialize Config

```bash
twx config init
```

This creates `~/.config/twx/config.yaml` with an empty workspace list.

### Step 2: Add a Workspace

```bash
twx workspace add backend-dev \
  --root ~/App/backend \
  --windows overview,codex-impl,gemini-review,test-watch,runtime-logs,git-diff,db,ssh \
  --env APP_ENV=development \
  --command overview="clear; pwd; git status"
```

### Step 3: Inspect Config

```bash
twx list
twx workspace show backend-dev
```

### Step 4: Start the Workspace

```bash
# Start without attaching
twx start backend-dev --no-attach

# Start and attach
twx start backend-dev
```

### Step 5: Inspect Tmux State

```bash
twx sessions
twx windows backend-dev
```

### Step 6: Manage the Workspace

```bash
# Attach to an existing session
twx attach backend-dev

# Restart: kill and recreate from config
twx restart backend-dev --no-attach

# Kill the tmux session
twx kill backend-dev

# Remove workspace from config
twx workspace remove backend-dev --force
```

> **Note:** Removing a workspace from config does not kill a running tmux session.

---

## Window Commands

### Add a Window

```bash
twx window add backend-dev jaeger \
  --command "echo jaeger logs"
```

### Update a Window Command

```bash
twx window set-command backend-dev overview "clear; pwd; git status"
```

### Remove a Window

```bash
twx window remove backend-dev jaeger --force
```

> **Note:** Window removal requires `--force`. Window config changes do not affect running tmux sessions until you run `twx restart`.

---

## TPM Setup

### Check TPM Status

```bash
twx tpm status
```

Shows detailed status of:

- tmux
- git availability
- TPM directory (`~/.tmux/plugins/tpm`)
- `~/.tmux.conf`
- TPM plugin configuration
- tmux-resurrect configuration
- tmux-continuum configuration

### Install TPM

```bash
twx tpm install
```

This command:

- Creates `~/.tmux/plugins` directory if needed
- Clones TPM into `~/.tmux/plugins/tpm`
- Creates `~/.tmux.conf` if needed
- Backs up existing `~/.tmux.conf` before modifying it
- Appends the twx TPM managed block
- Is idempotent and safe to run multiple times

**After installation,** open tmux and press your prefix followed by capital `I` to install plugins:

```bash
# Default tmux prefix
Ctrl+b I

# Common custom prefix (from project docs)
Ctrl+a I
```

---

## Configuration

The YAML config schema:

```yaml
version: 1

defaults:
  attach: true
  create_if_missing: true
  base_index: 1

workspaces:
  backend-dev:
    root: ~/App/backend
    env:
      APP_ENV: development
    windows:
      - name: overview
        command: clear; pwd; git status
      - name: codex-impl
      - name: gemini-review
      - name: test-watch
        command: echo 'Run tests here'
      - name: runtime-logs
      - name: git-diff
        command: git status
      - name: db
      - name: ssh
```

### Config Rules

- Workspace names map directly to tmux session names
- Window names map directly to tmux window names
- `root` supports `~` expansion during execution and is preserved as-is in YAML
- `env` supports workspace-level environment variables
- `command` is optional per window
- Empty `workspaces: {}` is valid after `twx config init`

---

## Development

### Setup

```bash
git clone https://github.com/jrarkaan/tmux-workspace.git
cd tmux-workspace
```

### Build

```bash
go mod tidy
go fmt ./...
go test ./...
go vet ./...
go build -o twx .
```

### Local Testing

```bash
# Show help and ASCII banner
./twx --help

# Check version
./twx version

# Run system diagnostics
./twx doctor

# Print default config without writing
./twx config init --print

# Validate example config
./twx --config ./examples/config.yaml config validate

# List from example config
./twx --config ./examples/config.yaml list
```

### Project Structure

```
cmd/              Cobra commands
internal/config/  Config loading, validation, writing, and mutation
internal/system/  OS and environment checks
internal/tmux/    tmux command wrapper and execution
internal/tpm/     TPM status and install helpers
scripts/          Build and install scripts
docs/             Command and feature documentation
examples/         Sample configuration (development-only)
```

---

## Release

This project uses [GoReleaser](https://goreleaser.com/) for Ubuntu-focused releases.

### Artifacts

Expected release artifacts:

- `twx_<version>_linux_amd64.tar.gz`
- `twx_<version>_linux_arm64.tar.gz`
- `twx_<version>_linux_amd64.deb`
- `twx_<version>_linux_arm64.deb`
- `checksums.txt`

### Local Snapshot Build

If GoReleaser is installed:

```bash
goreleaser check
goreleaser release --snapshot --clean
```

Snapshot artifacts are written to `dist/`.

### Creating a Release

Push a Git tag to trigger the GitHub Actions release workflow:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Artifacts are published to [GitHub Releases](https://github.com/jrarkaan/tmux-workspace/releases).

---

## Roadmap

Planned future features:

- Homebrew tap support
- Shell completions for bash, zsh, and fish
- Release polishing and pre-release support
- Pane layout DSL
- Interactive workspace wizard
- Shell script import/export
- Optional APT repository publishing

---

## Documentation

For more details, see:

- [`docs/commands.md`](docs/commands.md) — Detailed command reference
- [`docs/config.md`](docs/config.md) — Configuration guide
- [`docs/tpm.md`](docs/tpm.md) — TPM setup guide
- [`docs/roadmap.md`](docs/roadmap.md) — Feature roadmap and development phases
- [`docs/release.md`](docs/release.md) — Release and packaging notes
- [`AGENTS.md`](AGENTS.md) — Developer guidance for contributors and agents

---

## Contributing

Contributions are welcome. Please read [`AGENTS.md`](AGENTS.md) for development guidelines, safety rules, and project direction.

---

## License

MIT

---

## Related Projects

This project grew from a personal tmux workspace workflow and an earlier Bash-script approach. See [`tmux-workflow`](https://github.com/jrarkaan/tmux-workflow) for related workflow ideas and historical context.