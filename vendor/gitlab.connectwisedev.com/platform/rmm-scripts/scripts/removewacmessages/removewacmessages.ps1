$RegPath = "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\WindowsBackup"
$KeyName = "DisableMonitoring"

IF(-Not(Test-Path -Path $RegPath)){
	New-Item -Path $RegPath -Force | Out-Null
}

$Details = Get-ItemProperty $RegPath | gm -ErrorAction SilentlyContinue | ?{$_.Name -eq $KeyName}

IF ($Details -eq $null){
New-ItemProperty -Path $Regpath -Name $KeyName -PropertyType Dword -Value 1 | Out-Null
Write-Output "Item Successfully create and Property also Set"
}

IF($Details.Name -contains 'DisableMonitoring')    {
        $Value = Get-ItemProperty $RegPath -Name 'DisableMonitoring' 
IF ($Value.DisableMonitoring -eq 0)
    {
            Set-ItemProperty $RegPath -Name 'DisableMonitoring' -Value 1
            Write-Output "Item Value is Updated"
    }
Else
        {
            Write-Output "Item Property is alreaedy set"
        }
}

######################### Script ##############################

<#
Script name : turn off the action center messages for Backup,   Messages : Back up your files.

Steps to Disable the Backup
To do so, open regedit and navigate to the following  location:
HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\WindowsBackup


Steps to Perform in Script

1. User should have access to create registry entry on above path.


Progress:

1. Script Created for creating new entry and add existing entry in registry
#>
