Clear-Host


#[string]$DriveToScan = $env:SystemDrive 
#[string]$DriveToScan = "C:"


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


function InitiateSentinelOneScan(){
   
   $SOObject = New-Object -TypeName psobject
   $SOObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "Initiate Sentinel OneScan"

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
    
     $SOObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $SOObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $SOObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $SOObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
     $SOObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $SOObject
 
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
    
     $SOObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $SOObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $SOObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $SOObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
     $SOObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $SOObject

   }

   
  if ([string]::IsNullOrEmpty($DriveToScan)){
  
      $StdErrArr = @()
      $stdOutArr = @()
      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Please enter drive name"
                          detail = "Please enter drive name"

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ("Please enter drive name")
    
     $SOObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $SOObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $SOObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $SOObject | Add-Member -MemberType NoteProperty -Name Result -Value "Please enter drive name"
     $SOObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

  
       return $SOObject
    }

 $Drive = Get-WmiObject -Class Win32_logicaldisk | Where-Object {$_.DeviceID -eq $DriveToScan}

  if(-not $Drive){
   
     $StdErrArr = @()
     $stdOutArr = @()
     $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Drive not found in the system" 
                          detail = "Drive : $($DriveToScan) not found in the system" 

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ("Drive : $($DriveToScan) not found in the system")
    
     $SOObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $SOObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $SOObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $SOObject | Add-Member -MemberType NoteProperty -Name Result -Value "Drive : $($DriveToScan) not found in the system" 
     $SOObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr 

    return $SOObject
    
 }
   
 $filter = "Sentinel Agent"
 $resultObj = $null
 $hive = [Microsoft.Win32.RegistryKey]::OpenRemoteBaseKey([Microsoft.Win32.RegistryHive]::LocalMachine,$env:COMPUTERNAME)
 $regPathList = "SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall",
                 "SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall"

  foreach($regPath in $regPathList) {
    if($key = $hive.OpenSubKey($regPath)) {
        if($subkeyNames = $key.GetSubKeyNames()) {
             foreach($subkeyName in $subkeyNames) {
                $productKey = $key.OpenSubKey($subkeyName)
                $productName = $productKey.GetValue("DisplayName")
                $productVersion = $productKey.GetValue("DisplayVersion")
                $productComments = $productKey.GetValue("Comments")
                $UninstallString = $productKey.GetValue("UninstallString")
                $InstallLocation = $productKey.GetValue("InstallLocation")
                if(($productName -match $filter) -or ($productComments -match $filter)) {
                   $resultObj = [PSCustomObject]@{
                      
                      Product = $productName
                      Version = $productVersion                      
                      UninstallString = $UninstallString
                      InstallLocation = $InstallLocation 
                    }                
                  }
              }
          }
      }
      $key.Close()
  }

 if($resultObj){
   
    $ScanString = $null
    $UninstallString = $resultObj.UninstallString
    $lastIndex = $UninstallString.LastIndexOf(' ')
    $UninstallString = $UninstallString.Substring(0,$lastIndex)
    $UninstallString = $UninstallString -replace '"', ""

    if(($UninstallString -like '*exe*') -or  ($UninstallString -like '*EXE*')) {
                               
        $lastIndex1 = $UninstallString.LastIndexOf('\')
        $UninstallString = $UninstallString.Substring(0,$lastIndex1) 
                                
     }
        
    $ExeName = "SentinelCtl.exe"
    $UninstallString = $UninstallString 
    $ScanString = join-path -path $($UninstallString) -childpath $($ExeName) 

    $ScanString = $ScanString
    $command = @'
    cmd.exe /C $($ScanString) is_scan_in_progress
'@
   $IsScanInProgressText = Invoke-Expression -Command:$command 
   $IsScanInProgress = $IsScanInProgressText.Split(':')[1]

   if($IsScanInProgress.trim() -eq "False"){
    
    $SOObject | Add-Member -MemberType NoteProperty -Name Message1 -Value "Scan are running and the SOC will take further actions as necessary"


    $command = @'
    cmd.exe /C $($ScanString) scan_folder -i $($DriveToScan)
'@
    Invoke-Expression -Command:$command 
           
    
     $stdOutArr = @()
             
     $StdErrArr += $StdErr
     $stdOutArr += ("Scan started for drive $($DriveToScan)")
    
     $SOObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
     $SOObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
     $SOObject | Add-Member -MemberType NoteProperty -Name Result -Value "Scan started for drive $($DriveToScan)"
     $SOObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr 

     return  $SOObject

   }else{

     <#$SOObject | Add-Member -MemberType NoteProperty -Name Message -Value "Scan is in progress and there is no need to re-initiate the scan"
     $SOObject | Add-Member -MemberType NoteProperty -Name Message1 -Value "SOC will take further actions as necessary"
     
     #>


     $StdErrArr = @()
     $stdOutArr = @()
     $StdErr1 = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Scan is in progress and there is no need to re-initiate the scan"
                          detail = "Scan is in progress and there is no need to re-initiate the scan"

               }
     $StdErr2 = New-Object PSObject -Property @{		       
		                  id = 1;
                          title =  "SOC will take further actions as necessary"
                          detail = "SOC will take further actions as necessary"

               } 
         
     $StdErrArr += $StdErr1
     $StdErrArr += $StdErr2

     $stdOutArr += ("Scan is in progress and there is no need to re-initiate the scan")
    
     $SOObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $SOObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $SOObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $SOObject | Add-Member -MemberType NoteProperty -Name Result -Value $stdOutArr += "Scan is in progress and there is no need to re-initiate the scan"
     $SOObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr 
     
     return $SOObject 

   }
    

 }else{
   
   
     $StdErrArr = @()
     $stdOutArr = @()
     $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Sentinel is not installed on the endpoint"
                          detail = "Sentinel is not installed on the endpoint"

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ("Sentinel is not installed on the endpoint")
    
     $SOObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $SOObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $SOObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $SOObject | Add-Member -MemberType NoteProperty -Name Result -Value "Sentinel is not installed on the endpoint"
     $SOObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr 


   return $SOObject
 }
}



if($PSVersionTable.PSVersion.Major -eq 2){

    InitiateSentinelOneScan | ConvertTo-JSONP2

}else{

    InitiateSentinelOneScan | ConvertTo-Json -Depth 10
}
