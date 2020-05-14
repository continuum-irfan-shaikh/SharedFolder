#!/bin/bash

if [ -f ../../settings.cfg ];then 
	. ../../settings.cfg
fi

# apply settings
sed -i 's|regularity|'$regularity'|g' script.lua
sed -i 's|template|'$runtime'|g' ./templates/task_"$regularity"_template.json

# wrk job
wrk -t"$threads" -c"$connections" -d"$duration s" -R"$frequency" -s ./script.lua $url_local

# cleaning up
sed -i 's|'$regularity'|regularity|g' script.lua
sed -i 's|'$runtime'|template|g' ./templates/task_"$regularity"_template.json
