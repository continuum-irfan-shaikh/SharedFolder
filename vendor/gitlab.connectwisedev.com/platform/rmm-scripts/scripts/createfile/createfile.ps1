$array = $path.Split("\")
$folderPath = $array[0..($array.length-2)]
$folderPath = $folderPath -join "\"

if (!$createFolder -and !(Test-Path $folderPath)){
  Write-Error "Folder doesn't exist"
  return}

if (New-Item $path -type file -value $value -Force:$force) {
  Write-Output "Successfully created file. Path: $path"}
