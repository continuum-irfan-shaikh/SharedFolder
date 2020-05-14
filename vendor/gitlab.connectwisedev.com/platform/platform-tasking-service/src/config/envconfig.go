package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/messaging/publisher"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
	"github.com/kelseyhightower/envconfig"
)

// ScriptTaskType describe a script type of Task
const ScriptTaskType = "script"

var (
	// Config is a package variable, which is populated during init() execution and shared to whole application
	Config Configuration
	// ConfigFilePath defines a path to JSON-config file
	ConfigFilePath = "config.json"

	// APIVersions stores slice of supported versions
	APIVersions = []string{"v1"}
	//Connections stores slice of all connections
	Connections = []string{"Cassandra", "Kafka"}

	// Repository is link into github project
	Repository = "https://gitlab.connectwisedev.com/platform/platform-tasking-service"

	// BuildCommitSHA is current commit
	BuildCommitSHA string

	// BuildNumber is current jenkins build
	BuildNumber string
)

// Version options
type Version struct {
	SolutionName    string `json:"SolutionName"    envconfig:"TKS_SOLUTION_NAME"`
	ServiceName     string `json:"ServiceName"     envconfig:"TKS_SERVICE_NAME"`
	ServiceProvider string `json:"ServiceProvider" envconfig:"TKS_SERVICE_PROVIDER"`
}

// CacheSettings specifies options for task definition templates cache
type CacheSettings struct {
	Size              int `json:"Size"`              // cache size in bytes
	ReloadIntervalSec int `json:"ReloadIntervalSec"` // reload interval in seconds
	GCPercent         int `json:"GCPercent"`         // GC target percentage (see documentation for `debug.SetGCPercent()`)
	DataTTLSec        int `json:"DataTTLSec"`        // Data TTL in seconds
}

// KafkaConfig options
type KafkaConfig struct {
	Brokers         []string       `json:"BrokerHosts"                  envconfig:"TKS_KAFKA_BROKERS"`
	ConsumerGroup   string         `json:"ConsumerGroup"                envconfig:"TKS_KAFKA_CONSUMERGROUP"`
	GoroutineLimits map[string]int `json:"GoroutineLimitsPerOneCore"    envconfig:"TKS_KAFKA_GOROUTINE_LIMITS"`
}

// KafkaProducerSettings contains config settings for a kafka producer
type KafkaProducerSettings struct {
	CompressionType           publisher.CompressionType `json:"CompressionType"`
	TimeoutInSecond           int                       `json:"TimeoutInSecond"`
	MaxRetry                  int                       `json:"MaxRetry"`
	MaxMessageBytes           int                       `json:"MaxMessageBytes"`
	ReconnectIntervalInSecond int                       `json:"ReconnectIntervalInSecond"`
}

// CircuitBreakerSettings contains config settings for a circuit breaker
type CircuitBreakerSettings struct {
	MaxConcurrentRequests         int  `json:"MaxConcurrentRequests"`
	ErrorPercentThreshold         int  `json:"ErrorPercentThreshold"`
	RequestVolumeThreshold        int  `json:"RequestVolumeThreshold"`
	SleepWindowInSecond           int  `json:"SleepWindowInSecond"`
	RetryAfterLowerBoundInMinutes int  `json:"RetryAfterLowerBoundInMinutes"`
	RetryAfterUpperBoundInMinutes int  `json:"RetryAfterUpperBoundInMinutes"`
	CommandTimeoutInMilliseconds  int  `json:"CommandTimeoutInMilliseconds"`
	Enabled                       bool `json:"EnableCircuitBreaker"`
}

// ScheduledJob is a struct defining the actual scheduled job
// implementing `ScheduledJob` interface
type ScheduledJob struct {
	Name     string `json:"Name"`
	Task     string `json:"Task"`
	Schedule string `json:"Schedule"`
}

// GetName is necessary for implementing `ScheduledJob` interface
func (s ScheduledJob) GetName() string {
	return s.Name
}

// GetTask is necessary for implementing `ScheduledJob` interface
func (s ScheduledJob) GetTask() string {
	return s.Task
}

// GetSchedule is necessary for implementing `ScheduledJob` interface
func (s ScheduledJob) GetSchedule() string {
	return s.Schedule
}

// RetryStrategy is a struct defining retry backoff for HTTP calls to another microservices,
type RetryStrategy struct {
	MaxNumberOfRetries    int `json:"MaxNumberOfRetries"`
	RetrySleepIntervalSec int `json:"RetrySleepIntervalSec"`
}

// MemcachedConfig is a struct defining configuration of memcache client
type MemcachedConfig struct {
	MemcachedURL      string `json:"MemcachedURL"`
	DefaultDataTTLSec int    `json:"DefaultDataTTLSec"`
	MaxIdleConns      int    `json:"MaxIdleConns"`
	TimeoutSec        int    `json:"TimeoutSec"`
}

// Configuration options
type Configuration struct {
	Log                    logger.Config          `json:"Log"`
	Kafka                  KafkaConfig            `json:"Kafka"`
	KafkaProducerSettings  KafkaProducerSettings  `json:"KafkaProducerSettings"`
	CircuitBreakerSettings CircuitBreakerSettings `json:"CircuitBreakerSettings"`

	ListenURL                     string            `json:"ListenURL"                     envconfig:"TKS_LISTEN_URL"                         default:":12121"`
	CassandraURL                  string            `json:"CassandraURL"                  envconfig:"TKS_CASSANDRA_URL"                      default:"localhost:9042"`
	CassandraTimeoutSec           int               `json:"CassandraTimeoutSec"           envconfig:"TKS_CASSANDRA_TIMEOUT_SEC"              default:"1"`
	CassandraKeyspace             string            `json:"CassandraKeyspace"             envconfig:"TKS_CASSANDRA_Keyspace"                 default:"platform_tasking_db"`
	CassandraConnNumber           int               `json:"CassandraConnNumber"           envconfig:"TKS_CASSANDRA_CONN_NUMBER"              default:"20"`
	CassandraConcurrentCallNumber int               `json:"CassandraConcurrentCallNumber" envconfig:"TKS_CASSANDRA_CONCURRENT_CALL_NUMBER"   default:"30"`
	CassandraBatchSize            int               `json:"CassandraBatchSize"            envconfig:"TKS_CASSANDRA_BATCH_SIZE"               default:"5"`
	KafkaBrokers                  string            `json:"KafkaBrokers"                  envconfig:"SCS_KAFKA_BROKERS"                      default:"localhost:9092"`
	ManagedEndpointChangeTopic    string            `json:"ManagedEndpointChangeTopic"    envconfig:"SCS_KAFKA_MANAGED_ENDPOINT_CHANGE_TOPIC" default:"managed-endpoint-change"`
	TaskingEventsTopic            string            `json:"TaskingEventsTopic"            envconfig:"SCS_KAFKA_TASKING_EVENTS_TOPIC"         default:"tasking-events"`
	KafkaConsumerGroup            string            `json:"KafkaConsumerGroup"            envconfig:"SCS_KAFKA_CONSUMER_GROUP"               default:"platform-tasking-service"`
	KafkaRetryIntervalSec         int               `json:"KafkaRetryIntervalSec"         envconfig:"SCS_KAFKA_RETRY_INTERVAL_SEC"           default:"10"`
	ScriptingMsURL                string            `json:"ScriptingMsURL"                envconfig:"TKS_SCRIPTING_MS_URL"                   default:"http://localhost:8888/scripting/v1"`
	TDTCacheSettings              CacheSettings     `json:"TDTCacheSettings"              envconfig:"TKS_TDT_CACHE_SETTINGS"`
	TaskingMsURL                  string            `json:"TaskingMsURL"                  envconfig:"TKS_TASKING_MS_URL"                     default:"http://localhost:12121/tasking/v1"`
	AssetMsURL                    string            `json:"AssetMsURL"                    envconfig:"TKS_ASSET_MS_URL"                       default:"http://127.0.0.1:8084/asset/v1"`
	EntitlementMsURL              string            `json:"EntitlementMsURL"              envconfig:"TKS_ENTITLEMENT_MS_URL"                 default:"http://internal-entitlement-partnerspecific-72155629.ap-south-1.elb.amazonaws.com/entitlement/v1"`
	EntitlementCacheSettings      CacheSettings     `json:"EntitlementCacheSettings"      envconfig:"TKS_ENTITLEMENT_CACHE_SETTINGS"`
	Memcached                     MemcachedConfig   `json:"Memcached"                     envconfig:"TKS_MEMCACHED"`
	DynamicGroupsMsURL            string            `json:"DynamicGroupsMsURL"            envconfig:"TKS_DYNAMIC_GROUPS_MS_URL"              default:"http://internal-continuum-dg-service-elb-int-1291521176.ap-south-1.elb.amazonaws.com/dg/v1"`
	DynamicGroupMonitoringTopic   string            `json:"DynamicGroupMonitoring"        envconfig:"TKS_KAFKA_DYNAMIC_GROUP_MONITORING_TOPIC"      default:"dynamic_group_monitoring"`
	GraphQLMsURL                  string            `json:"GraphQLMsURL"                  envconfig:"TKS_GRAPHQL_MS_URL"                     default:"http://127.0.0.1:8080/GraphQL"`
	SitesMsURL                    string            `json:"SitesMsURL"                    envconfig:"TKS_SITES_MS_URL"                       default:"https://rmmitswebapi.dtitsupport247.net/v1"`
	SitesNoTokenURL               string            `json:"SitesNoTokenURL"               envconfig:"TKS_SITES_NO_TOKEN_URL"                 default:"http://rmmitsapi.dtitsupport247.net/rmmitsapi/v1"`
	AgentConfigMsURL              string            `json:"AgentConfigMsURL"              envconfig:"TKS_AGENT_CONFIG_MS_URL"                default:"http://internal-intplatformagentconfigurationser-2019980644.ap-south-1.elb.amazonaws.com/agent-configuration/v1"`
	AutomationEngineMSURL         string            `json:"AutomationEngineURL"           envconfig:"TKS_AUTOMATION_ENGINE_MS_URL"           default:"http://internal-continuum-ae-service-elb-int-561715420.ap-south-1.elb.amazonaws.com"`
	AssetCacheEnabled             bool              `json:"AssetCacheEnabled"             envconfig:"TKS_ASSET_CACHE_ENABLED"                default:"false"`
	Version                       Version           `json:"Version"                       envconfig:"TKS_VERSION"`
	ZookeeperHosts                string            `json:"ZookeeperHosts"                envconfig:"TKS_ZOOKEEPER_HOSTS"                    default:"localhost:2181"`
	ZookeeperBasePath             string            `json:"ZookeeperBasePath"             envconfig:"TKS_ZOOKEEPER_BASE_PATH"                default:"/tasking-service"`
	ScheduledJobs                 []ScheduledJob    `json:"ScheduledJobs"                 envconfig:"TKS_SCHEDULED_JOBS"`
	JobSchedulingInterval         int               `json:"JobSchedulingInterval"         envconfig:"TKS_JOB_SCHEDULING_INTERVAL"            default:"5"`
	JobListeningInterval          int               `json:"JobListeningInterval"          envconfig:"TKS_JOB_LISTENING_INTERVAL"             default:"3"`
	DefaultLanguage               string            `json:"DefaultLanguage"               envconfig:"TKS_DEFAULT_LANGUAGE"                   default:"en-US"`
	RetryStrategy                 RetryStrategy     `json:"RetryStrategy"                 envconfig:"TKS_RETRY_STRATEGY"`
	TaskTypes                     map[string]string `json:"TaskTypes"                     envconfig:"TKS_TASK_TYPES"`
	HTTPClientTimeoutSec          int               `json:"HTTPClientTimeoutSec"          envconfig:"TKS_HTTP_CLIENT_TIMEOUT_SEC"            default:"60"`
	HTTPClientResultsTimeoutSec   int               `json:"HTTPClientResultsTimeoutSec"   envconfig:"TKS_HTTP_CLIENT_RESULTS_TIMEOUT_SEC"    default:"120"`
	HTTPClientMaxIdleConnPerHost  int               `json:"HTTPClientMaxIdleConnPerHost"  envconfig:"TKS_HTTP_CLIENT_MAX_IDLE_CONN_PER_HOST" default:"100"`
	ConcurrentRESTCalls           int               `json:"ConcurrentRESTCalls"           envconfig:"TKS_TASK_CONCURRENT_REST_CALLS"         default:"50"`
	DataRetentionIntervalDay      int               `json:"DataRetentionIntervalDay"      envconfig:"TKS_DATA_RETENTION_INTERVAL_DAY"        default:"90"`
	DataRetentionRemoveBatchSize  int               `json:"DataRetentionRemoveBatchSize"  envconfig:"TKS_DATA_RETENTION_REMOVE_BATCH_SIZE"   default:"100"`
	FeaturesForRoutes             map[string]string `json:"FeaturesForRoutes"             envconfig:"TKS_FEATURES_FOR_ROUTES"`
	InMemoryCacheSize             int               `json:"InMemoryCacheSize"             envconfig:"TKS_IN_MEMORY_CACHE_SIZE"               default:"1073741274"`
	WSIncomeQueueSize             int               `json:"WSIncomeQueueSize"             envconfig:"TKS_WS_INCOME_QUEUE_SIZE"               default:"10"`
	WSOutcomeQueueSize            int               `json:"WSOutcomeQueueSize"            envconfig:"TKS_WS_OUTCOME_QUEUE_SIZE"              default:"100"`
	RecalculateTime               int               `json:"RecalculateTime"               envconfig:"RECALCULATE_TIME"                       default:"60"`
	ExecutionResultKafkaTopic     string            `json:"ExecutionResultKafkaTopic"     envconfig:"EXECUTION_RESULT_KAFKA_TOPIC"           default:"script_execution_result"`
	TriggerReloadInterval         int               `json:"TriggersReloadCacheIntervalSec"   envconfig:"TKS_TRIGGER_RELOAD"                  default:"60"`
	PartnerSitesCacheExpiration   int               `json:"PartnerSitesCacheExpiration"   envconfig:"TKS_PARTNER_SITES_CACHE_EXPIRATION"     default:"7200"`
	ClosestTasksWorkersTimeoutSec int               `json:"ClosestTasksWorkersTimeoutSec"   envconfig:"TKS_CLOSEST_TASKS_WORKERS_TIMEOUT_SEC"     default:"20"`
	AgentServiceURL               string            `json:"AgentServiceURL"               envconfig:"TKS_AGENT_SERVICE_URL"                       default:"https://integration.agent.exec.itsupport247.net/agent/v1"`
	EncryptionKey                 string            `json:"EncryptionKey"                 envconfig:"INT_ENCRYPTION_KEY"              validate:"required"`
}

// Load reads and loads configuration to Config variable
func Load() {
	var err error

	confLen := len(ConfigFilePath)
	if confLen != 0 {
		err = readConfigFromJSON(ConfigFilePath)
	}
	if confLen == 0 || err != nil {
		err = readConfigFromENV()
	}
	if err != nil {
		panic(fmt.Sprintf(`Configuration not found. Please specify configuration. err: %s`, err.Error()))
	}
}

// isMissing validates Configuration
func (c *Configuration) isMissing() bool {
	return c.ListenURL == "" || c.CassandraURL == ""
}

// nolint: gosec
func readConfigFromJSON(configFilePath string) error {
	log.Printf("Looking for JSON config file (%s)", configFilePath)

	cfgFile, err := os.Open(configFilePath)
	if err != nil {
		log.Printf("Reading configuration from JSON (%s) failed: %v\n", configFilePath, err)
		return err
	}
	defer func() {
		err := cfgFile.Close()
		if err != nil {
			log.Printf("Cannot close the configuration file [%s]: %v\n", cfgFile.Name(), err)
		}
	}()

	err = json.NewDecoder(cfgFile).Decode(&Config)
	if err != nil {
		log.Printf("Reading configuration from JSON (%s) failed: %s\n", configFilePath, err)
		return err
	}

	log.Printf("Configuration has been read from JSON (%s) successfully\n", configFilePath)
	return nil
}

// readConfigFromENV reads data from environment variables
func readConfigFromENV() (err error) {
	log.Println("Looking for ENV configuration")

	err = envconfig.Process("TKS", &Config)

	if err == nil && Config.isMissing() {
		err = errors.New("configuration is missing")
	} else {
		Config.TDTCacheSettings = CacheSettings{Size: 104857600, ReloadIntervalSec: 600, GCPercent: 20}
		Config.EntitlementCacheSettings = CacheSettings{Size: 104857600, DataTTLSec: 600, GCPercent: 20}
		Config.ScheduledJobs = []ScheduledJob{
			{
				Name:     "Check for tasks",
				Task:     "checkForScheduledTasks",
				Schedule: "@every 1m",
			},
			{
				Name:     "Check for retained data",
				Task:     "checkForRetainedData",
				Schedule: "@every 24h",
			},
			{
				Name:     "Check for expired executions",
				Task:     "checkForExpiredExecutions",
				Schedule: "@every 1m",
			},
			{
				Name:     "Recalculate counts of Tasks",
				Task:     "recalculateTasks",
				Schedule: "@every 24h",
			},
			{
				Name:     "Recalculate devices in trigger tasks",
				Task:     "activeTriggersReopening",
				Schedule: "@every 2m",
			},
		}
		Config.RetryStrategy = RetryStrategy{MaxNumberOfRetries: 10, RetrySleepIntervalSec: 10}
		Config.TaskTypes = map[string]string{ScriptTaskType: Config.ScriptingMsURL}
		Config.Memcached = MemcachedConfig{MemcachedURL: "localhost:11211", DefaultDataTTLSec: 86400, MaxIdleConns: 10, TimeoutSec: 30}
		log.Println("ENV configuration has been read successfully")
	}

	return err
}
