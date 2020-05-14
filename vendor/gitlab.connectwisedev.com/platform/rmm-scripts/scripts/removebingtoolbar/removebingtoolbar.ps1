<#
    .SYNOPSIS
        Remove Bing toolbar
    .DESCRIPTION
        Remove Bing toolbar from the system
    .Help
        To get more details refer below command. 
        HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall, HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall
        MsiExec.exe uninstallpath /quiet /norestart
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>

function get-toolbar {
if ((gwmi win32_operatingsystem | select osarchitecture).osarchitecture -eq "64-bit")
{
    $a =  Get-ChildItem -Path HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object {$_.DisplayName -match "Bing Bar"}| Select-Object -ExpandProperty UninstallString
}
else
{
    $a = Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object {$_.DisplayName -match "Bing Bar"} | Select-Object -expandProperty UninstallString
}
return $a
}

function Stop-Webbrowser {
    Get-Process | Where {$_.Name -eq "iexplore"} | kill -Force
    Get-Process | Where {$_.Name -eq "firefox"} | kill  -Force
    Get-Process | Where {$_.Name -eq "microsoftedge"} | kill  -Force 
}

Try 
{
   $toolbar= get-toolbar
   if (!$toolbar)
    {
        Write-Output "`nBing toolbar not installed on this system $ENV:ComputerName"
    }
    else
    {
        Stop-Webbrowser
        
        Start-Sleep 3
        
        $a = $toolbar.split(' ')[1] 
        
        MsiExec.exe $a /quiet /norestart
        
        Start-Sleep 5
        
        Stop-Webbrowser
        
        $toolbar1= get-toolbar
  
    if(!$toolbar1)
        {
          Write-Output "`nBing toolbar removed from the system $ENV:ComputerName" 
        } 
    else
        {
        Write-Error "`nThere is some issue with uninstallation. Kindly check manually."
        }            
    }
}
catch
{
    Write-Error $_.Exception.Message  
}

