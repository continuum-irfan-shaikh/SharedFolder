####################################################
###-------Enable LMI ---------------------##########
####################################################
 




##########################################################
###------Variable Declaration-----------------------------
##########################################################

$ComputerName = $env:computername

$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Enable LMI" 

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
###------Check and run Service----------------------------
##########################################################

Function CheckAndStartService([String] $ServiceName , $ServiceDisplayName)
{
    $IsContinued = $false

    Write-Host "-------------------------------"
    Write-Host "Checking $ServiceDisplayName Service"
    Write-Host "    " 

    try
    {
        Get-WmiObject -Class Win32_Service -Filter "Name='$ServiceName'" -ErrorAction Stop| Select Name, Status, StartMode, State | foreach {
            
            $ServiceOutput = New-Object -TypeName psobject
            $ServiceOutput | Add-Member -MemberType NoteProperty -Name 'Service Display Name' -Value $($ServiceDisplayName)
            $ServiceOutput | Add-Member -MemberType NoteProperty -Name 'Service' -Value $($_.Name)
            $ServiceOutput | Add-Member -MemberType NoteProperty -Name 'Status' -Value $($_.Status)
            $ServiceOutput | Add-Member -MemberType NoteProperty -Name 'StartMode' -Value $($_.StartMode)

            $global:ServiceOutputArray  += $ServiceOutput

            $global:SuccessMessageArray += "Service Display Name: $($ServiceDisplayName),  Service: $($_.Name),  Status: $($_.Status),  StartMode: $($_.StartMode)"

            
            if($_.StartMode -eq "Disabled")
            {
                set-service $ServiceName -startuptype automatic
            }
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray+= "Error while Starting $ServiceName Service: $_.exception.message"

        return $IsContinued
    }

    CheckServiceStatus -ServiceName $ServiceName

    try
    {
        if(Get-Service $ServiceName -ErrorAction Stop)
        {
            if( -not((Get-Service | Where {$_.name -eq $ServiceName}).Status -eq "Running"))
            { 
                try
                {
                    Get-Service $ServiceName -ErrorAction Stop| Where {$_.Status -eq 'Stopped'} | Start-Service -ErrorAction Stop 
                    
                    Start-Sleep -s 9

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
                    $global:ErrorMessageArray+= "Error while Starting $ServiceDisplayName Servic: $_.exception.message"
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
        $global:ErrorMessageArray+= "Error while Starting $ServiceDisplayName Service: $_.exception.message"
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
        #-----Service Check-----------------------
        Get-WmiObject Win32_Service -ComputerName .| Where-Object {$_.Name -eq $ServiceName  -and  $_.State -eq 'running'}|foreach {
            #write-host "State: " $_.State
            #write-host "Status: " $_.Status
            #write-host "Startup Type:  $($_.StartMode)"
            
            if($_.State -eq "Running")
            {
                if($_.Status -eq "Degraded")
                {
                    $global:Code = 3
                    $global:ErrorMessageArray+= "$($ServiceName) Service is not responding"
                    exit
                }
            }
        }
        
        #-----Registory Check-----------------------
        Get-Service -Name $ServiceName | Select-Object -Property * |foreach {
            #write-host "Status: " $_.Status
            #write-host "RequiredServices: " $_.ServicesDependedOn

            <#
            IF ([string]::IsNullOrEmpty($_.ServicesDependedOn))
            {
                 Write-Warning "$ServiceName Service is not responding"
                 exit
            } 
            #>
        }       
        #-------------------------------------------
    }
    else
    {
        $global:Code = 3
        $global:ErrorMessageArray+= "$($ServiceName) Service is not responding"
        #exit
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
        $successString = "Success: " + "Services are running"

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value "success"
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value $successString
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $global:SuccessMessageArray


        $OutputObject = New-Object -TypeName psobject
        $OutputObject | Add-Member -MemberType NoteProperty -Name 'Services' -Value $($global:ServiceOutputArray)

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'dataObject' -Value $OutputObject
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
    $ServiceUp =$false


    if(CheckAndStartService -ServiceName 'Winmgmt' -ServiceDisplayName 'Windows Management Instrumentation')
    {
    	if(CheckAndStartService -ServiceName 'LogMeIn' -ServiceDisplayName 'LogMeIn')
    	{	
            if(CheckAndStartService -ServiceName 'LMIGuardianSvc' -ServiceDisplayName 'LMIGuardianSvc')
            {
                if(CheckAndStartService -ServiceName 'LMIMaint' -ServiceDisplayName 'LogMeIn Maintenance Service')
                {
                   $ServiceUp =$true
                }
            }
    	}
    }


    if($ServiceUp -eq $true)
    {
        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Message' -Value "Services are running"
    }
       

    SetResult
    DisplayResult 
}
else
{
    SetResult
    DisplayResult
}
