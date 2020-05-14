<#
   .Script
   Renew IP Address
   .Author
   Nirav Sachora
   .Description
   Script will renew the DHCP lease.
   .Requirements
#>
try
{
$ethernet = Get-WmiObject -Class Win32_NetworkAdapterConfiguration | Where { $_.IpEnabled -eq $true -and $_.DhcpEnabled -eq $true} -ErrorAction Stop 
if(!$ethernet)
    {
    Write-output "`nDHCP is not Enabled on this system"
    Exit
    }

Write-output "`nRenewing IP Addresses"
Write-Output "Please Find the List of renewed IP Addresses`n"

foreach ($lan in $ethernet) { 
   # $lan.ReleaseDHCPLease() | out-Null  
    $value = $lan.RenewDHCPLease() 
    if($value.ReturnValue -eq 0)
        {
        Write-Output ""$lan.IPaddress"`n"
        }
    else
        {
        Write-Error "IP Address Could not be renewed"
        }
    }
}
catch
{
    $_.Exception.Message
}
