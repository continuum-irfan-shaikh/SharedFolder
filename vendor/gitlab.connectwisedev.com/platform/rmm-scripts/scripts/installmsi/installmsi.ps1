#As per Microsoft Article :: https://docs.microsoft.com/en-us/windows/desktop/msi/standard-installer-command-line-options
$Switch = "/install", "/uninstall", "/NoRestart", "/Quiet","/passive", "/promptrestart","/forcerestart", "/update","/uninstall","/package". "/log"

$UI = foreach ($Par in $Parameter) { if ($Switch -notcontains $Par) { "Incorrect Switch Please enter correct Switch for :: $Par"; break } }

if(!$UI) {

if (!(Test-Path $Path)) {
  Write-Error "Incorrect path $Path. File doesn't exist"
  return
}
$MSIArguments = @(
    "/i"
    $Path
    "/quiet"
    $Parameter
    'AllUsers="1"'
)
$exitCode = (Start-Process -FilePath "msiexec.exe" -ArgumentList $MSIArguments -Wait -Passthru).ExitCode
if($exitCode -ne 0) {
    Write-Error "Installation process returned error code: $exitCode"
} else {
    Write-Output "Successful installation process"
    }
}
