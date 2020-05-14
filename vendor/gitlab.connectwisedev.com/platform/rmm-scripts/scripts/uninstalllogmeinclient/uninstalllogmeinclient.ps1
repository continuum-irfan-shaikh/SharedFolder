<#
    .SYNOPSIS
       Uninstall LogMeIn Client 
    .DESCRIPTION
       Uninstall LogMeIn Client 
    .Help
        To get more details refer below command. 
        Use MSI uninstallation method.
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>

$program = "LogMeIn Client"

Try {    
        
    if (Get-WmiObject -Class Win32_Product | where-Object { $_.Name -eq $program }) {
        $ErrorActionPreference = 'stop'
        Get-WmiObject -Class Win32_Product | where-Object { $_.Name -eq $program } | foreach { $_.Uninstall() } | out-null
        
        $check = Get-WmiObject -Class Win32_Product | where-Object { $_.Name -eq $program }
                
        if ([string]::IsNullOrEmpty($check)) {
            Write-Output "`nLogMeIn Client uninstalled from the system $ENV:COMPUTERNAME"
        }
        else {
            Write-Error "Uninstalling LogMeIn Client failed on system $ENV:COMPUTERNAME"
        }
    }
    else {
        Write-Output "`nLogMeIn Client not installed on this system $ENV:COMPUTERNAME"
    }
}
catch {
    Write-Error $_.Exception.Message
}
