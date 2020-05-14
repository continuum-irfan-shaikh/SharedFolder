$DomainGroups = get-wmiobject win32_group -Filter "Name='Domain Admins'"

if (!$DomainGroups) {
  Write-Output "Domain users are not found"
  return
}

$listUsers = @();
foreach ($group in $DomainGroups) {
  $query = "GroupComponent = `"Win32_Group.Domain='$($group.domain)'`,Name='$($group.name)'`""
  $list = Get-WmiObject win32_groupuser -Filter $query
  $listUsers += $list | %{$_.PartComponent} | % {$_.substring($_.lastindexof("Domain=") + 7).replace("`",Name=`"","\")}
}

if (!$listUsers) {
  Write-Output "Domain users are not found"
  return
}

$users = Get-WmiObject -Class Win32_UserAccount -Filter  "LocalAccount!='True'"
foreach ($user in $users) {
  $admin = $false
  foreach ($useradmin in $listUsers) {

  $useradmin = $useradmin.Substring(0,$useradmin.Length-1)
  $useradmin = $useradmin.Substring(1)
    if ($useradmin -eq $user.Caption){
      $admin = $true
    }
  }
  Write-Output "$(($user).Caption), Admin : $admin"
}
