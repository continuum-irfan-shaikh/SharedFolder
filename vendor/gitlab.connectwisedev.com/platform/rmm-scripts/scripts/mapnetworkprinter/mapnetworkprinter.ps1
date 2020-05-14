# $action = "Connect" # STRING
# #$type = 'TCP/IP' # STRING
# $type = 'Shared' # STRING
# #$PrinterShare = "\\10.2.19.116\my HP Printer"
# $PrinterShare = "\\10.2.19.10\HPPrinter" # STRING
# $Default = $true # BOOLEAN
# $username = "Grtdc\prateek" # STRING
# $password = "India@1234" # STRING


if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}


$ErrorActionPreference = 'Stop'
$URI = [uri] $PrinterShare
if ($URI.IsUnc) {
    $ShareName = $URI.absolutepath -replace '/', ''
}

function ConnectPrinter {
    try {
        if ($username -and $password) {
            net use $PrinterShare /USER:$username $password /persistent:no
        }
        $net = New-Object -ComObject WScript.Network
        $net.AddWindowsPrinterConnection("$PrinterShare") | Out-Null

        if ($default) {
            $net.SetDefaultPrinter($PrinterShare)
        }

        $obj = New-Object PSObject -Property @{
            'isPrinterMapped' = if ((ListPrinters).ShareName -contains $ShareName) { $True }else { $false }
            'setToDefault'    = if (Get-WmiObject Win32_Printer | Where-Object { $_.Default -eq $true -and $_.ShareName -eq $ShareName }) { $true }else { $false }
        }
        return $obj
    }
    catch {
        $Message = "Failed to map the Printer `'$PrinterShare`'"
        if ($Error[0].exception -like "*error 86*") {
            Write-Output "$Message. Please check the username\password and try again."
        }
        else {
            Write-Error "$Message. $($Error[0].exception.Message)`nPlease check the Printer network path and credentials and try again."
        }
    }
}

function DisconnectPrinter {
    $net = New-Object -ComObject WScript.Network
    $net.RemovePrinterConnection($PrinterShare)
    if ($?) { return $true } else { return $false }
    
}

function DisconnectAllPrinter {
    ListPrinters | ForEach-Object { $_.delete() }
    if (ListPrinters) { $false }
    else { return $true }
}

Function ListPrinters {
    Get-WmiObject Win32_Printer -Filter "Network=True"
}

try {
    switch ($type) {
        'Shared' {
            switch ($action) {
                "Connect" { 
                    if ((ListPrinters).ShareName -contains $ShareName) {
                        Write-Output "Printer `'$PrinterShare`' already Exists. No action performed."
                    }
                    else {
                        $result = ConnectPrinter
                        if ($result.isPrinterMapped) { Write-Output "Printer:`'$PrinterShare`' added successfully, it may take few minutes to populate in 'Devices and Printers'" } else { Write-Output "Failed to add printer:`'$PrinterShare`'" }
                        if ($result.setToDefault) { Write-Output "Printer:`'$PrinterShare`' is now default printer" } else { Write-Output "Unable to set `'$PrinterShare`' as default printer." }
                    }
                    
                }
                "Disconnect" { 
                    if (DisconnectPrinter) {
                        Write-Output "Printer:`'$PrinterShare`' removed successfully"   
                    }
                    else {
                        Write-Output "Failed to remove Printer:`'$PrinterShare`'"
                    }
                }
                "Disconnect All" {
                    $all = ListPrinters
                    if ($all) {
                        if (DisconnectAllPrinter) {
                            Write-Output "All network printers removed successfully."   
                        }
                        else {
                            Write-Output "Failed to remove all network printers."
                        }
                    }
                    else {
                        Write-Output "No network printer found."
                    }

        
                }
            }
        }
        'TCP/IP' { }
    }
        
}
catch {
    Write-Output "`n"$_.Exception.Message
}
