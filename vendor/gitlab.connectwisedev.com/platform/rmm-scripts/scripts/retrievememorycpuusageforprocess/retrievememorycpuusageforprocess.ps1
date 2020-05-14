## $MemoryUtilization = 10
# $CPUUtilization = 5

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

try {
    $ErrorActionPreference = 'stop'
    $NumOfCores = (Get-WmiObject Win32_Processor).NumberOfLogicalProcessors
    $Process = Get-WmiObject Win32_PerfFormattedData_PerfProc_Process |
               Select-Object -property Name,
                                       @{Name = "CPU"; Expression = { ($_.PercentProcessorTime / $NumOfCores) } }, 
                                       @{Name = "PID"; Expression = { $_.IDProcess } },   
                                       @{"Name" = "Memory(MB)"; Expression = { [int]($_.WorkingSetPrivate / 1mb) } } |  
               Where-Object { $_.Name -notmatch "^(idle|_total|system)$" }

    if ($MemoryUtilization -or $CPUUtilization) {
        $Process = $Process | Where-Object { $_.'Memory(MB)' -ge $MemoryUtilization -and $_.CPU -ge $CPUUtilization }
        If(!$Process){Write-Output "`nNo process found matching the specified criteria";exit}
    }
    Write-Output "`nList of processes matching the specified criteria"
    ForEach ($Item in $Process) {
        Write-Output "`nProcessId   : $($Item.PID)"
        Write-Output "Name        : $($Item.Name)"
        Write-Output "CPU         : $([Math]::Round(($Item.CPU),2))"
        Write-Output "Memory (MB) : $([Math]::Round(($Item.'Memory(MB)'),2))"
    }
}
catch {
    Write-Error $_
}
