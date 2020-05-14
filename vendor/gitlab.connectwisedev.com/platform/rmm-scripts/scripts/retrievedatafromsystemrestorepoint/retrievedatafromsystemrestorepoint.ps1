<#
    .SYNOPSIS
       Retrieve data from system restore point
    .DESCRIPTION
       Retrieve data from system restore point
    .Author
       Santosh.Dakolia@continuum.net    
#>
Try{
$computername= $env:computername

$OS = (Get-WMIObject -Class Win32_OperatingSystem).Caption
If ($OS -like "*Server*") {
           Write-Warning "Current OS : $OS"
           Write-Warning "This functionality is not supported on this operating system."
                           }
Else {
$RestoreData = Get-ComputerRestorePoint 
$output = @()

ForEach($RD in $RestoreData){
        $RP = Switch ($RD.RestorePointType){
                0 { "Application installation" }
                1 { "Application uninstall" }
                6 { "Restore" }
                7 { "Checkpoint" }
                10 { "Device drive installation" }
                11 { "First run" }
                12 { "Modify settings" }
                13 {"Cancelled operation" }
                14 { "Backup recovery" }
                Default { "Unknown"}
        }
$CT = $rd.ConvertToDateTime($rd.CreationTime)
     $output +=  New-Object psobject -Property @{
     "Description" = $RD.Description
     "Sequence Number" = $RD.SequenceNumber
     "Restore Point Type" = $RP
     "Creation Time" = $CT
     }
     }
$output |FL Description, "Sequence Number", "Restore Point Type", "Creation Time"
}
}Catch{
    Write-Error "Error occured while retrieving Data..!! $_.Exception.Message"
}
