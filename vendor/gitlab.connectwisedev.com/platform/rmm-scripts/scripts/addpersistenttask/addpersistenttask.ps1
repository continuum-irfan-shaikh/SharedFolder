# $Action = 'create' # 'delete'
# $TaskName = 'New Task'
# $Overwrite = $false
# $ApplicationName = "C:\Windows\explorer.exe"
# $Arguments = '/q /norestart /b'
# $RunAsUser = 'admin'
# #$Password = 'lic@123'
# $StartDateTime = "2019-05-14T19:29:00.000Z" 

###########################
# $Schedule = 'Daily'
# $Frequency = 'Every day'

# $Frequency = 'Week days'
# $days = 5

# $Frequency = 'Every'
# $days = 17
#############################
# $Schedule = 'Weekly'
# $Frequency = 'Every'
# $Weeks = 13
# $Weekdays = 'Monday'#,'Wednesday', 'Friday'
#############################
# $Schedule = 'Monthly'
# $Frequency = 'Day'
# $Day = 4
# $Months = 'January','April','July'

# $Schedule = 'Monthly'
# $Frequency = 'Week'
# $Week = 'third'
# $Weekday = 'wednesday'
# $Months = 'January','April','July','September'
#############################
# $Schedule = 'Once'
#############################
# $Schedule = 'At System Startup


#############################
# $Schedule = 'At Logon'
#############################
# $Schedule = 'When Idle'
# $Minutes = 50

try {
    if($StartDateTime){
        $datetime = [datetime] ($StartDateTime.TrimEnd(":.000Z"))
        $StartTime = $datetime.ToString('HH:mm')
        $StartDate = $datetime.ToString('MM/dd/yyyy')
    }

    $ErrorActionPreference = 'Stop'
    $AllWeekdays = 'MON', 'TUE', 'WED', 'THU', 'FRI'

    if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
        if ($myInvocation.Line) {
            &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
        }
        else {
            &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
        }
        exit $lastexitcode
    }

    Function DeleteTask {
        [CmdletBinding()]
        Param( [String] $TaskName )

        if (!(Get-Task $TaskName)) {
            Write-Error "Failed to delete:`'$TaskName`' because it doesn't exists."
        }
   
        if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64" -and $env:PROCESSOR_ARCHITECTURE -eq 'x86') { $EXE = 'c:\windows\sysnative\Schtasks.exe' }else { $EXE = 'c:\windows\System32\Schtasks.exe' }
        
        $Expression = "$EXE /Delete /TN `"$TaskName`" /F" 
        $result = Invoke-Expression $Expression
        if ($Result -like "SUCCESS*") { return $true }
        else { return $false }
    }

    Function CreateTask {
        [CmdletBinding()]
        Param( 
            [String] $TaskName,
            [String] $ApplicationName,
            [String] $Arguments,
            [String] $Username,
            [String] $Password,
            [String] $Trigger,
            [Boolean] $Overwrite
        )

        if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64" -and $env:PROCESSOR_ARCHITECTURE -eq 'x86') {
            $EXE = 'c:\windows\sysnative\Schtasks.exe'
        }
        else {
            $EXE = 'c:\windows\System32\Schtasks.exe'
        }
    
        $AppAndArgs = "`'$ApplicationName`' $Arguments".trim()
        $Expression = "$EXE /Create /TN `"$TaskName`" /TR `"$AppAndArgs`" /RU `"$RunAsUser`" $Trigger" 

        if ($Password) { $Expression = $Expression + " /RP `"$Password`"" }
        if ($Overwrite) { 
            $Expression = $Expression + " /F" 
        }
        else {
            if (Get-Task $TaskName) {
                Write-Error "Failed to create:`'$TaskName`' because a task with this name already exists.`nPlease use the `"Overwrite`" option and run the script again."
            }
        }

        $result = Invoke-Expression $Expression -ErrorAction Stop
        if ($Result -like "SUCCESS*") { return $true }
        else { return $false }
    }

    Function Get-Task ($TaskName) {
        try {
            $data = Invoke-Expression "Schtasks.exe /Query /TN `"$TaskName`" "
            if ($data) { return $true }
        }
        catch {
            return $false
        }

    }

    Switch ($Action) {
        'Create' {

            $DisableDomainCreds = Get-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Lsa" -Name DisableDomainCreds | Select-Object -ExpandProperty disabledomaincreds
 
            if ([bool]$DisableDomainCreds -and $Password) {
                Write-Output "`nFailed to create the task because following `"Local Security Policy`" [Computer Configuration > Windows Settings > Security Settings > Local Policies > Security Options] :`n`n'Network access: Do not allow storage of credentials or .NET Passports for network authentication' = 'Enabled'.`n`nPlease 'Disable' it and run the script again." 
            }
            else{
                
                $ScheduleHash = @{
                    'Daily'             = 'DAILY'
                    'Weekly'            = 'WEEKLY'
                    'Monthly'           = 'MONTHLY'
                    'Once'              = 'ONCE'
                    'At System Startup' = 'ONSTART'
                    'At Logon'          = 'ONLOGON'
                    'When Idle'         = 'ONIDLE'
                }
            
                # calculates the schedule switches and build the task trigger
                $SC = $ScheduleHash[$Schedule]
                $Trigger = Switch ($SC) {
                    'ONCE' {
                        "/SC $SC /st $StartTime /sd $StartDate"
                    }
                    'DAILY' {
                        Switch ($Frequency) {
                            'Every Day' { "/SC `"$SC`" /st `"$StartTime`" /sd `"$StartDate`"" }
                            'Week Days' { "/SC WEEKLY /MO 1 /d `"$($AllWeekdays -join ',')`" /st `"$StartTime`" /sd `"$StartDate`"" }
                            'Every' { "/SC `"$SC`" /MO `"$Days`" /st `"$StartTime`" /sd `"$StartDate`"" }
                        }
                    }
                    'WEEKLY' {
                        Switch ($Frequency) {
                            'Every' {
                                $Weekdays = $Weekdays | ForEach-Object { $_.SubString(0, 3).toUpper() }
                                "/SC `"$SC`" /MO `"$Weeks`" /d `"$($Weekdays -join ',')`" /st `"$StartTime`" /sd `"$StartDate`""
                            }
                        }
                    }
                
                    'MONTHLY' {
                        $Months = $Months | ForEach-Object { $_.SubString(0, 3).toUpper() }
                        Switch ($Frequency) {
                            'Day' {
                                if ($Day -and $Months) {
                                    "/SC `"$SC`" /d `"$Day`" /m `"$($Months -Join ',')`" /st `"$StartTime`" /sd `"$StartDate`""
                                }
                            }
                            'Week' {
                                if ($Week -and $Weekday -and $Months) {
                                    "/SC `"$SC`" /MO `"$week`" /d `"$($Weekday.SubString(0,3).toUpper())`" /m `"$($Months -Join ',')`" /st `"$StartTime`" /sd `"$StartDate`""
                                }
                            }
                        }
                    }
                    'ONSTART' {
                        "/SC `"$SC`""
                    }
                    'ONLOGON' {
                        "/SC `"$SC`""
                    }
                    'ONIDLE' {
                        "/SC `"$SC`" /I $Minutes"
                    }
                }
            
                if (CreateTask $TaskName $ApplicationName $Arguments $RunAsUser $Password $Trigger $Overwrite -ErrorAction Stop) {
                    Write-Output "`nSuccessfuly created task: `'$TaskName`'" 
                }
                else {
                    Write-Output "`nFailed to created task: `'$TaskName`'" 
                }
            }
        }
        'Delete' {
            if (DeleteTask $TaskName -ErrorAction Stop) {
                Write-Output "`nSuccessfuly deleted task: `'$TaskName`'" 
            }
            else {
                Write-Output "`nFailed to delete task: `'$TaskName`'" 
            }
        }
    }
}
catch {
    Write-Error $_.Exception.message
}
