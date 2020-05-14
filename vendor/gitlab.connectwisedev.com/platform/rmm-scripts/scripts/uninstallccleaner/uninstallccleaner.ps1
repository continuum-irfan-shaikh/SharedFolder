<#
    .Script
    Uninstall CCleaner
    .Description
    Script will check for the version provided through portal and uninstall if same version is installed on the system.
    .Requirements
    Run script with highest privileges
    
#>

<# 
$action       ===> Type: Dropdown ===> ENUM: "Select Action", "Uninstall" ===>  Title : "Action"
$version      ===> Type: Dropdown ===> ENUM: "Select SubVersion", "5.30", "5.31", "5.32", "5.33", "5.34", "5.35", "5.36", "5.37", "5.38", "5.39"  ===> 
#>


if (Test-Path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\CCleaner") {
    $registryvalue = Get-ItemProperty -Path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\CCleaner" | Select -ExpandProperty DisplayVersion
    if ($version -ne $registryvalue) {
        Write-Output "provided Version of CCleaner is not installed."
        Exit;
    }
    else {
        Function Execute_process($path, $bit) {
            $pinfo = New-Object System.Diagnostics.ProcessStartInfo
            if ($bit -eq 32 -and (test-path "$path\uninst.exe")) {
                $pinfo.FileName = "$path\uninst.exe"
            }
            elseif ($bit -eq 64 -and (test-path "$path\uninst.exe")) {
                $pinfo.FileName = "$path\uninst.exe"
            }
            else {
                return -1
            }
            $pinfo.RedirectStandardError = $true
            $pinfo.RedirectStandardOutput = $true
            $pinfo.UseShellExecute = $false
            $pinfo.Arguments = "/S"
            $p = New-Object System.Diagnostics.Process
            $p.StartInfo = $pinfo
            $p.Start() | Out-Null
            $p.WaitForExit()
            return $p.ExitCode
        }
        
        
        
        if ((Test-path "C:\Program Files\CCleaner") -or (Test-Path "C:\Program Files (x86)\CCleaner")) {
            Write-Output "CCleaner execution is started, it may take time to complete the execution"
            $osarchitecture = (get-wmiobject -Class win32_operatingsystem).Osarchitecture
            if ($osarchitecture -eq "32-bit") {
                $path = Join-Path "C:\Program Files" "Ccleaner"
                if (test-path $path) {
                    $success = Execute_process -path $path -bit 32
                }
            }
        
            elseif ($osarchitecture -eq "64-bit") {
                $path32 = Join-Path "C:\Program Files (x86)" "Ccleaner"
                $path64 = Join-Path "C:\Program Files" "Ccleaner"
                if (Test-Path $path64) {
                    $success = Execute_process -path $path64 -bit 64
                } 
                elseif (Test-Path $path32) {
                    $success = Execute_process -path $path32 -bit 32
                }
            }
            if ($success -eq 0) {
                Write-Output "CCleaner has been removed"
            }
            elseif ($success -eq -1) {
                Write-Error "uninst.exe file is missing from C:\Program files\CCleaner"
            }
            else {
                Write-Error "Action Could not be completed"
            }
        }
    }
}
else {
    Write-Output "CCleaner is not installed on this system"
    Exit;
}
