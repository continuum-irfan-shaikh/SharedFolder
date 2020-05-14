<#  
.SYNOPSIS  
    Sophos AutoUpdate uninstallat
.DESCRIPTION  
    Sophos AutoUpdate uninstallat
.NOTES  
    File Name  : SophosAutoUpdateUninstall.ps1
    Author     : Durgeshkumar Patel  
    Requires   : PowerShell V2 or greater.   
.PARAMETERS
    
.HELP
#>

<# JSON SCHEMA
#$version    #Drop Down
2.6.0
2.9.0.344
3.1.1.18
4.1.0.273
5.2.0.276
5.7.220
5.10.139
5.11.206

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

$program = "Sophos AutoUpdate"
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
            if ($restart -eq $true) { $RestartArgument = '/forcerestart' }
            if ($restart -eq $false) { $RestartArgument = '/norestart' }
                                
            $process = Start-Process "msiexec.exe" -arg "/X $ProductGUID /qn $RestartArgument" -Wait -PassThru -ErrorAction 'Stop'
            #check exit code
            If (($process.exitcode -eq '3010') -or ($process.exitcode -eq '0')) {
    
                Write-Output "`n'$program' uninstalled from the system $ENV:ComputerName"
                bootrequiredmsg #System Reboot Message through function
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
