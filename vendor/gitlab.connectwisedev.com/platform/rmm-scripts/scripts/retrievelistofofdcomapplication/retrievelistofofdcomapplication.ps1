try {
    $ErrorActionPreference = 'Stop'
    foreach ($application in $(Get-WmiObject -class "Win32_DCOMApplication")) { 
        Write-Output "`nAppID       : $($application.AppID)"
        Write-Output "Name        : $($application.Name)"
        Write-Output "Description : $($application.Description)"
    } 
}
catch {
   Write-Error $_
}
