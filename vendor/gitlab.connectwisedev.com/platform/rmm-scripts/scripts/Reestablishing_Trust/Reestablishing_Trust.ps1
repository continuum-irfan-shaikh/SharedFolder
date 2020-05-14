###################################################
###-------Reestablishing Trust -------------########
####################################################
 
<#
$username = "Administrator"
$Password = "Abcd@1234"
#>



##########################################################
###------Variable Declaration-----------------------------
##########################################################

$timeoutSeconds = 30
$ReconnectAfterFail=$true
$ComputerName = $env:computername


$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Reestablishing Trust" 

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
###------ Check Connectivity With ADServer----------------
##########################################################

Function Check_Connectivity_With_ADServer()
{
    $IsContinued = $false

    Write-Host "-------------------------------"
    Write-Host "Verifying connectivity"
    Write-Host "    " 


    try
    {

        Get-WmiObject -Class Win32_ComputerSystem | Select Domain |ForEach-Object{
            $domainName = $_.Domain
        }
        write-host -ForegroundColor 8 "`t Domain Name : $domainName"
        write-host ""

        if (Test-Connection -ComputerName $domainName -Quiet)
        {
            write-host -ForegroundColor 8 "`t Connection to domain server successful"
            $IsContinued = $true
            return $IsContinued 
        }
        else
        {
            $global:Code = 1
            $global:ErrorMessageArray += "Connection to domain server failed $($domainName)"

            $IsContinued = $false
            return $IsContinued 
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Error while Connecting AD Server: $_.exception.message"

        return $IsContinued 
    }
}


##########################################################
###------ Get DNS Configuration---------------------------
##########################################################

Function Get_DNS_Configuration()
{
    #Write-Host "-------------------------------"
    #Write-Host "Getting DNS Configuration"

    try
    {
        $DNSInfo = Get-WmiObject -Class Win32_ComputerSystem

        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Domain' -Value $($DNSInfo.Domain)
        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Manufacturer' -Value $($DNSInfo.Manufacturer)
        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Model' -Value $($DNSInfo.Model)
        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Name' -Value $($DNSInfo.Name)
        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'PrimaryOwnerName' -Value $($DNSInfo.PrimaryOwnerName)
        $SuccessObject | Add-Member -MemberType NoteProperty -Name 'TotalPhysicalMemory' -Value $($DNSInfo.TotalPhysicalMemory)

        $global:SuccessMessageArray += "Domain: $($DNSInfo.Domain)"
        $global:SuccessMessageArray += "Manufacturer: $($DNSInfo.Manufacturer)"
        $global:SuccessMessageArray += "Model: $($DNSInfo.Model)"
        $global:SuccessMessageArray += "Name: $($DNSInfo.Name)"
        $global:SuccessMessageArray += "PrimaryOwnerName: $($DNSInfo.PrimaryOwnerName)"
        $global:SuccessMessageArray += "TotalPhysicalMemory: $($DNSInfo.TotalPhysicalMemory)"


        Get-WmiObject -Class Win32_NetworkAdapterConfiguration -Filter IPEnabled=TRUE  -ComputerName  $ComputerName  -ErrorAction Stop | Select-Object -Property [a-z]* -ExcludeProperty IPX*,WINS*| foreach{
        	#write-host "DNS Servers: $($_.DNSServerSearchOrder)"
            $SuccessObject | Add-Member -MemberType NoteProperty -Name 'DNS Servers' -Value $($_.DNSServerSearchOrder)

            $global:SuccessMessageArray += "DNS Servers: $($_.DNSServerSearchOrder)"
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Error while Getting DNS Configuration: $_.exception.message"
    }
}



##########################################################
###------ Repair Trust------------------------------------
##########################################################

Function Repair_Trust()
{
    $Result = 0
    
    Write-Host "-------------------------------"
    Write-Host "Repairing Trust"
     Write-Host ""

    try
    {
        $Pass = $Password | ConvertTo-SecureString -asPlainText -Force
        $credential = New-Object System.Management.Automation.PSCredential($username,$Pass)

        if(Test-ComputerSecureChannel -Credential $credential  -repair)
        {
            write-host -ForegroundColor 8 "`t Repairing Trust successful"
            $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Message' -Value "Repairing Trust successful"
            $global:SuccessMessageArray += "Message: Repairing Trust successful"

            $Result = 1
        }
        else
        {
            write-host -ForegroundColor 8 "`t Repairing Trust faild"

            $global:Code = 1
            $global:ErrorMessageArray += "Repairing Trust faild"
            $Result = 0
        }
    }
    catch [System.Security.Authentication.InvalidCredentialException] 
    {
        write-host -ForegroundColor 8 "`t Invalid Credentials!"

        $global:Code = 1
        $global:ErrorMessageArray += "Invalid Credentials"

        $Result = 2
        return $Result
    }
    catch [System.Security.Authentication.AuthenticationException] 
    {
        write-host -ForegroundColor 8 "`t Invalid Credentials"

        $global:Code = 1
        $global:ErrorMessageArray += "Invalid Credentials"

        $Result = 2
        return $Result
    }
    catch
    {
        IF ($Error[0].Exception.ToString().contains('Invalid username or password'))
        {
            write-host -ForegroundColor 8 "`t Invalid Credentials"

            $global:Code = 1
            $global:ErrorMessageArray += "Invalid Credentials"

            $Result = 2
        }
        elseIF ($Error[0].Exception.ToString().contains('The user name or password is incorrect'))
        {
            write-host -ForegroundColor 8 "`t Invalid Credentials"

            $global:Code = 1
            $global:ErrorMessageArray += "Invalid Credentials"

            $Result = 2
        }
        elseIF ($Error[0].Exception.ToString().contains('server is not operational'))
        {
            write-host -ForegroundColor 8 "`t server is not operational"

            $global:Code = 1
            $global:ErrorMessageArray += "server is not operational"

            $Result = 3
        }
        else
        {
             Write-Warning "Error while Repairing Trust: $_.exception.message"

            $global:Code = 2
            $global:ErrorMessageArray += "Error while Repairing Trust: $_.exception.message"

            $Result = 0
        }

       return $Result 
    }

    return $Result 
}


##########################################################
###------ Repair Trust for PS2----------------------------
##########################################################

Function Repair_TrustForPS2()
{
    Write-Host "-------------------------------"
    Write-Host "Repairing Trust by rejoining"
    
    $IsContinued = $false
    
    try
    {
        Get-WmiObject -Class Win32_ComputerSystem | Select Domain |ForEach-Object{
            $domainName = $_.Domain
        }
        #write-host -ForegroundColor 8 "`t Domain Name : $domainName"
        write-host ""
        
        $Count= $domainName.Length
        #write-host -ForegroundColor 8 "`t Count : $Count"
        
        if($Count -gt 0)
        {
            #-----Create User Credential----------------------
            $user = $domainName +"\" + $username
            $Pass = $Password | ConvertTo-SecureString -asPlainText -Force
            $credential = New-Object System.Management.Automation.PSCredential($user,$Pass)
            
            #write-host -ForegroundColor 8 "`t User Name : $user"
            #write-host -ForegroundColor 8 "`t Password : $Password"
            
            $Continue = $false
            try
            {
                #write-host -ForegroundColor 8 "`t Checking Connectivity to domain"
                Add-Computer -Credential $credential -DomainName $domainName  -errorAction stop    #-verbose
                $Continue = $true
                
                write-host ""
                $IsContinued=$true
            }
            catch
            {
                IF ($_.exception.message.ToString().contains('Access is denied'))
                {
                    $global:Code = 1
                    $global:ErrorMessageArray += "Invalid Credentials"
                }
                elseIF ($_.exception.message.ToString().contains('network path was not found'))
                {
                    $global:Code = 1
                    $global:ErrorMessageArray += "server is not operational"
                }
                else
                {
                    $global:Code = 1
                    $global:ErrorMessageArray += "Error while Repairing Trust : $_.exception.message"
                }

		        $IsContinued = $false
		        return $IsContinued
            }
            
            
            if($Continue)
            {
                #----Remove Computer from Domain-----
                write-host -ForegroundColor 8 "`t Removing computer from domain"
                Remove-Computer -credential $credential -passthru  #-verbose
                write-host ""
                
                
                #----Add Computer to Domain-----
                write-host -ForegroundColor 8 "`t Adding computer to domain"
                Add-Computer -Credential $credential -DomainName $domainName -passthru  #-verbose
                write-host ""

                $SuccessObject | Add-Member -MemberType NoteProperty -Name 'Message' -Value "Repairing Trust successful"
                $global:SuccessMessageArray += "Message: Repairing Trust successful"
            }

        }
        else
        {
            $global:Code = 1
            $global:ErrorMessageArray += "Error while Repairing Trust: Domain Name not found"

	        $IsContinued = $false
	        return $IsContinued
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Error while Repairing Trust : $_.exception.message"

	    $IsContinued = $false
	    return $IsContinued
    }
    
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
###------Set Result --------------------------------------
##########################################################

function SetResult 
{
    $ResultObject | Add-Member -MemberType NoteProperty -Name 'Code' -Value $global:Code

    if($global:Code -eq 0)
    {
        $successString = "Success: " + "Repairing Trust successful"

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
        
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value "Error: Repairing Trust was not successful" #$errorString

        if($global:SuccessMessageArray.Length -gt 0)
        {
            $OutputObject = New-Object -TypeName psobject
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'DNSConfiguration' -Value $SuccessObject

            $ResultObject | Add-Member -MemberType NoteProperty -Name 'dataObject' -Value $OutputObject
        }


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
    if(Check_Connectivity_With_ADServer)
    {
        if($PSVersionTable.PSVersion.Major -eq 2)
        {
            if(Repair_TrustForPS2)
            {
                write-host -ForegroundColor 8 "`t Repairing Trust Completed"
            }
            else
            {
                Get_DNS_Configuration
            }
        }
        else
        {
            $Res = Repair_Trust

            if($Res -eq 1)
            {

            }
            else
            {
                Get_DNS_Configuration
            
                if($Res -eq 0)
                {
                    if($ReconnectAfterFail -eq $true)
                    {
                        Write-host ""

                        if(Repair_TrustForPS2) {}
                    }
                }
            }
        }
    }
    else
    {
        Get_DNS_Configuration
    }


    SetResult
    DisplayResult
}
else
{
    SetResult
    DisplayResult
}