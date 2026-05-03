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
- `twx attach <workspace>` - planned
- `twx kill <workspace>` - planned
- `twx restart <workspace>` - planned
- `twx workspace add <workspace>` - planned
- `twx workspace remove <workspace>` - planned
- `twx workspace show <workspace>` - planned

## Window Management

- `twx add-window <workspace> <window>` - planned
- `twx remove-window <workspace> <window>` - planned

## TPM

- `twx tpm status` - planned
- `twx tpm install` - planned
