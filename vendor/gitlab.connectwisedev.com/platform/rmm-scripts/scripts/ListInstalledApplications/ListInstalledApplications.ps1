clear-Host

<#

cat - Application

[string]$ApplicationName = "123Realtek High Definition Audio Driver123" -optional

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


Function Get-InstalledSoftware  {

  [OutputType('System.Software.Inventory')]
  [Cmdletbinding()]
  Param( 
      [Parameter(ValueFromPipeline=$True,ValueFromPipelineByPropertyName=$True)] 
      [String[]]$Computername=$env:COMPUTERNAME
  )         
    
  Process  { 
  
   $InsObject = New-Object -TypeName psobject
   $InsObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "Get Installed Software"  
  
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
    
     $InsObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $InsObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $InsObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $InsObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
     $InsObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
    
     return $InsObject
 
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
    
     $InsObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $InsObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $InsObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $InsObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
     $InsObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

    
     return $InsObject

   }
       
  
    $Paths  = @("SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",
               "SOFTWARE\\Wow6432node\\Microsoft\\Windows\\CurrentVersion\\Uninstall")         
    
    $InstalledSoftware = @()
    $InstalledSoftwareStr = @()

   ForEach($Path in $Paths) {
      Try  {

            $reg=[microsoft.win32.registrykey]::OpenRemoteBaseKey('LocalMachine',$Computername) 

      } Catch  { 
        Write-Error $_ 
        Continue
      } 
  Try  {

    $regkey=$reg.OpenSubKey($Path)
    $subkeys=$regkey.GetSubKeyNames()

    ForEach ($key in $subkeys){

      $thisKey=$Path+"\\"+$key 
      Try {  

          $thisSubKey=$reg.OpenSubKey($thisKey)
          $DisplayName =  $thisSubKey.getValue("DisplayName")
          If ($DisplayName  -AND $DisplayName  -notmatch '^Update  for|rollup|^Security Update|^Service Pack|^HotFix|Update for') {
          
          $Date = $thisSubKey.GetValue('InstallDate') 
                  
          if($Date -like "*/*"){

            $DatePartArr = @()
            $DatePartArr = $date.Split('/')
            $Date = $DatePartArr[2]+$DatePartArr[0]+$DatePartArr[1]
            
          } 
                 
          If ($Date) {
              Try {
                $Date = [datetime]::ParseExact($Date, 'yyyyMMdd', $Null)

              } Catch{
                  Write-Warning "$($Computer): $_ <$($Date)>"
                  $Date = $Null
              }
          } 


          $Publisher =  Try {
            $thisSubKey.GetValue('Publisher').Trim()
          }
          Catch {
            $thisSubKey.GetValue('Publisher')
          }

          $Version = Try {
            $thisSubKey.GetValue('DisplayVersion').TrimEnd(([char[]](32,0)))
          }
          Catch {
            $thisSubKey.GetValue('DisplayVersion')
          }

          $UninstallString =  Try {
            $thisSubKey.GetValue('UninstallString').Trim()
          }
          Catch {
            $thisSubKey.GetValue('UninstallString')
          }

          $InstallLocation =  Try {
            $thisSubKey.GetValue('InstallLocation').Trim()
          }
          Catch {
             $thisSubKey.GetValue('InstallLocation')
          }

          $InstallSource =  Try {
            $thisSubKey.GetValue('InstallSource').Trim()
          }
          Catch {
             $thisSubKey.GetValue('InstallSource')
          }

          $HelpLink = Try {
            $thisSubKey.GetValue('HelpLink').Trim()
          } 

          Catch {
            $thisSubKey.GetValue('HelpLink')
          }

          $Object = [pscustomobject]@{
                                
              Name = $DisplayName
              Publisher = $Publisher
              Version  = $Version
              InstallDate = [string]$Date
          } 

          $Object.pstypenames.insert(0,'System.Software.Inventory')
          $InstalledSoftware += $Object
          $DateStr = $Object.InstallDate
          $InstalledSoftwareStr += "Name : $DisplayName, Publisher : $Publisher, Version : $Version, InstallDate : $DateStr"
          
        }

      } Catch {
        Write-Warning "$Key : $_"
      }   

  }

  } Catch  {}  
       $reg.Close()
  } # end foreach path

  # Case 1 : if application name is empty , display all installed application
  IF ([string]::IsNullOrEmpty($ApplicationName)){
      
       $InsSoftsObj =  New-Object psobject
       $InsSoftsObj | Add-Member -MemberType NoteProperty -Name InstalledSoftware -Value $InstalledSoftware 
              
       $InsObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
       $InsObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
       $InsObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : Installed software retrived"

       $InsObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $InstalledSoftwareStr
       $InsObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $InsSoftsObj  
     
     return $InsObject
   }

  else{
      # Case 2 : if application name is not empty and provided application name is found. 
      $SortedApp = @()
      
      $InstalledSoftware | Where-Object {$_.Name -eq $ApplicationName} | ForEach-Object{$SortedApp += $_}
       
      $SingleObjStrArr = @()
      $SingleObjArr = @()
      if($SortedApp.count -gt 0){
        
        foreach($soft in $SortedApp){

             $Object1 = New-Object PSObject -Property @{		       
		                   Name = $soft.Name;            
                           Publisher = $soft.Publisher;
                           Version  = $soft.Version;   
                           InstallDate = $soft.InstallDate;   
               }

            $Name = $soft.Name;            
            $Publisher = $soft.Publisher;
            $Version  = $soft.Version;   
            $InstallDate = $soft.InstallDate; 

            $SingleObjArr += $Object1 
            $SingleObjStrArr += "Name : $Name, Publisher : $Publisher,  Version : $Version, InstallDate : $InstallDate"
                
          }


         $InstalledSoftware2 =  New-Object psobject
         $InstalledSoftware2 | Add-Member -MemberType NoteProperty -Name InstalledSoftware -Value $SingleObjArr
              
         $InsObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
         $InsObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
         $InsObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : Application name : $($ApplicationName) found"
         $InsObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $SingleObjStrArr
         $InsObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $InstalledSoftware2  
         return $InsObject

       
      }else{ 
      
         # Case 3 : if application name is not empty and provided application name is not found. 
         # display all application started with first letter   
       
         $SortedApp = @()
         $ApplicationName1 = '*'+$ApplicationName+'*'
         $AppFirstLetter = $ApplicationName.ToCharArray()[0]      
         $AppFirstLetter1 = $AppFirstLetter+'*'
         
         $InstalledSoftware | Where-Object {$_.Name -like $AppFirstLetter1} | ForEach-Object{$SortedApp += $_}
       
         $FLObjStrArr = @()
         $FLObjArr = @()
         $FLObjErrArr = @()
         if($SortedApp.count -gt 0){
        
           foreach($soft in $SortedApp){

             $Object1 = New-Object PSObject -Property @{		       
		                   Name = $soft.Name;            
                           Publisher = $soft.Publisher;
                           Version  = $soft.Version;   
                           InstallDate = $soft.InstallDate;   
               }

            $Name = $soft.Name;            
            $Publisher = $soft.Publisher;
            $Version  = $soft.Version;   
            $InstallDate = $soft.InstallDate; 

            $FLObjArr += $Object1 
            $FLObjStrArr += "Name : $Name, Publisher : $Publisher,  Version : $Version, InstallDate : $InstallDate"
                
          }


          $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = "Application name : $($ApplicationName) not found"
                          detail = "Displaying all the application(s) starting with first letter $($AppFirstLetter)"

               }

         $FLObjErrArr += $StdErr
         $InstalledSoftwareStr =  New-Object psobject
         $InstalledSoftwareStr | Add-Member -MemberType NoteProperty -Name InstalledSoftware -Value $FLObjArr
              
         $InsObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
         $InsObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
         $InsObject | Add-Member -MemberType NoteProperty -Name stderr -Value $FLObjErrArr
         $InsObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error : Application name : $($ApplicationName) not found"
        
         $InsObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $FLObjStrArr
         $InsObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $InstalledSoftwareStr  
         return $InsObject

      }else{

          # Case 4 : if application name is not empty and provided application name is not found. 
          # Applications started with first letter not found 
         
          $ApplicationName1 = '*'+$ApplicationName+'*'
          $AppFirstLetter = $ApplicationName.ToCharArray()[0]      
          $AppFirstLetter1 = $AppFirstLetter+'*'
          $FLObjErrArr = @()
          
          $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = "Application name : $($ApplicationName) not found"
                          detail = "Application(s) starting with first letter $($AppFirstLetter) not found , displaying all installed applications"

               }

          $FLObjErrArr += $StdErr

          $InsSoftsObj =  New-Object psobject
          $InsSoftsObj | Add-Member -MemberType NoteProperty -Name InstalledSoftware -Value $InstalledSoftware 
              
          $InsObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
          $InsObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
          $InsObject | Add-Member -MemberType NoteProperty -Name stderr -Value $FLObjErrArr
          $InsObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error : Application name : $($ApplicationName) not found"
        

          $InsObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $InstalledSoftwareStr
          $InsObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $InsSoftsObj  
     
          return $InsObject

      }
     } 
      
   } 
   
  } 

}  


if($PSVersionTable.PSVersion.Major -eq 2){

    Get-InstalledSoftware | ConvertTo-JSONP2

}else{

    Get-InstalledSoftware | ConvertTo-Json -Depth 10

}
