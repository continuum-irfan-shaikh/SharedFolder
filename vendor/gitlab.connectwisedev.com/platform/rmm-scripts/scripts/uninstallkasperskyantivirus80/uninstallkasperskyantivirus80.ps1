<#   
    Name: Kaspersky Anti-Virus 8.0 for Windows Servers Enterprise Edition
    Category : Application
    .Description
    Script will uninstall Kaspersky Antivirus from the system
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
$kaspersky = Get-WmiObject -Class Win32_Product | ? {$_.Name -like "Kaspersky Anti-Virus*"}

if ($kaspersky -eq $null) {
    Write-Output "Kaspersky Antivirus is not installed on this system."
    Exit;
}
elseif (($kaspersky.Version) -ne $version) {
    "Kaspersky Antivirus version {0} is not installed on this system." -f $version
    Exit;
}
else {
    $result = $kaspersky.Uninstall()
    if ($result.ReturnValue -eq 0) {
        Write-Output "Kaspersky Antivirus has been uninstalled from the system."
        Exit;
    }
    else {
        Write-Error "Failed to uninstall Kaspersky Antivirus."
    }
}
