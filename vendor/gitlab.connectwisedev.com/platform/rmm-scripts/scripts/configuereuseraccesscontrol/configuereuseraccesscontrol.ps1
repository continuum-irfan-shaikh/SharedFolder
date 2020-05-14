<#
 $uacSettings =  $DefaultNotify/$ConditionalNotify/$AlwaysNotify/$NeverNotify
#>
if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

$executionlog = @()
if($uacSettings -eq "DefaultNotify") {
$level1 = "Disabled"
$level2 = "Disabled"
$level3 = "Enabled"
$level4 = "Disabled"} elseif ($uacSettings -eq "ConditionalNotify") {
$level1 = "Disabled"
$level2 = "Enabled"
$level3 = "Disabled"
$level4 = "Disabled"} elseif ($uacSettings -eq "AlwaysNotify") {
$level1 = "Disabled"
$level2 = "Disabled"
$level3 = "Disabled"
$level4 = "Enabled"} elseif ($uacSettings -eq "NeverNotify") {
$level1 = "Enabled"
$level2 = "Disabled"
$level3 = "Disabled"
$level4 = "Disabled"}
try {

    $Comp=get-wmiobject win32_computersystem
    $computer = $comp.name ; $DomainNamee = $comp.domain
    $ExecutionLog += "ComputerName : $computer"
    $ExecutionLog += "Domain/Workgroup : $DomainNamee"
    $path = "HKLM:\Software\Microsoft\Windows\CurrentVersion\Policies\System"

   if($level1 -eq "Enabled")
        {
        try
        {
        Set-ItemProperty -Path $path -Name "FilterAdministratorToken" -Value 0
        Set-ItemProperty -Path $path -Name "EnableUIADesktopToggle" -Value 0
        Set-ItemProperty -Path $path -Name "ConsentPromptBehaviorAdmin" -Value 0
        Set-ItemProperty -Path $path -Name "ConsentPromptBehaviorUser" -Value 3
        Set-ItemProperty -Path $path -Name "EnableInstallerDetection" -Value 1
        Set-ItemProperty -Path $path -Name "ValidateAdminCodeSignatures" -Value 0
        Set-ItemProperty -Path $path -Name "EnableSecureUIAPaths" -Value 1
        Set-ItemProperty -Path $path -Name "EnableLUA" -Value 0
        Set-ItemProperty -Path $path -Name "PromptOnSecureDesktop" -Value 0
        Set-ItemProperty -Path $path -Name "EnableVirtualization" -Value 1
        $executionlog += "The UAC settings successfully changed to $uacSettings, please restart the computer in order to apply the changes on the system"
        Write-Output $executionlog
        exit;
        }
        Catch
        {
        $executionlog += "Not able to find the reuqired properties to change the UAC Settings on the system, Please contact the Administartor Team for the same."
        Write-Output $executionlog
        exit;
        }
    }
   elseif($level2 -eq "Enabled")
        {
        try
        {
        Set-ItemProperty -Path $path -Name "FilterAdministratorToken" -Value 0
        Set-ItemProperty -Path $path -Name "EnableUIADesktopToggle" -Value 0
        Set-ItemProperty -Path $path -Name "ConsentPromptBehaviorAdmin" -Value 5
        Set-ItemProperty -Path $path -Name "ConsentPromptBehaviorUser" -Value 3
        Set-ItemProperty -Path $path -Name "EnableInstallerDetection" -Value 1
        Set-ItemProperty -Path $path -Name "ValidateAdminCodeSignatures" -Value 0
        Set-ItemProperty -Path $path -Name "EnableSecureUIAPaths" -Value 1
        Set-ItemProperty -Path $path -Name "EnableLUA" -Value 1
        Set-ItemProperty -Path $path -Name "PromptOnSecureDesktop" -Value 0
        Set-ItemProperty -Path $path -Name "EnableVirtualization" -Value 1
        $executionlog += "The UAC settings successfully changed to $uacSettings, please restart the computer in order to apply the changes on the system"
        Write-Output $executionlog
        exit;
        }
        Catch
        {
        $executionlog += "Not able to find the reuqired properties to change the UAC Settings on the system, Please contact the Administartor Team for the same."
        Write-Output $executionlog
        exit;
        }
        }
   elseif($level3 -eq "Enabled")
        {
        try
        {
        Set-ItemProperty -Path $path -Name "FilterAdministratorToken" -Value 0
        Set-ItemProperty -Path $path -Name "EnableUIADesktopToggle" -Value 0
        Set-ItemProperty -Path $path -Name "ConsentPromptBehaviorAdmin" -Value 5
        Set-ItemProperty -Path $path -Name "ConsentPromptBehaviorUser" -Value 3
        Set-ItemProperty -Path $path -Name "EnableInstallerDetection" -Value 1
        Set-ItemProperty -Path $path -Name "ValidateAdminCodeSignatures" -Value 0
        Set-ItemProperty -Path $path -Name "EnableSecureUIAPaths" -Value 1
        Set-ItemProperty -Path $path -Name "EnableLUA" -Value 1
        Set-ItemProperty -Path $path -Name "PromptOnSecureDesktop" -Value 1
        Set-ItemProperty -Path $path -Name "EnableVirtualization" -Value 1
        $executionlog += "The UAC settings successfully changed to $uacSettings, please restart the computer in order to apply the changes on the system"
        Write-Output $executionlog
        exit;
        }
        Catch
        {
        $executionlog += "Not able to find the reuqired properties to change the UAC Settings on the system, Please contact the Administartor Team for the same."
        Write-Output $executionlog
        exit;
        }
        }
   elseif($level4 -eq "Enabled")
        {
        try
        {
        Set-ItemProperty -Path $path -Name "FilterAdministratorToken" -Value 0
        Set-ItemProperty -Path $path -Name "EnableUIADesktopToggle" -Value 0
        Set-ItemProperty -Path $path -Name "ConsentPromptBehaviorAdmin" -Value 2
        Set-ItemProperty -Path $path -Name "ConsentPromptBehaviorUser" -Value 3
        Set-ItemProperty -Path $path -Name "EnableInstallerDetection" -Value 1
        Set-ItemProperty -Path $path -Name "ValidateAdminCodeSignatures" -Value 0
        Set-ItemProperty -Path $path -Name "EnableSecureUIAPaths" -Value 1
        Set-ItemProperty -Path $path -Name "EnableLUA" -Value 1
        Set-ItemProperty -Path $path -Name "PromptOnSecureDesktop" -Value 1
        Set-ItemProperty -Path $path -Name "EnableVirtualization" -Value 1
        $executionlog += "The UAC settings successfully changed to $uacSettings, please restart the computer in order to apply the changes on the system"
        Write-Output $executionlog
        exit;
        }
        Catch
        {
        $executionlog += "Not able to find the reuqired properties to change the UAC Settings on the system, Please contact the Administartor Team for the same."
        Write-Output $executionlog
        exit;
        }
        }
}
Catch {
        $ExecutionLog += "Not able to reach the computer remotly through WMI"
        Write-Output $executionlog
              exit;}

