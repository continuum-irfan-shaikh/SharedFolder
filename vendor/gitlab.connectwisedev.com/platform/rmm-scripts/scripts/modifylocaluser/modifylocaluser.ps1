$systemRole = (Get-WmiObject -Class Win32_ComputerSystem).DomainRole
if ($systemRole -eq 4 -or $systemRole -eq 5){ 
    Write-Output "This TASK is not supported on a Domain Controller."
    Exit
}
$modifications = @{}
$adsi = [ADSI]"WinNT://$env:computername"
$existing = $adsi.Children | where {$_.SchemaClassName -eq 'user' -and $_.Name -eq $accountName }
if ($existing){
    if ($newAccountName){
        $user = get-wmiObject win32_userAccount -filter "Name='$accountName'"
        $user.Rename($newAccountName) > $null
        $accountName = $newAccountName
        $modifications.add("Account name changed to", $newAccountName )
    }
    if ($newPassword){
        net user $accountName $newPassword > $null
        $modifications.add("Password changed", "Yes")
    }
    if ($active){
        net user $accountName /active:"yes" > $null
        $modifications.add("Account is", "Active")
    } else {
        net user $accountName /active:"no" > $null
        $modifications.add("Account is ", "Disabled")
    }
    if ($type){
        $user = (Get-WmiObject -Class Win32_UserAccount | where-object {$_.Name -eq "$accountName"})
        $matchedgroup = ((Get-WmiObject -Class Win32_Group | where-object {$_.Name -eq "Administrators"}).Caption)
        $user | foreach {
            $query = "Associators Of {Win32_UserAccount.Domain='" `
            + $_.Domain + "',Name='" + $_.Name `
            + "'} WHERE AssocClass=Win32_GroupUser"
        }
        $memberOf = Get-WmiObject -Query $query | select -ExpandProperty Caption
        if ($type -eq "Administrators"){
            if ($memberOf -notcontains $matchedgroup){
            try { 
               net localgroup $type $accountName /add > $null}catch{Write-Host "Error : $_.Exception.Message"}
               $modifications.add("Member of group", $type )
            }
        } elseif ($memberOf -contains $matchedgroup){
            net localgroup "Administrators" $accountName /delete > $null
            $modifications.add("Member of group", "Administrators")
        }
    }
    if ($expire){
        net user $accountName /expires:"$expire" > $null
        $modifications.add("Account expires", $expire)
    }
    if ($passwordchg){
        net user $accountName /passwordchg:"Yes" > $null

        $AccountObj = Get-WmiObject win32_USERACCOUNT | Where-Object {$_.NAME -eq $accountName }
        $AccountObj.PasswordExpires = $True
        $AccountObj.Put() | Out-Null

        $strComputer=$env:COMPUTERNAME
        $user= [ADSI]"WinNT://$strComputer/$accountName,user"
        $user.Put("PasswordExpired", 1)
        $user.SetInfo()
        $modifications.add("User must change password at next logon", "True")
   } 

} else {
    Write-Error "Account name $accountName does not exist!"
    Exit
}
$result = New-Object -TypeName psobject -Property $modifications
Write-Output "Following properties modified for the account $accountName"  
$result | fl

