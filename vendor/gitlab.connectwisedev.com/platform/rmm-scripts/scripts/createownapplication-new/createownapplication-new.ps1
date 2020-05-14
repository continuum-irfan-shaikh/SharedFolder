if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

#region Functions
function Get-MD5Hash ($path) {
    $md5 = New-Object -TypeName System.Security.Cryptography.MD5CryptoServiceProvider
    return [System.BitConverter]::ToString($md5.ComputeHash([System.IO.File]::ReadAllBytes($path))) -replace '-', ''
}

function Download-FromURL ($URL, $LocalFilePath) {
    try {
        # Write-Host "Downloading: $URL"
        $WebClient = New-Object System.Net.WebClient
        $WebClient.DownloadFile($URL, $LocalFilePath)
        if (-not(Test-Path $LocalFilePath)) { Write-Error "Unable to download the file from URL: $URL" }
    }
    catch {
        Write-Error "Unable to download the file from URL: $URL`n$($_.Exception.Message)"
    }
}

function Download-FromFTP ($RemoteFTPFilePath, $LocalFilePath, $UserName, $Password) {
    try {
        # Write-Host "Downloading: $RemoteFTPFilePath"
        $ftprequest = [System.Net.FtpWebRequest]::create($RemoteFTPFilePath)
        $ftprequest.Credentials = New-Object System.Net.NetworkCredential($username, $password)
        $ftprequest.Method = [System.Net.WebRequestMethods+Ftp]::DownloadFile
        $ftprequest.UseBinary = $true
        $ftprequest.KeepAlive = $false
        $ftpresponse = $ftprequest.GetResponse()
        $responsestream = $ftpresponse.GetResponseStream()
        $targetfile = New-Object IO.FileStream ($LocalFilePath, [IO.FileMode]::Create)
        [byte[]]$readbuffer = New-Object byte[] 1024

        do {
            $readlength = $responsestream.Read($readbuffer, 0, 1024)
            $targetfile.Write($readbuffer, 0, $readlength)
        }
        while ($readlength -ne 0)
        $targetfile.close()

        if (-not (Test-Path $LocalFilePath)) { Write-Error "Unable to download the file from FTP location: $RemoteFTPFilePath" }

    }
    catch {
        Write-Error "Unable to download the file from FTP location: $RemoteFTPFilePath`n$($_.Exception.Message)"
    }
}

function Download-FromUNC ($RemoteUNCFilePath, $PackageName, $LocalFilePath, $UserName, $Password) {
    try {
        # Write-Host "Downloading: $(Join-Path $RemoteUNCFilePath $PackageName)"
        $FreeDriveLetter = Get-RandomFreeDriveLetter
        MapNetworkDrive -DriveLetter $FreeDriveLetter -Directory $RemoteUNCFilePath -Persistent $false -Username $UserName -Password $Password
        Copy-Item $(Join-Path $FreeDriveLetter $PackageName) -Destination $LocalFilePath
        if (-not (Test-Path $LocalFilePath)) { Write-Error "Unable to download the file from Network location: $RemoteUNCFilePath" }
    }
    catch {
        Write-Error "Unable to download the file from Network location: $RemoteUNCFilePath`n$($_.Exception.Message)"
    }
    finally {
        try {
            $Network = New-Object -ComObject WScript.Network
            $Network.RemoveNetworkDrive($FreeDriveLetter, $True)
        }
        catch {
            # no action
        }
    }
}

Function MapNetworkDrive($DriveLetter, $Directory, $Persistent, $Username, $Password) {
    try {
        if (Test-Path $DriveLetter) {
            # Write-Output "`n[Map Network Drive] Failed to Map the network drive: `'$DriveLetter`' because it already exists."
            Write-Error "Failed to Map the network drive: `'$DriveLetter`' because it already exists."
        }
        else {
            $Network = New-Object -ComObject WScript.Network
            $Network.MapNetworkDrive($DriveLetter, "$Directory", $Persistent, $UserName, $Password)
        }
    }
    catch {
        Write-Error "Failed to Map the network drive.`n$($_.Exception.Message)"
    }
}

function Join-URLParts {
    param ([string[]] $Parts, [string] $Seperator = '')
    $search = '(?<!:)' + [regex]::Escape($Seperator) + '+'  #Replace multiples except in front of a colon for URLs.
    $replace = $Seperator
    ($Parts | Where-Object { $_ -and $_.Trim().Length }) -join $Seperator -replace $search, $replace
}

function Get-RandomFreeDriveLetter {
    Get-ChildItem function:[d-z]: -n | Where-Object { !(test-path $_) } | Get-Random
}


function Write-SuccessOrFail ([System.Diagnostics.Process]$Process, $PackageName, $Action) {
    If ($process.exitcode -eq 0) {
        Write-Output "`n$Action of '$PackageName' Successful."
    }
    else {
        Write-Output "`n$Action of '$PackageName' Failed. Exitcode: $($process.exitcode)"
    }
}

function Test-Registry ($Key, $Value, $Data) {
    $Path = "REGISTRY::$key"
    $KeyExists = Test-Path $Path
    $Result = Get-ItemProperty "REGISTRY::$key" -Name $Value -ErrorAction SilentlyContinue
    $ValueExists = [bool](($Result).$value)
    $DataMatches = $data -eq ($Result).$value

    return $KeyExists, $ValueExists, $DataMatches
}

function Test-FileVersion ($Path, [System.Version]$UserSpecifiedVersion) {
    $FileExists = Test-Path $Path
    $versioninfo = (Get-Item $path -ErrorAction SilentlyContinue).versioninfo
    $FileVersion = [System.Version] ("{0}.{1}.{2}.{3}" -f $versioninfo.FileMajorPart, $versioninfo.FileMinorPart, $versioninfo.FileBuildPart, $versioninfo.FilePrivatePart)
    $FileVersionMatches = $FileVersion -eq $UserSpecifiedVersion
    $FileVersionGreater = $FileVersion -gt $UserSpecifiedVersion
    $FileVersionLower = $FileVersion -lt $UserSpecifiedVersion

    $FileExists, $FileVersionMatches, $FileVersionGreater, $FileVersionLower
}

Function Test-AddRemoveProgram ($Program) {
    $Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall', 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'
    $product = Get-ChildItem $Registry -ErrorAction 'SilentlyContinue' | Get-ItemProperty | Where-Object { $_.Displayname -eq $program }
    return [bool] $product
}
#endregion Functions

#region Main
$ErrorActionPreference = 'Stop'

try {
    #region PreCheck
    if ($PerformPreCheck) {
        $Name = 'Pre-Check'
        Switch ($PreCheckType) {
            'Check if file exists' {
                $FileExists = Test-Path $PreCheckFileName
                $PreCheckAction = $PreCheckAction_for_PreCheckType_1  # handle Narayan's request to have multiple $PreCheckAction_for_PreCheckType_[n] variables substituted through JSON
                Switch ($PreCheckAction) {
                    'If file exists then Continue' { if ($FileExists) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If file does not exists then Continue' { if (-not $FileExists) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If file exists then Abort' { if ($FileExists) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If file does not exists then Abort' { if (-not $FileExists) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                }
                break;
            }
            'Check if registry information exists' {
                $KeyExists, $ValueMatches, $DataMatches = Test-Registry -Key $PreCheckKey -Value $PreCheckValue -Data $PreCheckData
                $PreCheckAction = $PreCheckAction_for_PreCheckType_2 # handle Narayan's request to have multiple $PreCheckAction_for_PreCheckType_[n] variables substituted through JSON
                Switch ($PreCheckAction) {
                    'If key exists then Continue' { if ($KeyExists) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If key does not exists then Continue' { if (-not $KeyExists) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If key exists then Abort' { if ($KeyExists) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If key not exists then Abort' { if (-not $KeyExists) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If key exists and value matches then Continue' { if ($KeyExists -and $ValueMatches) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If key exists and value does not matches then Continue' { if ($KeyExists -and (-not $ValueMatches)) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If key exists and value matches then Abort' { if ($KeyExists -and $ValueMatches) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If key exists and value does not matches then Abort' { if ($KeyExists -and (-not $ValueMatches)) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If key exists and value exists and data matches then Continue' { if ($KeyExists -and $datamatches) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If key exists and value exists and data does not matches then Continue' { if ($KeyExists -and $ValueMatches -and (-not $datamatches)) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If key exists and value exists and data matches then Abort' { if ($KeyExists -and $ValueMatches -and $datamatches) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If key exists and value exists and data does not matches then Abort' { if ($KeyExists -and $ValueMatches -and (-not $datamatches)) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                }
                break;
            }
            'Check if Add / Remove program exists' {
                $ProgramExists = Test-AddRemoveProgram -Program $PreCheckProgram
                $PreCheckAction = $PreCheckAction_for_PreCheckType_3 # handle Narayan's request to have multiple $PreCheckAction_for_PreCheckType_[n] variables substituted through JSON

                Switch ($PreCheckAction) {
                    'If the program exists then continue' { if ($ProgramExists) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If the program does not exists then continue' { if (!$ProgramExists) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If the program exists then abort' { if ($ProgramExists) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If the program does not exists then abort' { if (!$ProgramExists) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                }
                break;
            }
            'Check for file version' {
                $FileExists, $FileVersionMatches, $FileVersionGreater, $FileVersionLower = Test-FileVersion -Path $PreCheckFileName -UserSpecifiedVersion $PreCheckVersion
                $PreCheckAction = $PreCheckAction_for_PreCheckType_4 # handle Narayan's request to have multiple $PreCheckAction_for_PreCheckType_[n] variables substituted through JSON

                Switch ($PreCheckAction) {
                    'If file exists and version equal to specified then Continue' { if ($FileExists -and $FileVersionMatches) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If file exists and version not equal to specified then Continue' { if ($FileExists -and (-not $FileVersionMatches)) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If file exists and version higher to specified then Continue' { if ($FileExists -and $FileVersionGreater) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If file exists and version lower to specified then Continue' { if ($FileExists -and $FileVersionLower) { continue } else { Write-Output "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If file exists and version equal to specified then Abort' { if ($FileExists -and $FileVersionMatches) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If file exists and version not equal to specified then Abort' { if ($FileExists -and (-not $FileVersionMatches)) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If file exists and version higher to specified then Abort' { if ($FileExists -and $FileVersionGreater) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If file exists and version lower to specified then Abort' { if ($FileExists -and $FileVersionLower) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If file does not exists then Abort' { if (-not $FileExists) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                }
                break;
            }
        }
    }
    #endregion Precheck

    #region Main
    Switch ($TypeOfInstaller) {
        'Installer' {
            $PackageLocalPath = Join-Path $env:TEMP $PackageName

            #region InputValidation
            switch ($LocationType) {
                'ftp' {
                    $location = $ftpPath + "/" + $PackageName
                    $Scheme = ([uri]$Location).scheme
                    if ($Scheme -and ($Scheme -ne 'ftp')) {
                        Write-Error "Invalid `'$LocationType`' URL: $Location"
                    }
                    elseif (!($Location.StartsWith("ftp://"))) {
                        $Location = "ftp://" + $Location
                    }
                    break;
                }
                'http' {
                    $location = $httpPath + "/" + $PackageName
                    $Scheme = ([uri]$Location).scheme
                    if ($Scheme -and ($Scheme -ne 'http')) {
                        Write-Error "Invalid `'$LocationType`' URL: $Location"
                    }
                    elseif (!($Location.StartsWith("http://"))) {
                        $Location = "http://" + $Location
                    }
                    break;
                }
                'https' {
                    $location = $httpsPath + "/" + $PackageName
                    $Scheme = ([uri]$Location).scheme
                    if ($Scheme -and ($Scheme -ne 'https')) {
                        Write-Error "Invalid `'$LocationType`' URL: $Location"
                    }
                    elseif (!($Location.StartsWith("https://"))) {
                        $Location = "https://" + $Location
                    }
                    break;
                }
                'network' {
                    $location = $networkPath + "/" + $PackageName
                    $Scheme = ([uri]$Location).scheme
                    $isUNC = ([uri]$Location).isUNC

                    if ($Scheme -and ($Scheme -ne 'file')) {
                        # condition to check wrong scheme, like 'http' in case of 'ftp'
                        Write-Error "Invalid `'$LocationType`' Path: $Location"
                    }
                    elseif ($Scheme -and ($Scheme -eq 'file') -and !$isUNC) {
                        # has the correct scheme, but a local path
                        Write-Error "Invalid `'$LocationType`' Path: $Location"
                    }
                    elseif (!($Location.StartsWith("\\"))) {
                        $Location = "\\" + $Location
                    }
                    break;
                }
                'local' {
                    $location = $localPath

                    $Scheme = ([uri]$Location).scheme
                    $isUNC = ([uri]$Location).isUNC

                    if ($Scheme -and ($Scheme -ne 'file')) {
                        # condition to check wrong scheme, like 'http' in case of 'ftp'
                        Write-Error "Invalid `'$LocationType`' Path: $Location"
                    }
                    elseif ($Scheme -and ($Scheme -eq 'file') -and $isUNC) {
                        # has the correct scheme, but a network path
                        Write-Error "Invalid `'$LocationType`' Path: $Location"
                    }

                    if (!(Test-Path $(Join-Path $Location $PackageName))) {
                        Write-Error "Local path not found: $(Join-Path $Location $PackageName)"
                    }

                    break;
                }
            }
            #endregion InputValidation

            #region DownloadPackage
            switch (([uri]$Location).scheme) {
                'ftp' {
                    $RemoteFTPFilePath = Join-URLParts ($Location, $PackageName) -Seperator '/'
                    Download-FromFTP $RemoteFTPFilePath $PackageLocalPath $Username $Password
                    Break;
                }
                'http' { Download-FromURL $Location $PackageLocalPath ; Break; }
                'https' {
                    try{
                        [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType] 'Ssl3, Tls, Tls11, Tls12'
                    }
                    catch {
                        $TLSEnforce = $false
                    }
                    try{
                        Download-FromURL $Location $PackageLocalPath ; Break;
                    }
                    catch {
                        if(!$TLSEnforce){
                            Write-Output "Download from the URL has failed. It might be an issue while establishing SSL\TLS secure channel between client and the server, as this machine only supports following version of secure protocols: $([enum]::GetNames([System.Net.SecurityProtocolType]) -join ', ') through the Task."
                            Write-Output "Workaround: Try downloading the file from HTTP variant of the URL, if it exists. Or copy the exe\msi on a UNC (Network Share) path and then execute the setup using this script"
                        }
                        else {
                            Write-Error $_.Exception.Message
                        }
                    }
                }
                'file' {
                    if (([uri]$Location).IsUnc) {
                        # if the location is a shared UNC path
                        Download-FromUNC $Location $PackageName $PackageLocalPath $UserName $Password
                    }
                    else {
                        # if the location is a local system path
                        $PackageLocalPath = Join-Path $Location $PackageName
                    }
                    Break;
                }
                default {
                    Write-Output "Failed. Not a valid user input for package 'Location', please validate the input and try again."
                }
            }
            #endregion DownloadPackage

            # validate checksum of donwloaded file
            if ($MD5CheckSum -and $(Get-MD5Hash $PackageLocalPath) -ne $MD5CheckSum) {
                Write-Error "MD5 Checksum: '$MD5CheckSum' did not match with the MD5 Checksum of Package:'$PackageName'"
            }

            #region Install

            $InstallExecuteMode = "Install" # as per Narayan's request hardcoded this value for now, and will implement others once we have a clarity on this
            Switch ($InstallExecuteMode){
                'Install' {
                    switch ($TypeOfApplication) {
                        'exe' {
                            # Write-Host "Installing: $PackageLocalPath $InstallationParameter"
                            if (![string]::IsNullOrEmpty($InstallationParameter)) {
                                $InstallationParameter = $InstallationParameter.trim()
                                $Process = Start-Process -FilePath $PackageLocalPath -ArgumentList $InstallationParameter -Wait -PassThru
                            }
                            else {
                                $Process = Start-Process -FilePath $PackageLocalPath -Wait -PassThru
                            }

                            if (!$PerformPostCheck) {
                                Write-SuccessOrFail -Process $Process -PackageName $PackageName -Action 'Installation'
                            }
                        }
                        'msi' {
                            # Write-Host "Installing: $env:systemroot\system32\msiexec.exe /i $PackageLocalPath $InstallationParameter"
                            if (![string]::IsNullOrEmpty($InstallationParameter)) {
                                $MSIArgs = @(
                                    "/i"
                                    $PackageLocalPath
                                    $InstallationParameter.Split(" ") | ForEach-Object {$_.trim()}
                                )
                            }
                            else {
                                $MSIArgs = @(
                                    "/i"
                                    $PackageLocalPath
                                )
                            }

                            $Process = Start-Process -FilePath "$env:systemroot\system32\msiexec.exe" -ArgumentList $MSIArgs -Wait -PassThru

                            if (!$PerformPostCheck) {
                                Write-SuccessOrFail -Process $Process -PackageName $PackageName -Action 'Installation'
                            }
                        }
                    }
                }
                'Install From Gateway' {
                    # no action defined for now
                }
                'Download to Gateway' {
                    # no action defined for now
                }
            }

            #endregion Install
            break;
        }

        'UnInstaller' {
            #region UnInstall
            $UninstallExecuteMode = "Uninstall" # as per Narayan's request hardcoded this value for now, and will implement others once we have a clarity on this
            Switch ($UninstallExecuteMode){
                'Uninstall' {
                    switch ($TypeOfApplication) {
                        'exe' {
                            # Write-Host "UnInstalling: $UninstallationPath $UnInstallationParameter"
                            if ($UnInstallationParameter) {
                                $Process = Start-Process -FilePath $UninstallationPath -ArgumentList "$UnInstallationParameter" -Wait -PassThru
                            }
                            else {
                                $Process = Start-Process -FilePath $UninstallationPath -Wait -PassThru
                            }
                            if (!$PerformPostCheck) {
                                Write-SuccessOrFail -Process $Process -PackageName $UninstallationPath -Action 'UnInstallation'
                            }
                        }
                        'msi' {
                            $Msg = $UninstallationParameter
                            $MSIArgs = ($UninstallationParameter -Replace "msiexec.exe", "" -replace "msiexec", "").trim()
                            # Write-Host "UnInstalling: $env:systemroot\system32\msiexec.exe /i $UnInstallationParameter"

                            $Process = Start-Process -FilePath "$env:systemroot\system32\msiexec.exe" -ArgumentList $MSIArgs -Wait -PassThru
                            if (!$PerformPostCheck) {
                                Write-SuccessOrFail -Process $Process -PackageName $Msg -Action 'UnInstallation'
                            }
                        }
                    }
                }
            }

            #endregion UnInstall
            break;
        }
    }
    #endregion Main

    #region PostCheck
    if ($PerformPostCheck) {
        $Name = 'Post-Check'
        Switch ($PostCheckType) {
            'Check if file exists' {
                $FileExists = Test-Path $PostCheckFileName
                $PostCheckAction = $PostCheckAction_for_PostCheckType_1  # handle Narayan's request to have multiple $PostCheckAction_for_PostCheckType_[n] variables substituted through JSON
                Switch ($PostCheckAction) {
                    'If the file exists then mark Success' { if ($FileExists) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the file does not exists then mark Success' { if (-not $FileExists) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the file exists then mark Fail' { if ($FileExists) { Write-Output "Failed" } else { Write-Output "Success" } }
                    'If the file does not exists then mark Fail' { if (-not $FileExists) { Write-Output "Failed" } else { Write-Output "Success" } }
                }
                break;
            }
            'Check if registry information exists' {
                $KeyExists, $ValueMatches, $DataMatches = Test-Registry -Key $PostCheckKey -Value $PostCheckValue -Data $PostCheckData
                $PostCheckAction = $PostCheckAction_for_PostCheckType_2  # handle Narayan's request to have multiple $PostCheckAction_for_PostCheckType_[n] variables substituted through JSON

                Switch ($PostCheckAction) {
                    'If the key exists then mark success' { if ($KeyExists) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the key does not exists then mark success' { if (-not $KeyExists) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the key exists then mark fail' { if ($KeyExists) { Write-Output "Failed" } else { Write-Output "Success" } }
                    'If the key not exists then mark fail' { if (-not $KeyExists) { Write-Output "Failed" } else { Write-Output "Success" } }
                    'If the key exists and value matches then mark success' { if ($KeyExists -and $ValueMatches) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the key exists and value does not match then mark success' { if ($KeyExists -and (-not $ValueMatches)) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the key exists and value matches then mark fail' { if ($KeyExists -and $ValueMatches) { Write-Output "Failed" } else { Write-Output "Success" } }
                    'If the key exists and value does not match then mark fail' { if ($KeyExists -and (-not $ValueMatches)) { Write-Output "Failed" } else { Write-Output "Success" } }
                    'If the key exists and value exists and data matches then mark success' { if ($KeyExists -and $ValueMatches -and $DataMatches) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the key exists and value exists and data does not match then mark success' { if ($KeyExists -and $ValueMatches -and (-not $DataMatches)) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the key exists and value exists and data matches then mark fail' { if ($KeyExists -and $ValueMatches -and $DataMatches) { Write-Output "Failed" } else { Write-Output "Success" } }
                    'If the key exists and value exists and data does not match then mark fail' { if ($KeyExists -and $ValueMatches -and (-not $DataMatches)) { Write-Output "Failed" } else { Write-Output "Success" } }
                }
                break;
            }
            'Check if Add / Remove program exists' {
                $ProgramExists = Test-AddRemoveProgram -Program $PostCheckProgram
                $PostCheckAction = $PostCheckAction_for_PostCheckType_3  # handle Narayan's request to have multiple $PostCheckAction_for_PostCheckType_[n] variables substituted through JSON

                Switch ($PostCheckAction) {
                    'If the program exists then mark success' { if ($ProgramExists) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the program does not exists then mark success' { if (-not $ProgramExists) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the program exists then mark fail' { if ($ProgramExists) { Write-Output "Failed" } else { Write-Output "Success" } }
                    'If the program does not exists then mark fail' { if (-not $ProgramExists) { Write-Output "Failed" } else { Write-Output "Success" } }
                }
                break;
            }
            'Check for file version' {
                $FileExists, $FileVersionMatches, $FileVersionGreater, $FileVersionLower = Test-FileVersion -Path $PostCheckFileName -UserSpecifiedVersion $PostCheckVersion
                $PostCheckAction = $PostCheckAction_for_PostCheckType_4  # handle Narayan's request to have multiple $PostCheckAction_for_PostCheckType_[n] variables substituted through JSON

                Switch ($PostCheckAction) {
                    'If the file exists and the version matches then mark success' { if ($FileExists -and $FileVersionMatches) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the file exists and version does not match then mark success' { if ($FileExists -and (-not $FileVersionMatches)) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the file exists and version is greater than specified then mark success' { if ($FileExists -and $FileVersionGreater) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the file exists and version is lower than specified then mark success' { if ($FileExists -and $FileVersionLower) { Write-Output "Success" } else { Write-Output "Failed" } }
                    'If the file exists and the version is matches then mark fail' { if ($FileExists -and $FileVersionMatches) { Write-Output "Failed" } else { Write-Output "Success" } }
                    'If the file exists and version does not match then mark fail' { if ($FileExists -and (-not $FileVersionMatches)) { Write-Output "Failed" } else { Write-Output "Success" } }
                    'If the file exists and version is greater than specified then mark fail' { if ($FileExists -and $FileVersionGreater) { Write-Output "Failed" } else { Write-Output "Success" } }
                    'If the file exists and version is lower than specified then mark fail' { if ($FileExists -and $FileVersionLower) { Write-Output "Failed" } else { Write-Output "Success" } }
                    'If the file does not exists then mark fail' { if (-not $FileExists) { Write-Output "Failed" } else { Write-Output "Success" } }
                }
                break;
            }
        }
    }
    #endregion Postcheck
}
catch {
    Write-Error $_.exception.message
}
finally {
    try {
        # to remove donwloaded files from local machine if required
        Remove-Item $PackageLocalPath -ErrorAction SilentlyContinue
    }
    catch {
        # no action
    }
}

