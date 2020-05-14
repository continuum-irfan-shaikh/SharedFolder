<#
    .SYNOPSIS
        Create Shared Folder
    .DESCRIPTION
        Create Shared Folder
    .Help
        
    .Author
        Santosh
    .Version
        1.0
#>

<#Create Shared Folder
Varialbles:-
$path
$ShareName

Radio Buttons True/False
$MaxAllowed   $connectionAllowed 

If $connectionAllowed true then
Text Box (Take Integer only)
$connectionAllowed

#Examples:-
$path = "c:\durgesh\New11"
$ShareName = "DurgeshShareFolder11"
$MaxAllowed = $true
#$connectionAllowed = '11'
$Description = "This is a Shared Folder1"

Maximumnumberofconnections = 
#>

$Computer = $env:COMPUTERNAME
 
If ($Maximumnumberofconnections -eq "MaxAllowed") {
    $MaxConnections = $null
}
else {
    $MaxConnections = if ($connectionAllowed) { [UInt32]$connectionAllowed }
}


If (!(Test-Path $path)) {

    New-Item -Path $path -ItemType Directory -Force | Out-Null
}


# Create the Share
$Share = [WmiClass]"\\$Computer\root\cimv2:Win32_share"
$InParams = $Share.psbase.GetMethodParameters("Create")
$InParams.Description = $Description
$InParams.MaximumAllowed = $MaxConnections
$InParams.Name = $shareName
$InParams.Path = $Path
 
$Result = $Share.PSBase.InvokeMethod("Create", $InParams, $Null)
   
# Check if it was successfull
$rvalue = Switch ($Result.ReturnValue) {
    0 { "Success" }
    2 { "Access Denied" }     
    8 { "Unknown Failure" }     
    9 { "Invalid Name" }     
    10 { "Invalid Level" }     
    21 { "Invalid Parameter" }     
    22 { "Duplicate Share" }     
    23 { "Redirected Path" }     
    24 { "Unknown Device or Directory" }
    25 { "Net Name Not Found" }
    Default { "Unknown Error" }
}

Write-Output $rvalue
