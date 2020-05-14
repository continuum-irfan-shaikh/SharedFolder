<#
    .SYNOPSIS
         Set Default Internet Explorer Search Provider
    .DESCRIPTION
         Set Default Internet Explorer Search Provider. Script will set search providers like google/yahoo/bing/wikipedia/graysonline/arena as default search provider.
    .Help
         
    .Author
        Durgeshkumar Patel
    .Version
        1.0
#>

#JSON Schema variable
#$SearchProvider = "google"      Examples:- google/yahoo/bing/wikipedia/graysonline/arena
############################################################

$ErrorActionPreference = "SilentlyContinue"

if ($env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    if ($myInvocation.Line) {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile $myInvocation.Line
    }
    else {
        &"$env:systemroot\sysnative\windowspowershell\v1.0\powershell.exe" -NonInteractive -NoProfile -file "$($myInvocation.InvocationName)" $args
    }
    exit $lastexitcode
}

if (!$SearchProvider) {

    Write-Output "`nKindly provide search provider name. It can't be null."
    EXIT;
}

# Get each user profile SID and Path to the profile
$UserProfiles = Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList\*" | where { $_.PSChildName -match "S-1-5-21-(\d+-?){4}$" } | Select-Object @{Name = "SID"; Expression = { $_.PSChildName } }, @{Name = "UserHive"; Expression = { "$($_.ProfileImagePath)\NTuser.dat" } }, @{name = "UserName"; Expression = { ForEach-Object { $_.ProfileImagePath -replace "c:\\users\\", "" } } }

$UserProfiles += @( New-Object PSObject -Property @{SID = ".DEFAULT"; UserHive = "C:\Users\Public\ntuser.dat"; UserName = ".DEFAULT" })

function check ($guid1, $guid2) { 

    if ($guid1 -eq $guid2) {
        return $true
    }
    else {
        return $false
    }
}

function google {
    If (Test-Path ($RegKey + "\Internet Explorer")) {
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        If (Test-Path ($RegKey + "\SearchScopes")) {
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
            
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "Google" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "https://www.google.com/search?q={searchTerms}&sourceid=ie7&rls=com.microsoft:{language}:{referrer:source}&ie={inputEncoding?}&oe={outputEncoding?}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null   
            New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "https://www.google.com/complete/search?q={searchTerms}&client=ie8&mw={ie:maxWidth}&sh={ie:sectionHeight}&rh={ie:rowHeight}&inputencoding={inputEncoding}&outputencoding={outputEncoding}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWilsM" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURL" -value "https://www.google.com/favicon.ico" -PropertyType String -Force | Out-Null
              
        }
        else { 
            New-Item -path $RegKey -name "SearchScopes" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
           
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "Google" -PropertyType String -Force | Out-Null  
            New-ItemProperty -path $RegKey -name "URL" -value "https://www.google.com/search?q={searchTerms}&sourceid=ie7&rls=com.microsoft:{language}:{referrer:source}&ie={inputEncoding?}&oe={outputEncoding?}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "https://www.google.com/complete/search?q={searchTerms}&client=ie8&mw={ie:maxWidth}&sh={ie:sectionHeight}&rh={ie:rowHeight}&inputencoding={inputEncoding}&outputencoding={outputEncoding}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWilsM" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURL" -value "https://www.google.com/favicon.ico" -PropertyType String -Force | Out-Null
                  
        }
    }
    else {
        New-Item -path $RegKey -name "Internet Explorer" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        New-Item -path $RegKey -name "SearchScopes" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
       
        New-Item -path $RegKey -name "{$guid}" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
        New-ItemProperty -path $RegKey -name "DisplayName" -value "Google" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "URL" -value "https://www.google.com/search?q={searchTerms}&sourceid=ie7&rls=com.microsoft:{language}:{referrer:source}&ie={inputEncoding?}&oe={outputEncoding?}" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "https://www.google.com/complete/search?q={searchTerms}&client=ie8&mw={ie:maxWidth}&sh={ie:sectionHeight}&rh={ie:rowHeight}&inputencoding={inputEncoding}&outputencoding={outputEncoding}" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWilsM" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "FaviconURL" -value "https://www.google.com/favicon.ico" -PropertyType String -Force | Out-Null 
        
    } 
}

function Wikipedia {
    If (Test-Path ($RegKey + "\Internet Explorer")) {
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        If (Test-Path ($RegKey + "\SearchScopes")) {
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
            
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "Wikipedia (en)" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "http://en.wikipedia.org/w/index.php?title=Special:Search&search={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null   
            New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "http://en.wikipedia.org/w/api.php?action=opensearch&format=xml&search={searchTerms}&namespace=0" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWihZ8" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "SuggestionsURL_JSON" -value "http://en.wikipedia.org/w/api.php?action=opensearch&search={searchTerms}&namespace=0" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://en.wikipedia.org/favicon.ico" -PropertyType String -Force | Out-Null
              
        }
        else { 
            New-Item -path $RegKey -name "SearchScopes" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
           
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "Wikipedia (en)" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "http://en.wikipedia.org/w/index.php?title=Special:Search&search={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null   
            New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "http://en.wikipedia.org/w/api.php?action=opensearch&format=xml&search={searchTerms}&namespace=0" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWihZ8" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "SuggestionsURL_JSON" -value "http://en.wikipedia.org/w/api.php?action=opensearch&search={searchTerms}&namespace=0" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://en.wikipedia.org/favicon.ico" -PropertyType String -Force | Out-Null
                  
        }
    }
    else {
        New-Item -path $RegKey -name "Internet Explorer" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        New-Item -path $RegKey -name "SearchScopes" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
       
        New-Item -path $RegKey -name "{$guid}" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
        New-ItemProperty -path $RegKey -name "DisplayName" -value "Wikipedia (en)" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "URL" -value "http://en.wikipedia.org/w/index.php?title=Special:Search&search={searchTerms}" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null   
        New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "http://en.wikipedia.org/w/api.php?action=opensearch&format=xml&search={searchTerms}&namespace=0" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWihZ8" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "SuggestionsURL_JSON" -value "http://en.wikipedia.org/w/api.php?action=opensearch&search={searchTerms}&namespace=0" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://en.wikipedia.org/favicon.ico" -PropertyType String -Force | Out-Null 
        
    } 
}

function Arena {
    If (Test-Path ($RegKey + "\Internet Explorer")) {
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        If (Test-Path ($RegKey + "\SearchScopes")) {
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
            
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "The-Arena" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "http://www.the-arena.co.uk/default.aspx?sc={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null   
            New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "http://www.the-arena.co.uk/autosuggestion.aspx?sc={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWidzo" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://www.the-arena.co.uk/favicon.ico" -PropertyType String -Force | Out-Null
            New-ItemProperty -path $RegKey -name "SuggestionsURL_JSON" -value "http://www.the-arena.co.uk/autosuggestionjson.aspx?sc={searchTerms}" -PropertyType String -Force | Out-Null 
              
        }
        else { 
            New-Item -path $RegKey -name "SearchScopes" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
           
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "The-Arena" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "http://www.the-arena.co.uk/default.aspx?sc={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null   
            New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "http://www.the-arena.co.uk/autosuggestion.aspx?sc={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWidzo" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://www.the-arena.co.uk/favicon.ico" -PropertyType String -Force | Out-Null
            New-ItemProperty -path $RegKey -name "SuggestionsURL_JSON" -value "http://www.the-arena.co.uk/autosuggestionjson.aspx?sc={searchTerms}" -PropertyType String -Force | Out-Null 
                  
        }
    }
    else {
        New-Item -path $RegKey -name "Internet Explorer" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        New-Item -path $RegKey -name "SearchScopes" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
       
        New-Item -path $RegKey -name "{$guid}" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
        New-ItemProperty -path $RegKey -name "DisplayName" -value "The-Arena" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "URL" -value "http://www.the-arena.co.uk/default.aspx?sc={searchTerms}" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null   
        New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "http://www.the-arena.co.uk/autosuggestion.aspx?sc={searchTerms}" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWidzo" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://www.the-arena.co.uk/favicon.ico" -PropertyType String -Force | Out-Null
        New-ItemProperty -path $RegKey -name "SuggestionsURL_JSON" -value "http://www.the-arena.co.uk/autosuggestionjson.aspx?sc={searchTerms}" -PropertyType String -Force | Out-Null  
        
    } 
}

function GraysOnline {
    If (Test-Path ($RegKey + "\Internet Explorer")) {
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        If (Test-Path ($RegKey + "\SearchScopes")) {
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
            
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "GraysOnline" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "http://www.graysonline.com/Search.aspx?q={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null   
            New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "http://www.graysonline.com/OpenSearch/SuggestionsXml.aspx?q={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWiaTt" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://images.graysonline.com/favicon.ico" -PropertyType String -Force | Out-Null
            
              
        }
        else { 
            New-Item -path $RegKey -name "SearchScopes" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
           
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "GraysOnline" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "http://www.graysonline.com/Search.aspx?q={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null   
            New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "http://www.graysonline.com/OpenSearch/SuggestionsXml.aspx?q={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWiaTt" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://images.graysonline.com/favicon.ico" -PropertyType String -Force | Out-Null 
                  
        }
    }
    else {
        New-Item -path $RegKey -name "Internet Explorer" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        New-Item -path $RegKey -name "SearchScopes" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
       
        New-Item -path $RegKey -name "{$guid}" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
        New-ItemProperty -path $RegKey -name "DisplayName" -value "GraysOnline" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "URL" -value "http://www.graysonline.com/Search.aspx?q={searchTerms}" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null   
        New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "http://www.graysonline.com/OpenSearch/SuggestionsXml.aspx?q={searchTerms}" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWiaTt" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://images.graysonline.com/favicon.ico" -PropertyType String -Force | Out-Null  
        
    } 
}

function yahoo {
    If (Test-Path ($RegKey + "\Internet Explorer")) {
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        If (Test-Path ($RegKey + "\SearchScopes")) {
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
            
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "Yahoo!" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "https://search.yahoo.com/search?p={searchTerms}&b={startPage?}&fr=ie8" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "https://ie.search.yahoo.com/os?appid=ie8&command={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWilsN" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://search.yahoo.com/favicon.ico" -PropertyType String -Force | Out-Null 
                          
        }
        else {
            New-Item -path $RegKey -name "SearchScopes" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
            
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "Yahoo!" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "https://search.yahoo.com/search?p={searchTerms}&b={startPage?}&fr=ie8" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null  
            New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "https://ie.search.yahoo.com/os?appid=ie8&command={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWilsN" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://search.yahoo.com/favicon.ico" -PropertyType String -Force | Out-Null  
            
        }
    }
    else {
        New-Item -path $RegKey -name "Internet Explorer" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        New-Item -path $RegKey -name "SearchScopes" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
        
        New-Item -path $RegKey -name "{$guid}" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
        New-ItemProperty -path $RegKey -name "DisplayName" -value "Yahoo!" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "URL" -value "https://search.yahoo.com/search?p={searchTerms}&b={startPage?}&fr=ie8" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null  
        New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "https://ie.search.yahoo.com/os?appid=ie8&command={searchTerms}" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWilsN" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://search.yahoo.com/favicon.ico" -PropertyType String -Force | Out-Null         
        
    } 
}

function Bing {
    If (Test-Path ($RegKey + "\Internet Explorer")) {
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        If (Test-Path ($RegKey + "\SearchScopes")) {
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
            
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "Bing" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "http://www.bing.com/search?q={searchTerms}&src=IE-SearchBox&FORM=IE8SRC" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "SuggestionsURLFallback" -value "http://api.bing.com/qsml.aspx?query={searchTerms}&maxwidth={ie:maxWidth}&rowheight={ie:rowHeight}&sectionHeight={ie:sectionHeight}&FORM=IE8SSC&market={language}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURLFallback" -value "http://www.bing.com/favicon.ico" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "http://api.bing.com/qsml.aspx?query={searchTerms}&maxwidth={ie:maxWidth}&rowheight={ie:rowHeight}&sectionHeight={ie:sectionHeight}&FORM=IESS02&market={language}" -PropertyType String -Force | Out-Null 
                  
        }
        else {
            New-Item -path $RegKey -name "SearchScopes" | Out-Null
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
            
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "Bing" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "http://www.bing.com/search?q={searchTerms}&src=IE-SearchBox&FORM=IE8SRC" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "SuggestionsURLFallback" -value "http://api.bing.com/qsml.aspx?query={searchTerms}&maxwidth={ie:maxWidth}&rowheight={ie:rowHeight}&sectionHeight={ie:sectionHeight}&FORM=IE8SSC&market={language}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURLFallback" -value "http://www.bing.com/favicon.ico" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "http://api.bing.com/qsml.aspx?query={searchTerms}&maxwidth={ie:maxWidth}&rowheight={ie:rowHeight}&sectionHeight={ie:sectionHeight}&FORM=IESS02&market={language}" -PropertyType String -Force | Out-Null     
                          
        }
    }
    else {
        New-Item -path $RegKey -name "Internet Explorer" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        New-Item -path $RegKey -name "SearchScopes" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
        
        New-Item -path $RegKey -name "{$guid}" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
        New-ItemProperty -path $RegKey -name "DisplayName" -value "Bing" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "URL" -value "http://www.bing.com/search?q={searchTerms}&src=IE-SearchBox&FORM=IE8SRC" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "ShowSearchSuggestions" -value '1' -PropertyType dword -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "SuggestionsURLFallback" -value "http://api.bing.com/qsml.aspx?query={searchTerms}&maxwidth={ie:maxWidth}&rowheight={ie:rowHeight}&sectionHeight={ie:sectionHeight}&FORM=IE8SSC&market={language}" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "FaviconURLFallback" -value "http://www.bing.com/favicon.ico" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "SuggestionsURL" -value "http://api.bing.com/qsml.aspx?query={searchTerms}&maxwidth={ie:maxWidth}&rowheight={ie:rowHeight}&sectionHeight={ie:sectionHeight}&FORM=IESS02&market={language}" -PropertyType String -Force | Out-Null           
        
    } 
}

function YouTube {
    If (Test-Path ($RegKey + "\Internet Explorer")) {
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        If (Test-Path ($RegKey + "\SearchScopes")) {
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
            
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "YouTube" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "http://www.youtube.com/results?search_query={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWiaTu" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://s.ytimg.com/yt/favicon-vfl1123.ico" -PropertyType String -Force | Out-Null
        }
        else { 
            New-Item -path $RegKey -name "SearchScopes" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
           
            New-Item -path $RegKey -name "{$guid}" | Out-Null 
            $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
            New-ItemProperty -path $RegKey -name "DisplayName" -value "YouTube" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "URL" -value "http://www.youtube.com/results?search_query={searchTerms}" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWiaTu" -PropertyType String -Force | Out-Null 
            New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://s.ytimg.com/yt/favicon-vfl1123.ico" -PropertyType String -Force | Out-Null
        }
    }
    else {
        New-Item -path $RegKey -name "Internet Explorer" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer"
        New-Item -path $RegKey -name "SearchScopes" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
       
        New-Item -path $RegKey -name "{$guid}" | Out-Null 
        $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\{$guid}"
        New-ItemProperty -path $RegKey -name "DisplayName" -value "YouTube" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "URL" -value "http://www.youtube.com/results?search_query={searchTerms}" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "OSDFileURL" -value "https://www.microsoft.com/cms/api/am/binary/RWiaTu" -PropertyType String -Force | Out-Null 
        New-ItemProperty -path $RegKey -name "FaviconURL" -value "http://s.ytimg.com/yt/favicon-vfl1123.ico" -PropertyType String -Force | Out-Null     
    } 
}

function setguid { 

    $SearchGUIDSet = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes"
    #Test-Path $SearchGUIDSet
    New-ItemProperty -path $SearchGUIDSet -name "DefaultScope" -value "{$guid}" -PropertyType String -Force | Out-Null
    New-ItemProperty -path "registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Internet Explorer\SearchScopes\" -name "DefaultScope" -value "{$guid}" -PropertyType String -Force | Out-Null
    
}


if ($SearchProvider -eq "bing") { 
    $global:guid = "0633EE93-D776-472f-A0FF-E1416B8B2E3A"
}
else { 
    $global:guid = [guid]::NewGuid().ToString().ToUpper() 
}


# Loop through each profile on the machine</p>
Foreach ($UserProfile in $UserProfiles) {
    # Load User ntuser.dat if it's not already loaded
    If (($ProfileWasLoaded = Test-Path Registry::HKEY_USERS\$($UserProfile.SID)) -eq $false) {
        Start-Process -FilePath "CMD.EXE" -ArgumentList "/C REG.EXE LOAD HKU\$($UserProfile.SID) $($UserProfile.UserHive)" -Wait -WindowStyle Hidden
    } 
    $RegKey = "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft"
    
    ######Code Here for REGISTRY
    if (Test-Path "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\*" -ErrorAction 'SilentlyContinue') {
        Remove-Item -Path "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes\*" -Recurse -Force -ErrorAction 'SilentlyContinue'
       
    }

    if (Test-Path "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\User Preferences") {
        
        Remove-Item -Path "registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\User Preferences" -Recurse -Force
    }
    
    if (Test-Path "registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Internet Explorer\SearchScopes\*" -ErrorAction 'SilentlyContinue') {
        Remove-Item -Path "registry::HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Internet Explorer\SearchScopes\*" -Recurse -Force -ErrorAction 'SilentlyContinue'
       
    }
     
    
    switch ($SearchProvider) {
    
        "google" { 
            #call google
            
            google
            setguid
            break;
        }
        "bing" {
            #call bing
            
            bing
            setguid
            break;
        }
        "yahoo" {
            #call yahoo
            
            yahoo
            setguid   
            break;
        }
        "yahoo!" {
            #call yahoo
            
            yahoo
            setguid
            break;
        }
        "Wikipedia" {
        
            #call wikipedia
            Wikipedia
            setguid
            break;
        }
        "GraysOnline" {

            #call GraysOnline
            GraysOnline
            setguid
            break; 
        }
        "Arena" {

            #call Arena
            Arena
            setguid
            break; 
        }
        "youtube" {

            #call Arena
            youtube
            setguid
            break;
        }
        Default {
            Write-Output "`n'$SearchProvider' search provider not supported."
            Write-Output "`Supported search providers are Bing,Yahoo,Google,Wikipedia,GraysOnline,YouTube."
            EXIT;
        }
    }
     

    if ((check -guid1 $("registry::HKEY_USERS\$($UserProfile.SID)\Software\Microsoft\Internet Explorer\SearchScopes" | Get-Item | Get-ItemProperty | select -ExpandProperty DefaultScope) -guid2 "{$guid}") -eq $true) {

        Write-Output "$SearchProvider set as a default search provider for user $($UserProfile.UserName)"
    }
    else {
        Write-Output "Failed to set $SearchProvider as a default search provider for user $($UserProfile.UserName)"
    }
    #####Code complete for registry

    # Unload NTuser.dat        
    If ($ProfileWasLoaded -eq $false) {
        [gc]::Collect()
        Start-Sleep 1
        Start-Process -FilePath "CMD.EXE" -ArgumentList "/C REG.EXE UNLOAD HKU\$($UserProfile.SID)" -Wait -WindowStyle Hidden | Out-Null
    }
  
} 

Write-Output "`nNote:- User will get the popup for deafult search provider change for Internet Explorer v11 or greater where systems are not a part of any domain."


