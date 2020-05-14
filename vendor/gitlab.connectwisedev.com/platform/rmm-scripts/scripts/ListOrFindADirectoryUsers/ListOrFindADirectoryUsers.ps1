Clear-Host

<#


cat - Local Domain AD

[boolean]$HideBuiltInUser = $true -optional
[string] $UserName = ""           -optional 
[string] $OUname = ""             -optional

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

function List-UserInActiveDirectoy{
        
    $Domain = New-Object System.DirectoryServices.DirectoryEntry
    $Searcher = New-Object System.DirectoryServices.DirectorySearcher
    $Searcher.SearchRoot = $Domain
    $Searcher.SearchScope = "subtree"

    $Searcher.PropertiesToLoad.Add("distinguishedName") > $Null
    $Searcher.PropertiesToLoad.Add("Name") > $Null
    $Searcher.Filter = "(&(objectCategory=person)(objectClass=user))"
    $Users = $Searcher.FindAll()

     $UserArr = @()
     $UserArrStr = @()
   
    ForEach ($User In $Users){
    
        if($HideBuiltInUser){
            
            if(($User.Properties.Item("Name") -eq "Guest") -or ($User.Properties.Item("Name") -eq "Administrator")){
                continue
            }
          }

       $Object1 = New-Object PSObject -Property @{		       
		                UserName = [string]$User.Properties.Item("Name"); 
                        OUName = [string]$User.Properties.Item("distinguishedName");
                    }
       $UserName = [string]$User.Properties.Item("Name"); 
       $OUName = [string]$User.Properties.Item("distinguishedName"); 
      
       $UserArr +=  $Object1
       $UserArrStr += ("UserName : $UserName, OUName : $OUName" )
   
    }

    $ADObj =  New-Object psobject
    $ADObj | Add-Member -MemberType NoteProperty -Name Users -Value $UserArr
              
    $ADObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
    $ADObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
    $ADObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : List of all user objects in Active Directory"
    $ADObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $UserArrStr
    $ADObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $ADObj  

    return $ADObject
}

function Find-UserInActiveDirectoy{
   
    [bool]$IsUserExist = $false
    [bool]$IsOUExist = $false

    $ADSISearcherUser = [ADSISearcher]'(objectclass=user)'   
    $ADSISearcherUser.Filter = “name=$UserName” 
    $User = $ADSISearcherUser.FindOne() 
    
    if($User){$IsUserExist = $true}    

    $ADSISearcherOU = [ADSISearcher]'(objectclass=OU)'
    $ADSISearcherOU.Filter = “OU=$OUname” 
    $OU = $ADSISearcherOU.FindOne() 

    if($OU){$IsOUExist = $true}

    # Case : 1 User name not found and OU found
    if(-not $IsUserExist -and $IsOUExist){
        
        $Domain = New-Object System.DirectoryServices.DirectoryEntry
        $Searcher = New-Object System.DirectoryServices.DirectorySearcher
        $Searcher.SearchRoot = $Domain
        $Searcher.SearchScope = "subtree"

        $Searcher.PropertiesToLoad.Add("distinguishedName") > $Null
        $Searcher.PropertiesToLoad.Add("Name") > $Null
        $Searcher.Filter = "(objectCategory=organizationalUnit)"
        $OUs = $Searcher.FindAll()

         $UserArr = @()
         $UserArrStr = @()
         $UserArrErr = @()
      
        ForEach ($OU In $OUs){
                   
        if($OU.Properties.Item("Name") -eq $OUname){

                $OUDN = $OU.Properties.Item("distinguishedName")
                $OUBase = New-Object System.DirectoryServices.DirectoryEntry "LDAP://$OUDN"
                $Searcher.SearchRoot = $OUBase

                $Searcher.Filter = "(&(objectCategory=person)(objectClass=user))"
                $Users = $Searcher.FindAll()

                ForEach ($User In $Users)
                {                
                  $Object1 = New-Object PSObject -Property @{		       
		                UserName = [string]$User.Properties.Item("Name"); 
                        OUName = [string]$OUDN;
                    }

                 $UserName = [string]$User.Properties.Item("Name"); 
                 $OUName = [string]$OUDN;

                 $UserArr += $Object1
                 $UserArrStr += ("UserName : $UserName, OUName : $OUName")
                 
                }
            }
        }
        
        $StdErr = New-Object PSObject -Property @{		       
		                          id = 0;
                                  title = "User name not found , OU found"
                                  detail = "List of all users in OU :$($OUname)"

                       }


      $UserArrErr += $StdErr
       
      $ADObj1 =  New-Object psobject
      $ADObj1 | Add-Member -MemberType NoteProperty -Name Users -Value $UserArr
              
      $ADObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
      $ADObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
      $ADObject | Add-Member -MemberType NoteProperty -Name stderr -Value $UserArrErr
      $ADObject | Add-Member -MemberType NoteProperty -Name Result -Value "User name not found , OU found"
       
      $ADObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $UserArrStr
      $ADObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $ADObj1  
       
      return $ADObject
    }

    # Case : 2 User name found OU not found
    if($IsUserExist -and  -not $IsOUExist){
        
        [int]$count = 0
        $OUnameCharArr = $OUname.ToCharArray()
        $OUnameFirstLetter = $OUnameCharArr[0]+'*'

        $ADSISearcherOU1 = [ADSISearcher]'(objectclass=organizationalUnit)'
        $OUs = $ADSISearcherOU1.FindAll() 
       
               
        foreach ($OU in $OUs){
              $x = $OU.Properties |select @{N="Name"; E={$_.name}} | Where-Object {$_.name -like $OUnameFirstLetter} 
              if($x){
                             
                $count++  
            }      
         }
       
        if($count -gt 0){
               
               $OUArr = @()
               $OUArrStr = @()
               $OUArrErr = @()

               foreach ($OU in $OUs){
                  $x = $OU.Properties |select @{N="Name"; E={$_.name}} | Where-Object {$_.name -like $OUnameFirstLetter} 
                  if($x){
                        
                      $Object1 = New-Object PSObject -Property @{		       
		                    OU = $x.Name;
                        }

                    $OU = $x.Name;

                    $OUArr += $Object1
                    $OUArrStr += ("OU :$OU")
                  
                }      
              }

               
               $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = "User name found , OU not found"
                          detail = "List of all OU Starting with:$($OUnameFirstLetter)"

               }
              
              $OUArrErr += $StdErr
                
              $ADObj =  New-Object psobject
              $ADObj | Add-Member -MemberType NoteProperty -Name OUnits -Value $OUArr
                
              $ADObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
              $ADObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
              $ADObject | Add-Member -MemberType NoteProperty -Name stderr -Value $OUArrErr
              $ADObject | Add-Member -MemberType NoteProperty -Name Result -Value "User name found , OU not found"
              
              $ADObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $OUArrStr
              $ADObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $ADObj  

        }

        if($count -eq 0){
          
         $OUArr1 = @()
         $OUArrStr1 = @()
         $OUArrErr1 = @()
          
          foreach ($OU in $OUs){ 
                       
            $Object2 = New-Object PSObject -Property @{		       
		                OU =[string]($OU.Properties).name;
                    }

             $OU = [string]($OU.Properties).name;

             $OUArr1 += $Object2
             $OUArrStr1 += ("OU :$OU")

          }

           $StdErr1 = New-Object PSObject -Property @{		       
		                  id = 0;
                          title = "List of all OU in domain"
                          detail = "List of all OU Starting with:$($OUnameFirstLetter) is not found"

               }
              
              $OUArrErr1 += $StdErr1
               
              $ADObj1 =  New-Object psobject
              $ADObj1 | Add-Member -MemberType NoteProperty -Name OUnits -Value $OUArr1
              
              $ADObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
              $ADObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
              $ADObject | Add-Member -MemberType NoteProperty -Name stderr -Value $OUArrErr1
              $ADObject | Add-Member -MemberType NoteProperty -Name Result -Value "User name found , OU not found"
        
              $ADObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $OUArrStr1
              $ADObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $ADObj1  

        }

      return $ADObject
    }

    # Case 3: User name not found OU not found
    if(-not $IsUserExist -and  -not $IsOUExist){
      
      # List of all users staring with first letter
             
        [int]$UserCount = 0
        $UserNameCharArr = $UserName.ToCharArray()
        $UserNameFirstLetter = $UserNameCharArr[0]+'*'
               
        $UserArr = @()
        $UserArrStr = @()
        $UserArrErr = @()
      
        $ADSISearcherUser = [ADSISearcher]'(objectclass=user)'
        $Users = $ADSISearcherUser.FindAll() 
        
       
        foreach ($User in $Users){

              $u = $User.Properties |select @{N="Name"; E={$_.name}} | Where-Object {$_.name -like $UserNameFirstLetter} 
              if($u){  
                            
                $UserCount++  
                      
              }    
         }
        
        if($UserCount -gt 0){

            foreach ($User in $Users){

              $u = $User.Properties |select @{N="Name"; E={$_.name}} | Where-Object {$_.name -like $UserNameFirstLetter} 
               if($u){  
                             
                    $Object1 = New-Object PSObject -Property @{		       
		                    UserName =[string]$u.name;
                        }
               
                   $UserName =[string]$u.name;

                   $UserArr += $Object1;
                   $UserArrStr +=  ( "UserName : $UserName")
     
                }
                   
             }
            
             $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "User name not found"
                          detail = "List of all users Starting with:$($UserNameFirstLetter)"

               }
              
              $UserArrErr += $StdErr

        }
       
        if($UserCount -eq 0){
         
          foreach ($User in $Users){ 

               $Object2 = New-Object PSObject -Property @{		       
		                    UserName =[string]($User.Properties).name;
                        }
              
                $UserName =[string]($User.Properties).name;

                $UserArr += $Object2;
                $UserArrStr +=  ( "UserName : $UserName")
                              
          }

          $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "List of all user Starting with:$($UserNameFirstLetter) is not found"
                          detail = "List of all users in domain"

               }
              
         $UserArrErr += $StdErr

        }


      # List of all OU staring with first letter
      
        [int]$count = 0
        $OUnameCharArr = $OUname.ToCharArray()
        $OUnameFirstLetter = $OUnameCharArr[0]+'*'
               
               
        $ADSISearcherOU1 = [ADSISearcher]'(objectclass=organizationalUnit)'
        $OUs = $ADSISearcherOU1.FindAll() 
       
        foreach ($OU in $OUs){
              $x = $OU.Properties |select @{N="Name"; E={$_.name}} | Where-Object {$_.name -like $OUnameFirstLetter} 
              if($x){  
                         
                $count++ 
                           
              }    
         }

         if($count -gt 0){
            
             foreach ($OU in $OUs){

              $x = $OU.Properties |select @{N="Name"; E={$_.name}} | Where-Object {$_.name -like $OUnameFirstLetter} 
              if($x){  
                             
                $Object3 = New-Object PSObject -Property @{		       
		                OU = $x.Name;
                    }

                 $OU = $x.Name;

                 $UserArr += $Object3;
                 $UserArrStr +=  ( "OU : $OU")
           
                 }    
              }


              $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "OU was not found"
                          detail = "List of all OU Starting with:$($OUnameFirstLetter)"

                  }

               $UserArrErr += $StdErr
           }
       
        if($count -eq 0){
                
         foreach ($OU in $OUs){

            $Object4 = New-Object PSObject -Property @{		       
		                OU =[string]($OU.Properties).name;
                    }
          
             $OU = [string]($OU.Properties).name;

             $UserArr += $Object4;
             $UserArrStr +=  ( "OU : $OU")
          }

          $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "List of all OU Starting with:$($OUnameFirstLetter) is not found"
                          detail = "List of all OU in domain"

                  }
           $UserArrErr += $StdErr
        }

                      
        $ADObj1 =  New-Object psobject
        $ADObj1 | Add-Member -MemberType NoteProperty -Name UsersAndOUs -Value $UserArr
              
        $ADObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
        $ADObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
        $ADObject | Add-Member -MemberType NoteProperty -Name stderr -Value $UserArrErr
        $ADObject | Add-Member -MemberType NoteProperty -Name Result -Value "User name not found , OU not found"
        
        $ADObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $UserArrStr
        $ADObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $ADObj1  


       return $ADObject
    }

    # User name found and OU name found
     if($IsUserExist -and $IsOUExist){ 
       

        $Domain = New-Object System.DirectoryServices.DirectoryEntry
        $Searcher = New-Object System.DirectoryServices.DirectorySearcher
        $Searcher.SearchRoot = $Domain
        $Searcher.SearchScope = "subtree"

        $Searcher.PropertiesToLoad.Add("distinguishedName") > $Null
        $Searcher.PropertiesToLoad.Add("Name") > $Null
        $Searcher.Filter = "(objectCategory=organizationalUnit)"
        $OUs = $Searcher.FindAll()

        $UserOUArr = @()
        $UserOUArrStr = @()
       
        ForEach ($OU In $OUs){
   
           if($OU.Properties.Item("Name") -eq $OUname){
    
                $OUDN = $OU.Properties.Item("distinguishedName")
                $OUBase = New-Object System.DirectoryServices.DirectoryEntry "LDAP://$OUDN"
                $Searcher.SearchRoot = $OUBase

                $Searcher.Filter = "(&(objectCategory=person)(objectClass=user))"
                $Users = $Searcher.FindAll()
                ForEach ($User In $Users)
                {
                   if($User.Properties.Item("Name") -eq $UserName ){
                                        
                      $Object1 = New-Object PSObject -Property @{	
                      	       
		                UserName =[string]$User.Properties.Item("Name");
                        OU = [string]$OUDN;
                      }
                     
                      $UserName =[string]$User.Properties.Item("Name");
                      $OU = [string]$OUDN;  
                       
                      $UserOUArr += $Object1
                      $UserOUArrStr += ("UserName : $UserName, OU : $OU")

                    }
                }
                
            }
        }


        $ADObj =  New-Object psobject
        $ADObj | Add-Member -MemberType NoteProperty -Name Users -Value $UserOUArr
              
        $ADObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
        $ADObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
        $ADObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : Both user name and OU name are found"
        $ADObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $UserOUArrStr
        $ADObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $ADObj 
      
      return $ADObject
     
     }


}

function ListOrFind-UserInActiveDirectoy{
   
   
   $ADObject = New-Object -TypeName psobject
   $ADObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "List or find user in active directoy"

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
    
     $ADObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $ADObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $ADObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $ADObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
     $ADObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $ADObject
 
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
    
     $ADObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $ADObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $ADObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $ADObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
     $ADObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

    
     return $ADObject

   }
    
   if(-not(Test-DomainNetworkConnection)){
   
      $StdErrArr = @()
      $stdOutArr = @()
      $StdErr = New-Object PSObject -Property @{		       
		                  id = 0;
                          title =  "System has not joined to a domain."
                          detail = "System has not joined to a domain."

               }
         
     $StdErrArr += $StdErr
     $stdOutArr += ("System has not joined to a domain.")
    
     $ADObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $ADObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $ADObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $ADObject | Add-Member -MemberType NoteProperty -Name Result -Value "System has not joined to a domain."
     $ADObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

   
    return $ADObject

   }

   try{
        # user name provide OU not provided
       if(-not([string]::IsNullOrEmpty($UserName)) -and [string]::IsNullOrEmpty($OUname)){
             
          $StdErrArr = @()
          $stdOutArr = @()
          $StdErr = New-Object PSObject -Property @{		       
		                      id = 0;
                              title =  "Please provide OU name"
                              detail = "Please provide OU name"

                   }
         
         $StdErrArr += $StdErr
         $stdOutArr += ("Please provide OU name")
    
         $ADObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
         $ADObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
         $ADObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
         $ADObject | Add-Member -MemberType NoteProperty -Name Result -Value "Please provide OU name"
         $ADObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
        
         return $ADObject

       }

       # user name not provide OU name provided
       if([string]::IsNullOrEmpty($UserName) -and -not([string]::IsNullOrEmpty($OUname))){
                  
          $StdErrArr = @()
          $stdOutArr = @()
          $StdErr = New-Object PSObject -Property @{		       
		                      id = 0;
                              title =  "Please provide user name"
                              detail = "Please provide user name"

                   }
         
         $StdErrArr += $StdErr
         $stdOutArr += ("Please provide OU name")
    
         $ADObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
         $ADObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
         $ADObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
         $ADObject | Add-Member -MemberType NoteProperty -Name Result -Value "Please provide user name"
         $ADObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   
        

         return $ADObject       

       }

       # Both user and ou is not provided , list all user with domain
       if([string]::IsNullOrEmpty($UserName) -and [string]::IsNullOrEmpty($OUname)){
         
           List-UserInActiveDirectoy
       }

        # Both user and ou is provided
       if(-not([string]::IsNullOrEmpty($UserName)) -and -not([string]::IsNullOrEmpty($OUname))){
          
           Find-UserInActiveDirectoy
       }

   }catch{
        
        # Write-Host "Message: [$($_.Exception.Message)"] -ForegroundColor Red -BackgroundColor White 
        $Message = $_.Exception.Message.Split(':')[1]
        Write-Warning "Message: $($Message)" 
        return

   }
   

}

if($PSVersionTable.PSVersion.Major -eq 2){

    ListOrFind-UserInActiveDirectoy |  ConvertTo-JSONP2

}else{

    ListOrFind-UserInActiveDirectoy |  ConvertTo-Json -Depth 10
}