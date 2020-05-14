####################################################
###------get set network category---------##########
####################################################
 
 #---------- $OperationType value can be Read or Edit-------
 
<#
[string]$OperationType = "Read"
[string]$AdapterName = "abcd"
[string]$NetworkCategory = "private"
 #>


##########################################################
###------Variable Declaration-----------------------------
##########################################################

$ComputerName = $env:computername

$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Get Set Network Category" 

$SuccessObject = New-Object -TypeName psobject

$global:Code = 0
$global:ErrorMessageArray= @()
$global:SuccessMessageArray= @()

$global:AdapterInfoObjectArray= @()
$global:EditedAdapterInfoObjectArray= @()

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
###------Read Network ------------------------------------
##########################################################

Function Read-Network()
{
    $global:adaptername = $null
    $global:network_category = $null

    if($OperationType -eq "Read")
    {
        if(![string]::IsNullOrEmpty($global:adaptername) -and ![string]::IsNullOrEmpty($global:network_category) )
        {
            $global:code = 1
            $global:ErrorMessageArray += "The adapter name and network category parameters are not required for type Read"
        }

        if(![string]::IsNullOrEmpty($global:adaptername) -and [string]::IsNullOrEmpty($global:network_category) )
        {
            $global:code = 1
            $global:ErrorMessageArray += "The adapter name is not required for type Read"
        }

        if([string]::IsNullOrEmpty($global:adaptername) -and ![string]::IsNullOrEmpty($global:network_category) )
        {
            $global:code = 1
            $global:ErrorMessageArray += "The network category is not required for type Read" 
        }

        
        if([string]::IsNullOrEmpty($global:adaptername) -and [string]::IsNullOrEmpty($global:network_category) )
        {
            try
            {
                $global:adpname = $null
                $global:netcategory = $null

                if($PSVersionTable.PSVersion.Major -eq 2)
                {
                    $global:networkListManager = [Activator]::CreateInstance([Type]::GetTypeFromCLSID([Guid]"{DCB00C01-570F-4A9B-8D69-199FDBA5723B}"))
                    $global:connections = $global:networkListManager.GetNetworkConnections()
                    
                    foreach($global:conn in $global:connections)
                    {
                        $global:nett = $global:conn.GetNetwork().GetCategory() 
                        $global:nettype = $global:nett | ?{$_ -match "1" -or $_ -match "0"}
                        
                        if(![string]::IsNullOrEmpty($global:nettype))
                        {
                            if($global:nettype -eq 1)
                            {
                                $global:adpname = $global:conn.GetNetwork().GetName() 
                                $global:netcategory = "Private"
                            }

                            if([int]$global:nettype -eq 0)
                            {
                                $global:adpname = $global:conn.GetNetwork().GetName() 
                                $global:netcategory = "Public"
                            }
                        }
                
                        $AdapterInfoObject = New-Object -TypeName psobject
                        $AdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'Adapter Name:' -Value $($global:adpName)
                        $AdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'Network Type:' -Value $($global:netcategory)
                        
                        $global:AdapterInfoObjectArray += $AdapterInfoObject

                        $global:SuccessMessageArray += "Adapter Name:" + $($global:adpName)
                        $global:SuccessMessageArray += "Network Type:" + $($global:netcategory)
                    }
                }

                if($PSVersionTable.PSVersion.Major -ge 3)
                {
                    $global:netconn = Get-NetConnectionProfile
                     
                    $AdapterInfoObject = New-Object -TypeName psobject    
                    $AdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'Adapter Name:' -Value $($global:netconn.Name)
                    $AdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'Network Type:' -Value $([string]$global:netconn.NetworkCategory)
                    $global:AdapterInfoObjectArray += $AdapterInfoObject
                    
                    $global:SuccessMessageArray += "Adapter Name:" + $($global:netconn.Name)
                    $global:SuccessMessageArray += "Network Type:" + $($global:netconn.NetworkCategory)
                }
            }
            catch
            {
                $global:Code = 1
                $global:ErrorMessageArray += "Unable to obtain Adapter Name or Network Type: $_.Exception.Message" 
            }
        }
    }
}


##########################################################
###------Edit Network ------------------------------------
##########################################################

Function Edit-Network($adaptername, $network_category)
{
    
    $global:adaptername =  $adaptername
    $global:network_category =  $network_category

    $global:sol = @()
    $global:networkListManager = [Activator]::CreateInstance([Type]::GetTypeFromCLSID([Guid]"{DCB00C01-570F-4A9B-8D69-199FDBA5723B}"))
    $global:connections = $global:networkListManager.GetNetworkConnections()
                        
    foreach($global:conn in $global:connections){
        $global:name = $global:conn.GetNetwork().GetName()
        $global:category = $global:conn.GetNetwork().GetCategory()
                 
        if($global:category -eq 0){
            $global:net = "Public"  
        }
        if($global:category -eq 1){
            $global:net = "Private"                  
        }
        $global:sol += $global:name + ":" + $global:net
    }

    if([string]::IsNullOrEmpty($global:adaptername) -and [string]::IsNullOrEmpty($global:network_category) )
    {
        $global:code = 2
        $global:ErrorMessageArray += "The Adapter name and Network category is mandatory.  The network adapter found in this machine are $global:sol"        
    }
    if(![string]::IsNullOrEmpty($global:adaptername) -and [string]::IsNullOrEmpty($global:network_category) )
    {
        $global:code = 2
        $global:ErrorMessageArray += "The Network category is mandatory.  The network adapter found in this machine are $global:sol"        
    }
    if([string]::IsNullOrEmpty($global:adaptername) -and ![string]::IsNullOrEmpty($global:network_category) )
    {
        $global:code = 2
        $global:ErrorMessageArray += "The Adapter name is mandatory.  The network adapter found in this machine are $global:sol"        
    }
   

    if(![string]::IsNullOrEmpty($global:adaptername) -and ![string]::IsNullOrEmpty($global:network_category) )
    {
         if($global:network_category -ne "private" -and $global:network_category -ne "public")
        {
            $global:code = 2
            $global:ErrorMessageArray += "The network category should be either Private or Public.  The network adapter found in this machine are $global:sol"        
        }
        if($global:network_category -eq "Private" -or $global:network_category -eq "Public")
        {
            if($PSVersionTable.PSVersion.major -eq 2)
            {
                $global:networkListManager = [Activator]::CreateInstance([Type]::GetTypeFromCLSID([Guid]"{DCB00C01-570F-4A9B-8D69-199FDBA5723B}"))
                $global:connections = $global:networkListManager.GetNetworkConnections()
                        
                foreach($global:conn in $global:connections)
                {
                    $global:name = $global:conn.GetNetwork().GetName() | ?{ $_ -match $global:adaptername }
         
                    if($global:name -notmatch $global:adaptername)
                    { 
                       $global:code = 2
                       $global:ErrorMessageArray += "Invalid adapter name $global:adaptername. The network adapter found in this machine are `n $global:sol"
                    } 
         
                    if($global:name -match $global:adaptername)
                    {
         
                        if($global:network_category -eq "Private")
                        {
                            $global:conn.GetNetwork().setCategory(1) |?{$global:conn.GetNetwork().GetName() -match $global:adaptername }    
                        }

                        if($global:network_category -eq "Public")
                        {
                            $global:conn.GetNetwork().setCategory(0) |?{$global:conn.GetNetwork().GetName() -match $global:adaptername }    
                        }
                    }
                }

                $global:networkListManager = $null
                $global:networkListManager = [Activator]::CreateInstance([Type]::GetTypeFromCLSID([Guid]"{DCB00C01-570F-4A9B-8D69-199FDBA5723B}"))
                $global:connections = $global:networkListManager.GetNetworkConnections()
                $ips = arp -a |Select-String "Interface" | ForEach-Object {$_ -split ":" -split "---" | select -Index 1}                        
           
                foreach($ip in $ips)
                {
                    foreach($global:conn in $global:connections)
                    {
                        $global:check = $global:conn.GetNetwork().GetCategory() |?{$global:conn.GetNetwork().GetName() -match $global:adaptername }
                    
                        if($global:check -match $global:netcategory)
                        {
                            $global:code = 1
                            $global:ErrorMessageArray = "$global:adaptername is already set to $global:networkcategory"
                        }
                            
                        $global:nett = $global:conn.GetNetwork().GetCategory() |?{$global:conn.GetNetwork().GetName() -match $global:adaptername }
                        $global:nettype = $global:nett | ?{$_ -match "1" -or $_ -match "0"}
                        
                        if(![string]::IsNullOrEmpty($global:nettype))
                        {
                            if($global:nettype -eq 1)
                            {
                                $global:adpname = $global:conn.GetNetwork().GetName() | ?{$_ -match $global:adaptername }
                                $global:netcategory = "Private"
    
                            }
                            if([int]$global:nettype -eq 0)
                            {
                                $global:adpname = $global:conn.GetNetwork().GetName() | ?{$_ -match $global:adaptername }
                                $global:netcategory = "Public"
                            }
                        }
    
                        $net_config = Get-WmiObject win32_networkadapterconfiguration 
                        $index = $net_config| select DNSDomain,Description,InterfaceIndex,Ipaddress | ?{$_.Ipaddress -ne $null -or $_.DNSDomain -match $global:adaptername}
                                    
                        foreach($index1 in $index)
                        {
                            $adr = $index1.ipaddress -notmatch ":" 
                            $ip_trim =[string]$ip

                            if([string]$adr -match $ip_trim.trim())
                            {
                                $index = $index1 | select -ExpandProperty InterfaceIndex
                            }                        
                        
                            $netadapter = Get-WmiObject win32_networkadapter -Filter Netconnectionstatus=2 | ?{$_.InterfaceIndex -eq $index}
                            $global:interface = $netadapter | select -ExpandProperty NetconnectionID
                                        
                            if($(get-wmiobject win32_operatingsystem | select -ExpandProperty Name) -match "7 Professional")
                            {
                                $global:ipv4 = "Internet"  
                                $global:ipv6 = "No internet access"
                            }
                            if($(get-wmiobject win32_operatingsystem | select -ExpandProperty Name) -match "2008 Standard")
                            {
                                $global:ipv4 = "Internet"  
                                $global:ipv6 = "Local"
                            }
                            $EditedAdapterInfoObject = New-Object -TypeName psobject
                            $EditedAdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'Name:' -Value $($global:adpname)
                            $EditedAdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'Interface Alias:' -Value $($global:interface)
                            $EditedAdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'Network Category:' -Value $($global:netcategory)
                            $EditedAdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'IPv4 Connectivity:' -Value $($global:ipv4)
                            $EditedAdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'IPv6 Connectivity:' -Value $($global:ipv6)
                            
                            $global:EditedAdapterInfoObjectArray += $EditedAdapterInfoObject

                            $global:SuccessMessageArray += 'Name:' + $($global:adpname)
                            $global:SuccessMessageArray += 'Interface Alias:' + $($global:Interface)
                            $global:SuccessMessageArray += 'Network Category:' + $($global:NetworkCategory)
                            $global:SuccessMessageArray += 'IPv4 Connectivity:' + $($global:IPv4)
                            $global:SuccessMessageArray += 'IPv6 Connectivity:' + $($global:IPv6)
                                    
                        } 
                                   
                        if($global:netconnection.NetworkCategory -eq $global:network_category)
                        {
                            $global:Code = 1
                            $global:ErrorMessageArray += "Network category for ($global:netconnection.Name) is already set to $global:network_category " 
                        }

                    }
                    break
                }
            } 
  

            if($PSVersionTable.PSVersion.Major -ge 3)
            {
                $global:netconnection = Get-NetConnectionProfile
                if($(Get-NetConnectionProfile | ?{$_.Name -match $global:adaptername}) -ne $null){
                    if( $global:netconnection.NetworkCategory -ne $global:network_category)
                    {
                        Set-NetConnectionProfile -NetworkCategory $global:network_category | ?{$_.Name -match "$global:adpname"}
                        $global:netconnection1 = Get-NetConnectionProfile

                        $EditedAdapterInfoObject = New-Object -TypeName psobject
                        $EditedAdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'Message:' -Value $($global:SuccessMessageArray)
                        $EditedAdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'Name:' -Value $($global:netconnection1.Name)
                        $EditedAdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'Interface Alias:' -Value $($global:netconnection1.InterfaceAlias)
                        $EditedAdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'Network Category:' -Value $([string]$global:netconnection1.NetworkCategory)
                        $EditedAdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'IPv4 Connectivity:' -Value $([string]$global:netconnection1.IPv4Connectivity)
                        $EditedAdapterInfoObject | Add-Member -MemberType NoteProperty -Name 'IPv6 Connectivity:' -Value $([string]$global:netconnection1.IPv6Connectivity)

                        $global:EditedAdapterInfoObjectArray += $EditedAdapterInfoObject
                        
                        $global:SuccessMessageArray += "Network type changed"
                        $global:SuccessMessageArray += 'Name:' + $($global:netconnection1.Name)
                        $global:SuccessMessageArray += 'Interface Alias:' + $($global:netconnection1.InterfaceAlias)
                        $global:SuccessMessageArray += 'Network Category:' + $($global:netconnection1.NetworkCategory)
                        $global:SuccessMessageArray += 'IPv4 Connectivity:' + $($global:netconnection1.IPv4Connectivity)
                        $global:SuccessMessageArray += 'IPv6 Connectivity:' + $($global:netconnection1.IPv6Connectivity)
                    }
                    if( $global:netconnection.NetworkCategory -eq $global:network_category)
                    {
                        $global:code = 1
                        $global:ErrorMessageArray = "The Network category for $($global:netconnection1.Name) is already set to $global:network_category"
                    }
                }
                if($(Get-NetConnectionProfile | ?{$_.Name -match $global:adaptername}) -eq $null){
                    $global:code = 2
                    $global:ErrorMessageArray += "Invalid adapter name $global:adaptername. The network adapter found in this machine are $global:sol"
                }
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

        if($OperationType -eq "Read")
        {
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'Adapter' -Value $global:AdapterInfoObjectArray
        }

        if($OperationType -eq "Edit")
        {
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'Modfied Network Adapter' -Value $global:EditedAdapterInfoObjectArray #$EditedAdapterInfoObject | select -Unique
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
    if($OperationType -eq "")
    {
        $global:Code = 2
        $global:ErrorMessageArray += "Type is requied" 
    }
    else
    {
        if($OperationType -eq "Read")
        {
            Read-Network
        }

        if($OperationType -eq "Edit")
        {
            Edit-Network -adaptername $AdapterName -network_category $NetworkCategory
        }
    }
          
    SetResult
    DisplayResult 
}
else
{
    SetResult
    DisplayResult
}
