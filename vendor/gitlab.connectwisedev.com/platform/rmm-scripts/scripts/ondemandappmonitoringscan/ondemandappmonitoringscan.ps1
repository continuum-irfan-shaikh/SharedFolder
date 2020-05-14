
<#
    .Script
    On Demand App Monitoring scan
    .Author
    Nirav Sachora
    .Requirements
    Script should run with highest privileges.
    zsrvdet.exe file should be present at "C:\PROGRA~2\SAAZOD" or "C:\PROGRA~1\SAAZOD"
    script should run on server OS only.
#>


function Execute_exe($path) {
    if (Test-Path "$path") {
        $pinfo = New-Object System.Diagnostics.ProcessStartInfo
        $pinfo.FileName = $path
        $pinfo.RedirectStandardError = $true
        $pinfo.RedirectStandardOutput = $true
        $pinfo.UseShellExecute = $false
        $p = New-Object System.Diagnostics.Process
        $p.StartInfo = $pinfo
        $p.Start() | Out-Null
        $p.WaitForExit()
        return $p.ExitCode
    }
    else {
        Write-Error "zsrvdet.exe' is not present on this server"
    }
    
}

$oscheck = Get-WmiObject -Class win32_operatingsystem | select -ExpandProperty Caption
if ($oscheck -notlike "*Server*") {
    Write-Output "Script can run on Server Operating Systems only"
    Exit;
}

$os = Get-WmiObject -Class win32_operatingsystem | select -ExpandProperty OSArchitecture
if ($os -eq '64-bit') {
    $filepath = "C:\PROGRA~2\SAAZOD\zsrvdet.exe"
}
else {
    $filepath = "C:\PROGRA~1\SAAZOD\zsrvdet.exe"
}

$result = Execute_exe -path $filepath

if ($result -eq 0) {
    Write-Output "App Monitoring Scan has been completed successfully"
}
else {
    Write-Error "Operation Failed"
}
