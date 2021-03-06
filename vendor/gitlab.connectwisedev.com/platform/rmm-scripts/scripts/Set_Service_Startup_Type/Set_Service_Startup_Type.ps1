
#######################################################################
###------------Set Service Startup Type----------------------------####
#######################################################################
 
    [CmdletBinding()]
    Param( 
        #[Parameter(Position=0, Mandatory=$true, ParameterSetName='StartType')]
        [ValidateSet(“Automatic”,”Delayed”,”Manual”,"Disabled")] 
        [String] 
        $starttype
    , 
        #[Parameter(Position=1, Mandatory=$true, ParameterSetName='Service')] 
        [string] 
        $servicenames

    )



##########################################################
###------Variable Declaration-----------------------------
##########################################################

    $Computer = $env:computername
    $ResultObject = New-Object -TypeName psobject
    $ResultObject | Add-Member -MemberType NoteProperty -Name 'Task Name' -Value "Set Service Startup Type" 
    $global:starttype = $null
    $global:servicenames = $null
    $global:starttype = $starttype
    $global:servicenames = $servicenames

    $global:Code = 0
    $global:ErrorMessageArray= @()
    $global:SuccessMessageArray= @()
    $global:delay = $null
    $global:startmode = $null

##########################################################
###------Checking Pre Condition---------------------------
##########################################################


    Function Check-PreCondition{

        $IsContinued = $true

        Write-Host "-------------------------------"
        Write-Host "Checking Preconditions"
        Write-Host "-------------------------------" 
       

        #####################################
        # Verify PowerShell Version
        #####################################

        write-host -ForegroundColor 10 "`t PowerShell Version : $($PSVersionTable.PSVersion.Major)" 

        if(-NOT($PSVersionTable.PSVersion.Major -ge 2))
        {
            $global:Code = 1
            $global:ErrorMessageArray += "PowerShell version below 2.0 is not supported"

            $IsContinued = $false
        }

        
        ####################################     
        # Verify opearating system Version
        ####################################

        write-host -ForegroundColor 10 "`t Operating System Version : $([System.Environment]::OSVersion.Version.major)" 
        if(-not([System.Environment]::OSVersion.Version.major -ge 6))
        {
            $global:Code = 1
            $global:ErrorMessageArray += "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"

            $IsContinued = $false
        }
        
        
        Write-Host ""
        Write-Host -ForegroundColor 8 "`t Checking Precondition Completed"
        Write-Host ""

        return $IsContinued
    }



    ####################################################
    ###------Load script and get resullt--------##########
    ####################################################

    #Function Change-Service{
    
        #param([ref]$global:starttype,[ref]$global:servicenames)
    
        $global:servicename = (Get-Service $global:servicenames).Name
        
        Clear-Host
    
        if($global:servicename -eq $global:servicenames){ 
            $global:startmode = (Get-WmiObject win32_service  | ?{$_.Name -match $global:servicenames}).StartMode  
            $global:delay = Get-ChildItem HKLM:\SYSTEM\CURRENTControlSet\Services | where-object {$_.Property -eq "DelayedAutostart" -and $_.Name -match "$global:servicenames"}
            $serv = Get-ItemProperty 'HKLM:\SYSTEM\CurrentControlSet\Services\MSSQLSERVER' -Name "wmiApSrv"
            $isDelayedAutostart = $serv.DelayedAutostart -eq 0
        if($global:starttype -notmatch $global:startmode -or $global:delay -ne $null ){
        
            #Set service for Automatic
            if($global:starttype -match "Automatic" -and $global:starttype -notmatch "Delayed" `            -and $global:starttype -notmatch "Manual" -and $global:starttype -notmatch "Disabled"){
               if($global:startmode -notmatch "Auto"){
                    write-Host "Setting Service $global:servicenames - Start Mode from $global:startmode to" $global:starttype
                    try{
                
                         Set-Service $global:servicenames -StartupType "Automatic" -Confirm:$false
                         $global:SuccessMessageArray += "Service $global:servicenames Start Mode is successfully changed from $global:startmode to $global:starttype"
                
                         $global:SuccessObject = New-Object -TypeName psobject
                         $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'Message' -Value $($global:SuccessMessageArray)
                         $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'Service' -Value $($global:servicenames)
                         $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'Previous Start Type' -Value $($global:startmode)
                         $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'Updated Start Type' -Value $($global:starttype)
                         $global:SuccessMessageArray += "Service" + $($global:servicenames)
                         $global:SuccessMessageArray += "Previous Start Type " + $($global:startmode)
                         $global:SuccessMessageArray += "Updated Start Type " + $($global:starttype)
                
                    }
           
                    catch{
                        $global:Code = 1
                        $global:ErrorMsg = $_.Exception.Message
                        $global:FailedItem = $_.Exception.ItemName
                        $global:ErrorMessageArray += "Unable to set service startup type `"$global:starttype`" due to error $global:FailedItem :  $global:ErrorMsg" 

                    }
                }
                else{
    
                    #Write-Host "The current service startup type `"$global:starttype`" is already set to the type selected."
                    $global:Code = 1
                    $global:ErrorMessageArray += "The current service startup type `"$global:starttype`" is already set to the type selected."
    
                }
        
            }
        
            #Set service for Automatic Delayed
            if($global:starttype -notmatch "Automatic" -and $global:starttype -match "Delayed" `           -and $global:starttype -notmatch "Manual" -and $global:starttype -notmatch "Disabled"){
                if($startmode -notmatch "Auto"){
                    write-Host "Setting Service $global:servicenames - Start Mode from $global:startmode to" $global:starttype
                   $command = "sc.exe \\$Computer config $global:servicenames start= delayed-auto"
                   $Output = Invoke-Expression -Command $Command -ErrorAction Stop
                        if($LASTEXITCODE -ne 0){
                              $global:Code = 1
                		      Write-Host "$Computer : Failed to set $global:servicenames - to delayed start. More details: $Output" -foregroundcolor red
                			  $failedcomputers += $Computer
                              $global:ErrorMessageArray += "Unable to set service startup type `"$global:starttype`" due to error $global:FailedItem :  $global:ErrorMsg" 
		                 } 
                    else {
              
		      Write-Host "$Computer : Successfully changed $global:servicenames service from $global:startmode to delayed start" -foregroundcolor green
			  $global:SuccessMessageArray += "Service $global:servicenames Start Mode is successfully changed from $global:startmode to $global:starttype"
              $global:SuccessMessageArray += "Service" + $($global:servicenames)
              $global:SuccessMessageArray += "Previous Start Type " + $($global:delay)
              $global:SuccessMessageArray += "Updated Start Type " + $($global:starttype)
                
		   }}
          else{
    
            Write-Host "The current service startup type `"$global:starttype`" is already set to the type selected."
            $global:Code = 1
            $global:ErrorMessageArray += "The current service startup type `"$global:starttype`" is already set to the type selected."
    
        }
        
        }
        
        #Set service for Manual
        if($global:starttype -notmatch "Automatic" -and $global:starttype -notmatch "Delayed" `           -and $global:starttype -match "Manual" -and $global:starttype -notmatch "Disabled"){
           if($global:startmode -ne "Manual"){
           try{
                Set-Service $global:servicenames -StartupType "Manual" -Confirm:$false
                $global:SuccessMessageArray += "Service $global:servicenames Start Mode is successfully changed from $global:startmode to $global:starttype"
                
                 $global:SuccessMessageArray += "Service" + $($global:servicenames)
                 $global:SuccessMessageArray += "Previous Start Type " + $($global:startmode)
                 $global:SuccessMessageArray += "Updated Start Type " + $($global:starttype)
                
           }
           
           catch{
                $global:Code = 1
                $global:ErrorMsg = $_.Exception.Message
                $global:FailedItem = $_.Exception.ItemName
                $global:ErrorMessageArray += "Unable to set service startup type `"$global:starttype`" due to error $global:FailedItem :  $global:ErrorMsg" 

           }}
           
           else{
    
        Write-Host "The current service startup type `"$global:starttype`" is already set to the type selected."
        $global:Code = 1
        $global:ErrorMessageArray += "The current service startup type `"$global:starttype`" is already set to the type selected."
    
        }
        
               
        }
        
        #Set service for Disabled
        if($global:starttype -notmatch "Automatic" -and $global:starttype -notmatch "Delayed" `           -and $global:starttype -notmatch "Manual" -and $global:starttype -match "Disabled"){
           if($global:startmode -ne "Disabled"){
           try{
                Set-Service $global:servicenames -StartupType "Disabled" -Confirm:$false
                $global:SuccessMessageArray += "Service $global:servicenames Start Mode is successfully changed from $global:startmode to $global:starttype"
                $global:SuccessMessageArray += "Service" + $($global:servicenames)
                $global:SuccessMessageArray += "Previous Start Type " + $($global:startmode)
                $global:SuccessMessageArray += "Updated Start Type " + $($global:starttype)
                
           }
           
           catch{
                $global:Code = 1
                $global:ErrorMsg = $_.Exception.Message
                $global:FailedItem = $_.Exception.ItemName
                $global:ErrorMessageArray += "Unable to set service startup type `"$global:starttype`" due to error $global:FailedItem :  $global:ErrorMsg" 

           }}
           
           else{
    
        Write-Host "The current service startup type `"$global:starttype`" is already set to the type selected."
        $global:Code = 1
        $global:ErrorMessageArray += "The current service startup type `"$global:starttype`" is already set to the type selected."
    
        }
        
        }
    } 
   
    else{
    
        Write-Host "The current service startup type `"$global:starttype`" is already set to the type selected."
        $global:Code = 1
        $global:ErrorMessageArray += "The current service startup type `"$global:starttype`" is already set to the type selected."
    
    }
    
    }
    
    else{
    
        #Write-Host "The given serive name `"$global:servicenames`" doesn't match with service present in the machine"
        $global:Code = 1
        $global:ErrorMessageArray += "The given serive name `"$global:servicenames`" doesn't match with service present in the machine"
    } 
    
   # } 

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

    if($global:Code -eq 0)
    {
        $successString = "Success: " + "Service startup type changed"
        $global:SuccessObject = New-Object -TypeName psobject
        $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'Message' -Value $($global:SuccessMessageArray)
        $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'Service' -Value $($global:servicenames)
        $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'Previous Start Type' -Value $($global:startmode)
        $global:SuccessObject | Add-Member -MemberType NoteProperty -Name 'Updated Start Type' -Value $($global:starttype)
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'Status' -Value "success"
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'result' -Value $successString
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'stdout' -Value $global:SuccessMessageArray
        $ResultObject | Add-Member -MemberType NoteProperty -Name 'dataObject' -Value $global:SuccessObject
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
        #$ErrorObjectAray = $global:servicenames 
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
        $JSONResult= $ResultObject|ConvertTo-Json
        $JSONResult
    }
    else
    {
        $JSONResult= $ResultObject|ConvertTo-JSONPS2
        $JSONResult
    }
}




##########################################################
###------Execute Functions -------------------------------
##########################################################
cls

if(Check-PreCondition)
{
    $Continue = $false

    $global:SuccessObject = New-Object -TypeName psobject

    if($global:starttype -eq $null)
    {
        $global:Code = 1
        $global:ErrorMessageArray += "Starttype is not provided"
    }

    if($global:servicenames -eq $null)
    {
        $global:Code = 1
        $global:ErrorMessageArray += "Service Name is not provided"
    }
   # Change-Service -global:starttype $starttype -global:starttype $service
    SetResult
    DisplayResult
}
else
{   
   # Change-Service -global:starttype $starttype -global:starttype $service
    SetResult
    DisplayResult
}

