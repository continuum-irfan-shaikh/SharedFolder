 <#
    .Synopsis
        Retrieve the number of PST files with path.
    .Author
        narayan.gouda@continuum.net
    .Name 
        Retrieve count of PST files
#>

$output = @()
try{
      $PSTs = Get-WmiObject -Query "Select * from CIM_DataFile Where Extension = 'pst'" -ErrorAction stop
}catch{
     Write-Error "Error occured while searching PST files!!"
     Write-Error $_.Exception.Message

}
if (-not $PSTs){
   Write-Error "No PST files found on this computer..!"
   Exit
  
}Else {
        $numberoffiles = ($PSTs | Measure-Object).count
        Foreach ( $pst in $PSTs){
        $output += New-Object PSObject -Property @{
                    "Name" = (Get-Item $pst.name).Name
                    "Path" = $pst.name
                    "Size in KB" = [math]::Round($((Get-Item $pst.name).Length / 1024),2)
                   }

      }
        Write-Output "Total number of PST files = $numberoffiles"
        Write-Output "----------------------------------"
        $output | select Name, Path, "Size in KB" | fl
}
