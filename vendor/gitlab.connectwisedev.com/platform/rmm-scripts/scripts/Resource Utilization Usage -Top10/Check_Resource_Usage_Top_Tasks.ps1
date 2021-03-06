CLS
#$Sort_CPU_Memory = 'CPU'
#$IsAggregated = $True


##############################
######Check PreCondition######
##############################
 
$hostVerSionMajor = ($PSVersionTable.PSVersion.Major).ToString()
$hostVerSionMinor = ($PSVersionTable.PSVersion.Minor).ToString()
$hostVersion = $hostVerSionMajor +'.'+ $hostVerSionMinor 
 
$osVersionMajor = ([System.Environment]::OSVersion.Version.major).ToString()
$osVersionMinor = ([System.Environment]::OSVersion.Version.minor).ToString()
$osVersion = $osVersionMajor +'.'+ $osVersionMinor
 
[boolean]$isPsVersionOk = ([version]$hostVersion -ge [version]'2.0')
[boolean]$isOSVersionOk = ([version]$osVersion -ge [version]'6.0')
      
Write-Host "`nPowershell Version : $($hostVersion)"
if(-not $isPsVersionOk){
   
  Write-Warning "PowerShell version below 2.0 is not supported"
  return 
 
}
 
Write-Host "OS Name : $((Get-WMIObject win32_operatingsystem).Name.ToString().Split("|")[0])`n"  
if(-not $isOSVersionOk){
 
   Write-Warning "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
   return 
 
}
######################################################################
#################################################################################################

function Get_username_Description
{
Param([int]$PID1)

$ErrorActionPreference = 'silentlycontinue'
$Get_PRocess = Get-Process | ? {$_.id -eq "$PID1"}

$Win32_PRocess = Get-WmiObject -Class Win32_Process | ? {$_.ProcessId -eq "$PID1"} 
$processOwner = $Win32_PRocess.GetOwner()
$Username = $processOwner.Domain + '\' + $processOwner.User

$aaa = New-Object psobject -Property @{

Description = $Get_PRocess.Description
Username = $Username

} | select Description,Username

if($aaa.Description -eq $null -or $aaa.Username -eq $null){
'Not Found'
'Not Found'
}
else{
$aaa.Description
$aaa.Username
}
}


function Aggregated_Format
{
param([string]$Sort_CPU_Memory,$IsAggregated)


if($Sort_CPU_Memory -eq 'CPU'){
$Sort_By = 'PercentProcessorTime'
}
if($Sort_CPU_Memory -eq 'Memory'){
$Sort_By = 'Memory Usage(MB)'
}

if($Sort_CPU_Memory -eq 'CPU' -and $IsAggregated -eq $true)
{
write-host "######################################"
write-host "Sort BY CPU - With Aggregrate" 
write-host "######################################"
}

if($Sort_CPU_Memory -eq 'CPU' -and $IsAggregated -eq $False)
{
write-host "######################################"
write-host "Sort BY CPU - No Aggregrate" 
write-host "######################################"
}

if($Sort_CPU_Memory -eq 'Memory' -and $IsAggregated -eq $True)
{
write-host "######################################"
write-host "Sort BY Memory - With Aggregrate" 
write-host "######################################"
}

if($Sort_CPU_Memory -eq 'Memory' -and $IsAggregated -eq $False)
{
write-host "######################################"
write-host "Sort BY Memory - No Aggregrate" 
write-host "######################################"
}

$Global:FinalResult = Get-WmiObject Win32_PerfFormattedData_PerfProc_Process | ? {$_.Name -notmatch "^(idle|_total|system)$"} |select-object -property @{Name= 'Process Name';exp ={$($_.Name).split('#')[0]}}, PercentProcessorTime, IDProcess, @{"Name" = "Memory Usage(MB)"; Expression = {[int]($_.WorkingSetPrivate/1mb)}} | ? {$_.Name -notmatch "^(idle|_total|system)$"} |
Sort-Object -Property "$Sort_By" -Descending | select -first 10 | select 'Process Name', PercentProcessorTime,IDProcess, "Memory Usage(MB)"

$Global:Data = @()

foreach($FinalResult1 in $Global:FinalResult)
{

$Process_Status1 = ''
if((get-process "$($FinalResult1.'Process Name')" -ea SilentlyContinue) -eq $Null){ 
        $Process_Status1 = "Not Running"}else{ 
    $Process_Status1 = "Running"}
    
   $Global:Data += New-Object psobject -Property @{
   
      'Process Name' = $FinalResult1.'Process Name'
      'CPU Usage' = $FinalResult1.PercentProcessorTime
      PID = $FinalResult1.IDProcess
      Status = $Process_Status1
      'Memory Usage(MB)' = $FinalResult1."Memory Usage(MB)"
      Description = (Get_username_Description -PID1 $FinalResult1.IDProcess)[0]
      Username = (Get_username_Description -PID1 $FinalResult1.IDProcess)[1]
   }

}



if($IsAggregated)
{
    $grouped = $Global:Data | group 'Process Name'

    $Global:AggregrateData =@()
    foreach($grouped1 in $grouped)
    {

        $Global:AggregrateData += New-Object psobject -Property @{
                'Process Name' = $grouped1.group | select -unique -expand 'Process Name'
                'Status' = $grouped1.group | select -unique -expand 'Status'                
                'Description' = $grouped1.group | select -unique -expand 'Description'
                'Username' = $grouped1.group | select -unique -expand 'Username'
                'CPU Usage' = ($grouped1.group | measure 'CPU Usage' -sum).sum
                'Memory Usage(MB)' = ($grouped1.group | measure 'Memory Usage(MB)' -sum).sum
        } | select 'Process Name',Status,'CPU Usage','Memory Usage(MB)',Description,Username | Sort-Object -Property "$Sort_By" -Descending 
       
    }

    if($Sort_By -eq 'PercentProcessorTime')
      {  $Sort_By = 'CPU Usage' }
  $OutputData =     $Global:AggregrateData | Select 'Process Name',Status,'CPU Usage','Memory Usage(MB)',Description,Username  | sort "$Sort_By" -Descending
}
else{

$OutputData = $Global:Data | select 'Process Name',Status,PID,'CPU Usage','Memory Usage(MB)',Description,Username 
}

return $OutputData
}


Aggregated_Format -Sort_CPU_Memory "$Sort_CPU_Memory" -IsAggregated $IsAggregated | fL 
