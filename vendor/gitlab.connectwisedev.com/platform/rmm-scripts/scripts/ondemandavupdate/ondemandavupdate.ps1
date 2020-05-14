<#
    .SYNOPSIS
       On Demand AntiVirus update
    .DESCRIPTION
       Updates antivirus (supported AV) installed on the machine. 
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

$OStype = (Get-WMIObject Win32_OperatingSystem).ProductType #Workstation = 1 Server = 3
$OSArch = (Get-WMIObject Win32_OperatingSystem).OSArchitecture

Switch($OSArch){
   "64-bit" {  if ($OStype -eq 1){$avupd = ${ENV:ProgramFiles(x86)}+"\SAAZOD\BaseComponents\AVUpd\zAVUPD.exe" }
               elseif($OStype -eq 3){ $avupd = ${ENV:ProgramFiles(x86)}+"\SAAZOD\zAVUPD.exe"}   
            }
   "32-bit" {  if ($OStype -eq 1){$avupd = ${ENV:ProgramFiles}+"\SAAZOD\BaseComponents\AVUpd\zAVUPD.exe" }
               elseif($OStype -eq 3){ $avupd = ${ENV:ProgramFiles}+"\SAAZOD\zAVUPD.exe"}            
            }               
}

try{
      start-process $avupd -ArgumentList FORCE -Wait -EA stop
      Write-Output "Antivirus update completed."

}catch{ 
    if ( $_.Exception.Message -like "*The system cannot find the file specified*" ) {
       Write-Error "Antivirus update executable not found..!!" 
    }Else { Write-Error $_.Exception.Message }
}
