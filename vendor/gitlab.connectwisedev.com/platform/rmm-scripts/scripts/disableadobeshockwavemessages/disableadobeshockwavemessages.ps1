$RegPaths = "HKLM:\SOFTWARE\Adobe","HKLM:\SOFTWARE\WOW6432Node\Adobe"
$ErrorCount = 0
ForEach ($RegPath in $RegPaths){
	If(-Not(Test-Path -Path $RegPath)){
		$ErrorCount++
	} Else {
		$RegKeys += (Get-ChildItem $RegPath -Recurse) | Where-Object {$_.Name -Match "Shockwave" -And ($_.Name -Match "AutoUpdate")}
	}
}
If (($ErrorCount -eq $RegPaths.Count) -Or ($RegKeys -eq $null)){
	Write-Output "Shockwave is not installed"
	Return
} Else {
	ForEach ($RegKey in $RegKeys){
		Set-ItemProperty -Path Registry::$RegKey -Name "(default)" -Value ""
		$Result = (Get-ItemProperty -Path Registry::$RegKey)
		Write-Output $Result
	}
}
