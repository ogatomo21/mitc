# Changelog

All notable changes to this project are documented in this file.

## Unreleased

- Add `mitc path add` and `mitc path remove` for managing the current user's Windows PATH from the directory containing `mitc.exe`.
- Make the NSIS installer use `mitc path add` and offer update/repair or uninstall when an existing installation is detected.

## v1.0

Initial public release.

- Generate an MIT License using the current year by default.
- Support a saved default copyright holder and one-off overrides.
- Write to `LICENSE` by default, with overwrite confirmation.
- Provide native standalone binaries for Windows, Linux, and macOS.
- Provide a Windows x64 NSIS installer with Japanese/English automatic UI selection.
