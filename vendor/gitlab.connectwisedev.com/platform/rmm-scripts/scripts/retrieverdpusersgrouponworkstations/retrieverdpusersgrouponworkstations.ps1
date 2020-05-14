$Name = 'Remote Desktop Users'

Function Get-LocalGroupMembers ($Name) {
    try {
        $server = $env:COMPUTERNAME
        # find the group membership
        $Group = [ADSI]"WinNT://$Server/$Name,group"
        $Members = Foreach ($Item in $Group.psbase.Invoke("Members")) {
            $Property = {$Item.GetType().InvokeMember($args[0], 'GetProperty', $null, $Item, $null)}
            ($Property.invoke('ADSPath') -split "/")[2] + '\' + $Property.invoke('Name')            
        }
        
        If ($Members) {
            # if a group has members then display on the console
            Write-Output "Server: $Server`n"
            Write-Output "List Users in 'Remote Desktop Users' Group:"
            Write-Output $Members
        }
        # else { # NOTE: This section is commented out, which gets executed when no member found in the group. Incase we want to log the information in a CSV file.
        #     $OutputFile = "$env:TEMP\NoMember.csv"
        #     # if group doesn't have any members write\append to a file
        #     $obj = New-Object -TypeName  PSObject -Property @{Server = $Server; Comment = "No Member found in `'$Name`' Group on this server"}
        #     if (Test-Path $OutputFile) {
        #         # if file already exist then 'Append' to the file
        #         $obj | ConvertTo-Csv -NoTypeInformation |Select-Object -Skip 1 | Out-File $OutputFile -Append -Encoding ASCII    
        #     }
        #     else {
        #         # if file doesn't exist 'Write' to a new file
        #         $obj | Export-Csv $OutputFile -NoTypeInformation
        #     }
        # } 
    }
    catch {
        Write-Error "Server: `"$env:COMPUTERNAME`" - $($_.Exception.Message)" -ErrorAction Stop
    }
}

Get-LocalGroupMembers $Name
