
<#
    .Script
    LogMeIn host software update.
    .Author
    Nirav Sachora
    .Description
    Script will update LogmeIn Host software.
    .Requirements
    Script should run with highest privileges.
#>

$Erroractionpreference = 'SilentlyContinue'
$exist = get-wmiobject win32_product | ?{$_.Name -eq "LogMeIn"}
if(!$exist){
$Erroractionpreference = 'Continue'
Write-Error "LogMeIn host software is not installed on this system."
Exit;
}


$currentversion = (get-wmiobject win32_product | ?{$_.Name -eq "LogMeIn"}).Version
$bits = Get-WmiObject -Class win32_Operatingsystem | Select -ExpandProperty OSArchitecture

if ($bits -eq "32-bit") {
    
    $execute = "C:\PROGRA~1\LogMeIn\x86\update\raupdate.exe"
    
}
else {

    $execute = "C:\PROGRA~2\LogMeIn\x64\update\raupdate.exe"

}

if (Test-path $execute) {

    $Erroractionpreference = "Stop"
    try {

        $pinfo = New-Object System.Diagnostics.ProcessStartInfo
        $pinfo.FileName = "$execute"
        $pinfo.RedirectStandardError = $true
        $pinfo.RedirectStandardOutput = $true
        $pinfo.UseShellExecute = $false
        $pinfo.Arguments = "/s"
        $p = New-Object System.Diagnostics.Process
        $p.StartInfo = $pinfo
        $p.Start() | Out-Null
        $p.WaitForExit()
       
        $updatedversion = (get-wmiobject win32_product | ?{$_.Name -eq "LogMeIn"}).Version
        if (($p.ExitCode -eq 0) -and ($currentversion -ne $updatedversion)) {
            Write-Output "LogMeIn Host Software has been updated successfully."
            Exit;
        }
        elseif(($p.ExitCode -eq 0) -and ($currentversion -eq $updatedversion)){
            Write-Output "LogmeIn version seems to be already updated`nCurrent version:$currentversion."
        }
        else {
            Write-Error "LogMeIn host software update failed."
            Exit;
        }
    }
    catch {
        Write-Error "LogMeIn host software update failed `n$_"
        Exit;
    }
}
else {

    Write-Error "LogMeIn Host software update .exe file not found"
    Exit;

}  

