$Global:Output = @()
$Global:Errorx  = @()
$global:modulePathx1 = @()
$Global:STDOUT = @()
$Global:HashData = @()

$erroractionpreference = 'silentlycontinue'

######################################### Convert to JSON Function
#####################################
################ Create STDerr data
######################################
function Get_stderr
{
param($Title,$details)
$info = @()
$ID = 0

 $details | %{
    $info += new-object psobject -Property @{
    "id" = $ID
    "title" = $Title
    "detail" = $_
    }
    $ID++
    }

return $info
}

#####################################
##########################################################################
#####################################
function FormatString {
    param(
        [String] $String)
    # removed: #-replace '/', '\/' `
    # This is returned 
    $String -replace '\\', '\\' -replace '\n', '\n' `
        -replace '\u0008', '\b' -replace '\u000C', '\f' -replace '\r', '\r' `
        -replace '\t', '\t' -replace '"', '\"'
}

function GetNumberOrString {
    param(
        $InputObject)
    if ($InputObject -is [System.Byte] -or $InputObject -is [System.Int32] -or `
        ($env:PROCESSOR_ARCHITECTURE -imatch '^(?:amd64|ia64)$' -and $InputObject -is [System.Int64]) -or `
        $InputObject -is [System.Decimal] -or $InputObject -is [System.Double] -or `
        $InputObject -is [System.Single] -or $InputObject -is [long] -or `
        ($Script:CoerceNumberStrings -and $InputObject -match $Script:NumberRegex)) {
        Write-Verbose -Message "Got a number as end value."
        "$InputObject"
    }
    else {
        Write-Verbose -Message "Got a string as end value."
        """$(FormatString -String $InputObject)"""
    }
}

function ConvertToJsonInternal {
    param(
        $InputObject, # no type for a reason
        [Int32] $WhiteSpacePad = 0)
    [String] $Json = ""
    $Keys = @()
    Write-Verbose -Message "WhiteSpacePad: $WhiteSpacePad."
    if ($null -eq $InputObject) {
        Write-Verbose -Message "Got 'null' in `$InputObject in inner function"
        $null
    }
    elseif ($InputObject -is [Bool] -and $InputObject -eq $true) {
        Write-Verbose -Message "Got 'true' in `$InputObject in inner function"
        $true
    }
    elseif ($InputObject -is [Bool] -and $InputObject -eq $false) {
        Write-Verbose -Message "Got 'false' in `$InputObject in inner function"
        $false
    }
    elseif ($InputObject -is [HashTable]) {
        $Keys = @($InputObject.Keys)
        Write-Verbose -Message "Input object is a hash table (keys: $($Keys -join ', '))."
    }
    elseif ($InputObject.GetType().FullName -eq "System.Management.Automation.PSCustomObject") {
        $Keys = @(Get-Member -InputObject $InputObject -MemberType NoteProperty |
            Select-Object -ExpandProperty Name)
        Write-Verbose -Message "Input object is a custom PowerShell object (properties: $($Keys -join ', '))."
    }
    elseif ($InputObject.GetType().Name -match '\[\]|Array') {
        Write-Verbose -Message "Input object appears to be of a collection/array type."
        Write-Verbose -Message "Building JSON for array input object."
        #$Json += " " * ((4 * ($WhiteSpacePad / 4)) + 4) + "[`n" + (($InputObject | ForEach-Object {
        $Json += "[`n" + (($InputObject | ForEach-Object {
            if ($null -eq $_) {
                Write-Verbose -Message "Got null inside array."
                " " * ((4 * ($WhiteSpacePad / 4)) + 4) + "null"
            }
            elseif ($_ -is [Bool] -and $_ -eq $true) {
                Write-Verbose -Message "Got 'true' inside array."
                " " * ((4 * ($WhiteSpacePad / 4)) + 4) + "true"
            }
            elseif ($_ -is [Bool] -and $_ -eq $false) {
                Write-Verbose -Message "Got 'false' inside array."
                " " * ((4 * ($WhiteSpacePad / 4)) + 4) + "false"
            }
            elseif ($_ -is [HashTable] -or $_.GetType().FullName -eq "System.Management.Automation.PSCustomObject" -or $_.GetType().Name -match '\[\]|Array') {
                Write-Verbose -Message "Found array, hash table or custom PowerShell object inside array."
                " " * ((4 * ($WhiteSpacePad / 4)) + 4) + (ConvertToJsonInternal -InputObject $_ -WhiteSpacePad ($WhiteSpacePad + 4)) -replace '\s*,\s*$' #-replace '\ {4}]', ']'
            }
            else {
                Write-Verbose -Message "Got a number or string inside array."
                $TempJsonString = GetNumberOrString -InputObject $_
                " " * ((4 * ($WhiteSpacePad / 4)) + 4) + $TempJsonString
            }
        #}) -join ",`n") + "`n],`n"
        }) -join ",`n") + "`n$(" " * (4 * ($WhiteSpacePad / 4)))],`n"
    }
    else {
        Write-Verbose -Message "Input object is a single element (treated as string/number)."
        GetNumberOrString -InputObject $InputObject
    }
    if ($Keys.Count) {
        Write-Verbose -Message "Building JSON for hash table or custom PowerShell object."
        $Json += "{`n"
        foreach ($Key in $Keys) {
            # -is [PSCustomObject]) { # this was buggy with calculated properties, the value was thought to be PSCustomObject
            if ($null -eq $InputObject.$Key) {
                Write-Verbose -Message "Got null as `$InputObject.`$Key in inner hash or PS object."
                $Json += " " * ((4 * ($WhiteSpacePad / 4)) + 4) + """$Key"": null,`n"
            }
            elseif ($InputObject.$Key -is [Bool] -and $InputObject.$Key -eq $true) {
                Write-Verbose -Message "Got 'true' in `$InputObject.`$Key in inner hash or PS object."
                $Json += " " * ((4 * ($WhiteSpacePad / 4)) + 4) + """$Key"": true,`n"            }
            elseif ($InputObject.$Key -is [Bool] -and $InputObject.$Key -eq $false) {
                Write-Verbose -Message "Got 'false' in `$InputObject.`$Key in inner hash or PS object."
                $Json += " " * ((4 * ($WhiteSpacePad / 4)) + 4) + """$Key"": false,`n"
            }
            elseif ($InputObject.$Key -is [HashTable] -or $InputObject.$Key.GetType().FullName -eq "System.Management.Automation.PSCustomObject") {
                Write-Verbose -Message "Input object's value for key '$Key' is a hash table or custom PowerShell object."
                $Json += " " * ($WhiteSpacePad + 4) + """$Key"":`n$(" " * ($WhiteSpacePad + 4))"
                $Json += ConvertToJsonInternal -InputObject $InputObject.$Key -WhiteSpacePad ($WhiteSpacePad + 4)
            }
            elseif ($InputObject.$Key.GetType().Name -match '\[\]|Array') {
                Write-Verbose -Message "Input object's value for key '$Key' has a type that appears to be a collection/array."
                Write-Verbose -Message "Building JSON for ${Key}'s array value."
                $Json += " " * ($WhiteSpacePad + 4) + """$Key"":`n$(" " * ((4 * ($WhiteSpacePad / 4)) + 4))[`n" + (($InputObject.$Key | ForEach-Object {
                    #Write-Verbose "Type inside array inside array/hash/PSObject: $($_.GetType().FullName)"
                    if ($null -eq $_) {
                        Write-Verbose -Message "Got null inside array inside inside array."
                        " " * ((4 * ($WhiteSpacePad / 4)) + 8) + "null"
                    }
                    elseif ($_ -is [Bool] -and $_ -eq $true) {
                        Write-Verbose -Message "Got 'true' inside array inside inside array."
                        " " * ((4 * ($WhiteSpacePad / 4)) + 8) + "true"
                    }
                    elseif ($_ -is [Bool] -and $_ -eq $false) {
                        Write-Verbose -Message "Got 'false' inside array inside inside array."
                        " " * ((4 * ($WhiteSpacePad / 4)) + 8) + "false"
                    }
                    elseif ($_ -is [HashTable] -or $_.GetType().FullName -eq "System.Management.Automation.PSCustomObject" `
                        -or $_.GetType().Name -match '\[\]|Array') {
                        Write-Verbose -Message "Found array, hash table or custom PowerShell object inside inside array."
                        " " * ((4 * ($WhiteSpacePad / 4)) + 8) + (ConvertToJsonInternal -InputObject $_ -WhiteSpacePad ($WhiteSpacePad + 8)) -replace '\s*,\s*$'
                    }
                    else {
                        Write-Verbose -Message "Got a string or number inside inside array."
                        $TempJsonString = GetNumberOrString -InputObject $_
                        " " * ((4 * ($WhiteSpacePad / 4)) + 8) + $TempJsonString
                    }
                }) -join ",`n") + "`n$(" " * (4 * ($WhiteSpacePad / 4) + 4 ))],`n"
            }
            else {
                Write-Verbose -Message "Got a string inside inside hashtable or PSObject."
                # '\\(?!["/bfnrt]|u[0-9a-f]{4})'
                $TempJsonString = GetNumberOrString -InputObject $InputObject.$Key
                $Json += " " * ((4 * ($WhiteSpacePad / 4)) + 4) + """$Key"": $TempJsonString,`n"
            }
        }
        $Json = $Json -replace '\s*,$' # remove trailing comma that'll break syntax
        $Json += "`n" + " " * $WhiteSpacePad + "},`n"
    }
    $Json
}

function ConvertTo-Json2 {
    [CmdletBinding()]
    #[OutputType([Void], [Bool], [String])]
    param(
        [AllowNull()]
        [Parameter(Mandatory=$true,
                   ValueFromPipeline=$true,
                   ValueFromPipelineByPropertyName=$true)]
        $InputObject,
        [Switch] $Compress,
        [Switch] $CoerceNumberStrings = $false)
    begin{
        $JsonOutput = ""
        $Collection = @()
        # Not optimal, but the easiest now.
        [Bool] $Script:CoerceNumberStrings = $CoerceNumberStrings
        [String] $Script:NumberRegex = '^-?\d+(?:(?:\.\d+)?(?:e[+\-]?\d+)?)?$'
        #$Script:NumberAndValueRegex = '^-?\d+(?:(?:\.\d+)?(?:e[+\-]?\d+)?)?$|^(?:true|false|null)$'
    }
    process {
        # Hacking on pipeline support ...
        if ($_) {
            Write-Verbose -Message "Adding object to `$Collection. Type of object: $($_.GetType().FullName)."
            $Collection += $_
        }
    }
    end {
        if ($Collection.Count) {
            Write-Verbose -Message "Collection count: $($Collection.Count), type of first object: $($Collection[0].GetType().FullName)."
            $JsonOutput = ConvertToJsonInternal -InputObject ($Collection | ForEach-Object { $_ })
        }
        else {
            $JsonOutput = ConvertToJsonInternal -InputObject $InputObject
        }
        if ($null -eq $JsonOutput) {
            Write-Verbose -Message "Returning `$null."
            return $null # becomes an empty string :/
        }
        elseif ($JsonOutput -is [Bool] -and $JsonOutput -eq $true) {
            Write-Verbose -Message "Returning `$true."
            [Bool] $true # doesn't preserve bool type :/ but works for comparisons against $true
        }
        elseif ($JsonOutput-is [Bool] -and $JsonOutput -eq $false) {
            Write-Verbose -Message "Returning `$false."
            [Bool] $false # doesn't preserve bool type :/ but works for comparisons against $false
        }
        elseif ($Compress) {
            Write-Verbose -Message "Compress specified."
            (
                ($JsonOutput -split "\n" | Where-Object { $_ -match '\S' }) -join "`n" `
                    -replace '^\s*|\s*,\s*$' -replace '\ *\]\ *$', ']'
            ) -replace ( # these next lines compress ...
                '(?m)^\s*("(?:\\"|[^"])+"): ((?:"(?:\\"|[^"])+")|(?:null|true|false|(?:' + `
                    $Script:NumberRegex.Trim('^$') + `
                    ')))\s*(?<Comma>,)?\s*$'), "`${1}:`${2}`${Comma}`n" `
              -replace '(?m)^\s*|\s*\z|[\r\n]+'
        }
        else {
            ($JsonOutput -split "\n" | Where-Object { $_ -match '\S' }) -join "`n" `
                -replace '^\s*|\s*,\s*$' -replace '\ *\]\ *$', ']'
        }
    }
}

function ConvertTo-Json20([object] $item){
    add-type -assembly system.web.extensions
    $ps_js=new-object system.web.script.serialization.javascriptSerializer
    return $ps_js.Serialize($item)
}

function ConvertFrom-Json20([object] $item){ 
    add-type -assembly system.web.extensions
    $ps_js=new-object system.web.script.serialization.javascriptSerializer

    #The comma operator is the array construction operator in PowerShell
    return ,$ps_js.DeserializeObject($item)
}


function HashTable_Array_JSON{
Param($HeadingUnder_dataObject_Array,$Array,$taskName,$status,$Code,$stdout,$stderr,$objects,$result)

if($psversiontable.PSVersion.Major -eq 2)
{
    if($stderr -eq $null){
    $axa = ConvertTo-Json2 -InputObject @{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @("$stdout");objects = "$objects";result = "$result"
    dataObject = @{"$HeadingUnder_dataObject_Array" = @($Array)}
    }
    }

    if($stderr -ne $null){
    $axa = ConvertTo-Json2 -InputObject @{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @("$stdout");result = "$result"
        stderr = @($stderr)
    }
    }
}else
{
#write-host 'its not Version 2'
    if($stderr -eq $null){
     $HashTable = [ordered]@{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @("$stdout");objects = "$objects";result = "$result"
        dataObject = @{"$HeadingUnder_dataObject_Array" = @($Array)}
        }
    $axa = ConvertTo-Json -InputObject $HashTable -Depth 100
    }

    if($stderr -ne $null){
    $HashTable = [ordered]@{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @("$stdout");result = "$result"
        stderr = @($stderr)
    }
    $axa = ConvertTo-Json -InputObject $HashTable -Depth 100
    }
}
#$BackTOHasTable = ConvertFrom-Json20 $axa
#$BackTOHasTable
$axa
}
######################################### Convert to JSON --END



################################################################ 
####The C# definition is stored in a PowerShell string and passed as a parameter to the Add-Type cmdlet
################################################################
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


##################################################################
########verify and pass CSharp Definition written in above code
##################################################################
if (-not ([System.Management.Automation.PSTypeName]'GetNuGetPackage').Type)
{

    Add-Type -TypeDefinition $NuGetPackageCode -Language CSharp
}


########################################################
########Check DotNetFrameWork version
######################################################## 

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
 

########################################################
########Unzip method for DOtNetFramwork 2.0
######################################################## 
 
function Unzip-20($zipfile, $destination)
{
    $shell = New-Object -ComObject Shell.Application
    $zip = $shell.NameSpace($zipfile)
    foreach($item in $zip.items())
    {       
        $shell.Namespace($destination).copyhere($item)
    }
}

########################################################
########Unzip method for DOtNetFramwork 4.5 and above
######################################################## 
 
function Unzip-45($zipfile, $destination)
{
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    param([string]$zipfile, [string]$destination)
    [System.IO.Compression.ZipFile]::ExtractToDirectory($zipfile, $destination)
}
 

#################################################################################
# Download Module Method, Works on 2.0 and above PSversion
#################################################################################
 
Function DownloadPowerShellModule([String] $ModuleName, [bool] $IsVersionFolderExist)
{
    #Write-Host "------------------------------------------"
    #Write-Host "Module ($ModuleName) Download and installation Status"
    #Write-Host "------------------------------------------"
    #Write-Host -ForegroundColor 10 "Attemp to Download $ModuleName ......"
    #Write-Host ""
    # https://psg-prod-eastus.azureedge.net/packages/azurerm.profile.5.8.3.nupkg
    # https://www.powershellgallery.com/packages/AzureRM.profile/5.8.3
 
    $packageObj = New-Object PSObject
    $packageObj | Add-Member Noteproperty -Name PackageAPI -value "https://www.powershellgallery.com/api"
    $packageObj | Add-Member Noteproperty -Name PackageVer -value "v2"
    $packageObj | Add-Member Noteproperty -Name Systeminfo -value "Systeminfo"
    $packageObj | Add-Member Noteproperty -Name DownloadString -value ""
 
    if($ModuleName -eq "Systeminfo"){
           
        $packageObj.DownloadString = $packageObj.PackageAPI +"/" + $packageObj.PackageVer +"/package/"+ $packageObj.Systeminfo
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
                        $Global:psModulePath = $psModulePath                          
                        $sourcePath = Join-Path -Path $unZipModulePath -ChildPath "\*"
                        Copy-Item -Path $sourcePath -Destination $psModulePath -recurse -Force
                    }    
                    catch
                    {                      
                        #Write-Host " Error : [$($_.Exception.Message)"] -ForegroundColor Red -BackgroundColor White    
                        return $false                  
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
                                $Global:psModulePath = $psModulePath
                                $sourcePath =  Join-Path -Path $folder.FullName -ChildPath (Get-ChildItem -Path $folder.FullName)                           
                                $sourcePathWithSubFolder = Join-Path -Path $sourcePath -ChildPath "\*"
 
                                Copy-Item -Path $sourcePathWithSubFolder -Destination $psModulePath -recurse -Force
                            }    
                            catch
                            {
                                #Write-Host " Error : [$($_.Exception.Message)"] -ForegroundColor Red -BackgroundColor White    
                                return $false                  
                           }
                        }
                    }                  
           
                }
           
            }
        }
       
       Remove-Item -LiteralPath $tempPath -Force -Recurse
       Remove-Item -LiteralPath $unZipModulePath -Force -Recurse
       #Write-Host -ForegroundColor 10 "Module is successfylly downloaded."


       return $true
               
    }catch{
      
        #################### Error when Unable to download ((systeminfo)) module
        #write-host "Please check internet accessibility- Unable to download ((systeminfo)) module from www.powershellgallery.com"
        #Write-Host " Error : [$($_.Exception.Message)"] -ForegroundColor Red -BackgroundColor White    

        $Global:Errorx  = 'Please check internet accessibility- Unable to download ((systeminfo)) module from www.powershellgallery.com'
        
        $STDErr = Get_stderr -Title 'Unable to process further,PS module systeminfo is not exist' -details "$Global:Errorx"
        $Global:HashData = HashTable_Array_JSON -taskName 'S.M.A.R.T' -status 'Failed' -Code '2' -stdout "$Global:Errorx" -result 'Error: Unable to perform action,PS module systeminfo is not exist' -stderr $STDErr
        $Global:HashData
        
        return $False
       
    }
 
}
 
 


#################################################################################
#Pass ModuleName To Start Retrieving S.M.A.R.T info
#################################################################################

Function Pass_ModuleName_To_Start_Retrieving_SMART_info
{
$ModuleName = 'systeminfo'
if(! (Get-Module -ListAvailable -Name $ModuleName))
{
    #write-host ''
    if(DownloadPowerShellModule -ModuleName "$ModuleName" -IsVersionFolderExist $false)
    {
    #write-host " "
    #write-host "Module is Successfully installed"
    #write-host " "
    #write-host "Retrieving the S.M.A.R.T information..."

    foreach($path in $env:PSModulePath.Split(';') )
      {if($path.contains("system32") -or $path.contains("Program Files") )
       {$modulePathx = Get-ChildItem "$path" | ? {$_.name -eq 'systeminfo'} | % {$_.FullName} 
       $WantFile = "$modulePathx"
       $FileExists = Test-Path $WantFile
       If ($FileExists -eq $True) {   $modulePathx1 = "$modulePathx"+"\"+"Systeminfo.psm1"
       if(Test-Path "$modulePathx1"){
       $global:modulePathx1 = $modulePathx1
       }}}}


    if([string]::IsNullOrEmpty($global:modulePathx1)){
    $Global:Errorx  = 'PSModule ((Systeminfo)) not exist, Hence unable to pull report'
    $STDErr = Get_stderr -Title 'Unable to process further' -details "$Global:Errorx"
    $Global:HashData = HashTable_Array_JSON -taskName 'S.M.A.R.T' -status 'Failed' -Code '1' -stdout "$Global:Errorx" -result 'Error: Unable to perform action' -stderr $STDErr
    $Global:HashData
    return
    }else{
       Import-Module -Name "$Global:psModulePath\Systeminfo.psm1" -Force
       $HDDInfo =  Get-SystemInfo -Properties HddSmart
    }
   }
   else{
   return
   }
}
else{
    #write-host ''
    #Write-Host "------------------------------------------"
    #Write-Host "Retrieving the S.M.A.R.T information..."
  


    foreach($path in $env:PSModulePath.Split(';') )
      {if($path.contains("system32") -or $path.contains("Program Files") )
       {$modulePathx = Get-ChildItem "$path" | ? {$_.name -eq 'systeminfo'} | % {$_.FullName} 
       $WantFile = "$modulePathx"
       $FileExists = Test-Path $WantFile
       If ($FileExists -eq $True) {   $modulePathx1 = "$WantFile"+"\"+"Systeminfo.psm1"
       if(Test-Path "$modulePathx1"){
       $global:modulePathx1 = $modulePathx1
       }}}}


    if([string]::IsNullOrEmpty($global:modulePathx1)){
    $Global:Errorx  = 'PSModule ((Systeminfo)) not exist, Hence unable to pull report'
    
    $STDErr = Get_stderr -Title 'Unable to process further' -details "$Global:Errorx"
    $Global:HashData = HashTable_Array_JSON -taskName 'S.M.A.R.T' -status 'Failed' -Code '1' -stdout "$Global:Errorx" -result 'Error: Unable to perform action' -stderr $STDErr
    $Global:HashData
    
    return
    }else{
    Import-Module -Name "$global:modulePathx1" -Force
       $HDDInfo =  Get-SystemInfo -Properties HddSmart
    }
   
    #$HDDInfo.HddSmart
   
}
 
$Data =    $HDDInfo.HddSmart| select SmartStatus,'SerialNumber','Model','Reallocated Sector Count','Offline Uncorrectable Sector Count','Temperature','Spin-Up Time','Reallocated Event Count'
 
 
#################### Error when No data Retrieved
#$Data = $null
if($Data -eq $null)
{
#write-host ""
#write-host "Error: Unable to Pull Data !"
#write-host ""
$Global:Errorx  = 'Unable to Pull Data !'

 $STDErr = Get_stderr -Title 'Unable to process further' -details "$Global:Errorx"
 $Global:HashData = HashTable_Array_JSON -taskName 'S.M.A.R.T' -status 'Failed' -Code '2' -stdout "$Global:Errorx" -result 'Error: Unable to perform action' -stderr $STDErr
 $Global:HashData

return
}
 
if( $Data.'Reallocated Sector Count' -eq $null){$Data.'Reallocated Sector Count' = 'N/A'}
if( $Data.'Offline Uncorrectable Sector Count' -eq $null){$Data.'Offline Uncorrectable Sector Count' = 'N/A'}
if( $Data.'Temperature' -eq $null){$Data.'Temperature' = 'N/A'}
if( $Data.'Spin-Up Time' -eq $null){$Data.'Spin-Up Time' = 'N/A'}
if( $Data.Model -eq $null){$Data.Model = 'N/A'}
if( $Data.Model -eq $null){$Data.Model = 'N/A'}
if( $Data.SerialNumber -eq $null){$Data.SerialNumber = 'N/A'}
if( $Data.SmartStatus -eq $null){$Data.SmartStatus = 'N/A'}
 
 
#write-host ''
#write-host '....................S.M.A.R.T information....................'
$Global:Output = New-Object psobject -Property @{
  'Disk Number' = ($Data.SerialNumber).trim()
    'Disk Friendly name' = $Data.Model
      'Reallocated Sector Count' = $Data.'Reallocated Sector Count'
        'Reported Uncorrectable Errors' = $Data.'Offline Uncorrectable Sector Count'
          'Temperature' = $Data.'Temperature'
            'Spin-Up Time' = $Data.'Spin-Up Time'
            'S.M.A.R.T. status' = $Data.SmartStatus
            } | select 'Disk Number','Disk Friendly name','S.M.A.R.T. status','Reallocated Sector Count','Reported Uncorrectable Errors','Temperature','Spin-Up Time'

    $Disk_Number = $Global:Output.'Disk Number'
    $Disk_Friendly_name = $Global:Output.'Disk Friendly name'
    $SMART_status = $Global:Output.'S.M.A.R.T. status'
    $Reallocated_Sector_Count = $Global:Output.'Reallocated Sector Count'
    $Reported_Uncorrectable_Errors = $Global:Output.'Reported Uncorrectable Errors'
    $Temperature = $Global:Output.'Temperature'
    $Spin_Up_Time = $Global:Output.'Spin-Up Time'
        
    $Global:STDOUT += "Disk Number : $Disk_Number`r`nDisk Friendly name : $Disk_Friendly_name`r`nS.M.A.R.T. status : $SMART_status`r`nReallocated Sector Count : $Reallocated_Sector_Count`r`nReported Uncorrectable Errors : $Reported_Uncorrectable_Errors`r`nTemperature : $Temperature`r`nSpin-Up Time : $Spin-Up Time"

    $ObjectCount = $Success.count
    $Global:HashData = HashTable_Array_JSON -HeadingUnder_dataObject_Array "S.M.A.R.T information" -Array $Global:Output -taskName 'S.M.A.R.T' -status 'Success' -Code '0' -stdout $Global:STDOUT -objects "$ObjectCount" -result 'Success: S.M.A.R.T Information retrived'
    $Global:HashData
}
 

######################################################
#### Verify Phisical / Virtual Machine
######################################################

function Verify_Phisical_VM
{
 
#-------------------------------------------------
   # Verify powershell version
   #-------------------------------------------------
  
   if(-not($PSVersionTable.PSVersion.Major -ge 2)){     
 
        #Write-Host "Powershell Version below 2.0 is not supported"
        $Global:Errorx  = 'Powershell Version below 2.0 is not supported'
      
      $STDErr = Get_stderr -Title 'Unable to process further' -details "$Global:Errorx"
      $Global:HashData = HashTable_Array_JSON -taskName 'S.M.A.R.T' -status 'Failed' -Code '1' -stdout "Unable to process further Due to PS version Not Supported" -result 'Error: Unable to perform action' -stderr $STDErr
      $Global:HashData
      return
 
   }else{ #Write-Host "Powershell Version : $($PSVersionTable.PSVersion.Major)" 
   }
 
   #-------------------------------------------------
   # Verify operating system
   #-------------------------------------------------
 
   if(-not([System.Environment]::OSVersion.Version.major -ge 6)){
 
      Write-Warning "Powershell script supports Window 7, Window 2008R2 and higher version operating system"
      $Global:Errorx  = 'Powershell script supports Window 7, Window 2008R2 and higher version operating system'

      $STDErr = Get_stderr -Title 'Unable to process further' -details "$Global:Errorx"
      $Global:HashData = HashTable_Array_JSON -taskName 'S.M.A.R.T' -status 'Failed' -Code '1' -stdout "Unable to process further Due to OS Version Not Supported" -result 'Error: Unable to perform action' -stderr $STDErr
      $Global:HashData

      return $false
 
    }else{#write-host  "Operating system : $((Get-WMIObject win32_operatingsystem).Name.ToString().Split("|")[0])" 
    }
 
$MachineType = get-wmiobject -computer LocalHost win32_computersystem | % {$_.Model}
 
if($MachineType -cmatch 'Virtual')
{
  #write-host ""
  #write-host ""
  #write-host "------------------------------------------------------"
  #write-host 'This is a virtual Machine - No further action required'
  #write-host "------------------------------------------------------"

  $Global:Errorx = 'This is a virtual Machine - No further action required'
  
  $STDErr = Get_stderr -Title 'Unable to process further' -details "$Global:Errorx"
  $Global:HashData = HashTable_Array_JSON -taskName 'S.M.A.R.T' -status 'Failed' -Code '1' -stdout "Unable to process further, Only works for Physical machine" -result 'Error: Unable to perform action' -stderr $STDErr
  $Global:HashData

  }else{
        Pass_ModuleName_To_Start_Retrieving_SMART_info
}
}
 
cls

############################
#### Executing Script
############################
Verify_Phisical_VM
