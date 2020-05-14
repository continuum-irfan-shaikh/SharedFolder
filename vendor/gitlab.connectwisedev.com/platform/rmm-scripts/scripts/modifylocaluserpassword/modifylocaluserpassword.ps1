if (-not $password){
  $password = '""'
}

[int]$Domain_Role=(Get-WmiObject -Class Win32_computersystem).Domainrole
if ($Domain_Role -ge 4) {
    Write-Error "Local account passowrd reset on domain controllers isn't supported!"
    Exit
}

$account_status=Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='true' and  name='$accountName'" |
                select Name,Lockout,Disabled
if ($account_status.Lockout){
    Write-Error "The user account is locked out! unlock the account first."
    Exit
}
if ($account_status.Disabled){
    Write-Error "The user account is Disabled! enable the account first."
    Exit
}

try {
        net user $accountName $password | Out-Null
        if ($LASTEXITCODE -eq 0){
                if ( $password -eq '""' ){
                    Write-Warning "Password changed to blank passowrd for account $accountName, it's not safe..! please reset to strong password." 
                }
                else {Write-Output "Password changed successfully for account $accountName"}
        }
}Catch{ Write-Error "$_.Exception.Message" }
