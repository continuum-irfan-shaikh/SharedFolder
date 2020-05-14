<#
    .SYNOPSIS
        Add Local Intranet Sites
    .DESCRIPTION
        Add Local Intranet Sites. Zones defined as 1 (Local Intranet), 2 (Trusted Sites), 3 (Internet), and 4 (Restricted Sites)
    .Help
        HKLM:\Software\Policies\Microsoft\Windows\CurrentVersion\Internet Settings\ZoneMap\Domains
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>
###Define $urls as variable in JSON SCHEMA. It can accept new line input
#$urls = "https://continuum.atlassian.net
#         http://abcdefgh.com"     
# $urls --> Data Type :- String --> Multiline Textbox  --> Mandatory --> Title :- "URL:"


if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
} 

$urlsArr = @()
$urlsArr = $urls -split [System.Environment]::NewLine
$zone = 1 # Options are 1 (Local Intranet), 2 (Trusted Sites), 3 (Internet), and 4 (Restricted Sites)
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
    $uri = [system.URI]$urlsArr[$i]
    $scheme = ($uri).Scheme
    $domainname = ($uri).host -replace '^www\.'
    $domainRegPath = "$regkey\$domainname"
    if (($uri).IsAbsoluteUri -eq $true) {
        if (Test-Path $domainRegPath -ErrorAction SilentlyContinue) {
            if (-not (isSettingExists -registryPath $domainRegPath -name $scheme)) {
                New-ItemProperty -Path $domainRegPath -Name $scheme -Value $zone -PropertyType "DWord" >$null
            }
            else {
                Set-ItemProperty -Path $domainRegPath -Name $scheme -Value $zone >$null
            }
            Write-Output "`nURL '$domainname' in the list of intranet sites updated."
        }
        else {
            New-Item -Path $domainRegPath -Force >$null
            New-ItemProperty -Path $domainRegPath -Name $scheme -Value $zone -PropertyType "DWord" >$null
            Write-Output "`nURL '$domainname' added to the list of intranet sites."
        }
    }
    else {
        Write-Output "`nURL '$uri' not added to intranet sites as it is not a absolute URL. "
    }
}
