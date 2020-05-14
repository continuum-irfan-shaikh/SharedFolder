<#
    .SYNOPSIS
         Configure file/folder permissions
    .DESCRIPTION
         Configure file/folder permissions
    .Help
         
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>
<#

#Text Box
Path
Variable;- $path

#Text Box   
User/Group Name
Variable:- $UGName

#Drop Down
Varialbe:- $Action 
Values   Append Overwrite Revoke

#Check Boxes. True or False values
FullControl
Modify  
Write  
ReadAndExecute  
Read  
    
#Variables:-
$FullControl
$Modify  
$Write  
$ReadAndExecute  
$Read

#Radio Buttons 
    Allow  $true/$false 
    Deny   $false/$false
Variable :- $PermissionType

Examples:-
$path = 'c:\durgesh\test\'
$UGName = "Durgesh"
$FullControl = $false
$Modify = $false
$Read = $false
$Write = $true
$ReadAndExecute = $true

$PermissionType = "Allow"   #Allow OR Deny
$Action = 'Revoke'   #Append or Overwrite or Revoke 

#>
#########################################


if (!(Test-Path $path)) {
    Write-Error "'$path' not a correct folder path/file. Kindly check and try again"
    Exit;
}

if ($PermissionType -eq $null) {
    Write-Error "Select allow/deny permission type"
    Exit;
}

if (($FullControl -eq $false) -and ($Modify -eq $false) -and ($Read -eq $false) -and ($Write -eq $false) -and ($ReadAndExecute -eq $false)) {
    Write-Error "Permission not selected. Kindly select permission."
    Exit;
}
 
if ((Get-Item $path) -is [System.IO.DirectoryInfo]) {
    [system.security.accesscontrol.InheritanceFlags] $inherit = "ContainerInherit, ObjectInherit"
    [system.security.accesscontrol.PropagationFlags] $propagation = "None"
}
else {
    [system.security.accesscontrol.InheritanceFlags] $inherit = "None"
    [system.security.accesscontrol.PropagationFlags] $propagation = "None"
}

if ($Read -eq $true) {
    $Permission = "Read"
}
if ($write -eq $true) {
    $Permission = "write"
}
if (($ReadAndExecute -eq $true) -or (($Read -eq $true) -and ($ReadAndExecute -eq $true))) {
    $Permission = "ReadAndExecute"
}
if (($Read -eq $true) -and ($Write -eq $true)) {
    $Permission = "Write, Read"
}
if ((($ReadAndExecute -eq $true) -and ($write -eq $true)) -or (($Read -eq $true) -and ($Write -eq $true) -and ($ReadAndExecute -eq $true))) {
    $Permission = "Write, ReadAndExecute"
}
if (($Modify -eq $true) -or (($Modify -eq $true) -and ($Read -eq $true) -and ($Write -eq $true) -and ($ReadAndExecute -eq $true))) {
    $Permission = "Modify"
}
if (($FullControl -eq $true) -or (($FullControl -eq $true) -and ($Modify -eq $true) -and ($Read -eq $true) -and ($Write -eq $true) -and ($ReadAndExecute -eq $true)) -or (($FullControl -eq $true) -and ($Read -eq $true))) {
    $Permission = "FullControl"
}

function UserValidation {

    $ErrorActionPreference = 'SilentlyContinue'
    $user = net user "$UGName"  
    if ($user) {
        return $true
    }
    else {
        return $false
    }   
}

function GroupValidation {

    $ErrorActionPreference = 'SilentlyContinue'
    $group = net group "$UGName" 
    if ($group) {
        return $true
    }
    else {
        return $false
    }   
}

function Overwrite-Permissions {
     
    $Acl = (Get-Item $path).GetAccessControl('Access')
    $Ar = New-Object System.Security.AccessControl.FileSystemAccessRule("$UGName", "$Permission", "$inherit", "$propagation", "$PermissionType")
    #Overwrite  SetAccessRule
    $Acl.SetAccessRule($Ar)
    Set-Acl -AclObject $Acl -Path $path  

    if ($?) {
        return $true
    }
    else {
        return $false
    }
}

function RevokeOldPermissions {
    
    #This function is used with Overwrite-Permissions   
    $Acl = (Get-Item $path).GetAccessControl('Access')  #| fl
    $ptypes = "Allow", "Deny"
    foreach ($ptype in $ptypes) {
        
        $Ar = New-Object system.security.AccessControl.FileSystemAccessRule("$UGName", "$Permission", "$inherit", "$propagation", "$ptype")
  
        #Revoke RemoveAccessRule ALL (Tested if any permissions there)   
        $Acl.RemoveAccessRuleAll($Ar) 
        Set-Acl $path $Acl 
    }
     
    if ((Get-Item $path).GetAccessControl('Access') | where { $_.IdentityReference -match "$UGName" }) {
        return $false
    }
    else {
        return $true
    }
}

function Append-Permissions {
           
    $Acl = (Get-Item $path).GetAccessControl('Access')
    $Ar = New-Object system.security.AccessControl.FileSystemAccessRule("$UGName", "$Permission", "$inherit", "$propagation", "$PermissionType")
    #Add/Append  AddAccessRule
    $Acl.AddAccessRule($Ar)
    Set-Acl $path $Acl

    if ($?) {
        return $true
    }
    else {
        return $false
    }
}

function Revoke-Permissions {
    
    $Acl = (Get-Item $path).GetAccessControl('Access') 
    $Ar = New-Object System.Security.AccessControl.FileSystemAccessRule("$UGName", "$Permission", "$inherit", "$propagation", "$PermissionType")
    #Revoke RemoveAccessRule    
    $Acl.RemoveAccessRule($Ar)
    Set-Acl $path $Acl

    if ($?) {
        return $true
    }
    else {
        return $false
    }
}

function Inheritance {
    $Acl = (Get-Item $path).GetAccessControl('Access')
    $isInherited = $Acl.Access | Where-Object { $_.IdentityReference -match $UGName } | Where-Object { $_.isInherited -eq $true -and $_.AccessControlType -eq $PermissionType }  
    if ($isInherited) {
        Write-Output "Can not change inherited permission for $UGName."
        Exit;   
    }
}

function Permission_Check {
    $Acl = (Get-Item $path).GetAccessControl('Access')
    $Permission_Check = $Acl.Access | Where-Object { $_.IdentityReference -match $UGName }
    if (!$Permission_Check) {
        Write-Output "Can't revoke '$Permission' permission as no permissions were configured for '$UGName' for '$path'"
        Exit;   
    }
}

try {

    if ((UserValidation) -or (GroupValidation)) {
      
        switch ($Action) {
        
            "Append" {
                
                if (Append-Permissions) {
                    Write-Output "'$Permission ($PermissionType)' Permission added for $UGName"
                }
                else {
                    Write-Output "Adding '$Permission ($PermissionType)' permission failed for $UGName" 
                } break;
            }  
            
            "Overwrite" {
                
                if (RevokeOldPermissions) {
                    start-sleep 1
                    if (Overwrite-Permissions) {
                        Write-Output "Overwritten '$Permission ($PermissionType)' permission for $UGName"
                    }
                    else {
                        Write-Output "Overwriting '$Permission ($PermissionType)' permissions failed for $UGName" 
                    } 
                }
                else {
                    Write-Output "Overwriting '$Permission ($PermissionType)' permissions failed for $UGName"
                } break;
            }

            "Revoke" {

                #Check for existing permissions
                Permission_Check
                #Check for inheritance
                Inheritance
            
                if (Revoke-Permissions) {
                    Write-Output "Revoked '$Permission ($PermissionType)' permissions for $UGName "
                }
                else {
                    Write-Output "Revoking '$Permission ($PermissionType)' permissions failed for $UGName" 
                } break;
            }  
        }  
    }
    else {
        Write-Output "User/Group:- $UGName doesn't exist on system $ENV:COMPUTERNAME"
        EXIT;
    }
}
catch {
    Write-Error "Applying '$Permission ($PermissionType)' Permission Failed for '$UGName' on $ENV:COMPUTERNAME"
    Write-Error $_.Exception.Message
} 
