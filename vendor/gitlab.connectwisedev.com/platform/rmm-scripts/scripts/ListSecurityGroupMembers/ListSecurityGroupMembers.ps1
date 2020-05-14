Clear-Host


<#

[string]$GroupName = "Power"  #"Administrators"
[boolean]$ISRecurse = $true
[boolean]$IsDomainContext = $true

#>

$SecurityGroupMembers = @()
$SecurityGroups = @()

$UserArr = New-Object System.Collections.ArrayList
$GroupArr = New-Object System.Collections.ArrayList
$EmptyGroupArr = New-Object System.Collections.ArrayList
$InUsegroupArr = New-Object System.Collections.ArrayList


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



Function Get-SecurityGroupMembers {

  [Cmdletbinding()]
  Param(
  [parameter()]
  [string[]]$Group
  )

  Begin {
      Function  ConvertTo-SID {

          Param([byte[]]$BinarySID)
          (New-Object  System.Security.Principal.SecurityIdentifier($BinarySID,0)).Value
     }
     
     Function  Get-LocalGroupMember {

          Param  ($Group)
          $group.Invoke('members')  | ForEach {
          
          $_.GetType().InvokeMember("Name",  'GetProperty',  $null,  $_, $null)}
     }

  }
  Process  {

          Try  {
                $adsi = $null

                 if($IsDomainContext){
                    $adsi  = [ADSI]"WinNT://$env:USERDOMAIN"
                 }else{
                    $adsi  = [ADSI]"WinNT://$env:COMPUTERNAME"
                 }

                 
                  If($PSBoundParameters.ContainsKey('Group')) {
                   
                      $groups  = ForEach  ($item in  $group) { 
                       
                        $adsi.Children.Find($Item, 'Group')
                      }

                  }Else{
                     $groups  = $adsi.Children | where {$_.SchemaClassName -eq  'group'}
                  }

                  If  ($groups) {
                      $groups  | ForEach {
                     
                        $Members = (Get-LocalGroupMember  -Group $_)
                      }

                  }Else{Throw  "No groups found!"}

          }Catch{ Write-Warning  "$($Computer): $_"}
 
    return $Members
   }

  }


Function Get-SecurityGroups{
     
    if($IsDomainContext){

        [ADSI]$S = "WinNT://$env:USERDOMAIN"

      }else{

        [ADSI]$S = "WinNT://$env:COMPUTERNAME"
      }
      
           
      foreach( $grp in $S.children){
       if($grp.class -eq 'group'){
        $grp.name.value
       }
      
      }

    }




function List-SecurityGroupMembers{
   try {


    $SGMObject = New-Object -TypeName psobject
    $SGMObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "List security group members"

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
    
     $SGMObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $SGMObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $SGMObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $SGMObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
     $SGMObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $SGMObject
 
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
    
     $SGMObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $SGMObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $SGMObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $SGMObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
     $SGMObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $SGMObject
   }
   
    
    $IsSecurityGrpFound = $False
    $SecurityGroups = Get-SecurityGroups

    foreach($sg in $SecurityGroups){

       if($sg.ToString() -eq $GroupName){

           $IsSecurityGrpFound = $True
       }

    }

    if(-not $IsSecurityGrpFound){
              
       $SecGrpExist = @()

       $Fl = ($GroupName.ToCharArray())[0]
       $FlSearch = $Fl +'*'

            
        $SecurityGroups | Where-Object { $_.ToString() -like $FlSearch }| foreach{
         
           $sg = $_
           $SecGrpExist += $sg
         
         }                 

       if($SecGrpExist.Count -gt 0){
        
              $StdErrArr = @()
              $stdOutArr = @()
              $StdErr = New-Object PSObject -Property @{		       
		                          id = 0;
                                  title =  "Security Group : $($GroupName) not found";
                                  detail = "List of all security group start with $($Fl)";

                       }
                                   

             $StdErrArr += $StdErr
             $stdOutArr += ("List of all security group start with $($Fl)")

             $SGMObj =  New-Object psobject
             $SGMObj | Add-Member -MemberType NoteProperty -Name SecurityGroups -Value $SecGrpExist
    
             $SGMObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $SGMObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $SGMObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
             $SGMObject | Add-Member -MemberType NoteProperty -Name Result -Value "List of all security group start with $($Fl)"
             $SGMObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr
             $SGMObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $SGMObj    

             return $SGMObject

       }

       
       if($SecGrpExist.Count -eq 0){
        
        
              $StdErrArr = @()
              $stdOutArr = @()
              $StdErr = New-Object PSObject -Property @{		       
		                          id = 0;
                                  title =  "Security Group : $($GroupName) not found";
                                  detail = "List of all security group start with $($Fl) is not found"

                       }
                 
             $stdOutArr += ("List of all security group start with $($Fl) is not found")                                                      
             
             $SGMObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $SGMObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $SGMObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErr
             $SGMObject | Add-Member -MemberType NoteProperty -Name Result -Value "List of all security group start with $($Fl) is not found"
             $SGMObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr
             
             return $SGMObject
  
       }

              
    }

    $SecurityGroupMembers = Get-SecurityGroupMembers -Group $GroupName

    foreach($sg in $SecurityGroups){

        foreach($sgm in $SecurityGroupMembers){

            if($sgm -eq $sg){
             
              $null = $GroupArr.Add($sgm)

            }

        }

    }

   
    foreach($sgm in $SecurityGroupMembers){

    if(!$GroupArr.Contains($sgm)){
                  
         $null = $UserArr.Add($sgm)
      }

    }

    foreach($grp in $GroupArr){

      $temparr = @()
      $temparr = Get-SecurityGroupMembers -Group  $grp

      if($temparr.Count -gt 0){
 
        $null = $InUsegroupArr.Add($grp)
    
      }else{

        $null = $EmptyGroupArr.Add($grp)

      }

    }
    
   
    if(!$ISRecurse){

       $SecGrpMem = @()
       
       foreach($user in $UserArr){
               
          $SecGrpMem += "|--"+"$GroupName\$user"
          
        }
        foreach($eg in $EmptyGroupArr){
                   
          $SecGrpMem += "|--"+"[Group]"+"$GroupName\$eg"+"(empty)" 
          
        }
        foreach($ug in $InUsegroupArr){

           $SecGrpMem += "|--"+"[Group]"+"$GroupName\$ug"
           
        }
       
         $StdArr = @()
         $stdOutArr = @()
           
         $SGMObj =  New-Object psobject
         $SGMObj | Add-Member -MemberType NoteProperty -Name SecurityGroupMembers -Value $SecGrpMem
    
         $SGMObject | Add-Member -MemberType NoteProperty -Name Status -Value "success" 
         $SGMObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
         $SGMObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success: List of security members"
         $SGMObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $SecGrpMem
         $SGMObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $SGMObj    

         return $SGMObject
    }

    if($ISRecurse){   
      
       $SecGrpMem = @()
       
       foreach($user in $UserArr){
               
          $SecGrpMem += "|--"+"$GroupName\$user"
          
        }
        foreach($eg in $EmptyGroupArr){
                   
          $SecGrpMem += "|--"+"[Group]"+"$GroupName\$eg"+"(empty)" 
          
        }

        foreach($ug in $InUsegroupArr){

            $SecGrpMem += "|--"+"[Group]"+"$GroupName\$ug"
           
            $temparr1 = @()
            $temparr1 = Get-SecurityGroupMembers -Group  $ug

            foreach($mem in $temparr1){
               
               $SecGrpMem += "  "+"|--"+"$GroupName\$mem"
                                            
            }
         
        }

         $StdArr = @()
         $stdOutArr = @()
           
         $SGMObj =  New-Object psobject
         $SGMObj | Add-Member -MemberType NoteProperty -Name SecurityGroupMembers -Value $SecGrpMem
    
         $SGMObject | Add-Member -MemberType NoteProperty -Name Status -Value "success" 
         $SGMObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
         $SGMObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success: List of security members"
         $SGMObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $SecGrpMem
         $SGMObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $SGMObj    

         return $SGMObject
     
    }

    }catch {
     
     #  $_
              $StdErrArr = @()
              $stdOutArr = @()
              $StdErr = New-Object PSObject -Property @{		       
		                          id = 0;
                                  title =  "Unable to pull list of group members due to error"
                                  detail = "$_"

                       }
                 
             $stdOutArr += "Unable to pull list of group members due to error"
             
             $SGMObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $SGMObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $SGMObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErr
             $SGMObject | Add-Member -MemberType NoteProperty -Name Result -Value "Unable to pull list of group members due to error"
             $SGMObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr
             
             return $SGMObject

    }
}


if($PSVersionTable.PSVersion.Major -eq 2){

    List-SecurityGroupMembers |  ConvertTo-JSONP2

}else{

    List-SecurityGroupMembers |  ConvertTo-Json -Depth 10
}