<#
    .SYNOPSIS
       Reinstall LogMeIn Host Software
    .DESCRIPTION
       Reinstall LogMeIn Host Software
    .Help
        Use msiexec.exe
        Use below path 
        #https://secure.logmein.com/LogMeIn.msi   
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>
$action = "reinstall"
#Get Product GUID
$Registry = 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall', 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall'

$Product = Get-ChildItem $Registry -ErrorAction 'SilentlyContinue' | Get-ItemProperty | Where-Object { $_.DisplayName -match "LogMeIn" } 

$ProductGUID = $Product | Select-Object -ExpandProperty PSChildName -First 1

#Get Product type and set DepoyID
$ProductType = (Get-WmiObject -Class Win32_OperatingSystem).producttype 
if ($ProductType -eq '1') {
    $Dpid = "DEPLOYID=zuar32f54kdnyal62abu0qm5d134cd81irm8baa8"
}
else {
    $Dpid = "DEPLOYID=86t1g7qkbkmgh9v7w8th8r44m8k8lfwqbynigoz5"
}

function Download_File($requrl, $dpath) {
    $dn = New-Object System.Net.WebClient
    $dn.DownloadFile("$requrl", "$dpath")
}

#Download LogMeIn.msi file if product available
if ($ProductGUID) {
    $url = "https://secure.logmein.com/LogMeIn.msi"
    $downloadpath = "C:\Windows\Temp\LogMeIn.msi"
    Download_File -requrl $url -dpath $downloadpath
}

Try {  

    if ($action -eq "reinstall") {

        if ($ProductGUID) {
            
            #Uninstall LogMeIn Process
            $uninstall = Start-Process "msiexec.exe" -arg "/X $ProductGUID /quiet" -Wait -PassThru -ErrorAction 'Stop'
            
            #ExitCode for uninstall process
            If (($uninstall.exitcode -eq '3010') -or ($uninstall.exitcode -eq '0')) {

                if (Test-Path $downloadpath) {
                    
                    #Install LogMeIn Process
                    $processinfo = New-object System.Diagnostics.ProcessStartInfo
                    $processinfo.CreateNoWindow = $true
                    $processinfo.UseShellExecute = $false
                    $processinfo.RedirectStandardOutput = $true
                    $processinfo.RedirectStandardError = $true
                    $processinfo.FileName = "msiexec.exe"
                    $processinfo.Arguments = "/quiet /i $downloadpath $Dpid"
                    $process = New-Object System.Diagnostics.Process
                    $process.StartInfo = $processinfo
                    [void]$process.Start()
                    $process.WaitForExit()

                    #ExitCode for install process
                    If (($process.exitcode -eq '3010') -or ($process.exitcode -eq '0')) {
                        
                        Download_File -requrl "http://update.itsupport247.net/lmi/numhostidwp.exe" -dpath "C:\Windows\Temp\numhostidwp.exe"
                        if (Test-Path "C:\Windows\Temp\numhostidwp.exe") {
                            Start-Process "C:\Windows\Temp\numhostidwp.exe" -Wait -PassThru -ErrorAction 'Stop' | Out-Null
                        }
                        
                        Write-Output "LogMeIn Host software reinstalled successfully on system $ENV:COMPUTERNAME"                
                    }
                    else {
                        Write-Output "LogMeIn Host software reinstallation failed on system $ENV:COMPUTERNAME"
                    }
                }
                else {
                    Write-Output "Download Failed for LogMeIn Host software. Kindly execute script again"
                }
            }
            else {
                Write-Output "LogMeIn Host software Uninstallation failed on system $ENV:COMPUTERNAME"
            }
        }       
        else {
            Write-Output "LogMeIn Host software not installed on this system $ENV:COMPUTERNAME"
        }
    }
    else {
        Write-Output "Kindly select action as reinstall"
    }
}
catch {
    Write-Error $_.Exception.Message
}
