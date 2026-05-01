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

The project is being developed in phases.

Phase 4-5 tmux client wrapper and read-only inspection commands have been implemented:

- Go module setup
- Cobra CLI skeleton
- `twx version`
- `twx doctor`
- `twx config path`
- `twx config validate`
- `twx list`
- `twx sessions`
- `twx windows <session>`
- YAML config structs, loading, and validation
- Read-only tmux client wrapper
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
twx config path
twx config validate
twx list
twx sessions
twx windows <session>
twx init <workspace>
twx start <workspace>
twx attach <workspace>
twx kill <workspace>
twx restart <workspace>
twx windows <workspace>
twx add-window <workspace> <window>
twx remove-window <workspace> <window>
twx tpm status
twx tpm install
```

The main source of truth should be:

```
~/.config/twx/config.yaml
```

The CLI should create and manage tmux sessions and windows from that config.

Near-term roadmap:

1. `twx start <workspace>`
2. `twx attach <workspace>`
3. Lifecycle commands: kill, restart
4. Config mutation commands: init, add-window, remove-window

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

| Command | May modify files? |
|---------|-------------------|
| `twx doctor` | No |
| `twx list` | No |
| `twx config validate` | No |
| `twx sessions` | No |
| `twx windows <session>` | No |
| `twx init` | Yes — creates/updates config file |
| `twx add-window` | Yes — updates config file |
| `twx tpm install` | Yes — may modify `~/.tmux.conf`, must back it up first |

`twx list` and `twx config validate` must never create tmux sessions or modify user files.
Read-only tmux inspection commands must not create or modify tmux sessions.

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

Use small packages with clear responsibilities:

```
cmd/                 Cobra commands
internal/config/     Config paths, loading, validation, writing
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

Rules:

- Workspace names should map directly to tmux session names
- Window names should map directly to tmux window names
- `root` should support `~` expansion
- `env` should support workspace-level environment variables
- `command` is optional per window
- Missing config should produce a friendly warning, not a panic

---

## Tmux Behavior Direction

Future `twx start <workspace>` should:

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

Do not execute tmux session creation in `doctor`. `doctor` may check `tmux -V`.

---

## TPM Behavior Direction

TPM support is optional. Future TPM commands should include:

```bash
twx tpm status
twx tpm install
```

Rules:

- Missing TPM is not fatal
- Never overwrite `~/.tmux.conf` without a backup
- Avoid duplicate TPM config blocks
- Print clear next steps after install

---

## Documentation Guidelines

Keep README practical and command-focused. When adding commands, update:

- `README.md`
- `docs/commands.md`
- `docs/roadmap.md`
- `examples/config.yaml`

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

---

## Commit Style

Use simple conventional commits:

```
chore: initialize twx cobra cli skeleton
feat: add config loader and list command
feat: add tmux workspace start command
fix: handle missing config path gracefully
docs: update quick start
```

---

## Current Near-Term Roadmap

Implement next phases in this order:

1. Config loader and validator — `twx config path`, `twx config validate`
2. `twx list`
3. tmux client wrapper
4. `twx start <workspace>`, `twx attach <workspace>`
5. Lifecycle commands — `kill`, `restart`, `sessions`, `windows`
6. Config mutation commands — `init`, `add-window`, `remove-window`
7. TPM status/install

Do not skip directly to advanced features before config loading and start are stable.

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
