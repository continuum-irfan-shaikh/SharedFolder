<#
    .SYNOPSIS
       On Demand Validate All Smart Attributes
    .DESCRIPTION
       On Demand Validate All Smart Attributes
    .Author
       Santosh.Dakolia@continuum.net    

       # File Name : zSmart.exe
       #Location : C:\Program Files (x86)\SAAZOD
#>
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

#Check Arcitecture 
$OStype = (Get-WMIObject Win32_OperatingSystem).ProductType #Workstation = 1 Server = 3
$OSArch = (Get-WMIObject Win32_OperatingSystem).OSArchitecture

#Get path based on Arcitecture 
Switch($OSArch){
   "64-bit" {  if ($OStype -eq 1){$smartupd = ${ENV:ProgramFiles(x86)}+"\SAAZOD\zSmart.exe" }
               elseif($OStype -eq 3){ Write-Output "Smart Attributes Validation is not applicable for Server OS"; exit}   
            }
   "32-bit" {  if ($OStype -eq 1){$smartupd = ${ENV:ProgramFiles}+"\SAAZOD\zSmart.exe" }
               elseif($OStype -eq 3){ Write-Output "Smart Attributes Validation is not applicable for Server OS"; exit}            
            }               
}

#Excecute EXE with argument
try{
      start-process $smartupd -ArgumentList SmartAll -Wait -EA stop
      Write-Output "Smart Attributes Validation completed."
}catch{ 
    if ( $_.Exception.Message -like "*The system cannot find the file specified*" ) {
       Write-Error "Smart Attributes Validation executable not found..!!" 
    }Else { Write-Error $_.Exception.Message }
}
