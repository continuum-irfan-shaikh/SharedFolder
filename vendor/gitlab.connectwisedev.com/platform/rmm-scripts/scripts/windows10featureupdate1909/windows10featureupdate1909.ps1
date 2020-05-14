<#
   Name :  "Windows 10 1909 feature update"
   Description": "Task will update earlier versions of Windows 10  to version 1909 feature update."
  
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

$ErrorActionPreference = 'Stop'
Try {
    $LocalPath = "$Env:Temp\Windows10Upgrade.exe"
    $URL = 'http://sdupdt.itsupport247.net/tpupdate/Windows10Upgrade1909.exe'
    $Parameters = "/Install /ClientID Win10Upgrade:VNL:NHV19:{} /SkipEULA /EosUi /quiet /norestart"

    # downloads file, if downloaa fails you get a failure message
    function Download-FromURL ($URL, $LocalFilePath) {
        $WebClient = New-Object System.Net.WebClient
        $WebClient.DownloadFile($URL, $LocalFilePath)
        if (-not(Test-Path $LocalFilePath)) { Write-Error "Unable to download the file from URL: $URL" }
    }

    Function IsUpdateRunning {
        (Get-Process Windows10UpgraderApp -ErrorAction SilentlyContinue) -as [Bool]
    }


    # condition 1 - checks operating system version
    $OS = Get-WmiObject win32_operatingsystem
    $IsWindows10 = ($OS | Where-Object { $_.caption -like "*Windows 10*" }) -as [bool]
    if (!$IsWindows10) {
        Write-Output "Not Applicable : Task applies to only Windows 10"
        Break;
    }

    # condition 2 - check release id
    [int]$ReleaseID = (Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion").ReleaseId
    if ($ReleaseID -ge 1909) {
        Write-Output "No Action Needed : Latest Feature update already installed"
        Break;
    }

    # condition 3 - check system drive free space
    $SystemDriveLetter = Get-Partition | Where-Object { $_.isboot } | Select-Object -ExpandProperty DriveLetter
    $SystemDriveFreeSpace = (Get-WmiObject -Class Win32_logicaldisk -Filter "DriveType = '3' AND DeviceID = '${SystemDriveLetter}:'").FreeSpace / 1GB
    if ($SystemDriveFreeSpace -lt 32) {
        Write-Error "Prerequisite check failed : Minimum 32GB free space is needed on system drive"
        Break;
    }

    # condition 4 - if a process is already running
    if (IsUpdateRunning) {
        Write-Error "Failed : One instance of update is already running. No action will be performed."
        Break;
    }


    $UpgradeFolder = "${SystemDriveLetter}:\Windows10Upgrade\"
    $LogFile = Join-Path $UpgradeFolder 'upgrader_default.log'

    # remove if folder exists
    Remove-Item $UpgradeFolder -Force -Confirm:$false -ErrorAction SilentlyContinue -Recurse

    # downloads the URL
    Download-FromURL -URL $URL -LocalFilePath $LocalPath
    $Process = Start-Process $LocalPath -ArgumentList $Parameters -PassThru

    # allow update folders and dependecies to be downloaded
    While (!(Test-Path "${SystemDriveLetter}:\Windows10Upgrade\Windows10UpgraderApp.exe")) {
        Start-Sleep -Seconds 1
    }
    Start-Sleep -Seconds 20 # additional delay to wait update exe to launch

    $Downloadflag = $false

    # wait for donwload status to appear in logs
    While (1) {
        $DownloadString = (Get-Content $LogFile -ErrorAction SilentlyContinue | Where-Object { $_ -like "*Downloader hresult*" }) -as [bool]
        If ($DownloadString) {
            break;
        }
        Start-Sleep -Seconds (60)
    }

    While (1) {
        # check donwload status
        $SuccessfulDownload = (Get-Content $LogFile | Where-Object { $_ -like "*Downloader hresult: 0x0*" }) -as [bool]
        If ($SuccessfulDownload -and !$Downloadflag) {
            Write-Output "Download completed"
            $Downloadflag = $True
        }
        elseif (!$Downloadflag) {
            Start-Sleep 10
            if ($null -ne $(Get-Process -Name Windows10UpgraderApp -ErrorAction Ignore) ) {
                Stop-Process -Name Windows10UpgraderApp -Force -Confirm:$false -ErrorAction Ignore
            }
            Remove-Item $LocalPath -ErrorAction Ignore
            Write-Error "Failed to download the update, please check the log file for more details:`"$LogFile`" "
            Exit;
        }

        # check log status
        $UpdateSuccesful = (Get-Content $LogFile | Where-Object { $_ -like "*WaitForRestartWindows*" }) -as [Bool]
        if ($UpdateSuccesful) {
            if ($UpdateSuccesful) {
                Write-Output "Installation of version 1909 feature update completed successfully. `nPlease reboot the system for changes to take affect, or system will reboot after 30 minutes."
            }
            else {
                Start-Sleep 10
                if ($null -ne $(Get-Process -Name Windows10UpgraderApp -ErrorAction Ignore)) {
                    Stop-Process -Name Windows10UpgraderApp -Force -Confirm:$false -ErrorAction Ignore
                }
                Remove-Item $LocalPath -ErrorAction Ignore
                Write-Error "Feature update has failed. See log `"$LogFile`" for more details."
	    
            }
            exit;
        }
        Start-Sleep -Seconds (5 * 60)
    }

}
catch {
    Write-Output "Failed."
    Write-Error $_
}

