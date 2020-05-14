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

$Software = 'IP Scanning Tool'
$URL = 'http://update.itsupport247.net/Patch/UNIPScan.exe'
$LocalPath = "$env:temp\UNIPScan.exe"
$ErrorActionPreference = 'Stop'
try {
    Download-fromURL $URL $LocalPath # download the zip file from datat-center
    $process = Start-Process $LocalPath -Wait -PassThru -ErrorAction 'Stop' -Verbose
    If ($process.exitcode -eq 0) {
        Write-Output "Successfuly uninstalled '$Software'."
    }
    else {
        Write-Output "Failed to uninstall '$Software'. Exitcode: $($process.exitcode)"        
    }
}
catch {
    Write-Output "Failed to uninstall '$Software'"
    Write-Error $_
}
