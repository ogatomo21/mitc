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
!include "FileFunc.nsh"
!include "StrFunc.nsh"
!include "WinMessages.nsh"

${Using:StrFunc} StrStr
${Using:StrFunc} UnStrRep

Var ExistingPath
Var UpdatedPath

LangString InstallSectionName ${LANG_ENGLISH} "Install mitc"
LangString InstallSectionName ${LANG_JAPANESE} "mitc をインストール"

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

  ReadRegStr $ExistingPath HKCU "Environment" "Path"
  ${StrStr} $0 $ExistingPath "$INSTDIR"
  ${If} $0 == ""
    ${If} $ExistingPath == ""
      StrCpy $UpdatedPath "$INSTDIR"
    ${Else}
      StrCpy $UpdatedPath "$ExistingPath;$INSTDIR"
    ${EndIf}
    WriteRegExpandStr HKCU "Environment" "Path" $UpdatedPath
    SendMessage ${HWND_BROADCAST} ${WM_SETTINGCHANGE} 0 "STR:Environment" /TIMEOUT=5000
  ${EndIf}
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\mitc.exe"
  Delete "$INSTDIR\uninstall.exe"
  RMDir "$INSTDIR"

  ReadRegStr $ExistingPath HKCU "Environment" "Path"
  ${UnStrRep} $UpdatedPath $ExistingPath ";$INSTDIR" ""
  ${UnStrRep} $UpdatedPath $UpdatedPath "$INSTDIR;" ""
  ${If} $UpdatedPath == "$INSTDIR"
    DeleteRegValue HKCU "Environment" "Path"
  ${Else}
    WriteRegExpandStr HKCU "Environment" "Path" $UpdatedPath
  ${EndIf}
  SendMessage ${HWND_BROADCAST} ${WM_SETTINGCHANGE} 0 "STR:Environment" /TIMEOUT=5000

  DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\mitc"
SectionEnd
