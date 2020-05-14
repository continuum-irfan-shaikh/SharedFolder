# $Entries = @'
# 127.0.0.1 localhost # Loopback
# 192.168.0.1 test.com # external websites
# '@
# $OverwriteTheHostFile =  $true # user input

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}
$hostentries = @()
$failedentries = @()
$Entries = $Entries -split "`n"
foreach($entry in $entries){
        if($entry -match "\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}[ ]{1,}[^#][a-z.]{3,}"){
        $hostentries += $entry
        }
        else{$failedentries += $entry}
    }

try {
        $Path = $env:windir + "\System32\Drivers\etc\"
        $Filename = 'hosts'
        $Fullname = Join-Path $Path $Filename
        if (Test-Path $Fullname) {
            if ($OverwriteTheHostFile) {
                $NewName = $Filename + '_old'
                $i = 1
                While (Test-Path $(Join-path $Path $NewName)) {
                    $NewName = $Filename + "_old($i)"
                    $i++
                }
                Rename-Item $Fullname -NewName $NewName
                New-Item $Fullname -ItemType File | Out-Null
            }
            $hostentries | Out-File $Fullname -Append -Encoding ASCII
            if (($?) -and ($hostentries -ne $null)) {
                Write-Output "Following host entries has been done successfully. `n";$hostentries}
            else {'Failed to create host entries'}
            if($failedentries -ne $null){
            Write-Output "`nFailed to do host entry for following entries due to improper format. `n";$failedentries
            }
            
        }
        else {
            Write-Error "Unable to find Host file at path: `'$Fullname`'"
        }
}
catch {
    Write-Error $_
}
