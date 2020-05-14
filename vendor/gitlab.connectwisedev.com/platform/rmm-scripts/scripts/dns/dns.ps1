$dnsServers = Get-WmiObject Win32_NetworkAdapterConfiguration | % { $_.DNSServerSearchOrder }
write-output $dnsServers | select @{N="DNS Servers";E={$_}}

