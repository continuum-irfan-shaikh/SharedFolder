<#
.SYNOPSIS
Developed By : GRT [HD-Automation] , Continuum
File Vesion : 1.0

.DESCRIPTION
The script is designed to delete the provided printer driver from the registry.
.PARAMETER 
Printer_Driver_Name is the variable required for deleting the entry from the registry.

.Example
$Printer_Driver_Name="HP LaserJet 1022n Class Driver,4,Windows x64"

#>
$executionlog=@()

[double]$OSVersion=[Environment]::OSVersion.Version.ToString(2)
# OS Version comparison.
if ($osversion -lt '6.1')
        {
            $executionlog += 'Script is design for windows 7 and obove OS only, Script Execution Stopped.'
            Write-Output $executionlog
            exit;
        } 
# PS Version check.   
if ($PSVersionTable.PSVersion.Major -lt '2') 
        {
            $executionlog += write-output "Powershell version running is lower version, terminating script for furter execution."
            Write-Output $executionlog
            exit;
        }      
# checking print spooler service installed state.  
$spooler_State=Get-WmiObject -Class Win32_Service -Filter "name='Spooler'" | select Name,Startmode,State 
        if ($spooler_State.name -notcontains "Spooler")             
            {
                   $executionlog += write-output "Spooler service not available, Stopping Script for further execution"
                   Write-Output $executionlog
                   exit;
            }   
# check status of print spooler service and start if not running start service..     
 if ($spooler_State.State -ne "Running") 
        {
           try {
                  write-output "Print Spooler Service Not Started, Starting Service."
                  Get-WmiObject win32_service -Filter "name='spooler'" | Set-Service -StartupType Automatic
                  Start-Service -Name "Spooler" -ErrorAction Stop
               }
         catch {               
                  $executionlog +=  Write-Output "Caught an exception:, Service Failed to start."                       
                  $executionlog+=$_
                  Write-Output $executionlog
                  exit;
               }
          }    
#getting list of installed drivers and compare with Printer name supplied, if not available create hash table to disply driver name and driver version info.
$driver_Name= Get-WmiObject -Class win32_printerdriver -ComputerName $env:COMPUTERNAME | select -ExpandProperty Name
if ($driver_Name -notcontains $Printer_Driver_Name) 
     {            
         $executionlog += write-output "$Printer_Driver_Name not Available in the list of installed printer drivers. List of Availble printer Drivers As below. Script Execution Stopped."
         $Output = @()
         ForEach ($Driver in (Get-WmiObject Win32_PrinterDriver -ComputerName $env:COMPUTERNAME))
            {	 
               $Drive = $Driver.DriverPath
	           $Output += New-Object PSObject -Property @{
		       'Printer Driver Name' = $Driver.Name
               'Printer Driver Version' = ((Get-Item $Drive).VersionInfo.ProductVersion) }
	        }                    
                $executionlog += ($Output | select 'Printer Driver Name','Printer Driver Version')
                Write-Output $executionlog |fl 
                Exit;                                
     }
# Stop script if printer driver is already in used by any other printer.
$Used_Driver=Get-WmiObject -Class win32_printer | select -ExpandProperty DriverName
    if ($used_driver -contains ($Printer_Driver_Name -split ",")[0])
            {         
                $executionlog += write-output "Printer Driver is in use, Script Execution Stopped, list of available printers and associate Driver are as below."
                $executionlog += (Get-WmiObject -Class win32_printer | select @{l='Printer Name';e={$_.name}},DriverName)
                Write-Output $executionlog | fl
                exit;
            }
# Delete Printer Driver, if it's not used by any other Printers. 
if ($used_driver -notcontains ($Printer_Driver_Name -split ",")[0])      
   {
         try
              {
                RUNDLL32 PRINTUI.DLL,PrintUIEntry /dd /m ($Printer_Driver_Name -split ",")[0]
              }
        catch {
                  $executionlog +=  Write-Output "Caught an exception, Printer Driver deletion Might Failed."                       
                  $executionlog += $_
                  Write-Output $executionlog
                  exit;
              }
   }
# Validate change after Printer Driver Deletion.
    if ((Get-WmiObject -Class win32_printerdriver).name -notcontains ($Printer_Driver_Name -split ",")[0])
  {
        $executionlog += write-output "$Printer_Driver_Name, Printer Driver Deleted. List of Available Printer Drivers are as below."
        sleep 3
        $Output = @()
        ForEach ($Driver in (Get-WmiObject Win32_PrinterDriver -ComputerName $env:COMPUTERNAME))
            {           
               $Drive = $Driver.DriverPath
	           $Output += New-Object PSObject -Property @{
		       'Printer Driver Name' = $Driver.Name
               'Printer Driver Version' = ((Get-Item $Drive).VersionInfo.ProductVersion) }	                
            }                    
               $executionlog += $Output | select 'Printer Driver Name','Printer Driver Version'
               Write-Output $executionlog |fl
               exit; 
  }
 
