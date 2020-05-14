<#
    .SYNOPSIS
        Remove Labtech Agent (ConnectWise Automate)
    .DESCRIPTION
        Remove Labtech Agent (ConnectWise Automate)
    .Help
        To get more details refer below command. 
        Get-WmiObject -Class Win32_Product
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}



$program = "LabTech Client"

$app = Get-WmiObject -Class Win32_Product | where-Object {$_.Name -match $program}
Try {
    If ($app -eq $null) {
        Write-Output "`nLabtech Agent(ConnectWise Automate) not installed on this system $ENV:ComputerName"
    }
    Else {
        $app.Uninstall()  2>&1 | Out-Null
    
        $check = Get-WmiObject -Class Win32_Product | where-Object {$_.Name -match $program}
        if ($check -eq $null) {
    
            Write-Output "`nLabtech Agent(ConnectWise Automate) removed from the system $ENV:ComputerName"
        }
        else {
            Write-Output "Something went wrong while uninstalling. Kindly check manually."
        }
    }
}
catch {
    write-output "`n"$_.Exception.Message
} 
