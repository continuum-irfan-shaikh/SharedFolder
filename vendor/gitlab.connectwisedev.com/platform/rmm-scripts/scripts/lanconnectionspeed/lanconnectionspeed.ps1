$Adapters = Get-WmiObject -Class Win32_NetworkAdapter |
Where { $_.Speed -ne $null -and $_.MACAddress -ne $null -and $_.AdapterType -match 'Ethernet'}
if ($Adapters -eq $null){
    Return "LAN Adapters were not founded"
} Else {
    ForEach ($Adapter in $Adapters){
        If ($Adapter.NetConnectionStatus -eq 2){
            If ($Adapter.Speed -ge 1000000000){
                $Speed = [math]::Round(($Adapter.Speed/1000000000),1)
                $SpeedType = 'gbps'
            } ElseIf ($Adapter.Speed -ge 1000000){
                $Speed = [math]::Round(($Adapter.Speed/1000000),1)
                $SpeedType = 'mbps'
            } Else {
                $Speed = [math]::Round(($Adapter.Speed/1000),1)
                $SpeedType = 'kbps'
            }
            $Adapter | Add-Member -MemberType NoteProperty -Name SpeedType -Value "$Speed $SpeedType"
        } Else {
            $Adapter | Add-Member -MemberType NoteProperty -Name SpeedType -Value "0 kbps"
        }
    }
}
$Adapters | Format-Table -Property Name,@{Label="Speed"; Expression = {($_.SpeedType)}}
