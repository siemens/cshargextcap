; ContainerShark external capture plugin installer for Windows. post.sh
; defines VERSION, FILEVERSION, ARCH, BINARY, LICENSE, and OUTFILE.

Unicode true
ManifestDPIAware true
RequestExecutionLevel admin
SetCompressor /SOLID lzma

!include "MUI2.nsh"
!include "LogicLib.nsh"
!include "FileFunc.nsh"
!include "x64.nsh"

!define /date YEAR "%Y"
!define PRODUCT "ContainerShark Extcap Plugin"
!define PUBLISHER "Edgeshark Project"
!define DESCRIPTION "ContainerShark external capture plugin for capturing container network traffic"
!define UNINSTKEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\ContainerShark"
!define URLKEY "Software\Classes\Packetflix"
; extcap binary name used by installers before 0.11.
!define LEGACYBINARY "cshargextcap-amd64.exe"

Name "ContainerShark External Capture Plugin ${VERSION} for Wireshark"
; The default captions contain the full name including the version and get clipped.
Caption "${PRODUCT} Setup"
UninstallCaption "${PRODUCT} Uninstall"
OutFile "${OUTFILE}"
BrandingText "ContainerShark Installer"

VIProductVersion "${FILEVERSION}"
VIFileVersion "${FILEVERSION}"
VIAddVersionKey "ProductName" "${PRODUCT}"
VIAddVersionKey "ProductVersion" "${VERSION}"
VIAddVersionKey "FileVersion" "${VERSION}"
VIAddVersionKey "FileDescription" "${DESCRIPTION}"
VIAddVersionKey "CompanyName" "${PUBLISHER}"
VIAddVersionKey "LegalCopyright" "Copyright © Siemens AG ${YEAR}"

; Overridden in .onInit by the location of an installed Wireshark.
InstallDir "$PROGRAMFILES64\Wireshark\extcap"

Function .onInit
    SetRegView 64
!if "${ARCH}" == "arm64"
    ${IfNot} ${IsNativeARM64}
        MessageBox MB_ICONSTOP "This installer requires Windows on Arm64." /SD IDOK
        Abort
    ${EndIf}
!else
    ${IfNot} ${RunningX64}
        MessageBox MB_ICONSTOP "This installer requires 64-bit Windows." /SD IDOK
        Abort
    ${EndIf}
!endif
    UserInfo::GetAccountType
    Pop $0
    ${If} $0 != "admin"
        MessageBox MB_ICONSTOP "Installing requires administrator rights." /SD IDOK
        Abort
    ${EndIf}

    ; Respect an explicit /D=... installation directory.
    ${If} $INSTDIR != "$PROGRAMFILES64\Wireshark\extcap"
        Return
    ${EndIf}
    ReadRegStr $0 HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Wireshark" "InstallLocation"
    ${If} $0 == ""
        ReadRegStr $0 HKLM "Software\Microsoft\Windows\CurrentVersion\App Paths\Wireshark.exe" ""
        ${GetParent} "$0" $0
    ${EndIf}
    ${If} $0 != ""
    ${AndIf} ${FileExists} "$0\Wireshark.exe"
        StrCpy $INSTDIR "$0\extcap"
    ${ElseIfNot} ${FileExists} "$PROGRAMFILES64\Wireshark\Wireshark.exe"
        MessageBox MB_ICONEXCLAMATION|MB_OKCANCEL "Wireshark doesn't seem to be installed.$\r$\n$\r$\nInstall the plugin anyway?" /SD IDOK IDOK +2
        Abort
    ${EndIf}
FunctionEnd

!define MUI_ICON "containershark.ico"
!define MUI_UNICON "containershark.ico"
!define MUI_ABORTWARNING
!define MUI_WELCOMEFINISHPAGE_BITMAP "welcomefinish.bmp"
!define MUI_UNWELCOMEFINISHPAGE_BITMAP "welcomefinish.bmp"

!define MUI_WELCOMEPAGE_TITLE "\
Welcome to the ContainerShark Extcap Plugin for Wireshark Setup"
!define MUI_WELCOMEPAGE_TEXT "\
Setup will guide you through the installation of the \
ContainerShark external capture plugin \
(version ${VERSION}) for Wireshark.$\r$\n\
$\r$\n\
This will install the capture plugin to capture from:$\r$\n\
– Industrial Edge and Docker hosts$\r$\n\
– packetflix:// capture URL handler$\r$\n\
$\r$\n\
Click Next to continue.\
"
!define MUI_DIRECTORYPAGE_TEXT_TOP "\
Setup will install the plugin into the following extcap folder of Wireshark. \
Wireshark needs to be restarted to pick up the plugin."

!define MUI_FINISHPAGE_TITLE "Completing the ContainerShark Extcap Plugin Setup"
!define MUI_FINISHPAGE_NOAUTOCLOSE
!define MUI_FINISHPAGE_RUN "$INSTDIR\..\Wireshark.exe"
!define MUI_FINISHPAGE_RUN_TEXT "Run Wireshark"

; Next stays disabled until the user explicitly accepts the license.
!define MUI_LICENSEPAGE_RADIOBUTTONS

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "${LICENSE}"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

; The default titles contain the full name including the version and get clipped.
!define MUI_WELCOMEPAGE_TITLE "Uninstall the ContainerShark Extcap Plugin"
!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!define MUI_FINISHPAGE_TITLE "Completing the ContainerShark Extcap Plugin Uninstall"
!insertmacro MUI_UNPAGE_FINISH

!insertmacro MUI_LANGUAGE "English"

Section "Install"
    SetRegView 64

    ; Remove leftovers from installers before 0.11, which installed per user
    ; and under a different binary name.
    Delete "$INSTDIR\${LEGACYBINARY}"
    Delete "$INSTDIR\containershark.ico"
    DeleteRegKey HKCU "${UNINSTKEY}"
    DeleteRegKey HKCU "${URLKEY}"
    DeleteRegKey HKCU "Software\containershark"

    SetOutPath "$INSTDIR"
    File "/oname=cshargextcap.exe" "${BINARY}"
    ; Wireshark runs any executable in its extcap folder, so the uninstaller
    ; must go into a subfolder.
    SetOutPath "$INSTDIR\containershark"
    File "containershark.ico"
    File "/oname=LICENSE.txt" "${LICENSE}"
    WriteUninstaller "$INSTDIR\containershark\containershark-uninstall.exe"

    ${GetSize} "$INSTDIR" "/M=cshargextcap.exe /S=0K /G=0" $0 $1 $2
    WriteRegStr HKLM "${UNINSTKEY}" "DisplayName" "ContainerShark ExtCap Plugin for Wireshark"
    WriteRegStr HKLM "${UNINSTKEY}" "DisplayVersion" "${VERSION}"
    WriteRegStr HKLM "${UNINSTKEY}" "Publisher" "${PUBLISHER}"
    WriteRegStr HKLM "${UNINSTKEY}" "URLInfoAbout" "https://github.com/siemens/cshargextcap"
    WriteRegStr HKLM "${UNINSTKEY}" "DisplayIcon" "$\"$INSTDIR\containershark\containershark.ico$\""
    WriteRegStr HKLM "${UNINSTKEY}" "InstallLocation" "$INSTDIR"
    WriteRegStr HKLM "${UNINSTKEY}" "UninstallString" "$\"$INSTDIR\containershark\containershark-uninstall.exe$\""
    WriteRegStr HKLM "${UNINSTKEY}" "QuietUninstallString" "$\"$INSTDIR\containershark\containershark-uninstall.exe$\" /S"
    WriteRegDWORD HKLM "${UNINSTKEY}" "EstimatedSize" $0
    WriteRegDWORD HKLM "${UNINSTKEY}" "NoModify" 1
    WriteRegDWORD HKLM "${UNINSTKEY}" "NoRepair" 1

    ${GetParent} "$INSTDIR" $0
    WriteRegStr HKLM "${URLKEY}" "" "URL:Packetflix Protocol"
    WriteRegStr HKLM "${URLKEY}" "URL Protocol" ""
    WriteRegStr HKLM "${URLKEY}\DefaultIcon" "" "$\"$INSTDIR\containershark\containershark.ico$\""
    WriteRegStr HKLM "${URLKEY}\shell\open\command" "" '"$0\Wireshark.exe" -k -i packetflix -o "extcap.packetflix.url:%1"'
SectionEnd

Section "Uninstall"
    SetRegView 64
    ; $INSTDIR is the uninstaller's subfolder, but we need the extcap folder,
    ; which we must not remove.
    ${GetParent} "$INSTDIR" $INSTDIR
    Delete "$INSTDIR\cshargextcap.exe"
    Delete "$INSTDIR\${LEGACYBINARY}"
    Delete "$INSTDIR\containershark\containershark.ico"
    Delete "$INSTDIR\containershark\LICENSE.txt"
    Delete "$INSTDIR\containershark\containershark-uninstall.exe"
    RMDir "$INSTDIR\containershark"
    DeleteRegKey HKLM "${UNINSTKEY}"
    DeleteRegKey HKLM "${URLKEY}"
SectionEnd
