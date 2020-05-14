<#
    .Script
    Script to retrieve pagefile information across all the drives in the system
    .Author
    Nirav Sachora
    .Requirements
    Script should run with highest privileges.
#>

$details = Get-WMIObject -Class Win32_PageFileUsage | Select-Object Caption,Name,AllocatedBaseSize,CurrentUsage,peakusage
if(!$details -and $?)
{
    Write-Output "Pagefiles are not configured on this system"
}
elseif($? -eq $false)
{
    Write-Error "Error While retrieving Pagefile information"
}
else
{
    if($details.count -gt 1){$count = $details.count}else{$count = 1}
    "`nNumber of Pagefiles in the system : {0}" -f $count
    foreach($detail in $details)
    {
     "`nPagefile     : {0}"   -f $detail.Name
     "Pagefile Size: {0} MB" -f $detail.AllocatedBaseSize
     "Current Usage: {0} MB" -f $detail.CurrentUsage
     ("Peak Usage   : {0} MB" -f $detail.peakusage).Trim()
    }
}
