

clear-host

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


function Get-NetWorkShare {

   $NetWorkShareObject = New-Object -TypeName psobject                    
   $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "List network share"

   $hostVerSionMajor = ($PSVersionTable.PSVersion.Major).ToString()
   $hostVerSionMinor = ($PSVersionTable.PSVersion.Minor).ToString()
   $hostVersion = $hostVerSionMajor +'.'+ $hostVerSionMinor 

   $osVersionMajor = ([System.Environment]::OSVersion.Version.major).ToString()
   $osVersionMinor = ([System.Environment]::OSVersion.Version.minor).ToString()
   $osVersion = $osVersionMajor +'.'+ $osVersionMinor

   [boolean]$isPsVersionOk = ([version]$hostVersion -ge [version]'2.0')
   [boolean]$isOSVersionOk = ([version]$osVersion -ge [version]'6.0')
         
  
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
    
     $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
     $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $NetWorkShareObject

    }

  
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
    
     $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
     $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $NetWorkShareObject

    }
   try{ 

       $GetAllNetShares = get-WmiObject -class Win32_Share -computer $env:COMPUTERNAME -errorAction Stop       
       if($GetAllNetShares){
          
           $StdOutArr = @()
           $StdOutArrStr = @()
          
           foreach($NS in $GetAllNetShares){
               
              $Object = New-Object PSObject -Property @{		       
		                    Name = $NS.Name;
                            Path = $NS.Path;
                            Description = $NS.Description;
                        }

              $Name = $NS.Name;
              $Path = $NS.Path;
              $Description = $NS.Description;
         
              $StdOutArr +=  $Object
              $StdOutArrStr += "Name : $Name, Path : $Path, Description : $Description "

    
          }

          $NetWorkShareObj =  New-Object psobject
          $NetWorkShareObj | Add-Member -MemberType NoteProperty -Name NetworkShare1 -Value $StdOutArr
              
          $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
          $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
          $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : List of network share"        
          $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $StdOutArrStr
          $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $NetWorkShareObj  
      
          return $NetWorkShareObject


       }else{

          $StdErrArr = @()
          $stdOutArr = @()
          $StdErr = New-Object PSObject -Property @{		       
		                      id = 0;
                              title =  "Error Occurred";
                              detail = "Unable to pull share information due to error";

                   }
         
         $StdErrArr += $StdErr
         $stdOutArr += ("Unable to pull share information due to error")
    
         $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
         $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
         $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
         $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Result -Value "Unable to pull share information due to error"
         $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

         return $NetWorkShareObject
       }
   }catch{

          $StdErrArr = @()
          $stdOutArr = @()
          $StdErr = New-Object PSObject -Property @{		       
		                      id = 0;
                              title =  "Error Occurred";
                              detail = "Unable to pull share information due to error";

                   }
         
         $StdErrArr += $StdErr
         $stdOutArr += ("Unable to pull share information due to error")
    
         $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
         $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Code -Value 2
         $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
         $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name Result -Value "Unable to pull share information due to error"
         $NetWorkShareObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

         return $NetWorkShareObject
   }
}


if($PSVersionTable.PSVersion.Major -eq 2){

    Get-NetWorkShare |  ConvertTo-JSONP2

}else{

    Get-NetWorkShare | ConvertTo-Json -Depth 2
}

