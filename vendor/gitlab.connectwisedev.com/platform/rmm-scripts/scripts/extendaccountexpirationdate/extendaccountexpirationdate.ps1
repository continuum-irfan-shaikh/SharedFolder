<#
	.SYNOPSIS
	Developed By : GRT [HD-Automation], Continuum
	File Version : 3.0
	.DESCRIPTION
	The script is designed to change the expiration data local or domain user.
	.PARAMETER
	UserID is the variable required for changing the expiration data.
	.PARAMETER
	IsDomain is a boolean value with TRUE as user concidered as Domain User else loccal user.
	.PARAMETER
	Expire is variable to get the data value on which account will gets deactivated. Example : If you have selected 20th Oct 2019 as expire date then the account will be active till 19th Oct 2018.

	$UserID = "hd"
	$IsDomain = $true
	$expire =  (Get-Date).AddDays(108)
#>
<# Architecture Check #>
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

[datetime]$expireinput = $expire
$OutputMessage = @() <# Variable to capture output for debugging purpose. #>
<# Compatibility check if found incompatible will exit #>
try{
	$OSVersion = ([system.environment]::OSVersion.Version).Major
	$PSVersion = (Get-Host).Version
	if(($OSVersion -lt 6) -or ($PSVersion.Major -lt 2))
	{
		Write-Output "[MSG: System is not compatible with the requirement. Either machine is below Windows 7 / Windows 2008R2 or Powershell version is lower than 2.0"
		Exit
	}
}catch { Write-Output "[MSG: ERROR : $_.Exception.message]"}
<# Compatibility check ends #>
If($IsDomain)
{
	try{
			$sysinfo = Get-WmiObject -Class Win32_ComputerSystem
	}catch{ Write-Error "Unable to check if machine is in Domain.  $_.Exception.message"}
	if($sysinfo.domain -eq "WORKGROUP")
	{
		$OutputMessage += "Machine Connected is not in Domain Environment."
		Write-Output $OutputMessage
		exit
	}
	else
	{
		$NameStatus = $UserID.Contains("@")
		if($NameStatus -eq $true)
		{
			$ProvidedName,$ProvidedDomainName = $UserID.split('@')
			if($ProvidedDomainName -eq $sysinfo.Domain)
			{
				$UserID = $ProvidedName
			}
			else
			{
				$OutputMessage += "[MSG: User and Machine belongs to different Domain/Sub-Domain , Try UserID only without domain name]"
				Write-Output $OutputMessage
				exit
			}
		}
		try{ $DC = (Get-WmiObject Win32_ComputerSystem).domainrole
		}catch{ Write-Error "Unable to check if machine is Domain Controller. $_.Exception.message"}
		if($DC -eq "4" -or $DC -eq "5")
		{
			try{ Import-Module ActiveDirectory -ErrorAction Stop
			}catch{ Write-Error "Error while importing AD Powershell module. Connect to other DC. $_.Exception.message"}
			try{	$ADuser = $null
					$ADuser = Get-ADuser -identity $UserID -properties Name,AccountExpirationDate | Select Name,AccountExpirationDate -ErrorAction Stop
			}catch [Microsoft.ActiveDirectory.Management.ADIdentityNotFoundException] { $OutputMessage += "[MSG: User Doesnot exist] FAILED" }
			catch{ Write-Error "$_.Exception.message"}
			if($ADuser)
			{
				try{
					Set-ADuser -Identity $UserID -AccountExpirationDate $expireinput
					$ADuser = Get-ADuser -identity $UserID -properties Name,SamAccountName,AccountExpirationDate | select Name,SamAccountName,AccountExpirationDate -ErrorAction Stop
					$OutputMessage += "Account Expiration Date is configured."
					$OutputMessage += $ADuser |fl
				}catch{ Write-Error "$_.Exception.message"}
			}
			else
			{
				try{
					<# User doesnot match, get the list of user based on first letter #>
					$ADuser = $UserID.Substring(0,1)+"*"
					$FilterStringDomain = "Name -like '$ADuser'"
					$ADUserListSimilarType = Get-ADuser -filter $FilterStringDomain | Select Name,SamAccountName,Enabled,AccountExpirationDate,LockedOut,PasswordExpired |fl
					$OutputMessage += "Below are the listof users starting with Letter :"+$UserID.Substring(0,1)
					$OutputMessage += $ADUserListSimilarType
				}catch { Write-Error "Unable to pull the list of users of similar name. $_.Exception.message"}	
			}	
		}
		else
		{
			try{ <# Collect list of DC in the environment #>
				$Command = "nltest /DSGETDC:"+ $sysinfo.domain
				$DCList = Invoke-expression $Command
				if($DCList)
				{
					$OutputMessage += "Machine is either not connected to DC or have network issue. Please connect to the mentioned List"
					$OutputMessage += $DCList |fl
				}
				else
				{
					$OutputMessage += "Machine is either not connected to DC or have network issue. Please connect admin or escalate."
				}
			}catch{ Write-Error "Unable to pull the list of Domain Controller. $_.Exception.message"}
		}
	}
}
else <# If provided user is a local user #>
{
	<# If $DC is 4 means BDC or 5 means PDC as there is no local user on Domain Controller #>
	try{ $DC = (Get-WmiObject Win32_ComputerSystem).domainrole
	}catch{ Write-Error "Unable to check if machine is Domain Controller. $_.Exception.message" }
	if($DC -eq "4" -or $DC -eq "5")
	{
		$OutputMessage += "No Local User on Domain Controller."
	}
	else
	{
		<# Check if User Exist #>
		try{
				$FilterStringLocal = "LocalAccount='True' and Name = '$UserID'"
				$Localuser = Get-WmiObject -Class Win32_UserAccount -Filter $FilterStringLocal | Select Name, PSComputername, LocalAccount, Status, Disabled,AccountType, Lockout, PasswordRequired, PasswordChangeable,PasswordExpires, SID
		}catch{ Write-Error "Unable to check Local User infomration. $_.Exception.message"}
		if($Localuser)
		{
			$useraccountconnection = [adsi]"WinNT://./$UserID, user"
			$useraccountconnection.Put("AccountExpirationDate",$expireinput)
			$useraccountconnection.SetInfo()
			$useraccountconnection.RefreshCache()
			$OutputMessage += "Updated Account Expiration Date is " + $useraccountconnection.AccountExpirationDate
		}
		else
		{
			try{
				$FilterStringList = "LocalAccount='True'"
				$Localuser = Get-WmiObject -Class Win32_UserAccount -Filter $FilterStringList | Select Name,Disabled,Lockout
				$OutputMessage += "[MSG: User Doesnot Exist] FAILED, Below are the list of users"
				for($i=0;$i -lt $Localuser.length; $i++)
				{
					$TempUser = $Localuser[$i].Name
					$useraccountexpiration = [adsi]"WinNT://./$TempUser, user" | Select AccountExpirationDate
					$OutputMessage += "Name : "+ $Localuser[$i].Name
					$OutputMessage += "Disabled: "+ $Localuser[$i].Disabled
					$OutputMessage += "Lockout : "+ $Localuser[$i].Lockout
					$OutputMessage += "AccountExpirationDate : "+ $useraccountexpiration.AccountExpirationDate
					$OutputMessage += " "
				}
			}catch{Write-Error "Unable to get list of local user account. $_.Expiration.message"}
		}
	}
}
Write-Output $OutputMessage |fl
