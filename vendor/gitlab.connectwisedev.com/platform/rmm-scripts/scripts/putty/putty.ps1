
<#
    .Script
    Install/Uninstall Putty.
    .Description
    Script will uninstall Putty from the system
    User has to specify the version which needs to uninstall.
    if version specified is not present in the system then script will exit.
    script will take two parameters "Action","Version"
    .Author
    Nirav Sachora
    .Requirements
    Script should run with Administrator privileges. 
    $version = "0.71"
#>
$Action = "install"
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

$puttypath = "C:\Program Files\PuTTY\putty.exe"
$ptpath = "C:\Program Files\PuTTY\unins000.exe"
$registrypath = "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Putty*"


Function Check_putty {
    $puttyexist = Test-Path $ptpath
    if ($puttyexist -eq $true) {
        return $true
    }
    else {
        return $false
    }
}

$bit = (Get-WmiObject Win32_operatingsystem).Osarchitecture

<#Function Uninstall {
    if ((Get-ItemProperty $registrypath -Name "DisplayVersion" | select -ExpandProperty "DisplayVersion") -ne $version) { Write-Output "Version of putty is not installed on this system"; Exit; }

    if (Check_putty -eq $true) {
        $ErrorActionPreference = "SilentlyContinue"
        & "taskkill.exe" /f /im putty.exe | Out-Null
        Remove-Item $puttypath
        & "$ptpath" /VERYSILENT /SUPPRESSMSGBOXES /NORESTART /SP-
        $ErrorActionPreference = "Continue"
    }
    else {
        Write-Output "Putty is not installed on this system."
        Exit;
    }


    if (Check_putty -eq $false) {
        Write-Output "Putty has been uninstalled from the system."
        Exit;
    }
    else {
        Write-Output "Failed to install putty from the system"
        Exit;
    }
}#>

Function Install {
    if(Check_putty){
        Write-Output "Putty is already installed on this system."
        Exit;
    }
    if ($bit -eq "64-bit") { $url = "https://the.earth.li/~sgtatham/putty/latest/w64/putty-64bit-0.71-installer.msi" }
    if ($bit -eq "32-bit") { $url = "https://the.earth.li/~sgtatham/putty/latest/w32/putty-0.71-installer.msi" }
    $downloadpath = "C:\Windows\Temp\putty.msi"
    $ErrorActionPreference = "SilentlyContinue"
    $wc = New-Object System.Net.WebClient
    $wc.DownloadFile("$url", "$downloadpath")
    if (!(Test-Path $downloadpath)) {
        Write-Output "putty download has been failed."
        Exit;
    }
    & msiexec.exe /i  C:\Windows\Temp\putty.msi /qn | Out-Null

    $Erroractionpreference = "Continue"
    if (Test-Path $puttypath) {
        Write-Output "Putty has been installed on this system."
        Exit;
    }
    else {
        Write-Output "Failed to install putty on this system."
        Exit;
    }
    Remove-Item $downloadpath | Out-Null
}

if($action -eq "Install"){
    Install
}
elseif($action -eq "Uninstall"){
    Uninstall
}

