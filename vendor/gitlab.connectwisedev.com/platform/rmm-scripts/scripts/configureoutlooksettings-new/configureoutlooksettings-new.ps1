#$Warnbeforepermanentlydeletingitems = 'Yes'
#$EmptytheDeletedItemsfolderuponexit = 'Yes'
# $DisplayaNewmailDesktopAlert = 'Yes'
# $Playasound = 'Yes'
# $RunAutoArchive = 'Yes' # with variables $days
# $Days = 55 # integer
# $PrompttoAutoArchive = 'Yes'
# $DeleteExpiredItemsemailfoldersonly = 'Yes'
# $Allowcommaasaddressseparator = 'Yes'
# $Automaticnamechecking = 'Yes'
# $ComposeinthisMessageFormat = 'HTML' # 'HTML', 'Rich Text', 'Plain Text'
# $SendacopyofthepicturesinsteadofthereferencetotheirlocationonlyforHTMLformat = 'Yes'
# $SavecopiesinSentitemsfolder = 'Yes'
# $Autosaveunsent = 'Yes' 
# $Alwayscheckspellingbeforesending = 'Yes'
# $IgnorewordsinUPPERCASE = 'Yes'
# $Ignorewordswithnumbers = 'Yes'
# $Ignoreoriginalmessageinreplies = 'Yes'
# $Moveolditemsto = "C:\Users\username\Archive.pst"

#####################
 
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

Function Get-CurrentUser {
    try {
        $ErrorActionPreference = 'Stop'
        if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64" -and $env:PROCESSOR_ARCHITECTURE -eq 'x86') { $Query = 'c:\windows\sysnative\query.exe' }else { $Query = 'c:\windows\System32\query.exe' }
        $LoggedOnUsers = if (($Users = (& $Query user 2>&1))) {
            $Users | ForEach-Object { (($_.trim() -replace ">" -replace "(?m)^([A-Za-z0-9]{3,})\s+(\d{1,2}\s+\w+)", '$1  none  $2' -replace "\s{2,}", "," -replace "none", $null)) } |
            ConvertFrom-Csv |
            Where-Object { $_.state -eq 'Active' } #-and $_.SESSIONNAME -like "rdp*" }
        }
        return $LoggedOnUsers
    }
    catch {
        # no action
    }
}

$uname = (Get-CurrentUser).Username

if ([string]::IsNullOrEmpty("$uname")) { 

    Write-Error "This script requires logon user and currently no user is logged in. `nNo action will be performed."
    exit;
}
        
$USID = Get-WmiObject -Class Win32_UserAccount | Where-Object { $_.Name -eq "$uname" } | Select-Object -ExpandProperty SID 

$ErrorActionPreference = 'Stop'

try {
    $IsOfficeInstalled = Test-Path "registry::HKEY_USERS\$USID\Software\Microsoft\Office"
    $Version = [System.Version](New-Object -ComObject outlook.application | Select-Object -ExpandProperty version)

    if ($IsOfficeInstalled -and $Version) {
           
        
        if (!($Version.Major -eq 14 -or $Version.Major -eq 15 -or $Version.Major -eq 16 -or $Version.Major -eq 12)) {
            Write-Error "`nOutlook version is not supported. Only supported versions are Outlook 2007 / 2010 / 2013 / 2016."
            exit;
        }
       
        $Office = "HKEY_USERS\$USID\Software\Microsoft\Office"
        $MailSettings = "$Office\$($Version.Major).0\Common\MailSettings"
        $ProofingTools = "HKEY_USERS\$USID\Software\Microsoft\Shared Tools\Proofing Tools\1.0\Office"
        $Outlook = "$Office\$($Version.Major).0\Outlook"
        $Profile_Name = (Get-ChildItem "Registry::$Outlook\Profiles").PSChildName | Select-Object -First 1
        
        if ([string]::IsNullOrEmpty("$Profile_Name")) { 

            Write-Error "`nOutlook not configured for current user."
            exit;
        }
        $Preferences = "$Outlook\Preferences"
        $Options = "$Outlook\Options"
        $StartReg = "$Outlook\Profiles\$Profile_Name\0a0d020000000000c000000000000046"
       

        $CSV = @"
"InputVariable","Setting","Registry","Key"
"Warnbeforepermanentlydeletingitems","Warn before permanently deleting items","$Options\General","WarnDelete"
"EmptytheDeletedItemsfolderuponexit","Empty the Deleted Items folder upon exit","$Preferences","EmptyTrash"
"DisplayaNewmailDesktopAlert","Display a New mail Desktop Alert","$Preferences","NewmailDesktopAlerts"
"Playasound","Play a sound","$Preferences","PlaySound"
"RunAutoArchive","Run AutoArchive","$Preferences","DoAging"
"PrompttoAutoArchive","Prompt to AutoArchive","$Preferences","PromptForAging"
"DeleteExpiredItemsemailfoldersonly","Delete Expired Items (e-mail folders only)","$Preferences","DeleteExpired"
"Allowcommaasaddressseparator","Allow comma as address separator","$Preferences","AllowCommasInRecip"
"Automaticnamechecking","Automatic name checking","$Preferences","AutoNameCheck"
"SendacopyofthepicturesinsteadofthereferencetotheirlocationonlyforHTMLformat","Send a copy of the pictures instead of the reference to their location (only for HTML format)","$Options\Mail","Send Pictures With Document"
"SavecopiesinSentitemsfolder","Save copies in Sent items folder","$Preferences","SaveSent"
"Autosaveunsent","Autosave unsent","$MailSettings","AutosaveTime"
"Alwayscheckspellingbeforesending","Always check spelling before sending","$Options\Spelling","Check"
"IgnorewordsinUPPERCASE","Ignore words in UPPERCASE","$ProofingTools","IgnoreUpperCase"
"Ignorewordswithnumbers","Ignore words with numbers","$ProofingTools","IgnoreWordsWithNumbers"
"Ignoreoriginalmessageinreplies","Ignore original message in replies","$MailSettings","IgnoreReplySpelling"
"Moveolditemsto","Move old items to","$StartReg","001f0324"
"@
        #"Alwayssuggestreplacementsformisspelledwords","Always suggest replacements for misspelled words","$Options\Spelling","SpellAlwaysSuggest"
        #"UseMicrosoftWordtoeditemailmessages","Use Microsoft Word to edit email messages","$Options\Mail","UseWordMail"
        #"ComposeinthisMessageFormat","Compose in this Message Format","$Options\Mail","EditorPreference"

        $UserInputs = $csv | ConvertFrom-Csv

        Foreach ($item in $UserInputs) {
            $Registry = $item.Registry
            $Key = $item.Key
            $UserEnteredValue = (Get-Variable $item.InputVariable).value
            $Value = switch ($UserEnteredValue) {
                'Yes' { 1; break }
                'No' { 0; break }
                'HTML' { [Convert]::ToString('0x20000', 10); break }
                'Plain Text' { [Convert]::ToString('0x10000', 10); break }
                'Rich Text' { [Convert]::ToString('0x30000', 10); break }
                'No Change' { 'NoAction'; break }
                "$UserEnteredValue" { $UserEnteredValue; break }
            }

            if ($Value -ne 'NoAction') {
                if (!(Test-Path "REGISTRY::$Registry")) {
                       
                    New-Item -Path "REGISTRY::$Registry" -Name $Key -ErrorAction SilentlyContinue -Force | Out-Null
                       
                }

                if ($item.'Setting' -eq "Move old items to" -and "$Moveolditemsto") {
                 
                    #Check if user has provided the path with correct file extension
                    $Parent_Path = Split-Path "$Moveolditemsto" -Parent -ErrorAction 'SilentlyContinue'
                    $Leaf = Split-Path "$Moveolditemsto" -Leaf -ErrorAction 'SilentlyContinue'
                    $Default_Path = (Get-ItemProperty "Registry::HKEY_USERS\$USID\Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders\").Personal  | Select-Object -First 1

                    if (!($Parent_Path) -and $Leaf -and ($Leaf -like "*.pst")) {
                        $value = Join-path $Default_Path $Leaf
                    }
                    elseif ($Parent_Path -and ($leaf -like "*.pst")) {
                        $value = $Moveolditemsto
                    }
                    else {
                        Write-Error "Move old items provided file name is not a PST file. Check file extension"
                    }
                    
                    if (Test-Path (Split-Path "$value" -Parent)) {
                        Set-ItemProperty "REGISTRY::$Preferences" -Name 'ArchiveDelete' -Value 0 -ErrorAction SilentlyContinue -Force
                        Set-ItemProperty "REGISTRY::$Preferences" -Name 'ArchiveOld' -Value 1 -ErrorAction SilentlyContinue -Force
                     
                        Set-ItemProperty -Path "REGISTRY::$Registry" -Name $Key -Value $(([system.Text.Encoding]::Unicode).GetBytes($value)) -Type Binary -ErrorAction SilentlyContinue  -Force
                    }
                    else { Write-Error "Move old items path '$Moveolditemsto' does not exist." }
                }
                else {
                    Set-ItemProperty "REGISTRY::$Registry" -Name $Key -Value $Value -ErrorAction SilentlyContinue -Force
                }
                
                ### START Special conditions
                if ($item.'Setting' -eq "Run AutoArchive" -and $Days) {
                        
                    Set-ItemProperty "REGISTRY::$Preferences" -Name 'EveryDays' -Value $Days -ErrorAction SilentlyContinue -Force
                    if ((Get-ItemProperty "REGISTRY::$Preferences" -Name 'EveryDays' -ErrorAction SilentlyContinue ).EveryDays -eq $Days) {
                        Write-Output "Run auto archive days : $Days"
                    }
                    else {
                        Write-Error "Run auto archive days : $Days configuration failed"
                    }
                }
                ### END Special conditions
                    
                    
                if (($item.'Setting' -eq "Move old items to") -and (Test-Path (Split-Path "$value" -Parent))) {
                    
                    [string]$Value1 = $(([system.Text.Encoding]::Unicode).GetBytes($Value)) 
    
                    If ([string](Get-ItemProperty "REGISTRY::$Registry" -Name $Key -ErrorAction SilentlyContinue )."$Key" -eq $value1) {
                       
                        Write-Output "Move old items to : $value"
                    }
                    else {
                        Write-Error "Move old items to configuration failed"
                    }
                } 
                else {
                        
                    $ReValue = switch ($value) {
                        '1' { 'Yes'; break }
                        '0' { 'No'; break }
                        '131072' { 'HTML'; break }
                        '65536' { 'Plain Text'; break }
                        '196608' { 'Rich Text'; break }
                        'NoAction' { 'No Change'; break }
                        "$value" { $value; break }
                    }
                    If ((Get-ItemProperty "REGISTRY::$Registry" -Name $Key -ErrorAction SilentlyContinue )."$Key" -eq $Value) {
                        
                        Write-Output "$($item.Setting) : $ReValue"
                    }
                    else {
                        
                        Write-Error "$($item.Setting) configuration failed"
                    }       
                }   
            }
        }   
    }
    else {
        Write-Error "`nOutlook is not installed on System."
    }    
}
catch {
    $_.exception.message
}
