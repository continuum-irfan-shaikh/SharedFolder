<#
    .SYNOPSIS
       Retrieve processor information
    .DESCRIPTION
       Retrieve processor information
    .Author
       Santosh.Dakolia@continuum.net
    .Reference 
        https://docs.microsoft.com/en-us/windows/desktop/cimwin32prov/win32-processor    
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


$Processer = Get-WmiObject -Class win32_processor -ComputerName $computername

$OutputObj = New-Object -TypeName PSobject

    $outputobj | Add-Member -MemberType NoteProperty -Name Manufacturer -Value $Processer.Manufacturer
    $OutputObj | Add-Member -MemberType NoteProperty -Name Name -Value $Processer.Name
    $OutputObj | Add-Member -MemberType NoteProperty -Name Description -Value $Processer.Description
    $OutputObj | Add-Member -MemberType NoteProperty -Name ProcessorID -Value $Processer.ProcessorId
    $OutputObj | Add-Member -MemberType NoteProperty -Name AddressWidth -Value $Processer.AddressWidth
    $OutputObj | Add-Member -MemberType NoteProperty -Name DataWidth -Value $Processer.DataWidth
    $OutputObj | Add-Member -MemberType NoteProperty -Name Family -Value $Processer.Family
    $OutputObj | Add-Member -MemberType NoteProperty -Name MaxClockSpeed -Value $Processer.MaxClockSpeed
    $outputobj | fl Manufacturer, Name, Description, ProcessorID, Addresswidth, Datawidth, Family, MaxClockSpeed


}Catch{
    Write-Error "Error occured while retrieving Data..!! $_.Exception.Message"
}
