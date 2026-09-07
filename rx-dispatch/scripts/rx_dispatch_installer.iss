; ===================================================================
; Inno Setup Script for RX DISPATCH
; ===================================================================

[Setup]
; Unique AppId for the application.
AppId={{A8B7C6D5-E4F3-4A21-B9C8-F7D6E5A4B3C2}}
AppName=RX DISPATCH System
AppVersion=3.0
AppPublisher=KLIK
DefaultDirName={autopf}\RX DISPATCH
DefaultGroupName=RX DISPATCH
DisableProgramGroupPage=yes
OutputBaseFilename=rx_dispatch_setup
Compression=lzma
SolidCompression=yes
WizardStyle=modern

; Require administrator privileges for installation
PrivilegesRequired=admin

[Languages]
Name: "spanish"; MessagesFile: "compiler:Languages\Spanish.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
; Source: The path to your compiled binaries and scripts.
; You must run build_all.bat BEFORE compiling this installer.
Source: "..\..\bin\*.exe"; DestDir: "{app}\bin"; Flags: ignoreversion
Source: "..\..\scripts\start_all.bat"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\..\scripts\stop_all.bat"; DestDir: "{app}"; Flags: ignoreversion

; NOTE: Add any other required files here (e.g., config files, assets)
; Example: Source: "..\..\config.json"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
; Start Menu Shortcuts
Name: "{group}\Start RX DISPATCH"; Filename: "{app}\start_all.bat"
Name: "{group}\Stop RX DISPATCH"; Filename: "{app}\stop_all.bat"
Name: "{group}\Uninstall RX DISPATCH"; Filename: "{uninstallexe}"

; Optional Desktop Shortcut
Name: "{autodesktop}\Start RX DISPATCH"; Filename: "{app}\start_all.bat"; Tasks: desktopicon

[Run]
; Optional: Run a program after installation is complete.
Filename: "{app}\start_all.bat"; Description: "{cm:LaunchProgram,RX DISPATCH System}"; Flags: nowait postinstall skipifsilent

[UninstallDelete]
Type: files; Name: "{app}\bin\*"
Type: files; Name: "{app}\*.bat"
Type: dir; Name: "{app}\bin"
Type: dir; Name: "{app}"

[Code]
procedure CurStepChanged(CurStep: TSetupStep);
var
  ResultCode: Integer;
begin
  if CurStep = ssInstall then
  begin
    // Before installing, stop any running instances of the services
    Log('Stopping any existing RX DISPATCH services...');
    Exec('taskkill.exe', '/F /IM rx-gateway.exe', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    Exec('taskkill.exe', '/F /IM rx-security.exe', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    Exec('taskkill.exe', '/F /IM rx-study.exe', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    Exec('taskkill.exe', '/F /IM rx-storage.exe', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    Exec('taskkill.exe', '/F /IM rx-image.exe', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    Exec('taskkill.exe', '/F /IM rx-reader.exe', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    Exec('taskkill.exe', '/F /IM rx-result.exe', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    Exec('taskkill.exe', '/F /IM rx-delivery.exe', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    Exec('taskkill.exe', '/F /IM rx-audit.exe', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
  end;
end;