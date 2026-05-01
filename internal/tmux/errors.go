package tmux

import "errors"

var (
	ErrTmuxNotInstalled = errors.New("tmux is not installed or not in PATH")
	ErrSessionNotFound  = errors.New("tmux session not found")
)
