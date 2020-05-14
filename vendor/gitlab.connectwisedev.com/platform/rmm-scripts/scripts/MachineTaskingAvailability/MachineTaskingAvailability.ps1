
Clear-Host

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


Function CheckMachineTaskingAvailability{ 

   $MacTaskingObject = New-Object -TypeName psobject
   $MacTaskingObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "Machine Tasking Availability"


   $hostVerSionMajor = ($PSVersionTable.PSVersion.Major).ToString()
   $hostVerSionMinor = ($PSVersionTable.PSVersion.Minor).ToString()
   $hostVersion = $hostVerSionMajor +'.'+ $hostVerSionMinor 

   $osVersionMajor = ([System.Environment]::OSVersion.Version.major).ToString()
   $osVersionMinor = ([System.Environment]::OSVersion.Version.minor).ToString()
   $osVersion = $osVersionMajor +'.'+ $osVersionMinor

   [boolean]$isPsVersionOk = ([version]$hostVersion -ge [version]'6.0')
   [boolean]$isOSVersionOk = ([version]$osVersion -ge [version]'6.0')

  
   if($isPsVersionOk -and $isOSVersionOk){ 
     
     $StdOutArr = @()    
     $StdOutArr += ("Server : $($env:COMPUTERNAME) is available and OK")

     $MacTaskingObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
     $MacTaskingObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
     $MacTaskingObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : Server : $($env:COMPUTERNAME) is available and OK"
     $MacTaskingObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $StdOutArr
    
     return  $MacTaskingObject
    
   }

   if($isPsVersionOk -and -not $isOSVersionOk){
        
     # Write-Warning "Server : $($env:COMPUTERNAME) is available and the OS version is not OK"

     $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = "Server : $($env:COMPUTERNAME) is available and the OS version is not OK"
                          detail = "Server : $($env:COMPUTERNAME) is available and the OS version is not OK"

               }

    
      $StdOutArr = @()

      $StdOutArr += "Server : $($env:COMPUTERNAME) is available and the OS version is not OK"
                  
                    
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErr
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error : Server : $($env:COMPUTERNAME) is available and the OS version is not OK"
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name stdout -Value $StdOutArr
           
      return $MacTaskingObject
    
   }

   if(-not $isPsVersionOk -and $isOSVersionOk){
        
   #  Write-Warning "Server : $($env:COMPUTERNAME) is available and the PS version is not OK"
     
     $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = "Server : $($env:COMPUTERNAME) is available and the PS version is not OK"
                          detail = "Server : $($env:COMPUTERNAME) is available and the PS version is not OK"

               }

    
      $StdOutArr = @()

      $StdOutArr += "Server : $($env:COMPUTERNAME) is available and the PS version is not OK"
                  
                    
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErr
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error : Server : $($env:COMPUTERNAME) is available and the PS version is not OK"
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name stdout -Value $StdOutArr
           
      return $MacTaskingObject
     
   }

   if(-not $isPsVersionOk -and -not $isOSVersionOk){
        
     # Write-Warning "Server : $($env:COMPUTERNAME) is available and the OS and PS version is not OK" 
     
      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = "Server : $($env:COMPUTERNAME) is available and the OS and PS version is not OK" 
                          detail = "Server : $($env:COMPUTERNAME) is available and the OS and PS version is not OK" 

               }

     
      $StdOutArr = @()

      $StdOutArr += "Server : $($env:COMPUTERNAME) is available and the OS and PS version is not OK" 
                  
                    
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErr
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error : Server : $($env:COMPUTERNAME) is available and the OS and PS version is not OK" 
      $MacTaskingObject | Add-Member -MemberType NoteProperty -Name stdout -Value $StdOutArr
           
      return $MacTaskingObject
     

   }

   
}
 
 if($PSVersionTable.PSVersion.Major -eq 2){

    CheckMachineTaskingAvailability |  ConvertTo-JSONP2

}else{

    CheckMachineTaskingAvailability |  ConvertTo-Json -Depth 10
} 


