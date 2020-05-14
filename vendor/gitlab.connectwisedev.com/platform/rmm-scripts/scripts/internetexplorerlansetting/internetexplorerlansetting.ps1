try{
$ErrorActionPreference = 'Stop'
$Key = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings"
$LANSettings = Get-ItemProperty $Key
$AutoProxyDetectMode = ((Get-ItemProperty -Path "$key\Connections" -Name  DefaultConnectionSettings).DefaultConnectionSettings)[8]
$Validate = {if (!$args[0]) {'Not Configured'}else {$args[0]}} # scriptblock to avoid writing same code more than once

if ($LANSettings) { 
    Write-Output "`nAuto-Configuration Proxy: $($Validate.invoke($LANSettings.AutoConfigProxy))"
    Write-Output "Auto-Configuration URL: $($Validate.invoke($LANSettings.AutoConfigURL))"
    Write-Output "Auto Proxy Detection Mode: $(if($AutoProxyDetectMode -ge 9){'Enabled'}Else{'Disabled'})"
    Write-Output "Proxy: $(if($LANSettings.ProxyEnable){'Enabled'}else{'Disabled'})"
    Write-Output "Proxy Override: $($Validate.invoke($LANSettings.ProxyOverride))"
    Write-Output "Proxy Server: $($Validate.invoke($LANSettings.ProxyServer))"
} 
}
catch{
    Write-Error $_.Exception.Message
}
