<#
    .SYNOPSIS
        Remove/Uninstall Kaseya
    .DESCRIPTION
        Remove/Uninstall Kaseya software from the windows system. 
    .Help
        To get more details refer below command. 
        HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall, HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall
        Start-Process -NoNewWindow -FilePath Path_of_Uninstallation -ArgumentList " /s" -wait
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>

function verify {
if ((gwmi win32_operatingsystem | select osarchitecture).osarchitecture -eq "64-bit")
{
    $a =  Get-ChildItem -Path HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object {$_.DisplayName -match "kaseya" } | Select-Object -ExpandProperty UninstallString
}
else
{
    $a = Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object {$_.DisplayName -match "kaseya" } | Select-Object -expandProperty UninstallString
}
return $a
}

function FetchServices {

     $ServiceNames = Get-Service -displayname kaseya* | Select-Object name -ExpandProperty name
    return $ServiceNames
}

Try 
{
   $ver = verify
   if (!$ver)
    {
        Write-Output "`nKaseya agent not installed in this system $ENV:ComputerName"
    }
    else
    {
        $Fetching = FetchServices
        Start-Sleep 10
        foreach ($fetch in $Fetching){
        Stop-Service -name $fetch -Force -WarningAction  SilentlyContinue
        Start-Sleep 10 
        $ServiceStatus  = Get-WmiObject win32_service | where  {$_.name  -eq $fetch} | select startmode -ExpandProperty startmode
        if ($ServiceStatus -ne "unknown")
        {
        Set-Service $fetch -StartupType Disabled
        }
        }

        Start-Sleep 10 
        Start-Process -NoNewWindow $ver -ArgumentList " /s"
        Start-Sleep 10 
       
        $ver1 = verify
  
    if(!$ver1)
        {
          Write-Output "`nKaseya agent is successfully removed from the system $ENV:ComputerName" 
           Write-Output "Kindly reboot the system $ENV:ComputerName" 
          
          $keypath = (Get-Service -DisplayName "Kaseya Agent Endpoint*").name
          if ($keypath -ne $null)
          {
          if (test-path "HKLM:SYSTEM\CurrentControlSet\Services\$keypath")
          {
          Remove-Item -path "HKLM:SYSTEM\CurrentControlSet\Services\$keypath" -recurse       
          }
          }
        } 
    else
        {
        Write-Error "`nThere is some issue with uninstallation. Kindly check manually."
        }            
    }
}
catch
{
    write-output "`n"$_.Exception.Message
}
