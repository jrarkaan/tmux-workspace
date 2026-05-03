# Commands

Many commands are planned and not implemented yet. Implemented commands are marked below.

## System

- `twx doctor` - implemented
- `twx version` - implemented (outputs multi-line version, commit, date)

## Config

- `twx config init` - implemented
  - `--force` backs up and overwrites an existing config.
  - `--print` prints the default config without writing files.
- `twx config path` - implemented
- `twx config validate` - implemented

## Workspace

- `twx list` - implemented
- `twx sessions` - implemented
- `twx windows <session>` - implemented
- `twx start <workspace>` - implemented
- `twx workspace add <workspace>` - implemented
  - `--root` sets the workspace root.
  - `--windows` sets comma-separated window names.
  - `--env` may be repeated as `KEY=VALUE`.
  - `--command` may be repeated as `WINDOW=COMMAND`.
  - `--force` replaces an existing workspace.
- `twx workspace show <workspace>` - implemented
- `twx workspace remove <workspace>` - implemented
  - `--force` confirms removal from config.
- `twx attach <workspace>` - implemented
- `twx kill <workspace>` - implemented
- `twx restart <workspace>` - implemented
  - `--no-attach` recreates the session without attaching.

## Window Management

- `twx window add <workspace> <window>` - implemented
  - `--command` sets the initial command for the window.
  - `--force` replaces an existing window.
- `twx window remove <workspace> <window>` - implemented
  - `--force` confirms removal from config.
- `twx window set-command <workspace> <window> <command>` - implemented

## TPM

- `twx tpm status` - implemented
- `twx tpm install` - implemented
