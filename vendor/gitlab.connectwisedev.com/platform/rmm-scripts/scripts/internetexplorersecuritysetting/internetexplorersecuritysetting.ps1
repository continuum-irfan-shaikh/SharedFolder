# Internet explorer security zone setiings registry values
# https://support.microsoft.com/en-us/help/182569/internet-explorer-security-zones-registry-entries-for-advanced-users
$Levels = @{
    '0' = 'Custom'
    '10000' = 'Low'
    '10500' = 'Medium Low'
    '11000' = 'Medium'
    '11500' = 'Medium High'
    '12000' = 'High'
}

try {
    $ErrorActionPreference = 'Stop'
    0..4 | Foreach {
        $Settings = Get-ItemProperty "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings\Zones\$_"
        If ($Settings) {
            Write-Output "Zone name: $($Settings.DisplayName)"
            Write-Output "Security level: $($levels[$('{0:x}' -f $Settings.currentlevel)])`n"
        }
    }
}
catch {
    Write-Error $_.exception.message
}
