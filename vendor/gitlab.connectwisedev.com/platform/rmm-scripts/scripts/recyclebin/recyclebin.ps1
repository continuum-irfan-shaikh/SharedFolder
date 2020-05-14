$length = 0
$files = 0
$shell = New-Object -Com Shell.Application
$shell.NameSpace(0xA).Items() | Remove-Item -Recurse

Try{
$drives = Get-WmiObject Win32_LogicalDisk -Filter "DriveType=3"
ForEach ($disk in $drives) {
    if (Test-Path '$($disk.DeviceID)\Recycle') {
        $content = Get-ChildItem "$($disk.DeviceID)\Recycle" -Force -Recurse -ErrorAction SilentlyContinue  | Where-Object { -not $_.PSIsContainer } | Measure-Object -property length -Sum
        Remove-Item "$($disk.DeviceID)\Recycle" -Force -Recurse -ErrorAction SilentlyContinue
        if ($content){
		    $length += $content.Sum
		    $files += $content.Count
	    }
	}
	else {
		$content = Get-ChildItem "$($disk.DeviceID)\`$Recycle.Bin" -Force -Recurse -ErrorAction SilentlyContinue  | Where-Object { -not $_.PSIsContainer } | Measure-Object -property length -Sum
        Remove-Item "$($disk.DeviceID)\`$Recycle.Bin" -Force -Recurse -ErrorAction SilentlyContinue
        if ($content){
		    $length += $content.Sum
		    $files += $content.Count
	    }
	}
}

   Write-Output "Cleared $($length) bytes, $($files) files"

}

Catch{
   Write-Output "Cleared $($length) bytes, $($files) files Recycle Bin Empty"
}
