<#  
.SYNOPSIS  
    Adobe Flash Player Uninstaller
.DESCRIPTION  
    Uninstall Adobe Flash Player
    Uninstall Adobe Flash Player ActiveX
.NOTES  
    File Name  : AdobeFlashPlayerUninstaller.ps1
    Author     : Durgeshkumar Patel  
    Requires   : PowerShell V2 or greater.   
    .PARAMETERS
            $action
            $version
    .HELP
            #Registry::HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Adobe Flash Player Active  
            #Registry::HKLMHKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Adobe Flash Player NPAPI
            #Registry::HKLMHKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Adobe Flash Player Plugin
#>
<# JSON Schema Variables
#$action = "uninstall"
#$version = "ALL"
#$version = "ActiveX"
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

$registry = "Registry::HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall", "HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall"

if (($version -eq "ALL") -and ($action -eq "uninstall")) {

    $Program = "Adobe Flash Player"
    
    $Product = Get-ChildItem $Registry -ErrorAction 'SilentlyContinue' | Get-ItemProperty | Where-Object { $_.DisplayName -like "Adobe Flash Player *" }   
}

if (($version -eq "ActiveX") -and ($action -eq "uninstall")) {
    
    $Program = "Adobe Flash Player ActiveX"
  
    $Product = Get-ChildItem $Registry -ErrorAction 'SilentlyContinue' | Get-ItemProperty | Where-Object { $_.DisplayName -like "Adobe Flash Player * ActiveX" } 
}

#Uninstall Code
if ($Product) { 

    foreach ($Prod in $Product) {
        $ProductGUID = $Prod | Select-Object -ExpandProperty UninstallString -First 1
    
        if ($ProductGUID -match 'MsiExec.exe') {
            $PGUID = $Prod | Select-Object -ExpandProperty PSChildName -First 1
            $process = Start-Process "msiexec.exe" -arg "/X $PGUID /qn /norestart" -Wait -PassThru -ErrorAction 'Stop'
        }
        else {
            $path = $ProductGUID.split(" -") | select -First 1
            if (Test-Path $path) {
                $process = Start-Process $path -arg "-uninstall" -Wait -PassThru -ErrorAction 'Stop'
            }
        }
            
        #check exit code
        If (($process.exitcode -eq '3010') -or ($process.exitcode -eq '0')) {
              
            Write-Output "`n'$($Prod.DisplayName)' uninstalled from the system $ENV:ComputerName"
        }
        else {
            Write-Warning "`nFailed to uninstall '$($Prod.DisplayName)' from the system $ENV:ComputerName. Exitcode: $($process.exitcode)"
        }
    }
}
else {
    Write-Output "`n'$program' not installed on this system $ENV:ComputerName"
}

