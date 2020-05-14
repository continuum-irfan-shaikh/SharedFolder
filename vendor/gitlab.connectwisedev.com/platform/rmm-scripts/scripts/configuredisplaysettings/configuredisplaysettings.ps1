<#
Name : Configure display settings
category : Setup

$Remove_Windows_Welcome = "Yes"
$Remove_My_Documents_Desktop_Icon = "Yes"
$Enable_Num_Lock_on_Boot = "Yes"
$Wallpaper_File = "c:\temp\download.jpg"
$Registered_Owner = "Conti" 
$Registered_Company = "Continuum"
$Rename_My_Computer = "Comp"
$Rename_My_Network_places = "Net"
$Rename_My_Documents = "Doc"
$Rename_Recycle_bin= "Bin"
$Command_Prompt_Here = "Yes" 
$Tab_Auto_Complete = "Yes"
$Desktop_Cleanup_Wizard = "Yes"
$Remove_Shortcut_To_Prefix = "Yes"
$Sort_Start_Menu_by_Name = "Yes"

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}


$Computer = $env:COMPUTERNAME
$AllUsers = query user /server:$computer 2>&1

$Users = $AllUsers | ForEach-Object {(($_.trim() -replace ">" -replace "(?m)^([A-Za-z0-9]{3,})\s+(\d{1,2}\s+\w+)", 
'$1  none  $2' -replace "\s{2,}", "," -replace "none", $null))} | ConvertFrom-Csv
$CurrentUsers = @()
ForEach ($User in $Users)
{
  $CUser = ($user | ?{$_.state -ne 'Disc'} | Select-Object username).username
  $CurrentUsers+= $CUser
}
#$CurrentUsers

$SID = ((New-Object System.Security.Principal.NTAccount($CurrentUsers[0])).Translate([System.Security.Principal.SecurityIdentifier]).Value)

IF((Get-PSDrive).name -eq "HKCR"){} else {New-PSDrive -Name HKCR -PSProvider Registry -Root HKEY_CLASSES_ROOT | Out-Null -ErrorAction SilentlyContinue}
IF((Get-PSDrive).name -eq "HKU"){} else {New-PSDrive -Name HKU  -PSProvider Registry -Root HKEY_USERS | Out-Null -ErrorAction SilentlyContinue}

function ActiveUsers(){
[boolean]$flag=$false
IF($CurrentUsers -ne ''){
       $flag = $true
}
Return $flag
}


Function CreateModifyRegistryKey($RegistryPath,$RegistryKey, $RegistryValue, $PropertyType){
 #check registry is present for that option
# IF(Test-Path $RegistryPath) {} 
 

 #"HKCU:\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer"
 
 IF(!(Test-Path $RegistryPath -ErrorAction SilentlyContinue)){
 $key = $RegistryPath.split('\')[-1]
 $Path1 = $RegistryPath.replace($key,"")
 
    IF(!(Test-Path $Path1)){
    $key1 = $Path1.split('\')[-1]
    $path2 = $Path1.replace($key1,"")
    Set-Location $path2

    New-Item $key1 -ErrorAction SilentlyContinue | Out-Null
        }
  
 IF(Test-Path $Path1){
     Set-Location $Path1
       New-Item $key -ErrorAction SilentlyContinue | Out-Null
    }
}
$isRegistryPresent=(Get-ItemProperty $RegistryPath).psobject.properties | where {$_.name -eq $RegistryKey} -ErrorAction SilentlyContinue
           try
            {
             if($isRegistryPresent)
             {
               #Update Item Property Internet Option
              
                  Set-ItemProperty -Path $RegistryPath -Name $RegistryKey -Value $RegistryValue  -ErrorAction Stop  #-PropertyType "String" 
                  
             }
             else
             {
               #Create New Registry Entry
              
                  New-ItemProperty -Path $RegistryPath -Name $RegistryKey -Value $RegistryValue -PropertyType $PropertyType -ErrorAction Stop   #-PropertyType "String"
                 
              }
              return "Successfully updated"
            }#try end
             catch
             {
                $err=$Error[0].Exception.Message
                return $err
             }
}

Function RenameIcon($IconName,$IconValue){
  
         $Shell = new-object -comobject shell.application
         $Namespace = $Shell.Namespace($IconValue)
         $Namespace.self.name = $IconName
}
      #Initializing variable to store all result 
       $Global:Report=@()
Function CreateReport($TaskName,$status,$message){
                    $result = new-object PSObject
                    $result | add-member -membertype NoteProperty -name "Task Name" -value "$TaskName"
                    $result | add-member -membertype NoteProperty -name "Status" -value "$status"
                    $result | add-member -membertype NoteProperty -name "Description" -value "$message"
                    $Global:Report+=$result

                    
        }#End CreateReport


#Calling Function
Try
{
    #Check Any user is logged in or no
    $isUserLoggedIn=ActiveUsers
    if(!$isUserLoggedIn)
    {
        $Host.UI.WriteErrorLine("This script requires logon user and currently no user is logged in")
        EXIT
    }
    

      #Check input value for Remove_Windows_Welcome and Create/Update accordingly
       if($Remove_Windows_Welcome -eq "Yes")
       {
Set-Location HKLM:
        $actionResult= CreateModifyRegistryKey -RegistryPath "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -PropertyType DWord -RegistryKey "LogonType" -RegistryValue "0" 
         if($actionResult -contains "Successfully updated")
         {
           CreateReport -TaskName "Remove Windows Welcome" -status "Completed" -message "$actionResult"
         }
         else
         {
           CreateReport -TaskName "Remove Windows Welcome" -status "Failure" -message "$actionResult"
         }
       }
       elseif($Remove_Windows_Welcome -eq "No")
       {
Set-Location HKLM:

         $actionResult= CreateModifyRegistryKey -RegistryPath "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -RegistryKey "LogonType" -RegistryValue "1"
         if($actionResult -contains "Successfully updated")
         {
           CreateReport -TaskName "Remove Windows Welcome" -status "Completed" -message "$actionResult"
         }
         else
         {
           CreateReport -TaskName "Remove Windows Welcome" -status "Failure" -message "$actionResult"
         }
       }

       #Check input value for Remove_My_Documents_Desktop_Icon and Create/Update accordingly
       if($Remove_My_Documents_Desktop_Icon -eq "Yes")
       {
Set-Location HKLM:
         $actionResult= CreateModifyRegistryKey -RegistryPath "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -RegistryKey "LogonType" -RegistryValue "0"
         if($actionResult -contains "Successfully updated")
         {
           CreateReport -TaskName "Remove_My_Documents_Desktop_Icon" -status "Completed" -message "$actionResult"
         }
         else
         {
           CreateReport -TaskName "Remove 'My Documents Desktop Icon'" -status "Failure" -message "$actionResult"
         }
       }
       elseif($Remove_My_Documents_Desktop_Icon -eq "No")
       {
Set-Location HKLM:
         CreateModifyRegistryKey -RegistryPath "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -RegistryKey "LogonType" -RegistryValue "1"
             if($actionResult -contains "Successfully updated")
             {
               CreateReport -TaskName "Remove 'My Documents Desktop Icon'" -status "Completed" -message "$actionResult"
             }
             else
             {
               CreateReport -TaskName "Remove 'My Documents Desktop Icon'" -status "Failure" -message "$actionResult"
             }
       }
       
       #Check input value for Remove_My_Documents_Desktop_Icon and Create/Update accordingly
       if($Enable_Num_Lock_on_Boot -eq "Yes")
       {
#Set-Location HKU:


         $actionResult= CreateModifyRegistryKey -RegistryPath "Registry::HKU\.DEFAULT\Control Panel\Keyboard" -RegistryKey "InitialKeyboardIndicators" -RegistryValue "2"
           if($actionResult -contains "Successfully updated")
             {
               CreateReport -TaskName "Enable Num Lock on Boot" -status "Completed" -message "$actionResult"
             }
             else
             {
               CreateReport -TaskName "Enable Num Lock on Boot" -status "Failure" -message "$actionResult"
             }
       }
       elseif($Enable_Num_Lock_on_Boot -eq "No")
       {
         CreateModifyRegistryKey -RegistryPath "Registry::HKU\.DEFAULT\Control Panel\Keyboard" -RegistryKey "InitialKeyboardIndicators" -RegistryValue "0"
         if($actionResult -contains "Successfully updated")
             {
               CreateReport -TaskName "Enable Num Lock on Boot" -status "Completed" -message "$actionResult"
             }
             else
             {
               CreateReport -TaskName "Enable Num Lock on Boot" -status "Failure" -message "$actionResult"
             }
       }

       #Check input value for Wallpaper_File and Create/Update accordingly
       if($Wallpaper_File)
       {
         if(Test-Path $Wallpaper_File)
         {
           try
           { 
Set-ItemProperty -path "HKU:\$sid\Control Panel\Desktop" -name wallpaper -value $Wallpaper_File


              rundll32.exe user32.dll, UpdatePerUserSystemParameters 1, True 
              CreateReport -TaskName "Add Wallpaper" -status "Completed" -message "Successfully Updated"
            }
            catch
            {
              $err=$Error[0].Exception.Message
              CreateReport -TaskName "Add Wallpaper" -status "Failure" -message "$err"
            }
         }
         else
         {
             CreateReport -TaskName "Add Wallpaper" -status "Failure" -message "File path is not correct"
             
         }
         
       }
      
       #Check input value for Registered_Owner and Create/Update accordingly
       if($Registered_Owner)
       {
Set-Location HKLM:
         $actionResult= CreateModifyRegistryKey -RegistryPath "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion" -RegistryKey "RegisteredOwner" -RegistryValue $Registered_Owner
         
         if($actionResult -contains "Successfully updated")
             {
               CreateReport -TaskName "Rename Registered Owner" -status "Completed" -message "$actionResult"
             }
             else
             {
               CreateReport -TaskName "Rename Registered Owner" -status "Failure" -message "$actionResult"
             }
       }

       #Check input value for Registered_Company and Create/Update accordingly
       if($Registered_Company)
       {
Set-Location HKLM:
         $actionResult= CreateModifyRegistryKey -RegistryPath "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion" -RegistryKey "RegisteredOrganization " -RegistryValue $Registered_Company
         if($actionResult -contains "Successfully updated")
             {
               CreateReport -TaskName "Rename Registered Company" -status "Completed" -message "$actionResult"
             }
             else
             {
               CreateReport -TaskName "Rename Registered Company" -status "Failure" -message "$actionResult"
             }
           }

       #Check input value for Rename_My_Computer and Create/Update accordingly
       if($Rename_My_Computer)
       {
         RenameIcon -IconName $Rename_My_Computer -IconValue 17
          CreateReport -TaskName "Rename My Computer" -status "Completed" -message "Successfully Updated"
       }

       #Check input value for My Network Places and Create/Update accordingly
       if($Rename_My_Network_places)
       {
         RenameIcon -IconName $Rename_My_Network_places -IconValue 18
         CreateReport -TaskName "Rename My Network places" -status "Completed" -message "Successfully Updated"
       }

       #Check input value for My Documents and Create/Update accordingly
       if($Rename_My_Documents)
       {
         RenameIcon -IconName $Rename_My_Documents -IconValue 5
         CreateReport -TaskName "Rename My Documents" -status "Completed" -message "Successfully Updated"
       }

       #Check input value for My Documents and Create/Update accordingly
       if($Rename_Recycle_bin)
       {
         RenameIcon -IconName $Rename_Recycle_bin -IconValue 10
         CreateReport -TaskName "Rename Recycle bin" -status "Completed" -message "Successfully Updated"
       }

       #Check input value for Command_Prompt_Here and Create/Update accordingly
       #check registry is present for that option
        $isRegistryPresent=(Get-ItemProperty "Registry::HKEY_CLASSES_ROOT\Directory\Background\shell\cmd").psobject.properties | where {$_.name -eq "Extended"}

        if($isRegistryPresent)
        {
           try
            {
             if($Command_Prompt_Here -eq "Yes")
             {
               #Remove Extended property, this action will add 'open command prompt here' in desktop context menu
              
                  Remove-ItemProperty -Path "Registry::HKEY_CLASSES_ROOT\Directory\Background\shell\cmd" -Name "Extended" -Force -ErrorAction Stop
                  CreateReport -TaskName "Command Prompt Here" -status "Completed" -message "Successfully Updated"
             }
             elseif($Command_Prompt_Here -eq "No")
             {
                   CreateReport -TaskName "Command Prompt Here" -status "Completed" -message "Successfully Updated"
              }
            }#try end
             catch
             {
              $err=$Error[0].Exception.Message
              CreateReport -TaskName "Command Prompt Here" -status "Failure" -message "$err"
             }
        }
        else
        {
         try
            {
             if($Command_Prompt_Here -eq "Yes")
             {
               #Remove Extended property 
              
                  #Do Nothing
                  $Result+="Updated"
             }
             elseif($Command_Prompt_Here -eq "No")
             { 
                 #after adding Extended key, it will remove by default 'open command prompt here' from desktop context menu
                 New-ItemProperty -Path "Registry::HKEY_CLASSES_ROOT\Directory\Background\shell\cmd" -Name "Extended" -PropertyType "String" -ErrorAction Stop
                  $Result+="Updated"
                  CreateReport -TaskName "Command Prompt Here" -status "Completed" -message "Successfully Updated"
              }
            }#try end
             catch
             {
              $err=$Error[0].Exception.Message
              CreateReport -TaskName "Command Prompt Here" -status "Failure" -message "$err"
             }
        
        }#End for Command Prompt Here
      
      #Check input value for Tab_Auto_Complete and Create/Update accordingly
       if($Tab_Auto_Complete -eq "Yes")
       {
Set-Location HKLM:
         $actionResult= CreateModifyRegistryKey -RegistryPath "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Command Processor" -RegistryKey "CompletionChar" -RegistryValue "9"
         $actionResult= CreateModifyRegistryKey -RegistryPath "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Command Processor" -RegistryKey "PathCompletionChar" -RegistryValue "9"
         if($actionResult -contains "Successfully updated")
             {
               CreateReport -TaskName "Tab Auto Complete" -status "Completed" -message "$actionResult"
             }
             else
             {
               CreateReport -TaskName "Tab Auto Complete" -status "Failure" -message "$actionResult"
             }
        
       }
       elseif($Tab_Auto_Complete -eq "No")
       {
Set-Location HKLM:
         $actionResult=CreateModifyRegistryKey -RegistryPath "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Command Processor" -RegistryKey "CompletionChar" -RegistryValue "40"
         $actionResult=CreateModifyRegistryKey -RegistryPath "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Command Processor" -RegistryKey "PathCompletionChar" -RegistryValue "40"
         if($actionResult -contains "Successfully updated")
             {
               CreateReport -TaskName "Tab Auto Complete" -status "Completed" -message "$actionResult"
             }
             else
             {
               CreateReport -TaskName "Tab Auto Complete" -status "Failure" -message "$actionResult"
             }
       }

       #Check input value for Desktop_Cleanup_Wizard and Create/Update accordingly
       if($Desktop_Cleanup_Wizard -eq "Yes")
       {
         $actionResult= CreateModifyRegistryKey -RegistryPath "HKU:\$SID\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer" -RegistryKey "NoDesktopCleanupWizard" -RegistryValue "40" 
         if($actionResult -contains "Successfully updated")
             {
               CreateReport -TaskName "Desktop Cleanup Wizard" -status "Completed" -message "$actionResult"
             }
             else
             {
               CreateReport -TaskName "Desktop Cleanup Wizard" -status "Failure" -message "$actionResult"
             }
         
       }
       elseif($Desktop_Cleanup_Wizard -eq "No")
       {
         $actionResult=CreateModifyRegistryKey -RegistryPath "HKU:\$sid\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer" -RegistryKey "NoDesktopCleanupWizard" -RegistryValue "9"
         if($actionResult -contains "Successfully updated")
             {
               CreateReport -TaskName "Desktop Cleanup Wizard" -status "Completed" -message "$actionResult"
             }
             else
             {
               CreateReport -TaskName "Desktop Cleanup Wizard" -status "Failure" -message "$actionResult"
             }
         
         
       }

       #Check input value for Remove_Shortcut_To_Prefix and Create/Update accordingly
       #check registry is present for that option
        $isRegistryPresent=(Get-ItemProperty "HKU:\$SID\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer").psobject.properties | where {$_.name -eq "link"}

        if($isRegistryPresent)
        {
           try
            {
             if($Remove_Shortcut_To_Prefix -eq "Yes")
             {
               
                 [byte[]]$byte=0,0,0,0
                  Set-ItemProperty -Path "HKU:\$SID\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer" -Name "link" -Value $byte -ErrorAction Stop  # -PropertyType "Binary" 
                  CreateReport -TaskName "Remove Shortcut To Prefix" -status "Completed" -message "Successfully Updated"
             }
             elseif($Remove_Shortcut_To_Prefix -eq "No")
             {
               
                  Remove-ItemProperty -Path "HKU:\$SID\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer" -Name "link" -Force -ErrorAction Stop
                  CreateReport -TaskName "Remove Shortcut To Prefix" -status "Completed" -message "Successfully Updated"
              }
            }#try end
             catch
             {
              $err=$Error[0].Exception.Message
              CreateReport -TaskName "Remove Shortcut To Prefix" -status "Failure" -message "$err"
             
             }
        }
        else
        {
         try
            {
             if($Remove_Shortcut_To_Prefix -eq "Yes")
             {
                  #Create New link property 
                  [byte[]]$byte=0,0,0,0
                  New-ItemProperty -Path "HKU:\$SID\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer" -Name "link" -Value $byte -PropertyType "Binary" -ErrorAction Stop
                  CreateReport -TaskName "Remove Shortcut To Prefix" -status "Completed" -message "Successfully Updated"
             }
             elseif($Remove_Shortcut_To_Prefix -eq "No")
             {
                 #Do nothing
                  $Result+="Updated"
              }
            }#try end
             catch
             {
              $err=$Error[0].Exception.Message
               CreateReport -TaskName "Remove Shortcut To Prefix" -status "Failure" -message "$err"
             
             }
        
        }#End for Remove_Shortcut_To_Prefix

        #Check input value for Sort_Start_Menu_by_Name and Create/Update accordingly
       if($Sort_Start_Menu_by_Name -eq "Yes")
       {
         $actionResult=CreateModifyRegistryKey -RegistryPath "HKU:\$SID\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer" -RegistryKey "NoStrCmpLogical" -RegistryValue "1"
         if($actionResult -contains "Successfully updated")
             {
               CreateReport -TaskName "Sort_Start_Menu_by_Name" -status "Completed" -message "$actionResult"
             }
             else
             {
               CreateReport -TaskName "Sort_Start_Menu_by_Name" -status "Failure" -message "$actionResult"
             }
         
       }
       elseif($Sort_Start_Menu_by_Name -eq "No")
       {
         $actionResult=CreateModifyRegistryKey -RegistryPath "HKU:\$SID\Software\Microsoft\Windows\CurrentVersion\Policies\Explorer" -RegistryKey "NoStrCmpLogical" -RegistryValue "0"
         if($actionResult -contains "Successfully updated")
             {
               CreateReport -TaskName "Sort Start Menu by Name" -status "Completed" -message "$actionResult"
             }
             else
             {
               CreateReport -TaskName "Sort Start Menu by Name" -status "Failure" -message "$actionResult"
             }
         
       } Set-Location c:

       
$Global:Report | ft -AutoSize
}
Catch
{
    $Host.UI.WriteErrorLine($Error[0].Exception.Message)
}
