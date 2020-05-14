
<#
Script Name : On demand shutdown system 
Category : Maintenance

    .Synopsis
        On demand shutdown system.
    .Author
        Santosh.Dakolia@continuum.net
    .Name 
        On demand shutdown system.
#>


$computername= $env:computername

Try{
$OS = get-wmiobject win32_operatingsystem -computername $computername
$OS.psbase.Scope.Options.EnablePrivileges = $true
$OS.win32shutdown(5) | Out-Null
Write-Host "System : $computername Shutting down...."

}Catch{
    Write-Error "Error occured..!! $_.Exception.Message"
}

<#

0 (0x0)
Log Off - Logs the user off the computer. Logging off stops all processes associated with the security context of the process that called the exit function, logs the current user off the system, and displays the logon dialog box.

4 (0x4)
Forced Log Off (0 + 4) - Logs the user off the computer immediately and does not notify applications that the logon session is ending. This can result in a loss of data.

1 (0x1)
Shutdown - Shuts down the computer to a point where it is safe to turn off the power. (All file buffers are flushed to disk, and all running processes are stopped.) 
Users see the message, It is now safe to turn off your computer.

During shutdown the system sends a message to each running application. The applications perform any cleanup while processing the message and return True to indicate that they can be terminated.

5 (0x5)
Forced Shutdown (1 + 4) - Shuts down the computer to a point where it is safe to turn off the power. 
(All file buffers are flushed to disk, and all running processes are stopped.) Users see the message, It is now safe to turn off your computer.

When the forced shutdown approach is used, all services, including WMI, are shut down immediately. 
Because of this, you will not be able to receive a return value if you are running the script against a remote computer.

2 (0x2)
Reboot - Shuts down and then restarts the computer.

6 (0x6)
Forced Reboot (2 + 4) - Shuts down and then restarts the computer.

When the forced reboot approach is used, all services, including WMI, are shut down immediately. Because of this, 
you will not be able to receive a return value if you are running the script against a remote computer.

8 (0x8)
Power Off - Shuts down the computer and turns off the power (if supported by the computer in question).

12 (0xC)
Forced Power Off (8 + 4) - Shuts down the computer and turns off the power (if supported by the computer in question).

When the forced power off approach is used, all services, including WMI, are shut down immediately. Because of this, 
you will not be able to receive a return value if you are running the script against a remote computer.

Reserved [in]
A means to extend Win32Shutdown. Currently, the Reserved parameter is ignored.

#>
