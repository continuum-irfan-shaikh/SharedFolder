# $action = 'remove' # 'remove'
# $WindowTitle = 'Test'
# $Message = 'Test Message'

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

try {
    $ErrorActionPreference = 'stop'
    if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64" -and $env:PROCESSOR_ARCHITECTURE -eq 'x86') { $Query = 'c:\windows\sysnative\query.exe' }else { $Query = 'c:\windows\System32\query.exe' }
    $LoggedOnUsers = if(($Users = (& $Query user))){
        $Users | ForEach-Object {(($_.trim() -replace ">" -replace "(?m)^([A-Za-z0-9]{3,})\s+(\d{1,2}\s+\w+)", '$1  none  $2' -replace "\s{2,}", "," -replace "none", $null))} |
                 ConvertFrom-Csv |
                 Where-Object {$_.state -ne 'Disc'} |
                 Select-Object -expandproperty username
    }
    if ($LoggedOnUsers) {
        $registry = 'HKLM:\Software\Microsoft\Windows\CurrentVersion\Policies\System'
        function verify {
            return (Get-ItemProperty $registry -Name legalnoticecaption -ErrorAction SilentlyContinue ) -and (Get-ItemProperty $registry -Name legalnoticetext -ErrorAction Silentlycontinue)
        }
        
        switch ($action) {
            'add' {
                Set-ItemProperty $registry -Name legalnoticecaption -Value $WindowTitle	
                Set-ItemProperty $registry -Name legalnoticetext -Value $Message
                If (Verify) { Write-Output "Custom legal notice added successfuly." }else { Write-Output "Failed to add Custom legal notice." }
            }
            'remove' {
                Remove-ItemProperty $registry -Name legalnoticecaption -ErrorAction SilentlyContinue
                Remove-ItemProperty $registry -Name legalnoticetext -ErrorAction SilentlyContinue
                If (!$(verify)) { Write-Output "Custom legal notice removed successfuly." }else { Write-Output "Failed to remove Custom legal notice." }
            }
        }

    }
    else { Write-Output "This script requires logon user and currently no user is logged in. No action will be performed."; exit; }
}
catch { Write-Error $_ }
