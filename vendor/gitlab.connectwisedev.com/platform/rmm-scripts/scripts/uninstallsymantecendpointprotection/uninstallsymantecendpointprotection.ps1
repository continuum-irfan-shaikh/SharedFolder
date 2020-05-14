
<#
    Name : Symantec Endpoint Protection
    Cetegory : Application

    .Description
        Script will uninstall Symantec Endpoint PRotection from the system
        User has to specify the version which needs to uninstall.
        if version specified is not present in the system then script will exit.
        script will take two parameters "Action","Version"
    .Author
       Nirav Sachora
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

$symantec = Get-WmiObject -Class Win32_Product | ? {($_.Name -eq "Symantec Endpoint Protection") -and ($_.Version -eq $version)}

if ($symantec -eq $null) {
    Write-Output "Provided version of Symantec Endpoint Protection is not installed on this system."
    Exit;
}
else {
    $result = $symantec.Uninstall()
    if ($result.ReturnValue -eq 0) {
        Write-Output "Symantec Endpoint Protection has been uninstalled from the system."
        Exit
    }
    else {
        Write-Error "Failed to uninstall Symantec Endpoint Protection"
    }
}

