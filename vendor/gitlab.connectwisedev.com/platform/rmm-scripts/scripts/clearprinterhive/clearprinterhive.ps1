###################################################
###-------Clear Printer Hive ---------------########
####################################################
 
<# 
 $bAllClearHive = $false
 $PrinterName = 'HP LaserJet Pro MFP M127-M128 PCLmS1'
 $RegBackupPathName ='C:\Temp2'
 $RegBackupFileName ='Backup'
 #>
 $timeoutSeconds = 30
 $ComputerName = $env:computername
 $RegPrintPath = 'HKLM:\SYSTEM\CurrentControlSet\Control\Print'
 $registry = 'HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Print'

##########################################################
###------Check Pre Condition------------------------------
##########################################################

Function Check-PreCondition{

    $IsContinued = $true
    Write-Host "Checking Preconditions started..."
    Write-Host "-------------------------------"
    Write-Host "    " 
   
    #####################################
    # Step 1
    # Verify PowerShell Version
    #####################################

    if(-NOT($PSVersionTable.PSVersion.Major -ge 2)){
      
       $IsContinued = $false
       Write-Warning "Powershell version below 2.0 is not supported"
       return $false

    }else{
             
         write-host -ForegroundColor 10 "PowerShell Version : $($PSVersionTable.PSVersion.Major)" 
         
    }
    
    ####################################     
    # Step 2
    # Verify operating system Version
    ####################################

    if( $IsContinued -and -not([System.Environment]::OSVersion.Version.major -ge 6)){

      $IsContinued = $false
      Write-Warning "Powershell Script suppoerts Window 7, Window 2008R2 and higher version operating system"
      return $false

    }else{
              
        write-host -ForegroundColor 10 "Operating Syatem Version : $([System.Environment]::OSVersion.Version.major)" 
       
    }
    
     ##################################
     # Step 3
     # Verify StartMode if Disabled
     ##################################

     if( $IsContinued -and (Get-WmiObject -Class Win32_Service -Property StartMode -Filter "Name='Spooler'").StartMode -eq "Disabled"){
        
        $IsContinued = $false
        write-host "Service Status :" (Get-Service | Where {$_.name -eq 'Spooler'}).Status
        Write-Warning "Spooler Service is disabled, Please enable it." 
        return $false

     }else{

         write-host -ForegroundColor 10 "Spooler Service StartupType :" (Get-WmiObject -Class Win32_Service -Property StartMode -Filter "Name='Spooler'").StartMode 
         
     }  
     
     ##################################
     # Step 4
     # Get dependent Services  
     ##################################

     CheckServiceStatus

     if($IsContinued){
      
        Get-Service -CN . | Where-Object {$_.name -eq 'Spooler'} | ForEach-Object { 
        
              write-host -ForegroundColor 10 "Service name : $($_.name)" 
              
              if($_.DependentServices)  { 
              
                  write-host -ForegroundColor 3 "`tServices that depend on $($_.name)" 
                  foreach($s in $_.DependentServices){
                  Write-Host "`t`t" + $s.name 
                   } 
                } #end if DependentServices 
                
              if($_.RequiredServices) {
              
                  Write-host -ForegroundColor 3 "`tServices required by $($_.name)" 
                  foreach($r in $_.RequiredServices) {

                   Write-Host "`t`t" + $r.name 
                        
                        if( -not((Get-Service | Where {$_.name -eq $r.name}).Status -eq "Running")){ 
      
                            Get-Service $r.name | Where {$_.Status -eq 'Stopped'} | Start-Service                             
                            Write-Host -ForegroundColor 8 "`t`t $($r.name) Service Started"
                        }
                        else{
        
                            Write-Host -ForegroundColor 10 "`t`t $($r.name) Service is running"

                        }

                   } 
                } #end if DependentServices 
            } #end foreach-object
     
     }
    
    ##########################################
    # Step 5
    # Verify Spooler Service running status
    ##########################################

    if($IsContinued -and -not((Get-Service | Where {$_.name -eq 'Spooler'}).Status -eq "Running")){ 
      
        Get-Service Spooler | Where {$_.Status -eq 'Stopped'} | Start-Service  
        
        Write-Host "    "
        Write-Host "Spooler Service Started"
    }
    else{
        
         Write-Host "    "
        Write-Host -ForegroundColor 10 "Spooler Service is running"

    }


    
    Write-Host "    "
    write-host -ForegroundColor 7 "`t Checking Precondition Completed"

    return $true
}

##########################################################
###------Get Printer List---------------------------------
##########################################################

Function GetPrinters([String] $PrinterName){

    Write-Host "-------------------------------"
    write-host "Retriving printers"
    Write-Host ""

    try
    {
        $PrinterCount =0

        if($bAllClearHive -eq $true)
        {
            Get-WMIObject -Class Win32_Printer -Computer $computerName | ForEach-Object { 
                write-host -ForegroundColor 8 "Printer name : $($_.name)" 
                #write-host -ForegroundColor 10 "Driver : $($_.DriverName)"

                $PrinterCount = $PrinterCount +1
            }  
        }
        else
        {
            if ((Get-WMIObject -Class Win32_Printer -Computer $computerName | Where-Object { $_.name -eq $PrinterName }) -eq $null) 
            {
                Get-WMIObject -Class Win32_Printer -Computer $computerName | ForEach-Object { 
                    write-host -ForegroundColor 8 "Printer name : $($_.name)" 
                    #write-host -ForegroundColor 10 "Driver : $($_.DriverName)"

                    $PrinterCount = $PrinterCount +1
                } 
            }
            else 
            {
                Get-WMIObject -Class Win32_Printer -Computer $computerName | Where-Object { $_.name -eq $PrinterName } | ForEach-Object { 
                    write-host -ForegroundColor 8 "Printer name : $($_.name)" 
                    #write-host -ForegroundColor 10 "Driver : $($_.DriverName)"

                    $PrinterCount = $PrinterCount +1
                } 
             }
        }

        #---------
        if($PrinterCount -eq 0)
         {
            Write-Host ""
            write-host -ForegroundColor 10 "Found no Printer to Delete"
            Write-Host ""

            return $flase
         }
         else
         {
            return $true
         }
    }
    catch
    {
        Write-Warning "Error while searching printer: $_.exception.message"
        return $flase
    }
}

Function GetPrinters_V8([String] $PrinterName){

    Write-Host "-------------------------------"
    write-host "Retriving printers"
    Write-Host ""

    try
    {
        if($bAllClearHive -eq $true)
        {
            Get-Printer |ForEach-Object {
                write-host -ForegroundColor 8 "Name : $($_.Name)" 
            }
        }
        else
        {
            $SelPrinter = Get-Printer -Name $PrinterName

            if($SelPrinter -eq $null)
            {
                Get-Printer |ForEach-Object {
                    write-host -ForegroundColor 8 "Name : $($_.Name)" 
                }
            }
            else
            {
                write-host -ForegroundColor 8 "Printer name : $PrinterName"
            }
        }
    }
    catch
    {
        Write-Warning "Error while searching printer: $_.exception.message"
        return $flase
    }
}

##########################################################
###------Take Backup of Registry--------------------------
##########################################################

Function BackupRegistryPath(){

    param (
        [parameter(Mandatory=$true)]
        [ValidateNotNullOrEmpty()]$Path,
        [parameter(Mandatory=$true)]
        [ValidateNotNullOrEmpty()]$BackupPath,
        [parameter(Mandatory=$true)]
        [ValidateNotNullOrEmpty()]$BackupFileName
        )

    try
    {
        Write-Host "-------------------------------"
        Write-Host "Checking Registry: $Path"

        if(Test-Path $Path)
        {
            write-host -ForegroundColor 7 "`t Registry Exists: true"

	    write-host -ForegroundColor 7 "`t Backup Path: $($BackupPath)"

            #New-Item -ItemType Directory -Force -Path $BackupPath -errorAction Continue

            if(Test-Path $BackupPath)
            {
                $FilePath= $BackupPath + "\"+ $BackupFileName + "_"+ (get-date -Format dd-MM-y-hh-mm-ss) + ".reg" 

                write-host -ForegroundColor 7 "`t Creating Backup File : $FilePath"

			    Get-ChildItem $Path -recurse |Export-Clixml $FilePath

                write-host -ForegroundColor 7 "`t Registry Backup Completed"
                return $true
            }
            else
            {
                Write-Warning "Registry backup Path not be created" 
                return $flase
            }
        }
        else
        {
            Write-Warning "Registry Path not Exists: $Path" 
            return $flase
        }
    }
    catch
    {
        Write-Warning "Error while taking registry backup: $_.exception.message"
        return $flase
    }
}

Function BackupRegistryPath_V8(){

    param (
        [parameter(Mandatory=$true)]
        [ValidateNotNullOrEmpty()]$Path,
        [parameter(Mandatory=$true)]
        [ValidateNotNullOrEmpty()]$BackupPath,
        [parameter(Mandatory=$true)]
        [ValidateNotNullOrEmpty()]$BackupFileName
        )

    try
    {
        Write-Host "-------------------------------"
        Write-Host "Checking Registry: $Path"

        if(Test-Path $Path)
        {
            write-host -ForegroundColor 7 "`t Registry Exists: true"

            write-host ""
            write-host -ForegroundColor 7 "`t Backup Path: $($BackupPath)"

            if(Test-Path $BackupPath)
            {
                write-host -ForegroundColor 7 "`t Backup Path Exists: true"
            }
            else
            {
                write-host -ForegroundColor 7 "`t Creating Backup Path : $BackupPath"
                New-Item -ItemType Directory -Force -Path $BackupPath -errorAction Stop
            }


            if(Test-Path $BackupPath)
            {
                $FilePath= $BackupPath + "\"+ $BackupFileName + "_"+ (get-date -Format dd-MM-y-hh-mm-ss) + ".reg" 

                write-host -ForegroundColor 7 "`t Creating Backup File : $FilePath"

                REG EXPORT $registry $FilePath

                write-host ""
                write-host -ForegroundColor 7 "`t Registry Backup Completed"
                return $true
            }
            else
            {
                Write-Warning "Registry backup Path not be created" 
                return $flase
            }
        }
        else
        {
            Write-Warning "Registry Path not Exists: $Path" 
            return $flase
        }
    }
    catch
    {
        Write-Warning "Error while taking registry backup: $_.exception.message"
        return $flase
    }
}

##########################################################
###------Delete Printers ---------------------------------
##########################################################

Function DeletePrinters(){

        $foundPrinter = $false

        Write-Host "-------------------------------"
        Write-Host "Deleting Printer, Driver, Port "
        Write-Host ""

        $PrinterCount=0
         
        get-wmiobject -class "win32_printer" -Computer $computerName | ForEach-Object { 
            
             #write-host "---------------------------------"
             #write-host -ForegroundColor 10 "Printer Name : $($_.name)" 
             #write-host -ForegroundColor 8 "Port Name : $($_.PortName)" 
             #write-host -ForegroundColor 8 "Deleting Driver: $($_.DriverName)"
             #write-host -ForegroundColor 8 "Print Processor : $($_.PrintProcessor)"
             #write-host -ForegroundColor 8 "Spool Enabled: $($_.SpoolEnabled)" 

            $PrinterDriverName = $_.DriverName

            if($bAllClearHive -eq $true)
            {
                write-host "  ---------------------------------"
                write-host -ForegroundColor 10 "  Printer Name : $($_.name)"
                write-host ""
                 
                if($_.PrintProcessor -ne 'winprint')
                {
                    RemoveThirdPartyPrintProcessorAssociation $_.PrintProcessor $_.name
                }

                $PrinterCount = $PrinterCount+1

                #----Deleting Printer----------------------
                try
                {
                    write-host -ForegroundColor 8 "`t Deleting Printer : $($_.name)" 
                    $SelPortName = $_.PortName

                    $_.delete()
                    VerifyPrinterDeleted($_.name)
                    #write-host -ForegroundColor 7 "`t`t Printer Deleted : true"
                }
                catch
                {
                    Write-Warning "`t Error while Deleting Printer" #: $_.exception.message"
                }


                #----Deleting Driver----------------------
                try
                {
                    write-host ""
                    write-host -ForegroundColor 8 "`t Deleting Driver: $($PrinterDriverName) "
                    get-wmiobject -class "win32_printerdriver" -namespace "root\CIMV2" | where{$_.name -match $PrinterDriverName} | ForEach-Object { 
                        
                        write-host -ForegroundColor 8 "`t`t Driver Name : $($_.name)"

                        $_.delete()
                        write-host -ForegroundColor 7 "`t`t Driver Deleted: true"
                    }
                }
                catch
                {
                    Write-Warning "`t Error while Deleting Printer Driver" #: $_.exception.message"
                }


                #----Deleting Port----------------------
                try
                {
                    write-host ""
                    write-host -ForegroundColor 8 "`t Deleting Port: $($SelPortName) "
                    $port=Get-WMIObject -Class Win32_tcpipprinterport -filter "name='$($SelPortName)'" -enableall
                    $port.Delete()

                    write-host -ForegroundColor 7 "`t`t Port Deleted : true"
                }
                catch
                {
                    Write-Warning "`t Error while Deleting Printer Port" #: $_.exception.message"
                }

            }
            else
            {
                if($PrinterName -eq $_.name)
                {
                    $PrinterCount = $PrinterCount+1

                    write-host "  ---------------------------------"
                    write-host -ForegroundColor 10 "  Printer Name : $($_.name)" 
                    write-host ""

                    if($_.PrintProcessor -ne 'winprint')
                    {
                        RemoveThirdPartyPrintProcessorAssociation $_.PrintProcessor $_.name
                    }


                    #----Deleting Printer----------------------
                    try
                    {
                        write-host -ForegroundColor 8 "`t Deleting Printer : $($_.name)" 
                        $SelPortName = $_.PortName
                        $_.delete()

                        VerifyPrinterDeleted($_.name)
                        #write-host -ForegroundColor 7 "`t`t Printer Deleted : true"

                    }
                    catch
                    {
                        Write-Warning "`t Error while Deleting Printer" #: $_.exception.message"
                    }


                    #----Deleting Driver----------------------
                    try
                    {
                        write-host ""
                        write-host -ForegroundColor 8 "`t Deleting Driver: $($PrinterDriverName) "
                        get-wmiobject -class "win32_printerdriver" -namespace "root\CIMV2" | where{$_.name -match $PrinterDriverName} | ForEach-Object { 
                            write-host -ForegroundColor 8 "`t`t Driver Name : $($_.name)"

                            $_.delete()
                            write-host -ForegroundColor 7 "`t`t Driver Deleted: true"
                        }
                    }
                    catch
                    {
                        Write-Warning "`t Error while Deleting Printer Driver" #: $_.exception.message"
                    }


                    #----Deleting Port----------------------
                    try
                    {
                        write-host ""
                        write-host -ForegroundColor 8 "`t Deleting Port: $($SelPortName)"
                        $port=Get-WMIObject -Class Win32_tcpipprinterport -filter "name='$($SelPortName)'" -enableall
                        $port.Delete()
                        write-host -ForegroundColor 7 "`t`t Port Deleted : true"
                    }
                    catch
                    {
                        Write-Warning "`t Error while Deleting Printer Port" #: $_.exception.message"
                    }
                }
            }

            write-host ""
        } 

        if($PrinterCount -eq 0)
         {
            $foundPrinter = $flase
            Write-Host ""
            write-host -ForegroundColor 10 "Found no Printer to Delete"
            Write-Host "-------------------------------"
         }
         else
         {
            $foundPrinter = $true
         }

    return $foundPrinter
}

##########################################################
###----Remove third party print processor association-----
##########################################################

Function RemoveThirdPartyPrintProcessorAssociation([string] $processorName , [string] $PrinterName){
    
    write-host -ForegroundColor 8 "`t Remove third party print processor association"

    try
    {
        $PrinterWithOtherProcessor = 0

        get-wmiobject -class "win32_printer" -Computer $computerName | ForEach-Object {

            #write-host -ForegroundColor 10 "Printer Name : $($_.name)" 
   
            if($_.PrintProcessor -eq $processorName)
            {
                $PrinterWithOtherProcessor = $PrinterWithOtherProcessor + 1
            }
         }

        #-----Removing Reference--------------------

        $PrinterRegPath = $RegPrintPath+ '\Printers\' + $PrinterName

        write-host -ForegroundColor 7 "`t`t Printer Reg Path:  $($PrinterRegPath)"

        Set-Itemproperty -path $PrinterRegPath -Name 'Print Processor' -value ''
        Set-Itemproperty -path $PrinterRegPath -Name 'Port' -value ''

        write-host -ForegroundColor 7 "`t`t Reference removed..."
        Write-Host ""

        #write-host -ForegroundColor 7 "`t`t Printer With Other Processor : $($PrinterWithOtherProcessor)"


        if($PrinterWithOtherProcessor -gt 1)
        {
            write-host -ForegroundColor 7 "`t`t Print processor is used by multiple Printer : $($processorName)"
        }
        elseif($PrinterWithOtherProcessor -eq 1)
        {
            $ProcessorRegPath1 = $RegPrintPath + '\Environments\Windows x64\Print Processors\' + $processorName
            if(Test-Path $ProcessorRegPath1)
            {
                write-host -ForegroundColor 7 "`t`t Print processor Name : $($processorName)"
                write-host -ForegroundColor 7 "`t`t Print processor Path : $($ProcessorRegPath1)"

                Remove-Item -Path $ProcessorRegPath1 -Recurse
                write-host -ForegroundColor 7 "`t`t Print processor Removed.."
            }
          

            $ProcessorRegPath2 = $RegPrintPath + '\Environments\Windows NT x86\Print Processors\' + $processorName
            if(Test-Path $ProcessorRegPath2)
            {
                write-host -ForegroundColor 7 "`t`t Print processor Name : $($processorName)"
                write-host -ForegroundColor 7 "`t`t Print processor Path : $($ProcessorRegPath2)"

                Remove-Item -Path $ProcessorRegPath2 -Recurse
                write-host -ForegroundColor 7 "`t`t Print processor Removed.."
            }


            $ProcessorRegPath3 = $RegPrintPath + '\Environments\Windows IA64\Print Processors\' + $processorName
            if(Test-Path $ProcessorRegPath3)
            {
                write-host -ForegroundColor 7 "`t`t Print processor Name : $($processorName)"
                write-host -ForegroundColor 7 "`t`t Print processor Path : $($ProcessorRegPath3)"

                Remove-Item -Path $ProcessorRegPath3 -Recurse
                write-host -ForegroundColor 7 "`t`t Print processor Removed.."
            }

            
            $ProcessorRegPath4 = $RegPrintPath + '\Environments\Windows 4.0\Print Processors\' + $processorName
            if(Test-Path $ProcessorRegPath4)
            {
                write-host -ForegroundColor 7 "`t`t Print processor Name : $($processorName)"
                write-host -ForegroundColor 7 "`t`t Print processor Path : $($ProcessorRegPath4)"

                Remove-Item -Path $ProcessorRegPath4 -Recurse
                write-host -ForegroundColor 7 "`t`t Print processor Removed.."
            }
            
        }
        
        Write-Host ""

    }
    catch
    {
        Write-Warning "Error while Remove third party print processor association: $_.exception.message"
    }
}

##########################################################
###------Delete Printers Que -----------------------------
##########################################################

Function DeletePrintQue(){
    try
    {
        Write-Host "-------------------------------"
        Write-Host "Deleting Print Queue"
        Write-Host ""

        $PrintJobs1 = get-wmiobject -class "Win32_PrintJob" -namespace "root\CIMV2" -computername $computerName 
        if($PrintJobs1 -eq $null)
        {
            Write-Host -ForegroundColor 8 "`t No Print Job to cleare"
        }
        else
        {
            $PrintJobs =Get-WmiObject Win32_Printer | ForEach-Object {$_.CancelAllJobs()}
            Start-Sleep -s $timeoutSeconds

            $PrintJobs2 = get-wmiobject -class "Win32_PrintJob" -namespace "root\CIMV2" -computername $computerName 
            if($PrintJobs2 -eq $null)
            {
                Write-Host -ForegroundColor 8 "`t Print Job cleared"
            }
            else
            {
                ResetSpooler $false
            }
        }
    }
    catch
    {
        Write-Warning "Error while Deleting Print Job: $_.exception.message"
        return $flase
    }
}

##########################################################
###------Reset Spooler -----------------------------------
##########################################################

Function ResetSpooler([bool] $Message){

    try
    {
        if($Message -eq $true)
        {
            Write-Host "-------------------------------"
            Write-Host "Stoping the Spooler Service!"
            Write-Host ""
        }
        

        $j = Start-Job -ScriptBlock {
            Restart-Service -Name Spooler -Force -ErrorAction Stop
        }

        Wait-Job $j -Timeout $timeoutSeconds | out-null

        if ($j.State -eq "Completed"){ 

            if($Message -eq $true)
            {
                Write-Host "Completed.."

                write-host ""
                write-host "Please restart the machine"
            }
             
        }
        elseif ($j.State -eq "Running"){ 
            Write-Host -ForegroundColor 10 "Spooler Service is not responding"
            exit
        }
        
        Remove-Job -force $j 

        if($Message -eq $true)
        {
            Write-Host "-------------------------------"
             break
        }
        
    }
    catch
    {
        Write-Warning "Error while Stoping the Spooler: $_.exception.message"
        return $flase
    }
}

##########################################################
###------Checking Service Status ----------------------------
##########################################################

Function CheckServiceStatus()
{
    Write-Host ""

    if(Get-Service | Where {$_.name -eq 'Spooler'})
    {
        #-----Spooler Check-----------------------
        Get-WmiObject Win32_Service -ComputerName .| Where-Object {$_.Name -eq 'Spooler'  -and  $_.State -eq 'running'}|foreach {
            #write-host "State: " $_.State
            #write-host "Status: " $_.Status
            
            if($_.State -eq "Running")
            {
                if($_.Status -eq "Degraded")
                {
                    Write-Warning "Spooler Service is not responding"
                    exit
                }
            }
        }
        
        #-----Registory Check-----------------------
        Get-Service -Name "Spooler" | Select-Object -Property * |foreach {
            #write-host "Status: " $_.Status
            #write-host "RequiredServices: " $_.ServicesDependedOn

            IF ([string]::IsNullOrEmpty($_.ServicesDependedOn))
            {
                 Write-Warning "Spooler Service is not responding"
                 exit
            } 
        }       
        #-------------------------------------------
    }
    else
    {
        Write-Warning "Spooler Service is not responding"
        exit
    }

}

##########################################################
###------Verify Printer Deleted---------------------------
##########################################################

Function VerifyPrinterDeleted([string] $PrinterName)
{
    $PrinterCount=0;

    Get-WMIObject -Class Win32_Printer -Computer "." | Where-Object { $_.name -eq $PrinterName } | ForEach-Object { 

     $PrinterCount = $PrinterCount+1
    } 

    if($PrinterCount -gt 0)
    {
        write-host -ForegroundColor 10 "`t`t Printer Deleted : False"
    }
    else
    {
         write-host -ForegroundColor 7 "`t`t Printer Deleted : True"
    }
}

##########################################################
###------Execute Functions -------------------------------
##########################################################
cls
if(Check-PreCondition){

    CheckServiceStatus

    if(GetPrinters -PrinterName $PrinterName)
    {
        if(BackupRegistryPath_V8 -Path $RegPrintPath -BackupPath $RegBackupPathName -BackupFileName $RegBackupFileName)
        {
            DeletePrintQue
            if(DeletePrinters)
            {
                #DeletePrintQue
                ResetSpooler $true

            }
        }
    }

}
