# $EventLog = 'System' # user input

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

$ErrorActionPreference = 'Stop'
try{
    $All = Get-EventLog -List -OutVariable obj | ForEach-Object{$_.LogDisplayName}
        if(!($All -contains $EventLog)){ Write-Output "The EventLog: `'$EventLog`' doesn't exist on the system." ;continue}
        if(![bool]($obj|Where{$_.log -eq $EventLog}| Select-Object -expand Entries)){ Write-Output "The EventLog: `'$EventLog`' is already clear with 'zero' entries." ;continue}
        Clear-EventLog -LogName $EventLog -Confirm:$false
        if($?){
            Write-Output "Succesfuly cleared Event log: `'$EventLog`'"
        }
        else{
            Write-Output "Failed to clear Event log: `'$EventLog`'"
        }
}
catch{
    $_
}
