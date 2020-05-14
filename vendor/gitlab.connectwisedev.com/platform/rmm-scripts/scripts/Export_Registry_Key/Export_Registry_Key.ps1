####################################################
###------Export Registry Key--------------##########
####################################################
 
<#
[string] $FilePath = "C:\Temp"
[string] $Filename = "RegBKP.reg"
[string] $RegistryPath = "HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Print"
[string] $Username = ""
[string] $Password = ""
#>


##########################################################
###------Variable Declaration-----------------------------
##########################################################

$ComputerName = $env:computername

$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Export registry key" 

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
###------Take Backup of Registry--------------------------
##########################################################

Function BackupRegistryPath(){

    Write-Host "-------------------------------"
    Write-Host "Taking Registry Backup"
    write-host ""

    $ErrorActionPreference = "stop"
    
    try
    {
        $continue= $true

        $RegistryTestPth = ""
        $RegistryBackupPath = ""
        $UserCredential = $null

        
        #-----Data Validation-----------------------
        if($RegistryPath.Length -le 0 )
        {
            $global:Code = 1
            $global:ErrorMessageArray+= "Registry Path is not provided"
            $continue = $false
        }

        if($FilePath.Length -le 0 )
        {
            $global:Code = 1
            $global:ErrorMessageArray+= "File Path is not provided"
            $continue = $false
        }

        if($Filename.Length -le 0 )
        {
            $global:Code = 1
            $global:ErrorMessageArray+= "File name is not provided"
            $continue = $false
        }


        #---Set Credential-------------------------------------
        if(($Username.Length -gt 0) -and ($Password.Length -gt 0))
        {
            $pass = $Password | ConvertTo-SecureString -asPlainText -Force
            $UserCredential = New-Object System.Management.Automation.PSCredential($Username,$pass)
        }


        #-------Recreate registory Path------------------------------
        if($continue -eq $true)
        {
            $RegistryTestPth = GetTestPath $RegistryPath
            $RegistryBackupPath = GetBackupPath $RegistryPath

            #write-host -ForegroundColor 9 "`t Registry Test Pth: $($RegistryTestPth)" 
            #write-host -ForegroundColor 9 "`t Registry Backup Pth: $($RegistryBackupPath)"

            if(-NOT(Test-Path $RegistryTestPth))
            {
                $global:Code = 1
                $global:ErrorMessageArray+= "Registry Path not Exists: $($RegistryPath)"
                $continue = $false
            }
        }


        #-------Create File Path if not exist------------------------------
        if($continue -eq $true)
        {
            if(-NOT(Test-Path $FilePath))
            {
                write-host -ForegroundColor 10 "`t Creating File Path : $FilePath"

                if($UserCredential -ne $null)
                {
                    New-Item -ItemType Directory -Force -Path $FilePath  -errorAction Stop #-Credential $UserCredential
                }
                else
                {
                    New-Item -ItemType Directory -Force -Path $FilePath  -errorAction Stop 
                }
            }
        }


        if($continue -eq $true)
        {
            if(-NOT(Test-Path $FilePath))
            {
                $global:Code = 1
                $global:ErrorMessageArray+= "File path could not be created"
                $continue = $false
            }
        }

       
        if($continue -eq $true)
        {
            $File= $FilePath + "\"+ $Filename
            #$Res = REG EXPORT $RegistryBackupPath $File /y 

            if(REG EXPORT $RegistryBackupPath $File /y )
            {
                $SuccessObject | Add-Member -MemberType NoteProperty -Name 'RegistryPath' -Value $($RegistryPath)
                $SuccessObject | Add-Member -MemberType NoteProperty -Name 'FilePath' -Value $($FilePath)
                $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Filename' -Value $($Filename)

                $global:SuccessMessageArray += "RegistryPath: $($RegistryPath)"
                $global:SuccessMessageArray += "FilePath: $($FilePath)"
                $global:SuccessMessageArray += "Filename: $($Filename)"
            }
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray+= "Error while taking registry backup: $_.exception.message"
    }
}



Function GetTestPath([string] $RegPath)
{
    $result= ""

    if($RegPath.Contains("HKEY_CLASSES_ROOT"))
    {
        $result=$RegPath.Replace("HKEY_CLASSES_ROOT", "HKCR:")
    }
    if($RegPath.Contains("HKEY_CURRENT_USER"))
    {
        $result=$RegPath.Replace("HKEY_CURRENT_USER", "HKCU:")
    }
    if($RegPath.Contains("HKEY_LOCAL_MACHINE"))
    {
        $result=$RegPath.Replace("HKEY_LOCAL_MACHINE", "HKLM:")
    }
    if($RegPath.Contains("HKEY_USERS"))
    {
        $result=$RegPath.Replace("HKEY_USERS", "HKU:")     
    }
    if($RegPath.Contains("HKEY_CURRENT_CONFIG"))
    {
        $result=$RegPath.Replace("HKEY_CURRENT_CONFIG", "HKCC:")
    }

    if($result.Length -eq 0)
    {
        $result = $RegPath
    }


    return $result
}


Function GetBackupPath([string] $RegPath)
{
    $result= ""

    if($RegPath.Contains("HKCR:"))
    {
        $result=$RegPath.Replace("HKCR:", "HKEY_CLASSES_ROOT")
    }
    if($RegPath.Contains("HKCU:"))
    {
        $result=$RegPath.Replace("HKCU:", "HKEY_CURRENT_USER")
    }
    if($RegPath.Contains("HKLM:"))
    {
        $result=$RegPath.Replace("HKLM:", "HKEY_LOCAL_MACHINE")
    }
    if($RegPath.Contains("HKU:"))
    {
        $result=$RegPath.Replace("HKU:", "HKEY_USERS")     
    }
    if($RegPath.Contains("HKCC:"))
    {
        $result=$RegPath.Replace("HKCC:", "HKEY_CURRENT_CONFIG")
    }

    if($result.Length -eq 0)
    {
        $result = $RegPath
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
        $successString = "Success: " + "The operation completed successfully"

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value "success"
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value $successString
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $global:SuccessMessageArray


        $OutputObject = New-Object -TypeName psobject
        $OutputObject | Add-Member -MemberType NoteProperty -Name 'RegistryBackupPath' -Value $($SuccessObject)

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
    BackupRegistryPath

    SetResult
    DisplayResult 
}
else
{
    SetResult
    DisplayResult
}
