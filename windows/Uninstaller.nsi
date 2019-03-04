!include MUI2.nsh

;--------------------------------
;Perform Machine-level install, if possible

!define MULTIUSER_EXECUTIONLEVEL Highest
;Add support for command-line args that let uninstaller know whether to
;uninstall machine- or user installation:
!define MULTIUSER_INSTALLMODE_COMMANDLINE
!include MultiUser.nsh
!include LogicLib.nsh

Function .onInit
  !insertmacro MULTIUSER_INIT
FunctionEnd

Function un.onInit
  !insertmacro MULTIUSER_UNINIT
FunctionEnd

;--------------------------------
;General

  Name "Aminal"

;--------------------------------
;Pages

  !insertmacro MUI_UNPAGE_CONFIRM
  !insertmacro MUI_UNPAGE_INSTFILES

;--------------------------------
;Languages

  !insertmacro MUI_LANGUAGE "English"

;--------------------------------
;Installer Sections

Section
  SetOutPath "$InstDir"
  WriteUninstaller "$InstDir\uninstall.exe"
SectionEnd

;--------------------------------
;Uninstaller Section

!define UNINST_KEY \
  "Software\Microsoft\Windows\CurrentVersion\Uninstall\Aminal"
!define ROOT_KEY "Software\Aminal"
!define UPDATE_KEY \
  "${ROOT_KEY}\Update\Clients\{35B0CF1E-FBB0-486F-A1DA-BE3A41DDC780}"

Section "Uninstall"

  RMDir /r "$InstDir\Versions"
  Delete "$InstDir\Aminal.exe"
  Delete "$InstDir\uninstall.exe"
  ;Omaha leaves this directory behind. Delete if empty:
  RMDir "$InstDir\CrashReports"
  RMDir "$InstDir"
  Delete "$SMPROGRAMS\Aminal.lnk"
  DeleteRegKey SHCTX "${UNINST_KEY}"
  DeleteRegKey SHCTX "${UPDATE_KEY}"
  DeleteRegKey /ifempty SHCTX "${ROOT_KEY}\Update\Clients"
  ;Try to speed up uninstall of Omaha:
  DeleteRegValue SHCTX "${ROOT_KEY}\Update" "LastChecked"
  DeleteRegKey /ifempty SHCTX "${ROOT_KEY}\Update"
  WriteRegStr SHCTX "${ROOT_KEY}" "" ""
  DeleteRegKey /ifempty SHCTX "${ROOT_KEY}"

SectionEnd