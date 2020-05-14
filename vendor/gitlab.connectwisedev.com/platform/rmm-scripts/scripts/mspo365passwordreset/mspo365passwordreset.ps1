Clear-Host

<#
[string]$UserID = "janet@InfiniteLocal.onmicrosoft.com"
[string]$NewPassword = "jan#12346"

$AdminID = "Helpdesk@InfiniteLocal.onmicrosoft.com"
$AdminPassword = "Honeyjain1982"
#>

$AdminPassword = $AdminPassword | ConvertTo-SecureString -asPlainText -Force 
$AdminCredential = New-Object System.Management.Automation.PSCredential($AdminID,$AdminPassword) 

$timeoutSeconds = 600


$SuccessObj = New-Object PSObject

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

Function Check-PreCondition{   
   
    #####################################    
    # Check if it has the correct PS Version
    #####################################

    if(-NOT($PSVersionTable.PSVersion.Major -ge 2)){   
           
       return "Error:Powershell version below 2.0 is not supported, Please upadet powershell"

    }else{ }
    
    ####################################
    # Verify whether MSonline is installed
    ####################################

    if(-NOT(Get-Module -ListAvailable -Name "Msonline")){

       DownloadPowerShellModule -ModuleName 'MSOnline' $False
      
      }else{ }        
      
    ####################################  
    # Verify whether AzureAD is installed
    ####################################

     if(-NOT(Get-Module -ListAvailable -Name "AzureAD")){

       DownloadPowerShellModule -ModuleName 'AzureAD' $False
      
      }else{  }
    
    ####################################  
    # Verify whether AzurerRM.profile is installed
    ####################################
    
     if(-NOT(Get-Module -ListAvailable -Name "AzureRM.profile")){

       DownloadPowerShellModule -ModuleName 'AzureRM.profile' $False
      
      }else{ }
         
    ##########################################   
    # Check Admin Inputs
    ##########################################  

    if(([string]::IsNullOrEmpty($AdminID)) -and -not([string]::IsNullOrEmpty($AdminPassword)))
    {       
        return "Error:Please provide admin login id" 

    }elseif(-not([string]::IsNullOrEmpty($AdminID)) -and ([string]::IsNullOrEmpty($AdminPassword))){

        return "Error:Please provide admin login password" 

    }elseif(([string]::IsNullOrEmpty($AdminID)) -and ([string]::IsNullOrEmpty($AdminPassword))){

        return "Error:Please provide admin login id and password"
  
    }

        
   ##########################################   
    # If MFA {Multi Factor Authentication} is enabled
    ##########################################  
    
    try{    
        $null = Connect-AzureRmAccount -Credential $AdminCredential -ErrorAction Stop    
        
     }catch{
          
               
         if(($_.Exception.Message).ToString().Contains("Due to a configuration change made by your administrator")){
            
           return "Error:Bad username or password"
            
         }else{

            return "Error:$($_.Exception.Message)"
         }
        
     } 
    ##################################   
    # Check if O365 Admin credentials work  
    ################################## 

    try
    {   
       $null = Connect-MsolService -Credential $AdminCredential -ErrorAction Stop
        $role = Get-MsolRole -RoleName "Helpdesk Administrator" -ErrorAction Stop

        if(Get-MsolRoleMember -RoleObjectId $role.ObjectId | Where-Object {$_.EmailAddress -eq $AdminID }){
          
        }else{
          return "Error:Admin credentials provided are incorrect and need updated"
        }

    }catch{
                
         return "Error:$($_.Exception.Message)" 
    
    }   
    

    #################################    
    # Check if it can reach O365
    ##################################

    if(-NOT(Get-MsolDomain -ErrorAction SilentlyContinue)){
   
       return "Error:Office 365 service is not reachable"

     }else{ }
   
    ##########################################
    # if user exists 
    ##########################################

    if(([string]::IsNullOrEmpty($UserID)) -and -not([string]::IsNullOrEmpty($NewPassword)))
    {
        return "Error:Please provide user name"

    }elseif(-not([string]::IsNullOrEmpty($UserID)) -and ([string]::IsNullOrEmpty($NewPassword))){
            
        return "Error:Please provide new password" 

    }elseif(([string]::IsNullOrEmpty($UserID)) -and ([string]::IsNullOrEmpty($NewPassword))){
        
        return "Error:Please provide user name and new password"
    }

    if($NewPassword.ToCharArray().Length -lt 8){
       return "Error:User password should be at least of 8 characters" 
    }
   
    <# if(Get-MsolUser | Where-Object {$_.UserPrincipalName -eq $UserID }){
             
        $AllLicenses = Get-MsolUser | Where-Object {$_.UserPrincipalName -match $UserID } | select Licenses | Out-Default

        $SuccessObj | Add-Member Noteproperty -Name AllLicenses -value  $AllLicenses
         
    }else{        
     
      if((Get-MsolUser | Where-Object {$_.UserPrincipalName -like  $($UserID.ToCharArray(0,1)) +'*' }) -eq $null){
           
            return "Error:User not found ,No User found starting with $(($UserID.ToCharArray(0,1))[0].ToString().ToUpper())"          

       }else{
            
             return "Error1:User not found, Users found starting with $(($UserID.ToCharArray(0,1))[0].ToString().ToUpper())"          
             
            <#  $k = Start-Job -ScriptBlock {

                  $user2 = Get-MsolUser | Where-Object {$_.UserPrincipalName -like  $($UserID.ToCharArray(0,1)) +'*' }
                  $user2 | select DisplayName , UserPrincipalName |  Out-Default 
                  
                }

                Wait-Job $k -Timeout $timeoutSeconds | out-null

                if ($k.State -eq "Completed"){ 

                  $user2 = Get-MsolUser | Where-Object {$_.UserPrincipalName -like  $($UserID.ToCharArray(0,1)) +'*' }
                  $user2 | select DisplayName , UserPrincipalName |  Out-Default 
                  $ErrorArr += $user2 
            
                }
                elseif ($k.State -eq "Running"){ 

                    Write-Warning "System timeout"
                    exit
                }
                else{

                    Write-Warning "System timeout"
                    exit
                 }
                Remove-Job -force $K 

                #>
          
     <#  }

      
    
    }#>
    
    ##########################################
    # Get Block Status
    ##########################################

         
     if((Get-MsolUser -UserPrincipalName $UserID | Select Userprincipalname,blockcredential).blockcredential){

        return "Error:Current Status - Blocked"
     }
     else { 
       $SuccessObj | Add-Member Noteproperty -Name CurrentStatus -value "Not Blocked"
      }
    

    ##########################################
    # Get Enable / Disable Status
    ##########################################

    
    try{

       $null = Connect-AzureAD -Credential $AdminCredential -ErrorAction Stop

       $null = $ADUser = Get-AzureADUser | Where-Object {$_.UserPrincipalName -eq $UserID }

    
        if($ADUser.AccountEnabled){ 
    
          $null = $SuccessObj | Add-Member Noteproperty -Name CurrentUserAccountStatus -value "Enabled"
        
         }
        else{
                 
    
           $null = $SuccessObj | Add-Member Noteproperty -Name CurrentUserAccountStatus -value "Disabled"
        
         }  
    
    }catch{
            
         return "Error:[$($_.Exception.Message)"

    }   
        
    return "Success:PreCondition executed successfully"  
   
}



##########################################################
###------Download MSOnline/ AzureAD from Gallery ---------
##########################################################

$NuGetPackageCode = @"
using System.Net;

public class GetNuGetPackage : WebClient
{   
    protected override WebRequest GetWebRequest(System.Uri address)
    {
        WebRequest request = base.GetWebRequest(address);
        if (request != null)
        {
            return request;
        }
        return request;
    }   
}
"@;
# check type if already added
if (-not ([System.Management.Automation.PSTypeName]'GetNuGetPackage').Type)
{
    Add-Type -TypeDefinition $NuGetPackageCode -Language CSharp
}

function IsDotNetFrameWork45OrAbove
{
    param(
        [string]$ComputerName = $env:COMPUTERNAME
    )
    
    $tempObj = New-Object PSObject
    $tempObj | Add-Member Noteproperty -Name tempVersion  -value "0.0"    
    
    $dotNetRegistry  = 'SOFTWARE\Microsoft\NET Framework Setup\NDP'
    
     if($regKey = [Microsoft.Win32.RegistryKey]::OpenRemoteBaseKey('LocalMachine', $ComputerName))
      {
         if ($netRegKey = $regKey.OpenSubKey("$dotNetRegistry"))
          {
             foreach ($versionKeyName in $netRegKey.GetSubKeyNames())
              {
                 if ($versionKeyName -match '^v[123]') {
                 
                     $versionKey = $netRegKey.OpenSubKey($versionKeyName)
                     $version = [System.Version]($versionKey.GetValue('Version', ''))
                     
                     $DNetMajor = $version.Major.ToString()
                     $DNetMinor = $version.Minor.ToString()
                     $DNetVersion = $DNetMajor + '.' + $DNetMinor
                     
                     $tempObj.tempVersion = [System.Version]"0.0"
                     
                     if([Version]$DNetVersion -gt [Version]$tempObj.tempVersion){
                                             
                       $tempObj.tempVersion = $DNetVersion
                       continue
                     }
                  }
               }
               
              if([Version]$tempObj.tempVersion -ge [Version]"4.5"){
                return $true
              }else{ return $false }
           }           
           
      }  
}


function Unzip-20($zipfile, $destination)
{
    $shell = New-Object -ComObject Shell.Application
    $zip = $shell.NameSpace($zipfile)
    foreach($item in $zip.items())
    {        
        $shell.Namespace($destination).copyhere($item)
    }
}


function Unzip-45($zipfile, $destination)
{
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    param([string]$zipfile, [string]$destination)
    [System.IO.Compression.ZipFile]::ExtractToDirectory($zipfile, $destination)
}

Function DownloadPowerShellModule([String] $ModuleName, [bool] $IsVersionFolderExist)
{
   
    $packageObj = New-Object PSObject
    $packageObj | Add-Member Noteproperty -Name PackageAPI -value "https://www.powershellgallery.com/api"
    $packageObj | Add-Member Noteproperty -Name PackageVer -value "v2"
    $packageObj | Add-Member Noteproperty -Name MsModule -value "MSOnline/1.1.183.17"
    $packageObj | Add-Member Noteproperty -Name AdModule -value "AzureAD/2.0.2.4"
    $packageObj | Add-Member Noteproperty -Name AzRMModule -value "azurerm.profile/5.8.3"
    $packageObj | Add-Member Noteproperty -Name DownloadString -value ""

    if($ModuleName -eq "MSOnline"){
            
        $packageObj.DownloadString = $packageObj.PackageAPI +"/" + $packageObj.PackageVer +"/package/"+ $packageObj.MsModule 
        
    }elseif($ModuleName -eq "AzureAD"){       
        
        $packageObj.DownloadString = $packageObj.PackageAPI +"/" + $packageObj.PackageVer +"/package/"+ $packageObj.AdModule 

    }elseif($ModuleName -eq "AzureRM.profile"){       
        
        $packageObj.DownloadString = $packageObj.PackageAPI +"/" + $packageObj.PackageVer +"/package/"+ $packageObj.AzRMModule 
     }

    try{
       
        $webClient = New-Object GetNuGetPackage
        $downloaded = $webClient.downloadString($packageObj.DownloadString) 
        $tempPath = join-path $([System.IO.Path]::GetTempPath()) $($ModuleName +".zip")    
        Add-Content -path $tempPath -value $downloaded
              
        $unZipModulePath = join-path $([System.IO.Path]::GetTempPath()) -ChildPath "nugetPackage"
        
        if(!(Test-Path -Path $unZipModulePath )){
        
           $dir = New-Item -ItemType directory -Path $unZipModulePath
        }
        
        if(IsDotNetFrameWork45OrAbove){
        
            Unzip-45 $tempPath $unZipModulePath
        
        }else{ 
        
            Unzip-20 $tempPath $unZipModulePath
        
        }               

        if(-not($IsVersionFolderExist))
        {
            foreach($path in $env:PSModulePath.Split(';') )
            {                
                if($path.contains("system32") -or $path.contains("Program Files") )
                {                       
                    try
                    {
                        $psModulePath = Join-Path -Path $path -ChildPath $ModuleName
                        if(!(Test-Path -Path $psModulePath )){    
                           $dir = New-Item -ItemType directory -Path $psModulePath                          
                        } 
                                                   
                        $sourcePath = Join-Path -Path $unZipModulePath -ChildPath "\*"
                        Copy-Item -Path $sourcePath -Destination $psModulePath -recurse -Force 
                    }     
                    catch
                    {                       
                        return "Error:$($_.Exception.Message)"
                                          
                    }
                }
            }  
           
         }
         else
         {
            foreach($folder in Get-ChildItem -Path $unZipModulePath -Recurse | ?{ $_.PSIsContainer })
            {        
                if($folder.Name -eq $ModuleName)
                {
                    foreach($path in $env:PSModulePath.Split(';') )
                    { 
                        if($path.contains("system32") -or $path.contains("Program Files") )
                        {
                            try
                            {
                                $psModulePath = Join-Path -Path $path -ChildPath $ModuleName
                                if(!(Test-Path -Path $psModulePath ))
                                {    
                                   $dir =  New-Item -ItemType directory -Path $psModulePath                          
                                } 
                             
                                $sourcePath =  Join-Path -Path $folder.FullName -ChildPath (Get-ChildItem -Path $folder.FullName)                            
                                $sourcePathWithSubFolder = Join-Path -Path $sourcePath -ChildPath "\*"
 
                                Copy-Item -Path $sourcePathWithSubFolder -Destination $psModulePath -recurse -Force
                            }     
                            catch
                            {
                                return "Error:$($_.Exception.Message)"
                                                 
                           }
                        }
                    }                   
            
                }
            
            }
        }
        
       Remove-Item -LiteralPath $tempPath -Force -Recurse
       Remove-Item -LiteralPath $unZipModulePath -Force -Recurse
      
       return "Success:Downloading $ModuleName completed"
                
    }catch{
       
        return "Error:$($_.Exception.Message)"
    }

}

# change user Password
Function PasswordReset {

       
        $PassResetErrObject = New-Object -TypeName psobject
        $PassResetErrObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "Password reset"
        $PassResetErrObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
        $PassResetErrObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
        $PassResetErrObject | Add-Member -MemberType NoteProperty -Name stderr -Value $null
        $PassResetErrObject | Add-Member -MemberType NoteProperty -Name Result -Value $null
        $PassResetErrObject | Add-Member -MemberType NoteProperty -Name stdout -Value $null  
        
        
       # $SuccessArr = @()
        $SuccessOutArr = @()        
        $SuccessResult = ""   

        $OfficeRepairSuccObject = New-Object -TypeName psobject
        $OfficeRepairSuccObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "Password reset"
        $OfficeRepairSuccObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
        $OfficeRepairSuccObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
        $OfficeRepairSuccObject | Add-Member -MemberType NoteProperty -Name Result -Value $null
        $OfficeRepairSuccObject | Add-Member -MemberType NoteProperty -Name stdout -Value $null  
        $OfficeRepairSuccObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $null  



  $Result =  Check-PreCondition

  Write-Host "$Result"

  if($Result.Split(':')[0] -eq 'Error'){
  
     $StdArr = @()
     $StdOutArr = @()
     $StdErrArr = @()
     $Result = ""
   
    $StdErrArr = $Result.Split(':')[1]
    $Result = "Precondition Failed"

    $PassResetErrObject.stderr = $StdErrArr
    $PassResetErrObject.Result = $Result
    $PassResetErrObject.stdout = $StdOutArr
        
    return $PassResetErrObject
  }

  if($Result.Split(':')[0] -eq 'Error1'){
 
    $PassResetErrObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $null  
    
    $StdOutArr = @()
    $StdErrArr = $Result.Split(':')[1]
    $Result = "Precondition Failed"

    $user2 = Get-MsolUser | Where-Object {$_.UserPrincipalName -like  $($UserID.ToCharArray(0,1)) +'*' }
    $user2 | select DisplayName , UserPrincipalName |  Out-Default 


    $PassResetErrObject.stderr = $StdErrArr
    $PassResetErrObject.Result = $Result
    $PassResetErrObject.stdout = $StdOutArr
    $PassResetErrObject.dataObject = $user2
        
    return $PassResetErrObject
  }

  if($Result.Split(':')[0] -eq 'Success'){
  
  try{
    
    $ChangedPass = Set-MsolUserPassword -UserPrincipalName $UserID -NewPassword $NewPassword -ForceChangePasswordOnly $false -ForceChangePassword $false -ErrorAction Stop
    $SuccessObj | Add-Member Noteproperty -Name ChangedPassword -value $ChangedPass
    

    $OfficeRepairSuccObject.Result = "Changed password for:$($UserID)"
    $OfficeRepairSuccObject.stdout =  $SuccessObj
    $OfficeRepairSuccObject.dataObject = $SuccessObj
        
    return $OfficeRepairSuccObject

    }catch{
      
      $StdOutArr += "Message: $($_.Exception.Message)"

      $PassResetErrObject.stderr =  "Message: $($_.Exception.Message)"
      $PassResetErrObject.Result = "Password is not changed.Please try again"
      $PassResetErrObject.stdout = $StdOutArr
        
      return $PassResetErrObject
    }
  }

 }

 if($PSVersionTable.PSVersion.Major -eq 2){

    PasswordReset |  ConvertTo-JSONP2

}else{

    PasswordReset |  ConvertTo-Json -Depth 10
}
