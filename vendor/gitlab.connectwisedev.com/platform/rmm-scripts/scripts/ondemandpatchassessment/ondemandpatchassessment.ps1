<#
    .SYNOPSIS
       To execute Patch Assessment
    .DESCRIPTION
       It excecutes patch assessment binary.
    .Author
       narayan.gouda@continuum.net    
#>
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

$ZPMgmt_running = Get-Process ZPMgmt -EA SilentlyContinue
if ($ZPMgmt_running) { 
     Write-Error "Patch assessment is Already Running..!"
     Exit
}

$OSver = [System.Environment]::OSVersion.Version
$OSArch = [intPtr]::Size

Switch($OSArch){
    4 { $saazodInstallPath = (Get-ItemProperty -path HKLM:\SOFTWARE\SAAZOD).InstallationPath }      
    8 { $saazodInstallPath = (Get-ItemProperty -path HKLM:\SOFTWARE\Wow6432Node\SAAZOD).InstallationPath }
}

$saazod_path = (Get-Item -LiteralPath $saazodInstallPath).FullName
$ZPMgmt = $saazod_path+"BaseComponents\PatchManagement\ZPMgmt.exe"
try{
      start-process $ZPMgmt -ArgumentList ASSESS -Wait -EA stop
      Write-Output "Patch assessment successfully completed."
}catch{ 
    if ( $_.Exception.Message -like "*The system cannot find the file specified*" ) {
       Write-Error "Patch assement executable not found..!!" 
    }Else { Write-Error $_.Exception.Message }
}
