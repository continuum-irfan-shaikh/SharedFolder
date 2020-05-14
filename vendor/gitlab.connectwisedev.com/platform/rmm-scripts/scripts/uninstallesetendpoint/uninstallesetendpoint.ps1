<# Parameter $version - type string
    list of versions = ["4.2.76.0","5.0.2214.4","5.0.2214.4(32 bit)","5.0.2214.4(64 bit)","5.0.2216.0","5.0.2226.0",
   "5.0.2228.1","5.0.2248.0","5.0.2254.0","6.2.2033.0","6.3.2016.0"]

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

Try {
    $Architecture = if ((Get-WmiObject Win32_OperatingSystem).OSArchitecture -eq '32-bit') {'32'}else {'64'}
    $data = @"
    "AppName","Unins_Param","Post_Ins_Check","ITSAppname","AppVersion","SubVersion","Reboot","BitType"
    "ESET Endpoint Antivirus","msiexec /X {D6AE00FF-2367-42A9-A60D-2997ED93ECFA} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{D6AE00FF-2367-42A9-A60D-2997ED93ECFA}","ESET Endpoint Antivirus 5.0.2228.without reboot","5.0.2228.1","5.0.2228.1","0","32"
    "ESET Endpoint Antivirus","msiexec /X {D6AE00FF-2367-42A9-A60D-2997ED93ECFA} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{D6AE00FF-2367-42A9-A60D-2997ED93ECFA}","ESET Endpoint Antivirus 5.0.2228 with reboot","5.0.2228.1","5.0.2228.1","1","32"
    "ESET Endpoint Antivirus","msiexec /X {29177C87-01ED-4D9A-9C7A-32D3061A93AD} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{29177C87-01ED-4D9A-9C7A-32D3061A93AD}","ESET Endpoint Antivirus 5.0.2228 64 bit without reboot","5.0.2228.1","5.0.2228.1","0","64"
    "ESET Endpoint Antivirus","msiexec /X {29177C87-01ED-4D9A-9C7A-32D3061A93AD} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{29177C87-01ED-4D9A-9C7A-32D3061A93AD}","ESET Endpoint Antivirus 5.0.2228 64 bit with reboot","5.0.2228.1","5.0.2228.1","1","64"
    "ESET NOD32 Antivirus","msiexec /X {9CEC1801-DB68-48CE-B74F-5733BBD3F729} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{9CEC1801-DB68-48CE-B74F-5733BBD3F729}","ESET NOD32 Antivirus without reboot","4.2.76.0","4.2.76.0","0","32"
    "ESET NOD32 Antivirus","msiexec /X {9CEC1801-DB68-48CE-B74F-5733BBD3F729} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{9CEC1801-DB68-48CE-B74F-5733BBD3F729}","ESET NOD32 Antivirus with reboot","4.2.76.0","4.2.76.0","0","64"
    "ESET Endpoint Antivirus","msiexec /X {3187B3B0-3620-4459-A983-4403FC481420} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{3187B3B0-3620-4459-A983-4403FC481420}","ESET Endpoint Antivirus 5.0.2214 without reboot","5.0.2214.4","5.0.2214.4","0","32"
    "ESET Endpoint Antivirus","msiexec /X {3187B3B0-3620-4459-A983-4403FC481420} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{3187B3B0-3620-4459-A983-4403FC481420}","ESET Endpoint Antivirus 5.0.2214 with reboot","5.0.2214.4","5.0.2214.4","0","64"
    "ESET NOD32 Antivirus","MsiExec.exe /X{A1A01D26-AF53-42C0-9DAE-1BC2FCC68812} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{A1A01D26-AF53-42C0-9DAE-1BC2FCC68812}","Uninstall ESET Antivirus 8.0.304.0 32 bit without reboot","8.0.304.0","8.0.304.0","0","32"
    "ESET NOD32 Antivirus","MsiExec.exe /X{7F39EB28-B9B7-41B8-8564-DB33284A010D} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{7F39EB28-B9B7-41B8-8564-DB33284A010D}","Uninstall ESET Antivirus 8.0.304.0 64 bit without reboot","8.0.304.0","8.0.304.0","0","64"
    "ESET NOD32 Antivirus","MsiExec.exe /X{A1A01D26-AF53-42C0-9DAE-1BC2FCC68812} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{A1A01D26-AF53-42C0-9DAE-1BC2FCC68812}","Uninstall ESET ntivirus 8.0.304.0 32 bit with reboot","8.0.304.0","8.0.304.0","1","32"
    "ESET NOD32 Antivirus","MsiExec.exe /X{7F39EB28-B9B7-41B8-8564-DB33284A010D} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{7F39EB28-B9B7-41B8-8564-DB33284A010D}","Uninstall ESET Antivirus 8.0.304.0 64 bit with reboot","8.0.304.0","8.0.304.0","1","64"
    "ESET NOD32 Antivirus","MsiExec.exe /x{B096B8AB-C3BD-4801-A731-D2B94643DA86} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{B096B8AB-C3BD-4801-A731-D2B94643DA86}","Uninstall ESETAntivirus 8.0.312.0 32 bit without reboot","8.0.312.0","8.0.312.0","0","32"
    "ESET NOD32 Antivirus","MsiExec.exe /X{D6885DDE-4632-4640-A3BB-13C9F02CE81C} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{D6885DDE-4632-4640-A3BB-13C9F02CE81C}","Uninstall ESET Antivirus 8.0.312.0 64 bit without reboot","8.0.312.0","8.0.312.0","0","64"
    "ESET NOD32 Antivirus","MsiExec.exe /x{B096B8AB-C3BD-4801-A731-D2B94643DA86} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{B096B8AB-C3BD-4801-A731-D2B94643DA86}","Uninstall ESETAntivirus 8.0.312.0 32 bit with reboot","8.0.312.0","8.0.312.0","1","32"
    "ESET NOD32 Antivirus","MsiExec.exe /X{D6885DDE-4632-4640-A3BB-13C9F02CE81C} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{D6885DDE-4632-4640-A3BB-13C9F02CE81C}","Uninstall ESET Antivirus 8.0.312.0 64 bit with reboot","8.0.312.0","8.0.312.0","1","64"
    "ESET NOD32 Antivirus","MsiExec.exe /X{1231238A-E793-4030-A068-0E0A2643B8E3} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{1231238A-E793-4030-A068-0E0A2643B8E3}","Uninstall ESET Antivirus 8.0.319.0 32 bit without reboot","8.0.319.0","8.0.319.0","0","32"
    "ESET NOD32 Antivirus","MsiExec.exe /X{5F2AE448-CD4B-40BD-B245-5F0CD06A09B0} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{5F2AE448-CD4B-40BD-B245-5F0CD06A09B0}","Uninstall ESET Antivirus 8.0.319.0 64 bit without reboot","8.0.319.0","8.0.319.0","0","64"
    "ESET NOD32 Antivirus","MsiExec.exe /X{1231238A-E793-4030-A068-0E0A2643B8E3} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{1231238A-E793-4030-A068-0E0A2643B8E3}","Uninstall ESET Antivirus 8.0.319.0 32 bit with reboot","8.0.319.0","8.0.319.0","1","32"
    "ESET NOD32 Antivirus","MsiExec.exe /X{5F2AE448-CD4B-40BD-B245-5F0CD06A09B0} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{5F2AE448-CD4B-40BD-B245-5F0CD06A09B0}","Uninstall ESET Antivirus 8.0.319.0 64 bit with reboot","8.0.319.0","8.0.319.0","1","64"
    "ESET Endpoint Antivirus","MsiExec.exe /X{A9DF0706-AF66-4300-8D75-E6751243A8D9} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{A9DF0706-AF66-4300-8D75-E6751243A8D9}","ESET Endpoint Antivirus 6.2.2033.0 without reboot (32 bit)","6.2.2033.0","6.2.2033.0","0","32"
    "ESET Endpoint Antivirus","MsiExec.exe /X{13189425-6C52-490A-9E5A-3B66DB545629} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{13189425-6C52-490A-9E5A-3B66DB545629}","ESET Endpoint Antivirus 6.2.2033.0 without reboot (64 bit)","6.2.2033.0","6.2.2033.0","0","64"
    "ESET Endpoint Antivirus","MsiExec.exe /X{A9DF0706-AF66-4300-8D75-E6751243A8D9} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{A9DF0706-AF66-4300-8D75-E6751243A8D9}","ESET Endpoint Antivirus 6.2.2033.0 with reboot (32 bit)","6.2.2033.0","6.2.2033.0","1","32"
    "ESET Endpoint Antivirus","MsiExec.exe /X{13189425-6C52-490A-9E5A-3B66DB545629} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{13189425-6C52-490A-9E5A-3B66DB545629}","ESET Endpoint Antivirus 6.2.2033.0 with reboot (64 bit)","6.2.2033.0","6.2.2033.0","1","64"
    "ESET Endpoint Antivirus","MsiExec.exe /X{0D1090DA-441D-4068-BEF4-0A9DBB493243} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{0D1090DA-441D-4068-BEF4-0A9DBB493243}","ESET Endpoint Antivirus 6.3.2016.0 without reboot (32 bit)","6.3.2016.0","6.3.2016.0","0","32"
    "ESET Endpoint Antivirus","MsiExec.exe /X{8FA06832-4954-42E1-AB68-1FC4EA923BA7} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{8FA06832-4954-42E1-AB68-1FC4EA923BA7}","ESET Endpoint Antivirus 6.3.2016.0 without reboot (64 bit)","6.3.2016.0","6.3.2016.0","0","64"
    "ESET Endpoint Antivirus","MsiExec.exe /X{0D1090DA-441D-4068-BEF4-0A9DBB493243} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{0D1090DA-441D-4068-BEF4-0A9DBB493243}","ESET Endpoint Antivirus 6.3.2016.0 with reboot (32 bit)","6.3.2016.0","6.3.2016.0","1","32"
    "ESET Endpoint Antivirus","MsiExec.exe /X{8FA06832-4954-42E1-AB68-1FC4EA923BA7} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{8FA06832-4954-42E1-AB68-1FC4EA923BA7}","ESET Endpoint Antivirus 6.3.2016.0 with reboot (64 bit)","6.3.2016.0","6.3.2016.0","1","64"
    "ESET Endpoint Antivirus","msiexec.exe /x {AC92FE0A-588B-4A4D-9D33-F139B94BEA45} /qn","HKEY_LOCAL_MACHINE\Software\Microsoft\Windows\CurrentVersion\Uninstall\{08716A83-6648-4CB3-901C-8B0E34579C2A}","ESET Endpoint Antivirus 5.0.2248.0","5.0.2248.0","5.0.2248.0","0","32"
    "ESET Endpoint Antivirus","msiexec /X {3187B3B0-3620-4459-A983-4403FC481420} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{3187B3B0-3620-4459-A983-4403FC481420}","Uninstall ESET Endpoint Antivirus 5.0.2214.4 64 bit","5.0.2214.4","5.0.2214.4 (64 bit)","0","64"
    "ESET Endpoint Antivirus","msiexec /X {E45ED219-0786-4576-BCE5-1639B73864AC} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{E45ED219-0786-4576-BCE5-1639B73864AC}","Uninstall ESET Endpoint Antivirus 5.0.2214.4 32 bit","5.0.2214.4","5.0.2214.4 (32 bit)","0","32"
    "ESET Endpoint Antivirus","msiexec /X {8CD0547B-8CE7-4CAA-A413-A43B2C4D11FF} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{8CD0547B-8CE7-4CAA-A413-A43B2C4D11FF}","Uninstall ESET Endpoint Antivirus 5.0.2225.0 64 bit","5.0.2225.0","5.0.2225.0","0","64"
    "ESET Endpoint Antivirus","msiexec /X {E0F72657-6B20-496B-A884-1B09453FE7E7} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{E0F72657-6B20-496B-A884-1B09453FE7E7}","Uninstall ESET Endpoint Antivirus 5.0.2225.0 32 bit","5.0.2225.0","5.0.2225.0","0","32"
    "ESET Endpoint Antivirus","msiexec /X {AC6D818C-8EA0-4C2D-AFF9-9B0FA3E0E105} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{AC6D818C-8EA0-4C2D-AFF9-9B0FA3E0E105}","Uninstall ESET Endpoint Antivirus 5.0.2254.0 64 bit","5.0.2254.0","5.0.2254.0","0","64"
    "ESET Endpoint Antivirus","msiexec /X {B654805A-15D7-4CB3-AB70-EF6A321DB2F9} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{B654805A-15D7-4CB3-AB70-EF6A321DB2F9}","Uninstall ESET Endpoint Antivirus 5.0.2254.0 32 bit","5.0.2254.0","5.0.2254.0","0","32"
    "ESET Endpoint Antivirus","msiexec /X {6B0C987E-EED3-4D07-828E-FCE96B14BD5C} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{6B0C987E-EED3-4D07-828E-FCE96B14BD5C}","Uninstall ESET Endpoint Antivirus 5.0.2126.0 32 bit","5.0.2126.0","5.0.2126.0","0","32"
    "ESET Endpoint Antivirus","msiexec /X {FF8AC853-B984-4C9A-937A-1F20FB6AA6B9} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{FF8AC853-B984-4C9A-937A-1F20FB6AA6B9}","Uninstall ESET Endpoint Antivirus 5.0.2126.0 64 bit","5.0.2126.0","5.0.2126.0","0","64"
    "ESET Endpoint Security","msiexec.exe /x {54A2D672-5A36-4573-A429-A7BC5F2B6B1F} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{54A2D672-5A36-4573-A429-A7BC5F2B6B1F}","ESET Endpoint Security 6.6.2046.0 with reboot (32 Bit)","6.6.2046.0","6.6.2046.0","1","32"
    "ESET Endpoint Security","msiexec.exe /x {E6DCEC5A-066D-41B4-8EF8-3BB97F864DC8} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{E6DCEC5A-066D-41B4-8EF8-3BB97F864DC8}","ESET Endpoint Security 6.6.2046.0 with reboot (64 Bit)","6.6.2046.0","6.6.2046.0","1","64"
    "ESET Endpoint Security","msiexec.exe /x {54A2D672-5A36-4573-A429-A7BC5F2B6B1F} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{54A2D672-5A36-4573-A429-A7BC5F2B6B1F}","ESET Endpoint Security 6.6.2046.0 without reboot (32 Bit)","6.6.2046.0","6.6.2046.0","0","32"
    "ESET Endpoint Security","msiexec.exe /x {E6DCEC5A-066D-41B4-8EF8-3BB97F864DC8} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{E6DCEC5A-066D-41B4-8EF8-3BB97F864DC8}","ESET Endpoint Security 6.6.2046.0 without reboot (64 Bit)","6.6.2046.0","6.6.2046.0","0","64"
    "ESET Endpoint Security","msiexec.exe /x {677F552D-ED35-4234-BDFF-09662A68D870} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{677F552D-ED35-4234-BDFF-09662A68D870}","ESET Endpoint Security 6.6.2052.0 with reboot (32 Bit)","6.6.2052.0","6.6.2052.0","1","32"
    "ESET Endpoint Security","msiexec.exe /x {32F4EA9B-2E7B-4384-9F3D-C2A692F1A814} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{32F4EA9B-2E7B-4384-9F3D-C2A692F1A814}","ESET Endpoint Security 6.6.2052.0 with reboot (64 Bit)","6.6.2052.0","6.6.2052.0","1","64"
    "ESET Endpoint Security","msiexec.exe /x {677F552D-ED35-4234-BDFF-09662A68D870} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{677F552D-ED35-4234-BDFF-09662A68D870}","ESET Endpoint Security 6.6.2052.0 without reboot (32 Bit)","6.6.2052.0","6.6.2052.0","0","32"
    "ESET Endpoint Security","msiexec.exe /x {32F4EA9B-2E7B-4384-9F3D-C2A692F1A814} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{32F4EA9B-2E7B-4384-9F3D-C2A692F1A814}","ESET Endpoint Security 6.6.2052.0 without reboot (64 Bit)","6.6.2052.0","6.6.2052.0","0","64"
    "ESET Endpoint Security","msiexec.exe /x {201FF105-BB06-46A9-97A3-71DF7EA5953E} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{201FF105-BB06-46A9-97A3-71DF7EA5953E}","ESET Endpoint Security 6.6.2064.0 with reboot (32 Bit)","6.6.2064.0","6.6.2064.0","1","32"
    "ESET Endpoint Security","msiexec.exe /x {15347789-1267-4CFA-8B8A-C7795631C007} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{15347789-1267-4CFA-8B8A-C7795631C007}","ESET Endpoint Security 6.6.2064.0 with reboot (64 Bit)","6.6.2064.0","6.6.2064.0","1","64"
    "ESET Endpoint Security","msiexec.exe /x {201FF105-BB06-46A9-97A3-71DF7EA5953E} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{201FF105-BB06-46A9-97A3-71DF7EA5953E}","ESET Endpoint Security 6.6.2064.0 without reboot (32 Bit)","6.6.2064.0","6.6.2064.0","0","32"
    "ESET Endpoint Security","msiexec.exe /x {15347789-1267-4CFA-8B8A-C7795631C007} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{15347789-1267-4CFA-8B8A-C7795631C007}","ESET Endpoint Security 6.6.2064.0 without reboot (64 Bit)","6.6.2064.0","6.6.2064.0","0","64"
    "ESET Endpoint Security","msiexec.exe /x {34CD32DE-2CFD-49B6-B68C-ACFB1C5A7559} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{34CD32DE-2CFD-49B6-B68C-ACFB1C5A7559}","ESET Endpoint Security 6.6.2068.1 with reboot (32 Bit)","6.6.2068.1","6.6.2068.1","1","32"
    "ESET Endpoint Security","msiexec.exe /x {071C2D74-B758-4FFF-9E06-84138EFB2675} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{071C2D74-B758-4FFF-9E06-84138EFB2675}","ESET Endpoint Security 6.6.2068.1 with reboot (64 Bit)","6.6.2068.1","6.6.2068.1","1","64"
    "ESET Endpoint Security","msiexec.exe /x {34CD32DE-2CFD-49B6-B68C-ACFB1C5A7559} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{34CD32DE-2CFD-49B6-B68C-ACFB1C5A7559}","ESET Endpoint Security 6.6.2068.1 without reboot (32 Bit)","6.6.2068.1","6.6.2068.1","0","32"
    "ESET Endpoint Security","msiexec.exe /x {071C2D74-B758-4FFF-9E06-84138EFB2675} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{071C2D74-B758-4FFF-9E06-84138EFB2675}","ESET Endpoint Security 6.6.2068.1 without reboot (64 Bit)","6.6.2068.1","6.6.2068.1","0","64"
    "ESET Endpoint Security","msiexec.exe /x {8E9AB319-74FD-490E-864F-21B29AEC6A90} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{8E9AB319-74FD-490E-864F-21B29AEC6A90}","ESET Endpoint Security 6.6.2072.4 with reboot (32 Bit)","6.6.2072.4","6.6.2072.4","1","32"
    "ESET Endpoint Security","msiexec.exe /x {632B1C53-E8E5-4F68-817C-EBA1E9098FB7} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{632B1C53-E8E5-4F68-817C-EBA1E9098FB7}","ESET Endpoint Security 6.6.2072.4 with reboot (64 Bit)","6.6.2072.4","6.6.2072.4","1","64"
    "ESET Endpoint Security","msiexec.exe /x {8E9AB319-74FD-490E-864F-21B29AEC6A90} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{8E9AB319-74FD-490E-864F-21B29AEC6A90}","ESET Endpoint Security 6.6.2072.4 without reboot (32 Bit)","6.6.2072.4","6.6.2072.4","0","32"
    "ESET Endpoint Security","msiexec.exe /x {632B1C53-E8E5-4F68-817C-EBA1E9098FB7} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{632B1C53-E8E5-4F68-817C-EBA1E9098FB7}","ESET Endpoint Security 6.6.2072.4 without reboot (64 Bit)","6.6.2072.4","6.6.2072.4","0","64"
    "ESET Endpoint Security","msiexec.exe /x {2C546AFD-6849-4BD7-AC8D-D5BDDF5B7ECF} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{2C546AFD-6849-4BD7-AC8D-D5BDDF5B7ECF}","ESET Endpoint Security 6.6.2078.5 with reboot (32 Bit)","6.6.2078.5","6.6.2078.5","1","32"
    "ESET Endpoint Security","msiexec.exe /x {209B659A-4692-4FD4-B2BD-F6AD65EB5300} /qn","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{209B659A-4692-4FD4-B2BD-F6AD65EB5300}","ESET Endpoint Security 6.6.2078.5 with reboot (64 Bit)","6.6.2078.5","6.6.2078.5","1","64"
    "ESET Endpoint Security","msiexec.exe /x {2C546AFD-6849-4BD7-AC8D-D5BDDF5B7ECF} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{2C546AFD-6849-4BD7-AC8D-D5BDDF5B7ECF}","ESET Endpoint Security 6.6.2078.5 without reboot (32 Bit)","6.6.2078.5","6.6.2078.5","0","32"
    "ESET Endpoint Security","msiexec.exe /x {209B659A-4692-4FD4-B2BD-F6AD65EB5300} /qn /norestart","HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\{209B659A-4692-4FD4-B2BD-F6AD65EB5300}","ESET Endpoint Security 6.6.2078.5 without reboot (64 Bit)","6.6.2078.5","6.6.2078.5","0","64"    
"@ -replace "msiexec.exe ", '' -replace "msiexec",""
  
    $data = $data | ConvertFrom-Csv | Where-Object {$_.SubVersion -eq $version -and $_.BitType -eq $Architecture}
    $GUID = if ($data.Post_Ins_Check -match '{\w{8}-\w{4}-\w{4}-\w{4}-\w{12}}') {$matches[0]}
    Function VerifyProduct($GUID) {
        $Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall', 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'
        $Product = Get-ChildItem $Registry -ErrorAction 'SilentlyContinue' | Get-ItemProperty | Where-Object { $_.UninstallString -like "*$GUID*"}
        Return [bool]($Product)
    }
        Start-Process Msiexec.exe -ArgumentList ($data.Unins_Param).trim() -Wait
        if (!(VerifyProduct $GUID)) {
            Write-Output "Successfuly Uninstalled software."
        }
        else {
            Write-Output "Failed to Uninstall software."
        }
}
catch {
    Write-Output "Failed to Uninstall."
    Write-Error $_
}
