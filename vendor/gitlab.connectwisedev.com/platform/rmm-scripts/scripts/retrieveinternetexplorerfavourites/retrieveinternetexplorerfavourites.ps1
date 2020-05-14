#$username = "Dipak"  parameter
$ExecutionLog = @()
[Bool]$assigneduser = $False
$ouput = @()
$domainnames = @()
$ouput = @()
$ouputs = @()
$names = @()
$global:ProgressPreference = 'SilentlyContinue'  
[double]$OSVersion = [Environment]::OSVersion.Version.ToString(2)
# OS Version and Powershell comparison.
if (($osversion -lt '6.1') -or ($PSVersionTable.PSVersion.Major -lt '2')) {
    $executionlog += 'Prerequisites to run the script is not valid, Hence Script Exceution stopped' #'Script is design for windows 7 and Above Members only, Script Execution Stopped.'
    Write-Output $executionlog
    #exit;
} 
#To check the status of the computer
if ([Bool]$assigneduser -eq $True) { $CurrentUsers = $username }
else { 
    $AllUsers = query user /server:$computer 2>&1
    $Users = $AllUsers | ForEach-Object {(($_.trim() -replace ">" -replace "(?m)^([A-Za-z0-9]{3,})\s+(\d{1,2}\s+\w+)", '$1  none  $2' -replace "\s{2,}", "," -replace "none", $null))} | ConvertFrom-Csv
    $CurrentUsers = @()
    ForEach ($Userr in $Users)
    {
    $CUser = ($userr | ?{$_.state -ne 'Disc'} | Select-Object username).username
    $CurrentUsers+= $CUser
    }
 }
try {
    $Comp = get-wmiobject win32_computersystem
    $computer = $comp.name ; $DomainNamee = $comp.domain
    $ExecutionLog += "ComputerName : $computer"
    $ExecutionLog += "Domain/Workgroup : $DomainNamee" 
    foreach ($user in $CurrentUsers)
    {
  
    if (4, 5 -contains $Comp.DomainRole) {
        Import-Module ActiveDirectory
        $domainnames = Get-ADUser -filter {Enabled -eq $True} | Where-object {$_.Name -like "*$user*"}
        if ($domainnames -ne $null) {
            foreach ($domainname in $domainnames) {
                $username = $domainname.Name
                $path = "C:\users\$username\Favorites"
                $ouputs += Get-ChildItem $path -Recurse  | select Name, @{Name = "Link"; Expression = {($_ | Select-String "^URL" ).Line.Trim("URL=")}}
                $ouput += $ouputs | ? {$_.Name -notlike "*Links*" }
                
                if ($ouput.count -ne 0) {
                    $executionlog += "Domain User : $username"   
                    $executionlog += "The Favorite sites for the user $username on the Internet Explorer are mentioned below" 
                    foreach ($output in $ouput) {
                        if ($output.link -ne $null) {
                            $executionlog += "`n"      
                            $executionlog += ($output | select Name, Link | fl | Out-String).Trim()
                        }
                    }
                }
                else { 
                    $executionlog += "There are no favorite sites in the Internet Explorer for the user $username the Machine"
                    $ExecutionLog += "`n"
                } 
            } 
        }
        else { 
            $executionlog += "We cannot able to find the user name on the Machine, Please check the user name and try again"
          
        }

        
    }
    else {
        $names = Get-WmiObject -Class Win32_UserAccount -filter {disabled = "False" } | Where-object {$_.Name -like "*$user*"}
        if ($names -ne $null) {
            foreach ($name in $names) {
                $username = $name.Name
                $ExecutionLog += "Loacl User : $username" 
                $path = "C:\users\$username\Favorites"
                $ouputs = Get-ChildItem $path -Recurse  | select Name, @{Name = "Link"; Expression = {($_ | Select-String "^URL" ).Line.Trim("URL=")}}
                $ouput = $ouputs| ? {$_.Name -notlike "*Links*" }
                
                if ($ouput.count -ne 0) {
                    $ExecutionLog += "The Favorite sites for the user $username on the Internet Explorer are mentioned below"
                    foreach ($output in $ouput) {
                        if ($output.link -ne $null) {
                         
                            #$ExecutionLog += "`n"   
                            $executionlog += "`n" + ($output | select Name, Link | fl | Out-String).Trim()
                            $ExecutionLog += "`n"
                        }
                    }

                }
                else { 
                    $executionlog += "There are no favorite sites in the Internet Explorer for the user $username the Machine"
                    $ExecutionLog += "`n"
                }
            }

          
        }
              
        else { 
            $executionlog += "We cannot able to find the user name on the Machine, Please check the user name and try again"
           
        }
    }
    }
     Write-output $executionlog 
            #exit;
}
Catch { 
    $ExecutionLog += "Not able to reach the computer remotelt through WMI"
    Write-Output $executionlog
    #exit;
}
       

