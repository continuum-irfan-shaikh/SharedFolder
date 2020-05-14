<#
Name : Configure display settings
Cetegory : Setup 

DISCRIPTION = "Set the DNS of a site."

$oldDNSIP  			===> Type:String 	===> TextBox(Mandatory) ===> Title:"Old DNS Server IP Address:"
$configType 		===> Type:String 	===> Radio Button (Static / Dynamic) ===> Title:"IP Configuration:"
$primaryServerIP 	===> Type:String 	===> TextBox(Mandatory) ===> Title:"Primary Server IP Address:" 
$secondaryServerIP 	===> Type:String 	===> TextBox ===> Title:"Secondary Server IP Address:" 
$updateIfPingFail 	===> Type:Boolean 	===> CheckBox ===> Title:"Update DNS Settings even if pinging the IPs fails"

#>

# Declare Variables

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}



# Validation for NOT NULL or empty for mandatory parameters
if(($oldDNSIP -eq $NULL) -or ($oldDNSIP -eq ''))
  {
     Write-Host "Old DNS server IP can not be empty"  
     EXIT
  }


   #Validate IP Format
   $validateOldDNSIP=$oldDNSIP -as [ipaddress]
   if($validateOldDNSIP -eq $NULL)
   {
       Write-Host "Old DNS IP Address is not in correct format: $oldDNSIP"
       EXIT
   }
   #####################################################################################

 if($configType -eq "static")
 {
   [string[]]$dns=@()
    $dns+=$primaryServerIP
   #Creating array of Secondary IP Address not null
   if(!(($secondaryServerIP -eq $null) -or ($secondaryServerIP -eq ""))) 
   {
      $dns+=$secondaryServerIP.Split(",")
     
   }
   
   #Validating IP address for Secondary IP Addresses. 
   $flag=$false
   foreach($ip in $dns)
     {
       Write-Host "IP=$ip"
         $validateIP=$ip -as [ipaddress]
          if($validateIP -eq $NULL)
          {
          $flag=$true
           }
     }
    if($flag -eq $true)
    {
     Write-Host "1 or more IP address is not in correct format in Secondary IP address"
     EXIT
    }

    #Ping IP address before changing the settings, check when updateIfPingFail = $False
     
  if($updateIfPingFail -eq $False)
    {
         $NotPingingIPs=@()
         foreach($ip in $dns)
         {
            $conn=Test-Connection $ip -Count 4 -Quiet
            if($conn)
            {
            continue
            } 
            else
            {
              $NotPingingIPs+=$ip
            }
         }

         if($NotPingingIPs -ne $NULL)
          {
            Write-Host "Unable to update the IPs, below IPs are not reachable"
            foreach($ip in $NotPingingIPs)
            {
               $ip
            }
            EXIT
          }
    }
     
 }
 else
 {
   $dns=$NULL
 } 



############################################################################
# Function to add Static/dynamic DNS server IPs
############################################################################

Function ConfigureDNSSettings($oldDNSIP,$dns)
{

Try
{
    $adapter = Get-WmiObject win32_networkadapterconfiguration -filter "ipenabled ='true'" |Where-Object{$_.DNSServerSearchOrder -like $oldDNSIP}
    
    if($adapter -eq $NULL)
    {
      Write-Host "Network adapter not found for given IP: $oldDNSIP"
      EXIT
    }
    if($dns -eq $Null)
    {
        # Configure the DNS Servers automatically
        $adapter.SetDNSServerSearchOrder() | Out-Null
        
        # Enable DHCP
        $adapter.EnableDHCP() | Out-Null
        write-host "Dynamic DNS server IP has been updated successfully"
    }
    else
    {
        $adapter.SetDNSServerSearchOrder($dns) | Out-Null;
        write-host "DNS server IP has been updated successfully"
    }
 }
 Catch
 {
 $Host.UI.WriteErrorLine($Error[0].Exception.Message)
 }
}
########## End Function ConfigureDNSSettings 

#Calling Function
Try
{

   ConfigureDNSSettings -oldDNSIP $oldDNSIP -dns $dns
}
Catch{

 $Host.UI.WriteErrorLine($Error[0].Exception.Message)
}
