#!/bin/bash
#Ubuntu Service Upstart script commands

service="platform-tasking-service"
cnt=$(ps -ef | grep -v grep| grep $service |wc -l)
if test $cnt -gt 0
then
service $service stop
fi

release=`lsb_release -r -s | grep -oE "^[0-9]+"`
syslogConf="/etc/rsyslog.d/$service.conf"

#Install latest unzip bin and extract continuum.zip
if ! which unzip ; then
    apt-get install unzip
fi

unzip -o continuum.zip -d /opt

#Update Configuration File
. environment.sh $@ -f /opt/continuum/config/ctm_tasking_cfg.json

if (( "$release" >= "18" )); then
    cp /opt/continuum/config/platform-tasking-service.service /etc/systemd/system/
    echo "if \$programname == '$service' then /var/log/$service/$service.log
& stop" > $syslogConf
    systemctl restart rsyslog
    systemctl daemon-reload
    #Start agent systemd service
    systemctl start $service
    systemctl enable $service

else
    #Create symlink
    ln -sf /opt/continuum/config/platform-tasking-service.conf /etc/init



    #Reload configuration
    initctl reload-configuration

    #Start service
    initctl start $service
fi