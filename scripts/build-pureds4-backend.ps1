[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$OutputDir,
    [string]$Version = '0.1.0-pureds4.1',
    [switch]$RequireClean
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$repoRoot = (Resolve-Path -LiteralPath (Join-Path $PSScriptRoot '..')).Path
$outputPath = [IO.Path]::GetFullPath($OutputDir)
$resourcePath = Join-Path $repoRoot 'cmd\viiper\resource.syso'
$versionPath = Join-Path $outputPath 'versioninfo.tmp.json'
$templatePath = Join-Path $outputPath 'licenses.rendered.tpl'
$binaryName = "VIIPER-$Version-x64.exe"
$binaryPath = Join-Path $outputPath $binaryName
$noticePath = Join-Path $outputPath "VIIPER-$Version-LICENSES.txt"

New-Item -ItemType Directory -Path $outputPath -Force | Out-Null

Push-Location $repoRoot
try {
    if ($RequireClean) {
        $dirty = @(git status --porcelain=v1)
        if ($LASTEXITCODE -ne 0 -or $dirty.Count -ne 0) {
            throw 'Release backend requires a clean, committed source checkout.'
        }
    }

    & (Join-Path $PSScriptRoot 'inject-version.ps1') `
        -Version $Version -InputJson 'versioninfo.json' `
        -OutputJson $versionPath
    if (-not (Test-Path -LiteralPath $versionPath)) {
        throw 'Windows version metadata was not generated.'
    }

    go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo@v1.7.0 `
        -64 -o $resourcePath $versionPath
    if ($LASTEXITCODE -ne 0) { throw 'Windows resource generation failed.' }

    $env:CGO_ENABLED = '0'
    go build -tags release -trimpath `
        -ldflags "-s -w -X main.Version=$Version -X github.com/Alia5/VIIPER/internal/codegen/common.Version=$Version" `
        -o $binaryPath ./cmd/viiper
    if ($LASTEXITCODE -ne 0) { throw 'Backend build failed.' }

    $fileVersion = (Get-Item -LiteralPath $binaryPath).VersionInfo
    if ($fileVersion.ProductVersion -cne $Version -or
            $fileVersion.ProductName -cne 'VIIPER') {
        throw 'Backend Windows version metadata does not match this build.'
    }

    $template = (Get-Content -LiteralPath 'scripts/licenses.tpl' -Raw).
        Replace('VERSION_PLACEHOLDER', $Version)
    [IO.File]::WriteAllText($templatePath, $template,
        [Text.UTF8Encoding]::new($false))
    $notices = @(go run github.com/google/go-licenses/v2@v2.0.1 `
        report ./cmd/viiper --ignore github.com/Alia5/VIIPER `
        --template $templatePath)
    if ($LASTEXITCODE -ne 0 -or
            ($notices -join "`n") -match '\* License: \[Unknown\]') {
        throw 'Dependency license inventory is incomplete.'
    }
    $notices | Set-Content -LiteralPath $noticePath -Encoding utf8

    $buildInfo = @(go version -m $binaryPath)
    if ($LASTEXITCODE -ne 0) { throw 'Go build metadata could not be read.' }
    if ($RequireClean -and
            ($buildInfo -join "`n") -notmatch 'vcs.modified=false') {
        throw 'Release backend does not identify a clean source revision.'
    }

    Write-Output ([pscustomobject]@{
        SourceCommit = (git rev-parse HEAD)
        Binary = $binaryPath
        Sha256 = (Get-FileHash -LiteralPath $binaryPath -Algorithm SHA256).Hash
        Notices = $noticePath
        Version = $fileVersion.ProductVersion
    })
}
finally {
    Pop-Location
}
