 <#
 .Synopsis
 Script will hide continuum agent from Controlpanel\Add remove programs
 .Author
 Nirav Sachora
 .Requirements
 Script Should run with admin privileges.
 Sript Name: "Hide the Continuum Agent from the installed programs list"
#>



function verification($path)
{
    $Value = Get-ItemProperty $path -Name 'SystemComponent'  | Select -ExpandProperty SystemComponent
    if($Value -eq 1)
    {
        return $true
    }
    else
    {
        return $false
    }
}

function set-registry($path)
    {
        $ITSagent = Get-ItemProperty $path | gm | select -ExpandProperty Name
        if($ITSagent -Contains 'SystemComponent')
        {
        Set-ItemProperty $path -Name 'SystemComponent' -Value 1
	return $true
        }
         else
        {
	try
	{
        New-ItemProperty -Path $path -Name 'SystemComponent' -PropertyType "DWord" -Value 1 | Out-Null
	return $true
        }
	catch
	{
	return $false
	}
	}
    }

$osarchitecture = get-wmiobject -Class win32_operatingsystem | select -ExpandProperty Osarchitecture
if($osarchitecture -eq "32-bit")
{
$path = "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\SAAZOD"
try
{
$ITSPlatform = (Get-item -path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall").GetSubKeyNames() | Where-Object {$_ -like "{*}"} | ForEach-Object {Get-ItemProperty "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\$_"} | Where-Object {($_.DisplayName-eq "ITSPlatform") -and (Get-Member -inputobject $_ -name "InstallDate")} | Select -ExpandProperty PSChildName
}
catch
{
write-output "ITSPLatform is not installed"
}
$path1 = "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\$ITSPlatform"
}
else
{

$path = "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\SAAZOD"
try
{
$ITSPlatform = (Get-item -path "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall").GetSubKeyNames() | Where-Object {$_ -like "{*}"} | ForEach-Object {Get-ItemProperty "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\$_"} | Where-Object {($_.DisplayName-eq "ITSPlatform") -and (Get-Member -inputobject $_ -name "InstallDate")} | Select -ExpandProperty PSChildName
}
catch
{
write-output "ITSPLatform is not installed"
}
$path1 = "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\$ITSPlatform"

}

$ITSupport247 = Test-path $path
$ITPlatform = Test-path $path1

 if($ITSupport247)
    {
        $setreg = set-registry -path $path
        if($setreg){$verify = verification -path $path} else {Write-Error "could not create registry entry ITSupport247"}
            if($verify)
            {
            Write-Output "ITSupport247 has been hidden from the system"
            }
    }
   else
        {
        Write-Output "ITSupport247 is not installed on this system"
        }
  
  if($ITPlatform)
    {
        $setreg = set-registry -path $path1
        if($setreg){$verify = verification -path $path1} else {Write-Error "could not create registry entry ITPlatform"}
            if($verify)
            {
            Write-Output "ITPlatform has been hidden from the system"
            }
    }
   else
        {
        Write-Output "ITPlatform is not installed on this system"
        }
