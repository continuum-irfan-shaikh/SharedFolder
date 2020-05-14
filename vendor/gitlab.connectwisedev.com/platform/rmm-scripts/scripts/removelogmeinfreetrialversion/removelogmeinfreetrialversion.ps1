<#
    .SYNOPSIS
       Remove LogMeIn trial version
    .DESCRIPTION
       Remove LogMeIn trial version
    .Help
        To get more details refer below command. 
        HKLM:\SOFTWARE\LogMeIn\V5  for License check
        Use MSI uninstallation method.
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>

$ErrorActionPreference = "SilentlyContinue"

function TrialVersion {
    
    $b = Get-ChildItem -Path HKLM:\SOFTWARE\LogMeIn\ | Get-ItemProperty | Select-Object -ExpandProperty license 

    $ErrorActionPreference = "Continue"

    if ($b -match "SERVICE_FREE" -or $b -match "SERVICERATRIAL") {
        return $true
    }
    else {
        return $false
    }
}

Try {
    if (TrialVersion) {
    
        $program = "LogMeIn"
        
        Get-WmiObject -Class Win32_Product | where-Object {$_.Name -match $program} | foreach { $_.Uninstall() } 2>&1 | out-null

        $check = Get-WmiObject -Class Win32_Product | where-Object {$_.Name -match $program}
                
        if ($check -eq $null) {
            Write-Output "`nLogMeIn Trial version removed from the system $ENV:ComputerName"
        }
        else {
            Write-Output "Something went wrong while uninstalling. Kindly check manually."
        }
    }
    else {
        Write-Output "`nLogMeIn not installed or licensed version already installed on this system $ENV:COMPUTERNAME"
    }
}
catch {
    Write-Error $_.Exception.Message
}
