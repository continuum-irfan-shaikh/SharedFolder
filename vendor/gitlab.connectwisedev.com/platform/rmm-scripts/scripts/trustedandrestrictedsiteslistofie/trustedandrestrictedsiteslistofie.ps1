 <#
    .Synopsis
     Retrieve list of Trusted and Restricted Sites for Internet Explorer
    .Description 
     To get the details of the Trusted and Restricted Sites for Internet Explorer on the machine
    .Help
     To get more information about the Trusted and Restricted Sites for Internet Explorer then refer below details
     Get-ChildItem -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings\ZoneMap\EscDomains\" | Get-Member
    .Author
     sushma.yerasi@continuum.net
    .Name 
     Retrieve list of Trusted and Restricted Sites for Internet Explorer
#>

$Trustedsites = @()
$Restrictedsites = @()
$ExecutionLog = @()

$CurrentKey = @()
$URLS = @()
$path = @()
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
$root = "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings\ZoneMap\EscDomains\" 
try { 
if(Test-Path $root) {
$paths = Get-ChildItem $root -recurse -ErrorAction SilentlyContinue
foreach ($path in $paths) {
$CurrentKey = Get-ItemProperty -Path $path.PsPath | select http,https
foreach ($key in $currentkey) { 
    if($key.http -eq 4)  { 
    $array = ($path.Name -split "\\EscDomains\\")[1] -split '\\';[array]::Reverse($array);
    $p = "http"
    $URLS = $p+'://'+($array -join '.')
    $Restrictedsites += $URLS 
    }
    if($key.https -eq 4) { 
    $array = ($path.Name -split "\\EscDomains\\")[1] -split '\\';[array]::Reverse($array);
    $p = "https"
    $URLS = $p+'://'+($array -join '.')
    $Restrictedsites += $URLS 
    }
    if($key.http -eq 2){ 
    $array = ($path.Name -split "\\EscDomains\\")[1] -split '\\';[array]::Reverse($array);
    $p = "http"
    $URLS = $p+'://'+($array -join '.')
    $Trustedsites += $URLS 
    }
    if($key.https -eq 2){ 
    $array = ($path.Name -split "\\EscDomains\\")[1] -split '\\';[array]::Reverse($array);
    $p = "https"
    $URLS = $p+'://'+($array -join '.')
    $Trustedsites += $URLS 
    }
}
}
$Tsites = $trustedsites | measure
If($Tsites.count -eq 0) {
$ExecutionLog +=  "`n"
$ExecutionLog += "There are no Trusted Sites on the machine" }
else{
$ExecutionLog +=  "`n"
$ExecutionLog += "Trusted Sites on the machine are mentioned below"
$ExecutionLog += $Trustedsites
}
$Rsites = $restrictedsites | measure
If($Rsites.count -eq 0) {
$ExecutionLog +=  "`n"
$ExecutionLog += "There are no Restircetd Sites on the machine" }
else{
$ExecutionLog +=  "`n"
$ExecutionLog += "Restricted Sites on the machine are mentioned below"
$ExecutionLog += $Restrictedsites
}
Write-Output $ExecutionLog 
exit;
}
else { $paths = Get-ChildItem $root -recurse -ErrorAction Stop}
}
Catch {
$ExecutionLog += "There are NO Trusted and Restricted sites on the Machine"
Write-output $ExecutionLog 
exit;
}
}
Catch { 
	$ExecutionLog += "Not able to reach the computer remotelt through WMI"
	Write-Output $executionlog
              exit;
 }
