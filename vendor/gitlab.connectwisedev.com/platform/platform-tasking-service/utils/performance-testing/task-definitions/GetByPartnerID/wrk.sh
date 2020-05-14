#!/bin/bash

if [ -f ../../settings.cfg ];then 
	. ../../settings.cfg
fi

######  wrk job  #########
wrk -t"$threads" -c"$connections" -d"$duration s" -R"$frequency" $url_local"/tasking/v1/partners/"$partner_id"/task-definitions"