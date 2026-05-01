# Product

## Problem

tmux is excellent for long-running development sessions, but repeated project setup often turns into a pile of custom shell scripts. Those scripts are usually hard to share, hard to inspect, and easy to let drift between machines.

## Goal

twx aims to make tmux workspaces declarative. A developer should be able to describe a workspace once in YAML, then use a small CLI to start, attach to, inspect, and manage that workspace.

## Why this exists

The project exists to replace repetitive tmux session and window setup scripts with a clean open-source tool focused on Ubuntu development environments.

## Non-goals for MVP

- Full tmux workspace lifecycle management.
- Automatic TPM installation.
- Advanced pane layout orchestration.
- Cross-platform support beyond Ubuntu.
- Importing existing shell scripts.
