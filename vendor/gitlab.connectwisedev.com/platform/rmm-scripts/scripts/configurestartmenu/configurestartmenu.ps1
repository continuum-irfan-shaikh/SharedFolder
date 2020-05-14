<#
    List of Variable :
    $Removelinksandaccess
    $Removecommonprogram
    $Prohibituser
    $Removeprograms
    $Removenetworkconnections
    $Removesearch
    $Removehelpmenu
    $Removerun
    $Addlogoff
    $Removelogoff
    $Removeshutdownaccess
    $Removedraganddrop
    $Preventchangestotaskbar
    $Removecontextmenu
    $Donotkeepthehistory
    $Turnoffpersonalizedmenus
    $runinseparatememoryspace

    Value to the Variable :
    No Change
    Enable
    Disable
    #>
#####################################################################################
<# Examples of variables. 
$Removelinksandaccess = $null
$Removecommonprogram = $null
$Prohibituser = $null
$Removeprograms = $null
$Removenetworkconnections = 'Enable'
$Removesearch = $null
$Removehelpmenu = $null
$Removerun = $null
$Addlogoff = 'enable'
$Removelogoff = $null
$Removeshutdownaccess = 'Disable'
$Removedraganddrop = $null
$Preventchangestotaskbar = 'enable'
$Removecontextmenu = $null
$Donotkeepthehistory = $null
$Turnoffpersonalizedmenus = $null
$runinseparatememoryspace = $null
#>
#####################################################################################

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64" -and $env:PROCESSOR_ARCHITECTURE -eq 'x86') { $Query = 'c:\windows\sysnative\query.exe' }else { $Query = 'c:\windows\System32\query.exe' }
$LoggedOnUsers = if (($Users = (& $Query user))) {
    $Users | ForEach-Object { (($_.trim() -replace ">" -replace "(?m)^([A-Za-z0-9]{3,})\s+(\d{1,2}\s+\w+)", '$1  none  $2' -replace "\s{2,}", "," -replace "none", $null)) } |
    ConvertFrom-Csv |
    Where-Object { $_.state -ne 'Disc' } |
    Select-Object -expandproperty username
}   

if (!($LoggedOnUsers -eq $null)) {
    Set-Location c:\
    $SID = ((New-Object System.Security.Principal.NTAccount($LoggedOnUsers)).Translate([System.Security.Principal.SecurityIdentifier]).Value)
    $root = "HKU:\$sid\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer" 
    $create = New-PSDrive -PSProvider Registry -Name HKU -Root HKEY_USERS -ErrorAction SilentlyContinue
    if (Test-Path $root) { } else { 
    
        Set-Location "HKU:\$sid\Software\Microsoft\Windows\CurrentVersion\Policies\"
        $create = New-Item Explorer -ErrorAction SilentlyContinue 
        Set-Location "HKU:\$sid\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer"
        $root = Get-Location
    }
    Set-Location $root

    #Remove links and access to Windows update
    if ($Removelinksandaccess -eq "Enable") {
        try {
            $create = New-ItemProperty . -Name NoWindowsUpdate -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = set-ItemProperty . -Name NoWindowsUpdate -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove links and access to Windows update => Enabled"
    }
            
    elseif ($Removelinksandaccess -eq "Disable") {
        try {
            $create = New-ItemProperty . -Name NoWindowsUpdate -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoWindowsUpdate -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove links and access to Windows update => Disabled"
    }
    else { Write-Output "Remove links and access to Windows update => No Changes" }

    #Remove common program groups from start menu
    if ($Removecommonprogram -eq "Enable") {
        try {
            $create = New-ItemProperty . -Name NoCommonGroups -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoCommonGroups -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove common program groups from start menu => Enabled"
    }
    elseif ($Removecommonprogram -eq "Disable") {
        try {
            $create = New-ItemProperty . -Name NoCommonGroups -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoCommonGroups -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove common program groups from start menu => Disabled"
    }
    else { Write-Output "Remove common program groups from start menu => No Changes" }

    #Prohibit user from changing My Documents path
    if ($Prohibituser -eq "Enable") {
        try {
            $create = New-ItemProperty . -Name DisablePersonalDirChange -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name DisablePersonalDirChange -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Prohibit user from changing My Documents path => Enabled"
    }
    elseif ($Prohibituser -eq "Disable") {
        try {
            $create = New-ItemProperty . -Name DisablePersonalDirChange -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name DisablePersonalDirChange -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Prohibit user from changing My Documents path => Disabled"
    }
    else { Write-Output "Prohibit user from changing My Documents path => No Changes" }    

    #Remove programs on settings menu
    if ($Removeprograms -eq "Enable") {
        try {
            $create = New-ItemProperty . -Name NoSetFolders -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoSetFolders -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove programs on settings menu => Enabled"
    }
    elseif ($Removeprograms -eq "Disable") {
        Try {
            $create = New-ItemProperty . -Name NoSetFolders -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoSetFolders -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove programs on settings menu => Disabled"
    }
    else { Write-Output "Remove programs on settings menu => No Changes" }       

    #Remove network connections from start menu
    if ($Removenetworkconnections -eq "Enable") {
        Try {
            $create = New-ItemProperty . -Name NoNetworkConnections -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoNetworkConnections -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove network connections from start menu => Enabled"
    }
    elseif ($Removenetworkconnections -eq "Disable") {
        Try {
            $create = New-ItemProperty . -Name NoNetworkConnections -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoNetworkConnections -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove network connections from start menu => Disabled"
    }
    else { Write-Output "Remove network connections from start menu => No Changes" }    
    
    #Remove search from start menu
    if ($Removesearch -eq "Enable") {
        Try {
            $create = New-ItemProperty . -Name NoFind -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoFind -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove search from start menu => Enabled"
    }
    elseif ($Removesearch -eq "Disable") {
        Try {
            $create = New-ItemProperty . -Name NoFind -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoFind -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove search from start menu => Disabled"
    }
    else { Write-Output "Remove search from start menu => No Changes" }  

    #Remove help menu from start menu
    if ($Removehelpmenu -eq "Enable") {
        Try {
            $create = New-ItemProperty . -Name NoSMHelp -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoSMHelp -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove help menu from start menu => Enabled"
    }
    elseif ($Removehelpmenu -eq "Disable") {
        Try {
            $create = New-ItemProperty . -Name NoSMHelp -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoSMHelp -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove help menu from start menu => Disabled"
    }
    else { Write-Output "Remove help menu from start menu => No Changes" }   

    #Remove run from start menu
    if ($Removerun -eq "Enable") {
        try {
            $create = New-ItemProperty . -Name NoRun -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoRun -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove run from start menu => Enabled"
    }
    elseif ($Removerun -eq "Disable") {
        try {
            $create = New-ItemProperty . -Name NoRun -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoRun -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove run from start menu => Disabled"
    }
    else { Write-Output "Remove run from start menu => No Changes" }     

    #Add logoff to the start menu
    if ($Addlogoff -eq "Enable") {
        Try {
            $create = New-ItemProperty . -Name ForceStartMenuLogOff -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name ForceStartMenuLogOff -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Add logoff to the start menu => Enabled"
    }
    elseif ($Addlogoff -eq "Disable") {
        Try {
            $create = New-ItemProperty . -Name ForceStartMenuLogOff -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name ForceStartMenuLogOff -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Add logoff to the start menu => Disabled"
    }
    else { Write-Output "Add logoff to the start menu => No Changes" }  

    #Remove logoff on the start menu
    if ($Removelogoff -eq "Enable") {
        Try {
            $create = New-ItemProperty . -Name StartMenuLogOff -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name StartMenuLogOff -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove logoff on the start menu => Enabled"
    }
    elseif ($Removelogoff -eq "Disable") {
        Try {
            $create = New-ItemProperty . -Name StartMenuLogOff -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name StartMenuLogOff -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove logoff on the start menu => Disabled"
    }
    else { Write-Output "Remove logoff on the start menu => No Changes" } 

    #Remove and prevent access to the shutdown command
    if ($Removeshutdownaccess -eq "Enable") {
        Try {
            $create = New-ItemProperty . -Name NoClose -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoClose -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove and prevent access to the shutdown command => Enabled"
    }
    elseif ($Removeshutdownaccess -eq "Disable") {
        Try {
            $create = New-ItemProperty . -Name NoClose -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoClose -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove and prevent access to the shutdown command => Disabled"
    }
    else { Write-Output "Remove and prevent access to the shutdown command => No Changes" } 

    #Remove drag-and-drop context menu on the start menu
    if ($Removedraganddrop -eq "Enable") {
        Try {
            $create = New-ItemProperty . -Name NoChangeStartMenu -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoChangeStartMenu -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove drag-and-drop context menu on the start menu => Enabled"
    }
    elseif ($Removedraganddrop -eq "Disable") {
        Try {
            $create = New-ItemProperty . -Name NoChangeStartMenu -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoChangeStartMenu -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove drag-and-drop context menu on the start menu => Disabled"
    }
    else { Write-Output "Remove drag-and-drop context menu on the start menu => No Changes" } 
   
    #Prevent changes to taskbar and start menu settings
    if ($Preventchangestotaskbar -eq "Enable") {
        Try {
            $create = New-ItemProperty . -Name NoSetTaskbar -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        catch {
            $create = Set-ItemProperty . -Name NoSetTaskbar -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Prevent changes to taskbar and start menu settings => Enabled"
    }
    elseif ($Preventchangestotaskbar -eq "Disable") {
        Try {
            $create = New-ItemProperty . -Name NoSetTaskbar -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoSetTaskbar -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Prevent changes to taskbar and start menu settings => Disabled"
    }
    else { Write-Output "Prevent changes to taskbar and start menu settings => No Changes" } 

    #Remove context menu for the taskbar
    if ($Removecontextmenu -eq "Enable") {
        try {
            $create = New-ItemProperty . -Name NoTrayContextMenu -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoTrayContextMenu -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove context menu for the taskbar => Enabled"
    }
    elseif ($Removecontextmenu -eq "Disable") {
        Try {
            $create = New-ItemProperty . -Name NoTrayContextMenu -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoTrayContextMenu -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Remove context menu for the taskbar => Disabled"
    }
    else { Write-Output "Remove context menu for the taskbar => No Changes" } 

    #Do not keep the history of recently opened documents
    if ($Donotkeepthehistory -eq "Enable") {
        Try {
            $create = New-ItemProperty . -Name NoRecentDocsHistory -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoRecentDocsHistory -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Do not keep the history of recently opened documents => Enabled"
    }
    elseif ($Donotkeepthehistory -eq "Disable") {
        try {
            $create = New-ItemProperty . -Name NoRecentDocsHistory -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name NoRecentDocsHistory -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Do not keep the history of recently opened documents => Disabled"
    }
    else { Write-Output "Do not keep the history of recently opened documents => No Changes" }

    #Turn off personalized menus
    if ($Turnoffpersonalizedmenus -eq "Enable") {
        try {
            $create = New-ItemProperty . -Name Intellimenus -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name Intellimenus -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Turn off personalized menus => Enabled"
    }
    elseif ($Turnoffpersonalizedmenus -eq "Disable") {
        Try {
            $create = New-ItemProperty . -Name Intellimenus -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name Intellimenus -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Turn off personalized menus => Disabled"
    }
    else { Write-Output "Turn off personalized menus => No Changes" }    
   
    #Add 'run in separate memory space' check box to run dialog box
    if ($runinseparatememoryspace -eq "Enable") {
        Try {
            $create = New-ItemProperty . -Name MemCheckBoxInRunDlg -Value 1 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name MemCheckBoxInRunDlg -Value 1 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Add 'run in separate memory space' check box to run dialog box => Enabled"
    }
    elseif ($runinseparatememoryspace -eq "Disable") {
        Try {
            $create = New-ItemProperty . -Name MemCheckBoxInRunDlg -Value 0 -PropertyType DWORD -ErrorAction Stop
        }
        Catch {
            $create = Set-ItemProperty . -Name MemCheckBoxInRunDlg -Value 0 -Type DWORD -ErrorAction SilentlyContinue
        }
        Write-Output "Add 'run in separate memory space' check box to run dialog box => Disabled"
    }
    else { Write-Output "Add 'run in separate memory space' check box to run dialog box => No Changes" }   

    Write-Output "`nNote:"
    Write-Output "You might need to reboot the system to make the changes effective"

}
else {
    Write-Output "This script requires logon user and currently no user is logged in."
}
