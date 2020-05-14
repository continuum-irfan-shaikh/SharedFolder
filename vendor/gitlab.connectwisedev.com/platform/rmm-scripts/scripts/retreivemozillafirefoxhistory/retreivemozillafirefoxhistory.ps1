<#
.Synopsis
Script will retreive Mozilla firefox history for each user
.Author
Nirav Sachora
.Requirements
Sqlite3.exe should present in system.
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
if(!$days){
$days = 90
}
$date = (get-date).AddDays(-$days)
$datestr = $date.ToString("yyyy-MM-dd")

$OSArch = (Get-WmiObject Win32_OperatingSystem).OSArchitecture
Switch ($OSArch) {
    "64-bit" { $sqlite3 = ${ENV:ProgramFiles(x86)} + "\ITSPlatform\plugin\scripting\sqlite3.exe" }                              
    "32-bit" { $sqlite3 = ${ENV:ProgramFiles} + "\ITSPlatform\plugin\scripting\sqlite3.exe" }               
}

$mozillaexist = Test-Path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Mozilla*"
if ($mozillaexist -eq $False) {
    Write-Output "Mozilla is not installed on this system"
    exit;
}

$sqlite = Test-Path "$sqlite3"
if ($sqlite) {
    if ($users -eq "allusers") {
        $profiledata = Get-ChildItem 'HKLM:\Software\Microsoft\Windows NT\CurrentVersion\ProfileList' | ForEach-Object { $_.GetValue('ProfileImagePath') } # Details about each User profile in the system
        $profilepath = @()
        foreach ($input in $profiledata) {
            if ($input -like "C:\Users\*") {
                $profilepath += "$input" + "\AppData\Roaming\Mozilla\Firefox\Profiles"
            }
        }
    }
    else {
        $profilepath = @()
        $notexist = @()
        $userprofiles = $usernames -split ","
        foreach ($profile1 in $userprofiles) {
            if (Test-Path "C:\Users\$profile1") {
                $profilepath += "C:\Users\$profile1" + "\AppData\Roaming\Mozilla\Firefox\Profiles"
            }
            else {
                $notexist += $profile1
            }
        }
    }
 
    foreach ($profile in $profilepath) {
        $Profileexist = Test-Path $profile
        if ($Profileexist) {
            $foldername = Get-ChildItem $profile | Where-Object { $_.Name -like "*.default" } | select -ExpandProperty Name
        
            $username = $profile -replace "\\AppData\\Roaming\\Mozilla\\Firefox\\Profiles", ""; $username = $username -replace "C:\\Users\\", ""
            $profile = $profile.replace("\", "\\")
            #if($name -ne $null)
                
            $path = $profile + "\\" + $foldername + "\\places.sqlite"
          [string]$data = ".open $path
                SELECT datetime(moz_historyvisits.visit_date/1000000,'unixepoch'), moz_places.url
                FROM moz_places, moz_historyvisits 
                WHERE moz_places.id = moz_historyvisits.place_id and date(moz_historyvisits.visit_date/1000000,'unixepoch') >= $datestr;
                "
           
            $history = $data | & $sqlite3
            #$history
            foreach ($record in $history) {
                $lastdate = (Get-Date).AddDays(-$days)
                $a, $b = $record.split("|")
                $c = [uri]$b; $c = $c.scheme + "://" + $c.Host
                $a = Get-Date $a
                if ($a -gt $lastdate) {
                    $obj = New-Object psobject -Property @{
                        Site     = $c
                        Username = $username
                        Date     = $a
                   
                    }
                    $obj | select Username, Date, Site
                }
                
            }
                
        }
    }
}

else {
    Write-Error "SQLite DB client not found. Unable to connect history database."
    Exit
}

if($notexist.length -gt 0){
Write-Output "`nBelow profiles does not exists on the system. `n$notexist"
  } 
