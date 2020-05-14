# $Restart = $true #[Boolean]

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

Try {      
    Function VerifyProduct($Name) {
        $Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall', 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'
        $Product = Get-ChildItem $Registry -ErrorAction 'SilentlyContinue' | Get-ItemProperty | Where-Object { $_.Displayname -like "*$Name*" }
        $Product = $Product | Select-Object -First 1
        if ($Product) {
            $GUID = if ($Product -and $Product.UninstallString -match '{\w{8}-\w{4}-\w{4}-\w{4}-\w{12}}') { $matches[0] }
        }
        else {
            $Product = Get-WmiObject win32_Product | Where-Object { $_.Name -like "*$Name*" }
            $GUID = $Product | Select-Object -ExpandProperty IdentifyingNumber
        }
        Return $GUID
    }
    if ($Restart) { $RestartArgument = '/forcerestart' }else { $RestartArgument = '/norestart' }
    $guid = VerifyProduct 'BitDefender Business Client'

    if ($guid) {
        Start-Process Msiexec.exe -ArgumentList "/X $GUID /qn $RestartArgument" -Wait
        if (!(VerifyProduct 'BitDefender Business Client')) {
            Write-Output "Successfuly Uninstalled software."
        }
        else {
            Write-Output "Failed to Uninstall software."
        }
    }
    else {
        Write-Output "Software is not installed on the system."
    }
}
catch {
    Write-Output "Failed to Uninstall."
    Write-Error $_
}
