try{
<# Number defines how many days logs you want to monitor#>
    $startDate = (Get-Date) - (New-TimeSpan -Day 5) 
<# Logon Type is the 2-Interactive,3-Network,7-Unlock,10-RemoteInteractive, 11-CachedInteractive #>
    $UserLoginTypes = 2,3,7,10,11 
    $LastLogonUser = Get-WinEvent  -FilterHashtable @{Logname='Security';ID=4624;StartTime=$startDate} | Where-Object {-not(($_.Properties[4].Value -like  "S-1-5-18") -or ($_.Properties[4].Value -like  "S-1-5-19") -or ($_.Properties[4].Value -like  "S-1-5-20"))}  | SELECT TimeCreated, @{N='Username'; E={$_.Properties[5].Value} }, @{N='Domain/Machine'; E={$_.Properties[6].Value} },@{N='SID'; E={$_.Properties[4].Value}}, @{N='LogonType'; E={$_.Properties[8].Value}}, @{N='IP Address'; E={$_.Properties[18].Value}} | WHERE {$UserLoginTypes -contains $_.LogonType}  | Sort-Object TimeCreated | Select -last 1
Write-Output $LastLogonUser
}
catch
{
    Write-Error "Information is not available at the moment : $($_.Exception.Message)"
    Exit
}
