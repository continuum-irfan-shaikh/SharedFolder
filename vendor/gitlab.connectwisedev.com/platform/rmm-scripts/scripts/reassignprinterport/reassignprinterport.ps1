<#
.SYNOPSIS
Developed By : GRT [HD-Automation], Continuum
File Version : 1.0

.DESCRIPTION
The Script is designed to reassign the port to the printer if port not in use.

.PARAMETER
Printer_Name is the string variable required for activity on printer.

.PARAMETER
Port_Name is the String Variable used for reassigning to the Printer.

.PARAMETER
Link_Status is a boolean. True Means to Link the Port to the Printer.

.Example
$Printer_Name="XXXXXXX"
$Port_Name="XXXXXXXX"
$Link_Status=$True
#>
$Executionlog=@()
[double]$OSVersion=[Environment]::OSVersion.Version.ToString(2)
#OS Version Comparison.
if ($OSVersion -lt '6.1')
     {
        $Executionlog += "Script is Design for Windows 7 and above OS only, Script Execition Stopped."
        Write-Output $Executionlog
        exit;
     }
# PS Version Check.
if ($PSVersionTable.PSVersion.Major -lt '2')
     {
        $Executionlog+= Write-Output "SCript is Design for Powershell version 2.0 and Above, Terminating Script for further Execution." 
        Write-Output $Executionlog
        Exit;
     }
#Check Status of Print Spooler Service and start if not running.
$Spooler_State=Get-WmiObject -Class win32_service -Filter "Name='Spooler'" | select Name,StartMode,State
if ($Spooler_State.name -notcontains "Spooler")
     {
        $Executionlog += Write-Output "Spooler Service Not Available, Stopping Script For Further Execution."
        Write-Output $Executionlog
        Exit;
     }
#Check Status of Print Spooler Service and start if not Running.
if ($Spooler_State.state -ne "Running")
     {
        try {
                Write-Output "Print Spooler Service Not Started, Starting Service."
                Get-WmiObject Win32_Service -Filter "Name='Spooler'" | Set-Service -StartupType Automatic
                Start-Service -Name "Spooler" -ErrorAction Stop
            }
      catch {
                 $Executionlog += Write-Output "Caught An Exeception; Print Spooler Service Faild to Start."
                 $Executionlog += $_
                 Write-Output $Executionlog
                 Exit;
            }
     }
# Get list of available TCP IP Port and Compare with Port Name Provided in Variable. If not Available list down the list of Available Ports.
$ListTCPIPPort=Get-WmiObject -Class win32_TCPIPPrinterPort | select -ExpandProperty Name
if ($ListTCPIPPort -notcontains $Port_Name)
     {
            $Executionlog += Write-Output "Port Name Entered Not Available in the List of Available TCP\IP Ports Or None of the TCP\IP Port Available in System.
            `nList Of Available TCP\IP Printer Ports is as below.(Note: if Output Blank then No TCP\IP Port Available in System.)"

            $Executionlog += $ListTCPIPPort
            Write-Output $Executionlog |fl
            Exit;
     }
#Query Specific Printer And Check Assign PortName to the Printer.
$PrinterConfig=Get-WmiObject win32_Printer | ? {$_.name -eq "$Printer_Name"}
if ($PrinterConfig.PortName -eq $Port_Name)
     {
            $Executionlog += Write-Output "Port Name Already Linked to the Printer Named $printer_Name, Script Execution Stopped."
            Write-Output $Executionlog |fl
            Exit;   
     }
#if List Status is Set to False.
if (-not($link_Status))
     {
            $Executionlog += Write-Output "No Changes has been made, Because 'Link_Status' Set to False"
            Write-Output $Executionlog
            Exit;
     }
#Query Specific Printer And Check if Printer is Available.
$AllPrinterName=Get-WmiObject win32_Printer | select -ExpandProperty Name
if ($AllPrinterName -notcontains $Printer_Name)
     {
            $Executionlog += Write-Output "$Printer_Name Printer not available in the list of Installed Printer.
Script Execution Stopped. List Of Available Printers are as below. 
(Note: If Output is Blank then most probabilty no Printer available.)"`n
            $Executionlog += $AllPrinterName
            Write-Output $Executionlog |fl
            Exit;   
     }
# if Link Status Set to True.
if ($link_Status)
     {
            try {
                   $Port_Change=gwmi win32_printer -filter "Name='$Printer_Name'"
                   $Port_Change.PortName=$Port_Name
                   $Port_Change.Put()
                }
           Catch{
                   $Executionlog += Write-Output "Caught An Exception While Assigning Port to Printer. Port Assignment Failed."
                   $Executionlog += $_
                   Write-Output $Executionlog
                   Exit;
                }
     }
#Remove PrinterConfig Variable to Get Fresh Printer Config Information.
Remove-Variable PrinterConfig
$PrinterConfig=Get-WmiObject Win32_Printer | Where-Object {$_.name -eq "$Printer_Name"}
#Status After Configuration Change.
if ($PrinterConfig.PortName -eq $Port_Name)
     {
          $Executionlog += Write-Output "Port Assignment to Printer Completed Successfully."
          Write-Output $Executionlog
          Exit;
     }
