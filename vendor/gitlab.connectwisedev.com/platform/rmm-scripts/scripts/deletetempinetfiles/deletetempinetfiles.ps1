$MajorOsVersion = ([environment]::OSVersion.Version).Major
$MinorOsVersion = ([environment]::OSVersion.Version).Minor
 
$LogOnUser = Get-ChildItem -Path $ENV:SystemDrive\Users\ | ?{$_ -notlike "*All*" -and $_ -notlike "*Default*" -and $_ -notlike "*public*" -and $_ -notlike "*Classic*"}
$LogOnUser = $LogOnUser | Select-Object NAME | ForEach-Object{$_.Name}

If ($LogOnUser){
    ForEach($AU in $LogOnUser){
    	$LiteralPathPart = $ENV:SYSTEMDRIVE + "\USERS\" + $AU + "\AppData\Local"
    	If (($MajorOsVersion -eq 6) -And ($MinorOsVersion -eq 1)){
    	    $PathPart = "Temporary Internet Files"
    	} ElseIf ($MajorOsVersion -eq 10){  $PathPart = "INetCache"	}  

    	$ItemsDeleted = 0
    	$Length = 0
        $LengthNew = 0
        $ItemsDeleted = 0
    	$LiteralPath = "$LiteralPathPart\Microsoft\Windows\$PathPart","$LiteralPathPart\Temp"
    	$Items = (Get-ChildItem -LiteralPath $LiteralPath -Recurse -Force -ErrorAction SilentlyContinue)

    	If (($Items.Count -ne 0) -And ($Items -ne $null)){
    		$ContentBefore = Get-ChildItem -Path $LiteralPath -Force -Recurse -ErrorAction SilentlyContinue | Where-Object { -not $_.PSIsContainer } | Measure-Object -Property Length -Sum
    		If ($ContentBefore){
    			$Length += $ContentBefore.Sum
    		}

    		$Items | ForEach {
        		Remove-Item $_.FullName -Recurse -Force -ErrorAction SilentlyContinue
    	    	If ((-Not(Test-Path $_.FullName -ErrorAction SilentlyContinue)) -And (-Not($Error[0].Exception -is [System.UnauthorizedAccessException]))){
    		    	$ItemsDeleted++
    	    	}
    	    }

    	    $ContentAfter = Get-ChildItem -Path $LiteralPath -Force -Recurse -ErrorAction SilentlyContinue | Where-Object { -not $_.PSIsContainer } | Measure-Object -Property Length -Sum
    	    If ($ContentAfter){
    		    $LengthNew += $ContentAfter.Sum
    	    }

    	    If($ItemsDeleted -eq 0){
                "User Name : " + $AU
                Write-Output "$($Items.Count) file(-s) for User $AU has(-ve) not been deleted because of access rights or they were using by another process"
                "`n"		  
    	    }

    	    If(($Items.Count-$ItemsDeleted) -eq 0){
                "User Name : " + $AU
    		    Write-Output "Cleared $($Length-$LengthNew) bytes, $ItemsDeleted file(-s)"
                "`n"		  
    	    }

    	    If(($Items.Count-$ItemsDeleted) -gt 0){
                "User Name : " + $AU
                Write-Output "Cleared $($Length-$LengthNew) bytes, $ItemsDeleted file(-s); $($Items.Count-$ItemsDeleted) file(-s) for User $AU has(-ve) not been deleted because of access rights or they were using by another process"
                "`n"		  
    	    }
        }
    }

} Else {
    Write-Error "No Currently Logged-On User"
}

