#$FolderPath = 'C:\PSTools'

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}


$ErrorActionPreference = "SilentlyContinue"

    $SubFolderPath = Get-ChildItem -Path $FolderPath -Force | Where-Object {$_.PSIsContainer} | Select -Expandproperty FullName
    $TargetFolders = @()
    $TargetFolders += $FolderPath
    $TargetFolders += $SubFolderPath

    Write-Output "`nFollowing are the details of folder and its sub-folders"
    ForEach ($Path in $TargetFolders) {
        
        If (Test-Path $Path) {
            $Size = (Get-ChildItem $path -Recurse -Force | Measure-Object -Property Length -Sum).Sum
            if ($Size -ge 1GB) { $SizeWithUnit = "{0:N2} GB" -f ($Size / 1gb) }elseif ($Size -ge 1Mb) { $SizeWithUnit = "{0:N2} MB" -f ($Size / 1mb)  }else {$SizeWithUnit = "{0:N2} KB" -f ($Size / 1kb)} 
            Write-Output "`nPath: $Path"
            Write-Output "Size: $SizeWithUnit"
        }
    }
    $ErrorActionPreference = "Continue"
