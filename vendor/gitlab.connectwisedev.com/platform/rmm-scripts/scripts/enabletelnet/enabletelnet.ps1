<#
    .SYNOPSIS
        Enable Telnet Feature
    .DESCRIPTION
        Enable Telnet Feature. User can have access the other system remotely in same network. 
    .Help
        To get more details refer below command. 
        dism /Online /Enable-Feature /?
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}
$status = dism /online /Get-Featureinfo /FeatureName:TelnetClient | findstr "Enabled"

If ($?)
{
Write-Output "`nTelnet feature already enabled on system $Env:Computername"
}
else 
{
    #Enable TelnetClient
    try 
    {
    dism /online /Enable-Feature /FeatureName:TelnetClient 2>&1 | Out-Null
    If ($?)
    {
    $check = dism /online /Get-Featureinfo /FeatureName:TelnetClient | findstr "Enabled"
        if ($check -eq "State : Enabled")
        {
            Write-Output "`nTelnet feature enabled on system $Env:Computername"
        }
        else
        {
        Write-Error "`nSomething went wrong. Telnet feature not enabled on system $Env:Computername"
        }
    }
    }
    catch
    {
     write-output "`n"$_.Exception.Message
    } 
}
