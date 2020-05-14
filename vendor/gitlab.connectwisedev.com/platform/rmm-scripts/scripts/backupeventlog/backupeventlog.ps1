<#
Name : Backup event log,
Category : Maintenance

$eventLogType = "Default"
$eventLog = "Application"
$logFileName = "C:\temp\san\san1\san2\san3\san4\san5\san6\New"
$logFileExtension = ".evt"
########################
$overwriteExistingFile = $False
$createDirIfRequired = $True
#>

# Declare Variables

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}


# Validation for NOT NULL or empty for mandatory paramaters

if(($eventLog -eq $NULL) -or ($eventLog -eq ''))
  {
     Write-Host "Event log can not be empty"  
     EXIT
  }
    
if(($logFileName -eq $NULL) -or ($logFileName -eq ''))
  {
     Write-Host "Log File Name can not be empty"  
     EXIT
  }

############################################################################
# Function to store backup of the Event Log
############################################################################


Function Create_Directory {

  [CmdletBinding(ConfirmImpact='Low')] 
    Param(
        [Parameter(Mandatory=$true,
                   ValueFromPipeLine=$true,
                   ValueFromPipeLineByPropertyName=$true,
                   Position=0)]
            [String]$FolderName, 
        [Parameter(Mandatory=$false,
                   Position=1)]
            [Switch]$NoCreate = $false
    )

   if ($FolderName.Length -gt 254) {
        Write-Error "Folder name '$FolderName' is too long - ($($FolderName.Length)) characters"
        break
    }
    if (Test-Path $FolderName) {
        Write-Verbose "Confirmed folder '$FolderName' exists"
        #$true
    } else {
        Write-Verbose "Folder '$FolderName' does not exist"
        if ($NoCreate) {
            $false
            break  
        } else {
            Write-Verbose "Creating folder '$FolderName'"
            try {
                New-Item -Path $FolderName -ItemType directory -Force -ErrorAction Stop | Out-Null
                Write-Verbose "Successfully created folder '$FolderName'"
                #$true
            } catch {
                Write-Error "Failed to create folder '$FolderName'"
                $false
            }
        }
    }
}


Function BackupEventLog($eventLog,$logFileName,$logFileExtension,$overwriteExistingFile,$createDirIfRequired)
{

Try
{
    #generating file name
    $FileName= $logFileName+$logFileExtension

    #check whether the directory/file path is present or not
    #For Directory
    if(!(Test-Path (Split-Path -Path $FileName)))
    {
      Write-Host "Directory is not present"
      if($createDirIfRequired -eq $True)
      {
        $Last = $FileName.Split('\')[-1]
        $Data = $FileName.Replace($Last,"")

        Create_Directory -FolderName $Data
      }
      else
      {
         Write-Host "Unable to create directory as user selected 'No' for creating new directory if not present"
         EXIT
      }
    }
    #For File
    if(Test-Path $FileName)
    {
       Write-Host "Unable to create Files $filename as user selected 'No' for Over Write Existing File"

    }else{
    
    #based on the  log file extension create log file
    if($logFileExtension -eq ".txt")
     {
        Get-EventLog -LogName $eventLog |Select-Object EventID,MachineName,Index,Category,CategoryNumber,EntryType, `
                     Message,Source,InstanceId,TimeGenerated,TimeWritten,UserName,Site,Container| `
                     Out-File $FileName -Force
        $host.UI.WriteLine("Log file has been generated successfully, File Name: $FileName")
     }
    elseif( $logFileExtension -eq ".csv")
    {
        Get-EventLog -LogName $eventLog |Select-Object EventID,MachineName,Index,Category,CategoryNumber,EntryType, `
                  Message,Source,InstanceId,TimeGenerated,TimeWritten,UserName,Site,Container| `
                  Export-Csv $FileName -NoTypeInformation -NoClobber -Force
        $host.UI.WriteLine("Log file has been generated successfully, File Name: $FileName")
    }

    elseif( $logFileExtension -eq ".evt")
    {
        $logFile = Get-WmiObject Win32_NTEventlogFile | Where-Object {$_.logfilename -eq $eventLog}
        $logFile.backupeventlog($FileName) | Out-Null
        $host.UI.WriteLine("Log file has been generated successfully, File Name: $FileName")
    }
    }

 }
 Catch
 {
 $Host.UI.WriteErrorLine($Error[0].Exception.Message)
 }
}
########## End Function Backup Event Log 

#Calling Function
Try
{

   BackupEventLog -eventLog $eventLog -logFileName $logFileName -logFileExtension $logFileExtension -overwriteExistingFile $overwriteExistingFile -createDirIfRequired $createDirIfRequired
}
Catch{

 $Host.UI.WriteErrorLine($Error[0].Exception.Message)
}
