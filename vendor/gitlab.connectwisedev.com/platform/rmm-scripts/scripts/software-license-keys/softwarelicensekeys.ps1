$States = @{0 = "Unlicensed"; 1 = "Licensed"; 2 = "OOBGrace"; 3 = "OOTGrace"; 4 = "NonGenuineGrace"; 5 = "Notification"; 6 = "ExtendedGrace"}
$Products = (Get-WmiObject -Class SoftwareLicensingProduct -Property Name,ApplicationID,LicenseStatus,ProductKeyID)
ForEach ($Product in $Products){
	$LicenseStatus = $States.Item([int]$Product.LicenseStatus)
	$Product | Format-List -Property Name,ApplicationID,ProductKeyID,@{Label="LicenseStatus"; Expression = {$LicenseStatus}}
}