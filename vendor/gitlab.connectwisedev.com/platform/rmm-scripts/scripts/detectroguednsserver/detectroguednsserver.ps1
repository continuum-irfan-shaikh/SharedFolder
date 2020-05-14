# List of Rogue Public DNS Server IP Addresses
# https://kc.mcafee.com/resources/sites/MCAFEE/content/live/PRODUCT_DOCUMENTATION/23000/PD23652/en_US/McAfee_Labs_Threat_Advisory_DNSChanger.pdf

Function Test-IsRogueDNS ($IPAddress) {
    $RogueDNSRanges = @(
        "85.255.112.0-85.255.127.255",
        "67.210.0.0-67.210.15.255",
        "93.188.160.0-93.188.167.255",
        "77.67.83.0-77.67.83.255",
        "213.109.64.0-213.109.79.255",
        "64.28.176.0-64.28.191.255"
    )
    foreach ($range in $RogueDNSRanges) {
        $from, $to = $range -split "-"
        return ([System.Version]$from -le [System.Version]$IPAddress-and [System.Version]$IPAddress -le [System.Version]$to)
    }
}       

try {
    $RogueDNSNics = $false
    Foreach($Network in $(Get-WmiObject -Class Win32_NetworkAdapterConfiguration -ErrorAction Stop)) {
        $RogueDNSIP = @() 
        Foreach($DNS in $Network.DNSServerSearchOrder){
            If(Test-IsRogueDNS -IPAddress $DNS){
                $RogueDNSNics = $true
                Write-Output "`nRogue DNS Found: `'$DNS`'`n"
                Write-Output "Network Config:Â "
                Write-Output "Name:Â $($Network.Description)" 
                Write-Output "IP Address: $($Network.IPAddress -join ', ')"
                Write-Output "DNS: $($Network.DNSServerSearchOrder -join ', ')"
            }
        }
    }  

    If(!$RogueDNSNics){
        Write-Output "`nRogue DNS not found"
    }
}
catch {
  Write-Error $_
}

