 <#
    .Synopsis
     Retrive List of Users whose passwords will expire shortly  
    .Description 
     To get the details of the users whose will get expired soon on the machine
    .Help
     To get more information about the users whose password will be going to expired soon then refer below details
     Get-WmiObject -Class Win32_UserAccount
    .Author
     sushma.yerasi@continuum.net
    .Name 
     Retrive List of Users whose passwords will expire shortly
#>

$ExecutionLog=@()
$Obj = @()
#$days = "5"
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

    if (4,5 -contains $Comp.DomainRole){	
			 Import-Module ActiveDirectory
			 $users = Get-ADUser -filter {Enabled -eq $True -and PasswordNeverExpires -eq $False}  –Properties "SamAccountName","msDS-UserPasswordExpiryTimeComputed" | Select-Object  @{Name="Account Name";Expression={$_."SamAccountName"}}, @{Name="Password Expiry Date";Expression={[datetime]::FromFileTime($_."msDS-UserPasswordExpiryTimeComputed")}} | Where-Object {$_.'Password expiry date' -lt (Get-Date).AddDays($days).ToString("d")} | fl
			 if ($users -ne $null) {
                         $ExecutionLog += $users
			 Write-Output $ExecutionLog
                         Exit; }
			else {    
				$ExecutionLog += "There are no users on the Machine whose passowrd will be going to be expired soon"
				Write-Output $executionlog
			        exit;}      }
    else {	
			$now = Get-Date
			$AllLocalAccounts = Get-WmiObject -Class Win32_UserAccount -filter {LocalAccount = "True" and disabled = "False" } 
			$Obj = $AllLocalAccounts | ForEach-Object {
                        $user = ([adsi]"WinNT://$computer/$($_.Name),user")
		        $pwAge    = $user.PasswordAge.Value
		        $maxPwAge = $user.MaxPasswordAge.Value
		        New-Object -TypeName PSObject -Property @{
		       'Account Name'         = $_.Name
		       'Password Expires'     = $_.PasswordExpires
		       'Password Last Set'    = $pwLastSet
		       'Password Expiry Date' = $now.AddSeconds($maxPwAge - $pwAge)
		          }
		         }
		       $Ouput = $Obj | select "Account Name", "Password Expiry Date", "Password Expires" | Where-Object {$_.'Password Expiry Date' -lt (Get-Date).AddDays($days).ToString("d") } | Where-Object { $_.'Password Expires' -eq $True}		
		       if ($Ouput -ne $null) {
                $ExecutionLog += $Ouput | select "Account Name", "Password Expiry Date" | fl
	       	       Write-Output $ExecutionLog
                       Exit;}
		       else {
				$ExecutionLog += "There are no users on the Machine whose passowrd will be going to be expired soon"
				Write-Output $executionlog
			        exit;}
		       	       }
}
Catch { 
	$ExecutionLog += "Not able to reach the computer remotly through WMI"
	Write-Output $executionlog
              exit;}
