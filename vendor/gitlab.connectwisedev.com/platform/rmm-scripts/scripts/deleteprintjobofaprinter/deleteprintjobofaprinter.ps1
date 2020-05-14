#Developed By :  GRT , Continuum
#File Vesion  :  1.0
try{
    $OutputMSG = @()
    $PrinterInfo = Get-WmiObject -Class Win32_printer  -namespace "root\CIMV2" -filter "name='$PrinterName'"
    If($PrinterInfo)
    {
        try
        {
            $PrinterInfo | Foreach{$_.CancelAllJobs()} | out-null
            $PrintQueueStatus = Get-WMIObject Win32_PerfFormattedData_Spooler_PrintQueue -filter "name='$PrinterName'" |  Select Name, @{Expression={$_.jobs};Label="CurrentJobs"}, JobErrors
            $OutputMSG += $PrintQueueStatus |ft
        }
        Catch
        {
            $OutputMSG += "Failed to cancel print jobs stuck in the print queue.Please try Force deletion of Print Job"
        }
    }
    Else
    {
        $OutputMSG += "Cannot find any printer name :  $PrinterName . Please check from the below list"
        $PrintQueueStatus = Get-WMIObject Win32_PerfFormattedData_Spooler_PrintQueue |  Select Name, @{Expression={$_.jobs};Label="CurrentJobs"}, JobErrors
        $OutputMSG += $PrintQueueStatus
    }
    Write-Output $OutputMSG
}catch{ Write-Output "ERROR while deleting print job of a printer $_.exception.message" }