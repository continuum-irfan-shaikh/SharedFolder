package appLoader

import (
	"log"
	"os"
	"strings"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gocql/gocql"
	"github.com/google/uuid"

	"github.com/ContinuumLLC/hystrix-go/hystrix/callback"
	"gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/messaging/publisher"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/distributed/scheduler"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/distributed/zookeeper"
	commonlibRest "gitlab.connectwisedev.com/platform/platform-common-lib/src/web/rest"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/asset"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/kafka"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/memcached"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
)

// AppLoader defines the whole configuration of the Application.
// It incorporates partial configurations of components
type AppLoader struct {
	ConfigService    config.Configuration
	CassandraService *gocql.ClusterConfig
	Scheduler        scheduler.Interface
}

// AppLoaderService is the whole configuration of the Application
var (
	AppLoaderService *AppLoader
)

// LoadApplicationServices loads all partial configurations of components
// and populates the AppLoaderService with the configuration data
func LoadApplicationServices(isTest bool) {
	config.Load()
	if err := logger.Load(config.Config.Log); err != nil {
		log.Println("LoadApplicationServices: error during loading logger: ", err)
	}

	models.LoadTemplatesCache()
	asset.Load()

	if err := translation.Load(); err != nil {
		log.Println("ERROR during loading translations: ", err)
		os.Exit(2)
	}

	if isTest {
		AppLoaderService = &AppLoader{
			ConfigService: config.Config,
		}
	}

	if config.Config.AssetCacheEnabled {
		memcached.Load()
	}

	if !isTest {
		cassandra.Load()
	}

	AppLoaderService = &AppLoader{
		ConfigService: config.Config,
		Scheduler:     zookeeper.Scheduler,
	}

	hystrix.ConfigureCommand(kafka.BrokerCommand, hystrix.CommandConfig{
		Timeout:                config.Config.KafkaProducerSettings.TimeoutInSecond * 1000,
		MaxConcurrentRequests:  config.Config.CircuitBreakerSettings.MaxConcurrentRequests,
		RequestVolumeThreshold: config.Config.CircuitBreakerSettings.RequestVolumeThreshold,
		SleepWindow:            config.Config.CircuitBreakerSettings.SleepWindowInSecond * 1000,
		ErrorPercentThreshold:  config.Config.CircuitBreakerSettings.ErrorPercentThreshold,
	})

	callback.Register(kafka.BrokerCommand, func(_ string, state callback.State) {
		transaction := uuid.New().String()
		if state == callback.Open {
			publisher.ReConnect(transaction, kafka.GetConfig(publisher.RegularKafkaProducer, config.Config))
		} else if state == callback.Close {
			publisher.Connected(transaction, kafka.GetConfig(publisher.RegularKafkaProducer, config.Config))
		}
	})

	registerVersion()
	registerHealth()
}

func registerVersion() {
	commonlibRest.RegistryVersion(&commonlibRest.Version{
		GeneralInfo:          common.GetGeneralInfo(),
		SupportedAPIVersions: config.APIVersions,
		BuildNumber:          config.BuildNumber,
		BuildCommitSHA:       config.BuildCommitSHA,
		Repository:           config.Repository,
	})
}

func registerHealth() {
	commonlibRest.RegistryHealth(&commonlibRest.Health{
		GeneralInfo: common.GetGeneralInfo(),
		ConnMethods: []commonlibRest.Statuser{
			&cassandra.ConnectionStatus{Session: cassandra.Session},
			publisher.Health(publisher.RegularKafkaProducer, kafka.GetConfig(publisher.RegularKafkaProducer, config.Config)),
			&zookeeper.ConnectionStatus{
				Path:  config.Config.ZookeeperBasePath,
				Hosts: strings.Split(config.Config.ZookeeperHosts, ","),
			},
		},
		ListenURL: config.Config.ListenURL,
	})
}
