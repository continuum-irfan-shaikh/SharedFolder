<#
    .SYNOPSIS
       On Demand Auto Patch Deployment
    .DESCRIPTION
       On Demand Auto Patch Deployment
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
$OStype = (Get-WMIObject Win32_OperatingSystem).ProductType
if ($OStype -eq 1 ){
        $ZPMgmt_running = Get-Process "ZPMAuto" -EA SilentlyContinue
        if ($ZPMgmt_running) {
            Write-Error "Auto Patch Deployment is Already Running..!"
            Exit
        }

        $OSver = [System.Environment]::OSVersion.Version
        $OSArch = [intPtr]::Size

        Switch($OSArch){
            4 { $saazodInstallPath = (Get-ItemProperty -path HKLM:\SOFTWARE\SAAZOD).InstallationPath }
            8 { $saazodInstallPath = (Get-ItemProperty -path HKLM:\SOFTWARE\Wow6432Node\SAAZOD).InstallationPath }
        }

        $saazod_path = (Get-Item -LiteralPath $saazodInstallPath).FullName
        $zAPMgmt = $saazod_path+"BaseComponents\PatchManagement\ZPMAuto.exe"
        try{
            start-process $zAPMgmt -Wait -EA stop
            Write-Output "Auto Patch Deployment completed."
        }catch{
            if ( $_.Exception.Message -like "*The system cannot find the file specified*" ) {
            Write-Error "Auto Patch Deployment executable not found..!!"
            }Else { Write-Error $_.Exception.Message }
        }
}
Else{ Write-Error "Auto atch deployment not supported on Server Operating Syetem." }   
