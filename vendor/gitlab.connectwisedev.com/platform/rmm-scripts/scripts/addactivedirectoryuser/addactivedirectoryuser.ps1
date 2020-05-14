# $UserLogonName = 'Will.smith'
# $Password = 'Test@123'
# $Domain = 'scriptdc1.local'
# $OrganizationalUnit = 'OU=TestOU,DC=scriptdc1,DC=local'
# $Group = 'Enterprise Admins'
# $FirstName = 'Will'
# $Initial = 'NS'
# $LastName = 'Smith'
# $DisplayName = 'Will smith'
# $Description = 'PowerSHell Developer'
# $TelephoneNumbers = ''
# $EMail = 'Will@gmail.com'
# $WebPage = 'https://technet.microsoft.com/en-gb/'
# $Street = 'Ghatkopar'
# $POBox = ''
# $City = 'New Delhi'
# $State = 'Maharashtra'
# $Zip = '400052'
# $Country = 'IN'
# $Phone = '2222222'
# $Mobile = '8800744559'
# $Fax = ''
# $IPPhone = '123123123123'
# $Notes = 'test note'
# $Jobtitle = 'Tech Support'
# $Department = 'GRT'
# $Company = 'Continuum'

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

try{

$ErrorActionPreference = 'Stop'

if(!(Get-Module -Name "ActiveDirectory")){
    if(!(Import-Module ActiveDirectory -PassThru)){
        Write-Error "Unable to import the module [ActiveDirectory]"
    }
}

if((Get-WmiObject -Class Win32_OperatingSystem).producttype -ne 2){
    Write-Output "`nScript can only run on Domain controllers, no action is performed."
    exit
}

#region verify-the-user-arguments
if (-not (Get-ADDomain -Identity $Domain)) {
    Write-Error "Unable to contact the Domain [$Domain], please verify the connectivty to Active Directory."
}
if (Get-ADUser -Filter "Name -eq '$UserLogonName'") {
    Write-Error "The User Logon Name [$UserLogonName] is not available, and already exists. Please try again with a different Name."
}
if (-not (Get-ADOrganizationalUnit -Filter "DistinguishedName -eq '$OrganizationalUnit'")) {
    Write-Error "The user Organizational Unit [$($OrganizationalUnit)] does not exist. Can't add a user there."    
}
if (-not (Get-ADGroup -Filter "Name -eq '$Group'")) {
    Write-Error "The group [$($Department)] does not exist. Can't add the user to this group."    
}
#endregion verify-the-user-arguments

#region build-attributes-for-splatting
$Attributes = @{  
    Enabled           = $true 
    SAMAccountName    = $UserLogonName
    userPrincipalName = $UserLogonName
    AccountPassword   = $Password | ConvertTo-SecureString -AsPlainText -Force
    Server            = $Domain
    Path              = $OrganizationalUnit
    GivenName         = $FirstName
    Initials          = $Initial
    SurName           = $LastName
    Name              = $(if($DisplayName){$DisplayName}else{"$FirstName $LastName".trim()})
    Description       = $Description
    Office            = $office
    EmailAddress      = $EMail
    HomePage          = $WebPage
    StreetAddress     = $Street
    POBox             = $POBox
    City              = $City
    State             = $State
    PostalCode        = $Zip 
    Country           = $Country
    HomePhone         = $Phone
    MobilePhone       = $Mobile
    Fax               = $Fax
    OfficePhone       = $TelephoneNumbers
    Title             = $Jobtitle
    Department        = $Department
    Company           = $Company
}
#endregion build-attributes-for-splatting

$user = New-ADUser @Attributes -PassThru

if ($user) {
    if ($Group) {
        Add-ADGroupMember -Identity $Group -Members $User.SamAccountName
    }
    if ($Notes) {
        Set-ADUser $User.SamAccountName -Replace @{info = $Notes}
    }
    if ($IPPhone) {
        Set-ADUser $User.SamAccountName -Replace @{IPPhone = $IPPhone}
    } 
   Write-Output "`nNew Active Directory user [$($user.UserPrincipalName)] added sucessfuly."
}

}
catch{
    Write-Output "`nFailed to add new Active Directory user [$($user.UserPrincipalName)]."
    Write-Error $_
}
