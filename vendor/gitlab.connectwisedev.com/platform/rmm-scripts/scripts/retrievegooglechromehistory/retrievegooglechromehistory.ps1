<#
Template Name: Retrieve Google Chrome history

Description:
Retrieves google chrome history. Script allows to retrive history for slected users or all users.
#>

$OSArch = (Get-WMIObject Win32_OperatingSystem).OSArchitecture
Switch($OSArch){
   "64-bit" {  $sqlite3 = ${ENV:ProgramFiles(x86)}+"\ITSPlatform\plugin\scripting\sqlite3.exe"}                              
   "32-bit" {  $sqlite3 = ${ENV:ProgramFiles}+"\ITSPlatform\plugin\scripting\sqlite3.exe" }               
}
if (!$days){ $days = 90 } 
$SqlyQuery = "SELECT datetime(last_visit_time / 1000000 + (strftime('%s', '1601-01-01')),
              'unixepoch','localtime' ) as visit_time, visit_count, 
               url FROM urls WHERE visit_time > (SELECT DATETIME('now', '-$days day')) ORDER by visit_time DESC;"
$profileDirs = @()
if($users -eq "allusers"){    
        Get-ChildItem -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList" | ForEach-Object {
        	$userProfile = (Get-ItemProperty -Path $_.PSPath).ProfileImagePath
        	
        	if ($userprofile.substring(0, $($env:windir).length) -eq $env:windir) {
        		# Skipping - Profile in Windows Folder
        	} elseif (!(Get-ChildItem $userProfile -ErrorAction SilentlyContinue)) {
        		# Skipping - Folder does not exist"
        	} else {
        		$profileDirs += $userProfile
        	}
        }
}else{ $profileDirs = $usernames -split(",") | %{"c:\users\" + $_ }}

$Result = @()
$Errors = @()
foreach ($profileDir in $profileDirs){
        $user = ($profileDir -split ('\\'))[2]
        $HistoryDB = Join-Path $profileDir "\AppData\Local\Google\Chrome\User Data\Default\History" #
        if (Test-Path $HistoryDB -EA SilentlyContinue ) {
         
            Try {
                  Copy-Item $HistoryDB $ENV:TEMP -ErrorAction Stop
                  sleep 2
                  $tempHistoryDB = Join-Path $ENV:TEMP "\History"            
                  $Histories = $SqlyQuery | & $sqlite3 $tempHistoryDB 2>&1
                  if ($LASTEXITCODE -ne 0) {
                      $Errors +=  "Failed to retrieve history for user : $user"
                      $Histories = $null
                  }else{            
                       if ($Histories -ne $null ) {
                           foreach ( $Data in $Histories ) {
                                $uri = [uri]($Data -split "\|")[2]  
                                $obj  = New-Object PSObject -Property @{ 
                                                                         "User Name" = $user
                                                                         "Date" = ($Data -split "\|")[0]
                                                                         "Site Name" = $uri.Scheme + "://" + $uri.Host
                                                                       }                      
                                $Result += $obj           
                       }
                       } Else { $Errors += "Histroy not found for user : $user" }
        
                  Remove-Item $tempHistoryDB -ErrorAction SilentlyContinue                              
                 }            
            } catch {
                 $Errors += "Failed to retrieve history for user : $user"
                  continue
                }
           }else {
               $Errors +=  "Histroy not found for user : $user"             
           }             
 }           
$Result | select "User Name" ,"Date", "Site Name" | ft -AutoSize
$Errors
