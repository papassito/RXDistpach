; RX-DISPATCH Inno Setup Script

#define MyAppName "RX DISPATCH"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "KLIK Soft PRO"

[Setup]
AppId={{F2A7A6B3-9F1E-4D7C-8A3D-1B6C0E7F8D4A}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
DefaultDirName={autopf}\KLIK Soft PRO\RX DISPATCH
DefaultGroupName={#MyAppName}
OutputDir=output
OutputBaseFilename=RX-DISPATCH-Setup
Compression=lzma2
SolidCompression=yes
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "spanish"; MessagesFile: "compiler:Languages\Spanish.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
Source: "..\dist\RX-DISPATCH\*"; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs

[Icons]
Name: "{group}\Start {#MyAppName}"; Filename: "{sys}\WindowsPowerShell\v1.0\powershell.exe"; Parameters: "-NoProfile -ExecutionPolicy Bypass -File ""{app}\scripts\start.ps1"""; WorkingDir: "{app}"
Name: "{group}\Stop {#MyAppName}"; Filename: "{sys}\WindowsPowerShell\v1.0\powershell.exe"; Parameters: "-NoProfile -ExecutionPolicy Bypass -File ""{app}\scripts\stop.ps1"""; WorkingDir: "{app}"
Name: "{group}\Open {#MyAppName} Folder"; Filename: "{app}"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{sys}\WindowsPowerShell\v1.0\powershell.exe"; Parameters: "-NoProfile -ExecutionPolicy Bypass -File ""{app}\scripts\start.ps1"""; WorkingDir: "{app}"; Tasks: desktopicon

[Run]
Filename: "{app}\README.txt"; Description: "{cm:LaunchProgram,View Readme}"; Flags: nowait postinstall shellexec skipifsilent

[UninstallDelete]
Type: files; Name: "{app}\.pids\*.pid"
Type: dirifempty; Name: "{app}\.pids"
Type: dirifempty; Name: "{app}\logs"

; Do not delete data/ on uninstall.