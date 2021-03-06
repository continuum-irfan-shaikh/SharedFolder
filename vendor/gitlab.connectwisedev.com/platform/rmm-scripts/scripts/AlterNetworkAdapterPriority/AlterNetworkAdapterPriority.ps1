
Clear-Host

<#

cat - Network

[string]$AdapterName = "Local Area Connection 2" -optional
[int]$InterfaceMetric = 14                       -optional
[boolean]$AutomaticMetric = $true                -optional

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


function Alter-NetworkAdapterPriority{

   $AdptObject = New-Object -TypeName psobject
   $AdptObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "Alter Network Adapter Priority"
    
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
    
     $AdptObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $AdptObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $AdptObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $AdptObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
     $AdptObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $AdptObject
 
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
    
     $AdptObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $AdptObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $AdptObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $AdptObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
     $AdptObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $AdptObject
   }
   
    # case 1 : Adapter name is not provided

    if([string]::IsNullOrEmpty($AdapterName)){   
     
      $AdptArr = @()
      $AdptArrStr = @()
      $AdptArrErr = @()
      
      $Adpts = Get-WmiObject win32_networkadapter | Where-Object {$_.NetConnectionID -ne $null} | select Name,NetConnectionID, NetEnabled 
            
      foreach($Adpt in $Adpts){
               
         $Object1 = New-Object PSObject -Property @{		       
		                Name = $Adpt.NetConnectionID;
                        Enabled = $Adpt.NetEnabled;
                    }

          $Name = $Adpt.NetConnectionID;
          $Enabled = $Adpt.NetEnabled;
         
          $AdptArr +=  $Object1
          $AdptArrStr += "Name : $Name,  Enabled : $Enabled"
    
      }

      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = "Adapter name is not provided"
                          detail = "Displaying all adapters"

               }

      $AdptArrErr += $StdErr
                  
      $AdapterStr =  New-Object psobject
      $AdapterStr | Add-Member -MemberType NoteProperty -Name Adapters -Value $AdptArr
              
      $AdptObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
      $AdptObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
      $AdptObject | Add-Member -MemberType NoteProperty -Name stderr -Value $AdptArrErr
      $AdptObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error : Adapter name is not provided"
        
      $AdptObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $AdptArrStr
      $AdptObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $AdapterStr  
      
      return $AdptObject
    }
    
    $index = (Get-WmiObject win32_networkadapter | ? {$_.NetConnectionID -eq "$AdapterName"}).index
    # Case 2 : Incorrect adapter name is provided
    if($index -lt 0){
         
      $AdptArr = @()
      $AdptArrStr = @()
      $AdptArrErr = @()
      
      $Adpts = Get-WmiObject win32_networkadapter | Where-Object {$_.NetConnectionID -ne $null} | select Name,NetConnectionID, NetEnabled 
            
      foreach($Adpt in $Adpts){
               
         $Object1 = New-Object PSObject -Property @{		       
		                Name = $Adpt.NetConnectionID;
                        Enabled = $Adpt.NetEnabled;
                    }

          $Name = $Adpt.NetConnectionID;
          $Enabled = $Adpt.NetEnabled;
         
          $AdptArr +=  $Object1
          $AdptArrStr += "Name : $Name,  Enabled : $Enabled"
    
      }

      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = "Adapter $($AdapterName) not found"
                          detail = "Displaying all adapters"

               }


      $AdptArrErr += $StdErr
                  
      $AdapterStr =  New-Object psobject
      $AdapterStr | Add-Member -MemberType NoteProperty -Name Adapters -Value $AdptArr
              
      $AdptObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
      $AdptObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
      $AdptObject | Add-Member -MemberType NoteProperty -Name stderr -Value $AdptArrErr
      $AdptObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error : Adapter $($AdapterName) not found"
        
      $AdptObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $AdptArrStr
      $AdptObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $AdapterStr  
      
      return $AdptObject
    
    }
    
    # Assign automatic matrix through script

    if(-not $AutomaticMetric){   
     $Output = netsh interface ipv4 set interface "$AdapterName" metric = $InterfaceMetric 

         if($Output -eq "Ok."){
              
              $AdptArr = @()
              $AdptArrStr = @()
              $AdptArrErr = @()
      
              $AdptName = Get-WmiObject win32_networkadapter | Where-Object {$_.NetConnectionID -eq $AdapterName} | select Name
           
              $IMatrix = get-wmiobject win32_networkadapterConfiguration `
              | Where-Object {$_.Description -eq $AdptName.Name} | select IPConnectionMetric

              $Object1 = New-Object PSObject -Property @{		       
		                Name = $AdapterName;
                        InterfaceMetric = $IMatrix.IPConnectionMetric;
                        AutoMaticMatric = $AutomaticMetric
                }
             
             $Name = $AdapterName;
             $InterfaceMetric = $IMatrix.IPConnectionMetric;
             $AutoMaticMatric = $AutomaticMetric
         
             $AdptArr +=  $Object1
             $AdptArrStr += "Name : $Name,  InterfaceMetric : $InterfaceMetric, AutoMaticMatric : $AutomaticMetric" 
                        
             $AdapterStr =  New-Object psobject
             $AdapterStr | Add-Member -MemberType NoteProperty -Name Adapter -Value $AdptArr
              
             $AdptObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
             $AdptObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
             $AdptObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : Adapter $($AdapterName) Interface matrix changed scuucessfully"  
             $AdptObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $AdptArrStr
             $AdptObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $AdapterStr  
                          
             return $AdptObject

         }else{
                            
              $AdptArrErr = @()
              $AdptArrStr = @()

              $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = $Output;
                          detail = $Output;

               }


              $AdptArrErr += $StdErr
              $AdptArrStr +=  ($Output) 
              
              $AdptObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
              $AdptObject | Add-Member -MemberType NoteProperty -Name Code -Value 2
              $AdptObject | Add-Member -MemberType NoteProperty -Name stderr -Value $AdptArrErr
              $AdptObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error : $Output"
        
              $AdptObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $AdptArrStr
              $AdptObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $AdapterStr  


              return $AdptObject

         }
          
    }
    
    if($AutomaticMetric){
            
         $Output = netsh interface ipv4 set interface "$AdapterName" metric = $null

         if($Output -eq "Ok."){
              
              $AdptArr = @()
              $AdptArrStr = @()
              $AdptArrErr = @()

              $AdptName = Get-WmiObject win32_networkadapter | Where-Object {$_.NetConnectionID -eq $AdapterName} | select Name
           
              $IMatrix = get-wmiobject win32_networkadapterConfiguration `
              | Where-Object {$_.Description -eq $AdptName.Name} | select IPConnectionMetric

              $Object1 = New-Object PSObject -Property @{		       
		                Name = $AdptName.Name;
                        InterfaceMetric = $IMatrix.IPConnectionMetric;
                        AutoMaticMatric = $AutomaticMetric
                }
                        
             $Name = $AdapterName;
             $InterfaceMetric = $IMatrix.IPConnectionMetric;
             $AutoMaticMatric = $AutomaticMetric
         
             $AdptArr +=  $Object1
             $AdptArrStr += "Name : $Name, InterfaceMetric : $InterfaceMetric,  AutoMaticMatric : $AutomaticMetric" 
                        
             $AdapterStr =  New-Object psobject
             $AdapterStr | Add-Member -MemberType NoteProperty -Name Adapter -Value $AdptArr
              
             $AdptObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
             $AdptObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
             $AdptObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : Adapter $($AdapterName) Interface matrix changed by system"  
             $AdptObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $AdptArrStr
             $AdptObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $AdapterStr  
                          
             return $AdptObject
         }else{
        
                            
              $AdptArrErr = @()
              $AdptArrStr = @()

              $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = $Output;
                          detail = $Output;

               }


              $AdptArrErr += $StdErr
              $AdptArrStr +=  ($Output) 
              
              $AdptObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
              $AdptObject | Add-Member -MemberType NoteProperty -Name Code -Value 2
              $AdptObject | Add-Member -MemberType NoteProperty -Name stderr -Value $AdptArrErr
              $AdptObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error : $Output"
        
              $AdptObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $AdptArrStr
              $AdptObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $AdapterStr  


              return $AdptObject

         }

     
    }
    
    
}


if($PSVersionTable.PSVersion.Major -eq 2){

    Alter-NetworkAdapterPriority |  ConvertTo-JSONP2

}else{

    Alter-NetworkAdapterPriority |  ConvertTo-Json -Depth 10
}