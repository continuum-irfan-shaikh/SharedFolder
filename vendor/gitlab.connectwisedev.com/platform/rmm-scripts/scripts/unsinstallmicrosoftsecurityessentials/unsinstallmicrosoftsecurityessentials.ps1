<#
	.Scriptname
	Uninstall microsoft security essentials.
	.Author.
	Nirav Sachora
	.Requirements
	Run script with highest privileges.
#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

function Uninstall_MSE{
$pinfo = New-Object System.Diagnostics.ProcessStartInfo
$pinfo.FileName = "C:\Program Files\Microsoft Security Client\Setup.exe"
$pinfo.RedirectStandardError = $true
$pinfo.RedirectStandardOutput = $true
$pinfo.UseShellExecute = $false
$pinfo.Arguments = "/X /s"
$p = New-Object System.Diagnostics.Process
$p.StartInfo = $pinfo
$p.Start() | Out-Null
$p.WaitForExit()
return $p.ExitCode
}

if((Test-path "C:\Program Files\Microsoft Security Client\Setup.exe") -eq $false){
Write-output "Microsoft security essentials is not installed on this system"
Exit;
}

if(((Uninstall_MSE) -eq 0) -and ($reboottype -eq "noreboot")){
Write-Output "Microsoft Security Essentials has been uninstalled from the system."
}
elseif(((Uninstall_MSE) -eq 0) -and ($reboottype -eq "reboot")){
Write-Output "Microsoft Security Essentials has been uninstalled from the system `n`n System will reboot now."
Restart-Computer
}
elseif(((Uninstall_MSE) -eq 0) -and ($reboottype -eq "forcereboot")){
Write-Output "Microsoft Security Essentials has been uninstalled from the system `n`n System will reboot now."
Restart-Computer -Force
}
else{
Write-Error "Failed to uninstall microsoft security essentials."
}
