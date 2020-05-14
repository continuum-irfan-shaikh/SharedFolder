<#
    .SYNOPSIS
       Retrieve Windows firewall properties
    .DESCRIPTION
       Retrieve Windows firewall properties
    .Author
       sushma yerasi  
#>

$computer = $env:COMPUTERNAME
$ExecutionLog=@()
$executionlog += "Computer Name : $computer"
[double]$OSVersion=[Environment]::OSVersion.Version.ToString(2)

# OS Version and Powershell comparison.
if (($osversion -lt '6.1') -or ($PSVersionTable.PSVersion.Major -lt '2'))
        {
            $executionlog += 'Prerequisites to run the script is not valid, Hence Script Exceution stopped' #'Script is design for windows 7 and Above Members only, Script Execution Stopped.'
            Write-Output $executionlog
            exit;
        }  

try {	
		if((get-service mpssvc).Status -eq 'Running')
		{ 
            $fwmgr=New-Object -Com HNetCfg.FwMgr
			$cp=$fwmgr.LocalPolicy.CurrentProfile
            if ($fwmgr.LocalPolicy.CurrentProfile.Type -eq 0){ $type = "Private"  }
            elseif($fwmgr.LocalPolicy.CurrentProfile.Type -eq 1){ $type = "Public" }
            elseif ($fwmgr.LocalPolicy.CurrentProfile.Type -eq 2){ $type = "Domain" }
            else { $type = "Unknown" }
			$properties = $cp | select @{N="Type";E={$type}},FirewallEnabled, ExceptionsNotAllowed, NotificationsDisabled | fl
			$executionlog += $properties 
			Write-Output $ExecutionLog
            Exit;
		}
	}
catch 			
	{ 
		$executionlog += "Windows Firewall Service is not running on the server"
		Write-Output $ExecutionLog
        Exit
	}
