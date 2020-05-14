<#
    .SYNOPSIS
       Remove LogMeIn from Start menu
    .DESCRIPTION
       Remove LogMeIn from Start menu
    .Help
        Use below path for start menu icons.
        $Env:ProgramData\Microsoft\Windows\Start Menu\Programs
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>

function get-logmein {
    if ((gwmi win32_operatingsystem | select osarchitecture).osarchitecture -eq "64-bit") {
        $a = Get-ChildItem -Path HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object {$_.DisplayName -match "LogMeIn"}| Select-Object -ExpandProperty UninstallString
    }
    else {
        $a = Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object {$_.DisplayName -match "LogMeIn"} | Select-Object -expandProperty UninstallString
    }
    return $a
}
    
Try {
    if (get-logmein) {
    
        $Path = "$Env:ProgramData\Microsoft\Windows\Start Menu\Programs"
        $values = Get-ChildItem $Path -Recurse -Include *.lnk | Where-Object {$_.Name -match 'LogMeIn'} | Remove-Item -Force
        
        if ($values -eq $null) {
            Write-Output "`nLogMeIn removed successfully from start menu on system $ENV:COMPUTERNAME."
        }
    }
    else {
        Write-Output "`nLogMeIn not installed on this system $ENV:COMPUTERNAME"
    }
}
catch {
    Write-Error $_.Exception.Message
}
