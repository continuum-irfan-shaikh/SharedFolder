<#
    .SYNOPSIS
        Enable file and printer sharing throguh windows firewall
    .DESCRIPTION
        Enable file and printer sharing throguh windows firewall. 
    .Help
       netsh advfirewall firewall set rule group="File and Printer Sharing" new enable=yes
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#> 

$rules=(New-object -ComObject HNetCfg.FWPolicy2).rules
$rule = $rules | Where-Object {$_.enabled -eq $false -and $_.direction -eq 1 -and $_.name -match "File and Printer Sharing" }

if($rule -eq $null)
{
    Write-Output "`nFile and Printer Sharing already enabled on system $env:COMPUTERNAME"
}
else
{
    netsh advfirewall firewall set rule group="File and Printer Sharing" new enable=yes
    if ($?)
    {
    Write-Output "File and Printer Sharing enabled on system $env:COMPUTERNAME"
    }
    else {
    Write-Error $_.Exception.Message
    Write-Error "`nFile and Printer Sharing not enabled on system $env:COMPUTERNAME. Kindly check manually."
    }
}
