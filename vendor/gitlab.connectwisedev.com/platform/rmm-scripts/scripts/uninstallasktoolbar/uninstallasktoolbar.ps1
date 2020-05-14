<#
    name : uninstall Ask Toolbar
    Category : Application 
#>

$ErrorActionPreference = 'Stop'
try{
    $Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall','HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'
    $Product = Get-ChildItem $Registry -ErrorAction 'SilentlyContinue' | Get-ItemProperty | Where-Object { $_.DisplayName -match "ask toolbar" -and $_.DisplayVersion -like "*$Version*" }
    $ProductGUID = [regex]::Matches($($Product.uninstallstring),'{([-0-9A-F]+?)}') | Select-Object -ExpandProperty value
    
    if ($Product) {
        $process = Start-Process "msiexec.exe" -arg "/X $ProductGUID /qn" -Wait -PassThru -ErrorAction 'Stop'
        If ($process.exitcode -eq 0) {
            Write-Output "Successfuly uninstalled 'Ask Toolbar v$Version'."
        }
        else {
            Write-Output "Failed to uninstall 'Ask Toolbar v$Version'. Exitcode: $($process.exitcode)"        
        }
    }
    else {
        Write-Output "'Ask Toolbar v$Version' is not installed on the system."
    }
}
catch{
    Write-Error $_
}
