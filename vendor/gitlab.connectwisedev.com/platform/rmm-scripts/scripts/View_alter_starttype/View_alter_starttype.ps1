####################################################
###------View or Alter Start Type---------##########
####################################################
 
 #$OperationType = "View"#"$null" #"Enable" # "Disable"
 #$StartTypePath = "" #(Optional)

##########################################################
###------Variable Declaration-----------------------------
##########################################################

$ResultObject = New-Object -TypeName psobject
$ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "View or Alter Start Type" 

$global:viewstartupArray = @()
$global:EnablestartupArray = @()
$global:DisablestartupArray = @()

$global:Code = 0
$global:ErrorMessageArray= @()
$global:SuccessMessageArray= @()



##########################################################
###------Checking Pre Condition---------------------------
##########################################################

Function Check-PreCondition{

    $IsContinued = $true

    Write-Host "-------------------------------"
    Write-Host "Checking Preconditions"
    Write-Host "" 
   

    #####################################
    # Verify PowerShell Version
    #####################################

    write-host -ForegroundColor 10 "`t PowerShell Version : $($PSVersionTable.PSVersion.Major)" 

    if(-NOT($PSVersionTable.PSVersion.Major -ge 2))
    {
        $global:Code = 2
        $global:ErrorMessageArray += "PowerShell version below 2.0 is not supported"

        $IsContinued = $false
    }

    
    ####################################     
    # Verify opearating system Version
    ####################################

    write-host -ForegroundColor 10 "`t Operating System Version : $([System.Environment]::OSVersion.Version.major)" 
    if(-not([System.Environment]::OSVersion.Version.major -ge 6))
    {
        $global:Code = 2
        $global:ErrorMessageArray += "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"

        $IsContinued = $false
    }
    
    
    Write-Host ""
    Write-Host -ForegroundColor 8 "`t Checking Precondition Completed"
    Write-Host ""
    Write-Host "-------------------------------"
    Write-Host ""

    return $IsContinued
}

##########################################################
###------View Startup Programs----------------------------
##########################################################
Function View-Startup{
    try{
    
        $startupitems = Get-WmiObject win32_startupcommand
        if($startupitems -ne $null){ 
            foreach($startupitem in $startupitems){
                $global:caption = $startupitem | select -ExpandProperty caption
                $global:command = $startupitem | select -ExpandProperty command
                $global:user =$startupitem | select -ExpandProperty user
                $global:status = "Enabled"
                
                $global:viewstartup = New-Object -TypeName psobject
                $global:viewstartup | Add-Member -MemberType NoteProperty -Name 'Name:' -Value $($global:caption)
                $global:viewstartup | Add-Member -MemberType NoteProperty -Name 'Path:' -Value $($global:command)
                $global:viewstartup | Add-Member -MemberType NoteProperty -Name 'User:' -Value $($global:user)
                $global:viewstartup | Add-Member -MemberType NoteProperty -Name 'Startup Type:' -Value $($global:status)
                $global:viewstartupArray += $global:viewstartup                   
                
                $global:SuccessMessageArray += 'Name:' + $($global:caption)
                $global:SuccessMessageArray += 'Path:' + $($global:command)
                $global:SuccessMessageArray += 'User:' + $($global:User)
                $global:SuccessMessageArray += 'Startup Type:' + $($global:status)
            }
        }
        if($startupitems -eq $null){
            $global:code = 2
            $global:ErrorMessageArray += "No Startup Items found"   
        }
    }
    
    catch{
        $global:code = 1
        $global:ErrorMessageArray += "Unable to fetch the startup items"   
    }
}

##########################################################
###------Enable Startup Programs--------------------------
##########################################################

function Enable-Startups {
    [CmdletBinding()]
    Param(
        #[parameter(DontShow = $true)]
        $32bit = "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Run",
        #[parameter(DontShow = $true)]
        $32bitRunOnce = "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnce",
        #[parameter(DontShow = $true)]
        $64bit = "HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Run",
        #[parameter(DontShow = $true)]
        $64bitRunOnce = "HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\RunOnce",
        #[parameter(DontShow = $true)]
        $currentLOU = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Run",
        #[parameter(DontShow = $true)]
        $currentLOURunOnce = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnce"
    )

    Begin {

        New-PSDrive -PSProvider Registry -Name HKU -Root HKEY_USERS | Out-Null
        $startups = Import-Csv "C:\Backup.csv"
    }
    
    Process {
        foreach ($startUp in $startUps){
            $number = ($startUp.Location).IndexOf("\")
            $location = ($startUp.Location).Insert("$number",":")
            $enableitem = Get-wmiobject Win32_StartupCommand | ?{$_.Name -match $startUp.name}
            if($enableitem -eq $null){
                try{                    
                    Set-ItemProperty -Path "$location" -Name $startUp.Name -Value $startup.command
                    $enable_start = Get-wmiobject Win32_StartupCommand | ?{$_.Name -match $startUp.name}
                    if($enable_start -eq $null){
                        $global:code = 1
                        $global:ErrorMessageArray += "Unable to enable the startup item : $($startUp.name)" 
                    }
                    if($enable_start -ne $null){
                        $global:name = $enable_start | select -ExpandProperty Name
                        $global:command = $enable_start | select -ExpandProperty Command
                        $global:username = $enable_start | select -ExpandProperty user                                              
                        $global:Enablestartup = New-Object -TypeName psobject
                        $global:Enablestartup | Add-Member -MemberType NoteProperty -Name 'Name:' -Value $($global:name)
                        $global:Enablestartup | Add-Member -MemberType NoteProperty -Name 'Path:' -Value $($global:command)
                        $global:Enablestartup | Add-Member -MemberType NoteProperty -Name 'User:' -Value $($global:username)
                        $global:Enablestartup | Add-Member -MemberType NoteProperty -Name 'Old Startup Type:' -Value "Disable"
                        $global:Enablestartup | Add-Member -MemberType NoteProperty -Name 'New Startup Type:' -Value "Enable"
                        $global:EnablestartupArray += $global:Enablestartup                  
            
                        $global:SuccessMessageArray += 'Name:' + $($global:name)
                        $global:SuccessMessageArray += 'Path:' + $($global:command)
                        $global:SuccessMessageArray += 'User:' + $($global:username)
                        $global:SuccessMessageArray += 'Old Startup Type:' + "Disable"
                        $global:SuccessMessageArray += 'New Startup Type:' + "Enable"
                    }                    
                }
                catch{
                    $global:code = 1
                    $global:ErrorMessageArray += "Unable to enable the startup items" 
                }
            }
            if($enableitem -ne $null){
               #Write-Host "$($startUp.name) is already in Enable state" 
               $global:code = 1
               $sname = $enableitem | select -ExpandProperty Name
               $global:ErrorMessageArray += "$($sname) is already in Enable state" 
            }
        }      
    }
}

##########################################################
###------Disable Startup Programs-------------------------
##########################################################

function Disable-Startups {
    [CmdletBinding()]
    Param(
        #[parameter(DontShow = $true)]
        $32bit = "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Run",
        #[parameter(DontShow = $true)]
        $32bitRunOnce = "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnce",
        #[parameter(DontShow = $true)]
        $64bit = "HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Run",
        #[parameter(DontShow = $true)]
        $64bitRunOnce = "HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\RunOnce",
        #[parameter(DontShow = $true)]
        $currentLOU = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Run",
        #[parameter(DontShow = $true)]
        $currentLOURunOnce = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnce"
    )

    Begin {

        New-PSDrive -PSProvider Registry -Name HKU -Root HKEY_USERS | Out-Null
        Copy-Item -Path Backup.csv -Destination C:\Backup.csv
        $startups = Get-wmiobject Win32_StartupCommand | Select-Object Name,Location
        Get-wmiobject Win32_StartupCommand | Select Name,Caption,User,Command,Location |
        Export-Csv Backup.csv -NoTypeInformation
    }
    
    Process {
        foreach ($startUp in $startUps){
            $number = ($startUp.Location).IndexOf("\")
            $location = ($startUp.Location).Insert("$number",":")
            $disableitem = Get-wmiobject Win32_StartupCommand | ?{$_.Name -match $startUp.name}
            if($disableitem -ne $null){                
                try{
                    #Write-Output "Disabling $($startUp.Name) from $location)"
                    Remove-ItemProperty -Path "$location" -Name "$($startUp.name)"
                    $disable_start = Get-wmiobject Win32_StartupCommand | ?{$_.Name -match $startUp.name}
                    if($disable_start -eq $null){
                        $global:name = $disableitem | select -ExpandProperty Name
                        $global:command = $disableitem | select -ExpandProperty Command
                        $global:username = $disableitem | select -ExpandProperty user 
                        $global:Disablestartup = New-Object -TypeName psobject
                        $global:Disablestartup | Add-Member -MemberType NoteProperty -Name 'Name:' -Value $($global:name)
                        $global:Disablestartup | Add-Member -MemberType NoteProperty -Name 'Path:' -Value $($global:command)
                        $global:Disablestartup | Add-Member -MemberType NoteProperty -Name 'User:' -Value $($global:username)
                        $global:Disablestartup | Add-Member -MemberType NoteProperty -Name 'Old Startup Type:' -Value "Enable"
                        $global:Disablestartup | Add-Member -MemberType NoteProperty -Name 'New Startup Type:' -Value "Disable"
                        $global:DisablestartupArray += $global:Disablestartup                 
            
                        $global:SuccessMessageArray += 'Name:' + $($global:name)
                        $global:SuccessMessageArray += 'Path:' + $($global:command)
                        $global:SuccessMessageArray += 'User:' + $($global:username)
                        $global:SuccessMessageArray += 'Old Startup Type:' + "Enable"
                        $global:SuccessMessageArray += 'New Startup Type:' + "Disable"
                    }
                    if($disable_start -ne $null){                       
                        #Write-Host "$($startUp.name) is not Disabled"
                        $global:code = 1
                        $global:ErrorMessageArray += "$($startUp.name) is not Disabled"
                    }
                }
                catch{
                    $global:code = 1
                    $global:ErrorMessageArray += "Unable to disable the startup items" 
                }
            }
            if($disableitem -eq $null){
               #Write-Host "$($startUp.name) is already in Disable state" 
               $global:code = 1  
               $dname = Import-csv C:\Backup.csv                         
               $global:ErrorMessageArray += "$($dname.Name) is already disabled" 
            }
        }      
    } 
}

##########################################################
###------Enable Startup Programs with Path----------------
##########################################################

function Enable-StartupTypePath {
    [CmdletBinding()]
    Param(
        #[parameter(DontShow = $true)]
        $32bit = "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Run",
        #[parameter(DontShow = $true)]
        $32bitRunOnce = "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnce",
        #[parameter(DontShow = $true)]
        $64bit = "HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Run",
        #[parameter(DontShow = $true)]
        $64bitRunOnce = "HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\RunOnce",
        #[parameter(DontShow = $true)]
        $currentLOU = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Run",
        #[parameter(DontShow = $true)]
        $currentLOURunOnce = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnce"
    )

    Begin {

        New-PSDrive -PSProvider Registry -Name HKU -Root HKEY_USERS | Out-Null
        $startups = Import-Csv "Backup.csv" 
    }
    
    Process {
        foreach ($startUp in $startUps){
            $number = ($startUp.Location).IndexOf("\")
            $location = ($startUp.Location).Insert("$number",":")
            $enableitem = Get-wmiobject Win32_StartupCommand | ?{$_.Name -match $startUp.name}
            if($enableitem -eq $null){
                try{
                    Write-Output "Enabling $($startUp.Name) to $location)"                
                    Set-ItemProperty -Path "$location" -Name $startUp.Name -Value $StartTypePath -Force
                    $enable_start = Get-wmiobject Win32_StartupCommand | ?{$_.Name -match $startUp.name}
                    if($enable_start -ne $null){                        
                        $global:Enablestartup = New-Object -TypeName psobject
                        $global:Enablestartup | Add-Member -MemberType NoteProperty -Name 'Name:' -Value $($enableitem.Name)
                        $global:Enablestartup | Add-Member -MemberType NoteProperty -Name 'Path:' -Value $($enableitem.command)
                        $global:Enablestartup | Add-Member -MemberType NoteProperty -Name 'User:' -Value $($enableitem.user)
                        $global:Enablestartup | Add-Member -MemberType NoteProperty -Name 'Old Startup Type:' -Value "Disable"
                        $global:Enablestartup | Add-Member -MemberType NoteProperty -Name 'New Startup Type:' -Value "Enable"
                        $global:EnablestartupArray += $global:Enablestartup                  
            
                        $global:SuccessMessageArray += 'Name:' + $($enableitem.Name)
                        $global:SuccessMessageArray += 'Path:' + $($enableitem.command)
                        $global:SuccessMessageArray += 'User:' + $($enableitem.user)
                        $global:SuccessMessageArray += 'Old Startup Type:' + "Disable"
                        $global:SuccessMessageArray += 'New Startup Type:' + "Enable"
                    }
                    if($enable_start -eq $null){
                        Write-Host "$($startUp.name) is not Enabled"
                        $global:code = 1
                        $global:ErrorMessageArray += "Unable to enable the startup item : $($startUp.name)" 
                    }
                }
                catch{
                    $global:code = 1
                    $global:ErrorMessageArray += "Unable to enable the startup items" 
                }
            }
            if($enableitem -ne $null){
               Write-Host "$($startUp.name) is already in Enable state" 
               $global:code = 1
               $global:ErrorMessageArray += "$($startUp.name) is already in Enable state" 
            }
        }      
    }
}

##########################################################
###------JSON for PS2 ------------------------------------
##########################################################

function Escape-JSONString($str)
{
	if ($str -eq $null) 
    {
        return ""
    }

	$str = $str.ToString().Replace('"','\"').Replace('\','\\').Replace("`n",'\n').Replace("`r",'\r').Replace("`t",'\t')

	return $str;
}


function ConvertTo-JSONPS2($maxDepth = 4,$forceArray = $false) 
{
	begin {
		$data = @()
	}
	process{
		$data += $_
	}
	
	end{
	
		if ($data.length -eq 1 -and $forceArray -eq $false) {
			$value = $data[0]
		} else {	
			$value = $data
		}

		if ($value -eq $null) {
			return "null"
		}

		

		$dataType = $value.GetType().Name
		
		switch -regex ($dataType) {
	            'String'  {
					return  "`"{0}`"" -f (Escape-JSONString $value )
				}

	            '(System\.)?DateTime'  {return  "`"{0:yyyy-MM-dd}T{0:HH:mm:ss}`"" -f $value}

	            'Int32|Double' {return  "$value"}

				'Boolean' {return  "$value".ToLower()}

	            '(System\.)?Object\[\]' { # array
					
					if ($maxDepth -le 0){return "`"$value`""}
					
					$jsonResult = ''
					foreach($elem in $value){
						#if ($elem -eq $null) {continue}
						if ($jsonResult.Length -gt 0) {$jsonResult +=', '}				
						$jsonResult += ($elem | ConvertTo-JSONPS2 -maxDepth ($maxDepth -1))
					}
					return "[" + $jsonResult + "]"
	            }

				'(System\.)?Hashtable' { # hashtable
					$jsonResult = ''
					foreach($key in $value.Keys){
						if ($jsonResult.Length -gt 0) {$jsonResult +=', '}
						$jsonResult += 
@"
	"{0}": {1}
"@ -f $key , ($value[$key] | ConvertTo-JSONPS2 -maxDepth ($maxDepth -1) )
					}
					return "{" + $jsonResult + "}"
				}

	            default { #object
					if ($maxDepth -le 0){return  "`"{0}`"" -f (Escape-JSONString $value)}
					
					return "{" +
						(($value | Get-Member -MemberType *property | % { 
@"
	"{0}": {1}
"@ -f $_.Name , ($value.($_.Name) | ConvertTo-JSONPS2 -maxDepth ($maxDepth -1) )			
					
					}) -join ', ') + "}"
	    		}
		}
	}
}











##########################################################
###------Set Result --------------------------------------
##########################################################

function SetResult 
{
    $ResultObject | Add-Member -MemberType NoteProperty -Name 'Code' -Value $global:Code

    if(($global:Code -eq 0) -or ($global:Code -eq 1))
    {
        $statusMSG = ""
        if($global:Code -eq 0)
        {
            $successString = "Success: " + "The operation completed successfully"
            $statusMSG = "success"
        }
        else
        {
            $successString = "Fail: " + "The operation faild with error"
            $statusMSG = "fail"
        }

        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value $statusMSG
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value $successString
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $global:SuccessMessageArray
 
        
        $OutputObject = New-Object -TypeName psobject
        if($OperationType -eq "View"){
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'View Startups' -Value $global:viewstartupArray
        } 
    
        if($OperationType -eq "Enable"){
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'Enable Startups' -Value $global:EnablestartupArray
        }
         
        if($OperationType -eq "Disable"){
            $OutputObject | Add-Member -MemberType NoteProperty -Name 'Disable Startups' -Value $global:DisablestartupArray
        }     
        
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'dataObject' -Value $OutputObject


        #------------Error with Success--------------------------------------------
        $errorCount= 0
        $ErrorObjectAray= @()

        foreach($ErrorMessage in $global:ErrorMessageArray)
        {
           
            $ErrObject = New-Object -TypeName psobject
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'id' -Value $errorCount
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'title' -Value $ErrorMessage
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'detail' -Value $ErrorMessage

            $ErrorObjectAray += $ErrObject

            $errorCount = $errorCount +1
        }

        if($ErrorObjectAray.Count -gt 0)
        {
            $ResultObject | Add-Member -MemberType NoteProperty -Name 'stderr' -Value $ErrorObjectAray
        }

    }
    else
    {
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value "fail"

        $ErrorObjectAray= @()
        $errorCount= 0
        $errorString = "Error: "

        foreach($ErrorMessage in $global:ErrorMessageArray)
        {
            $errorString += $ErrorMessage + ", "
            
            $ErrObject = New-Object -TypeName psobject
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'id' -Value $errorCount
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'title' -Value $ErrorMessage
            $ErrObject | Add-Member -MemberType NoteProperty -Name 'detail' -Value $ErrorMessage

            $ErrorObjectAray += $ErrObject

            $errorCount = $errorCount +1
        }
        
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value $errorString
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $global:ErrorMessageArray
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stderr' -Value $ErrorObjectAray
    }
}


##########################################################
###------Display Result ----------------------------------
##########################################################

function DisplayResult 
{

    if($PSVersionTable.PSVersion.Major -gt 2)
    {
        $JSONResult= $ResultObject|ConvertTo-Json -Depth 6
        $JSONResult
    }
    else
    {
        $JSONResult= $ResultObject|ConvertTo-JSONPS2 -maxDepth 6
        $JSONResult
    }
}



##########################################################
###------Execute Functions -------------------------------
##########################################################
cls

if(Check-PreCondition)
{
    if($OperationType -eq "View"){
        View-Startup -ErrorAction SilentlyContinue
    } 
    
    if($OperationType -eq "Enable"){
        Enable-Startups -ErrorAction SilentlyContinue
    }
     
    if($OperationType -eq "Disable"){
        Disable-Startups -ErrorAction SilentlyContinue
    }
    if($OperationType -eq $null -or $OperationType -eq ""){
        $global:code = 2
        $global:ErrorMessageArray += "Any of the Operation Type (View/Enable/Disable) is mandatory"
    } 
<#
    if($OperationType -eq "Enable" -and $StartTypePath -ne $null){
        Enable-StartupTypePath -ErrorAction SilentlyContinue
    } 
 #>       
    SetResult
    DisplayResult 
}
else
{
    SetResult
    DisplayResult
}
