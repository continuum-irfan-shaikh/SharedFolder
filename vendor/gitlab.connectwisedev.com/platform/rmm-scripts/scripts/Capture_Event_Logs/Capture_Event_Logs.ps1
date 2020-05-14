####################################################
###-------Capture Event Logs -------------##########
####################################################
 
<#
[string] $LogNames = "Application"
[string] $StartDate= "05/06/2019"
[string] $EndDate="05/24/2019"
[int] $LogCount =10

[string] $EventIds = ""
[string] $EventType = "information"
[string] $EventSource = ""
[bool] $EventDetails = $true
#>



##########################################################
###------Variable Declaration-----------------------------
##########################################################


if(-not $LogCount -gt 0)
{
    $LogCount =1
}

$EventDetails = [System.Convert]::ToBoolean($EventDetails)
$ComputerName = $env:computername


$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Capture Event Logs" 

$SuccessObject = New-Object -TypeName psobject


$global:Code = 0
$global:ErrorMessageArray= @()
$global:SuccessMessageArray= @()


$global:EventLogOutputArray = @()


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
    Write-Host "-------------------------------"
    Write-Host ""

    return $IsContinued
}


##########################################################
###-----Get Event Log by Log Name-------------------
##########################################################

Function GetEventLogByLogName
{

    #Write-Host "-------------------------------"
    #Write-Host "Checking Event Log for - $($LogNames)"
    #Write-Host "" 

    try
    {
        #-----LogNames------------------------
        $LogNames=$LogNames.Replace(" ", "")
        $AllLogNames = $LogNames.Split(",")


        #-----EventIds------------------------
        $EventIdList =$null
        $EventIDsString =""
        if($EventIds.Length -gt 0)
        {
            $EventIDsString = GetEventIdString $EventIds
            $EventIdList = $EventIDsString.split(",")
        }

        #-----Date Range------------------------
        [Nullable[datetime]]$AfterDate=$StartDate
        [Nullable[datetime]]$BeforeDate=$EndDate
        $BeforeDate= $BeforeDate.AddDays(1)


        #-----Event Type List------------
        $EventTypeList =$null
        if($EventType.Length -gt 0)
        {
            $EventType = $EventType.Replace(" ", "")
            $EventTypeList = $EventType.Split(",")
        }

        #-----Event Source List------------
        $EventSourceList= $null
        if($EventSource.Length -gt 0)
        {
            #$EventSource = $EventSource.Replace(" ", "")
            $EventSourceList = $EventSource.Split(",")
        }



        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Event Start Date' -Value $($AfterDate.ToString())
        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Event End Date' -Value $($BeforeDate.ToString())

        $global:SuccessMessageArray += "Event Start Date: $($AfterDate.ToString())"
        $global:SuccessMessageArray += "Event End Date: $($BeforeDate.ToString())"


        if($EventIds.Length -gt 0)
        {
            $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Event Ids' -Value $($EventIds)
            $global:SuccessMessageArray += "Event Ids: $($EventIds)"
        }
        if($EventType.Length -gt 0)
        {
            $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Event Type' -Value $($EventType)
            $global:SuccessMessageArray += "Event Type: $($EventType)"
        }
        if($EventSource.Length -gt 0)
        {
            $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Event Source' -Value $($EventSource)
            $global:SuccessMessageArray += "Event Source: $($EventSource)"
        }

        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'EventDetails' -Value $($EventDetails)
        $global:SuccessMessageArray += "EventDetails: $($EventDetails)"



        #$global:SuccessMessageArray += "TimeGenerated,    EntryType,      EventId,       Source"

        if($EventDetails -eq $true)
        {
            if(($EventIdList.Count -gt 0) -and ($EventTypeList.Count -gt 0) -and ($EventSourceList.Count -gt 0))
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -InstanceId $EventIdList -EntryType $EventTypeList -Source $EventSourceList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop } | where {$_.InstanceId -gt 0} | Select  TimeGenerated, EntryType, EventId, Source, Message
            }
            elseif(($EventIdList.Count -gt 0) -and ($EventTypeList.Count -gt 0))
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -InstanceId $EventIdList -EntryType $EventTypeList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop } | where {$_.InstanceId -gt 0} | Select  TimeGenerated, EntryType, EventId, Source, Message
            }
            elseif(($EventIdList.Count -gt 0)  -and ($EventSourceList.Count -gt 0))
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -InstanceId $EventIdList -Source $EventSourceList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop } | where {$_.InstanceId -gt 0} | Select  TimeGenerated, EntryType, EventId, Source, Message
            }
            elseif(($EventTypeList.Count -gt 0) -and ($EventSourceList.Count -gt 0))
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -EntryType $EventTypeList -Source $EventSourceList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop } | where {$_.InstanceId -gt 0} | Select  TimeGenerated, EntryType, EventId, Source, Message
            }
            elseif($EventIdList.Count -gt 0)
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -InstanceId $EventIdList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop } | where {$_.InstanceId -gt 0} | Select  TimeGenerated, EntryType, EventId, Source, Message
            }
            elseif($EventTypeList.Count -gt 0)
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -EntryType $EventTypeList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop} | Select  TimeGenerated, EntryType, EventId, Source , Message
            }
            elseif($EventSourceList.Count -gt 0)
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -Source $EventSourceList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop} | Select  TimeGenerated, EntryType, EventId, Source, Message
            }
            else
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop} | Select  TimeGenerated, EntryType, EventId, Source, Message
            }

            foreach($log in $EVLog)
            {
                $LogOutput = New-Object -TypeName psobject
                $LogOutput | Add-Member -MemberType NoteProperty -Name 'TimeGenerated' -Value $($log.TimeGenerated.tostring())
                $LogOutput | Add-Member -MemberType NoteProperty -Name 'EntryType' -Value $($log.EntryType.tostring())
                $LogOutput | Add-Member -MemberType NoteProperty -Name 'EventId' -Value $($log.EventId)
                $LogOutput | Add-Member -MemberType NoteProperty -Name 'Source' -Value $($log.Source)
                $LogOutput | Add-Member -MemberType NoteProperty -Name 'Message' -Value $($log.Message)

                $global:EventLogOutputArray += $LogOutput


                #$global:SuccessMessageArray += "TimeGenerated: $($log.TimeGenerated.tostring()),  EntryType: $($log.EntryType),  EventId: $($log.EventId),  Source: $($log.Source),   Message: $($log.Message)"
                $global:SuccessMessageArray += "$($log.TimeGenerated.tostring()),  $($log.EntryType),  $($log.EventId),  $($log.Source),   $($log.Message.tostring())"
            }
        }
        else
        {
            if(($EventIdList.Count -gt 0) -and ($EventTypeList.Count -gt 0) -and ($EventSourceList.Count -gt 0))
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -InstanceId $EventIdList -EntryType $EventTypeList -Source $EventSourceList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop } | where {$_.InstanceId -gt 0} | Select  TimeGenerated, EntryType, EventId, Source
            }
            elseif(($EventIdList.Count -gt 0) -and ($EventTypeList.Count -gt 0))
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -InstanceId $EventIdList -EntryType $EventTypeList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop } | where {$_.InstanceId -gt 0} | Select  TimeGenerated, EntryType, EventId, Source
            }
            elseif(($EventIdList.Count -gt 0)  -and ($EventSourceList.Count -gt 0))
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -InstanceId $EventIdList -Source $EventSourceList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop } | where {$_.InstanceId -gt 0} | Select  TimeGenerated, EntryType, EventId, Source
            }
            elseif(($EventTypeList.Count -gt 0) -and ($EventSourceList.Count -gt 0))
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -EntryType $EventTypeList -Source $EventSourceList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop } | where {$_.InstanceId -gt 0} | Select  TimeGenerated, EntryType, EventId, Source
            }
            elseif($EventIdList.Count -gt 0)
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -InstanceId $EventIdList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop } | where {$_.InstanceId -gt 0} | Select  TimeGenerated, EntryType, EventId, Source
            }
            elseif($EventTypeList.Count -gt 0)
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -EntryType $EventTypeList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop} | Select  TimeGenerated, EntryType, EventId, Source 
            }
            elseif($EventSourceList.Count -gt 0)
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -Source $EventSourceList -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop} | Select  TimeGenerated, EntryType, EventId, Source
            }
            else
            {
                $EVLog = $AllLogNames | foreach  { Get-EventLog -LogName $_  -After $AfterDate -Before $BeforeDate -Newest $LogCount -ErrorAction Stop} | Select  TimeGenerated, EntryType, EventId, Source
            }
  
            foreach($log in $EVLog)
            {
                $LogOutput = New-Object -TypeName psobject
                $LogOutput | Add-Member -MemberType NoteProperty -Name 'TimeGenerated' -Value $($log.TimeGenerated.tostring())
                $LogOutput | Add-Member -MemberType NoteProperty -Name 'EntryType' -Value $($log.EntryType.tostring())
                $LogOutput | Add-Member -MemberType NoteProperty -Name 'EventId' -Value $($log.EventId)
                $LogOutput | Add-Member -MemberType NoteProperty -Name 'Source' -Value $($log.Source)

                $global:EventLogOutputArray += $LogOutput


                #$global:SuccessMessageArray += "TimeGenerated: $($log.TimeGenerated.tostring()),  EntryType: $($log.EntryType),  EventId: $($log.EventId),  Source: $($log.Source)"
                $global:SuccessMessageArray += "$($log.TimeGenerated.tostring()),  $($log.EntryType),  $($log.EventId),  $($log.Source)"
            }
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Error while retriving event log information: $($_.exception.message)"

    }
}


##########################################################
###-----Set EventId String -------------------------------
##########################################################

Function GetEventIdString([string] $EventIdsString)
{
    $Result= ""

    try
    {
        $EventIdsString= $EventIdsString.Replace(" ", "")

        $EventIdArray = $EventIdsString.Split(",")

        foreach($EvntId in $EventIdArray)
        {
            if ($EvntId.Contains("-"))
            {
                $NewRes= GetAllNumberBetween $EvntId

                $Result = $Result + $NewRes 
            }
            else
            {
                $Result = $Result + $EvntId + ","
            }
        }

    }
    catch
    {
        $Result= ""
    }

    return $Result
}


Function GetAllNumberBetween([string] $NumberString)
{
    $Result= ""

    try
    {
        $NumberString= $NumberString.Replace(" ", "")

        [int] $FirstNo, [int] $SecondNo = $NumberString.Split("-")
        
        while($FirstNo -ne $SecondNo+1)
        {
            $Result = $Result + $FirstNo + ","
            $FirstNo = $FirstNo + 1 
        }
    }
    catch
    {
        $Result= ""
    }

    return $Result
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
        $successString = "Success: " + "Event log capture was successful"

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value "success"
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value $successString
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $global:SuccessMessageArray


        $OutputObject = New-Object -TypeName psobject
        $OutputObject | Add-Member -MemberType NoteProperty -Name 'Filter' -Value $($SuccessObject)
        $OutputObject | Add-Member -MemberType NoteProperty -Name 'Event Log' -Value $($global:EventLogOutputArray)

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
        $JSONResult= $ResultObject|ConvertTo-Json -Depth 4
        $JSONResult
    }
    else
    {
        $JSONResult= $ResultObject|ConvertTo-JSONPS2 -maxDepth 4
        $JSONResult
    }
}



##########################################################
###------Execute Functions -------------------------------
##########################################################
cls

if(Check-PreCondition)
{


    if($LogNames.Length -le 0 )
    {
        $global:Code = 1
        $global:ErrorMessageArray += "Log Name is not provided"
    }
    elseif(($StartDate.Length -le 0) -and ($EndDate.Length -le 0))
    {
        $global:Code = 1
        $global:ErrorMessageArray += "Date range is not provided"
    }
    elseif($StartDate.Length -le 0)
    {
        $global:Code = 1
        $global:ErrorMessageArray += "Start Date is not provided"
    }
    elseif($EndDate.Length -le 0)
    {
        $global:Code = 1
        $global:ErrorMessageArray += "End Date is not provided"
    }
    else
    {
        GetEventLogByLogName
    } 


    SetResult
    DisplayResult 
}
else
{
    SetResult
    DisplayResult
}

