<#
Usecases
 
    .Script
    Delete Log files from folder "C:\ProgramData\LogMeIn"
    Author
    Nirav Sachora 
 
#>
 

$folderpath1 = "C:\ProgramData\LogMeIn\Dumps"
 

$ErrorActionPreference = "SilentlyContinue"

    if (!(Test-path "C:\ProgramData\LogMeIn\Dumps")) {
        "Could not find the path $folderpath1"
        return
    }
    $filecount1 = (Get-ChildItem C:\ProgramData\LogMeIn\Dumps).Count
 
    if (!$filecount1) {
        "`nNo Files present at the directory $folderpath1`n"
        Exit;
    }
    else {
        Get-ChildItem $folderpath1 | Remove-item -Recurse -Force 
        $postfilecount1 = (Get-ChildItem C:\ProgramData\LogMeIn\Dumps).Count
    }
    
    if (!((Get-ChildItem C:\ProgramData\LogMeIn\Dumps).Count)) {
        "`nDump files deleted successfully."
	Exit;
    }
    elseif ($postfilecount1 -eq $filecount1) {
        "`nOperation Failed`n"
        Exit;
    }
    elseif ($postfilecount1) {
        $deletedfilecount1 = $filecount1 - $postfilecount1
        "`nNumber of files before deletion : $filecount1`nNote :Files which are in use will not be deleted.`nNumber of files deleted : $deletedfilecount1`nFailed to delete files : $postfilecount1`n"
        Exit;
    }
    


