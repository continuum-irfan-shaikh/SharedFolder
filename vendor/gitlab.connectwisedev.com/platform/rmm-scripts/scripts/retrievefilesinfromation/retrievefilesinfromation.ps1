#$File = "C:\Prateek\AcroRdrDC1801120040_en_US.exe,C:\Prateek\vipre-business.7.0.3.12.exe"

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

$File = $File -split "," | foreach { $_.trim() }

Foreach($Item in $File){

    if(Test-Path $Item){
        $VersionInfo = Get-ItemProperty $Item -OutVariable FileProperty | Select -ExpandProperty VersionInfo        
        Write-Output "`nCompany name    : $($VersionInfo.CompanyName)"
        Write-Output "Description     : $($VersionInfo.FileDescription)"
        Write-Output "File Name       : $($FileProperty | Select-Object -ExpandProperty Fullname)"
        Write-Output "File type       : $($FileProperty | Select-Object -ExpandProperty Extension)"
        Write-Output "File version    : $($VersionInfo.FileVersion)"
        Write-Output "Internal name   : $($VersionInfo.InternalName)"
        Write-Output "Legal copyright : $($VersionInfo.LegalCopyright)"    
        Write-Output "Legal trademark : $($VersionInfo.LegalTrademarks)"    
        Write-Output "Original name   : $($VersionInfo.OriginalFilename)" 
        Write-Output "Product Name    : $($VersionInfo.ProductName)"
        Write-Output "Product version : $($VersionInfo.ProductVersion)"   
        Write-Output "Last Modified   : $($FileProperty | Select-Object -ExpandProperty LastWriteTime)"
        Write-Output "Size of File    : $($FileProperty | Select-Object -ExpandProperty Length)"
    }
    else{
        Write-Output "File: `'$Item`' does not exist on the system."
    }
}
