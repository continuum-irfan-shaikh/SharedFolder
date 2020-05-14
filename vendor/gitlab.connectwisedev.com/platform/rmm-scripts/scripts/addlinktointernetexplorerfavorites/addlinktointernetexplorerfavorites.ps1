# $DisplayName = 'Technet' # userinput
# $URL = 'https://microsoft.com' # userinput

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

Function Create-IEFavorite ($DisplayName, $URL) {
    if($env:PROCESSOR_ARCHITEW6432 -eq "AMD64" -and $env:PROCESSOR_ARCHITECTURE -eq 'x86') {$Query = 'c:\windows\sysnative\query.exe'}else{$Query = 'c:\windows\System32\query.exe'}
    $LoggedOnUsers = if(($Users = (& $Query user))){
        $Users | ForEach-Object {(($_.trim() -replace ">" -replace "(?m)^([A-Za-z0-9]{3,})\s+(\d{1,2}\s+\w+)", '$1  none  $2' -replace "\s{2,}", "," -replace "none", $null))} |
                 ConvertFrom-Csv |
                 Where-Object {$_.state -ne 'Disc'} |
                 Select-Object -expandproperty username
    }

    if ($LoggedOnUsers) {
        foreach($Item in $LoggedOnUsers){
            $LoggedOnUserProfilePath = Get-WmiObject win32_userprofile | Where-object {$_.SID -like "S-1-5-21*" -and $_.localpath -like "*\$Item" } | Select-Object -ExpandProperty localpath
            $IEFavouritePath = "$LoggedOnUserProfilePath\Favorites"
            if (Test-Path  $IEFavouritePath) {
                $Shell = New-Object -ComObject WScript.Shell
                $FullPath = Join-Path -Path $IEFavouritePath -ChildPath "$($DisplayName).url"
                $Shortcut = $Shell.CreateShortcut($FullPath)
                $Shortcut.TargetPath = $Url
                $Shortcut.Save()
                $ContentCheck = Get-Content $FullPath | Where-Object {$_ -like "URL=*"}
                if ((Test-Path $FullPath) -and ($ContentCheck -like "*URL=$URL*")) {
                    Write-Output "Successfuly added Link:`'$URL`' with Name:`'$DisplayName`' for User: `'$Item`'"
                }
                else {
                    Write-Output "Failed to add link:`'$URL`' with Name:`'$DisplayName`' for User: `'$Item`'"
                }
            }
            else {
                Write-Error "Path: `'$IEFavouritePath`' doesn't exist. Failed to add link. for User: `'$Item`'"
            }
        }
    }
    else { Write-Output 'No users logged-on the machine. No action will be performed.'; exit; }
}
    
$ErrorActionPreference = 'stop'
try {
    Create-IEFavorite -DisplayName $DisplayName -URL $URL
}
catch {
    Write-Error $_
}
