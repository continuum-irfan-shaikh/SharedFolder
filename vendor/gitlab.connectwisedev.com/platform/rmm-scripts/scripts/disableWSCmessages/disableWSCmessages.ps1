$MajorOsVersion = ([environment]::OSVersion.Version).Major
$MinorOsVersion = ([environment]::OSVersion.Version).Minor


If (($MajorOsVersion -eq 6) -And ($MinorOsVersion -eq 1)){
	$RegPath = "HKLM:\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer"
	$KeyName = "HideSCAHealth"
} elseif(($MajorOsVersion -eq 6) -And ($MinorOsVersion -eq 2)){

	$RegPath = "HKLM:\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer"
	$KeyName = "HideSCAHealth"

} Elseif (($MajorOsVersion -eq 6) -And ($MinorOsVersion -eq 3)){

	$RegPath = "HKLM:\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer"
	$KeyName = "HideSCAHealth"
        New-Item -Path $RegPath -Force | Out-Null

} ElseIf ($MajorOsVersion -eq 10){
	$RegPath = "HKCU:\Software\Policies\Microsoft\Windows\Explorer"
	$KeyName = "DisableNotificationCenter"
}


If(-Not(Test-Path -Path $RegPath)){
	New-Item -Path $RegPath -Force | New-ItemProperty -Name $KeyName -PropertyType DWord -Value 1 | Out-Null
	}
Set-ItemProperty -Path $RegPath -Name $KeyName -Value 1
Return Get-ItemProperty $RegPath


<#
    


https://superuser.com/questions/470212/change-action-center-settings-via-command-line


You can hide the Action Center icon with the following command:

"reg add HKLM\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer /v HideSCAHealth /t REG_DWORD /d 0x1"


It won't disable messages (to do that you'd have to edit a binary registry value), but they won't be displayed anymore.





Script is not design to Handle OS : Win8 32/64 , Win8.1 32/64  As the Major and Minor Version is ::

Win8  Major : 6  and Minor : 2
Win8.1 Major : 6 and Minor : 3


#>
