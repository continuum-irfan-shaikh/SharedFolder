####################################################
###-------Repair LMI ---------------------##########
####################################################
 





##########################################################
###------Variable Declaration-----------------------------
##########################################################

$Key= "DEPLOYID=86t1g7qkbkmgh9v7w8th8r44m8k8lfwqbynigoz5"

$InstallerMSI1 = "C:\Program Files\SAAZOD\LMI\logmein.msi"
$InstallerMSI2 = "C:\Program Files (x86)\SAAZOD\LMI\logmein.msi"

$LogFile1 ="C:\Program Files\SAAZOD\LMI\ApplicationLog\msilog.log" 
$LogFile2 ="C:\Program Files (x86)\SAAZOD\LMI\ApplicationLog\msilog.log"   

$LogFile1 = "c:\Temp\msilog.log" 
$LogFile2 = "c:\Temp\msilog.log"


$ComputerName = $env:computername

$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Repair LMI" 

$SuccessObject = New-Object -TypeName psobject


$global:Code = 0
$global:ErrorMessageArray= @()
$global:SuccessMessageArray= @()


$global:ServiceOutputArray = @()

    
##########################################################
###------Checking Pre Condition---------------------------
##########################################################

Function Check-PreCondition{

    $IsContinued = $true

    Write-Host "-------------------------------"
    Write-Host "Checking Preconditions"
    Write-Host "" 
   

    #####################################
    # Verify PowerShell Version
    #####################################

    write-host -ForegroundColor 10 "`t PowerShell Version : $($PSVersionTable.PSVersion.Major)" 

    if(-NOT($PSVersionTable.PSVersion.Major -ge 2))
    {
        $global:Code = 1
        $global:ErrorMessageArray += "PowerShell version below 2.0 is not supported"

        $IsContinued = $false
    }

    
    ####################################     
    # Verify opearating system Version
    ####################################

    write-host -ForegroundColor 10 "`t Operating System Version : $([System.Environment]::OSVersion.Version.major)" 
    if(-not([System.Environment]::OSVersion.Version.major -ge 6))
    {
        $global:Code = 1
        $global:ErrorMessageArray += "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"

        $IsContinued = $false
    }
    
    
    Write-Host ""
    Write-Host -ForegroundColor 8 "`t Checking Precondition Completed"
    Write-Host ""

    return $IsContinued
}



##########################################################
###------Get installer MSI-----------------------------
##########################################################

Function GetInstallerMSI
{
    $result = "0"

    Write-Host "-------------------------------"
    Write-Host "Checking installer MSI"
    Write-Host "" 
    
    try
    {
        if(Test-Path $InstallerMSI1)
        {
            write-host -ForegroundColor 8 "`t File exist : $($InstallerMSI1)" 
            $result = 1
        }
        else
        {
            if(Test-Path $InstallerMSI2)
            {
                write-host -ForegroundColor 8 "`t File exist : $($InstallerMSI2)" 
                $result = 2
            }
            else
            {
                $global:Code = 1
                $global:ErrorMessageArray+= "No MSI file found"
            }
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray+= "Error while locating MSI file: $_.exception.message"
        return $result
    }

    return $result
}


##########################################################
###------Install/Unstall/Repair LMI-----------------------
##########################################################

Function InstallLMI([string] $InstallerFile, [string] $logFile, [string] $InstallationType)
{
    $Result= -1
    
    Write-Host "-------------------------------"
    if($InstallationType -eq "/x")
    {
        Write-Host "Uninstalling LMI" 
        Write-Host ""
        write-host -ForegroundColor 8 "`t Uninstalling LMI started, Please wait........ " 
    }
    elseif($InstallationType -eq "/i")
    {
        Write-Host "Installing LMI"
        Write-Host "" 
        write-host -ForegroundColor 8 "`t Installing LMI started, Please wait........ " 
    }
    elseif($InstallationType -eq "/fa")
    {
        Write-Host "Repairing LMI"
        Write-Host ""
        write-host -ForegroundColor 8 "`t Repairing LMI started, Please wait........ " 
    }
    Write-Host ""
     
    #Write-host -ForegroundColor 8 "`t Work: $($InstallationType) " 
    #Write-host -ForegroundColor 8 "`t MSI: $($InstallerFile) " 
    #Write-host -ForegroundColor 8 "`t Log: $($logFile) " 

    try
    {
        $MSIArguments = @(
            $InstallationType
            ('"{0}"' -f $InstallerFile)
            "/qn"
            $Key
            "/norestart"
            "/L*v"
            $logFile
        )
        $Result = (Start-Process "msiexec.exe" -ArgumentList $MSIArguments -Wait -Passthru).ExitCode
        
        #write-host -ForegroundColor 8 "`t Unstalling LMI Completed Code:  $($Result)" 
    }
    catch
    {
        $global:Code = 2

        if($InstallationType -eq "/x")
        {
            $global:ErrorMessageArray+= "Error while Uninstalling LMI: $_.exception.message"
        }
        elseif($InstallationType -eq "/i")
        {
           $global:ErrorMessageArray+=  "Error while Installing LMI: $_.exception.message"
        }
        elseif($InstallationType -eq "/fa")
        {
            $global:ErrorMessageArray+= "Error while Repairing LMI: $_.exception.message"
        }
        return $result
    }
    
    return $result
}



##########################################################
###------GetMessage for ExitCode--------------------------
##########################################################

Function GetMessageforExitCode([int] $ExtCode)
{
    $result= ""
    
    try
    {
        if ( $ExtCode -eq 0 ) { $result = 'The action completed successfully.'}
        elseif ( $ExtCode -eq 13 ) { $result = 'The data is invalid.'}
        elseif ( $ExtCode -eq 87 ) { $result = 'One of the parameters was invalid.'}
        elseif ( $ExtCode -eq 120 ) { $result = 'This value is returned when a custom action attempts to call a function that cannot be called from custom actions. The function returns the value ERROR_CALL_NOT_IMPLEMENTED. Available beginning with Windows Installer version 3.0.'}
        elseif ( $ExtCode -eq 1259 ) { $result = 'If Windows Installer determines a product may be incompatible with the current operating system, it displays a dialog box informing the user and asking whether to try to install anyway. This error code is returned if the user chooses not to try the installation.'}
        elseif ( $ExtCode -eq 1601 ) { $result = 'The Windows Installer service could not be accessed. Contact your support personnel to verify that the Windows Installer service is properly registered.'}
        elseif ( $ExtCode -eq 1602 ) { $result = 'The user cancels installation.'}
        elseif ( $ExtCode -eq 1603 ) { $result = 'A fatal error occurred during installation.'}
        elseif ( $ExtCode -eq 1604 ) { $result = 'Installation suspended, incomplete.'}
        elseif ( $ExtCode -eq 1605 ) { $result = 'This action is only valid for products that are currently installed.'}
        elseif ( $ExtCode -eq 1606 ) { $result = 'The feature identifier is not registered.'}
        elseif ( $ExtCode -eq 1607 ) { $result = 'The component identifier is not registered.'}
        elseif ( $ExtCode -eq 1608 ) { $result = 'This is an unknown property.'}
        elseif ( $ExtCode -eq 1609 ) { $result = 'The handle is in an invalid state.'}
        elseif ( $ExtCode -eq 1610 ) { $result = 'The configuration data for this product is corrupt. Contact your support personnel.'}
        elseif ( $ExtCode -eq 1611 ) { $result = 'The component qualifier not present.'}
        elseif ( $ExtCode -eq 1612 ) { $result = 'The installation source for this product is not available. Verify that the source exists and that you can access it.'}
        elseif ( $ExtCode -eq 1613 ) { $result = 'This installation package cannot be installed by the Windows Installer service. You must install a Windows service pack that contains a newer version of the Windows Installer service.'}
        elseif ( $ExtCode -eq 1614 ) { $result = 'The product is uninstalled.'}
        elseif ( $ExtCode -eq 1615 ) { $result = 'The SQL query syntax is invalid or unsupported.'}
        elseif ( $ExtCode -eq 1616 ) { $result = 'The record field does not exist.'}
        elseif ( $ExtCode -eq 1618 ) { $result = 'Another installation is already in progress. Complete that installation before proceeding with this install.For information about the mutex, see _MSIExecute Mutex.'}
        elseif ( $ExtCode -eq 1619 ) { $result = 'This installation package could not be opened. Verify that the package exists and is accessible, or contact the application vendor to verify that this is a valid Windows Installer package.'}
        elseif ( $ExtCode -eq 1620 ) { $result = 'This installation package could not be opened. Contact the application vendor to verify that this is a valid Windows Installer package.'}
        elseif ( $ExtCode -eq 1621 ) { $result = 'There was an error starting the Windows Installer service user interface. Contact your support personnel.'}
        elseif ( $ExtCode -eq 1622 ) { $result = 'There was an error opening installation log file. Verify that the specified log file location exists and is writable.'}
        elseif ( $ExtCode -eq 1623 ) { $result = 'This language of this installation package is not supported by your system.'}
        elseif ( $ExtCode -eq 1624 ) { $result = 'There was an error applying transforms. Verify that the specified transform paths are valid.'}
        elseif ( $ExtCode -eq 1625 ) { $result = 'This installation is forbidden by system policy. Contact your system administrator.'}
        elseif ( $ExtCode -eq 1626 ) { $result = 'The function could not be executed.'}
        elseif ( $ExtCode -eq 1627 ) { $result = 'The function failed during execution.'}
        elseif ( $ExtCode -eq 1628 ) { $result = 'An invalid or unknown table was specified.'}
        elseif ( $ExtCode -eq 1629 ) { $result = 'The data supplied is the wrong type.'}
        elseif ( $ExtCode -eq 1630 ) { $result = 'Data of this type is not supported.'}
        elseif ( $ExtCode -eq 1631 ) { $result = 'The Windows Installer service failed to start. Contact your support personnel.'}
        elseif ( $ExtCode -eq 1632 ) { $result = 'The Temp folder is either full or inaccessible. Verify that the Temp folder exists and that you can write to it.'}
        elseif ( $ExtCode -eq 1633 ) { $result = 'This installation package is not supported on this platform. Contact your application vendor.'}
        elseif ( $ExtCode -eq 1634 ) { $result = 'Component is not used on this machine.'}
        elseif ( $ExtCode -eq 1635 ) { $result = 'This patch package could not be opened. Verify that the patch package exists and is accessible, or contact the application vendor to verify that this is a valid Windows Installer patch package.'}
        elseif ( $ExtCode -eq 1636 ) { $result = 'This patch package could not be opened. Contact the application vendor to verify that this is a valid Windows Installer patch package.'}
        elseif ( $ExtCode -eq 1637 ) { $result = 'This patch package cannot be processed by the Windows Installer service. You must install a Windows service pack that contains a newer version of the Windows Installer service.'}
        elseif ( $ExtCode -eq 1638 ) { $result = 'Another version of this product is already installed. Installation of this version cannot continue. To configure or remove the existing version of this product, use Add/Remove Programs in Control Panel.'}
        elseif ( $ExtCode -eq 1639 ) { $result = 'Invalid command line argument. Consult the Windows Installer SDK for detailed command-line help.'}
        elseif ( $ExtCode -eq 1640 ) { $result = 'The current user is not permitted to perform installations from a client session of a server running the Terminal Server role service.'}
        elseif ( $ExtCode -eq 1641 ) { $result = 'The installer has initiated a restart. This message is indicative of a success.'}
        elseif ( $ExtCode -eq 1642 ) { $result = 'The installer cannot install the upgrade patch because the program being upgraded may be missing or the upgrade patch updates a different version of the program. Verify that the program to be upgraded exists on your computer and that you have the correct upgrade patch.'}
        elseif ( $ExtCode -eq 1643 ) { $result = 'The patch package is not permitted by system policy.'}
        elseif ( $ExtCode -eq 1644 ) { $result = 'One or more customizations are not permitted by system policy.'}
        elseif ( $ExtCode -eq 1645 ) { $result = 'Windows Installer does not permit installation from a Remote Desktop Connection.'}
        elseif ( $ExtCode -eq 1646 ) { $result = 'The patch package is not a removable patch package. Available beginning with Windows Installer version 3.0.'}
        elseif ( $ExtCode -eq 1647 ) { $result = 'The patch is not applied to this product. Available beginning with Windows Installer version 3.0.'}
        elseif ( $ExtCode -eq 1648 ) { $result = 'No valid sequence could be found for the set of patches. Available beginning with Windows Installer version 3.0.'}
        elseif ( $ExtCode -eq 1649 ) { $result = 'Patch removal was disallowed by policy. Available beginning with Windows Installer version 3.0.'}
        elseif ( $ExtCode -eq 1650 ) { $result = 'The XML patch data is invalid. Available beginning with Windows Installer version 3.0.'}
        elseif ( $ExtCode -eq 1651 ) { $result = 'Administrative user failed to apply patch for a per-user managed or a per-machine application that is in advertise state. Available beginning with Windows Installer version 3.0.'}
        elseif ( $ExtCode -eq 1652 ) { $result = 'Windows Installer is not accessible when the computer is in Safe Mode. Exit Safe Mode and try again or try using System Restore to return your computer to a previous state. Available beginning with Windows Installer version 4.0.'}
        elseif ( $ExtCode -eq 1653 ) { $result = 'Could not perform a multiple-package transaction because rollback has been disabled. Multiple-Package Installations cannot run if rollback is disabled. Available beginning with Windows Installer version 4.5.'}
        elseif ( $ExtCode -eq 1654 ) { $result = 'The app that you are trying to run is not supported on this version of Windows. A Windows Installer package, patch, or transform that has not been signed by Microsoft cannot be installed on an ARM computer.'}
        elseif ( $ExtCode -eq 3010 ) { $result = 'A restart is required to complete the install. This message is indicative of a success. This does not include installs where the ForceReboot action is run.'}
        elseif ( $ExtCode -eq -1 ) { $result = 'Unknown Error.'}
        else { $result = 'No message found for Exit Code: $($ExtCode)'}
    }
    catch
    {
        $result ="No message found for Exit Code: $($ExtCode)"
    }
    
    
    return $result
}




##########################################################
###------JSON for PS2 ------------------------------------
##########################################################

function Escape-JSONString($str)
{
	if ($str -eq $null) 
    {
        return ""
    }

	$str = $str.ToString().Replace('"','\"').Replace('\','\\').Replace("`n",'\n').Replace("`r",'\r').Replace("`t",'\t')

	return $str;
}


function ConvertTo-JSONPS2($maxDepth = 4,$forceArray = $false) 
{
	begin {
		$data = @()
	}
	process{
		$data += $_
	}
	
	end{
	
		if ($data.length -eq 1 -and $forceArray -eq $false) {
			$value = $data[0]
		} else {	
			$value = $data
		}

		if ($value -eq $null) {
			return "null"
		}

		

		$dataType = $value.GetType().Name
		
		switch -regex ($dataType) {
	            'String'  {
					return  "`"{0}`"" -f (Escape-JSONString $value )
				}

	            '(System\.)?DateTime'  {return  "`"{0:yyyy-MM-dd}T{0:HH:mm:ss}`"" -f $value}

	            'Int32|Double' {return  "$value"}

				'Boolean' {return  "$value".ToLower()}

	            '(System\.)?Object\[\]' { # array
					
					if ($maxDepth -le 0){return "`"$value`""}
					
					$jsonResult = ''
					foreach($elem in $value){
						#if ($elem -eq $null) {continue}
						if ($jsonResult.Length -gt 0) {$jsonResult +=', '}				
						$jsonResult += ($elem | ConvertTo-JSONPS2 -maxDepth ($maxDepth -1))
					}
					return "[" + $jsonResult + "]"
	            }

				'(System\.)?Hashtable' { # hashtable
					$jsonResult = ''
					foreach($key in $value.Keys){
						if ($jsonResult.Length -gt 0) {$jsonResult +=', '}
						$jsonResult += 
@"
	"{0}": {1}
"@ -f $key , ($value[$key] | ConvertTo-JSONPS2 -maxDepth ($maxDepth -1) )
					}
					return "{" + $jsonResult + "}"
				}

	            default { #object
					if ($maxDepth -le 0){return  "`"{0}`"" -f (Escape-JSONString $value)}
					
					return "{" +
						(($value | Get-Member -MemberType *property | % { 
@"
	"{0}": {1}
"@ -f $_.Name , ($value.($_.Name) | ConvertTo-JSONPS2 -maxDepth ($maxDepth -1) )			
					
					}) -join ', ') + "}"
	    		}
		}
	}
}




##########################################################
###------Set Result --------------------------------------
##########################################################

function SetResult 
{
    $ResultObject | Add-Member -MemberType NoteProperty -Name 'Code' -Value $global:Code

    if($global:Code -eq 0)
    {
        $successString = "Success: " + "LMI repair successful"

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value "success"
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value $successString
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $global:SuccessMessageArray

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'dataObject' -Value $SuccessObject
    }
    else
    {
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value "fail"

        $ErrorObjectAray= @()
        $errorCount= 0
        $errorString = "Error: "

        foreach($ErrorMessage in $global:ErrorMessageArray)
        {
            $errorString += $ErrorMessage + ", "
            
            $ErrObject = New-Object -TypeName psobject
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'id' -Value $errorCount
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'title' -Value $ErrorMessage
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'detail' -Value $ErrorMessage

            $ErrorObjectAray += $ErrObject

            $errorCount = $errorCount +1
        }
        
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value $errorString
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $global:ErrorMessageArray
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stderr' -Value $ErrorObjectAray
    }
}


##########################################################
###------Display Result ----------------------------------
##########################################################

function DisplayResult 
{

    if($PSVersionTable.PSVersion.Major -gt 2)
    {
        $JSONResult= $ResultObject|ConvertTo-Json
        $JSONResult
    }
    else
    {
        $JSONResult= $ResultObject|ConvertTo-JSONPS2
        $JSONResult
    }
}




##########################################################
###------Execute Functions -------------------------------
##########################################################
cls

if(Check-PreCondition)
{
    $NextStep ="R"

    $Installer = GetInstallerMSI
    
    $InstallerFile= ""
    $LogFile= ""
    
    if($Installer -eq 1)
    {
        $InstallerFile= $InstallerMSI1
        $LogFile= $LogFile1
    }
    elseif($Installer -eq 2)
    {
        $InstallerFile= $InstallerMSI2
        $LogFile= $LogFile2
    }
    
    
    if($InstallerFile.length -gt 0)
    {
        #-----Repair----------------
        if($NextStep -eq "R")
        {
            $Res1 = InstallLMI $InstallerFile $LogFile "/fa"
        
            write-host -ForegroundColor 8 "`t Result Code:  $($Res1)" 
            $Message1=  GetMessageforExitCode $Res1

            $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Repair Status' -Value $($Message1)
            $global:SuccessMessageArray += "Repair Status: $($Message1)"

            if($Res1 -eq 0)
            {
                $NextStep =""
            }
            elseif($Res1 -eq 1605)
            {
                $NextStep ="I"
            }
            else
            {
                $NextStep ="U"
                
                write-host ""
                write-host -ForegroundColor 10 "`t Repair was not sucessful, trying to Uninstall and reinstall LMI"
                write-host ""
                $global:SuccessMessageArray += "Repair was not sucessful, trying to Uninstall and reinstall LMI"
            }
        }
        
        #-----Uninstall----------------
        if($NextStep -eq "U")
        {
            $Res2 = InstallLMI $InstallerFile $LogFile "/x"
            
            write-host -ForegroundColor 8 "`t Result Code:  $($Res2)" 
            $Message2=  GetMessageforExitCode $Res2
           
            $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Uninstall Status' -Value $($Message2)
            $global:SuccessMessageArray += "Uninstall Status: $($Message2)"

                
            if($Res2 -eq 0)
            {
                $NextStep ="I"
                Start-Sleep -s 10
            }
            elseif($Res2 -eq 1605)
            {
                $NextStep ="I"
            }
            else
            {
                $NextStep =""
                
                $global:Code = 1
                $global:ErrorMessageArray+= "LMI Could not be Uninstalled"
                $global:SuccessMessageArray += "LMI Could not be Uninstalled"
            }
        }
        
        #-----Install----------------
        if($NextStep -eq "I")
        {
            
            $res3= InstallLMI $InstallerFile $LogFile "/i"
            
            write-host -ForegroundColor 8 "`t Result Code:  $($res3)" 
            $Message3 =  GetMessageforExitCode $res3

            $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Install Status' -Value $($Message2)
            $global:SuccessMessageArray += "Install Status: $($Message2)"
            
            if($Res3 -eq 0)
            {
                $NextStep =""
            }
            else
            {
                $global:Code = 1
                $global:ErrorMessageArray+= "LMI Could not be Installed"
            }
        }
    }



    SetResult
    DisplayResult 

}
else
{
    SetResult
    DisplayResult
}
