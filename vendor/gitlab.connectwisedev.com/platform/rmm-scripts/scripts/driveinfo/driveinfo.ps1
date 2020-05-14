Get-WmiObject Win32_DiskDrive | % {
  $disk = $_
  $partitions = "ASSOCIATORS OF " +
                "{Win32_DiskDrive.DeviceID='$($disk.DeviceID)'} " +
                "WHERE AssocClass = Win32_DiskDriveToDiskPartition"
  Get-WmiObject -Query $partitions | % {
    $partition = $_
    $drives = "ASSOCIATORS OF " +
              "{Win32_DiskPartition.DeviceID='$($partition.DeviceID)'} " +
              "WHERE AssocClass = Win32_LogicalDiskToPartition"
    Get-WmiObject -Query $drives | % {
      New-Object -Type PSCustomObject -Property  @{
      
            'Disk Model'   = $disk.Model
            'Disk Serial Number' = $_.volumeserialnumber
            'Disk Size'    = [String]([Math]::Round($disk.size/1GB, 2)) + ' GB'
            'Disk'= $disk.DeviceID
            'Drive Letter' = $_.DeviceID
            'FreeSpace'   = [String]([Math]::Round($_.FreeSpace/1GB, 2)) + ' GB'
            'Partition'   = $partition.Name
            'Raw Size'     = [String]([Math]::Round($partition.size/1GB, 2)) + ' GB'
            'Volume Name'  = $_.VolumeName
            'Volume Size'= [String]([Math]::Round($_.Size/1GB, 2)) + ' GB'
            
      } |Select-object 
    }
  }
}
