<#
    .SYNOPSIS
        Create Shortcut
    .DESCRIPTION
        Create Shortcut
    .HELP

    .AUTHOR
        Durgeshkumar Patel
    .VERSION
        1.0
#>

<#############GUI Parameters########
CheckBox 
     CreateShortcutForAllUsers
Type
     Shell ShortCut
     Web Link
Actiona
     Create
     Delete
CheckBox
    OverwriteExistingShortcut
Shortcut Name
Shortcut Location
Create Shortcut for Target
Target Arguments
Start In Folder
Shortcut Comments
#>
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

####################GUI Parameters#########################
<#$CreateShortcutForAllUsers = $true  #$true/$false
$Type = "Shell Shortcut"   #Shell Shortcut / Web Link
$Action = "Create"    #Create / Delete
$OverwriteExistingShortcut = $true    #$true or $false
$ShortcutName = "Notepad6.lnk"
$ShortcutLocation = "User Desktop"
$CreateShortcutforTarget = "$env:SystemRoot\System32\notepad.exe"   #"$env:SystemRoot\System32\notepad.exe"   #http://google.com
$TargetArguments = ""
$StartInFolder = "C:\Users\Labadmin\Desktop\test"
$ShortcutComments = "Some Comment"    #Comments #>
########################################################
# Create Shell Shortcut
function Create_Shell_Shortcut() { 
    
    try {
        $ErrorActionPreference = "Stop"
        #Join Path
        $ShortcutLocation_path = join-path -path "$($Hash_Result["$ShortcutLocation"])" -ChildPath "$ShortcutName" -ErrorAction "SilentlyContinue"  
        $word_result = CheckAlreadyExist
        #-ComObject WScript.Shell: This creates an instance of the COM object that represents the WScript.Shell for invoke CreateShortCut
        $WScriptShell = New-Object -ComObject WScript.Shell
        $Shortcut = $WScriptShell.CreateShortcut($ShortcutLocation_path)
        $Shortcut.TargetPath = "$CreateShortcutforTarget"
        $Shortcut.DESCRIPTION = "$ShortcutComments"
        $Shortcut.Arguments = "$TargetArguments"
        $Shortcut.WorkingDirectory = "$StartInFolder"
        #Save the Shortcut to the TargetPath
        $Shortcut.Save()
        Write-Output "$Type $ShortcutName $word_result successfully for user $U."
    }
    catch {
        $ErrorActionPreference = "Continue"
        Write-Error "$Type $ShortcutName creation failed for user $U."
        Write-Error $_.Excpetion.Message

    }
}

# Create Web Link
function Create_Web_Link() { 
    
    try {
        $ErrorActionPreference = "Stop"
        #Join Path
        $ShortcutLocation_path = join-path -path "$($Hash_Result["$ShortcutLocation"])" -ChildPath "$ShortcutName" -ErrorAction "SilentlyContinue"
        $word_result = CheckAlreadyExist
        #-ComObject WScript.Shell: This creates an instance of the COM object that represents the WScript.Shell for invoke CreateShortCut
        $WScriptShell = New-Object -ComObject WScript.Shell
        $Shortcut = $WScriptShell.CreateShortcut($ShortcutLocation_path)
        $Shortcut.TargetPath = "$CreateShortcutforTarget"
        #Save the Shortcut to the TargetPath
        $Shortcut.Save()
        Write-Output "'$Type' $ShortcutName $word_result successfully for user '$U'."
    }
    catch {
        $ErrorActionPreference = "Continue"
        Write-Error "'$Type' $ShortcutName creation failed for user '$U'."
        Write-Error $_.Excpetion.Message
    }
}

#Hash table for user default paths
function HashTable($PM) {
    $ShortcutLocation_Hash = @{ 
        "User Desktop"          = "C:\Users\$U\Desktop";  
        "User Favourites"       = "C:\Users\$U\Favorites";
        "User Start Menu"       = "C:\Users\$U\AppData\Roaming\Microsoft\Windows\Start Menu";
        "User Programs Group"   = "C:\Users\$U\AppData\Roaming\Microsoft\Windows\Start Menu\Programs";
        "User Startup Group"    = "C:\Users\$U\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup";
        "User Send To"          = "C:\Users\$U\AppData\Roaming\Microsoft\Windows\SendTo";
        "User Quick Launch Bar" = "C:\Users\$U\AppData\Roaming\Microsoft\Internet Explorer\Quick Launch" 
    }
    return $ShortcutLocation_Hash
} 

#Check shortcut already exist.
function CheckAlreadyExist {
    $word = "created"
    if ((Test-Path "$ShortcutLocation_path") -and ($OverwriteExistingShortcut -eq $true)) {
        $word = "overwritten"
    }
    return $word      
}

#Append .lnk if user does not provide shortcut name correctly for Shell Shortcut
if ($ShortcutName -notmatch ".lnk") {
    $ShortcutName = "$ShortcutName" + ".lnk"
}

if ($Action -ne "Delete") {
    #Exit if CreateShortcutforTarget is not correct for Shell Shortcut
    if ((!(Test-Path "$CreateShortcutforTarget")) -and ($Type -eq "Shell Shortcut")) {
        Write-Output "Shortcut for Target '$CreateShortcutforTarget' path is not correct."
        exit;
    }

    #Exit if CreateShortcutforTarget URI is not correct for Web Link
    if ((!(([System.URI]$CreateShortcutforTarget).AbsoluteURI -ne $null -and ([System.URI]$CreateShortcutforTarget).Scheme -match 'http|https')) -and ($Type -eq "Web Link") ) {
        Write-Output "Shortcut for Target '$CreateShortcutforTarget' path is not correct."
        exit;
    }

    #Exit if Start in folder path is not correct
    if (![string]::IsNullOrEmpty("$StartInFolder")) {
        if ((!(Test-Path "$StartInFolder")) -and ($Type -eq "Shell Shortcut")) {
            Write-Output "Start in folder '$StartInFolder' path is not correct."
            exit; 
        }
    }
}



#Check for all users or a logged in user.
if ($CreateShortcutForAllUsers -eq $true) {
    #Find all local users
    $LoggedOnUsers = Get-WmiObject -Class Win32_UserAccount -Filter "LocalAccount=True" | where-Object { $_.name -ne "Guest" } | Select-Object -ExpandProperty name
}
else {
    #Find currently logged in user
    if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64" -and $env:PROCESSOR_ARCHITECTURE -eq 'x86') { $Query = 'c:\windows\sysnative\query.exe' } else { $Query = 'c:\windows\System32\query.exe' }
    $LoggedOnUsers = if (($Users = (& $Query user))) {
        $Users | ForEach-Object { (($_.trim() -replace ">" -replace "(?m)^([A-Za-z0-9]{3,})\s+(\d{1,2}\s+\w+)", '$1  none  $2' -replace "\s{2,}", "," -replace "none", $null)) } |
        ConvertFrom-Csv |
        Where-Object { $_.state -ne 'Disc' } |
        Select-Object -expandproperty username
    }   
}

#Create Shell Shortcut
if (($Action -eq "Create") -and ($Type -eq "Shell Shortcut")) {

    foreach ($U in $LoggedOnUsers) {
        $Hash_Result = HashTable -PM $U
        if (($OverwriteExistingShortcut -ne $true) -and (Test-Path $(join-path -path "$($Hash_Result["$ShortcutLocation"])" -ChildPath "$ShortcutName" -ErrorAction "SilentlyContinue"))) {

            Write-Output "$ShortcutName already exist for user $U. Kindly choose overwrite option."
        }
        else {
            Create_Shell_Shortcut    
        }   
    }
}

#Create Web Link
if (($Action -eq "Create") -and ($Type -eq "Web Link")) {

    foreach ($u in $LoggedOnUsers) {
        $Hash_Result = HashTable -PM $U
        if (($OverwriteExistingShortcut -ne $true) -and (Test-Path $(join-path -path "$($Hash_Result["$ShortcutLocation"])" -ChildPath "$ShortcutName" -ErrorAction "SilentlyContinue"))) {

            Write-Output "$ShortcutName already exist for user $U. Kindly choose overwrite option."
        }
        else {
            Create_Web_Link    
        }
    }
}

#Delete Shortcut
if ($Action -eq "Delete") {
    try {
        $ErrorActionPreference = "Stop"    
        foreach ($U in $LoggedOnUsers) {
            $Hash_Result = HashTable -PM $U
            $ShortcutLocation_path = join-path -path "$($Hash_Result["$ShortcutLocation"])" -ChildPath "$ShortcutName" -ErrorAction "SilentlyContinue"
            if (Test-Path "$ShortcutLocation_path") {
                Remove-Item -Path "$ShortcutLocation_path" -Force
                Write-Output "$Type $ShortcutName deleted successfully for user '$U'."
            }
            else {
                Write-Output "Failed to delete $Type. Check '$ShortcutLocation_path' or '$ShortcutName' for user $U."
            }
        }
    }
    catch {
        $ErrorActionPreference = "Continue"
        Write-Error "Failed to delete '$Type' for user '$U'"
        Write-Error $_.Excpetion.Message
    }
}
