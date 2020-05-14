# $Method = "Close Chrome"

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}
if (!$Method) { 
    Write-output "Update option not selected. Please select a task option before executing again."
    Exit
}

$Scriptblock = {
    try {
        $AllowException = 'logmein.com'
        $ErrorActionPreference = 'Stop'
        $Profiles = Get-WmiObject win32_userprofile | Where-Object { $_.sid -like "S-1-5-21*" } | Select-Object -ExpandProperty localpath
        $ProgramFiles = if ($env:PROCESSOR_ARCHITECTURE -eq 'AMD64') { ${Env:ProgramFiles(x86)} } else { $Env:ProgramFiles }
        $DateTimeNow = (Get-Date).tostring('dd-MMM-yyyy_HH-mm-ss')
        $LogDir = "C:\Windows\Logs\LMICookieException\$DateTimeNow"
        $LogFile = Join-Path $LogDir "Log_$DateTimeNow.log"
        $ChromeMasterPreferenceBackupFile = Join-Path $LogDir "master_preferences_BACKUP_$DateTimeNow"
        $ChromeMasterPreferenceFile = "$ProgramFiles\Google\Chrome\Application\master_preferences"
        $Padding = ($Profiles | ForEach-Object { (Split-Path $_ -Leaf).Length } | Measure-Object -Maximum).Maximum
        $TextInfo = (Get-Culture).TextInfo
        Function Write-Log {
            Param(
                $Message,
                [Switch] $WriteOutput,
                [Switch] $AddNewLine
            ) 
        
            $Log = "[{0}] {1}" -f [datetime]::Now , $Message
            $Log | Out-File $LogFile -Append
            if ($WriteOutput) {
                Write-Output $Message
            }
        }
        function ConvertTo-Json([psobject] $item) {
            [System.Reflection.Assembly]::LoadWithPartialName("System.Web.Extensions") | out-null
            $ser = New-Object System.Web.Script.Serialization.JavaScriptSerializer 
            $hashed = @{ }
            $item.psobject.properties | ForEach-Object { $hashed.($_.Name) = $_.Value }
            write-output $ser.Serialize($hashed) 
        }
        function ConvertFrom-Json([string] $json) {
            [System.Reflection.Assembly]::LoadWithPartialName("System.Web.Extensions") | out-null
            $ser = New-Object System.Web.Script.Serialization.JavaScriptSerializer
            write-output (new-object -type PSObject -property $ser.DeserializeObject($json))
        }
        Function AddAllowCookieException {
            param(
                $Path,
                $URL
            )
        
            $LastModified = [int64](([datetime]::UtcNow) - (get-date "1/1/1970")).TotalMilliseconds
            $Obj = ConvertFrom-Json (Get-Content $Path)
        
            if (!$Obj.profile) {
                $Dictionary = New-Object 'System.Collections.Generic.Dictionary[string,long]'
                $Dictionary.add('last_modified', $LastModified)
                $Dictionary.add('setting', 1)
            
                $cookies = New-Object 'System.Collections.Generic.Dictionary[System.string,System.Object]'
                $cookies.add("$URL,*", $dictionary)
            
                $exceptions = New-Object 'System.Collections.Generic.Dictionary[System.string,System.Object]'
                $exceptions.add('cookies', $cookies)
            
                $content_settings = New-Object 'System.Collections.Generic.Dictionary[System.string,System.Object]'
                $content_settings.add('exceptions', $exceptions)
            
                $pref_profile = New-Object 'System.Collections.Generic.Dictionary[System.string,System.Object]'
                $pref_profile.add('content_settings', $content_settings)
            
                $Obj | Add-Member -MemberType NoteProperty -name "profile" -Value ([Hashtable]$pref_profile)
            }
            elseif (!$obj.profile.content_settings) {
                $Dictionary = New-Object 'System.Collections.Generic.Dictionary[string,long]'
                $Dictionary.add('last_modified', $LastModified)
                $Dictionary.add('setting', 1)
            
                $cookies = New-Object 'System.Collections.Generic.Dictionary[System.string,System.Object]'
                $cookies.add("$URL,*", $dictionary)
            
                $exceptions = New-Object 'System.Collections.Generic.Dictionary[System.string,System.Object]'
                $exceptions.add('cookies', $cookies)
            
                $content_settings = New-Object 'System.Collections.Generic.Dictionary[System.string,System.Object]'
                $content_settings.add('exceptions', $exceptions)
                
                $obj.profile.add('content_settings', $content_settings)
            }
            elseif (!$obj.profile.content_settings.exceptions) {
                $Dictionary = New-Object 'System.Collections.Generic.Dictionary[string,long]'
                $Dictionary.add('last_modified', $LastModified)
                $Dictionary.add('setting', 1)
            
                $cookies = New-Object 'System.Collections.Generic.Dictionary[System.string,System.Object]'
                $cookies.add("$URL,*", $dictionary)
            
                $exceptions = New-Object 'System.Collections.Generic.Dictionary[System.string,System.Object]'
                $exceptions.add('cookies', $cookies)
                
                $obj.profile.content_settings.add('exceptions', $exceptions)
            }
            elseif (!$obj.profile.content_settings.exceptions.cookies) {
                $Dictionary = New-Object 'System.Collections.Generic.Dictionary[string,long]'
                $Dictionary.add('last_modified', $LastModified)
                $Dictionary.add('setting', 1)
            
                $cookies = New-Object 'System.Collections.Generic.Dictionary[System.string,System.Object]'
                $cookies.add("$URL,*", $dictionary)
                
                $obj.profile.content_settings.exceptions.add('cookies', $cookies)
            }
            elseif (!($obj.profile.content_settings.exceptions.cookies."$URL,*")) {
                $Dictionary = New-Object 'System.Collections.Generic.Dictionary[string,long]'
                $Dictionary.add('last_modified', $LastModified)
                $Dictionary.add('setting', 1)
                $obj.profile.content_settings.exceptions.cookies.add("$URL,*", $dictionary)
            }
        
            return $Obj
        }
        function ValidateCookieException {
            param (
                $Path,
                $URL
            )
            return ((ConvertFrom-Json (Get-Content $Path)).profile.content_settings.exceptions).cookies."$URL,*"
        }
        
        if (![IO.Directory]::Exists($LogDir)) { [IO.Directory]::CreateDirectory($LogDir) | Out-Null } # create folder if that doesn't exists
        
        #region CHROME_SETTINGS
        Write-Log "[BEGIN: STEP-1] Loop and Add cookie exception ($AllowException) to all Chrome user profiles"
        Write-Output "`nAdding cookie exception to following Chrome user profiles:`n"
        # adding google chrome cookie exception on each user profile
        foreach ($Item in $Profiles) { 
            try {
                $Username = $TextInfo.ToTitleCase($(Split-Path $Item -Leaf))
                $ChromeUserPreferenceBackupFile = Join-Path $LogDir "$(Split-Path $Username -Leaf)_preferences_BACKUP_$DateTimeNow"
                $LocalAppData = "$Item\AppData\Local"
                $Path = "$LocalAppData\Google\Chrome\User Data\Default\preferences"
                Write-Log "[ATTEMPT] Add cookie exception to the user profile [$Item]"
                if (Test-Path $Path) {
                    Write-Log "[CONTINUE] Preferences file path found for the user [$Path]"
                    $LastModified = [int64](([datetime]::UtcNow) - (get-date "1/1/1970")).TotalMilliseconds
                    if (ValidateCookieException -Path $Path -URL $AllowException) {
                        Write-Log "[SKIP] Cookie exception in User Preferences already exists [$AllowException,*]"
                        Write-Output $("{0} : Already exists" -f $Username.PadRight($Padding, ' '))
                    }
                    else {  
                        $UserPreferencesObject = AddAllowCookieException -Path $Path -URL $AllowException
                        Write-Log "[INFO] Backup Chrome User Preferences File at location [$ChromeUserPreferenceBackupFile]"
                        Copy-Item $Path $ChromeUserPreferenceBackupFile
                        ConvertTo-Json $UserPreferencesObject | Out-File $Path -Encoding UTF8
        
                        $ValidateSetting = ValidateCookieException -Path $Path -URL $AllowException
                        if ($ValidateSetting.setting -eq 1) {
                            Write-Log "[SUCCESS] Add Cookie exception [$AllowException,*] on user profile [$Item]"
                            Write-Output $("{0} : Success" -f $Username.PadRight($Padding, ' '))
                        }
                        else {
                            Write-Log "[FAILURE] Add Cookie exception [$AllowException,*] on user profile [$Item]"
                            Write-Output $("{0} : Failure" -f $Username.PadRight($Padding, ' '))
                        }
                    }
                    
                }
                else {
                    Write-Log "[SKIP] Preference file not found [$Path]"
                    Write-Output $("{0} : Cookies configuration file not found. Please check may be Chrome is not installed." -f $Username.PadRight($Padding, ' '))
                }
            }
            catch {
                Write-Log "[ERROR] $($_.exception.message)"
                Write-Log "[FAILURE] Unable to add Cookie exception [$AllowException,*] on user profile [$Item]"
                Write-Output $("{0} : Failure" -f $Username.PadRight($Padding, ' '))
            }
            
        }
        
        Write-Log "[END: STEP-1] Loop and Add cookie exception ($AllowException) to all Chrome user profiles"
        
        # adding google chrome cookie exception on master preference settings for all new profiles that would be created
        
        Write-Log "[BEGIN: STEP-2] Add cookie exception ($AllowException) to Chrome's Master Preferences File"
        
        Try {
            if (Test-Path $ChromeMasterPreferenceFile) {
                if (ValidateCookieException -Path $ChromeMasterPreferenceFile -URL $AllowException) {
                    Write-Log "[SKIP] Cookie exception in 'Master Preferences' already exists [$AllowException,*]"
                    # Write-output "`nAdding cookie to Chrome's Master Preferences : Already exists"
                }
                else {
                    Write-Log "[ATTEMPT] Add Cookie exception in 'Master Preferences' [$AllowException,*]"
                    $ChromePreferencesObject = AddAllowCookieException -Path $ChromeMasterPreferenceFile -URL $AllowException
                    Write-Log "[INFO] Backup Chrome Master File at location [$ChromeMasterPreferenceBackupFile]"
                    Copy-Item "$ChromeMasterPreferenceFile" "$ChromeMasterPreferenceBackupFile"
                    Convertto-Json $ChromePreferencesObject | Out-File $ChromeMasterPreferenceFile -Encoding UTF8
        
                    # validate setting after updating it
                    $ValidateSetting = ValidateCookieException -Path $ChromeMasterPreferenceFile -URL $AllowException
                    if ($ValidateSetting.setting -eq 1) {
                        Write-Log "[SUCCESS] Cookie exception in 'Master Preferences' exists [$AllowException,*]"
                        # Write-output "`nAdding cookie to Chrome's Master Preferences : Success"
                    }
                    else {
                        Write-Log "[FAILURE] Cookie exception in 'Master Preferences' doesn't exists [$AllowException,*]"
                        # Write-output "`nAdding cookie to Chrome's Master Preferences : Failure"
                    }
                }
            }
            else {
                Write-Log "[SKIP] Preference file not found [$ChromeMasterPreferenceFile]"
                # Write-output "`nAdding cookie to Chrome's Master Preferences : Preferences file not found"
            }
        }
        catch {
            Write-Log "[ERROR] $($_.exception.message)"
            Write-Log "[FAILURE] Unable to add Cookie exception in 'Master Preferences' [$AllowException,*]"
            # Write-output "`nAdding cookie to Chrome's Master Preferences : Failure"
        }
        Write-Log "[END: STEP2] Add cookie exception ($AllowException) to Chrome's Master Preferences File"
        #endregion CHROME_SETTINGS
    }
    catch {
        # NOTE: intentionally avoid printing errors here    
    }
    finally {
        try {
            Start-Process -FilePath PowerShell.exe -ArgumentList 'schtasks.exe /Delete /TN LogMeInCookieExceptionOnRestart /F' -WindowStyle Hidden
        }
        catch {
            # NOTE: no action required, as this try..catch block is to suppress any errors due to non-existent task
        }
    }
}

$ExceptionIsSet = @()
$ExceptionIsNotSet = @()
$NoPreferenceFile = @()
$AllowException = 'logmein.com'
$ErrorActionPreference = 'Stop'
$Profiles = Get-WmiObject win32_userprofile | Where-Object { $_.sid -like "S-1-5-21*" } | Select-Object -ExpandProperty localpath
$ProgramFiles = if ($env:PROCESSOR_ARCHITECTURE -eq 'AMD64') { ${Env:ProgramFiles(x86)} } else { $Env:ProgramFiles }
$Padding = ($Profiles | ForEach-Object { (Split-Path $_ -Leaf).Length } | Measure-Object -Maximum).Maximum
$TextInfo = (Get-Culture).TextInfo
if (![IO.Directory]::Exists('C:\Windows\Logs\LMICookieException\')) { [IO.Directory]::CreateDirectory('C:\Windows\Logs\LMICookieException\') | Out-Null } # create folder if that doesn't exists
    
function ConvertFrom-Json([string] $json) {
    [System.Reflection.Assembly]::LoadWithPartialName("System.Web.Extensions") | out-null
    $ser = New-Object System.Web.Script.Serialization.JavaScriptSerializer
    write-output (new-object -type PSObject -property $ser.DeserializeObject($json))
}
    
function ValidateCookieException ($Path, $URL) {
    return ((ConvertFrom-Json (Get-Content $Path)).profile.content_settings.exceptions).cookies."$URL,*"
}


#if (Test-Path "$ProgramFiles\Google\Chrome\Application\chrome.exe") {
$output = foreach ($Item in $Profiles) { 
    try {
        $Username = $TextInfo.ToTitleCase($(Split-Path $Item -Leaf))
        $LocalAppData = "$Item\AppData\Local"
        $Path = "$LocalAppData\Google\Chrome\User Data\Default\preferences"
        if (Test-Path $Path) {
            if (ValidateCookieException -Path $Path -URL $AllowException) {
                Write-Output $("{0} : Already exists" -f $Username.PadRight($Padding, ' '))
                $ExceptionIsSet += $Username
            }
            else {  
                Write-Output $("{0} : Not Set" -f $Username.PadRight($Padding, ' '))
                $ExceptionIsNotSet += $Username
            }
        }
        else {
            Write-Output $("{0} : Cookies configuration file not found. Please check may be Chrome is not installed." -f $Username.PadRight($Padding, ' '))
            $ExceptionIsNotSet += $Username
        }
    }
    catch {
        Write-Output $("{0} : Failure" -f $Username.PadRight($Padding, ' '))
    }
}
    
if ($ExceptionIsNotSet.Count -ge 1) {
    # execute script only if
    Write-Output ("`nFound {0} out of {1} {2}" -f $ExceptionIsNotSet.count, ($ExceptionIsSet + $ExceptionIsNotSet).count, "user profiles that don't have the cookie exception.") 

    switch ($Method) {
        'System Reboot' {
            "`n"
            $output
            $ScriptFile = "C:\Windows\Logs\LMICookieException\Script.ps1"
            $Scriptblock.ToString() | Out-File $ScriptFile -Force
            $Task = "PowerShell.exe -executionpolicy bypass -NoExit -noprofile -WindowStyle Hidden -command '. $ScriptFile'"
            schtasks.exe /create /s $($env:COMPUTERNAME) /tn LogMeInCookieExceptionOnRestart /sc ONSTART /tr $Task /ru 'System' /F | Out-Null
            schtasks.exe /End /TN LogMeInCookieExceptionOnRestart | Out-Null
            
            # special condition to handle the case if - preference file is not found and only one profile exists on the system
            # then no need to schedule the script for next system reboot, and script can exit
            $Count = $output| Measure-Object |Select-Object -expand count
            if($Count -eq 1){
                if($Output -like "*Please check may be Chrome is not installed.*"){
                    break
                }
            }
            
            if (Invoke-Expression "Schtasks.exe /Query /TN LogMeInCookieExceptionOnRestart /v") {
                Write-Output "`nSuccessfully scheduled the script to add cookie exception at next system reboot."
            }
            else {
                Write-Output "`nFailed to schedule the script"
            }    
        }
        'Close Chrome' {
            Get-Process Chrome -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
            Start-Sleep -s 5
            if (!(Get-Process Chrome -ErrorAction SilentlyContinue)) {
                & $Scriptblock
            }
            else {
                Write-Output 'Chrome.exe is still running on the system. Failed to perform the operation.'
            }
        }
    }

}
else {
    # exit no action is required
    Write-Output ("`nSummary: {0} out of {1} {2}" -f $ExceptionIsSet.count, ($ExceptionIsSet + $ExceptionIsNotSet).count, "user profiles have the cookie exception, no action is required.") 
}
# }
# else {
#     Write-Output "`nChrome is not installed on the system. Hence the task is exited as no action is needed."
# } 
