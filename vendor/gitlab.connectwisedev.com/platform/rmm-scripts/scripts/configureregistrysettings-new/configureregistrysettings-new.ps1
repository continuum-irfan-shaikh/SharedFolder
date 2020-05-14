<#
$Action = "Delete Key"

    Write Value Name,
    Delete Value,
    Add Key,
    Delete Key

$Hive = "HKEY_LOCAL_MACHINE"

    HKEY_CLASSES_ROOT,
    HKEY_CURRENT_USER,
    HKEY_LOCAL_MACHINE,
    HKEY_USERS,
    HKEY_CURRENT_CONFIG

$key = "HARDWARE\DESCRIPTION\System\Test\Test1"
$ForceFor64Bit
$Type = "REG_BINARY"

    REG_BINARY,
    REG_DWORD,
    REG_EXPAND_SZ,
    REG_MULTI_SZ,
    REG_SZ

$ValueName = "Test2"
$data = 000000000
$DeleteSubKeys = 
#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

if ($Hive -eq "HKEY_CURRENT_USER") {
    $loggedinuser = ((quser) -replace '^>', '') -replace '\s{2,}', ',' | ConvertFrom-Csv | ? { $_.SESSIONNAME -eq 'Console' } | Select -ExpandProperty USERNAME
    if (!($loggedinuser)) { Write-Error "No user is loggedin to this system"; Exit; }
    $SID = Get-WmiObject -Class Win32_UserAccount | ? { $_.Name -eq $loggedinuser } | Select -ExpandProperty SID
    $Hive = "HKEY_USERS" + "\" + "$SID"
    $Hive1 = "HKEY_CURRENT_USER"
}

$bit = (Get-WmiObject Win32_operatingsystem).Osarchitecture

Function Create_NewPath {
    if ($key -like "*WOW6432Node*") {
        $pathtocheck = "Registry::" + "$Hive" + "\" + "$key"
        return $pathtocheck 
    }
    $bit = (Get-WmiObject Win32_operatingsystem).Osarchitecture
    $Pathtocheck = "Registry::" + "$Hive" + "\" + $key.split("\")[0]
    if (($bit -eq "64-bit") -and (!($ForceFor64Bit)) -and ($Pathtocheck -ne "Registry::HKEY_LOCAL_MACHINE\SOFTWARE")) {
        Write-Error "Invalid Parameters."
        Exit;
    }
    if (($bit -eq "64-bit") -and (!($ForceFor64Bit)) -and ($Pathtocheck -eq "Registry::HKEY_LOCAL_MACHINE\SOFTWARE")) {
        $pathtocheck = "Registry::" + "$Hive" + "\" + "SOFTWARE" + "\" + "WOW6432Node" + "\" + ($key.Split("\", 2)[1])
        return $pathtocheck
    }
}

Function Convert_ValueBinary($Value) {
    $value = $value -replace "\s+", ""; $value = $value.ToCharArray()
    $final = @()
    $place = "odd"
    $len = $value.Count

    foreach ($Num in $value) {
        if ($place -eq "odd") {
            $final += $num;
            $place = "even"; Continue
        }
        if ($place -eq "even") {
            $final += $num
            $final += ","
            $place = "odd"
        }
        $i++
    }

    if ($final[$final.count - 1] -eq ",") {
        for ($j = 0; $j -lt $final.count - 1; $j++) {
            $final1 += $final[$j]
        }
    }

    if (($final[$final.count - 1] -like "*") -and ($final[$final.count - 2] -eq ",")) {
        for ($k = 0; $k -lt $final.count - 2; $k++) {
            $final1 += $final[$k]
        }
    }
    $Value = $final1
    $hexified = $Value -Split ',' | % { "0x$_" }
    Return $hexified
}

Function New_ItemProperty {
    $RegistryPath = "Registry::" + "$Hive" + "\" + "$key"
    switch ($Type) {
        "REG_SZ" { $Type = "String" }
        "REG_EXPAND_SZ" { $Type = "ExpandString" }
        "REG_MULTI_SZ" { $Type = "MultiString" }
        "REG_DWORD" { $Type = "DWord" }
        "REG_BINARY" {
            $Type = "Binary";
            $data = $data -Replace "`n",""
            if ($data -match '^[0-9 ]+$') {
                $Data = Convert_ValueBinary -Value $Data 
            }
        }
    }
    if (($Hive -eq "HKEY_LOCAL_MACHINE") -and (!($ForceFor64Bit )) -and ($bit -eq "64-bit")) {
        $RegistryPath = Create_NewPath   
    }
    if (!(Test-Path $RegistryPath)) {
        if ($Hive1 -eq "HKEY_CURRENT_USER") { $Hive = $Hive1; $RegistryPath = "Registry::" + "$Hive" + "\" + "$key" }
        Write-Error "Registry path does not exist.`nHive : $Hive `n Registry : $RegistryPath`nRegistry Entry : $ValueName"
        Exit;
    }
   
    try {
        $ErrorActionPreference = "SilentlyContinue"
        $AlreadyExist = Get-Itemproperty -Path $RegistryPath | Get-Member | Select -ExpandProperty Name
        $ErrorActionPreference = "Continue"
        if ($AlreadyExist -Contains $ValueName) {
            $Createdentry = New-ItemProperty -Path $RegistryPath -Name $ValueName -Value $data -PropertyType $Type -Force -ErrorAction Stop
            if ($Createdentry) {
                if ($Hive1 -eq "HKEY_CURRENT_USER") { $Hive = $Hive1; $RegistryPath = "Registry::" + "$Hive" + "\" + "$key" }
                "Registry entry updated successfully.`nHive : $Hive `nRegistry : $RegistryPath `nRegistry Entry : $ValueName"
            }
        }
        else {
            $Createdentry = New-ItemProperty -Path $RegistryPath -Name $ValueName -Value $data -PropertyType $Type -ErrorAction Stop
            if ($Createdentry) {
                if ($Hive1 -eq "HKEY_CURRENT_USER") { $Hive = $Hive1; $RegistryPath = "Registry::" + "$Hive" + "\" + "$key" }
                "Registry entry created successfully.`nHive : $Hive `nRegistry : $RegistryPath `nRegistry Entry : $ValueName"
            }
        }
    }
    catch {
        $_.Exception.Message
        Write-Error "Failed to create registry entry"
    }
}

Function Remove_ItemProperty {
    $RegistryPath = "Registry::" + "$Hive" + "\" + "$key"
    if (($Hive -eq "HKEY_LOCAL_MACHINE") -and (!($ForceFor64Bit)) -and ($bit -eq "64-bit")) {
        $RegistryPath = Create_NewPath
    }
    if (!(Test-Path $RegistryPath)) {
        Write-Error "Registry path does not exist."
        Exit;
    }
    try {
        Remove-ItemProperty -Path $RegistryPath -Name $ValueName -ErrorAction Stop
        if ($?) {
            if ($Hive1 -eq "HKEY_CURRENT_USER") { $Hive = $Hive1; $RegistryPath = "Registry::" + "$Hive" + "\" + "$key" }
            "Registry entry deleted successfully.`nHive : $Hive `nRegistry : $RegistryPath `nRegistry Entry : $ValueName"
        }
    }
    catch {
        $_.Exception.Message
        Write-Error "Failed to delete registry entry"
    }
}

Function New_Item {
    $RegistryPath = "Registry::" + "$Hive" + "\" + "$key"
    if (($Hive -eq "HKEY_LOCAL_MACHINE") -and (!($ForceFor64Bit)) -and ($bit -eq "64-bit")) {
        $RegistryPath = Create_NewPath
    }
    if (Test-Path $RegistryPath) {
        if ($Hive1 -eq "HKEY_CURRENT_USER") { $Hive = $Hive1; $RegistryPath = "Registry::" + "$Hive" + "\" + "$key" }
        Write-Error "Registry key $RegistryPath is already exist"
        Exit;
    }
    try {        
        $Createdkey = New-Item -Path $RegistryPath -Force -ErrorAction Stop
        if ($Createdkey) {
            if ($Hive1 -eq "HKEY_CURRENT_USER") { $Hive = $Hive1; $RegistryPath = "Registry::" + "$Hive" + "\" + "$key" }
            "Registry key created successfully.`nHive : $Hive `nRegistry : $RegistryPath"
            Exit;
        }
    }
    catch {
        $_.Exception.Message
        Write-Error "Failed to create registry key"
        Exit;
    }
}

function Remove_Item {
    $RegistryPath = "Registry::" + "$Hive" + "\" + "$key"
    if (($Hive -eq "HKEY_LOCAL_MACHINE") -and (!($ForceFor64Bit)) -and ($bit -eq "64-bit")) {
        $RegistryPath = Create_NewPath
    }
    if (!(Test-Path $RegistryPath)) {
        if ($Hive1 -eq "HKEY_CURRENT_USER") { $Hive = $Hive1; $RegistryPath = "Registry::" + "$Hive" + "\" + "$key" }
        Write-Error "Registry path does not exist.`nHive : $Hive `nRegistry : $RegistryPath"
        Exit;
    }
    $childitem = Get-ChildItem $RegistryPath
    try {
        if (($childitem -and ($DeleteSubKeys)) -or (!$childitem)) { 
            Remove-Item $RegistryPath -Recurse -ErrorAction Stop
            if (!(Test-Path $RegistryPath)) {
                "Registry Deleted Successfully."
            }
        }
        elseif ($childitem -and (!($DeleteSubKeys))) { "Registry contains one or more keys"; Exit }
    }
    catch {
        $_.Exception.Message
        Write-Error "Failed to delete registry key"
    }
}

switch ($Action) {
    "Write Value Name" { New_ItemProperty }
    "Delete Value" { Remove_ItemProperty }
    "Add Key" { New_Item }
    "Delete Key" { Remove_Item }
}

