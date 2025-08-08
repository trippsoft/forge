#!/usr/bin/env pwsh
# This script is used to discover OS information on Windows systems.
# It returns 1 if FIPS mode is enabled, otherwise it returns 0.
$registryPath = 'HKLM:\SYSTEM\CurrentControlSet\Control\Lsa\FipsAlgorithm'
$value = Get-ItemPropertyValue -LiteralPath $registryPath \
            -Name 'Enabled' \
            -ErrorAction 'SilentlyContinue'

Write-Host $value
