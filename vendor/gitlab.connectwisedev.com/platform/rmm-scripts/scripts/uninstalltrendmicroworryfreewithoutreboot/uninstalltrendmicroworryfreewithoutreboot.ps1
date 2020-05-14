<#
    .Script
    Uninstall Trend Micro Worry-Free Business Security Agent from the system.
    .Description
    Script will uninstall Trend Micro Worry-Free Business Security Agent from the system
    User has to specify the version which needs to uninstall.
    if version specified is not present in the system then script will exit.
    script will take one parameters ,"Version"
    .Author
    Nirav Sachora
    .Requirements
    Script should run with Administrator privileges. 
#>
<#

Versions

8.0
9.0

#>

#[string]$version

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

$Trend = Get-WmiObject -Class Win32_Product | ?{($_.Name -like "Trend Micro Worry-Free Business Security Agent") -and ($_.Version -eq $version)}
if(!$Trend){
    Write-Output "Trend Micro Worry-Free Business Security Agent is not installed on this system."
    Exit;
}

    $url = "http://dcmdwld.itsupport247.net/trend_uninstall.zip"
    $downloadpath = "C:\Windows\Temp\trend_uninstall.zip"
    $ErrorActionPreference = "SilentlyContinue"
    $wc = New-Object System.Net.WebClient
    $wc.DownloadFile("$url", "$downloadpath")
    $ErrorActionPreference = "Continue"
    
    $verifytools = Test-path "C:\Windows\Temp\trend_uninstall.zip"
    if(!$verifytools){
    Write-Error "Uninstallation binaries is missing from the system."
    Exit;
    }
    
    
    $shell = new-object -com shell.application
    $shell.Namespace("C:\Windows\Temp").copyhere($shell.NameSpace("C:\Windows\Temp\trend_uninstall.zip").Items(),4)
    
    
    $verifytools = Test-path "C:\Windows\Temp\Trend_Uninstall"
    if(!$verifytools){
    Write-Error "Uninstallation binaries is missing from the system."
    Exit;
    }
    
    $pinfo = New-Object System.Diagnostics.ProcessStartInfo
    $pinfo.FileName = "C:\Windows\Temp\Trend_Uninstall\Uninstall.bat"
    $pinfo.CreateNoWindow = $true;
    $pinfo.RedirectStandardError = $true
    $pinfo.RedirectStandardOutput = $true
    $pinfo.UseShellExecute = $false
    $p = New-Object System.Diagnostics.Process
    $p.StartInfo = $pinfo
    $p.Start() | Out-Null
    $p.WaitForExit()
    $result = $p.ExitCode

    if($result -eq 0){
    Write-Output "Trend Micro Worry-Free Business Security Agent has been uninstalled from the system."
    }
    else{
    Write-Output "Failed to uninstall Trend Micro Worry-Free Business Security Agent."
    }

    Remove-item "C:\Windows\Temp\trend_uninstall.zip" -Force
    Remove-item "C:\Windows\Temp\Trend_uninstall" -recurse 
