<#  
.SYNOPSIS  
    Adobe Acrobat Reader DC installation
.DESCRIPTION  
    Adobe Acrobat Reader DC installation
.NOTES  
    File Name  : InstallAdobeAcrobatReaderDC.ps1
    Author     : Durgeshkumar Patel  
    Requires   : PowerShell V2 or greater.   
.PARAMETERS
    
.HELP
#>

<# JSON SCHEMA
$action = 'install'  #Drop Down
Version 
19

$subversion    #Drop Down
19.010.20098

Examples:-
$action = "install"
$version = "19"  #Not used in script. Just used for JSON UI for user
$subversion = "19.010.20098"

#>

if ("$action" -eq "uninstall") {
    $MyApp = Get-WmiObject -Class Win32_Product | Where-Object{$_.Name -eq "Adobe Acrobat Reader DC"}
    if (!$MyApp) {
        Write-Error "Adobe Acrobat Reader DC not found on the machine"
        Exit 1;
    }


    $MyApp = Get-WmiObject -Class Win32_Product | Where-Object{$_.Name -eq "Adobe Acrobat Reader DC"}
    $MyApp.Uninstall()
    Write-Output "Adobe Acrobat Reader DC has successfully been uninstalled"
    Exit;
}

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

$product_name = "Adobe Acrobat Reader DC"

#Get already installed product
function get-product {
    if ((gwmi win32_operatingsystem | select osarchitecture).osarchitecture -eq "64-bit") {
        $a = Get-ChildItem -Path HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object { $_.DisplayName -match $product_name } | Select-Object -ExpandProperty DisplayVersion
    }
    else {
        $a = Get-ChildItem -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall | Get-ItemProperty | Where-Object { $_.DisplayName -match $product_name } | Select-Object -expandProperty DisplayVersion
    }
    return $a
}

$installedversion = get-product

#Exit is same version is installed
if ("$installedversion" -eq "$subversion") {
    Write-Output "Adobe Acrobat Reader DC v$subversion already installed on system $ENV:COMPUTERNAME"
    Exit;
}

#Hash table for version and related downlowad link
$version = @{
    '19.010.20098' = "http://ardownload.adobe.com/pub/adobe/reader/win/AcrobatDC/1901020098/AcroRdrDC1901020098_en_US.exe"
    #'1801120058' = "http://ardownload.adobe.com/pub/adobe/reader/win/AcrobatDC/1801120058/AcroRdrDC1801120058_en_US.exe"

}

$ErrorActionPreference = 'Stop'

#Url to download required exe
$url = $version["$subversion"]

#Download path where exe will be downloaded.
$downloadpath = "C:\Windows\Temp\AdobeAcroRdrDC_en_US.exe"

#This function downloads the file
function downloadfile {

    $ErrorActionPreference = "Stop"
    $wc = New-Object System.Net.WebClient
    try {
        $wc.DownloadFile("$url", "$downloadpath")
    }
    catch {
        $ErrorActionPreference = "Continue"
        $_.Exception.Message
    }

    if (Test-Path $downloadpath) {
        return $true
    }
    else {
        Write-Output "Adobe Acrobat Reader DC v$subversion installation file download has been failed. Check internet connection and try again."
        Exit;
    }
}

try {
    #Download the file for Adobe Reader DC
    if (downloadfile) {

        #Install process with silent parameters
        $process = Start-Process $downloadpath -arg "/sAll /rs" -Wait -PassThru -ErrorAction 'Stop'

        #Succesful or Error output
        if (($process.exitcode -eq 0) -or ($process.exitcode -eq 1641) -or ($process.exitcode -eq 3010)  ) {
            Write-Output "Adobe Acrobat Reader DC v$subversion installation successful on system $ENV:COMPUTERNAME."
        }
        else {
            Write-Error "Failed to install Adobe Acrobat Reader DC v$subversion on system $ENV:COMPUTERNAME"
            Write-Error "Exit Code:- $($process.exitcode)"
        }
        #Remove downloaded file
        if (Test-Path $downloadpath) {
            Remove-Item $downloadpath -Force -ErrorAction SilentlyContinue
        }
    }
}
catch {
    $_.Exception.Message
}
