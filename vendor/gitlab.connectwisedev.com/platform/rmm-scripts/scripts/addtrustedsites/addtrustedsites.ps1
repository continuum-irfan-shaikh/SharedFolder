$urlsArr = @()
$urlsArr = $urls -split "`n"
$zone = 2 # Options are 1 (Local Intranet), 2 (Trusted Sites), 3 (Internet), and 4 (Restricted Sites)
$regkey = "HKLM:\Software\Policies\Microsoft\Windows\CurrentVersion\Internet Settings\ZoneMap\Domains"
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
For ($i = 0; $i -lt $urlsArr.Length; $i++) {
        
    try {
        $ErrorActionPreference = 'Stop'
        $uri = [system.uri]($urlsArr[$i]).replace("://*.", "://") 
        $scheme = ($uri).Scheme
        $domainname = ($uri).host -replace '^www\.'
        $domainRegPath = "$regkey\$domainname"
        
        if (Test-Path $domainRegPath -ErrorAction SilentlyContinue) {
            if (-not (isSettingExists -registryPath $domainRegPath -name $scheme)) {
                New-ItemProperty -Path $domainRegPath -Name $scheme -Value $zone -PropertyType "DWord" >$null
            }
            else {
                Set-ItemProperty -Path $domainRegPath -Name $scheme -Value $zone >$null
            }
            Write-Output "`n Domain $domainname in the list of trusted sites was updated."
        }
        else {
            New-Item -Path $domainRegPath -Force >$null
            New-ItemProperty -Path $domainRegPath -Name $scheme -Value $zone -PropertyType "DWord" >$null
            Write-Output "`n Domain $domainname was added to the list of trusted sites."
        }
        
    } 
    catch { 
        Write-Output  "`n '$($urlsArr[$i])' is Invalid URL. Kindly provide correct URL" 
    }
} 
