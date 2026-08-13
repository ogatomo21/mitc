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

LangString InstallSectionName ${LANG_ENGLISH} "Install mitc"
LangString InstallSectionName ${LANG_JAPANESE} "mitc をインストール"
LangString PathAddFailed ${LANG_ENGLISH} "mitc was installed, but the user PATH could not be updated. Your existing PATH was left unchanged."
LangString PathAddFailed ${LANG_JAPANESE} "mitc はインストールされましたが、ユーザー PATH を更新できませんでした。既存の PATH は変更していません。"
LangString PathRemoveFailed ${LANG_ENGLISH} "mitc was uninstalled, but its PATH entry could not be removed. Other PATH entries were not changed."
LangString PathRemoveFailed ${LANG_JAPANESE} "mitc はアンインストールされましたが、PATH エントリを削除できませんでした。他の PATH エントリは変更していません。"

# -----------------------------------------------------------------------------
# PATH helper
#
# IMPORTANT:
#   Do NOT use ReadRegStr/WriteRegExpandStr to round-trip PATH here.
#   Normal NSIS variables are limited to NSIS_MAX_STRLEN (1024 by default).
#   A long PATH can therefore be read as an empty string/truncated and then
#   accidentally overwritten.
#
#   Instead, generate a small PowerShell helper and let .NET read/write the
#   registry value directly. The PATH value itself never enters an NSIS string.
# -----------------------------------------------------------------------------

Function WritePathHelper
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
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "Publisher" "Tomoya Ogawa"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "UninstallString" '"$INSTDIR\uninstall.exe"'
  WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "NoModify" 1
  WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc" "NoRepair" 1

  # Add only this install directory to the user PATH.
  # The existing PATH is never stored in an NSIS variable.
  Call WritePathHelper
  ExecWait '"$SYSDIR\WindowsPowerShell\v1.0\powershell.exe" -NoLogo -NoProfile -NonInteractive -ExecutionPolicy Bypass -File "$PLUGINSDIR\mitc-path.ps1" -Action Add -Entry "$INSTDIR"' $PathHelperResult
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
