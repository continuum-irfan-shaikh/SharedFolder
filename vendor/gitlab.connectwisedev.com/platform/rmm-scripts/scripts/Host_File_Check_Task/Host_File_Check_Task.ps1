#$Query_IP_Host = 'artifact.corp.continuum.net'

$SResult = @()
$Global:Success     = @()
$Global:Success_MSG = @()
$Query_Result = @()
$Global:Error_Msg = @()
$Global:STDErr = @()

$file = "$env:SystemDrive\Windows\System32\Drivers\etc\hosts"
$File_Contents_Exclude_COmments = Get-Content -Path $file |where {!$_.StartsWith("#") -and $_ -notlike ''}

$Global:FData = @()
foreach($File_Contents_Exclude_COmments1 in $File_Contents_Exclude_COmments)
{ $CombineData = ($File_Contents_Exclude_COmments1 -replace '(^\s+|\s+$)','' -replace '\s+',',').split(',')
$First = $CombineData[0]
$Second = $CombineData[1]
if("$First" -like "$Query_IP_Host" -or "$Second" -like "$Query_IP_Host"){$Global:FData += $File_Contents_Exclude_COmments1}
}

$daa= $Global:FData



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
    $axa = ConvertTo-Json2 -InputObject @{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @($stdout);objects = "$objects";result = "$result"
    dataObject = @{"$HeadingUnder_dataObject_Array" = @($Array)}
    }
    }

    if($stderr -ne $null){
    $axa = ConvertTo-Json2 -InputObject @{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @($stdout);result = "$result"
        stderr = @($stderr)
    }
    }
}else
{
##write-host 'its not Version 2'
    if($stderr -eq $null){
     $HashTable = [ordered]@{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @($stdout);objects = "$objects";result = "$result"
        dataObject = @{"$HeadingUnder_dataObject_Array" = @($Array)}
        }
    $axa = ConvertTo-Json -InputObject $HashTable -Depth 100
    }

    if($stderr -ne $null){
    $HashTable = [ordered]@{taskName = "$taskName";status = "$status";Code = "$Code";result = "$result";stdout = @($stdout)
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

function Query_Host_Entries
{ 
Param($Host_Contents_Data)

if($Host_Contents_Data.GetType().BaseType.name -eq 'Object'){  
#write-host "$($Host_Contents_Data)"
$Query_Result += $Host_Contents_Data | % {$_ -replace '(^\s+|\s+$)','' -replace '\s+','  '}
}



if($Host_Contents_Data.GetType().BaseType.name -eq 'Array'){  
    $Max = $Host_Contents_Data.count
    $Max1 = $Max -1

    if($Max -gt 0){
        #write-host "Found Two or more query result"
        #write-host "=============================="
    }
    
    for($i=0;$i -le $Max1; $i++)
    {$Index = $i +1
    #write-host "Query Result $Index : $($Host_Contents_Data[$i])"

    $Query_Result += $Host_Contents_Data[$i] | % {$_ -replace '(^\s+|\s+$)','' -replace '\s+','  '} 
    }

    $d12 = @()
    foreach($d1 in $Query_Result){

    $d12 += @($d1)

    }

    $ObjectCount = $Query_Result.count
    $Global:HashData = HashTable_Array_JSON -HeadingUnder_dataObject_Array "Query Result" -Array $Query_Result -taskName 'Host File Check Task' -status 'Success' -Code '0' -stdout $d12 -objects "$ObjectCount" -result "Success: Query result has found"

  }



}

function Check_Hosts_Contents_Lines
{ 
Param($File_Contents_Exclude_COmments)

    if($File_Contents_Exclude_COmments.count -le 10){
    
    #write-host "Number of Hosts lines Found excluding Comments & WhiteSpace: ($($File_Contents_Exclude_COmments.count))"
    #write-host "============================================================="

    $aa = $File_Contents_Exclude_COmments | % {$_ -replace '(^\s+|\s+$)','' -replace '\s+','  '}
    #$Global:Success += @{'Host Entries' = @($aa)}
    
    $d12 = @()
    foreach($d1 in $aa){

    $d12 += @($d1)

    }

    $ObjectCount = $aa.count
    $Global:HashData = HashTable_Array_JSON -HeadingUnder_dataObject_Array "Retrived Host Entries" -Array $aa -taskName 'Host File Check Task' -status 'Success' -Code '0' -stdout $d12 -objects "$ObjectCount" -result "Success: Number of Host entries below 10, Hence Displaying the Host Entries"

    }else{

    #write-host "Number of Hosts lines Found excluding Comments & WhiteSpace: ($($File_Contents_Exclude_COmments.count))"
    #write-host "==========================================================="
    #write-host 'Recommend technician check manually'
    #write-host 'Hosts File Path : %WINDIR%\system32\drivers\etc\hosts'

  $Global:Error_Msg  += "Recommend technician check manually |  Host path ; %WINDIR%\system32\drivers\etc\hosts"

  $Global:STDErr += Get_stderr -Title 'Message | Found more than 10 Host entries' -details $Global:Error_Msg
  $Global:HashData = HashTable_Array_JSON -stderr $Global:STDErr -taskName 'Host File Check Task' -status 'Failed' -Code '1' -stdout "Number of Host entries - $($File_Contents_Exclude_COmments.count)" -objects "$ObjectCount" -result "Success: Number of Host entries above 10, Hence Recommend technician check manually"
    
    }

}

if($daa.length -ne 0)
{
#write-host "`n##########################"
#write-host "Search query found"
#write-host "##########################`n"

Query_Host_Entries -Host_Contents_Data $daa
  

}if($daa.length -eq 0){

if($daa.length -eq 0){
#write-host "`n##########################"
#write-host "Search query is Blank"
#write-host "##########################`n"

Check_Hosts_Contents_Lines -File_Contents_Exclude_COmments $File_Contents_Exclude_COmments



}else{

#write-host "`n##########################"
#write-host "Search query not found"
#write-host "##########################`n"

Check_Hosts_Contents_Lines -File_Contents_Exclude_COmments $File_Contents_Exclude_COmments
}

}


$Global:HashData