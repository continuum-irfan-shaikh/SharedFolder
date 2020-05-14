$hotfixes = Get-HotFix
Write-Output "Hotfixes detected"
Write-Output $hotfixes
$CPUInfo = Get-WmiObject Win32_Processor
$OSTotalVirtualMemory = [math]::round($OSInfo.TotalVirtualMemorySize / 1MB, 2)
$OSTotalVisibleMemory = [math]::round(($OSInfo.TotalVisibleMemorySize  / 1MB), 2)
$PhysicalMemory = Get-WmiObject CIM_PhysicalMemory | Measure-Object -Property capacity -Sum | % {[math]::round(($_.sum / 1GB),2)}

Write-Output "Hardware detected"
$infoObject = New-Object PSObject
Add-Member -inputObject $infoObject -memberType NoteProperty -name "CPU_Name" -value $CPUInfo.Name
Add-Member -inputObject $infoObject -memberType NoteProperty -name "CPU_Cores" -value $CPUInfo.NumberOfCores
Add-Member -inputObject $infoObject -memberType NoteProperty -name "TotalPhysical_Memory_GB" -value $PhysicalMemory
Add-Member -inputObject $infoObject -memberType NoteProperty -name "TotalVirtual_Memory_MB" -value $OSTotalVirtualMemory
Add-Member -inputObject $infoObject -memberType NoteProperty -name "TotalVisable_Memory_MB" -value $OSTotalVisibleMemory
Write-Output $infoObject

Write-Output "Software detected"
$OSInfo = Get-WmiObject Win32_OperatingSystem #Get OS Information
Write-Output $OSInfo
Write-Output "OS Name: $($OSInfo.Name)"
$SoftInfo = Get-WmiObject -Class Win32_Product
Write-Output $SoftInfo
