$ErrorActionPreference = "Stop"

$ProjectRoot = $PSScriptRoot
$EnvFile = Join-Path $ProjectRoot ".env"

Write-Host ""
Write-Host "========================================="
Write-Host " RuangWali Backend"
Write-Host "========================================="
Write-Host ""

# =========================================================
# VALIDATE PROJECT
# =========================================================

$GoModFile = Join-Path $ProjectRoot "go.mod"

if (-not (Test-Path $GoModFile)) {
    Write-Error "go.mod tidak ditemukan: $GoModFile"
    exit 1
}

# =========================================================
# VALIDATE .ENV
# =========================================================

if (-not (Test-Path $EnvFile)) {
    Write-Error "File .env tidak ditemukan: $EnvFile"
    exit 1
}

# =========================================================
# LOAD ENVIRONMENT VARIABLES
# =========================================================

Write-Host "[1/3] Loading environment variables..."

Get-Content $EnvFile | ForEach-Object {
    $Line = $_.Trim()

    if (
    $Line -and
            -not $Line.StartsWith("#")
    ) {
        $Parts = $Line -split "=", 2

        if ($Parts.Count -ne 2) {
            Write-Error "Format .env tidak valid: $Line"
            exit 1
        }

        $Name = $Parts[0].Trim()
        $Value = $Parts[1].Trim()

        if ([string]::IsNullOrWhiteSpace($Name)) {
            Write-Error "Nama environment variable tidak boleh kosong"
            exit 1
        }

        [Environment]::SetEnvironmentVariable(
                $Name,
                $Value,
                "Process"
        )
    }
}

Write-Host "      Environment variables loaded."

# =========================================================
# VALIDATE REQUIRED ENVIRONMENT VARIABLES
# =========================================================

Write-Host "[2/3] Validating required configuration..."

$RequiredVariables = @(
    "APP_ENV",
    "HTTP_ADDR",
    "DATABASE_URL",
    "JWT_ISSUER",
    "JWT_AUDIENCE",
    "JWT_SECRET"
)

foreach ($VariableName in $RequiredVariables) {
    $Value = [Environment]::GetEnvironmentVariable(
            $VariableName,
            "Process"
    )

    if ([string]::IsNullOrWhiteSpace($Value)) {
        Write-Error "Environment variable wajib belum tersedia: $VariableName"
        exit 1
    }
}

Write-Host "      Required configuration valid."

# =========================================================
# RUN API
# =========================================================

Write-Host "[3/3] Starting API..."
Write-Host ""
Write-Host "Project : $ProjectRoot"
Write-Host "Env     : $env:APP_ENV"
Write-Host "HTTP    : $env:HTTP_ADDR"
Write-Host ""

Set-Location $ProjectRoot
$env:GOEXPERIMENT = "jsonv2"
go run ./cmd/api

if ($LASTEXITCODE -ne 0) {
    Write-Error "RuangWali API berhenti dengan exit code $LASTEXITCODE"
    exit $LASTEXITCODE
}