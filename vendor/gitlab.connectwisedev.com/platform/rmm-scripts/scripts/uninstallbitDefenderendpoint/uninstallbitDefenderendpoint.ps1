#non-parameter
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

function Download-FromURL ($URL, $LocalFilePath) {
    $WebClient = New-Object System.Net.WebClient
    $WebClient.DownloadFile($URL, $LocalFilePath)
    if (-not(Test-Path $LocalFilePath)) { Write-Error "Unable to download the file from URL: $URL"}
}

$Software = 'Bitdefender Endpoint Security Tools'
$URL = 'http://dcmdwld.itsupport247.net/bittool_uninstall.zip'
$LocalPath = "$env:temp\bit_uninstall.zip"
$Destination = Join-Path $env:temp 'BitDefenderUninstaller'
$ErrorActionPreference = 'Stop'
try{

    Download-fromURL $URL $LocalPath # download the zip file from datat-center
    # expand the zip file in a new folder inside Temp folder
    if (!(Test-Path $Destination)) {
        [void] (New-Item -Path $Destination -ItemType Directory -Force)
        $Shell = new-object -com Shell.Application
        $Shell.Namespace($destination).copyhere($Shell.NameSpace("$LocalPath\bittool_uninstall").Items(),4) 
    }
       
    $Service = [bool](Get-Service EPSecurityService -ErrorAction silentlycontinue)
    $Process = [bool](Get-Process EPSecurityService -ErrorAction silentlycontinue)
    if ($Service -and $Process) {
        $process = Start-Process "$Destination\UninstallTool.exe" -arg "/silent /force:Endpoint Security by Bitdefender" -Wait -PassThru -ErrorAction 'Stop'
        If ($process.exitcode -eq 0) {
            Write-Output "Successfuly uninstalled '$Software'."
        }
        else {
            Write-Output "Failed to uninstall '$Software'. Exitcode: $($process.exitcode)"        
        }
    }
    else {
        Write-Output "'$Software' is not installed on the system."
    }
}
catch{
    Write-Output "Failed to uninstall '$Software'"
    Write-Error $_
}
