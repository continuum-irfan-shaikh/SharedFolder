
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}



$ErrorActionPreference = 'SilentlyContinue'
$UserPaths = Get-WmiObject win32_userprofile | Where-Object { $_.sid -like "S-1-5-21*"} | Select-Object -ExpandProperty localpath
ForEach ($UserPath in $UserPaths) {
    $Size = 0
    $Size = (Get-ChildItem "$UserPath\Documents\" -Recurse -Force | Where-Object {!$_.PSIsContainer} | select -ExpandProperty Length | Measure-Object -Sum).sum
    if($Size -ge 1GB){ $SizeWithUnit = "{0:N2} GB" -f ($Size/1gb) }elseif($Size -ge 1Mb){ $SizeWithUnit = "{0:N2} MB" -f ($Size/1mb)  }else{$SizeWithUnit = "{0:N2} KB" -f ($Size/1kb)} 
    $username = (Split-Path $UserPath -Leaf)
    Write-Output "`nUser : $(if($env:USERNAME -eq $username){"$username (Current logged-on User)"}else{$username})"
    Write-Output "Path : $UserPath\Documents\"
    Write-Output "Size : $SizeWithUnit"
}
