# mitc project guide

## Overview

`mitc` is a Go command-line tool that generates the canonical MIT License text.
It writes to `LICENSE` by default, can print to standard output, and supports saved or one-off copyright holder names.

## Repository layout

- `main.go`: executable entry point and release-version variable.
- `cli.go`: argument parsing, configuration handling, and license generation.
- `cli_test.go`: CLI behavior tests.
- `installer/mitc.nsi`: per-user Windows NSIS installer for the native x64 executable.
- `.github/workflows/ci.yml`: cross-platform Go build, tests, and smoke test.
- `.github/workflows/release.yml`: native release asset and installer publishing.
- `README.md`, `CHANGELOG.md`, `CONTRIBUTING.md`: public project documentation.

## CLI behavior

- `mitc`: write `LICENSE`; ask `[y/N]` before overwriting an existing file.
- `-p` / `--print`: write only to standard output.
- `-y` / `--year`: set the copyright year.
- `-u` / `--user`: set the copyright holder for one run.
- `--set-user`: save the default user to `.mitc.toml` in the user's home directory.
- `-f` / `--filename`: change the output file name.
- `-v` / `--version`: show the version; release builds inject the Git tag.

Keep the MIT License wording canonical and preserve the default `John Doe` user.

## Build and validation

Requires Go 1.26 or later.

```powershell
go build ./...
go test ./...
go run . -p -y 2026 -u "Test User"
```

Build a native Windows x64 executable:

```powershell
$env:CGO_ENABLED = '0'
$env:GOOS = 'windows'
$env:GOARCH = 'amd64'
go build -trimpath -ldflags '-s -w -X main.version=v1.0' -o .\publish-win-x64\mitc.exe .
```

Build the installer after installing NSIS:

```powershell
makensis /DPRODUCT_VERSION=1.0 installer\mitc.nsi
```

The installer places `mitc.exe` in `%LOCALAPPDATA%\mitc`, adds that directory to the user PATH, and removes it on uninstall. It does not require or bundle an additional runtime.

## Releases

Pushing a tag beginning with `v` (for example `v1.0.0`) triggers the release workflow.
The workflow creates the GitHub Release titled `mitc v1.0.0`, builds native binaries, and attaches:

- `mitc-win-x64.exe`
- `mitc-win-x86.exe`
- `mitc-win-setup.exe`
- `mitc-linux-x64`
- `mitc-macos-x64`
- `mitc-macos-arm64`

Do not commit generated `publish/`, `publish-win-x64/`, `artifacts/`, or installer EXEs.
