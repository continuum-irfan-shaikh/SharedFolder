<#
.SYNOPSIS
Developed By : GRT [HD-Automation] , Continuum
File Vesion : 2.0
.DESCRIPTION
The script is designed to change the password of domain user whose account is locked/disabled or both also change it to Enable and Unlock if domain User account is Locked/Disabled or both.
.PARAMETER
DomainUserID is to provide the Domain User ID .
.PARAMETER
ChangePassword is to change the password .
.PARAMETER
ChangePwdNextLogon is boolean for opting if user should change its password at next logon. If Password never expires or User can't change password is opted it will not change this option.
#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

$OutputMessage = @() <# Variable to capture Output for debugging purpose.#>
<# Compatibility Check if found incompatiblt with exit#>
try{
    $OSVersion=([system.environment]::OSVersion.Version).Major
    $PSVersion =(Get-Host).version
    if(($OSVersion -lt 6) -or ($PSVersion.Major -lt 2))
    {
        Write-Output += "[MSG : Incompatible System Either OS is below windows 7/2008R2 or powershell version less than 2.0]"
        Exit
    }
}catch{ Write-Output "[MSG: $_.Exception.message]"}
<# Compatibility Check Ends #>
try{
    $sysinfo = Get-WmiObject -Class Win32_ComputerSystem
    $NameStatus =  $DomainUserID.Contains("@")
    if($NameStatus -eq $True)
    {
        $ProvidedName,$ProvidedDomainName = $DomainUserID.split('@')
        If($ProvidedDomainName -eq $sysinfo.Domain)
        {
            $DomainUserID = $ProvidedName
        }
        else
        {
            $OutputMessage += "[MSG: Domain name mis-match ]"
        }
    }
}catch{Write-Error "$_.Exception.message"}
<# Check if machine is in Domain or Workgroup Environment #>
if($sysinfo.domain -eq "WORKGROUP")
{
    $OutputMessage += "Machine Connected is not in domain environment."
}
else
{
    <# If $DC is equall to 4(backup DC) or 5(primary DC) #>
    try{
        $DC = (Get-WmiObject Win32_ComputerSystem).domainrole
    }catch{ Write-error " Unable to Check if machine is Domain Controller $_.Exception.message" }
    if($DC -eq "4" -or $DC -eq "5" )
    {
        try{
            Import-Module ActiveDirectory -ErrorAction Stop
            $OutputMessage += "[MSG: AD Powershell Module Imported] "
        }catch{ Write-error "[MSG: Exception encountered while Importing AD module.Please Connect to Other Domain Controller. ] $_.Exception.message" }

        try{   $ADuser=$null
        $ADuser =  Get-ADUser $DomainUserID -properties * -ErrorAction Stop
        }catch [Microsoft.ActiveDirectory.Management.ADIdentityNotFoundException] {$OutputMessage += "[MSG: User Doesnot Exist] FAILED "}
        catch { write-error "$_.Exception.message" }

        If($ADuser)
        {
            try{
                if($ADUser.Enabled -eq $False)
                {
                    Enable-ADAccount -Identity $DomainUserID
                    $OutputMessage += "Account Enabled"
                }
                if($ADuser.LockedOut -eq $True)
                {
                    Unlock-ADAccount -Identity $DomainUserID
                    $OutputMessage += "Account Unlocked"
                }
                Set-adaccountpassword $DomainUserID -reset -newpassword (ConvertTo-SecureString -AsPlainText $ChangePassword -Force)
                if($ChangePwdNextLogon -eq $true)
                {
                    if(($ADUser.PasswordNeverExpires -eq $False) -And ($ADUser.CannotChangePassword -eq $False ))
                    {
                        Set-ADUser $DomainUserID -ChangePasswordAtLogon $true
                    }
                    else
                    {
                        $OutputMessage += "MSG: ChangePasswordAtLogon is not changed for $DomainUserID as user is Configuration for :"
                        $OutputMessage += "Password Never Expire : "+$ADUser.PasswordNeverExpires
                        $OutputMessage += "CannotChangePassword : "+$ADUser.CannotChangePassword
                    }
                }
                else
                {
                    $OutputMessage += "[MSG: Change Password at next logon is not selected  for user $DomainUserID] "
                }
                $OutputMessage += "[MSG: Password changed Successful]"
            }catch{ Write-Error "[MSG: Exception encountered while retreiving User Details ] $_.Exception.message"}
        }
        else
        {
            try{
                <#User does not match , get the list of user based on first 2 letters #>
                $OutputMessage +="Below are the list of similar user in domain"
                $ADuser = $DomainUserID.Substring(0,1)+"*"
                $ADUserListSimilarType = Get-ADUser -filter 'Name -like $ADuser' -ErrorAction stop |select Name, UserPrincipalName, Enabled|ft
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
            $OutputMessage += $DCList
        }catch{ Write-Error "Unable to pull the list of Domain Controllers.  $_.Exception.message"}
    }
}
Write-output $OutputMessage
