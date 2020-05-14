<#
.SYNOPSIS
    The purpose of this script is to uninstall KB articles.
  
.DESCRIPTION
    This script will uninstalls KB4103718 & KB4103712. 
    This script supports  vesrions of Windows 7 higher and Windows 2008 R2 higher
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
$Hostname = $env:COMPUTERNAME

if(-not($OSver.Major -eq 6 -and $OSver.Minor -eq 1 -and $OSver.Build -ge 7601)) {
   Write-Error "Not supported on this operating system!"
   Exit 
}

$KBsToRemove = ("KB4103718", "KB4103712")
$ComputerName = $env:COMPUTERNAME

$Result = @()
Foreach($KB in $KBsToRemove ){
        $Hotfix = Get-HotFix -Id $KB -ErrorAction SilentlyContinue
        if($Hotfix) {  
                $PackageName = (DISM /Online /Get-packages | ?{ $_ -match $KB }) -Replace("Package Identity : ", "")
                $cmdoutput = DISM /Online /Remove-Package /PackageName:$PackageName /quiet /norestart 2>&1
                Start-Sleep -s 3
                Wait-Process "dism" -ErrorAction SilentlyContinue
                if($cmdoutput){
                   $UninstallStatus = "Failed $cmdoutput"
                }
                Elseif( Get-HotFix -Id $KB -ErrorAction SilentlyContinue ){
	                $UninstallStatus = "Failed"
                    
                }Else {
                    $UninstallStatus = "Success"
                }
               
		}Else { $UninstallStatus = "Not found" }

        $obj = New-Object PSObject -Property @{
           										Computer = $ComputerName;
                                                HotFixID = $KB;
        										Status = $UninstallStatus
                                              }
        $Result += $obj
}
Write-Output $Result | Select Computer, HotfixID, Status | FT -AutoSize

