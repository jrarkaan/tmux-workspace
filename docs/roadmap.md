# Roadmap

## v0.1.0 Goals

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
- Documentation and examples placeholders.
- Local install script.
- GitHub Actions CI.

## Phase 3 Implemented

- YAML config structs.
- Config loading from disk.
- Config validation with aggregated errors.
- Read-only workspace listing.

## Phase 4-5 Implemented

- tmux client wrapper.
- Read-only active session inspection.
- Read-only window inspection for an existing tmux session.

## Phase 6 Implemented

- `twx start <workspace>`.
- Workspace session/window creation from YAML config.
- `--no-attach` support.
- `--force` session recreation.

## Phase 7 Implemented

- `twx config init`.
- Safe default runtime config creation at `~/.config/twx/config.yaml`.
- `--print` support for previewing default config.
- `--force` backup before overwrite.
- Empty workspace configs validate and list cleanly.

## Next Phase

- `twx workspace add <workspace>`.
- `twx workspace remove <workspace>`.
- `twx workspace show <workspace>`.
- Window add/remove commands.
- Then `twx attach`, `twx kill`, and `twx restart`.

## Later Ideas

- TPM status and install.
- Interactive workspace wizard.
- Pane layout support.
- Shell script import and export.
- GitHub release binaries.
