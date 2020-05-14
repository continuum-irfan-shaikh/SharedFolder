# $Target = 'domain' # STRING
# $Users = @" # MULTILINE STRING
# test100
#  test2000sd
# test300sd
# test400
# "@

# $Target = 'local'
# $Users = @"
# test100
#  test200
# test300
# test400
# "@


if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

$ErrorActionPreference = 'Stop'
try {

    $Users = $Users -split "`n" | Where-Object { ![string]::IsNullOrWhiteSpace($_.trim())}
    $SystemRole = (Get-WmiObject -Class Win32_ComputerSystem).DomainRole
    
    switch ($Target) {
        'local' {
            if (4, 5 -contains $SystemRole) {
                Write-Output "Script can not run on a Domain Controller. Please choose a Local machine and retry."
                Break
            }
            
            ForEach ($User in $Users) {
                $Obj = Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount='True' AND Name='$User'"
                if ($Obj) {
                    if (!$obj.disabled) {
                        if ($Obj.PasswordExpires -eq $false) {
                            $Obj.PasswordExpires = $true
                            $Obj.Put() | Out-Null
                        }
                        $ErrorActionPreference = 'SilentlyContinue'
                        $result = net user $User /logonpasswordchg:yes 2>&1
                        $ErrorActionPreference = 'Stop'
                        if ($result -like "*Success*") {
                            Write-Output "$User : Password expired succesfully."
                        }
                        else {
                            Write-Output "$User : Failed to expire the password."   
                        }
                    }
                    else {
                        Write-Output "$User : User is disabled." 
                    }
                }
                else {
                    Write-Output "$User : The user name could not be found."                                     
                }
            }
        }
        
        'domain' {
            if (0..3 -contains $SystemRole) {
                Write-Output "Script can not run on a Local Machine. Please choose a Domain controller and retry."
                Break
            }

            Import-Module ActiveDirectory

            Foreach ($User in $Users) {
                $Obj = Get-ADUser -Filter { name -eq $User } -Properties PasswordExpired -erroraction silentlycontinue
                if ($obj) {
                    if ($obj.Enabled) {
                        Set-ADUser $Obj -ChangePasswordAtLogon $true -PasswordNeverExpires $false -ErrorAction Silentlycontinue
                        if ((Get-ADUser $Obj -Properties PasswordExpired).PasswordExpired) {
                            Write-Output "$User : Password expired succesfully."                 
                        }
                        else {
                            Write-Output "$User : Failed to expire the password."                                     
                        }
                    }
                    else {
                        Write-Output "$User : User is disabled."                                     
                    }
                }
                else {
                    Write-Output "$User : The user name could not be found."                                     
                }
            }
        }

    }
}
catch {
    Write-Error $($_.exception.message)
}
