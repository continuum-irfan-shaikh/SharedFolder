if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

if ( $mode -eq "delayed-auto") {
    $modename = "Automatic(Delayed Start)"
}
if ( $mode -eq "Auto") {
    $modename = "Automatic"
}
if ( $mode -eq "demand") {
    $modename = "Manual"
}
if ( $mode -eq "disabled") {
    $modename = "Disabled"
}

if ( ((Get-WmiObject win32_service -Filter "name='$service'").StartMode -eq "Disabled" ) -and ($action -like "*start*") -and ($mode -notmatch "auto`|demand`|delayed-auto")) {
        Write-Output "The specified service is disabled and cannot be started/restared."
        Exit
}

if ($mode) {
    $result = sc.exe config $service start= $mode
    if ($result -like "*SUCCESS*") {
        "Service start mode changed to $modename."
        if ($mode -eq "disabled" -and $action -like "*start*") {
            Write-Output "The specified service is disabled and cannot be started/restarted."
            Exit  
        }
    }
    else {
        "ERROR: Failed to change service mode to $modename.`n$result"
    }
}


try {
    $status = (get-service $service).Status
    switch ($action) {
        "start" {
            if ($status -ne "Running") {
                Start-Service $service -ErrorAction Stop
            }
            else {
                Write-Output "Service is already running."
                Get-WmiObject win32_service -Filter "name='$service'" | select Name, Status, State, StartMode | fl
                Return
            }
        }
        "restart" {
            Restart-Service $service -ErrorAction Stop
        }
        "stop" {
            if ($status -ne "Stopped") {
                Stop-Service $service -ErrorAction Stop
            }
            else {
                Write-Output "Service is already stopped."
                Get-WmiObject win32_service -Filter "name='$service'" | select Name, Status, State, StartMode | fl
                Return
            }
        }
        Default { 
            Write-Output "Invalid Action: $action. Should be one of: Start, Stop, Restart"
        }
    }
    if ($?) {
        $res = Get-WmiObject win32_service -Filter "name='$service'" | select Name, Status, State, StartMode | fl
        if ($action -eq "stop") {
            Write-Output "Service stopped successfully."
            $res
        }
        else {
            Write-Output "Service $action`ed successfully."
            $res
        }
    }
}
catch {

    if ( $($_.Exception.message)  -like "*it has dependent services*" ) {
     Write-Output "The specified service cannot be stopped, it has following dependent services." 
     $dep = Get-Service $service -DependentServices | select @{N = "Dependent services:"; E = { $_.DisplayName }} | ft
     $res = Get-WmiObject win32_service -Filter "name='$service'" | select Name, Status, State, StartMode | fl
     Write-Output $dep,$res
    
    }else { Write-Output "ERROR: $_.Exception.message" }

}
