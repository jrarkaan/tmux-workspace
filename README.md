# twx - Declarative tmux Workspace Manager

`twx` is a Go CLI for managing tmux workspaces declaratively on Ubuntu.

## Problem

Development workspaces often start as useful tmux shell scripts and slowly turn into repeated, project-specific setup commands. `twx` is intended to replace those scripts with a readable YAML configuration and a focused CLI.

## Target Ubuntu Versions

- Ubuntu 20.04
- Ubuntu 22.04
- Ubuntu 24.04

## MVP Scope

The first milestone establishes the project foundation:

- Go module setup.
- Cobra CLI skeleton.
- `twx version`.
- `twx doctor`.
- Documentation and examples placeholders.
- Local install script.
- GitHub Actions CI.

Full workspace start, list, attach, tmux session creation, and TPM installation are planned for later phases.

## Development

```sh
go mod tidy
go run . --help
go run . version
go run . doctor
```

## Config Path

The default config path is:

```text
~/.config/twx/config.yaml
```

## Future Roadmap

Planned work includes a YAML config loader, workspace listing, tmux start and attach workflows, window management, TPM status and install helpers, pane layout support, an interactive workspace wizard, shell script import/export, and GitHub release binaries.

## License

MIT
