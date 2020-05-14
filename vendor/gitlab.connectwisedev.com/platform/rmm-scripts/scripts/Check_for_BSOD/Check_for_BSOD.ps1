####################################################
###------Check for BSOD-------------------##########
####################################################
 
<#
[int] $TimeFrame = 2
#>



##########################################################
###------Variable Declaration-----------------------------
##########################################################

if(-not $TimeFrame -gt 0)
{
    $TimeFrame = 2
}


$ComputerName = $env:computername

$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Check for BSOD" 

$SuccessObject = New-Object -TypeName psobject

$global:Code = 0
$global:ErrorMessageArray= @()
$global:SuccessMessageArray= @()

$global:EventLogOutputArray = @()
$global:MiniDumpFileArray= @()

$CrashBehaviourObject = New-Object -TypeName psobject

$global:CrashBehaviour= ""
$EventCount =100

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
        $global:Code = 2
        $global:ErrorMessageArray += "PowerShell version below 2.0 is not supported"

        $IsContinued = $false
    }

    
    ####################################     
    # Verify opearating system Version
    ####################################

    write-host -ForegroundColor 10 "`t Operating System Version : $([System.Environment]::OSVersion.Version.major)" 
    if(-not([System.Environment]::OSVersion.Version.major -ge 6))
    {
        $global:Code = 2
        $global:ErrorMessageArray += "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"

        $IsContinued = $false
    }
    
    
    Write-Host ""
    Write-Host -ForegroundColor 8 "`t Checking Precondition Completed"
    Write-Host ""
    Write-Host "-------------------------------"
    Write-Host ""

    return $IsContinued
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
###------get Event Detail --------------------------------
##########################################################

function GetEventDetail([string] $MessageLike)
{
    try
    {
        $TimeGap = (Get-Date).AddDays(- $TimeFrame)

        $EVLog = Get-EventLog -LogName application -Newest $EventCount -Source 'Windows Error*' |select EventId, Source, timewritten, message | Where-Object {$_.timewritten -gt $($TimeGap)} | where message -match $MessageLike 

        foreach($log in $EVLog)
        {
            $LogOutput = New-Object -TypeName psobject
            $LogOutput | Add-Member -MemberType NoteProperty -Name 'EventId' -Value $($log.EventId)
            $LogOutput | Add-Member -MemberType NoteProperty -Name 'Source' -Value $($log.Source)
            $LogOutput | Add-Member -MemberType NoteProperty -Name 'Time' -Value $($log.timewritten.tostring())
            $LogOutput | Add-Member -MemberType NoteProperty -Name 'Description' -Value $($log.message)

            $global:EventLogOutputArray += $LogOutput

            $global:SuccessMessageArray += "EventId: $($log.EventId),  Source: $($log.Source),  Time: $($log.timewritten.tostring()),  Description: $($log.Message.tostring())"
        }

        if($global:EventLogOutputArray.Count -le 0)
        {
            $global:Code = 1
            $global:ErrorMessageArray += "Within $($TimeFrame) days BlueScreen events not found"

            $global:SuccessMessageArray += "Within $($TimeFrame) days BlueScreen events not found"
        }
    }
    catch
    {
        $global:Code = 1
        $global:ErrorMessageArray+= "Error while retriving event log detail: $_.exception.message"
    }
}




##########################################################
###------get Event Detail --------------------------------
##########################################################
function GetCrashBehaviour
{
    try
    {
        <#
        $CrashBeh = Get-WmiObject Win32_OSRecoveryConfiguration -EnableAllPrivileges
        $global:CrashBehaviour = $CrashBeh #| Format-List *
        #>
                
        $DebugInfo = (Get-WmiObject -Class Win32_OSRecoveryConfiguration).DebugInfoType
        $path = (Get-WmiObject -Class Win32_OSRecoveryConfiguration).PATH
        $Name = (Get-WmiObject -Class Win32_OSRecoveryConfiguration).Name
        $MemoryDumpLocation = (Get-WmiObject -Class Win32_OSRecoveryConfiguration).ExpandedDebugFilePath
        $MiniDumLocation = (Get-WmiObject -Class Win32_OSRecoveryConfiguration).ExpandedMiniDumpDirectory


        $CBehv = ""
        if($DebugInfo -eq 0)
        {
            $CBehv = "None"

            $global:Code = 1
            $global:ErrorMessageArray+= "Memory dump info could not be obtained because memory dumps are not enabled"
        }
        elseif($DebugInfo -eq 1)
        {
            $CBehv = "Complete memory dump"
        }
        elseif($DebugInfo -eq 2)
        {
            $CBehv = "Kernel memory dump"
        }
        elseif($DebugInfo -eq 3)
        {
            $CBehv = "Small memory dump"
        }


        $CrashBehaviourObject | Add-Member -MemberType NoteProperty -Name 'Debug Info Type' -Value "$DebugInfo ($($CBehv))"
        $CrashBehaviourObject | Add-Member -MemberType NoteProperty -Name 'Memory Dump Location' -Value $MemoryDumpLocation
        $CrashBehaviourObject | Add-Member -MemberType NoteProperty -Name 'Mini Dump Location' -Value $MiniDumLocation
        $CrashBehaviourObject | Add-Member -MemberType NoteProperty -Name 'PATH' -Value $path
        $CrashBehaviourObject | Add-Member -MemberType NoteProperty -Name 'Name' -Value $Name
    }
    catch
    {
        $global:Code = 1
        $global:ErrorMessageArray+= "Error while retriving Crash Behaviour: $_.exception.message"
    }
}


##########################################################
###------Get Mini Dump Files -----------------------------
##########################################################
function GetMiniDumpFiles
{
    try
    {
        $MiniDumLocation1 = (Get-WmiObject -Class Win32_OSRecoveryConfiguration).ExpandedMiniDumpDirectory
        $AllFiles = Get-ChildItem -Path $MiniDumLocation1 -errorAction stop

        foreach($miniDumpFile in $AllFiles)
        {
            $global:MiniDumpFileArray += $miniDumpFile.Name
        }
    }
    catch
    {
        $global:Code = 1
        $global:ErrorMessageArray+= "Error while retriving MiniDumpFiles: $_.exception.message"
    }
}


##########################################################
###------Set Result --------------------------------------
##########################################################

function SetResult 
{
    $ResultObject | Add-Member -MemberType NoteProperty -Name 'Code' -Value $global:Code

    if(($global:Code -eq 0) -or ($global:Code -eq 1))
    {
        $statusMSG = ""
        if($global:Code -eq 0)
        {
            $successString = "Success: " + "The operation completed successfully"
            $statusMSG = "success"
        }
        else
        {
            $successString = "Fail: " + "The operation faild with error"
            $statusMSG = "fail"
        }

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value $statusMSG
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value $successString
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $global:SuccessMessageArray


        $OutputObject = New-Object -TypeName psobject
        $OutputObject | Add-Member -MemberType NoteProperty -Name 'CrashBehaviour' -Value $($CrashBehaviourObject)


        if($global:EventLogOutputArray.Count -gt 0)
        {
            $msg= "Within $($TimeFrame) days $($global:EventLogOutputArray.Count)  BSODs were discovered" 

            $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Message' -Value $msg    #"BlueScreen events found"
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'Result' -Value $($SuccessObject)

            $OutputObject | Add-Member -MemberType NoteProperty -Name 'BSOD' -Value $($global:EventLogOutputArray)
        }
       
        
        if($global:MiniDumpFileArray.Count -gt 0)
        {
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'MiniDumpFiles' -Value $($global:MiniDumpFileArray)
        }

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'dataObject' -Value $OutputObject


        #------------Error with Success--------------------------------------------
        $errorCount= 0
        $ErrorObjectAray= @()

        foreach($ErrorMessage in $global:ErrorMessageArray)
        {
           
            $ErrObject = New-Object -TypeName psobject
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'id' -Value $errorCount
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'title' -Value $ErrorMessage
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'detail' -Value $ErrorMessage

            $ErrorObjectAray += $ErrObject

            $errorCount = $errorCount +1
        }

        if($ErrorObjectAray.Count -gt 0)
        {
            $ResultObject | Add-Member -MemberType NoteProperty -Name 'stderr' -Value $ErrorObjectAray
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
        $JSONResult= $ResultObject|ConvertTo-Json -Depth 6
        $JSONResult
    }
    else
    {
        $JSONResult= $ResultObject|ConvertTo-JSONPS2 -maxDepth 6
        $JSONResult
    }
}



##########################################################
###------Execute Functions -------------------------------
##########################################################
cls

if(Check-PreCondition)
{
    GetCrashBehaviour

    GetEventDetail -MessageLike "blue"
    
    GetMiniDumpFiles
       
    SetResult
    DisplayResult 
}
else
{
    SetResult
    DisplayResult
}
