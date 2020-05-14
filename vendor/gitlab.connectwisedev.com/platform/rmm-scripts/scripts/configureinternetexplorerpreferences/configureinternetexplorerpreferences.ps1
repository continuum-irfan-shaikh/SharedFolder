<#$homepageblank = $true            #Boolean
$homepage = $null                 #String
$downloaddir = "D:\Test1"         #String
$autoscript = "abcd.com"                  #String
$proxyenable = $false             #Boolean
$proxyaddress = "mmknd.corp.com"  #String
$proxyport = "80"                 #int
$bypassforlocaladdress = $true    #Boolean
$exceptions = "continuum.net"     #String Proxy exceptions
#>

$os = (Get-WmiObject -Class win32_operatingsystem).Name
if($os -like "*Server*"){
Write-output "Script cannot be executed on server"
Exit;
}

Function Create_registry($path, $Name, $Value, $propertytype) {
    # Function will set registry value and if item is not present will create registry entry.
    $Details = Get-ItemProperty -Path $path
    if ($Details -ne $null) {
        $Details = Get-ItemProperty -Path $path | gm | select -ExpandProperty Name
        if ($Details -contains $Name) {
            Set-ItemProperty $path -Name $Name -Value $Value
        }
        else {
            New-ItemProperty -Path $path -Name $Name -PropertyType $propertytype -Value $Value | Out-Null
        }
	
    } #End if statement
    else {
        New-ItemProperty -Path $path -Name $Name -PropertyType $propertytype -Value $Value | Out-Null
    }  # End else statement

    if ((Get-ItemProperty $path -Name $Name | select -ExpandProperty $Name) -eq $Value) { return $true } else { return $false }
}


$profilelist = Get-Item "REgistry::HKU\S-1-5-21-*" | Where-Object {$_.Name -notlike '*_Classes'} | select -ExpandProperty name
    Write-Output "`nMaking changes for logged in users"
    foreach ($profile in $profilelist) {
        if (Test-path "Registry::$profile\Volatile Environment") {
            $username = Get-ItemProperty -Path "Registry::$profile\Volatile Environment" -Name Username | Select -ExpandProperty Username

           if(Test-path "Registry::HKU\TEMP\Software\Microsoft\Internet Explorer\Main"){
               if($homepageblank -eq $true){
                $result = Create_registry -path "Registry::$profile\Software\Microsoft\Internet Explorer\Main" -Name "Start Page" -Value "" -propertytype "String"
                if ($result) {
                    Write-Output "Homepage has been removed for $username"
                }
               }
               if ($homepage -ne $null) {
                $result = Create_registry -path "Registry::$profile\Software\Microsoft\Internet Explorer\Main" -Name "Start Page" -Value $homepage -propertytype "String"
                if ($result) {
                    Write-Output "Homepage has been set for $username"
                }
            }
            if ($downloaddir -ne $null) {
                $result = Create_registry -path "Registry::$profile\Software\Microsoft\Internet Explorer\Main" -Name "Default Download Directory" -Value $downloaddir
                if ($result) {
                    Write-Output "Download directory has been set for $username."
                }
            }}else{Write-Output "Path not found for homepage and downloaddir"}
            if(Test-path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings"){
            if ($proxyenable -eq $true) {
                $result = Create_registry -path "Registry::$profile\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyEnable" -Value 1
                if ($result) {
                    Write-Output "Proxy has been enabled for $username"
                }
            }
            elseif($proxyenable -eq $false){
                $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyEnable" -Value 0
                if ($result) {
                    Write-Output "Proxy has been disabled for $username"
                }
            }
            if ($autoscript -ne $null) {
            $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "AutoConfigURL" -Value $autoscript -propertytype "String"
            if ($result) {
            Write-Output "Automatic configuration script has been enabled."
                }
            }
            if ($proxyaddress -ne $null) {
                $value = "$proxyaddress" + ":" + "$proxyport"
                $result = Create_registry -path "Registry::$profile\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyServer" -Value $value -propertytype "String"
                if ($result) {
                    Write-Output "Proxyaddress has been enabled for $username"
                }
            }
            if (($bypassforlocaladdress -eq $true) -and ($exceptions -eq $null)) {
                $result = Create_registry -path "Registry::$profile\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyOverride" -Value "<local>" -propertytype "String"
                if ($result) {
                    Write-Output "Bypass without exceptions has been set for $username."
                }
            }
            elseif (($bypassforlocaladdress -eq $true) -and ($exceptions -ne $null)) {
                $value = "$exceptions" + ";" + "<local>"
                $result = Create_registry -path "Registry::$profile\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyOverride" -Value "$value" -propertytype "String"
                if ($result) {
                    Write-Output "Bypass with exceptions has been set for $username."
                }
            }
           } else{Write-Output "Failed to update proxy status,proxyaddress and exceptions"}
        }
        else {
            Continue
        }
    }


$profiles = Get-ChildItem C:\Users | select -ExpandProperty Name
Write-Output "`nMaking changes for user profiles"
foreach ($profile in $profiles) {
    
    $ntuser = Test-Path C:\Users\$profile\NTUSER.DAT
    if ($ntuser) {
        try{
            $FileStream = [System.IO.File]::Open("C:\Users\$profile\NTUSER.DAT", 'Open', 'Write')
            $FileStream.Close()
            $FileStream.Dispose()
            reg.exe LOAD HKU\TEMP "C:\Users\$profile\NTUSER.DAT" | Out-Null
            if(test-path "Registry::HKU\TEMP\Software\Microsoft\Internet Explorer\Main"){
                if($homepageblank -eq $true){
                    $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Internet Explorer\Main" -Name "Start Page" -Value "" -propertytype "String"
                    if ($result) {
                        Write-Output "Homepage has been removed for $profile"
                    }
                   }
            if ($homepage -ne $null) {
                $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Internet Explorer\Main" -Name "Start Page" -Value $homepage -propertytype "String"
                if ($result) {
                    Write-Output "Homepage has been set for $profile"
                }
            }
            if ($downloaddir -ne $null) {
                $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Internet Explorer\Main" -Name "Default Download Directory" -Value $downloaddir
                if ($result) {
                    Write-Output "Download directory has been set for $profile."
                }
            }}else{Write-Output "Failed to update Homepage and Download directory"}
            if(Test-path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings"){
            if ($proxyenable -eq $true) {
                $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyEnable" -Value 1
                if ($result) {
                    Write-Output "Proxy has been enabled for $profile"
                }
            }
            elseif($proxyenable -eq $false){
                $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyEnable" -Value 0
                if ($result) {
                    Write-Output "Proxy has been disabled for $profile"
                }
            }
            if ($autoscript -ne $null) {
            $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "AutoConfigURL" -Value $autoscript -propertytype "String"
            if ($result) {
            Write-Output "Automatic configuration script has been enabled."
             }
            }
            if ($proxyaddress -ne $null) {
                $value = "$proxyaddress" + ":" + "$proxyport"
                $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyServer" -Value $value -propertytype "String"
                if ($result) {
                    Write-Output "Proxyaddress has been enabled for $profile"
                }
            }
            if (($bypassforlocaladdress -eq $true) -and ($exceptions -eq $null)) {
                $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyOverride" -Value "<local>" -propertytype "String"
                if ($result) {
                    Write-Output "Bypass without exceptions has been set for$profile ."
                }
            }
            elseif (($bypassforlocaladdress -eq $true) -and ($exceptions -ne $null)) {
                $value = "$exceptions" + ";" + "<local>"
                $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyOverride" -Value "$value" -propertytype "String"
                if ($result) {
                    Write-Output "Bypass with exceptions has been set for $profile."
                }
            }}else{Write-Output "Failed to update proxy status,proxyaddress and exceptions"}
            Start-Sleep -Seconds 10
            [gc]::collect()
            reg.exe UNLOAD HKU\TEMP | Out-Null           
    } 
    catch{
Continue
}  
}

}

    $output = REG LOAD HKU\TEMP "C:\Users\default\NTUSER.DAT"
if ($?) {
    Write-Output "`nMaking changes for default users"
    if(Test-path "Registry::HKU\TEMP\Software\Microsoft\Internet Explorer\Main"){
        if($homepageblank -eq $true){
            $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Internet Explorer\Main" -Name "Start Page" -Value "" -propertytype "String"
            if ($result) {
                Write-Output "Homepage has been removed for default profiles"
            }
           }
    if ($homepage -ne $null) {
        $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Internet Explorer\Main" -Name "Start Page" -Value $homepage -propertytype "String"
        if ($result) {
            Write-Output "Homepage has been set for default profiles."
        }
    }
    if ($downloaddir -ne $null) {
        $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Internet Explorer\Main" -Name "Default Download Directory" -Value $downloaddir
        if ($result) {
            Write-Output "Download directory has been set for default profiles."
        }
    }}else{Write-Output "Path not found for homepage and downloaddir"}
    if(Test-path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings"){
    if ($proxyenable -eq $true) {
        $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyEnable" -Value 1
        if ($result) {
            Write-Output "Proxy has been enabled for default profiles."
        }
    }
    elseif($proxyenable -eq $false){
                $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyEnable" -Value 0
                if ($result) {
                    Write-Output "Proxy has been disabled for default profiles."
                }
            }
    if ($autoscript -ne $null) {
        $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "AutoConfigURL" -Value $autoscript -propertytype "String"
        if ($result) {
            Write-Output "Automatic configuration script has been enabled."
        }
    }
    if ($proxyaddress -ne $null) {
        $value = "$proxyaddress" + ":" + "$proxyport"
        $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyServer" -Value $value -propertytype "String"
        if ($result) {
            Write-Output "Proxyaddress has been enabled for default profiles."
        }
    }
    if (($bypassforlocaladdress -eq $true) -and ($exceptions -eq $null)) {
        $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyOverride" -Value "<local>" -propertytype "String"
        if ($result) {
            Write-Output "Bypass without exceptions has been set for default profiles."
        }
    }
    elseif (($bypassforlocaladdress -eq $true) -and ($exceptions -ne $null)) {
        $value = "$exceptions" + ";" + "<local>"
        $result = Create_registry -path "Registry::HKU\TEMP\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name "ProxyOverride" -Value "$value" -propertytype "String"
        if ($result) {
            Write-Output "Bypass with exceptions has been set for default profiles."
        }
    }}else{Write-Output "Failed to set proxyaddress"}                    
    
}
else {
    Write-Error "Failed to load registry for default profile"
}
start-sleep -Seconds 10
[gc]::collect()
REG UNLOAD HKU\TEMP

#$output = REG UNLOAD HKU\TEMP 2>&1



