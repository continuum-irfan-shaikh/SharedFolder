<#
    .SYNOPSIS
        Disable the windows update notification in system tray for windows 10 for currently logged on users
    .DESCRIPTION
        Disable the windows update notification in system tray for windows 10 for currently logged on users
    .Help
        To get more details refer below registry. 
        HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Notifications\Settings\Windows.SystemToast.SecurityAndMaintenance\Enabled = 0 (DWORD)
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>

#Check win10 OS
if (!(((Get-WmiObject -Class Win32_OperatingSystem).caption) -match "Microsoft Windows 10")) {
    Write-Output "This task is supported only for windows 10 operating systems"  
    Exit;
}

Function Create_registry($path, $Name, $Value, $propertytype) {
    # Function will set registry value and if item is not present will create registry entry.
    $Details = Get-ItemProperty -Path $path
    if ($Details -ne $null) {
        $Details = Get-ItemProperty -Path $path | gm | select -ExpandProperty Name
        if ($Details -contains $Name) {
            Set-ItemProperty $path -Name $Name -Value $Value
        }
        else {
            New-ItemProperty -Path $path -Name $Name -PropertyType $propertytype -Value $Value | Out-Null
        }
    
    } #End if statement
    else {
        New-ItemProperty -Path $path -Name $Name -PropertyType $propertytype -Value $Value | Out-Null
    }  # End else statement

    if ((Get-ItemProperty $path -Name $Name | select -ExpandProperty $Name) -eq $Value) { return $true } else { return $false }
}

$profilelist = Get-Item "REgistry::HKU\S-1-5-21-*" | Where-Object { $_.Name -notlike '*_Classes' } | select -ExpandProperty name
 
if ($profilelist) {

    foreach ($userprofile in $profilelist) {

        if (Test-path "Registry::$userprofile\Volatile Environment") {
            $username = Get-ItemProperty -Path "Registry::$userprofile\Volatile Environment" -Name Username | Select -ExpandProperty Username
                    
            $RegKey = "registry::$userprofile\SOFTWARE\Microsoft\Windows\CurrentVersion\Notifications"

            if (!(Test-Path $RegKey)) {
                New-Item -Path "registry::$userprofile\SOFTWARE\Microsoft\Windows\CurrentVersion" -Name "Notifications" | Out-Null
            }
            if (!(Test-Path "$RegKey\Settings")) {
                New-Item -path "$RegKey" -name "\Settings" | Out-Null
            }
 
            if (!(Test-Path "$RegKey\Settings\Windows.SystemToast.SecurityAndMaintenance")) {
                New-Item "$RegKey\Settings" -name "Windows.SystemToast.SecurityAndMaintenance" | Out-Null
            }
   
            Start-Sleep 1

            $RegistryKey = "registry::$userprofile\SOFTWARE\Microsoft\Windows\CurrentVersion\Notifications\Settings\Windows.SystemToast.SecurityAndMaintenance"
    
            #Call create_registry function
            if (Create_registry -path "$RegistryKey" -Name "Enabled" -Value '0' -propertytype DWORD) {
                Write-Output "Windows update notification disabled successfully for user '$username' on system $ENV:COMPUTERNAME"
            }
            else {
                Write-Output "Failed to disable Windows update notification for user '$username' on system $ENV:COMPUTERNAME"
            }
        }
        else {
            Continue
        }
    }
}

else {
    Write-Output "Could not disable windows update notification as currently no user was logged on to system $ENV:COMPUTERNAME."
}
