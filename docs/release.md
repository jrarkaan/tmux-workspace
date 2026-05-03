# Release Process

This document outlines the release scope, validation, and process for `twx`.

## Scope

Currently, `twx` targets an **Ubuntu-first** environment (20.04, 22.04, 24.04). The release pipeline generates `tar.gz` and `.deb` artifacts for Linux `amd64` and `arm64`.

*Note: macOS support, Homebrew taps, Snap, and Apt repository publishing are planned for future phases.*

## Local Validation

Before pushing any changes or triggering a release, ensure the codebase is clean:

```bash
go mod tidy
go fmt ./...
go test ./...
go vet ./...
go build -o twx .
```

### Version Injection Testing

The `version` command defaults to local development values:
```bash
./twx version
# twx version dev
# commit: none
# date: unknown
```

To test the multi-line injected version metadata, build with `ldflags`:
```bash
go build -ldflags "-X github.com/jrarkaan/tmux-workspace/cmd.version=0.1.0-test \
  -X github.com/jrarkaan/tmux-workspace/cmd.commit=testcommit \
  -X github.com/jrarkaan/tmux-workspace/cmd.date=2026-05-03T00:00:00Z" -o twx .

./twx version
# twx version 0.1.0-test
# commit: testcommit
# date: 2026-05-03T00:00:00Z
```

### GoReleaser Local Snapshot

If you have [GoReleaser](https://goreleaser.com/) installed locally, you can verify the configuration and simulate a release:

```bash
goreleaser check
goreleaser release --snapshot --clean
```

## Creating a Release

Releases are fully automated via GitHub Actions when a version tag (`v*`) is pushed.

1. Ensure your local `main` branch is up to date:
   ```bash
   git checkout main
   git pull origin main
   ```
2. Create an annotated tag for the release version:
   ```bash
   git tag v0.1.0
   ```
3. Push the tag to GitHub:
   ```bash
   git push origin v0.1.0
   ```

GitHub Actions will automatically run the validation checks (`fmt`, `test`, `vet`) and execute GoReleaser to publish the release.

## Expected Artifacts

The GitHub Release will include the following files:
- `twx_<version>_linux_amd64.tar.gz`
- `twx_<version>_linux_arm64.tar.gz`
- `twx_<version>_linux_amd64.deb`
- `twx_<version>_linux_arm64.deb`
- `checksums.txt`

## Testing .deb Packages Locally

After running a local `goreleaser release --snapshot --clean`, you can test the generated `.deb` package:

```bash
sudo dpkg -i dist/twx_<version>_linux_amd64.deb
twx version
```
