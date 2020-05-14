<#
Name : Retrieve shared folder information
Category : Data Collection

    .Synopsis
        Retrieve Shared Folder Information.
    .Author
        Santosh.Dakolia@continuum.net
    .Name 
        Retrieve shared folder information
#>

Try{
$computername= $env:computername

    $ShareData = get-WmiObject -class Win32_Share -ComputerName $computername
    IF(!$sharedata) {"No Data available for Shared Folders"}
    $output = @()

ForEach($SD in $ShareData){
    $MA = IF($sd.AllowMaximum -eq 'True'){"Allow Maximum is Set to True"} Else {$sd.MaximumAllowed}

        $output +=  New-Object psobject -Property @{
            "Allow Maximum" = $sd.allowmaximum
            "Caption" = $sd.Caption
            "Maximum Allowed" = $MA
            "Name" = $sd.Name
            "Path" = $sd.Path
       }
}

    $output | FL "Allow Maximum", Caption, "Maximum Allowed", Name, Path 
    $count = ($sharedata | Measure-Object).count
    Write-output "Shared Folder Count :: $count" 

}Catch{
    Write-Error "Error occured while retrieving Data..!! $_.Exception.Message"
}
