<#

    $action = "restart" # "start", ""stop", "restart"
    $processName = "notepad++" # any process name

#>
switch ($action) {
    "start" {
        try {
            Start-Process -FilePath $processName -ErrorAction Stop
            Write-Output "Process $processName successfully started." 
        }
        catch {
            if ( $_.Exception.message -match "The system cannot find the file specified.") {
                Write-Error "The system cannot find the file/process name $processName specified."
        
            }
            Else { Write-Error "Error occured while statrting the $processName : $_.Exception.message" }
        }
        
    }

    "stop" {
        try {
            if ($forceStop -eq $true) {
                Stop-Process -Name $processName -Force -EA Stop
                Write-Output "Process $processName forcefully stopped." 
            }
            Else {
                Stop-Process -Name $processName -EA Stop
                Write-Output "Process $processName successfully stopped!!."       
            }
        }
        catch {
            Write-error "Error occured while stopping the $processName  : $_.Exception.message"
      
        }
    }

    "restart" {
        if ( $null -eq $(Get-Process $processName -ErrorAction SilentlyContinue ) ) {
            return "Process $processName currently not running."
        }

        try {
            if ($forceStop -eq $true) {
                Stop-Process -Name $processName -Force -EA Stop
                Write-Output "Process $processName forcefully stopped." 
            }
            Else {
                Stop-Process -Name $processName -EA Stop
                Write-Output "Process $processName successfully stopped."       
            }
        }
        catch {
            Write-error "Error occured while stopping the $processName  : $_.Exception.message"
      
        }

        Start-Sleep -s 5
        try {
            Start-Process -FilePath $processName -ErrorAction Stop
            Write-Output "Process $processName successfully restarted." 
        }
        catch {
            if ( $_.Exception.message -match "The system cannot find the file specified.") {
                Write-Error "The system cannot find the file/process name $processName specified."
        
            }
            Else { Write-Error "ERROR : $_.Exception.message" }
        }


    }
    default {
        Write-Error "Invalid action: $action. Should be one of: start, stop, restart"
        return
    }
}
