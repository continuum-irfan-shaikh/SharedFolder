$RegPath = "HKLM:\Software\Policies\Adobe\Acrobat Reader\DC\FeatureLockdown"
$KeyName = "bUpdater"
If(-Not(Test-Path -Path $RegPath)){
    Return "Adobe Reader is currently not installed"
} Else {
    New-Item -Path $RegPath -Force | New-ItemProperty -Name $KeyName -PropertyType DWord -Value 0 | Out-Null
}
Set-ItemProperty -Path $RegPath -Name $KeyName -Value 0
Return Get-ItemProperty $RegPath
