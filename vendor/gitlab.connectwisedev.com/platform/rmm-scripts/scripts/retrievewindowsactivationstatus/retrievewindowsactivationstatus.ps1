<#
    .Script
    Retrieve Windows activation status.
    .Author
    Nirav Sachora.
    .Description
    Script will fetch winows activation status using WMI class SoftwareLicensingProduct and display result.
    .Requirements
    Script should run with highest privileges.
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

$ErrorActionPreference = "SilentlyContinue"
$status = Get-WmiObject -Class SoftwareLicensingProduct -Filter "Name like 'Windows%'" | where {$_.PartialProductKey -ne $null}  | Select Name,licensestatus,ProductKeyID,description,@{Name = "Evaluationenddate"; Expression ={$_.ConvertToDateTime($_.EvaluationEndDate)}},Graceperiodremaining
if($? -eq $false){
    $ErrorActionPreference = "Continue"
    Write-Error "Windows activation status could not be fetched."
    Exit;
}
if($status.licensestatus -ne 0){
    switch ($status.licensestatus) {
        1 { $licence = "Licensed, {0}" -f 1 ;break}
        2 { $licence = "OOBGrace, {0}" -f 2 ;break}
        3 { $licence = "OOTGrace, {0}" -f 3 ;break}
        4 { $licence = "NonGenuineGrace, {0}" -f 4 ;break}
        5 { $licence = "Notification, {0}" -f 5 ;break}
        6 { $licence = "ExtendedGrace, {0}" -f 6 ;break}
        Default {$licence = "Error fetching license status";break}
    }
    $evdate = {$_.ConvertToDatetime($status.Evaluationenddate)}
    Write-Output "Windows is activated on your system."
    $object = New-Object -TypeName psobject
    $object | Add-Member -MemberType NoteProperty -Name 'Name' -Value $status.Name
    $object | Add-Member -MemberType NoteProperty -Name 'Description' -Value $status.Description
    $object | Add-Member -MemberType NoteProperty -Name 'License Status' -Value $licence
    $object | Add-Member -MemberType NoteProperty -Name 'Product Key' -Value $status.ProductKeyID
    #$object | Add-Member -MemberType NoteProperty -Name 'Evaluation period end date' -Value $status.Evaluationenddate
    $object | Add-Member -MemberType NoteProperty -Name 'Remaining grace period' -Value $status.Graceperiodremaining
    Write-Output "`n"
    ($object | fl | Out-String).trim()
}
else {
    Write-Output "Windows is not activated on your system."
}
$ErrorActionPreference = "Continue"
