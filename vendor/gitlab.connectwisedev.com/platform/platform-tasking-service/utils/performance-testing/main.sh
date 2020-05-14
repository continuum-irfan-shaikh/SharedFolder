#!/bin/bash

root=$( pwd )

methods=("task-definitions/Create"
        "task-definitions/GetByID"
        "task-definitions/GetByPartnerID"
        "tasks/Create"
        "tasks/GetByPartnerAndID"
        "tasks/GetByPartnerAndTarget"
        "execution-results/Count"
        "execution-results/Get"
        "execution-results/History"
        "templates/GetAll"
        "templates/GetByOriginID"
        "templates/GetByType"
        "task-definitions/DeleteByID")

source settings.cfg

for method in ${methods[@]}
do
        echo "Started testing of "$method
        echo $method"" >> "./stat_"$frequency".log"
        echo "----------------------" >> "./stat_"$frequency".log"

        pushd $method
        sh ./wrk.sh >> $root/"stat_"$frequency".log"
        popd

        echo "----------------------" >> "./stat_"$frequency".log"
        echo "" >> "./stat_"$frequency".log"
       	echo "Cassandra is resting for 15 sec"

        # wait for cassandra gets calm
        sleep 15
done