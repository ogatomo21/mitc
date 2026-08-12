# mitc

`mitc` is a small command-line tool that creates an MIT License file. It uses the current year by default and remembers an optional default copyright holder.

> 日本語の説明はこのREADMEに続きます。

## Install

Download an asset from GitHub Releases.

- `mitc-win-x64.exe`: self-contained Windows x64 executable.
- `mitc-win-x86.exe`: self-contained Windows x86 executable.
- `mitc-win-setup.exe`: Windows installer for x86, x64, and Arm64. Requires the .NET 10 Runtime.
- `mitc-linux-x64`: self-contained Linux x64 executable.
- `mitc-macos-x64` / `mitc-macos-arm64`: self-contained macOS executables.

The self-contained executables do not require a separate .NET Runtime installation.

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

Running `mitc` creates `LICENSE` in the current directory. If the file already exists, only `y` confirms overwriting it.

```powershell
mitc
mitc -y 2026 -u "Tomoya Ogawa"
mitc -f LICENSE.txt
mitc -p
mitc --set-user "Tomoya Ogawa"
```

The default user is `John Doe`. `--set-user` stores the value in `~/.mitc.toml` (or `%USERPROFILE%\.mitc.toml` on Windows).

```toml
user = "Tomoya Ogawa"
```

## Development

Requires the .NET 10 SDK.

```powershell
dotnet build -c Release
dotnet run -- --print -y 2026 -u "Test User"
```

To build a self-contained Windows x64 executable:

```powershell
dotnet publish -c Release -r win-x64 --self-contained true -o .\publish\win-x64
```

To build the Windows Any CPU installer, install NSIS and run:

```powershell
dotnet publish -c Release --self-contained false -p:PublishSingleFile=false -o .\publish-anycpu
makensis /DPRODUCT_VERSION=1.0 installer\mitc.nsi
```

## Releases

Pushing a `v`-prefixed tag, such as `v1.0.0`, creates a GitHub Release and attaches all platform assets automatically. The workflow sets the Release title to `mitc v1.0.0` and uses `1.0.0` for .NET and NSIS product versions.

See [CHANGELOG.md](CHANGELOG.md) for release history and [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

---

## 日本語

`mitc` はMIT Licenseを生成するC#製CLIツールです。現在年を既定で使用し、著作権者名の既定値も保存できます。

### インストール

GitHub Releasesから利用環境に合うファイルをダウンロードしてください。

- `mitc-win-x64.exe`：Windows x64向け、.NET Runtime不要の単一EXE。
- `mitc-win-x86.exe`：Windows x86向け、.NET Runtime不要の単一EXE。
- `mitc-win-setup.exe`：Windows x86/x64/Arm64共通のセットアップ。 .NET 10 Runtimeが必要。
- `mitc-linux-x64`：Linux x64向け、.NET Runtime不要の単一実行ファイル。
- `mitc-macos-x64` / `mitc-macos-arm64`：macOS向け、.NET Runtime不要の単一実行ファイル。

### 使い方

引数なしではカレントディレクトリの`LICENSE`へ保存します。既存の`LICENSE`は`y`を入力した場合だけ上書きします。

```powershell
mitc
mitc -y 2026 -u "Tomoya Ogawa"
mitc -f LICENSE.txt
mitc -p
mitc --set-user "Tomoya Ogawa"
```

| オプション | 内容 |
| --- | --- |
| `-y`, `--year <YEAR>` | 年を指定します（既定：現在年）。 |
| `-u`, `--user <NAME>` | 今回だけ著作権者名を指定します。 |
| `--set-user <NAME>` | 既定の著作権者名を保存します。 |
| `-f`, `--filename <FILE>` | 出力ファイル名を指定します。 |
| `-p`, `--print` | ファイルへ保存せず標準出力へ表示します。 |
| `-h`, `--help` | ヘルプを表示します。 |
| `-v`, `--version` | バージョンを表示します。 |

既定ユーザーは`John Doe`です。`--set-user`で保存すると、ユーザーホームの`.mitc.toml`に設定されます。

```toml
user = "Tomoya Ogawa"
```

### 開発とリリース

.NET 10 SDKが必要です。

```powershell
dotnet build -c Release
dotnet run -- --print -y 2026 -u "Test User"
```

`v1.0.0`のようなタグをpushすると、GitHub Releaseを作成して全プラットフォーム向け成果物を自動添付します。Release名は`mitc v1.0.0`へ自動設定され、.NETとNSISには`1.0.0`が渡されます。

更新履歴は[CHANGELOG.md](CHANGELOG.md)、貢献方法は[CONTRIBUTING.md](CONTRIBUTING.md)を参照してください。
