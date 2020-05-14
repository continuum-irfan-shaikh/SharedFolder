<#
   .Script
   Disable Apple update messages
   .Author
   Nirav Sachora
   .Requirements
   Script should run with highest privileges
   .Description
   Script will disable the task which is created in task scheduler "AppleSoftwareUpdate".
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

$pre = $erroractionpreference

$erroractionpreference = 'silentlycontinue'

$service = new-object -ComObject("Schedule.Service")
$service.Connect()
$rootFolder = $service.GetFolder("\Apple")
$task = $rootFolder.GetTask("AppleSoftwareUpdate")
if (!$task) {
    Write-Output "AppleSoftwareupdate is not installed in this system `n"
    $erroractionpreference = $pre
    Exit;
}
else {
    $task.Enabled = $false
    $task = $rootFolder.GetTask("AppleSoftwareUpdate")

    if ($task.Enabled -eq $false) {
        Write-Output "AppleSoftwareUpdate task has been disabled"
    }
    else {
        Write-Error "Operation Failed, run scipt with highest privileges"
    }
    $erroractionpreference = $pre
}
