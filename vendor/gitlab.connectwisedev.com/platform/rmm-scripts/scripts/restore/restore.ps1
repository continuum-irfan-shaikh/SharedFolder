$OS = (Get-WMIObject -Class Win32_OperatingSystem).Caption
If ($OS -like "*Server*") {
            Write-Error "Current OS : $OS. This functionality is not supported on this operating system."
            Exit
} Else {
    Try     { 
            Enable-ComputerRestore -Drive (Get-WmiObject Win32_OperatingSystem).SystemDrive
 IF ($?){ 
                Write-Output "System restore enabled"
        }
            }
Catch {
           Write-Error $_.Exception.Message
           Exit 
      }
}
