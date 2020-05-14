# $Warnbeforepermanentlydeletingitems = 'Yes'
# $EmptytheDeletedItemsfolderuponexit = 'Yes'
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

Function Get-RDPSessions {
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

Function AppendToCSV ($Type, $Message) {
    New-Object -TypeName psobject -Property @{Type = $Type; Message = $Message } | ConvertTo-Csv -NoTypeInformation | Select-Object -Skip 1 | Out-File C:\temp\log.csv -Append
}

$Scriptblock = {
    $ErrorActionPreference = 'Stop'
    Function AppendToCSV ($Type, $Message) {
        New-Object -TypeName psobject -Property @{Type = $Type; Message = $Message } | ConvertTo-Csv -NoTypeInformation | Select-Object -Skip 1 | Out-File C:\temp\log.csv -Append
    }

    try {
        $IsOfficeInstalled = Test-Path "Registry::HKEY_CURRENT_USER\Software\Microsoft\Office"
        $Version = [System.Version](New-Object -ComObject outlook.application | Select-Object -ExpandProperty version)

        if ($IsOfficeInstalled -and $Version) {
        
            if (!($Version.Major -eq 14 -or $Version.Major -eq 15 -or $Version.Major -eq 16)) {
                AppendToCSV "Output" "`nOutlook version is not supported. Only supported versions are Outlook 2010 / 2013 / 2016."
                exit;
            }

            $Office = "HKEY_CURRENT_USER\Software\Microsoft\Office"
            $MailSettings = "$Office\$($Version.Major).0\Common\MailSettings"
            $ProofingTools = "HKEY_CURRENT_USER\Software\Microsoft\Shared Tools\Proofing Tools\1.0\Office"
            $Outlook = "$Office\$($Version.Major).0\Outlook"
            $Preferences = "$Outlook\Preferences"
            $Options = "$Outlook\Options"
            $Success = $Failure = @()

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
"ComposeinthisMessageFormat","Compose in this Message Format","$Options\Mail","EditorPreference"
"SendacopyofthepicturesinsteadofthereferencetotheirlocationonlyforHTMLformat","Send a copy of the pictures instead of the reference to their location (only for HTML format)","$Options\Mail","Send Pictures With Document"
"SavecopiesinSentitemsfolder","Save copies in Sent items folder","$Preferences","SaveSent"
"Autosaveunsent","Autosave unsent","$MailSettings","AutosaveTime"
"Alwayscheckspellingbeforesending","Always check spelling before sending","$Options\Spelling","Check"
"IgnorewordsinUPPERCASE","Ignore words in UPPERCASE","$ProofingTools","IgnoreUpperCase"
"Ignorewordswithnumbers","Ignore words with numbers","$ProofingTools","IgnoreWordsWithNumbers"
"Ignoreoriginalmessageinreplies","Ignore original message in replies","$MailSettings","IgnoreReplySpelling"
"@
            #"Alwayssuggestreplacementsformisspelledwords","Always suggest replacements for misspelled words","$Options\Spelling","SpellAlwaysSuggest"
            #"UseMicrosoftWordtoeditemailmessages","Use Microsoft Word to edit email messages","$Options\Mail","UseWordMail"

            $UserInputs = $csv | ConvertFrom-Csv

            Foreach ($item in $UserInputs) {
                $Registry = $item.Registry
                $Key = $item.Key
                $UserEnteredValue = (Get-Variable $item.InputVariable).value
                $Value = switch ($UserEnteredValue) {
                    'Yes' { 1 }
                    'No' { 0 }
                    'HTML' { [Convert]::ToString('0x20000', 10) }
                    'Plain Text' { [Convert]::ToString('0x10000', 10) }
                    'Rich Text' { [Convert]::ToString('0x30000', 10) }
                    'No Change' { 'NoAction' }
                }

                if ($Value -ne 'NoAction') {
                    if (!(Test-Path "REGISTRY::$Registry")) {
                        New-Item -Path "REGISTRY::$Registry" -Name $Key -ErrorAction SilentlyContinue -Force | Out-Null
                        #New-ItemProperty -Path "REGISTRY::$Registry" -Name $Key -ErrorAction SilentlyContinue -Force | Out-Null
                    }

                    Set-ItemProperty "REGISTRY::$Registry" -Name $Key -Value $Value -ErrorAction SilentlyContinue -Force
                    
                    ### START Special conditions
                    if ($item.'Setting' -eq "Run AutoArchive" -and $Days) {
                        'here'
                        Set-ItemProperty "REGISTRY::$Preferences" -Name 'EveryDays' -Value $Days -ErrorAction SilentlyContinue -Force
                    }
                    ### END Special conditions

                    If ((Get-ItemProperty "REGISTRY::$Registry" -Name $Key -ErrorAction SilentlyContinue )."$Key" -eq $Value) {
                        $Success += New-Object PSOBject -Property @{ 'Setting' = $item.Setting; 'Value' = $UserEnteredValue }
                    }
                    else {
                        $Failure += New-Object PSOBject -Property @{ 'Setting' = $item.Setting; 'Value' = $UserEnteredValue }
                    }
                }
            }

            $ResultCount = $Success.count + $Failure.count
            if ($ResultCount -ge 1) {
                $i = 1
                if ($Success) {
                    AppendToCSV "Output" $("`n[{0}/{1}] Outlook configurations successful.`n" -f $Success.count, $ResultCount)
                    $Padding = ($Success | ForEach-Object { $_.setting.length } | Measure-Object -Maximum).Maximum + 3
                    $Success | ForEach-Object {
                        AppendToCSV "Output" $($("{0}. {1}" -f $i, $_.Setting).PadRight($Padding, " ") + " = `'$($_.Value)`'")
                        $i++
                    }
                }
        
                $j = 1
                if ($Failure) {
                    AppendToCSV "Output" $("`n[{0}/{1}] Outlook configurations failed.`n" -f $Failure.count, $ResultCount)
                    $Padding = ($Failure | ForEach-Object { $_.setting.length } | Measure-Object -Maximum).Maximum + 3
                    $Failure | ForEach-Object {
                        AppendToCSV "Output" $($("{0}. {1}" -f $j, $_.Setting).PadRight($Padding, " ") + " = `'$($_.Value)`'")
                        $j++
                    }
                }
            }
            else {
                AppendToCSV "Output" "`nSuccessful. No Configuration was selected in the user input. No Action Performed."
            }    
        }
        else {
            AppendToCSV "Error" "Outlook is not installed on System."
        }    
    }
    catch {
        AppendToCSV "Error" $($_.exception.message)
    }
}

$file = @"
`$Warnbeforepermanentlydeletingitems = "$Warnbeforepermanentlydeletingitems"
`$EmptytheDeletedItemsfolderuponexit = "$EmptytheDeletedItemsfolderuponexit"
`$DisplayaNewmailDesktopAlert = "$DisplayaNewmailDesktopAlert"
`$Playasound = "$Playasound"
`$RunAutoArchive = "$RunAutoArchive"
`$Days = "$days"
`$PrompttoAutoArchive = "$PrompttoAutoArchive"
`$DeleteExpiredItemsemailfoldersonly = "$DeleteExpiredItemsemailfoldersonly"
`$Allowcommaasaddressseparator = "$Allowcommaasaddressseparator"
`$Automaticnamechecking = "$Automaticnamechecking"
`$ComposeinthisMessageFormat = "$ComposeinthisMessageFormat"
`$SendacopyofthepicturesinsteadofthereferencetotheirlocationonlyforHTMLformat = "$SendacopyofthepicturesinsteadofthereferencetotheirlocationonlyforHTMLformat"
`$SavecopiesinSentitemsfolder = "$SavecopiesinSentitemsfolder"
`$Autosaveunsent = "$Autosaveunsent"
`$Alwayscheckspellingbeforesending = "$Alwayscheckspellingbeforesending"
`$IgnorewordsinUPPERCASE = "$IgnorewordsinUPPERCASE"
`$Ignorewordswithnumbers = "$Ignorewordswithnumbers"
`$Ignoreoriginalmessageinreplies = "$Ignoreoriginalmessageinreplies"
"@ + $Scriptblock.ToString()

$ErrorActionPreference = 'Stop'
$TaskName = "OutlookConfig"
$LogDir = 'C:\temp\'
$LogFilePath = 'C:\temp\Log.csv'
$MainFilePath = "C:\temp\$TaskName.ps1"

if (![IO.Directory]::Exists($LogDir)) { [IO.Directory]::CreateDirectory($LogDir) | Out-Null } # create folder if that doesn't exists
$File | Out-File $MainFilePath -Encoding utf8 # copy main script which will be later executed from scheduled task

try {
    if (($Session = Get-RDPSessions)) {
        '"Message","Type"' | Out-File $LogFilePath -Append
        $User = $Session | Select-Object -First 1 -ExpandProperty Username # select one user profile from active RDP sessions
        
        Start-Sleep -Seconds 5
        $Task = "PowerShell.exe -executionpolicy bypass -NoExit -noprofile -WindowStyle Hidden -command '. $MainFilePath '"
        $StartTime = (Get-Date).AddMinutes(2).ToString('HH:mm') # time in 24hr format
        
        schtasks.exe /create /s $($env:COMPUTERNAME) /tn $TaskName /sc once /tr $Task /st $StartTime /ru $User /F | Out-Null
        Start-Sleep -Seconds 5

        schtasks.exe /End /TN $TaskName | Out-Null
        schtasks.exe /Run /TN $TaskName | Out-Null
        Start-Sleep -Seconds 200

        # Sending logs to PowerShell Output\Error streams so that agent can capture it
        if (Test-Path $LogFilePath) {
            $Logs = Import-Csv $LogFilePath
            if ($Logs | Select-Object -ExpandProperty Message -ErrorAction SilentlyContinue | Foreach-Object { $_.trim() }) {
                foreach ($item in $Logs) {
                    switch ($item.type) {
                        'Output' { Write-Output "$($item.message)" }
                        'Error' { Write-Error "$($item.message)" }
                    }
                }
            }
            else {
                Write-Output "`nNo action performed, because the task exited due to timeout. This may be due to slowness of the System."
            }
        }
    }
    else {
        Write-Output "This script requires logon user and currently no user is logged in. `nNo action will be performed."; exit;
    }
}
catch {
    Write-Error $_
}
finally {
    try {
        #schtasks.exe /Run /TN $TaskName | Out-Null
        schtasks.exe /Delete /TN $TaskName /F | Out-Null # scheduled task cleanup
        Remove-Item -Path $MainFilePath, $LogFilePath -ErrorAction SilentlyContinue # file cleanup
    }
    catch {
        # no action
    }
}
