if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

####################validating input path####################
if ((Test-path $NewFileLocation) -and (!$OverwriteExistingFile)) {
    Write-Error "File name provided is already exist."
    Exit;
}

if (($FileContents) -and (($FileContents.Length -gt 1024000 ))) {
    Write-Error "File Contents more than 1 MB is not allowed"
    Exit
} 

$testpath = Split-Path -Path $NewFileLocation -Parent
if (!(test-path $testpath) -and ($CreateDirectoryIfRequired)) {
    $newtestpath = Split-Path -Path $testpath -Parent
    if (!(Test-Path $newtestpath)) {
        Write-Error "Path provided is incorrect"
        Exit;
    }
    else {
        New-Item -ItemType Directory -Path (split-path $testpath -parent) -Name (split-path $testpath -Leaf) | Out-Null
        if (!(Test-Path $testpath)) {
            Write-error "Failed to create folder"
            Exit
        }
    }
}
elseif (!(test-path $testpath)) {
    Write-Error "Path provided is incorrect"
    Exit;
}

$name = Split-Path $NewFileLocation -Leaf

try {
    New-Item -path (Split-Path $NewFileLocation -Parent) -Name $name -ItemType File -Value $FileContents  -Force -ErrorAction Stop | Out-Null
    if ( (Test-Path $NewFileLocation) -And (-not $FileContents)) {
        Write-Output "File content was not provided hence empty file is created successfully.`n$NewFileLocation"
        Exit
    }
    Elseif ( (Test-Path $NewFileLocation) -And ($FileContents)){
       Write-Output "File created successfully.`n$NewFileLocation"
         Exit
    }
}
catch {
    Write-Error "Failed to create the file`n. $_.Exception.Message" 
    Exit
}    
