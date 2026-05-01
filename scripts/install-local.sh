#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="$HOME/bin"
BASHRC="$HOME/.bashrc"
PATH_LINE='export PATH="$HOME/bin:$PATH"'

cd "$ROOT_DIR"

echo "Building twx..."
go build -o twx .

mkdir -p "$BIN_DIR"
cp twx "$BIN_DIR/twx"
chmod +x "$BIN_DIR/twx"

touch "$BASHRC"
if ! grep -Fxq "$PATH_LINE" "$BASHRC"; then
  printf '\n%s\n' "$PATH_LINE" >> "$BASHRC"
fi

echo
echo "Installed twx to $BIN_DIR/twx"
echo
echo "Next steps:"
echo "  source ~/.bashrc"
echo "  twx --help"
