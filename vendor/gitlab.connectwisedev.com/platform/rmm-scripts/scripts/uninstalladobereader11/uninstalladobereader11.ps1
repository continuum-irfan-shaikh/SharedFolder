<#
    name : Uninstall Adobe Reader-11
    Ctegory : Application
    .Script
    Uninstall Adobe Reader from the system.
    .Description
    Script will uninstall Adobe Reader from the system
    User has to specify the version which needs to uninstall.
    if version specified is not present in the system then script will exit.
    script will take two parameters "Action","Version"
    .Author
    Nirav Sachora
    .Requirements
    Script should run with Administrator privileges. 

    Versions:
        11.0.00
        11.0.10
        11.0.9
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

$action = "Uninstall"
$Adobe = Get-WmiObject -Class Win32_Product | ? {$_.Name -like "Adobe Reader*"}

if ($Adobe -eq $null) {
    Write-Output "Adobe Reader is not installed on this system."
    Exit;
}
elseif (($Adobe.Version) -ne $version) {
    "Adobe Reader version {0} is not installed on this system." -f $version
    Exit;
}
else {
    $result = $Adobe.Uninstall()
    if ($result.ReturnValue -eq 0) {
        Write-Output "Adobe Reader has been uninstalled from the system."
        Exit;
    }
    else {
        Write-Error "Failed to uninstall Adobe Reader."
    }
}
