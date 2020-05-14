<#
    .SYNOPSIS
       Retrieve Local Time on a Computer
    .DESCRIPTION
       Retrieve Local Time on a Computer
    .Author
       Santosh.Dakolia@continuum.net
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



Try{
$computername= $env:computername

$dt = gwmi win32_operatingsystem -computer $computername
$dt_str = $dt.converttodatetime($dt.localdatetime)#.ToLongTimeString()

$d=get-date
$aa=$d.datetime
$stime= [System.TimeZone]::CurrentTimeZone
$timezone = $stime.standardname

$OutputObj = New-Object -TypeName PSobject

    $outputobj | Add-Member -MemberType NoteProperty -Name MachineName -Value $computername
    $outputobj | Add-Member -MemberType NoteProperty -Name LocalTime -Value $dt_str 
    $outputobj | Add-Member -MemberType NoteProperty -Name TimeZone -Value $timezone
    $outputobj | FL MachineName, Localtime, TimeZone

    }Catch{
    Write-Error "Error occured while retrieving Data..!! $_.Exception.Message"
}
