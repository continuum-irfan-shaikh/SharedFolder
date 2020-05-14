
cls

<#

Boolean : IsQuickRepairSelected [True , Default] / IsFullRepairSelected [False] 
[Note:  Check Box if selected, powershell 
script will opt the quick repair option  otherwise script will go for full repair option]

#>

<#

#Dynamic variable
[String]$RepairMSOffice = "Microsoft Office Professional Plus 2019 - en-us"
[bool] $IsQuickRepairSelected = $true

[string]$Admin = "Kumarg@InfiniteLocal.onmicrosoft.com"
[string]$AdminPass = "Honeyjain1982"
[string]$User = "janet3@InfiniteLocal.onmicrosoft.com"

#optional
[bool] $DisplayRepairModal = $false
[bool] $StartAndStopPrintSpooler = $false

#>

#Static variable
[bool] $WaitForInstallToFinish = $true

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


Function RepairMSOfficeXml {

if($DisplayRepairModal){
@"
<Configuration Product="ProPlus">  

  <!--  <Display Level="none" CompletionNotice="no" SuppressModal="yes" AcceptEula="yes" /> -->
    <Display Level="full" CompletionNotice="yes" SuppressModal="no" AcceptEula="no" />
   </Configuration>
"@
}
else{

@"
<Configuration Product="ProPlus">  

    <Display Level="none" CompletionNotice="no" SuppressModal="yes" AcceptEula="yes" /> 
   <!-- <Display Level="full" CompletionNotice="yes" SuppressModal="no" AcceptEula="no" /> -->
   </Configuration>
"@

}

}

function TocheckRunningStatusofWORD 
{

$ErrorActionPreference = 'silentlycontinue'
$DAtaWithout = [Runtime.Interopservices.Marshal]::GetActiveObject('Word.Application')

    if($DAtaWithout -eq $null)
    {
       #Write-Host 'No any running instance found'
       return $true
    }
    else{
           # Write-Host "Saving and Closing opened instance"
           $DAtaWithout.Documents | % { $_.Saved = $true ; $_.Close()} -ErrorVariable ABC
           
           if($ABC.count -eq 0)
           {
             # write-host 'Successfully Saved and closed instance' -ForegroundColor Green
             return $true
           }else{
             # Write-Host "User error : [$($ABC)"] -ForegroundColor Red -BackgroundColor White  
             return $false 
           }
    }

}
Function VerifyAndUpdateModule([String] $ModuleName)
{
    $IsContinued = $false

    Write-Host "-------------------------------"
    Write-Host "Verifying Module : $ModuleName"
    Write-Host "    " 

    # If module is imported say that and do nothing
    if (Get-Module | Where-Object {$_.Name -eq $ModuleName}) 
    {
        write-host -ForegroundColor 8 "`t $ModuleName is already installed" 
        $IsContinued = $true
    } 
    else 
    {
        # If module is not imported, but available on disk then import
        if (Get-Module -ListAvailable | Where-Object {$_.Name -eq $ModuleName}) 
        {
            try
          {
            Import-Module $ModuleName -Verbose -ErrorAction Stop
            write-host -ForegroundColor 8 "`t $ModuleName is already installed" 
            $IsContinued = $true

          }catch{
            
             Write-Host "Message: [$($_.Exception.Message)"] -ForegroundColor Red -BackgroundColor White         
             return $false

          }
        } 
        else 
        {
            # If module is not imported, not available on disk, but is in online gallery then install and import
            if (Find-Module -Name $ModuleName | Where-Object {$_.Name -eq $ModuleName}) 
            {
               try{

                    write-host "Installing: $ModuleName "

                    Install-Module -Name $ModuleName -Force -Verbose -Scope CurrentUser -ErrorAction Stop
                    Import-Module $ModuleName -Verbose -ErrorAction Stop

                    write-host -ForegroundColor 8 "`t $ModuleName installed sucessfully"
                    $IsContinued = $true

                }catch {

                     Write-Host "Message: [$($_.Exception.Message)"] -ForegroundColor Red -BackgroundColor White         
                     return $false

                }
            } 
            else 
            {
                # If module is not imported, not available and not in online gallery then abort
                 Write-Warning "Error while Installing Module $ModuleName : not imported, not available and not in online gallery, exiting."
            }
        }
    }

    return $IsContinued
}

Function Check-O365PreCondition{

    if(([string]::IsNullOrEmpty($Admin)) -and -not([string]::IsNullOrEmpty($AdminPass)))
    {
        return "Error:Please provide admin login id" 

    }elseif(-not([string]::IsNullOrEmpty($Admin)) -and ([string]::IsNullOrEmpty($AdminPass))){

        return "Error:Please provide admin login password" 

    }elseif(([string]::IsNullOrEmpty($Admin)) -and ([string]::IsNullOrEmpty($AdminPass))){

        return "Error:Please provide admin login password" 
  
    }

    $AdminPass = $AdminPass | ConvertTo-SecureString -asPlainText -Force
    $AdminCredential = New-Object System.Management.Automation.PSCredential($Admin,$AdminPass)
        
    ###### Verify whether MSonline is installed ########   

    if(VerifyAndUpdateModule -ModuleName 'MSOnline'){    
    }else{ return "Error:MSOnline download failed " }      
     
    ###### Verify whether AzureAD is installed #########    

    if(VerifyAndUpdateModule -ModuleName 'AzureAD'){    
    }else{ return "Error:AzureAD download failed" }

     ##### Check if domain exist ########

    $AdminDomain = $Admin.Split('@')[1]  
    $UserDomain = $User.Split('@')[1]

    
    ##### Check if O365 Admin credentials work ########
    ##### Check if user exist ########
   
    try{
        Connect-MsolService -Credential $AdminCredential -ErrorAction Stop

        if($?){
    
            if([string]::IsNullOrEmpty($User))
            {              
                return "Error:Please provide user name" 
            } 

                 try{
                        $d = Get-MsolDomain -DomainName $UserDomain -ErrorAction Stop
      
                        if($?){                          
                              
                          }else{

                             return "Error:Domain name $($UserDomain) is not found"
                           }

                    }catch{

                        return "Error:$($_.Exception.Message)"
                    }

                    if(!($AdminDomain -eq $UserDomain)){

                       return "Error:Admin Domain $($AdminDomain) mismatched with user domain $($UserDomain)"

                    }

            
           try{

             Get-MsolUser -UserPrincipalName $User -ErrorAction Stop

               if($?){                    
                   
               }else{
                                  
                    return "Error:User not found"
               }

           }catch{
                           
                return "Error:[$($_.Exception.Message)]" 
           }          
                    
        }
        else {
                    
            return "Error:MS online service is not connected"
        }

    }catch {
 
      return "Error:[$($_.Exception.Message)]"
    }

   return "Success:PreCondition executed successfully"
}

Function Repair-MSOffice {

[CmdletBinding()]
Param(
    [string] $ComputerName = $env:COMPUTERNAME,

    [Parameter()]
    [bool] $WaitForInstallToFinish = $True,

    [Parameter(ValueFromPipelineByPropertyName=$true)]
    [string] $TargetFilePath = $NULL
        
)

     Process{
        
        $StdArr = @()
        $StdOutArr = @()
        $StdErrArr = @()
        $Result = ""

        $OfficeRepairErrObject = New-Object -TypeName psobject
        $OfficeRepairErrObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "MS office repair"
        $OfficeRepairErrObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
        $OfficeRepairErrObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
        $OfficeRepairErrObject | Add-Member -MemberType NoteProperty -Name stderr -Value $null
        $OfficeRepairErrObject | Add-Member -MemberType NoteProperty -Name Result -Value $null
        $OfficeRepairErrObject | Add-Member -MemberType NoteProperty -Name stdout -Value $null  
        
        
        $SuccessArr = @()
        $SuccessOutArr = @()        
        $SuccessResult = ""   

        $OfficeRepairSuccObject = New-Object -TypeName psobject
        $OfficeRepairSuccObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "MS office repair"
        $OfficeRepairSuccObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
        $OfficeRepairSuccObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
        $OfficeRepairSuccObject | Add-Member -MemberType NoteProperty -Name Result -Value $null
        $OfficeRepairSuccObject | Add-Member -MemberType NoteProperty -Name stdout -Value $null  
        $OfficeRepairSuccObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $null  

        $currentFileName = Get-CurrentFileName
        Set-Alias -name LINENUM -value Get-CurrentLineNumber 
           
        [bool] $isInPipe = $true
        if (($PSCmdlet.MyInvocation.PipelineLength -eq 1) -or `
            ($PSCmdlet.MyInvocation.PipelineLength -eq $PSCmdlet.MyInvocation.PipelinePosition)) {

            $isInPipe = $false
        }

        ###########  Verify Powershell Version ##############

        $GetPSVersion = $PSVersionTable.PSVersion 
        $PSMajor = $GetPSVersion.Major.ToString()
        $PSMinor = $GetPSVersion.Minor.ToString()
        $PSVersion = $PSMajor + '.' + $PSMinor

         if([Version]$PSVersion -gt [Version]"2.0" -or [Version]$PSVersion -eq [Version]"2.0"){
           # Write-Host "PS Version : $([Version]$PSVersion)"    
        }else {        
           # Write-Host "Installed PS Version : $([Version]$PSVersion) is not supported"

           $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "PowerShell version below 2.0 is not supported";
                          detail = "PowerShell version below 2.0 is not supported";

               }

           $StdErrArr += $StdErr
           $StdOutArr += ("PowerShell version below 2.0 is not supported")
           $Result = "PowerShell version below 2.0 is not supported"

           $OfficeRepairErrObject.stderr = $StdErrArr
           $OfficeRepairErrObject.Result = $Result          
           $OfficeRepairErrObject.stdout = $StdOutArr

           return $OfficeRepairErrObject
        }

        ###########  Verify Operating System Version ##############

        $GetOSVersion = [Version](Get-ItemProperty -Path "$($Env:Windir)\System32\hal.dll" `
                   -ErrorAction SilentlyContinue).VersionInfo.FileVersion.Split()[0]

        $OSMajor = $GetOSVersion.Major.ToString()
        $OSMinor = $GetOSVersion.Minor.ToString()
        $OSVersion = $OSMajor + '.' + $OSMinor

        if([Version]$OSVersion -gt [Version]"6.0" -or [Version]$OSVersion -eq [Version]"6.0"){         
           
         }else {
           
           $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system";
                          detail = "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system";

               }

           $StdErrArr += $StdErr
           $StdOutArr += ("PowerShell Script supports Window 7, Window 2008R2 and higher version operating system")
           $Result = "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"

           $OfficeRepairErrObject.stderr = $StdErrArr
           $OfficeRepairErrObject.Result = $Result          
           $OfficeRepairErrObject.stdout = $StdOutArr

           return $OfficeRepairErrObject
                      
        }

        ###########  Get Operating System Architecture ##############

        $OSArchitecture = (Get-WmiObject Win32_OperatingSystem).OSArchitecture
        # Write-Host  "OS Architecture : $OSArchitecture"

        $SuccessArr += "OS Architecture : $OSArchitecture"
        $SuccessOutArr += "OS Architecture : $OSArchitecture"

        ###########  Get Office Versions ##############

        $MSOfficeVersion = Get-OfficeVersion  
                       
        if($MSOfficeVersion -eq $NULL){            
             # Write-Host "MS Office is not installed on system"
             # return

           $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "MS Office is not installed on system"
                          detail = "MS Office is not installed on system"

               }

           $StdErrArr += $StdErr
           $StdOutArr += ("MS Office is not installed on system")
           $Result = "MS Office is not installed on system"

           $OfficeRepairErrObject.stderr = $StdErrArr
           $OfficeRepairErrObject.Result = $Result          
           $OfficeRepairErrObject.stdout = $StdOutArr

           return $OfficeRepairErrObject
         } 
         
                
         $OfficeVersion = $MSOfficeVersion.Version.Split('.')
         $OMajor = $OfficeVersion[0]
         $OMinor = $OfficeVersion[1]
         $OVersion = $OMajor + '.' + $OMinor

               
         if([Version]$OVersion -gt [Version]"14.0" -or [Version]$OVersion -eq [Version]"14.0"){  
               
            # Write-Host "MS Office Version : $([Version]$OVersion)" 
            $SuccessArr += "MS Office Version : $([Version]$OVersion)" 
            $SuccessOutArr += "MS Office Version : $([Version]$OVersion)"
               
         }else {

           # Write-Host "Installed MS office Version : $($OfficeVersion) is not supported"
           # return

           $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Installed MS office Version issue"
                          detail = "Installed MS office Version : $($OfficeVersion) is not supported"

               }

           $StdErrArr += $StdErr
           $StdOutArr += ("Installed MS office Version : $($OfficeVersion) is not supported")
           $Result = "Installed MS office Version : $($OfficeVersion) is not supported"

           $OfficeRepairErrObject.stderr = $StdErrArr
           $OfficeRepairErrObject.Result = $Result          
           $OfficeRepairErrObject.stdout = $StdOutArr

           return $OfficeRepairErrObject

         }      

         if($MSOfficeVersion.Count -gt 1) {            
            
            <# Write-Host  "    "
             Write-Host "Multiple versions of MS Office is installed on the system. Uninstall `n all versions of office and installed desired version" 
             Write-Host  "    "
             Write-Host "List of all installed MS Office"
             Write-Host "-------------------------------"
             Write-Host "$($MSOfficeVersion.DisplayName)  $($MSOfficeVersion.Version)"
             return  #>

             $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Multiple versions of MS Office is installed on the system. Uninstall all versions of office and installed desired version" 
                          detail = "List of all installed MS Office"

               }

               $StdErrArr += $StdErr
               $StdOutArr += $MSOfficeVersion.DisplayName
               
               $Result = "Multiple versions of MS Office is installed on the system. Uninstall all versions of office and installed desired version" 

               $OfficeRepairErrObject.stderr = $StdErrArr
               $OfficeRepairErrObject.Result = $Result          
               $OfficeRepairErrObject.stdout = $StdOutArr

               $OfficeRepairErrObj =  New-Object psobject
               $OfficeRepairErrObj | Add-Member -MemberType NoteProperty -Name InstalledOffices -Value $MSOfficeVersion.DisplayName

               $OfficeRepairErrObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $OfficeRepairErrObj

               return $OfficeRepairErrObject

          }
          
          if(!($RepairMSOffice -eq $MSOfficeVersion.DisplayName)){
            
            <#
             
             Write-Host  "    "
             Write-Host "Selected MS office is not found"
             Write-Host "List of all installed MS Office"
             Write-Host "-------------------------------"
             Write-Host "$($MSOfficeVersion.DisplayName)"
             
             #>

             $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Selected MS office is not found"
                          detail = "List of all installed MS Office"

               }

               $StdErrArr += $StdErr
               $StdOutArr += $MSOfficeVersion.DisplayName
               
               $Result = "Selected MS office is not found"

               $OfficeRepairErrObject.stderr = $StdErrArr
               $OfficeRepairErrObject.Result = $Result          
               $OfficeRepairErrObject.stdout = $StdOutArr

               $OfficeRepairErrObj =  New-Object psobject
               $OfficeRepairErrObj | Add-Member -MemberType NoteProperty -Name InstalledOffices -Value $MSOfficeVersion.DisplayName

               $OfficeRepairErrObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $OfficeRepairErrObj

               return $OfficeRepairErrObject

          }

         if($MSOfficeVersion.Count -eq 1 -and $RepairMSOffice -eq "Microsoft Office Professional Plus 2010") {  
            if(-not(TocheckRunningStatusofWORD)){
                return
            }
        }
   
     if($MSOfficeVersion.DisplayName.Contains("365")){
                                    
             $O365PreCheck = Check-O365PreCondition  

             $O365PreCheckArr = $O365PreCheck.Split(':')

             if( $O365PreCheckArr[0] -eq "Error"){

                $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "Error occured" ;
                          detail = $O365PreCheckArr[0] ;

                }

               $StdErrArr += $StdErr
                             
               $Result = "Error occured"

               $OfficeRepairErrObject.stderr = $StdErrArr
               $OfficeRepairErrObject.Result = $Result          
               $OfficeRepairErrObject.stdout = $StdOutArr
                              
               return $OfficeRepairErrObject
             }

             <#                   
                 if(!$O365PreCheck){
                   return
                 }

             #>
          }
          
        
          ################ Repair Office ##########################
                
           if(!$MSOfficeVersion.ClickToRun){
              
              $InstallationType = "MSI"
              $RepairTypeString = "Repair"

              $RepairXmlPath = "$env:PUBLIC\Documents\RepairMSOfficeXml.xml"
              RepairMSOfficeXml | Out-File $RepairXmlPath

              $RepairString = $MSOfficeVersion.UninstallString.Split('/')[0]  
             
              $cmdLine = $($RepairString)
              $cmdSubArgs = "/repair PROPLUSR /config $RepairXmlPath" 

           }else{             
               
               $InstallationType = "ClickToRun"  
                     
               $index = $MSOfficeVersion.ModifyPath.LastIndexOf('"')
               $RepairString = $MSOfficeVersion.ModifyPath.Substring(0,$index + 1)
               $LanguageString = $MSOfficeVersion.ModifyPath.Substring($index + 1).Split(" ")[3]
               $PlatformString = $MSOfficeVersion.ModifyPath.Substring($index + 1).Split(" ")[2]
                                
               $cmdLine = $($RepairString)

               if(!$IsQuickRepairSelected){
                 $RepairTypeString = "FullRepair"
               }else{
                 $RepairTypeString = "QuickRepair"
               }

               $cmdSubArgs = "scenario=Repair $($PlatformString) $($LanguageString) RepairType=$($RepairTypeString) DisplayLevel=$($DisplayRepairModal)" 
                       
           } 
       
        $MSOfficeVersion =  $MSOfficeVersion[0]
        $MSOfficeName = $MSOfficeVersion.DisplayName                
             
        if($MSOfficeVersion) {
            if(!($isInPipe)) {  
                          
               <# Write-Host "  " 
                Write-Host "Please wait while MS Office is being repaired..."  
                Write-Host "MS Office Version : $($MSOfficeName)"
                Write-Host "Installation Type : $($InstallationType)"
                Write-Host "Repair Type : $($RepairTypeString)" #>

               # $SuccessResult = "Please wait while MS Office is being repaired..."  

                $SuccessArr += "MS Office Version : $($MSOfficeName)"
                $SuccessArr += "Installation Type : $($InstallationType)"
                $SuccessArr += "Repair Type : $($RepairTypeString)" 
                               
                $SuccessOutArr += "MS Office Version : $($MSOfficeName)"
                $SuccessOutArr += "Installation Type : $($InstallationType)"
                $SuccessOutArr += "Repair Type : $($RepairTypeString)" 
               
               ###########################
                          
            }            
        }
       
         if($StartAndStopPrintSpooler){ # stop Spooler            
           Stop-Service -Name Spooler -Force
         }
          
         StartProcess -execFilePath $cmdLine -execParams $cmdSubArgs -WaitForExit $true  
         
         if($StartAndStopPrintSpooler){ # start Spooler
           Start-Service -Name Spooler 
         }                       
       
        if($MSOfficeVersion){           
                                                      
           if (!($isInPipe)) { 
              
               if($RepairTypeString -eq "FullRepair" ){

                # Write-Host "Please stay online while office downloads and install "
                $SuccessResult = "Please stay online while office downloads and install "

               }else{

                # Write-Host "MS Office has been repaired successfully"
                 $SuccessResult = "MS Office has been repaired successfully"
                    
               }
               
            }
          
        }                                      
                                                                               
        if ($isInPipe) {
            $results = new-object PSObject[] 0;
            $Result = New-Object -TypeName PSObject 
            Add-Member -InputObject $Result -MemberType NoteProperty -Name "TargetFilePath" -Value $TargetFilePath
            $Result
        }

        $SuccessArr = @()
        $SuccessOutArr = @()        
        $SuccessResult = ""   

        $OfficeRepairSuccObj =  New-Object psobject
        $OfficeRepairSuccObj | Add-Member -MemberType NoteProperty -Name MsOfficeRepair -Value $SuccessArr
       
        $OfficeRepairSuccObject.Result = $SuccessResult
        $OfficeRepairSuccObject.stdout = $SuccessOutArr
        $OfficeRepairSuccObject.dataObject = $OfficeRepairSuccObj  

        return $OfficeRepairSuccObject

    }
}

Function Get-OfficeVersion  {

      [OutputType('System.Software.Inventory')]
      [Cmdletbinding()] 

  Param( 

    )         

  Begin { 
    
    [String]$Computername=$env:COMPUTERNAME
          
    $installKeys = 'SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall',
                   'SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall'
                       
    $defaultDisplaySet = 'DisplayName','Version','ComputerName','ClientCulture'

    $defaultDisplayPropertySet = New-Object System.Management.Automation.PSPropertySet('DefaultDisplayPropertySet',[string[]]$defaultDisplaySet)
    $PSStandardMembers = [System.Management.Automation.PSMemberInfo[]]@($defaultDisplayPropertySet)
  
  }

  Process  {   
  
     $array = @()
     $MSexceptionList = "mui","visio","project","proofing","visual","Shared","Access","Excel","PowerPoint","Publisher"`
                        ,"Outlook","Proof","Components";  
           
       If  (Test-Connection -ComputerName  $Computername -Count  1 -Quiet) {
                  
            ForEach($Path in $installKeys) { 

                  Write-Verbose  "Checking Path: $Path"
                                    
                  Try  { $reg=[microsoft.win32.registrykey]::OpenRemoteBaseKey('LocalMachine',$Computername )}
                  Catch  { 
                           Write-Error $_  
                           Continue
                         } 
                                            
                   Try  {                         
                          $regkey=$reg.OpenSubKey($Path)                          
                          $subkeys=$regkey.GetSubKeyNames()      
                                                    
                          ForEach ($key in $subkeys){   

                                  Write-Verbose "Key: $Key"
                                  $thisKey=$Path+"\\"+$key 

                                  Try {  

                                          $thisSubKey=$reg.OpenSubKey($thisKey)
                                          $DisplayName =  $thisSubKey.getValue("DisplayName")

                                          If ($DisplayName  -AND $DisplayName  -notmatch '^Update  for|rollup|^Security Update|^Service Pack|^HotFix') {

                                              $Date = $thisSubKey.GetValue('InstallDate')
                                                                                           
                                              $Publisher =  Try { $thisSubKey.GetValue('Publisher').Trim()}
                                                            Catch { $thisSubKey.GetValue('Publisher')} 

                                              $Version = Try { $thisSubKey.GetValue('DisplayVersion').TrimEnd(([char[]](32,0)))} 
                                                         Catch { $thisSubKey.GetValue('DisplayVersion')}

                                              $UninstallString =  Try {$thisSubKey.GetValue('UninstallString').Trim()} 
                                                                  Catch {$thisSubKey.GetValue('UninstallString')}

                                              $InstallLocation =  Try {$thisSubKey.GetValue('InstallLocation').Trim()} 
                                                                  Catch {$thisSubKey.GetValue('InstallLocation')}

                                              $InstallSource =  Try {$thisSubKey.GetValue('InstallSource').Trim()}
                                                                Catch {$thisSubKey.GetValue('InstallSource')}

                                              $ModifyPath =  Try {$thisSubKey.GetValue('ModifyPath').Trim()}
                                                                Catch {$thisSubKey.GetValue('ModifyPath')}
                                              
                                              $VersionCount = 0  
                                              
                                              if(!$DisplayName.ToUpper().Contains("MICROSOFT OFFICE")){continue}
                                              
                                              if ($DisplayName.ToUpper().Contains("MICROSOFT OFFICE")) {
                                                 $isOfficePrimaryProduct = $True
                                                  foreach($exception in $MSexceptionList){
                                                    
                                                     if($DisplayName.ToLower().Contains($exception.ToLower())){

                                                        $isOfficePrimaryProduct = $false                                                        
                                                     }
                                                  }
                                                  if(!$isOfficePrimaryProduct){continue}
                                               }

                                              $clickToRunComponent = $false

                                              if(!($UninstallString.Contains("OfficeClickToRun") -or                                              
                                              $UninstallString.Contains("setup"))){continue}

                                              $clickToRunComponent = $false

                                              if ($UninstallString.Contains("OfficeClickToRun")) {
                                                     $clickToRunComponent = $true
                                                 }

                                              $VersionCount += 1

                                              $Object = [pscustomobject]@{

                                                  Computername = $Computername
                                                  DisplayName = $DisplayName
                                                  Version  = $Version
                                                  Publisher = $Publisher
                                                  UninstallString = $UninstallString                                                  
                                                  InstallLocation = $InstallLocation
                                                  InstallSource  = $InstallSource
                                                  ClickToRun = $clickToRunComponent
                                                  ModifyPath = $ModifyPath
                                                  Count = $VersionCount
                                              }

                                              $Object.pstypenames.insert(0,'System.Software.Inventory')     
                                              $object | Add-Member MemberSet PSStandardMembers $PSStandardMembers
                                                                                     
                                            $array += $Object 

                                    }# End Display if

                                } # End Display try
                                Catch {  Write-Warning "$Key : $_" }   

                           } #End subkey foreach 

                       } # try Drill down into the Uninstall key using the OpenSubKey Method 
                       Catch  {}   

                $reg.Close() 

             } # End paths foreach                 

         } # end test connection if         
         Else  { Write-Error  "$($Computername): unable to reach remote system!" }

     return $array 

   } # end process 

}

Function StartProcess {
	Param
	(
        [Parameter()]
		[String]$execFilePath,

        [Parameter()]
        [String]$execParams,

        [Parameter()]
        [bool]$WaitForExit = $false,

        [Parameter()]
        [string]$LogFilePath
	)
    
    $currentFileName = Get-CurrentFileName
    Set-Alias -name LINENUM -value Get-CurrentLineNumber 
    $timeoutSeconds = 2700000
    Try
    {
        $startExe = new-object System.Diagnostics.ProcessStartInfo
        $startExe.FileName = $execFilePath
        $startExe.Arguments = $execParams
        $startExe.CreateNoWindow = $false
        $startExe.UseShellExecute = $false

        $execStatement = [System.Diagnostics.Process]::Start($startExe) 

        
        if ($WaitForExit) {
          $execStatement.WaitForExit($timeoutSeconds)
        }
               
    }
    Catch
    {
        Write-Host $_.Exception.Message       
    }
}

Function IsDotSourced() {
  [CmdletBinding(SupportsShouldProcess=$true)]
  param(
    [Parameter(ValueFromPipelineByPropertyName=$true)]
    [string]$InvocationLine = ""
  )
  $cmdLine = $InvocationLine.Trim()
  Do {
    $cmdLine = $cmdLine.Replace(" ", "")
  } while($cmdLine.Contains(" "))

  $dotSourced = $false
  if ($cmdLine -match '^\.\\') {
     $dotSourced = $false
  } else {
     $dotSourced = ($cmdLine -match '^\.')
  }

  return $dotSourced
}

function Get-CurrentLineNumber {
    $MyInvocation.ScriptLineNumber
}

function Get-CurrentFileName{
    $MyInvocation.ScriptName.Substring($MyInvocation.ScriptName.LastIndexOf("\")+1)
}

$dotSourced = IsDotSourced -InvocationLine $MyInvocation.Line

if (!($dotSourced)) {

  
  if($PSVersionTable.PSVersion.Major -eq 2){

         $r = Repair-MSOffice -ComputerName $ComputerName  -WaitForInstallToFinish $WaitForInstallToFinish -TargetFilePath $TargetFilePath 
         $r |  ConvertTo-JSONP2

    }else{

       $r = Repair-MSOffice -ComputerName $ComputerName  -WaitForInstallToFinish $WaitForInstallToFinish -TargetFilePath $TargetFilePath 
       $r |  ConvertTo-Json -Depth 10

    }
}

