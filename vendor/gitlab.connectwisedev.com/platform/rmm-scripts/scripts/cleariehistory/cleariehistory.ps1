$dataToClear = @($dataToClear)

if ($dataToClear -contains "history") {$History = $true}
if ($dataToClear -contains "cookies") {$Cookies = $true}
if ($dataToClear -contains "TempIEFiles") {$TempIEFiles = $true}

$MajorOsVersion = ([environment]::OSVersion.Version).Major
$MinorOsVersion = ([environment]::OSVersion.Version).Minor

function removeItems {
	Param($itemsPath, $itemsName)

    $items = Get-ChildItem -Path $itemsPath -Exclude *.dat -Recurse -Force -ErrorAction SilentlyContinue
	if ($items -ne $null) {
        $count = $items.Length
        foreach ($item in $items) {
            Remove-item $item.FullName -Recurse -Force -ErrorAction SilentlyContinue
        }
	} else {
	    $count = 0
	}
	Write-Output "Deleted $itemsName : $count elements."
}

function clearData {
    Param($userName, $userSID)

    $HistoryPath = "$env:SystemDrive\Users\$userName\AppData\Local\Microsoft\Windows\History"
    If (($MajorOsVersion -eq 6) -And ($MinorOsVersion -eq 1)){
        $CookiesPath = "$env:SystemDrive\Users\$userName\AppData\Roaming\Microsoft\Windows\Cookies"
        $TempFilesPath = "$env:SystemDrive\Users\$userName\AppData\Local\Microsoft\Windows\Temporary Internet Files"
    } ElseIf ($MajorOsVersion -eq 10){
        $TempFilesPath = "$env:SystemDrive\Users\$userName\AppData\Local\Microsoft\Windows\INetCache"
        $CookiesPath = "$env:SystemDrive\Users\$userName\AppData\Local\Microsoft\Windows\INetCookies"
    }

    if ($History) {
        if (-not (Get-PSDrive -name "HKU" -ErrorAction SilentlyContinue)) {
                    New-PSDrive -Name HKU -PSProvider Registry -Root HKEY_USERS >$null
        }

        $regPath = "HKU:\$userSID\Software\Microsoft\Internet Explorer\TypedURLs"
        Get-Item -literalPath $regPath -ErrorAction SilentlyContinue | ForEach-Object {
            $key = $_
            $key.GetValueNames() | ForEach-Object {
                $value = $_
                Remove-ItemProperty -Path $regPath -Name $value
            }
        }
        removeItems -itemsPath $HistoryPath -itemsName "history"
		Invoke-Expression -Command "$env:SystemDrive\Windows\System32\RunDll32.exe InetCpl.cpl, ClearMyTracksByProcess 1"
    }

    if ($Cookies) {
    	removeItems -itemsPath $CookiesPath -itemsName "cookies"
    }

    if ($TempIEFiles) {
    	removeItems -itemsPath $TempFilesPath -itemsName "temporary files"
    }
}

$PatternSID = 'S-1-5-21-\d+-\d+\-\d+\-\d+$'
$users = @()
$users = Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList\*' | Where-Object { $_.PSChildName -match $PatternSID } |
    Select  @{ name = "SID"; expression = { $_.PSChildName } },
			@{ name = "Name"; expression = { $_.ProfileImagePath -replace '^(.*[\\\/])', '' } }

foreach ($user in $users) {
    $uName = $user.Name
    $uSID = $user.SID
    Write-Output "User $uName :"
    clearData -userName $uName -userSID $uSID
    Write-Output "`n"
}
