$regDrive, $regPathTail = $registryPath.split(':')

if (!(Get-PSDrive -name $regDrive -ErrorAction SilentlyContinue)) {
	switch($regDrive) {
		"HKCR" {
			New-PSDrive -Name HKCR -PSProvider Registry -Root HKEY_CLASSES_ROOT >$null
		}
		"HKCU" {
			New-PSDrive -Name HKCU -PSProvider Registry -Root HKEY_CURRENT_USER >$null
		}
		"HKLM" {
			New-PSDrive -Name HKLM -PSProvider Registry -Root HKEY_LOCAL_MACHINE >$null
		}
		"HKU" {
			New-PSDrive -Name HKU -PSProvider Registry -Root HKEY_USERS >$null
		}
		"HKCC" {
			New-PSDrive -Name HKCC -PSProvider Registry -Root HKEY_CURRENT_CONFIG >$null
		}
		default {
			Write-Error "Invalid Registry drive: $regDrive. Should be one of: HKCR, HKCU, HKLM, HKU or HKCC"
			return
		}
	}
}

function isSettingExists {
	if (Test-Path $registryPath -PathType container) {
		$key = Get-Item -LiteralPath $registryPath
        if ($key.GetValue($name, $null) -ne $null) {
			$true
        }
	}
	return $false
}

switch ($action) {
    "create" {
		if (!(Test-Path $registryPath -PathType container)) {
			New-Item -Path $registryPath -Force:$force >$null
		}
        if (!(isSettingExists)) {
            New-ItemProperty -Path $registryPath -Name $name -Value $value -PropertyType $type -Force:$force >$null
        } else {
            Write-Error "Invalid action: $action. RegistrySetting already exists. The 'update' action should be used."
            return
        }
    }
    "update" {
        if (isSettingExists) {
            Set-ItemProperty -Path $registryPath -Name $name -Value $value -Force:$force
        } else {
            Write-Error "Invalid action: $action. RegistrySetting doesn't exist. The 'create' action should be used."
            return
        }
    }
    "delete" {
        if (isSettingExists) {
            Remove-ItemProperty -Path $registryPath -Name $Name -Force:$force
        } else {
            Write-Error "Invalid action: $action. RegistrySetting doesn't exist."
            return
        }
    }
    default {
        Write-Error "Invalid action: $action. Should be one of: create, update, delete"
        return
    }
}

IF ($?){
return "
	action: $action,
	hive: $registryPath,
	key: $Name,
	regtype: $type,
	regvalue: $value,
	force: $false
	"
}
