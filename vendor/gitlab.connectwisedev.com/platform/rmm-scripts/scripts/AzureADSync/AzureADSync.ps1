####################################################
###-------AzureAD Sync --------------------#########
####################################################

#----Sync Type can be : None, DeltaSync, FullSync

<#

[string]$username = "Administrator"
[string]$Pass = "Abcd@1234"
[string]$SyncType = "DeltaSync"

#>


##########################################################
###------Variable Declaration-----------------------------
##########################################################

$timeoutSeconds = 30

$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "AzureAD Sync" 

$SuccessObject = New-Object -TypeName psobject

$CompnyInformationObject = New-Object -TypeName psobject

$global:Code = 0
$global:ErrorMessageArray= @()
$global:SuccessMessageArray= @()

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
###------Verify-NewAzure-Service -------------------------
##########################################################

Function Verify-Azure-Service
{
    $IsContinued = $false

    if(CheckAndStartService -ServiceName 'ADSync' -ServiceDisplayName 'Microsoft Azure AD Sync')
    {
        $IsContinued = $true
    }

    return $IsContinued
}


##########################################################
###------Check and run Service----------------------------
##########################################################

Function CheckAndStartService([String] $ServiceName , $ServiceDisplayName)
{
    CheckServiceStatus -ServiceName $ServiceName

    $IsContinued = $false

    Write-Host "-------------------------------"
    Write-Host "Checking $ServiceDisplayName Service"
    Write-Host "    " 

    try
    {
        if(Get-Service $ServiceName -ErrorAction Stop)
        {
            if( -not((Get-Service | Where {$_.name -eq $ServiceName}).Status -eq "Running"))
            { 
                try
                {
                    Get-Service $ServiceName -ErrorAction Stop| Where {$_.Status -eq 'Stopped'} | Start-Service -ErrorAction Stop 
                    
                    Start-Sleep -s 7

                    if((Get-Service -ErrorAction Stop| Where {$_.name -eq $ServiceName}).Status -eq "Running")
                    {
                        Write-Host -ForegroundColor 8 "`t $ServiceDisplayName Service started"
                        $IsContinued = $true
                    }
                    else
                    {
                        $global:Code = 1
                        $global:ErrorMessageArray+= "$($ServiceDisplayName) Service could not be started"
                    }
                }
                catch
                {
                    $global:Code = 2
                    $global:ErrorMessageArray+= "Error while Starting $($ServiceDisplayName) Servic: $_.exception.message"
                }
            }
            else
            {
                Write-Host -ForegroundColor 8 "`t $ServiceDisplayName Service is running"
                $IsContinued = $true
            } 
        }
        else
        {
            $global:Code = 1
            $global:ErrorMessageArray+= "$($ServiceDisplayName) Service Not Found"
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray+= "Error while Starting $($ServiceDisplayName) Service: $_.exception.message"
    }

    return $IsContinued
}


##########################################################
###------Verify Azure Module-----------------------------
##########################################################

Function VerifyModule([String] $ModuleName)
{
    $IsContinued = $false

    Write-Host "-------------------------------"
    Write-Host "Verifying Module : $ModuleName"
    Write-Host "    " 

    if (Get-Module -ErrorAction Stop | Where-Object {$_.Name -eq $ModuleName}) 
    {
        write-host -ForegroundColor 8 "`t $ModuleName exist" 
        $IsContinued = $true
    } 
    else 
    {
        # If module is not imported, but available on disk then import
        if (Get-Module -ListAvailable -ErrorAction Stop| Where-Object {$_.Name -eq $ModuleName}) 
        {
			try
            {
				Import-Module $ModuleName -ErrorAction Stop #-Verbose

				write-host -ForegroundColor 8 "`t $ModuleName exist" 
				$IsContinued = $true
			}
			catch
            {
                $global:Code = 2
                $global:ErrorMessageArray+= "Error while importing Module $($ModuleName): $_.exception.message"
			}
        } 
    }

    return $IsContinued
}



##########################################################
###------Check If any Sync is running --------------------
##########################################################

Function check-ADSyncConnectorRunStatus
{
    $IsContinued = $false

    Write-Host "-------------------------------"
    Write-Host "Checking for Active Session"
    Write-Host "    " 

    if(Get-ADSyncConnectorRunStatus)
    {
        $global:Code = 1
        $global:ErrorMessageArray+= "Sync engine is Busy"
    }
    else
    {
        write-host -ForegroundColor 8 "`t Sync engine is free"
        $IsContinued= $true
    }

    return $IsContinued
}


##########################################################
###------Login to Azure Account --------------------------
##########################################################

Function Login-AzureAccount 
{
    $IsContinued = $false

    Write-Host "-------------------------------"
    Write-Host "Checking Azure Account"
    Write-Host "    " 

    $password = $Pass | ConvertTo-SecureString -asPlainText -Force
    $UserCredential = New-Object System.Management.Automation.PSCredential($username,$password)

    try
    {
        if(Connect-AzureAD -Credential $UserCredential -ErrorAction Stop)
        {
             write-host -ForegroundColor 8 "Azure Login Sucessfull"
             $IsContinued= $true
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray+= "Error while Azure Login: $_.exception.message"
    }

    return $IsContinued
}



##########################################################
###------Get Company Information --------------------------
##########################################################

Function Get-CompanyInformation
{
    Write-Host "-------------------------------"
    Write-Host "Checking Company Information"
    Write-Host "    " 

    $password = $Pass | ConvertTo-SecureString -asPlainText -Force
    $UserCredential = New-Object System.Management.Automation.PSCredential($username,$password)

    try
    {
        Connect-MsolService -Credential $UserCredential

        Get-MsolCompanyInformation | Select-Object DisplayName, DirSyncServiceAccount, LastPasswordSyncTime, LastDirSyncTime |ForEach-Object {
            
            #write-host -ForegroundColor 8 "Name : $($_.DisplayName)" 
            #write-host -ForegroundColor 8 "Dir SyncService Account : $($_.DirSyncServiceAccount)"
            #write-host -ForegroundColor 8 "Last Password sync : $($_.LastPasswordSyncTime)" 
            #write-host -ForegroundColor 8 "Last dir sync : $($_.LastDirSyncTime)"
            
            $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'Name' -Value $($_.DisplayName)
            $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'Dir SyncService Account' -Value $($_.DirSyncServiceAccount) 
            $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'Last Password sync' -Value $($_.LastPasswordSyncTime) 
            $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'Last dir sync' -Value $($_.LastDirSyncTime) 
            
            $global:SuccessMessageArray += "Name: $($_.DisplayName)"
            $global:SuccessMessageArray += "Dir SyncService Account: $($_.DirSyncServiceAccount)" 
            $global:SuccessMessageArray += "Last Password sync: $($_.LastPasswordSyncTime)" 
            $global:SuccessMessageArray += "Last dir sync: $($_.LastDirSyncTime)" 
        }

        #write-host -ForegroundColor 10 "Next sync"

        $schdl = Get-ADSyncScheduler
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'AllowedSyncCycleInterval' -Value $($schdl.AllowedSyncCycleInterval)
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'CurrentlyEffectiveSyncCycleInterval' -Value $($schdl.CurrentlyEffectiveSyncCycleInterval)
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'CustomizedSyncCycleInterval' -Value $($schdl.CustomizedSyncCycleInterval)
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'NextSyncCyclePolicyType' -Value $($schdl.NextSyncCyclePolicyType)
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'NextSyncCycleStartTimeInUTC' -Value $($schdl.NextSyncCycleStartTimeInUTC)
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'PurgeRunHistoryInterval' -Value $($schdl.PurgeRunHistoryInterval)
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'SyncCycleEnabled' -Value $($schdl.SyncCycleEnabled)
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'MaintenanceEnabled' -Value $($schdl.MaintenanceEnabled)
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'StagingModeEnabled' -Value $($schdl.StagingModeEnabled)
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'SchedulerSuspended' -Value $($schdl.SchedulerSuspended)
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'SyncCycleInProgress' -Value $($schdl.SyncCycleInProgress)

        $global:SuccessMessageArray += "AllowedSyncCycleInterval: $($schdl.AllowedSyncCycleInterval)"
        $global:SuccessMessageArray += "CurrentlyEffectiveSyncCycleInterval: $($schdl.CurrentlyEffectiveSyncCycleInterval)"
        $global:SuccessMessageArray += "CustomizedSyncCycleInterval: $($schdl.CustomizedSyncCycleInterval)"
        $global:SuccessMessageArray += "NextSyncCyclePolicyType: $($schdl.NextSyncCyclePolicyType)"
        $global:SuccessMessageArray += "NextSyncCycleStartTimeInUTC: $($schdl.NextSyncCycleStartTimeInUTC)"
        $global:SuccessMessageArray += "PurgeRunHistoryInterval: $($schdl.PurgeRunHistoryInterval)"
        $global:SuccessMessageArray += "SyncCycleEnabled: $($schdl.SyncCycleEnabled)"
        $global:SuccessMessageArray += "MaintenanceEnabled: $($schdl.MaintenanceEnabled)"
        $global:SuccessMessageArray += "StagingModeEnabled: $($schdl.StagingModeEnabled)"
        $global:SuccessMessageArray += "SchedulerSuspended: $($schdl.SchedulerSuspended)"
        $global:SuccessMessageArray += "SyncCycleInProgress: $($schdl.SyncCycleInProgress)"


        #Start-Sleep -s 10

        $OuInfo = Get-MSOlUser -ALL | Select-Object UserPrincipalName, LastDirSyncTime
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'UserPrincipalName' -Value $($OuInfo.UserPrincipalName)
        $CompnyInformationObject | Add-Member -MemberType NoteProperty -Name 'LastDirSyncTime' -Value $($OuInfo.LastDirSyncTime)

        $global:SuccessMessageArray += "UserPrincipalName: $($OuInfo.UserPrincipalName)"
        $global:SuccessMessageArray += "LastDirSyncTime: $($OuInfo.LastDirSyncTime)"
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray+= "Error while getting Company Information: $_.exception.message"
    }

}


##########################################################
###------Sync AD -----------------------------------------
##########################################################

Function Sync-AD 
{
    $IsContinued = $false

    Write-Host "-------------------------------"
    Write-Host "Syncing Active Directory.. ($SyncType)"
    Write-Host "    " 
    
    try
    {
        if($SyncType -eq "DeltaSync")
        {
            Start-ADSyncSyncCycle -PolicyType delta
            Write-Host -ForegroundColor 8 "`t Sync strted Sucessfully"
            $IsContinued = $true
        }
        elseif($SyncType -eq "FullSync")
        {
            Start-ADSyncSyncCycle -PolicyType initial
            Write-Host -ForegroundColor 8 "`t Sync strted Sucessfully"
            $IsContinued = $true
        }
        else
        {
            $global:Code = 1
            $global:ErrorMessageArray+= "Sync type Not Found it should be ('FullSync' or 'DeltaSync')"

            #Write-Host -ForegroundColor 10 "`t Sync type Not Found it should be ('FullSync' or 'DeltaSync')"
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray+= "Error while Syncing Active Directory: $_.exception.message"

        #Write-Warning "Error while Azure Login"#: $_.exception.message"
    }

    return $IsContinued
}


##########################################################
###------Check Service Status ----------------------------
##########################################################

Function CheckServiceStatus([String] $ServiceName)
{
    Write-Host ""

    if(Get-Service | Where {$_.name -eq $ServiceName})
    {
        #-----Spooler Check-----------------------
        Get-WmiObject Win32_Service -ComputerName .| Where-Object {$_.Name -eq $ServiceName  -and  $_.State -eq 'running'}|foreach {
            #write-host "State: " $_.State
            #write-host "Status: " $_.Status
            
            if($_.State -eq "Running")
            {
                if($_.Status -eq "Degraded")
                {
                    $global:Code = 3
                    $global:ErrorMessageArray+= "$($ServiceName) Service is not responding"

                    #Write-Warning "$ServiceName Service is not responding"
                    exit
                }
            }
        }
        
        #-----Registory Check-----------------------
        Get-Service -Name $ServiceName | Select-Object -Property * |foreach {
            #write-host "Status: " $_.Status
            #write-host "RequiredServices: " $_.ServicesDependedOn

            IF ([string]::IsNullOrEmpty($_.ServicesDependedOn))
            {
                $global:Code = 3
                $global:ErrorMessageArray+= "$($ServiceName) Service is not responding"

                #Write-Warning "$ServiceName Service is not responding"
                exit
            } 
        }       
        #-------------------------------------------
    }
    else
    {
        $global:Code = 3
        $global:ErrorMessageArray+= "$($ServiceName) Service is not responding"

        #Write-Warning "$ServiceName Service is not responding"
        exit
    }

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
        if($SyncType -ne "None")
        {
            $successString = "Success: " + "Sync completed Sucessfully"
        }
        else
        {
            $successString = "Success: " + "completed Sucessfully"
        }


        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value "success"
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value $successString
        

        if($SyncType -eq "None")
        {
            $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $global:SuccessMessageArray

            $OutputObject = New-Object -TypeName psobject
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'CompnyInformation' -Value $($CompnyInformationObject)
            $ResultObject | Add-Member -MemberType NoteProperty -Name 'dataObject' -Value $OutputObject
        }
        else
        {
            $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $SuccessObject

            $OutputObject = New-Object -TypeName psobject
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'Services' -Value $($SuccessObject)
            $ResultObject | Add-Member -MemberType NoteProperty -Name 'dataObject' -Value $OutputObject
        }

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
    if($SyncType.Length -le 0)
    {
        $global:Code = 1
        $global:ErrorMessageArray+= "Sync Type is not provided"
    }
    if($Pass.Length -le 0)
    {
        $global:Code = 1
        $global:ErrorMessageArray+= "Password is not provided"
    }
    if($username.Length -le 0)
    {
        $global:Code = 1
        $global:ErrorMessageArray+= "UserName is not provided"
    }
   

    if($global:Code -eq 0)
    {

        $Proceed = $false

        if(Verify-Azure-Service)
        {
            $Proceed = $true
        }

        if($Proceed -eq $true)
        {
            $Proceed = $false

            if(VerifyModule -ModuleName 'AzureAD')
            {
                if(VerifyModule -ModuleName 'MSOnline')
                {
                    $Proceed = $true
                }
            }
        }

        if($Proceed -eq $true)
        {
            if(check-ADSyncConnectorRunStatus)
            {
                if(Login-AzureAccount)
                {
                    write-host -ForegroundColor 8 "`t Azure Login Sucessfull"

                    if($SyncType -eq "None")
                    {
                        Get-CompanyInformation
                    }
                    else
                    {
                        Sync-AD 

                        $SyncDone= $true
                        while($SyncDone)
                        {
                            Write-Host "`t Syncing, Please Wait.."
                            Start-Sleep -s 20
                            $SyncDone= Get-ADSyncConnectorRunStatus
                        }

                        Start-Sleep -s 20
                        
                        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Result' -Value "Sync completed Sucessfully"
                    }
                }
                else
                {
                    $global:Code = 1
                    $global:ErrorMessageArray+= "Azure Login Fail"
                }
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