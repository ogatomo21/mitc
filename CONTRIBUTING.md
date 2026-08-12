# Contributing

Thanks for helping improve mitc.

## Development setup

Install Go 1.26 or later.

```powershell
go build ./...
go test ./...
go run . -p -y 2026 -u "Test User"
```

## Before opening a pull request

- Keep the CLI output compatible with the canonical MIT License text.
- Update `README.md` when user-facing behavior or options change.
- Add an entry to `CHANGELOG.md` for release-visible changes.
- Run `go build ./...` and `go test ./...` with no errors.

## Release versioning

Push a tag beginning with `v`, such as `v1.0.0`, to create the GitHub Release automatically.
The release title is set to `mitc v1.0.0`; the version without the leading `v` is passed to NSIS, while the original tag is embedded in each native binary.
