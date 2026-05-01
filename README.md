# twx - Declarative tmux Workspace Manager

`twx` is a Go CLI for managing tmux workspaces declaratively on Ubuntu.

## Problem

Development workspaces often start as useful tmux shell scripts and slowly turn into repeated, project-specific setup commands. `twx` is intended to replace those scripts with a readable YAML configuration and a focused CLI.

## Target Ubuntu Versions

- Ubuntu 20.04
- Ubuntu 22.04
- Ubuntu 24.04

## MVP Scope

The current milestone establishes the project foundation and read-only config commands:

- Go module setup.
- Cobra CLI skeleton.
- `twx version`.
- `twx doctor`.
- `twx config path`.
- `twx config validate`.
- `twx list`.
- Documentation and examples placeholders.
- Local install script.
- GitHub Actions CI.

Full workspace start, attach, tmux session creation, and TPM installation are planned for later phases.

## Development

`twx --help` includes a compact ASCII banner and command overview.

```sh
go mod tidy
go run . --help
go run . version
go run . doctor
go run . config path
go run . --config ./examples/config.yaml config validate
go run . --config ./examples/config.yaml list
```

## Config Path

The default config path is:

```text
~/.config/twx/config.yaml
```

Use the included example config for local read-only checks:

```sh
./twx --config ./examples/config.yaml config validate
./twx --config ./examples/config.yaml list
```

## Future Roadmap

Planned work includes tmux start and attach workflows, window management, TPM status and install helpers, pane layout support, an interactive workspace wizard, shell script import/export, and GitHub release binaries.

## License

MIT
