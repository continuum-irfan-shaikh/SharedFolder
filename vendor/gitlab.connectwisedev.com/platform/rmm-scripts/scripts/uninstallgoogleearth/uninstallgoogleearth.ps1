# $Reboot = $False # user input

$Software = 'Google Earth'
$ErrorActionPreference = 'Stop'
try{
   # if ($Reboot) {$RestartArgument = '/forcerestart'}else{$RestartArgument = '/norestart'}
    $Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall','HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'
    
    $Product = Get-ChildItem $Registry -ErrorAction 'SilentlyContinue' | Get-ItemProperty | Where-Object { $_.DisplayName -match $Software -and $_.DisplayVersion -like "*$Version*" }
    $ProductGUID = $Product | Select-Object -ExpandProperty PSChildName -First 1

  if ($Product) {
        $process = Start-Process "msiexec.exe" -arg "/X $ProductGUID /qn" -Wait -PassThru -ErrorAction 'Stop'
        If ($process.exitcode -eq 0) {
            Write-Output "Successfuly uninstalled 'Google Earth v$Version'."
        }
        else {
            Write-Output "Failed to uninstall 'Google Earth v$Version'. Exitcode: $($process.exitcode)"        
        }
    }
    else {
        Write-Output "'Google Earth v$Version' is not installed on the system."
    }
}
catch{
    Write-Error $_
}
