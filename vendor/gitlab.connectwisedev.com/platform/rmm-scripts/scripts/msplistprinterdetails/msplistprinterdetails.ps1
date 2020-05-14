<#
.SYNOPSIS
Developed By : GRT [HD-Automation] , Continuum
File Vesion : 2.0
.DESCRIPTION
The script is designed to pull printer details[Name, DriverName, Queued, PrinterStatus, PrinterState,
Local, Network, Shared, ShareName, ServerName, PortName, WorkOffline] installed on Per Machine .
.PARAMETER
No Parameter needed .
#>
$serviceName = "Spooler"
$timeoutSeconds = 60
$ServiceStatus =$FALSE
$OutputMSG = @()
$hostAddresses = @{}
$PrinterStatusCode = @{"1"="Other"; "2"="Unknown" ; "3"="Idle"; "4" = "Printing"; "5" = "Warmup"; "6" = "Stopped Printing"; "7" = "Offline"}
try{ $service = Get-Service $serviceName
if ( -not $service )
{     OutputMSG += "SPOOLER_SERVICE_DOES_NOT_EXIST"
Write-Output $OutputMSG
Exit   }
if ( $service.Status -eq [ServiceProcess.ServiceControllerStatus]::Running )
{     $OutputMSG += "SPOOLER_SERVICE_ALREADY_RUNNING"  }
else{   $timeSpan = New-Object Timespan 0,0,$timeoutSeconds
try {
    $service.Start()
    $service.WaitForStatus([ServiceProcess.ServiceControllerStatus]::Running, $timeSpan)
    $OutputMSG += "Spooler Service Started [Was not running]."
} catch [Management.Automation.MethodInvocationException],[ServiceProcess.TimeoutException] {
    $OutputMSG += "INFORMATION_SERVICE_REQUEST_TIMEOUT[60 SEC]"
}catch{   Write-Error "ERROR_SPOOLER_SERVICE_ISSUE $_.Exception.message"
Exit   }
}
#WMI Query to collect TCPIPPrinterPort infomration to store in Hash Table
Get-WmiObject Win32_TCPIPPrinterPort | ForEach-Object {$hostAddresses.Add($_.Name, $_.HostAddress)}
$PrinterList = Get-WmiObject -class "Win32_Printer" -namespace "root\CIMV2" | ForEach-Object {
    New-Object PSObject -Property @{
        "Name" = $_.Name
        "DriverName" = $_.DriverName
        "Queued" = $_.Queued
        "PrinterStatus" = $PrinterStatusCode[[String]$_.PrinterStatus]
        "PrinterState" = $_.PrinterState
        "Local" = $_.Local
        "Network" = $_.Network
        "Shared" = $_.Shared
        "ShareName" = $_.ShareName
        "ServerName" = $_.ServerName
        "PortName" = $_.PortName
        "HostAddress" = $hostAddresses[$_.PortName]
        "WorkOffline" = $_.WorkOffline
    }
}| Format-List -Property Name, DriverName, Queued, PrinterStatus, PrinterState, Local, Network, Shared, ShareName, ServerName, PortName, HostAddress, WorkOffline

$OutputMSG += $PrinterList|fl
if($PrinterList)
{   Write-Output $OutputMSG |fl }
else
{ Write-Output "No Printer in List"}
}catch { Write-error "Error in retrieving the printer information. $_.exception.message" }
