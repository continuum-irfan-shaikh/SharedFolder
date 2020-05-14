#$SizeInMB = 10 # user input

$ErrorActionPreference = 'SilentlyContinue'
$UserPaths = Get-WmiObject win32_userprofile | Where-Object { $_.sid -like "S-1-5-21*"} | Select-Object -ExpandProperty localpath
$result = ForEach ($UserPath in $UserPaths) {
    If (Test-Path $UserPath) {
        $Size = 0
        $Size = (Get-ChildItem $UserPath -Recurse -Force | Where-Object {!$_.PSIsContainer} | select -ExpandProperty Length | Measure-Object -Sum).sum
        
        # only list profiles greater in size as per the user input
        if ($Size -gt $($SizeInMB * 1MB)) { 
            # choose the right unit (GB/MB/KB) depending upon the size
            if ($Size -ge 1GB) { $SizeWithUnit = "{0:N2} GB" -f ($Size / 1gb) }elseif ($Size -ge 1Mb) { $SizeWithUnit = "{0:N2} MB" -f ($Size / 1mb)  }else {$SizeWithUnit = "{0:N2} KB" -f ($Size / 1kb)} 
            $username = (Split-Path $UserPath -Leaf)
            New-Object -TypeName PSOBject -Property @{
                User = $(if ($env:USERNAME -eq $username) {"$username (Current logged-on User)"}else {$username})
                Size = $SizeWithUnit
                Path = $UserPath
            }
        }
    }
}

if ($result) {
    Write-Output "`nFollowing are the user profiles greater than '$SizeInMB MB'"
    $result |ForEach-Object {
        Write-Output "`nUser : $($_.User)"
        Write-Output "Path : $($_.Path)"
        Write-Output "Size : $($_.Size)"
    }
}
else {
    Write-Output "`nNo user profiles found which are greater than '$SizeInMB MB'"
}

