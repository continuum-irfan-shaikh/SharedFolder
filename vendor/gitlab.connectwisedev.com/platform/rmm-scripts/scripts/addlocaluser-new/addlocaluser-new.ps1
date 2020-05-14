<#
    Actions
        AddNewUser
        ChangePassword
        RemoveUser
        ModifyUser
#>


<#
Date :  13/01/2020
Version : 1.0

$action = "AddNewUser"
$username = "nirav2"
$fullname = "Nirav Sachora updated"
$description = "this is my updated account"
$password = "pass@123"
$usermustchangepasswordatnextlogon = 'Yes'
$usercannotchangepassword = 'Yes'
$passwordneverexpires = 'Yes'
$addtheusertothesegroups = "group1,group2,group3"
$removetheuserfromthesegroups = "group1,group2,group3"
$AccountIsDisabled = 'No'
$profilepath = 'C:\USers\nirav'
$logonscript = 'C:\profiles'
$homefolder = Conect/Local
$driveletter
$path
$HideAccountOnTheWelcomeScreen = $true
$OverwriteIfTheUserAlreadyExists = $true
#>
#########################Preparing inputs######################

###############OS INFO#########################################

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

$OSInfo = (Get-WmiObject -Class Win32_Operatingsystem).ProductType

if ($OSInfo -eq $null) {

    "[ERROR]:Failed to check System type, script will now exit"
    Exit

}
elseif ($OSInfo -eq 2) {

    "[MSG]:Script cannot be run at Domain Controller"
    Exit

}


if ($addtheusertothesegroups) {
    $addtheusertothesegroups = $addtheusertothesegroups -split ","
    $grlist = @()
    $erroractionpreference = "SilentlyContinue"
    foreach ($group in $addtheusertothesegroups) {
        $exist = Net Localgroup $group 2>&1
        if ($exist) { $grlist += $group }else { "$group does not exist"; Continue; }
    }
    $addtheusertothesegroups = $grlist
}
if ($removetheuserfromthesegroups) {
    $removetheuserfromthesegroups = $removetheuserfromthesegroups -split ","
    $grlist = @()
    foreach ($group in $removetheuserfromthesegroups) {
        $exist = Net Localgroup $group 2>&1
        if ($exist) { $grlist += $group }else { "$group does not exist"; Continue; }
    }
    $removetheuserfromthesegroups = $grlist
}
$erroractionpreference = "Continue"

################################################################

function change_password {

    $ErrorActionpreference = "SilentlyContinue"
    $check_exist = net user $username 2>&1
    if ($check_exist) {

        net user $username $password 2>&1 | out-null
    }
    else {
        "[ERROR]:User is not present in the system."
        Exit;
    }

    if ($?) {
        "-"*30+"`n[MSG]:Password updated successfully.`n"+"-"*30

    }
    else {
        "-"*30+"`n[ERROR]:Failed to update password.`n"+"-"*30
    }
}

function common_function {
    if ($description) {
        net user $username /comment:$description 2>&1 | out-null
        if ($?) { $tracker.Add("Description", "SUCCESS") }else { $tracker.Add("Description", "FAILED") }
    }
    if($action -eq "Modifyuser" -and ($password)){
        change_password    
    }
    if ($usermustchangepasswordatnextlogon -ne 'Yes') {
        if ($usercannotchangepassword -eq 'Yes') {
            net user $username /passwordchg:no 2>&1 | out-null
            if ($?) { $tracker.Add("Password change permission", "SUCCESS") }else { $tracker.Add("Password change permission", "FAILED") }
        }
        elseif (($usercannotchangepassword -eq 'No') -and (($action -eq "ModifyUser") -or ($OverwriteIfTheUserAlreadyExists -eq $true))) {
            net user $username /passwordchg:yes 2>&1 | out-null
            if ($?) { $tracker.Add("Password change permission", "SUCCESS") }else { $tracker.Add("Password change permission", "FAILED") }
        }
        if ($passwordneverexpires -eq 'Yes') {
            $userdetails = net user $username 2>&1
            if ($userdetails[9] -like '*never*') { $tracker.Add("Password never expire", "SUCCESS") }
            else {
                $user = [adsi]"WinNT://$env:computername/$username"
                $user.UserFlags.value = $user.UserFlags.value -bor 0x10000
                $user.CommitChanges()
                if ($?) { $tracker.Add("Password never expire", "SUCCESS") }else { $tracker.Add("Password never expire", "FAILED") }
            }
        }
        elseif (($passwordneverexpires -eq 'No') -and (($action -eq "ModifyUser") -or ($OverwriteIfTheUserAlreadyExists -eq $true))) {
            $userdetails = net user $username 2>&1
            if ($userdetails[9] -like '*never*') {
                $user = [adsi]"WinNT://$env:computername/$username"
                $ADS_UF_DONT_EXPIRE_PASSWD = 0x00010000
                $adsi = [ADSI]"WinNT://$env:COMPUTERNAME"
                $user.UserFlags = $user.UserFlags.Value -bxor $ADS_UF_DONT_EXPIRE_PASSWD
                $user.SetInfo()
                if ($?) { $tracker.Add("Password never expire", "SUCCESS") }else { $tracker.Add("Password never expire", "FAILED") }
            }
            else {
                $tracker.Add("Password never expire", "SUCCESS")
            }
        }
    }
    if ($usermustchangepasswordatnextlogon -eq 'Yes') {
        net user $username /passwordchg:"Yes" 2>&1 | out-null
        $acc = Get-WmiObject win32_USERACCOUNT | ? { $_.NAME -eq $username }
        $acc.PasswordExpires = $True
        $acc.Put() | out-null
        
        $str = $env:COMPUTERNAME
        $user = [ADSI]"WinNT://$str/$username,user"
        $user.Put("PasswordExpired", 1)
        $user.SetInfo()
        if ($?) { $tracker.Add("User must change password at next logon", "SUCCESS") }else { $tracker.Add("User must change password at next logon", "FAILED") }
    }
    elseif (($usermustchangepasswordatnextlogon -eq 'No') -and (($action -eq "ModifyUser") -or ($OverwriteIfTheUserAlreadyExists -eq $true))) {
        $acc = Get-WmiObject win32_USERACCOUNT | ? { $_.NAME -eq $username }
        $str = $env:COMPUTERNAME
        $user = [ADSI]"WinNT://$str/$username,user"
        $user.Put("PasswordExpired", 0)
        $user.SetInfo()
        if ($?) { $tracker.Add("User must change password at next logon", "SUCCESS") }else { $tracker.Add("User must change password at next logon", "FAILED") }
    }
    if ($action -eq "ModifyUser") {
        if ($AccountIsDisabled -eq 'No') {
            net user $username /active:yes 2>&1 | out-null
            if ($?) { $tracker.Add("Enabled", "SUCCESS") }else { $tracker.Add("Enabled", "FAILED") }
        }
    }
    if ($AccountIsDisabled -eq 'Yes') {
        net user $username /active:no 2>&1 | out-null
        if ($?) { $tracker.Add("Disabled", "SUCCESS") }else { $tracker.Add("Disabled", "FAILED") }
    }
    if ($profilepath) {
        net user $username /profilepath:$profilepath 2>&1 | out-null
        if ($?) { $tracker.Add("Profile Path", "SUCCESS") }else { $tracker.Add("Profile Path", "FAILED") }
    }
    if ($logonscript) {
        net user $username /scriptpath:$logonscript 2>&1 | out-null
        if ($?) { $tracker.Add("Logon Script", "SUCCESS") }else { $tracker.Add("Logon Script", "FAILED-Invalid logon path") }
    }
    if ($path) {
        if (!$driveletter -and $homefolder -eq "Local") {
            net user $username /homedir:$path 2>&1 | out-null
            if ($?) { $tracker.Add("Logon Path", "SUCCESS") }else { $tracker.Add("Logon Path", "FAILED") }
        }
        else {
            if (!$driveletter) {
                "Please provide drive letter"
                Exit;
            }
            net user $username /homedir $driveletter :$path 2>&1 | out-null
            if ($?) { $tracker.Add("Logon Path", "SUCCESS") }else { $tracker.Add("Logon Path", "FAILED") }
        }
    }
    if ($addtheusertothesegroups) {
        $useraddition = @{ }
        foreach ($i in $addtheusertothesegroups) {
            net localgroup $i $username /add  2>&1 | out-null
            if ($?) { $useraddition.Add("$i", "Added") }else { $useraddition.Add("$i", "Failed to add") }
        }
        $useraddition.GetEnumerator() | Foreach-object {
            Write-output "Adding user to $($_.key)"
            Write-output "$($_.value)"
            "-----------------------------------"
        }
    }
    if ($removetheuserfromthesegroups) {
        $userdeletion = @{ }
        foreach ($i in $removetheuserfromthesegroups) {
            net localgroup $i $username /delete 2>&1 | out-null
            if ($?) { $userdeletion.Add("$i", "Removed") }else { $userdeletion.Add("$i", "Failed to remove") }
        }
        $userdeletion.GetEnumerator() | Foreach-object {
            Write-output "Removing user from $($_.key)"
            Write-output "$($_.value)"
            "-----------------------------------"
        }
    }
    if ($HideAccountOnTheWelcomeScreen) {
        if (!(Test-Path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\SpecialAccounts")) {
            New-Item -Path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -Name "SpecialAccounts" | out-null
            New-Item -Path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\SpecialAccounts" -Name "UserList" | out-null
            New-ItemProperty -Path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\SpecialAccounts\UserList" -Name $username -PropertyType "Dword" | out-null
        }
        elseif (Test-Path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\SpecialAccounts") {
            if (!(Test-path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\SpecialAccounts\UserList")) {
                New-Item -Path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\SpecialAccounts" -Name "UserList" | out-null
                New-ItemProperty -Path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\SpecialAccounts\UserList" -Name $username -PropertyType "Dword" | out-null
            }
            elseif ((((Get-itemproperty -path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\SpecialAccounts\UserList").$username) -eq 0) -eq $false) {
                New-ItemProperty -Path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\SpecialAccounts\UserList" -Name $username -PropertyType "Dword" | out-null
            }
        }
        if ((Get-itemproperty -path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\SpecialAccounts\UserList").$username -eq 0) {
            $tracker.Add("Hidden From Welcome Screen", "SUCCESS`n`n[Note]:This user is hidden from welcome screen")
        }
        else { $tracker.Add("Hidden from Welcome screen", "FAILED") }
    }
    switch ($action) {
        "AddNewUser" { "-" * 40 + "`nProperty update status for $username`n`nNote: If no properties are provided, output will be blank`n" + "-" * 40 }
        "ModifyUser" { "-" * 40 + "`nProperty update status for $username`n`nNote: If no properties are provided, output will be blank`n" + "-" * 40 }
    }
    $tracker.GetEnumerator() | Foreach-object {
        Write-output "Updating..$($_.key)"
        Write-output "$($_.value)"
        "-----------------------------------"
    }
}

function modify_user {
    $tracker = @{ }
    $ErrorActionpreference = "SilentlyContinue"
    $check_exist = net user $username 2>&1
    if (!$check_exist) {
        "[ERROR]:User is not present in the system"
        Exit;
    }
    if ($fullname) {
        net user $username /fullname:$fullname 2>&1 | out-null
        if ($?) { $tracker.Add("Display Name", "SUCCESS") }else { $tracker.Add("Display Name", "FAILED") }
    }
    common_function
}

function adduser {
    $tracker = @{ }
    $ErrorActionpreference = "SilentlyContinue"
    $check_exist = net user $username 2>&1
    if (!$check_exist) {
        net user $username $password /add /fullname:$fullname  2>&1 | out-null
        if (!$?) {
            "[ERROR]:User creation failed"
            Exit;
        }
        else {
            "-" * 35 + "`n[MSG]:User $username created successfully`n`nChecking for properties to update..`n" + "-" * 35
        }
    }
    elseif ($check_exist -and ($OverwriteIfTheUserAlreadyExists -eq $true)) {
        modify_user
        return
    }
    else {
        "[ERROR]:User already present in the system."
        Exit;
    }
    common_function
}

function remove_user {
    $ErrorActionpreference = "SilentlyContinue"
    $check_exist = net user $username 2>&1
    if ($check_exist) {

        net user $username /delete 2>&1 | out-null
    }
    else {
        "[ERROR]:User is not present in the system"
        Exit;
    }
    if ($?) {
        "[MSG]:User deleted successfully."

    }
    else {
        "[ERROR]:Failed to delete user."
    }
}



switch ($action) {
    "AddNewUser" { adduser }
    "ChangePassword" { change_password }
    "RemoveUser" { remove_user }
    "ModifyUser" { modify_user }
}


