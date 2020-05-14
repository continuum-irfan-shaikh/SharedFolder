<#
    .SYNOPSIS
       Retrieve Network Adapter Information
    .DESCRIPTION
       Retrieve Network Adapter Information
    .Author
       Santosh.Dakolia@continuum.net    
#>

$computer = $env:COMPUTERNAME

Try{
function Netinfo ($strIN)
{
 $num = $strIN.length
 for($i=1 ; $i -le $num ; $i++)
  { $Fline = $Fline + "=" }
}
Write-Host "Network adapter settings on $computer"
Get-WmiObject -Class win32_NetworkAdapterSetting -computername $computer |
Foreach-object  {
  If( ([wmi]$_.element).netconnectionstatus -eq 2)
      {
    $NetStatus = switch (([wmi]$_.element).netconnectionstatus )
    {
        0 { "Disconnected" }
        1 { "Connecting" }
        2 { "Connected" }
        3 { "Disconnecting" }
        4 { "Hardware not present" }
        5 { "Hardware disabled" }
        6 { "Hardware malfunction" }
        7 { "Media disconnected" }
        8 { "Authenticating" }
        9 { "Authentication succeeded" }
        10 { "Authentication failed" }
        11 { "Invalid address" }
        12 { "Credentials required" }
    }
    Write-Output "Network Status :: $NetStatus"
     [wmi]$_.setting
     [wmi]$_.element
          }
 } 
 }Catch{
     Write-Error "Error occured..!! $_.Exception.Message"

 }
