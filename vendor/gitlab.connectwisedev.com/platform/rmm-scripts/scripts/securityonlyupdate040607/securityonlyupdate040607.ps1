<#
   Description:
   Security-only update installer script.
   This script is supported for only Windows7 SP1 and higher Windows 2008 R2 SP1 and higher versions.
   Name : Security-Only Updates Script for April 2018 to September 2019 
   JIRA : GRT-5203
#>

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }else{
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
exit $lastexitcode
}

$OSver = [System.Environment]::OSVersion.Version
$OSName = (Get-WMIObject win32_operatingsystem).Caption
$OSArch = [intPtr]::Size
$Hostname = $env:COMPUTERNAME
$DownloadPath = $env:SystemRoot+"\temp"
$logfilePath = $env:SystemRoot+"\Logs\WindowsUpdate\"

if(-not($OSver.Major -eq 6 -and $OSver.Minor -eq 1 -and $OSver.Build -ge 7601)) {
   Write-Output "The security update is not applicable for OS - $OSName. No action is needed."
   Exit 
}
$MRollupKB4516065 = Get-HotFix "KB4516065" -ErrorAction SilentlyContinue
if($MRollupKB4516065){
       Write-Output "This security-only update is not applicable on this computer as the Monthly Rollup KB4516065 is already installed"
       Exit
}

Switch($OSArch){
    4 { 
        #Windows 7 32 Bit KBs    
        $MSUurls = ("http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/04/windows6.1-kb4093108-x86_f0e2e9c3d7cb132c358aa790b891eed37253fa36.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2018/04/windows6.1-kb4103712-x86_9e9ca80634e4f94e95cc3a02eaea374b328f0f9d.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/06/windows6.1-kb4284867-x86_e841ad96c7b70dd96aae7720ed10fe56b70fd884.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/07/windows6.1-kb4338823-x86_4b18056251ec97112381473c933b7964c778d4ed.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2018/08/windows6.1-kb4343899-x86_0d9ef0cfeca3da376193a05088c6d172774378dc.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2018/09/windows6.1-kb4457145-x86_7b0fbc85360a4117eaae84c7088ccd09eac7527f.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/09/windows6.1-kb4462915-x86_bfcfa4c0997862cd2c0f8cd3df6f38bdacf6d07b.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/11/windows6.1-kb4467106-x86_e50f03c417cbe4ec2acbdcd6cd609c23bbc656fc.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2018/12/windows6.1-kb4471328-x86_bf41fb711ea06d24ac27361bba39d3bee0aa56a2.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/01/windows6.1-kb4480960-x86_dc89957c2ba506cef74cdf6760dc73237a067b9a.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/01/windows6.1-kb4486564-x86_4b0702863cf9aeea96f06ebb99778922019b7ff4.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/03/windows6.1-kb4489885-x86_8078e687b908bf6319d77d48fc2f70e0f67dfcf5.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2019/04/windows6.1-kb4493448-x86_187831a12093488fb2fc5be81af26f8f8d0fb386.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2019/06/windows6.1-kb4503269-x86_525652cb7e59c7ec922ff4e7efc60426d10cbe14.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2019/06/windows6.1-kb4507456-x86_41556c1452fcaadce2984d9e4ee9fe6068f38e29.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/09/windows6.1-kb4516655-x86_47655670362e023aa10ab856a3bda90aabeacfe6.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/09/windows6.1-kb4474419-v3-x86_0f687d50402790f340087c576886501b3223bec6.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2019/08/windows6.1-kb4512486-x86_4c88f71af8e9d07e5fb141d7aed0bcc7f532781e.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/09/windows6.1-kb4516033-x86_a6edd53b7a26c1dd6f74e10a57554dd082e7b4ae.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2019/08/ie11-windows6.1-kb4516046-x86_db07434ee4fbab99dc3522f02bef58b8a2cc30d3.msu"
                    )
      }
    8 {
        #Windows 7 or Windows 2008 R2 64 Bit KBs
        $MSUurls = ("http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/04/windows6.1-kb4093108-x64_fe804365f849cc61b133fda1efae299c534b830f.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/04/windows6.1-kb4103712-x64_44bc3455369066d70f52da47c30ca765f511cf68.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/06/windows6.1-kb4284867-x64_c2ecdf5620a36f257537e2e10c797f3ab572a7fe.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/07/windows6.1-kb4338823-x64_a141926d69d13f84e280086cb70b9b37dd590219.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/08/windows6.1-kb4343899-x64_09b367dfef2423a314f52325ce82d8675b2c5611.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/09/windows6.1-kb4457145-x64_b9404d9790106da7b6ee732a406f9d15a1b5242e.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/09/windows6.1-kb4462915-x64_63d42d3fb635f643f43e87d762b6077998735469.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/11/windows6.1-kb4467106-x64_ee54f25e11ccbb5d9eea964bbed2838583169ee5.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2018/12/windows6.1-kb4471328-x64_f9ae741bb45b98421d159469e57d765451a4d950.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/01/windows6.1-kb4480960-x64_bd23adfd0d82403d58aa8cd649636d136cf77700.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/01/windows6.1-kb4486564-x64_ad686ee44cfd554e461c55d1975d377b68af5eca.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/03/windows6.1-kb4489885-x64_3456932a9c8da3cde6a436d26f502126188332a0.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2019/04/windows6.1-kb4493448-x64_26274aef6de2f6b66e71f4a8ae51539238d1ec2d.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2019/06/windows6.1-kb4503269-x64_d518b12868bb1202a03fbc33c2d716092ae9c2e2.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/06/windows6.1-kb4507456-x64_6aa110cb2d01b8f291d1ea2c3cdc5e82204686ed.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/09/windows6.1-kb4516655-x64_8acf6b3aeb8ebb79973f034c39a9887c9f7df812.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/09/windows6.1-kb4474419-v3-x64_b5614c6cea5cb4e198717789633dca16308ef79c.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2019/08/windows6.1-kb4512486-x64_547fe7e4099c11d494c95d1f72e62a693cd70441.msu",
                    "http://download.windowsupdate.com/c/msdownload/update/software/secu/2019/09/windows6.1-kb4516033-x64_976486f9defe12ce403bdacfd932cb6f97540f0e.msu",
                    "http://download.windowsupdate.com/d/msdownload/update/software/secu/2019/08/ie11-windows6.1-kb4516046-x64_58cac9692fd89d6c502a7b78369fc993cb2bda7f.msu"
                    )
      }
}

Function download_msu{
       param($url)
       $wc = new-object System.Net.WebClient
       $MSUDest = "$DownloadPath\"+$($msuregex.Match($url).Groups[0].Value)
       try{
            $wc.DownloadFile($url,$MSUDest)
            return "Success"
       }catch { return "Could not download Security-only update $KB : $_.Exception.Message" }
}

Function Install_msu{
       param($msu,$kbid)
       if(!(test-path $logfilePath)){ New-Item -Itemtype Directory -Force -Path $logfilePath  | Out-Null }
       $logfile = $logfilePath+$kbid+".evtx"
       $ouput = start-process "wusa.exe" -ArgumentList "$msu /quiet /norestart /log:$logfile" -Wait -PassThru
       start-sleep -s 3
       if(Get-HotFix $kbid -ErrorAction SilentlyContinue){ return 'Success'}
       Else{ Write-output "FAILED : Security-only update $KB installation. See logfile $logfile for more details" }
}

$Mrollups = ("KB4093118","KB4103718","KB4284826","KB4338818","KB4343900",
             "KB4457144","KB4462923","KB4467107","KB4471318","KB4480970",
             "KB4486563","KB4489878","KB4493472","KB4499164","KB4503292",
             "KB4507449","KB4512506")

foreach ($mrollup in $Mrollups) {
    $rollup = Get-HotFix $mrollup -ErrorAction SilentlyContinue
    if ($rollup){ break }
}

$MSUurls = [System.Collections.Generic.List[System.Object]]$MSUurls
if ( $rollup.HotFixID -eq "KB4093118" ) { $MSUurls.RemoveRange(0,1) }
if ( $rollup.HotFixID -eq "KB4103718" ) { $MSUurls.RemoveRange(0,2) }
if ( $rollup.HotFixID -eq "KB4284826" ) { $MSUurls.RemoveRange(0,3) }
if ( $rollup.HotFixID -eq "KB4338818" ) { $MSUurls.RemoveRange(0,4) }
if ( $rollup.HotFixID -eq "KB4343900" ) { $MSUurls.RemoveRange(0,5) }
if ( $rollup.HotFixID -eq "KB4457144" ) { $MSUurls.RemoveRange(0,6) }
if ( $rollup.HotFixID -eq "KB4462923" ) { $MSUurls.RemoveRange(0,7) }
if ( $rollup.HotFixID -eq "KB4467107" ) { $MSUurls.RemoveRange(0,8) }
if ( $rollup.HotFixID -eq "KB4471318" ) { $MSUurls.RemoveRange(0,9) }
if ( $rollup.HotFixID -eq "KB4480970" ) { $MSUurls.RemoveRange(0,10) }
if ( $rollup.HotFixID -eq "KB4486563" ) { $MSUurls.RemoveRange(0,11) }
if ( $rollup.HotFixID -eq "KB4489878" ) { $MSUurls.RemoveRange(0,12) }
if ( $rollup.HotFixID -eq "KB4493472" ) { $MSUurls.RemoveRange(0,13) }
if ( $rollup.HotFixID -eq "KB4499164" ) { $MSUurls.RemoveRange(0,13) }
if ( $rollup.HotFixID -eq "KB4503292" ) { $MSUurls.RemoveRange(0,14) }
if ( $rollup.HotFixID -eq "KB4507449" ) { $MSUurls.RemoveRange(0,15) }
if ( $rollup.HotFixID -eq "KB4512506" ) { $MSUurls.RemoveRange(0,15); $MSUurls.Removeat(2) } 

$kbregex = [regex]"-(kb.*?)-"
$msuregex = [regex]"(?!.*\/).+"
$MSUs = @()

# START Internet Explorer Patch Cutomization 
$IEVersion =  ([System.Version][System.Diagnostics.FileVersionInfo]::GetVersionInfo("$env:ProgramFiles\Internet Explorer\iexplore.exe").ProductVersion).Major # Major vesion of intenet explorer
$include = @()
foreach ($item in $MSUurls) {
    $KB = $kbregex.Match($item).Groups[1].Value
    # Condition: Checks the Internet Explorer version and only install patches for IE v11
    if ($IEVersion -ne 11 -and $Item -like "*ie11-*") {
        Write-Output "Security-only update $KB is only supported for Internet Explorer v11"
    }
    else {
        $include += $item
    }
}
$MSUurls = $include
# END Internet Explorer Patch Cutomization

foreach ($msuurl in $MSUurls){
   $KB = $kbregex.Match($msuurl).Groups[1].Value
   if(-not(Get-HotFix $KB -ErrorAction SilentlyContinue)){
           $status = download_msu -url $msuurl
           if($status -eq "Success"){
              $MSUs += "$DownloadPath\"+$($msuregex.Match($msuurl).Groups[0].Value)
           }Else {Write-output $status}
    }Else{
            
            #Additional check added for re-relesed KB4474419
            if ($KB -eq "KB4474419"){
                    if( -not ($(dism /online /get-packages | ? { $_ -match "KB4474419"}) -contains "Package Identity : Package_for_KB4474419~31bf3856ad364e35~amd64~~6.1.3.2") ){
                
                           $status = download_msu -url $msuurl
                           if($status -eq "Success"){
                              $MSUs += "$DownloadPath\"+$($msuregex.Match($msuurl).Groups[0].Value)
                           }Else {Write-output $status}
                
                    }else{Write-Output "Security-only update $KB installation : Already installed"}
                
            }else{Write-Output "Security-only update $KB installation : Already installed"}
  
           
         }

}

Foreach ($MSU in $MSUs){
   $KB = $kbregex.Match($msu).Groups[1].Value
   $IStatus = Install_msu -msu $MSU -kbid $KB
   if($IStatus -eq "Success"){Write-Output "Security-only update $KB installation : Success"}
   Else {Write-Output "Security-only update $KB installation : Failed"}
   if (Test-Path $MSU) { Remove-Item $MSU -ErrorAction SilentlyContinue }
}

