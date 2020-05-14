
<#
    .Script
    Configure on-demand shutdown.
    .Author
    Nirav Sachora
    .Description
    Script will shutdown the system after checking the conditions provided
    .Requirements
    Script should run with admin privileges.
#>
#$shutdowntype = ask/notify

<#[string]$processname = "powershell"
[bool]$donotinitiate = $true
[bool]$executecmd = executecmd/executeprocess
    [string]$cmd
    [string]$process = "notepad.exe"
    [string]$parameter
 

[string]$displaymsg = "Hello"
$timeout = 10
$processtimeout = 1#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}


[int]$conf

$os = (Get-WMIObject Win32_OperatingSystem).Caption
if ($os -like "*server*") {
    $loggedinusers = (quser) -replace '\s{2,}', ','
    if (!$loggedinusers.count) {
        Write-Output "No users are logged on to this system."
        Exit;
    }
    elseif ($loggedinusers.count -gt 2) {
        Write-Output "Multiple users are logged on to this machine.`nWindows cannot be shutdown."
        Exit;
    }
}

if ($timeout -gt 30) {
    Write-Error "Maximum timeout allowed is 30 minutes"
    Exit;
}

if ($processname -ne $null) {
    $processdetails = Get-Process -Name $processname -ErrorAction "SilentlyContinue"
    if ($processdetails -ne $null) {
        if ($donotinitiate -eq $true) {
            Write-Output "$processname is running in the memory`n`nWindows cannot be shutdown."
            Exit;
        }
        elseif ($executecmd -eq "executecmd") {
            try {
                $pinfo = New-Object System.Diagnostics.ProcessStartInfo
                $pinfo.FileName = "$cmd"
                $pinfo.RedirectStandardError = $true
                $pinfo.RedirectStandardOutput = $true
                $pinfo.UseShellExecute = $false
                if (!$parameter) {
                    $pinfo.Arguments = "$parameter"
                }
                $p = New-Object System.Diagnostics.Process
                $p.StartInfo = $pinfo
                $ErrorActionPreference = "Stop"
                $p.Start()
                $ErrorActionPreference = "Continue"
                if (!$time) { $exitstatus = $p.WaitForExit() }
                else {
                    $exittime = $processtimeout * 60000
                    $exitstatus = $p.WaitForExit($exittime)
                    if ($exitstatus) { $conf = 1 }
                    else { $p.Kill() }
                }
            }
            catch {
                Write-Error "Process could not be started"
                Exit;
            }
        }
    }
}

elseif ($conf -eq 1) {
    $wshell = New-Object -ComObject Wscript.Shell -ErrorAction Stop
    if ($shutdowntype -eq "ask") {
        $userinput = $wshell.Popup("$displaymsg`nPop-up will auto close after $timeout minutes ", $timeout * 60, "Shutdown Confirmation", 32 + 1)
        switch ($userinput) {
            1 { Stop-Computer }
            2 { Write-Output "User has declined the request to shutdown computer."; Exit; }
            -1 { Write-Output "No Confirmation from user."; Exit; }
        }
    }
    if ($shutdowntype -eq "notify") {
        $userinput = $wshell.Popup("$displaymsg`n`nPress OK or Close the dialogue box to shutdown`n`nSystem will automatically shutdown after $timeout minutes", $timeout * 60, "Shutdown Notification", 64 + 0)
        switch ($userinput) {
            1 { Stop-Computer }
            -1 { Stop-Computer }
        }
    }
}
