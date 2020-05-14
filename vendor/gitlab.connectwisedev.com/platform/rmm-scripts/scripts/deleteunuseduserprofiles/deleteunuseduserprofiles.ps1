<#    
    .Synopsis
    Script will check user profiles in computer and delete if profile is unused from last 180 days.
    .Author
    Nirav Sachora
    .Requirements
    Script should run with highest privilege on computer
    
#>
#$notfound = Get-WMIObject -class Win32_UserProfile | Where {($_.LastUsetime -eq $null)} | Select -ExpandProperty LocalPath
$profiles = Get-WMIObject -class Win32_UserProfile | Where {($_.LastUsetime -ne $null) -and (Test-Path $_.Localpath) -and (!$_.Loaded) -and (!$_.Special) -and ($_.ConvertToDateTime($_.LastUseTime) -lt (Get-Date).AddDays(-180)) -and ($_.Localpath -ne 'C:\Users\administrator') -and ($_.Localpath -ne 'C:\Users\Remote Support')}
$count = $profiles | Measure-Object | select -ExpandProperty count
if($count -eq 0)
{
Write-output "No Profiles found for Deletion"
exit
}
<#if($notfound){
"-"*40 + "`n[Error]Error while fetching last login time for below profiles`n$notfound`n" + "-"*40
}#>
write-output "Number of profiles : $count`n"
$profiledeleted = @()
$profileerror = @()
foreach($profile in $profiles)
{
    $profile | Remove-WmiObject
        if($?)
        {
        $profiledeleted += $profile.Localpath 
        continue
        }
        else
        {
        $profileerror += $profile.Localpath
        }
 }
 if($profiledeleted.length -gt 0)
 {
 write-output "Localpath of removed profiles"
 $profiledeleted
 }
 if($profileerror.length -gt 0)
 {
 write-output "Error While removing below profiles"
 $profileerror
 }
