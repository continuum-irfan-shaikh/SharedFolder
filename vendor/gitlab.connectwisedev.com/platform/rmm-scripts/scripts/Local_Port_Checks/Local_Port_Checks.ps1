$Specific_port = $null
$Multiple_port = $null
$Port_range = $null

#$RemoteServer = "localhost"

$Global:Closed1 = @()
$Global:Open = @()
$Global:Closed = @()
$Global:Invalid123 = @()
$da = @()
$Range = @()
$TotalPortsToCheck = @()
$PortsortStatus = @()
$Global:STDOUT = @()
$rangePort = @()
$WithoutRange = @()

#$Specific_port = ''
#$Multiple_port = ''
#$Port_range = '130-140'
#$All_Format = ''



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
    "detail" = "Port ($($_.Port)) Port is closed"
    }
    $ID++
    }

return $info
}



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
    $Portss_js=new-object system.web.script.serialization.javascriptSerializer
    return $Portss_js.Serialize($item)
}

function ConvertFrom-Json20([object] $item){ 
    add-type -assembly system.web.extensions
    $Portss_js=new-object system.web.script.serialization.javascriptSerializer

    #The comma operator is the array construction operator in PowerShell
    return ,$Portss_js.DeserializeObject($item)
}


function HashTable_Array_JSON{
Param($HeadingUnder_dataObject_Array,$Array,$taskName,$status,$Code,$stdout,$stderr,$objects,$result)

##write-host ($stderr -eq $null)

    if($stderr -eq $null -and $HeadingUnder_dataObject_Array -ne $null){
        ##write-host 'NOstrerror -and YESHeadingUnder_dataObject_Array'
    $global:axa = ConvertTo-Json2 -InputObject @{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @($stdout);objects = "$objects";result = "$result"
    dataObject = @{"$HeadingUnder_dataObject_Array" = @($Array)}
    }
    }
    

    if($stderr -ne $null -and $HeadingUnder_dataObject_Array -eq $null){
        ##write-host 'YESstrerror -and NOHeadingUnder_dataObject_Array'
    $global:axa = ConvertTo-Json2 -InputObject @{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @($stdout);result = "$result"
        stderr = @($stderr)
    }
    }


    if($stderr -ne $null -and $HeadingUnder_dataObject_Array -ne $null){
        ##write-host 'YESstrerror -and YESHeadingUnder_dataObject_Array'
    $global:axa = ConvertTo-Json2 -InputObject @{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @($stdout);result = "$result"
        stderr = @($stderr)
        dataObject = @{"$HeadingUnder_dataObject_Array" = @($Array)}
    }
    }    

#$BackTOHasTable = ConvertFrom-Json20 $axa
#$BackTOHasTable
#$status
$global:axa 
}
#####################################
##########################################################################


$Input12 = @()

if(![string]::IsNullOrEmpty($Specific_port)){
$Input12 += $Specific_port
}
if(![string]::IsNullOrEmpty($Multiple_port)){
$Input12 += $Multiple_port
}
if(![string]::IsNullOrEmpty($Port_range)){
$Input12 += $Port_range
}
if(![string]::IsNullOrEmpty($All_Format)){
$Input12 += $All_Format
}


if($Input12.count -eq 1){

if(![string]::IsNullOrEmpty($All_Format)){
#Check-LocalPort -Port_range $Port_range

write-host 'Mixed'

$da = $All_Format -split(',')
$WithoutRange = $da | ? {$_ -notmatch '-'}
$Range = ($da | ? {$_ -match '-'}) -split ('-')
$rangePort = $Range[0]..$Range[1]
}else{

if(![string]::IsNullOrEmpty($Specific_port)){
#Check-LocalPort -Specific_port $Specific_port

write-host 'Spc'

if($Specific_port -notlike '*,*' -and $Specific_port -notlike '*-*'){
'Valid Format'

$da = $Specific_port -split(',')
$WithoutRange = $da | ? {$_ -notmatch '-' -and $_ -notmatch ','}

}else{
"invalid Specific Spec"

$Global:Closed1 += @{Id = 0;Title = "submitted Format: $Specific_port";details = "Expected formatting for Specific port: 80"}
HashTable_Array_JSON -taskName 'Local Check Open Ports' -status 'Failed' -Code '1' -stdout 'Found invalid Port Format' -stderr $Global:Closed1 -result 'Error: Found invalid Ports'

return
}

}
if(![string]::IsNullOrEmpty($Multiple_port)){
#Check-LocalPort -Multiple_port $Multiple_port

write-host 'Multi'

if($Multiple_port -like '*,*' -and $Multiple_port -notlike '*-*'){
'Valid Format'

$da = $Multiple_port -split(',')
$WithoutRange = $da | ? {$_ -notmatch '-'}

}else{
'InValid Multiple Format'

$Global:Closed1 += @{Id = 0;Title = "submitted Format: $Multiple_port";details = "Expected formatting for Multiple ports: 80,81,82,83"}
HashTable_Array_JSON -taskName 'Local Check Open Ports' -status 'Failed' -Code '1' -stdout 'Found invalid Port Format' -stderr $Global:Closed1 -result 'Error: Found invalid Ports'

return
}

}
if(![string]::IsNullOrEmpty($Port_range)){
#Check-LocalPort -Port_range $Port_range

write-host 'Range'

if($Port_range -notlike '*,*' -and $Port_range -like '*-*'){
'Valid Format'

$Range = ($Port_range | ? {$_ -match '-'}) -split ('-')
$rangePort = $Range[0]..$Range[1]

}else{
'InValid rangePort Format'

#"submitted Format: $Port_range"


$Global:Closed1 += @{Id = 0;Title = "submitted Format: $Port_range";details = "Expected formatting for Port Range: 80-83"}
HashTable_Array_JSON -taskName 'Local Check Open Ports' -status 'Failed' -Code '1' -stdout 'Found invalid Port Format' -stderr $Global:Closed1 -result 'Error: Found invalid Ports'

return
}

}
}

}else{
write-host 'Please specify only one input'

return
}


if($rangePort -ne $null)
{
$TotalPortsToCheck = @($rangePort)
}

if($WithoutRange -ne $null)
{
$TotalPortsToCheck += @($WithoutRange)
}

$AllPorts = $TotalPortsToCheck | select -Unique



$inputValue = 0
ForEach ($Item  in $AllPorts)  {
$inputValid = [int]::TryParse(($Item), [ref]$inputValue)
    if (-not $inputValid) {
    write-host "This $Item is invalid" -ForegroundColor red
    
    $Global:Invalid123 = @($Item)
    return

    }else{
      #write-host "$Item Proceed for next port" -ForegroundColor Green

    }

}

if($Global:Invalid123.count -eq 0){

foreach($AllPorts1 in $AllPorts)
{
    $test = New-Object System.Net.Sockets.TcpClient;
    Try {
      Write-Host "Connecting to "$RemoteServer":"$AllPorts1" (TCP)..";
      $test.Connect($RemoteServer, $AllPorts1);
      Write-Host "$AllPorts1 Connection successful" -ForegroundColor Green; 
      
      $Global:Open += New-Object psobject -Property @{
      Hostname = "$RemoteServer"
      Port  = $AllPorts1
      Status= 'Open'
      } | Select Hostname,Port,Status

      }
    Catch { 
    
    Write-Host "$AllPorts1 Connection failed"-ForegroundColor Red ; 

      $Global:Closed += New-Object psobject -Property @{
      Hostname = "$RemoteServer"
      Port  = $AllPorts1
      Status= 'Closed'
      } | Select Hostname,Port,Status

    
    }
    Finally { $test.Dispose(); }
}

}


if($Global:Open.count -gt 0 -and $Global:Closed.count -eq 0){

##write-host 'Open'
#$Global:Open

$Object = $Global:Open.count
HashTable_Array_JSON -HeadingUnder_dataObject_Array "Tested ports report" -Array $Global:Open -taskName 'Local Check Open Ports' -status 'Success' -Code '0' -stdout $Global:Open -objects "$Object" -result 'Success: Retrieved port status'

}

if($Global:Open.count -eq 0 -and $Global:Closed.count -gt 0){

##write-host 'Closed'
#$Global:Closed

$Global:Closed1 = Get_stderr -Title 'Closed Port' -details $Global:Closed
HashTable_Array_JSON -taskName 'Local Check Open Ports' -status 'Failed' -Code '1' -stdout 'Found Closed Ports' -stderr $Global:Closed1 -result 'Error: Found Closed Ports'


}

if($Global:Open.count -gt 0 -and $Global:Closed.count -gt 0){
##write-host 'Open+Closed'
#$Global:Open
#$Global:Closed

$Global:Closed1 = Get_stderr -Title 'Closed Port' -details $Global:Closed
HashTable_Array_JSON -HeadingUnder_dataObject_Array "Open Ports result" -Array $Global:Open -taskName 'Local Check Open Ports' -status 'Failed' -Code '1' -stdout 'Found Closed and Open Ports' -stderr $Global:Closed1 -result 'Error: Found Closed and Open Ports'

}