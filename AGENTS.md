# mitc project guide

## Overview

`mitc` is a .NET 10 command-line tool that generates the canonical MIT License text.
It saves to `LICENSE` by default, can print to standard output, and supports saved or one-off copyright holder names.

## Repository layout

- `Program.cs`: CLI argument parsing, configuration handling, and license text generation.
- `mitc.csproj`: .NET 10 project and single-file publishing defaults.
- `installer/mitc.nsi`: per-user Windows NSIS installer for the framework-dependent Any CPU distribution.
- `.github/workflows/ci.yml`: cross-platform build and CLI smoke test.
- `.github/workflows/release.yml`: release asset publishing.
- `README.md`, `CHANGELOG.md`, `CONTRIBUTING.md`: public project documentation.

## CLI behavior

- `mitc`: write `LICENSE`; ask `[y/N]` before overwriting an existing file.
- `-p` / `--print`: write only to standard output.
- `-y` / `--year`: set the copyright year.
- `-u` / `--user`: set the copyright holder for one run.
- `--set-user`: save the default user to `.mitc.toml` in the user's home directory.
- `-f` / `--filename`: change the output file name.
- `-v` / `--version`: show the informational version, such as `v1.0`.

Keep the MIT License wording canonical and preserve the default `John Doe` user.

## Build and validation

Requires the .NET 10 SDK specified by `global.json`.

```powershell
dotnet build mitc.csproj --configuration Release
dotnet run --project mitc.csproj --configuration Release -- --print --year 2026 --user "Test User"
```

Build a self-contained single Windows executable:

```powershell
dotnet publish -c Release -r win-x64 --self-contained true -o .\publish\win-x64
```

Build the Any CPU Windows installer after installing NSIS:

```powershell
dotnet publish -c Release --self-contained false -p:PublishSingleFile=false -o .\publish-anycpu
makensis /DPRODUCT_VERSION=1.0 installer\mitc.nsi
```

The installer intentionally does not bundle .NET. It installs the Any CPU files under `%LOCALAPPDATA%\mitc`, adds that directory to the user PATH, and removes it on uninstall.

## Releases

Publishing a GitHub Release from a tag beginning with `v` (for example `v1.0.0`) triggers the release workflow.
The workflow sets the Release title to `mitc v1.0.0`, uses `1.0.0` for .NET and NSIS product versions, and attaches:

- `mitc-win-x64.exe`
- `mitc-win-x86.exe`
- `mitc-win-setup.exe`
- `mitc-linux-x64`
- `mitc-macos-x64`
- `mitc-macos-arm64`

Do not commit generated `bin/`, `obj/`, `publish/`, `publish-anycpu/`, or installer EXEs.
