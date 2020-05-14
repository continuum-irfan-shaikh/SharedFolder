$app = Get-WmiObject -Class win32_product | Where-object {$_.Name -like 'Vipre*' -and $_.Version -like "9.6*"}
if(!$app)
    {
    Write-Error "Vipre is not Installed On this System"
    exit
    }
Else
    {
    
    $result = $app.uninstall()
        if($result.ReturnValue -eq 0)
        {
        Write-Output "Vipre is Uninstalled from the System"
        }
        elseif($result.ReturnValue -eq 3010)
        {
        Write-Output "Vipre Uninstalled, Please Restart the System"
        }
        else
        {
        Write-Error "Uninstallation Failed with return Value" $($result.ReturnValue)
        }
    }
