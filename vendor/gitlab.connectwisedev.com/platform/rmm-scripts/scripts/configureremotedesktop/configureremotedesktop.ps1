$RegPath = "HKLM:\SYSTEM\CurrentControlSet\Control\Terminal Server"
If(-Not(Test-Path -Path $RegPath)){
    New-Item -Path $RegPath -Force | New-ItemProperty -Name "fDenyTSConnections" -PropertyType DWord -Value 0 | Out-Null
}
Switch ($Action){
	"Enable"{
		Set-ItemProperty -Path $RegPath -Name "fDenyTSConnections" -Value 0
		$Val = Get-ItemProperty -Path $RegPath -Name "fDenyTSConnections"
		if($Val.fDenyTSConnections -eq 0){
			Return "Successfully Enabled Remote Desktop"
		} Else {
			Write-Error "Error While Enabling Remote Desktop"
            Return
		}
	}
	"Disable"{
		Set-ItemProperty -Path $RegPath -Name "fDenyTSConnections" -Value 1
		$Val = Get-ItemProperty -Path $RegPath -Name "fDenyTSConnections"
		if($Val.fDenyTSConnections -eq 1){
			Return "Successfully Disabled Remote Desktop"
		} Else {
			Write-Error "Error While Disabling Remote Desktop"
            Return
		}
	}
}
