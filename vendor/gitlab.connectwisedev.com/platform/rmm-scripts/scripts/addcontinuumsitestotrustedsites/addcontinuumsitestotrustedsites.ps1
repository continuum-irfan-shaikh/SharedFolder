<#
    .SYNOPSIS
        Add Continuum Sites to Trusted Sites.
    .DESCRIPTION
        Add Continuum Sites to Trusted Sites. Zones defined as 1 (Local Intranet), 2 (Trusted Sites), 3 (Internet), and 4 (Restricted Sites)
        Default sites to be added to trusted sites.
        http://webpost.itsupport247.net
        http://update1.itsupport247.net/
        http://update.itsupport247.net/
        https://update.itsupport247.net/
        http://rc.itsupport247.net/ 

    .Help
        HKLM:\Software\Policies\Microsoft\Windows\CurrentVersion\Internet Settings\ZoneMap\Domains
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>
###Define $urls as variable in JSON SCHEMA. It can take multiple values saperated by comma(,)
#$urls = ("http://googles.com", "https://www.abcd.com")

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
} 

$urls = "http://webpost.itsupport247.net,http://update1.itsupport247.net/,http://update.itsupport247.net/,https://update.itsupport247.net/,http://rc.itsupport247.net/"

#Function to check if heky entry exists or not. 
function isSettingExists {
    Param($registryPath, $name)
    if (Test-Path $registryPath -PathType container) {
        $key = Get-Item -LiteralPath $registryPath
        if ($key.GetValue($name, $null) -ne $null) {
            return $true
        }
    }
    return $false
}

#Array to add URL
$urlsArr = @()
$urlsArr = $urls -split ","
$zone = 2 # Options are 1 (Local Intranet), 2 (Trusted Sites), 3 (Internet), and 4 (Restricted Sites)

#Delete Domains RegKey under HKLM Policies to prevent override over user's registry. 
$regtodel = "HKLM:\Software\Policies\Microsoft\Windows\CurrentVersion\Internet Settings\ZoneMap\Domains"
if (get-childitem $regtodel -ErrorAction SilentlyContinue) {
    remove-item $regtodel -recurse -ErrorAction SilentlyContinue
}

#Delete Domains RegKey under HKLM to prevent override over user's registry. 
$regtodel1 = "HKLM:\Software\Microsoft\Windows\CurrentVersion\Internet Settings\ZoneMap\Domains"
if (get-childitem $regtodel1 -ErrorAction SilentlyContinue) {
    remove-item $regtodel1 -recurse -ErrorAction SilentlyContinue
}

# Get each user profile Pschildname and Path to the profile
$UserProfiles = (Get-ChildItem "Registry::HKEY_USERS" | where { $_.Name -match 'S-\d-\d+-(\d+-){1,14}\d+$' }).Pschildname


if (![string]::isnullorempty($UserProfiles)) {
    #loop 
    For ($i = 0; $i -lt $urlsArr.Length; $i++) { 

        # Loop through each profile on the machine</p>
        Foreach ($UserProfile in $UserProfiles) {
       
            $uname = Get-ItemProperty "Registry::HKEY_USERS\$UserProfile\Volatile Environment" | select -ExpandProperty USERNAME

            #Delete Registry under policies to prevent override.
            $regtodel2 = "Registry::HKEY_USERS\$UserProfile\Software\Policies\Microsoft\Windows\CurrentVersion\Internet Settings\ZoneMap\Domains"
            if (get-childitem $regtodel2 -ErrorAction SilentlyContinue) {
                remove-item $regtodel2 -recurse -ErrorAction SilentlyContinue
            }

            $RegKey = "Registry::HKEY_USERS\$UserProfile\Software"

            #Adding missing registries. 
            if (!(Test-Path "Registry::HKEY_USERS\$UserProfile\Software")) {
                New-Item -path "Registry::HKEY_USERS\$UserProfile" -name "Software" | Out-Null
            }
        
            if (!(Test-Path "$RegKey\Microsoft")) {
                New-Item -path "$RegKey" -name "Microsoft" | Out-Null
            }
        
            if (!(Test-Path "$RegKey\Microsoft\Windows")) {
                New-Item -path "$RegKey\Microsoft" -name "Windows" | Out-Null
            }

            if (!(Test-Path "$RegKey\Microsoft\Windows\CurrentVersion")) {
                New-Item -path "$RegKey\Microsoft\Windows" -name "CurrentVersion" | Out-Null
            }
        
            if (!(Test-Path "$RegKey\Microsoft\Windows\CurrentVersion\Internet Settings")) {
                New-Item -path "$RegKey\Microsoft\Windows\CurrentVersion" -name "Internet Settings" | Out-Null
            }     

            if (!(Test-Path "$RegKey\Microsoft\Windows\CurrentVersion\Internet Settings\ZoneMap")) {
                New-Item -path "$RegKey\Microsoft\Windows\CurrentVersion\Internet Settings" -name "ZoneMap" | Out-Null
            }

            if (!(Test-Path "$RegKey\Microsoft\Windows\CurrentVersion\Internet Settings\ZoneMap\Domains")) {
                New-Item -path "$RegKey\Microsoft\Windows\CurrentVersion\Internet Settings\ZoneMap" -name "Domains" | Out-Null
            }

            Start-Sleep 1

            $regkey1 = "Registry::HKEY_USERS\$UserProfile\Software\Microsoft\Windows\CurrentVersion\Internet Settings\ZoneMap\Domains"

            #Check URI. 
            $uri = [system.uri]($urlsArr[$i]).replace("://*.", "://") 
            $scheme = ($uri).Scheme
            $domainname = ($uri).host -replace '^www\.'
            $domainRegPath = "$regkey1\$domainname"
        
            if (![string]::isnullorempty($scheme)) {  
                if (Test-Path $domainRegPath -ErrorAction SilentlyContinue) {
                    #Create/Update HKEY entries. 
                    if (-not (isSettingExists -registryPath $domainRegPath -name $scheme)) {
                        New-ItemProperty -Path $domainRegPath -Name $scheme -Value $zone -PropertyType "DWord" >$null
                    }
                    else {
                        Set-ItemProperty -Path $domainRegPath -Name $scheme -Value $zone >$null
                    }
                    Write-Output "`n Domain $domainname updated to trusted sites for $uname."
                }
                else {
                    New-Item -Path $domainRegPath -Force >$null
                    New-ItemProperty -Path $domainRegPath -Name $scheme -Value $zone -PropertyType "DWord" >$null
                    Write-Output "`n Domain $domainname added to trusted sites for $uname."
                }
            }
            else {
                Write-Output  "`n Domain not added for user $uname as '$($urlsArr[$i])' is Invalid URL. Kindly provide correct URL" 
            }
        }
    }
}
else {
    Write-Output "No user currently logged in to system $ENV:COMPUTERNAME"
}
