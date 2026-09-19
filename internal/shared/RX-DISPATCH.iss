; RX-DISPATCH Inno Setup Script
; Ubicación esperada de este archivo: Raíz del proyecto (RXDistpach\installer.iss)

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
Compression=lzma2/ultra64
SolidCompression=yes
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible

[Languages]
Name: "spanish"; MessagesFile: "compiler:Languages\Spanish.isl"
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
Source: "..\..\dist\RX-DISPATCH\*"; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs

[Icons]
Name: "{group}\Iniciar {#MyAppName}"; Filename: "{sys}\WindowsPowerShell\v1.0\powershell.exe"; Parameters: "-NoProfile -ExecutionPolicy Bypass -WindowStyle Hidden -File ""{app}\scripts\start.ps1"""; WorkingDir: "{app}"
Name: "{group}\Detener {#MyAppName}"; Filename: "{sys}\WindowsPowerShell\v1.0\powershell.exe"; Parameters: "-NoProfile -ExecutionPolicy Bypass -WindowStyle Hidden -File ""{app}\scripts\stop.ps1"""; WorkingDir: "{app}"
Name: "{group}\Abrir Carpeta de {#MyAppName}"; Filename: "{app}"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{sys}\WindowsPowerShell\v1.0\powershell.exe"; Parameters: "-NoProfile -ExecutionPolicy Bypass -WindowStyle Hidden -File ""{app}\scripts\start.ps1"""; WorkingDir: "{app}"; Tasks: desktopicon

[Run]
Filename: "{app}\README.txt"; Description: "{cm:LaunchProgram,README.txt}"; Flags: nowait postinstall shellexec skipifsilent skipifdoesntexist

[UninstallRun]
Filename: "{sys}\WindowsPowerShell\v1.0\powershell.exe"; Parameters: "-NoProfile -ExecutionPolicy Bypass -WindowStyle Hidden -File ""{app}\scripts\stop.ps1"""; WorkingDir: "{app}"; Flags: runhidden

[UninstallDelete]
Type: filesandordirs; Name: "{app}\.pids"
Type: filesandordirs; Name: "{app}\logs"