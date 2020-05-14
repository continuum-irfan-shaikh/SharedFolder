$Global:Stderr = @()
$Global:Changes_Result_stdout = @()
$Global:Format_Check = @()
$Global:Changed_Data  = @()
$Global:Changes_Result = @()
$Global:HashTable_Data = @()
$Global:Invalid1 = @{}
$GLobal:GatewayOffline = $null

#40.100.137.82
#Active IP in Public
#40.100.141.162

#Active IPV6
#fe80::555:1b2e:8b61:eb13

# Correct Format
#2001:4860:4860::8888

$ErrorActionPreference = 'silentlycontinue'

$Adapter_Name = $null
$Address_Type_IPv4_IPv6 = $null
$IPv4_OR_IPv6_Address = $NULL
$IPv4_OR_IPv6_Gateway = $NULL
$IPv4_Subnet_Mask_OR_IPv6_Subnet_Prefix_Length = $NULL

#$Adapter_Name = 'Wireless Network Connection'
#$Address_Type_IPv4_IPv6 = 'IPV4'

#$IPv4_OR_IPv6_Address = '40.100.137.82'
#$IPv4_OR_IPv6_Gateway = '40.100.141.162'

#$IPv4_Subnet_Mask_OR_IPv6_Subnet_Prefix_Length = '255.0.0.0'


$IPv4_Subnet_Mask = $IPv4_Subnet_Mask_OR_IPv6_Subnet_Prefix_Length
$IPv6_Subnet_Prefix_Length = $IPv4_Subnet_Mask_OR_IPv6_Subnet_Prefix_Length



##################### JSON Convert Start
function JSONOutPutCall{
$Global:Invalid1 = @{}
###### if no return

if($GLobal:AdapterValidation -eq 'YES')
{
#write-host 'Valid Adapter'
if($GLobal:GatewayOffline -eq 'YES')
{

$Stderr = Get_stderr -Title 'Provide Gateway is offline' -details "Provided Gateway Address [$IPv4_OR_IPv6_Gateway]"

$HashData = HashTable_Array_JSON -taskName 'Set static IP information' -status 'Failed' -Code '1' -stdout "Provide Gateway is offline hence can't process further" -stderr $Stderr -result 'Error: Provide Gateway is offline'
$JSON2 = $HashData

}else{

if($Global:HashTable_Data.count -eq 0)
{

###### if Format is correct
    if($Global:Format_Check | % {$_ -notlike '*invalid*'})
    {
        ###### if all command success executed
        $DataSuccess = $Global:Changes_Result | % {$_.result} | ? {$_}
        
        if(($DataSuccess -clike 'Successfully Changed').count -gt 0 -or  (($DataSuccess -clike 'Successfully Changed') -eq $true) -or ($DataSuccess -eq $Null))
        {##### Final Adapter Settings
        #$Global:Changed_Data
        $Object = $Global:Changed_Data.count
        $HashData = HashTable_Array_JSON -HeadingUnder_dataObject_Array "Post Changes Adapter configuration" -Array $Global:Changed_Data -taskName 'Set static IP information' -status 'Success' -Code '0' -stdout $Global:Changes_Result_stdout -objects "$Object" -result 'Successfully changed the Adapter configuration'
        #$JSON2 = ConvertTo-Json20 $HashData
        $JSON2 = $HashData
        }

###### if any of the command not successfull
        else{
         ##### Command UnSuccess result
        #$Global:Changes_Result
        
        $Title = $Global:Changes_Result | % {$_.Action}
        $details = $Global:Changes_Result | % {$_.result}
        
        $stderr = Get_stderr -Title $Title -details $details
        
        $HashData = HashTable_Array_JSON -HeadingUnder_dataObject_Array "Error while changing Adapter setting" -taskName 'Set static IP information' -status 'Failed' -Code '1' -stdout 'Failed to change Adapter configuration' -stderr $stderr -result 'Error: Failed to change Adapter configuration'
        #$JSON2 = ConvertTo-Json20 $HashData
        $JSON2 = $HashData

       }
       
    }
###### if any of the given Format is not correct
    if($Global:Format_Check | % {$_ -like '*InValid*'}){
        $Global:Invalid1 = $Global:Format_Check | Get-Member| % {if($_.Definition.split('=')[1] -like '*invalid*'){

        $Keyx = $_.name;$Valuex = ($_.Definition.split('=')[1]).split('|')[0]
        [string]$Name12 = $_.name

        if($Keyx -Like 'IPV4Address' -or $Keyx -Like 'IPV4Gateway' -or $Keyx -Like 'IPV4Subnet')
        {
        
        if($Keyx -notLike 'IPV4Subnet') 
        {
            $Key_Values = @{$Keyx = @(
            "Submitted $Name12 : [$Valuex]"
            "Expected format: [1.102.103.104]"
            )}
            $Global:stderr += "Submitted $($Name12): [$Valuex] and Expected format: [1.102.103.104]"        
        }else
        {
            $Key_Values = @{$Keyx = @(
            "Submitted $Name12 : [$Valuex]"
            "Expected format: [1.102.103.104]"
            )}
            $Global:stderr += "Submitted $($Name12): [$Valuex] and Expected format: [255.255.255.0 or 255.0.0.0]"        
         }

        
        }


        if($Keyx -Like 'IPV6Address' -or $Keyx -Like 'IPV6Gateway')
        {
        [string]$Name12 = $_.name
        
        $Key_Values = @{$Keyx = @(
        "Submitted IPv6 address: [$Valuex]",
        "Expected format: [2001:4860:4860::8888] or [2001:4860:4860:0000:0000:0000:0000:8888]"
        )}
        $Global:stderr += "Submitted $($Name12): [$Valuex] and Expected format: [2001:4860:4860::8888] or [2001:4860:4860:0000:0000:0000:0000:8888]"
        } 
        if($Keyx -Like 'IPv6_Subnet_Prefix_Length1')
        {

        $Key_Values = @{$Keyx = @(
        "Submitted IPv6 Subnet Prefix Length: [$Valuex]",
        "Expected format: [8-128]"
        )}
        
        $Global:stderr += "Submitted IPv6 Subnet Prefix Length: [$Valuex] and Expected format: [8-128]"
        }
        
        $Key_Values
         
        }}
        
        $stderr = Get_stderr -Title 'invalid input provided' -details $Global:stderr
        
        $AddressType = $Global:Format_Check.'Address Type (IPv4/IPv6)'
              
        $HashData = HashTable_Array_JSON -HeadingUnder_dataObject_Array "invalid Input passed for $AddressType" -taskName 'Set static IP information' -status 'Failed' -Code '1' -stdout 'invalid Input provided' -result 'Error: invalid Input provided' -stderr $stderr
        #$JSON2 = ConvertTo-Json20 $HashData
        $JSON2 = $HashData
    }

}

if($Global:HashTable_Data.count -ne 0){
#$Global:HashTable_Data

$HashData = HashTable_Array_JSON -taskName 'Set static IP information' -status 'Failed' -Code '1' -stdout 'Provide valid inputs' -stderr $Global:Stderr -result 'Error: Incorrect inputs'
$JSON2 = $HashData
}

}
}else{

#write-host 'Invalid Adapter'

$Global:Stderr = New-Object psobject -Property @{
'Adapter Name' = $Adapter_Name
'Existance Status' = 'Not exists'
}

$HashData = HashTable_Array_JSON -taskName 'Set static IP information' -status 'Failed' -Code '1' -stdout 'provided Invalid Adapter name' -stderr $Global:Stderr -result 'Error: Invalid Adapter name'
$JSON2 = $HashData

}
#OUTPUT
##################
$JSON2

}
####################### JSON Convert END

########################## GET_Stderr
function Get_stderr
{
param($Title,$details)
$info = @()
$ID = 0

 $details | %{
    $info += new-object psobject -Property @{
    "id" = $ID
    "title" = $Title
    "detail" = $_
    }
    $ID++
    }

return $info
}

#################################
if($IPv4_OR_IPv6_Address -eq $IPv4_OR_IPv6_Gateway){

$Stderr = Get_stderr -Title 'IP and Gateway address should not be same' -details "Provided Address are same IP Address [$IPv4_OR_IPv6_Address] and Gateway [$IPv4_OR_IPv6_Gateway]"

$HashData = HashTable_Array_JSON -taskName 'Set static IP information' -status 'Failed' -Code '1' -stdout "Provided Address are same IP Address [$IPv4_OR_IPv6_Address] and Gateway [$IPv4_OR_IPv6_Gateway]" -stderr $Stderr -result 'Error: Provided Same IP and Gateway address'
$JSON2 = $HashData
$JSON2 
return
}

####################################################################  Display Adapter Settings 

function DisplayAdapterConfiguration{
param($Which_Version_Setting_IPV6_OR_IPV4,$Adapter_Name)
    $Adapter_Properties = (Gwmi win32_networkadapter | ? {$_.NetConnectionID -eq "$Adapter_Name"})
    $adapterSettings = gwmi win32_networkadapterconfiguration | ? {$_.index -eq $Adapter_Properties.index}
    
  if($Which_Version_Setting_IPV6_OR_IPV4 -eq 'IPV6'){

   #################### IPV6 Current Settings
    $IPAddress_Index = ($adapterSettings.ipaddress.count -1)
    $IPSubnet_Index =($adapterSettings.IPSubnet.count -1)

    $Defaultx1 = ($adapterSettings.DefaultIPGateway | ? {$_ -match ':'}) -join ','
    
    
    $IPV6_Settings = New-Object psobject -Property @{
    'Adapter Name' = $Adapter_Properties.NetConnectionID
    'Adapter Enabled/Disabled Status' = 'Enabled'
    'IPV6 Address' = $adapterSettings.ipaddress[$IPAddress_Index]
    'IPV6 Gateway' = $Defaultx1
    'IPV6 Subnet Prefix Length' = $adapterSettings.IPSubnet[$IPSubnet_Index]
    } | select 'Adapter Name','Adapter Enabled/Disabled Status','IPV6 Address','IPV6 Subnet Prefix Length','IPV6 Gateway'
    $IPV6_Settings
    
    $NetConnectionID = $Adapter_Properties.NetConnectionID
    $Subnet_Prefix = $IPV6_Settings.'IPV6 Subnet Prefix Length'
    $IPAddress = $IPV6_Settings.'IPV6 Address'
    
    $Global:Changes_Result_stdout = "Adapter Name: $NetConnectionID`r`n Adapter Enabled/Disabled Status: Enabled`r`n IPV6 Address: $IPAddress`r`n IPV6 Subnet Prefix Length: $Subnet_Prefix`r `n IPV6 Gateway:$DefaultGateway"    
    }
    
  if($Which_Version_Setting_IPV6_OR_IPV4 -eq 'IPV4'){
  
    if($adapterSettings.DHCPEnabled -eq $False)
    {    
        $IPAddress  = ($adapterSettings.IpAddress[0]  | ? {$_ -notmatch ':'})
        $SubnetMask  = $adapterSettings.IPSubnet[0]
        $DefaultGateway = ($adapterSettings.DefaultIPGateway[0] | ? {$_ -notmatch ':'})
   }else{
        $IPAddress  = $null
        $SubnetMask  = $null
        $DefaultGateway = $null   
   }
    $OutputObj = New-Object psobject -Property @{
    'Adapter Name' = $Adapter_Properties.NetConnectionID
    'Adapter Enabled/Disabled Status' = 'Enabled'
    'IPV4 Address' = $IPAddress
    'IPV4 SubnetMask' = $SubnetMask
    'IPV4 Gateway' = $DefaultGateway
     } | select 'Adapter Name','Adapter Enabled/Disabled Status','IPV4 Address','IPV4 SubnetMask','IPV4 Gateway'
     $OutputObj
     
     $NetConnectionID = $Adapter_Properties.NetConnectionID
     
     $Global:Changes_Result_stdout = "Adapter Name: $NetConnectionID`r`n Adapter Enabled/Disabled Status: Enabled`r`n IPV4 Address: $IPAddress`r`n IPV4 SubnetMask: $SubnetMask`r `n IPV4 Gateway:$DefaultGateway"
     
    }
 }
 
######################################################################## Display Adapter Settings END


function Check_ReturnValue
{ param($Return_Value)
  
    switch ($Return_Value) 
    {
        -1 {'Successfully Changed'}
        0  {'Successfully Changed'; break}
        1  {'Successfully Changed'; break}
        64 {'Method not supported on this platform'; break}
        65 {'Unknown failure'; break}
        66 {'Invalid subnet mask'; break}
        67 {'An error occurred while processing an Instance that was returned'; break}
        68 {'Invalid input parameter'; break}
        69 {'More than 5 gateways specified'; break}
        70 {'Invalid IP address'; break}
        71 {'Invalid gateway IP address'; break}
        72 {'An error occurred while accessing the Registry for the requested information'; break}
        73 {'Invalid domain name'; break}
        74 {'Invalid host name'; break}
        75 {'No primary/secondary WINS server defined'; break}
        76 {'Invalid file'; break}
        77 {'Invalid system path'; break}
        78 {'File copy failed'; break}
        79 {'Invalid security parameter'; break}
        80 {'Unable to configure TCP/IP service'; break}
        81 {'Unable to configure DHCP service'; break}
        82 {'Unable to renew DHCP lease'; break}
        83 {'Unable to release DHCP lease'; break}
        84 {'IP not enabled on adapter'; break}
        85 {'IPX not enabled on adapter'; break}
        86 {'Frame/network number bounds error'; break}
        87 {'Invalid frame type'; break}
        88 {'Invalid network number'; break}
        89 {'Duplicate network number'; break}
        90 {'Parameter out of bounds'; break}
        91 {'Access denied'; break}
        92 {'Out of memory'; break}
        93 {'Already exists'; break}
        94 {'Path, file or object not found'; break}
        95 {'Unable to notify service'; break}
        96 {'Unable to notify DNS service'; break}
        97 {'Interface not configurable'; break}
        98 {'Not all DHCP leases could be released/renewed'; break}
        100 {'DHCP not enabled on adapter'; break}
        2147786788 {"Write lock not enabled"; break}
        2147749891 {"Must be run with admin privileges"; break}
        default {"Faild with error code $($Return_Value)"; break}
    }

}


################################################ Changes/Action Part

Function Validate_AdapterName{
Param($Adapter_Name)
$adapter = Gwmi win32_networkadapter | where {$_.NetConnectionID -eq "$Adapter_Name"}
if($adapter){
  return $True
}
    else{
      return $False
    }

}

Function Change_Adapter_Settings{
param($Adapter_Name,$Which_Version_Setting_IPV6_OR_IPV4,$IPv4_OR_IPv6_Address,$IPv4_OR_IPv6_Gateway,$IPv4_Subnet_Mask,$IPv6_Subnet_Prefix_Length)


$ErrorActionPreference = 'Silentlycontinue'

$adapter_Index = (Gwmi win32_networkadapter | ? {$_.NetConnectionID -eq "$Adapter_Name"}).index
######## Get Adapter Settings    
$adapterSettings = Get-WmiObject win32_networkAdapterConfiguration | where {$_.index -eq $adapter_Index} | select *


if($Which_Version_Setting_IPV6_OR_IPV4 -eq 'IPV6')
{

   #################### IPV6 Current Settings
    $IPV4_IP_Index = ($adapterSettings.ipaddress.count -1)
    $IPSubnet_Index =($adapterSettings.IPSubnet.count -1)

    $Defaultx1 = ($adapterSettings.DefaultIPGateway | ? {$_ -match ':'}) -join ','

    
    $IPV6_Settings = New-Object psobject -Property @{
    'IPV6 IP-Address' = $adapterSettings.ipaddress[$IPV4_IP_Index]
    'IPV6 DefaultIPGateway' = $Defaultx1
    'IPV6 Subnet Prefix Length' = $adapterSettings.IPSubnet[$IPSubnet_Index]
    } | select 'IPV6 IP-Address','IPV6 DefaultIPGateway','IPV6 Subnet Prefix Length'
    #$IPV6_Settings
    
    $IPV6_IP = $IPV6_Settings.'IPV6 IP-Address'
    $IPV6_Gateway = $IPV6_Settings.'IPV6 DefaultIPGateway'
    $IPV6_Subnet_Prefix = $IPV6_Settings.'IPV6 Subnet Prefix Length'
    
# IP Found then Remove IP Values
    if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -or ![string]::IsNullOrEmpty($IPv6_Subnet_Prefix_Length))
    {
      $IPx = $IPV6_IP | ? {$_ -match ':'}
      if(![string]::IsNullOrEmpty($IPx))
      {
      #write-host 'Removing the existing IPV6 Address'
      netsh interface ipv6 delete address $Adapter_Name $IPV6_IP
      #write-host ip $IPV6_IP
      
      }
     }
    
# Gateway Found then Remove gateway Values
    if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway))
    { 
      if(![string]::IsNullOrEmpty($Defaultx1))
        {  
           #write-host 'Removing the existing IPV6 Gatway'
           netsh interface ipv6 delete route ::/0 "$Adapter_Name" "$IPV6_Gateway"
         }
     }
      

# IPAddress and Suffix given
    if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and ![string]::IsNullOrEmpty($IPv6_Subnet_Prefix_Length))
    {
      if(![string]::IsNullOrEmpty($IPV6_Settings.'IPV6 IP-Address')){
      #write-host "Adding the new IPV6 Address along with the SubnetPrefix length $IPv4_OR_IPv6_Address/$IPv6_Subnet_Prefix_Length"
      
      $Action = 'Initiate changes for -IPV6 Address and -SubnetPrefix length'
      $Result = netsh interface ipv6 add address $Adapter_Name "$IPv4_OR_IPv6_Address/$IPv6_Subnet_Prefix_Length"
       $Result = $Result | ? {$_}
      if(($Result -eq $null) -or ($Result -eq '*OK*')){$Result = 'Successfully Changed'}
       
      $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}    
      $Global:Changes_Result_stdout += "Action: $Action`r`nresult = $result"
# Gateway not given
      if([string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway))
        {  
             if(![string]::IsNullOrEmpty($Defaultx1))
             {  
               #write-host 'Adding the existing Gatway'
               netsh interface ipv6 add route ::/0 "$Adapter_Name" "$IPV6_Gateway"
             }
         }

# Gateway given
      if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway))
        {  
               #write-host 'Adding the new Gatway'      
               $Action = 'Initiate changes for -IPV6 Gateway'
               $Result = netsh interface ipv6 add route ::/0 "$Adapter_Name" "$IPv4_OR_IPv6_Gateway"
            $Result = $Result | ? {$_}
      if(($Result -eq $null) -or ($Result -eq '*OK*')){$Result = 'Successfully Changed'}
               $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}                 
               $Global:Changes_Result_stdout += "Action: $Action`r`nresult = $result"
         }         
      }

    }else{

# Only IPAddress and Not Suffix given
    if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and [string]::IsNullOrEmpty($IPv6_Subnet_Prefix_Length))
    {
      #write-host 'Adding the new IPV6 Address along with the default 64 SubnetPrefix length'    
      
      $Action = 'Initiate changes for -IPV6 Address with Default SubnetPrefix length'
      $result = netsh interface ipv6 add address $Adapter_Name "$IPv4_OR_IPv6_Address"
      $Result = $Result | ? {$_}
      if(($Result -eq $null) -or ($Result -eq '*OK*')){$Result = 'Successfully Changed'}
      $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}
      $Global:Changes_Result_stdout += "Action: $Action`r`nresult = $result"
      }

# Only Suffix and Not IP given
    if([string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and ![string]::IsNullOrEmpty($IPv6_Subnet_Prefix_Length))
    {
      #write-host 'Appending only SubnetPrefix length'
      $Action = 'Adding SubnetPrefix length'      
      $Result = netsh interface ipv6 add address $Adapter_Name "$IPV6_IP/$IPv6_Subnet_Prefix_Length"
            $Result = $Result | ? {$_}
      if(($Result -eq $null) -or ($Result -eq '*OK*')){$Result = 'Successfully Changed'}

      $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}                     
      $Global:Changes_Result_stdout += "Action: $Action`r`nresult = $result"
     }
    
    if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway))
    {
     #write-host "Changing the IPV6 Gatway 117 $IPv4_OR_IPv6_Gateway"
     $Action = 'Initiate changes for -IPV6 Gateway'
     $Result = netsh interface ipv6 add route ::/0 "$Adapter_Name" "$IPv4_OR_IPv6_Gateway"  
            $Result = $Result | ? {$_}
           if(($Result -eq $null) -or ($Result -eq '*OK*')){$Result = 'Successfully Changed'}

      $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}                 
          $Global:Changes_Result_stdout += "Action: $Action`r`nresult: $result"
    }  
}
##############################################################################################################################

}

if($Which_Version_Setting_IPV6_OR_IPV4 -eq 'IPV4')
{    
    $Adapter_Properties = (Gwmi win32_networkadapter | ? {$_.NetConnectionID -eq "$Adapter_Name"})
    $ethernet = gwmi win32_networkadapterconfiguration | ? {$_.index -eq $Adapter_Properties.index}
    $NetConnectionID = $Adapter_Properties.NetConnectionID
    $InterfaceName = $Adapter_Properties.Name
    
    if($ethernet.DHCPEnabled -eq $False)
    {
        $IPV4_IP  = $adapterSettings.IpAddress[0]
        $SubnetMask  = $adapterSettings.IPSubnet[0]
        $DefaultGateway = $adapterSettings.DefaultIPGateway[0]

    $OutputObj  = New-Object -Type PSObject
    $OutputObj | Add-Member -MemberType NoteProperty -Name IPAddress -Value $IPV4_IP
    $OutputObj | Add-Member -MemberType NoteProperty -Name SubnetMask -Value $SubnetMask
    $OutputObj | Add-Member -MemberType NoteProperty -Name Gateway -Value ($DefaultGateway -join ",")
    $IPV4_Settings = $OutputObj
    
    $IPV4_IP = $IPV4_Settings.'IPAddress'
    $IPV4_Gateway = $IPV4_Settings.'Gateway'
    $IPV4_Subnet = $IPV4_Settings.'SubnetMask'
            
      }else{
      
     if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and [string]::IsNullOrEmpty($IPv4_Subnet_Mask) -and ![string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway))
      {
       #write-host "Please Provide IPV4 address and its Subnetmask combination to proceed"
       $Global:HashTable_Data += 'Please Provide IP address and Subnetmask combination to proceed'
       $Global:Stderr += Get_stderr -Title 'Unable to process request' -details $Global:HashTable_Data
       return        
      }

      
        $IPV4_IP  = $null
        $SubnetMask  = $null
        $DefaultGateway = $null
      }
      
# IPAddress and subnet given
    if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and ![string]::IsNullOrEmpty($IPv4_Subnet_Mask) -and ![string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway))
    {

      #write-host 'Adding the new IPV4 Address +SubnetMask+Gateway'
      #netsh interface ipv6 add address $Adapter_Name "$IPv4_OR_IPv6_Address/$IPv4_Subnet_Mask"
        $ReturnValue = ($ethernet | % {
            $_.EnableStatic($IPv4_OR_IPv6_Address, $IPv4_Subnet_Mask)}).ReturnValue

        $Changesinfo = Check_ReturnValue -Return_Value $ReturnValue
        #write-host  "inf: $Changesinfo - $ReturnValue"
        
        $Action = 'initiate to setting up IPV4 Address | SubnetMask | Gateway'
        $result = $Changesinfo
        
      $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}
      $Global:Changes_Result_stdout += "Action: $Action`r`nresult = $result"       
     sleep 5
     #netsh interface ipv6 add route ::/0 "$Adapter_Name" "$IPv4_OR_IPv6_Gateway"      

        if($result -eq 'Successfully Changed'){
         #write-host 'Changing the IPV4 Gatway'        
        $ReturnValue = ($ethernet | % {
            $_.SetGateways($IPv4_OR_IPv6_Gateway)}).ReturnValue
     
        $Changesinfo = Check_ReturnValue -Return_Value $ReturnValue
        #write-host  "inf: $Changesinfo - $ReturnValue"

        $Action = 'initiate to setting Gateway'
        $result = $Changesinfo
        
      $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}
      $Global:Changes_Result_stdout += "Action: $Action`r`nresult = $result"
      } 
     }
    else{

#write-host 'ElsePart IPV4'
sleep 5
# Only IPAddress and Not SubnetMask given
    if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and [string]::IsNullOrEmpty($IPv4_Subnet_Mask))
    {
      if(![string]::IsNullOrEmpty($IPV4_Subnet)){
      #write-host "Adding the new IPV4 Address along with the Previous SubnetMask $IPv4_OR_IPv6_Address | $IPV4_Subnet"
      
           $ReturnValue = ($ethernet | % {
           $_.EnableStatic($IPv4_OR_IPv6_Address, $IPV4_Subnet)}).ReturnValue

      $Changesinfo = Check_ReturnValue -Return_Value $ReturnValue
      #write-host  "inf: $Changesinfo - $ReturnValue"

      $Action = 'initiate to setting IP Address with the Previous SubnetMask'
      $result = $Changesinfo
        
      $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}
                       $Global:Changes_Result_stdout += "Action: $Action`r`nresult = $result"
        }else{
        #write-host "Please enter SubnetMask also change SET the IP Address (($IPv4_OR_IPv6_Address))"
        $Global:HashTable_Data = 'Please enter SubnetMask also change SET the IP Address (($IPv4_OR_IPv6_Address))'
        $Global:Stderr += Get_stderr -Title 'Unable to process request' -details $Global:HashTable_Data        
        return
        }
        

        
        if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway))
        {
         #write-host 'Changing the IPV6 Gatway'    
         #netsh interface ipv6 add route ::/0 "$Adapter_Name" "$IPv4_OR_IPv6_Gateway"      

        $ReturnValue = ($ethernet | % {
            $_.SetGateways($IPv4_OR_IPv6_Gateway)}).ReturnValue
        $Changesinfo = Check_ReturnValue -Return_Value $ReturnValue
        #write-host  "inf: $Changesinfo - $ReturnValue"        

        $Action = 'initiate to setting Gateway'
        $result = $Changesinfo
        
        $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}        
              $Global:Changes_Result_stdout += "Action: $Action`r`nresult = $result"
        }
        
      }

  

# Only Suffix and Not IP given
    if(!$ethernet.DHCPEnabled -and [string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and ![string]::IsNullOrEmpty($IPv4_Subnet_Mask))
    {
      #write-host 'Changing SubnetMask with Previous IP'
      #netsh interface ipv6 add address $Adapter_Name "$IPv4_Subnet_Mask"
      
        $ReturnValue = ($ethernet | % {
            $_.EnableStatic($IPV4_IP, $IPv4_Subnet_Mask)}).ReturnValue
        
        $Changesinfo = Check_ReturnValue -Return_Value $ReturnValue
        #write-host  "inf: $Changesinfo - $ReturnValue"

        $Action = 'initiate to setting IP Address and Gateway'
        $result = $Changesinfo
        
        $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}    
        
        if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway))
        {
         #write-host 'Changing the IPV6 Gatway'
         #netsh interface ipv6 add route ::/0 "$Adapter_Name" "$IPv4_OR_IPv6_Gateway"      

        $ReturnValue = ($ethernet | % {
            $_.SetGateways($IPv4_OR_IPv6_Gateway)}).ReturnValue
        $Changesinfo = Check_ReturnValue -Return_Value $ReturnValue
        #write-host  "inf: $Changesinfo - $ReturnValue"        

        $Action = 'initiate to setting IP Gateway'
        $result = $Changesinfo
        
        $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}   
                      $Global:Changes_Result_stdout += "Action: $Action`r`nresult = $result"
        }
     }
 
 
    if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway) -and [string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and [string]::IsNullOrEmpty($IPv4_Subnet_Mask))
    {
         #write-host "Changing the IPV6 Gatway $IPv4_OR_IPv6_Gateway"
         #netsh interface ipv6 add route ::/0 "$Adapter_Name" "$IPv4_OR_IPv6_Gateway"      

        $ReturnValue = ($ethernet | % {
            $_.SetGateways($IPv4_OR_IPv6_Gateway)}).ReturnValue
        $Changesinfo = Check_ReturnValue -Return_Value $ReturnValue
        #write-host  "inf: $Changesinfo - $ReturnValue"        

        $Action = 'initiate to setting IP Gateway'
        $result = $Changesinfo
        
        $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}  
      $Global:Changes_Result_stdout += "Action: $Action`r`nresult = $result"
     }    

    if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and [string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway) -and ![string]::IsNullOrEmpty($IPv4_Subnet_Mask))
    {
      #write-host 'Adding the new IPV4 Address+SubnetMask'
      #netsh interface ipv6 add address $Adapter_Name "$IPv4_OR_IPv6_Address/$IPv4_Subnet_Mask"
        $ReturnValue = ($ethernet | % {
            $_.EnableStatic($IPv4_OR_IPv6_Address, $IPv4_Subnet_Mask)}).ReturnValue

        $Changesinfo = Check_ReturnValue -Return_Value $ReturnValue
        #write-host  "inf: $Changesinfo - $ReturnValue"
        
        $Action = 'initiate to setting up IPV4 Address | SubnetMask'
        $result = $Changesinfo
        
      $Global:Changes_Result += New-Object psobject -Property @{ Action = $Action;result = $result}
      $Global:Changes_Result_stdout += "Action: $Action`r`nresult = $result"
     }

##############################################################################################################################
}



}
sleep 4
$Global:Changed_Data = @(DisplayAdapterConfiguration -Which_Version_Setting_IPV6_OR_IPV4 "$Which_Version_Setting_IPV6_OR_IPV4" -Adapter_Name "$Adapter_Name")
}
################################################ Changes/Action Part END




################################# JSON Convert

#####################################
#################Convert_HashToString
#####################################
function Convert_HashToString
{param([System.Collections.Hashtable]$Hash)
$hashstr = ""; $keys = $Hash.keys; foreach ($key in $keys) { $v = $Hash[$key]; 
if ($key -match "\s") { $hashstr += "`"$key`"" + ":" + "`"$v`"" + ";" }
else { $hashstr += $key + ":" + "`"$v`"" + ";" } }; $hashstr += "";
return $hashstr}


function FormatString {
    param(
        [String] $String)
    # removed: #-replace '/', '\/' `
    # This is returned 
    $String -replace '\\', '\\' -replace '\n', '\n' `
        -replace '\u0008', '\b' -replace '\u000C', '\f' -replace '\r', '\r' `
        -replace '\t', '\t' -replace '"', '\"'
}

function GetNumberOrString {
    param(
        $InputObject)
    if ($InputObject -is [System.Byte] -or $InputObject -is [System.Int32] -or `
        ($env:PROCESSOR_ARCHITECTURE -imatch '^(?:amd64|ia64)$' -and $InputObject -is [System.Int64]) -or `
        $InputObject -is [System.Decimal] -or $InputObject -is [System.Double] -or `
        $InputObject -is [System.Single] -or $InputObject -is [long] -or `
        ($Script:CoerceNumberStrings -and $InputObject -match $Script:NumberRegex)) {
        Write-Verbose -Message "Got a number as end value."
        "$InputObject"
    }
    else {
        Write-Verbose -Message "Got a string as end value."
        """$(FormatString -String $InputObject)"""
    }
}

function ConvertToJsonInternal {
    param(
        $InputObject, # no type for a reason
        [Int32] $WhiteSpacePad = 0)
    [String] $Json = ""
    $Keys = @()
    Write-Verbose -Message "WhiteSpacePad: $WhiteSpacePad."
    if ($null -eq $InputObject) {
        Write-Verbose -Message "Got 'null' in `$InputObject in inner function"
        $null
    }
    elseif ($InputObject -is [Bool] -and $InputObject -eq $true) {
        Write-Verbose -Message "Got 'true' in `$InputObject in inner function"
        $true
    }
    elseif ($InputObject -is [Bool] -and $InputObject -eq $false) {
        Write-Verbose -Message "Got 'false' in `$InputObject in inner function"
        $false
    }
    elseif ($InputObject -is [HashTable]) {
        $Keys = @($InputObject.Keys)
        Write-Verbose -Message "Input object is a hash table (keys: $($Keys -join ', '))."
    }
    elseif ($InputObject.GetType().FullName -eq "System.Management.Automation.PSCustomObject") {
        $Keys = @(Get-Member -InputObject $InputObject -MemberType NoteProperty |
            Select-Object -ExpandProperty Name)
        Write-Verbose -Message "Input object is a custom PowerShell object (properties: $($Keys -join ', '))."
    }
    elseif ($InputObject.GetType().Name -match '\[\]|Array') {
        Write-Verbose -Message "Input object appears to be of a collection/array type."
        Write-Verbose -Message "Building JSON for array input object."
        #$Json += " " * ((4 * ($WhiteSpacePad / 4)) + 4) + "[`n" + (($InputObject | ForEach-Object {
        $Json += "[`n" + (($InputObject | ForEach-Object {
            if ($null -eq $_) {
                Write-Verbose -Message "Got null inside array."
                " " * ((4 * ($WhiteSpacePad / 4)) + 4) + "null"
            }
            elseif ($_ -is [Bool] -and $_ -eq $true) {
                Write-Verbose -Message "Got 'true' inside array."
                " " * ((4 * ($WhiteSpacePad / 4)) + 4) + "true"
            }
            elseif ($_ -is [Bool] -and $_ -eq $false) {
                Write-Verbose -Message "Got 'false' inside array."
                " " * ((4 * ($WhiteSpacePad / 4)) + 4) + "false"
            }
            elseif ($_ -is [HashTable] -or $_.GetType().FullName -eq "System.Management.Automation.PSCustomObject" -or $_.GetType().Name -match '\[\]|Array') {
                Write-Verbose -Message "Found array, hash table or custom PowerShell object inside array."
                " " * ((4 * ($WhiteSpacePad / 4)) + 4) + (ConvertToJsonInternal -InputObject $_ -WhiteSpacePad ($WhiteSpacePad + 4)) -replace '\s*,\s*$' #-replace '\ {4}]', ']'
            }
            else {
                Write-Verbose -Message "Got a number or string inside array."
                $TempJsonString = GetNumberOrString -InputObject $_
                " " * ((4 * ($WhiteSpacePad / 4)) + 4) + $TempJsonString
            }
        #}) -join ",`n") + "`n],`n"
        }) -join ",`n") + "`n$(" " * (4 * ($WhiteSpacePad / 4)))],`n"
    }
    else {
        Write-Verbose -Message "Input object is a single element (treated as string/number)."
        GetNumberOrString -InputObject $InputObject
    }
    if ($Keys.Count) {
        Write-Verbose -Message "Building JSON for hash table or custom PowerShell object."
        $Json += "{`n"
        foreach ($Key in $Keys) {
            # -is [PSCustomObject]) { # this was buggy with calculated properties, the value was thought to be PSCustomObject
            if ($null -eq $InputObject.$Key) {
                Write-Verbose -Message "Got null as `$InputObject.`$Key in inner hash or PS object."
                $Json += " " * ((4 * ($WhiteSpacePad / 4)) + 4) + """$Key"": null,`n"
            }
            elseif ($InputObject.$Key -is [Bool] -and $InputObject.$Key -eq $true) {
                Write-Verbose -Message "Got 'true' in `$InputObject.`$Key in inner hash or PS object."
                $Json += " " * ((4 * ($WhiteSpacePad / 4)) + 4) + """$Key"": true,`n"            }
            elseif ($InputObject.$Key -is [Bool] -and $InputObject.$Key -eq $false) {
                Write-Verbose -Message "Got 'false' in `$InputObject.`$Key in inner hash or PS object."
                $Json += " " * ((4 * ($WhiteSpacePad / 4)) + 4) + """$Key"": false,`n"
            }
            elseif ($InputObject.$Key -is [HashTable] -or $InputObject.$Key.GetType().FullName -eq "System.Management.Automation.PSCustomObject") {
                Write-Verbose -Message "Input object's value for key '$Key' is a hash table or custom PowerShell object."
                $Json += " " * ($WhiteSpacePad + 4) + """$Key"":`n$(" " * ($WhiteSpacePad + 4))"
                $Json += ConvertToJsonInternal -InputObject $InputObject.$Key -WhiteSpacePad ($WhiteSpacePad + 4)
            }
            elseif ($InputObject.$Key.GetType().Name -match '\[\]|Array') {
                Write-Verbose -Message "Input object's value for key '$Key' has a type that appears to be a collection/array."
                Write-Verbose -Message "Building JSON for ${Key}'s array value."
                $Json += " " * ($WhiteSpacePad + 4) + """$Key"":`n$(" " * ((4 * ($WhiteSpacePad / 4)) + 4))[`n" + (($InputObject.$Key | ForEach-Object {
                    #Write-Verbose "Type inside array inside array/hash/PSObject: $($_.GetType().FullName)"
                    if ($null -eq $_) {
                        Write-Verbose -Message "Got null inside array inside inside array."
                        " " * ((4 * ($WhiteSpacePad / 4)) + 8) + "null"
                    }
                    elseif ($_ -is [Bool] -and $_ -eq $true) {
                        Write-Verbose -Message "Got 'true' inside array inside inside array."
                        " " * ((4 * ($WhiteSpacePad / 4)) + 8) + "true"
                    }
                    elseif ($_ -is [Bool] -and $_ -eq $false) {
                        Write-Verbose -Message "Got 'false' inside array inside inside array."
                        " " * ((4 * ($WhiteSpacePad / 4)) + 8) + "false"
                    }
                    elseif ($_ -is [HashTable] -or $_.GetType().FullName -eq "System.Management.Automation.PSCustomObject" `
                        -or $_.GetType().Name -match '\[\]|Array') {
                        Write-Verbose -Message "Found array, hash table or custom PowerShell object inside inside array."
                        " " * ((4 * ($WhiteSpacePad / 4)) + 8) + (ConvertToJsonInternal -InputObject $_ -WhiteSpacePad ($WhiteSpacePad + 8)) -replace '\s*,\s*$'
                    }
                    else {
                        Write-Verbose -Message "Got a string or number inside inside array."
                        $TempJsonString = GetNumberOrString -InputObject $_
                        " " * ((4 * ($WhiteSpacePad / 4)) + 8) + $TempJsonString
                    }
                }) -join ",`n") + "`n$(" " * (4 * ($WhiteSpacePad / 4) + 4 ))],`n"
            }
            else {
                Write-Verbose -Message "Got a string inside inside hashtable or PSObject."
                # '\\(?!["/bfnrt]|u[0-9a-f]{4})'
                $TempJsonString = GetNumberOrString -InputObject $InputObject.$Key
                $Json += " " * ((4 * ($WhiteSpacePad / 4)) + 4) + """$Key"": $TempJsonString,`n"
            }
        }
        $Json = $Json -replace '\s*,$' # remove trailing comma that'll break syntax
        $Json += "`n" + " " * $WhiteSpacePad + "},`n"
    }
    $Json
}

function ConvertTo-JSON2 {
    [CmdletBinding()]
    #[OutputType([Void], [Bool], [String])]
    param(
        [AllowNull()]
        [Parameter(Mandatory=$true,
                   ValueFromPipeline=$true,
                   ValueFromPipelineByPropertyName=$true)]
        $InputObject,
        [Switch] $Compress,
        [Switch] $CoerceNumberStrings = $false)
    begin{
        $JsonOutput = ""
        $Collection = @()
        # Not optimal, but the easiest now.
        [Bool] $Script:CoerceNumberStrings = $CoerceNumberStrings
        [String] $Script:NumberRegex = '^-?\d+(?:(?:\.\d+)?(?:e[+\-]?\d+)?)?$'
        #$Script:NumberAndValueRegex = '^-?\d+(?:(?:\.\d+)?(?:e[+\-]?\d+)?)?$|^(?:true|false|null)$'
    }
    process {
        # Hacking on pipeline support ...
        if ($_) {
            Write-Verbose -Message "Adding object to `$Collection. Type of object: $($_.GetType().FullName)."
            $Collection += $_
        }
    }
    end {
        if ($Collection.Count) {
            Write-Verbose -Message "Collection count: $($Collection.Count), type of first object: $($Collection[0].GetType().FullName)."
            $JsonOutput = ConvertToJsonInternal -InputObject ($Collection | ForEach-Object { $_ })
        }
        else {
            $JsonOutput = ConvertToJsonInternal -InputObject $InputObject
        }
        if ($null -eq $JsonOutput) {
            Write-Verbose -Message "Returning `$null."
            return $null # becomes an empty string :/
        }
        elseif ($JsonOutput -is [Bool] -and $JsonOutput -eq $true) {
            Write-Verbose -Message "Returning `$true."
            [Bool] $true # doesn't preserve bool type :/ but works for comparisons against $true
        }
        elseif ($JsonOutput-is [Bool] -and $JsonOutput -eq $false) {
            Write-Verbose -Message "Returning `$false."
            [Bool] $false # doesn't preserve bool type :/ but works for comparisons against $false
        }
        elseif ($Compress) {
            Write-Verbose -Message "Compress specified."
            (
                ($JsonOutput -split "\n" | Where-Object { $_ -match '\S' }) -join "`n" `
                    -replace '^\s*|\s*,\s*$' -replace '\ *\]\ *$', ']'
            ) -replace ( # these next lines compress ...
                '(?m)^\s*("(?:\\"|[^"])+"): ((?:"(?:\\"|[^"])+")|(?:null|true|false|(?:' + `
                    $Script:NumberRegex.Trim('^$') + `
                    ')))\s*(?<Comma>,)?\s*$'), "`${1}:`${2}`${Comma}`n" `
              -replace '(?m)^\s*|\s*\z|[\r\n]+'
        }
        else {
            ($JsonOutput -split "\n" | Where-Object { $_ -match '\S' }) -join "`n" `
                -replace '^\s*|\s*,\s*$' -replace '\ *\]\ *$', ']'
        }
    }
}

function ConvertTo-Json20([object] $item){
    add-type -assembly system.web.extensions
    $ps_js=new-object system.web.script.serialization.javascriptSerializer
    return $ps_js.Serialize($item)
}

function ConvertFrom-Json20([object] $item){ 
    add-type -assembly system.web.extensions
    $ps_js=new-object system.web.script.serialization.javascriptSerializer

    #The comma operator is the array construction operator in PowerShell
    return ,$ps_js.DeserializeObject($item)
}


function HashTable_Array_JSON{
Param($HeadingUnder_dataObject_Array,$Array,$taskName,$status,$Code,$stdout,$stderr,$objects,$result)

if($psversiontable.PSVersion.Major -eq 2)
{
    if($stderr -eq $null){
    $axa = ConvertTo-JSON2 -InputObject @{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @($stdout);objects = "$objects";result = "$result"
    dataObject = @{"$HeadingUnder_dataObject_Array" = @($Array)}
    }
    }

    if($stderr -ne $null){
    $axa = ConvertTo-JSON2 -InputObject @{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @($stdout);result = "$result"
        stderr = @($stderr)
    }
    }
}else
{
#write-host 'its not Version 2'
    if($stderr -eq $null){
     $HashTable = [ordered]@{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @($stdout);objects = "$objects";result = "$result"
        dataObject = @{"$HeadingUnder_dataObject_Array" = @($Array)}
        }
    $axa = ConvertTo-Json -InputObject $HashTable -Depth 100
    }

    if($stderr -ne $null){
    $HashTable = [ordered]@{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @($stdout);result = "$result"
        stderr = @($stderr)
    }
    $axa = ConvertTo-Json -InputObject $HashTable -Depth 100
    }
}
#$BackTOHasTable = ConvertFrom-Json20 $axa
#$BackTOHasTable
$axa
}

########################################################### JSON Convert END





$adapter = (Gwmi win32_networkadapter | ? {$_.NetConnectionID -eq "$Adapter_Name"})
if($adapter.NetEnabled -eq $False){
#write-host 'Provided Adapter is in Disabled State'

$Stderr = Get_stderr -Title 'Provided Adapter is in Disabled State' -details "Provided Adapter is in Disabled State [$Adapter_Name]"

$HashData = HashTable_Array_JSON -taskName 'Set static IP information' -status 'Failed' -Code '1' -stdout "Provided Adapter is in Disabled State hence can't process further" -stderr $Stderr -result 'Error: Provided Adapter is in Disabled State'
$JSON2 = $HashData

}else{

$Which_Version_Setting_IPV6_OR_IPV4 = $null

function ValidateSubnetMask
{param ($strSubnetMask)
	$bValidMask = $true
	$arrSections = @()
	$arrSections +=$strSubnetMask.split(".")
	#firstly, make sure there are 4 sections in the subnet mask
	if ($arrSections.count -ne 4) {$bValidMask =$false}
	
	#secondly, make sure it only contains numbers and it's between 0-255
	if ($bValidMask)
	{
		[reflection.assembly]::LoadWithPartialName("'Microsoft.VisualBasic") | Out-Null
		foreach ($item in $arrSections)
		{
			if (!([Microsoft.VisualBasic.Information]::isnumeric($item))) {$bValidMask = $false}
		}
	}
	
	if ($bValidMask)
	{
		foreach ($item in $arrSections)
		{
			$item = [int]$item
			if ($item -lt 0 -or $item -gt 255) {$bValidMask = $false}
		}
	}
	
	#lastly, make sure it is actually a subnet mask when converted into binary format
	if ($bValidMask)
	{
		foreach ($item in $arrSections)
		{
			$binary = [Convert]::ToString($item,2)
			if ($binary.length -lt 8)
			{
				do {
				$binary = "0$binary"
				} while ($binary.length -lt 8)
			}
			$strFullBinary = $strFullBinary+$binary
		}
		if ($strFullBinary.contains("01")) {$bValidMask = $false}
		if ($bValidMask)
		{
			$strFullBinary = $strFullBinary.replace("10", "1.0")
			if ((($strFullBinary.split(".")).count -ne 2)) {$bValidMask = $false}
		}
	}
	

        if($bValidMask){
        Return "$strSubnetMask|Valid Format"
        }else{
        Return "$strSubnetMask|InValid Format"
        }
}

function IPV6_IPV4_Format_Validation{
param($Address,$IP_Type,$Address_type)

$Ip = "$Address"

if([ipaddress]::tryparse($ip,[ref]$null) ){
    
    if(([ipaddress]"$ip").AddressFamily -eq 'InterNetworkV6'){
    #"$Ip-IPV6|Valid Format"
        if($IP_Type -eq 'IPV6'){
        "$Ip|Valid Format"        
        }else{
        "$Ip|InValid Format"
        }
          
    }if(([ipaddress]"$ip").AddressFamily -eq 'InterNetwork'){
          #"$Ip-IPV4|Valid Format"
        if($IP_Type -eq 'IPV4'){
        "$Ip|Valid Format"        
        }else{
        "$Ip|InValid Format"
        }

    }

}else{

"$Ip|InValid Format"
}

}


if($Address_Type_IPv4_IPv6 -eq 'IPV4'){

    if([string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and [string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway) -and [string]::IsNullOrEmpty($IPv4_Subnet_Mask))
    { 
        #write-host 'Please provide IPV4 Network Adapter Information'
    $Global:HashTable_Data = 'Please provide IPV4 Network Adapter Information'

    $Global:Stderr += Get_stderr -Title 'Unable to process request' -details $Global:HashTable_Data
    
    JSONOutPutCall
    return          
    }

    if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and ![string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway) -and ![string]::IsNullOrEmpty($IPv4_Subnet_Mask))
    { 

       if(Validate_AdapterName -Adapter_Name $Adapter_Name){
       $GLobal:AdapterValidation = 'YES'
        #write-host 'To be change IPV4 Address--IPV4 Gateway--IPV4 SubnetMask'


            $IPV4Address = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Address -IP_Type $Address_Type_IPv4_IPv6
            $IPV4Gateway = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Gateway -IP_Type $Address_Type_IPv4_IPv6
            $IPV4Subnet = ValidateSubnetMask -strSubnetMask $IPv4_Subnet_Mask
            
            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPV4Address = $IPV4Address
            IPV4Gateway = $IPV4Gateway
            IPV4Subnet = $IPV4Subnet
            } | select 'Address Type (IPv4/IPv6)',IPV4Address,IPV4Gateway,IPV4Subnet
          

            if($Global:Format_Check.IPV4Address -like '*InValid*'){$IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address ;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV4Gateway -like '*InValid*'){$IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway ;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPV4Subnet -like '*InValid*'){$IPv4_Subnet_Mask_Stored = $IPv4_Subnet_Mask ;$IPv4_Subnet_Mask = $null}
        
                if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}

        if(!(Test-Connection $IPv4_OR_IPv6_Gateway -Quiet -Count 1)){ $GLobal:GatewayOffline = 'YES';JSONOutPutCall;return}            
        Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_OR_IPv6_Address $IPv4_OR_IPv6_Address -IPv4_OR_IPv6_Gateway $IPv4_OR_IPv6_Gateway -IPv4_Subnet_Mask $IPv4_Subnet_Mask
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    
    }
    
    
    if(! ([string]::IsNullOrEmpty($IPv4_OR_IPv6_Address)) -and [string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway) -and [string]::IsNullOrEmpty($IPv4_Subnet_Mask))
    { 
        #write-host 'To be change IPV4 Address'  

        if(Validate_AdapterName -Adapter_Name $Adapter_Name){
       $GLobal:AdapterValidation = 'YES'
            $IPV4Address = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Address -IP_Type $Address_Type_IPv4_IPv6
            
            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;          
            IPV4Address = $IPV4Address
            } | select 'Address Type (IPv4/IPv6)', IPV4Address

            if($Global:Format_Check.IPV4Address -like '*InValid*'){$IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address ;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV4Gateway -like '*InValid*'){$IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway ;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPV4Subnet -like '*InValid*'){$IPv4_Subnet_Mask_Stored = $IPv4_Subnet_Mask ;$IPv4_Subnet_Mask = $null}
            
         if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}  
        
        Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_OR_IPv6_Address $IPv4_OR_IPv6_Address
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    
    }
    
    if([string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and !([string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway)) -and [string]::IsNullOrEmpty($IPv4_Subnet_Mask))
    { 
        #write-host 'To be change IPV4 Gateway'
        if(Validate_AdapterName -Adapter_Name $Adapter_Name){
       $GLobal:AdapterValidation = 'YES'

        
                        $IPV4Gateway = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Gateway -IP_Type $Address_Type_IPv4_IPv6
            
            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPV4Gateway = $IPV4Gateway
            } | select 'Address Type (IPv4/IPv6)',IPV4Gateway

            if($Global:Format_Check.IPV4Address -like '*InValid*'){$IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address ;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV4Gateway -like '*InValid*'){$IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway ;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPV4Subnet -like '*InValid*'){$IPv4_Subnet_Mask_Stored = $IPv4_Subnet_Mask ;$IPv4_Subnet_Mask = $null}     
         
                        if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}
                  
        if(!(Test-Connection $IPv4_OR_IPv6_Gateway -Quiet -Count 1)){ $GLobal:GatewayOffline = 'YES';JSONOutPutCall;return}            
        Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_OR_IPv6_Gateway $IPv4_OR_IPv6_Gateway        
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    
    }    
    if([string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and [string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway) -and !([string]::IsNullOrEmpty($IPv4_Subnet_Mask)))
    { 
        #write-host 'To be change IPV4 SubnetMask'  

        if(Validate_AdapterName -Adapter_Name $Adapter_Name){ 
       $GLobal:AdapterValidation = 'YES'
            $IPV4Subnet = ValidateSubnetMask -strSubnetMask $IPv4_Subnet_Mask
            
            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPV4Subnet = $IPV4Subnet
            } | select 'Address Type (IPv4/IPv6)', IPV4Subnet
            
            if($Global:Format_Check.IPV4Address -like '*InValid*'){$IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address ;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV4Gateway -like '*InValid*'){$IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway ;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPV4Subnet -like '*InValid*'){$IPv4_Subnet_Mask_Stored = $IPv4_Subnet_Mask ;$IPv4_Subnet_Mask = $null}
            
                          if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}

       
        Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_Subnet_Mask $IPv4_Subnet_Mask                
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    
    }
    



    if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and [string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway) -and ![string]::IsNullOrEmpty($IPv4_Subnet_Mask))
    { 
        #write-host 'To be change IPV4 Address--Subnet_Mask' 

         if(Validate_AdapterName -Adapter_Name $Adapter_Name){
                $GLobal:AdapterValidation = 'YES'
        $IPV4Address = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Address -IP_Type $Address_Type_IPv4_IPv6
        $IPV4Subnet = ValidateSubnetMask -strSubnetMask $IPv4_Subnet_Mask
            
            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPV4Address = $IPV4Address
            IPV4Subnet = $IPV4Subnet
            } | select 'Address Type (IPv4/IPv6)', IPV4Address,IPV4Subnet

            if($Global:Format_Check.IPV4Address -like '*InValid*'){$IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address ;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV4Gateway -like '*InValid*'){$IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway ;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPV4Subnet -like '*InValid*'){$IPv4_Subnet_Mask_Stored = $IPv4_Subnet_Mask ;$IPv4_Subnet_Mask = $null}
          
                          if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}  
                      

         Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_OR_IPv6_Address $IPv4_OR_IPv6_Address -IPv4_Subnet_Mask $IPv4_Subnet_Mask         
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    
    }
    if([string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and !([string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway)) -and !([string]::IsNullOrEmpty($IPv4_Subnet_Mask)))
    { 
        #write-host 'To be change IPV4 Gateway--Subnet_Mask'  
        if(Validate_AdapterName -Adapter_Name $Adapter_Name){
               $GLobal:AdapterValidation = 'YES'

                        $IPV4Gateway = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Gateway -IP_Type $Address_Type_IPv4_IPv6
$IPV4Subnet = ValidateSubnetMask -strSubnetMask $IPv4_Subnet_Mask
            
            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPV4Gateway = $IPV4Gateway
            IPV4Subnet = $IPV4Subnet
            } | select 'Address Type (IPv4/IPv6)', IPV4Gateway,IPV4Subnet

            if($Global:Format_Check.IPV4Address -like '*InValid*'){$IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address ;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV4Gateway -like '*InValid*'){$IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway ;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPV4Subnet -like '*InValid*'){$IPv4_Subnet_Mask_Stored = $IPv4_Subnet_Mask ;$IPv4_Subnet_Mask = $null}
            
                          if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}
                      
        if(!(Test-Connection $IPv4_OR_IPv6_Gateway -Quiet -Count 1)){ $GLobal:GatewayOffline = 'YES';JSONOutPutCall;return}            
        Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_OR_IPv6_Gateway $IPv4_OR_IPv6_Gateway -IPv4_Subnet_Mask $IPv4_Subnet_Mask         
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    
    }
    if(![string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and !([string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway)) -and ([string]::IsNullOrEmpty($IPv4_Subnet_Mask)))
    { 
        #write-host 'To be change IPV4 Address--Gateway'  
         if(Validate_AdapterName -Adapter_Name $Adapter_Name){
                $GLobal:AdapterValidation = 'YES'
                
                        $IPV4Address = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Address -IP_Type $Address_Type_IPv4_IPv6
                        $IPV4Gateway = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Gateway -IP_Type $Address_Type_IPv4_IPv6
            
            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPV4Address = $IPV4Address
            IPV4Gateway = $IPV4Gateway
            } | select 'Address Type (IPv4/IPv6)', IPV4Address,IPV4Gateway

            if($Global:Format_Check.IPV4Address -like '*InValid*'){$IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address ;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV4Gateway -like '*InValid*'){$IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway ;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPV4Subnet -like '*InValid*'){$IPv4_Subnet_Mask_Stored = $IPv4_Subnet_Mask ;$IPv4_Subnet_Mask = $null}
                      
                 if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}else{} 
                                    
        if(!(Test-Connection $IPv4_OR_IPv6_Gateway -Quiet -Count 1)){ $GLobal:GatewayOffline = 'YES';JSONOutPutCall;return}            
         Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_OR_IPv6_Gateway $IPv4_OR_IPv6_Gateway -IPv4_OR_IPv6_Address $IPv4_OR_IPv6_Address          
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    
    }        
}

if($Address_Type_IPv4_IPv6 -eq 'IPV6'){


    if([string]::IsNullOrEmpty($IPv4_OR_IPv6_Address) -and [string]::IsNullOrEmpty($IPv4_OR_IPv6_Gateway) -and [string]::IsNullOrEmpty($IPv6_Subnet_Prefix_Length))
    { 
        #write-host 'Please provide IPV6 Network Adapter Information'
    $Global:HashTable_Data = 'Please provide IPV6 Network Adapter Information'
    $Global:Stderr += Get_stderr -Title 'Unable to process request' -details $Global:HashTable_Data    
    JSONOutPutCall
    return                    
    }
    
    if(![string]::IsNullOrEmpty($IPV4_OR_IPv6_Address) -and ![string]::IsNullOrEmpty($IPV4_OR_IPv6_Gateway) -and ![string]::IsNullOrEmpty($IPv6_Subnet_Prefix_Length))
    { 
        #write-host 'To be change IPV6 Address--IPV6 Gateway--IPV6 Subnet_Prefix_Length'  
         if(Validate_AdapterName -Adapter_Name $Adapter_Name){
                $GLobal:AdapterValidation = 'YES'
            $IPV6Address = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Address -IP_Type $Address_Type_IPv4_IPv6
            $IPV6Gateway = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Gateway -IP_Type $Address_Type_IPv4_IPv6
            $IPv6_Subnet_Prefix_Length1 = if(@(8..128) -contains $IPv6_Subnet_Prefix_Length){ 'Valid IPv6 Subnet Prefix Length' }else{ "$IPv6_Subnet_Prefix_Length|InValid IPv6 Subnet Prefix Length"}

            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPV6Address = $IPV6Address
            IPV6Gateway = $IPV6Gateway
            IPv6_Subnet_Prefix_Length1 = $IPv6_Subnet_Prefix_Length1
            } | select 'Address Type (IPv4/IPv6)', IPV6Address,IPV6Gateway,IPv6_Subnet_Prefix_Length1

            if($Global:Format_Check.IPV6Address -like '*InValid*'){ $IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV6Gateway -like '*InValid*'){ $IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPv6_Subnet_Prefix_Length1 -like '*InValid*'){ $IPv6_Subnet_Prefix_Length_Stored = $IPv6_Subnet_Prefix_Length;$IPv6_Subnet_Prefix_Length = $null}

           if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}else{}
         
        if(!(Test-Connection $IPv4_OR_IPv6_Gateway -Quiet -Count 1)){ $GLobal:GatewayOffline = 'YES';JSONOutPutCall;return}            
        Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_OR_IPv6_Address $IPv4_OR_IPv6_Address -IPv4_OR_IPv6_Gateway $IPv4_OR_IPv6_Gateway -IPv6_Subnet_Prefix_Length $IPv6_Subnet_Prefix_Length        
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    }
    
    if(! ([string]::IsNullOrEmpty($IPV4_OR_IPv6_Address)) -and [string]::IsNullOrEmpty($IPV4_OR_IPv6_Gateway) -and [string]::IsNullOrEmpty($IPv6_Subnet_Prefix_Length))
    { 
        #write-host 'To be change IPV6 Address'  
         if(Validate_AdapterName -Adapter_Name $Adapter_Name){
                $GLobal:AdapterValidation = 'YES'
            $IPV6Address = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Address -IP_Type $Address_Type_IPv4_IPv6

            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPV6Address = $IPV6Address
            } | select 'Address Type (IPv4/IPv6)', IPV6Address

            if($Global:Format_Check.IPV6Address -like '*InValid*'){ $IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV6Gateway -like '*InValid*'){ $IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPv6_Subnet_Prefix_Length1 -like '*InValid*'){ $IPv6_Subnet_Prefix_Length_Stored = $IPv6_Subnet_Prefix_Length;$IPv6_Subnet_Prefix_Length = $null}
            
                          if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}

        Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_OR_IPv6_Address $IPv4_OR_IPv6_Address
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    
    }
    
    if([string]::IsNullOrEmpty($IPV4_OR_IPv6_Address) -and !([string]::IsNullOrEmpty($IPV4_OR_IPv6_Gateway)) -and [string]::IsNullOrEmpty($IPv6_Subnet_Prefix_Length))
    { 
        #write-host 'To be change IPV6 Gateway'  
         if(Validate_AdapterName -Adapter_Name $Adapter_Name){
                $GLobal:AdapterValidation = 'YES'

            $IPV6Gateway = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Gateway -IP_Type $Address_Type_IPv4_IPv6

            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPV6Gateway = $IPV6Gateway
            } | select 'Address Type (IPv4/IPv6)', IPV6Gateway

            if($Global:Format_Check.IPV6Address -like '*InValid*'){ $IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV6Gateway -like '*InValid*'){ $IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPv6_Subnet_Prefix_Length1 -like '*InValid*'){ $IPv6_Subnet_Prefix_Length_Stored = $IPv6_Subnet_Prefix_Length;$IPv6_Subnet_Prefix_Length = $null}
         
                         if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}   
                      
        if(!(Test-Connection $IPv4_OR_IPv6_Gateway -Quiet -Count 1)){ $GLobal:GatewayOffline = 'YES';JSONOutPutCall;return}            
        Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_OR_IPv6_Gateway $IPv4_OR_IPv6_Gateway
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    
    }
    if([string]::IsNullOrEmpty($IPV4_OR_IPv6_Address) -and [string]::IsNullOrEmpty($IPV4_OR_IPv6_Gateway) -and !([string]::IsNullOrEmpty($IPv6_Subnet_Prefix_Length)))
    { 
        #write-host 'Please specify the IPV6 Address also to set Subnet_Prefix_Length'
         if(Validate_AdapterName -Adapter_Name $Adapter_Name){
                $GLobal:AdapterValidation = 'YES'
        $Global:HashTable_Data = 'Please specify the IPV6 Address also to set Subnet_Prefix_Length'
        $Global:Stderr += Get_stderr -Title 'Unable to process request' -details $Global:HashTable_Data
        JSONOutPutCall
        return     
            $IPv6_Subnet_Prefix_Length1 = if(@(8..128) -contains $IPv6_Subnet_Prefix_Length){ 'Valid IPv6 Subnet Prefix Length' }else{ "$IPv6_Subnet_Prefix_Length|InValid IPv6 Subnet Prefix Length"}
            
            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPv6_Subnet_Prefix_Length1 = $IPv6_Subnet_Prefix_Length1
            } | select 'Address Type (IPv4/IPv6)',IPv6_Subnet_Prefix_Length1

            if($Global:Format_Check.IPV6Address -like '*InValid*'){ $IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV6Gateway -like '*InValid*'){ $IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPv6_Subnet_Prefix_Length1 -like '*InValid*'){ $IPv6_Subnet_Prefix_Length_Stored = $IPv6_Subnet_Prefix_Length;$IPv6_Subnet_Prefix_Length = $null}
            
                          if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}
                      

        Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv6_Subnet_Prefix_Length $IPv6_Subnet_Prefix_Length  
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
        
        }
        
    if(!([string]::IsNullOrEmpty($IPV4_OR_IPv6_Address)) -and [string]::IsNullOrEmpty($IPV4_OR_IPv6_Gateway) -and !([string]::IsNullOrEmpty($IPv6_Subnet_Prefix_Length)))
    { 
        #write-host 'To be change IPV6 Address--Subnet_Prefix_Length' 

         if(Validate_AdapterName -Adapter_Name $Adapter_Name){
                $GLobal:AdapterValidation = 'YES'
            $IPV6Address = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Address -IP_Type $Address_Type_IPv4_IPv6
            $IPv6_Subnet_Prefix_Length1 = if(@(8..128) -contains $IPv6_Subnet_Prefix_Length){ 'Valid IPv6 Subnet Prefix Length' }else{ "$IPv6_Subnet_Prefix_Length|InValid IPv6 Subnet Prefix Length"}
            
            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPV6Address = $IPV6Address
            IPv6_Subnet_Prefix_Length1 = $IPv6_Subnet_Prefix_Length1
            } | select 'Address Type (IPv4/IPv6)', IPV6Address,IPv6_Subnet_Prefix_Length1
            if($Global:Format_Check.IPV6Address -like '*InValid*'){ $IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV6Gateway -like '*InValid*'){ $IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPv6_Subnet_Prefix_Length1 -like '*InValid*'){ $IPv6_Subnet_Prefix_Length_Stored = $IPv6_Subnet_Prefix_Length;$IPv6_Subnet_Prefix_Length = $null}
            
                 if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}else{}
                                       

        Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_OR_IPv6_Address $IPv4_OR_IPv6_Address -IPv6_Subnet_Prefix_Length $IPv6_Subnet_Prefix_Length        
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    
    }
    if([string]::IsNullOrEmpty($IPV4_OR_IPv6_Address) -and !([string]::IsNullOrEmpty($IPV4_OR_IPv6_Gateway)) -and !([string]::IsNullOrEmpty($IPv6_Subnet_Prefix_Length)))
    { 
        #write-host 'To be change IPV6 Gateway--Subnet_Prefix_Length'  
        if(Validate_AdapterName -Adapter_Name $Adapter_Name){
               $GLobal:AdapterValidation = 'YES'


            $IPV6Gateway = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Gateway -IP_Type $Address_Type_IPv4_IPv6
            $IPv6_Subnet_Prefix_Length1 = if(@(8..128) -contains $IPv6_Subnet_Prefix_Length){ 'Valid IPv6 Subnet Prefix Length' }else{ "$IPv6_Subnet_Prefix_Length|InValid IPv6 Subnet Prefix Length"}

            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPV6Gateway = $IPV6Gateway
            IPv6_Subnet_Prefix_Length1 = $IPv6_Subnet_Prefix_Length1
            } | select 'Address Type (IPv4/IPv6)', IPV6Gateway,IPv6_Subnet_Prefix_Length1

            if($Global:Format_Check.IPV6Address -like '*InValid*'){ $IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV6Gateway -like '*InValid*'){ $IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPv6_Subnet_Prefix_Length1 -like '*InValid*'){ $IPv6_Subnet_Prefix_Length_Stored = $IPv6_Subnet_Prefix_Length;$IPv6_Subnet_Prefix_Length = $null}
            
                          if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}
        if(!(Test-Connection $IPv4_OR_IPv6_Gateway -Quiet -Count 1)){ $GLobal:GatewayOffline = 'YES';JSONOutPutCall;return}            
        Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_OR_IPv6_Gateway $IPv4_OR_IPv6_Gateway -IPv6_Subnet_Prefix_Length $IPv6_Subnet_Prefix_Length        
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    
    }
    if(![string]::IsNullOrEmpty($IPV4_OR_IPv6_Address) -and !([string]::IsNullOrEmpty($IPV4_OR_IPv6_Gateway)) -and ([string]::IsNullOrEmpty($IPv6_Subnet_Prefix_Length)))
    { 
        #write-host 'To be change IPV6 Address--Gateway 333'
        
        if(Validate_AdapterName -Adapter_Name $Adapter_Name){
               $GLobal:AdapterValidation = 'YES'

        
            $IPV6Address = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Address -IP_Type $Address_Type_IPv4_IPv6
            $IPV6Gateway = IPV6_IPV4_Format_Validation -Address $IPv4_OR_IPv6_Gateway -IP_Type $Address_Type_IPv4_IPv6


            $Global:Format_Check = new-object psobject -Property @{'Address Type (IPv4/IPv6)' = $Address_Type_IPv4_IPv6;
            IPV6Address = $IPV6Address
            IPV6Gateway = $IPV6Gateway
            } | select 'Address Type (IPv4/IPv6)', IPV6Address,IPV6Gateway
            
            if($Global:Format_Check.IPV6Address -like '*InValid*'){ $IPv4_OR_IPv6_Address_Stored = $IPv4_OR_IPv6_Address;$IPv4_OR_IPv6_Address = $null}      
            if($Global:Format_Check.IPV6Gateway -like '*InValid*'){ $IPv4_OR_IPv6_Gateway_Stored = $IPv4_OR_IPv6_Gateway;$IPv4_OR_IPv6_Gateway = $null}
            if($Global:Format_Check.IPv6_Subnet_Prefix_Length1 -like '*InValid*'){ $IPv6_Subnet_Prefix_Length_Stored = $IPv6_Subnet_Prefix_Length;$IPv6_Subnet_Prefix_Length = $null}
         
                         if($Global:Format_Check | % {$_ -like '*InValid*'}){JSONOutPutCall;return}
        
            if(!(Test-Connection $IPv4_OR_IPv6_Gateway -Quiet -Count 1)){ $GLobal:GatewayOffline = 'YES';JSONOutPutCall;return}
            
         Change_Adapter_Settings -Adapter_Name $Adapter_Name -Which_Version_Setting_IPV6_OR_IPV4 "$Address_Type_IPv4_IPv6" -IPv4_OR_IPv6_Address $IPv4_OR_IPv6_Address -IPv4_OR_IPv6_Gateway $IPv4_OR_IPv6_Gateway
        }else{$GLobal:AdapterValidation = 'NO';JSONOutPutCall;return}
    
    }
}


sleep 5


}


$JSONOUT = JSONOutPutCall
$JSONOUT