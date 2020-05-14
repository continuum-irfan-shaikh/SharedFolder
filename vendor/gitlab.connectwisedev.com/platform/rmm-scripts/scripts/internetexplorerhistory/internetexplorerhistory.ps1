<#
Script to capture IE History from all user profiles on a workstation or a server.
#>

try {
    $ErrorActionPreference = 'stop'
    $Shell = New-Object -ComObject Shell.Application
    $UserPaths = Get-WmiObject win32_userprofile | Where-Object { $_.sid -like "S-1-5-21*"}| Select-Object -ExpandProperty localpath
    $computer = $env:COMPUTERNAME
}
catch {
    Write-Error $_
}

# sometimes namespace is not able to open the folder, 'SilentlyContinue' to skip those errors
# and capture IE history from remaining user profiles
$ErrorActionPreference = 'SilentlyContinue' 
ForEach ($UserPath in $UserPaths) {
    $IEHistory = $Shell.NameSpace("$UserPath\AppData\Local\Microsoft\Windows\History")
    if ($IEHistory) {
        # nested loops to iterate through nested history folder hierarchy        
        Foreach ($Item in $IEHistory.Items()) {
            Foreach ($WebSiteItem in $Item.GetFolder.Items()) {
                If ($WebSiteItem.IsFolder) {
                    $SiteFolder = $WebSiteItem.GetFolder
                    ForEach ($Site in $SiteFolder.Items()) {
                        $Title = $($SiteFolder.GetDetailsOf($Site, 1))
                        $URL = $($SiteFolder.GetDetailsOf($Site, 0))
                        $Date = $($SiteFolder.GetDetailsOf($Site, 2))
                        $username = (Split-Path $UserPath -Leaf)
                        if ($Url -like "http*") { # to exclude file paths in history folder and include URLs starting with HTTP or HTTPS
                            Write-Output "`nComputer     : $computer"
                            Write-Output "User         : $(if($env:USERNAME -eq $username){"$username (Current logged-on User)"}else{$username})"
                            Write-Output "Last Visited : $Date"
                            Write-Output "URL          : $URL"
                        }
                    }
                }
            }
        }
    }
}
