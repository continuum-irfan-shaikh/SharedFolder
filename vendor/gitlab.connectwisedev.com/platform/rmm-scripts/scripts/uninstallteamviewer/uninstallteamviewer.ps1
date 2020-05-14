<#
    name : Uninstall TeamViewer
    Category : Application
   
    .SYNOPSIS
         TeamViewer uninstallation from all desktops and servers
    .DESCRIPTION
         TeamViewer uninstallation from all desktops and servers
    .Help
         Use uninstallation HKEY
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>

########
# Define $version in JSON schema file. Use drop down and put values like 5,6,7,8,9,10,11,12,13
# Example :- $version = "11"
########

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}


function verify {
    if ((gwmi win32_operatingsystem | select osarchitecture).osarchitecture -eq "64-bit") {
        $a = Get-ChildItem -Path HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object {$_.DisplayName -match "TeamViewer" } | Select-Object -ExpandProperty UninstallString
    }
    else {
        $a = Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object {$_.DisplayName -match "TeamViewer" } | Select-Object -expandProperty UninstallString
    }
    return $a
}


function versioncheck {

    if ((gwmi win32_operatingsystem | select osarchitecture).osarchitecture -eq "64-bit") {
        $b = Get-ChildItem -Path HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object {$_.DisplayName -match "TeamViewer" } | Select-Object -ExpandProperty VersionMajor
    }
    else {
        $b = Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object {$_.DisplayName -match "TeamViewer" } | Select-Object -expandProperty VersionMajor
    }
    return $b
    
}


Try {
    $uninstallkey = verify
    $versionnum = versioncheck
    if ((!$uninstallkey) -or ($versionnum -ne $version)) {
        Write-Output "`nTeamViewer $version not installed on this system $ENV:ComputerName"
    }
    else {
    
        $process = Start-Process $uninstallkey -arg "/S" -Wait -PassThru -ErrorAction 'Stop'
        Start-Sleep 10
       
        $ver1 = verify
        if (!($ver1)) {
           
            Write-Output "`nTeamViewer $version is successfully uninstalled on the system $ENV:ComputerName" 
                      
        } 
        else {
            Write-Output "`nFailed to uninstall TeamViewer $version on the system $ENV:ComputerName" 
        }
    }
}
catch {
    
    write-output "`n"$_.Exception.Message
}
