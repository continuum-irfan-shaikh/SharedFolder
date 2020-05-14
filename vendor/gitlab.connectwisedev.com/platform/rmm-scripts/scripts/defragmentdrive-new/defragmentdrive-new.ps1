<#
	.Script name
	Drive Defragmentation
	.Author
	Nirav Sachora
	.Version
	1.0
	.Variable List
		$DefragDrive = "All Drives","Windows Drives","Specify Drives"
		$A   to    $Z
	.Last updated date
	14-1-2020
#>
<#
$DefragDrive = "All Drives"

$c = "true"
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

function Returnvalue_generator($value) {

    switch ($value) {
        0 { $resultMsg = 'Success'; break }
        1 { $resultMsg = 'Access Denied'; break }
        2 { $resultMsg = 'Not Supported'; break }
        3 { $resultMsg = 'Volume Dirty Bit Set'; break }
        4 { $resultMsg = 'Not Enough Free Space'; break }
        5 { $resultMsg = 'Corrupt MFT Detected'; break }
        6 { $resultMsg = 'Call Cancelled'; break }
        7 { $resultMsg = 'Cancellation Request Requested Too Late'; break }
        8 { $resultMsg = 'Defrag In Progress'; break }
        9 { $resultMsg = 'Defrag Engine Unavailable'; break }
        10 { $resultMsg = 'Defrag Engine Error'; break }
        default { $resultMsg = 'Unknown Error' }
    }
    return $resultMsg
}


function defrag_drive($drivename) {
    if($DefragDrive -eq "Specify Drives"){
        $alldrivenames = Get-WmiObject win32_volume -Filter "drivetype = 3" | Where-Object { $_.DriveLetter -ne $null } | Select-Object -ExpandProperty DriveLetter
        if(!($alldrivenames -contains $drivename)){
            "-"*30+"`n[Not Found]`"$drivename`" drive is not present in the system.`n"+"-"*30
            return
        }  
    }
    try {
        $ErrorActionPreference = "Stop"
        $result = Get-WmiObject win32_volume -filter "driveletter='$drive'" | Invoke-WmiMethod -Name Defrag | Select-Object -ExpandProperty ReturnValue
    }
    catch {
        $ErrorActionPreference = "Continue"
        "[ERROR]:Unknown error occured"  #can be changed, rather that exit, return from function
        Exit;
    }

    $resultMsg = Returnvalue_generator -value $result

    if ($result -eq 0) {
        "-"*30+"`n[MSG]:Defragmentation completed for $drive`n"+"-"*30
    }
    else {
        "[ERROR]$resultMsg on $drive"
    }
}

function alldrives {

    $alldrivenames = Get-WmiObject win32_volume -Filter "drivetype = 3" | Where-Object { $_.DriveLetter -ne $null } | Select-Object -ExpandProperty DriveLetter

    if (!$alldrivenames) {
        "[ERROR]:Unknown Error occured, Script will now exit"
        Exit;
    }
    foreach ($drive in $alldrivenames) {
       defrag_drive -drivename $drive
    }
}

function Windowsdrive {
    $drive = (Split-Path ($env:windir) -Parent) -replace "\\", ""
    defrag_drive -drivename $drive
}

function selecteddrives {
    $alldriveoptions = 'a','b','c','d','e','f','g','h','i','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z'
    $NuOfdrives = 0
    foreach ($driveentry in $alldriveoptions){
        $Erroractionpreference = 'SilentlyContinue'
        if(((get-variable $driveentry).Value) -eq 'True'){
	    $NuOfdrives += 1
            $drive = "$driveentry" + ":"
            defrag_drive -drivename $drive
        } 
    }
    if($NuOfdrives -eq 0){
	"-"*30+"`n[ERROR]:No drives selected`n"+"-"*30
	Exit;
    }
}

switch ($DefragDrive){
    "All Drives"{alldrives}
    "Windows Drives"{Windowsdrive}
    "Specify Drives"{selecteddrives}
}
