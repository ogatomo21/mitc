# -*- coding: utf-8 -*-
Unicode True
ManifestDPIAware True
RequestExecutionLevel user
SetCompressor /SOLID lzma

# WindowsのUI言語に合わせて自動選択（未対応言語は英語）
LoadLanguageFile "${NSISDIR}\Contrib\Language files\English.nlf"
LoadLanguageFile "${NSISDIR}\Contrib\Language files\Japanese.nlf"

!ifndef PRODUCT_VERSION
  !define PRODUCT_VERSION "1.0"
!endif

!define PRODUCT_NAME "mitc"
!define PUBLISH_DIR "..\publish-win-x64"
!define INSTALL_DIR "$LOCALAPPDATA\mitc"

Name "${PRODUCT_NAME} ${PRODUCT_VERSION}"
OutFile "mitc-win-setup.exe"
InstallDir "${INSTALL_DIR}"
ShowInstDetails hide
ShowUninstDetails hide

!include "LogicLib.nsh"
!include "WinMessages.nsh"

Var PathHelperResult
Var ExistingInstallDir
Var ExistingUninstaller

LangString InstallSectionName ${LANG_ENGLISH} "Install mitc"
LangString InstallSectionName ${LANG_JAPANESE} "mitc をインストール"
LangString ExistingInstall ${LANG_ENGLISH} "mitc is already installed. Choose Yes to update/repair, No to uninstall, or Cancel to exit."
LangString ExistingInstall ${LANG_JAPANESE} "mitc はすでにインストールされています。はいで更新（修復）、いいえでアンインストール、キャンセルで終了します。"
LangString UninstallerFailed ${LANG_ENGLISH} "The existing mitc uninstaller could not be started successfully."
LangString UninstallerFailed ${LANG_JAPANESE} "既存の mitc アンインストーラーを正常に起動できませんでした。"
LangString PathAddFailed ${LANG_ENGLISH} "mitc was installed, but its user PATH could not be updated or reloaded. Existing PATH entries were preserved."
LangString PathAddFailed ${LANG_JAPANESE} "mitc はインストールされましたが、ユーザー PATH を更新または再読み込みできませんでした。既存の PATH エントリは保持されています。"
LangString PathRemoveFailed ${LANG_ENGLISH} "mitc was uninstalled, but its PATH entry could not be removed or reloaded. Other PATH entries were preserved."
LangString PathRemoveFailed ${LANG_JAPANESE} "mitc はアンインストールされましたが、PATH エントリを削除または再読み込みできませんでした。他の PATH エントリは保持されています。"

Function .onInit
  ReadRegStr $ExistingInstallDir HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "InstallLocation"
  ReadRegStr $ExistingUninstaller HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "UninstallString"
  StrCmp $ExistingUninstaller "" 0 ExistingInstallFound
  IfFileExists "${INSTALL_DIR}\uninstall.exe" ExistingInstallFound NoExistingInstall

ExistingInstallFound:
  StrCmp $ExistingInstallDir "" UseDefaultInstallDir
  StrCpy $INSTDIR $ExistingInstallDir
  Goto ChooseExistingInstallAction

UseDefaultInstallDir:
  StrCpy $ExistingInstallDir "${INSTALL_DIR}"
  StrCpy $INSTDIR $ExistingInstallDir

ChooseExistingInstallAction:
  MessageBox MB_YESNOCANCEL|MB_ICONQUESTION "$(ExistingInstall)" IDYES ExistingInstallUpdate IDNO ExistingInstallUninstall IDCANCEL ExistingInstallCancel

ExistingInstallUpdate:
  Return

ExistingInstallUninstall:
  StrCmp $ExistingUninstaller "" UseDefaultUninstaller
  Goto RunExistingUninstaller

UseDefaultUninstaller:
  StrCpy $ExistingUninstaller '"$ExistingInstallDir\uninstall.exe"'

RunExistingUninstaller:
  ExecWait $ExistingUninstaller $PathHelperResult
  StrCmp $PathHelperResult "0" ExistingInstallUninstallDone
  MessageBox MB_ICONEXCLAMATION|MB_OK "$(UninstallerFailed)"

ExistingInstallUninstallDone:
  Quit

ExistingInstallCancel:
  Quit

NoExistingInstall:
FunctionEnd

# The uninstaller helper below is retained so uninstallers from releases before
# the path subcommands were available can still remove their PATH entry.
Function un.WritePathHelper
  InitPluginsDir
  FileOpen $0 "$PLUGINSDIR\mitc-path.ps1" w
  FileWrite $0 "param([ValidateSet('Add','Remove')][string]$$Action,[Parameter(Mandatory=$$true)][string]$$Entry)$\r$\n"
  FileWrite $0 "$$ErrorActionPreference='Stop'$\r$\n"
  FileWrite $0 "$$key=[Microsoft.Win32.Registry]::CurrentUser.CreateSubKey('Environment')$\r$\n"
  FileWrite $0 "try {$\r$\n"
  FileWrite $0 "  $$path=[string]$$key.GetValue('Path','',[Microsoft.Win32.RegistryValueOptions]::DoNotExpandEnvironmentNames)$\r$\n"
  FileWrite $0 "  $$norm=[Environment]::ExpandEnvironmentVariables($$Entry).TrimEnd('\')$\r$\n"
  FileWrite $0 "  $$parts=@($$path -split ';')$\r$\n"
  FileWrite $0 "  $$match={ param($$x) [string]::Equals([Environment]::ExpandEnvironmentVariables($$x).TrimEnd('\'),$$norm,[StringComparison]::OrdinalIgnoreCase) }$\r$\n"
  FileWrite $0 "  $$exists=@($$parts | Where-Object { & $$match $$_ }).Count -gt 0$\r$\n"
  FileWrite $0 "  if($$Action -eq 'Add' -and -not $$exists){$\r$\n"
  FileWrite $0 "    if([string]::IsNullOrEmpty($$path)){ $$new=$$Entry } elseif($$path.EndsWith(';')){ $$new=$$path+$$Entry } else { $$new=$$path+';'+$$Entry }$\r$\n"
  FileWrite $0 "    $$key.SetValue('Path',$$new,[Microsoft.Win32.RegistryValueKind]::ExpandString)$\r$\n"
  FileWrite $0 "  } elseif($$Action -eq 'Remove' -and $$exists){$\r$\n"
  FileWrite $0 "    $$new=@($$parts | Where-Object { -not (& $$match $$_) }) -join ';'$\r$\n"
  FileWrite $0 "    $$key.SetValue('Path',$$new,[Microsoft.Win32.RegistryValueKind]::ExpandString)$\r$\n"
  FileWrite $0 "  }$\r$\n"
  FileWrite $0 "} finally { if($$null -ne $$key){ $$key.Dispose() } }$\r$\n"
  FileClose $0
FunctionEnd

# ページ
Page directory
Page instfiles

Section "$(InstallSectionName)" SecInstall
  SetOutPath "$INSTDIR"
  SetOverwrite on
  File "${PUBLISH_DIR}\mitc.exe"

  WriteUninstaller "$INSTDIR\uninstall.exe"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "DisplayName" "mitc"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "DisplayVersion" "${PRODUCT_VERSION}"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "InstallLocation" "$INSTDIR"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "Publisher" "Tomoya Ogawa"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "UninstallString" '"$INSTDIR\uninstall.exe"'
  WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "NoModify" 1
  WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "NoRepair" 1

  # Let the installed executable manage its own directory in the user PATH.
  ExecWait '"$INSTDIR\mitc.exe" path add' $PathHelperResult
  ${If} $PathHelperResult != 0
    MessageBox MB_ICONEXCLAMATION|MB_OK "$(PathAddFailed)"
  ${Else}
    SendMessage ${HWND_BROADCAST} ${WM_SETTINGCHANGE} 0 "STR:Environment" /TIMEOUT=5000
  ${EndIf}
SectionEnd

Section "Uninstall"
  # Remove only the exact mitc directory from PATH before deleting the files.
  # Never delete HKCU\Environment\Path itself, even if it becomes empty.
  Call un.WritePathHelper
  ExecWait '"$SYSDIR\WindowsPowerShell\v1.0\powershell.exe" -NoLogo -NoProfile -NonInteractive -ExecutionPolicy Bypass -File "$PLUGINSDIR\mitc-path.ps1" -Action Remove -Entry "$INSTDIR"' $PathHelperResult
  ${If} $PathHelperResult != 0
    MessageBox MB_ICONEXCLAMATION|MB_OK "$(PathRemoveFailed)"
  ${Else}
    SendMessage ${HWND_BROADCAST} ${WM_SETTINGCHANGE} 0 "STR:Environment" /TIMEOUT=5000
  ${EndIf}

  Delete "$INSTDIR\mitc.exe"
  Delete "$INSTDIR\uninstall.exe"
  RMDir "$INSTDIR"

  DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc"
SectionEnd
