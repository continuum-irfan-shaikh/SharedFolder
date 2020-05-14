$DomainRole = (Get-WmiObject Win32_ComputerSystem).DomainRole
$adsi = [ADSI]"WinNT://$env:COMPUTERNAME"
$users = $adsi.Children | where {$_.SchemaClassName -eq 'User'} | % {$_.Name}
$Password = ""
$nopassUsers = @()

If($DomainRole -eq 0 -or $DomainRole -eq 1 -or $DomainRole -eq 2 -or $DomainRole -eq 3){

    Add-Type -AssemblyName System.DirectoryServices.AccountManagement
    $userObj = New-Object System.DirectoryServices.AccountManagement.PrincipalContext('machine', $env:computername)

    foreach ($user in $users) {

        try {
            if ($userObj.ValidateCredentials($User, $Password)){ $nopassUsers += $user }
        }catch{
            $LimitBlankPasswordUse = (Get-ItemProperty "HKLM:\SYSTEM\CurrentControlSet\Control\Lsa" -ErrorAction SilentlyContinue).LimitBlankPasswordUse
            if( $_.Exception -like "*blank passwords*" -And $LimitBlankPasswordUse -eq 1 ){ $nopassUsers += $user }
        }
    }
}

ElseIf ( $DomainRole -eq 4 -or $DomainRole -eq 5) {

    $CurrentDomain = "LDAP://" + ([ADSI]"").distinguishedName
    foreach ($user in $users){
        $domain = New-Object System.DirectoryServices.DirectoryEntry($CurrentDomain,$user,$Password)
        If ($domain.name) { $nopassUsers += $user }
    }
} Else {
    Write-Error "Unable to detect accounts with blank password as machine role is Unknown!"
    Exit
}

If($nopassUsers){
    Write-Output "Account with blank passsword`n============================"
    Write-Output $nopassUsers}
Else{ "Not found Account with no password!" }