
#######################################################################
###---------------Get Public IP Address----------------------------####
#######################################################################


##########################################################
###------Variable Declaration-----------------------------
##########################################################

$Computer = $env:computername
$global:ips = @()
$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Get Public IP Address" 

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



####################################################
###------Load form and get resullt--------##########
####################################################
 
 Function Find-PublicIP{
    try{

        if($PSVersionTable.PSVersion.Major -eq 2){
        $global:ips = @()
        $url = "http://checkip.dyndns.com"
        $WebRequest = [System.Net.WebRequest]::Create($url)
        $WebRequest.Method = "GET"
        $WebRequest.ContentType = "application/json"
        $Response = $WebRequest.GetResponse()
        $ResponseStream = $Response.GetResponseStream()
        $ReadStream = New-Object System.IO.StreamReader $ResponseStream
        $Data=$ReadStream.ReadToEnd()
        $HtmlObject = New-Object -ComObject "HTMLfile"
        $HtmlObject.IHTMLDocument2_Write($Data)
        $result = $HtmlObject.body.innerHTML
        $global:ip = $result.Split(":")[1].Trim()
        #$global:ipadd = $global:ip.IPAddressToString
        Write-Host "The public IP address of the local computer is $global:ip"
        }
        if($PSVersionTable.PSVersion.Major -ne 2.0){
         $global:ip = ((Invoke-WebRequest ifconfig.me/ip).content.Trim())
        }
    
    }
    catch{
        $global:Code = 1
        $global:ErrorMsg = $_.Exception.Message
        $global:FailedItem = $_.Exception.ItemName
        $global:ErrorMessageArray += "Unable to obtain public IP address(es) due to error $global:FailedItem :  $global:ErrorMsg" 

    }
                        
    
        $global:ips += $global:ip
    
    #$global:ips
    
    #$global:ips += "255.255.255.255"
    #$global:ips = $null
    
    $global:count = $global:ips.count
    
    if($global:count -eq 1){
        $computer = $env:COMPUTERNAME
        $global:ipv = Test-Connection -ComputerName $computer -Count 1 |             
        Select-Object -Property IPv*Address
        $global:ipv4= $ipv.IPv4Address.IPAddressToString 
        $global:ipv6= $ipv.IPv6Address
        $global:SuccessMessageArray += "Success Public IP Obtained (Single)"
        #$global:ipadd = $global:ips.IPAddressToString        
        $global:SuccessObject = New-Object -TypeName psobject
        $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'Message:' -Value $($global:SuccessMessageArray)
        $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'Public IP:' -Value $($global:ip)
        $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'IPV4 Address:' -Value $($global:ipv4)
        $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'IPV6 Address:' -Value $([string]$global:ipv6)
        $global:SuccessMessageArray += "Public IP:" + $("$global:ip")
        $global:SuccessMessageArray += "IPV4 Address:" + $("$global:ipv4")
        $global:SuccessMessageArray += "IPV6 Address:" + $($global:ipv6) 
    }

    if($global:count -gt 1){
        $address = gwmi Win32_NetworkAdapterConfiguration |
        Where { $_.IPAddress } | # filter the objects where an address actually exists
        Select -Expand IPAddress
    
        $global:ipv6= $address -match "::"
        $global:ipv4= $address -notmatch "::"
        $global:SuccessMessageArray += "Success Public IP Obtained (Multiple)"
                
        $global:SuccessObject = New-Object -TypeName psobject
        $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'Message:' -Value $($global:SuccessMessageArray)
        $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'Public IP:' -Value $($global:ips)
        $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'IPV4 Address:' -Value $($global:ipv4)
        $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'IPV6 Address:' -Value $($global:ipv6)
        $global:SuccessMessageArray += "Public IP:" + $($global:ips)
        $global:SuccessMessageArray += "IPV4 Address:" + $($global:ipv4)
        $global:SuccessMessageArray += "IPV6 Address:" + $($global:ipv6) 
    }
    if($global:count -eq $null){
        $global:Code = 1
        $global:ErrorMsg = $_.Exception.Message
        $global:FailedItem = $_.Exception.ItemName
        $global:ErrorMessageArray += "Unable to obtain public IP address(es) due to error: Public IP not found" 

    
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
        $successString = "Success: " + "Information capture was successful"

        <#
        $stdoutMessageString= "message :"
        foreach($SuccessMessage in $global:SuccessMessageArray)
        {
            $stdoutMessageString += $SuccessMessage + ", "
        }
        #>

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
    Find-PublicIP
    SetResult
    DisplayResult
}
else
{
    Find-PublicIP
    SetResult
    DisplayResult
}



