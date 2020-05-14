Clear-Host

<#
cat - Anti-Malware
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


Function CheckAntiVirusStatus([string] $ServiceName)
{
   $AVStatus = Get-WmiObject -Class Win32_Service | Where-Object {$_.name -eq $ServiceName }
   if($AVStatus.Status -ne 'OK' -or $AVStatus.State -ne 'Running' -or $AVStatus.StartMode -eq 'Disabled'){return 1}     
   return 0 
}


function GetAllInstalledAntivirus(){

   $AVObject = New-Object -TypeName psobject
   $AVObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "Check for all AV Solutions"

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
    
     $AVObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $AVObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $AVObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $AVObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
     $AVObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
    
     return $AVObject 
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
     
     return $AVObject
   }

    # $array -match 'DEF'

    $filters = "Sentinel Agent","Webroot SecureAnywhere","Windows Defender"

    $AntiVirusList = @()

    $hive = [Microsoft.Win32.RegistryKey]::OpenRemoteBaseKey([Microsoft.Win32.RegistryHive]::LocalMachine, $env:COMPUTERNAME)
    $regPathList = "SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall",
                   "SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall",
                   "SOFTWARE\Microsoft"                 

    foreach($regPath in $regPathList) {
        if($key = $hive.OpenSubKey($regPath)) {
            if($subkeyNames = $key.GetSubKeyNames()) {
               
                foreach($subkeyName in $subkeyNames) {
                  
                  # Window defender 
                   if($subkeyName -eq "Windows Defender")  {                       
                      $productKey = $key.OpenSubKey($subkeyName) 
                      $productName = $subkeyName       
                      $InstallLocation = $productKey.GetValue("InstallLocation") 
                      $ProductStatusValue =  $productKey.GetValue("ProductStatus") 
                      $ProductStatus = "Inactive"
                      if($ProductStatusValue -eq 0){$ProductStatus = "Active"}
                       
                      $antiVirusObj = [PSCustomObject]@{
                       DisplayName = $productName                          
                       InstallLocation = $InstallLocation
                       ProductStatus = $ProductStatus
                     }
                    $AntiVirusList += $antiVirusObj
                   
                   }else{
                                         
                    $productKey = $key.OpenSubKey($subkeyName)   
                    $productName = $productKey.GetValue("DisplayName")
                    $productVersion = $productKey.GetValue("DisplayVersion")
                    $productComments = $productKey.GetValue("Comments")
                    $UninstallString = $productKey.GetValue("UninstallString")
                    $InstallLocation = $productKey.GetValue("InstallLocation") 
                    

                    foreach($filter in $filters){
                        if(($productName -match $filter) -or ($productComments -match $filter)) {
                            
                            $ProductStatusValue = $null
                            $ProductStatus = "Inactive"
                            $ServiceName = $null
                              
                            if($productName -match "Webroot SecureAnywhere"){$ServiceName = 'WRSVC'}
                            if($productName -match "Sentinel Agent"){$ServiceName = 'sedsvc'} 
                            $ProductStatusValue = CheckAntiVirusStatus -ServiceName $ServiceName
                            if($ProductStatusValue -eq 0){$ProductStatus = "Active"}

                            $lastIndex = $UninstallString.LastIndexOf(' ')
                            $UninstallString = $UninstallString.Substring(0,$lastIndex)
                            $UninstallString = $UninstallString -replace '"', ""

                            if(($UninstallString -like '*exe*') -or  ($UninstallString -like '*EXE*')) {
                               
                                 $lastIndex1 = $UninstallString.LastIndexOf('\')
                                 $UninstallString = $UninstallString.Substring(0,$lastIndex1)                                
                            }

                            $antiVirusObj = [PSCustomObject]@{
                               DisplayName = $productName                          
                               InstallLocation = $UninstallString
                               ProductStatus = $ProductStatus
                            }
                            $AntiVirusList += $antiVirusObj
                        }

                    }

                 }


                }
            }
        }
        if($key -ne $null){
        $key.Close()}
    }


    if($AntiVirusList.Count -eq 0){
      
      $StdErrArr = @()
      $stdOutArr = @()
      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Antivirus is not installed in the system"
                          detail = "Antivirus is not installed in the system"

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ("Antivirus is not installed in the system")
    
     $AVObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $AVObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $AVObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $AVObject | Add-Member -MemberType NoteProperty -Name Result -Value "Antivirus is not installed in the system"
     $AVObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr  
     
      return $AVObject
    }

    
    <#
    $AVObject | Add-Member -MemberType NoteProperty -Name Status -Value 0
    $AVObject | Add-Member -MemberType NoteProperty -Name Message -Value "List of antiVirus installed in the system"
    #>

    $AVArr = @()
    $AVArrStr = @()
    $AVArrErr = @()

   
    foreach($AV in $AntiVirusList){
        
        $Object1 = New-Object PSObject -Property @{		       
		               DisplayName = $AV.DisplayName;
                       InstallLocation = $AV.InstallLocation;
                       ProductStatus = $AV.ProductStatus;

                    }
       $DisplayName = $AV.DisplayName;
       $InstallLocation = $AV.InstallLocation;
       $ProductStatus = $AV.ProductStatus;
         
       $AVArr +=  $Object1
       $AVArrStr += "DisplayName : $DisplayName, InstallLocation : $InstallLocation, AutoMaticMatric : $ProductStatus" 

    }    
    
   $AVObj =  New-Object psobject
   $AVObj | Add-Member -MemberType NoteProperty -Name Adapter -Value $AVArr
              
   $AVObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
   $AVObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
   $AVObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : List of antiVirus installed in the system"
   $AVObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $AVArrStr
   $AVObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $AVObj  
       
  return $AVObject
}


if($PSVersionTable.PSVersion.Major -eq 2){

   GetAllInstalledAntivirus | ConvertTo-JSONP2

}else{

    GetAllInstalledAntivirus | ConvertTo-Json -Depth 10
}