<#
    .SYNOPSIS
        Remove/Uninstall ScreenConnect
    .DESCRIPTION
        Remove/Uninstall ScreenConnect software from the windows system. 
    .Help
        To get more details refer below command. 
        Get-WmiObject -Class Win32_Product
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>
$program  = "ScreenConnect"
$app = Get-WmiObject -Class Win32_Product | where-Object {$_.Name -match $program}
Try 
{
    If ($app -eq $null)
    {
    Write-Output "`nScreenConnect is not installed in the system $ENV:ComputerName"
    }
    Else
    {
    $app.Uninstall()  2>&1 | Out-Null
    $check = Get-WmiObject -Class Win32_Product | where-Object {$_.Name -match $program}
    
        if ($check -eq $null)
        {
        Write-Output "`nScreenConnect removed successfully from the system $ENV:ComputerName"
        }
    }
}
catch
{
    write-output "`n"$_.Exception.Message
}
