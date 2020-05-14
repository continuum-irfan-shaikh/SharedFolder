<#
Template Name : Disable Java update messages
Description : Turns off the auto updates for Java
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

$OS = [intPtr]::Size

if ( $OS -eq 8 ) {
   $regPath = "HKLM:\SOFTWARE\Wow6432Node\JavaSoft\Java Update\Policy"
}

if ( $OS -eq 4 ) {
   $regPath = "HKLM:\SOFTWARE\JavaSoft\Java Update\Policy"
}

if ([intPtr]::Size -eq 8 ) {
   
   try {
        $KeyName = Get-ItemProperty $RegPath -Name EnableJavaUpdate -ea stop
        
        if ( $KeyName.EnableJavaUpdate -ne 0 ) {
            
            Set-ItemProperty -Path $RegPath -Name EnableJavaUpdate -Value 0
            
            if ((Get-ItemProperty $RegPath -Name EnableJavaUpdate -ea stop).EnableJavaUpdate -eq 0 ) {
                 Write-Output "Java update disabled successfully" 
            }     
            
        } elseif ( $KeyName.EnableJavaUpdate -eq 0  ) {
          
            Write-output "Java update already disabled"
        }
        
                  
   } catch {
     
        Write-Error "Java not found on this machine..! : $($_.Exception.Message)"
       
   }

}
