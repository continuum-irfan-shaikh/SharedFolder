<#
    .SYNOPSIS
       Uninstall Google toolbar for internet explorer.
    .DESCRIPTION
       Uninstall Google toolbar for internet explorer.
    .Author
       narayan.gouda@continuum.net    
#>
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

$OSver = [System.Environment]::OSVersion.Version
$OSArch = [intPtr]::Size
try{
    Switch($OSArch){
        4 { $UninstallString = (Get-ItemProperty -path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\'{2318C2B1-4965-11d4-9B18-009027A5CD4F}' -ErrorAction Stop).UninstallString }      
        8 { $UninstallString = (Get-ItemProperty -path HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\'{2318C2B1-4965-11d4-9B18-009027A5CD4F}' -ErrorAction Stop).UninstallString }
    }
}catch{
      if ( $_.Exception.Message -like  "*Cannot find path*"){
         Write-Error "Google Toolbar not found on this computer..!!"
         Exit
      }else { 
         Write-Error $_.Exception.Message
         Exit     
      }
}
$exe_path, $arg = $UninstallString.Split('/')
$cmd = $exe_path -replace '"', ""

try{
    start-process $cmd -ArgumentList /uninstall -Wait -ErrorAction Stop
    Write-Output "Google toolbar uninstalled successfully."
}catch{ 
    Write-Error "Error occured while uninstalling Google toolbar..!! $_.Exception.Message"
}
