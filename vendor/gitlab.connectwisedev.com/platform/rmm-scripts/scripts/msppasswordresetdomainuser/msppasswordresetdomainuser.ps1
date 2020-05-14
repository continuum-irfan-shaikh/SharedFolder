<#
.SYNOPSIS
Developed By : GRT [HD-Automation] , Continuum
File Vesion : 2.0
.DESCRIPTION
The script is designed to change the password of domain user, if Local User account is Locked/Disabled or both it will not do anything.
.PARAMETER 
DomainUserID is the variable required for changing the password.
.PARAMETER 
ChangePassword is the password which needs to be set as password for the user account.
 
.PARAMETER
ChangePwdNextLogon is a boolean. True means to set Change Password on next logon.
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
        Write-Output "[MSG : Incompatible System Either OS is below windows 7/2008R2 or powershell version less than 2.0]"
        Exit
    }
    <# else{Write-output "[MSG: Good to Go]"} #>
}catch{ Write-error "[MSG: $_.Exception.message]"}
<# Check if machine is in Domain or Workgroup Environment #>
try{
    $sysinfo = Get-WmiObject -Class Win32_ComputerSystem }
catch{Write-Error "Unable to check if machine is in Domain. $_.Exception.message"}

if($sysinfo.domain -eq "WORKGROUP")
{
    $OutputMessage += "Machine Connected is not in domain environment."
}
Else
{   <# If $DC is equall to 4(backup DC) or 5(primary DC) #>
    try{
        $DC = (Get-WmiObject Win32_ComputerSystem).domainrole
    }catch{ Write-error " Unable to Check if machine is Domain Controller $_.Exception.message" }

    if($DC -eq "4" -or $DC -eq "5" )
    {   <#Check if Powershell module is present#>
        try{
            Import-Module ActiveDirectory -ErrorAction Stop
            $OutputMessage += "Active Directory powershell module imported"
        }catch{ Write-Error "Exception encountered while Importing AD module. Please Connect to Other Domain Controller.  $_.Exception.message"}
        <#Check if User Exist#>
        try{
            $ADuser=$null
            $ADuser =  Get-ADUser $DomainUserID -properties * -ErrorAction stop
        }catch [Microsoft.ActiveDirectory.Management.ADIdentityNotFoundException] {$OutputMessage += "[MSG: User Doesnot Exist] FAILED "}
        catch { write-error "$_.Exception.message" }
        If($ADuser)
        {
            if(($ADUser.Enabled -eq $true) -And ($ADuser.LockedOut -eq $False))
            {
                Set-adaccountpassword $DomainUserID -reset -newpassword (ConvertTo-SecureString -AsPlainText $ChangePassword -Force)

                if($ChangePwdNextLogon -eq $true)
                {
                    if(($ADUser.PasswordNeverExpires -eq $False) -And ($ADUser.CannotChangePassword -eq $False ))
                    {
                        Set-ADUser $DomainUserID -ChangePasswordAtLogon $true
                    }
                    else
                    {
                        $OutputMessage += "Unable to change password at Next logon as $DomainUserID Account is Configuration for Password Exipration :$ADUser.PasswordNeverExpires And CannotChangePassword : $ADUser.CannotChangePassword"
                    }
                }
                else
                {
                    $OutputMessage += "Technician has not selected to change the password at next logon for user $DomainUserID"
                }
                $OutputMessage += "Password Changed"
            }
            else
            {
                $OutputMessage += "Account is Locked :"+$ADUser.LockedOut
                $OutputMessage += "Account is Enabled : "+$ADUser.Enabled
            }

        }
        else
        {
            try{
                <#User does not match , get the list of user based on first 2 letters #>
                $ADuser = $DomainUserID.Substring(0,1)+"*"
                $ADUserListSimilarType = Get-ADUser -Filter 'name -like $ADuser' | select Name, Enabled, UserprincipalName
                $OutputMessage += $ADUserListSimilarType
            }catch{Write-Error "Unable to pull the list of Users of similar Name.  $_.Exception.message"}
        }
    }
    else
    {
        try{
            <#Collect List of DC in environment#>
            $Command = "nltest /dclist:"+$sysinfo.domain
            $DCList = Invoke-expression $Command
            $OutputMessage += "Machine Connected is not DC. Please connect to the mentioned list ."
            $OutputMessage +=$DCList
        }catch{ Write-Error "Unable to pull the list of Domain Controllers.  $_.Exception.message"}
    }

}
Write-output $OutputMessage|Format-List
