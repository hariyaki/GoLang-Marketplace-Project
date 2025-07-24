$schema = "ProjectSchema.txt"
$ignore = @(
  'node_modules', 'dist', 'uploads', 'docs\swagger.*',
  '.git', '.idea', '.vscode', '.DS_Store'
) -join '|'

# 'tree' is present on every Windows; use /F (files) /A (ASCII)
$raw = tree /F /A | Out-String

# simple regex filter
$filtered = $raw -split "`n" | Where-Object { $_ -notmatch $ignore }

$filtered | Out-File $schema -Encoding utf8
Write-Host "ProjectSchema.txt updated."