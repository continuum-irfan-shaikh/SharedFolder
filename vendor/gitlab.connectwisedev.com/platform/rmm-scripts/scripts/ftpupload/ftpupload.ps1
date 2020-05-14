# $Connection = '10.2.130.27' # string
# $Username = "FTPUser" # string 
# $Password = "test@123" # string
# $Type = 'File' # string
# $SourceDirectoryPath = 'C:\Prateek\test.ps1' # string
# $DestinationDirectoryPath = 'zeus\alpha\beta\' # string
# $RecursiveUpload = $true # boolean
# $CreateDirectory = $true # boolean

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}  

#Startregion Functions
function Join-URLParts {
    param ([string[]] $Parts, [string] $Seperator = '')
    $search = '(?<!:)' + [regex]::Escape($Seperator) + '+'  #Replace multiples except in front of a colon for URLs.
    $replace = $Seperator
    ($Parts | Where-Object { $_ -and $_.Trim().Length }) -join $Seperator -replace $search, $replace
}

Function CreateDirectoryStructure ([System.IO.DirectoryInfo[]]$SrcFolders) {
    $WebClient = New-Object System.Net.WebClient 
    $WebClient.Credentials = New-Object System.Net.NetworkCredential($Username, $Password)  
    
    foreach ($folder in $Srcfolders) {    
        $SrcFolderPath = $SourceDirectoryPath -replace "\\", "\\" -replace "\:", "\:"   
        $DesFolder = $folder.Fullname -replace $SrcFolderPath, $(Join-URLParts -Parts $($Connection, $DestinationDirectoryPath) -Seperator '/')
        $DesFolder = $DesFolder -replace "\\", "/"
        #"Creating Directory: $DesFolder"
        
        try {
            $makeDirectory = [System.Net.WebRequest]::Create($DesFolder);
            $makeDirectory.Credentials = New-Object System.Net.NetworkCredential($Username, $Password);
            $makeDirectory.Method = [System.Net.WebRequestMethods+FTP]::MakeDirectory;
            $makeDirectory.GetResponse() | Out-Null
            #folder created successfully
        }
        catch [Net.WebException] {
            try {
                #if there was an error returned, check if folder already existed on server
                $checkDirectory = [System.Net.WebRequest]::Create($DesFolder);
                $checkDirectory.Credentials = New-Object System.Net.NetworkCredential($Username, $Password);
                $checkDirectory.Method = [System.Net.WebRequestMethods+FTP]::PrintWorkingDirectory;
                $response = $checkDirectory.GetResponse() | Out-Null
                #folder already exists!
            }
            catch [Net.WebException] {
                #if the folder didn't exist
            }
        }
    }
}
Function CreateRootDirectoryIfRequired ($Path) {
    try {
        If (!(CheckFTPDirectory $Path)) {
            CreateNestedFTPDirectory $Path $Connection
        }
    }
    catch {
        Write-Error $_
    }
}

function CreateNestedFTPDirectory ($DirectoryPath, $ServerPath) {
    $DirectoryParts = ([uri]$DirectoryPath).segments | ForEach-Object { $_ -replace "/", "" } | Where-Object { ![string]::IsNullOrEmpty($_) }
    
    If ($DirectoryParts.Count -gt 1) {
        $Position = "";
        for ($i = 0; $i -lt $DirectoryParts.Length; $i++) {
            try {
                $Position += $DirectoryParts[$i] + "/";	
                $AbsoluteTemporaryPath = New-Object System.Uri($(Join-URLParts -Seperator '/' -Parts ($ServerPath, $Position) ));
                $WebRequest = [System.Net.WebRequest]::Create($AbsoluteTemporaryPath);
                $WebRequest.KeepAlive = $false;
                $WebRequest.Credentials = New-Object System.Net.NetworkCredential($Username, $Password);
                $WebRequest.Method = [System.Net.WebRequestMethods+Ftp]::MakeDirectory;
                $WebRequest.GetResponse() | Out-Null
                
            }
            catch [Net.WebException] { 
                continue;
            }
        }
    }
    else {	
        $AbsoluteTemporaryPath = New-Object System.Uri($(Join-URLParts -Seperator '/' -Parts ($ServerPath, $(($DirectoryParts + "/"))) ));
        $WebRequest = [System.Net.WebRequest]::Create($AbsoluteTemporaryPath);
        $WebRequest.KeepAlive = $false;
        $WebRequest.Credentials = New-Object System.Net.NetworkCredential($Username, $Password);
        $WebRequest.Method = [System.Net.WebRequestMethods+Ftp]::MakeDirectory;
        $WebRequest.GetResponse() | Out-Null    
    }
}

Function CheckFTPDirectory ($Path) {
    try {
        $CheckDirectory = [System.Net.WebRequest]::Create($Path);
        $CheckDirectory.Credentials = New-Object System.Net.NetworkCredential($Username, $Password);
        $CheckDirectory.Method = [System.Net.WebRequestMethods+FTP]::PrintWorkingDirectory;
        $Response = $CheckDirectory.GetResponse()
        If ($Response) { Return $true }
    }
    catch {
        return $false
    }
}

Function UploadFiles ([System.IO.FileInfo[]]$SrcFiles) {
    $WebClient = New-Object System.Net.WebClient 
    $WebClient.Credentials = New-Object System.Net.NetworkCredential($Username, $Password)  
    
    foreach ($entry in $SrcFiles) {
        $SrcFullname = $entry.fullname
        #"Uploading File: $SrcFullname"
        $SrcName = $entry.Name
        if ($type -eq 'File') {
            $SrcFilePath = (Split-Path $SourceDirectoryPath) -replace "\\", "\\" -replace "\:", "\:"
        }
        else {
            $SrcFilePath = (($SourceDirectoryPath) -replace "\\", "\\" -replace "\:", "\:").trimend('\')
        }
        
        $DesFile = $SrcFullname -replace $SrcFilePath, $($(Join-URLParts -Parts $($Connection, $DestinationDirectoryPath) -Seperator '/'))
        $DesFile = $DesFile -replace "\\", "/"
        #Write-Output "Uploading `'$($entry.fullname)`' to localtion: `'$DesFile`'"
        
        $uri = New-Object System.Uri($DesFile) 
        $WebClient.UploadFile($uri, $SrcFullname)
    }
}
#endregion Functions

#Startregion Main
try {
    $ErrorActionPreference = 'Stop'
    if (!($Connection -match "ftp://")) { $Connection = "ftp://" + $Connection }
    
    $DestinationDirectoryPath = $DestinationDirectoryPath.TrimStart('\').TrimEnd('\')
    Switch ($Type) {
        'Folder' {
            if (!((Get-Item $SourceDirectoryPath).gettype().name -eq 'DirectoryInfo')) { Write-Output "Not a valid folder path: $SourceDirectoryPath"; break; }
            $DestinationPath = ($(Join-URLParts ($Connection, $DestinationDirectoryPath) -Seperator '/') -replace "\\", "/")
            
            if ($RecursiveUpload) { $SrcEntries = Get-ChildItem $SourceDirectoryPath -Recurse }
            else { $SrcEntries = Get-ChildItem $SourceDirectoryPath }
            
            If ($CreateDirectory) { CreateRootDirectoryIfRequired $DestinationPath }
            
            If (CheckFTPDirectory $DestinationPath) {
                $Srcfolders = $SrcEntries | Where-Object { $_.PSIsContainer }
                $SrcFiles = $SrcEntries | Where-Object { !$_.PSIsContainer }
                if ($Srcfolders) {
                    CreateDirectoryStructure $Srcfolders
                }
                if ($SrcFiles) {
                    UploadFiles $SrcFiles
                    if ($?) { Write-Output "`nFTP Upload Successful" }
                }                    
            }
            else {
                Write-Output "Unable to find destination directory on FTP: `'$DestinationPath`'"
            }
        }
        
        'File' {
            if (!((Get-Item $SourceDirectoryPath).gettype().name -eq 'FileInfo')) { Write-Output "Not a valid file path: $SourceDirectoryPath"; break; }
            
            $DestinationPath = ($(Join-URLParts ($Connection, $DestinationDirectoryPath) -Seperator '/') -replace "\\", "/")
            $SrcFiles = Get-ChildItem $SourceDirectoryPath
            
            If ($CreateDirectory) { CreateRootDirectoryIfRequired $DestinationPath }
            
            If (CheckFTPDirectory $DestinationPath) {
                if ($SrcFiles) {
                    UploadFiles $SrcFiles
                    if ($?) { Write-Output "`nFTP Upload Successful" }
                } 
            }
            else {
                Write-Output "Unable to find destination directory on FTP: `'$DestinationPath`'"
            }
            
        }
    }
}
catch {
    Write-Output "`nFTP Upload Failed"
    Write-Error $_.Exception.message
}
#endregion Main
