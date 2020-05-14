if ((Get-WMIObject win32_operatingsystem).name -like "*Server*" ){
  Write-Error "Not supported on Windows Server operating system!"
  Exit
}

$fw=New-object -comObject HNetCfg.FwPolicy2    
if (($fw.IsRuleGroupCurrentlyEnabled("Network Discovery")) -eq 1 ){
     Write-Output "Network Discovery already enabled on this computer!"
}Else{
    $output = netsh advfirewall firewall set rule group="Network Discovery" new enable=yes
    if ($LASTEXITCODE -eq 0){
         if (($fw.IsRuleGroupCurrentlyEnabled("Network Discovery")) -eq 1 ){
              Write-output "Network Discovery successfully enabled."
         } Else { Write-Error "Failed to enable Network Discovery. The dependency services for Network Discovery might not be running."}
    }else{ Write-Error "Error occured while enabling Network Discovery. $output" }
}
