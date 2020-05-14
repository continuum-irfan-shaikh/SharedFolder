<#

List of Variable :
===================
$overwrites	     ===> Boolean ===> CheckBox ===> Title :"Overwrite existing file"
$CreateDirectory ===> Boolean ===> CheckBox ===> Title :"Create directory if required"
$Inputpath	     ===> String  ===> TextBox  ===> Title :"New File Location:"
$value           ===> String  ===> MultiLineTextBox ===> Title :"File Contents:"

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

$Location = @()
$path = @()
$finalpath = @()
$output = @()
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
Write-Output $executionlog
$path = Split-Path -Path $Inputpath -parent 
$file = Split-Path -Path $Inputpath -leaf
$drive = Split-Path -Path $Inputpath -Qualifier
Set-Location $drive\
if($CreateDirectory -eq $True) {

$finalpath = $path.Split("{\}")
foreach ($final in $finalpath) {
if((test-path $final) -eq $True) { 
Set-Location $final 
}
else {
$Location  = Get-Location | select path 
$output += New-Item -Path $Location.path -name $final -ItemType "directory"
Set-Location $final
}
}}

if($overwrites -eq $True) {

If((Test-Path $Inputpath) -eq $True) {
Remove-Item $Inputpath
$output = New-Item -Path $path -Name $file -ItemType File -ErrorAction Stop
Set-Content -Path $Inputpath -value $value
Write-Output "Created the Text File on the specified path on the Machine"
}
else {
try{
$output = New-Item -Path $path -Name $file -ItemType File -ErrorAction Stop
Set-Content -Path $Inputpath -value $value
}
Catch {
Write-Output "The Path which was given as an input doesn't exist on the Machine"
}
}
}
else {
If((Test-Path $Inputpath) -eq $False) {
try{
$output = New-Item -Path $path -Name $file -ItemType File -ErrorAction stop
Set-Content -Path $Inputpath -value $value
Write-Output "Created the Text File on the specified path on the Machine"
}
Catch {
Write-Output "The Path which was given as an input doesn't exist on the Machine"
}
}
else {
Write-Output "The Path which was given as an input already exist on the Machine"
}}
}
Catch { 
	$ExecutionLog += "Not able to reach the computer remotelt through WMI"
	Write-Output $executionlog
              exit;
 }
	
