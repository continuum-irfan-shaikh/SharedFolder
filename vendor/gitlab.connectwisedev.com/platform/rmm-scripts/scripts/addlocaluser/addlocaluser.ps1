$adsi = [ADSI]"WinNT://$env:computername"
$existing = $adsi.Children | where {$_.SchemaClassName -eq 'user' -and $_.Name -eq $accountName }
if (-Not($existing)){
    net user $accountName $password /add /y > $null
} else {
    Write-Error "$accountName already exists"
    exit
}
if ($expire){
    net user $accountName /expires:"$expire" > $null
}
if ($active){
    net user $accountName /active:"yes" > $null
} else {
    net user $accountName /active:"no" > $null
}
if ($passwordchg){
    net user $accountName /passwordchg:"yes" > $null
} else {
    net user $accountName /passwordchg:"no" > $null
}
net user $accountName
