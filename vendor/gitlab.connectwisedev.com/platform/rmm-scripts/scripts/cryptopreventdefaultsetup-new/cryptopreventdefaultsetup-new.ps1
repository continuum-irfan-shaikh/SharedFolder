<#
    .SYNOPSIS
       Install/Uninstall Crypto Prevent
    .DESCRIPTION
       Install/Uninstall Crypto Prevent
       When Installing use URL "http://dcmdwld.itsupport247.net/CryptoPreventBulkSetupV9.exe" to donwload the setup.
       Script will support Version 9.0.0.0 for install and 7,8 for uninstall
       Install Parameter --  /VERYSILENT /NORESTART /SUPPRESSMSGBOXES
       Uninstall String for Version 7 and 8 :  -----      C:\Program Files\Foolish IT\CryptoPrevent\unins000.exe /SILENT
    .Help
        Use msiexec.exe
        Use below path 
        http://dcmdwld.itsupport247.net/CryptoPreventBulkSetupV9.exe   
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>
<#JSON Variable
$Action = "Uninstall"   #Uninstall
#$InstallMajorVersion = 9
$UninstallMajorVersion = 9
#$InstallSubVersion = "9.0.0.0"
$UninstallSubVersion = "9.0.0.0"
#$InstallExecutionMode = "Install from continuum"
$UninstallExecutionMode = "Uninstall"
#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

$ErrorActionPreference = "Stop"

#Get Product
function Get_Prod($MVersion, $SubVersion) {
    $Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall', 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'

    $Product = Get-ChildItem $Registry -ErrorAction SilentlyContinue | Get-ItemProperty | Where-Object { $_.DisplayName -match "CryptoPrevent" -and $_.DisplayVersion -eq "$SubVersion" -and $_.MajorVersion -eq "$MVersion" }

    $Item = $Product | Select-Object -ExpandProperty UninstallString -First 1 -ErrorAction SilentlyContinue
    return $Item
}

function Get_Product() {
    $Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall', 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'

    $Product = Get-ChildItem $Registry -ErrorAction SilentlyContinue | Get-ItemProperty | Where-Object { $_.DisplayName -match "CryptoPrevent" }

    $Item = $Product | Select-Object -ExpandProperty UninstallString -First 1 -ErrorAction SilentlyContinue
    return $Item
}

#Respected version's URL Hash table
$URL_Hash = @{
    '9.0.0.0' = "http://dcmdwld.itsupport247.net/CryptoPreventBulkSetupV9.exe"
}

#Download CryptoPrevent Setup
function Download_File($requrl, $dpath) {
    try {
        $dn = New-Object System.Net.WebClient
        $dn.DownloadFile("$requrl", "$dpath")
    }
    catch {
        Write-Error $_.Exception.Message
    }
}

if ($Action -eq "Install" -and $InstallExecutionMode -eq "Install from continuum") {
    Try {  
        $ErrorActionPreference = "Stop"
        $UninstallString = Get_Prod -MVersion "$InstallMajorVersion" -SubVersion "$InstallSubVersion"

        if (!$UninstallString) {
            
            #Download CryptoPrevent setup    
            $ddpath = [System.IO.Path]::GetFileName("$($URL_Hash[$InstallSubVersion])")
            $downloadpath = join-path "C:\Windows\Temp" "$ddpath"
            Download_File -requrl "$($URL_Hash[$InstallSubVersion])" -dpath "$downloadpath"
    
            if (Test-Path $downloadpath) {
                    
                #Install CryptoPrevent Process
                $processinfo = New-object System.Diagnostics.ProcessStartInfo
                $processinfo.CreateNoWindow = $true
                $processinfo.UseShellExecute = $false
                $processinfo.RedirectStandardOutput = $true
                $processinfo.RedirectStandardError = $true
                $processinfo.FileName = "$downloadpath"
                $processinfo.Arguments = "/VERYSILENT /SUPPRESSMSGBOXES /NORESTART"
                $process = New-Object System.Diagnostics.Process
                $process.StartInfo = $processinfo
                [void]$process.Start()
                $process.WaitForExit()

                #ExitCode for install process
                If ($process.exitcode -eq '0') {      
                    Write-Output "CryptoPrevent $InstallMajorVersion Installed successfully on system $ENV:COMPUTERNAME" 
                    if (Test-Path $downloadpath) {
                        Remove-Item $downloadpath -Force -ErrorAction SilentlyContinue
                    }                   
                }
                else {
                    if (Test-Path $downloadpath) {
                        Remove-Item $downloadpath -Force -ErrorAction SilentlyContinue
                    }
                    Write-Error "CryptoPrevent $InstallMajorVersion Installation failed on system $ENV:COMPUTERNAME"
                    Exit;
                }
            }
            else {
                Write-Error "Download CryptoPrevent $InstallMajorVersion setup failed. Kindly try again"
                Exit;
            }
        }         
        else {
            Write-Output "CryptoPrevent $InstallMajorVersion already installed on this system $ENV:COMPUTERNAME"
            Exit;
        }
    }
    catch {
        Write-Error $_.Exception.Message
    }
}   
#Uninstall Process    
elseif ($Action -eq "Uninstall" -and $UninstallExecutionMode -eq "Uninstall") {
    Try {  
        $ErrorActionPreference = "Stop"
        $UninstallString = Get_Prod -MVersion "$UninstallMajorVersion" -SubVersion "$UnInstallSubVersion"
        $prd = Get_Product
        
        if ($UninstallString -and $prd) {       
            #Install CryptoPrevent Process
            $processinfo = New-object System.Diagnostics.ProcessStartInfo
            $processinfo.CreateNoWindow = $true
            $processinfo.UseShellExecute = $false
            $processinfo.RedirectStandardOutput = $true
            $processinfo.RedirectStandardError = $true
            $processinfo.FileName = "$UninstallString"
            $processinfo.Arguments = "/VERYSILENT /SUPPRESSMSGBOXES /NORESTART"
            $process = New-Object System.Diagnostics.Process
            $process.StartInfo = $processinfo
            [void]$process.Start()
            $process.WaitForExit()

            #ExitCode for install process
            If ($process.exitcode -eq '0') {     
                Write-Output "CryptoPrevent $UninstallMajorVersion uninstalled successfully on system $ENV:COMPUTERNAME"                    
            }
            else {
                Write-Error "CryptoPrevent $UninstallMajorVersion uninstallation failed on system $ENV:COMPUTERNAME"
                Exit;
            }   
        }
        elseif ((!$UninstallString) -and ($prd)) {
            Write-Error "CryptoPrevent version information not found on this system $ENV:COMPUTERNAME"
            Write-Error "Kindly uninstall manually."
            Exit;
        }          
        else {
            Write-Output "CryptoPrevent $UninstallMajorVersion not installed on this system $ENV:COMPUTERNAME"
            Exit;
        }
    }
    catch {
        Write-Error $_.Exception.Message
    }
}
else {
    Write-Output "Kindly select action and Execution mode"
    Exit;
}
