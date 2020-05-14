<#
    .SYNOPSIS
         Disable User Accounts
    .DESCRIPTION
         It will disable the local user accounts.
         On member server machine and AD controller, it will disable the domain user account.
    .Help
         work station = 1
         #Domain Controller = 2
         #server = 3
         Use net user command for workstation. 
         Use ActiveDirectory cmdlets for DC and Member Server
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>

#Define $users in JSON Schema

$users1 = @()
$users1 = $users.split(',')

$ProductType = (Get-WmiObject -Class Win32_OperatingSystem).producttype 

function DisableAdUser {
   
    foreach ($user in $users1) {
        if ([bool] (Get-ADUser -Filter { SamAccountName -eq $user } -ErrorAction SilentlyContinue)) {
               
            Disable-ADAccount $user -ErrorAction SilentlyContinue
            $enabled = (Get-ADUser $user).Enabled
            if ($enabled -eq $false) {
                Write-Output "User account '$user' disabled successfully."
            }
            else {
                Write-Output "User account '$user' not disabled"
            }
   
        }
        else {
            Write-Output "User account '$user' doesn't exist"
        }
    
    }
    return
}

function DisableLocalUser {
    
    foreach ($user in $users1) {
            
        net user $user > null 2>&1
        
        if ($?) {
            net user $user /active:no > null 2>&1
   
            if ($?) {
                Write-Output "User account '$user' disabled successfully."
            }
            else {
                Write-Output "User account '$user' not disabled"
            }
        }
        else {
            Write-Output "User account '$user' doesn't exist"
        }
    }
    return
}

try {
    if ($ProductType -eq 1) {
        
        DisableLocalUser
    }
    else {
        if (Get-Module activedirectory) {
           
            DisableAdUser
        }
        else {
            if (Import-Module activedirectory) {
               
                DisableAdUser
            }
            else {
                Write-Output "`nUser accounts not disabled."
                Write-Error "`n"$_.Exception.Message
                EXIT
            }
          
        }
    }
}
catch {
    write-error "`n"$_.Exception.Message
}
