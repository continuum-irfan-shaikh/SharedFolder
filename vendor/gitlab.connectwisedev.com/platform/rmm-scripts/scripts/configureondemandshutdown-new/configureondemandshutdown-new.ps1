<#

    .Name
    Configure On demand shutdown
    
    .Description
        1)Script will Ask/Notify loggedIn user for initiating shutdown.
        2)User can set criteria to shutdown the machine 
        3)User can Pass Command or Process to execute before initiating shutdown if criteria doesn't met.
        4)User can Deny shutdown if Critera doesn't met.
    
    .Requirements
        1)Script should run in user context.
        2)Script will run only when user is loggedIn to the system.
        3)LoggedIn user should have permission to shutdown the Machine.

    .Author
        Nirav Sachora


#>

<#

$IfUserIsLoggedOn
$AddProcess
$processname
$action
$AddCmdProcess
$command
#$process
#$parameter
$displaymsg
$ShutdownTimeout
#$Allow

#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

Function Main {
    Add-Type -AssemblyName System.Windows.Forms
    Add-Type -AssemblyName System.Drawing
    ###############################################################################################################  
    if ($IfUserIsLoggedOn -eq "Ask the user for confirmation before initiating Shutdown") {
        Add-Type -AssemblyName System.Windows.Forms
        [System.Windows.Forms.Application]::EnableVisualStyles()

        Function Auto_Close {
            $Timer.Stop(); 
            $Form.Close(); 
            $Form.Dispose();
            $Timer.Dispose();
            &{shutdown /s /f /t 15}
        }

        Function Timer_Tick {
            $Label3.Text = "Confirm Shutdown in $Script:CountDown seconds"
            --$Script:CountDown
            if ($Script:CountDown -lt 0) {
                Auto_Close
            }
        }

        $form = New-Object system.Windows.Forms.Form
        $form.ClientSize = '546,237'
        $form.text = "Shutdown Manager"
        $form.TopMost = $false
        $form.StartPosition = "CenterScreen"

        $Button1 = New-Object system.Windows.Forms.Button
        $Button1.BackColor = "#9b9b9b"
        $Button1.text = "Shutdown Later"
        $Button1.width = 131
        $Button1.height = 39
        $Button1.location = New-Object System.Drawing.Point(355, 85)
        $Button1.Font = 'Microsoft Sans Serif,10'
        $Button1.DialogResult = 1

        $Button2 = New-Object system.Windows.Forms.Button
        $Button2.BackColor = "#9b9b9b"
        $Button2.text = "Shutdown Now"
        $Button2.width = 123
        $Button2.height = 40
        $Button2.location = New-Object System.Drawing.Point(241, 173)
        $Button2.Font = 'Microsoft Sans Serif,10'
        $Button2.DialogResult = 6

        $Button3 = New-Object system.Windows.Forms.Button
        $Button3.BackColor = "#9b9b9b"
        $Button3.text = "Abort"
        $Button3.width = 121
        $Button3.height = 39
        $Button3.location = New-Object System.Drawing.Point(393, 173)
        $Button3.Font = 'Microsoft Sans Serif,10'
        $Button3.DialogResult = 7

        $Label1 = New-Object system.Windows.Forms.Label
        $Label1.text = "Select Postpone for"
        $Label1.AutoSize = $true
        $Label1.width = 25
        $Label1.height = 10
        $Label1.location = New-Object System.Drawing.Point(26, 90)
        $Label1.Font = 'Microsoft Sans Serif,10'

        $Label2 = New-Object system.Windows.Forms.Label
        $Label2.AutoSize = $true
        $Label2.width = 25
        $Label2.height = 10
        $Label2.location = New-Object System.Drawing.Point(10, 20)
        $Label2.Font = 'Microsoft Sans Serif,10'

        $Label3 = New-Object system.Windows.Forms.Label
        $Label3.text = "Confirm shutdown in "
        $Label3.AutoSize = $true
        $Label3.width = 25
        $Label3.height = 10
        $Label3.location = New-Object System.Drawing.Point(21, 21)
        $Label3.Font = 'Microsoft Sans Serif,10'

        $Label4 = New-Object system.Windows.Forms.Label
        $Label4.AutoSize = $true
        $Label4.width = 25
        $Label4.height = 10
        $Label4.location = New-Object System.Drawing.Point(213, 20)
        $Label4.Font = 'Microsoft Sans Serif,10'
        
        $Label5 = New-Object system.Windows.Forms.Label
        $Label5.AutoSize = $true
        $Label5.width = 500
        $Label5.height = 10
        $Label5.Text = $displaymsg
        $Label5.location = New-Object System.Drawing.Point(21, 45)
        $Label5.Font = 'Microsoft Sans Serif,10'

        $SelectOne = New-Object system.Windows.Forms.ComboBox
        $SelectOne.text = "Select One"
        $SelectOne.width = 105
        $SelectOne.height = 22
        @('15 Minutes', '30 Minutes', '60 Minutes') | ForEach-Object { [void] $SelectOne.Items.Add($_) }
        $SelectOne.location = New-Object System.Drawing.Point(158, 89)
        $SelectOne.Font = 'Microsoft Sans Serif,10'

        $Script:countdown = $ShutdownTimeout
        $Form.controls.AddRange(@($Button1, $Button2, $Button3, $Label1, $Label2, $Label3, $Label4, $SelectOne,$Label5))
        
        $Timer = New-Object System.Windows.Forms.Timer
        $Timer.Interval = 1000
        $Button1.Add_Click( {
                $Global:Delaytime = 0
                switch ($SelectOne.text) {
                    "15 Minutes" {shutdown /s /f /t 900;$Global:Delaytime = 15}
                    "30 Minutes" {shutdown /s /f /t 1800;$Global:Delaytime = 30}
                    "60 Minutes" {shutdown /s /f /t 3600;$Global:Delaytime = 60}
                    Default { [System.Windows.MessageBox]::Show('Please select an option to postpone shutdown.') }
                }
            })
            
                
            
        $Button2.Add_Click({&{shutdown /s /f /t 15}})
        $Button3.Add_Click( { $form.Close() })
        $Timer.Add_Tick( { Timer_Tick })
        $Timer.Start()
        $Output = $form.ShowDialog()
        Switch($Output){
        "Yes"{"Shutdown Initiated."}
        "OK"{"Shutdown Scheduled after $Global:Delaytime Minutes."}
        "No"{"User Denied to shutdown."}
        "Cancel"{"Shutdown Initiated"}
        }
    }

    ########################################################################################################################

    ########################################################################################################################
    if ($IfUserIsLoggedOn -eq "Notify the user that the machine will Shutdown") {
        Add-Type -AssemblyName System.Windows.Forms
        [System.Windows.Forms.Application]::EnableVisualStyles()

        Function Auto_Close {
            $Timer.Stop(); 
            $Form.Close(); 
            $Form.Dispose();
            $Timer.Dispose();
            &{shutdown /s /f /t 15}
        }

        Function Timer_Tick {
            $Label1.Text = "This Device will shutdown.$Script:CountDown seconds"
            --$Script:CountDown
            if ($Script:CountDown -lt 0) {
                Auto_Close
            }
        }

        $Form = New-Object system.Windows.Forms.Form
        $Form.ClientSize = '488,234'
        $Form.text = "Shutdown Manager"
        $Form.TopMost = $false
        $form.StartPosition = "CenterScreen"

        $Label1 = New-Object system.Windows.Forms.Label
        $Label1.AutoSize = $true
        $Label1.width = 500
        $Label1.height = 10
        $Label1.location = New-Object System.Drawing.Point(20, 16)
        $Label1.Font = 'Microsoft Sans Serif,10'

        $Label2 = New-Object system.Windows.Forms.Label
        $Label2.AutoSize = $true
        $Label2.width = 500
        $Label2.height = 10
        $Label2.location = New-Object System.Drawing.Point(20, 35)
        $Label2.Font = 'Microsoft Sans Serif,10'

        $Button1 = New-Object system.Windows.Forms.Button
        $Button1.text = "OK"
        $Button1.width = 138
        $Button1.height = 47
        $Button1.location = New-Object System.Drawing.Point(170, 170)
        $Button1.Font = 'Microsoft Sans Serif,10'
        $Button1.DialogResult = 1

        $Script:countdown = $ShutdownTimeout
        $Label2.Text = $displaymsg
        $Label1.Text = "This Device will shutdown in "
        $Form.controls.AddRange(@($Label1, $Label2, $Button1))
        $Timer = New-Object System.Windows.Forms.Timer
        $Timer.Interval = 1000

        $Button1.Add_Click({&{shutdown /s /f /t 15}})
        $Timer.Add_Tick( { Timer_Tick })
        $Timer.Start()
        $result = $form.ShowDialog()
        switch($result){
            "Ok"{"Shutdown Initiated"}
            "Cancel"{
                "Shutdown Initiated"
                $Timer.Stop()
                $Timer.Dispose()

            }
        }
    }
    ########################################################################################################################
}

########################################################################################################################

if ($AddProcess) {  
    $processdetails = Get-Process -Name $processname -ErrorAction "SilentlyContinue"
    if ($processdetails -ne $null) {
        if ($action -eq "Do not initiate Shutdown") {
            Write-Output "$processname is running in the memory`n`nWindows cannot be shutdown."
            Exit;
        }
        elseif ($action -eq "Execute the following commands and process the Shutdown") {
            if ($AddCmdProcess -eq "Command") {
                $Commandtoexecute = "$command"
                try {
                    $Erroractionpreference = "Stop"
                    Invoke-Command -ScriptBlock { & { cmd.exe /c $Commandtoexecute } } | Out-Null
                    if ($?) {
                        Main
                    }
                    Else {
                        Write-Error "Command Failed to Execute"
                    }
                }
                catch {
                    $_.Exception.Message
                    Exit;
                }
            }
            if ($AddCmdProcess -eq "Process") {
                $pinfo = New-Object System.Diagnostics.ProcessStartInfo
                $pinfo.FileName = "$process"
                $pinfo.RedirectStandardError = $true
                $pinfo.RedirectStandardOutput = $true
                $pinfo.UseShellExecute = $false
                if ($parameter) {
                    $pinfo.Arguments = "$parameter"
                }
                $p = New-Object System.Diagnostics.Process
                $p.StartInfo = $pinfo
                $ErrorActionPreference = "Stop"
                $p.Start() | Out-Null
                $ErrorActionPreference = "Continue"
                if (!$Allow) { $exitstatus = $p.WaitForExit() }
                else {
                    $exittime = $Allow * 60000
                    $exitstatus = $p.WaitForExit($exittime)
                }
                if ($exitstatus) { Main }
                else { $p.Kill(); Write-Error "Process could not be executed in given time." }
            }
        }
    }
    else{Main}
}
else{Main}
