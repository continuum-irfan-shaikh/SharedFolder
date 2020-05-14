<#
    .Synopsis
     Retrive UserName and Domain  
    .Description 
     To get the details of the users who are logged in to the server along with their domain name and profile type
    .Help
     To get more information about the userdetails, please refer below details
     Get-WmiObject -Class Win32_UserProfile
    .Author
     sushma.yerasi@continuum.net
    .Name 
     Retrive UserName and Domain 
#>

$ExecutionLog=@()
[double]$OSVersion=[Environment]::OSVersion.Version.ToString(2)
# OS Version and Powershell comparison.
if (($osversion -lt '6.1') -or ($PSVersionTable.PSVersion.Major -lt '2'))
        {
            $executionlog += 'Prerequisites to run the script is not valid, Hence Script Exceution stopped' #'Script is design for windows 7 and Above Members only, Script Execution Stopped.'
            Write-Output $executionlog
            exit;
        } 

#To check the role of the computer
try {

    $Comp=get-wmiobject win32_computersystem
    $computer = $comp.name ; $DomainNamee = $comp.domain
    $ExecutionLog += "ComputerName : $computer"
    $ExecutionLog += "Domain/Workgroup : $DomainNamee"
    $ExecutionLog += "Currently Loggedin user names are mentioned below"    
    $AllUsers = query user /server:$computer 2>&1
    $Users = $AllUsers | ForEach-Object {(($_.trim() -replace ">" -replace "(?m)^([A-Za-z0-9]{3,})\s+(\d{1,2}\s+\w+)", '$1  none  $2' -replace "\s{2,}", "," -replace "none", $null))} | ConvertFrom-Csv
    $CurrentUsers = @()
    ForEach ($User in $Users)
    {
    $CUser = ($user | ?{$_.state -ne 'Disc'} | Select-Object username).username
    $CurrentUsers+= $CUser
    }
    
    $ExecutionLog += "$CurrentUsers"
# Commands to check the Local User Profiles

  $Profiles = Get-WmiObject -Class Win32_UserProfile -Computer $Computer -ea 0
  $ExecutionLog += "All user profiles are mentioned below" 
  foreach ($profile in $profiles) {
  try {
      $objSID = New-Object System.Security.Principal.SecurityIdentifier($profile.sid)
      $objuser = $objsid.Translate([System.Security.Principal.NTAccount])
      $objusername = $objuser.value
  } catch {
        $objusername = $profile.sid
  }
  switch($profile.status){
   1 { $profileType="Temporary" }
   2 { $profileType="Roaming" }
   4 { $profileType="Mandatory" }
   8 { $profileType="Corrupted" }
   default { $profileType = "LOCAL" }
  }
  $User = $objUser.Value | ?{$_ -notlike "*All*" -and $_ -notlike "*Default*" -and $_ -notlike "*public*" -and $_ -notlike "*Classic*"}   
  if (($User -notlike '*NT Authority*') -and ($User -notlike $null)) { 
  $ExecutionLog += " UserName: $user"
  $ExecutionLog += " ProfileType: $ProfileType"
   }
 }  
   Write-Output $executionlog
              exit;
}
Catch { 
	$ExecutionLog += "Not able to reach the computer remotelt through WMI"
	Write-Output $executionlog
              exit;
 }
