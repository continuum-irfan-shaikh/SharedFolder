<# .SYNOPSIS
	Developed By : GRT [HD-Automation] , Continuum
	File Vesion : 1.0

	.DESCRIPTION
	The script is designed to change the locked/disabled to Unlock/Enable if Local User account is Locked/Disabled or both based on technicians selection.

	.PARAMETER 
	AccountUserID is the variable required for changing the password.

	.PARAMETER 
	IsDomain is the variable required for knowing if we are working on Domain Account or Local Account. If True means Domain Account . 

	.PARAMETER 
	AccountUnlock is a boolean. True means to unlock the account.
	 
	.PARAMETER
	AccountEnable is a boolean. True means to enable the account.
	.Example
	$AccountUserID = "test"
	$IsDomain = $false
	$AccountUnlock = $true
	$AccountEnable=$true
#>
<# Architecture check and if 32bit powershell open in 64 bit OS will take care of it. #>
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

<# Compatibility Check if found incompatible will exit#>
try{
    $OSVersion=([system.environment]::OSVersion.Version).Major
    $PSVersion =(Get-Host).version
    if(($OSVersion -lt 6) -or ($PSVersion.Major -lt 2))
    {
        Write-Output "[MSG : Incompatible System Either OS is below windows 7/2008R2 or powershell version less than 2.0]"
        Exit 
    }
}catch{ Write-Output "[MSG: $_.Exception.message]"}
$OutputMessage = @() <# Variable to capture Output for debugging purpose.#>
 <# If $DC is equall to 4(backup DC) or 5(primary DC) as there is no local user on Domain Controller #>   
try{
    $DC = (Get-WmiObject Win32_ComputerSystem).domainrole
   }catch{ Write-error " Unable to Check if machine is Domain Controller $_.Exception.message" }
        
if($IsDomain -eq $true) #Section for Domain Account
    {
        try{
            $sysinfo = Get-WmiObject -Class Win32_ComputerSystem 
            $NameStatus =  $AccountUserID.Contains("@")
            if($NameStatus -eq $True)
            {
                $ProvidedName,$ProvidedDomainName = $AccountUserID.split('@')
                If($ProvidedDomainName -eq $sysinfo.Domain)
                {
                    $AccountUserID = $ProvidedName 
                }
                else
                {
                    $OutputMessage += "[MSG: Domain name mis-match. Please confirm the Domain if User belongs to ]: " + $sysinfo.Domain
                }
             }
            }catch{Write-Error "$_.Exception.message"}
        
        if($DC -eq "4" -or $DC -eq "5")
        {
         try{
                Import-Module ActiveDirectory -ErrorAction Stop
                $OutputMessage += "[MSG: AD Powershell Module Imported] "                        
            }catch{ Write-error "[MSG: Exception encountered while Importing AD module.Please Connect to Other Domain Controller. ] $_.Exception.message" }

         try{   $ADuser=$null
                $ADuser =  Get-ADUser $AccountUserID -properties * -ErrorAction Stop
            }catch [Microsoft.ActiveDirectory.Management.ADIdentityNotFoundException] {$OutputMessage += "[MSG: User does not Exist] in the domain: "+$sysinfo.Domain } 
            catch { write-error "$_.Exception.message" }
                
            If($ADuser)
                {
                  try{  
                       if(($AccountEnable -eq $True) -and ($ADUser.Enabled -eq $false))
                       {
                           Enable-ADAccount -Identity $AccountUserID
                           $OutputMessage += "Account Enabled"
                       }
                       if(($AccountUnlock -eq $True) -and ($ADuser.LockedOut -eq $True))
                       {
                            Unlock-ADAccount -Identity $AccountUserID
                            $OutputMessage += "Account Unlocked"
                       }
                       $UpdatedADuser =  Get-ADUser $AccountUserID -properties Name,SamAccountName,Enabled,Lockedout | Select Name,SamAccountName,Enabled,Lockedout  -ErrorAction Stop 
                       $OutputMessage += $UpdatedADuser|fl                          
                     }catch{ Write-Error "[MSG: Exception encountered while Enabling/Unlocking the User Account ] $_.Exception.message"}   
                  }
            else
                {
                      try{
                      <#User does not match , get the list of user based on first 2 letters #> 
                      $OutputMessage +="Below are the list of similar user in domain"
                      $ADuser = $AccountUserID.Substring(0,1)+"*"
                      $ADUserListSimilarType = Get-ADUser -filter 'Name -like $ADuser' -ErrorAction stop |select Name,SamAccountName, UserPrincipalName, Enabled|fl
                      $OutputMessage += $ADUserListSimilarType                       
                      }catch{Write-Error "[MSG: Exception encountered while pulling the user list] $_.Exception.message"}
                 }
        }
        else
        {
            try{
               <#Collect List of DC in environment#>
               $Command = "nltest /dclist:"+$sysinfo.domain
               $DCList = Invoke-expression $Command 
               $OutputMessage += "Please connect to Domain Controller for executing the script from below list :" 
               $OutputMessage += $DCList |fl
               }catch{ Write-Error "Unable to pull the list of Domain Controllers.  $_.Exception.message"}
        }
    }
else # Section for Local Account
    {
        if($DC -eq "4" -or $DC -eq "5")
        {
            $OutputMessage += "No Local User on Domain Controller"
        }
        else
        {
             try{
                    $Localuser =  Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='True' and Name='$AccountUserID'" | Select PSComputername, LocalAccount, Status,Disabled, AccountType, Lockout, PasswordRequired, PasswordChangeable, PasswordExpires, SID
                }catch{ Write-error " Unable to check Local user information $_.Exception.message" }
                if ($Localuser)
                 {
                    try{
                            $MachineName = $env:COMPUTERNAME
                            $UsrCommand =[ADSI] "WinNT://$MachineName/$AccountUserID, user"
                            <# If account is Lockedout#>
                            if(($AccountUnlock -eq $True) -and ($Localuser.Lockout -eq $True))
                            {
                                $UsrCommand.IsAccountLocked = $False
                                $UsrCommand.Setinfo()
                                $OutputMessage += " Account is Unlocked" 
                            }
                            <# If account is Disabled#>                      
                            if(($AccountEnable -eq $True) -and ($Localuser.Disabled -eq $True))
                            {
                                $UserAccountProperties = Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='True' and Name='$AccountUserID'" #| %{$_.disabled=$false;$_.put()} | out-null
                                $UserAccountProperties.Disabled = $false
                                $UserAccountProperties.Put() | out-null
                                $OutputMessage += " Account is Enabled"
                            }
                            $UpdatedStatus=Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='True' and Name='$AccountUserID'" | Select Name,Disabled,Lockout
                            $OutputMessage += $UpdatedStatus |fl
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
    }
Write-Output $OutputMessage |fl

