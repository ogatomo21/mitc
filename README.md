# mitc

[![CI](https://github.com/ogatomo21/mitc/actions/workflows/ci.yml/badge.svg)](https://github.com/ogatomo21/mitc/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`mitc` is a tiny, native command-line tool for creating the canonical MIT License. It uses the current year by default and can remember your copyright holder name.

> 日本語の説明は[こちら](#日本語)です。

## Features

- Native standalone binaries for Windows, Linux, and macOS — no additional runtime required.
- Writes a `LICENSE` file by default and asks before overwriting one.
- Supports a saved default copyright holder, as well as one-off overrides.
- Prints the license to standard output for scripts and pipelines.

## Install

Download the appropriate file from [GitHub Releases](https://github.com/ogatomo21/mitc/releases).

| File | Platform | Notes |
| --- | --- | --- |
| `mitc-win-setup.exe` | Windows x64 | Recommended. Adds `mitc` to the current user's `PATH`; Japanese or English is selected automatically. |
| `mitc-win-x64.exe` | Windows x64 | Portable executable. |
| `mitc-win-x86.exe` | Windows x86 | Portable executable. |
| `mitc-linux-x64` | Linux x64 | Make executable with `chmod +x mitc-linux-x64`. |
| `mitc-macos-x64` | macOS Intel | Make executable with `chmod +x mitc-macos-x64`. |
| `mitc-macos-arm64` | macOS Apple silicon | Make executable with `chmod +x mitc-macos-arm64`. |

For a portable executable, rename it to `mitc` (or `mitc.exe` on Windows) and place it in a directory on your `PATH`.

## Usage

```text
mitc [options]

-y, --year <YEAR>      Set the copyright year (default: current year)
-u, --user <NAME>      Use a copyright holder for this run only
    --set-user <NAME>  Save the default copyright holder
-f, --filename <FILE>  Change the output file name
-p, --print            Write to standard output instead of a file
-h, --help             Show help
-v, --version          Show the version
```

```powershell
# Create LICENSE using the current year and saved/default user.
mitc

# Create a license for one project without changing the saved user.
mitc -y 2026 -u "Tomoya Ogawa"

# Save the default copyright holder in ~/.mitc.toml.
mitc --set-user "Tomoya Ogawa"

# Print instead of writing a file.
mitc -p
```

The default user is `John Doe`. `--set-user` saves the value to `~/.mitc.toml` (or `%USERPROFILE%\.mitc.toml` on Windows).

```toml
user = "Tomoya Ogawa"
```

## Development

Go 1.26 or later is required.

```powershell
go build ./...
go test ./...
go vet ./...
go run . -p -y 2026 -u "Test User"
```

Build a Windows x64 binary and its installer:

```powershell
$env:CGO_ENABLED = '0'
$env:GOOS = 'windows'
$env:GOARCH = 'amd64'
go build -trimpath -ldflags '-s -w -X main.version=v1.0' -o .\publish-win-x64\mitc.exe .
makensis /DPRODUCT_VERSION=1.0 installer\mitc.nsi
```

## Releases

Pushing a `v`-prefixed tag, such as `v1.0.0`, creates or updates the corresponding GitHub Release. Its title becomes `mitc v1.0.0` and the workflow attaches every platform binary and the Windows installer.

See [CHANGELOG.md](CHANGELOG.md) for release history, [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines, and [LICENSE](LICENSE) for license terms.

---

## 日本語

`mitc` は、標準的なMIT Licenseを生成する小さなネイティブCLIツールです。既定では現在年を使い、著作権者名を保存できます。

### 特長

- Windows、Linux、macOS向けの単体バイナリ。追加ランタイムは不要です。
- 既定で`LICENSE`へ保存し、既存ファイルは確認してから上書きします。
- 著作権者名の保存と、その場限りの上書き指定に対応しています。
- `--print`で標準出力へ出せるため、スクリプトからも利用できます。

### インストール

[GitHub Releases](https://github.com/ogatomo21/mitc/releases)から環境に合うファイルをダウンロードしてください。

| ファイル | 対象 | 補足 |
| --- | --- | --- |
| `mitc-win-setup.exe` | Windows x64 | 推奨。`mitc`を現在のユーザーの`PATH`へ追加します。画面は日本語・英語を自動選択します。 |
| `mitc-win-x64.exe` | Windows x64 | 持ち運び用EXEです。 |
| `mitc-win-x86.exe` | Windows x86 | 持ち運び用EXEです。 |
| `mitc-linux-x64` | Linux x64 | `chmod +x mitc-linux-x64`を実行してから使います。 |
| `mitc-macos-x64` | Intel Mac | `chmod +x mitc-macos-x64`を実行してから使います。 |
| `mitc-macos-arm64` | AppleシリコンMac | `chmod +x mitc-macos-arm64`を実行してから使います。 |

持ち運び版は`mitc`（Windowsでは`mitc.exe`）へリネームして、`PATH`上のフォルダーへ置くと便利です。

### 使い方

引数なしではカレントディレクトリの`LICENSE`へ保存します。既存の`LICENSE`は`y`を入力した場合だけ上書きします。

```powershell
# 現在年と保存済み（または既定）のユーザーでLICENSEを作成
mitc

# この実行だけ著作権者名と年を指定
mitc -y 2026 -u "Tomoya Ogawa"

# 既定の著作権者名を保存
mitc --set-user "Tomoya Ogawa"

# ファイルを作らず画面へ表示
mitc -p
```

既定ユーザーは`John Doe`です。`--set-user`で保存すると、ユーザーホームの`.mitc.toml`に設定されます。

```toml
user = "Tomoya Ogawa"
```

### 開発とリリース

Go 1.26以降が必要です。

```powershell
go build ./...
go test ./...
go vet ./...
go run . -p -y 2026 -u "Test User"
```

`v1.0.0`のようなタグをpushすると、GitHub Releaseを自動で作成・更新します。リリース名は`mitc v1.0.0`になり、各OS向けバイナリとWindowsセットアップが添付されます。
