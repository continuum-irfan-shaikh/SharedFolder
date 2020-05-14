<#
$eventlogType
$DefalutEventLog
$UserDefinedEventLog
#>

try {
    $erroractionpreference = 'stop'
	if($eventlogType -eq "Default"){
    		Clear-EventLog -logname $DefalutEventLog
	}
	elseif($eventlogType -eq "User Defined"){
		Clear-EventLog -logname $UserDefinedEventLog
	}
    if($?){
        "-" * 30 + "`n[MSG]$eventlog logs deleted successfully`n" + "-" * 30
    }
}
catch {
    "-" * 30 + "`n[ERROR]Error Occured`n" + "-" * 30
    Write-Error $_.Exception.Message
    Exit;
}
