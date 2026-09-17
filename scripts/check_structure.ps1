$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
$problems = @()
$files = @(Get-ChildItem -LiteralPath $root -File -Recurse -Force | Where-Object { $_.FullName -notmatch '\\(\.git|bin|\.pids|node_modules|build|dist|coverage)\\' })
if (Test-Path -LiteralPath (Join-Path $root 'rx-dispatch')) { $problems += 'Nested project root: rx-dispatch/' }
$modules = @($files | Where-Object { $_.Name -eq 'go.mod' })
if ($modules.Count -ne 1 -or $modules[0].DirectoryName -ne $root) { $problems += 'Expected exactly one go.mod at repository root.' }
foreach ($file in $files) {
    if ($file.Extension -eq '.md' -and $file.DirectoryName -eq $root -and $file.Name -notin @('README.md','AGENTS.md')) { $problems += "Place project documentation in docs/: $($file.Name)" }
}
$duplicates = $files | Where-Object { $_.Length -gt 0 } | Get-FileHash -Algorithm SHA256 | Group-Object Hash | Where-Object { $_.Count -gt 1 }
foreach ($group in $duplicates) { $problems += 'Review identical files: ' + (($group.Group.Path | ForEach-Object { $_.Substring($root.Length + 1) }) -join ', ') }
if ($problems.Count) { $problems | ForEach-Object { Write-Output "ERROR: $_" }; exit 1 }
Write-Output 'Structure OK: one Go root; no identical nonempty source files.'
