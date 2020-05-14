if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}


try {
    if ((Get-WmiObject Win32_OperatingSystem).producttype -eq 1) {
        $ErrorActionPreference = 'stop'
        # makes sure the Windows Firewall service is running
        # otherwise, next step fails to set the firewall profiles to 'ON'    
        Start-Service -DisplayName '*Firewall*'
        $Status = Invoke-Command {
            NetSh Advfirewall set privateprofile state on 
            NetSh Advfirewall set publicprofile state on
            # intentionally excluded 'domain' profile 
            # as it needs to be set only on private and public profiles
        }
        If ($Status -like "*OK*") {
            Write-Output "Firewall is Enabled."
        }
        else {
            throw $Status
        }
    }
    else {
        Write-Output "`nOnly suppported on Work Station(s)"
    }
}
catch {
    Write-Error $_.Exception.Message
}
