if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}
# https://support.microsoft.com/en-sg/help/968264/error-message-when-you-try-to-map-to-a-network-drive-of-a-dfs-share-by
Function Get-RDPSessions {
    try{
        $ErrorActionPreference = 'Stop'
        if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64" -and $env:PROCESSOR_ARCHITECTURE -eq 'x86') { $Query = 'c:\windows\sysnative\query.exe' }else { $Query = 'c:\windows\System32\query.exe' }
        $LoggedOnUsers = if (($Users = (& $Query user 2>&1))) {
            $Users | ForEach-Object { (($_.trim() -replace ">" -replace "(?m)^([A-Za-z0-9]{3,})\s+(\d{1,2}\s+\w+)", '$1  none  $2' -replace "\s{2,}", "," -replace "none", $null)) } |
            ConvertFrom-Csv |
            Where-Object { $_.state -eq 'Active'} #-and $_.SESSIONNAME -like "rdp*" }
        }
        return $LoggedOnUsers
    }
    catch{
        # no action
    }
}

Function TestRegistry {
    $value = Get-ItemProperty 'HKLM:\System\CurrentControlSet\Control\Lsa\' -Name DisableDomainCreds -ErrorAction SilentlyContinue | Select-Object -ExpandProperty disabledomaincreds
    if ($value -eq 1) {
        return $true
    }
    else { 
        return $false
    }
}

Function AppendToCSV ($Type, $Message) {
    New-Object -TypeName psobject -Property @{Type = $Type; Message = $Message } | ConvertTo-Csv -NoTypeInformation | Select-Object -Skip 1 | Out-File C:\temp\log.csv -Append
}

# validate input parameters
if($DriveLetter -and !$DriveLetter.EndsWith(':')){$DriveLetter = $DriveLetter + ':'}
if($Path -and $Path.EndsWith('\')){$Path = $Path.TrimEnd('\')}

$Scriptblock = {
    $ErrorActionPreference = 'Stop'
    try {    
        Function MapNetworkDrive($DriveLetter, $Directory, $Persistent, $Username, $Password) {
            try {
                if (Test-Path $DriveLetter) {
                    # Write-Output "`n[Map Network Drive] Failed to Map the network drive: `'$DriveLetter`' because it already exists."
                    AppendToCSV 'Output' "[Map Network Drive] Failed to Map the network drive: `'$DriveLetter`' because it already exists."
                }
                else {
                    $Network = New-Object -ComObject WScript.Network
                    $Network.MapNetworkDrive($DriveLetter, "$Directory", $Persistent, $UserName, $Password)
                }
                #Stop-Process -ProcessName explorer -Confirm:$false -Force
            }
            catch {
                # Write-Output "`n[Map Network Drive] $($_.Exception.Message)"
                AppendToCSV 'Output' "[Map Network Drive] $($_.Exception.Message)"
            }
        }
    
        Function RemoveNetworkDrive($DriveLetter) {
            try {
                if (Test-Path $DriveLetter) {
                    $Network = New-Object -ComObject WScript.Network
                    $Network.RemoveNetworkDrive($DriveLetter, $True)
                }
                else {
                    AppendToCSV 'Output' "[Remove Network Drive] Failed to Remove the network drive: `'$DriveLetter`' because it does not exists."
                }
            }
            catch {
                # Write-Output "`n[Remove Network Drive] $($_.Exception.Message)"
                AppendToCSV 'Output' "[Remove Network Drive] $($_.Exception.Message)"
            }
        }
    
        Function ChangeLabel($DriveLetter, $NewLabel) {
            try {
                if (Test-Path $DriveLetter) {
                    $Application = New-Object -ComObject shell.application
                    $Application.NameSpace($DriveLetter).self.name = $NewLabel
                }
                else {
                    #Write-Output "`n[Change Drive Label] Failed to change the label of drive: `'$DriveLetter`' because it does not exists."
                    AppendToCSV 'Output' "[Change Drive Label] Failed to change the label of drive: `'$DriveLetter`' because it does not exists."
                }
            }
            catch {
                #Write-Output "`n[Change Drive Label] $($_.Exception.Message)"
                AppendToCSV 'Output' "[Change Drive Label] $($_.Exception.Message)"
            }
        }
    
        Function HideDriveFromExplorer($DriveLetter) { 
            # logoff session to see the changes
            try {
                $Registry = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer"
                $DriveNumber = [Int][Char]$($DriveLetter -replace ':', '').ToUpper() - 65
                $Mask += [System.Math]::Pow(2, $DriveNumber)
            
                if (Test-Path $Registry) {
                    Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer"-Name NoDrives -Value $Mask -Type DWORD -Force
                }
                else {
                    New-Item -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Policies\" -Name 'Explorer' -Force | Out-Null
                    Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer"-Name NoDrives -Value $Mask -Type DWORD -Force
                }

                $NoDrivesValue = Get-ItemProperty "HKCU:\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer" -Name NoDrives | Select-Object -ExpandProperty NoDrives
                if ($NoDrivesValue -eq $Mask) {
                    AppendToCSV 'Output' "Drive $DriveLetter is now hidden, but it is required to log-off and log-in again to see the changes."
                }

            }
            catch {
                #Write-Output "`n[Hide Drive From Explorer] $($_.Exception.Message)"
                AppendToCSV 'Output' "[Hide Drive From Explorer] $($_.Exception.Message)"
            }
        }
    
        Function SuccessfulIfScriptBlockReturnsTrue {
            param(
                [ScriptBlock] $ScriptBlock
            )
            
            if ( & $ScriptBlock ) { 
                AppendToCSV 'Output' "Overall Task is Successful"
            }
            else { 
                AppendToCSV 'Output' "Overall Task has failed"
            }
        }
    
        Function ListNetworkDrives {
            $Data = (New-Object -ComObject wscript.network).EnumNetworkDrives() | Where-Object { $_.typename -notlike "*Comobject*" }
            if ($Data.Length -ge 1) {
                $obj = For ($i = 0; $i -lt $($Data.length - 1); $i = $i + 2) {
                    New-Object PSOBject -Property @{ DriveLetter = $data[$i]; Path = $Data[$i + 1] }
                }
            }
    
            return $obj
        }

        Function AppendToCSV ($Type, $Message) {
            New-Object -TypeName psobject -Property @{Type = $Type; Message = $Message } | ConvertTo-Csv -NoTypeInformation | Select-Object -Skip 1 | Out-File C:\temp\log.csv -Append
        }
        
        if ($MapPersistent -eq 'True') { $MapPersistent = $true } else { $MapPersistent = $false }
    
        Switch ($Action) {
            'Map Network Drive' {
                if ($DisconnectBeforeMapping -eq 'True') { RemoveNetworkDrive $DriveLetter }
                MapNetworkDrive $DriveLetter $Path $MapPersistent $Username $Password
                Start-Sleep -Seconds 2
                if ($ExplorerLabel) { ChangeLabel $DriveLetter $ExplorerLabel }
                if ($HideFromWindowsExplorer -eq 'True') { HideDriveFromExplorer $DriveLetter }
                SuccessfulIfScriptBlockReturnsTrue { [bool](ListNetworkDrives | Where-Object { $_.DriveLetter -eq $DriveLetter -and $_.Path -eq $Path }) }
            }
            'Map with WebDAV' {
                if ($DisconnectBeforeMapping -eq 'True') { RemoveNetworkDrive $DriveLetter }
                MapNetworkDrive $DriveLetter $Path $MapPersistent $Username $Password
                if ($ExplorerLabel) { ChangeLabel $DriveLetter $ExplorerLabel }
                if ($HideFromWindowsExplorer -eq 'True') { HideDriveFromExplorer $DriveLetter }        
                SuccessfulIfScriptBlockReturnsTrue { [bool](ListNetworkDrives | Where-Object { $_.DriveLetter -eq $DriveLetter -and $_.Path -eq $Path }) }
            }
            'Disconnect' {
                RemoveNetworkDrive $DriveLetter
                SuccessfulIfScriptBlockReturnsTrue { -not (ListNetworkDrives | Where-Object { $_.DriveLetter -eq $DriveLetter }) }
            }
            'Disconnect All' {
                Foreach ($item in $(((New-Object -ComObject wscript.network).EnumNetworkDrives() -match "^([a-z]):"))) { RemoveNetworkDrive $item }
                SuccessfulIfScriptBlockReturnsTrue { -not((ListNetworkDrives | Select-Object -ExpandProperty DriveLetter) -match "^([a-z]):") } 
            }
        }
    }
    catch {
        $_
    }
}

$file = @"
`$Action = "$Action"
`$DriveLetter = "$DriveLetter"
`$DisconnectBeforeMapping = `'$DisconnectBeforeMapping`'
`$MapPersistent = `'$MapPersistent`'
`$HideFromWindowsExplorer = `'$HideFromWindowsExplorer`'
`$Path = "$Path"
`$UserName = "$UserName"
`$Password = "$Password"
`$ExplorerLabel = "$ExplorerLabel"
"@ + $Scriptblock.ToString()

$ErrorActionPreference = 'Stop'
$LogDir = 'C:\temp\'
$LogFilePath = 'C:\temp\Log.csv'
$MainFilePath = 'C:\temp\MapNetworkDrive.ps1'

if (![IO.Directory]::Exists($LogDir)) { [IO.Directory]::CreateDirectory($LogDir) | Out-Null } # create folder if that doesn't exists
$File | Out-File $MainFilePath -Encoding utf8 # copy main script which will be later executed from scheduled task

try {
    if (($Session = Get-RDPSessions)) {
        '"Message","Type"' | Out-File $LogFilePath -Force
        # if(!(TestRegistry)){
        #     AppendToCSV 'Output' "Following registry setting is mandatory. Please configure and reboot the system before running this script again.

        # Location: HKEY_LOCAL_MACHINE\System\CurrentControlSet\Control\Lsa\
        # Name: DisableDomainCreds
        # Value: 1 (DWORD)"
        # }
        # else{
        $User = $Session | Select-Object -First 1 -ExpandProperty Username # select one user profile from active RDP sessions
        Start-Sleep -Seconds 5
        $TaskName = "MapNet"
        $Task = "PowerShell.exe -executionpolicy bypass -NoExit -noprofile -WindowStyle Hidden -command '. $MainFilePath '"
        $StartTime = (Get-Date).AddMinutes(2).ToString('HH:mm') # time in 24hr format

        schtasks.exe /create /s $($env:COMPUTERNAME) /tn $TaskName /sc once /tr $Task /st $StartTime /ru $User /F | Out-Null
        schtasks.exe /End /TN $TaskName | Out-Null
        schtasks.exe /Run /TN $TaskName | Out-Null
        Start-Sleep -Seconds (120)
        # }

        # Sending logs to PowerShell Output\Error streams so that agent can capture it
        if (Test-Path $LogFilePath) {
            $Logs = Import-Csv $LogFilePath
            if($Logs | Select-Object -ExpandProperty Message -ErrorAction SilentlyContinue | Foreach-Object {$_.trim()}){
                foreach ($item in $Logs) {
                    switch ($item.type) {
                        'Output' { Write-Output "`n$($item.message)" }
                        'Error' { Write-Error "`n$($item.message)" }
                    }
                }
            }
            else{
                Write-Output "`nNo action performed, because the task exited due to timeout. This may be due to slowness of the System."
            }
        }

        schtasks.exe /Delete /TN "Mapnet" /F | Out-Null # scheduled task cleanup
        Remove-Item -Path $MainFilePath, $LogFilePath -ErrorAction SilentlyContinue # file cleanup
    }
    else {
        Write-Output "This script requires logon user and currently no user is logged in. `nNo action will be performed."; exit;
    }
}
catch {
    Write-Error $_
}
