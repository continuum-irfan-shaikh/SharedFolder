<#
    .SYNOPSIS
        Remove GFI MAX Agent (Advanced Monitoring Agent)
    .DESCRIPTION
        Remove GFI MAX Agent (Advanced Monitoring Agent)
    .Help
        Use hardcoced uninstall string.
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>
function get-product {
    if ((gwmi win32_operatingsystem | select osarchitecture).osarchitecture -eq "64-bit") {
        $a = Get-ChildItem -Path HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object { $_.DisplayName -match "Advanced Monitoring Agent" } | Select-Object -ExpandProperty UninstallString
    }
    else {
        $a = Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object { $_.DisplayName -match "Advanced Monitoring Agent" } | Select-Object -expandProperty UninstallString
    }
    return $a
}
 
Try {
    $product = get-product
    if (!$product) {
        Write-Output "`nAdvanced Monitoring Agent(GFI MAX Agent) not installed on this system $ENV:ComputerName"
    }
    else {
            
        $process = Start-Process $product -arg "/VERYSILENT /SILENT" -Wait -PassThru -ErrorAction 'Stop'      
     
        If (($process.exitcode -eq '3010') -or ($process.exitcode -eq '0')) {
            
            Write-Output "`nAdvanced Monitoring Agent(GFI MAX Agent) uninstalled from the system $ENV:ComputerName"
        } 
        else {
            Write-Warning "`nFailed to uninstall Advanced Monitoring Agent(GFI MAX Agent) from the system $ENV:ComputerName."
        }             
    }
}
catch {
    Write-Error $_.Exception.Message  
}
