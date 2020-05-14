 <#
 Name : Disable standby mode
 Category : Green Computing 

    .Synopsis
        Disable standby mode.
    .Author
        Prateek.Singh@Continuum.net
    .Name 
        Disable standby mode
#>


try { 
    $ErrorActionPreference = 'Stop'
    # get the name of Active Power Scheme
    $Name = Invoke-Command -ScriptBlock {powercfg.exe /list} | foreach {
        if($_ -match "GUID.*?:\s+(?<GUID>\S*)\s+\((?<Name>.*?)\)\s*(?<Active>\*?)" -and [bool]$matches.Active){
            $matches.Name
        }
    }
    # stream redirection from error(2) to output stream(1) in case of errors
    # without this (2>&1) some versions of Windows OS will not catch an error
    Invoke-Command -ScriptBlock {
        powercfg.exe -x -monitor-timeout-ac 0 2>&1
        powercfg.exe -x -monitor-timeout-dc 0 2>&1
        powercfg.exe -x -disk-timeout-ac 0 2>&1
        powercfg.exe -x -disk-timeout-dc 0 2>&1
        powercfg.exe -x -standby-timeout-ac 0 2>&1
        powercfg.exe -x -standby-timeout-dc 0 2>&1
        powercfg.exe -x -hibernate-timeout-ac 0 2>&1
        powercfg.exe -x -hibernate-timeout-dc 0 2>&1
    }
    Write-Output "Disabled Standby Mode Successfully on Active Power Scheme: `'$Name`'"
}
catch {
    Write-Error $_.exception.Message
}
