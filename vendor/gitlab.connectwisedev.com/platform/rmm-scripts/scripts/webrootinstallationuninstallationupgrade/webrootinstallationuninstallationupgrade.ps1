<#
    .Synopsis
        Script will Install/Uninstall/Upgrade Webroot from the system.
    .Author
        Nirav Sachora
    .Requirements
        Script Should run with admin privileges.
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

$bit = (Get-WmiObject Win32_operatingsystem).Osarchitecture

if ($bit -eq "64-bit") {
    $wrpath = "C:\Program Files (x86)\Webroot"
}
else {
    $wrpath = "C:\Program Files\Webroot"
}


Function Check_Webroot {
    
$bit = (Get-WmiObject Win32_operatingsystem).Osarchitecture
    if ($bit -eq "64-bit") {
        $wrpath = "C:\Program Files (x86)\Webroot"
    }
    else {
        $wrpath = "C:\Program Files\Webroot"
    }

    $wrexist = Test-Path $wrpath
    if ($wrexist -eq $true) {
        return $true
    }
    else {
        return $false
    }
}


Function Uninstall {
    if (-not (Check_Webroot)) {
        Write-Output "Webroot is not installed on this system."
        Exit;
    }
    $url = "http://update.itsupport247.net/webroot/UnstWebroot.exe"
    $downloadpath = "C:\Windows\Temp\UnstWebroot.exe"
    $ErrorActionPreference = "SilentlyContinue"
    $wc = New-Object System.Net.WebClient
    $wc.DownloadFile("$url", "$downloadpath")
    $pinfo = New-Object System.Diagnostics.ProcessStartInfo
    $pinfo.FileName = "C:\Windows\Temp\UnstWebroot.exe"
    $pinfo.RedirectStandardError = $true
    $pinfo.RedirectStandardOutput = $true
    $pinfo.UseShellExecute = $false
    $pinfo.Arguments = "Uninstall"
    $p = New-Object System.Diagnostics.Process
    $p.StartInfo = $pinfo
    $p.Start() | Out-Null
    $p.WaitForExit()
    $result = $p.ExitCode
    if($result -eq 0){
        Write-Output "Webrrot has been unistalled from the system."   
    }
    Else{
        Write-Output "Failed to uninstall webroot."
    }
    Remove-Item "C:\Windows\Temp\UnstWebroot.exe"
}


Function Webroot_Install($key) {
    $url = "http://update.itsupport247.net/webroot/wsasme.msi"
    $downloadpath = "C:\Windows\Temp\wsasme.msi"
    $wc = New-Object System.Net.WebClient
    $wc.DownloadFile("$url", "$downloadpath")
    
    $pinfo = New-Object System.Diagnostics.ProcessStartInfo
    $pinfo.FileName = "msiexec"
    $pinfo.RedirectStandardError = $true
    $pinfo.RedirectStandardOutput = $true
    $pinfo.UseShellExecute = $false
    $pinfo.Arguments = "/i  C:\Windows\Temp\wsasme.msi GUILIC=$key /qn"
    $p = New-Object System.Diagnostics.Process
    $p.StartInfo = $pinfo
    $p.Start() | Out-Null
    $p.WaitForExit()
    $result = $p.ExitCode
    if($result -eq 0) {
        Write-Output "Webroot has been installed to the system."  
    }
    Else{
        Write-Output "Failed to install webroot."
    }
    Remove-Item "C:\Windows\Temp\wsasme.msi"
}


switch($operation){
    "Install"{Webroot_Install -key $keycode}
    "Uninstall"{Uninstall}
}
