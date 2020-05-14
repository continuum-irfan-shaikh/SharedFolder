<#  
.SYNOPSIS  
    AVG Anti-Virus Uninstallation
.DESCRIPTION  
    AVG Anti-Virus Uninstallation
.NOTES  
    File Name  : AVGAntiVirusUninstall.ps1
    Author     : Durgeshkumar Patel  
    Requires   : PowerShell V2 or greater.   
.PARAMETERS
    
.HELP
#> 

<# JSON SCHEMA
#$action = 'uninstall'  #Drop Down
#$version    #Drop Down
2016.0.4756
2015.0.4342
2014.0.4311
2013.0.3495
2012.0.4311
2012.1.2265
8
9

#radio buttons
#$restart = $false  #true

Examples:-
$action = "uninstall"
$version = "2012.1.2265"
$restart = $false

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

$program = "AVG"

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

$ProductGUID = $Product | Select-Object -ExpandProperty UninstallString -First 1

if ($ProductGUID) {

    $path = $ProductGUID.split("/") 
    
    $arg1 = $path[1]
    $arg2 = $path[2]
}


Try {
    if ($Product) {
           
        if ($restart -eq $false) { $RestartArgument = '/dontrestart' }
            
        $uninstallinfo = New-Object System.Diagnostics.ProcessStartInfo
        $uninstallinfo.CreateNoWindow = $true
        $uninstallinfo.UseShellExecute = $false
        $uninstallinfo.RedirectStandardOutput = $true
        $uninstallinfo.RedirectStandardError = $true
        $uninstallinfo.FileName = $path[0]
        $uninstallinfo.Arguments = "/$arg1 /$arg2 /uilevel=Silent $RestartArgument"
        $uninstall = New-Object System.Diagnostics.Process
        $uninstall.StartInfo = $uninstallinfo
        [void]$uninstall.Start()
        $uninstall.WaitForExit()

        #check exit code
        If (($uninstall.exitcode -eq '0') -or ($uninstall.exitcode -eq '3010') -or ($uninstall.exitcode -eq '1641')) {
    
            Write-Output "`n'$program' v$version uninstalled from the system $ENV:ComputerName"
            bootrequiredmsg #System Reboot Message through function
             if ($restart -eq $true) {
            
                Restart-Computer -Force
            }
        }
        else {
            Write-Warning "`nFailed to uninstall '$program' v$version from the system $ENV:ComputerName. Exitcode: $($process.exitcode)"
                
        }
    }
    else {
        Write-Output "`n'$program' v$version not installed on this system $ENV:ComputerName"
    }
   
}
catch {
    Write-Output "`n"$_.Exception.Message
} 


