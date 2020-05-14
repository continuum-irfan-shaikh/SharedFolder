# $fileName = 'zSaazAccess.exe'
# $fileDateOption = 'modifiedbefore'
# $date = "2019-03-20T21:27:00.000Z"
# $ltORgt = 'Lessthan'
# $sizeInKB = 200
# $ReadOnly = $false
# $archive = $false
# $system = $false
# $hidden = $false
# $encrypted = $false
# $compressed = $false
# $recursiveSearch = $true
# $searchIn = 'path'
# $searchPath = 'C:\Program Files\SAAZOD'

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}  
    
$Substitution = @{
    'LessThan'           = '$_.length -lt $sizeInKB'
    'GreaterThan'        = '$_.length -gt $sizeInKB'
    'ModifiedBefore'     = '$_.LastWriteTime -lt $Date'
    'ModifiedAfter'      = '$_.LastWriteTime -gt $Date'
    'CreatedBefore'      = '$_.CreationTime -lt $Date'
    'CreatedAfter'       = '$_.CreationTime -gt $Date'
    'LastAccessedBefore' = '$_.LastAccessTime -lt $Date'
    'LastAccessedAfter'  = '$_.LastAccessTime -gt $Date'
}
try {
    if($date){$date = [Datetime]$date}
    # Build attributes array
    $Attributes = @()
    if ($ReadOnly) {$Attributes += 'ReadOnly'}
    if ($Archive) {$Attributes += 'Archive'}
    if ($System) {$Attributes += 'System'}
    if ($Hidden) {$Attributes += 'Hidden'}
    if ($Encrypted) {$Attributes += 'Encrypted'}
    if ($Compressed) {$Attributes += 'Compressed'}

    # Build the 'Get-ChildItem' Command and Filter as per the user input
    $String = ''

    switch ($SearchIn) {
        'Path' {
            if (Test-Path $searchPath) {$String = ("Get-ChildItem", "`'$searchPath`'" -join ' ').trim()}
            else {Write-Output "`nPath doesn't exist: $searchPath"; continue }
            break;
        }
        'AllLocalDrive' {
            $String = ("Get-ChildItem", $((Get-WmiObject win32_logicaldisk -filter "DriveType = 3"| ForEach-Object { $_.deviceid, '\' -join ''}) -join ',') -join ' ').trim()
            break;
        }
    }

    $String = $String, '-Force' -join ' '
    if ($recursiveSearch) { $String = $String, "-Recurse" -join ' ' }
    $String = $String, " | Where-Object {!`$_.PSIsContainer -and `$_.name -eq `'$Filename`'" -join ''
    if ($Attributes) { ForEach ($attribute in $Attributes) { $attribute = $attribute -replace '-', ''; $String = $String, "`$_.attributes -like '*$Attribute*'" -join ' -and ' } }
    if ($fileDateOption -and $date) { $String = $String, $Substitution[$fileDateOption] -join ' -and '}
    if ($ltORgt -and $sizeInKB) {$sizeInKB=$sizeInKB*1kb;$String = $String, $Substitution[$ltORgt] -join ' -and '}
    $String = $String, '}' -join ''

    # Invoke the command string
    $ErrorActionPreference = 'SilentlyContinue'
    $List = Invoke-Expression -Command $String 
    if ($List) { Write-Output "`nList of Files:`n"; $List | ForEach-Object { Write-Output $_.Fullname } }
    else { Write-Output "`nNo files found with specified criteria.`n" }
}
catch { Write-Error $_ }
