Function Execute_process($path,$bit)
{
$pinfo = New-Object System.Diagnostics.ProcessStartInfo
if($bit -eq 32 -and (test-path "$path\CCleaner.exe"))
{
$pinfo.FileName = "$path\CCleaner.exe"
}
elseif($bit -eq 64 -and (test-path "$path\CCleaner64.exe"))
{
$pinfo.FileName = "$path\CCleaner64.exe"
}
else
{
return -1
}
$pinfo.RedirectStandardError = $true
$pinfo.RedirectStandardOutput = $true
$pinfo.UseShellExecute = $false
$pinfo.Arguments = "/auto"
$p = New-Object System.Diagnostics.Process
$p.StartInfo = $pinfo
$p.Start() | Out-Null
$p.WaitForExit()
return $p.ExitCode
}



if((Test-path "C:\Program Files\CCleaner") -or (Test-Path "C:\Program Files (x86)\CCleaner"))
{
    Write-Output "CCleaner execution is started, it may take time to complete the execution"
    $osarchitecture = (get-wmiobject -Class win32_operatingsystem).Osarchitecture
    if($osarchitecture -eq "32-bit")
    {
    $path = Join-Path "C:\Program Files" "Ccleaner"
    if(test-path $path)
    {
    $success = Execute_process -path $path -bit 32
    }
  }

    elseif($osarchitecture -eq "64-bit")

    {
    $path32 = Join-Path "C:\Program Files (x86)" "Ccleaner"
    $path64 = Join-Path "C:\Program Files" "Ccleaner"
    if(Test-Path $path64) 
        {
        $success = Execute_process -path $path64 -bit 64
        } 
        elseif (Test-Path $path32) 
        {
        $success = Execute_process -path $path32 -bit 32
        }
    }
    if($success -eq 0)
    {
    Write-Output "CCleaner process executed successfully"
    }
    elseif($success -eq -1)
    {
    Write-Error "CCleaner.exe file is missing from C:\Program files\CCleaner"
    }
    else
    {
    Write-Error "Action Could not be completed"
    }
}
Else
{
Write-Output "CCleaner is not installed on this system"
exit
}
