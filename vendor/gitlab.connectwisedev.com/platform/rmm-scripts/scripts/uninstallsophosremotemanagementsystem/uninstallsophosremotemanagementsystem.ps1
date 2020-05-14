<#  
.SYNOPSIS  
    Sophos Remote Management System Uninstallation
.DESCRIPTION  
    Sophos Remote Management System Uninstallation
.NOTES  
    File Name  : SophosRemoteManagementSystemUninstall.ps1
    Author     : Durgeshkumar Patel  
    Requires   : PowerShell V2 or greater.   
.PARAMETERS
    
.HELP
#> 

<# JSON SCHEMA
#$version    #Drop Down
4.0.2
4.1.1

#radio button
#$restart = $false       $true = restart  $false = norestart
#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

$program = "Sophos Remote Management System"
$action = 'uninstall'

$ErrorActionPreference = 'Stop'
function bootrequiredmsg {
        
    if ($restart -eq $true) {
        Write-Output "`nSystem $ENV:ComputerName will reboot now"
    }
      
    if ($restart -eq $false) {
        Write-Output "`nSystem $ENV:ComputerName will not reboot"
    }   
}

$Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall', 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'

$Product = Get-ChildItem $Registry -ErrorAction 'SilentlyContinue' | Get-ItemProperty | Where-Object { $_.DisplayName -match $program -and $_.DisplayVersion -like "*$version*" }

$ProductGUID = $Product | Select-Object -ExpandProperty PSChildName -First 1

Try {
    if ($action -eq "uninstall") {
        if ($Product) {
            
            $process = Start-Process "msiexec.exe" -arg "/X $ProductGUID /qn /norestart" -Wait -PassThru -ErrorAction 'Stop'

            #check exit code
            If (($process.exitcode -eq '3010') -or ($process.exitcode -eq '0') -or ($process.exitcode -eq '1641')) {
    
                Write-Output "`n'$program' uninstalled from the system $ENV:ComputerName"
                
                bootrequiredmsg #System Reboot Message through function
                Start-Sleep -Seconds 10
                if ($restart -eq $true) {
                
                    Restart-Computer -Force
                }

            }
            else {
                Write-Warning "`nFailed to uninstall '$program' from the system $ENV:ComputerName. Exitcode: $($process.exitcode)"
                
            }
        }
        else {
            Write-Output "`n'$program' v$version not installed on this system $ENV:ComputerName"
        }
    }
    else {
        Write-Output "`nKindly select action as 'Uninstall'"
    }
}
catch {
    write-output "`n"$_.Exception.Message
}
