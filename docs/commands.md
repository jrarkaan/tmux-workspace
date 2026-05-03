# Commands

Many commands are planned and not implemented yet. Implemented commands are marked below.

## System

- `twx doctor` - implemented
- `twx version` - implemented

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

- `twx add-window <workspace> <window>` - planned
- `twx remove-window <workspace> <window>` - planned

## TPM

- `twx tpm status` - planned
- `twx tpm install` - planned
