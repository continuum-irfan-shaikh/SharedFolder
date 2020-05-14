<#
    .SYNOPSIS
        Install .Net Framework v4.8
    .DESCRIPTION
        Install .Net Framework v4.8
    .Help
        Refer below links for more details.
        https://go.microsoft.com/fwlink/?linkid=2088631
        https://docs.microsoft.com/en-us/dotnet/framework/deployment/guide-for-administrators
    .Author
        Durgeshkumar Patel
    .Version
        1.0
    .Note 
        This is script currently designed to install only .net 4.8 version.
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

$url = "https://go.microsoft.com/fwlink/?linkid=2088631" 
$downloadpath = "C:\Windows\Temp\DotNetFramework4.8.exe"

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
        Write-Output ".Net Framework v4.8 installation file download has been failed. Check internet connection and try again."
        Exit;
    }
}

Function GetDotNetVersion {
    $version = $null
    if (Test-Path -Path "HKLM:SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full") {
   
        switch ((Get-ItemProperty -Path "HKLM:SOFTWARE\Microsoft\NET Framework Setup\NDP\v4\Full" -ErrorAction SilentlyContinue).Release) {
            528040 { $Version = "4.8" }
            528049 { $Version = "4.8" }
        }  
    }
    return $Version
}

 #Check if installed version of .Net Framework is v4.8
 $GetVersion = GetDotNetVersion
    if ($GetVersion) {
        ".Net Framework v{0} is already installed on system {1}" -f $GetVersion, $Env:COMPUTERNAME
        Exit;
    }

try {
        #Download the file for .Net Framework v4.8
        if (downloadfile) {

            #Install process
            $process = Start-Process $downloadpath -arg "/q /norestart" -Wait -PassThru -ErrorAction 'Stop'
            
            switch ($process.exitcode) {

                #ExitCodes 
                0 { Write-Output ".Net Framework v4.8 Installation completed successfully on system $ENV:COMPUTERNAME." }
                1602 { Write-Output  "The user canceled installation of .Net Framework v4.8 on system $ENV:COMPUTERNAME." }
                1603 { Write-Output "A fatal error occurred during installation of .Net Framework v4.8 on system $ENV:COMPUTERNAME." }
                1641 { Write-Output "A restart is required to complete the installation of .Net Framework v4.8 on $ENV:COMPUTERNAME." }
                3010 { Write-Output "A restart is required to complete the installation of .Net Framework v4.8 on $ENV:COMPUTERNAME." }
                5100 { Write-Output "The computer $ENV:COMPUTERNAME does not meet system requirements to install .Net Framework v4.8." }
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
