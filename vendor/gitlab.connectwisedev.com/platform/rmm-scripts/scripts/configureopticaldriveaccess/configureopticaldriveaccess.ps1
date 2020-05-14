<#
    .Script
    Script will enable or disable optical drive access.
    .Author
    Nirav Sachora
    .Description
    Script will Enable or disable optical drive access, script have option to exclude administrators..
    .Requirements
    Script should run with admin privileges.
#>

<#$OpticalDriveAccess = "Enable"
$ExcludeAdmin = $false#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
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

foreach ($profile in $profilelist) {
    if (Test-Path "Registry::$profile\Volatile Environment") {
        $username = Get-ItemProperty -Path "Registry::$profile\Volatile Environment" -Name Username | select -ExpandProperty Username
        if ($ExcludeAdmin -eq $true) {
            $isadmin = net.exe localgroup administrators | where { $_ -eq "$username" }
        }
        else { $isadmin = $null }
        if ($isadmin -eq $null) {

            $Path = "Registry::$profile\Software\Policies\Microsoft\Windows\RemovableStorageDevices"
            $CD_DVD_GUID = "{53f56308-b6bf-11d0-94f2-00a0c91efb8b}"
            if ((Test-Path "Registry::$profile\Software\Policies\Microsoft\Windows\RemovableStorageDevices") -eq $false) {
                New-Item -Path "Registry::$profile\Software\Policies\Microsoft\Windows" -Name "RemovableStorageDevices" -ErrorAction Stop | Out-Null
            }

            $CD_DVD = Get-ChildItem -Path "Registry::$profile\Software\Policies\Microsoft\Windows\RemovableStorageDevices" | Where-Object { $_.Name -like "*$CD_DVD_GUID" } -ErrorAction SilentlyContinue
            if (!$CD_DVD) {   
                #Creating new entry '$CD_DVD_GUID'
                New-Item -Path $Path -Name $CD_DVD_GUID -ErrorAction Stop | Out-Null
            }
            if ($OpticalDriveAccess -eq "Enable") { $value = 0 } elseif ($OpticalDriveAccess -eq "Disable") { $value = 1 }
            $read = Create_registry -path "Registry::$profile\Software\Policies\Microsoft\Windows\RemovableStorageDevices\$CD_DVD_GUID" -Name "Deny_Read" -Value "$value" -propertytype "DWORD"
            $write = Create_registry -path "Registry::$profile\Software\Policies\Microsoft\Windows\RemovableStorageDevices\$CD_DVD_GUID" -Name "Deny_Write" -Value "$value" -propertytype "DWORD"
            $execute = Create_registry -path "Registry::$profile\Software\Policies\Microsoft\Windows\RemovableStorageDevices\$CD_DVD_GUID" -Name "Deny_Execute" -Value "$value" -propertytype "DWORD"
            if ($read -and $write -and $execute) {
                switch ($OpticalDriveAccess) {
                    "Enable" { "Optical Drive has been Enabled for $username" }
                    "Disable" { "Optical Drive has been Disabled for $username" }
                }
            }
        }
        else {
            Write-Output "$username account is member of Administrators group."
        }
    }
}
