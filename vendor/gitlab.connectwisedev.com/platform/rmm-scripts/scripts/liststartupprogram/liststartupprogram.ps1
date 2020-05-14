try
{
    $StartupProgram = Get-WMIobject Win32_StartupCommand | Select-Object Name, Description, command, Location, User, SettingID | Format-List
    if(!$StartupProgram)
    {Write-Output "No Startup Programs are Enabled On this System"}
    else
    {Write-Output $StartupProgram }
}
catch
{
    Write-Error "WMI is not handleing Call properly : $($_.Exception.Message)"
    Exit
}
