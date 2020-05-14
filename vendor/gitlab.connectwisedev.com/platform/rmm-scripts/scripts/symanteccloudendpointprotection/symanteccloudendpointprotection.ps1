<#
    .Script
    Uninstall Symantec cloud from the system.
    .Author
    Nirav Sachora
    .Description
    Script will uninstall Symantec cloud software from the system.
    .Requirements
    Script should run with admin privileges.
#>
$restart = $true

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

    Function Execute_process($path) {
        $pinfo = New-Object System.Diagnostics.ProcessStartInfo
        $pinfo.FileName = "C:\Program Files\Symantec.cloud\PlatformAgent\Uninstall.exe"
        $pinfo.RedirectStandardError = $true
        $pinfo.RedirectStandardOutput = $true
        $pinfo.UseShellExecute = $false
        $pinfo.Arguments = "/S"
        $p = New-Object System.Diagnostics.Process
        $p.StartInfo = $pinfo
        $p.Start() | Out-Null
        $p.WaitForExit()
        return $p.ExitCode
    }

    if((Test-path "C:\Program Files\Symantec.cloud\PlatformAgent\Uninstall.exe") -eq $false){
        Write-Output "Symantec cloud is not installed on this system."
        Exit;
    }

    $result = Execute_process -path $path
    
    if(($result -eq 0) -and ($restart -eq $true)){
        Write-Output "Symantec cloud has been uninstalled from the system."
        Restart-computer -Force
    }
    elseif(($result -eq 0) -and ($restart -eq $false)) {
        Write-Output "Symantec cloud has been uninstalled from the system.`nPlease manually restart the system."
    }
    else{
        Write-Output "Error while uninstalling Symantec.cloud."
    }

