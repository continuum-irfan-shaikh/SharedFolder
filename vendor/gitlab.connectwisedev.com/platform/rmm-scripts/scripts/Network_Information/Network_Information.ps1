####################################################
###-------Network Information----------##########
####################################################
 
<#
[bool] $ActiveNetworkOnly= $true
[bool] $ActiveWirelessOnly= $true
[bool] $FullProfile= $false
#>


##########################################################
###------Variable Declaration-----------------------------
##########################################################

$ActiveNetworkOnly = [System.Convert]::ToBoolean($ActiveNetworkOnly)
$ActiveWirelessOnly = [System.Convert]::ToBoolean($ActiveWirelessOnly)
$FullProfile = [System.Convert]::ToBoolean($FullProfile)

$ComputerName = $env:computername



$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Network Information" 

$SuccessObject = New-Object -TypeName psobject

$global:Code = 0
$global:ErrorMessageArray= @()
$global:SuccessMessageArray= @()

$global:AdapterOutputArray = @()
$global:WirelessOutputArray = @()
$global:ProfileOutputArray = @()


$global:FullProfileOutputArray = @()


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
###------Get Network Information--------------------------
##########################################################

Function GetNetworkInformation
{
    #Write-Host "-------------------------------"
    #write-host -ForegroundColor 10 "Retriving Network Information"
    #Write-Host "" 
    
    try
    {
        if($ActiveNetworkOnly -eq $true)
        {
            Get-WmiObject -Class Win32_NetworkAdapterConfiguration -Filter IPEnabled=TRUE  -ComputerName  $ComputerName  -ErrorAction Stop | Select-Object -Property [a-z]* -ExcludeProperty IPX*,WINS*| foreach{

                $AdapterOutput = New-Object -TypeName psobject 
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'Adapter Name' -Value $($_.Description)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'DNS Suffix' -Value $($_.DNSDomainSuffixSearchOrder)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'DHCP Status' -Value $($_.DHCPEnabled)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'Mac Address' -Value $($_.MACAddress)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'IP Address' -Value $($_.IPAddress)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'Subnet Mask' -Value $($_.IPSubnet)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'Default Gateway' -Value $($_.DefaultIPGateway)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'DNS Servers' -Value $($_.DNSServerSearchOrder)

                $global:AdapterOutputArray += $AdapterOutput 

                $global:SuccessMessageArray += "Adapter Name: $($_.Description)"
                $global:SuccessMessageArray += "DNS Suffix: $($_.DNSDomainSuffixSearchOrder)"
                $global:SuccessMessageArray += "DHCP Status: $($_.DHCPEnabled)"
                $global:SuccessMessageArray += "Mac Address: $($_.MACAddress)"
                $global:SuccessMessageArray += "IP Address: $($_.IPAddress)"
                $global:SuccessMessageArray += "Subnet Mask: $($_.IPSubnet)"
                $global:SuccessMessageArray += "Default Gateway: $($_.DefaultIPGateway)"
                $global:SuccessMessageArray += "DNS Servers: $($_.DNSServerSearchOrder)"

                $global:SuccessMessageArray += "---------------------"

            }
        }
        else
        {

            Get-WmiObject -Class Win32_NetworkAdapterConfiguration  -ComputerName  $ComputerName  -ErrorAction Stop | Select-Object -Property [a-z]* -ExcludeProperty IPX*,WINS* | Sort-Object -Property @{Expression = "IPAddress"; Descending = $true},  @{Expression = "MACAddress"; Descending = $true}, @{Expression = "DNSServerSearchOrder"; Descending = $true}| foreach{

                $AdapterOutput = New-Object -TypeName psobject
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'Adapter Name' -Value $($_.Description)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'DNS Suffix' -Value $($_.DNSDomainSuffixSearchOrder)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'DHCP Status' -Value $($_.DHCPEnabled)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'Mac Address' -Value $($_.MACAddress)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'IP Address' -Value $($_.IPAddress)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'Subnet Mask' -Value $($_.IPSubnet)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'Default Gateway' -Value $($_.DefaultIPGateway)
                $AdapterOutput | Add-Member -MemberType NoteProperty -Name 'DNS Servers' -Value $($_.DNSServerSearchOrder)

                $global:AdapterOutputArray += $AdapterOutput

                $global:SuccessMessageArray += "Adapter Name: $($_.Description)"
                $global:SuccessMessageArray += "DNS Suffix: $($_.DNSDomainSuffixSearchOrder)"
                $global:SuccessMessageArray += "DHCP Status: $($_.DHCPEnabled)"
                $global:SuccessMessageArray += "Mac Address: $($_.MACAddress)"
                $global:SuccessMessageArray += "IP Address: $($_.IPAddress)"
                $global:SuccessMessageArray += "Subnet Mask: $($_.IPSubnet)"
                $global:SuccessMessageArray += "Default Gateway: $($_.DefaultIPGateway)"
                $global:SuccessMessageArray += "DNS Servers: $($_.DNSServerSearchOrder)"

                $global:SuccessMessageArray += "---------------------"

            }
        }

    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Error while retriving network information: $_.exception.message"

        return $result
    }
}




##########################################################
###------Get All Wireless Network Information-------------
##########################################################

function Get-AllWifiNetwork 
{
    #Write-Host "-------------------------------"
    #write-host -ForegroundColor 10 "Retriving wireless Information"
    #Write-Host "" 

    try
    {

      netsh wlan sh net mode=bssid | % -process {
        if ($_ -match '^SSID (\d+) : (.*)$') {
            $current = @{}
            $networks += $current
            $current.Index = $matches[1].trim()
            $current.SSID = $matches[2].trim()
        
        } 
        else 
        {
            if ($_ -match '^\s+(.*)\s+:\s+(.*)\s*$') {
                $current[$matches[1].trim()] = $matches[2].trim()
            }
        }
      } -begin { $networks = @() } -end { $networks|% { 
  
        $network= new-object psobject -property $_ 
    
            $WirelessOutput = New-Object -TypeName psobject
            $WirelessOutput | Add-Member -MemberType NoteProperty -Name 'SSID' -Value $($network.SSID)
            $WirelessOutput | Add-Member -MemberType NoteProperty -Name 'Signal' -Value $($network.Signal)
            $WirelessOutput | Add-Member -MemberType NoteProperty -Name 'Authentication' -Value $($network.Authentication)
            #$WirelessOutput | Add-Member -MemberType NoteProperty -Name 'Cipher' -Value $($network.Cipher)
            $WirelessOutput | Add-Member -MemberType NoteProperty -Name 'Cipher' -Value $($network.Encryption)
            $global:WirelessOutputArray  += $WirelessOutput

            $global:SuccessMessageArray += "SSID: $($network.SSID)"
            $global:SuccessMessageArray += "Signal: $($network.Signal)"
            $global:SuccessMessageArray += "Authentication: $($network.Authentication)"
            $global:SuccessMessageArray += "Cipher: $($network.Encryption)"
            $global:SuccessMessageArray += "---------------------"
   
        } }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Error while retriving Connected wireless Information: $_.exception.message"
    }
  
 }



##########################################################
###------Get Connected Wireless Network Information-------
##########################################################

Function Get-Connected-WifiNetwork
{
    #Write-Host "-------------------------------"
    #write-host -ForegroundColor 10 "Getting Connected wireless Information"
    #Write-Host "" 

    try
    {
        $network = netsh wlan show interfaces | Select-String -Pattern SSID , Signal, Authentication, Cipher | %{ ($_ -split ":")[-1].Trim() };

        if($network.length -gt 0)
        {
            

            $WirelessOutput = New-Object -TypeName psobject
            $WirelessOutput | Add-Member -MemberType NoteProperty -Name 'SSID' -Value $($network[0])
            $WirelessOutput | Add-Member -MemberType NoteProperty -Name 'Authentication' -Value $($network[2])
            $WirelessOutput | Add-Member -MemberType NoteProperty -Name 'Signal' -Value $($network[4])
            #$WirelessOutput | Add-Member -MemberType NoteProperty -Name 'Cipher' -Value $($network[3])
            $WirelessOutput | Add-Member -MemberType NoteProperty -Name 'Cipher' -Value $($network[3])
            $global:WirelessOutputArray  += $WirelessOutput


            $global:SuccessMessageArray +="SSID: $($network[0])"
            $global:SuccessMessageArray +="Authentication: $($network[2])"
            $global:SuccessMessageArray +="Signal: $($network[4])"
            $global:SuccessMessageArray +="Cipher: $($network[3])"
            $global:SuccessMessageArray += "---------------------"
        }
    }
    catch
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Error while retriving Connected wireless Information: $_.exception.message"
    }
}


##########################################################
###------Get saved profile Information--------------------
##########################################################

Function GetsavedProfileInformation
{
    #Write-Host "-------------------------------"
    #write-host -ForegroundColor 10 "Getting saved profile Information"
    #Write-Host "" 

    try
    {
        $listProfiles = netsh wlan show profiles | Select-String -Pattern "All User Profile" | %{ ($_ -split ":")[-1].Trim() };

        $listProfiles | foreach {

            $profileInfo = netsh wlan show profiles name=$_ key="clear";

	        $SSID = $profileInfo | Select-String -Pattern "SSID Name" | %{ ($_ -split ":")[-1].Trim() };
	        $Key = $profileInfo | Select-String -Pattern "Key Content" | %{ ($_ -split ":")[-1].Trim() };

            if($SSID.Length -gt 0)
            {
                $ProfileOutput = New-Object -TypeName psobject
                $ProfileOutput | Add-Member -MemberType NoteProperty -Name 'Profile Name' -Value $($SSID)
                $ProfileOutput | Add-Member -MemberType NoteProperty -Name 'Password' -Value $($Key)
                $global:ProfileOutputArray  += $ProfileOutput

                $global:SuccessMessageArray += "Profile Name: $($SSID)"
                $global:SuccessMessageArray += "Password: $($Key)"
                $global:SuccessMessageArray += "---------------------"
            }

            if($FullProfile -eq $true)
            {
                $FullInfo=  netsh wlan show profiles name=$($_) key="clear";
 
                $global:FullProfileOutputArray += $FullInfo
            }
        }
    }
    catch
    {
        $global:Status = 2
        $global:Message= "Error while Getting saved progile Information: $_.exception.message"
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
    

    if($global:Code -eq 0)
    {
        $successString = "Success: " + "The operation completed successfully"


        #------Change the Code if any of the three fails-----------------
        $FailMSG = ""
        if($global:AdapterOutputArray.Count -le 0)
        {
            $global:Code = 1
            $FailMSG += "No Ethernet adapter found, "
        }
        if($global:WirelessOutputArray.Count -le 0)
        {
            $global:Code = 1
            $FailMSG += "No wireless adapter found, "
        }
        if($FullProfile -eq $false)
        {
            if($global:ProfileOutputArray.Count -le 0)
            {
                $global:Code = 1
                $FailMSG += "No Profile found, "
            }
        }
        if($FullProfile -eq $true)
        {
            if($global:FullProfileOutputArray.Count -le 0)
            {
                $global:Code = 1
                $FailMSG += "No full Profile found "
            }
        }
        
        if( $FailMSG.Length -gt 0)
        {
            $successString = "Fail: " + $FailMSG
        }

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Code' -Value $global:Code

        if($global:Code -eq 0)
        {
            $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value "success"
        }
        else
        {
            $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value "Fail"
        }


        #----------------------------------------------------------------
        
        $TotalObject = $global:AdapterOutputArray.Count + $global:WirelessOutputArray.Count + $global:ProfileOutputArray.Count

        

        


        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value $successString
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Objects' -Value $TotalObject

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $global:SuccessMessageArray


        $OutputObject = New-Object -TypeName psobject

        if($global:AdapterOutputArray.Count -gt 0)
        {
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'Ethernet Adapter' -Value $($global:AdapterOutputArray)
        }

        if($global:WirelessOutputArray.Count -gt 0)
        {
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'Wireless Network' -Value $($global:WirelessOutputArray)
        }

        if($global:ProfileOutputArray.Count -gt 0)
        {
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'Profile' -Value $($global:ProfileOutputArray)
        }

        if($global:FullProfileOutputArray.Count -gt 0)
        {
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'FullProfile' -Value $($global:FullProfileOutputArray)
        }

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'dataObject' -Value $OutputObject


        #------------Error with Success--------------------------------------------
        $errorCount= 0
        $ErrorObjectAray= @()

        if($global:AdapterOutputArray.Count -le 0)
        {
            $ErrObject = New-Object -TypeName psobject
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'id' -Value $errorCount
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'title' -Value "No Ethernet adapter found"
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'detail' -Value "No Ethernet adapter found"
            $ErrorObjectAray += $ErrObject

            $errorCount = $errorCount +1
        }
        if($global:WirelessOutputArray.Count -le 0)
        {
            $ErrObject = New-Object -TypeName psobject
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'id' -Value $errorCount
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'title' -Value "No wireless adapter found"
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'detail' -Value "No wireless adapter found"
            $ErrorObjectAray += $ErrObject

            $errorCount = $errorCount +1
        }
        if($FullProfile -eq $false)
        {
            if($global:ProfileOutputArray.Count -le 0)
            {
                $ErrObject = New-Object -TypeName psobject
                $ErrObject | Add-Member -MemberType NoteProperty -Name 'id' -Value $errorCount
                $ErrObject | Add-Member -MemberType NoteProperty -Name 'title' -Value "No Profile found"
                $ErrObject | Add-Member -MemberType NoteProperty -Name 'detail' -Value "No Profile found"
                $ErrorObjectAray += $ErrObject

                $errorCount = $errorCount +1
            }
        }
        if($FullProfile -eq $true)
        {
            if($global:FullProfileOutputArray.Count -le 0)
            {
                $ErrObject = New-Object -TypeName psobject
                $ErrObject | Add-Member -MemberType NoteProperty -Name 'id' -Value $errorCount
                $ErrObject | Add-Member -MemberType NoteProperty -Name 'title' -Value "No full Profile found"
                $ErrObject | Add-Member -MemberType NoteProperty -Name 'detail' -Value "No full Profile found"
                $ErrorObjectAray += $ErrObject

                $errorCount = $errorCount +1
            }
        }


        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stderr' -Value $ErrorObjectAray
    }
    else
    {
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Code' -Value $global:Code
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
    GetNetworkInformation


    if($ActiveWirelessOnly -eq $true)
    {
        Get-Connected-WifiNetwork
    }
    else
    {
        Get-AllWifiNetwork| select ssid, signal,Authentication, Encryption |sort signal -desc |Format-Table
    }

    GetsavedProfileInformation
   

    SetResult
    DisplayResult

}
else
{
    SetResult
    DisplayResult
}
