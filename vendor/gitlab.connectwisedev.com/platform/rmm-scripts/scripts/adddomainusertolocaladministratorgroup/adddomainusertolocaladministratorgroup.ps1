<#  
.SYNOPSIS  
    Add Domain User to Local Administrators Group
.DESCRIPTION  
    Add Domain User to Local Administrators Group. Scritp will check correct Domain Name. Script will validate the domain user and also check if domain user already a member of local administrators group.
.NOTES  
    File Name  : AddDomainUserToLocalAdministatorGroup.ps1
    Author     : Durgeshkumar Patel  
    Requires   : PowerShell V2 or greater.   
.PARAMETERS
    
.HELP
#>

<# Variables to define in JSON Schema
$DomainUser              Example "smith7"  #"SamAccountName"
$DomainName              Example  "Scriptdc3.LOCAL"   #Full Domain Name
#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}


$Computer = $env:computername
$LocalGroup = "Administrators"
function domaincheck {
    
    if (!$DomainName) {
        Write-Output "`nKindly provide the Domain Name. Provided value is null"
        Exit;
    }
    else {
        if ((Get-WmiObject -class Win32_ComputerSystem).domain -eq $DomainName) {
            return $true
        }
        else {
            return $false
        }
    }
}


function validate {
    
    $ErrorActionPreference = "SilentlyContinue"
    net user $DomainUser /domain > null 2>&1
    if ($?) {
        return $true
    }
    else {
        return $false
    }
}

function member {

    $member = net localgroup administrators | Where {$_ -eq $DomainUser}
    if ($member -ne $null) {
        return $true
    }
    else {
        return $false
    }
}

if (!$DomainUser) {
    Write-Output "`nKindly provide the Username. Provided value is null"
    Exit;
}

if (!(domaincheck)) {
    "Provided Domain Name '{0}' is not correct" -f $DomainName
    Exit;
}

try {

    if (member) {
        Write-Output "`n$DomainUser is already a member of Local Administrators Group"
        Exit;
    }
    else {
        if (validate) {
            $target = ([ADSI]"WinNT://$Computer/$LocalGroup,group")
            $target.psbase.Invoke("Add", ([ADSI]"WinNT://$DomainName/$DomainUser").path)
            if (member) {
                Write-Output "`n$DomainUser added to Local Administrators group"
            }
            else {
                Write-Output "`n$DomainUser not added to Local Administrators group"
            }
        }
        else {
            Write-Output "`n$DomainUser is not a domain user"
        }
    }
}
catch {
    Write-Output "`nFailed to add $DomainUser to Local Administrators group"
    Write-Output $_.Exception.Message
}


