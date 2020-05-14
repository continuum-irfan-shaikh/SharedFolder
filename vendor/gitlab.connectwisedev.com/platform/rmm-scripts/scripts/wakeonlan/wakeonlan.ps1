try {
    $ErrorActionPreference = 'Stop'
    # capture all network adapters that are physcal, active and not for wireless or bluetooth
    $nics = Get-WmiObject Win32_NetworkAdapter -filter "netconnectionstatus = 2 AND AdapterTypeID = '0' AND PhysicalAdapter = 'true' AND NOT Description LIKE '%Centrino%' AND NOT Description LIKE '%wireless%' AND NOT Description LIKE '%virtual%' AND NOT Description LIKE '%WiFi%' AND NOT Description LIKE '%Bluetooth%'"

    # capture all network adapters that have 
    $WakeOnLan = Get-WmiObject MSPower_DeviceWakeEnable -Namespace root\wmi | where {$_.enable -eq $true}
    
    If ($WakeOnLan) {
        $MagicPackets = Get-WmiObject MSNdis_DeviceWakeOnMagicPacketOnly -Namespace root\wmi| where { $_.EnableWakeOnMagicPacketOnly -eq $true}
        $MagicPacketsOutput = $WakeonLanOutput = @()
        
        # iterate each Network adapter and identify status of `WakeOnLan` and ability to wakeup on 'MagicPacketsOnly'
        Foreach ($nic in $nics) {
            Foreach ($item in $WakeOnLan) {
                If ($item.enable -and $item.instancename -match [regex]::escape($nic.PNPDeviceID)) {
                    $WakeonLanOutput += $nic |Select-Object Name, status, AdapterType, DeviceID, MacAddress, @{n = 'WakeOnLanEnabled?'; e = {$true}}
                } 
            }
            Foreach ($item in $MagicPackets) {
                If ($item.EnableWakeOnMagicPacketOnly -and $item.instancename -match [regex]::escape($nic.PNPDeviceID)) {
                    $MagicPacketsOutput += $nic |Select-Object Name, status, AdapterType, DeviceID, MacAddress, @{n = 'WakeOnMagicPacketOnly?'; e = {$true}}
                } 
            }
    
        } 
        If ($WakeonLanOutput) {
            Write-Output "`n'WakeOnLan' is Enabled on following 'Active' Network Adpaters"
            Write-Output $WakeonLanOutput
        }
        If ($MagicPacketsOutput) {
            Write-Output "'WakeOnMagicPacketOnly' is Enabled on following 'Active' Network Adpaters`n"
            Write-Output $MagicPacketsOutput
        }
    }
    else {
        Write-Output "`nThis feature is not supported on the System"
    }
}
catch {
    Write-Error $_.Exception.message
}
