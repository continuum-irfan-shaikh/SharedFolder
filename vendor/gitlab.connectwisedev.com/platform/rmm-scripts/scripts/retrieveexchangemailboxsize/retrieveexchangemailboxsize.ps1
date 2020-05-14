
<#
    .SYNOPSIS
       Retrieve Exchange mailbox size
    .DESCRIPTION
       Retrieve Exchange mailbox size
    .Author
       Santosh.Dakolia@continuum.net
#>


if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    write-warning "Excecuting the script under 64 bit powershell"
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

$ProductType = (Get-WmiObject -Class Win32_OperatingSystem).ProductType 
$computername= $env:computername

Try{

    IF($ProductType -eq 1){"'$computername' is a Client Machine hence unable to fetch Information"; exit}
    IF($ProductType -eq 2 -or $ProductType -eq 3){

$warningpreference = "SilentlyContinue"
$Service = Get-Service -Name MSExchangeIS 
$Stauts = $Service.Status
$Exch2007 = @()
$Exch2010 = @()

        IF(!$Service) {“MSExchangeIS service not found”}
        IF(($Service) -and ($Stauts -eq 'stopped')){"MSExchangeIS Service is Not Running"}
        IF($Service -and $Stauts -eq 'running') {


    Add-PSSnapin *Exch* -ErrorAction SilentlyContinue 

    $data = Get-Mailbox -ResultSize unlimited | Get-MailboxStatistics

    $ver = (Get-ExchangeServer -Identity $env:COMPUTERNAME).admindisplayversion.major
IF($ver -eq 8){

    ForEach($MB in $data){
    $SizeM = $MB.TotalItemsize.Value
            IF($SizeM -ge 1GB){$Sizeu = "{0:N2} GB" -f ($MB.TotalItemSize.Value.ToGB()) } 
                elseif($SizeM -ge 1Mb) {$Sizeu = "{0:N2} MB" -f ($MB.TotalItemSize.Value.ToMB())  } 
                    else {$Sizeu = "{0:N2} BYtes" -f ($MB.TotalItemSize.Value.ToKB())} 
          
        $outputData = New-Object -TypeName psobject 
        $outputData | Add-Member -MemberType NoteProperty -Name ServerName -Value $MB.servername
        $outputData | Add-Member -MemberType NoteProperty -Name StorageGroup -Value $MB.StorageGroupName
        $outputData | Add-Member -MemberType NoteProperty -Name DatabaseName -Value $MB.DatabaseName
        $outputData | Add-Member -MemberType NoteProperty -Name MailboxDisplayName -Value $MB.DisplayName
        $outputData | Add-Member -MemberType NoteProperty -Name Size -Value $Sizeu
        $outputData | Add-Member -MemberType NoteProperty -Name TotalItems -Value $MB.Itemcount
        $Exch2007 += $outputdata
            }
           $Exch2007
        }
    }

IF($Ver -ge 14){
    ForEach($Mailbox in $data){
    $SizeMB = $Mailbox.TotalItemsize.Value
            IF($SizeMB -ge 1GB){$SizeWithUnit = "{0:N2} GB" -f ($Mailbox.TotalItemSize.Value.ToGB()) } 
                elseif($SizeMB -ge 1Mb) {$SizeWithUnit = "{0:N2} MB" -f ($Mailbox.TotalItemSize.Value.ToMB())  } 
                    else {$SizeWithUnit = "{0:N2} KB" -f ($Mailbox.TotalItemSize.Value.ToKB())} 

        $output = New-Object -TypeName psobject 
        $output | Add-Member -MemberType NoteProperty -Name ServerName -Value $mailbox.servername
        $output | Add-Member -MemberType NoteProperty -Name DataBaseName -Value $mailbox.databasename
        $output | Add-Member -MemberType NoteProperty -Name MailboxDisplayName -Value $mailbox.DisplayName
        $output | Add-Member -MemberType NoteProperty -Name Size -Value $SizeWithUnit
        $output | Add-Member -MemberType NoteProperty -Name TotalItems -Value $mailbox.Itemcount
        $Exch2010 += $output
            }
           $Exch2010
        }
    }
}
Catch{

    Write-Error "Error occured while retrieving Data..!! $_.Exception.Message"
}
