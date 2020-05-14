if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

$start_time = Get-Date
$PathArray = $output -split "\\"
$DirectoryPath = ($PathArray | Where-Object { $_ -ne $PathArray[-1] }) -join "\"


if ($onErrorContinue) {
    $ErrorActionPreference = 'Continue'
}
else {
    $ErrorActionPreference = 'Stop'
}

if ((-not $overwrite) -and (Test-Path $output)) {
    Write-Error "File exists. Delete existing file or check overwrite option"
    Exit
}
 
if ((-not $createDirectory) -and (-not (Test-Path $DirectoryPath))) {
    Write-Error "Output path is not valid, download directory does not exist. Input valid output path or check createDirectory option"
    return
} 

if ($createDirectory -and (-not (Test-Path $DirectoryPath))) {
    New-Item -Path $DirectoryPath -ItemType "directory" -Force >$null
}


$Powershell_Version = ($PSVersionTable.PsVersion).Major

if ($Powershell_Version -eq 2) {

    #call WebClient
    $client = New-Object System.Net.WebClient
    if ($username -and $password) {
        $pair = "${username}:${password}"
        $base64 = [System.Convert]::ToBase64String([System.Text.Encoding]::ASCII.GetBytes($pair))
        $client.Headers.add('Authorization', "Basic $base64")
    }
    try {
        $ErrorActionPreference = 'stop'
        $client.DownloadFile($url, $output)
        Write-Output "Download successful. Time taken: $((Get-Date).Subtract($start_time).Seconds) second(s)"
    }
    catch { Write-Error "Download failed! $_.Exception.Message" }

}
else {
    #call invoke-webrequest
    try {
        $ErrorActionPreference = 'Stop'
        if ($username -and $password) {
            $pair = "${username}:${password}"
            $bytes = [System.Text.Encoding]::ASCII.GetBytes($pair)
            $base64 = [System.Convert]::ToBase64String($bytes)
            $basicAuthValue = "Basic $base64"
            $headers = @{ Authorization = $basicAuthValue }
            Invoke-WebRequest -uri "$url" -Headers $headers -OutFile "$output"
            Write-Output "Download successful. Time taken: $((Get-Date).Subtract($start_time).Seconds) second(s)"
        }
        else {
            Invoke-WebRequest -Uri "$url" -OutFile "$output" 
            Write-Output "Download successful. Time taken: $((Get-Date).Subtract($start_time).Seconds) second(s)"
        }
    }
    catch { 
        Write-Error "Download failed! $_.Exception.Message"
    }
    
}


if ($md5) {
    $md5Obj = New-Object -TypeName System.Security.Cryptography.MD5CryptoServiceProvider
    $filemd5 = [System.BitConverter]::ToString($md5Obj.ComputeHash([System.IO.File]::ReadAllBytes($output)))
    $filemd5 = $filemd5 -replace "-", ""
    if ($filemd5 -ne $md5) {
        Remove-Item -path $output
        Write-Error "md5 mismatch ($md5!=$filemd5) for downloaded file, file removed"
        Exit
    }
    Else { "MD5 checksum verification passed." }
}
