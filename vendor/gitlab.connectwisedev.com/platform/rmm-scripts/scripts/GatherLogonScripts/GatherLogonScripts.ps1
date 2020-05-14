
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


function Gather-LoginScripts{

   $LogOnObject = New-Object -TypeName psobject
   $LogOnObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "Gather Login Scripts"
    
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
    
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $LogOnObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
     $LogOnObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $LogOnObject
 
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
    
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $LogOnObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
     $LogOnObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $LogOnObject
   }


   [boolean]$IsDomainContaoller = $false
   $IsDomainContaoller = (Get-WmiObject -Class Win32_ComputerSystem).PartOfDomain

   #---------- Is machine part of domain ------------------
   if( -not $IsDomainContaoller){

      $StdErrArr = @()
      $stdOutArr = @()
      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Server is not a domain controller";
                          detail = "Error, script must be run on a Domain Controller";

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ("Error, script must be run on a Domain Controller")
    
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $LogOnObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error, script must be run on a Domain Controller"
     $LogOnObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
 

     return $LogOnObject

   }
   
try{

  $LogonScriptsArr = New-Object System.Collections.ArrayList
  $DomainName = (Get-WmiObject -Class Win32_ComputerSystem).Domain
  $LogonScriptPath = [IO.Path]::Combine('C:\Windows\SYSVOL\sysvol\', $DomainName, 'scripts')
 
  $LogonScripts = Get-ChildItem $LogonScriptPath | select name -ErrorAction Stop

  foreach($s in $LogonScripts){

   $null = $LogonScriptsArr.Add($s)

  }

  if($LogonScriptsArr.Count -eq 0){

      $StdErrArr = @()
      $stdOutArr = @()
      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Logon are Scripts not found";
                          detail = "No logon scripts available within $DomainName";

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ("No logon scripts available within $DomainName")
    
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $LogOnObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Result -Value "No login scripts available within $DomainName"
     $LogOnObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $LogOnObject

  }

  if($LogonScriptsArr.Count -gt 0){

      $StdOutArr = @()
      $stdStrArr = @()

      foreach($script in $LogonScriptsArr){

          $Object = New-Object PSObject -Property @{		       
		                Name = $script;
                    }

          $Name = $script.Name;
         
          $StdOutArr +=  $Object.Name ;
          $stdStrArr += "Name : $Name" ;
      }

      $LogOnObj =  New-Object psobject
      $LogOnObj | Add-Member -MemberType NoteProperty -Name LogonScripts -Value $StdOutArr
              
      $LogOnObject | Add-Member -MemberType NoteProperty -Name Status -Value "success" 
      $LogOnObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
      $LogOnObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : List of logon scripts"
        
      $LogOnObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdStrArr
      $LogOnObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $LogOnObj  
      
      return $LogOnObject
  }

  }Catch{

      $StdErrArr = @()
      $stdOutArr = @()
      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Unable to obtain login scripts due to error";
                          detail = $_.Exception.Message;

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ($_.Exception.Message)
    
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Code -Value 2
     $LogOnObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $LogOnObject | Add-Member -MemberType NoteProperty -Name Result -Value $_.Exception.Message
     $LogOnObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $LogOnObject
  


  }

}


if($PSVersionTable.PSVersion.Major -eq 2){

    Gather-LoginScripts |  ConvertTo-JSONP2

}else{

    Gather-LoginScripts |  ConvertTo-Json -Depth 10
}
