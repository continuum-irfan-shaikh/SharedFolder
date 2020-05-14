<#
    .SYNOPSIS
        Install Defraggler
    .DESCRIPTION
        Install Defraggler
    .Help
        
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#> 
<#
JSON Schema
$action           Example:- $action = "uninstall"
$version
$subversion

SubVersion
2.18
2.19
2.20
2.21

Only subversion variable will be used in script. 
#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

$action = "Install"

if ($action -eq "install") {

    $subversionpath = @{
        '2.18' = "http://dcmdwld.itsupport247.net/dfsetup218_slim.exe"
        '2.19' = "http://dcmdwld.itsupport247.net/dfsetup219_slim.exe"
        '2.20' = "http://dcmdwld.itsupport247.net/dfsetup220_slim.exe"
        '2.21' = "http://dcmdwld.itsupport247.net/dfsetup221_slim.exe"
    }

    Function Prod_Install {
    
        $url = $subversionpath["$version"]
        #$url = "https://download.ccleaner.com/dfsetup222.exe"
        $downloadpath = "C:\Windows\Temp\dfsetup_abcd.exe"
        $wc = New-Object System.Net.WebClient
        $wc.DownloadFile("$url", "$downloadpath")
        $ps = New-Object System.Diagnostics.Process
        $ps.StartInfo.Filename = $downloadpath
        $ps.StartInfo.Arguments = "/S"
        $ps.StartInfo.RedirectStandardOutput = $True
        $ps.StartInfo.UseShellExecute = $false
        $ps.start()
        $ps.WaitForExit()
        return $ps.ExitCode
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



 $Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall', 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'

    $Product_Check = Get-ChildItem $Registry -ErrorAction 'SilentlyContinue' | Get-ItemProperty | Where-Object { $_.DisplayName -match "Defraggler" } | select -First 1

    if ($Product_Check) {
        Write-Output "`n"
        Write-Output "Defraggler is already installed on system $ENV:COMPUTERNAME"
        EXIT;
    }

    $OSArch = (Get-WmiObject Win32_OperatingSystem).OSArchitecture 
    Switch ($OSArch) {
        "64-bit" { $RegKey = "registry::HKLM\SOFTWARE\Wow6432Node" }
        "32-bit" { $RegKey = "registry::HKLM\Software" }
        } 
    
    $items = @("Google", "No Chrome Offer Until")
    
    $ErrorActionPreference = "SilentlyContinue"
    foreach ($item in $items) {
    
        New-Item -path $RegKey -name $item -Force | Out-Null; $RegKey = $RegKey + "\" + "$item"
        
    }
    $ErrorActionPreference = "Continue"
    
    if (Test-Path $RegKey) {
    
        $result = Create_registry -path $RegKey -name "Piriform Ltd" -Value 20991231 -propertytype DWORD
    
    }
     
    if ($result) {
    
        if (Prod_Install -eq 0) {
            Write-Output "`n"
            Write-Output "Defraggler v$version installed successfully on system $ENV:COMPUTERNAME"
        
        }
        else {
            Write-Output "`n"
            Write-Output "Failed to install Defraggler v$version on system $ENV:COMPUTERNAME"
         
        }   
    }

}
else {
    Write-Output "`nKindly select 'Install' as an action"
}
