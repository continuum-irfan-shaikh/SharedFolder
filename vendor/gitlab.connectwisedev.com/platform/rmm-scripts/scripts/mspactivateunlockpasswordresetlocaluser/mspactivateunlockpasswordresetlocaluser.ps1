<#
.SYNOPSIS
Developed By : GRT [HD-Automation] , Continuum
File Vesion : 2.0
.DESCRIPTION
The script is designed to change the password of local user whose account is locked/disabled or both also change it to Enable and Unlock if Local User account is Locked/Disabled or both.
.PARAMETER 
LocalUserID is the variable required for changing the password.
.PARAMETER 
ChangePassword is the password which needs to be set as password for the user account.
.PARAMETER
ChangePwdNextLogon is a boolean. True means to set Change Password on next logon.
.Example
$LocalUserID = "hd1"
$ChangePassword = "ED3%V"
$ChangePwdNextLogon=$false
#>
<# Compatibility Check if found incompatible will exit#>
try{
    $OSVersion=([system.environment]::OSVersion.Version).Major
    $PSVersion =(Get-Host).version
    if(($OSVersion -lt 6) -or ($PSVersion.Major -lt 2))
    {
    Write-Output "[MSG : Incompatible System Either OS is below windows 7/2008R2 or powershell version less than 2.0]"
    Exit 
    }
    <# else{Write-output "[MSG: Good to Go]"} #>
    }catch{ Write-Output "[MSG: $_.Exception.message]"}
    $OutputMessage = @() <# Variable to capture Output for debugging purpose.#>
     <# If $DC is equall to 4(backup DC) or 5(primary DC) as there is no local user on Domain Controller #>   
            try{
                 $DC = (Get-WmiObject Win32_ComputerSystem).domainrole
                }catch{ Write-error " Unable to Check if machine is Domain Controller $_.Exception.message" }
            if($DC -eq "4" -or $DC -eq "5")
                {   
                    $OutputMessage += "No Local User on Domain Controller."                
                }
            else
                {
                    <#Check if User Exist#>
                   try {
                    $Localuser =  Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='True' and Name='$LocalUserID'" | Select PSComputername, LocalAccount, Status,Disabled, AccountType, Lockout, PasswordRequired, PasswordChangeable, PasswordExpires, SID
                     }catch{ Write-error " Unable to check get Local user information $_.Exception.message" }
                    if ($Localuser)
                     {
                        try{
                                $MachineName = $env:COMPUTERNAME
                                $UsrCommand =[ADSI] "WinNT://$MachineName/$LocalUserID, user"
                                <# If account is Lockedout#>
                                if($Localuser.Lockout -eq $True)
                                {
                                    $UsrCommand.IsAccountLocked = $False
                                    $UsrCommand.Setinfo() 
                                }
                                <# If account is Disabled#>                      
                                if($Localuser.Disabled -eq $True) 
                                {
                                    $UserAccountProperties = Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='True' and Name='$LocalUserID'" #| %{$_.disabled=$false;$_.put()} | out-null
                                    $UserAccountProperties.Disabled = $false
                                    $UserAccountProperties.Put() | out-null
                                } 
                                <# reset the password#>
                                $UsrCommand.SetPassword($ChangePassword)
                                $UsrCommand.Setinfo()
                                $OutputMessage += "Password Changed Successfully."
                                <# If User needs to change at next logon#>
                                if($ChangePwdNextLogon -eq $true)
                                {
                                    if(($Localuser.PasswordChangeable -eq $True) -And ($Localuser.PasswordExpires -eq $true))
                                    {
                                    <# Use to check for force pwd change at next logon. If PasswordChangeable and PasswordExpires is False the option or Password change at next logon is grayed out.#>
                                        $UsrCommand.PasswordExpired = 1  
                                        $UsrCommand.Setinfo() 
                                    }
                                    Else 
                                    {
                                         $OutputMessage += "User Can't Change password at next logon because Either it is Set to never Expired or Cannot Change Password."
                                    }
                                }     
                                else
                                {
                                   $OutputMessage += "Next Logon Password Change is not configrued by Technician."
                                }
                              }catch {Write-error " Unable to change password $_.Exception.message"}
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
                            }catch{ Write-error " Unable to get list of Local User Account " }
                        }
                }
    Write-Output $OutputMessage
