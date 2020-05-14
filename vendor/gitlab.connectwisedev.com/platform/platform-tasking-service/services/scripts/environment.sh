#!/bin/bash
# Tasking Micro-Service Configuration Update script for Environment based on parameters

ListenURL="{{service_port}}"
URLPrefix="{{service_uri}}"
APIVersion="{{service_version}}"
LogLevel="{{log_level}}"
LogMaxFileSizeInMB="{{log_size | default('100')}}"
LogOldFileToKeep="{{log_files_to_keep | default('5')}}"

CassandraURL="{{cassandra_hosts}}"
ZookeeperHosts="{{zookeeper_hosts}}"
KafkaBrokers="{{kafka_host}}"

TaskingMsURL="{{tasking_elb}}"
ScriptingMsURL="{{scripting_elb}}"
AgentConfigMsURL="{{agent_config_elb}}"
SequenceMsURL="{{sequence_elb}}"
PatchingMsURL="{{patching_elb}}"
AlertingMsURL="{{alerting_elb}}"
ProfilingMsURL="{{profiling_elb}}"
WebrootMsURL="{{webroot_elb}}"
AssetMsURL="{{asset_elb}}"
DynamicGroupsMsURL="{{dg_elb}}"
GraphQLMsURL="{{graphql_elb}}"
MemcachedURL="{{memcached_hosts}}"
EntitlementMsURL="{{entitlement_elb}}"
SitesMsURL="{{sites_elb}}"
SitesNoTokenURL="{{sites_no_token_host}}"
AutomationEngineURL="{{automation_engine_elb}}"
AgentServiceURL="{{agent_service_elb}}"
EncryptionKey="{{encryption_key}}"

CassandraTimeoutSec="{{ cassandra_timeout_sec | default('30') }}"
CassandraConnNumber="{{ cassandra_conn_number | default('20') }}"
CassandraConcurrentCallNumber="{{ cassandra_concurrent_call_number | default('30') }}"
CassandraBatchSize="{{ cassandra_batch_size | default('5') }}"
InMemoryCacheSize="{{ inmemory_cache_size | default('1073741274') }}"
HTTPClientResultsTimeoutSec="{{ http_timeout_results_sec | default('120') }}"
ExecutionResultKafkaTopic="{{ execution_results_topic | default('script_execution_result') }}"
ClosestTasksWorkersTimeoutSec="{{ closest_tasks_workers_timeout_sec | default('15') }}"

LogFile=/opt/continuum/log/ctm_tasking_service.log
FILE_PATH=/opt/continuum/config/ctm_tasking_cfg.json
BrokerHosts=$(echo $KafkaBrokers | sed 's/,/","/g')


echo "====== Environment Variables ====="
echo "Log Level           = " $LogLevel
echo "LogMaxFileSizeInMB  = " $LogMaxFileSizeInMB
echo "LogOldFileToKeep    = " $LogOldFileToKeep
echo "Listen Port         = " $ListenURL
echo "URLPrefix           = " $URLPrefix
echo "API Version         = " $APIVersion
echo "CassandraURL        = " $CassandraURL
echo "Kafka Brokers       = " $KafkaBrokers
echo "File Path           = " $FILE_PATH
echo "Zookeeper Hosts     = " $ZookeeperHosts
echo "ScriptingMsURL      = " $ScriptingMsURL
echo "SequenceMsURL       = " $SequenceMsURL
echo "TaskingMsURL        = " $TaskingMsURL
echo "AssetMsURL          = " $AssetMsURL
echo "DynamicGroupsMsURL  = " $DynamicGroupsMsURL
echo "GraphQLMsURL        = " $GraphQLMsURL
echo "MemcachedURL        = " $MemcachedURL
echo "EntitlementMsURL    = " $EntitlementMsURL
echo "AlertingMsURL       = " $AlertingMsURL
echo "PatchingMsURL       = " $PatchingMsURL
echo "ProfilingMsURL      = " $ProfilingMsURL
echo "WebrootMsURL        = " $WebrootMsURL
echo "SitesMsURL          = " $SitesMsURL
echo "SitesNoTokenURL     = " $SitesNoTokenURL
echo "AutomationEngineURL = " $AutomationEngineURL
echo "AgentConfigMsURL    = " $AgentConfigMsURL
echo "ExecutionResultKafkaTopic = " $ExecutionResultKafkaTopic
echo "AgentServiceURL     = " $AgentServiceURL
echo "EncryptionKey       = " $EncryptionKey
echo "===== Environment Variables ===="

function var_replace() {
  key=$1
  value=$2
  awk -v val="$value" "/$key/{\$2=val}1" $FILE_PATH > tmp && mv tmp $FILE_PATH
}

var_replace ListenURL "\":$ListenURL\","
var_replace URLPrefix "\"/$URLPrefix\","
var_replace APIVersion "\"/$APIVersion\","
var_replace CassandraURL "\"$CassandraURL\","
var_replace KafkaBrokers "\"$KafkaBrokers\","
var_replace BrokerHosts "[\"$BrokerHosts\"],"
var_replace ZookeeperHosts "\"$ZookeeperHosts\","
var_replace logLevel "\"$LogLevel\","
var_replace filename "\"$LogFile\","
var_replace script "\"$ScriptingMsURL\""
var_replace sequence "\"$SequenceMsURL\","
var_replace patching "\"$PatchingMsURL\","
var_replace alerting "\"$AlertingMsURL\","
var_replace profiling "\"$ProfilingMsURL\","
var_replace webroot "\"$WebrootMsURL\","
var_replace TaskingMsURL "\"$TaskingMsURL\","
var_replace ScriptingMsURL "\"$ScriptingMsURL\","
var_replace AssetMsURL "\"$AssetMsURL\","
var_replace DynamicGroupsMsURL "\"$DynamicGroupsMsURL\","
var_replace GraphQLMsURL "\"$GraphQLMsURL\","
var_replace MemcachedURL "\"$MemcachedURL\","
var_replace EntitlementMsURL "\"$EntitlementMsURL\","
var_replace SitesMsURL "\"$SitesMsURL\","
var_replace SitesNoTokenURL "\"$SitesNoTokenURL\","
var_replace AutomationEngineURL "\"$AutomationEngineURL\","
var_replace AgentConfigMsURL "\"$AgentConfigMsURL\","
var_replace ExecutionResultKafkaTopic "\"$ExecutionResultKafkaTopic\","
var_replace AgentServiceURL "\"$AgentServiceURL\","
var_replace EncryptionKey "\"$EncryptionKey\","

var_replace CassandraConcurrentCallNumber "$CassandraConcurrentCallNumber,"
var_replace CassandraConnNumber "$CassandraConnNumber,"
var_replace CassandraTimeoutSec "$CassandraTimeoutSec,"
var_replace CassandraBatchSize "$CassandraBatchSize,"
var_replace InMemoryCacheSize "$InMemoryCacheSize,"
var_replace LogMaxFileSizeInMB "$LogMaxFileSizeInMB,"
var_replace LogOldFileToKeep "$LogOldFileToKeep,"
var_replace HTTPClientResultsTimeoutSec "$HTTPClientResultsTimeoutSec,"
var_replace ClosestTasksWorkersTimeoutSec "$ClosestTasksWorkersTimeoutSec,"

echo "########### Updated Config File with Environment Variables #############"
