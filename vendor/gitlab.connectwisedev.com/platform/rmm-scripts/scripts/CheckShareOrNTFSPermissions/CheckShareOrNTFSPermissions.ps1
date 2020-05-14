
Clear-Host

<#

[string]$Path = "D:\temp2"        
[string]$Username  = ""         # INFICS\\kumarg  
[bool]$SharePermissions = $true
[bool]$NTFSPermissions  = $true
sss
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


 
function Get-SharePermissions{ 

[cmdletbinding( 
    DefaultParameterSetName = 'computer', 
    ConfirmImpact = 'low' 
)] 
    Param( 
        [Parameter( 
            Mandatory = $True, 
            Position = 0, 
            ParameterSetName = 'computer', 
            ValueFromPipeline = $True)] 
            [array]$computer                       
            ) 
Begin {                 
    #Process Share report 
    $sharereport = @() 
    } 
Process { 
    #Iterate through comptuers 
    ForEach ($c in $computer) { 
        Try {     
            Write-Verbose "Computer: $($c)" 
            #Retrieve share information from comptuer 
            $ShareSec = Get-WmiObject -Class Win32_LogicalShareSecuritySetting -ComputerName $c -ea stop 
            ForEach ($Shares in $sharesec) { 
                Write-Verbose "Share: $($Shares.name)" 
                    #Try to get the security descriptor 
                    $SecurityDescriptor = $ShareS.GetSecurityDescriptor() 
                    #Iterate through each descriptor 
                    ForEach ($DACL in $SecurityDescriptor.Descriptor.DACL) { 

                        $arrshare = New-Object PSObject 
                        $arrshare | Add-Member NoteProperty Computer $c 
                        $arrshare | Add-Member NoteProperty Name $Shares.Name 
                        $arrshare | Add-Member NoteProperty ID $DACL.Trustee.Name 

                        #Convert the current output into something more readable 
                        Switch ($DACL.AccessMask) { 
                            2032127 {$AccessMask = "FullControl"} 
                            1179785 {$AccessMask = "Read"} 
                            1180063 {$AccessMask = "Read, Write"} 
                            1179817 {$AccessMask = "ReadAndExecute"} 
                            -1610612736 {$AccessMask = "ReadAndExecuteExtended"} 
                            1245631 {$AccessMask = "ReadAndExecute, Modify, Write"} 
                            1180095 {$AccessMask = "ReadAndExecute, Write"} 
                            268435456 {$AccessMask = "FullControl (Sub Only)"} 
                            default {$AccessMask = $DACL.AccessMask} 
                            } 
                        $arrshare | Add-Member NoteProperty AccessMask $AccessMask 
                        #Convert the current output into something more readable 
                        Switch ($DACL.AceType) { 
                            0 {$AceType = "Allow"} 
                            1 {$AceType = "Deny"} 
                            2 {$AceType = "Audit"} 
                            } 
                        $arrshare | Add-Member NoteProperty AceType $AceType 
                        #Add to existing array 
                        $sharereport += $arrshare 
                        } 
                    } 
                } 
            #Catch any errors                 
            Catch { 
                $arrshare | Add-Member NoteProperty Computer $c 
                $arrshare | Add-Member NoteProperty Name "NA" 
                $arrshare | Add-Member NoteProperty ID "NA"  
                $arrshare | Add-Member NoteProperty AccessMask "NA"           
                }  
            Finally { 
                #Add to existing array 
                $sharereport += $arrshare 
                }                                                    
            }  
        }                        
End { 
        
        $SharePermission = @()
        $ShareArr = @()

         if(([uri]$Path).IsUnc){
              
             $FileOrFolder = ($Path -split '\\')[-1]
            
          }Else{

             $FileOrFolder = split-path $Path -leaf -resolve
          }

        $SharePermission = $Sharereport | Where-Object {$_.Name -eq "$FileOrFolder"}

        if( $SharePermission.Count -gt 0 ){

            foreach($SP in $SharePermission){

               $SP.ID +' '+ '(' + $SP.AceType +')' +' '+ $SP.AccessMask
            }
                       
        }Else{

           # Write-Warning "File or Folder : $($FileOrFolder) not found"
        }

    } # End
}  # function   
 
function Get-ShareNTFSPermissions{ 

[cmdletbinding( 
    DefaultParameterSetName = 'computer', 
    ConfirmImpact = 'low' 
)] 
    Param( 
        [Parameter( 
            Mandatory = $True, 
            Position = 0, 
            ParameterSetName = 'computer', 
            ValueFromPipeline = $True)] 
            [array]$computer                       
            )   
Begin {               
    #Process NTFS Share report                 
    $ntfsreport = @()       
    } 
Process { 
$arrntfs = @()
    #Iterate through each computer 
    ForEach ($c in $computer) {  
        Try {                  
            Write-Verbose "Computer: $($c)"  
            #Gather share information 
            $shares = Gwmi -comp $c Win32_Share -ea stop | ? {$_.Name -ne 'ADMIN$'-AND $_.Name -ne 'C$' -AND $_.Name -ne 'IPC$'} | Select Name,Path 
            ForEach ($share in $shares) { 
                #Iterate through shares 
                Write-Verbose "Share: $($share.name)" 
                If ($share.path -ne "") { 
                    #Retrieve ACL information from Share   
                    $remoteshare = $share.path -replace ":","$"  
                    Try { 
                        #Gather NTFS security information from each share 
                        $acls = Get-ACL "\\$computer\$remoteshare" 
                        #iterate through each ACL 
                        ForEach ($acl in $acls.access) {
                           $AccessMask = $acl.FileSystemRights 

                           $arrntfs += New-Object PSObject -Property @{                             
                           Computer = $c               
                           ShareName = $Share.name 
                           Path =  $share.path 
                           NTFS_User = $acl.IdentityReference 
                           NTFS_Rights = $AccessMask 
                           NTFS_ControlType = $acl.AccessControlType
                           NTFS_IsInherited = $acl.IsInherited
                            }
                            #$arrntfs | select ShareName,Path,NTFS_User,NTFS_Rights,NTFS_ControlType
                        }  # end foreach    
                        } 
                    Catch { 
                        $arrntfs = New-Object PSObject                     
                        #Process NTFS Report          
                        $arrntfs | Add-Member NoteProperty Computer $c               
                        $arrntfs | Add-Member NoteProperty ShareName "NA" 
                        $arrntfs | Add-Member NoteProperty Path "NA" 
                        $arrntfs | Add-Member NoteProperty NTFS_User "NA" 
                        $arrntfs | Add-Member NoteProperty NTFS_Rights "NA" 
                        $arrntfs | Add-Member NoteProperty NTFS_ControlType "NA" 
                        $arrntfs | Add-Member NoteProperty NTFS_IsInherited "NA"                   
                        } 
                    Finally { 
                        #Add to existing array 

                        $ntfsreport = $arrntfs  
                        }                                                                                    
                    }                                
                } 
            } 
        Catch { 
            $arrntfs | Add-Member NoteProperty Computer $c               
            $arrntfs | Add-Member NoteProperty ShareName "NA"  
            $arrntfs | Add-Member NoteProperty Path "NA"  
            $arrntfs | Add-Member NoteProperty NTFS_User "NA"  
            $arrntfs | Add-Member NoteProperty NTFS_Rights "NA" 
            $arrntfs | Add-Member NoteProperty NTFS_ControlType "NA" 
            $arrntfs | Add-Member NoteProperty NTFS_IsInherited "NA"            
            } 
        Finally { 
            #Add to existing array 
            $ntfsreport = $arrntfs          
            }                                         
        }             
    } 
End { 
        
    $NTFSPermission = @()
    $NTFSArr = @()
     
   #  $ntfsreport

    # Filter without user
    if([string]::IsNullOrEmpty($Username)){
      
      
          if(([uri]$Path).IsUnc){
            
            $FileOrFolder = ($Path -split '\\')[-1]
            $FileOrFolder = $FileOrFolder
            $NTFSPermission = $ntfsreport | Where-Object {$_.ShareName -eq $FileOrFolder} 

          }Else{

            $NTFSPermission = $ntfsreport | Where-Object {$_.Path -eq $Path} 
          }


     
        if($NTFSPermission){

            if($NTFS.NTFS_IsInherited){   
                         
                foreach($NTFS in $NTFSPermission){
                              
                  $NTFSArr += '[' + "IsInherited" + ']'+  ($NTFS.NTFS_User).ToString() +' '+ '(' + $NTFS.NTFS_ControlType +')' +' '+ $NTFS.NTFS_Rights

                }
            }else{
               foreach($NTFS in $NTFSPermission){
                        
                 $NTFSArr += ($NTFS.NTFS_User).ToString() +' '+ '(' + $NTFS.NTFS_ControlType +')' +' '+ $NTFS.NTFS_Rights
                }
            }

            return $NTFSArr

        }else{
          
         # Write-Warning "Specified folder not found"

        }
    }

    # Filter with user
    if(-not [string]::IsNullOrEmpty($Username)){

        
          if(([uri]$Path).IsUnc){
            
            $FileOrFolder = ($Path -split '\\')[-1]
            $FileOrFolder = $FileOrFolder
            $NTFSPermission = $ntfsreport | Where-Object {($_.ShareName -eq $FileOrFolder) -and ($_.NTFS_User -eq $Username)}

          }Else{

            $NTFSPermission = $ntfsreport | Where-Object {($_.Path -eq $Path) -and ($_.NTFS_User -eq $Username)}
          }
               
       
        if($NTFSPermission){

            if($NTFS.NTFS_IsInherited){
               
                foreach($NTFS in $NTFSPermission){
            
                $NTFSArr += '[' + "IsInherited" + ']'+  ($NTFS.NTFS_User).ToString() +' '+ '(' + $NTFS.NTFS_ControlType +')' +' '+ $NTFS.NTFS_Rights
                }
            }else{
               
                foreach($NTFS in $NTFSPermission){
            
                 $NTFSArr +=  ($NTFS.NTFS_User).ToString() +' '+ '(' + $NTFS.NTFS_ControlType +')' +' '+ $NTFS.NTFS_Rights
                }
            }
        }else{
          
         # Write-Warning "Specified user not found"

        }

    }


    }   # End           
} # End function
 

function Get-shareNTFSPermission{
   
   $SharedNTFSObject = New-Object -TypeName psobject
   $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name TaskName -Value "Check shared or NTFS permission"

   $Comp = $env:COMPUTERNAME

   $hostVerSionMajor = ($PSVersionTable.PSVersion.Major).ToString()
   $hostVerSionMinor = ($PSVersionTable.PSVersion.Minor).ToString()
   $hostVersion = $hostVerSionMajor +'.'+ $hostVerSionMinor 

   $osVersionMajor = ([System.Environment]::OSVersion.Version.major).ToString()
   $osVersionMinor = ([System.Environment]::OSVersion.Version.minor).ToString()
   $osVersion = $osVersionMajor +'.'+ $osVersionMinor

   [boolean]$isPsVersionOk = ([version]$hostVersion -ge [version]'2.0')
   [boolean]$isOSVersionOk = ([version]$osVersion -ge [version]'6.0')
        
  
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
    
         $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
         $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
         $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
         $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell version below 2.0 is not supported"
         $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

         return $SharedNTFSObject

    }

  
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
    
     $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
     $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
     $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
     $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Result -Value "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
     $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

     return $SharedNTFSObject

    }

   #-------------------------------------------------
   # Verify UNC or absolute path
   #-------------------------------------------------

   if(([uri]$Path).IsUnc){
    
     #Write-Host "UNC Path : $($Path)"

     if(-not($Path -match '^(\\)(\\[A-Za-z0-9-_.]+){2,2}(\\?)$')){
        
          $StdErrArr = @()
          $stdOutArr = @()
          $StdErr = New-Object PSObject -Property @{		       
		                      id = 0;
                              title =  "Entered path is not correct";
                              detail = "The path should be of the form \\server\share"

                   }
         
         $StdErrArr += $StdErr
         $stdOutArr += ("Entered path is not correct")
    
         $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
         $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
         $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
         $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Result -Value "Entered path is not correct"
         $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

         return $SharedNTFSObject
     } 
   }else{

     
      $driveletter = $Path.Split(':')[0] 
      $DriveLetterRange = 'A-Z'
      $driveletter = $driveletter.ToUpper();

      if (-not($driveletter -notmatch '[A-Z](-[A-Z])?') -and ($driveletter.Length -ne 1)) {

              $StdErrArr = @()
              $stdOutArr = @()
              $StdErr = New-Object PSObject -Property @{		       
		                          id = 0;
                                  title =  "Entered drive not correct";
                                  detail = "Drive Letter $driveletter is not in the range of A-Z"

                       }
         
             $StdErrArr += $StdErr
             $stdOutArr += ("Entered drive not correct")
    
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Result -Value "Entered drive not correct"
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr   

             return $SharedNTFSObject
      }
         
      $IsDriveLetterExist = Get-WmiObject -Class Win32_logicaldisk| Where-Object{ $_.DeviceID -eq ($driveletter+':') }

      if(-not $IsDriveLetterExist){

                             
              $StdErrArr = @()
              $stdOutArr = @()

              $StdErr = New-Object PSObject -Property @{		       
		                          id = 0;
                                  title =  "Error: Volume $driveletter does not exist";
                                  detail = "List of all avaliable drives"

                       }
             
             $Drives = Get-WmiObject -Class Win32_logicaldisk | select DeviceID,VolumeName

             $SharedNTFSObj =  New-Object psobject
             $SharedNTFSObj | Add-Member -MemberType NoteProperty -Name Drives -Value $Drives
         
             $StdErrArr += $StdErr
             $stdOutArr += ("Error: Volume $driveletter does not exist")
    
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error: Volume $driveletter does not exist"
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $SharedNTFSObj     

             return $SharedNTFSObject

       }

      if(-not (Test-Path -Path $Path)){

        
              $StdErrArr = @()
              $stdOutArr = @()

              $StdErr = New-Object PSObject -Property @{		       
		                          id = 0;
                                  title =  "Incorrect Path";
                                  detail = "Error: Folder $Path does not exist";

                       }
             
                      
             $StdErrArr += $StdErr
             $stdOutArr += ("Incorrect Path")
    
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error: Incorrect Path"
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr
            
             return $SharedNTFSObject

      }
   } # end path
   
  $Comp = ""
  
  if(([uri]$Path).IsUnc){
     $uri = new-object System.Uri($Path)
     $HostName = $uri.host
     $Comp = $HostName

  }Else{

    $Comp = $env:COMPUTERNAME
  }
  
  if($SharePermissions -and $NTFSPermissions){

    $Shared =  Get-SharePermissions -computer $Comp
    $NTFS = Get-ShareNTFSPermissions -computer $Comp

    $SharedObject =  New-Object psobject
    $SharedObject |  Add-Member -MemberType NoteProperty -Name SharedPermissions -Value $Shared
    $SharedObject |  Add-Member -MemberType NoteProperty -Name NTFSPermissions -Value $NTFS

    $arrStr = @()

    $arrStr += $Shared
    $arrStr += $NTFS
              
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : Shared and NTFS permissions"  
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $arrStr
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $SharedObject
                          
    return $SharedNTFSObject

  }

  if($SharePermissions -and (-not $NTFSPermissions)){
    
    $Shared = Get-SharePermissions -computer $Comp

    $SharedObject =  New-Object psobject
    $SharedObject |  Add-Member -MemberType NoteProperty -Name SharedPermissions -Value $Shared
   
    $arrStr = @()

    $arrStr += $Shared
    
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : Shared permissions"  
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $arrStr
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $SharedObject
                          
    return $SharedNTFSObject
  }

  if(( -not $SharePermissions) -and $NTFSPermissions){
    
    
    $NTFS = Get-ShareNTFSPermissions -computer $Comp

    $SharedObject =  New-Object psobject
    $SharedObject |  Add-Member -MemberType NoteProperty -Name NTFSPermissions -Value $NTFS
   
    $arrStr = @()

    $arrStr += $NTFS
    
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Status -Value "Success" 
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Code -Value 0
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Result -Value "Success : NTFS permissions"  
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $arrStr
    $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name dataObject -Value $SharedObject
                          
    return $SharedNTFSObject
  }

  if((-not $SharePermissions) -and (-not $NTFSPermissions)){
    
              $StdErrArr = @()
              $stdOutArr = @()

              $StdErr = New-Object PSObject -Property @{		       
		                          id = 0;
                                  title =  "Please select atleast one permission";
                                  detail = "Please select atleast one permission";

                       }
             
                      
             $StdErrArr += $StdErr
             $stdOutArr += ("Please select atleast one permission")
    
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Status -Value "fail" 
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Code -Value 1
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stderr -Value $StdErrArr
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name Result -Value "Error: Please select atleast one permission";
             $SharedNTFSObject | Add-Member -MemberType NoteProperty -Name stdout -Value  $stdOutArr
            
             return $SharedNTFSObject    
  }

}


if($PSVersionTable.PSVersion.Major -eq 2){

    Get-shareNTFSPermission |  ConvertTo-JSONP2

}else{

    Get-shareNTFSPermission |  ConvertTo-Json -Depth 10
}
