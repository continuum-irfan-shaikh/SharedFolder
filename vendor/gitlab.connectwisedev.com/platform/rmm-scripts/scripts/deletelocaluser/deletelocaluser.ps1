function isValidUser {
	Param($userName, $users)
	foreach ($iUser in $users) {
		if ($iUser.Name -eq $user) {
			return $true
		}
	}
	return $false
}

function isLoggedOnUser {
	Param($userName)
	Get-WmiObject -Class Win32_LoggedOnUser | ForEach-Object {
		$loggedOnUser = ($_.Antecedent.ToString() -split '[\=\"]')[5]
		if ($loggedOnUser -eq $userName) {
			return $true
		}
	}
	return $false
}

$users = Get-WmiObject -Class Win32_UserAccount
if (-not (isValidUser -userName $user -users $users)) {
	return Write-Error "$user profile is not recognized"
}

if (isLoggedOnUser -userName $user) {
	return Write-Error "Please log off $user for deleting"
}

net user $user /delete >$null
Write-Output "$user was deleted successfully"
