# $Sourcetype = 'Network'
# $Filepath = "\\SCRIPT-WIN2K8R2\Prateek\prateek.reg,\\SCRIPT-WIN2K8R2\Prateek\Prateek.reg"
# $UserName = 'Prateek'
# $Password = 'test@123'

$ErrorActionPreference = 'Stop'
$Filepath = $Filepath -split "," | foreach { $_.trim() }

ForEach ($File in $Filepath) {
    try {
        $Path = ''
        if ($File -notlike "*.reg") {
            Write-Error "File provided:`'$File`' doesn't have '.reg' extension. Please verify it is valid registry file name."
        }
        
        if ($Sourcetype -eq 'Network' -and !([uri]$file).IsUnc) {
            Write-Error "Source Type: `'$SourceType`' but `'$File`' is a 'local' path. Please provide `'$SourceType'` paths only." -ErrorAction Stop
        }
        
        if ($Sourcetype -eq 'Local' -and ([uri]$file).IsUnc) {
            Write-Error "Source Type: `'$SourceType`' but `'$File`' is a 'Network' Path. Please provide `'$SourceType'` paths only." -ErrorAction Stop
        }

        Switch (([uri]$File).IsUnc) {
            $false { 
                $Path = $File 
                if (Test-Path $Path) {
                    $Process = Start-Process -FilePath "$Env:SystemRoot\REGEDIT.exe" -ArgumentList '/s', $Path -Verb 'RunAs' -PassThru -Wait
                    If ($Process.ExitCode -eq 0) { Write-Output "Registry import successful. File Path: $File" }
                    Else { Write-Output "[$File] Registry import failed with Exit code: $($Process.ExitCode). Source Path:$File" }
                }
                else {
                    Write-Output "File: $Path doesn't exist."
                }
                Break;
            }
            $True {
                $Directory = [System.IO.Path]::GetDirectoryName($File)
                $Leaf = [System.IO.Path]::GetFileName($File)
                $Path = "$env:TEMP\registry.reg"
                $DriveLetter = Get-ChildItem function:[g-z]: -n | Where-Object { !(test-path $_) } | random
                $Net = New-Object -ComObject WScript.Network
                $Net.MapNetworkDrive($DriveLetter, "$Directory", $false, $UserName, $Password) 
                Start-Sleep -s 4
                If (Test-Path "$DriveLetter\") {
                    Copy-Item "$DriveLetter\$Leaf" $Path
                    if (Test-Path $Path) {
                        $Process = Start-Process -FilePath "$Env:SystemRoot\REGEDIT.exe" -ArgumentList '/s', $Path -Verb 'RunAs' -PassThru -Wait
                        If ($Process.ExitCode -eq 0) { Write-Output "Registry import successful. File Path: $File" }
                        Else { Write-Output "[$File] Registry import failed with Exit code: $($Process.ExitCode). File Path: $File" }
                    }
                    else {
                        Write-Output "Unable to copy file: $File to local system."
                    }
                }
                else {
                    Write-Output "Unable to map network location: $File to local system."
                }
                Break;
            }

            $null {
                Write-Error "`'$File`' is not a valid Path"
                Break;
            }
        }
    }
    catch {
        Write-Error $_
    }
}
