####################################################
###-------Domain Group Policy Report -----##########
####################################################

<#
[bool]$FullReport = $false
[string]$userInfo = "infinite\mahendran"
#>



##########################################################
###------Variable Declaration-----------------------------
##########################################################

$FullReport = [System.Convert]::ToBoolean($FullReport)

$ComputerName = $env:computername

$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Domain Group Policy Report" 

$SuccessObject = New-Object -TypeName psobject


$global:Code = 0
$global:ErrorMessageArray= @()
$global:SuccessMessageArray= @()

$ReportsObject = New-Object -TypeName psobject

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
###-----Check If Computer is in domaint-------------------
##########################################################

Function IsMachineInDomain()
{
    $result = $false

    #Write-Host "-------------------------------"
    #Write-Host "Checking If machine is part of domain"
    #Write-Host "" 


    try
    {
        if((Get-WmiObject -Class Win32_ComputerSystem).PartOfDomain)
        {
            Get-WmiObject -Class Win32_ComputerSystem -ComputerName $ComputerName |  foreach {

                #write-host -ForegroundColor 8 "`t Machine Name: $($_.Name)"
                #write-host -ForegroundColor 8 "`t Domain Name: $($_.Domain)"

                $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Machine Name' -Value $($_.Name) 
                $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Domain Name' -Value $($_.Domain) 

                $result = $true
            }
        }
        else
        {
            $global:Code = 1
            $global:ErrorMessageArray += "Machine is not a domain member, Local Group Policy items will have to be checked"
        }

    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Error while retriving domain information: $_.exception.message"

        return $result
    }

    return $result
}



##########################################################
###-----retriving Group policy----------------------------
##########################################################

Function GetGroupPolicy()
{
    <#
    Write-Host "-------------------------------"
    if($FullReport)
    {
        Write-Host -ForegroundColor 10 "Retriving full group policy"
    }
    else
    {
        Write-Host -ForegroundColor 10 "Retriving GPO names"
    }
    Write-Host "" 
    #>

    try
    {
        if($FullReport)
        {
            if($userInfo.Length -gt 0)
            {
                $Report = gpresult /z /scope:user /user:$userInfo
                $ReportsObject | Add-Member -MemberType NoteProperty -Name 'Domain Group Policy Report' -Value $Report
            }
            else
            {
                $Report= gpresult /z /scope:user
                $ReportsObject | Add-Member -MemberType NoteProperty -Name 'Domain Group Policy Report' -Value $Report
            }
        }
        else
        {
	        if($userInfo.Length -gt 0)
            {
                $Report = gpresult /r /scope:computer /user:$userInfo
                $ReportsObject | Add-Member -MemberType NoteProperty -Name 'Domain Group Policy Report' -Value $Report
            }
            else
            {
                $Report = gpresult /r /scope:computer
                $ReportsObject | Add-Member -MemberType NoteProperty -Name 'Domain Group Policy Report' -Value $Report
            }
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Error while retriving Group policy, please run the script on the DC to check user specific GPO : $_.exception.message"

        return $result
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
        $successString = "Success: " + "The operation completed successfully"

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value "success"
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value $successString
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $SuccessObject #$global:SuccessMessageArray


        $OutputObject = New-Object -TypeName psobject
        $OutputObject | Add-Member -MemberType NoteProperty -Name 'Domain' -Value $($SuccessObject)

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'dataObject' -Value $ReportsObject
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
    if(IsMachineInDomain)
    {
        GetGroupPolicy
    }

    SetResult
    DisplayResult
}
else
{
    SetResult
    DisplayResult
}
