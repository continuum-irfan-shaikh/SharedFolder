#!/bin/bash

if [ -f ../../settings.cfg ];then 
	. ../../settings.cfg
fi

# wrk job
wrk -t"$threads" -c"$connections" -d"$duration s" -R"$frequency" -s ./script.lua $url_local