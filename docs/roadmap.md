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
- `twx workspace add <workspace>`.
- `twx workspace remove <workspace>`.
- `twx workspace show <workspace>`.
- `twx attach <workspace>`.
- `twx kill <workspace>`.
- `twx restart <workspace>`.
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

## Phase 8 Implemented

- `twx workspace add <workspace>`.
- `twx workspace remove <workspace>`.
- `twx workspace show <workspace>`.
- Safe config backups before workspace writes.
- Workspace inspection without tmux side effects.

## Phase 9 Implemented

- `twx attach <workspace>`.
- `twx kill <workspace>`.
- `twx restart <workspace>`.
- `twx restart --no-attach`.
- Lifecycle commands operate on tmux sessions without modifying config.

## Phase 10 Implemented

- `twx window add <workspace> <window>`.
- `twx window remove <workspace> <window>`.
- `twx window set-command <workspace> <window> <command>`.
- Safe config backups before window writes.
- Window inspection and mutation without tmux side effects.

## Phase 11 Implemented

- `twx tpm status`.
- `twx tpm install`.
- Read-only TPM detection.
- Idempotent cloning and config updating.

## Next Phase

- Phase 12: Release Packaging
  - version injection
  - GitHub Releases
  - `.deb` package
  - Homebrew tap

## Later Ideas

- TPM status and install.
- Interactive workspace wizard.
- Pane layout support.
- Shell script import and export.
- GitHub release binaries.
- Release packaging.
