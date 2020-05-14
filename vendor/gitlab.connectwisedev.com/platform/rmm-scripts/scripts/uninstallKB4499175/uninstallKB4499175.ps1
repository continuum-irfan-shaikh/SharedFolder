<#
.SYNOPSIS
      KB4499175 uninstaller.

.DESCRIPTION
      This script will uninstall KB4499175 if install status is "Install Pending".
      This script supports  vesrions of Windows 7 SP1 higher and Windows 2008 R2 SP1 higher.
.Author
    GRT
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
$OSName = (Get-WMIObject win32_operatingsystem).Caption
if(-not($OSver.Major -eq 6 -and $OSver.Minor -eq 1 -and $OSver.Build -ge 7601)) {
   Write-Output "KB4499175 is not applicable for the OS : $OSName. No action needed."
   Exit
}

$KB = "KB4499175"
$vObj = [System.Environment]::OSVersion.Version
[int]$OSVer = -join ($vObj.Major,$vObj.Minor,$vObj.Build)
$ComputerName = $env:COMPUTERNAME
$KBInstalled = (DISM /Online /Get-packages | ?{ $_ -match $KB }) -Replace("Package Identity : ", "")

if($KBInstalled) {
     foreach ($package in $KBInstalled){
          $State = (DISM /Online /Get-PackageInfo /PackageName:$package | ?{$_ -match "State"}).Split(':')[-1].trim()

          if ($State -eq "Install Pending"){
                $cmdoutput = DISM /Online /Remove-Package /PackageName:$package /quiet /norestart 2>&1
                Start-Sleep -s 3
                Wait-Process "dism" -ErrorAction SilentlyContinue
          }ElseIf($State -eq "Installed") {
                  Write-Output "KB4499175 is already installed. No action needed."
                  Exit
          }Else{ Write-Error "Action is not defined for the state : $State"
                 Exit
               }
     }     
     $KBExist = DISM /Online /Get-packages | ?{ $_ -match $KB }
     if (!($KBExist)){ Write-Output "KB4499175 is successfully uninstalled."}
     Else{ Write-Error "KB4499175 unistallation failed! $cmdoutput" }
}
Else { Write-output "KB4499175 is not installed. No action needed." }

