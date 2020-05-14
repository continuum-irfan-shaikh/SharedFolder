# $version = '4.5.1' # user input

# 4
# 4.5.1
# 4.5.2
# 4.6
# 4.6.01055
# 4.6.1
# 4.6.2
# 4.7
# 4.7.02558
# 4.7.2

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

Function GetDotNetVersion {
    $version = $null
    if (Test-Path -Path "HKLM:SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full") {
   
        switch ((Get-ItemProperty -Path "HKLM:SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full" -ErrorAction SilentlyContinue).Release) {
            378389 { $Version = "4.5" }
            378675 { $Version = "4.5.1" }
            378758 { $Version = "4.5.1" }
            379893 { $Version = "4.5.2" } 
            393295 { $Version = "4.6" }
            393297 { $Version = "4.6" }
            394254 { $Version = "4.6.1" } 
            394271 { $Version = "4.6.1" } 
            394802 { $Version = "4.6.2" } 
            394806 { $Version = "4.6.2" } 
            460798 { $Version = "4.7" } 
            460805 { $Version = "4.7" } 
            461308 { $Version = "4.7.1" } 
            461310 { $Version = "4.7.1" } 
            461808 { $Version = "4.7.2" } 
            461814 { $Version = "4.7.2" } 
        }
    }
       
    if (!$version -and (Test-Path -Path "HKLM:SOFTWARE\Microsoft\NET Framework Setup\NDP\")) {
        $dotnetvers = (Get-ChildItem "HKLM:SOFTWARE\Microsoft\NET Framework Setup\NDP\v*" -Name)
        foreach ($dotnetver in $dotnetvers) {
            switch ($dotnetver) {
                "v2.0" { $Version = "2" }
                "v3.0" { $Version = "3.0" }
                "v3.5" { $Version = "3.5" }
                "v4" { $Version = "4.0" }
            }
        }
    }
         
    return $Version
}

Try {
    $Architecture = if ((Get-WmiObject Win32_OperatingSystem).OSArchitecture -eq '32-bit') { '32' }else { '64' }

    $Data = @"
"AppName","SetUpName","Ins_Param","Unins_Param","Download_Link","UserName","Post_Ins_Check","BitTypeFull","Password","MD5CheckSum","ITSAppname","Description","AppVersion","BaseAppName","MajorVersion","SubVersion","Reboot","BitType"
"Microsoft .NET Framework 4.0","dotNetFx40_Full_x86_x64.exe","/q /norestart","MsiExec.exe /X{0A0CADCF-78DA-33C4-A350-CD51849B9702} /q /norestart","http://dcmdwld.itsupport247.net/dotNetFx40_Full_x86_x64.exe","","HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{0A0CADCF-78DA-33C4-A350-CD51849B9702}","Common","","251743DFD3FDA414570524BAC9E55381","Microsoft .NET Framework 4.0","Installs or uninstalls the application on Windows desktops and servers","","Microsoft .NET Framework","4","4","","0"
"Microsoft .NET Framework 4.5.1","NDP451-KB2858728-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{4903D172-DCCB-392F-93A3-34CA9D47FE3D} /q","http://dcmdwld.itsupport247.net/NDP451-KB2858728-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{4903D172-DCCB-392F-93A3-34CA9D47FE3D}","32 Bit","","607d3c9b399e3d303a1b14eb4326bd1e","Microsoft .NET Framework 4.5.1 (32-bit)","Installs or uninstalls the application on Windows desktops and servers","4.5.1","Microsoft .NET Framework","4","4.5.1","","32"
"Microsoft .NET Framework 4.5.1","NDP451-KB2858728-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{7DEBE4EB-6B40-3766-BB35-5CBBC385DA37} /q","http://dcmdwld.itsupport247.net/NDP451-KB2858728-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{7DEBE4EB-6B40-3766-BB35-5CBBC385DA37}","64 Bit","","607d3c9b399e3d303a1b14eb4326bd1e","Microsoft .NET Framework 4.5.1 (64-bit)","Installs or uninstalls the application on Windows desktops and servers","4.5.1","Microsoft .NET Framework","4","4.5.1","","64"
"Microsoft .NET Framework 4.5.2","NDP452-KB2901907-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{3911CF56-9EF2-39BA-846A-C27BD3CD0685} /q","http://dcmdwld.itsupport247.net/NDP452-KB2901907-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{3911CF56-9EF2-39BA-846A-C27BD3CD0685}","32 Bit","","ee01fc4110c73a8e5efc7cabda0f5ff7","Microsoft .NET Framework 4.5.2 (32-bit)","Installs or uninstalls the application on Windows desktops and servers","4.5.51209","Microsoft .NET Framework","4","4.5.2","","32"
"Microsoft .NET Framework 4.5.2","NDP452-KB2901907-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{26784146-6E05-3FF9-9335-786C7C0FB5BE} /q","http://dcmdwld.itsupport247.net/NDP452-KB2901907-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{26784146-6E05-3FF9-9335-786C7C0FB5BE}","64 Bit","","ee01fc4110c73a8e5efc7cabda0f5ff7","Microsoft .NET Framework 4.5.2 (64-bit)","Installs or uninstalls the application on Windows desktops and servers","4.5.51209","Microsoft .NET Framework","4","4.5.2","","64"
"Microsoft .NET Framework 4.6","NDP452-KB2901907-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{444C5574-6BE0-323E-9BDD-922F6C3C4A04} /q","http://dcmdwld.itsupport247.net/NDP46-KB3045557-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{444C5574-6BE0-323E-9BDD-922F6C3C4A04}","32 Bit","","88bc05e20114a4506f40c36911de92fa","Microsoft .NET Framework Ver.4.6 (32 bit)","Installs or uninstalls the application on Windows desktops and servers","4.6.00081","Microsoft .NET Framework","4","4.6","","32"
"Microsoft .NET Framework 4.6","NDP452-KB2901907-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{94A631D5-B30A-3DD8-B65C-1117C09DA73E} /q","http://dcmdwld.itsupport247.net/NDP46-KB3045557-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{94A631D5-B30A-3DD8-B65C-1117C09DA73E}","64 Bit","","88bc05e20114a4506f40c36911de92fa","Microsoft .NET Framework Ver.4.6 (64 bit)","Installs or uninstalls the application on Windows desktops and servers","4.6.00081","Microsoft .NET Framework","4","4.6","","64"
"Microsoft .NET Framework 4.6.1","NDP461-KB3102436-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{30500C7C-2206-3DC6-9792-96E95A04669D} /q","http://dcmdwld.itsupport247.net/NDP461-KB3102436-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{30500C7C-2206-3DC6-9792-96E95A04669D}","32 Bit","","864056903748706e251fec9f5d887ef9","Microsoft .NET Framework Ver.4.6.1 (32 bit)","Installs or uninstalls the application on Windows desktops and servers","4.6.01055","Microsoft .NET Framework","4","4.6.01055","","32"
"Microsoft .NET Framework 4.6.1","NDP461-KB3102436-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{BD6F5371-DAC1-30F0-9DDE-CAC6791E28C3} /q","http://dcmdwld.itsupport247.net/NDP461-KB3102436-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{BD6F5371-DAC1-30F0-9DDE-CAC6791E28C3}","64 Bit","","864056903748706e251fec9f5d887ef9","Microsoft .NET Framework Ver.4.6.1 (64 bit)","Installs or uninstalls the application on Windows desktops and servers","4.6.01055","Microsoft .NET Framework","4","4.6.01055","","64"
"Microsoft .NET Framework 4.6.1","NDP461-KB3102436-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{30500C7C-2206-3DC6-9792-96E95A04669D} /q","http://dcmdwld.itsupport247.net/NDP461-KB3102436-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{30500C7C-2206-3DC6-9792-96E95A04669D}","32 Bit","","864056903748706e251fec9f5d887ef9","Microsoft .NET Framework 4.6.1 (32-bit)","Installs or uninstalls the application on Windows desktops and servers","4.6.1","Microsoft .NET Framework","4","4.6.1","","32"
"Microsoft .NET Framework 4.6.1","NDP461-KB3102436-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{BD6F5371-DAC1-30F0-9DDE-CAC6791E28C3} /q","http://dcmdwld.itsupport247.net/NDP461-KB3102436-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{BD6F5371-DAC1-30F0-9DDE-CAC6791E28C3}","64 Bit","","864056903748706e251fec9f5d887ef9","Microsoft .NET Framework 4.6.1 (64-bit)","Installs or uninstalls the application on Windows desktops and servers","4.6.1","Microsoft .NET Framework","4","4.6.1","","64"
"Microsoft .NET Framework 4.6.2","NDP462-KB3151800-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{476F88BA-08DD-32D8-A8B0-E85EE28CB27F} /qn","http://dcmdwld.itsupport247.net/NDP462-KB3151800-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{476F88BA-08DD-32D8-A8B0-E85EE28CB27F}","32 Bit","","9a5d647ee710af2b1aede329c40bbe1a","Microsoft .NET Framework 4.6.2 (32 Bit)","Installs or uninstalls the application on Windows desktops and servers","4.6.01590","Microsoft .NET Framework","4","4.6.2","","32"
"Microsoft .NET Framework 4.6.2","NDP462-KB3151800-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{63DF5C4B-E3BF-3346-A033-C57B22F44C9E} /qn","http://dcmdwld.itsupport247.net/NDP462-KB3151800-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{63DF5C4B-E3BF-3346-A033-C57B22F44C9E}","64 Bit","","9a5d647ee710af2b1aede329c40bbe1a","Microsoft .NET Framework 4.6.2 (64 Bit)","Installs or uninstalls the application on Windows desktops and servers","4.6.01590","Microsoft .NET Framework","4","4.6.2","","64"
"Microsoft .NET Framework 4.6.2","NDP462-KB3151800-x86-x64-AllOS-ENU.exe","/q /norestart","","http://dcmdwld.itsupport247.net/NDP462-KB3151800-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full","Common","","9a5d647ee710af2b1aede329c40bbe1a","Microsoft .NET Framework 4.6.2 For Windows 8.1 and Above","Installs or uninstalls the application on Windows desktops and servers","4.6.01590","Microsoft .NET Framework","4","4.6.2 (8.1 and above)","","0"
"Microsoft .NET Framework 4.7","NDP47-KB3186497-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{2A842F3F-CE6D-3DFD-9ECB-9CC3C5150A67} /qn","http://dcmdwld.itsupport247.net/NDP47-KB3186497-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{2A842F3F-CE6D-3DFD-9ECB-9CC3C5150A67}","32 Bit","","b59040e489bff55433508438627d11e6","Microsoft .NET Framework 4.7 (32 Bit)","Installs or uninstalls the application on Windows desktops and servers","4.7","Microsoft .NET Framework","4","4.7","","32"
"Microsoft .NET Framework 4.7","NDP47-KB3186497-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{BCF0C1F7-671C-3922-A7EA-8AC11F4FC0EB} /qn","http://dcmdwld.itsupport247.net/NDP47-KB3186497-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{BCF0C1F7-671C-3922-A7EA-8AC11F4FC0EB}","64 Bit","","b59040e489bff55433508438627d11e6","Microsoft .NET Framework 4.7 (64 Bit)","Installs or uninstalls the application on Windows desktops and servers","4.7","Microsoft .NET Framework","4","4.7","","64"
"Microsoft .NET Framework 4.7","NDP47-KB3186497-x86-x64-AllOS-ENU.exe","/q /norestart","","http://dcmdwld.itsupport247.net/NDP47-KB3186497-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full","Common","","b59040e489bff55433508438627d11e6","Microsoft .Net Framework 4.7 For Windows 8.1 and Above","Installs or uninstalls the application on Windows desktops and servers","4.7","Microsoft .NET Framework","4","4.7 (8.1 and above)","","0"
"Microsoft .NET Framework 4.7.1","NDP471-KB4033342-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{4CB05D36-1518-395D-8C39-A102343CF661} /q","http://dcmdwld.itsupport247.net/NDP471-KB4033342-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\Software\Microsoft\Windows\CurrentVersion\Uninstall\{4CB05D36-1518-395D-8C39-A102343CF661}","32 Bit","","660e1a104f209f3cdb55b6d4e9ffa475","Microsoft .NET Framework 4.7.1 32 Bit","Installs or uninstalls the application on Windows desktops and servers","4.7.02558","Microsoft .NET Framework","4","4.7.02558","","32"
"Microsoft .NET Framework 4.7.1","NDP471-KB4033342-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{E0C7523C-686B-3EE6-8FB1-CB4339E30EDD} /q","http://dcmdwld.itsupport247.net/NDP471-KB4033342-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{E0C7523C-686B-3EE6-8FB1-CB4339E30EDD}","64 Bit","","660e1a104f209f3cdb55b6d4e9ffa475","Microsoft .NET Framework 4.7.1 64 Bit","Installs or uninstalls the application on Windows desktops and servers","4.7.02558","Microsoft .NET Framework","4","4.7.02558","","64"
"Microsoft .NET Framework 4.7.1","NDP471-KB4033342-x86-x64-AllOS-ENU.exe","/q /norestart","","http://dcmdwld.itsupport247.net/NDP471-KB4033342-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full","Common","","660e1a104f209f3cdb55b6d4e9ffa475","Microsoft .Net Framework 4.7.1 For Windows 8.1 and Above","Installs or uninstalls the application on Windows desktops and servers","4.7.02558","Microsoft .NET Framework","4","4.7.1 (8.1 and above)","","0"
"Microsoft .NET Framework 4.7.2","NDP472-KB4054530-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{10C4E843-C226-3FDF-9DD6-F4E3275E734D} /q","http://dcmdwld.itsupport247.net/NDP472-KB4054530-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{10C4E843-C226-3FDF-9DD6-F4E3275E734D}","32 Bit","","F1F3EA28AD5F41DD366E7067BE8B5124","Microsoft .NET Framework 4.7.2 32 Bit","Installs or uninstalls the application on Windows desktops and servers","4.7.03062","Microsoft .NET Framework","4","4.7.2","","32"
"Microsoft .NET Framework 4.7.2","NDP472-KB4054530-x86-x64-AllOS-ENU.exe","/q /norestart","MsiExec.exe /X{09CCBE8E-B964-30EF-AE84-6537AB4197F9} /q","http://dcmdwld.itsupport247.net/NDP472-KB4054530-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{09CCBE8E-B964-30EF-AE84-6537AB4197F9}","64 Bit","","F1F3EA28AD5F41DD366E7067BE8B5124","Microsoft .NET Framework 4.7.2 64 Bit","Installs or uninstalls the application on Windows desktops and servers","4.7.03062","Microsoft .NET Framework","4","4.7.2","","64"
"Microsoft .NET Framework 4.7.2","NDP472-KB4054530-x86-x64-AllOS-ENU.exe","/q /norestart","","http://dcmdwld.itsupport247.net/NDP472-KB4054530-x86-x64-AllOS-ENU.exe","","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full","Common","","F1F3EA28AD5F41DD366E7067BE8B5124","Microsoft .Net Framework 4.7.2 for Windows 8.1 and Above","Installs or uninstalls the application on Windows desktops and servers","4.7.03062","Microsoft .NET Framework","4","4.7.2 (8.1 and Above)","","0"
"@ -replace "msiexec.exe ", '' -replace "msiexec", "" | ConvertFrom-Csv | Where-Object {$_.SubVersion -and $_.BitType -and $_.Unins_Param -and $_.Post_Ins_Check}

    $SupportedVersions = $Data | Select-Object -ExpandProperty SubVersion -Unique
    if (!($SupportedVersions -contains $version)) { Write-Output "Uninstallation of .Net Framework (v$version) is not supported."; exit; }

    $data = $data | Where-Object { $_.SubVersion -eq $version -and $_.BitType -eq $Architecture }
    $ValidateProductReg = Test-Path "Registry::$($data.Post_Ins_Check)"
    $InstalledVersion = GetDotNetVersion

    if ($InstalledVersion -and !($ValidateProductReg)) {
        Write-Output "Failed to Uninstall .Net Framework (v$version) either because it is not installed or shipped with the Operating System.";
    }
    elseif ($InstalledVersion -and $ValidateProductReg) {
        $RestartArgument = '/norestart'
        $Arguments = "$(($data.Unins_Param).trim()) $RestartArgument"
        # "Version: $InstalledVersion"
        # "Arguments: $Arguments"
        # "Post isntall registry: $($data.Post_Ins_Check)"
        Start-Process Msiexec.exe -ArgumentList $Arguments  -Wait 
        if (!(Test-Path "Registry::$($data.Post_Ins_Check)")) {
            Write-Output "Successfuly Uninstalled .Net Framework (v$version)."
        }
        else {
            #if ($data.Post_Ins_Check -match '{\w{8}-\w{4}-\w{4}-\w{4}-\w{12}}') {$matches[0]}
            Write-Output "Failed to Uninstall .Net Framework (v$version).`nPost Install Check: Registry `'$($data.Post_Ins_Check)`' doesn't exists"
        }
    }
    else{
        Write-Output "Failed to Uninstall."
    }
}
catch {
    Write-Output "Failed to Uninstall."
    Write-Error $_
}
