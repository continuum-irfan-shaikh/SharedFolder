<#
    .SYNOPSIS
       Retrieve battery information
    .DESCRIPTION
       Retrieve battery information
    .Author
       Santosh.Dakolia@continuum.net
    .Reference 
        https://docs.microsoft.com/en-us/windows/desktop/cimwin32prov/win32-battery    
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

$computername= $env:computername


Try {

$Chasis = (Get-WmiObject win32_systemenclosure -ComputerName $computername).chassistypes

switch ($chasis) 
{
1 {$Show = "Other" }
2 {$Show = "Unknown" }
3 {$Show = "Desktop" }
4 {$Show = "Low Profile Desktop" }
5 {$Show = "Pizza Box" }
6 {$Show = "Mini Tower" }
7 {$Show = "Tower" }
8 {$Show = "Portable" }
9 {$Show = "Laptop" }
10 {$Show = "Notebook" }
11 {$Show = "Hand Held" }
12 {$Show = "Docking Station" }
13 {$Show = "All in One" }
14 {$Show = "Sub Notebook" }
15 {$Show = "Space-Saving" }
16 {$Show = "Lunch Box" }
17 {$Show = "Main System Chassis" }
18 {$Show = "Expansion Chassis" }
19 {$Show = "SubChassis" }
20 {$Show = "Bus Expansion Chassis" }
21 {$Show = "Peripheral Chassis" }
22 {$Show = "RAID Chassis" }
23 {$Show = "Rack Mount Chassis" }
24 {$Show = "Sealed-case PC" }
25 {$Show = "Multi-system chassis" }
26 {$Show = "Compact PCI" }
27 {$Show = "Advanced TCA" }
28 {$Show = "Blade" }
29 {$Show = "Blade Enclosure" }
30 {$Show = "Tablet" }
31 {$Show = "Convertible" }
32 {$Show = "Detachable" }
33 {$Show = "IoT Gateway" }
34 {$Show = "Embedded PC" }
35 {$Show = "Mini PC" }
36 {$Show = "Stick PC" }
 }

#Write-Host "Current System '$computername' Chassis Type is :: $Show"
}Catch{
    Write-Error "Error occured while retrieving Data..!! $_.Exception.Message"
}

Try{

$B = (Get-WmiObject -Class win32_battery)

IF (!$B) {"Current System '$computername' Chassis Type is : $Show Hence Unable to Fetch Battery Information"}
IF ($B) {

$batterydetails = get-wmiobject -class "Win32_Battery" -ComputerName $computername
$BatteryA = (get-wmiobject -class "Win32_Battery").Availability
switch ( $BatteryA )
    {
1 {$BatteryAvailability =  "Other " }
2 {$BatteryAvailability =  "Unknown " }
3 {$BatteryAvailability =  "Running or Full Power " }
4 {$BatteryAvailability =  "Warning " }
5 {$BatteryAvailability =  "In Test " }
6 {$BatteryAvailability =  "Not Applicable " }
7 {$BatteryAvailability =  "Power Off " }
8 {$BatteryAvailability =  "Off Line " }
9 {$BatteryAvailability =  "Off Duty " }
10 {$BatteryAvailability =  "Degraded " }
11 {$BatteryAvailability =  "Not Installed " }
12 {$BatteryAvailability =  "Install Error " }
13 {$BatteryAvailability =  "The device is known to be in a power save mode, but its exact status is unknown." }
14 {$BatteryAvailability =  "The device is in a power save state but still functioning, and may exhibit degraded performance " }
15 {$BatteryAvailability =  "The device is not functioning, but could be brought to full power quickly " }
16 {$BatteryAvailability =  "Power Cycle " }
17 {$BatteryAvailability =  "The device is in a warning state, though also in a power save mode" }
18 {$BatteryAvailability =  "The device is paused" }
19 {$BatteryAvailability =  "The device is not ready" }
20 {$BatteryAvailability =  "The device is not configured " }
21 {$BatteryAvailability =  "The device is quiet" }
    }
    
$BatteryS = (get-wmiobject -class "Win32_Battery" -namespace "root\CIMV2" -ComputerName $computername).BatteryStatus
switch ( $BatteryS )
    {
1 {$BatteryStatus =  "The battery is discharging" }
2 {$BatteryStatus =  "The system has access to AC so no battery is being discharged. However, the battery is not necessarily charging" }
3 {$BatteryStatus =  "Fully Charged" }
4 {$BatteryStatus =  "Low" }
5 {$BatteryStatus =  "Critical" }
6 {$BatteryStatus =  "Charging" }
7 {$BatteryStatus =  "Charging and High" }
8 {$BatteryStatus =  "Charging and Low" }
9 {$BatteryStatus =  "Charging and Critical" }
10 {$BatteryStatus =  "Undefined" }
11 {$BatteryStatus =  "Partially Charged" }
    }

$BatteryC = (Get-WmiObject -Class "Win32_Battery" -ComputerName $computername).chemistry
Switch ($BatteryC){
1 {$BatteryChemistry =  "Other"}
2 {$BatteryChemistry =  "Unknown"}
3 {$BatteryChemistry =  "Lead Acid"}
4 {$BatteryChemistry =  "Nickel Cadmium"}
5 {$BatteryChemistry =  "Nickel Metal Hydride"}
6 {$BatteryChemistry =  "Zinc air"}
7 {$BatteryChemistry =  "Lithium Polymer"}
}


    $time = new-timespan -minutes $batterydetails.EstimatedRunTime
    $t = $time.Hours 
    $s = $time.Minutes
    $FTime = Write-Output "$t Hours $s Minutes"

    $OutputObj = New-Object -TypeName PSobject

    $outputobj | Add-Member -MemberType NoteProperty -Name Availability -Value $BatteryAvailability
    $outputobj | Add-Member -MemberType NoteProperty -Name Status -Value $batterystatus
    $OutputObj | Add-Member -MemberType NoteProperty -Name Chemistry -Value $BatteryChemistry
    $OutputObj | Add-Member -MemberType NoteProperty -Name Description -Value $batterydetails.Description
    $OutputObj | Add-Member -MemberType NoteProperty -Name DesignVoltage -Value $batterydetails.DesignVoltage
    $outputobj | Add-Member -MemberType NoteProperty -Name DeviceID -Value $batterydetails.deviceid
    $OutputObj | Add-Member -MemberType NoteProperty -Name EstimatedRunTime -Value $FTime
    $outputobj | Add-Member -MemberType NoteProperty -Name Name -Value $batterydetails.name
    $OutputObj | Add-Member -MemberType NoteProperty -Name PowerManagementSupported -Value $batterydetails.PowerManagementSupported
    $outputobj | fl Availability, Status, Chemistry, Description, DesignVoltage, DeviceID, EstimatedRunTime, Name, PowerManagementSupported

}
}Catch{
    Write-Error "Error occured while retrieving Data..!! $_.Exception.Message"
}
