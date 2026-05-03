# twx - Declarative tmux Workspace Manager

`twx` is a Go CLI for managing tmux workspaces declaratively on Ubuntu.

## Problem

Development workspaces often start as useful tmux shell scripts and slowly turn into repeated, project-specific setup commands. `twx` is intended to replace those scripts with a readable YAML configuration and a focused CLI.

## Target Ubuntu Versions

- Ubuntu 20.04
- Ubuntu 22.04
- Ubuntu 24.04

## MVP Scope

The current milestone establishes the project foundation, config commands, tmux inspection, and workspace start:

- Go module setup.
- Cobra CLI skeleton.
- `twx version`.
- `twx doctor`.
- `twx config init`.
- `twx config path`.
- `twx config validate`.
- `twx list`.
- `twx sessions`.
- `twx windows <session>`.
- `twx start <workspace>`.
- `twx workspace add <workspace>`.
- `twx workspace show <workspace>`.
- `twx workspace remove <workspace>`.
- `twx window add <workspace> <window>`.
- `twx window remove <workspace> <window>`.
- `twx window set-command <workspace> <window> <command>`.
- `twx attach <workspace>`.
- `twx kill <workspace>`.
- `twx restart <workspace>`.
- `twx tpm status`.
- `twx tpm install`.
- YAML config structs, loading, validation, writing, and mutation.
- Read-only tmux client wrapper.
- Workspace session/window creation from config.
- Safe runtime config initialization with backup-on-force.
- Workspace and window config mutation with backup-on-write.
- Workspace lifecycle commands for attach, kill, and restart.
- TPM detection and installation configuration management.
- Documentation and examples placeholders.
- Local install script.
- GitHub Actions CI.

Release packaging is planned for later phases.

## Development

`twx --help` includes a compact ASCII banner and command overview.

```sh
go mod tidy
go run . --help
go run . version
go run . doctor
go run . config init --print
go run . config path
go run . --config ./examples/config.yaml config validate
go run . --config ./examples/config.yaml list
go run . sessions
go run . windows backend-dev
go run . --config ./examples/config.yaml start backend-dev --no-attach
```

## Config Path

The default config path is:

```text
~/.config/twx/config.yaml
```

Initialize the runtime config with:

```sh
twx config init
twx config path
twx config validate
twx list
```

`twx config init --print` prints the default config without writing files. `twx config init --force` backs up an existing config before overwriting it.

`examples/config.yaml` is sample/development data only. It is not the default runtime config.

## Quickstart

```sh
twx config init

twx workspace add backend-dev \
  --root ~/App/backend \
  --windows overview,codex-impl,gemini-review,test-watch,runtime-logs,git-diff,db,ssh \
  --env APP_ENV=development \
  --command overview="clear; pwd; git status"

twx list
twx start backend-dev --no-attach
twx attach backend-dev
twx restart backend-dev --no-attach
twx kill backend-dev
twx windows backend-dev
```

Removing a workspace from config does not kill any running tmux session.
`twx kill` kills the tmux session only. It does not remove the workspace from config; use `twx workspace remove <workspace> --force` to remove the config entry.

Use the included example config for local checks:

```sh
./twx --config ./examples/config.yaml config validate
./twx --config ./examples/config.yaml list
./twx sessions
./twx windows backend-dev
./twx --config ./examples/config.yaml start backend-dev --no-attach
```

`--no-attach` creates or restarts the tmux session without attaching. `--force` recreates an existing session from config.

## Future Roadmap

Planned work includes window management, TPM status and install helpers, pane layout support, an interactive workspace wizard, shell script import/export, and GitHub release binaries.

## License

MIT
