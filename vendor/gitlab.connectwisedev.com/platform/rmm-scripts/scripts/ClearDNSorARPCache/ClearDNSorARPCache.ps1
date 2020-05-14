Clear-Host


<#
[bool]$ClearDNS = $true 
[bool]$ClearARP = $true
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



function Clear-DNSOrARPCache{
   
   $DNSARPObj = New-Object -TypeName psobject
   $DNSARPObj | Add-Member -MemberType NoteProperty -Name TaskName -Value "Clear DNS Or ARP Cache"  

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
    
     $DNSARPObj | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $DNSARPObj | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $DNSARPObj | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $DNSARPObj | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
     $DNSARPObj | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
     
     return $DNSARPObj
 
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
    
     $DNSARPObj | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $DNSARPObj | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $DNSARPObj | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $DNSARPObj | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
     $DNSARPObj | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
    
     return $DNSARPObj

   }

    try{


        if(($ClearDNS) -and (-not $ClearARP)){
          
            ipconfig /flushDns | Out-Null           
            sleep 2
           
                        
            $StdErrArr = @()
            $stdOutArr = @()
            
            $stdOutArr += ("DNS Cache Flushed")
    
            $DNSARPObj | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
            $DNSARPObj | Add-Member -MemberType NoteProperty -Name Code -Value 0
            $DNSARPObj | Add-Member -MemberType NoteProperty -Name Result -Value "DNS Cache Flushed"
            $DNSARPObj | Add-Member -MemberType NoteProperty -Name stdout -Value $stdOutArr    
                      
            return $DNSARPObj
        }

        if(($ClearARP) -and (-not $ClearDNS)){
                     
            # Flush the ARP Cache
            netsh interface ip delete arpcache | Out-Null 
            sleep 2
                       
            $StdErrArr = @()
            $stdOutArr = @()
            
            $stdOutArr += ("ARP Cache Flushed")
    
            $DNSARPObj | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
            $DNSARPObj | Add-Member -MemberType NoteProperty -Name Code -Value 0
            $DNSARPObj | Add-Member -MemberType NoteProperty -Name Result -Value "ARP Cache Flushed"
            $DNSARPObj | Add-Member -MemberType NoteProperty -Name stdout -Value $stdOutArr    
           
            return $DNSARPObj
        }

        if( $ClearARP -and $ClearDNS){

            ipconfig /flushDns | Out-Null           
             sleep 2
           
             netsh interface ip delete arpcache | Out-Null 
             sleep 2
                        
            $StdErrArr = @()
            $stdOutArr = @()
            
            $stdOutArr += ("DNS Cache Flushed")
            $stdOutArr += ("ARP Cache Flushed")
    
            $DNSARPObj | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
            $DNSARPObj | Add-Member -MemberType NoteProperty -Name Code -Value 0
            $DNSARPObj | Add-Member -MemberType NoteProperty -Name Result -Value "DNS Cache Flushed"
            $DNSARPObj | Add-Member -MemberType NoteProperty -Name stdout -Value $stdOutArr    
                      
            return $DNSARPObj
        }
        if((-not $ClearARP) -and (-not $ClearDNS)){

              $StdErrArr = @()
              $stdOutArr = @()
              $StdErr = New-Object PSObject -Property @{		       
		                          id = 0;
                                  title =  "No options are selected";
                                  detail = "Select either one of the options(ClearARP,ClearDNS) or both";

                       }
         
             $StdErrArr += $StdErr
             $stdOutArr += ("Select either one of the options(ClearARP,ClearDNS) or both")
    
             $DNSARPObj | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $DNSARPObj | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $DNSARPObj | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
             $DNSARPObj | Add-Member -MemberType NoteProperty -Name Result -Value "Select either one of the options(ClearARP,ClearDNS) or both"
             $DNSARPObj | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
    
             return $DNSARPObj

        }

    }catch{
    
          $StdErrArr = @()
          $stdOutArr = @()
          $StdErr = New-Object PSObject -Property @{		       
		                      id = 0;
                              title =  "Exception occured";
                              detail = "Message: [$($_.Exception.Message)]"

                   }
         
         $StdErrArr += $StdErr
         $stdOutArr += ("Exception occured")
    
         $DNSARPObj | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
         $DNSARPObj | Add-Member -MemberType NoteProperty -Name Code -Value 2
         $DNSARPObj | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
         $DNSARPObj | Add-Member -MemberType NoteProperty -Name Result -Value "Exception occured"
         $DNSARPObj | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
    

          return $DNSARPObj

    }


}
 

 if($PSVersionTable.PSVersion.Major -eq 2){

     Clear-DNSOrARPCache | ConvertTo-JSONP2

}else{

     Clear-DNSOrARPCache | ConvertTo-Json -Depth 10
}
