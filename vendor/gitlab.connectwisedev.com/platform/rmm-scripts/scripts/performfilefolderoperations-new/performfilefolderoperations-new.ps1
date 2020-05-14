# Move Operation On Error Exit, Ignore Attribute criteria

<#
     .Scriptname 
     Perform File folder operation
     .Author
     Nirav Sachora
     .Description
     Script will perform file folder operation as mentioned in below parameters.
     Script can perform file folder operation over the network as well, to perform over network user has to provide valid share name and username,password.
     .Requirement
     Script should run with highest privileges.
#>

<#
$fileAction = "deletefiles"
    "copyafile
    copyafolder
    renameafile
    moveafile
    renameafolder
    movefolder
    renameafolder
    deleteafile
    deletefiles
    deleteafolder
    createfolder"
$sourcepath = "Network"
$sourcefile = "\\10.2.19.127\Nirav\Test\abc.txt,\\10.2.19.127\Nirav\Test\pqr.txt,\\10.2.19.127\Nirav\Test\xyz.txt"#\abc.txt,D:\Nirav\pqr.txt"
$destinationpath = "Local"
$destinationfile = "E:\Destination\Move"#\\10.2.130.81\Shared Files\PSTools" #or \\192.168.2.1\sharename\folder\file.txt
$newname = "pqr.txt"
$username = "admin"
$password = "Welcome@123"
$includesub = "false"
$includeSystemFile = $false
$includeReadOnlyFile = $false
$includeHiddenFile = $true
$Overwrite = $true
$CreateDestDirectory = $true
    
$continueonerror = "false"
$fileModification = "lastAccessedbefore"
$modificationDate = "09/25/2019"
$modificationDays = -1
$startDate = "08/01/2019"
$endDate = "08/31/2019"
<#modifiedBefore                 modificationDate         
    modifiedAfter                  modificationDate  
    modifiedOn                     modificationDate  
    modifiedBetween                startDate    endDate
    modifiedOlderthanXdays         modificationDays
    lastAccessedBefore             modificationDate
    lastAccessedAfter              modificationDate 
    lastAccessedOn                 modificationDate  
    lastAccessedBetween            startDate    endDate
    createdOlderthanXdays          modificationDays
  #>  
<#   
    if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
        if ($myInvocation.Line) {
            &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
        }
        else {
            &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
        }
        exit $lastexitcode
    }
 #> #>
 

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
} 

#########################Function to map network drive#######################################################################################
$global:newdrive | Out-Null
$global:ogdrivedetails | out-Null
function Map_networkdrive($action, $Drive) {
    $Net = New-Object -ComObject WScript.Network
    $DriveLetter = Get-ChildItem function:[g-z]: -n | Where-Object { !(Test-Path $_) } | random 
    
    if ($action -eq "Create") {
        if (($sourcepath -eq 'Network') -and (($fileAction -eq "copyfilestoafolder") -or ($fileAction -eq "deletefiles"))) {
            $DriveDetails = (([regex]::Match($sourcefile[0], "\\\\[\S\s]*?\\[\S\s]*?\\")).Value).TrimEnd("\\")
            $ErrorActionpreference = "SilentlyContinue"
            $Net.MapNetworkDrive($DriveLetter, "$DriveDetails", $false, $UserName, $Password)
            $global:ogdrivedetails = $DriveDetails
            $ErrorActionpreference = "Continue"
            if (Get-WmiObject -Class Win32_MappedLogicalDisk | Where-Object { $_.name -eq "$DriveLetter" }) {
                $updatedsourcefile = foreach ($file in $sourcefile) {
                    -join ($DriveLetter, "\", ($file.Replace($DriveDetails, "").TrimStart("\\")))
                }
                $global:newdrive = $DriveLetter
                return $updatedsourcefile
            }
            else {
                "-" * 30 + "`nDrive mapping failed`nUsername or password is incorrect or remote system is not accessible over network`n" + "-" * 30
                Exit;
            }
            
            
        }
        elseif ($sourcepath -eq 'Network') {
            $DriveDetails = (([regex]::Match($sourcefile, "\\\\[\S\s]*?\\[\S\s]*?\\")).Value).TrimEnd("\\")
            #$ErrorActionpreference = "SilentlyContinue"
            $Net.MapNetworkDrive($DriveLetter, "$DriveDetails", $false, $UserName, $Password)
            #$ErrorActionpreference = "Continue"
            if (Get-WmiObject -Class Win32_MappedLogicalDisk | Where-Object { $_.name -eq "$DriveLetter" }) {
                $updatedsourcefile = -join ($DriveLetter, "\", ($sourcefile.Replace($DriveDetails, "").TrimStart("\\")))
                $global:ogdrivedetails = $DriveDetails
                $global:newdrive = $DriveLetter
                return $updatedsourcefile
            }
            else {
                "-" * 30 + "`nDrive mapping failed`nUsername or password is incorrect or remote system is not accessible over network`n" + "-" * 30
                Exit;
            }
            
        }
        elseif ($destinationpath -eq 'Network') {
            
            $DriveDetails = (([regex]::Match($destinationfile, "\\\\[\S\s]*?\\[\S\s]*?\\")).Value).TrimEnd("\\")
            if (!$DriveDetails) { $DriveDetails = $destinationfile }
            $ErrorActionpreference = "SilentlyContinue"
            $Net.MapNetworkDrive($DriveLetter, "$DriveDetails", $false, $UserName, $Password)
            $ErrorActionpreference = "Continue"
            if (Get-WmiObject -Class Win32_MappedLogicalDisk | Where-Object { $_.name -eq "$DriveLetter" }) {
                $updateddestinationfile = -join ($DriveLetter, "\", ($destinationfile.Replace($DriveDetails, "").TrimStart("\\")), "\")
                $global:newdrive = $DriveLetter
                return $updateddestinationfile
            }
            else {
                "-" * 30 + "`nDrive mapping failed`nUsername or password is incorrect or remote system is not accessible over network`n" + "-" * 30
                Exit;
            }
            
        }
   
    }
    elseif ($action -eq "Remove") {
        $ErrorActionpreference = "Stop"
        if (($sourcepath -eq "Local") -and ($destinationpath -eq "Local")) { return }
        if (($sourcepath -eq "Local") -and !($destinationpath)) { return }
        try {
            $Net.RemoveNetworkDrive($Drive) 
        }
        catch {
            $_.Exception.Message
        }
    }
}

######################################################################################################################################################

#########################Function to check attributes based on input parameters given######################################################## 
    
function verify_Attributes($sourcepath, $readonly, $hiddens, $system) {
    $fileattributes = ((Get-Item $sourcepath -Force).Attributes) -split ", "
    if (($readonly -eq $false) -and ($hiddens -eq $false) -and ($system -eq $false)) { $attributes = "000" }
    elseif (($readonly -eq $true) -and ($hiddens -eq $false) -and ($system -eq $false)) { $attributes = "100" }
    elseif (($readonly -eq $false) -and ($hiddens -eq $true) -and ($system -eq $false)) { $attributes = "010" }
    elseif (($readonly -eq $false) -and ($hiddens -eq $false) -and ($system -eq $true)) { $attributes = "001" }
    elseif (($readonly -eq $true) -and ($hiddens -eq $true) -and ($system -eq $false)) { $attributes = "110" }
    elseif (($readonly -eq $false) -and ($hiddens -eq $true) -and ($system -eq $true)) { $attributes = "011" }
    elseif (($readonly -eq $true) -and ($hiddens -eq $false) -and ($system -eq $true)) { $attributes = "101" }
    elseif (($readonly -eq $true) -and ($hiddens -eq $true) -and ($system -eq $true)) { $attributes = "111" }
        
    switch ($attributes) {
        "000" {
            if (($fileattributes -contains "ReadOnly") -or ($fileattributes -contains "Hidden") -or ($fileattributes -contains "System")) {
                return $false
            }
            else { return $true }
        }
        "100" {
            if (($fileattributes -contains "Hidden") -or ($fileattributes -contains "System")) {
                return $false 
            }
            else { return $true }
        }
        "010" {
            if (($fileattributes -contains "ReadOnly") -or ($fileattributes -contains "System")) {
                return $false   
            }
            else { return $true }
        }
        "001" {
            if (($fileattributes -contains "ReadOnly") -or ($fileattributes -contains "Hidden")) {
                return $false 
            }
            else { return $true }
        }
        "110" {
            if ($fileattributes -contains "System") {
                return $false 
            }
            else { return $true }
        }
        "011" {
            if ($fileattributes -contains "ReadOnly") {
                return $false 
            }
            else { return $true }
        }
        "101" {
            if ($fileattributes -contains "Hidden") {
                return $false 
            }
            else { return $true }
        }
        "111" {
            return $false 
        }
        
    }
}

###############################################################################################################################

##########################Function to check modification based on criteria provided##########################################################
    
function filemodification_check($ffc_sourcepath, $ffc_criteria, $ffc_modificationdate, $ffc_modificationdays, $ffc_startdate, $ffc_enddate) {
        
    if ($ffc_modificationdate -ne $null) { $modificationdate = get-date -date $ffc_modificationdate }
    if ($ffc_startdate -ne $null) { $startdate = get-date -date $ffc_startdate }
    if ($ffc_enddate -ne $null) { $enddate = get-date -date $ffc_enddate }
    
    switch ($ffc_criteria) {
        "modifiedBefore" {
            if ((get-item $ffc_sourcepath -Force).LastWriteTime -lt $modificationdate) { return $true }else { return $false }
        }
        "modifiedAfter" {
            if ((get-item $ffc_sourcepath -Force).LastWriteTime -gt $modificationdate) { return $true }else { return $false }
        }
        "modifiedOn" {
            if ((((get-item $ffc_sourcepath -Force).LastWriteTime | Get-Date).Date) -eq (($modificationdate | Get-Date).Date)) { return $true }else { return $false }
        }
        "modifiedBetween" {
            if (((get-item $ffc_sourcepath -Force).LastWriteTime -gt $startdate) -and ((get-item $ffc_sourcepath -Force).LastWriteTime -lt $enddate)) { return $true }else { return $false }
        }
        "modifiedOlderthanXdays" {
            $lastwritetime = get-date -date (Get-Item $ffc_sourcepath -Force).LastWriteTime
            $currentdate = Get-Date
            $days = $currentdate - $lastwritetime
            if ($days.Days -gt $ffc_modificationdays) {
                return $true
            }
            else { return $false }
        }
        "lastAccessedBefore" {
            if ((get-item $ffc_sourcepath -Force).LastAccessTime -lt $modificationdate) { return $true }else { return $false }
        }
        "lastAccessedAfter" {
            if ((get-item $ffc_sourcepath -Force).LastAccessTime -gt $modificationdate) { return $true }else { return $false }
        }
        "lastAccessedOn" {
            if ((((get-item $ffc_sourcepath -Force).LastAccessTime | Get-Date).Date) -eq (($modificationdate | Get-Date).Date)) { return $true }else { return $false }
        }
        "lastAccessedBetween" {
            if (((get-item $ffc_sourcepath -Force).LastAccessTime -gt $startdate) -and ((get-item $ffc_sourcepath -Force).LastAccessTime -lt $enddate)) { return $true }else { return $false }
        }
        "createdOlderthanXdays" {
            $lastaccesstime = get-date -date (Get-Item $ffc_sourcepath -Force).CreationTime
            $currentdate = Get-Date
            $days = $currentdate - $lastaccesstime
            if ($days.Days -gt $ffc_modificationdays) {
                return $true
            }
            else { return $false }
        }
    }
}
 
############################################################################################################################################# 

########################################### Precheck ######################################

if (($sourcepath -eq "Network") -and ($destinationpath -eq "Network")) {
    Write-Error "Input Error:Source and destination path can't be network"
    Exit;
}

########################################### Precheck Multiple files ######################################

if (($fileAction -eq "copyfilestoafolder") -or ($fileAction -eq "deletefiles")) {
    $sourcefile = $sourcefile -split ","

    if (-not($sourcefile.Gettype().Name -eq "String[]")) {
        Write-Error "Error: Unknown Error"
        Exit
    }

    if ($sourcepath -eq "Network") {
        $sourcefile = Map_networkdrive -Action "Create"
    }

    if ($destinationpath -eq "Network") {
        $destinationfile = Map_networkdrive -Action "Create"
    }

    ################Check whether sourcefiles name is successfully converted to an array#############

    $verifiedpath = @()
    $discardedwhilesourcecheck = @()
    foreach ($validatepath in $sourcefile) {
        
        if (((Test-path $validatepath -IsValid) -eq $true) -And ((Test-path $validatepath) -eq $true)) {    
            $verifiedpath += $validatepath
        }
        else { $discardedwhilesourcecheck += $validatepath }
    }

    $sourcefile = $verifiedpath

    
    ########################Completed Sourcefiles check##############################################

    #########################Check valid file#########################################################
    $isvalidfilename = @()
    $isnotvalidfilename = @()

    foreach ($file in $sourcefile) {

        if (((Get-Item $file -Force).Attributes) -ne "Directory") {
    
            $isvalidfilename += $file
    
        }
        else { $isnotvalidfilename += $file }

    }

    $sourcefile = $isvalidfilename
    ##################################################################################################

    ######################################################################################################

    if (($fileAction -eq "copyfilestoafolder") -and (!$Overwrite) -and (-not(!$sourcefile))) {
        $filetobecopied = @()
        $filenottobecopied = @()
        foreach ($file in $sourcefile) {
            $filename = Split-Path $file -Leaf
            $temppath = "$destinationfile" + "\" + "$filename"
            if (!(Test-path $temppath)) { $filetobecopied += $file }else { $filenottobecopied += $file }
        }
        $sourcefile = $filetobecopied
    }

    ########################################################################################################
    
    ####################Checking Attribute criteria for each file####################################
    
    if (((!$includeSystemFile) -or (!$includeReadOnlyFile) -or (!$includeHiddenFile)) -and (-not(!$sourcefile))) {
        $attributeverifiedfiles = @()
        $attributenotverifiedfiles = @()
        foreach ($verifyfile in $sourcefile) {
    
            if (verify_Attributes -sourcepath $verifyfile -readonly $includeReadOnlyFile -hiddens $includeHiddenFile -system $includeSystemFile) {
    
                $attributeverifiedfiles += $verifyfile
    
            }
            else { $attributenotverifiedfiles += $verifyfile }
        }
        $sourcefile = $attributeverifiedfiles
    }

    #####################################################################################################

    ####################Checking Modification criteria for each file#####################################

    if (($fileModification) -and (-not(!$sourcefile))) {
        $filemodificationverifiedfiles = @()
        $filemodificationnotverifiedfiles = @()
        foreach ($verifyfile in $sourcefile) {
    
            if (filemodification_check -ffc_sourcepath $verifyfile -ffc_criteria $fileModification -ffc_modificationdate $modificationDate -ffc_modificationdays $modificationDays -ffc_startdate $startDate -ffc_enddate $endDate) {
    
                $filemodificationverifiedfiles += $verifyfile
    
            }
            else {
    
                $filemodificationnotverifiedfiles += $verifyfile

            }

        }
        $sourcefile = $filemodificationverifiedfiles
    }
 

    if (!$sourcefile) {

        Write-Error "Condition mismatched for all the files"
        if ($sourcepath -eq "Network") { Map_networkdrive -Action "Remove" -Drive $global:newdrive }
        if ($destinationpath -eq "Network") { Map_networkdrive -Action "Remove" -Drive $global:newdrive }
        Exit;
    
    }
    else {
    
        function rollback_source($filestoupdate) {
    
            $rollbacked = foreach ($file in $filestoupdate) {
                Join-Path $global:ogdrivedetails (split-path $file -NoQualifier)
            }
            return $rollbacked

        }

        if ($discardedwhilesourcecheck) {
            if ($sourcepath -eq "Network") {
                $Global:discardedwhilesourcecheck = rollback_source -filestoupdate $discardedwhilesourcecheck
            }
             
        }
        if ($isnotvalidfilename) {
            if ($sourcepath -eq "Network") {
                $global:isnotvalidfilename = rollback_source -filestoupdate $isnotvalidfilename
            }
             
        }
        if ($attributenotverifiedfiles) { 
            if ($sourcepath -eq "Network") {
                $global:attributenotverifiedfiles = rollback_source -filestoupdate $attributenotverifiedfiles
            }
            
        }
        if ($filemodificationnotverifiedfiles) { 
            if ($sourcepath -eq "Network") {
                $global:filemodificationnotverifiedfiles = rollback_source -filestoupdate $filemodificationnotverifiedfiles
            }
            
        }
        if ($filenottobecopied -and !($fileaction -eq "deletefiles")) {
            if ($sourcepath -eq "Network") {
                $global:filenottobecopied = rollback_source -filestoupdate $filenottobecopied
            }
            
        }
    } 

}
elseif (($fileAction -eq "copyafile") -or ($fileAction -eq "moveafile") -or ($fileAction -eq "deleteafile") -or ($fileAction -eq "renameafile")) {
    
    function rollback_source1($filetoupdate) {  
        return  (Join-Path $global:ogdrivedetails (split-path $filetoupdate -NoQualifier))
    }

    function remove_net_drive {
        if ($sourcepath -eq "Network") { $tempdrive = Split-path $sourcefile -Qualifier; Map_networkdrive -Action "Remove" -Drive $tempdrive }
        if ($destinationpath -eq "Network") { $tempdrive = Split-path $destinationfile -Qualifier; Map_networkdrive -Action "Remove" -Drive $tempdrive }
    }

    if ($sourcepath -eq "Network") {
        $sourcefile = Map_networkdrive -Action "Create"
    }
    
    if ($destinationpath -eq "Network") {
        $destinationfile = Map_networkdrive -Action "Create"
    }
    if (!(((Test-path $sourcefile -IsValid) -eq $true) -And ((Test-path $sourcefile) -eq $true))) {
        remove_net_drive
        if ($sourcepath -eq "Network") {
            $sourcefile = rollback_source1 -filetoupdate $sourcefile
        }    
        Write-Error "Source File not found" 
        Exit
    }
    
    if (((Get-Item $sourcefile -ErrorAction SilentlyContinue).Attributes) -eq "Directory") {
        remove_net_drive
        if ($sourcepath -eq "Network") {
            $sourcefile = rollback_source1 -filetoupdate $sourcefile
        }
        Write-Error "Invalid file"
        Exit;
    }
    
    if ((($fileAction -eq "copyafile") -or ($fileAction -eq "moveafile")) -And (!$Overwrite)) {
        $filename = Split-Path $sourcefile -Leaf
        $temppath = "$destinationfile" + "\" + "$filename"
        if (Test-path $temppath) {
            remove_net_drive 
            if ($sourcepath -eq "Network") {
                $sourcefile = rollback_source1 -filetoupdate $sourcefile
            }
            Write-Error "File already present at destination path" 
            Exit; 
        }
    }
    
    
    
    if ($fileModification) {
        if (!(filemodification_check -ffc_sourcepath $sourcefile -ffc_criteria $fileModification -ffc_modificationdate $modificationDate -ffc_modificationdays $modificationDays -ffc_startdate $startDate -ffc_enddate $endDate)) {
            remove_net_drive
            if ($sourcepath -eq "Network") {
                $sourcefile = rollback_source1 -filetoupdate $sourcefile
            }
            Write-Error "File modification criteria mismatched"
            Exit;
        }
    }
    
    if (((!$includeSystemFile) -or (!$includeReadOnlyFile) -or (!$includeHiddenFile)) -and ($fileAction -ne "renameafile")) {
        if (!(verify_Attributes -sourcepath $sourcefile -readonly $includeReadOnlyFile -hiddens $includeHiddenFile -system $includeSystemFile)) { 
            remove_net_drive
            if ($sourcepath -eq "Network") {
                $sourcefile = rollback_source1 -filetoupdate $sourcefile
            }
            Write-Error "Attribute criteria mismatched"
            Exit; 
        }
    }
    
    
}

elseif (($fileAction -eq "copyafolder") -or ($fileAction -eq "movefolder") -or ($fileAction -eq "renameafolder") -or ($fileAction -eq "deleteafolder")) {
    
    function rollback_source2($foldertoupdate) {  
        return  (Join-Path $global:ogdrivedetails (split-path $foldertoupdate -NoQualifier))
    }

    if (($fileAction -eq "movefolder") -and (($sourcepath -eq "Network") -or ($destinationpath -eq "Network"))) {
        Write-Error "Move Folder operation is not allowed over the network."
        Exit;
    }

    function remove_net_drive {
        if ($sourcepath -eq "Network") { $tempdrive = Split-path $sourcefile -Qualifier; Map_networkdrive -Action "Remove" -Drive $tempdrive }
        if ($destinationpath -eq "Network") { $tempdrive = Split-path $destinationfile -Qualifier; Map_networkdrive -Action "Remove" -Drive $tempdrive }
    }

    if ($sourcepath -eq "Network") {
        $sourcefile = Map_networkdrive -Action "Create"
    }
    
    if ($destinationpath -eq "Network") {
        $destinationfile = Map_networkdrive -Action "Create"
    }

    if (!(((Test-path $sourcefile -IsValid) -eq $true) -And ((Test-path $sourcefile) -eq $true))) { 
        remove_net_drive
        if ($sourcepath -eq "Network") {
            $sourcefile = rollback_source2 -foldertoupdate $sourcefile
        }   
        Write-Error "Sourcefolder not found"     
        Exit
    }

    if (((Get-Item $sourcefile).Attributes) -notlike "Directory*") {
        remove_net_drive
        if ($sourcepath -eq "Network") {
            $sourcefile = rollback_source2 -foldertoupdate $sourcefile
        } 
        Write-Error "Invalid Folder"
        Exit;
    }


    if ($fileModification) {
        if (!(filemodification_check -ffc_sourcepath $sourcefile -ffc_criteria $fileModification -ffc_modificationdate $modificationDate -ffc_modificationdays $modificationDays -ffc_startdate $startDate -ffc_enddate $endDate)) {
            remove_net_drive
            if ($sourcepath -eq "Network") {
                $sourcefile = rollback_source2 -foldertoupdate $sourcefile
            } 
            Write-Error "Folder modification criteria mismatched"
            Exit;
        }
    }
}


#####################################Precheck completed multiple files#######################################################



#####################################Standard Function to create destination path############################################

function Create_Destinationpath($destinationpath) {
    $ErrorActionPreference = "SilentlyContinue"
    New-Item $destinationpath -ItemType Directory | Out-Null
    $ErrorActionPreference = "Continue"
    if (Test-path $destinationpath) {
        return $true
    }
    else {
        Write-Error "Failed to create destination path;"
        Exit;
    }
}

##############################################################################################################################

###################################Copy files to folder ######################################################################

<# 
   Calling a function requires 4 parameter
   1. $sourcefilespath :       $sourcefile
   2. $destinationfolderpath : $destinationfile
   3. $cf_continueonerror:     $continueonerror
   4. $cf_destinationdirectory: $CreateDestDirectory

 #>       
function copyfilestoafolder($sourcefilespath, $destinationfolderpath, $cf_continueonerror, $cf_destinationdirectory) {
    $copysuccessful = @()
    $copyfailed = @() 
    if ((Test-path $destinationfolderpath) -eq $false) {
        #it will check for destination path if it does not exist.
        if (($cf_destinationdirectory -eq $false)) {
            # and create destination directory is true than new directory will be created.
            Write-Error "Destination path provided does not exist"
            Exit;
        }
        else {
            Create_Destinationpath -destinationpath $destinationfolderpath | Out-Null
        }
    }
       
    foreach ($path in $sourcefilespath) {
        Copy-item $path -Destination $destinationfolderpath -ErrorAction SilentlyContinue | out-null
        $destinationpath1 = Join-path $destinationfolderpath (split-path $path -Leaf)
        if (Test-path $destinationpath1) {
            $copysuccessful += $path
        }
        else {
            $copyfailed += $path
        }
    }

    if ($sourcepath -eq "Network" -and $copysuccessful.length -gt 0) {
        $op_success = @()
        foreach ($file in $copysuccessful) {
            $op_success += Join-Path $global:ogdrivedetails (split-path $file -NoQualifier)
        }
        "`nCopy Successful:"; "-" * 30; $op_success; "-" * 30
    }

    if ($sourcepath -eq "Network" -and $copyfailed.length -gt 0) {
        $op_failed = @()
        foreach ($file in $copyfailed) {
            $op_failed += Join-Path $global:ogdrivedetails (split-path $file -NoQualifier)
        }
        Write-Error "`nCopy Failed:`n$op_failed"
    }
    
    if ($sourcepath -eq "Local") {
        if ($copysuccessful.length -gt 0) {
            "`nCopy Successful:"; "-" * 30; $copysuccessful; "-" * 30
        }
        if ($copyfailed.length -gt 0) {
            Write-Error "`nCopy Failed:`n$copyfailed"
        }
    }
}


############################################################################################################################

function copyfiletoafolder($sourcefilespath, $destinationfolderpath, $cf_destinationdirectory) {
   
    if ((Test-path $destinationfolderpath) -eq $false) {
        #it will check for destination path if it does not exist.
        if (($cf_destinationdirectory -eq $false)) {
            # and create destination directory is true than new directory will be created.
            Write-Error "Destination path provided does not exist"
            Exit;
        }
        else {
            Create_Destinationpath -destinationpath $destinationfolderpath | Out-Null
        }
    }

    Copy-item $sourcefilespath -Destination $destinationfolderpath | out-null
    $desti_path = Join-path $destinationfolderpath (split-path $sourcefilespath -Leaf)
    if (Test-path $desti_path) {
        if ($sourcepath -eq "Network") {
            $sourcefilespath = Join-Path $global:ogdrivedetails (split-path $sourcefilespath -NoQualifier) 
        }
        "`nCopy Successful:"; "-" * 30; $sourcefilespath; "-" * 30
    }
    else {
        if ($sourcepath -eq "Network") {
            $sourcefilespath = Join-Path $global:ogdrivedetails (split-path $sourcefilespath -NoQualifier) 
        }
        Write-Error "Copy operation failed."
    }
}

###############################################################################################################################

###############################################################################################################################

function movefiletoafolder($sourcefilespath, $destinationfolderpath, $cf_destinationdirectory) {
   
    if ((Test-path $destinationfolderpath) -eq $false) {
        #it will check for destination path if it does not exist.
        if (($cf_destinationdirectory -eq $false)) {
            # and create destination directory is true than new directory will be created.
            Write-Error "Destination path provided does not exist."
            Exit;
        }
        else {
            Create_Destinationpath -destinationpath $destinationfolderpath | Out-Null
        }
    }

    move-item $sourcefilespath -Destination $destinationfolderpath -Force | out-null
    #$desti_path = Join-path $destinationfolderpath (split-path $sourcefilespath -Leaf)
    if ($?) {
        if ($sourcepath -eq "Network") {
            $sourcefilespath = Join-Path $global:ogdrivedetails (split-path $sourcefilespath -NoQualifier) 
        }
        "`nMove Successful:"; "-" * 30; $sourcefilespath; "-" * 30
    }
    else {
        if ($sourcepath -eq "Network") {
            $sourcefilespath = Join-Path $global:ogdrivedetails (split-path $sourcefilespath -NoQualifier) 
        }
        Write-Error "Move operation Failed"
    }
}    
    
############################################################################################################################### 
   
Function Rename_file($renamepath, $newname) {
    if (($newname -notlike "*.*") -and ($fileAction -eq "renameafile")) {
        Write-Error "Please provide file extension with filename."
        Exit;
    }
     
    Rename-Item -Path $renamepath -NewName $newname
    if ($?) {

        if ($sourcepath -eq "Network") {
            $renamepath = Join-Path $global:ogdrivedetails (split-path $renamepath -NoQualifier)
        }

        "`nRename Successful:"; "-" * 30; $renamepath; "-" * 30
    }         
}
###############################################################################################################################
    
function Delete_files($sourcefilepath, $cont_error) {

    $filedeleted = @()
    $failedtodelete = @()
    $display_op = @()
    $display_opfail = @()

    foreach ($file in $sourcefilepath) {

        Get-item $file -Force | Remove-Item -Force -ErrorAction Stop
        
        if ((Test-path $file) -eq $false) {
            $filedeleted += $file          
        }    
        else {
            $failedtodelete += $file        
        }
    }

    if ($sourcepath -eq "Network") {
        foreach ($file in $filedeleted) {     
            $display_op += Join-Path $global:ogdrivedetails (split-path $file -NoQualifier)
        }
        foreach ($file in $failedtodelete) {     
            $display_opfail += Join-Path $global:ogdrivedetails (split-path $file -NoQualifier)
        }
        if ($display_op.length -gt 0) {
            "`nDeletion Successful:"; "-" * 30; $display_op; "-" * 30
        }
        if ($display_opfail.length -gt 0) {
            Write-Error "`nDeletion Failed: $display_opfail"
        }

    }
    else {

        if ($filedeleted.length -gt 0) {
            "`nDeletion Successful:"; "-" * 30; $filedeleted; "-" * 30
        }
        if ($failedtodelete.length -gt 0) {
            Write-Error "`nDeletion Failed: $failedtodelete"
        }
    }
}
  
function Delete_file($sourcefilepath) {

    Get-item $sourcefilepath -Force | Remove-Item -Force -ErrorAction Stop

    if ($sourcepath -eq "Network") {
        $sourcefilepath = Join-Path $global:ogdrivedetails (split-path $sourcefilepath -NoQualifier)
    }
        
    if ((Test-path $sourcefilepath) -eq $false) {
        "-" * 30 + "`nFile deleted successfully : $sourcefilepath`n" + "-" * 30
            
    } 
    
    else {
        Write-Error "Error while deleting file $sourcefilepath."
        
    }
}   
function Delete_folder($sourcefolderpath) {
        
    try {
        Get-item $sourcefolderpath -Force | Remove-Item -force -Recurse -ErrorAction Stop

        
            
        if ((Test-path $sourcefolderpath) -eq $false) {
            if ($sourcepath -eq "Network") {
                $sourcefolderpath = Join-Path $global:ogdrivedetails (split-path $sourcefolderpath -NoQualifier)
            }
            "`nFolder Deleted:"; "-" * 30; $sourcefolderpath; "-" * 30
        } 
    }
    catch {
        $ErrorActionPreference = "Continue"
        Write-Error "Failed to delete folder."
    }
}
        
function Copy_foldertofolder($sourcefolderpath, $destinationfolderpath, $createdest) {
    $foldername = split-path $sourcefolderpath -leaf
    $temppath = "$destinationfolderpath" + "\" + "$foldername"
    $Resulttestpath = Test-Path $temppath
    if(($Resulttestpath) -and !($Overwrite)){
        $ErrorActionPreference = "Continue"
        Write-Error "Destination folder already exist."
        Exit;
    }
    if ((Test-path $destinationfolderpath) -eq $false) {
        #it will check for destination path if it does not exist.
        if (($createdest -eq $false)) {
            # and create destination directory is true than new directory will be created.
            Write-Error "Destination path provided does not exist."
            Exit;
        }
        else {
            Create_Destinationpath -destinationpath $destinationfolderpath | Out-Null
        }
    }
    if ($includesub -eq $true) {
        if (!($includeSystemFile) -or !($includeReadOnlyFile) -or !($includeHiddenFile)) {
            try {
                $ErrorActionPreference = "Stop"
                $SourceFolderName = Split-Path $sourcefolderpath -Leaf
                $NewFolderPath = "$destinationfile" + "\" + "$SourceFolderName"
                if(!($Resulttestpath)){New-Item $NewFolderPath -ItemType Directory | Out-Null}
                Get-ChildItem $sourcefolderpath | ? { $_.PSISContainer -eq $true } | Select -ExpandProperty FullName | Copy-Item  -Destination $NewFolderPath -Recurse | Out-Null 
                $FilesToCheck = Get-ChildItem $sourcefolderpath | where { ! $_.PSIsContainer } | Select -ExpandProperty FullName
                foreach ($File in $FilesToCheck) {
                    $AttCheckResult = verify_Attributes -sourcepath $file -readonly $includeReadOnlyFile -hiddens $includeHiddenFile -system $includeSystemFile
                    if ($AttCheckResult) {
                        Copy-Item $file $NewFolderPath | Out-Null
                    }
                }
            }
            catch {
                $ErrorActionPreference = "SilentlyContinue"
                if(!($Resulttestpath)){Remove-Item $NewFolderPath -Recurse | Out-Null}
                $_.Exception.Message
                Write-Error "Operation Failed"; 
                Exit
            }
        }
        Else {
            Copy-item $sourcefolderpath -Destination $destinationfolderpath -Recurse -Force | Out-Null
        }
    }
    else {
        try {
            Copy-item $sourcefolderpath -Destination $destinationfolderpath -ErrorAction Stop | Out-Null
            $FilesToCheck = Get-ChildItem $sourcefolderpath | where { ! $_.PSIsContainer } | Select -ExpandProperty FullName
            Function tocopyfiles($FileList){
                if($continueonerror -eq $True){$ErrorActionPreference = "SilentlyContinue"}Else{$ErrorActionPreference = "Stop"}
                foreach ($File in $FileList) {
                    $AttCheckResult = verify_Attributes -sourcepath $file -readonly $includeReadOnlyFile -hiddens $includeHiddenFile -system $includeSystemFile
                    if ($AttCheckResult -eq $true) {
                        Copy-Item $file $temppath | Out-Null
                    }
                }
            }  
            tocopyfiles -FileList $FilesToCheck
        }
        catch {
            Remove-Item $temppath -Recurse -Force | Out-Null
            $ErrorActionPreference = "Continue"
            Write-Error "Error while copying files."
        }
    }
    if ($sourcepath -eq "Network") {
        $sourcefolderpath = Join-Path $global:ogdrivedetails (split-path $sourcefolderpath -NoQualifier)
    }
    
    if (Test-path $temppath) {
        "`nCopy Successful:"; "-" * 30; $sourcefolderpath; "-" * 30
    } 
    else {
        Write-Error "Copy operation failed"
    } 
}
    
####### Ignoring Attribute criteria and moving folder############################################

#On Error Script will Terminate#
   
function Move_foldertofolder($sourcefolderpath, $destinationfolderpath, $createdest) {
    try{
    $ErrorActionPreference = "Stop"
    $foldername = split-path $sourcefolderpath -leaf
    $temppath = "$destinationfolderpath" + "\" + "$foldername"

    if ((Test-path $destinationfolderpath) -eq $false) {
        #it will check for destination path if it does not exist.
        if (($createdest -eq $false)) {
            # and create destination directory is true than new directory will be created.
            Write-Error "Destination path provided does not exist."
            Exit;
        }
        else {
            Create_Destinationpath -destinationpath $destinationfolderpath | Out-Null
        }
    }
    if((Test-Path $temppath) -and ($Overwrite -eq $true)){
        try{
            Copy-Item $sourcefolderpath $destinationfolderpath -Recurse -Force -ErrorAction Stop | Out-Null
            Remove-Item $sourcefolderpath -Recurse -Force | Out-Null
        }
        catch{
            $ErrorActionPreference = "Continue"
            Write-Error "Move operation failed"
        }
    }
    elseif((Test-Path $temppath) -and !($Overwrite)){Write-Error "Folder Already exist."}
    else{try{Move-item $sourcefolderpath -Destination $destinationfolderpath -ErrorAction Stop}catch{$ErrorActionPreference = "Continue";Write-Error "Move operation failed"}}
    
    if ($sourcepath -eq "Network") {
        $sourcefolderpath = Join-Path $global:ogdrivedetails (split-path $sourcefolderpath -NoQualifier)
    }

    if (Test-path $temppath) {
        "`nMove Successful:"; "-" * 30; $sourcefolderpath; "-" * 30
    } 
    else {
        Write-Error "Move operation failed."
    } 
    }
    catch{
        $ErrorActionPreference = "Continue"
        $_.Exception.Message
        Write-Error "Error While moving folder"
    }
}
Function Create_Folder {
    
    if ($SourcePath -eq "Network") {
        try {
            $ErrorActionPreference = "Stop"
            $sourcefile = Map_networkdrive -Action "Create"
            $ErrorActionPreference = "Continue"
        }
        catch {
            $ErrorActionPreference = "Continue"
            Write-Error "Failed to map network drive"
            Exit;
        }
    }
    $FullPath = $SourceFile
    try {
        if (Test-path (Split-Path $SourceFile -Parent)) {
            New-Item -Path $FullPath -ItemType "Directory" -ErrorAction "Stop" | Out-Null
            if (Test-path $FullPath) { "Folder Created Successfully" }else { Write-Error "Failed to create folder"; Exit; }
        }
        else {
            Write-Error "Source path does not exist."
            Exit;
        }
    }
    catch {
        Write-Error "Failed to create folder"
        Exit;
    }
    
}

Function Select_Operation {

    switch ($fileaction) {

        "copyfilestoafolder" { copyfilestoafolder -sourcefilespath $sourcefile -destinationfolderpath $destinationfile -cf_continueonerror $continueonerror -cf_destinationdirectory $CreateDestDirectory }
        "copyafile" { copyfiletoafolder -sourcefilespath $sourcefile -destinationfolderpath $destinationfile -cf_destinationdirectory $CreateDestDirectory }
        "copyafolder" { Copy_foldertofolder -sourcefolderpath $sourcefile -destinationfolderpath $destinationfile -createdest $CreateDestDirectory }
        "renameafile" { Rename_file -renamepath $sourcefile -newname $newname }
        "moveafile" { movefiletoafolder -sourcefilespath $sourcefile -destinationfolderpath $destinationfile -cf_destinationdirectory $CreateDestDirectory }
        "movefolder" { Move_foldertofolder -sourcefolderpath $sourcefile -destinationfolderpath $destinationfile -createdest $CreateDestDirectory }
        "renameafolder" { Rename_file -renamepath $sourcefile -newname $newname }
        "deleteafile" { Delete_file -sourcefilepath $sourcefile }
        "deletefiles" { Delete_files -sourcefilepath $sourcefile -cont_error $continueonerror }
        "deleteafolder" { Delete_folder -sourcefolderpath $sourcefile }
        "CreateFolder"{Create_Folder}
    }

}
         
try {
    Select_Operation
}
catch { }
Finally { Map_networkdrive -Action "Remove" -Drive $global:newdrive }




