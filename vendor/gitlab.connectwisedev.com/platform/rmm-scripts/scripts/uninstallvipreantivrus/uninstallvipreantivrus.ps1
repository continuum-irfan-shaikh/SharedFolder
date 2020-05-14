# $Reboot = $False # user input
# $Version = '9.6' # user input

$Software = 'Vipre'
$ErrorActionPreference = 'Stop'
try{
    if ($Reboot) {$RestartArgument = '/forcerestart'}else{$RestartArgument = '/norestart'}
    $Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall','HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'
    
    $Product = Get-ChildItem $Registry -ErrorAction 'SilentlyContinue' | Get-ItemProperty | Where-Object { $_.DisplayName -match $Software -and $_.DisplayVersion -like "*$Version*" }
    $ProductGUID = $Product | Select-Object -ExpandProperty PSChildName -First 1
    
    if ($Product) {
        $process = Start-Process "msiexec.exe" -arg "/X $ProductGUID /qn $RestartArgument" -Wait -PassThru -ErrorAction 'Stop'
        If ($process.exitcode -eq 0) {
            Write-Output "Successfuly uninstalled 'Vipre Antivirus v$Version'."
        }
        else {
            Write-Output "Failed to uninstall 'Vipre Antivirus v$Version'. Exitcode: $($process.exitcode)"        
        }
    }
    else {
        Write-Output "'Vipre Antivirus v$Version' is not installed on the system."
    }
}
catch{
    Write-Error $_
}
