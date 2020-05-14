<#
      .Script Name
      Configure Firewall Exceptions.
      .Author
      Nirav Sachora
      .Requirements
      Script should run with highest privileges.
#>

<#$task = "Create"   #Create/Delete.
$direction = "Inbound" #Inbound/outbound  this will required for both action.if deletion Inbound/outbound/All
$ruletype = "Port" #Program/Port.
$programpath = "C:\Users" #this will be radio button All/Specify path of .exe file.
$action = "Allow"  #Allow/Block
$profiles = "Domain","Public","Private"#,"Private" #All/Domain/Private/Public.
$protocol = "UDP" #TCP/UDP.    For Deletion TCP/UDP/All
$portnumbers = "6559987" #All/specific ports.
$taskname = "Test12"  # Name of the rule.
$remotescope = "10.2.2.2,10.2.3.4,10.2.2.89-10.2.2.95"
$localscope = "10.2.2.6,10.2.3.4,10.2.2.89-10.2.2.95"
$description # Description of the rule, this will be optional.#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

if($remotescope){
    $ErrorActionPreference = "Stop"
$remoteaddress = $remotescope -split ","
foreach ($ip in $remoteaddress) {
    if ($ip -match "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])-(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$") {    
        try {
            [ipaddress]$first, [ipaddress]$second = $ip -split "-"
                $IPAddresses = @(

                    [System.Net.IPAddress]$first

                    [System.Net.IPAddress]$second

                )
                $sorted = $IPAddresses | Sort | Select IPAddressToString
            if($sorted[0].IPAddressToString -eq $first.IPAddressToString){
                
                Continue;
            }
            else{
                Write-Error "Invalid IP Address range $ip"
            }
        }
        catch {
            Write-Error "Invalid IP Address range $ip"
        }
    }
    elseif ($ip -match "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$") {
        try {
            [ipaddress]$validateip = $ip
            Continue;
        } 
        catch {
            Write-Error "Invalid IP Address $ip"
        }
    }
    else{
    Write-Error "Invalid ip address $ip"
    }
}
$ErrorActionPreference = "Continue"
}

if($localscope){
    $ErrorActionPreference = "Stop"
$localaddress = $localscope -split ","
foreach ($ip in $localaddress) {
    if ($ip -match "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])-(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$") {    
        try {
            [ipaddress]$first, [ipaddress]$second = $ip -split "-"
                $IPAddresses = @(

                    [System.Net.IPAddress]$first

                    [System.Net.IPAddress]$second

                )
                $sorted = $IPAddresses | Sort | Select IPAddressToString
            if($sorted[0].IPAddressToString -eq $first.IPAddressToString){
                
                Continue;
            }
            else{
                Write-Error "Invalid IP Address range $ip"
            }
        }
        catch {
            Write-Error "Invalid IP Address range $ip"
        }
    }
    elseif ($ip -match "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$") {
        try {
            [ipaddress]$validateip = $ip
            Continue;
        } 
        catch {
            Write-Error "Invalid IP Address $ip"
        }
    }
    else{
    Write-Error "Invalid ip address $ip"
    }
}
$ErrorActionPreference = "Continue"
}


function Validate_ports($ports) {

    $validateport = $ports -split ","
    foreach ($port in $validateport) {
        if ($port -like "*-*") {
            [int]$range1, [int]$range2 = $port -split "-"
            if (($range1 -le 65535) -and ($range2 -le 65535) -and ("$port" -match "^[0-9]{1,5}\-[0-9]{1,5}$")) {
                Continue;
            }
            else {
                Write-error "Given port number is not valid port, Please enter valid port number."
                Exit;
            }
        }
        else {
            
            if (([int]$port -le 65535)) {
                Continue;    
            }
            else {
                Write-error "Given port number is not valid port, Please enter valid port number."
                Exit;
            }
        }

    }
    return $true
}

function set_profiles {
    $length = $profiles.length
    if ($length -eq 3) {
        return 2147483647
    }
    elseif ($length -eq 2) {
        if (($profiles -contains "Domain") -and ($profiles -contains "Private")) { return 3 }
        if (($profiles -contains "Domain") -and ($profiles -contains "Public")) { return 5 }
        if (($profiles -contains "Public") -and ($profiles -contains "Private")) { return 6 }
    }
    else {
        switch ($profiles) {
            "Domain" { return 1 }
            "Private" { return 2 }
            "Public" { return 4 }
        }
    }
}
function Add_Program {
    $firewallrule = New-Object -ComObject HNetCfg.FWRule
    $firewallrule.Name = $taskname
    $firewallrule.ApplicationName = $programpath
    if (((Test-path $programpath) -eq $false) -or (((Get-Item "$programpath") -is [System.IO.DirectoryInfo]) -eq $true)) { Write-Error "Program path provided is incorrect, please provide valid file path."; Exit; }
    switch ($direction) {
        #switch statement for Inbound and outbound direction
        "Inbound" { $firewallrule.Direction = 1; break }
        "Outbound" { $firewallrule.Direction = 2; break }
    }
    switch ($action) {
        "Allow" { $firewallrule.Action = 1; break }
        "Block" { $firewallrule.Action = 0; break }
    }
    
    $setprofile = set_profiles
    $firewallrule.Profiles = $setprofile
    if($remotescope){
    $firewallrule.RemoteAddresses = $remotescope
    }
    if($localscope){
    $firewallrule.LocalAddresses = $localscope
    }
    $firewallrule.Enabled = 1
    return $firewallrule
}

function Add_port {
    $firewallrule = New-Object -ComObject HNetCfg.FWRule
    $firewallrule.Name = $taskname
    switch ($protocol) {
        "TCP" { $firewallrule.Protocol = 6; break }
        "UDP" { $firewallrule.Protocol = 17; break }
    }
    switch ($direction) {
        "Inbound" {
            $firewallrule.Direction = 1;
            if ($portnumbers -ne "") {
                $verifyports = Validate_ports -ports $portnumbers
                if($verifyports){
                $firewallrule.LocalPorts = $portnumbers
               } 
            }
            break;
        }
        "Outbound" {
            $firewallrule.Direction = 2;
            if ($portnumbers -ne ""){
                $verifyports = Validate_ports -ports $portnumbers
                if($verifyports){
                $firewallrule.RemotePorts = $portnumbers
                }
            } 
            break
        }
    }
    
    $setprofile = set_profiles
    $firewallrule.Profiles = $setprofile
    
    switch ($action) {
        "Allow" { $firewallrule.Action = 1; break }
        "Block" { $firewallrule.Action = 0; break }
    }
    if($remotescope){
    $firewallrule.RemoteAddresses = $remotescope
    }
    if($localscope){
    $firewallrule.LocalAddresses = $localscope
    }
    $firewallrule.Enabled = 1
    return $firewallrule       
}

function delete_firewallrule($deleterulefor) {
    switch ($direction) {
        "Inbound" { $value = 1 }
        "Outbound" { $value = 2 }
    }
    switch ($deleterulefor) {
        "Program" {
            $fwPolicy = New-Object -ComObject HNetCfg.FwPolicy2
            $rule = $fwpolicy.Rules | ? { ($_.ApplicationName -eq $programpath) -and ($_.Direction -eq $value) } | Select -ExpandProperty Name
            netsh advfirewall firewall delete rule name=$rule | Out-Null
            if ($LASTEXITCODE -eq 0) {
                return 0
            }
        }
        "Port" {
            $fwPolicy = New-Object -ComObject HNetCfg.FwPolicy2
            switch ($protocol) {
                "TCP" { $portvalue = 6 }
                "UDP" { $portvalue = 17 }
            }
            switch ($value) {
                1 {
                    if(($portnumbers.gettype()).Name -ne "String"){
                        $verifyport = $portnumbers -join ","
                    }
                    else{
                        $verifyport = $portnumbers
                    }
                    $rule = $fwpolicy.Rules | ? { ($_.Name -eq $taskname) -and ($_.Direction -eq $value) -and  ($_.Protocol -eq $portvalue) -and ($_.LocalPorts -eq $verifyport) } | Select -ExpandProperty Name

                    if (!$rule) {
                        $rule = $fwpolicy.Rules | ? { ($_.Name -eq $taskname) -and ($_.Direction -eq $value) -and ($_.Protocol -eq $portvalue) }; if (!$rule) { Write-output "Rule not found"; Exit }
                        $updateport = @()
                        if ($rule.LocalPorts -match ",") {
                            $temp = $rule.LocalPorts -split ","
                        }
                        if ($temp -contains $verifyport) {
                            foreach ($localport in $temp) {
                                if ($localport -ne $verifyport) {
                                    $updateport += $localport
                                }
                            }
                            $updateport = $updateport -join ","
                            $rule.localports = $updateport
                            return 1
                        }
                    }
                }
                2 { 
                    if(($portnumbers.gettype()).Name -ne "String"){
                        $verifyport = $portnumbers -join ","
                    }
                    else{
                        $verifyport = $portnumbers
                    }
                    $rule = $fwpolicy.Rules | ? { ($_.Name -eq $taskname) -and ($_.Direction -eq $value) -and  ($_.Protocol -eq $portvalue) -and ($_.RemotePorts -eq $verifyport) } | Select -ExpandProperty Name
                    
                    if (!$rule) {
                        $rule = $fwpolicy.Rules | ? { ($_.Name -eq $taskname) -and ($_.Direction -eq $value) -and ($_.Protocol -eq $portvalue) }; if (!$rule) { Write-output "Rule not found"; Exit }
                        $updateport = @()
                        if ($rule.RemotePorts -match ",") {
                            $temp = $rule.RemotePorts -split ","
                        }
                        if ($temp -contains $verifyport) {
                            foreach ($remoteport in $temp) {
                                if ($remoteport -ne $verifyport) {
                                    $updateport += $remoteport
                                }
                            }
                            $updateport = $updateport -join ","
                            $rule.remoteports = $updateport
                            return 2
                        }
                    }
                }
            }        
            netsh advfirewall firewall delete rule name=$rule | Out-null
            if($LASTEXITCODE -eq 0){return 0}
            else{return -1}
        }
    }
}

    if (($ruletype -eq "Program") -and ($task -eq "Create")) {
        $fwPolicy = New-Object -ComObject HNetCfg.FwPolicy2
        $fwPolicy.Rules.Add((Add_Program))
        if ($?) {
            Write-output "Rule has been added successfully for program."
            Exit;
        }
    }
    elseif (($ruletype -eq "Port") -and ($task -eq "Create")) {
        $fwPolicy = New-Object -ComObject HNetCfg.FwPolicy2
        $fwPolicy.Rules.Add((Add_port))
        if ($?) {
            Write-output "Rule has been added successfully for port."
            Exit;
        }
    }


$ErrorActionPreference = "Continue"
if ($task -eq "Delete") {
    if (!$portnumbers) {
        switch ($direction) {
            "Inbound" { $value = 1 }
            "Outbound" { $value = 2 }
        }
        switch ($protocol) {
            "TCP" { $portvalue = 6 }
            "UDP" { $portvalue = 17 }
        }
        $fwPolicy = New-Object -ComObject HNetCfg.FwPolicy2    
        $rule1 = $fwpolicy.Rules | ? { ($_.Name -eq $taskname) -and ($_.Direction -eq $value) -and ($_.Protocol -eq $portvalue) }
 
        if ($rule1) {
            $rulename = $rule1.Name
            netsh advfirewall firewall delete rule name=$rulename | Out-null
            if($LASTEXITCODE -eq 0){Write-output "Rule has been deleted successfully"}
            Else{Write-Error "Failed to delete rule."}
            Exit;
        }
    }
    $deleteresult = delete_firewallrule -deleterulefor $ruletype
    switch ($deleteresult) {
        0 { Write-output "Rule has been deleted successfully" }
        1 { Write-output "Local Port has been removed" }
        2 { Write-output "Remote Port has been removed" }
        Default { Write-output "Rule not found" }
    }
}



