
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

$OSBit = [IntPtr]::Size 


IF($OSBit -eq '4'){

$cfgDirPath32 = "$env:SystemRoot\System32\Macromed\Flash"
IF (-not (Test-Path -Path $cfgDirPath32)) {
			Write-Output "Adobe Flash Player is not installed"
				Return
}

$cfgFilePath32 = $cfgDirPath32 + "\mms.cfg"
IF (-not (Test-Path -Path $cfgFilePath32)) {
		Write-Output "Config file '$cfgFilePath32' was created"
			New-Item -Path $cfgDirPath32 -Name mms.cfg -ItemType "file" -Value "`n AutoUpdateDisable=1"
} else {
		Write-Output "Config file '$cfgFilePath32' was updated"
			(Get-Content -Path $cfgFilePath32) -Replace 'AutoUpdateDisable=.*$', "AutoUpdateDisable=1" | Set-Content -Path $cfgFilePath32
}
Write-Output (Get-Content -Path $cfgFilePath32)

}

IF($OSBit -eq '8'){

$cfgDirPath = "$env:SystemRoot\SysWOW64\Macromed\Flash"
IF (-not (Test-Path -Path $cfgDirPath)) {
		Write-Output "Adobe Flash Player is not installed"
			Return
}

$cfgFilePath = $cfgDirPath + "\mms.cfg"
IF (-not (Test-Path -Path $cfgFilePath)) {
		Write-Output "Config file '$cfgFilePath' was created"
			New-Item -Path $cfgDirPath -Name mms.cfg -ItemType "file" -Value "`n AutoUpdateDisable=1"
} else {
		Write-Output "Config file '$cfgFilePath' was updated"
			(Get-Content -Path $cfgFilePath) -Replace 'AutoUpdateDisable=.*$', "AutoUpdateDisable=1" | Set-Content -Path $cfgFilePath
}
Write-Output (Get-Content -Path $cfgFilePath)



}


function isSettingExists {
    Param($registryPath, $name)
	IF(Test-Path $registryPath -PathType container) {
		$key = Get-Item -LiteralPath $registryPath
			IF ($key.GetValue($name, $null) -ne $null) {
				return $true
        }
	}
	return $false
}

IF (-not (Get-PSDrive -name "HKU" -ErrorAction SilentlyContinue)) {
		New-PSDrive -Name HKU -PSProvider Registry -Root HKEY_USERS >$null
}

$PatternSID = 'S-1-5-21-\d+-\d+\-\d+\-\d+$'
$users = @()
$users = Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList\*' | Where-Object { $_.PSChildName -match $PatternSID } |
    Select  @{ name = "SID"; expression = { $_.PSChildName } },
			@{ name = "Username"; expression = { $_.ProfileImagePath -replace '^(.*[\\\/])', '' } }

ForEach ($user in $users) {
	$userSID = $($user.SID)
    IF (Test-Path "HKU:\$userSID\Software\Macromedia\FlashPlayerUpdate" -PathType container) {
			Remove-Item -Path "HKU:\$userSID\Software\Macromedia\FlashPlayerUpdate" -Recurse
			if ($?) {
					Write-Output "`n Registry setting 'HKU:\$userSID\Software\Macromedia\FlashPlayerUpdate' was found and deleted"
			} else {
					Write-Error "Error: have not managed to delete registry setting 'HKU:\$userSID\Software\Macromedia\FlashPlayerUpdate'"
		}
    }
    IF (isSettingExists -registryPath "HKU:\$userSID\Software\Microsoft\Windows\CurrentVersion\RunOnce" -name FlashPlayerUpdate) {
        Remove-ItemProperty -Path "HKU:\$userSID\Software\Microsoft\Windows\CurrentVersion\RunOnce" -Name FlashPlayerUpdate
		IF ($?) {
			Write-Output "`n Registry setting key 'HKU:\$userSID\Software\Microsoft\Windows\CurrentVersion\RunOnce > FlashPlayerUpdate' was found and deleted"
		} else {
			Write-Error "Error: have not managed to delete registry setting key 'HKU:\$userSID\Software\Microsoft\Windows\CurrentVersion\RunOnce > FlashPlayerUpdate'"
		}
    }
}


<# FlashInstall log file
 
32 Bit Windows:

C:\Windows\system32\Macromed\Flash\FlashInstall32.log
64 Bit Windows:

C:\Windows\system32\Macromed\Flash\FlashInstall64.log
C:\Windows\syswow64\Macromed\Flash\FlashInstall32.log
Note: Both log files are required from a 64-bit OS

The file names changed from FlashInstall.log to FlashInstall32.log and FlashInstall64.log in the Flash Player 26 release.  The contents of the FLashInstall.log file are appended to the new file(s).  If requested to provide log files and you have a FlashInstall.log file in the directory, provide that file as well
#>
