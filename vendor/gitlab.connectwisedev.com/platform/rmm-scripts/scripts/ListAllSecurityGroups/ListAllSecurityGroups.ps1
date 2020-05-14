Clear-Host
 
<#

[string]$Enviroment = "Local" # "Domain" Local
[string]$SecurityGroup =  "Administrators"       #"Administrators123"
 
[bool]$HideEmptyGroups = $false
[bool]$HideBuiltInGroups = $true

#>

function Escape-JSONString($str){
    if ($str -eq $null) {return ""}
         $str = $str.ToString().Replace('"','\"').Replace('\','\\').Replace("`n",'\n').Replace("`r",'\r').Replace("`t",'\t')
    return $str;
}


function ConvertTo-JSONP2($maxDepth = 10,$forceArray = $false) {
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
$jsonResult += ($elem | ConvertTo-JSONP2 -maxDepth ($maxDepth -1))
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
"@ -f $key , ($value[$key] | ConvertTo-JSONP2 -maxDepth ($maxDepth -1) )
}
return "{" + $jsonResult + "}"
}
            default { #object
if ($maxDepth -le 0){return  "`"{0}`"" -f (Escape-JSONString $value)}
return "{" +
(($value | Get-Member -MemberType *property | % { 
@"
"{0}": {1}
"@ -f $_.Name , ($value.($_.Name) | ConvertTo-JSONP2 -maxDepth ($maxDepth -1) ) 
}) -join ', ') + "}"
    }
}
}
}
 
 
function Test-DomainNetworkConnection
{
    $strOSVersion = (Get-WmiObject -Query "Select Version from Win32_OperatingSystem").Version
    $arrStrOSVersion = $strOSVersion.Split(".")
    $intOSMajorVersion = [UInt16]$arrStrOSVersion[0]
    if ($arrStrOSVersion.Length -ge 2)
    {
        $intOSMinorVersion = [UInt16]$arrStrOSVersion[1]
    } `
    else
    {
        $intOSMinorVersion = [UInt16]0
    }
        
    if (($intOSMajorVersion -gt 6) -or (($intOSMajorVersion -eq 6) -and ($intOSMinorVersion -gt 1)))
    {        
        $domainNetworks = Get-NetConnectionProfile | Where-Object {$_.NetworkCategory -eq "Domain"}
    } `
    else
    {
        $domainNetworks = ([Activator]::CreateInstance([Type]::GetTypeFromCLSID([Guid]"{DCB00C01-570F-4A9B-8D69-199FDBA5723B}"))).GetNetworkConnections() | `
            ForEach-Object {$_.GetNetwork().GetCategory()} | Where-Object {$_ -eq 2}
    }
    return ($domainNetworks -ne $null)
    
}
 
Function Get-LocalORDomainGroup  {
 
  [Cmdletbinding()]
  Param(
  [Parameter(ValueFromPipeline=$True, ValueFromPipelineByPropertyName=$True)]
  [String[]]$Computername =  $Env:COMPUTERNAME,
  [parameter()]
  [string[]]$Group
  )
 
  Begin {
 
      Function  ConvertTo-SID {
 
          Param([byte[]]$BinarySID)
         (New-Object  System.Security.Principal.SecurityIdentifier($BinarySID,0)).Value
         
     }
 
      Function  GetType {
 
          Param([byte[]]$Type)
         (New-Object  System.Security.Principal.SecurityIdentifier($Type,0)).Value
         
     }
     
     Function  Get-LocalGroupMember {
 
          Param  ($Group)
          $group.Invoke('members')  | ForEach {
          $_.GetType().InvokeMember("Name",  'GetProperty',  $null,  $_, $null)}
     }
 
  }
  Process  {
            
           $LocalGroups = @()
           $adsi = ""
 
          Try  {
              
              if([string]::IsNullOrEmpty($Enviroment)){
                
                Write-Warning "Please provide scope as 'Local' or 'Domain' "
                return
              }
 
              if($Enviroment -eq "Local"){$adsi  = [ADSI]"WinNT://$Computername"}
 
              if($Enviroment -eq "Domain"){$adsi  = [ADSI]"WinNT://$env:USERDOMAIN"}
 
              if([string]::IsNullOrEmpty($adsi)){return}
             
                 If($PSBoundParameters.ContainsKey('Group')) {
                      try{
                          $Groups  = ForEach  ($item in  $group) { 
                          $adsi.Children.Find($Item, 'Group')}
                      }catch{ }
 
                  }Else{$groups  = $adsi.Children | where {$_.SchemaClassName -eq  'group'}}
 
                  If  ($groups) {                               
                        $groups  | ForEach {
                              [pscustomobject]@{
                             Name = $_.Name[0]
                             Members = ((Get-LocalGroupMember  -Group $_))  -join ', '                               
                             SID = (ConvertTo-SID -BinarySID $_.ObjectSID[0])                           
                          }  
                      }
 
                  }Else{}#Else{Throw  "No groups found!"}
             
          }Catch{ Write-Warning  "$($Computer): $_"} 
 
   }
 
  }
 
  Function Get_AllSecurityGroups{

   $SecGrpObject = New-Object -TypeName psobject
   $SecGrpObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "List all security groups"
            
   $hostVerSionMajor = ($PSVersionTable.PSVersion.Major).ToString()
   $hostVerSionMinor = ($PSVersionTable.PSVersion.Minor).ToString()
   $hostVersion = $hostVerSionMajor +'.'+ $hostVerSionMinor 
 
   $osVersionMajor = ([System.Environment]::OSVersion.Version.major).ToString()
   $osVersionMinor = ([System.Environment]::OSVersion.Version.minor).ToString()
   $osVersion = $osVersionMajor +'.'+ $osVersionMinor
 
   [boolean]$isPsVersionOk = ([version]$hostVersion -ge [version]'2.0')
   [boolean]$isOSVersionOk = ([version]$osVersion -ge [version]'6.0')
   
   #---------- Check for Powershell version ------------------
   if(-not $isPsVersionOk){
      
      $StdErrArr = @()
      $stdOutArr = @()
      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "PowerShell version below 2.0 is not supported";
                          detail = "PowerShell version below 2.0 is not supported";

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ("PowerShell version below 2.0 is not supported")
    
     $SecGrpObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $SecGrpObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $SecGrpObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $SecGrpObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
     $SecGrpObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $SecGrpObject
 
    }
 
   #---------- Check OS Version ------------------
   if(-not $isOSVersionOk){
   
      $StdErrArr = @()
      $stdOutArr = @()
      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system";
                          detail = "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system";

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ("PowerShell Script supports Window 7, Window 2008R2 and higher version operating system")
    
     $SecGrpObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $SecGrpObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $SecGrpObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $SecGrpObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
     $SecGrpObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $SecGrpObject
    }
 
   #------- Check whether system is part of doamin or not-----------------
   # Test-DomainNetworkConnection
 
   if($Enviroment -eq "Domain"){
     if(-not(Test-DomainNetworkConnection)){
             
              $StdErrArr = @()
              $stdOutArr = @()
              $StdErr = New-Object PSObject -Property @{		       
		                          id = 0;
                                  title =  "Domain scope selected on a machine not joined to a domain.";
                                  detail = "Domain scope selected on a machine not joined to a domain.";

                       }
         
             $StdErrArr += $StdErr
             $stdOutArr += ("Domain scope selected on a machine not joined to a domain.")
    
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name Result -Value "Domain scope selected on a machine not joined to a domain."
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

             return $SecGrpObject
      }
   }
 
   if($Enviroment -eq "Domain"){
       $PartOfDomain = (Get-WmiObject -Class Win32_ComputerSystem).PartOfDomain
       if(-not $PartOfDomain){
                          
              $StdErrArr = @()
              $stdOutArr = @()
              $StdErr = New-Object PSObject -Property @{		       
		                          id = 0;
                                  title =  "Domain scope selected on a machine not joined to a domain.";
                                  detail = "Domain scope selected on a machine not joined to a domain.";

                       }
         
             $StdErrArr += $StdErr
             $stdOutArr += ("Domain scope selected on a machine not joined to a domain.")
    
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name Result -Value "Domain scope selected on a machine not joined to a domain."
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

             return $SecGrpObject
         }
   }
   
   if($Enviroment -eq "Local"){
      $osInfo = Get-WmiObject -Class Win32_OperatingSystem     
      if( $osInfo.ProductType -eq 2){

        Write-Warning "Local scope selected on a Domain Controller";
        return

              $StdErrArr = @()
              $stdOutArr = @()
              $StdErr = New-Object PSObject -Property @{		       
		                          id = 0;
                                  title =  "Local scope selected on a Domain Controller";
                                  detail = "Local scope selected on a Domain Controller";

                       }
         
             $StdErrArr += $StdErr
             $stdOutArr += ("Local scope selected on a Domain Controller")
    
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name Result -Value "Local scope selected on a Domain Controller";
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

             return $SecGrpObject

      }
     
   }
 
    [string[]]$SecurityGroupsArr = @()
    if(-not([string]::IsNullOrEmpty($SecurityGroup))){ $SecurityGroupsArr = $SecurityGroup.Split(',') }
 
    $LocalDomainGroups = @()
    $LocalDomainGroups = Get-LocalORDomainGroup -Computername  $env:COMPUTERNAME
 
    # if($Enviroment -eq "Local"){ Write-Host  "Pulling groups from system: $($env:COMPUTERNAME)"}
    # if($Enviroment -eq "Domain"){Write-Host  "Pulling groups from system: $($env:USERDOMAIN)"}
 
           
          
   if($SecurityGroupsArr.Count -gt 0){  
           
        $Code = 0
        $Status = "Success"    
        $ResultMsg = "Success : Pulling groups: $($SecurityGroupsArr -join ',')"
       
        $SIDArr = @()
        $SID = ""
        $LastIndexValue = ""
        $LastIndexValueArr = @()
        $Uses = "In Use"
        $Type = "Custom"
        $GroupsExist = @()
        $GroupsDoesNotExist = @()
        $GroupNames = @()
        
        $StdArr = @()
        $StdArrStr = @()
        $StdErrArr = @()

        foreach($gn in $SecurityGroupsArr){  
                 
           $x = $LocalDomainGroups | Where-Object {$_.Name -eq $gn}
 
           if($x){$GroupsExist += $gn}
           else{$GroupsDoesNotExist += $gn}           
        }
        
        if($GroupsExist.Count -gt 0){
 
         # Write-Host "Pulling group(s): $($GroupsExist -join ',')" 
         # Write-Host "-------------------------------------------"
          $LocalGroups = Get-LocalORDomainGroup -Computername  $env:COMPUTERNAME -Group $GroupsExist
 
          if($LocalGroups){
              $LocalGroups | ForEach {
                    
                   
                     if([string]::IsNullOrEmpty($_.Members)){$Uses = "Empty"}
                     $SID = $_.SID
 
                     $SIDArr = $SID.Split('-')
                     $LastIndexValue = $SIDArr[$SIDArr.length -1]
            
                     $LastIndexValueArr = $LastIndexValue.ToCharArray()           
                     if(($LastIndexValueArr.length -eq "3") -and ($LastIndexValueArr[0] -eq "5")){$Type = "Built In"} 
                     
                     # New-Object psobject -Property @{Name = $_.Name;Status = $Uses;Type = $Type} | select Name,Status,Type

                     $Object = New-Object psobject -Property @{Name = $_.Name;Status = $Uses;Type = $Type}
                 
                     $Name = $_.Name;
                     $Status = $Uses;
                     $Type = $Type

                     $StdArr += $Object
                     $StdArrStr += "Name : $Name, Status : $Status, Type : $Type"
                      
                 }
           }else {
     
           # Write-Warning "No Group found"
           # Write-Host ""
     
         }
        }
 
        if($GroupsDoesNotExist.Count -gt 0){
           
          # Write-Host ""
          # Write-Warning "No group(s) found containing : $($GroupsDoesNotExist -join ',')"  

           # $SuccessMsg = "Pulling groups: $($SecurityGroupsArr -join ',')"

           $ErrorTitle = "No group(s) found containing : $($GroupsDoesNotExist -join ',')" 
                   
           $FL1 = ""
           $FL = ""
           foreach($sg in $GroupsDoesNotExist){

            $SecurityGroupsCharArr = $sg.ToCharArray()
            $FL1 +=  $SecurityGroupsCharArr[0]+',' 
            
            }
 
          $LI = $FL1.LastIndexOf(',')
          $FL1 = $FL1.SubString(0,$LI)

         $ErrorDetail = "List of all group starting with $($FL1) "
          #  Write-Host "-----------------------------------------"

         $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = $ErrorTitle;
                          detail = $ErrorDetail;

               }


          $StdErrArr += $StdErr
          
          foreach($sg in $GroupsDoesNotExist){
             
             $SecurityGroupsCharArr = $sg.ToCharArray()
             $FL =  $SecurityGroupsCharArr[0]+'*'             
             $FL = $FL.Trim()         
           
             $LocalGroups =  Get-LocalORDomainGroup -Computername  $env:COMPUTERNAME | Where-Object{ $_.Name -like  $FL }
             
             if($LocalGroups){
                 $LocalGroups | ForEach {

                     if([string]::IsNullOrEmpty($_.Members)){$Uses = "Empty"}
                     $SID = $_.SID
 
                     $SIDArr = $SID.Split('-')
                     $LastIndexValue = $SIDArr[$SIDArr.length -1]
            
                     $LastIndexValueArr = $LastIndexValue.ToCharArray()           
                     if(($LastIndexValueArr.length -eq "3") -and ($LastIndexValueArr[0] -eq "5")){$Type = "Built In"} 

                     # New-Object psobject -Property @{Name = $_.Name;Status = $Uses;Type = $Type}|select Name,Status,Type

                     $Object = New-Object psobject -Property @{Name = $_.Name;Status = $Uses;Type = $Type}
                 
                     $Name = $_.Name;
                     $Status = $Uses;
                     $Type = $Type

                     $StdArr += $Object
                     $StdArrStr += "Name : $Name, Status : $Status, Type : $Type"
                    
               }

           } else {


            if($GroupsExist.Count -eq 0){
              $StdArrStr += "No Group Found"
            }

            $SecGrpObj =  New-Object psobject
            $SecGrpObj | Add-Member -MemberType NoteProperty -Name SecurityGroups -Value $StdArr
     
            $SecGrpObject | Add-Member -MemberType NoteProperty -Name Status -Value "Fail"
            $SecGrpObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
            $SecGrpObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error : No Group Found"
            $SecGrpObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $StdArrStr
            $SecGrpObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr

            if($GroupsExist.Count -gt 0){

                $SecGrpObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $SecGrpObj  
            }                                     
            return $SecGrpObject
     
         }
 
          }# end foreach      
 
        }

        # $StdErrArr

        if($GroupsDoesNotExist.Count -gt 0){
          
           $Code = 1
           $Status = "Fail"
           $ResultMsg = "Fail : Failed to pull group(s): $($SecurityGroupsArr -join ',')"
          
        }


        $SecGrpObj =  New-Object psobject
        $SecGrpObj | Add-Member -MemberType NoteProperty -Name SecurityGroups -Value $StdArr
              
        $SecGrpObject | Add-Member -MemberType NoteProperty -Name Status -Value $Status
        $SecGrpObject | Add-Member -MemberType NoteProperty -Name Code -Value $Code
        $SecGrpObject | Add-Member -MemberType NoteProperty -Name Result -Value "$ResultMsg"
        $SecGrpObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $StdArrStr
        $SecGrpObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $SecGrpObj  

        if($GroupsDoesNotExist.Count -gt 0){

            $SecGrpObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr

        }
                          
        return $SecGrpObject

  } 
   
   
   if([string]::IsNullOrEmpty($SecurityGroup) -and $SecurityGroupsArr.Count -eq 0){
 
     if($LocalDomainGroups){
          
          $StdArr = @()
          $StdArrStr = @()
          
          $LocalDomainGroups | ForEach {
              
           $SIDArr = @()
           $SID = ""
           $LastIndexValue = ""
           $LastIndexValueArr = @()
 
           $Uses = "In Use"
           $Type = "Custom"
 
           if([string]::IsNullOrEmpty($_.Members)){$Uses = "Empty"}
       
           $SID = $_.SID                                 
           $SIDArr = $SID.Split('-')
           $LastIndexValue = $SIDArr[$SIDArr.length -1]
           $LastIndexValueArr = $LastIndexValue.ToCharArray()
 
           if(($LastIndexValueArr.length -eq "3") -and ($LastIndexValueArr[0] -eq "5")){$Type = "Built In"}  
              
           if($HideEmptyGroups -and $HideBuiltInGroups){                
             if( ($Uses -eq "In Use") -and ( $Type -eq "Custom")){

                  #New-Object psobject -Property @{Name = $_.Name;Status = $Uses;Type = $Type}|select Name,Status,Type
                  $Object = New-Object psobject -Property @{Name = $_.Name;Status = $Uses;Type = $Type}
                 
                  $Name = $_.Name;
                  $Status = $Uses;
                  $Type = $Type

                  $StdArr += $Object
                  $StdArrStr += "Name : $Name, Status : $Status, Type : $Type" 
          

                }
             } 
              
            if($HideEmptyGroups -and -not($HideBuiltInGroups)){
              if($Uses -eq "In Use"){
                                 
                 # New-Object psobject -Property @{Name = $_.Name;Status = $Uses;Type = $Type}|select Name,Status,Type}

                 $Object = New-Object psobject -Property @{Name = $_.Name;Status = $Uses;Type = $Type}
                 
                  $Name = $_.Name;
                  $Status = $Uses;
                  $Type = $Type

                  $StdArr += $Object
                  $StdArrStr += "Name : $Name, Status : $Status, Type : $Type" 
          
               }    
              }

            if(-not($HideEmptyGroups) -and $HideBuiltInGroups){
               if($Type -eq "Custom"){

                   # New-Object psobject -Property @{Name = $_.Name;Status = $Uses;Type = $Type}|select Name,Status,Type}
                   $Object = New-Object psobject -Property @{Name = $_.Name;Status = $Uses;Type = $Type}
                 
                  $Name = $_.Name;
                  $Status = $Uses;
                  $Type = $Type

                  $StdArr += $Object
                  $StdArrStr += "Name : $Name, Status : $Status, Type : $Type" 
          
               }    
              }

            if(-not($HideEmptyGroups) -and -not($HideBuiltInGroups)){

                  #  New-Object psobject -Property @{Name = $_.Name;Status = $Uses;Type = $Type}|select Name,Status,Type

                  $Object = New-Object psobject -Property @{Name = $_.Name;Status = $Uses;Type = $Type}
                 
                  $Name = $_.Name;
                  $Status = $Uses;
                  $Type = $Type

                  $StdArr += $Object
                  $StdArrStr += "Name : $Name, Status : $Status, Type : $Type" 
          
              } 
         
         }#end foreach 

             $SecGrpObj =  New-Object psobject
             $SecGrpObj | Add-Member -MemberType NoteProperty -Name SecurityGroups -Value $StdArr
              
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : Pulling all groups,Hide empty groups : $($HideEmptyGroups),Hide built-in groups : $($HideBuiltInGroups)"
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $StdArrStr
             $SecGrpObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $SecGrpObj  
                          
             return $SecGrpObject


     } # end if
     else {
     
        Write-Warning "No Group found"
        Write-Host ""
     
     }
   }
     
  }


if($PSVersionTable.PSVersion.Major -eq 2){

    Get_AllSecurityGroups |  ConvertTo-JSONP2

}else{

    Get_AllSecurityGroups |  ConvertTo-Json -Depth 10
}
