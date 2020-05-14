<#
    .SYNOPSIS
        Retrieve Terminal Service Mode 
    .DESCRIPTION
        Get startup mode of service termservice. It can be automatic/manual/disabled
    .Help
        To get more details refer below details
        Get-WmiObject win32_service
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>

try {

$value = Get-WmiObject win32_service -filter "name='termservice'" | select @{l='Service Name';e={$_.name}},@{l='StartUP Mode';e={$_.startmode}},State | fl -ErrorAction stop
    
    if(!$value)
    {
    Write-Output "No Terminal Service Found."
    }
    else {
    Write-Output $value
    }
}
catch {

Write-Error $_.Exception.Message

}

