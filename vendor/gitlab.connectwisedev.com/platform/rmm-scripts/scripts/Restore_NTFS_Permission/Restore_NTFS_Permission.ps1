####################################################
###-------Restore NTFS Permission---------##########
####################################################
 
 <#
    ----------------------------------
    Parameter Example
    -------------------------
    $TargetFolderPath = "D:\" If while taking the backup it was "D:\Test2" it should be one folder before
    $BackupFilePath = "C:\Temp\NTFS_Permission_BKP.txt"

    This will restore NTFS permission to  'Test2' and all Sub folder.
    -----------------------------------
#>


<#
[string] $TargetFolderPath =  "D:\"
[string] $BackupFilePath =  "C:\Temp\NTFS_Permission_BKP.txt"
#>

##########################################################
###------Variable Declaration-----------------------------
##########################################################

$ComputerName = $env:computername

$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Restore NTFS Permission" 

$SuccessObject = New-Object -TypeName psobject

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
    Write-Host "-------------------------------"
    Write-Host ""

    return $IsContinued
}



##########################################################
###------Backup NTFS Permission----------------------------
##########################################################

Function RestoreNTFSPermission
{
    $result = $false

    if(-not (Test-Path $TargetFolderPath))
    {
        $global:Code = 1
        $global:ErrorMessageArray += "Target Folder was not found at $($TargetFolderPath)" 
        return $result
    }


    if(-not(Test-Path $BackupFilePath))
    {
        $global:Code = 1
        $global:ErrorMessageArray += "Backup file not found at $($BackupFilePath)" 
        return $result 
    }


    $ErrorActionPreference = "stop"
    try
    {
	    icacls $TargetFolderPath /restore $BackupFilePath 

        $global:SuccessMessageArray += "TargetFolderPath: $($TargetFolderPath)"
        $global:SuccessMessageArray += "BackupFilePath: $($BackupFilePath)"
        $global:SuccessMessageArray += "Restoring NTFS Permission successful."

        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'TargetFolderPath' -Value $($TargetFolderPath)
        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'BackupFilePath' -Value $($BackupFilePath)
        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Message' -Value "Restoring NTFS Permission successful."

        $result = $true
    }
    Catch
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Error while Restoring NTFS Permission : $_.exception.message"
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
        $successString = "Success: " + "Restoring NTFS Permission Completed"

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
        $JSONResult= $ResultObject|ConvertTo-Json  -Depth 4
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
    if($TargetFolderPath.Length -le 0) 
    {
        $global:Code = 1
        $global:ErrorMessageArray += "TargetFolderPath is not provided"
    }

    if($BackupFilePath.Length -le 0) 
    {
        $global:Code = 1
        $global:ErrorMessageArray += "BackupFilePath is not provided"
    }


    if($global:Code -eq 0)
    {
        if(RestoreNTFSPermission) {}
    }
    
    SetResult
    DisplayResult
}
else
{
    SetResult
    DisplayResult
}
