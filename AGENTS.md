# AGENTS.md

Guidance for AI coding agents working on this repository.

---

## Project Overview

This repository contains `twx`, a Go CLI for managing tmux workspaces declaratively on Ubuntu.

The goal is to replace repetitive project-specific tmux shell scripts with a single CLI and a YAML configuration file.

| Item | Value |
|------|-------|
| Repository | `tmux-workspace` |
| Binary | `twx` |
| Go module | `github.com/jrarkaan/tmux-workspace` |
| Target OS | Ubuntu 20.04, 22.04, 24.04 |
| CLI framework | Cobra |
| Config format | YAML |
| Default config path | `~/.config/twx/config.yaml` |

---

## Current Phase

The project is being developed in phases. Phase 9 (lifecycle commands) has been completed.

### Implemented Features

- Go module setup
- Cobra CLI skeleton
- `twx version`
- `twx doctor`
- `twx config init`
- `twx config path`
- `twx config validate`
- `twx list`
- `twx sessions`
- `twx windows <session>`
- `twx start <workspace>`
- `twx workspace add <workspace>`
- `twx workspace show <workspace>`
- `twx workspace remove <workspace>`
- `twx attach <workspace>`
- `twx kill <workspace>`
- `twx restart <workspace>`
- YAML config structs, loading, validation, writing, and mutation
- Read-only tmux client wrapper
- Workspace session/window creation from config
- Safe runtime config initialization with backup-on-force
- Workspace config mutation with backup-on-write
- Workspace lifecycle commands for attach, kill, and restart
- Basic docs
- Examples placeholder
- Local install script
- GitHub Actions CI

Do not assume later-phase features are already implemented unless the code exists.

---

## Product Direction

`twx` should eventually support:

```bash
twx --help
twx doctor
twx version

twx config init
twx config path
twx config validate

twx list

twx workspace add <workspace>
twx workspace remove <workspace>
twx workspace show <workspace>

twx window add <workspace> <window>
twx window remove <workspace> <window>
twx window set-command <workspace> <window> <command>

twx sessions
twx windows <session>

twx start <workspace>
twx attach <workspace>
twx kill <workspace>
twx restart <workspace>

twx tpm status
twx tpm install
```

### Source of Truth

The main source of truth is: `~/.config/twx/config.yaml`

- Runtime config lives at `~/.config/twx/config.yaml`
- `examples/config.yaml` is documentation/testing only and must not be treated as the default runtime config
- The CLI creates and manages tmux sessions and windows from that config

### Near-Term Roadmap

1. Window mutation commands
   - `twx window add <workspace> <window>`
   - `twx window remove <workspace> <window>`
   - `twx window set-command <workspace> <window> <command>`
2. TPM status/install
   - `twx tpm status`
   - `twx tpm install`
3. Release packaging
   - Version injection
   - GitHub Releases
   - `.deb` package
   - Homebrew tap

---

## Design Principles

Keep the project:

- Simple
- Idiomatic Go
- Easy to run locally
- Safe for user machines
- Ubuntu-first
- Declarative where possible
- Idempotent where possible
- Friendly to DevOps/platform workflows

Avoid over-engineering. Do not introduce unnecessary frameworks, background daemons, databases, or network services.

---

## Safety Rules

Be careful with commands that modify user state. Do not modify user files unless the relevant command explicitly exists for that purpose.

### Command Safety Matrix

| Command | May modify files? | May modify tmux? |
|---------|-------------------|------------------|
| `twx doctor` | No | No |
| `twx list` | No | No |
| `twx config path` | No | No |
| `twx config validate` | No | No |
| `twx sessions` | No | No |
| `twx windows <session>` | No | No |
| `twx workspace show <workspace>` | No | No |
| `twx tpm status` | No | No |
| `twx config init` | Yes — creates config file | No |
| `twx workspace add <workspace>` | Yes — updates config file | No |
| `twx workspace remove <workspace>` | Yes — updates config file | No |
| `twx window add <workspace> <window>` | Yes — updates config file | No |
| `twx window remove <workspace> <window>` | Yes — updates config file | No |
| `twx window set-command <workspace> <window> <command>` | Yes — updates config file | No |
| `twx start <workspace>` | No | Yes — creates sessions/windows |
| `twx attach <workspace>` | No | No (attaches only) |
| `twx kill <workspace>` | No | Yes — kills session only |
| `twx restart <workspace>` | No | Yes — kills/recreates session |
| `twx tpm install` | Yes — may modify `~/.tmux.conf` | No |

### Rules by Category

**Read-only inspection commands** (`list`, `sessions`, `windows`, `config validate`, etc.):

- Must never create tmux sessions
- Must never modify user files
- Must never modify tmux config

**Config read commands** (`config path`, `workspace show`):

- Must not create tmux sessions
- Must not modify config files
- Must not modify tmux

**Workspace start** (`start`):

- Allowed to create tmux sessions/windows
- Must not modify config files
- Must not modify tmux config
- Idempotent: safe to run multiple times on same workspace

**Config init** (`config init`):

- May create config files
- `--force` must back up before overwrite
- `--print` must not write files

**Workspace mutation** (`workspace add`, `workspace remove`):

- May modify config files
- `workspace remove` must not kill tmux sessions
- Every config write must validate before writing
- Every config write must back up existing config first

**Window mutation** (`window add`, `window remove`, `window set-command`):

- May modify config files
- Must not create, kill, restart, or attach to tmux sessions
- Every config write must validate before writing
- Every config write must back up existing config first
- Removing/editing a window in config must not affect running tmux sessions
- Users can apply config changes later with `twx restart <workspace>`

**Lifecycle** (`attach`, `kill`, `restart`):

- Operate on tmux sessions only
- `kill` must not remove workspace config
- `restart` may kill and recreate sessions but must not modify config

**TPM** (`tpm install`):

- May modify `~/.tmux.conf`
- Must back up before modifying
- Must be idempotent

### General File Safety

When writing files:

- Create parent directories if needed
- Preserve existing config where possible
- Avoid overwriting user data
- Make backups before destructive or broad changes
- Prefer clear error messages over silent behavior

---

## Go Coding Guidelines

Use standard Go style.

Before completing work, run:

```bash
go mod tidy
go fmt ./...
go test ./...
go vet ./...
go build -o twx .
```

If command behavior changed, also run relevant CLI checks:

```bash
./twx --help
./twx version
./twx doctor
```

### Package Structure

Use small packages with clear responsibilities:

```
cmd/                 Cobra commands
internal/config/     Config paths, loading, validation, writing, and mutation
internal/system/     OS and environment checks
internal/tmux/       tmux command wrapper and workspace execution
internal/tpm/        TPM status/install helpers
```

Cobra commands should stay thin. Business logic should live under `internal/`.

---

## Dependency Guidelines

Current intended dependencies:

```
github.com/spf13/cobra
gopkg.in/yaml.v3
```

Do not add new dependencies unless they clearly reduce complexity. Avoid large UI, TUI, or prompt libraries during MVP unless explicitly requested.

---

## Config Schema Direction

The intended v1 YAML shape:

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

### Schema Rules

- Workspace names should map directly to tmux session names
- Window names should map directly to tmux window names
- `root` should support `~` expansion during tmux execution
- `root` should be preserved as provided when writing YAML
- `env` should support workspace-level environment variables
- `command` is optional per window
- Missing config should produce a friendly warning, not a panic
- Empty `workspaces: {}` is valid after `twx config init`

---

## Window Mutation Direction

Phase 10 will implement window mutation commands with this namespace:

```bash
twx window add <workspace> <window>
twx window remove <workspace> <window>
twx window set-command <workspace> <window> <command>
```

### Product Rules for Window Mutation

- Operate only on YAML config
- Do not create tmux windows directly
- Do not kill tmux windows directly
- Do not attach to tmux
- Do not restart tmux sessions
- Validate config before and after mutation
- Back up existing config before writing
- Preserve workspace order and window order as much as practical
- Reject duplicate window names inside a workspace
- Return clear errors for missing config, missing workspace, missing window, duplicate window, and invalid command usage

### Suggested Behavior

```bash
twx window add backend-dev jaeger
twx window add backend-dev logs --command "kubectl logs -f deploy/backend"
twx window set-command backend-dev overview "clear; pwd; git status"
twx window remove backend-dev jaeger --force
```

**Removal:** Require `--force` unless interactive confirmation is explicitly implemented. For MVP, prefer `--force` over prompts.

**Config changes:** Removing or editing a window in config must not affect any currently running tmux session. Users can apply config changes later with:

```bash
twx restart <workspace>
```

---

## Tmux Behavior Direction

### `twx start <workspace>`

1. Load config
2. Resolve workspace
3. Expand `root` path
4. Check if tmux session exists
5. If session exists, attach unless flags say otherwise
6. If session does not exist, create it
7. Create configured windows
8. Send optional initial commands
9. Select the first window
10. Attach unless `--no-attach` is used

### `twx attach <workspace>`

Attach only to an existing session.

### `twx kill <workspace>`

Kill only the tmux session. Must not remove config.

### `twx restart <workspace>`

May kill and recreate a tmux session from config, but must not modify config.

### `twx doctor`

Do not execute tmux session creation. May check `tmux -V`.

---

## TPM Behavior Direction

TPM support is optional. Future TPM commands should include:

```bash
twx tpm status
twx tpm install
```

### TPM Rules

- Missing TPM is not fatal
- Never overwrite `~/.tmux.conf` without a backup
- Avoid duplicate TPM config blocks
- Print clear next steps after install
- `twx tpm status` must be read-only
- `twx tpm install` may clone TPM and update `~/.tmux.conf`, but must be idempotent and backup first

---

## Documentation Guidelines

Keep README practical and command-focused. When adding commands, update:

- `README.md`
- `docs/commands.md`
- `docs/roadmap.md`
- `examples/config.yaml`
- `AGENTS.md`

Prefer copy-pasteable examples. Use concise explanations.

---

## Testing Guidelines

For pure logic, prefer unit tests. For tmux behavior, avoid brittle tests that require a live tmux server during normal CI unless explicitly designed.

Initial CI should stay simple:

```bash
test -z "$(gofmt -l .)"
go test ./...
go vet ./...
go build ./...
```

### Test Best Practices

**Config mutation commands:** Use temp directories in tests. Never modify real `~/.config/twx/config.yaml` in tests.

**Tmux lifecycle commands:** Prefer fake runners for unit tests. Use separate manual integration checks for real tmux behavior.

---

## Commit Style

Use simple conventional commits:

```
chore: initialize twx cobra cli skeleton
feat: add config loader and list command
feat: add tmux workspace start command
feat: add workspace lifecycle commands
feat: add window mutation commands
feat: add tpm management commands
fix: handle missing config path gracefully
docs: update quick start
```

---

## Current Near-Term Roadmap

Implement next phases in this order:

### Phase 10: Window Mutation

- `twx window add <workspace> <window>`
- `twx window remove <workspace> <window>`
- `twx window set-command <workspace> <window> <command>`

### Phase 11: TPM Management

- `twx tpm status`
- `twx tpm install`

### Phase 12: Release Packaging

- Version injection
- GitHub Releases
- `.deb` package
- Homebrew tap

Do not skip directly to advanced features before config loading, workspace mutation, start, and lifecycle commands are stable.

---

## Things Not To Do Yet

Do not implement these during MVP unless explicitly requested:

- Cloud sync
- Remote workspace registry
- Pane layout DSL
- Shell script import/export
- Interactive TUI
- Background daemon
- Plugin system
- Non-Ubuntu platform support
- Automatic execution of development servers without explicit config