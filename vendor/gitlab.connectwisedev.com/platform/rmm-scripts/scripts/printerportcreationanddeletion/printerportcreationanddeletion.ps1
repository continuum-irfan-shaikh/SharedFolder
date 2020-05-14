<#
.SYNOPSIS
Developed By : GRT [HD-Automation] , Continuum
File Vesion : 1.0

.DESCRIPTION
TCP/IP Printer Port Creation and Deletion.

.PARAMETER Name
$PrinterPortName Specify Priniter Port Name which you want to Create or Delete.

.PARAMETER Name
PrinterIP Specify Printer IP, Used while creation of Printer Port.


.PARAMETER Name
$SNMPEnabled SNMP Status for the Port Created.

.PARAMETER Name
$SNMPString SNMP String Used by Port.

.PARAMETER Name
$Port_Deletion If Port Deletion is True Than only Port Deletion Starts.

.PARAMETER Name
$PortCreation If port Creation is True Then only Port Creation Starts.

Note:

1) $PrinterPortName, $PrinterIP is Mandatory Paramenter for creating Printer Port.
2) For Printer port deletion $PrinterPortName is Mandatory Paramenter.
3) Make sure only one Paramenter Set to True among $Port_Deletion and $PortCreation. if both set to true, By Default it will Execute Port Deletion only.

.EXAMPLE

[string]$PrinterPortName="NewPrinterPort"
[string]$PrinterIP="10.137.2.15"
[bool]$SNMPEnabled=$false
[string]$SNMPString="public"
[bool]$Port_Deletion=$False
[bool]$PortCreation=$True

#>

$outlog=@()
[double]$OSVersion=[Environment]::OSVersion.Version.ToString(2)
# OS Version comparison.
if ($osversion -lt '6.1')
        {
            $outlog += 'Script is design for windows 7 and obove OS only, Script Execution Stopped.'
            Write-Output $outlog
            exit;
        } 
# PS Version check.   
if ($PSVersionTable.PSVersion.Major -lt '2') 
        {
            $outlog += write-output "Powershell version running is lower version, terminating script for furter execution."
            Write-Output $outlog
            exit;
        }      
# checking print spooler service installed state.  
$spooler_State=Get-WmiObject -Class Win32_Service -Filter "name='Spooler'" | select Name,Startmode,State 
        if ($spooler_State.name -notcontains "Spooler")             
            {
                   $outlog += write-output "Spooler service not available, Stopping Script for further execution"
                   Write-Output $outlog
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
                  $outlog +=  Write-Output "Caught an exception:, Service Failed to start."                       
                  $outlog+=$_
                  Write-Output $outlog
                  exit;
               }
          }  
# if Port Deletion status set to True, execute following Loop.     
if ($Port_Deletion)
{
        try   {
# if provided port name not available in System.
              $listport=Get-WmiObject -class win32_tcpipprinterport | select -ExpandProperty Name
              
              if ($listport -notcontains $PrinterPortName) {
                   $outlog += Write-Output "$PrinterPortName TCP/IP Printer port does not exist, list of Available TCP/IP Ports is as below.`n"
                   $outlog += Write-Output ($listport -join (", "))
                   Write-Output $outlog |fl
                   exit;
                  }
# if provided port name already mapped with printer.
$mappedports=Get-WmiObject -Class win32_printer -Filter "portname='$PrinterPortName'"

              if ($mappedports.portname -contains $PrinterPortName) 
              {                 
                   $outlog += "$PrinterPortName is in Mapped with Device\Printer. List of Printer and Associcate ports are as below."
                   $data = (Get-WmiObject win32_printer | select @{l="Priter Name";e={$_.Name}},PortName)
                   $outlog += $data
                   Write-Output $outlog |fl
                   exit;
              }
               # Delete Port Name.  
                        $port = [wmiclass]"Win32_TcpIpPrinterPort" 
                        $port.psbase.scope.options.EnablePrivileges = $true 
                        $delete_port= $port.CreateInstance() 
                        $delete_port.name = $PrinterPortName
                        $delete_port.delete()
               }

       catch    {$outlog += "Cought Exeception While Deleting Printer Port."
                 $outlog += $_
                 Write-Output $outlog
                }
# Varify port deletion status after changes and disply change.
                 $mappedports=Get-WmiObject -Class win32_printer | select -ExpandProperty PortName
           if ($mappedports -notcontains $PrinterPortName)
                {                
                    $listport=Get-WmiObject -class win32_tcpipprinterport | select -ExpandProperty Name
                    $outlog += "$PrinterPortName Port Deleted, revised list of Available TCP\IP ports as below.`n"                    
                    $outlog +=  $listport
                    Write-Output $outlog |fl
                    exit;
                }
}
# If Port Creation is true.
if ($PortCreation) 
{
            $listport=Get-WmiObject -class win32_tcpipprinterport | select -ExpandProperty Name
# if provided port name already available.
            if ($listport -contains $PrinterPortName)
        {
            $outlog += "$PrinterPortName Printer Port Name Already exist, list of Printer TCP/IP Ports is as below.`n"
            $outlog += $listport -join (", ")
            Write-Output $outlog |fl
            exit;
        }
#create printer port
    else {
            try {   
                    $port = [wmiclass]"Win32_TcpIpPrinterPort" 
                    $port.psbase.scope.options.EnablePrivileges = $true 
                    $newPort = $port.CreateInstance() 
                    $newport.name = $PrinterPortName
                    $newport.Protocol = 1 
                    $newport.HostAddress = $PrinterIP 
                    $newport.PortNumber = "9100" 
                    $newport.SnmpEnabled = $SNMPEnabled
                    $newPort.SNMPCommunity=$SNMPString
                    $newport.Put()
                }
           catch {
                    $outlog += "Error While Creating Printer Port."
                    $outlog += $_
                    Write-Output $outlog                        
                 }
#Verify Port Status After creating Printer Port.

            $listport=Get-WmiObject -class win32_tcpipprinterport | select -ExpandProperty name
           
                  if ($listport -contains $PrinterPortName)
                     {
                        $outlog += Write-Output "$PrinterPortName Port Created, revised list of Available Printer and TCP\IP ports as below."`n
                        $outlog += $listport -join (", ")
                        Write-Output $outlog |fl
                        exit;         
                     }             
         }
}
