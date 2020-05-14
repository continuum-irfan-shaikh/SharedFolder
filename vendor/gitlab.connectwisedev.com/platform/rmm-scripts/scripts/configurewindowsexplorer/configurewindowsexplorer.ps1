<#
  .Script
  Configure windows explorer.
  .Author
  Nirav Sachora.
  .Version
  2.0 and above.
  .Requirements
  Script should run with highest privileges.
#>

<#$removefolderoptions
$removefilemenu
$disconnectnetworkdrive
$removecontextmenu#>

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

if($removefolderoptions -eq "Enable"){
    if(Create_registry -path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\policies\Explorer" -Name "NoFolderOptions" -Value 0 -propertytype "DWord"){
        Write-Output "Folder options has been enabled."
    }
}
elseif($removefolderoptions -eq "Disable"){
    if(Create_registry -path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\policies\Explorer" -Name "NoFolderOptions" -Value 1 -propertytype "DWord"){
        Write-Output "Folder options has been disabled."
    }
}

if($removefilemenu -eq "Enable"){
    if(Create_registry -path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\policies\Explorer" -Name "NoFileMenu" -Value 0 -propertytype "DWord"){
        Write-Output "File menu has been enabled."
    }
}
elseif($removefilemenu -eq "Disable"){
    if(Create_registry -path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\policies\Explorer" -Name "NoFileMenu" -Value 1 -propertytype "DWord"){
        Write-Output "File menu has been disabled."
    }
}

if($disconnectnetworkdrive -eq "Enable"){
    if(Create_registry -path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\policies\Explorer" -Name "NoNetConnectDisconnect" -Value 0 -propertytype "DWord"){
        Write-Output "network drive has been enabled."
    }
}
elseif($disconnectnetworkdrive -eq "Disable"){
    if(Create_registry -path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\policies\Explorer" -Name "NoNetConnectDisconnect" -Value 1 -propertytype "DWord"){
        Write-Output "network drive has been disabled."
    }
}

if($removecontextmenu -eq "Enable"){
    if(Create_registry -path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\policies\Explorer" -Name "NoViewContextMenu" -Value 0 -propertytype "DWord"){
        Write-Output "Context menu has been enabled."
    }
}
elseif($removecontextmenu -eq "Disable"){
    if(Create_registry -path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\policies\Explorer" -Name "NoViewContextMenu" -Value 1 -propertytype "DWord"){
        Write-Output "context menu has been disabled."
    }
}




