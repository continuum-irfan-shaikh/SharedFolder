$Items = Get-Childitem -Path $env:TEMP -Recurse

$ItemsDeleted = 0
$Length = 0

If (!$Items)
{
    Write-Output "Temp folder is empty"
}
if ($Items)
{

    If ($Items.Count -ne 0)
    {
        $Content = Get-ChildItem -Path $env:TEMP -Force -Recurse -ErrorAction SilentlyContinue | Measure-Object -Property Length -Sum
        If ($Content)
        {
            $Length += $Content.Sum
        }
        $LengthNew = 0
        $Items | ForEach {
            Remove-Item $_.FullName -Recurse -Force -ErrorAction SilentlyContinue
            If ((-Not(Test-Path $_.FullName -ErrorAction SilentlyContinue)) -And (-Not($Error[0].Exception -is [System.UnauthorizedAccessException])))
            {
                $ItemsDeleted++
            }
        }
        $Content = Get-ChildItem -Path $env:TEMP -Force -Recurse -ErrorAction SilentlyContinue | Measure-Object -Property Length -Sum
        If ($Content)
        {
            $LengthNew += $Content.Sum
        }
        If ($ItemsDeleted -eq 0)
        {
            Write-output "$( $Items.Count ) file(-s) has(-ve) not been deleted because of access rights or they were using by another process"
            Return
        }
        If (($Items.Count - $ItemsDeleted) -eq 0)
        {
            Write-Output "Cleared $( $Length - $LengthNew ) bytes, $ItemsDeleted file(-s)"
            Return
        }
        If (($Items.Count - $ItemsDeleted) -gt 0)
        {
            Write-Output "Cleared $( $Length - $LengthNew ) bytes, $ItemsDeleted file(-s), $( $Items.Count - $ItemsDeleted ) file(-s) has(-ve) not been deleted because of access rights or they were using by another process"
            Return
        }
    }
}