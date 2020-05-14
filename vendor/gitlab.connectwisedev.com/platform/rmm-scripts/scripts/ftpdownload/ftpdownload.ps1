# region InputVariables
# $ConnectionPath = 'ftp://10.2.130.27'   # String 
# $Username = "FTPUser"   # String  
# $Password = "test@123"   # String 
# $Type = 'File'  # String
# $Source = 'ChromeSetup.exe'   # String 
# $Recurse = $true   # Boolean
# $Destination = 'C:\partner\'   # String
# $CreateDestinationFolder = $True # Boolean
# $PerformPreCheck = $false   # Boolean
# $PreCheckType = 'Check if file exists'   # String
# $PreCheckAction = 'If file does not exists then Continue'   # String
# $PreCheckFileName = 'C:\x\csetup.exe'   # String
# $PreCheckVersion = '7.0.4.3'   # String
# $PreCheckMD5Checksum = 'D41D8CD98F00B204E9800998ECF8427Eu'   # String
# $PreCheckKey = 'HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\DB Browser for SQLite'   # String
# $PreCheckValue = 'nomodify'   # String
# $PreCheckData = '1'   # String
# $PerformPostCheck = $false   # Boolean
# $PostCheckFileName = 'C:\x\csetupd.exe'   # String
# $PostCheckMD5Checksum = '907427099DE2E504BB27BF4CED5EDC5E'   # String
# $PostCheckAction = 'If file exists and MD5 checksum matches then Success'   # String
#endregion InputVariables

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
Function ListFTPItems { 
    Param(
        [String]$Path = "",
        $Credentials,
        [Switch]$Recurse,    
        [Int]$Depth = 0,        
        [String]$Filter = "*"
    )
    
    Process {
        if (!($Path -match "ftp://")) { $Path = "ftp://" + $Path }
        $Path = [regex]::Replace($Path, '/$', '')
        $Path = [regex]::Replace($Path, '/+', '/')
        $Path = [regex]::Replace($Path, '^ftp:/', 'ftp://')

        if ($Depth -gt 0) {
            $CurrentDepth = [regex]::matches($Path, "/").count
            if ((Get-Variable -Scope Script -Name MaxDepth -ErrorAction SilentlyContinue) -eq $null) {
                New-Variable -Scope Script -Name MaxDepth -Value ([Int]$CurrentDepth + $Depth)
            }
            $Recurse = $true
        }
  
        if ((GetFTPItemSize $Path -Credentials $Credentials -Silent) -eq -1) {
            $ParentPath = $Path
        }
        else {
            $LastIndex = $Path.LastIndexOf("/")
            $ParentPath = $Path.SubString(0, $LastIndex)
        }
                        
        [System.Net.FtpWebRequest]$Request = [System.Net.WebRequest]::Create($Path)
        $Request.Credentials = $Credentials
            
        $Request.Method = [System.Net.WebRequestMethods+FTP]::ListDirectoryDetails
        Try {
            $Mode = "Unknown"
            $Response = $Request.GetResponse()
            [System.IO.StreamReader]$Stream = New-Object System.IO.StreamReader($Response.GetResponseStream(), [System.Text.Encoding]::Default)
            $DirectoryList = @()
            $ItemsCollection = @()
            Try {
                [string]$Line = $Stream.ReadLine()
            }
            Catch {
                $Line = $null
            }
                
            While ($Line) {
                if ($Mode -eq "Compatible" -or $Mode -eq "Unknown") {
                    $null, [string]$IsDirectory, [string]$Flag, [string]$Link, [string]$UserName, [string]$GroupName, [string]$Size, [string]$Date, [string]$Name = [regex]::split($Line, '^([d-])([rwxt-]{9})\s+(\d{1,})\s+([.@A-Za-z0-9-]+)\s+([A-Za-z0-9-]+)\s+(\d{1,})\s+(\w+\s+\d{1,2}\s+\d{1,2}:?\d{2})\s+(.+?)\s?$', "SingleLine,IgnoreCase,IgnorePatternWhitespace")

                    if ($IsDirectory -eq "" -and $Mode -eq "Unknown") {
                        $Mode = "IIS6"
                    }
                    elseif ($Mode -ne "Compatible") {
                        $Mode = "Compatible" #IIS7/Linux
                    }

                    #Write-Host $Mode -ForegroundColor Red
                        
                    if ($Mode -eq "Compatible") {
                        $DatePart = $Date -split "\s+"
                        $NewDateString = "$($DatePart[0]) $('{0:D2}' -f [int]$DatePart[1]) $($DatePart[2])"
                            
                        Try {
                            if ($DatePart[2] -match ":") {
                                $Month = ([DateTime]::ParseExact($DatePart[0], "MMM", [System.Globalization.CultureInfo]::InvariantCulture)).Month
                                if ((Get-Date).Month -ge $Month) {
                                    $NewDate = [DateTime]::ParseExact($NewDateString, "MMM dd HH:mm", [System.Globalization.CultureInfo]::InvariantCulture)
                                }
                                else {
                                    $NewDate = ([DateTime]::ParseExact($NewDateString, "MMM dd HH:mm", [System.Globalization.CultureInfo]::InvariantCulture)).AddYears(-1)
                                }
                            }
                            else {
                                $NewDate = [DateTime]::ParseExact($NewDateString, "MMM dd yyyy", [System.Globalization.CultureInfo]::InvariantCulture)
                            }
                        }
                        Catch {
                        }                            
                    }
                }
                    
                #if ($Mode -eq "IIS6") {
                #    $null, [string]$NewDate, [string]$IsDirectory, [string]$Size, [string]$Name = [regex]::split($Line, '^(\d{2}-\d{2}-\d{2}\s+\d{2}:\d{2}[AP]M)\s+<*([DIR]*)>*\s+(\d*)\s+(.+).*$', "SingleLine,IgnoreCase")
                #}
                if ($Mode -eq "IIS6") {
                    $var = $Line.split(" ") | ForEach-Object { if ($_ -ne "") { return $_ } }
                    [string]$NewDate = $var[0] + " " + $var[1] 
                    if ($var[2] -imatch "DIR") {
                        [string]$IsDirectory = $var[2]
                    }
                    else {
                        [string]$IsDirectory = ""
                        [string]$Size = $var[2]
                    }
                    [string]$Name = $var[3]
                }
                    
                    
                if ($ParentPath -match "\*|\?") {
                    $LastIndex = $ParentPath.LastIndexOf("/")
                    $ParentPath = $ParentPath.SubString(0, $LastIndex)
                    $ParentPath.Trim() + "/" + $Name.Trim()
                }
                    
                $LineObj = New-Object PSObject -Property @{
                    Type         = if ($IsDirectory -eq '') { 'File' }else { 'Directory' }
                    Right        = $Flag
                    Ln           = $Link
                    User         = $UserName
                    Group        = $GroupName
                    Size         = $Size
                    ModifiedDate = $NewDate
                    Name         = $Name.Trim()
                    FullName     = $ParentPath.Trim() + "/" + $Name.Trim()
                    Parent       = $ParentPath.Trim()
                }
                    
                if ($Recurse -and ($LineObj.Type -eq 'Directory')) {
                    $DirectoryList += $LineObj
                }
                    
                if ($LineObj.Type) {
                    if ($LineObj.Name -like $Filter) {
                        $ItemsCollection += $LineObj
                    }
                }
                $Line = $Stream.ReadLine()
            }
                
            $Response.Close()
                
            if ($Recurse -and ($CurrentDepth -lt $Script:MaxDepth -or $Depth -eq 0)) {
                $RecurseResult = @()
                $DirectoryList | ForEach-Object {
                    $RecurseResult += ListFTPItems -Path ($_.FullName) -Credentials $Credentials -Recurse -Filter $Filter -Depth $Depth
                        
                }    

                $ItemsCollection += $RecurseResult
            }    
                
            if ($ItemsCollection.count -eq 0) {
                Return 
            }
            else {
                Return $ItemsCollection | Sort-Object -Property @{Expression = "Parent"; Descending = $false }, @{Expression = "Type"; Descending = $true }, @{Expression = "Name"; Descending = $false } 
            }
        }
        Catch {
            Write-Error $_.Exception.Message -ErrorAction Stop 
        }
        
        if ($CurrentDepth -ge $Script:MaxDepth) {
            Remove-Variable -Scope Script -Name CurrentDepth 
        }        
    }
}

Function GetFTPItemSize {
    Param(
        [parameter(Mandatory = $true)]
        [String]$Path = "",
        [Switch]$Silent = $False,
        $Credentials
    )
    Process {
        if (!($Path -match "ftp://")) {
            $Path = "ftp://" + $Path
        }
        $Path = [regex]::Replace($Path, '/$', '')
        $Path = [regex]::Replace($Path, '/+', '/')
        $Path = [regex]::Replace($Path, '^ftp:/', 'ftp://')
           
        [System.Net.FtpWebRequest]$Request = [System.Net.WebRequest]::Create($Path)
        $Request.Credentials = $Credentials
            
        $Request.Method = [System.Net.WebRequestMethods+FTP]::GetFileSize 
        Try {
            $Response = $Request.GetResponse()
            $Length = $Response.ContentLength
            $Response.Close()
            Return $Length
        }
        Catch {
            if (!$Silent) {
                Write-Error $_.Exception.Message -ErrorAction Stop  
            }    
            Return -1
        }
    }           
}

Function DownloadFtpDirectory {
    Param(
        [String]$Path = "",
        [String]$localPath,
        $Credentials,
        [Switch]$Recurse,    
        [Int]$Depth = 0,        
        [String]$Filter = "*"
    )
    $Result = ListFTPItems -Path $Path -Credentials $credentials -Recurse:$Recurse -Depth $Depth -Filter $Filter
  
    Foreach ($item in $result) {
        $localFilePath = Join-Path $localPath $(($Item.fullname -replace $Path).trim())
        
        # Write-Host "`$localFilePath:  $localFilePath" -fore blue

        if (($Item.Type -eq 'Directory') -and $Recurse) {
            # name with no file extension is considered a directory
            if (!(Test-Path $localFilePath -PathType container)) {
                # Write-Host "Creating directory $localFilePath"
                New-Item $localFilePath -Type directory | Out-Null
            }
            #DownloadFtpDirectory -Path ($Item.FullName + "/")  -Credentials $credentials -localPath $($localFilePath)
            
        }
        elseif ($item.Type -eq 'File') {

            DownloadFTPFile -Path $Item.FullName -localPath $localFilePath -Credentials $Credentials
        }
    }
}

Function DownloadFTPFile {
    Param(
        [String]$Path = "",
        [String]$localPath,
        $Credentials
    )
    # Write-Host "Downloading $Path to $localPath"
    $downloadRequest = [Net.WebRequest]::Create($Path)
    $downloadRequest.Method = [System.Net.WebRequestMethods+Ftp]::DownloadFile
    $downloadRequest.Credentials = $credentials
    
    $downloadResponse = $downloadRequest.GetResponse()
    $sourceStream = $downloadResponse.GetResponseStream()
    $targetStream = [System.IO.File]::Create($localPath)
    $buffer = New-Object byte[] 10240
    while (($read = $sourceStream.Read($buffer, 0, $buffer.Length)) -gt 0) {
        $targetStream.Write($buffer, 0, $read);
    }
    $targetStream.Dispose()
    $sourceStream.Dispose()
}

function Join-URLParts {
    param ([string[]] $Parts, [string] $Seperator = '')
    $search = '(?<!:)' + [regex]::Escape($Seperator) + '+'  #Replace multiples except in front of a colon for URLs.
    $replace = $Seperator
    ($Parts | Where-Object { $_ -and $_.Trim().Length }) -join $Seperator -replace $search, $replace
}

function Get-MD5Hash ($path) {
    $md5 = New-Object -TypeName System.Security.Cryptography.MD5CryptoServiceProvider
    return [System.BitConverter]::ToString($md5.ComputeHash([System.IO.File]::ReadAllBytes($path))) -replace '-', ''
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
    $FileVersion = (Get-Item $Path -ErrorAction SilentlyContinue).versioninfo.fileversion
    $FileVersionMatches = $FileVersion -eq $UserSpecifiedVersion
    $FileVersionGreater = $FileVersion -gt $UserSpecifiedVersion
    $FileVersionLower = $FileVersion -lt $UserSpecifiedVersion

    return $FileExists, $FileVersionMatches, $FileVersionGreater, $FileVersionLower
}
#endregion Functions

try {
    $ErrorActionPreference = 'Stop'

    #region PreCheck
    if ($PerformPreCheck) {
        $Name = 'Pre-Check'
        Switch ($PreCheckType) {
            'Check if file exists' {
                $FileExists = Test-Path $PreCheckFileName
                Switch ($PreCheckAction) {
                    'If file exists then Continue' { if ($FileExists) { continue } else { Write-Error "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If file does not exists then Continue' { if (-not $FileExists) { continue } else { Write-Error "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If file exists then Abort' { if ($FileExists) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If file does not exists then Abort' { if (-not $FileExists) { Write-Error "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If file exists and MD5 checksum matches then Continue' { 
                        if ($PreCheckMD5Checksum) {
                            if ($FileExists -and ((Get-MD5Hash $PreCheckFileName) -eq $PreCheckMD5Checksum)) { continue }
                            else { Write-Error "Aborting because $Name condition was met `'$PreCheckAction`'"; return }
                        }
                        else { Write-Error "Please provide a MD5 Checksum to validate the $Name condition" }
                    }
                    'If file exists and MD5 checksum does not matches then Continue' { 
                        if ($PreCheckMD5Checksum) {
                            if ($FileExists -and !((Get-MD5Hash $PreCheckFileName) -eq $PreCheckMD5Checksum)) { continue }
                            else { Write-Error "Aborting because $Name condition was met `'$PreCheckAction`'"; return }
                        }
                        else { Write-Error "Please provide a MD5 Checksum to validate the $Name condition" }
                    }
                    'If file exists and MD5 checksum matches then Abort' { 
                        if ($PreCheckMD5Checksum) {
                            if ($FileExists -and ((Get-MD5Hash $PreCheckFileName) -eq $PreCheckMD5Checksum)) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return }
                            else { continue }
                        }
                        else { Write-Error "Please provide a MD5 Checksum to validate the $Name condition" }
                    }
                    'If file exists and MD5 checksum does not matches then Abort' { 
                        if ($PreCheckMD5Checksum) {
                            if ($FileExists -and !((Get-MD5Hash $PreCheckFileName) -eq $PreCheckMD5Checksum)) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return }
                            else { continue }
                        }
                        else { Write-Error "Please provide a MD5 Checksum to validate the $Name condition" }
                    }
                }
                break;
            }
            'Check if registry information exists' {
                $KeyExists, $ValueMatches, $DataMatches = Test-Registry -Key $PreCheckKey -Value $PreCheckValue -Data $PreCheckData
                Switch ($PreCheckAction) {
                    'If key exists then Continue' { if ($KeyExists) { continue } else { Write-Error "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If key does not exists then Continue' { if (-not $KeyExists) { continue } else { Write-Error "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If key exists then Abort' { if ($KeyExists) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If key not exists then Abort' { if (-not $KeyExists) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If key exists and value matches then Continue' { if ($KeyExists -and $ValueMatches) { continue } else { Write-Error "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If key exists and value does not matches then Continue' { if ($KeyExists -and (-not $ValueMatches)) { continue } else { Write-Error "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If key exists and value matches then Abort' { if ($KeyExists -and $ValueMatches) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If key exists and value does not matches then Abort' { if ($KeyExists -and (-not $ValueMatches)) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If key exists and value exists and data matches then Continue' { if ($KeyExists -and $datamatches) { continue } else { Write-Error "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If key exists and value exists and data does not matches then Continue' { if ($KeyExists -and $ValueMatches -and (-not $datamatches)) { continue } else { Write-Error "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If key exists and value exists and data matches then Abort' { if ($KeyExists -and $ValueMatches -and $datamatches) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                    'If key exists and value exists and data does not matches then Abort' { if ($KeyExists -and $ValueMatches -and (-not $datamatches)) { Write-Output "Aborting because $Name condition was met `'$PreCheckAction`'"; return } else { continue } }
                }
                break;
            }
            'Check for file version' {
                $FileExists, $FileVersionMatches, $FileVersionGreater, $FileVersionLower = Test-FileVersion -Path $PreCheckFileName -UserSpecifiedVersion $PreCheckVersion
                Switch ($PreCheckAction) {
                    'If file exists and version equal to specified then Continue' { if ($FileExists -and $FileVersionMatches) { continue } else { Write-Error "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If file exists and version not equal to specified then Continue' { if ($FileExists -and (-not $FileVersionMatches)) { continue } else { Write-Error "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If file exists and version higher to specified then Continue' { if ($FileExists -and $FileVersionGreater) { continue } else { Write-Error "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
                    'If file exists and version lower to specified then Continue' { if ($FileExists -and $FileVersionLower) { continue } else { Write-Error "Aborting because $Name condition was not met `'$PreCheckAction`'"; return } }
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
 
    # validate connection and source string
    if (!($ConnectionPath -match "ftp://")) { $ConnectionPath = "ftp://" + $ConnectionPath }
    $Source = $Source.TrimStart('\').TrimEnd('\') -replace '\\', '/'

    # build credentials and source URL
    $credentials = New-Object System.Net.NetworkCredential($username, $password) 
    $Url = Join-URLParts -Parts $($ConnectionPath, $Source) -Seperator '/'

    # validate destination folder
    if (!$(Test-Path $Destination) -and ($CreateDestinationFolder -eq $false)) {
        Write-Error "Not a valid folder path: $Destination";
        exit;
    }
    # create destination directory
    if ($CreateDestinationFolder -and !$(Test-Path $Destination)) {
        New-Item $Destination -Type directory | Out-Null
    }

    if (Test-Path $Destination) {
        # validate destination folder
        if (!((Get-Item $Destination -ErrorAction silentlycontinue) -is [System.IO.DirectoryInfo])) { Write-Error "Not a valid folder path: $Destination"; break; }
        Switch ($Type) {
            'Folder' {
                DownloadFtpDirectory -Path $Url -Credentials $credentials -Recurse:$Recurse -localPath $Destination #| Format-Table -AutoSize
                $status = $?
                if (!$PerformPostCheck) {
                    if ($Status) { Write-Output 'FTP Download Successful' }else { Write-Error 'FTP Download Failed' }
                }                    
            }
            'File' {
                DownloadFtpFile -Path $Url -Credentials $credentials -localPath $(Join-Path $Destination $(Split-Path $Source -Leaf))
                $status = $?
                if (!$PerformPostCheck) {
                    if ($Status) { Write-Output 'FTP Download Successful' }else { Write-Error 'FTP Download Failed' }
                }                
            }
        }
    }
    else {
        Write-Error "Unable to find destination directory: `'$Destination`'"
    }

    #endregion Main

    #region PostCheck
    if ($PerformPostCheck) {
        $Name = 'Post-Check'
        $FileExists = Test-Path $PostCheckFileName
        Switch ($PostCheckAction) {
            'If file exists then Success' { if ($FileExists) { Write-Output "Successful for $Name '$PostCheckAction'" } else { Write-Output "Failed for $Name '$PostCheckAction'" } }
            'If file does not exists then Success' { if (-not $FileExists) { Write-Output "Successful for $Name '$PostCheckAction'" } else { Write-Output "Failed for $Name '$PostCheckAction'" } }
            'If file exists then Fail' { if ($FileExists) { Write-Output "Failed for $Name '$PostCheckAction'" } else { Write-Output "Successful for $Name '$PostCheckAction'" } }
            'If file does not exists then Fail' { if (-not $FileExists) { Write-Output "Failed for $Name '$PostCheckAction'" } else { Write-Output "Successful for $Name '$PostCheckAction'" } }
            'If file exists and MD5 checksum matches then Success' { 
                if ($PostCheckMD5Checksum) {
                    if ($FileExists -and ((Get-MD5Hash $PostCheckFileName) -eq $PostCheckMD5Checksum)) { Write-Output "Successful for $Name '$PostCheckAction'" }
                    else { Write-Output "Failed for $Name '$PostCheckAction'" }
                }
                else { Write-Error "Please provide a MD5 Checksum to validate the $Name condition" }
            }
            'If file exists and MD5 checksum does not matches then Success' { 
                if ($PostCheckMD5Checksum) {
                    if ($FileExists -and !((Get-MD5Hash $PostCheckFileName) -eq $PostCheckMD5Checksum)) { Write-Output "Successful for $Name '$PostCheckAction'" }
                    else { Write-Output "Failed for $Name '$PostCheckAction'" }
                }
                else { Write-Error "Please provide a MD5 Checksum to validate the $Name condition" }
            }
            'If file exists and MD5 checksum matches then Fail' { 
                if ($PostCheckMD5Checksum) {
                    if ($FileExists -and ((Get-MD5Hash $PostCheckFileName) -eq $PostCheckMD5Checksum)) { Write-Output "Failed for $Name '$PostCheckAction'" }
                    else { Write-Output "Successful for $Name '$PostCheckAction'" }
                }
                else { Write-Error "Please provide a MD5 Checksum to validate the $Name condition" }
            }
            'If file exists and MD5 checksum does not matches then Fail' { 
                if ($PostCheckMD5Checksum) {
                    if ($FileExists -and !((Get-MD5Hash $PostCheckFileName) -eq $PostCheckMD5Checksum)) { Write-Output "Failed for $Name '$PostCheckAction'" }
                    else { Write-Output "Successful for $Name '$PostCheckAction'" }
                }
                else { Write-Error "Please provide a MD5 Checksum to validate the $Name condition" }
            }
        }
        break;
    }
    #endregion Postcheck
}
catch {
    Write-Output "`nFTP Download Failed"
    Write-Error $_.Exception.message
}
