<#
    .SYNOPSIS
        Uninstall Internet Explorer.
    .DESCRIPTION
        Uninstall Internet Explorer. Degrade Internet Explorer to default version
    .Help
        To get more details refer below command. 
        Dism /Online /Get-Packages
        Dism /Online /Remove-Package /packagename:$packagename 
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>
$action = 'uninstall'
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

if (($action -eq "uninstall") -and ($version)) {
    
    $ErrorActionPreference = 'SilentlyContinue'
    
    $ieversion = ([System.Version][System.Diagnostics.FileVersionInfo]::GetVersionInfo("$env:ProgramFiles\Internet Explorer\iexplore.exe").ProductVersion).Major

    if ($version -eq $ieversion) {

        if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64"-and $env:PROCESSOR_ARCHITECTURE -eq 'x86') {
            $Query = 'C:\windows\sysnative\dism.exe'
        }
        else {
            $Query = 'C:\windows\System32\dism.exe'
        } 

        #Get installed updates for IE
        $packages = Dism /Online /Get-Packages | Select-String "InternetExplorer"
        if ($packages) {
            foreach ($package in $packages) {

                $packagename = $package.ToString().Replace("Package Identity : ", '')

                #Process Object
                $uninstallinfo = New-object System.Diagnostics.ProcessStartInfo
                $uninstallinfo.CreateNoWindow = $true
                $uninstallinfo.UseShellExecute = $false
                $uninstallinfo.RedirectStandardOutput = $true
                $uninstallinfo.RedirectStandardError = $true
                $uninstallinfo.FileName = "$Query"
                $uninstallinfo.Arguments = "/Online /Remove-Package /packagename:$packagename /norestart"
                $uninstall = New-Object System.Diagnostics.Process
                $uninstall.StartInfo = $uninstallinfo
                [void]$uninstall.Start()
                $uninstall.WaitForExit()

            }
        }
        else {
            Write-Output "Can not degrade default Internet Explorer $version."
            Exit;
        }
        Start-Sleep 5
        $check = Dism /Online /Get-Packages | Select-String "InternetExplorer"
        
        if ($check) {
            Write-Output "Failed to degrade internet explorer $version on system $ENV:COMPUTERNAME. Exit Code: $($uninstall.exitcode)"
        }
        else {
        
            Write-Output "Internet Explorer $version is degraded now." 
            write-Output "Reboot the system $ENV:COMPUTERNAME to make the changes effective."
        }

    }
    else {
        Write-Output "Internet Explorer $version not installed on system $ENV:COMPUTERNAME"
    }
}
else {
    Write-Output "Uninstall/Version input selection missing."
}
