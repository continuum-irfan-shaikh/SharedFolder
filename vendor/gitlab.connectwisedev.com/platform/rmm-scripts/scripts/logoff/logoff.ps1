If((gwmi win32_operatingsystem -ComputerName $env:COMPUTERNAME).Win32Shutdown(4)){
    Return "Log off successfully done"
}