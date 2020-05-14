$args = @($args)
$drives = @($drives)

[bool]$FixErrors = $false
[bool]$VigorousIndexCheck = $false
[bool]$SkipFolderCycle = $false
[bool]$ForceDismount = $false
[bool]$RecoverBadSectors = $false
[bool]$OkToRunAtBootUp = $false

if ($args -contains "FixErrors") {$FixErrors = $true}
if ($args -contains "VigorousIndexCheck") {$VigorousIndexCheck = $true}
if ($args -contains "SkipFolderCycle") {$SkipFolderCycle = $true}
if ($args -contains "ForceDismount") {$ForceDismount = $true}
if ($args -contains "RecoverBadSectors") {$RecoverBadSectors = $true}
if ($args -contains "OkToRunAtBootUp") {$OkToRunAtBootUp = $true}

for ($i = 0; $i -le $drives.length-1; $i++) {
	$drives[$i] = $drives[$i] + ":"
}

$localDiskType = 3
$drivesObjects = Get-WmiObject win32_volume | Where-Object { $_.DriveType -eq $localDiskType -and $_.DriveLetter -ne $null}
$availableDrives = @()
ForEach ($obj in $drivesObjects) {
	$availableDrives += $obj.DriveLetter
}

ForEach ($drive in $drives) {
	if ($availableDrives -contains $drive) {
		$driveObj = Get-WmiObject win32_volume | Where-Object {$_.DriveLetter -eq $drive}

		#Chkdsk method of the Win32_Volume class according to https://msdn.microsoft.com/en-us/library/aa384915(v=vs.85).aspx
		$result = $driveObj.chkdsk($FixErrors,
							$VigorousIndexCheck,
							$SkipFolderCycle,
							$ForceDismount,
							$RecoverBadSectors,
							$OkToRunAtBootUp)
		$resultCode = $result.ReturnValue
		If ($resultCode -ge 0 -and $resultCode -le 5) {
			Switch ($resultCode) {
				0{Write-Output "Drive $drive Success - Chkdsk Completed"}
				1{Write-Output "Drive $drive Success - Volume Locked and Chkdsk Scheduled on Reboot"}
				2{Write-Error "Drive $drive Failure - Unsupported File System"}
				3{Write-Error "Drive $drive Failure - Unknown File System"}
				4{Write-Error "Drive $drive Failure - No Media In Drive"}
				5{Write-Error "Drive $drive Failure - Unknown Error"}
			}
		}
		Else {
			Write-Error "Drive $drive Failure - *Invalid Result Code* == $resultCode"
			Continue
		}
	}
	Else {
		Write-Error "Drive $drive does not exist"
		Continue
	}
}
