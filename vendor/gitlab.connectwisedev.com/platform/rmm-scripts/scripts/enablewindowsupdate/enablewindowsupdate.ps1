$ExecutionLog = @()
[double]$OSVersion=[Environment]::OSVersion.Version.ToString(2)
# OS Version and Powershell comparison.
if (($osversion -lt '6.1') -or ($PSVersionTable.PSVersion.Major -lt '2'))
        {
            $executionlog += 'Prerequisites to run the script is not valid, Hence Script Exceution stopped' #'Script is design for windows 7 and Above Members only, Script Execution Stopped.'
            Write-Output $executionlog
            exit;
        } 

#To check the status of the computer
try {

    $Comp=get-wmiobject win32_computersystem
    $computer = $comp.name ; $DomainNamee = $comp.domain
    $ExecutionLog += "ComputerName : $computer"
    $ExecutionLog += "Domain/Workgroup : $DomainNamee"
    $ExecutionLog += "`n"

$root = "HKLM:\Software\Policies\Microsoft\Windows\WindowsUpdate\AU\" 

$service = Get-Service wuauserv | select status
if($service.status -eq "Running") 
{
try { 
if(Test-Path $root) {
$value = Get-ItemProperty -path $root | select -ExpandProperty NoAutoUpdate
If($value -eq 0) {
$ExecutionLog += "The Windows Update is already enabled on the Machine"
Write-Output $ExecutionLog 
}
else {
$output = reg add "HKEY_LOCAL_MACHINE\Software\Policies\Microsoft\Windows\WindowsUpdate\AU\" /v NoAutoUpdate /t REG_DWORD /d 0 /f 
$ExecutionLog += "The Windows Update is enabled on the Machine"
Write-Output $ExecutionLog 
}
}
else { $paths = Get-ChildItem $root -recurse -ErrorAction Stop}
}
Catch {
$ExecutionLog += "Windows Update Installer is not installed successfully on the machine"
Write-output $ExecutionLog 
}
}
else { 
$ExecutionLog += "The windows update service is not in running state on the server"
	Write-output $executionlog
}}
Catch { 
	$ExecutionLog += "Not able to reach the computer remotelt through WMI"
	Write-Error $executionlog       
 }
