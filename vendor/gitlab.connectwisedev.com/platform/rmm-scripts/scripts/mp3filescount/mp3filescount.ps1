 <#

 Name : Retrieve count of MP3 files
 Category : Data Collection


    .Synopsis
        Retrieve count of MP3 files.
    .Author
        Santosh.Dakolia@continuum.net
    .Name 
        Retrieve count of MP3 files
#>

$output = @()
try{
      $EXT = Get-WmiObject -Query "Select * from CIM_DataFile Where Extension = 'MP3'" -ErrorAction stop
}catch{
     Write-Error "Error occured while searching MP3 files!!"
     Write-Error $_.Exception.Message

}
if (-not $EXT){
   Write-Error "No MP3 files found on this computer..!"
   Exit
  
}Else {
        $numberoffiles = ($EXT | Measure-Object).count
        Foreach ( $File in $EXT){
        $output += New-Object PSObject -Property @{
                    "Name" = (Get-Item $File.name).Name
                    "Path" = $File.name
                    "Size in KB" = [math]::Round($((Get-Item $File.name).Length / 1024),2)
                   }

      }
        Write-Output "Total number of MP3 files = $numberoffiles"
        Write-Output "----------------------------------"
        $output | select Name, Path, "Size in KB" | fl
}
