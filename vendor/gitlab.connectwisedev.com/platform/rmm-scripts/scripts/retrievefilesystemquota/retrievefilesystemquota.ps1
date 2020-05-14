$drives = Get-WmiObject -Class Win32_LogicalDisk | ? {$_.Description -eq 'Local Fixed Disk'}
if ($drives.Count -eq $null) {$count = 1}else {$count = $drives.count}
"`nNumber of volumes on this system : {0}`n" -f $count

try {
    $volumes = Get-WmiObject -Class Win32_QuotaSetting | select State, VolumePath, Caption, DefaultLimit, DefaultWarningLimit, ExceededNotification, WarningExceededNotification -ErrorAction Stop

    if (!$volumes) {
        Write-Output "Failed to retrieve quota information, make sure you are running script with highest privileges. `n"
    }

    foreach ($volume in $volumes) {
        if ($volume.State -eq 0) {
            $state = "False"
        }
        else {
            $state = "True"
        }
        "`nDrive Path             :  {0}" -f $volume.VolumePath
        "Drive Name             :  {0}" -f $volume.Caption
        "Quota Enabled          :  {0}" -f $state
        ######## Quota Limit ###### 
        if ($volume.DefaultLimit -eq -1) {
            $Defaultlimit = "No Limit"
            "Quota Limit            :  {0}" -f $Defaultlimit
        }
        else {
            if (($volume.Defaultlimit / 1GB) -ge 1) {    
                "Quota Limit            :  {0} {1}" -f ("{0:N2}" -f ($volume.Defaultlimit / 1GB)), "GB"
            }
            else {
                "Quota Limit            :  {0} {1}" -f ("{0:N2}" -f ($volume.Defaultlimit / 1MB)), "MB"
            }
        }
   
        ################################
        ######## Warning Limit #########
        if ($volume.DefaultWarningLimit -eq -1) {
            $Defaultwarninglimit = "No Limit"
            "Warning Limit          :  {0}" -f $Defaultwarninglimit
        }
        else {
            if (($volume.Defaultwarninglimit / 1GB) -ge 1) {
                "Warning Limit          :  {0} {1}" -f ("{0:N2}" -f ($volume.Defaultwarninglimit / 1GB)), "GB"
            }
            else {
                "Warning Limit          :  {0} {1}" -f ("{0:N2}" -f ($volume.Defaultwarninglimit / 1MB)), "MB"
            }
        }
    
        ################################
        "Log event when user exceeds their quota limit   :  {0}" -f $volume.ExceededNotification
        "Log event when user exceeds their warning level :  {0}`n" -f $volume.WarningExceededNotification
    }

}
catch {
    "`n"
    Write-Error "Failed To Retrieve Quota Information"
}

