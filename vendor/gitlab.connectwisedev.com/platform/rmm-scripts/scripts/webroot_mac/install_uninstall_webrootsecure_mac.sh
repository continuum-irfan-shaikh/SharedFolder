#!/bin/sh -e

if [ "$action" = uninstall ]
then
    sudo launchctl unload /Library/LaunchDaemons/com.webroot.security.mac.plist
    sudo kextunload /System/Library/Extensions/SecureAnywhere.kext
    sudo rm /usr/local/bin/WSDaemon
    sudo rm /usr/local/bin/WFDaemon
    sudo killall -9 WSDaemon
    sudo killall -9 WfDaemon
    sudo killall -9 "Webroot SecureAnywhere"
    sudo rm -rf /System/Library/Extensions/SecureAnywhere.kext
    sudo rm -rf "/Applications/Webroot SecureAnywhere.app"
    sudo rm /Library/LaunchAgents/com.webroot.WRMacApp.plist
    sudo rm /Library/LaunchDaemons/com.webroot.security.mac.plist
    echo "Webroot SecureAnywhere has successfully been uninstalled"
    exit 0
fi

BASE_DIR="/opt/local/share/continuum_mac"
webroot_path="$BASE_DIR/downloads"
dmg_path="$webroot_path/webroot"
web_atch_dtch="/Volumes/Webroot SecureAnywhere/"

#checking internet availability
if eval "ping -c 1 www.google.com >/dev/null"; then
echo "  "
else
echo "Could not install webroot, as there is no active internet connection."
exit 1
fi

#Checking webroot if it is already installed
if [ -f /Library/LaunchDaemons/com.webroot.security.mac.plist ] || [ -f /Library/LaunchAgents/com.webroot.WRMacApp.plist ] || [ -f /Applications/Webroot\ SecureAnywhere.app ]
then
    echo "Installation Failed : webroot already installed"
    exit 1
fi

if [ -d /Volumes/Webroot\ SecureAnywhere/ ]
then
echo " Installation Failed : Webroot already attach kindly detach it and run the script again"
exit 1
fi

#to check whether directory is present or not
if [ ! -d $dmg_path ]
then
    mkdir $dmg_path
fi

#Download WSA mac client
cd $dmg_path; curl -O http://anywhere.webrootcloudav.com/zerol/wsamacsme.dmg

#mount the dmg
hdiutil attach $dmg_path/wsamacsme.dmg

ditto /Volumes/Webroot\ SecureAnywhere /Applications/

#Trigger the webroot Installation.
sudo "/Applications/Webroot SecureAnywhere.app/Contents/MacOS/Webroot SecureAnywhere" install -keycode="$licensecode" -silent

retVal=$?
echo $retVal

if [ $retVal -ne 0 ]
then
    echo "Installation Failed : webroot installation failed"
    exit 5
#else
# echo " Installation Succesfully : Webroot installation successful"
fi


webrootfile="/Applications/Webroot*.app"

if [ -e $webrootfile ]
then
    echo "Installation Succesfully : Webroot installed successfully"
else
    exit 6
fi

#Unmount the DMG

hdiutil detach /Volumes/Webroot\ SecureAnywhere/

rm -rf $dmg_path/*

exit 0