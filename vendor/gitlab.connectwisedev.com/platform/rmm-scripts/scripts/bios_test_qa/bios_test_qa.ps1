$biosInfo = Get-WmiObject win32_bios
$infoObject = New-Object PSObject
Add-Member -inputObject $infoObject -memberType NoteProperty -name "ReleaseDate" -value $biosinfo.ConvertToDateTime($biosinfo.ReleaseDate)
Add-Member -inputObject $infoObject -memberType NoteProperty -name "BuildNumber" -value $biosInfo.BuildNumber
Add-Member -inputObject $infoObject -memberType NoteProperty -name "CurrentLanguage" -value $biosInfo.CurrentLanguage
Add-Member -inputObject $infoObject -memberType NoteProperty -name "InstallableLanguages" -value $biosInfo.InstallableLanguages
Add-Member -inputObject $infoObject -memberType NoteProperty -name "Manufacturer" -value $biosInfo.Manufacturer
Add-Member -inputObject $infoObject -memberType NoteProperty -name "Name" -value $biosInfo.Name
Add-Member -inputObject $infoObject -memberType NoteProperty -name "PrimaryBIOS" -value $biosInfo.PrimaryBIOS
Add-Member -inputObject $infoObject -memberType NoteProperty -name "SerialNumber" -value $biosInfo.SerialNumber
Add-Member -inputObject $infoObject -memberType NoteProperty -name "SMBIOSBIOSVersion" -value $biosInfo.SMBIOSBIOSVersion
Add-Member -inputObject $infoObject -memberType NoteProperty -name "SMBIOSMajorVersion" -value $biosInfo.SMBIOSMajorVersion
Add-Member -inputObject $infoObject -memberType NoteProperty -name "SMBIOSMinorVersion" -value $biosInfo.SMBIOSMinorVersion
Add-Member -inputObject $infoObject -memberType NoteProperty -name "SMBIOSPresent" -value $biosInfo.SMBIOSPresent
Add-Member -inputObject $infoObject -memberType NoteProperty -name "Status" -value $biosInfo.Status
Add-Member -inputObject $infoObject -memberType NoteProperty -name "Version" -value $biosInfo.Version

Write-Output $infoObject
