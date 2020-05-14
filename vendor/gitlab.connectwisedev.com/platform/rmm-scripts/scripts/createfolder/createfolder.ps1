# Create Folder
# sourcePath Network
if ($sourceType -eq "Network"){
    $pathendchr = $sourcePath.Substring($sourcePath.get_Length()-1)
  if ( $pathendchr -eq "\") {
       $sourcePath = $sourcePath.Substring(0,$sourcePath.Length-1)
  }
  if (!(Test-Path $sourcePath+'\')){
    $DrvLtrs = [char[]]([char]'C'..[char]'Z')
    $InuseDrvLtrs = get-wmiobject win32_logicaldisk | select -expand DeviceID | % { $_.Trim(":") }
    $FreeDrvLtr = $DrvLtrs | ?{ $InuseDrvLtrs -notcontains $_ }
    $net = new-object -ComObject WScript.Network
    try {
           $net.MapNetworkDrive($FreeDrvLtr[0]+':',$sourcePath, $false, $userName, $password)
        }Catch{
           Write-Error "Could not connect to network path : $_.Exception.Message" 
        }
  }
  if (Test-Path $sourcePath){
    $folder = New-Item $($sourcePath + "\" + $folderName) -type directory 
    $net.RemoveNetworkDrive($FreeDrvLtr[0]+':',$True)
  } else {
    Write-Error "Path $sourcePath was not found"
  }
# sourcePath Local
} else {
  if (Test-Path $sourcePath){
    $pathendchr = $sourcePath.Substring($sourcePath.get_Length()-1)
    if ( $pathendchr -ne "\") {
        $folder = New-Item $($sourcePath + "\" + $folderName) -type directory
    }Else{
        $folder = New-Item $($sourcePath + $folderName) -type directory
    }

  } else {
    Write-Error "Path $sourcePath was not found"
  }
}
$folder |  Select @{Name="Folder Name";Expression={$_.Name}},
                  @{Name="Folder Path";Expression={$_.FullName}},
                  @{Name="Creation Time";Expression={$_.CreationTime}} | FL
