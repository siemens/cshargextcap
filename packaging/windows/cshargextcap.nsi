; ContainerShark external capture plugin NSIS installer script for Windows

; Create a Unicode-based installer; requires NSIS 3+.
Unicode true

!include "MUI2.nsh"
!include 'LogicLib.nsh'
!include "FileFunc.nsh"

!include "pluginversion.nsh"
!define DESCRIPTION "ContainerShark external capture plugin for capturing container network traffic"

Name "ContainerShark External Capture Plugin ${VERSION} for Wireshark"
OutFile "${INSTALLERBINARY}"
Icon "containershark.ico"
VIProductVersion "${FILEVERSION}.0"
VIAddVersionKey "FileDescription" "${DESCRIPTION}"
VIAddVersionKey "LegalCopyright" "${COPYRIGHT}"
VIAddVersionKey "FileVersion" "${VERSION}.0"
VIAddVersionKey "CompanyName" "Edgeshark Project"

; Wireshark doesn't necessarily indicate its installation directory
; in the registry; probably in case it's in the default location?
; We first set an educated guess...
InstallDir "$PROGRAMFILES\Wireshark\extcap"
Function .onInit
    ${If} ${FileExists} "$PROGRAMFILES64\Wireshark\*.*"
        StrCpy $INSTDIR "$PROGRAMFILES64\Wireshark\extcap"
    ${EndIf}
FunctionEnd

!define MUI_ICON "containershark.ico"
!define MUI_UNICON "containershark.ico"
BrandingText "ContainerShark Installer"

!define MUI_PAGE_HEADER_TEXT "ContainerShark"
!define MUI_PAGE_HEADER_SUBTEXT "\
A deep dive into your virtual container networks."

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

!define MUI_WELCOMEFINISHPAGE_BITMAP "welcomefinish.bmp"
!define MUI_FINISHPAGE_NOAUTOCLOSE

!define MUI_FINISHPAGE_RUN "$INSTDIR\..\wireshark.exe"
!define MUI_FINISHPAGE_RUN_TEXT "Run Wireshark"

!define MUI_ABORTWARNING

; Installer pages...
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "license.txt"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

; Uninstaller pages...
!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

!insertmacro MUI_LANGUAGE "English"

Section "Install"
    ; We need to write the uninstaller using a unique name, as we might otherwise
    ; conflict with other extcap plugin installers. Not that this is really a topic
    ; well thought out by Wireshark for Windows.
    SetOutPath "$INSTDIR"
    File "${BINARYPATH}/${BINARYNAME}"
    File "containershark.ico"
    WriteRegStr HKCU "Software\containershark" "" "$INSTDIR"
    ; We need to put the uninstaller into a subdir, because otherwise Wireshark
    ; will run it thinking it might be an extcap plugin...
    CreateDirectory "$INSTDIR\containershark"
    WriteUninstaller "$INSTDIR\containershark\containershark-uninstall.exe"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\ContainerShark" "DisplayName" "ContainerShark ExtCap Plugin for Wireshark"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\ContainerShark" "Publisher" "Edgeshark Project"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\ContainerShark" "DisplayIcon" "$\"$INSTDIR\containershark.ico$\""
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\ContainerShark" "Uninstallstring" "$\"$INSTDIR\containershark\containershark-uninstall.exe$\""
    ; Registering packetflix: URL scheme handler
    WriteRegStr HKCU "Software\Classes\Packetflix" "" "URL:Packetflix Protocol"
    WriteRegStr HKCU "Software\Classes\Packetflix" "URL Protocol" ""
    WriteRegStr HKCU "Software\Classes\Packetflix\DefaultIcon" "" "$\"$INSTDIR\containershark.ico$\""
    WriteRegStr HKCU "Software\Classes\Packetflix\shell\open\command" "" 'wireshark -k -i packetflix -o "extcap.packetflix.url:%1"'
SectionEnd

Section "Uninstall"
    ; Please remember that we use an uninstaller with a dedicated unique name in
    ; order to avoid uninstall.exe name clashes in case other extcaps were installed
    ; using their own installers and uninstallers too.
    ;
    ; Note #1: we must not remove the installation directory as this will be
    ; Wireshark's extcap directory.
    ;
    ; Note #2: since we had to move the uninstaller into a subfolder below the
    ; plugin installation place, we now need to go up to this directory, because
    ; the uninstaller sets $INSTDIR to point to the directory where it was written
    ; to...
    ${GetParent} "$INSTDIR" $INSTDIR
    Delete "$INSTDIR\${BINARYNAME}"
    Delete "$INSTDIR\containershark.ico"
    Delete "$INSTDIR\containershark\*.*"
    RMDir "$INSTDIR\containershark"
    DeleteRegKey /ifempty HKCU "Software\containershark"
    DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\ContainerShark"
    DeleteRegKey HKCU "Software\Classes\Packetflix"
SectionEnd
