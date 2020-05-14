# $Filename = "$($env:COMPUTERNAME).jpg" # string
# $Type = 'network' # string
# $Path = '\\10.2.19.25\e$\Antivirus Setups\Prateek' # string
# $UserName = 'administrator' # string
# $Password = 'Grt@2018' # string

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

Function Get-RDPSessions {
    try {
        $ErrorActionPreference = 'Stop'
        if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64" -and $env:PROCESSOR_ARCHITECTURE -eq 'x86') { $Query = 'c:\windows\sysnative\query.exe' }else { $Query = 'c:\windows\System32\query.exe' }
        $LoggedOnUsers = if (($Users = (& $Query user 2>&1))) {
            $Users | ForEach-Object { (($_.trim() -replace ">" -replace "(?m)^([A-Za-z0-9]{3,})\s+(\d{1,2}\s+\w+)", '$1  none  $2' -replace "\s{2,}", "," -replace "none", $null)) } |
            ConvertFrom-Csv |
            Where-Object { $_.state -eq 'Active' } #-and $_.SESSIONNAME -like "rdp*" }
        }
        return $LoggedOnUsers
    }
    catch {
        # no action
    }
}

Function AppendToCSV ($Type, $Message) {
    New-Object -TypeName psobject -Property @{Type = $Type; Message = $Message } | ConvertTo-Csv -NoTypeInformation | Select-Object -Skip 1 | Out-File C:\temp\log.csv -Append
}


$Scriptblock = {

    Function AppendToCSV ($Type, $Message) {
        New-Object -TypeName psobject -Property @{Type = $Type; Message = $Message } | ConvertTo-Csv -NoTypeInformation | Select-Object -Skip 1 | Out-File C:\temp\log.csv -Append
    }

    try {
        $ErrorActionPreference = 'Stop'
        $Extension = [IO.Path]::GetExtension($FileName)
        if (!($Extension -eq '.jpg' -or $Extension -eq '.jpeg')) { Write-Error "File extension must be `'.jpg`'" }
                    
        [System.Reflection.Assembly]::LoadWithPartialName("System.Drawing") | Out-Null # load assemblies
        [System.Reflection.Assembly]::LoadWithPartialName("System.Windows.Forms") | Out-Null
        $Screen = [System.Windows.Forms.Screen]::PrimaryScreen.Bounds
        $Bitmap = New-Object System.Drawing.Bitmap $Screen.width, $Screen.height
        $Size = New-Object System.Drawing.Size $Screen.width, $Screen.height
        $FromImage = [System.Drawing.Graphics]::FromImage($Bitmap)
        $FromImage.copyfromscreen(0, 0, 0, 0, $Size, ([System.Drawing.CopyPixelOperation]::SourceCopy))
        
        $TempLoc = 'C:\Temp\'
        
        switch ($type) {
            'local' {
                $TempTarget = Join-Path $TempLoc $FileName
                $Target = Join-Path $Path $FileName
                if (![IO.Directory]::Exists($Path)) { [IO.Directory]::CreateDirectory($Path) | Out-Null } # create folder if that doesn't exists
                $Bitmap.Save("$TempTarget", ([system.drawing.imaging.imageformat]::Jpeg))
            }
            'network' {
                $DriveLetter = Get-ChildItem function:[g-z]: -n | Where-Object { !(Test-Path $_) } | Get-Random
                $Network = New-Object -ComObject WScript.Network
                $Network.MapNetworkDrive($DriveLetter, "$Path", $false, $UserName, $Password) 
                Start-Sleep -s 4
                if (Test-Path "$DriveLetter\") {
                    $Target = Join-Path $DriveLetter $FileName
                    $Bitmap.Save($Target, ([system.drawing.imaging.imageformat]::Jpeg))
                }         
            }
        }
    }
    catch {  AppendToCSV "Output" "Error: $($_.exception.message)" }
    finally {
        try { $Network.RemoveNetworkDrive($DriveLetter, $True) }catch { }
    }

}

$file = @"
`$Filename = "$FileName"
`$Type = "$Type"
`$Path = "$Path"
`$UserName = "$Username"
`$Password = "$Password"
"@ + $Scriptblock.ToString()

$ErrorActionPreference = 'Stop'
$LogDir = 'C:\temp\'
$LogFilePath = 'C:\temp\Log.csv'
$MainFilePath = 'C:\temp\CaptureScreenshot.ps1'
$TaskName = "CaptureScreen"

if (![IO.Directory]::Exists($LogDir)) { [IO.Directory]::CreateDirectory($LogDir) | Out-Null } # create folder if that doesn't exists
$File | Out-File $MainFilePath -Encoding utf8 # copy main script which will be later executed from scheduled task

try {
    if (($Session = Get-RDPSessions)) {
        '"Message","Type"' | Out-File $LogFilePath -Append
        $User = $Session | Select-Object -First 1 -ExpandProperty Username # select one user profile from active RDP sessions
        
        Start-Sleep -Seconds 5
        $Task = "PowerShell.exe -executionpolicy bypass -NoExit -noprofile -WindowStyle Hidden -command '. $MainFilePath '"
        $StartTime = (Get-Date).AddMinutes(2).ToString('HH:mm') # time in 24hr format
        
        schtasks.exe /create /s $($env:COMPUTERNAME) /tn $TaskName /sc once /tr $Task /st $StartTime /ru $User /F | Out-Null
        Start-Sleep -Seconds 5

        schtasks.exe /End /TN $TaskName | Out-Null
        schtasks.exe /Run /TN $TaskName | Out-Null
        Start-Sleep -Seconds 120
        
        $Target = Join-Path $Path $Filename
        if($Type -eq 'Local'){
            $TempTarget = Join-Path 'C:\Temp\' $Filename
            Copy-Item $TempTarget $Target -ErrorAction SilentlyContinue
        }    
        
        if (Test-Path $Target) { 
            AppendToCSV "Output" "Successfully captured the screenshot and saved to location: $(Join-Path $Path $Filename)"
        }
        else { 
            AppendToCSV "Output" "Failed to capture the screenshot."
        }

        # Sending logs to PowerShell Output\Error streams so that agent can capture it
        if (Test-Path $LogFilePath) {
            $Logs = Import-Csv $LogFilePath
            foreach ($item in $Logs) {
                switch ($item.type) {
                    'Output' { Write-Output "`n$($item.message)" }
                    'Error' { Write-Error "`n$($item.message)" }
                }
            }
        }

        schtasks.exe /Delete /TN $TaskName /F | Out-Null # scheduled task cleanup
        Remove-Item -Path $MainFilePath, $LogFilePath, $TempTarget -ErrorAction SilentlyContinue # file cleanup
    }
    else {
        Write-Output "This script requires logon user and currently no user is logged in. `nNo action will be performed."; exit;
    }
}
catch {
    Write-Error $_
}
