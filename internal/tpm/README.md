# tpm

This package handles the detection and installation of the Tmux Plugin Manager (TPM).

It checks for `git` and `tmux` presence, evaluates `~/.tmux.conf` for existing plugins, and can install TPM by cloning it and safely appending a managed block to `~/.tmux.conf`.
