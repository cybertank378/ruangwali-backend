function Get-CleanTree {
    param (
        [string]$Path = ".",

        [string[]]$Exclude = @(
    # Version Control
        ".git",

        # IDE / Editor
        ".idea",
        ".vscode",

        # Go Build Output
        "bin",
        "dist",
        "build",

        # Test / Coverage Output
        "coverage",
        "test-results",

        # Generated Code
        "generated",
        "gen",

        # Temporary / Cache
        "tmp",
        "temp",
        ".cache",

        # Vendor Dependencies
        "vendor"
    ),

        [string[]]$ExcludeExtensions = @(
        ".exe",
        ".dll",
        ".so",
        ".dylib",
        ".out",
        ".test",
        ".prof",
        ".cover"
    ),

        [string[]]$ExcludeFiles = @(
        "clean-tree.txt"
    ),

        [int]$Depth = 20,

        [string]$Indent = ""
    )

    try {
        $items = Get-ChildItem `
      -LiteralPath $Path `
      -Force `
      -ErrorAction Stop |
                Where-Object {
                    $Exclude -notcontains $_.Name -and
                            $ExcludeFiles -notcontains $_.Name -and
                            (
                            $_.PSIsContainer -or
                                    $ExcludeExtensions -notcontains $_.Extension
                            )
                } |
                Sort-Object `
        @{ Expression = "PSIsContainer"; Descending = $true },
                @{ Expression = "Name"; Descending = $false }

        foreach ($item in $items) {
            Write-Output "$Indent├── $($item.Name)"

            if ($item.PSIsContainer -and $Depth -gt 0) {
                Get-CleanTree `
          -Path $item.FullName `
          -Exclude $Exclude `
          -ExcludeExtensions $ExcludeExtensions `
          -ExcludeFiles $ExcludeFiles `
          -Depth ($Depth - 1) `
          -Indent "$Indent│   "
            }
        }
    }
    catch {
        Write-Output "$Indent├── [ACCESS DENIED] $Path"
    }
}

Get-CleanTree |
        Out-File `
    -FilePath "clean-tree.txt" `
    -Encoding utf8