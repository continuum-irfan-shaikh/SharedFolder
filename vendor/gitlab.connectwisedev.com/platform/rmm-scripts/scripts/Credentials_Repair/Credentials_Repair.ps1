$Global:Error_Msg = @()
$Global:Success = @()

$Global:Output = @()
$Global:Errorx  = @()
$Global:STDOUT = @()
$Global:HashData = @()

$ErrorActionPreference = 'silentlycontinue'

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

function Get_DataObject
{
param($Title,$details)
$info = @()

 $details | %{
    $info += new-object psobject -Property @{
    "title" = $Title
    "Status" = $_
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
##write-host 'its not Version 2'
    if($stderr -eq $null){
     $HashTable = [ordered]@{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @("$stdout");objects = "$objects";result = "$result"
        dataObject = @{"$HeadingUnder_dataObject_Array" = @($Array)}
        }
    $axa = ConvertTo-Json -InputObject $HashTable -Depth 100
    }

    if($stderr -ne $null){
    $HashTable = [ordered]@{taskName = "$taskName";status = "$status";Code = "$Code";result = "$result";stdout = @("$stdout")
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


Function Unlock_Reset{
param($Unlock_Reset)
$Unlock_Reset = 'Unlock & Reset'

##############################
######Check PreCondition######
##############################

$hostVerSionMajor = ($PSVersionTable.PSVersion.Major).ToString()
$hostVerSionMinor = ($PSVersionTable.PSVersion.Minor).ToString()
$hostVersion = $hostVerSionMajor +'.'+ $hostVerSionMinor 
 
$osVersionMajor = ([System.Environment]::OSVersion.Version.major).ToString()
$osVersionMinor = ([System.Environment]::OSVersion.Version.minor).ToString()
$osVersion = $osVersionMajor +'.'+ $osVersionMinor
 
[boolean]$isPsVersionOk = ([version]$hostVersion -ge [version]'2.0')
[boolean]$isOSVersionOk = ([version]$osVersion -ge [version]'6.0')
      
#write-host "Powershell Version : $($hostVersion)"
if(-not $isPsVersionOk){
   
  Write-Warning "PowerShell version below 2.0 is not supported"
  return 
 
}
 
#write-host "OS Name : $((Get-WMIObject win32_operatingsystem).Name.ToString().Split("|")[0])"  
if(-not $isOSVersionOk){
 
   Write-Warning "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
   return 
 
}
######################################################################
#################################################################################################

#write-host ""
#write-host ""

# OS check and set the path
$OSArchitecture = (Get-WmiObject Win32_OperatingSystem).OSArchitecture
if($OSArchitecture -like '32-bit'){
$PasswordVault_FileLocation =  'C:\Program Files\SAAZOD\Executables\PasswordVault.exe'
$UserOptin_FileLocation = 'C:\Program Files\SAAZOD\UserOptin.exe'
}
if($OSArchitecture -like '64-bit'){
$PasswordVault_FileLocation =  'C:\Program Files (x86)\SAAZOD\Executables\PasswordVault.exe'
$UserOptin_FileLocation = 'C:\Program Files (x86)\SAAZOD\UserOptin.exe'
}

function CheckFileExistance
{
  param([string]$PasswordVault_FileLocation,[string]$UserOptin_FileLocation,[String]$Unlock_Reset)

   $Split_Path = Split-Path -Parent $UserOptin_FileLocation

    $PasswordVault_Existance_Checks = Test-Path "$PasswordVault_FileLocation"
    $UserOptin_Existance_Checks = Test-Path "$UserOptin_FileLocation"

  if(Test-path "$Split_Path")
  {

    if($PasswordVault_Existance_Checks -eq $True -and $UserOptin_Existance_Checks -eq $True)
    {
               ##write-host 'Both exe are exist- Executing it'
               
               #write-host "`n#################################"
               #write-host 'Execution Status..'
               #write-host "#################################`n"
               
                if(Execution -ExecutionPath $PasswordVault_FileLocation)
                {  
                   #write-host 'PasswordVault.exe successfully Executed' -ForegroundColor Green
                   
                   $Global:Success_MSG = 'PasswordVault.exe successfully Executed'
                   $Global:Success += Get_DataObject -Title 'Successfully Execution' -details "$Global:Success_MSG"

                     
                     if(Execution -ExecutionPath $UserOptin_FileLocation -Param $Unlock_Reset)
                     {
                        #write-host "UserOptin.exe successfully Executed" -ForegroundColor Green

                   $Global:Success_MSG = 'UserOptin.exe successfully Executed'
                   $Global:Success += Get_DataObject -Title 'Successfully Execution' -details "$Global:Success_MSG"
 
                           #write-host "`n#################################"
                           #write-host 'Credentials Repair status..'
                           #write-host "#################################`n"
                                      
                        #write-host "Credentials has been repaired" -ForegroundColor Green
                   
                   $Global:Success_MSG = 'Credentials has been repaired'
                   $Global:Success += Get_DataObject -Title 'Successfully repaired' -details "$Global:Success_MSG"                        
                   
                      }
                      else
                      {
                        #write-host 'UserOptin.exe Failed to Execute' -ForegroundColor Red
                        
                        $Global:Error_Msg = 'UserOptin.exe Failed to Execute'
                        $Global:STDErr += Get_stderr -Title 'Execution Error' -details "$Global:Error_Msg"
                      }
                }
                else
                {
                   #write-host 'PasswordVault.exe Failed to Execute' -ForegroundColor red

                  $Global:Error_Msg  = "PasswordVault.exe Failed to Execute"
                  $Global:STDErr += Get_stderr -Title 'PasswordVault.exe & UserOptin.exe Existance Error' -details "$Global:Error_Msg"

                }    
                
                
    }
    if($PasswordVault_Existance_Checks -eq $false -and $UserOptin_Existance_Checks -eq $false)
    {
                           #write-host "`n###############################################"
                           #write-host 'PasswordVault.exe UserOptin.exe existance Checks..'
                           #write-host "#################################################`n"
                               
               #write-host 'PasswordVault.exe and UserOptin.exe are not exist -- No further action' -ForegroundColor Red

                  $Global:Error_Msg  = "PasswordVault.exe and UserOptin.exe are not exist, Hence Can not process further Action"
                  $Global:STDErr += Get_stderr -Title 'PasswordVault.exe & UserOptin.exe Existance Error' -details "$Global:Error_Msg"

               return
    }
    if($PasswordVault_Existance_Checks -eq $False -and $UserOptin_Existance_Checks -eq $True)
    {
                           #write-host "`n###############################################"
                           #write-host 'PasswordVault.exe UserOptin.exe existance Checks..'
                           #write-host "#################################################`n"
               #write-host 'PasswordVault.exe not Found | - Downloading(PasswordVault.exe) and extruct to respective path' -ForegroundColor Yellow

               $Global:Success_MSG = "PasswordVault.exe not Found | - Downloading(PasswordVault.exe) and extruct to $($PasswordVault_FileLocation)"
               $Global:Success += Get_DataObject -Title 'Downloading(PasswordVault.exe) and extruct to respective path' -details "$Global:Success_MSG"
                Download_Extruct_PasswordVault -PasswordVault_FileLocation $PasswordVault_FileLocation -UserOptin_FileLocation $UserOptin_FileLocation -Unlock_Reset $Unlock_Reset             
    }
    if($PasswordVault_Existance_Checks -eq $True -and $UserOptin_Existance_Checks -eq $false)
    {
                           #write-host "`n###############################################"
                           #write-host 'PasswordVault.exe UserOptin.exe existance Checks..'
                           #write-host "#################################################`n"
                      
                      #write-host 'PasswordVault.exe Found -- UserOptin.exe not Found| -- No further action' -ForegroundColor red

                  $Global:Error_Msg  = "PasswordVault.exe and UserOptin.exe are not exist, Hence Can not process further Action"
                  $Global:STDErr += Get_stderr -Title 'PasswordVault.exe & UserOptin.exe Existance Error' -details "$Global:Error_Msg"

    }
}else
{
  #write-host "Directory : Not Found (($Split_Path))" -ForegroundColor red

  $Global:Error_Msg  = "Directory : Not Found ($($Split_Path))"
  $Global:STDErr += Get_stderr -Title 'Directory Existance Error' -details "$Global:Error_Msg"
}

}


#######################---------------------#########################
#Execute Programs in sequence -PasswordVault.exe and UserOptin.exe-
#######################---------------------#########################

function Execution
{
param($ExecutionPath,$Param)


    if($Param -eq $null)
    {
        $returnfromexe = Start-Process -FilePath "$ExecutionPath" -NoNewWindow -Wait -PassThru

            if(! ($returnfromexe.exitcode -eq 0))
            {
                 #write-host $returnfromexe.exitcode "failed"

                $Global:Error_Msg  = "Error while Executing executable Execution Failure Code :($($returnfromexe.exitcode))"
                $Global:STDErr += Get_stderr -Title 'Execution Error' -details "$Global:Error_Msg"

            }
            else
            {
               ##write-host 'Success'
               Return $True
            }
        }
      if("$param" -eq 'Unlock')
      { #write-host 'Unlocking.....'
            $returnfromexe = Start-Process -FilePath "$ExecutionPath" -ArgumentList 4 -NoNewWindow -PassThru
            return $True
        }
      if("$param" -eq 'Reset')
      {  #write-host 'Resetting....'
            $returnfromexe = Start-Process -FilePath "$ExecutionPath" -ArgumentList 5 -NoNewWindow -PassThru
            return $True
        }

      if("$param" -eq 'Unlock & Reset')
      {
            Start-Process -FilePath "$ExecutionPath" -ArgumentList 4 -NoNewWindow -PassThru
            Start-Process -FilePath "$ExecutionPath" -ArgumentList 5 -NoNewWindow -PassThru
         return $True
        }
}

#######################---------------------#########################
#Download PasswordVault and extruct it to respective location
#######################---------------------#########################

function Download_Extruct_PasswordVault{

param([string]$PasswordVault_FileLocation,[string]$UserOptin_FileLocation)

$CurrentPSVersion = $PSVersionTable.PSVersion.Major

#write-host "$CurrentPSVersion"
$url = "http://update.itsupport247.net/PasswordVault/PasswordVault.exe"

Try{

    if($CurrentPSVersion -lt 4)
    {
    (New-Object Net.WebClient).DownloadFile($url, $PasswordVault_FileLocation) 
    }

    if($CurrentPSVersion -ge 4)
    {
    Invoke-WebRequest -Uri $url -OutFile "$PasswordVault_FileLocation"
    }
 CheckFileExistance -PasswordVault_FileLocation $PasswordVault_FileLocation -UserOptin_FileLocation $UserOptin_FileLocation -Unlock_Reset $Unlock_Reset
    
}
catch
    {
  #write-host "Unable to Download Passwordvault.exe" -ForegroundColor red

  $Global:Error_Msg  = "Unable to Download Passwordvault.exe"
  $Global:STDErr += Get_stderr -Title 'Passwordvault.exe Downloading Error' -details "$Global:Error_Msg"
  
 }

}

#######################---------------------#########################
#Existance checks for PasswordVault.exe and UserOptin.exe 
#######################---------------------#########################



CheckFileExistance -PasswordVault_FileLocation $PasswordVault_FileLocation -UserOptin_FileLocation $UserOptin_FileLocation -Unlock_Reset $Unlock_Reset


}

Unlock_Reset


if($Global:Success -ne $null){

$ObjectCount = $Global:Success.Count

$d12 = @()
foreach($d1 in $Global:Success){

$d12 += @("$($d1.Status): $($d1.title)")

}
$stdout = $d12 -join '`r`n'

$Global:HashData = HashTable_Array_JSON -HeadingUnder_dataObject_Array "Results" -Array $Global:Success -taskName 'Credential Repair' -status 'Success' -Code '0' -stdout $stdout -objects "$ObjectCount" -result 'Success: Credential is Successfully reparied'
$Global:HashData

}

if($Global:STDErr -ne $null){

$Global:HashData = HashTable_Array_JSON -taskName 'Credential Repair' -status 'Failed' -Code '1' -stdout "Unable to process further, due to error occured" -result 'Error: Unable to perform action' -stderr $Global:STDErr
$Global:HashData
}