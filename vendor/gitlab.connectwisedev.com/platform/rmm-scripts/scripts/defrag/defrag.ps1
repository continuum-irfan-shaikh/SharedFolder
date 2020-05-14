$drives = @($drives)

for ($i=0;$i -le $drives.length-1;$i++) {
	$drives[$i] = $drives[$i] + ":"
}

$drivesObjects = Get-WmiObject win32_volume | Where-Object { $_.DriveType -eq 3 -and $_.DriveLetter -ne $null}
$availableDrives = @()
ForEach ($obj in $drivesObjects) {
	$availableDrives += $obj.DriveLetter
}

ForEach ($drive in $drives) {
	if ($availableDrives -contains $drive) {
		$result = Get-WmiObject win32_volume -filter "driveletter='$drive'" | Invoke-WmiMethod -Name Defrag | Select-Object -ExpandProperty ReturnValue
		switch ($result) {
			0  {$resultMsg = 'Success'; break}
			1  {$resultMsg = 'Access Denied'; break}
			2  {$resultMsg = 'Not Supported'; break}
			3  {$resultMsg = 'Volume Dirty Bit Set'; break}
			4  {$resultMsg = 'Not Enough Free Space'; break}
			5  {$resultMsg = 'Corrupt MFT Detected'; break}
			6  {$resultMsg = 'Call Cancelled'; break}
			7  {$resultMsg = 'Cancellation Request Requested Too Late'; break}
			8  {$resultMsg = 'Defrag In Progress'; break}
			9  {$resultMsg = 'Defrag Engine Unavailable'; break}
			10 {$resultMsg = 'Defrag Engine Error'; break}
			default {$resultMsg = 'Unknown Error'}
		}
		if ($result -eq 0) {
			Write-Output "Volume $drive $resultMsg"
		} else {
			Write-Error "Volume $drive $resultMsg"
		}
	} else {
		Write-Error "There is no disk with volume $drive"
	}
}
