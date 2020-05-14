<#  
.SYNOPSIS  
    Disable LogMeINPro remote printing. 
.DESCRIPTION  
    LogMeInPro user will not be able to print thorugh LogMeInPro application from remote host. 
.NOTES  
    File Name  : DisableLogMeInProRemotePrintingv1.0.ps1
    Author     : Ratnesh Mishra  
    Modified   : Durgeshkumar Patel
    Requires   : PowerShell V2 or greater.   
.PARAMETERS
    
.HELP
#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

$RegPath = "HKLM:\SOFTWARE\LogMeIn\V5\RemoteControl"
<#===Check if LogMeIn is installed or not from registry key====#>
try
{
$KeyName = Get-ItemProperty $RegPath 
}
catch
{
Write-Error "LogMeIn Pro not installed : $($_.Exception.Message)"
Exit
}
<#=======Check registry value if exist will change it else will create it ======#>
 try
 {
        if ((Get-ItemProperty $RegPath -Name EnableRemotePrinting -ErrorAction SilentlyContinue) -ne $Null)
        {
            Set-ItemProperty -Path $RegPath -Name EnableRemotePrinting -Value 0
        } 
        else
        {
            New-ItemProperty -Path $RegPath -Name EnableRemotePrinting -Value 0 -PropertyType DWORD -Force | Out-Null
        }
        
        if ((Get-ItemProperty $RegPath -Name ForceBitmapPrinting -ErrorAction SilentlyContinue) -ne $Null)
        {
            Set-ItemProperty -Path $RegPath -Name ForceBitmapPrinting -Value 0
        } 
        else
        {
            New-ItemProperty -Path $RegPath -Name ForceBitmapPrinting -Value 0 -PropertyType DWORD -Force | Out-Null
        }   
        Write-Output "Task [Disable LogMeIn Pro Remote Printing] Completed Successfully"       
   } 
catch
    {     
        Write-Error "LogMeIn Pro Remote printing configuration issue : $($_.Exception.Message)"
        Exit       
    }
