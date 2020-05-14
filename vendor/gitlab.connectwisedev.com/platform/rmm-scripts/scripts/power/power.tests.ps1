$here = (Split-Path -Parent $MyInvocation.MyCommand.Path)
. $here\power.ps1

function removePowerPlan {
    param([string]$name)

    $GuidString = Get-WmiObject -Namespace 'root\cimv2\power' -Class Win32_PowerPlan | Where{($_.ElementName -eq "$name")} | Select -ExpandProperty InstanceID
    $Regex = [regex]"{(.*?)}$"
    $PowerSchemeGuid = $Regex.Match($GuidString).Groups[1].Value
    Invoke-Expression "Powercfg /d $PowerSchemeGuid" | Out-Null
}

function createPowerPlan {
    param([string]$name)

    $NewPlanGuid = [guid]::NewGuid()
    $BasePowerPlan = getPowerPlanIDByName -name "Balanced"
    Invoke-Expression "Powercfg DuplicateScheme $BasePowerPlan $NewPlanGuid" | Out-Null
    Invoke-Expression "Powercfg ChangeName $NewPlanGuid '$name'" | Out-Null
}

Describe 'ConfigPowerPlan' {
    It "Invalid Action, error: Invalid action: <action>. Should be one of: create, update, delete" {
        [string]$NewPowerPlan = -join ((65..90) | Get-Random -Count 32 | % {[char]$_})
        [string]$Action = 'coding'
        [string]$BasePowerPlan = 'Balanced'
        [int]$TurnOffDisplayAC = 120
        [int]$TurnOffDisplayDC = 360
        [bool]$SetActive = $false

        $result = configPowerPlan -Action $Action -PowerPlan $NewPowerPlan -BasePowerPlan $BasePowerPlan -TurnOffDisplayAC $TurnOffDisplayAC -TurnOffDisplayDC $TurnOffDisplayDC -SetActive $SetActive
        $result | Should -Be $null
        $error[0].ToString() | Should -Be "Invalid action: $Action. Should be one of: create, update, delete"
    }

    It "Create Power Plan, success" {
        [string]$NewPowerPlan = -join ((65..90) | Get-Random -Count 32 | % {[char]$_})
        [string]$Action = 'create'
        [string]$BasePowerPlan = 'Balanced'
        [int]$TurnOffDisplayAC = 120
        [int]$TurnOffDisplayDC = 360
        [bool]$SetActive = $false

        $result = configPowerPlan -Action $Action -PowerPlan $NewPowerPlan -BasePowerPlan $BasePowerPlan -TurnOffDisplayAC $TurnOffDisplayAC -TurnOffDisplayDC $TurnOffDisplayDC -SetActive $SetActive
        $result | Should -Be "Successfully created the power plan:'$NewPowerPlan'"

        # cleaning up
        removePowerPlan -name $NewPowerPlan
    }

    It "Create Power Plan, error: The power plan with name '<name>' already exists" {
        [string]$NewPowerPlan = 'Balanced'
        [string]$Action = 'create'
        [string]$BasePowerPlan = 'Balanced'
        [int]$TurnOffDisplayAC = 120
        [int]$TurnOffDisplayDC = 360
        [bool]$SetActive = $false

        $result = configPowerPlan -Action $Action -PowerPlan $NewPowerPlan -BasePowerPlan $BasePowerPlan -TurnOffDisplayAC $TurnOffDisplayAC -TurnOffDisplayDC $TurnOffDisplayDC -SetActive $SetActive
        $result | Should -Be $null
        $error[0].ToString() | Should -Be "The power plan with name '$NewPowerPlan' already exists"
    }

    It "Create Power Plan, error: The base power plan with name '<name>' was not found, please input the correct value" {
        [string]$Action = 'create'
        [string]$NewPowerPlan = -join ((65..90) | Get-Random -Count 32 | % {[char]$_})
        [string]$BasePowerPlan = -join ((65..90) | Get-Random -Count 32 | % {[char]$_})
        [int]$TurnOffDisplayAC = 120
        [int]$TurnOffDisplayDC = 360
        [bool]$SetActive = $false

        $result = configPowerPlan -Action $Action -PowerPlan $NewPowerPlan -BasePowerPlan $BasePowerPlan -TurnOffDisplayAC $TurnOffDisplayAC -TurnOffDisplayDC $TurnOffDisplayDC -SetActive $SetActive
        $result | Should -Be $null
        $error[0].ToString() | Should -Be "The base power plan with name '$BasePowerPlan' was not found, please input the correct value"
    }

    It "Delete Power Plan, success" {
        [string]$Action = 'delete'
        [string]$PowerPlan = -join ((65..90) | Get-Random -Count 32 | % {[char]$_})
        [bool]$SetActive = $false

        # creating test power plan
        createPowerPlan -name $PowerPlan

        $result = configPowerPlan -Action $Action -PowerPlan $PowerPlan -BasePowerPlan $BasePowerPlan -TurnOffDisplayAC $TurnOffDisplayAC -TurnOffDisplayDC $TurnOffDisplayDC -SetActive $SetActive
        $result | Should -Be "Successfully removed the power plan:'$PowerPlan'"
    }

    It "Delete Power Plan, error: You can't delete the plans, which are provided by the PC manufacturer('Balanced', 'High Performance', 'Power saver'), ..." {
        [string]$Action = 'delete'
        [string]$PowerPlan = 'Balanced'
        [bool]$SetActive = $false

        $result = configPowerPlan -Action $Action -PowerPlan $PowerPlan -BasePowerPlan $BasePowerPlan -TurnOffDisplayAC $TurnOffDisplayAC -TurnOffDisplayDC $TurnOffDisplayDC -SetActive $SetActive
        $result | Should -Be $null
        $error[0].ToString() | Should -Be "You can't delete the plans, which are provided by the PC manufacturer('Balanced', 'High Performance', 'Power saver'), or the plan that you're currently using"
    }

    It "Delete Power Plan, error: You can't delete the plans, ..., or the plan that you're currently using" {
        $CurrentPowerPlan = Get-WmiObject -Namespace 'root\cimv2\power' -Class Win32_PowerPlan | Where{$_.IsActive -eq $true} | Select -ExpandProperty ElementName
        [string]$PowerPlan = $CurrentPowerPlan
        [string]$Action = 'delete'
        [bool]$SetActive = $false

        $result = configPowerPlan -Action $Action -PowerPlan $PowerPlan -BasePowerPlan $BasePowerPlan -TurnOffDisplayAC $TurnOffDisplayAC -TurnOffDisplayDC $TurnOffDisplayDC -SetActive $SetActive
        $result | Should -Be $null
        $error[0].ToString() | Should -Be "You can't delete the plans, which are provided by the PC manufacturer('Balanced', 'High Performance', 'Power saver'), or the plan that you're currently using"
    }

    It "Delete Power Plan, error: The power plan with name '<name>' was not found, please input the correct value" {
        [string]$Action = 'delete'
        [string]$PowerPlan = -join ((65..90) | Get-Random -Count 32 | % {[char]$_})
        [bool]$SetActive = $false

        $result = configPowerPlan -Action $Action -PowerPlan $PowerPlan -BasePowerPlan $BasePowerPlan -TurnOffDisplayAC $TurnOffDisplayAC -TurnOffDisplayDC $TurnOffDisplayDC -SetActive $SetActive
        $result | Should -Be $null
        $error[0].ToString() | Should -Be "The power plan with name '$PowerPlan' was not found, please input the correct value"
    }

    It "Update Power Plan, success with SetActive == false" {
        [string]$Action = 'update'
        [string]$PowerPlan = -join ((65..90) | Get-Random -Count 32 | % {[char]$_})
        [int]$TurnOffDisplayAC = 120
        [int]$TurnOffDisplayDC = 360
        [bool]$SetActive = $false

        # creating test power plan
        createPowerPlan -name $PowerPlan

        $result = configPowerPlan -Action $Action -PowerPlan $PowerPlan -BasePowerPlan $BasePowerPlan -TurnOffDisplayAC $TurnOffDisplayAC -TurnOffDisplayDC $TurnOffDisplayDC -SetActive $SetActive
        $result | Should -Be "Successfully updated the power plan:'$PowerPlan'"

        # checking whether the updated power plan is not an active plan
        $activePowerPlan = Get-WmiObject -Namespace 'root\cimv2\power' -Class Win32_PowerPlan | Where{$_.IsActive -eq $true} | Select -ExpandProperty ElementName
        $activePowerPlan | Should -Not -Be $PowerPlan

        # cleaning up
        removePowerPlan -name $PowerPlan
    }

    It "Update Power Plan, success with SetActive == true" {
        $CurrentPowerPlan = Get-WmiObject -Namespace 'root\cimv2\power' -Class Win32_PowerPlan | Where{$_.IsActive -eq $true} | Select -ExpandProperty ElementName
        $CurrentPowerPlanID = getPowerPlanIDByName -name $CurrentPowerPlan
        [string]$Action = 'update'
        [string]$PowerPlan = -join ((65..90) | Get-Random -Count 32 | % {[char]$_})
        [int]$TurnOffDisplayAC = 120
        [int]$TurnOffDisplayDC = 360
        [bool]$SetActive = $true

        # creating test power plan
        createPowerPlan -name $PowerPlan

        $result = configPowerPlan -Action $Action -PowerPlan $PowerPlan -BasePowerPlan $BasePowerPlan -TurnOffDisplayAC $TurnOffDisplayAC -TurnOffDisplayDC $TurnOffDisplayDC -SetActive $SetActive
        $result | Should -Be "Successfully updated the power plan:'$PowerPlan'"

        # checking whether the updated power plan is an active plan
        $activePowerPlan = Get-WmiObject -Namespace 'root\cimv2\power' -Class Win32_PowerPlan | Where{$_.IsActive -eq $true} | Select -ExpandProperty ElementName
        $activePowerPlan | Should -Be $PowerPlan

        # cleaning up
        Invoke-Expression "Powercfg SetActive $CurrentPowerPlanID" | Out-Null
        removePowerPlan -name $PowerPlan
    }

    It "Update Power Plan, error: The power plan with name '<name>' was not found, please input the correct value" {
        [string]$Action = 'update'
        [string]$PowerPlan = -join ((65..90) | Get-Random -Count 32 | % {[char]$_})
        [int]$TurnOffDisplayAC = 120
        [int]$TurnOffDisplayDC = 360
        [bool]$SetActive = $false

        $result = configPowerPlan -Action $Action -PowerPlan $PowerPlan -BasePowerPlan $BasePowerPlan -TurnOffDisplayAC $TurnOffDisplayAC -TurnOffDisplayDC $TurnOffDisplayDC -SetActive $SetActive
        $result | Should -Be $null
        $error[0].ToString() | Should -Be "The power plan with name '$PowerPlan' was not found, please input the correct value"
    }
}
