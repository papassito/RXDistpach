; RX-DISPATCH Inno Setup Script

; -- Defines --
#define MyAppName "RX DISPATCH"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "KLIK Soft PRO"
#define MyAppExeName "start.ps1"

[Setup]
AppId={{F2A7A6B3-9F1E-4D7C-8A3D-1B6C0E7F8D4A}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
DefaultDirName={autopf}\{#MyAppPublisher}\{#MyAppName}
DefaultGroupName={#MyAppName}
OutputDir=output
OutputBaseFilename=RX-DISPATCH-Setup
Compression=lzma
SolidCompression=yes
WizardStyle=modern
ArchitecturesInstallIn64BitMode=x64

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "spanish"; MessagesFile: "compiler:Languages\Spanish.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
; Source path is relative to the .iss file. Go up one level to the project root, then into dist/.
Source: "..\dist\RX-DISPATCH\*"; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs

[Icons]
Name: "{group}\Start {#MyAppName}"; Filename: "{app}\scripts\start.ps1"
Name: "{group}\Stop {#MyAppName}"; Filename: "{app}\scripts\stop.ps1"
Name: "{group}\Open {#MyAppName} Folder"; Filename: "{app}"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\scripts\start.ps1"; Tasks: desktopicon

[Run]
Filename: "{app}\README.txt"; Description: "{cm:LaunchProgram,View Readme}"; Flags: nowait postinstall shellexec skipifsilent

[UninstallDelete]
Type: files; Name: "{app}\.pids\*.pid"
Type: dirifempty; Name: "{app}\.pids"
Type: dirifempty; Name: "{app}\logs"
; We intentionally do not delete the 'data' directory on uninstall to preserve user-generated data.