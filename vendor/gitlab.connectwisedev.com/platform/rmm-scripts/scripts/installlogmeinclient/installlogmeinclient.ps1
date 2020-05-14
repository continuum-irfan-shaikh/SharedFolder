<#
   .SYNOPSIS  
        LogMeIn Remote Control Client.  
   .DESCRIPTION  
        The script installs LogMeInIgnition.msi applicatrion.
        This script supports only Windows-7/Windows2008-R2 and higher version of OS.         
   .NOTES  
        File Name  : LogMeInIgnition.ps1  
        Author     : GRT
        Requires   : PowerShell V2 and higher.
        Version    : 1.0
        Date       : 05/21/2019
#>
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

$OSver = [System.Environment]::OSVersion.Version

if($OSver.Build -le 7600) {
   Write-Output "Not supported on this operating system."
   Exit
}


$oscheck = (Get-WmiObject win32_operatingsystem).OSarchitecture

$precheck = Get-WmiObject win32_product | ? { $_.Name -eq "LogMeIn Client" }
if($precheck){
    Write-output "LogMeIn Remote Client is already installed on this system."
    Exit
}

$url = "https://secure.logmein.com/LogMeInIgnition.msi"
$downloadpath = Join-Path $env:SystemRoot "\Temp\LogMeInIgnition.msi"

$ErrorActionPreference = "Stop"
try{
    $wc = New-Object System.Net.WebClient
    $wc.DownloadFile("$url", "$downloadpath")
}
catch{
    
    Write-Error "LogMeIn Remote Control Client download failed. $_"
}

$Erroractionpreference = "Continue"
$pinfo = New-Object System.Diagnostics.ProcessStartInfo
$pinfo.FileName = "msiexec.exe"
$pinfo.RedirectStandardError = $true
$pinfo.RedirectStandardOutput = $true
$pinfo.UseShellExecute = $false
$pinfo.Arguments = "/i  C:\Windows\Temp\LogMeInIgnition.msi /quiet"
$p = New-Object System.Diagnostics.Process
$p.StartInfo = $pinfo
$p.Start() | Out-Null
$p.WaitForExit()
        
if ($p.ExitCode -eq 0) {
    $installresult = Get-WmiObject win32_product | ? { $_.Name -eq "LogMeIn Client" }
    if ($installresult) {
        Write-Output "LogMeIn Remote Control Client has been successfully installed."
        
    }
}
else {
    Write-Error "LogMeIn Remote Control Client installation failed."
}
Remove-Item $downloadpath -ErrorAction SilentlyContinue | Out-Null
