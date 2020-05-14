<#
    .SYNOPSIS
       On Demand Warranty Deployment
    .DESCRIPTION
       On Demand Warranty Deployment
    .Author
       Santosh.Dakolia@continuum.net    

       # File Name : zWrnPDtls.exe
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
   "64-bit" {  if ($OStype -eq 1){$wrupd = ${ENV:ProgramFiles(x86)}+"\SAAZOD\zWrnPDtls.exe" }
               elseif($OStype -eq 3){ $wrupd = ${ENV:ProgramFiles(x86)}+"\SAAZOD\zWrnPDtls.exe"}   
            }
   "32-bit" {  if ($OStype -eq 1){$wrupd = ${ENV:ProgramFiles}+"\SAAZOD\zWrnPDtls.exe" }
               elseif($OStype -eq 3){ $wrupd = ${ENV:ProgramFiles}+"\SAAZOD\zWrnPDtls.exe"}            
            }               
}

#Excecute EXE with argument
try{
      start-process $wrupd -Wait -EA stop
      Write-Output "On Demand Warranty Deployment completed."

}catch{ 
    if ( $_.Exception.Message -like "*The system cannot find the file specified*" ) {
       Write-Error "On Demand Warranty Deployment executable not found..!!" 
    }Else { Write-Error $_.Exception.Message }
}
