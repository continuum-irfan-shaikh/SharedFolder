<#

	.ScriptName
	Nirav Sachora
	.Description
	Script will create restore point.
	Author
	Nirav Sachora
	.Requirements
	Script should run with highest privileges.
	
#>

<#

$RestorePointType = "Application Install"
$EventType = "Begin System Change"

#>

$name = "Restorepoint" + (Get-Date -Format "ddMMyyyyHHmm")

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

if((Get-WMIObject win32_operatingsystem).Caption -like "*server*"){
    Write-Error "Script is not supported on Server operating system."
    Exit;
}

switch($RestorePointType){
    "Application Install"{$RestorePointType = 0;Break;}
    "Application Uninstall"{$RestorePointType = 1;Break}
    "Device Driver Install"{$RestorePointType = 10;Break}
    "Modify Settings"{$RestorePointType = 12;Break}
    "Cancelled Operation"{$RestorePointType = 13;Break}
}

switch($EventType){
    "Begin System Change"{$EventType = 100;Break;}
    "End System Change"{$EventType = 101;Break}
    "Begin Nested System Change"{$EventType = 102;Break}
    "End Nested System Change"{$EventType = 103;Break}
}

$command = "Wmic.exe /Namespace:\\root\default Path SystemRestore Call CreateRestorePoint $name, $EventType, $RestorePointType"
#$ErrorActionpreference = "SilentlyContinue"
$result = &{cmd.exe /c "$command" 2>&1}


    if ($result -like "*ReturnValue = 0;*") {
        "Restore point created successfully.`nRestore Point Name : $name"    
    }
    else{
        Write-Error "Failed to create restore point." 
        Exit;
    }
