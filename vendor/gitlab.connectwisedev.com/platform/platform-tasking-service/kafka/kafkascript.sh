#!/bin/bash
#This Scipt will be used to deploy kafka from Jenkins

AllParams=$@

AllReqParams=""

echo "$AllParams"
echo "$@"

while [ -n "$1" ]
do
 case "$1" in
     --topic)Topic=$2;AllReqParams+=" "$1; shift;;
     --zookeeper)Zookeeper=$2;AllReqParams+=" "$1; shift;;
     --kafkapath)KafkaPath=$2;shift;shift;;
     --partitions)Partitions=$2;AllReqParams+=" "$1;shift;;
     --replication-factor)RepFactor=$2;AllReqParams+=" "$1;shift;;
     *)AllReqParams+=" "$1; shift;;

  esac
done

#DON'T CHANGE NEXT TWO LINES MANUALLY
#Topic must be set during build (not deployment) time
Topic="managed_endpoint_change"

if [ "$Topic" = "" ]
        then
                echo "Command Parameter missing, please provide 'topic'"
                exit 1
fi

if [ "$Zookeeper" = "" ]
        then
                echo "Command Parameter missing, please provide 'zookeeper'"
                exit 1
fi

if [ "$KafkaPath" = "" ]
        then
                echo "Command Parameter missing, please provide 'kafkapath'"
                exit 1
fi

if [ "$Partitions" = "" ]
	then
		echo  "Command Parameter missing, please provide 'partitions'"
		exit 1
fi

if [ "$RepFactor" = "" ]
        then
                echo  "Command Parameter missing, please provide 'replication-factor'"
		exit 1
fi


Path="/bin/kafka-topics.sh"

ListTopics=$KafkaPath$Path" --list --zookeeper "$Zookeeper

SearchTopicsResult=`$ListTopics|grep -w "$Topic"`

echo "Searching topic..." $SearchTopicsResult

Empty=""
CreateNewTopic=$KafkaPath$Path" --create "$AllReqParams

if [ $SearchTopicsResult ]
        then
                echo "Topic already exists"
else
	echo "Topic not found, creating Topic"
	CreateTopicResult=`$CreateNewTopic`
	echo "$CreateTopicResult"
fi
