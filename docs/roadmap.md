# Roadmap

## v0.1.0 Goals

- Go module setup.
- Cobra CLI skeleton.
- `twx version`.
- `twx doctor`.
- `twx config path`.
- `twx config validate`.
- `twx list`.
- `twx sessions`.
- `twx windows <session>`.
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

## Next Phase

- `twx start <workspace>`.
- `twx attach <workspace>`.
- `twx kill <workspace>`.
- `twx restart <workspace>`.

## Later Ideas

- TPM status and install.
- Interactive workspace wizard.
- Pane layout support.
- Shell script import and export.
- GitHub release binaries.
