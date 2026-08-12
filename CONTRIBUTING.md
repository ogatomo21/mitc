# Contributing

Thanks for helping improve mitc.

## Development setup

Install the .NET 10 SDK. The supported SDK policy is in `global.json`.

```powershell
dotnet build -c Release
dotnet run -- --print -y 2026 -u "Test User"
```

## Before opening a pull request

- Keep the CLI output compatible with the canonical MIT License text.
- Update `README.md` when user-facing behavior or options change.
- Add an entry to `CHANGELOG.md` for release-visible changes.
- Run `dotnet build -c Release` with no warnings or errors.

## Release versioning

Publish GitHub Releases from a tag beginning with `v`, such as `v1.0.0`.
The release workflow removes the `v` for .NET and NSIS product versions, and changes the Release title to `mitc v1.0.0`.
