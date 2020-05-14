<#
    .SYNOPSIS
        Configure Power Management
    .DESCRIPTION
        Configure Power Management
    .HELP
        https://docs.microsoft.com/en-us/windows-hardware/customize/power-settings/configure-power-settings
        https://docs.microsoft.com/en-us/windows-hardware/design/device-experiences/powercfg-command-line-options
    .AUTHOR
        Durgeshkumar Patel
    .VERSION
        1.1
#>

<#
Action :- Create Scheme
          Modify Scheme
          Delete Scheme

SchemeName:-   String
SchemeBasedon :- Balanced
                   High Performance
                   Power Saver
SchemeDescription:- String

 If a scheme with this name already exists overwrite it
 
 Make this scheme as active power scheme

 Enable Hibernate Support

When Cumputer is Idle
Battery
Power Button and lid
Processor Power Management
PCI Express
Wireless Adapter Setting
Additional Options
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

try {
    $ErrorActionPreference = "stop"
    
    #Additional Change due to Json changes. Assigning values
    if (($Action -eq "Modify Scheme") -and ($ModifyScheme -eq "Default Scheme")) {
        switch ($SchemeBasedOn) {
            "1" { $SchemeName = "Balanced"; break }
            "2" { $SchemeName = "High Performance"; break }
            "3" { $SchemeName = "Power Saver"; break }
        }
    }

    $name = (Get-WmiObject -Namespace "root\cimv2\power" -class Win32_PowerPlan | Where-Object { $_.ElementName -eq "$SchemeName" } | Select -ExpandProperty InstanceID).substring(20).trim("{/}")
}
catch { $ErrorActionPreference = 'continue' }

$currnt_scheme = (Get-WmiObject -namespace "root\cimv2\power" -class Win32_powerplan | where { $_.IsActive } | Select -ExpandProperty InstanceID).substring(20).trim("{/}")

#Additional Change due to Json changes. Assigning values
switch ($SchemeBasedOn) {

    "1" { $SchemeBasedOn = "Balanced"; break }
    "2" { $SchemeBasedOn = "High Performance"; break }
    "3" { $SchemeBasedOn = "Power Saver"; break }
}

#Hardcoded GUIDs for Default Power Schemes
$Scheme_GUIDs = @{ "Balanced" = "381b4222-f694-41f0-9685-ff5bb260df2e"; "High performance" = "8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c"; "Power saver" = "a1841308-3541-4fab-bc81-f71556f20b4a" }

#Create Power Scheme
function Create_Scheme { 
    #If a scheme with this name already exists overwrite it
    if ((($IfSchemeExistsOverwrite -eq $false) -or ([string]::IsNullOrEmpty("$IfSchemeExistsOverwrite"))) -and (![string]::IsNullOrEmpty($name))) {
            
        Write-Output "`nPower scheme '$SchemeName' already exist, Kindly select Overwrite option to avoid duplication of power scheme."
        exit;
    }
    elseif (($IfSchemeExistsOverwrite -eq $true) -and (![string]::IsNullOrEmpty($name))) {
        $global:create_modify_scheme = $name
        Write-Output "`nPower scheme '$SchemeName' already exist and configuration will be overwritten."
    }
    else {
        $newguid = [guid]::NewGuid()
        try {
            $ErrorActionPreference = 'stop'
            POWERCFG /DUPLICATESCHEME $($Scheme_GUIDs["$SchemeBasedOn"]) $newguid | Out-Null

            powercfg -changename $newguid "$SchemeName" "$SchemeDescription"
          
            Start-Sleep 1 
            $global:create_modify_scheme = (Get-WmiObject -Namespace "root\cimv2\power" -class Win32_PowerPlan | Where-Object { $_.ElementName -eq "$SchemeName" } | Select -ExpandProperty InstanceID).substring(20).trim("{/}")
            Write-Output "`n'$SchemeName' power scheme created successfully with GUID $newguid"
        }
        catch { 
            $ErrorActionPreference = 'continue'    
            write-output "`nFailed to create '$SchemeName' power scheme"
            Write-Error $_.exception.message
        }
    }
}

#Delete Power Scheme
function Delete_Scheme {

    if ( [string]::IsNullOrEmpty("$name") ) {
        Write-Output "`nPower scheme '$SchemeName' not found"
        exit;
    }
    else {
        if ($name -ne $currnt_scheme) {    
            try {
                $ErrorActionPreference = 'stop'
                powercfg -delete "$name"
                Write-Output "`nPower scheme '$SchemeName' deleted successfully"
            }
            catch {
                $ErrorActionPreference = 'continue'    
                Write-Output "`nFailed to delete power scheme '$SchemeName'"
                Write-Error $_.exception.message
            }
        }
        else {
            Write-Output "`nThe active power scheme '$SchemeName' cannot be deleted"
            exit;
        }
    }
}

#Modify Power Scheme
function Modify_Scheme {

    if ( [string]::IsNullOrEmpty("$name") ) {
        Write-Output "`nPower scheme '$SchemeName' not found"
        exit;
    }
       
    if (($IfSchemeExistsOverwrite -eq $true) -or ($ModifyScheme -eq "User defined Scheme") -or ($ModifyScheme -eq "Default Scheme")) {
        try {
            $ErrorActionPreference = 'stop'
            $global:create_modify_scheme = $name 
           
            powercfg -changename $create_modify_scheme "$SchemeName" "$SchemeDescription"
            Write-Output "`nPower scheme '$SchemeName' modified successfully"
        }
        catch {
            $ErrorActionPreference = 'continue'    
            Write-Output "`nFailed to modified power scheme '$SchemeName'"
            Write-Error $_.exception.message
        }
    }
}

#unfriendly values
function Frnd_Calc ($Var1) {

    $F_Result = switch ($Var1) {
    
        "On" { "001"; break }
        "Off" { "000"; break }
        "Do Nothing" { "000"; break }
        "Sleep" { "001"; break }
        "Hibernate" { "002"; break }
        "Shut down" { "003"; break }
        "Active" { "001"; break }
        "Passive" { "000"; break }
        "Moderate power savings" { "001"; break }
        "Maximum power savings" { "002"; break }
        "Maximum Performance" { "000"; break }
        "Low Power Saving" { "001"; break }
        "Medium Power Saving" { "002"; break }
        "Maximum Power Saving" { "003"; break }
        "Yes" { "001"; break }
        "No" { "000"; break }
    } 
    return $F_Result
}
    
#Start menu power button action unfriednly value
function Frnd_Calc1 ($Stmp_Var) {
    
    $Stmp_Var_Result = switch ($Stmp_Var) {
    
        "Sleep" { "000"; break }
        "Hibernate" { "001"; break }
        "Shut down" { "002"; break }
    }
    return $Stmp_Var_Result
}

function Frnd_Calc2($Some_Val) {
    #Sleep action
    $Some_Result = switch ($Some_Val) {
        "1" { "Do Nothing"; break }
        "4" { "Sleep"; break }
        "5" { "Hibernate"; break }
        "6" { "Shut down"; break }
    }
    return $Some_Result
}

function Frnd_Calc3($Some_Val1) {
    $Some_Result1 = switch ($Some_Val1) {

        #On/Off 
        "1" { "On"; break }
        "0" { "Off"; break }
    }
    return $Some_Result1
}

function Frnd_Calc4 ($Sm) {

    $Sm_result = Switch ($Sm) {
        #Active/Passive
        "1" { "Active"; break }
        "0" { "Passive"; break }
    }
    return $Sm_result
}

function Frnd_Calc5 ($lnst) {
    #Power Saving
    $lnst_result = switch ($lnst) {
        "0" { "Off"; break }
        "1" { "Moderate power savings"; break }
        "2" { "Maximum power savings"; break }
    }
    return $lnst_result
}

function Frnd_Calc6($wapsm) {
    #Wireless power savings
    $wapsm_result = switch ($wapsm) {
        "1" { "Maximum Performance"; break }
        "2" { "Low Power Saving"; break }
        "3" { "Medium Power Saving"; break }
        "4" { "Maximum Power Saving"; break }
    }
    return $wapsm_result
}

function Frnd_Calc7($Some_V) {
    $Some_Re = switch ($Some_V) {

        #On/Off 
        "1" { "Yes"; break }
        "0" { "No"; break }
    }
    return $Some_Re
}

$Aftermin = @{
    "60"    = "After 1 min"
    "120"   = "After 2 mins"
    "180"   = "After 3 mins"
    "300"   = "After 5 mins"
    "600"   = "After 10 mins"
    "900"   = "After 15 mins"
    "1200"  = "After 20 mins"
    "1500"  = "After 25 mins"
    "1800"  = "After 30 mins"
    "2700"  = "After 45 mins"
    "3600"  = "After 1 hour"
    "7200"  = "After 2 hours"
    "10800" = "After 3 hours"
    "14400" = "After 4 hours"
    "18000" = "After 5 hours"
    "0"     = "Never"
}

#Action 
if ( $Action -eq "Create Scheme") {
    Create_Scheme
}
elseif ( $Action -eq "Modify Scheme" ) {
    Modify_Scheme
}
elseif ( $Action -eq "Delete Scheme") {
    Delete_Scheme
}
else {
    Write-Output "`nSelect Action first"
    exit;
}

#Update values when computer is plugged in/on Battery
if (($Action -eq "Create Scheme") -or ($Action -eq "Modify Scheme")) { 

    ##Predefined SubGUIDs ans SettingGUIDs  
    
    #If the computer is idle
    $Hard_Disk_SubGroup = "0012ee47-9041-4b5d-9b77-535fba8b1442"   #Hard Disk
    $Hard_Disk_SettingGUID = "6738e2c4-e8a5-4a42-b16a-e040e769756e"  #(Turn off hard disk after)
    
    $Sleep_SubGroup = "238c9fa8-0aad-41ed-83f4-97be242c8f20"  #(Sleep   Standby)
    $Sleep_SettingGUID = "29f6c1db-86da-48c5-9fdb-f2b67b1f44da"  #(Sleep after)
    $Hibernate_SettingGUID = "9d7815a6-7ee4-497e-8888-515a05f02364"  #(Hibernate after)
    
    $Display_SubGroup = "7516b95f-f776-4464-8c53-06167f40cc99"  #(Display)
    $TurnOfDisplay_SettingGUID = "3c0bc021-c8a8-4e07-a973-6b14cbcb2b7e"  #(Turn off display after)
    
    #Battery
    $Battery_Subgroup = "e73a048d-bf27-4f12-9731-8b2076e8891f"  #(Battery)
    $Critical_battery_level_SettingGUID = "9a66d8d7-4ff7-4ef9-b5a2-5a326ca2a469"  #(Critical battery level)
    $Critical_battery_action_SettingGUID = "637ea02f-bbcb-4015-8e2c-a1c7b9c0b546"  #(Critical battery action  -reaches critical)
    $Low_battery_level_SettingGUID = "8183ba9a-e910-48da-8769-14ae6dc1170a"  #(Low battery level)
    $Low_battery_action_SettingGUID = "d8742dcb-3e6a-4b3c-b3fe-374623cdcf06"  #(Low battery action   -reached low)
    $Low_battery_notification_SettingGUID = "bcded951-187b-4d05-bccc-f7e51960c258"  #(Low battery notification)
    
    #Power Button and Lid
    $Power_buttons_and_lid_Subgroup = "4f971e89-eebd-4455-a8de-9e59040e7347"  #(Power buttons and lid)
    $Lid_close_action_SettingGUID = "5ca83367-6e45-459f-a27b-476b1d01c936"  #(Lid close action)
    $Power_button_action_SettingGUID = "7648efa3-dd9c-4e3e-b566-50f929386280"  #(Power button action)
    $Sleep_button_action_SettingGUID = "96996bc0-ad50-47ec-923b-6f41874dd9eb"  #(Sleep button action)
    $Start_menu_power_button_SettingGUID = "a7066653-8d6c-40a8-910e-a1f54b84c7e5"  #(Start menu power button)
    
    #Processor Power Management
    $Processor_power_management_Subgroup = "54533251-82be-4824-96c1-47b60b740d00" #Processor power management)
    $Minimum_processor_state_SettingGUID = "893dee8e-2bef-41e0-89c6-b55d0929964c"  #(Minimum processor state)
    $Maximum_processor_state_SettingGUID = "bc5038f7-23e0-4960-96da-33abaf5935ec"  #(Maximum processor state)
    $System_cooling_policy_SettingGUID = "94d3a615-a899-4ac5-ae2b-e4d8f634367f"  #(System cooling policy)
    
    ##PCI Express
    $PCI_Express_Subgroup = "501a4d13-42af-4429-9fd1-a8218c268e20"  #(PCI Express)
    $Link_State_Power_Management_SettingGUID = "ee12f906-d277-404b-b6da-e5fa1a576df5"  #(Link State Power Management)
    
    ##Wireless Adapter Setting
    $Wireless_Adapter_Settings_Subgroup = "19cbb8fa-5279-450e-9fac-8a3d5fedd0c1"  #(Wireless Adapter Settings)
    $Power_Saving_Mode_SettingGUID = "12bbebe6-58d6-4636-95bb-3217ef867c1a"  #(Power Saving Mode)
    
    ##Additional Options
    $Additional_Options_Require_password_SubGroup = "fea3413e-7e05-4911-9a71-700331f1c294"    #SubGroup GUID Require a password on wakeup
    $Require_password_SettingGUID = "0e796bdb-100d-47d6-a2d5-f7d2daa51f51"   #Require a password on wakeup
    
    $Additional_Options_Subgroup = "238c9fa8-0aad-41ed-83f4-97be242c8f20"  #Additional_Options
    $Allow_hybrid_sleep_SettingGUID = "94ac6d29-73ce-41a6-809f-6363ba21b47e"  # Allow hybrid sleep
    $Allow_wake_timers_SettingGUID = "bd3b718a-0680-4d9d-8ab2-e1d2b4ac806d"  #Allow wake timers
        
    $USB_SubgroupGUID = "2a737441-1930-4402-8d77-b2bebba308a3"  #(USB settings)
    $USB_selective_SettingGUID = "48e6b7a6-50f5-4782-a5d4-53bb8f07e226"  #(USB selective suspend setting)
    
    #If Computer is idle  
    #Note:- Values for time is in seconds. Multiply (value*60)   
    # when plugged in
    Write-Output "`nIf Computer is idle :"
    try {
        $ErrorActionPreference = 'stop'
        #Turn of monitor when in AC mode
        powercfg /SETACVALUEINDEX $create_modify_scheme $Display_SubGroup $TurnOfDisplay_SettingGUID $TurnOffMonitorAC 
        Write-Output "Turn off monitor value set to '$($Aftermin["$TurnOffMonitorAC"])' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Turn off monitor value set to '$($Aftermin["$TurnOffMonitorAC"])' failed when plugged in" 
        Write-Error $_.exception.message
    }
        
    try {
        $ErrorActionPreference = 'stop'
        #Turn of hard disk when in AC mode
        powercfg /SETACVALUEINDEX $create_modify_scheme $Hard_Disk_SubGroup $Hard_Disk_SettingGUID $TurnOffHardDiskAC
        Write-Output "Turn off hard disk value set to '$($Aftermin["$TurnOffHardDiskAC"])' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Turn off hard disk value set to '$($Aftermin["$TurnOffHardDiskAC"])' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    try {
        $ErrorActionPreference = 'stop' 
        #Standby timeout whne in AC mode
        powercfg /SETACVALUEINDEX $create_modify_scheme $Sleep_SubGroup $Sleep_SettingGUID $StandbyAC
        Write-Output "Standby value set to '$($Aftermin["$StandbyAC"])' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Standby value set to '$($Aftermin["$StandbyAC"])' failed when plugged in"
        Write-Error $_.exception.message
    }
        
    try {
        $ErrorActionPreference = 'stop'
        #Hibernate when in AC mode
        powercfg /SETACVALUEINDEX $create_modify_scheme $Sleep_SubGroup $Hibernate_SettingGUID $HibernateAC
        Write-Output "Hibernate value $($Aftermin["$HibernateAC"]) updated when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Hibernate value set to '$($Aftermin["$HibernateAC"])' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    #when on battery
    try {
        $ErrorActionPreference = 'stop' 
        #Turn of monitor when in DC mode
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Display_SubGroup $TurnOfDisplay_SettingGUID $TurnOffMonitorDC    
        Write-Output "`nTurn off monitor value set to '$($Aftermin["$TurnOffMonitorDC"])' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Turn off monitor value set to '$($Aftermin["$TurnOffMonitorDC"])' failed when on battery"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        #Turn of hard disk when in DC mode
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Hard_Disk_SubGroup $Hard_Disk_SettingGUID $TurnOffHardDiskDC
        Write-Output "Turn off hard disk value set to '$($Aftermin["$TurnOffHardDiskDC"])' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Turn off hard disk value set to '$($Aftermin["$TurnOffHardDiskDC"])' failed when on battery"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        #Standby timeout when in DC mode
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Sleep_SubGroup $Sleep_SettingGUID $StandbyDC
        Write-Output "Standby value set to '$($Aftermin["$StandbyDC"])' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Standby value set to '$($Aftermin["$StandbyDC"])' failed when on battery"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        #Hibernate when in DC mode
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Sleep_SubGroup $Hibernate_SettingGUID $HibernateDC
        Write-Output "Hibernet value set to '$($Aftermin["$HibernateDC"])' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Hibernet value set to '$($Aftermin["$HibernateDC"])' failed when on battery"
        Write-Error $_.exception.message    
    }
        
    #Battery Settings
    Write-Output "`n`nBattery Settings :"
    # when plugged in
    try {
        $ErrorActionPreference = 'stop'
        #CriticalBatteryLevel  when plugged in               %value
        powercfg /SETACVALUEINDEX $create_modify_scheme $Battery_Subgroup $Critical_battery_level_SettingGUID $CriticalBatteryLevelAC
        Write-Output "Critical battery level value set to '$CriticalBatteryLevelAC %' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Critical battery level value set to '$CriticalBatteryLevelAC %' failed when plugged in"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        #IfBatteryLevelReachesCritical when plugged in     Do nothing/Sleep/Hibernate/Shutdown
        $BatteryLevelReachesCriticalAC_Name = Frnd_Calc2 $BatteryLevelReachesCriticalAC
        $BatteryLevelReachesCriticalAC_Value = Frnd_Calc $BatteryLevelReachesCriticalAC_Name
        powercfg /SETACVALUEINDEX $create_modify_scheme $Battery_Subgroup $Critical_battery_action_SettingGUID $BatteryLevelReachesCriticalAC_Value
        Write-Output "Critical battery level action value set to '$BatteryLevelReachesCriticalAC_Name' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Critical battery level action value set to '$BatteryLevelReachesCriticalAC_Name' failed when plugged in"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        #LowBatteryLevel  when plugged in  %value
        powercfg /SETACVALUEINDEX $create_modify_scheme $Battery_Subgroup $Low_battery_level_SettingGUID $LowBatteryLevelAC
        Write-Output "Low battery level value set to '$LowBatteryLevelAC %' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Low battery level value set to '$LowBatteryLevelAC %' failed when plugged in"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        #IfBatteryLevelReachesLow when on plugged in    Do nothing/Sleep/Hibernate/Shutdown
        $BatteryLevelReachesLowAC_Name = Frnd_Calc2 $BatteryLevelReachesLowAC
        $BatteryLevelReachesLowAC_Value = Frnd_Calc $BatteryLevelReachesLowAC_Name
        powercfg /SETACVALUEINDEX $create_modify_scheme $Battery_Subgroup $Low_battery_action_SettingGUID $BatteryLevelReachesLowAC_Value
        Write-Output "Low battery level action value set to '$BatteryLevelReachesLowAC_Name' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Low battery level action value set to '$BatteryLevelReachesLowAC_Name' failed when plugged in"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        #LowBatteryNotifucation when plugged in      On/Off
        $LowBatteryNotificationAC_Name = Frnd_Calc3 $LowBatteryNotificationAC
        $LowBatteryNotificationAC_Value = Frnd_Calc $LowBatteryNotificationAC_Name
        powercfg /SETACVALUEINDEX $create_modify_scheme $Battery_Subgroup $Low_battery_notification_SettingGUID $LowBatteryNotificationAC_Value
        Write-Output "Low battery notification value set to '$LowBatteryNotificationAC_Name' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Low battery notification value set to '$LowBatteryNotificationAC_Name' failed when plugged in"
        Write-Error $_.exception.message
    }    
    
    # when on battery
    try {
        $ErrorActionPreference = 'stop'
        #CriticalBatteryLevel  on Battery               %value
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Battery_Subgroup $Critical_battery_level_SettingGUID $CriticalBatteryLevelDC
        Write-Output "`nCritical battery level value set to '$CriticalBatteryLevelDC %' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Critical battery level value set to '$CriticalBatteryLevelDC %' failed when on battery"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        #IfBatteryLevelReachesCritical when on battery    Do nothing/Sleep/Hibernate/Shutdown
        $BatteryLevelReachesCriticalDC_Name = Frnd_Calc2 $BatteryLevelReachesCriticalDC
        $BatteryLevelReachesCriticalDC_Value = Frnd_Calc $BatteryLevelReachesCriticalDC_Name
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Battery_Subgroup $Critical_battery_action_SettingGUID $BatteryLevelReachesCriticalDC_Value
        Write-Output "Critical battery level action value set to '$BatteryLevelReachesCriticalDC_Name' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Critical battery level action value set to '$BatteryLevelReachesCriticalDC_Name' failed when on battery"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        #LowBatteryLevel  when on battery  %value
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Battery_Subgroup $Low_battery_level_SettingGUID $LowBatteryLevelDC
        Write-Output "Low battery level value set to '$LowBatteryLevelDC %' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Low battery level value set to '$LowBatteryLevelDC %' failed when on battery"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        #IfBatteryLevelReachesLow when on battery    Do nothing/Sleep/Hibernate/Shutdown
        $BatteryLevelReachesLowDC_Name = Frnd_Calc2 $BatteryLevelReachesLowDC
        $BatteryLevelReachesLowDC_Value = Frnd_Calc $BatteryLevelReachesLowDC_Name
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Battery_Subgroup $Low_battery_action_SettingGUID $BatteryLevelReachesLowDC_Value
        Write-Output "Low battery level action value set to '$BatteryLevelReachesLowDC_Name' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Low battery level action value set to '$BatteryLevelReachesLowDC_Name' failed when on battery"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        #LowBatteryNotification on Battery       On/Off
        $LowBatteryNotificationDC_Name = Frnd_Calc3 $LowBatteryNotificationDC
        $LowBatteryNotificationDC_Value = Frnd_Calc $LowBatteryNotificationDC_Name
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Battery_Subgroup $Low_battery_notification_SettingGUID $LowBatteryNotificationDC_Value
        Write-Output "Low battery level notification value set to '$LowBatteryNotificationDC_Name' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Low battery level notification value set to '$LowBatteryNotificationDC_Name' failed when on battery"
        Write-Error $_.exception.message
    }
    
    #Power Button and lid
    Write-Output "`n`nPower Button and lid :"
    #when plugged in
    try {
        $ErrorActionPreference = 'stop'
        $WhenCloseTheLidAC_Name = Frnd_Calc2 $WhenCloseTheLidAC
        $WhenCloseTheLidAC_Value = Frnd_Calc $WhenCloseTheLidAC_Name
        # Close the lid when plugged in 
        powercfg /SETACVALUEINDEX $create_modify_scheme $Power_buttons_and_lid_Subgroup $Lid_close_action_SettingGUID $WhenCloseTheLidAC_Value
        Write-Output "When close the lid value set to '$WhenCloseTheLidAC_Name' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "When close the lid value set to '$WhenCloseTheLidAC_Name' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    try {
        $ErrorActionPreference = 'stop'
        $WhenPressThePowerButtonAC_Name = Frnd_Calc2 $WhenPressThePowerButtonAC
        $WhenPressThePowerButtonAC_Value = Frnd_Calc $WhenPressThePowerButtonAC_Name
        #Press Power Button when plugged in
        powercfg /SETACVALUEINDEX $create_modify_scheme $Power_buttons_and_lid_Subgroup $Power_button_action_SettingGUID $WhenPressThePowerButtonAC_Value
        Write-Output "Press power button value set to '$WhenPressThePowerButtonAC_Name' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Press power button value set to '$WhenPressThePowerButtonAC_Name' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    try {
        $ErrorActionPreference = 'stop'
        $WhenPressTheSleepButtonAC_Name = Frnd_Calc2 $WhenPressTheSleepButtonAC
        $WhenPressTheSleepButtonAC_Value = Frnd_Calc $WhenPressTheSleepButtonAC_Name
        #Press Sleep Button when plugged in
        powercfg /SETACVALUEINDEX $create_modify_scheme $Power_buttons_and_lid_Subgroup $Sleep_button_action_SettingGUID $WhenPressTheSleepButtonAC_Value
        Write-Output "Press sleep button value set to '$WhenPressTheSleepButtonAC_Name' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Press sleep button value set to '$WhenPressTheSleepButtonAC_Name' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    try {
        $ErrorActionPreference = 'stop'
        if ($WhenPressTheStartMenuPowerButtonAC -ne "1") { 
            $WhenPressTheStartMenuPowerButtonAC_Name = Frnd_Calc2 $WhenPressTheStartMenuPowerButtonAC
            $WhenPressTheStartMenuPowerButtonAC_Val = Frnd_Calc1 $WhenPressTheStartMenuPowerButtonAC_Name
            #Press Start Menu Power Button when plugged in
            powercfg /SETACVALUEINDEX $create_modify_scheme $Power_buttons_and_lid_Subgroup $Start_menu_power_button_SettingGUID $WhenPressTheStartMenuPowerButtonAC_Val
            Write-Output "Press start menu power button value set to '$WhenPressTheStartMenuPowerButtonAC_Name' when plugged in"
        }
        else {
            Write-Output "Press start menu power button value set to 'Do Nothing' when plugged in"
        }
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Press start menu power button value set to '$WhenPressTheStartMenuPowerButtonAC_Name' failed when plugged in"
        Write-Error $_.exception.message
    }
        
    # when on battery
    try {
        $ErrorActionPreference = 'stop'
        $WhenCloseTheLidDC_Name = Frnd_Calc2 $WhenCloseTheLidDC
        $WhenCloseTheLidDC_Val = Frnd_Calc $WhenCloseTheLidDC_Name
        # Close the lid when on battery 
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Power_buttons_and_lid_Subgroup $Lid_close_action_SettingGUID $WhenCloseTheLidDC_Val
        Write-Output "`nClose the lid value set to '$WhenCloseTheLidDC_Name' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Close the lid value set to '$WhenCloseTheLidDC_Name' failed when on battery"
        Write-Error $_.exception.message
    }
    
    try {
        $ErrorActionPreference = 'stop'
        $WhenPressThePowerButtonDC_Name = Frnd_Calc2 $WhenPressThePowerButtonDC
        $WhenPressThePowerButtonDC_Val = Frnd_Calc $WhenPressThePowerButtonDC_Name
        #Press Power Button when on battery
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Power_buttons_and_lid_Subgroup $Power_button_action_SettingGUID $WhenPressThePowerButtonDC_Val
        Write-Output "Press power button value set to '$WhenPressThePowerButtonDC_Name' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Press power button value set to '$WhenPressThePowerButtonDC_Name' failed when on battery"
        Write-Error $_.exception.message
    }
        
    try {
        $ErrorActionPreference = 'stop'
        $WhenPressTheSleepButtonDC_Name = Frnd_Calc2 $WhenPressTheSleepButtonDC
        $WhenPressTheSleepButtonDC_Val = Frnd_Calc $WhenPressTheSleepButtonDC_Name
        #Press Sleep Button when on battery
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Power_buttons_and_lid_Subgroup $Sleep_button_action_SettingGUID $WhenPressTheSleepButtonDC_Val
        Write-Output "Press sleep button value set to '$WhenPressTheSleepButtonDC_Name' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Press sleep button value set to '$WhenPressTheSleepButtonDC_Name' failed when on battery"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        if ($WhenPressTheStartMenuPowerButtonDC -ne "1") { 
            $WhenPressTheStartMenuPowerButtonDC_Name = Frnd_Calc2 $WhenPressTheStartMenuPowerButtonDC
            $WhenPressTheStartMenuPowerButtonDC_val = Frnd_Calc1 $WhenPressTheStartMenuPowerButtonDC_Name
            #Press Start Menu Power Button when on battery
            powercfg /SETDCVALUEINDEX $create_modify_scheme $Power_buttons_and_lid_Subgroup $Start_menu_power_button_SettingGUID $WhenPressTheStartMenuPowerButtonDC_val
            Write-Output "Press start menu power button value set to '$WhenPressTheStartMenuPowerButtonDC_Name' when on battery"
        }
        else {
            Write-Output "Press start menu power button value set to 'Do Nothing' when on battery"
        }
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Press start menu power button value set to '$WhenPressTheStartMenuPowerButtonDC_Name' failed when on battery"
        Write-Error $_.exception.message
    }    
    
    #Processor Power Management
    Write-Output "`n`nProcessor Power Management :"
    # when plugged in
    try {
        $ErrorActionPreference = 'stop'
        # Minimum Processor State when pluged in
        powercfg /SETACVALUEINDEX $create_modify_scheme $Processor_power_management_Subgroup $Minimum_processor_state_SettingGUID $MinimumProcessorStateAC
        Write-Output "Minimum processor state value set to '$MinimumProcessorStateAC %' when plugged in"
    }

    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Minimum processor state value set to '$MinimumProcessorStateAC %' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    try {
        $ErrorActionPreference = 'stop'
        # Maximum Processor State when pluged in
        powercfg /SETACVALUEINDEX $create_modify_scheme $Processor_power_management_Subgroup $Maximum_processor_state_SettingGUID $MaximumProcessorStateAC
        Write-Output "Maximum processor state value set to '$MaximumProcessorStateAC %' when plugged in"
    }
    
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Maximum processor state value set to '$MaximumProcessorStateAC %' failed when plugged in"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        $SystemCoolingPolicyAC_Name = Frnd_Calc4 $SystemCoolingPolicyAC
        $SystemCoolingPolicyAC_Val = Frnd_Calc $SystemCoolingPolicyAC_Name
        #SystemCoolingPolicy when plugged in
        powercfg /SETACVALUEINDEX $create_modify_scheme $Processor_power_management_Subgroup $System_cooling_policy_SettingGUID $SystemCoolingPolicyAC_Val
        Write-Output "System cooling policy value set to '$SystemCoolingPolicyAC_Name' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "System cooling policy value set to '$SystemCoolingPolicyAC_Name' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    # when on battery
    try {
        $ErrorActionPreference = 'stop'
        # Minimum Processor State when on battery
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Processor_power_management_Subgroup $Minimum_processor_state_SettingGUID $MinimumProcessorStateDC
        Write-Output "`nMinimum processor state value set to '$MinimumProcessorStateDC %' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Minimum processor state value set to '$MinimumProcessorStateDC %' failed when on battery"
        Write-Error $_.exception.message
    }
    
    try {
        $ErrorActionPreference = 'stop'
        # Maximum Processor State when on battery
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Processor_power_management_Subgroup $Maximum_processor_state_SettingGUID $MaximumProcessorStateDC
        Write-Output "Maximum processor state value set to '$MaximumProcessorStateDC %' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Maximum processor state value set to '$MaximumProcessorStateDC %' failed when on battery"
        Write-Error $_.exception.message
    }
    
    try {
        $ErrorActionPreference = 'stop'
        $SystemCoolingPolicyDC_Name = Frnd_Calc4 $SystemCoolingPolicyDC
        $SystemCoolingPolicyDC_Val = Frnd_Calc $SystemCoolingPolicyDC_Name
        #SystemCoolingPolicy when on battery
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Processor_power_management_Subgroup $System_cooling_policy_SettingGUID $SystemCoolingPolicyDC_Val
        Write-Output "System cooling policy value set to '$SystemCoolingPolicyDC_Name' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "System cooling policy value set to '$SystemCoolingPolicyDC_Name' failed when on battery"
        Write-Error $_.exception.message
    }
    
    #PCI Express
    Write-Output "`n`nPCI Express :"
    # when plugged in 
    try {
        $ErrorActionPreference = 'stop'
        $LinkStatePowerManagementAC_Name = Frnd_Calc5 $LinkStatePowerManagementAC
        $LinkStatePowerManagementAC_Val = Frnd_Calc $LinkStatePowerManagementAC_Name
        #Link State Power Management when plugged in
        powercfg /SETACVALUEINDEX $create_modify_scheme $PCI_Express_Subgroup $Link_State_Power_Management_SettingGUID $LinkStatePowerManagementAC_Val 
        Write-Output "Link state power management value set to '$LinkStatePowerManagementAC_Name' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Link state power management value set to '$LinkStatePowerManagementAC_Name' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    # when on battery
    try {
        $ErrorActionPreference = 'stop'   
        $LinkStatePowerManagementDC_Name = Frnd_Calc5 $LinkStatePowerManagementDC
        $LinkStatePowerManagementDC_Val = Frnd_Calc $LinkStatePowerManagementDC_Name
        #Link State Power Management when on battery
        powercfg /SETDCVALUEINDEX $create_modify_scheme $PCI_Express_Subgroup $Link_State_Power_Management_SettingGUID $LinkStatePowerManagementDC_Val
        Write-Output "`nLink state power management value set to '$LinkStatePowerManagementDC_Name' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Link state power management value set to '$LinkStatePowerManagementDC_Name' failed when on battery"
        Write-Error $_.exception.message
    }
    
    #Wireless Adapter Setting
    Write-Output "`n`nWireless Adapter Setting :"
    # when plugged in 
    try {
        $ErrorActionPreference = 'stop'  
        $WirelessAdapterPowerSavingModeAC_Name = Frnd_Calc6 $WirelessAdapterPowerSavingModeAC
        $WirelessAdapterPowerSavingModeAC_V = Frnd_Calc $WirelessAdapterPowerSavingModeAC_Name
        #Wireless Power Savin Mode when plugged in
        powercfg /SETACVALUEINDEX $create_modify_scheme $Wireless_Adapter_Settings_Subgroup $Power_Saving_Mode_SettingGUID $WirelessAdapterPowerSavingModeAC_V
        
        Write-Output "Wireless power saving value set to '$WirelessAdapterPowerSavingModeAC_Name' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Wireless power saving value set to '$WirelessAdapterPowerSavingModeAC_Name' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    # when on battery
    try {
        $ErrorActionPreference = 'stop'  
        $WirelessAdapterPowerSavingModeDC_Name = Frnd_Calc6 $WirelessAdapterPowerSavingModeDC
        $WirelessAdapterPowerSavingModeDC_V = Frnd_Calc $WirelessAdapterPowerSavingModeDC_Name
        #Wireless Power Savin Mode when plugged in
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Wireless_Adapter_Settings_Subgroup $Power_Saving_Mode_SettingGUID $WirelessAdapterPowerSavingModeDC_V
        
        Write-Output "`nWireless power saving value set to '$WirelessAdapterPowerSavingModeDC_Name' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Wireless power saving value set to '$WirelessAdapterPowerSavingModeDC_Name' failed when on battery"
        Write-Error $_.exception.message
    }
    
    #Additional Options
    Write-Output "`n`nAdditional Options :"
    # when plugged in 
    try {
        $ErrorActionPreference = 'stop'
        $RequirePasswordOnWakeupAC_N = Frnd_Calc7 $RequirePasswordOnWakeupAC
        $RequirePasswordOnWakeupAC_V = Frnd_Calc $RequirePasswordOnWakeupAC_N
        # Require a password on wakeup when plugged in
        powercfg /SETACVALUEINDEX $create_modify_scheme $Additional_Options_Require_password_SubGroup $Require_password_SettingGUID $RequirePasswordOnWakeupAC_V
        Write-Output "Require a password on wakeup value set to '$RequirePasswordOnWakeupAC_N' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Require a password on wakeup value set to '$RequirePasswordOnWakeupAC_N' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    try {
        $ErrorActionPreference = 'stop'
        $AllowHybridSleepAC_N = Frnd_Calc7 $AllowHybridSleepAC
        $AllowHybridSleepAC_V = Frnd_Calc $AllowHybridSleepAC_N
        # Allow hybrid sleep when plugged in
        powercfg /SETACVALUEINDEX $create_modify_scheme $Additional_Options_Subgroup $Allow_hybrid_sleep_SettingGUID $AllowHybridSleepAC_V
        Write-Output "Allow hybrid sleep value set to '$AllowHybridSleepAC_N' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Allow hybrid sleep value set to '$AllowHybridSleepAC_N' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    try {
        $ErrorActionPreference = 'stop'
        $AllowWakeTimersAC_N = Frnd_Calc7 $AllowWakeTimersAC
        $AllowWakeTimersAC_V = Frnd_Calc $AllowWakeTimersAC_N
        #Allow wake timers when plugged in
        powercfg /SETACVALUEINDEX $create_modify_scheme $Additional_Options_Subgroup $Allow_wake_timers_SettingGUID $AllowWakeTimersAC_V
        Write-Output "Allow wake timers value set to '$AllowWakeTimersAC_N' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Allow wake timers value set to '$AllowWakeTimersAC_N' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    try {
        $ErrorActionPreference = 'stop'
        $USBSelectiveSuspendSettingAC_N = Frnd_Calc7 $USBSelectiveSuspendSettingAC
        $USBSelectiveSuspendSettingAC_V = Frnd_Calc $USBSelectiveSuspendSettingAC_N
        #USB selective suspend setting when plugged in
        powercfg /SETACVALUEINDEX $create_modify_scheme $USB_SubgroupGUID $USB_selective_SettingGUID $USBSelectiveSuspendSettingAC_V
        Write-Output "USB selective suspend setting value set to '$USBSelectiveSuspendSettingAC_N' when plugged in"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "USB selective suspend setting value set to '$USBSelectiveSuspendSettingAC_N' failed when plugged in"
        Write-Error $_.exception.message
    }
    
    # when on battery
    try {
        $ErrorActionPreference = 'stop'
        $RequirePasswordOnWakeupDC_N = Frnd_Calc7 $RequirePasswordOnWakeupDC
        $RequirePasswordOnWakeupDC_V = Frnd_Calc $RequirePasswordOnWakeupDC_N
        # Require a password on wakeup when on battery
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Additional_Options_Require_password_SubGroup $Require_password_SettingGUID $RequirePasswordOnWakeupDC_V
        Write-Output "`nRequire a password on wakeup value set to '$RequirePasswordOnWakeupDC_N' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Require a password on wakeup value set to '$RequirePasswordOnWakeupDC_N' failed when on battery"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        $AllowHybridSleepDC_N = Frnd_Calc7 $AllowHybridSleepDC
        $AllowHybridSleepDC_V = Frnd_Calc $AllowHybridSleepDC_N
        # Allow hybrid sleep when on battery
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Additional_Options_Subgroup $Allow_hybrid_sleep_SettingGUID $AllowHybridSleepDC_V
        Write-Output "Allow hybrid sleep value set to '$AllowHybridSleepDC_N' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Allow hybrid sleep value set to '$AllowHybridSleepDC_N' failed when on battery"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        $AllowWakeTimersDC_N = Frnd_Calc7 $AllowWakeTimersDC
        $AllowWakeTimersDC_V = Frnd_Calc $AllowWakeTimersDC_N
        #Allow wake timers when on battery
        powercfg /SETDCVALUEINDEX $create_modify_scheme $Additional_Options_Subgroup $Allow_wake_timers_SettingGUID $AllowWakeTimersDC_V
        Write-Output "Allow wake timers value set to '$AllowWakeTimersDC_N' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "Allow wake timers value set to '$AllowWakeTimersDC_N' failed when on battery"
        Write-Error $_.exception.message
    }

    try {
        $ErrorActionPreference = 'stop'
        $USBSelectiveSuspendSettingDC_N = Frnd_Calc7 $USBSelectiveSuspendSettingDC
        $USBSelectiveSuspendSettingDC_V = Frnd_Calc $USBSelectiveSuspendSettingDC_N
        #USB selective suspend setting when on battery
        powercfg /SETDCVALUEINDEX $create_modify_scheme $USB_SubgroupGUID $USB_selective_SettingGUID $USBSelectiveSuspendSettingDC_V
        Write-Output "USB selective suspend setting value set to '$USBSelectiveSuspendSettingDC_N' when on battery"
    }
    catch {
        $ErrorActionPreference = 'continue'    
        Write-Output "USB selective suspend setting value set to '$USBSelectiveSuspendSettingDC_N' failed when on battery"
        Write-Error $_.exception.message
    }  

    #Set active power scheme
    if ($MakeSchemeActive -eq $true) {
        if (![string]::IsNullOrEmpty("$create_modify_scheme")) {
        
            try {
                $ErrorActionPreference = 'stop'
                powercfg -SETACTIVE "$create_modify_scheme"
                Write-Output "`nScheme '$SchemeName' is now active power scheme"
            }
            catch {
                $ErrorActionPreference = 'continue'    
                Write-Output "`nFailed to set active '$SchemeName' power scheme"
                Write-Error $_.exception.message
            }
        }
    }

    #Enable Hibernate Support
    if ($EnableHibernate -eq $true) {

        $ProcessInfo = New-Object System.Diagnostics.ProcessStartInfo 
        $ProcessInfo.FileName = "powercfg.exe" 
        $ProcessInfo.Arguments = "/hibernate on "
        $ProcessInfo.RedirectStandardError = $true 
        $ProcessInfo.RedirectStandardOutput = $true 
        $ProcessInfo.CreateNoWindow = $true
        $ProcessInfo.UseShellExecute = $false 
        $Process = New-Object System.Diagnostics.Process 
        $Process.StartInfo = $ProcessInfo 
        $Process.Start() | Out-Null
        $Process.WaitForExit() 
        $Err = (($Process.StandardError.ReadToEnd()).split("`n"))[0]    
        if ($Process.ExitCode -eq 0) {
            Write-Output "`nHibernate Enabled" 
        }     
        elseif ($Err -like "*The request is not supported*") {
            Write-Output "`n$Err"
        }
        else { Write-Output "`nSomething went wrong while enabling the Hibernation. Kindly check manually" }     
    }

}

