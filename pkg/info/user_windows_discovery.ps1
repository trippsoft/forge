#!/usr/bin/env pwsh
# This script is used to discover user information on Windows systems.
# It returns a JSON object with the user's name, ID, and home directory.

$userName = $env:USERNAME
$userId = [Security.Principal.WindowsIdentity]::GetCurrent().User.Value
$userHomeDir = $env:USERPROFILE

$output = @{
    user_name = $userName
    user_id = $userId
    user_home_dir = $userHomeDir
}

$json = $output | ConvertTo-Json -Depth 3
Write-Host $json
