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
        1.1
#>

<########## Variables to define in JSON #################
$action = "Enable"     ## Enable  OR Disable
$timeout = 1    #Timeout in minutes  Default is 1 minute
$PasswordProtect = "Enable"  ##Enable OR Disable

$ChangeScreenSaver = $true   ## True or False

$ImageFileFrom = "Local"  # Local OR Network     
$scrpath = "C:\Durgesh\scr"       #Local screensaver path where screensaver file available
$scrfile = "Bubbles111.scr"

#Network drive screensaver details
$NetworkPath = '\\SCRIPT-WIN10-64\Durgesh\scr' #Network path start with \\
$scrfile = "Bubbles111.scr"  # Filename with extension .scr
$UserName = 'labadmin'   #'administrator'  #Username to connect network path
$Password = 'lic@123'   #'Grt@2018'        #Password to connect network path
###################>

$time_out = $timeout * 60 
$LocalScreenSaverPath = "C:\Windows\System32"
$ErrorActionPreference = 'SilentlyContinue'
    
function connect {

    $global:DriveLetter = Get-ChildItem function:[g-z]: -n | Where-Object { !(Test-Path $_) } | random

    $global:Net = New-Object -ComObject WScript.Network
 
    $ErrorActionPreference = 'SilentlyContinue'

    if ([string]::IsNullOrEmpty($UserName) -and [string]::IsNullOrEmpty($Password)) {
        $Net.MapNetworkDrive($DriveLetter, "$NetworkPath", $false)
        if (Get-WmiObject -Class Win32_MappedLogicalDisk | Where-Object { $_.name -eq "$DriveLetter" }) {
            return $true
        }
        else {
            return $false
        }
    }
    else {
        $Net.MapNetworkDrive($DriveLetter, "$NetworkPath", $false, $UserName, $Password)
        if (Get-WmiObject -Class Win32_MappedLogicalDisk | Where-Object { $_.name -eq "$DriveLetter" }) {
            return $true
        }
        else {
            return $false
        }
    }
    
}

function removedrive {

    $Net.RemoveNetworkDrive($DriveLetter)
    Start-Sleep 2

}

function validation {
    switch ($ImageFileFrom) {
        "Local" {
            $global:localpath = Join-Path "$scrpath" "$scrfile"
            if ((Test-Path $localpath) -eq $false) {
                Write-Output "`nLocalpath '$localpath' of the screen saver is incorrect or screen saver file does not exist."
                Exit;
            }
        }
        "Network" {
            if (connect) {
                Start-Sleep 2
                $mypath = Join-path "$DriveLetter" "$scrfile"
                if (Test-Path $mypath) {
                    if (Test-Path $LocalScreenSaverPath) {
                
                        Copy-Item -path "$DriveLetter\$scrfile" "$LocalScreenSaverPath" -Force
                
                        Start-Sleep 2 
                 
                        $global:localpath = join-path "$LocalScreenSaverPath" "$scrfile"
                    
                        removedrive
                    }
                    else {
                        Write-Output "`nLocal screen saver path '$LocalScreenSaverPath' is incorrect"
                        removedrive
                        Exit;
                    }
                } 
                else {
                    Write-Output "`nScreenSaver file '$scrfile' does not exist at '$NetworkPath'"
                    removedrive
                    Exit;
                }
            }
            else {
                Write-Output "`nNetwork Drive '$NetworkPath' not mapped. Kindly provide correct network path or credentials"
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
    if ($action -eq "Enable") {
        $EnableValue = 1
   
        if (SetKey) {
            Write-Output "Screen Saver($action) configured successfully for user '$($UserProfile.username)'"
        } 
        else {
            Write-Error "Screen Saver($action) configuration failed for user '$($UserProfile.username)'" 
        }
        
        #set timeout for screensaver
        Timeoutset | Out-Null

        #set password protect
        if ($PasswordProtect -eq "Enable") {
            $value = 1
            SetPasswordSecureKey | Out-Null
        }
        if ($PasswordProtect -eq "Disable") {
            $value = 0
            SetPasswordSecureKey | Out-Null
        }
        
        #change screen saver 
        If ($ChangeScreenSaver -eq $true) {

            if ($ImageFileFrom -eq "Local") {
        
                SetScreenSaver -localpath $localpath | Out-Null
            }
            else {
                           
                SetScreenSaver -localpath $localpath | Out-Null     
            }
        } 
    }

    if ($action -eq "Disable") {
        $EnableValue = 0

        if (SetKey) {
            Write-Output "Screen Saver($action) configured successfully for user '$($UserProfile.username)'"
        } 
        else {
            Write-Error "Screen Saver($action) configuration failed for user '$($UserProfile.username)'" 
        }
    }
   
    # Unload NTuser.dat        
    If ($ProfileWasLoaded -eq $false) {
        [gc]::Collect()
        Start-Sleep 1
        Start-Process -FilePath "CMD.EXE" -ArgumentList "/C REG.EXE UNLOAD HKU\$($UserProfile.SID)" -Wait -WindowStyle Hidden | Out-Null
    }
        
} #Foreach closed
