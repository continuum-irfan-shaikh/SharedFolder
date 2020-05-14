<#
  .Script
  Disable always show for icons and notification.
  .Author
  Nirav Sachora.
  .Version
  2.0 and above.
  .Requirements
  Script should run with highest privileges.
#>

$users = @()

Function Create_registry($path)    # Function will set registry value and if item is not present will create registry entry.
{
$Details = Get-ItemProperty -Path $path
if($Details -ne $null)
    	{
    $Details = Get-ItemProperty -Path $path | gm | select -ExpandProperty Name
    if($Details -contains 'EnableAutoTray')
                        {
                        Set-ItemProperty $path -Name 'EnableAutoTray' -Value 1
                        }
                        else
                        {
                        New-ItemProperty -Path $path -Name 'EnableAutoTray' -PropertyType "DWord" -Value 1 | Out-Null
                        }
	
	} #End if statement
else
	{
	New-ItemProperty -Path $path -Name 'EnableAutoTray' -PropertyType "DWord" -Value 1 | Out-Null
	}  # End else statement

	if((Get-ItemProperty $path -Name 'EnableAutoTray'  | Select -ExpandProperty EnableAutoTray) -eq 1){return $true} else {return $false}
}  # End Function Create_registry

#registry edit for currently logged uin users
try
{
$profilelist = Get-Item "REgistry::HKU\S-1-5-21-*" | Where-Object {$_.Name -notlike '*_Classes'} | select -ExpandProperty name
foreach($profile in $profilelist)
    {
    if(Test-path "Registry::$profile\Volatile Environment")
    {
    $username = Get-ItemProperty -Path "Registry::$profile\Volatile Environment" -Name Username | Select -ExpandProperty Username

    $test = "Registry::$profile\Software\Microsoft\Windows\CurrentVersion\Explorer"
    if($test)
        {
        $cu = Create_registry -path "Registry::$profile\Software\Microsoft\Windows\CurrentVersion\Explorer"
        if($cu -eq $false)
           {
           Write-Error "Error $profile"
           }
        else
            {
        $Users += $username
        }
        }
    }
    else
    {
    Continue
    }
    }
}
catch
{
    $_.Exception.Message
}

# Registry edit for default profiles
    
    REG LOAD HKU\TEMP "C:\Users\default\NTUSER.DAT" | Out-Null
    if($?)
    {
        $defaultpath = "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Explorer"
        $defaultresult = Create_registry -path $defaultpath
            if($defaultresult -eq $false)
            {
                Write-Error "Error Default"
            }                      
       REG UNLOAD HKU\TEMP | Out-Null
    }
    else
    {
        Write-Error "Failed to load registry for default profile"
    }

#Registry edit for current profiles

$profiles = get-childitem C:\Users | select -ExpandProperty Name
foreach($profile in $profiles)
{
    $ntuser = Test-Path C:\Users\$profile\NTUSER.DAT
    if($ntuser)
        {
        try
            {
                    $FileStream = [System.IO.File]::Open("C:\Users\$profile\NTUSER.DAT",'Open','Write')
                    $FileStream.Close()
                    $FileStream.Dispose()
                    REG LOAD HKU\TEMP "C:\Users\$profile\NTUSER.DAT" | Out-Null
                    $path = "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Explorer"
                    $registry = Test-path $path
                    if($registry)
                    {
                        $result = Create_registry -path $path
                        if($result -eq $false)
                        {
                        Write-Output "Error $profile"
                        }
                        else
                        {
                        $Users += $profile
                        }
                    } #Closing IF Line 23

                    REG UNLOAD HKU\TEMP | Out-Null
            } #Closing Try
        catch
           {
                    Continue
           }
            
        } #Closing IF Statement Line 12  

} # Closing Foreach Loop
Write-Output "Always show for icons and notifications has been disabled for below users on this system`n"
Write-Output $Users
Write-Output "`nPlease logoff and login, if changes are not reflected."
