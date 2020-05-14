Try{

$startDate = (Get-Date) - (New-TimeSpan -Day 2)
$UserLoginTypes = 2,3,7,10,11 
$LastLogonUser = Get-WinEvent  -FilterHashtable @{Logname='Security';ID=4624;StartTime=$startDate} `
	| Where-Object {-not(($_.Properties[4].Value -like  "S-1-5-18") -or ($_.Properties[4].Value -like  "S-1-5-19") 	-or ($_.Properties[4].Value -like  "S-1-5-20")`
     -or ($_.properties[8].value -like '3') -or ($_.properties[8].value -like '4') -or ($_.properties[8].value -like '5') -or ($_.properties[8].value -like '8') `
      -or ($_.properties[8].value -like '9') -or ($_.properties[8].value -like '11'))}  | SELECT TimeCreated, @{N='Username'; E={$_.Properties[5].Value} }, 
	@{N='Domain/Machine'; E={$_.Properties[6].Value} },@{N='SID'; E={$_.Properties[4].Value}}, @{N='LogonType'; E={$_.Properties[8].Value}}, 
	@{N='IP Address'; E={$_.Properties[18].Value}} | WHERE {$UserLoginTypes -contains $_.LogonType}  | Sort-Object TimeCreated | Select -last 1
}
Catch{
        Write-Error "Information is not available at the moment : $($_.Exception.Message)"
        Exit
     }

if (!$allusers){
    $LogOnUser =  $lastlogonuser.Username
    If (!$LogOnUser){
        Write-Error "User should be logged on to create private shortcut"
        return
    }
} else {
    if ("desktop", "favorites" -notcontains $location){
        Write-Error "Wrong location for all users shortcut"
        return
    }
} 

Write-Host $LogOnUser

switch ($location){
    "desktop"{
        if (!$allusers) {
            $DestinationPath = "$env:homedrive\Users\$logonuser\Desktop\$name"
        } else {
             $DestinationPath = "$env:public\Desktop\$name"
        }
    }
    "favorites"{
        $DestinationPath = "$env:homedrive\Users\Default\Favorites\$name"
        if (!$allusers) {
            $DestinationPath = "$env:homedrive\Users\$logonuser\Favorites\$name"
        }
    }
    "menu"{  
        $DestinationPath = "$env:homedrive\Users\$logonuser\AppData\Roaming\Microsoft\Windows\Start Menu\$name"
    }
    "programs"{
        $DestinationPath = "$env:homedrive\Users\$logonuser\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\$name"
    }
    "startup"{
        $DestinationPath = "$env:homedrive\Users\$logonuser\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\StartUp\$name"
    }
    "sendto"{
        $DestinationPath = "$env:homedrive\Users\$logonuser\AppData\Roaming\Microsoft\Windows\SendTo\$name"
    }
    "quickbar"{
        $DestinationPath = "$env:homedrive\Users\$logonuser\AppData\Roaming\Microsoft\Internet Explorer\Quick Launch\$name"
    }
    default {
        Write-Error "Unsupported shortcut location"
        return
    }
}

if ($type -eq "shell"){
	$DestinationPath += ".lnk"
} else {
	$DestinationPath += ".url"
}

switch ($action){
	"create" {
		if (!$overwrite -and (Test-Path $DestinationPath)){
			Write-Error "Shortcut exists. Delete existing or check overwrite option"
			return
		}
		$WshShell = New-Object -comObject WScript.Shell
		$Shortcut = $WshShell.CreateShortcut($DestinationPath)
		$Shortcut.TargetPath = "$target"
		if ($arguments -ne $null){
			$Shortcut.Arguments = $arguments
		}
		if ($startInFolder -ne $null){
			$Shortcut.WorkingDirectory = $startInFolder
		}
		if ($comments -ne $null){
			$Shortcut.Description = $comments
		}
		$Shortcut.Save()
        if ($?){
            Write-Output "Shortcut was successfully created"
        }Else {"Unable to Save the Shortcut"}
	}
	"delete" {
		Remove-Item $DestinationPath -Force
        IF ($?){
            Write-Output "Shortcut was successfully deleted"
        }
	}
	default {
		Write-Error "Unsupported action $action, should be one of create, delete"
	}
}
