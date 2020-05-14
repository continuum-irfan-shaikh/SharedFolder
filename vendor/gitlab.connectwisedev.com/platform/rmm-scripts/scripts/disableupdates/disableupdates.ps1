$RegPath = "HKLM:\SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate\AU"
If(-Not(Test-Path -Path $RegPath)){
    New-Item -Path $RegPath -Force | New-ItemProperty -Name "NoAutoUpdate" -PropertyType DWord -Value 1 | Out-Null
}
Set-ItemProperty -Path $RegPath -Name "NoAutoUpdate" -Value 1
Return Get-ItemProperty $RegPath