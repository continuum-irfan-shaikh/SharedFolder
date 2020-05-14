$Global:HashData = @()
$STDErr = @()
$Global:OutputObj = @()

$Global:STD_Error = @()
$Primary_DNS = ''
$Secondry_DNS = ''


#$networkName  = "Local Area Connection 3"
#$DNS_Server = "8.8.8.4"
#$Protocol_IPV4_IPV6 = 'IPV4'


$dnsServers   = ("$DNS_Server")
$dnsServers = $dnsServers.split(',')

if($dnsServers[0] -ne $null -or $dnsServers[0] -eq ""){
   $Primary_DNS = $($dnsServers[0])
}
if($dnsServers[1] -ne $null -or $dnsServers[1] -eq ""){
   $Secondry_DNS = $($dnsServers[1])
}
$dnsServers   = ("$Primary_DNS","$Secondry_DNS")

$dnsServers = $dnsServers | ? {$_}

#$ip_address   = '10.0.0.1'
#$subnetMask   = '255.0.0.0'
#$gateway      = '10.0.0.2'
$Global:STD_Error = $null

$Global:FinalDNS1 = @()
cls
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
      
#Write-Host "Powershell Version : $($hostVersion)"
if(-not $isPsVersionOk){
   
  Write-Warning "PowerShell version below 2.0 is not supported"
  return 
 
}
 
#Write-Host "OS Name : $((Get-WMIObject win32_operatingsystem).Name.ToString().Split("|")[0])"  
if(-not $isOSVersionOk){
 
   Write-Warning "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
   return 
 
}

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

###Write-Host ($stderr -eq $null)

    if($stderr -eq $null -and $HeadingUnder_dataObject_Array -ne $null){
        ###Write-Host 'NOstrerror -and YESHeadingUnder_dataObject_Array'
    $global:axa = ConvertTo-Json2 -InputObject @{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @($stdout);objects = "$objects";result = "$result"
    dataObject = @{"$HeadingUnder_dataObject_Array" = @($Array)}
    }
    }
    

    if($stderr -ne $null -and $HeadingUnder_dataObject_Array -eq $null){
        ###Write-Host 'YESstrerror -and NOHeadingUnder_dataObject_Array'
    $global:axa = ConvertTo-Json2 -InputObject @{taskName = "$taskName";status = "$status";Code = "$Code";stdout = @($stdout);result = "$result"
        stderr = @($stderr)
    }
    }


    if($stderr -ne $null -and $HeadingUnder_dataObject_Array -ne $null){
        ###Write-Host 'YESstrerror -and YESHeadingUnder_dataObject_Array'
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


function CheckAnyofUnreach{
$PrimTest = $Global:PRoceed_Data | ? {$_.Prim_Second -eq 'Prim'} | select * 
$SecondTest = $Global:PRoceed_Data | ? {$_.Prim_Second -eq 'Second'} | select * 

if($PrimTest.'Process for DNS Setting Changes' -eq 'No' -and $SecondTest.'Process for DNS Setting Changes' -eq 'No')
{
$Global:STD_Error = "Primary ($($PrimTest.'Address')) and Secondary $($SecondTest.'Address') DNS Address are not Reachable" 

return $false
    
}else{
    if($PrimTest.'Process for DNS Setting Changes' -eq 'No'){

$Global:STD_Error = "Primary DNS Address $($PrimTest.'Address') is not Reachable" 

    return $false
    }
    if($SecondTest.'Process for DNS Setting Changes' -eq 'No'){

    $Global:STD_Error = "Secondary DNS Address $($SecondTest.'Address') is not Reachable" 
    return $false
    }
return $True    
}}


#Write-Host "`n##########################################"
#Write-Host 'Validating the Network Adapter existance'
#Write-Host "########################################## `n"

############################################################################################
##################### Process for Setting Changes #######################
############################################################################################
function Process_For_Change_Setting
{
param($networkName)

        ######## Get Adapter Settings
        $adapter = Gwmi win32_networkadapter | where {$_.NetConnectionID -eq $networkName}

        ######## Get Adapter Settings    
        $adapterSettings = Get-WmiObject win32_networkAdapterConfiguration | where {$_.index -eq $adapter.index}
        
        ######### Get Target Adapter DNS config
        $global:Current_Adapter_Config = Get_Current_Adapter_Info -adapterSettings $adapterSettings

        $Sec = $Global:FinalDNS | ? {$_.Prim_Second -eq 'Second'} | % {$_.Address}
        $Prim = $Global:FinalDNS | ? {$_.Prim_Second -eq 'Prim'} | % {$_.Address}
        $Prim_Current_config = $global:Current_Adapter_Config.DNSServers.split(',')[0]
        $Sec_Current_config = $global:Current_Adapter_Config.DNSServers.split(',')[1]

        $dnsServers = ProcessForChanges -Primary_DNS $Prim -Secondry_DNS $Sec -Primary_DNS_Current_Config $Prim_Current_config -Secondry_DNS_Current_Config $Sec_Current_config
        
        if($Sec -ne $null -and $prim -eq $null){
        #Write-Host "process for changes : $($Sec)"
        }
        if($prim -ne $null -and $Sec -eq $null){
        #Write-Host "process for changes : $($prim)"
        }
        if($Sec -ne $null -and $prim -ne $null){        
        #Write-Host "process for changes : $($prim,$Sec)"
        }                
        #Write-Host "`nAttemp to change DNS configuration.." -ForegroundColor yellow
               
        $Specified_Setting = 'DNS Server_Address Changes'
        $dnsChange_ReturnValue = ($adapterSettings.SetDNSServerSearchOrder($dnsServers)).ReturnValue
        #sleep 2
       $Return_Info = Check_ReturnValue -Return_Value $dnsChange_ReturnValue -Specified_Setting $Specified_Setting
       
       if($dnsChange_ReturnValue -eq 0 -or $dnsChange_ReturnValue -eq 1)
       {
          $DNSIPs = ($dnsServers) -join ','
          #Write-Host "DNS is successfully changed`n" -ForegroundColor Green


      #Write-Host "`n#######################################################################"
      #Write-Host "Post Changes of Network Adapter: ($($networkName))"
      #Write-Host "#######################################################################`n"
    
    $adapterSettings = Get-WmiObject win32_networkAdapterConfiguration | where {$_.index -eq $adapter.index}  
    $global:Current_Adapter_Config1 = Get_Current_Adapter_Info -adapterSettings $adapterSettings
    #Write-Host "IPAddress     : $($global:Current_Adapter_Config1.IPAddress)"
    #Write-Host "SubnetMask    : $($global:Current_Adapter_Config1.SubnetMask)"
    #Write-Host "Gateway       : $($global:Current_Adapter_Config1.Gateway)" 
    #Write-Host "IsDHCPEnabled : $($global:Current_Adapter_Config1.IsDHCPEnabled)" 
    #Write-Host "DNSServers    : $($global:Current_Adapter_Config1.DNSServers)"           
    
    $global:Current_Adapter_Config1          
       }

}
#########################################################
function ProcessForChanges{
param($Primary_DNS,$Secondry_DNS,$Primary_DNS_Current_Config,$Secondry_DNS_Current_Config)

  #Write-Host "`n#########################################"
  #Write-Host "Changes in Progress of Network_Adapter Setting"
  #Write-Host "#########################################`n"
  
#$Primary_DNS   = '8.8.4.4'
#$Secondry_DNS  = '8.8.8.8'


#$Primary_DNS_Current_Config  = $null
#$Secondry_DNS_Current_Config  = '8.8.8.8'


if($Primary_DNS -ne $null -and $Secondry_DNS -ne $null)
{
   $DNSConfig = @("$Primary_DNS","$Secondry_DNS")  
        
}


if($Primary_DNS -ne $null -and $Secondry_DNS -eq $null)
{
    if($Secondry_DNS_Current_Config -ne $null)
    {
        $DNSConfig = @("$Primary_DNS","$Secondry_DNS_Current_Config")
    }

    if($Secondry_DNS_Current_Config -eq $null)
    {
        $DNSConfig = @("$Primary_DNS")
    }
        
}

if($Primary_DNS -eq $null -and $Secondry_DNS -ne $null)
{
    if($Primary_DNS_Current_Config -ne $null)
    {
        $DNSConfig = @("$Primary_DNS_Current_Config","$Secondry_DNS")
    }

    if($Primary_DNS_Current_Config -eq $null)
    {
        #$DNSConfig = @($null,"$Secondry_DNS")
        #Write-Host "No Primary found hence Valid Secondry address going to set as Primary Address $Secondry_DNS"          
        $DNSConfig = @("$Secondry_DNS")
        
    }    
}

return $DNSConfig

}


############################################################################################
##################### Validate the Adapter Setting Changes #######################
############################################################################################
function Check_ReturnValue
{ param($Return_Value,$Specified_Setting)


    ##Write-Host "$Specified_Setting Status"

    switch ($Return_Value) 
    {
        -1 { ''}
        0  {'Successful completion, no reboot required'; break}
        1  {'Successful completion, reboot required'; break}
        64 { 'Method not supported on this platform'; break}
        65 { 'Unknown failure'; break}
        66 { 'Invalid subnet mask'; break}
        67 { 'An error occurred while processing an Instance that was returned'; break}
        68 { 'Invalid input parameter'; break}
        69 { 'More than 5 gateways specified'; break}
        70 { 'Invalid IP address'; break}
        71 { 'Invalid gateway IP address'; break}
        72 { 'An error occurred while accessing the Registry for the requested information'; break}
        73 { 'Invalid domain name'; break}
        74 { 'Invalid host name'; break}
        75 { 'No primary/secondary WINS server defined'; break}
        76 { 'Invalid file'; break}
        77 { 'Invalid system path'; break}
        78 { 'File copy failed'; break}
        79 { 'Invalid security parameter'; break}
        80 { 'Unable to configure TCP/IP service'; break}
        81 { 'Unable to configure DHCP service'; break}
        82 { 'Unable to renew DHCP lease'; break}
        83 { 'Unable to release DHCP lease'; break}
        84 { 'IP not enabled on adapter'; break}
        85 { 'IPX not enabled on adapter'; break}
        86 { 'Frame/network number bounds error'; break}
        87 { 'Invalid frame type'; break}
        88 { 'Invalid network number'; break}
        89 { 'Duplicate network number'; break}
        90 { 'Parameter out of bounds'; break}
        91 { 'Access denied'; break}
        92 { 'Out of memory'; break}
        93 { 'Already exists'; break}
        94 { 'Path, file or object not found'; break}
        95 { 'Unable to notify service'; break}
        96 { 'Unable to notify DNS service'; break}
        97 { 'Interface not configurable'; break}
        98 { 'Not all DHCP leases could be released/renewed'; break}
        100 { 'DHCP not enabled on adapter'; break}
        2147786788 { "Write lock not enabled"; break}
        2147749891 { "Must be run with admin privileges"; break}
        default { "Faild with error code $($Return_Value)"; break}
    }

}
##############################################################################################################################

############################################################################################
##################### Validate DNS Address format and provided Protocol match ####
############################################################################################  
function Check_DNSIP_Format_Protocol
{
  param($Primary_DNS,$secondary_dns,$Protocol)

if($Primary_DNS -ne $null){
$Primary_DNS = "$Primary_DNS|Prim"
}
if($secondary_dns -ne $null){
$secondary_dns = "$secondary_dns|Second"
}

    if($Primary_DNS -ne $null -and $secondary_dns -ne $null){
    $dnsServers   = ("$Primary_DNS","$secondary_dns")
    }
    if($Primary_DNS -ne $null -and $secondary_dns -eq $null){
    $dnsServers   = ("$Primary_DNS")
    }
    if($Primary_DNS -eq $null -and $secondary_dns -ne $null){
    $dnsServers   = ("$secondary_dns")
    }

##Write-Host $dnsServers
      
    $testAddresses = $dnsServers
    $Protocol = "$Protocol"


function Test-IsValidIPv6Address 
{
    param(
        [Parameter(Mandatory=$true,HelpMessage='Enter IPv6 address to verify')] [string] $IP)
    $IPv4Regex = '(((25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2}))'
    $G = '[a-f\d]{1,4}'
    # In a case sensitive regex, use:
    #$G = '[A-Fa-f\d]{1,4}'
    $Tail = @(":",
        "(:($G)?|$IPv4Regex)",
        ":($IPv4Regex|$G(:$G)?|)",
        "(:$IPv4Regex|:$G(:$IPv4Regex|(:$G){0,2})|:)",
        "((:$G){0,2}(:$IPv4Regex|(:$G){1,2})|:)",
        "((:$G){0,3}(:$IPv4Regex|(:$G){1,2})|:)",
        "((:$G){0,4}(:$IPv4Regex|(:$G){1,2})|:)")
    [string] $IPv6RegexString = $G
    $Tail | foreach { $IPv6RegexString = "${G}:($IPv6RegexString|$_)" }
    $IPv6RegexString = ":(:$G){0,5}((:$G){1,2}|:$IPv4Regex)|$IPv6RegexString"
    $IPv6RegexString = $IPv6RegexString -replace '\(' , '(?:' # make all groups non-capturing
    [regex] $IPv6Regex = $IPv6RegexString
    if ($IP -imatch "^$IPv6Regex$") {
        $true
    } else {
        $false
    }
}

Function Test-IPv4Address($ipAddress) {
 if($testAddress -match "\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b") {
  $addressValid = $true
 } else {
  $addressValid = $false
 }
 return $addressValid
}
#--------------------------------------------------------------------------------------------------#

$IPAddress_Protocol_Format = @()
#--------------------------------------------------------------------------------------------------#
foreach($testAddress in $testAddresses) {
$Prim_Second = ($testAddress.split('|'))[1]
$testAddress = ($testAddress.split('|'))[0]


 if((Test-IPv4Address $testAddress)) {
  ##Write-Host "$testAddress is a valid formatted IPv4 Address" -foregroundColor Green

            if($Protocol -notmatch 'IPv4')
            {
              $matchs = 'Not matched'
            }
            if($Protocol -match 'IPv4')
            {
              $matchs = 'Matched'
            }            
             $IPAddress_Protocol_Format += New-Object psobject -Property @{
             "IPAddress"    = $testAddress
             "Provided Protocol"     = "$Protocol"
             "Protocol Detected"     = "IPv4"
             "Valid Format" = "YES"
             Prim_Second    = "$Prim_Second"
             "Protocol_IP Match" = "$matchs"}
             
  
 } else {
        
        if($testAddress -match ':')
         {
            if( (Test-IsValidIPv6Address $testAddress)) {
            ##Write-Host "$testAddress is a valid formatted IPv6 Address" -foregroundColor Green

            if($Protocol -notmatch 'IPv6')
            {
              $matchs = 'Not matched'
            }
            if($Protocol -match 'IPv6')
            {
              $matchs = 'Matched'
            }            
             $IPAddress_Protocol_Format += New-Object psobject -Property @{
             "IPAddress"    = $testAddress
             "Provided Protocol"     = "$Protocol"
             "Protocol Detected"     = "IPv6"
             "Valid Format" = "YES"
             Prim_Second    = "$Prim_Second"
             "Protocol_IP Match" = "$matchs"}
            
            }
            else{
            ##Write-Host "$testAddress is not a valid formatted IPV6 Address" -foregroundColor Red

            if($Protocol -notmatch 'IPv6')
            {
              $matchs = 'Not matched'
            }
            if($Protocol -match 'IPv6')
            {
              $matchs = 'Matched'
            }            
             $IPAddress_Protocol_Format += New-Object psobject -Property @{
             "IPAddress"    = $testAddress
             "Provided Protocol"     = "$Protocol"
             "Protocol Detected"     = "IPv6"
             "Valid Format" = "No"
             Prim_Second    = "$Prim_Second"
             "Protocol_IP Match" = "$matchs"}
                        
            }
          }
        
        if($testAddress -notmatch ':')
        {  
        ##Write-Host "$testAddress is not a valid formatted IPv4 Address" -foregroundColor Red

            if($Protocol -notmatch 'IPv4')
            {
              $matchs = 'Not matched'
            }
            if($Protocol -match 'IPv4')
            {
              $matchs = 'Matched'
            }            
             $IPAddress_Protocol_Format += New-Object psobject -Property @{
             "IPAddress"    = $testAddress
             "Provided Protocol"     = "$Protocol"
             "Protocol Detected"     = "IPv4"
             "Valid Format" = "No"
             Prim_Second    = "$Prim_Second"
             "Protocol_IP Match" = "$matchs"}
                    
        }
    }
}

$global:aa = $IPAddress_Protocol_Format | select IPAddress,'Provided Protocol','Protocol Detected','Protocol_IP Match','Valid Format',Prim_Second

$Protocol_Matched_Valid_Formata0 = $global:aa | ? {$_.'Valid Format' -eq 'Yes' -and $_.'Protocol_IP Match' -eq 'Matched' -and $_.IPAddress -ne ''}
$InValid_Format = $global:aa | ? {$_.'Valid Format' -eq 'No' -and $_.IPAddress -ne ''} 
$Protocol_IP_Match = $global:aa | ? {$_.'Protocol_IP Match' -eq 'Not matched' -and $_.IPAddress -ne ''}

if($InValid_Format -ne $null){


if( ($InValid_Format | % {$_.Prim_Second}) -contains 'Prim')
{
  if(($InValid_Format).GetType().BaseType.name -eq 'Array')
   {
    #Write-Host "Invalid format of Primary DNS Address : $(($InValid_Format | % {$_.IPAddress})[0]). | Expected format is: 192.168.129.1, 192.168.129.254"
    $Global:STD_Error += @("Invalid format of Primary DNS Address : $(($InValid_Format | % {$_.IPAddress})[0]). | Expected format is: 192.168.129.1, 192.168.129.254")

    }else{
    #Write-Host "Invalid format of Primary DNS Address : $($InValid_Format.IPAddress). | Expected format is: 192.168.129.1, 192.168.129.254"    
    $Global:STD_Error += @("Invalid format of Primary DNS Address : $($InValid_Format.IPAddress). | Expected format is: 192.168.129.1, 192.168.129.254")


    }

$STDErr = Get_stderr -Title 'Invalid DNS Address Format' -details $Global:STD_Error
$Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to Invalid DNS Address' -result 'Error: Unable to perform action' -stderr $STDErr

}

if(($InValid_Format | % {$_.Prim_Second}) -contains 'Second')
{
  if(($InValid_Format).GetType().BaseType.name -eq 'Array')
   {
    #Write-Host "Invalid format of Secondary DNS Address : $(($InValid_Format | % {$_.IPAddress})[1]). | Expected format is: 192.168.129.1, 192.168.129.254"
    $Global:STD_Error += @("Invalid format of Secondary DNS Address : $(($InValid_Format | % {$_.IPAddress})[1]). | Expected format is: 192.168.129.1, 192.168.129.254")



   }
   else{
    #Write-Host "Invalid format of Secondary DNS Address : $($InValid_Format.IPAddress). | Expected format is: 192.168.129.1, 192.168.129.254"
    $Global:STD_Error += @("Invalid format of Secondary DNS Address : $($InValid_Format.IPAddress). | Expected format is: 192.168.129.1, 192.168.129.254")


}


}
$STDErr = Get_stderr -Title 'Invalid DNS Address Format' -details $Global:STD_Error
$Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to Invalid DNS Address' -result 'Error: Unable to perform action' -stderr $STDErr
$Global:HashData
return
}


if($Protocol_IP_Match -ne $null){


if( ($Protocol_IP_Match | % {$_.Prim_Second}) -contains 'Prim'){

    #Write-Host  "Primary: protocol version does not match address format. Ex: Submitted an IPv4 Primary address with IPv6 protocol selected"

$Global:STD_Error += @("Protocol version does not match address format. Ex: Submitted an IPv4 Primary address with IPv6 protocol selected")

 }

if(($Protocol_IP_Match | % {$_.Prim_Second}) -contains 'Second'){
#Write-Host  "Second: protocol version does not match address format. Ex: Submitted an IPv4 Secondary address with IPv6 protocol selected"

$Global:STD_Error += @("Protocol version does not match address format. Ex: Submitted an IPv4 Secondary address with IPv6 protocol selected")

 }


$STDErr = Get_stderr -Title 'DNS Server Protocol mismatch' -details $Global:STD_Error
$Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address protocol mismatch' -result 'Error: Unable to perform action' -stderr $STDErr

$Global:HashData
 return
}



function InValid_Format
{
param($InValid_Format,$Primary_DNS,$Secondary_DNS)
    
    if($InValid_Format -ne $null)
    {
        if(($InValid_Format).GetType().BaseType.name -eq 'Array')
        {   if($InValid_Format.IPAddress -ne ""){ 
             #Write-Host "Invalid format of DNS Address (Primary : $("$Primary_DNS") and (Secondary : $("$Secondary_DNS").  Expected format is: 192.168.129.1, 192.168.129.254 `n" -ForegroundColor Red

            $Global:STD_Error = "Invalid format of DNS Address (Primary : $("$Primary_DNS") and (Secondary : $("$Secondary_DNS").  Expected format is: 192.168.129.1, 192.168.129.254"

             }
           }

          if(($InValid_Format).GetType().BaseType.name -eq 'Object')
           {  
              if($InValid_Format.IPAddress -ne ""){  
              
              $prim = $InValid_Format | ? {$_.Prim_Second -eq 'Prim'} | % {$_.IPAddress}
              $Sec = $global:aa | ? {$_.Prim_Second -eq 'Second'} | % {$_.IPAddress}
              
                if($Sec -ne $null -and $prim -eq $null){
                $Prim_Secc =  "Secondary"}
                if($prim -ne $null -and $Sec -eq $null){
                $Prim_Secc =  "Primary"}
                      
              #Write-Host "Invalid format of $Prim_Secc DNS Address : $($InValid_Format.IPAddress). | Expected format is: 192.168.129.1, 192.168.129.254 `n" -ForegroundColor Red
              
              $Global:STD_Error = "Invalid format of $Prim_Secc DNS Address : $($InValid_Format.IPAddress). | Expected format is: 192.168.129.1, 192.168.129.254"

              }
           }
     }
 }
 
if($Protocol_Matched_Valid_Formata0 -ne $null)
{
    if( ($Protocol_Matched_Valid_Formata0).GetType().BaseType.name -eq 'Object')
    {

       $1st_Phase_Result = $Protocol_Matched_Valid_Formata0

       if($1st_Phase_Result.Prim_Second -eq 'Prim')
       {
          
           InValid_Format -InValid_Format $InValid_Format -Primary_DNS $Primary_DNS 
           #Write-Host "Process for Reachiblity test against Primary DNS Address: $($Protocol_Matched_Valid_Formata0.IPAddress)`n" -ForegroundColor Yellow

           $Global:FinalDNS1 = Check_DNS_Server_Reachability -Primary_DNS $Primary_DNS

           if(!(CheckAnyofUnreach))
           {
           $STDErr = Get_stderr -Title 'Unable to process further' -details $Global:STD_Error
           $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr
           $Global:HashData
             return
           }
            if(!$Global:FinalDNS1){
           $STDErr = Get_stderr -Title 'Unable to process further' -details $Global:STD_Error
           $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr
           $Global:HashData
                   return ;
                   }

            if($Global:FinalDNS1 -eq $null){ 

           $STDErr = Get_stderr -Title 'Unable to process further' -details $Global:STD_Error
           $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr
           $Global:HashData
           Return ;

           }
           else{
           $Global:FinalDNS = $Global:FinalDNS1
           if(($Global:PRoceed_Data | % {$_.'Process for DNS Setting Changes'}) -contains 'NO'){

           #Write-Host $Global:STD_Error
           
           $STDErr = Get_stderr -Title 'Unable to process further' -details $Global:STD_Error
           $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr

           }else{
           Process_For_Change_Setting -networkName "$networkName"   
           }
           
           }

       }
       
       if($1st_Phase_Result.Prim_Second -eq 'Second')
       {
           InValid_Format -InValid_Format $InValid_Format -Secondary_DNS $Secondary_DNS 
           #Write-Host "`nProcess for Reachiblity test against Secondary DNS Address: $($Protocol_Matched_Valid_Formata0.IPAddress)`n" -ForegroundColor Yellow  
           
           $Global:FinalDNS1 = Check_DNS_Server_Reachability  -Secondary_DNS $Secondary_DNS          

           if(! (CheckAnyofUnreach))
           {
           $STDErr = Get_stderr -Title 'Unable to process further' -details $Global:STD_Error
           $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr
           $Global:HashData
             return
           }

             if(!$Global:FinalDNS1){

           $STDErr = Get_stderr -Title 'Unable to process further' -details $Global:STD_Error
           $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr
           $Global:HashData 
               return ;
               }      


        if($Global:FinalDNS1 -eq $null){ 

           $STDErr = Get_stderr -Title 'Unable to process further' -details $Global:STD_Error
           $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr
           $Global:HashData
           Return ;  
       }
       else{
           $Global:FinalDNS = $Global:FinalDNS1
           if(($Global:PRoceed_Data | % {$_.'Process for DNS Setting Changes'}) -contains 'NO'){

           #Write-Host $Global:STD_Error
           
           
           $STDErr = Get_stderr -Title 'Unable to process further' -details $Global:STD_Error
           $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr
           
           }else{
           Process_For_Change_Setting -networkName "$networkName"   
           }

           }

       }       



       
}
    if(($Protocol_Matched_Valid_Formata0).GetType().BaseType.name -eq 'Array')
    {
       #$aa | select IPAddress,'Provided Protocol','Protocol Detected','Protocol_IP Match','Valid Format' | ft
       InValid_Format -InValid_Format $Valid_Format -Primary_DNS $Primary_DNS -Secondary_DNS $Secondary_DNS
       #Write-Host "`nBoth DNS_Server Address are in valid format and Protocol also matched || Primary:($($Protocol_Matched_Valid_Formata0[0].IPAddress)) Secondary:($($Protocol_Matched_Valid_Formata0[1].IPAddress))`n" -ForegroundColor Yellow
       

       #Write-Host "Process for checking Reachiblity for Both DNS_Server Address : ($($Protocol_Matched_Valid_Formata0[0].IPAddress)) and ($($Protocol_Matched_Valid_Formata0[1].IPAddress))" -ForegroundColor Yellow
       #Write-Host ""
       $Global:FinalDNS1 = Check_DNS_Server_Reachability  -Primary_DNS "$Primary_DNS" -Secondary_DNS "$Secondary_DNS"
       
           if(! (CheckAnyofUnreach))
           {
           $STDErr = Get_stderr -Title 'Unable to process further' -details $Global:STD_Error
           $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr
           $Global:HashData
             return
           }

       if(!$Global:FinalDNS1){

           $STDErr = Get_stderr -Title 'Unable to process further' -details $Global:STD_Error
           $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr
           $Global:HashData
       return ;
       }

        if($Global:FinalDNS1 -eq $null){ 
           $STDErr = Get_stderr -Title 'Unable to process further' -details $Global:STD_Error
           $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr
           $Global:HashData
           Return ;  
       }
       else{
           $Global:FinalDNS = $Global:FinalDNS1
           if(($Global:PRoceed_Data | % {$_.'Process for DNS Setting Changes'}) -contains 'NO'){
           
           
           #Write-Host $Global:STD_Error
           sleep 100  
           $STDErr = Get_stderr -Title 'Unable to process further' -details $Global:STD_Error
           $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr

           }else{
           Process_For_Change_Setting -networkName "$networkName"   
           }

           }
    }

}
else{
       $Protocol_Matched_Valid_Formata1 | select IPAddress,'Provided Protocol','Protocol Detected','Protocol_IP Match','Valid Format' | ft
       #Write-Host "`nInvalid DNS Address format OR DNS Address Protocol and Provided Protocol doesn't match" -ForegroundColor red

       $Global:STD_Error = "Invalid DNS Address format OR DNS Address Protocol and Provided Protocol doesn't match"

       $STDErr = Get_stderr -Title 'DNS Server Protocol mismatch' -details $Global:STD_Error
       $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -status 'Failed' -Code '1' -stdout 'Unable to process further, Due to DNS Address unreachable issue' -result 'Error: Unable to perform action' -stderr $STDErr
       $Global:HashData 


       Return $False

}

}

############################################################################################
##################### Reachability test ####
############################################################################################  
function Check_DNS_Server_Reachability
{
   param($Primary_DNS,$secondary_dns,$Protocol)

    if($Global:STD_Error -ne $null){
    $Global:STD_Error
    return $false
    }

        #Write-Host "########################"
        #Write-Host "Reachiblity Valaidation..."
        #Write-Host "########################"
        #Write-Host ""
        
    if($Primary_DNS -ne $null -and $secondary_dns -ne $null){
    $dnsServers   = ("$Primary_DNS","$secondary_dns")
    }
    if($Primary_DNS -ne $null -and $secondary_dns -eq $null){
    $dnsServers   = ("$Primary_DNS")
    }
    if($Primary_DNS -eq $null -and $secondary_dns -ne $null){
    $dnsServers   = ("$secondary_dns")
    }
    
      
    $servers = $dnsServers

     $ping = new-object system.net.networkinformation.ping
     $pingreturns = @()

    foreach ($entry in $servers) {

        $Prim_Second = ($entry.split('|'))[1]
        $testAddress = ($entry.split('|'))[0]

      $pingreturns += $ping.send($testAddress) | select Status,@{name = 'Address';exp = {"$testAddress"}},@{name = 'Prim_Second';exp = {"$Prim_Second"}}
      #$pingreturns
    }

    $Success_To_contact = ($pingreturns | where {$_.Status -eq 'Success'}) | select Status,Address,Prim_Second

    $Failed_To_contact = ($pingreturns | where {$_.Status -ne 'Success'}) | select Status,Address,Prim_Second

    $Global:PRoceed_Data = @()

    if($Success_To_contact -ne $null)
    {
    $Global:PRoceed_Data += $Success_To_contact | select Address,@{name = 'Reachable Status';exp = {'YES'}},@{name = 'Process for DNS Setting Changes';exp = {'YES'}},Prim_Second
    }

    if($Failed_To_contact -ne $null)
    {
    $Global:PRoceed_Data += $Failed_To_contact | select Address,@{name = 'Reachable Status';exp = {'No'}},@{name = 'Process for DNS Setting Changes';exp = {'NO'}},Prim_Second

    }

    $Global:Process_Data_Changes     = $Global:PRoceed_Data | ? {$_.'Process for DNS Setting Changes' -eq 'YES'} | select Address,Prim_Second
    $Global:NOT_Process_Data_Changes = $Global:PRoceed_Data | ? {$_.'Reachable Status' -eq 'NO'} | % {$_.Address} | select Address,Prim_Second
    

    
    if($Global:Process_Data_Changes -ne $null)
    {
        if(($Global:Process_Data_Changes).GetType().BaseType.name -eq 'Array')
        {    
           #Write-Host "Process for DNS Address changes (Primary : $("$Primary_DNS")) and (Secondary : $("$Secondary_DNS"))" -ForegroundColor Yellow
           $Global:DNS_TO_Process = $dnsServers -join ','
        }

        if(($Global:Process_Data_Changes).GetType().BaseType.name -eq 'Object')
        {  
          if($Global:Process_Data_Changes.Prim_Second -eq 'Prim')
          {
            $prim1 = ($Primary_DNS.split('|'))[0]
           #Write-Host "Process for DNS Address changes (Primary : $("$prim1"))" -ForegroundColor Yellow
           $Global:DNS_TO_Process = $Process_Data_Changes    
          }
          
          if($Global:Process_Data_Changes.Prim_Second -eq 'Second')
          {
             $Second1 = ($secondary_dns.split('|'))[0]         
           #Write-Host "Process for DNS Address changes (Secondary : $("$Second1"))" -ForegroundColor Yellow
           $Global:DNS_TO_Process = $Process_Data_Changes    
          }
          
        }
    }

 
 
 if($Global:Process_Data_Changes -eq $null)
 {
      
     if($Secondary_DNS -ne $null)
     {
      #Write-Host "(Secondary : $("$Secondary_DNS")) DNS unreachable" -ForegroundColor red
     }
     if($Primary_DNS -ne $null)
     {
      #Write-Host "(Primary : $("$Primary_DNS")) DNS unreachable" -ForegroundColor red
     }
 }

    $Global:Process_Data_Changes
    
    #return $DNS_TO_Process
}

############################################################################################
##################### Get Current Adapter Info ####
############################################################################################  
function Get_Current_Adapter_Info
{
param($adapterSettings)

    foreach($Network in $adapterSettings)
    {
    $AdapterName = $networkName
    $IPAddress  = $Network.IpAddress[0]
    $SubnetMask  = $Network.IPSubnet[0]
    $DefaultGateway = $Network.DefaultIPGateway
    $DNSServers  = $Network.DNSServerSearchOrder
    $IsDHCPEnabled = $false
    If($network.DHCPEnabled) {
     $IsDHCPEnabled = $true
    }
    $Global:OutputObj  = New-Object -Type PSObject
    $Global:OutputObj | Add-Member -MemberType NoteProperty -Name AdapterName -Value $AdapterName
    $Global:OutputObj | Add-Member -MemberType NoteProperty -Name IPAddress -Value $IPAddress
    $Global:OutputObj | Add-Member -MemberType NoteProperty -Name SubnetMask -Value $SubnetMask
    $Global:OutputObj | Add-Member -MemberType NoteProperty -Name Gateway -Value ($DefaultGateway -join ",")      
    $Global:OutputObj | Add-Member -MemberType NoteProperty -Name IsDHCPEnabled -Value $IsDHCPEnabled
    $Global:OutputObj | Add-Member -MemberType NoteProperty -Name DNSServers -Value ($DNSServers -join ",")     
    $Global:OutputObj

    $ObjectCount = 1
    $Global:HashData = HashTable_Array_JSON -HeadingUnder_dataObject_Array "Post Changes of Network Adapter configuration" -Array $Global:OutputObj -taskName 'Change DNS Configuration task' -status 'Success' -Code '0' -stdout $Global:OutputObj -objects "$ObjectCount" -result 'Success: DNS Configuration has changed'
    ############################################################
    ############### Re-enable network adapter   
    ############################################################
    #$networkAdapter_Index = $network.index
    #Enable_Disable_Adapter -networkAdapter_Index $networkAdapter_Index

}

}

$adapter = Gwmi win32_networkadapter | where {$_.NetConnectionID -eq $networkName}
if($adapter)
    {

      #Write-Host "Found Network Adapter: ($($networkName))" -ForegroundColor Green

      $adapterSettings = Get-WmiObject win32_networkAdapterConfiguration | where {$_.index -eq $adapter.index}

        if($adapterSettings.DHCPEnabled -eq $True){
        
     #Write-Host "`n#############################################################"
     #Write-Host "DHCP Status for Adapter ($($networkName))"
     #Write-Host "#############################################################`n"

        #Write-Host "`n`nDHCP is Enabled, can't perform changes for DNS`n `n" -ForegroundColor Yellow
        return
        }


     #Write-Host "`n#############################################################"
     #Write-Host "Valaidation- Match DNS-Address protocol & DNS-Address Format"
     #Write-Host "#############################################################`n"

    $Check_Format_Protocol = Check_DNSIP_Format_Protocol -Primary_DNS "$Primary_DNS" -Secondary_DNS "$Secondry_DNS" -Protocol "$Protocol_IPV4_IPV6"

    if(!$Check_Format_Protocol){
      return
    }

    if($Global:FinalDNS -eq $null){
      return
    }
    else{     
      
   }
      #Write-Host '#######################################################################'
      #Write-Host "Network Adapter Setting Changes for Network Adapter: ($($networkName))"
      #Write-Host '#######################################################################'

    }
else
{
  #Write-Host "Network Adapter (($networkName)) does not exist, Hence retrieving available Network Adapters......." -ForegroundColor Yellow
  #Write-Host ''

  #Write-Host '#########################################'
  #Write-Host "List of available Network Adapters"
  #Write-Host '#########################################'
$Avail  =   get-wmiobject win32_networkadapter | where {$_.netconnectionid -ne $null}| select NetConnectionID,Name,Description,NetEnabled,index | Select NetConnectionID,Name,Description,NetEnabled,index


  $STDErr = Get_stderr -Title 'Network Adapter not exist' -details 'Network Adapter not exist,Hence retrieving available Network Adapters'
  $Global:HashData = HashTable_Array_JSON -taskName 'Change DNS Configuration task' -HeadingUnder_dataObject_Array "Available Network Adapter information" -Array $Avail -status 'Failed' -Code '1' -stdout "Network Adapter ($($networkName)) does not exist, Hence retrieving available Network Adapters" -result "Error: Network Adapter ($($networkName)) does not exist" -stderr $STDErr

}

$Global:HashData