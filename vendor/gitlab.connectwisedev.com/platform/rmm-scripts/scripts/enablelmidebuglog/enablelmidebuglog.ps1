# $EnableDisableLog = "enable"
# $LocationOfEventLogs = ''
# $KeepEventLogFor = '90'

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}  

#region Default-Values-Input-Validation

$Action = @{
    $true  = 'Enabled'
    $false = 'Disabled'
}


$warningPreference = 'SilentlyContinue'
$already = $false
$Path = 'HKEY_LOCAL_MACHINE\SOFTWARE\LogMeIn\V5\Log'

if (!$LocationOfEventLogs) {
    $LocationOfEventLogs = 'C:\ProgramData\LogmeIn\'
}

if ($EnableDisableLog -eq 'Enable') {
    $Enable = $true
}
else {
    $Enable = $false
}

if (!$KeepEventLogFor) {
    $KeepEventLogFor = 90
}

try {
    [void][int]$KeepEventLogFor
}
catch {
    Write-Output "Failed to `"$($Action[$Enable].trimend('d'))`" LogMeIn Debug Logs"
    Write-Output "ArchivalDays can only set to a number. Current value = `"$KeepEventLogFor`""
    break
}

# removes and formatting issues in input variable like single forward slash or double backward slash
$LocationOfEventLogs = ([System.IO.DirectoryInfo]$LocationOfEventLogs).Fullname

$expectedValues = @"
Path,Name,Type,Data
$path,Debug,Reg_DWord,$([int] $enable)
$path,ArchivalDays,Reg_DWord,$KeepEventLogFor
$path,Directory,Reg_String,$LocationOfEventLogs
"@ | ConvertFrom-Csv
#endregion Default-Values-Input-Validation


#region functions-and-scriptblocks

function Test-Registry {
    [cmdletbinding()]
    param($Path, $Name, $Type, $Data)

    try {
        $TypeHash = @{
            String       = 'REG_SZ'
            ExpandString = 'REG_EXPAND_SZ'
            Binary       = 'REG_BINARY'
            DWord        = 'REG_DWORD'
            MultiString  = 'REG_MULTI_SZ'
            Qword        = 'REG_QWORD'
        }
    
        $Path = "REGISTRY::$Path"
        $PathExists = Test-Path $Path
        $Result = Get-ItemProperty $Path -Name $Name -ErrorAction SilentlyContinue
        $DataExists = [bool](($Result).$Name -as [String])
        $DataMatches = $Data -eq ($Result).$Name
        $Kind = Get-Item $Path -ErrorAction SilentlyContinue
        if ($kind) {
            $PropertyMatches = $Type -eq $TypeHash["$($Kind.GetValueKind($Name))"]
        }
        else {
            $PropertyMatches = $false
        }
    }
    catch {

    }
    return New-Object PSObject -Property @{
        'PathExists'      = if ($PathExists) { $true }else { $false }
        'DataExists'      = if ($DataExists) { $true }else { $false }
        'DataMatches'     = if ($DataMatches) { $true }else { $false }
        'PropertyMatches' = if ($PropertyMatches) { $true }else { $false }
    }
}

$StopService = {
    $Service = Get-Service Logmein 
    if ($Service.Status -eq 'Running') {
        Stop-Service -Name $Service.Name -Confirm:$false
    }
}
$StartService = { 
    $Service = Get-Service Logmein 
    if ($Service.Status -eq 'Stopped') {
        Start-Service -Name $Service.Name -Confirm:$false
    }

}
#endregion functions-and-scriptblocks

#region main
Try { 
    $ErrorActionPreference = 'Stop'
    $Results = @()
    $Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall', 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'
    $Product = Get-ChildItem $Registry -ErrorAction 'silentlycontinue' | Get-ItemProperty | Select Displayname | Where { $_.displayname -eq "logmein" } 

    if ($product) {
        Foreach ($item in $ExpectedValues) {
            #Write-host "Working on $($item.name)" -ForegroundColor Red
            $PreValidate = Test-Registry -Path $item.Path -Name $Item.Name -data $Item.data -type $Item.Type -ErrorAction SilentlyContinue

            if ($PreValidate.PathExists) {
                if (!$PreValidate.datamatches) {
                    & $StopService
                    Set-ItemProperty -Path "Registry::$($item.Path)" -Name $item.Name -Value $item.data
                    & $StartService
                }
                else {
                    if ($item.name -eq 'Debug') {
                        $Results += Write-Output "Debug log is already $($action[[bool][int]$item.Data])"
                        $already = $true
                    }
                }
            }
            else {
                & $StopService
                New-Item -Path "Registry::$Path"
                New-ItemProperty -Path "Registry::$($item.Path)" -Name $Item.Name -PropertyType $item.type -Value $item.data -Verbose
                & $StartService
            }
            
            $PostValidate = Test-Registry -Path $item.Path -Name $Item.Name -data $Item.data -type $Item.Type
            if ($PostValidate.datamatches) {
                if ($item.name -eq 'Debug' -and !$already) {
                    $Results += Write-Output "Debug log is now $($action[[bool][int]$item.Data])"
                }
                elseif ($item.name -ne 'Debug') {
                    if ($item.name -eq 'Directory' -and !(Test-Path $item.Data)) {
                        $Results += "{0} set to {1}: Success, but the folder doesn't exists, please create it." -f $item.Name, $item.Data
                    }
                    else {
                        $Results += "{0} set to {1}: Success" -f $item.Name, $item.Data
                    }
                }             
            }
            else {
                $Results += "{0} set to {1}: Failed" -f $item.Name, $item.Data                
            }
            
            if ($item.name -eq 'Debug' -and $Item.Data -eq 0) { break; }

        }

        $Results

    }
    else {
        Write-Output "LogMeIn is not installed. No action is performed."
    }
}
catch {
    Write-Output "Failed to `"$($Action[$Enable].trimend('d'))`" LogMeIn Debug Logs"
    $_
}
#endregion main
