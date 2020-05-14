<#
EventType:

BEGIN_NESTED_SYSTEM_CHANGE = 102
BEGIN_SYSTEM_CHANGE = 100
END_NESTED_SYSTEM_CHANGE = 103
END_SYSTEM_CHANGE = 101

RestorePointType:

APPLICATION_INSTALL = 0
APPLICATION_UNINSTALL = 1
CANCELLED_OPERATION = 13
DEVICE_DRIVER_INSTALL = 10
MODIFY_SETTINGS = 12
#>

$OS = (Get-WMIObject -Class Win32_OperatingSystem).Caption

If ($OS -like "*Server*") {
           Write-Error "Current OS : $OS. This functionality is not supported on this operating system."
           Exit
} Else {
     try {  
           Enable-ComputerRestore -Drive (Get-WmiObject Win32_OperatingSystem).SystemDrive
           Checkpoint-Computer -Description "Scripted restore" -RestorePointType $restoreType -WA stop -EA Stop
           $restorepoints = Get-ComputerRestorePoint -EA stop | Sort-Object CreationTime -Descending
     }catch{
           Write-Error $_.Exception.Message
           Exit
     }

}
$result = @()
$new_restore_point = ($restorepoints[0]).SequenceNumber
Write-Output "System restore point (SequenceNumber = $new_restore_point) successfully created."

foreach ($item in $restorepoints) {
     if ($item.EventType -eq 102 ) { $EventType = "BEGIN_NESTED_SYSTEM_CHANGE" }
         elseif ($item.EventType -eq 100 ) { $EventType = "BEGIN_SYSTEM_CHANGE" }    
         elseif ($item.EventType -eq 103 ) { $EventType = "END_NESTED_SYSTEM_CHANGE" }
         elseif ($item.EventType -eq 101 ) { $EventType = "END_SYSTEM_CHANGE" }
     else { $EventType = $item.EventType }
     
     if ($item.RestorePointType -eq 0) { $RestorePointType = "APPLICATION_INSTALL" }
         elseif ($item.RestorePointType -eq 1) { $RestorePointType = "APPLICATION_UNINSTALL" }
         elseif ($item.RestorePointType -eq 13) { $RestorePointType = "CANCELLED_OPERATION" }
         elseif ($item.RestorePointType -eq 10) { $RestorePointType = "DEVICE_DRIVER_INSTALL" }
         elseif ($item.RestorePointType -eq 12) { $RestorePointType = "MODIFY_SETTINGS" }
     else { $RestorePointType = $item.RestorePointType }
    
     $result += New-Object PSObject -Property @{
                    "CreationTime" = $item.ConvertToDateTime($item.CreationTime)
                    "Description" = $item.Description
                    "SequenceNumber" = $item.SequenceNumber
                    "EventType" = $EventType
                    "RestorePointType" = $RestorePointType 
     }
}
$result | Format-List CreationTime, Description, SequenceNumber, EventType, RestorePointType

