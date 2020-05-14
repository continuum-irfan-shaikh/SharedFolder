<#

Name : Enable remote administration
Category : Security

    .SYNOPSIS
        Enable Remote Administration
    .DESCRIPTION
        Enable Remote Administration. By enabling the remote administration, we can access the system remotely
    .Help
        HKey 'HKLM:\SYSTEM\CurrentControlSet\Control\Terminal Server'
        HKey Property 'f
        DenyTSConnections'
        0 = Enabled
        1 = Disabled
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>
Function Create_registry($path, $key, $value)    
{
$Details = Get-ItemProperty -Path $path
if($Details -ne $null)
                {
    $Details = Get-ItemProperty -Path $path | gm | select -ExpandProperty Name
    if($Details -contains $key)
                        {
                        Set-ItemProperty $path -Name $key -Value $value
                        }
                        else
                        {
                        New-ItemProperty -Path $path -Name $key -PropertyType "DWord" -Value $value | Out-Null
                        }
                } 
                else
                {
                New-ItemProperty -Path $path -Name $key -PropertyType "DWord" -Value $value | Out-Null
                }  
        if((Get-ItemProperty $path -Name $key  | Select -ExpandProperty $key) -eq $value){return $true} else {return $false}
}  

$path = 'HKLM:\SYSTEM\CurrentControlSet\Control\Terminal Server'
$name = 'fDenyTSConnections'
$value = 0
    
       if (Create_registry -path $path -key $name -value $value) {
        Write-Output "`nRemote Administration enabled on this system $env:COMPUTERNAME"
        }
        else
        {
        Write-Error $_.Exception.Message
        Write-Error "`nSomething went wrong while updating the Registry. Kindly check manually"
        }
