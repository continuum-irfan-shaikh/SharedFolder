####################################################
###------Get Servers and Roles------------##########
####################################################
 
<#
[string]$Status = "Installed"
[string]$Status = "UnInstalled"
#>

##########################################################
###------Variable Declaration-----------------------------
##########################################################

$ComputerName = $env:computername

$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Get Servers and Roles" 

$SuccessObject = New-Object -TypeName psobject

$global:Code = 0
$global:ErrorMessageArray= @()
$global:SuccessMessageArray= @()

$global:RoleArray = @()

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
###-----Get OS Information--------------------------------
##########################################################

Function GetOSInformation()
{
    try
    {
        <#
        $os = Get-CimInstance Win32_OperatingSystem |select Caption , CSDVersion, CreationClassName -ErrorAction Stop
        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'OS' -Value $($os)
        $global:SuccessMessageArray += "OS Caption: $($os.Caption)"
        $global:SuccessMessageArray += "OS CSDVersion: $($os.CSDVersion)"
        $global:SuccessMessageArray += "OS CreationClassName: $($os.CreationClassName)"
        #>

        $OS =(Get-WMIObject win32_operatingsystem).name
        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'OS' -Value $($OS)
        $global:SuccessMessageArray += "OS : $($OS)"
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Error while retriving OS information: $_.exception.message"

        return $result
    }
}



##########################################################
###-----Is Machine Domain Controller----------------------
##########################################################

Function IsMachineDomainController()
{
    try
    {
        $osInfo = Get-WmiObject -Class Win32_OperatingSystem

        if($osInfo.ProductType -eq 1)
        {
            $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Machine Type' -Value "workstation"
            $global:SuccessMessageArray += "Machine Type: workstation"
        }
        elseif($osInfo.ProductType -eq 2)
        {
            $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Machine Type' -Value "Domain Controller"
            $global:SuccessMessageArray += "Machine Type: Domain Controller"
        }
        elseif($osInfo.ProductType -eq 3)
        {
            $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Machine Type' -Value "Server that is not a Domain Controller"
            $global:SuccessMessageArray += "Machine Type: Server that is not a Domain Controller"
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Error while retriving domain information: $_.exception.message"

        return $result
    }
}


##########################################################
###-----Check If Computer is in domaint-------------------
##########################################################

Function IsMachineInDomain()
{
    $result = $false

    try
    {
        if((Get-WmiObject -Class Win32_ComputerSystem).PartOfDomain)
        {
            Get-WmiObject -Class Win32_ComputerSystem -ComputerName $ComputerName |  foreach {

                #write-host -ForegroundColor 8 "`t Machine Name: $($_.Name)"
                #write-host -ForegroundColor 8 "`t Domain Name: $($_.Domain)"

                $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Machine Name' -Value $($_.Name) 
                $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Domain Name' -Value $($_.Domain) 

                $global:SuccessMessageArray += "Machine Name: $($_.Name)"
                $global:SuccessMessageArray += "Domain Name: $($_.Domain)"

                $result = $true
            }
        }
        else
        {
            $global:Code = 1
            $global:ErrorMessageArray += "Machine is not a domain member, Script must be run on a domain controller"
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
###------Get Server Role ---------------------------------
##########################################################

function GetServerRole
{
    try
    {
        $M = Get-Module -ListAvailable ServerManager
        Import-Module -ModuleInfo $M

        #$Res= Get-WindowsFeature
        #$Res= Get-WindowsFeature |Where Installed
        #$Res= Get-WindowsFeature |select DisplayName,Name,FeatureType,SubFeatures,Installed

        if($Status -eq "Installed")
        {
            $Res= Get-WindowsFeature |Where Installed |select DisplayName,Name,FeatureType,SubFeatures,Installed
        }
        elseif($Status -eq "UnInstalled")
        {
            $Res= Get-WindowsFeature |Where Installed -EQ $false |select DisplayName,Name,FeatureType,SubFeatures,Installed
        }
        else
        {
            $global:Code = 2
            $global:ErrorMessageArray += "Status is not provided"
        }
        

        $global:RoleArray += $Res
    }
    catch
    {
        $global:Code = 1
        $global:ErrorMessageArray += "Error while retriving Server Role: $_.exception.message"
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
        $OutputObject | Add-Member -MemberType NoteProperty -Name 'Result' -Value $($SuccessObject)

        if($global:RoleArray.Count -gt 0)
        {
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'Role' -Value $($global:RoleArray)
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
    IsMachineDomainController

    GetOSInformation

    if(IsMachineInDomain )
    {
        GetServerRole
    }

    SetResult
    DisplayResult 
}
else
{
    SetResult
    DisplayResult
}
