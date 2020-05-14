$computernames = $env:computername

IF ($force) {
    Write-Output "Successfully restarted"
    Restart-Computer -ComputerName $computernames -Force:$true
} ELSE {
    shutdown.exe -r -t 300 -c 'This system will be restarted in 5 minutes... Please save all work in progress. Any unsaved changes will be lost'
    Write-Output "Restart is scheduled and will be performed in 5 minutes"
}
