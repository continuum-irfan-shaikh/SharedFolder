<#
.SYNOPSIS
   Copy INF file from Driverstore to c:\Windows\INF

.DESCRIPTION
  This script check for missing INF file for actiev NIC if missing it will copy the inf file form Driver store to C:\windows\INF.
  Supported Operating Systems : Windows 7 and Windows 2008 R2
  Version : 8.0.0
#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

$Computer = $env:COMPUTERNAME
$OSver = [System.Environment]::OSVersion.Version
$OSName = (Get-WMIObject win32_operatingsystem).Caption
if (-not($OSver.Major -eq 6 -and $OSver.Minor -eq 1 -and $OSver.Build -ge 7601)) {
    Write-Output "Not applicable for $OSName. No action needed."
    Exit
}
#Logging
$LogFolder = $env:SystemDrive + "\windows\Logs"
$LogFile = $LogFolder + "\INF-Copy-" + (Get-Date -UFormat "%d-%m-%Y") + ".log"
Function Write-Log {
    param (
        [Parameter(Mandatory = $True, Position = 1)]
        [string]$LogOutput,
        [Parameter(Mandatory = $false, Position = 2)]
        [string]$Path = $LogFile,
        [Parameter(Mandatory = $false, Position = 0)]
        [ValidateSet("Error", "Warn", "Info")]
        [string]$Level = "Info"
    )
    switch ($Level) {
        'Error' {
            Write-Error $LogOutput
            $LevelText = 'ERROR:'
        }
        'Warn' {
            Write-Warning $LogOutput
            $LevelText = 'WARNING:'
        }
        'Info' {
            Write-Verbose $LogOutput
            $LevelText = 'INFO:'
        }
    }

    $currentDate = (Get-Date -UFormat "%d-%m-%Y")
    $currentTime = (Get-Date -UFormat "%T")
    $logOutput = $logOutput -join (" ")
    "[$currentDate $currentTime] [$LevelText] $logOutput" | Out-File $Path -Append
}
$MissingOemInfNics = @()
$MissingInfNics = @()
$InfPresent = @()
try {
    $NICObjs = @()
    $activeNICs = get-wmiobject win32_networkadapter -ErrorAction Stop -filter "netconnectionstatus = 2"
    foreach ( $activeNIC in $activeNICs ) {
        Write-Log -Level Info -LogOutput ("Connected NIC : {0}" -f $activeNIC.name)
        $NICObjs += get-wmiobject win32_PnPSignedDriver -ErrorAction Stop |
        where { $_.DeviceName -eq $activeNIC.name -or
            $_.Description -eq $activeNIC.name -or
            $_.FriendlyName -eq $activeNIC.name }
    }
}
catch {
    Write-Error "Exception : $_.Exception.Message"
    Write-Log -Level Error -LogOutput ("WMI Exception : {0}" -f $_.Exception.Message )
    Exit
}
foreach ($nic in $NICObjs ) {
    if ($nic.InfName) {
        Write-Log -Level Info -LogOutput ("Retrieved INF name [{0}] from device manager for {1}." -f $nic.InfName, $nic.DeviceName )
        if ( $nic.InfName -match "^oe[m]") {
            $infPath = Join-Path $env:windir "\inf\$($nic.InfName)"
            if (Test-Path  $infPath) {
                Write-Log -Level Info -LogOutput ("Checked whether the INF file [{0}] exists in the location C:\windows\Inf for {1}. The INF file is present." -f $nic.InfName, $nic.DeviceName )
                $InfPresent += New-Object -TypeName PSObject -Property @{
                    "Network Adapter Name" = $nic.DeviceName
                    "INF Name"             = $($nic.InfName)
                    "INF Status"           = "Present"
                }
            }
            Else {
                $MissingInfNics += New-Object -TypeName PSObject -Property @{
                    "Network Adapter Name" = $nic.DeviceName
                    "INF Name"             = $($nic.InfName)
                    "HardwareID"           = $($nic.HardWareID)
                    "INF Status"           = "Missing"
                    "ServiceName"          = $(( $ActiveNICs | ? { $_.Name -eq $nic.DeviceName -or $_.Name -eq $nic.FriendlyName } ).ServiceName )
                }
                Write-Log -Level Info -LogOutput ("Checked whether the INF file [{0}] exists in the location C:\windows\Inf  for {1}. The INF file is not present." -f $nic.InfName, $nic.DeviceName )
            }
        }
        Else {
            $infPath = Join-Path $env:windir "\inf\$($nic.InfName)"
            if (Test-Path  $infPath) {
                Write-Log -Level Info -LogOutput ("Checked whether the INF file [{0}] exists in the location C:\windows\Inf  for {1}. The INF file is present." -f $nic.InfName, $nic.DeviceName )
                $InfPresent += New-Object -TypeName PSObject -Property @{
                    "Network Adapter Name" = $nic.DeviceName
                    "INF Name"             = $($nic.InfName)
                    "INF Status"           = "Present"
                }
            }
            Else {
                $MissingInfNics += New-Object -TypeName PSObject -Property @{
                    "Network Adapter Name" = $nic.DeviceName
                    "INF Name"             = $($nic.InfName)
                    "HardwareID"           = $($nic.HardWareID)
                    "INF Status"           = "Missing"
                    "ServiceName"          = $(( $ActiveNICs | ? { $_.Name -eq $nic.DeviceName -or $_.Name -eq $nic.FriendlyName } ).ServiceName )
                }
                Write-Log -Level Info -LogOutput ("Checked whether the INF file [{0}] exists in the location C:\windows\Inf  for {1}. The INF file is not present." -f $nic.InfName, $nic.DeviceName )
            }
        }

    }
    Else {
        Write-output "INF Name registered in the Device Manager not found for Network Adapter $($nic.DeviceName)"
        Write-Log -Level Info -LogOutput ("INF name not found in device manager for : {0}" -f $($nic.DeviceName))

    }
}

Function Get-DriverDetails($ServiceName) {
    $ServicePath = "HKLM:\\SYSTEM\CurrentControlSet\Services\" + $ServiceName
    $driverName = Split-Path -path $((Get-ItemProperty -Path $ServicePath -Name ImagePath).ImagePath) -leaf
    try{
        $driverVersion = (Get-Item $ENV:SystemRoot\System32\drivers\$driverName -EA Stop).VersionInfo.FileVersion
    }catch{
       if ($_.Exception.Message -like "*Cannot find path*" -or $_.Exception.Message -like "*Could not find item*"){
           $driverVersion = "Error : Unable to find $driverName in System driver store. Driver re-installation recommended"
       }else{$driverVersion = "Error : $_.Exception.Message" }
    }
    return $driverName, $driverVersion
}
Function SearchINF($driverName, $driverVersion, $missingInf, $hardWareID) {
    $DriverStore = $Env:SystemRoot + "\System32\DriverStore\FileRepository"
    $driverFiles = Get-ChildItem -Path $DriverStore -Filter $driverName -Recurse | % { $_.FullName }
    if ($driverFiles) {
        foreach ($file in $driverFiles) {
            if ($(get-item $file).VersionInfo.FileVersion -eq $driverVersion ) {

                $inf = $file -replace (".{3}$", "inf")
                if ( Test-path $inf ) {
                    $infFile = $inf
                    break
                }
                if ((Test-path $($file -replace ( $(split-path -path $file -leaf), $missingInf) ))) {
                    $inf = $file -replace ( $(split-path -path $file -leaf), $missingInf)
                    if ( Test-path $inf ) {
                        $infFile = $inf
                        break
                    }
                }
                If ($hardWareID) {
                    $driverPath = split-path -Path $file
                    $infs = Get-Item -Path $driverPath\*.inf
                    foreach ($inf in $infs) {
                        if (Select-String -Path $inf -Pattern ([regex]::escape($hardWareID)) -Quiet) {
                            $infFile = $inf
                            break
                        }
                    }
                }
                if ($driverName) {
                    $driverPath = split-path -Path $file
                    $infs = Get-Item -Path $driverPath\*.inf
                    foreach ($inf in $infs) {
                        if (Select-String -Path $inf -SimpleMatch $driverName -Quiet) {
                            $infFile = $inf
                            break
                        }
                    }
                }
            }
        }
    }
    return $infFile
}
if ($MissingInfNics) {
    foreach ( $MissingInfNic in $MissingInfNics) {
        $curDriverFile, $curDriverVersion = Get-DriverDetails -ServiceName $($MissingInfNic.ServiceName)
        if ($curDriverVersion -match "Error") {
        
            Write-Output "$curDriverVersion for $($MissingInfNic."Network Adapter Name")"
            continue
        }
        Write-Log -Level Info -LogOutput ("Retrieved driver details for missing inf [{0}]. Driver : {1}, DriverVersion: {2}, Hardware ID: {3}" -f $MissingInfNic."INF Name", $curDriverFile, $curDriverVersion, $MissingInfNic.HardWareID)
        if ( $curDriverFile -and $curDriverVersion) {
            $driverStoreInfFile = SearchINF -driverName $curDriverFile -driverVersion $curDriverVersion -missingInf $($MissingInfNic."INF Name") -hardWareID $($MissingInfNic.HardWareID)
            $curDriverInfPath = $Env:SystemRoot + "\INF\$($MissingInfNic."INF Name")"
            if ($driverStoreInfFile) {
                Write-Log -Level Info -LogOutput ("Searched for missing INF file [{0}] in Driver Store. Found INF file [$(split-path -path $driverStoreInfFile -leaf)] in {1} for {2}." -f $($MissingInfNic."INF Name"), $(split-path -path $driverStoreInfFile), $MissingInfNic."Network Adapter Name" )
                try {
                    Copy-Item -Path $driverStoreInfFile -Destination $curDriverInfPath -force -ErrorAction Stop
                    Write-Output "Missing $($MissingInfNic."INF Name") file copied from Driver Store to $curDriverInfPath for $($MissingInfNic."Network Adapter Name")"
                    Write-Log -Level Info -LogOutput ("Copied INF file [$(split-path -path $driverStoreInfFile -leaf)] from {1} to C:\windows\INF\{0} for {2}." -f $($MissingInfNic."INF Name"), $(split-path -path $driverStoreInfFile), $MissingInfNic."Network Adapter Name" )
                }
                catch {
                    Write-Output "Error occured while copying INF file. $(split-path -path $driverStoreInfFile -leaf) for $($MissingInfNic."Network Adapter Name")`n  $_.Exception.Message"
                    Write-Log -Level Error -LogOutput ("Failed to copy INF file [$(split-path -path $driverStoreInfFile -leaf)] from {1} to C:\windows\INF\{0} for {2}.`n Exception : $_.Exception.Message" -f $($MissingInfNic."INF Name"), $(split-path -path $driverStoreInfFile), $MissingInfNic."Network Adapter Name" )

                }
            }
            Else {
                Write-Output "Missing $($MissingInfNic."INF Name") file not found in Driver Store for $($MissingInfNic."Network Adapter Name"). Driver re-installation recommended."
                Write-Log -Level Info -LogOutput ("Searched for missing INF file [{0}]. Not found INF file in Driver Store for {1}." -f $($MissingInfNic."INF Name"), $MissingInfNic."Network Adapter Name" )

            }
        }
    }
}
Else {
    Write-Output "INF file(s) already present on the system. No action needed."
    Write-Log -Level Info -LogOutput ("INF file not missing on the system hence no action is needed.")
}
$InfPresent + $MissingInfNics | fl "Network Adapter Name", "INF Name", "INF Status"
