
Clear-Host


<# 

cat-Network

[string]$DNSQuery = "google.com" -mandatory  

[bool]$DNS_TYPE_A = $true         -All are optiional
[bool]$DNS_TYPE_AAAA = $true
[bool]$DNS_TYPE_Cname = $true
[bool]$DNS_TYPE_SRV = $true
[bool]$DNS_TYPE_MX = $true
[bool]$DNS_TYPE_NS = $true
[bool]$DNS_TYPE_PTR = $true
[bool]$DNS_TYPE_TXT = $true

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


function Get-DnsAddressList
{
    param(
        [parameter(Mandatory=$true)][Alias("Host")]
          [string]$HostName)

    try {
        return [System.Net.Dns]::GetHostEntry($HostName).AddressList
    }
    catch [System.Net.Sockets.SocketException] {
        if ($_.Exception.ErrorCode -ne 11001) {
            throw $_
        }
        return = @()
    }
}

function Get-DnsMXQuery
{
    param(
        [parameter(Mandatory=$true)]
          [string]$DomainName)

    if (-not $Script:global_dnsquery) {
        $Private:SourceCS = @'
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Runtime.InteropServices;

namespace PM.Dns {
  public class MXQuery {
    [DllImport("dnsapi", EntryPoint="DnsQuery_W", CharSet=CharSet.Unicode, SetLastError=true, ExactSpelling=true)]
    private static extern int DnsQuery(
        [MarshalAs(UnmanagedType.VBByRefStr)]
        ref string pszName, 
        ushort     wType, 
        uint       options, 
        IntPtr     aipServers, 
        ref IntPtr ppQueryResults, 
        IntPtr pReserved);

    [DllImport("dnsapi", CharSet=CharSet.Auto, SetLastError=true)]
    private static extern void DnsRecordListFree(IntPtr pRecordList, int FreeType);

    public static string[] Resolve(string domain)
    {
        if (Environment.OSVersion.Platform != PlatformID.Win32NT)
            throw new NotSupportedException();

        List<string> list = new List<string>();

        IntPtr ptr1 = IntPtr.Zero;
        IntPtr ptr2 = IntPtr.Zero;
        int num1 = DnsQuery(ref domain, 15, 0, IntPtr.Zero, ref ptr1, IntPtr.Zero);
        if (num1 != 0)
            throw new Win32Exception(num1);
        try {
            MXRecord recMx;
            for (ptr2 = ptr1; !ptr2.Equals(IntPtr.Zero); ptr2 = recMx.pNext) {
                recMx = (MXRecord)Marshal.PtrToStructure(ptr2, typeof(MXRecord));
                if (recMx.wType == 15)
                    list.Add(Marshal.PtrToStringAuto(recMx.pNameExchange));
            }
        }
        finally {
            DnsRecordListFree(ptr1, 0);
        }

        return list.ToArray();
    }

    [StructLayout(LayoutKind.Sequential)]
    private struct MXRecord
    {
        public IntPtr pNext;
        public string pName;
        public short  wType;
        public short  wDataLength;
        public int    flags;
        public int    dwTtl;
        public int    dwReserved;
        public IntPtr pNameExchange;
        public short  wPreference;
        public short  Pad;
    }
  }
}
'@

        Add-Type -TypeDefinition $Private:SourceCS -ErrorAction Stop
        $Script:global_dnsquery = $true
    }

    [PM.Dns.MXQuery]::Resolve($DomainName) | % {
        $rec = New-Object PSObject
        Add-Member -InputObject $rec -MemberType NoteProperty -Name "MailExchange"        -Value $_
        Add-Member -InputObject $rec -MemberType NoteProperty -Name "InternetAddress" -Value $(Get-DnsAddressList $_)
        $rec
    }
}

function New-ObjectWithAddPropertyScriptMethod
{   
    $record = New-Object -TypeName PSObject;
    Add-Member -Name AddProperty -InputObject $record -MemberType ScriptMethod -Value {
        if ($args.Count -eq 2)
        {
            ($name, $value) = $args;
            if (Get-Member -InputObject $this -Name $name){$this.$name = $value;} 
            else { Add-Member -InputObject $this -MemberType NoteProperty -Name $name -Value $value; }
        } 
    } 
    $record;
} 
function Parse-NslookupTyperecord
{
    param( 
    [parameter(ValueFromPipeline=$true)][string[]]$FQDN,
    [parameter(ValueFromPipeline=$true)][string]$CMDString  
    
    );
    
    begin{ } 

    process
    {
        foreach ($_fqdn in $FQDN)
        {
           $cmd = $CMDString

            Write-Verbose "$($MyInvocation.MyCommand.Name) -FQDN $_fqdn"; 

            $record = $null;
            $data = @()

            $cmd | nslookup.exe 2>&1 |
            ? {
                $_ -and
                $_ -notmatch '^Address:' -and
                $_ -notmatch '^Server:' -and
                $_ -notmatch '^Default Server:' -and
                $_ -notmatch '^>'
            } |
            % {                
                switch -Regex ($_)
                {
                    "^[^\s]"
                    {
                        if ($record){$data += $_;}
                        else
                        {                
                            $record = New-ObjectWithAddPropertyScriptMethod;
                            $record.AddProperty('FQDN', ($_ -replace '\s.*'));
                        } 

                    }
                    "^\s"
                    {
                        if ($_ -match ' = ')
                        {                           
                            $name = $_ -replace '^\s*' -replace '\s.* = .*';
                            $value = $_ -replace '.* = ' -replace '\s*$';
                            $record.AddProperty($name, $value);

                        } 
                        else
                        {
                            Write-Host $_
                            $data += $_;
                        } 
                    } 
                }
            } 

            $record.AddProperty('Text', $data)
            $record
        } 
    }      
} 



Function DNS-LookUp {  
   
   $DNSObject = New-Object -TypeName psobject
   $DNSObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "DNS LookUp"  

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
    
     $DNSObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $DNSObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $DNSObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $DNSObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
     $DNSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
         
     return $DNSObject
 
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
    
     $DNSObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $DNSObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $DNSObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $DNSObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
     $DNSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $DNSObject

   }  
    
    $Error.Clear()
    $QueryDNS = ""
    $QueryResult = ""
    $QueryDNS = ""
    [bool]$IsFromIP = $false
        
    try{

        if([bool]($DNSQuery -as [ipaddress])){
            
          $DNSQuery =  [System.Net.Dns]::GetHostEntry($DNSQuery).HostName
          $IsFromIP = $true
         
        }else{

           $Pattern = '(?=^.{1,254}$)(^(?:(?!\d+\.|-)[a-zA-Z0-9_\-]{1,63}(?<!-)\.?)+(?:[a-zA-Z]{2,})$)' 
           if(-not($DNSQuery -match $Pattern)){
                        
             $DNSObject | Add-Member -MemberType NoteProperty -Name Status -Value 1
             $DNSObject | Add-Member -MemberType NoteProperty -Name Message -Value "Invalid character in DNSQuery or multi record selection"
             $DNSObject | Add-Member -MemberType NoteProperty -Name Message -Value "Expected format: www.microsoft.com or 10.2.19.75"
             
             $StdErrArr = @()
             $stdOutArr = @()
             $StdErr = New-Object PSObject -Property @{		       
		                      id = 0;
                              title =  "Invalid character in DNSQuery or multi record selection"
                              detail = "Expected format: www.microsoft.com or 10.2.19.75"

                   }
         
             $StdErrArr += $StdErr
             $stdOutArr += ("Invalid character in DNSQuery or multi record selection")
    
             $DNSObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $DNSObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $DNSObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
             $DNSObject | Add-Member -MemberType NoteProperty -Name Result -Value "Invalid character in DNSQuery or multi record selection"
             $DNSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
      
             return $DNSObject
           }

        }

        $QueryDNS = [System.Net.Dns]::gethostentry($DNSQuery)
        $QueryResult = $QueryDNS.AddressList | select IPAddressToString
        $Value = $QueryResult.IPAddressToString
       
        if($Error) {
             
             $StdErrArr = @()
             $stdOutArr = @()
             $StdErr = New-Object PSObject -Property @{		       
		                      id = 0;
                              title =  "DNS Host/IP not found"
                              detail = "DNS Host/IP not found"

                   }
         
             $StdErrArr += $StdErr
             $stdOutArr += ("DNS Host/IP not found")
    
             $DNSObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $DNSObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $DNSObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
             $DNSObject | Add-Member -MemberType NoteProperty -Name Result -Value "DNS Host/IP not found"
             $DNSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
      
            
             return $DNSObject

        } else {
            # Write-Host "$Value" -ForegroundColor Green
        }

         $DNSArr = @()
         $DNSArrStr = @()
                
        if($DNS_TYPE_A){
         
          $Object1 = New-Object PSObject -Property @{	
                       DNSType = "A";	       
		               FQDN = $QueryDNS.HostName;            
                       IPv4Address = $QueryDNS.AddressList;
                    }
           $DNSType = "A";	       
		   $FQDN = $QueryDNS.HostName;            
           $IPv4Address = $QueryDNS.AddressList;
          
          $DNSArr += $Object1
          $DNSArrStr += "DNSType : $DNSType, FQDN : $FQDN, IPv4Address = $IPv4Address" 

       }

       if($DNS_TYPE_AAAA){
            
         $ns = $null
                        
         if($ns = (nslookup -type=aaaa $DNSQuery 2>$null)){
            # Write-Host "Non-authoritative answer"
          }else{
           $ns = (nslookup -type=aaaa $DNSQuery)
          }
              
         $AAAAlookup = [PSCustomObject]@{
           FQDN = ($ns[3] -split ‘:’)[1]
           IPv6 = ($ns[4] -split ‘:  ’)[1]
         }
        
         $Object1 = New-Object PSObject -Property @{
                       DNSType = "AAAA";		       
		               FQDN = $QueryDNS.HostName;            
                       IPv6Address = $AAAAlookup.IPV6;
                    }
         
                  
          $DNSType = "AAAA";	       
		  $FQDN = $QueryDNS.HostName;            
          $IPv6Address = $QueryDNS.AddressList;
          
          $DNSArr += $Object1
          $DNSArrStr += "DNSType : $DNSType, FQDN : $FQDN, IPv6Address = $IPv6Address" 
                  
      }

      if($DNS_TYPE_Cname){
                     
          $Cname = Parse-NslookupTyperecord -FQDN $QueryDNS.HostName -CMDString "set type=Cname`n$DNSQuery"

          $Object1 = New-Object PSObject -Property @{
                       DNSType = "CName";		       
		               Cname = $Cname;
                    }

          $DNSType = "CName";	       
		  $Cname = $Cname;
          
          $DNSArr += $Object1
          $DNSArrStr += "DNSType : $DNSType, Cname : $Cname" 
       }

       if($DNS_TYPE_SRV){
                    
          $SRV =  Parse-NslookupTyperecord -FQDN $QueryDNS.HostName -CMDString "set type=srv`n$DNSQuery"
          $Object1 = New-Object PSObject -Property @{
                       DNSType = "SRV";		       
		               SRV = $SRV;
                    }
          $DNSType = "SRV";	       
		  $SRV = $SRV;
          
          $DNSArr += $Object1
          $DNSArrStr += "DNSType : $DNSType, SRV : $SRV"
        
       }

       if($DNS_TYPE_MX){
           
           if($IsFromIP){

               $DomainName = $QueryDNS.HostName.ToString()
               $DomainNameArr = $DomainName.Split('.')
               $DomainName1 = $DomainNameArr[1]+'.'+$DomainNameArr[2]

              $MX =  Get-DnsMXQuery -DomainName $DomainName1
                            
              foreach($rec in $MX){
                 
                 $Object1 = New-Object PSObject -Property @{
                        DNSType = "Mx";		       
		                InternetAddress = $rec.InternetAddress;
                        MailExchange = $rec.MailExchange;
                    }

                 $DNSType = "Mx";		       
		         $InternetAddress = $rec.InternetAddress;
                 $MailExchange = $rec.MailExchange;

                 $DNSArr += $Object1
                 $DNSArrStr += "DNSType : $DNSType, InternetAddress : $InternetAddress, MailExchange : $MailExchange "
               
              }

           }else{

              $MX = Get-DnsMXQuery -DomainName $QueryDNS.HostName
                            
              foreach($rec in $MX){

                 $Object1 = New-Object PSObject -Property @{
                        DNSType = "Mx";		       
		                InternetAddress = $rec.InternetAddress;
                        MailExchange = $rec.MailExchange;
                    }
                 
                 $DNSType = "Mx";		       
		         $InternetAddress = $rec.InternetAddress;
                 $MailExchange = $rec.MailExchange;

                 $DNSArr += $Object1
                 $DNSArrStr += "DNSType : $DNSType, InternetAddress : $InternetAddress, MailExchange : $MailExchange "
                
              }
           }
       
       }
       if($DNS_TYPE_NS){
           
          $NS = nslookup -type=NS $DNSQuery 2>$null
                    
          $Object1 = New-Object PSObject -Property @{
                       DNSType = "NS";		       		       
		               NS = $NS;
                    }
          
          $DNSType = "NS";		       
		  $NS = $NS;
         
          $DNSArr += $Object1
          $DNSArrStr += "DNSType : $DNSType, NS : $NS"
          
         }

       if($DNS_TYPE_PTR){
          
          $PTR = Parse-NslookupTyperecord -FQDN $QueryDNS.HostName -CMDString "set type=PTR`n$DNSQuery"
          $Object1 = New-Object PSObject -Property @{
                       DNSType = "PTR";		       
		               PTR = $PTR;
                    }

         $DNSType = "PTR";		       
		 $PTR = $PTR;
         
         $DNSArr += $Object1
         $DNSArrStr += "DNSType : $DNSType, PTR : $PTR"
      
      }

      if($DNS_TYPE_TXT){
           #-----------------'TXT Type Records -------------------
          
          $TXT =Parse-NslookupTyperecord -FQDN $QueryDNS.HostName -CMDString "set type=txt`n$DNSQuery"
          $Object1 = New-Object PSObject -Property @{	
                       DNSType = "TXT";	       
		               TXT = $TXT;
                    }

         $DNSType = "TXT";		       
		 $TXT = $TXT;
         
         $DNSArr += $Object1
         $DNSArrStr += "DNSType : $DNSType, TXT : $TXT"
         
       }
     
     $DNSObj =  New-Object psobject
     $DNSObj | Add-Member -MemberType NoteProperty -Name DNSType -Value $DNSArr
              
     $DNSObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
     $DNSObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
     $DNSObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : The Query result for $($Search)"
     $DNSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $DNSArrStr
     $DNSObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $DNSObj  
     
     return $DNSObject

    }catch{
        
         [string]$error = $_.Exception.Message.ToString()
         [string]$errorStr = ""        
         if($error.Contains("No such host is known")){
             $errorStr = “No information found for the requested record(s)”
        
         }else{
           $errorStr = "$($error.Split(':')[1])"
        }
        

      $StdErrArr = @()
      $stdOutArr = @()
      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Exception occured"
                          detail = $errorStr

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ($errorStr)
    
     $DNSObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $DNSObject | Add-Member -MemberType NoteProperty -Name Code -Value 2
     $DNSObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $DNSObject | Add-Member -MemberType NoteProperty -Name Result -Value "Exception occured"
     $DNSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr 
     return $DNSObject

    }


}


if($PSVersionTable.PSVersion.Major -eq 2){

    DNS-LookUp |  ConvertTo-JSONP2

}else{

    DNS-LookUp | ConvertTo-Json -Depth 10
}



