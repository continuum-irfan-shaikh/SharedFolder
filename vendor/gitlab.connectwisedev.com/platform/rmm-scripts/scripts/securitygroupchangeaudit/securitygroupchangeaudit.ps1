<#
.SYNOPSIS
 Script to fetch more information for monitoring rule 'Security Notification - A user account added to Security Group'.

.DESCRIPTION
 This script fetch more information for monitoring rule 'Security Notification - A user account added to Security Group'.

.NOTES    
 Name: Get-SGMembershipInformation.ps1
 Author: Imran Khan  
 Version: 1.2
 DateCreated: 2019-09-20 
 DateUpdated: 2019-10-01 
.PARAMETER 
 
.EXAMPLE
 .\Get-SGMembershipInformation.ps1 
#>

############################################################################
# Function to execute non powershell commands
############################################################################
Function ExecuteCMDCommands($Command) {
    try {
        if ($Command) {
            $finalCmd = "cmd /c " + $Command + " '2>&1'"
            $msg = (Invoke-Expression -Command $finalCmd -ErrorAction Stop)
            #$cmdOutput = $finalCmd | % { $_.ToString() };
            Write-Output $msg
       
        }
        else { Write-Output "Command can not be null" }
    }
    catch
    { Write-Output "Error:" $msg }

}#End function ExecuteCMDCommands

Function GetNameFromSID($Sid) {
    #Default SID Mappings
    $SIDMappings = @{
        "S-1-1-0"      = "Everyone"
        "S-1-2"        = "Local Authority"
        "S-1-5"        = "NT Authority"
        "S-1-5-11"     = "Authenticated Users"
        "S-1-5-18"     = "Local System"
        "S-1-5-19"     = "NT Authority (Local Service)"
        "S-1-5-20"     = "NT Authority (Network Service)"
        "S-1-5-32-551" = "Backup Operators"
            
    }
      
    $UserAccount = Get-WmiObject Win32_useraccount -Filter "sid='$($sid)'"
    if ($UserAccount) {
        $return = ("{0}\{1}" -f $UserAccount.domain, $UserAccount.Name)
    }
    else {
        $Group = Get-WmiObject Win32_Group -Filter "sid='$($sid)'"
        if ($Group) {
            $return = ("{0}\{1}" -f $Group.domain, $Group.Name)
        }
        elseif ($SIDMappings["$sid"]) {
            $return = $SIDMappings["$sid"]
        }
        else {
            $return = $sid
        }
    }
    return $return
} #Function End

############################################################################
# Function to get event information
############################################################################
Function GetEventInfo($LogName, $EventID) {
    try {
        
        #Get Latest Event log for system
        $isEventPresent = Get-EventLog -List | Where-Object { $_.Log -eq "$LogName" }
        
        if ($isEventPresent) {
            #check event ID latest trigger time and number of occurances withing 24 hours
            $event = Get-EventLog -logname "$LogName" -After (get-date).AddHours(-24) -ErrorAction Stop | Where-Object { $_.EventID -eq $EventID }
            #$EventInfo = $event | Sort-Object -Descending TimeGenerated | Select-Object -First 1 -Property *
            if($null -ne $event){
            Write-Output $("-" * 40)
            Write-Output $("-" * 40)
                $eventCount = ($event | Measure-Object).count
                if ($eventCount) {
                    Write-Output "Event $EventID has been occurred $($eventCount) times in last 24 hours"
                    Write-Output "Details for all event ID $EventID generated in 24 hours:"
                }
                else {
                    Write-Output "Event $EventID has not been occurred in last 24 hours"
                }
            foreach($eventinfo in $event){
            if ($eventInfo) {
            
                Write-Output $("-" * 40)
                Write-Output "Event $EventID  occurred at $($eventInfo.TimeGenerated.ToString())."
                
                Write-Output $("-" * 40)
                
               
                #Member
                $MemberSecurityID = GetNameFromSID $EventInfo.ReplacementStrings[1]        
                #Group
                $GroupName = $EventInfo.ReplacementStrings[2]
                $GroupDomain = $EventInfo.ReplacementStrings[3]
                $GroupSecurityID = GetNameFromSID $EventInfo.ReplacementStrings[4]    
                #Subject
                $SubjectSecurityID = GetNameFromSID $EventInfo.ReplacementStrings[5]
                $SubjectAccountName = $EventInfo.ReplacementStrings[6]
                $SubjectAccountDomain = $EventInfo.ReplacementStrings[7]
                $SubjectLogonID = $EventInfo.ReplacementStrings[8]    
                Write-Output "User details responsible for adding member to the Group:"
                Write-Output "Security ID   : $SubjectSecurityID"
                Write-Output "Account Name  : $SubjectAccountName"
                Write-Output "Account Domain: $SubjectAccountDomain"
                Write-Output "Logon ID      : $SubjectLogonID"
                Write-Output $("-" * 40)
                Write-Output "Member details which has been added to the Group:"
                Write-Output "Security ID   : $MemberSecurityID"
                Write-Output $("-" * 40)
                Write-Output "Group details in which User has been added:"
                Write-Output "Group Name    : $GroupName"
                Write-Output "Group Domain  : $GroupDomain"
                Write-Output "Security ID   : $GroupSecurityID"
                Write-Output $("-" * 40)
                #For Group
                $GroupInfo = Get-WmiObject -Query "select * from win32_group where Name='$($GroupName)'"
                $GroupMembers = Get-WmiObject -Query "select * from win32_groupuser where GroupComponent=`"Win32_Group.Domain='$($GroupInfo.Domain)',Name='$($GroupInfo.Name)'`""           
                $GroupList = @()
                foreach ($GroupMember in $GroupMembers) {
                    if ($GroupMember.PartComponent -match 'Name="(.+)"') {               
                        $GroupList += $Matches[1]
                    } 
                }
                Write-Output "Member List for Group: $GroupName :"
                Write-Output $GroupList
                Write-Output $("-" * 40)
            
              
            }
            }
              Write-Output $("-" * 40)
              Write-Output "Printing more information for all the occurrences of event ID $EventID during the last 24 hours"
                $event | ForEach-Object {
                    Write-Output "Time Generated:$($_.TimeGenerated)"
                    Write-Output "Message       :"$($_.Message)
                    Write-Output $("*" * 50)
                }
            }else {
                Write-Output  "No entry found for event $EventID under $LogName logs for alert condition in last 24 hours"
            }

        }
        else {
            Write-Output "$LogName log file is not present on current server"
        }  
    }
    catch {
        Write-Output "[MSG: ERROR : $($_.Exception.message)]"
    }
 
}#Function GetEventInfo End


#Actual execution will start from here    
try {
    <# Architecture check started and PS changed to the OS compatible #>
    if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    
        if ($myInvocation.Line) {
            &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
        }
        else {
            &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
        }
        exit $lastexitcode
    }
    <#Architecture check completed #>

    <# Compatibility check if found incompatible will exit #>
    try {
        [double]$OSVersion = [Environment]::OSVersion.Version.ToString(2)
        $PSVersion = (Get-Host).Version
        if (($OSVersion -lt 6.1) -or ($PSVersion.Major -lt 2)) {
            Write-Output "[MSG: System is not compatible with the requirement. Either machine is below Windows 7 / Windows 2008R2 or Powershell version is lower than 2.0]"
            Exit
        }
    }
    catch { Write-Output "[MSG: ERROR : $($_.Exception.message)]" ; EXIT } 
    <# Compatibility Check Code Ends #>

    #$OSInfo = Get-WmiObject -Class Win32_OperatingSystem -Property * -ErrorAction Stop
     
    #Retrieve Event Information
    $isEventPresent = Get-EventLog -List | Where-Object { $_.Log -eq "Security" }
        
    if ($isEventPresent) {
        #check event ID latest trigger time and number of occurances withing 24 hours
       
        GetEventInfo -LogName "Security" -EventID "4728"
        GetEventInfo -LogName "Security" -EventID "4732"
        GetEventInfo -LogName "Security" -EventID "4756"
        }
    
    else {
        Write-Output "$LogType log file is not present on current server"
        EXIT
    }
    
    #GPO Verification
    #Get GPO information fro command gpresult
    if (!(Test-Path "C:\temp")) {
        New-Item C:\temp -ItemType Directory -ErrorAction Stop | Out-Null
    }
    
    #Get GPO Information from gpresult /z
    $GPOAllInfo = ExecuteCMDCommands -Command "gpresult /USER `"Administrator`" /SCOPE Computer /z"
    if ($GPOAllInfo -match "does not have RSOP data") {
        $GPOAllInfo = ExecuteCMDCommands -Command "gpresult /USER `"nochelpdesk`" /SCOPE Computer /z"
    }
    if ($GPOAllInfo -match "does not have RSOP data") {
        Write-Output "$GPOAllInfo"
    }
    else {
        $startIndex = ($GPOAllInfo -join "`n").IndexOf("Restricted Groups") + 18
        $IndexCounter = ($GPOAllInfo -join "`n").IndexOf("System Services") - $startIndex
        $requiredData = ($GPOAllInfo -join "`n").subString($startIndex, $IndexCounter)
        
        if ($requiredData -like "*GPO*") {
            $GPOInfo = $requiredData
        }
        Write-Output $("-" * 40)
        if ($GPOInfo) {
            Write-Output "Below configuration has been done on restricted groups using GPO: "$($GPOInfo -replace "--", "")
        }
        else {
            Write-Output "Did not find 'restricted groups' group policy configured"
        }
        Write-Output $("-" * 40)
    }
    Write-Output $("-" * 40)
    $GPOxmlFile = "C:\temp\GPOResult_$((Get-Date).ToString("ddMMyyhhss")).xml"
    $GPOResult = ExecuteCMDCommands -Command "gpresult /USER `"Administrator`" /SCOPE:computer /x $GPOxmlFile"
    if ($GPOResult -match "does not have RSOP data") {
        $GPOResult = ExecuteCMDCommands -Command "gpresult /USER `"nochelpdesk`" /SCOPE:computer /x $GPOxmlFile"
    }
    
    if ($GPOResult -match "does not have RSOP data") {
        Write-Output "$GPOResult"
        EXIT
    }
    if (Test-Path $GPOxmlFile) {
        [xml]$xmlData = Get-Content $GPOxmlFile
        $ExtensionData = $xmlData.rsop.ComputerResults.ExtensionData
        $xmlInfo = @()
        
        $ExtensionData | ForEach-Object {
            if (($_.Extension.ChildNodes | Select-Object BaseInstanceXml).BaseInstanceXml.INSTANCE.Members.Member | Select-Object name -ExpandProperty name) {
                $GPODomain = ($_.Extension.ChildNodes | Select-Object GPO).GPO.Domain.'#text'
                $GPOName = ($_.Extension.ChildNodes | Select-Object BaseInstanceXml).BaseInstanceXml.PROPERTY | Where-Object { $_.Name -eq "polmkrBaseGpoDisplayName" } | Select-Object -ExpandProperty Value
                $GroupName = ($_.Extension.ChildNodes | Select-Object BaseInstanceXml).BaseInstanceXml.INSTANCE.PROPERTY | Where-Object { $_.Name -eq "polmkrGroupName" } | Select-Object -ExpandProperty Value
                $MembersName = ($_.Extension.ChildNodes | Select-Object BaseInstanceXml).BaseInstanceXml.INSTANCE.Members.Member | Select-Object name -ExpandProperty name
                $xmlInfo += Write-Output "GPO Domain   :$GPODomain"
                $xmlInfo += Write-Output "GPO Name     :$GPOName"
                $xmlInfo += Write-Output "Security Group where the members are configured to be added   :$GroupName"
                $xmlInfo += Write-Output "Members Name who are configured to be added to above group    :"
                $xmlInfo += $MembersName | ForEach-Object { "               $_" }
            }
        }
        if ($xmlInfo) {
            Write-Output "Information  from GPO XML:"
            Write-Output $xmlInfo
        } 
        else {
            Write-Output "There are no Preference Policies has been configured to update the Security Groups"
        }
    }
    else {
        Write-Output "xml file $GPOxmlFile is not present"
    }
}#try close
catch {
    Write-Output "[MSG: ERROR : $($_.Exception.message)]"
}
finally {
    if (Test-Path $GPOxmlFile) {
        Remove-Item $GPOxmlFile -Force -Confirm:$false -ErrorAction SilentlyContinue
    }
}
