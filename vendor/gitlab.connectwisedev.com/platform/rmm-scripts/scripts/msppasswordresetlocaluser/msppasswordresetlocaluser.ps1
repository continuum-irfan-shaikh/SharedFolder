<#
	.SYNOPSIS
	Developed By : GRT [HD-Automation] , Continuum
	File Vesion : 2.0
	.DESCRIPTION
	The script is designed to change the password of local user, if Local User account is Locked/Disabled or both it will not do anything.
	.PARAMETER 
	LocalUserID is the variable required for changing the password.
	.PARAMETER 
	ChangePassword is the password which needs to be set as password for the user account.
	.PARAMETER
	ChangePwdNextLogon is a boolean. True means to set Change Password on next logon.
	.Example for Test
	$LocalUserID = "hd1"
	$ChangePassword = "ED3%V"
	$ChangePwdNextLogon=$true
#>
<# Compatibility Check if found incompatiblt with exit#>
try{
    $OSVersion=([system.environment]::OSVersion.Version).Major
    $PSVersion =(Get-Host).version
    if(($OSVersion -lt 6) -or ($PSVersion.Major -lt 2))
    {
    Write-Output "[MSG : Incompatible System Either OS is below windows 7/2008R2 or powershell version less than 2.0]"
    Exit 
    }}catch{ Write-Output "[MSG: $_.Exception.message]"}
    
    $OutputMessage = @() <# Variable to capture Output for debugging purpose.#>
    <# If $DC is equall to 4(backup DC) or 5(primary DC) as there is no local user on Domain Controller #>   
    try{ $DC = (Get-WmiObject Win32_ComputerSystem).domainrole
    }catch{ Write-error "Unable to Check if machine is Domain Controller $_.Exception.message" }
    if($DC -eq "4" -or $DC -eq "5")
    { 
        $OutputMessage += "No Local User on Domain Controller."
        Write-Output $OutputMessage
        Exit
    }
    else
    {
    <#Check if User Exist#>
    try {
            $Localuser =  Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='True' and Name='$LocalUserID'" | Select PSComputername, LocalAccount, Status,Disabled, AccountType, Lockout, PasswordRequired, PasswordChangeable, PasswordExpires, SID
        }catch{ Write-error " Unable to retrieve Local user information" }
        if ($Localuser)
        {
            if(($Localuser.Disabled -eq $False) -And ($Localuser.Lockout -eq $False))
            {
                $MachineName = $env:COMPUTERNAME
                <# reset the password#>
                try{
                    $UsrCommand =[ADSI] "WinNT://$MachineName/$LocalUserID, user"
                    $UsrCommand.SetPassword($ChangePassword) 
                    if($ChangePwdNextLogon -eq $true)
                    {
                        if(($Localuser.PasswordChangeable -eq $True) -And ($Localuser.PasswordExpires -eq $true))
                        {
                        <# Use to check for force pwd change at next logon. If PasswordChangeable and PasswordExpires is False the option or Password change at next logon is grayed out.#>
                            $UsrCommand.PasswordExpired = 1  
                            $UsrCommand.Setinfo() 
                        }
                        Else 
                        {   $OutputMessage += "Password is Either Set to never Expired or Cannot Change Password."  }
                     }
                     else
                     {  $OutputMessage += "Next Logon Password Change is not configrued by Technician."   }
                     $OutputMessage += "Password Changed Successfully."
                     }catch {Write-error " Unable to change password $_.Exception.message"}
              }
         else
             { $OutputMessage +="Account Status Locked : "+$Localuser.Lockout 
               $OutputMessage +="Account Status Disabled :"+$Localuser.Disabled
             }
                        }
                    else
                        {
                            try{
                            $Localuser =  Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='True'" | Select Name,Disabled,Lockout
                            for($i=0; $i -lt $Localuser.length ; $i++)
                            {   $TempUser=$Localuser[$i].Name
                                $useraccountexpiration = [adsi]"WinNT://./$TempUser, user" | Select AccountExpirationDate
                                $OutputMessage += "Name :"+$Localuser[$i].Name 
                                $OutputMessage += "Disabled :"+$Localuser[$i].Disabled
                                $OutputMessage += "Lockout :"+$Localuser[$i].Lockout
                                $OutputMessage += "AccountExpirationDate :"+$useraccountexpiration.AccountExpirationDate
                                $OutputMessage += " "
                            }
                            }catch{ Write-error "Unable to get list of Local User Account." }
                       }
                }
    Write-Output $OutputMessage |fl
