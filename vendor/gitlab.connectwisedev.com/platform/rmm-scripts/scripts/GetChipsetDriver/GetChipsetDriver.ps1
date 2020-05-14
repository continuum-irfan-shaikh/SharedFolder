
Clear-Host

<#

cat - 

[string]$ChipsetDriverName = "Microsoft Virtual WiFi Miniport Adapter" -optional 

#>

function Escape-JSONString($str){
    if ($str -eq $null) {return ""}
         $str = $str.ToString().Replace('"','\"').Replace('\','\\').Replace("`n",'\n').Replace("`r",'\r').Replace("`t",'\t')
    return $str;
}


function ConvertTo-JSONP2($maxDepth = 10,$forceArray = $false) {
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
$jsonResult += ($elem | ConvertTo-JSONP2 -maxDepth ($maxDepth -1))
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
"@ -f $key , ($value[$key] | ConvertTo-JSONP2 -maxDepth ($maxDepth -1) )
}
return "{" + $jsonResult + "}"
}
            default { #object
if ($maxDepth -le 0){return  "`"{0}`"" -f (Escape-JSONString $value)}
return "{" +
(($value | Get-Member -MemberType *property | % { 
@"
"{0}": {1}
"@ -f $_.Name , ($value.($_.Name) | ConvertTo-JSONP2 -maxDepth ($maxDepth -1) ) 
}) -join ', ') + "}"
    }
}
}
}


function Get-ChipsetDriver{
    
   $ChipsetObject = New-Object -TypeName psobject
   $ChipsetObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "Get Chipset Driver"

    $PnpIDs = @()
    foreach( $c in Get-WmiObject Win32_PNPEntity){
        $PnpIDs +=  $c.PNPDeviceID
    }      
 
    $root = Get-Childitem "HKLM:\SYSTEM\CurrentControlSet\Enum"
    $classroot = Get-item "HKLM:\SYSTEM\CurrentControlSet\Control\Class"
    $alldevices = @()

    foreach ($i in $root)
    {
        $Class = $i.PSChildName
 
        # Build the subkey path    
        $subkeylevel1 = Join-Path -Path $i.PSParentPath -ChildPath $i.PSChildName
 
        # Get its properties
        $subkeyslevel2 = Get-Childitem -Path $subkeylevel1
 
        foreach ($j in $subkeyslevel2)
        {
            $subkeylevel3 = Join-Path -Path $j.PSParentPath -ChildPath $j.PSChildName
            $HardwareID = $j.PSChildName
            $subkeyslevel4 = Get-Childitem -LiteralPath $subkeylevel3
            foreach ($k in $subkeyslevel4)
            {
                $properties = $FriendlyName = $desc = $null
                $properties = Get-ItemProperty $k.PSPath
                if ($properties.FriendlyName -eq $null)
                {
                    if ($properties.DeviceDesc -match "^@")
                    {
                        $FriendlyName = ($properties.DeviceDesc -split ";")[1]
                        if ($FriendlyName -eq $null)
                        {
                            $FriendlyName = ($properties.DeviceDesc -split ";")[0]
                            switch($FriendlyName)
                            {
                                "@%systemroot%\system32\drivers\afd.sys,-1000"     { $desc = "Ancillary Function Driver for Winsock" }
                                "@%systemroot%\system32\appidsvc.dll,-102"         { $desc = "AppID Driver" }
                                "@%systemroot%\system32\browser.dll,-102"          { $desc = "Browser Support Driver" }
                                "@%SystemRoot%\system32\clfs.sys,-100"             { $desc = "Common Log (CLFS)" }
                                "@%systemroot%\system32\cscsvc.dll,-202"           { $desc = "Offline Files Driver" }
                                "@%systemroot%\system32\drivers\dfsc.sys,-101"     { $desc = "DFS Namespace Client Driver" }
                                "@%systemroot%\system32\drivers\discache.sys,-102" { $desc = "System Attribute Cache" }
                                "@%SystemRoot%\system32\drivers\fileinfo.sys,-100" { $desc = "File Information FS MiniFilter" }
                                "@%SystemRoot%\system32\drivers\fltmgr.sys,-10001" { $desc = "FltMgr" }
                                "@%SystemRoot%\system32\drivers\fvevol.sys,-100"   { $desc = "Bitlocker Drive Encryption Filter Driver" }
                                "@%SystemRoot%\system32\drivers\http.sys,-1"       { $desc = "HTTP" }
                                "@%systemroot%\system32\drivers\hwpolicy.sys,-101" { $desc = "Hardware Policy Driver" }
                                "@%systemroot%\system32\drivers\luafv.sys,-100"    { $desc = "UAC File Virtualization" }
                                "@%SystemRoot%\system32\drivers\mountmgr.sys,-100" { $desc = "Mount Point Manager" }
                                "@%SystemRoot%\system32\FirewallAPI.dll,-23092"    { $desc = "Windows Firewall Authorization Driver" }
                                "@%systemroot%\system32\webclnt.dll,-104"          { $desc = "WebDav Client Redirector Driver" }
                                "@%systemroot%\system32\wkssvc.dll,-1002"          { $desc = "SMB MiniRedirector Wrapper and Engine" }
                                "@%systemroot%\system32\wkssvc.dll,-1004"          { $desc = "SMB 1.x MiniRedirector" }
                                "@%systemroot%\system32\wkssvc.dll,-1006"          { $desc = "SMB 2.0 MiniRedirector" }
                                "@%systemroot%\system32\drivers\mup.sys,-101"      { $desc = "MUP" }
                                "@%SystemRoot%\system32\drivers\ndis.sys,-200"     { $desc = "NDIS System Driver" }
                                "@%SystemRoot%\system32\drivers\netbt.sys,-2"      { $desc = "NETBT" }
                                "@%SystemRoot%\system32\drivers\nsiproxy.sys,-2"   { $desc = "NSI proxy service driver." }
                                "@%SystemRoot%\System32\drivers\pacer.sys,-101"    { $desc = "QoS Packet Scheduler" }
                                "@%systemroot%\system32\wkssvc.dll,-1000"          { $desc = "Redirected Buffering Sub Sysytem" }
                                "@%systemroot%\system32\DRIVERS\RDPCDD.sys,-100"   { $desc = "RDPCDD" }
                                "@%systemroot%\system32\drivers\RDPENCDD.sys,-101" { $desc = "RDP Encoder Mirror Driver" }
                                "@%systemroot%\system32\drivers\RdpRefMp.sys,-101" { $desc = "Reflector Display Driver used to gain access to graphics data" }
                                "@%systemroot%\system32\srvsvc.dll,-102"           { $desc = "Server SMB 1.xxx Driver" }
                                "@%systemroot%\system32\srvsvc.dll,-104"           { $desc = "Server SMB 2.xxx Driver" }
                                "@%SystemRoot%\system32\vmstorfltres.dll,-1000"    { $desc = "Disk Virtual Machine Bus Acceleration Filter Driver" }
                                "@%SystemRoot%\system32\tcpipcfg.dll,-50003"       { $desc = "TCP/IP Protocol Driver" }
                                "@%SystemRoot%\system32\tcpipcfg.dll,-50004"       { $desc = "NetIO Legacy TDI Support Driver" }
                                "@%SystemRoot%\System32\DRIVERS\tssecsrv.sys,-101" { $desc = "Remote Desktop Services Security Filter Driver" }
                                "@%SystemRoot%\system32\drivers\volmgrx.sys,-100"  { $desc = "Dynamic Volume Manager" }
                                "@%systemroot%\system32\rascfg.dll,-32012"         { $desc = "Remote Access IPv6 ARP Driver" }
                                default                                            { $desc = ""}
                            }
                            $FriendlyName  = $desc
                        }
 
                    } else {
                        $FriendlyName = $properties.DeviceDesc
                    }
                } else {
                    if ($properties.FriendlyName -match "^@")
                    {
                        $FriendlyName = ($properties.FriendlyName -split ";")[1]
                    } else {
                        $FriendlyName = $properties.FriendlyName
                    }
                }
             
                $Object = $null
                # Build an object to store all the properties we are interested in
                $Object = New-Object -TypeName PSObject -Property @{
                   
                    ChipsetName = $FriendlyName
                    HwID = $Class + "\" + $HardwareID + "\" + $k.PSChildName                   
                    CompatibleIDs = $properties.CompatibleIDs
                }
                # Now that we have a link to its driver info, we can gather additional info from the registry
                if ($properties.Driver -ne $null)
                {
                    $driverinfopath = Join-Path -Path $classroot.PSPath -ChildPath $properties.Driver
                    if (Test-Path $driverinfopath)
                    {
                        $driverproperties = Get-ItemProperty -Path $driverinfopath

                       # Write-Host $driverproperties `n
 
                        $DriverProperties = New-Object -TypeName PSObject -Property @{            
                            InfFilePath = "$env:systemroot\inf\" + $driverproperties.InfPath
                            InfSection = $driverproperties.InfSection
                            DriverDescription = $driverproperties.DriverDesc
                            Manufacturer = $driverproperties.ProviderName
                            DriverDate = [string]$driverproperties.DriverDate
                            DriverVersion = [version]$driverproperties.DriverVersion
                        }
                        $Object | add-member Noteproperty -Name DriverVersion -Value $DriverProperties.DriverVersion
                        $Object | add-member Noteproperty -Name DriverDate -Value $DriverProperties.DriverDate
                        $Object | add-member Noteproperty -Name VendorID -Value $DriverProperties.Manufacturer
 
                    }
                }
                # Add the object to our array    
                $alldevices += $Object
            } # end of foreach level 4
        } # end of foreach level 2
        # break
    } # end of foreach root
 
   
   $hostVerSionMajor = ($PSVersionTable.PSVersion.Major).ToString()
   $hostVerSionMinor = ($PSVersionTable.PSVersion.Minor).ToString()
   $hostVersion = $hostVerSionMajor +'.'+ $hostVerSionMinor 
 
   $osVersionMajor = ([System.Environment]::OSVersion.Version.major).ToString()
   $osVersionMinor = ([System.Environment]::OSVersion.Version.minor).ToString()
   $osVersion = $osVersionMajor +'.'+ $osVersionMinor
 
   [boolean]$isPsVersionOk = ([version]$hostVersion -ge [version]'2.0')
   [boolean]$isOSVersionOk = ([version]$osVersion -ge [version]'6.0')
   
   #---------- Check for Powershell version ------------------         
   
   if(-not $isPsVersionOk){
       
      $StdErrArr = @()
      $stdOutArr = @()
      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "PowerShell version below 2.0 is not supported";
                          detail = "PowerShell version below 2.0 is not supported";

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ("PowerShell version below 2.0 is not supported")
    
     $ChipsetObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $ChipsetObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $ChipsetObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $ChipsetObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
     $ChipsetObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
    
     return $ChipsetObject
    }
 
   #---------- Check OS Version ------------------
  
   if(-not $isOSVersionOk){
        
      $StdErrArr = @()
      $stdOutArr = @()
      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system";
                          detail = "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system";

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ("PowerShell Script supports Window 7, Window 2008R2 and higher version operating system")
    
     $ChipsetObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $ChipsetObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $ChipsetObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $ChipsetObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
     $ChipsetObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr
    
     return $ChipsetObject
 
    }

    if([string]::IsNullOrEmpty($ChipsetDriverName)){
    
      $ChipsetObjectArr = @()
      $ChipsetObjectStrArr = @()

      foreach($d in $alldevices){        
        foreach($p in $PnpIDs){
           if($d.HwID -eq $p){

             $ChipSetDriver ="ChipSet Driver " + $Counter         
             $Object1 = New-Object PSObject -Property @{		       
		                 ChipsetName = $d.ChipsetName;
                         DriverVersion = $d.DriverVersion;
                         DriverDate = $d.DriverDate;
                         VendorID = $d.VendorID;
                         HwID = $d.HwID
                    }
                                
             $ChipsetName = $d.ChipsetName;
             $DriverVersion = $d.DriverVersion;
             $DriverDate = $d.DriverDate;
             $VendorID = $d.VendorID;
             $HwID = $d.HwID

            $ChipsetObjectArr += $Object1 
            $ChipsetObjectStrArr += "ChipsetName : $ChipsetName, DriverVersion : $DriverVersion,`
                                     DriverDate : $DriverDate, VendorID : $VendorID,`
                                     HwID : $HwID "
      
           }
        }
      }

       $InstalledChipset =  New-Object psobject
       $InstalledChipset | Add-Member -MemberType NoteProperty -Name ChipsetDrivers -Value $ChipsetObjectArr

       $ChipsetObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
       $ChipsetObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
       $ChipsetObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : ChipSet Driver retrived"

       $ChipsetObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $ChipsetObjectStrArr
       $ChipsetObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $InstalledChipset  
          
      
      return $ChipsetObject

    } # empty string if ends
    else{
           
      $ChipsetDriver = $alldevices | Where-Object { $_.ChipsetName -match [regex]::escape($ChipsetDriverName) } `
                | Select-Object -Property ChipsetName,DriverVersion,DriverDate,VendorID,HwID -ErrorAction SilentlyContinue
                      
      $ChipsetObjectArr = @()
      $ChipsetObjectStrArr = @()

      if($ChipsetDriver){
                      
        foreach($d in $ChipsetDriver){  
            
            $ChipSetDriver ="ChipSet Driver " + $Counter 
            $Object1 = New-Object PSObject -Property @{		       
		                     ChipsetName = $d.ChipsetName;
                             DriverVersion = $d.DriverVersion;
                             DriverDate = $d.DriverDate;
                             VendorID = $d.VendorID;
                             HwID = $d.HwID
                        } 
                               
            $ChipsetName = $d.ChipsetName;
            $DriverVersion = $d.DriverVersion;
            $DriverDate = $d.DriverDate;
            $VendorID = $d.VendorID;
            $HwID = $d.HwID

            $ChipsetObjectArr += $Object1 
            $ChipsetObjectStrArr += "ChipsetName : $ChipsetName, DriverVersion : $DriverVersion,`
                                     DriverDate : $DriverDate, VendorID : $VendorID, `
                                     HwID : $HwID "             
        }
        
       $InstalledChipset =  New-Object psobject
       $InstalledChipset | Add-Member -MemberType NoteProperty -Name ChipsetDrivers -Value $ChipsetObjectArr

       $ChipsetObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
       $ChipsetObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
       $ChipsetObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : Details of chipset drivers :$($ChipsetDriverName)"

       $ChipsetObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $ChipsetObjectStrArr
       $ChipsetObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $InstalledChipset   
        return $ChipsetObject

      }else{

     
      $ChipsetObjectArr = @()
      $ChipsetObjectStrArr = @()
      $ChipsetObjectErrArr = @()

      foreach($d in $alldevices){        
        foreach($p in $PnpIDs){
           if($d.HwID -eq $p){

             $ChipSetDriver ="ChipSet Driver " + $Counter         
             $Object1 = New-Object PSObject -Property @{		       
		                 ChipsetName = $d.ChipsetName;
                         DriverVersion = $d.DriverVersion;
                         DriverDate = $d.DriverDate;
                         VendorID = $d.VendorID;
                         HwID = $d.HwID
                    }
                                
             $ChipsetName = $d.ChipsetName;
             $DriverVersion = $d.DriverVersion;
             $DriverDate = $d.DriverDate;
             $VendorID = $d.VendorID;
             $HwID = $d.HwID

            $ChipsetObjectArr += $Object1 
            $ChipsetObjectStrArr += "ChipsetName : $ChipsetName, DriverVersion : $DriverVersion,`
                                     DriverDate : $DriverDate, VendorID : $VendorID,`
                                     HwID : $HwID "
      
           }
        }
      }

        $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = "chipset drivers :$($ChipsetDriverName) is not installed"
                          detail = "List of all installed chipset drivers"

               }
       
       $ChipsetObjectErrArr += $StdErr
       $InstalledChipset =  New-Object psobject
       $InstalledChipset | Add-Member -MemberType NoteProperty -Name ChipsetDrivers -Value $ChipsetObjectArr

       $ChipsetObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
       $ChipsetObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
       $ChipsetObject | Add-Member -MemberType NoteProperty -Name stderr -Value $ChipsetObjectErrArr
       $ChipsetObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error : chipset drivers :$($ChipsetDriverName) is not installed"

       $ChipsetObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $ChipsetObjectStrArr
       $ChipsetObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $InstalledChipset  
          
      
      return $ChipsetObject

      }

    } # empty string if ends

}


 if($PSVersionTable.PSVersion.Major -eq 2){

     Get-ChipsetDriver | ConvertTo-JSONP2

}else{

     Get-ChipsetDriver | ConvertTo-Json -Depth 10

}
