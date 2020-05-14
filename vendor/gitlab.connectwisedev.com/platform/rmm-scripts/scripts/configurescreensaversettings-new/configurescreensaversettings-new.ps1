<#
    .SYNOPSIS
         Configure Screen Saver
    .DESCRIPTION
         Configure Screen Saver. It will enable/disable the screen saver. It will set timeout, password protection enable/disable. 
         It can set screen saver by using local path or network path.
    .Help
         Use HKEY 
    .Author
        Durgeshkumar Patel
    .Version
        1.3
#>

<########## Variables to define in JSON #################
$Action = "Enable"     ## Enable  OR Disable
$TimeOut = 1    #Timeout in minutes  Default is 1 minute
$OnResumePasswordProtected = "Enable"  ##Enable OR Disable OR No Change

$ChangeScreenSaver = $true   ## True or False

$StoreTheImageFileOn = "Local Machine"  # Local Machine OR Network Machine
$SourcePath = "C:\Durgesh\scr"       #Local screensaver path where screensaver file available
$FileName = "Bubbles111.scr"

#Network drive screensaver details
$SourcePath = '\\SCRIPT-WIN10-64\Durgesh\scr' #Network path start with \\
$FileName = "Bubbles111.scr"  # Filename with extension .scr
$UserName = 'labadmin'   #'administrator'  #Username to connect network path
$Password = 'lic@123'   #'Grt@2018'        #Password to connect network path
$LocalScreenSaverPath = "C:\Windows\System32"
###################>

$ErrorActionPreference = 'Stop'
$time_out = $TimeOut * 60 
    
function connect {
    try {
        $ErrorActionPreference = 'Stop'
        $global:DriveLetter = Get-ChildItem function:[g-z]: -n | Where-Object { !(Test-Path $_) } | random

        $global:Net = New-Object -ComObject WScript.Network

        if ([string]::IsNullOrEmpty($UserName) -and [string]::IsNullOrEmpty($Password)) {
            $Net.MapNetworkDrive($DriveLetter, "$SourcePath", $false)
       
        }
        else {
            $Net.MapNetworkDrive($DriveLetter, "$SourcePath", $false, $UserName, $Password)
        
        }
        if (Get-WmiObject -Class Win32_MappedLogicalDisk | Where-Object { $_.name -eq "$DriveLetter" }) {
            return $true
        }
        else {
            return $false
        }
    }
    catch {
        Write-Error $_.Exception.Message
    }
} #function close

function removedrive {

    $Net.RemoveNetworkDrive($DriveLetter)
    Start-Sleep 2
}

function validation {
    switch ($StoreTheImageFileOn) {
        "Local Machine" {
            $global:localpath = Join-Path "$SourcePath" "$FileName"
            if ((Test-Path "$localpath") -eq $false) {
                Write-Error "`nLocalpath '$localpath' of the screen saver is incorrect or screen saver file does not exist."
                Exit;
            }
        }
        "Network Machine" {
            if (connect) {
                Start-Sleep 2
                $mypath = Join-path "$DriveLetter" "$FileName"
                if (Test-Path "$mypath") {
                    if (Test-Path "$LocalScreenSaverPath") {
                
                        Copy-Item -path "$DriveLetter\$FileName" "$LocalScreenSaverPath" -Force
                
                        Start-Sleep 2 
                 
                        $global:localpath = join-path "$LocalScreenSaverPath" "$FileName"
                    
                        removedrive
                    }
                    else {
                        removedrive
                        Write-Error "`nLocal screen saver path '$LocalScreenSaverPath' is incorrect"
                        Exit;
                    }
                } 
                else {
                    removedrive
                    Write-Error "`nScreenSaver file '$FileName' does not exist at '$SourcePath'"
                    Exit;
                }
            }
            else {
                Write-Error "`nNetwork Drive '$SourcePath' connection issue. Check network connectivity."
                Exit;
            }
        }
    }
}  

function SetScreenSaver ($localpath) {
    
    If (Test-Path ($RegKey + "\Control Panel")) {
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel"
        If (Test-Path ($RegKey + "\Desktop")) {
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel\Desktop"
            ##Enabled
            New-ItemProperty -path $RegKey -name SCRNSAVE.EXE -value $localpath -PropertyType String -Force  
        }
        else {
            New-Item -path $RegKey -name "Desktop" -Force
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel\Desktop"
            ##Enabled
            New-ItemProperty -path $RegKey -name SCRNSAVE.EXE -value $localpath -PropertyType String -Force   
        }
    }
    else {
        New-Item -path $RegKey -name "Control Panel" -Force
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel"
        New-Item -path $RegKey -name "Desktop" -Force
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel\Desktop"
        ##Enabled
        New-ItemProperty -path $RegKey -name SCRNSAVE.EXE -value $localpath -PropertyType String  -Force
    }
}

function SetKey {

    If (Test-Path ($RegKey + "\Control Panel")) {
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel"
        If (Test-Path ($RegKey + "\Desktop")) {
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel\Desktop"
            ##Enabled
            New-ItemProperty -path $RegKey -name ScreenSaveActive -value $EnableValue -PropertyType String -Force | Out-Null
        }
        else {
            New-Item -path $RegKey -name "Desktop" -Force | Out-Null
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel\Desktop"
            ##Enabled
            New-ItemProperty -path $RegKey -name ScreenSaveActive -value $EnableValue -PropertyType String -Force | Out-Null
        }
    }
    else {
        New-Item -path $RegKey -name "Control Panel" -Force | Out-Null
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel"
        New-Item -path $RegKey -name "Desktop" -Force | Out-Null
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel\Desktop"
        ##Enabled
        New-ItemProperty -path $RegKey -name ScreenSaveActive -value $EnableValue -PropertyType String -Force | Out-Null
    }

    if ((Get-ItemProperty $RegKey -Name ScreenSaveActive | select -ExpandProperty ScreenSaveActive) -eq "$EnableValue") { return $true } else { return $false }
}

function Timeoutset {

    If (Test-Path ($RegKey + "\Control Panel")) {
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel"
        If (Test-Path ($RegKey + "\Desktop")) {
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel\Desktop"
            #ScreenSaveTimeOut
            New-ItemProperty -path $RegKey -name ScreenSaveTimeOut -value $time_out -PropertyType String -Force
        }
        else {
            New-Item -path $RegKey -name "Desktop" -Force
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel\Desktop"
            #ScreenSaveTimeOut
            New-ItemProperty -path $RegKey -name ScreenSaveTimeOut -value $time_out -PropertyType String -Force 
        }
    }
    else {
        New-Item -path $RegKey -name "Control Panel" -Force
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel"
        New-Item -path $RegKey -name "Desktop" -Force
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel\Desktop"
        #ScreenSaveTimeOut
        New-ItemProperty -path $RegKey -name ScreenSaveTimeOut -value $time_out -PropertyType String -Force
    }
}

function SetPasswordSecureKey {

    If (Test-Path ($RegKey + "\Control Panel")) {
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel"
        If (Test-Path ($RegKey + "\Desktop")) {
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel\Desktop"
            ##Enabled
            New-ItemProperty -path $RegKey -name ScreenSaverIsSecure -value $value -PropertyType String -Force  
        }
        else {
            New-Item -path $RegKey -name "Desktop"
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel\Desktop"
            ##Enabled
            New-ItemProperty -path $RegKey -name ScreenSaverIsSecure -value $value -PropertyType String     
        }
    }
    else {
        New-Item -path $RegKey -name "Control Panel"
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel"
        New-Item -path $RegKey -name "Desktop"
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Control Panel\Desktop"
        ##Enabled
        New-ItemProperty -path $RegKey -name ScreenSaverIsSecure -value $value -PropertyType String  
    }
}

#Check for validation
validation

# Get each user profile SID and Path to the profile
[array]$UserProfiles = Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList\*" | Where { $_.PSChildName -match "S-1-5-21-(\d+-?){4}$" } | Select-Object @{Name = "SID"; Expression = { $_.PSChildName } }, @{Name = "UserHive"; Expression = { "$($_.ProfileImagePath)\NTuser.dat" } }, @{name = "Username"; expression = { $_.ProfileImagePath -replace '^(.*[\\\/])', '' } }

$UserProfiles += @( New-Object PSObject -Property @{SID = ".DEFAULT"; UserHive = "C:\Users\Public\ntuser.dat"; username = "DEFAULT" })

# Loop through each profile on the machine</p>
Foreach ($UserProfile in $UserProfiles) {
    # Load User ntuser.dat if it's not already loaded
    If (($ProfileWasLoaded = Test-Path Registry::HKEY_USERS\$($UserProfile.SID)) -eq $false) {
        Start-Process -FilePath "CMD.EXE" -ArgumentList "/C REG.EXE LOAD HKU\$($UserProfile.SID) $($UserProfile.UserHive)" -Wait -WindowStyle Hidden
        Start-Sleep 1
    }
    $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)"
   
    #enable/disable screen saver
    if ($Action -eq "Enable") {
        $EnableValue = 1
   
        if (SetKey) {
            Write-Output "Screen Saver($Action) configured successfully for user '$($UserProfile.username)'"
        } 
        else {
            Write-Error "Screen Saver($Action) configuration failed for user '$($UserProfile.username)'" 
        }
        
        #set timeout for screensaver
        Timeoutset | Out-Null

        #set password protect
        if ($OnResumePasswordProtected -eq "Enable") {
            $value = 1
            SetPasswordSecureKey | Out-Null
        }
        if ($OnResumePasswordProtected -eq "Disable") {
            $value = 0
            SetPasswordSecureKey | Out-Null
        }
        
        #change screen saver 
        If ($ChangeScreenSaver -eq $true) {

            if ($StoreTheImageFileOn -eq "Local Machine") {
        
                SetScreenSaver -localpath $localpath | Out-Null
            }
            else {
                           
                SetScreenSaver -localpath $localpath | Out-Null     
            }
        } 
    }

    if ($Action -eq "Disable") {
        $EnableValue = 0

        if (SetKey) {
            Write-Output "Screen Saver($Action) configured successfully for user '$($UserProfile.username)'"
        } 
        else {
            Write-Error "Screen Saver($Action) configuration failed for user '$($UserProfile.username)'" 
        }
    }
   
    # Unload NTuser.dat        
    If ($ProfileWasLoaded -eq $false) {
        [gc]::Collect()
        Start-Sleep 1
        Start-Process -FilePath "CMD.EXE" -ArgumentList "/C REG.EXE UNLOAD HKU\$($UserProfile.SID)" -Wait -WindowStyle Hidden | Out-Null
    }
        
} #Foreach closed
