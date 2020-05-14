<#
    .SYNOPSIS
        Enable Continuum Port
    .DESCRIPTION
        Enable Continuum port 443 for all TCP protocol and all user profile. Direction is outbound.
    .Help
        Use HNetCfg to get/add configuration of firewall.  
       
        #Default windows contstants for firewall policy
        $NET_FW_PROFILE2_DOMAIN = 1
        $NET_FW_PROFILE2_PRIVATE = 2
        $NET_FW_PROFILE2_PUBLIC = 4
        $NET_FW_IP_PROTOCOL_UDP = 17
        $NET_FW_IP_PROTOCOL_ICMPv4 = 1
        $NET_FW_IP_PROTOCOL_ICMPv6 = 58
        $NET_FW_PROFILE2_ALL = 2147483647
        $NET_FW_IP_PROTOCOL_TCP = 6
        $NET_FW_ACTION_ALLOW = 1
        $NET_FW_RULE_DIR_OUT = 2
        $NET_FW_RULE_DIR_IN = 1
        $NET_FW_ACTION_BLOCK = 0
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#> 
$NET_FW_IP_PROTOCOL_TCP = 6
$NET_FW_PROFILE2_ALL = 2147483647
$NET_FW_ACTION_ALLOW = 1
$NET_FW_RULE_DIR_OUT = 2
$PortNumber = 443

function Add-FirewallRule
{
$fwPolicy = New-Object -ComObject HNetCfg.FwPolicy2

$rule = New-Object -ComObject HNetCfg.FWRule
$rule.Name = 'Continuum Port'
$rule.Profiles = $NET_FW_PROFILE2_ALL
$rule.Enabled = $true
$rule.Action = $NET_FW_ACTION_ALLOW
$rule.Direction = $NET_FW_RULE_DIR_OUT
$rule.Protocol = $NET_FW_IP_PROTOCOL_TCP
$rule.RemotePorts = $PortNumber

$fwPolicy.Rules.Add($rule)

if ((New-Object -comObject HNetCfg.FwPolicy2).rules | ?{$_.Direction -eq 2 -and $_.RemotePorts -eq 443 -and $_.Profiles -eq 2147483647 -and $_.protocol -eq 6 -and $_.Action -eq 1 -and $_.Name -eq "Continuum Port" })
{
return $true
}
else {
return $false
}
}

function Test-Port($PortNumber)
{
    $rulescheck = (New-Object -comObject HNetCfg.FwPolicy2).rules | ?{$_.Direction -eq 2 -and $_.RemotePorts -eq 443 -and $_.Profiles -eq 2147483647 -and $_.protocol -eq 6 -and $_.Action -eq 1 -and $_.Name -eq "Continuum Port" }
    
    if($rulescheck -ne $null)
    {
        $msg = "`nPort $PortNumber is already open on system $env:computername"
      
    }
    else
    {
        if (Add-FirewallRule)
        {
        $msg = "`nPort $PortNumber on system $env:computername is opened now "                                
        }
        else
        {
        $msg = "`nThere is some issue while enabliing the port on system $env:computername. Kindly check manually."
        }
    }
    
    Write-Output $msg
}

Test-Port -PortNumber $PortNumber
