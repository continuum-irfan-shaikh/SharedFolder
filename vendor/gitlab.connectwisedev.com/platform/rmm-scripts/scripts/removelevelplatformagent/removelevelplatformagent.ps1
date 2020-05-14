<#
    .SYNOPSIS
        Remove Level Platform Agent
    .DESCRIPTION
        Remove Level Platform Agent
    .Help
        To get more details refer below command. 
        HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall, HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall
        MsiExec.exe uninstallpath /quiet /norestart
    .Author
        Durgeshkumar Patel
    .Version
        1.0
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

$program = "Managed Workplace Device Manager"
$RestartArgument = '/norestart'

$ErrorActionPreference = 'Stop'

$Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall', 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'

$Product = Get-ChildItem $Registry -ErrorAction 'SilentlyContinue' | Get-ItemProperty | Where-Object { $_.DisplayName -match $program }

$ProductGUID = $Product | Select-Object -ExpandProperty PSChildName -First 1

Try {
     if ($Product) {
                                            
            $process = Start-Process "msiexec.exe" -arg "/X $ProductGUID /qn $RestartArgument" -Wait -PassThru -ErrorAction 'Stop'
            #check exit code
            If (($process.exitcode -eq '3010') -or ($process.exitcode -eq '0')) {
    
                Write-Output "`n'$program' uninstalled from the system $ENV:ComputerName"
                
            }
            else {
                Write-Warning "`nFailed to uninstall '$program' from the system $ENV:ComputerName. Exitcode: $($process.exitcode)"
                
            }
        }
        else {
            Write-Output "`n'$program' not installed on this system $ENV:ComputerName"
        }
   }
catch {
    write-output "`n"$_.Exception.Message
} 

