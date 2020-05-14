<#
    .SYNOPSIS
       Retrieve data from system restore point
    .DESCRIPTION
       Retrieve data from system restore point
    .Author
       Santosh.Dakolia@continuum.net    

    .Reference URL
        # List of Chassis :: https://blogs.technet.microsoft.com/brandonlinton/2017/09/15/updated-win32_systemenclosure-chassis-types/ 
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

Try {
$computername= $env:computername

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

Write-Host "Current System '$computername' Chassis Type is :: $Show"
}Catch{
    Write-Error "Error occured while retrieving Data..!! $_.Exception.Message"
}
