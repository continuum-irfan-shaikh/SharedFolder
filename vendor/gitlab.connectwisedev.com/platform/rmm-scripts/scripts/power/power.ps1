$TurnOffDisplayAC = [int]$TurnOffDisplayAC * 60
$TurnOffDisplayDC = [int]$TurnOffDisplayDC * 60

function getPowerPlanIDByName {
    param([string]$name)

    $GuidString = Get-WmiObject -Namespace 'root\cimv2\power' -Class Win32_PowerPlan | Where{($_.ElementName -eq "$name")} | Select -ExpandProperty InstanceID
    if ($GuidString -ne $null) {
        $Regex = [regex]"{(.*?)}$"
        $PowerSchemeGuid = $Regex.Match($GuidString).Groups[1].Value
    }
    return $PowerSchemeGuid
}

function setUpPowerPlan {
    param(
        [guid]$guid,
        [int]$valueForAC,
        [int]$valueForDC,
        [bool]$setActiveValue
    )

    Invoke-Expression "Powercfg SetACValueIndex $guid SUB_VIDEO VIDEOIDLE $valueForAC" | Out-Null
    Invoke-Expression "Powercfg SetDCValueIndex $guid SUB_VIDEO VIDEOIDLE $valueForDC" | Out-Null
    if ($setActiveValue) {
        Invoke-Expression "Powercfg SetActive $guid" | Out-Null
    }
}

function configPowerPlan {
    param(
        [ValidateNotNullOrEmpty()] [string]$Action,
        [ValidateNotNullOrEmpty()] [string]$PowerPlan,
        [string]$BasePowerPlan,
        [ValidateScript({$_ -ge 0})] [int]$TurnOffDisplayAC,
        [ValidateScript({$_ -ge 0})] [int]$TurnOffDisplayDC,
        [bool]$SetActive
    )

    $BasePowerPlans = @("Balanced", "High Performance", "Power saver")

    switch ($Action) {
        "delete" {
            $CurrentPowerPlan = Get-WmiObject -Namespace 'root\cimv2\power' -Class Win32_PowerPlan | Where{$_.IsActive -eq $true} | Select -ExpandProperty ElementName
            If(($BasePowerPlans -contains "$PowerPlan") -or ($PowerPlan -eq "$CurrentPowerPlan")) {
                Write-Error "You can't delete the plans, which are provided by the PC manufacturer('Balanced', 'High Performance', 'Power saver'), or the plan that you're currently using"
                return
            }

            $PowerSchemeGuid = getPowerPlanIDByName -name $PowerPlan
            If([string]::IsNullOrEmpty($PowerSchemeGuid)) {
                Write-Error "The power plan with name '$PowerPlan' was not found, please input the correct value"
                return
            }
            Invoke-Expression "Powercfg /d $PowerSchemeGuid" | Out-Null
            Write-Output "Successfully removed the power plan:'$PowerPlan'"
            return
        }
        "create" {
            $PowerSchemeGuid = getPowerPlanIDByName -name $PowerPlan
            If($PowerSchemeGuid -ne $null) {
                Write-Error "The power plan with name '$PowerPlan' already exists"
                return
            }

            $PowerSchemeGuid = getPowerPlanIDByName -name $BasePowerPlan
            If([string]::IsNullOrEmpty($PowerSchemeGuid)) {
                Write-Error "The base power plan with name '$BasePowerPlan' was not found, please input the correct value"
                return
            }

            $NewPlanGuid = [guid]::NewGuid()
            Invoke-Expression "Powercfg DuplicateScheme $PowerSchemeGuid $NewPlanGuid" | Out-Null
            Invoke-Expression "Powercfg ChangeName $NewPlanGuid '$PowerPlan'" | Out-Null
            setUpPowerPlan -guid $NewPlanGuid -valueForAC $TurnOffDisplayAC -valueForDC $TurnOffDisplayDC -setActiveValue $SetActive
            Write-Output "Successfully created the power plan:'$PowerPlan'"
            return
        }
        "update" {
            $PowerSchemeGuid = getPowerPlanIDByName -name $PowerPlan
            If([string]::IsNullOrEmpty($PowerSchemeGuid)) {
                Write-Error "The power plan with name '$PowerPlan' was not found, please input the correct value"
                return
            }

            setUpPowerPlan -guid $PowerSchemeGuid -valueForAC $TurnOffDisplayAC -valueForDC $TurnOffDisplayDC -setActiveValue $SetActive
            Write-Output "Successfully updated the power plan:'$PowerPlan'"
            return
        }
        default {
            Write-Error "Invalid action: $Action. Should be one of: create, update, delete"
            return
        }
    }
}

if ($MyInvocation.InvocationName -ne '.')
{
    if ($SetActive -eq $null) {
        $SetActive = $false
    }
    configPowerPlan -Action $Action -PowerPlan $PowerPlan -BasePowerPlan $BasePowerPlan -TurnOffDisplayAC $TurnOffDisplayAC -TurnOffDisplayDC $TurnOffDisplayDC -SetActive $SetActive
}
