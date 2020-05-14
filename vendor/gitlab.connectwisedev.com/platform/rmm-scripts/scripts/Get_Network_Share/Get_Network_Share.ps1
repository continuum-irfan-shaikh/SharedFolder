

clear-host

function Get-NetWorkShare {

  
   $hostVerSionMajor = ($PSVersionTable.PSVersion.Major).ToString()
   $hostVerSionMinor = ($PSVersionTable.PSVersion.Minor).ToString()
   $hostVersion = $hostVerSionMajor +'.'+ $hostVerSionMinor 

   $osVersionMajor = ([System.Environment]::OSVersion.Version.major).ToString()
   $osVersionMinor = ([System.Environment]::OSVersion.Version.minor).ToString()
   $osVersion = $osVersionMajor +'.'+ $osVersionMinor

   [boolean]$isPsVersionOk = ([version]$hostVersion -ge [version]'2.0')
   [boolean]$isOSVersionOk = ([version]$osVersion -ge [version]'6.0')
         
   Write-Host "Powershell Version : $($hostVersion)"
   if(-not $isPsVersionOk){
      
     Write-Warning "PowerShell version below 2.0 is not supported"
     return 

    }

   Write-Host "OS Name : $((Get-WMIObject win32_operatingsystem).Name.ToString().Split("|")[0])"  
   if(-not $isOSVersionOk){
   
      Write-Warning "PowerShell Script supports Window 7, Window 2008R2 and higher version operating system"
      return 

    }

    Write-Host " " 
   try{ 

       $GetAllNetShares = get-WmiObject -class Win32_Share -computer $env:COMPUTERNAME -errorAction Stop       
       if($GetAllNetShares){
           Write-Host "List of all network shares"
           Write-Host "---------------------------"
           $GetAllNetShares
       }else{

         Write-Warning "Unable to pull share information due to error"
       }
   }catch{

        Write-Warning "Unable to pull share information due to error"
   }
}

Get-NetWorkShare



# #get-WmiObject -class Win32_Share -computer $env:COMPUTERNAME
