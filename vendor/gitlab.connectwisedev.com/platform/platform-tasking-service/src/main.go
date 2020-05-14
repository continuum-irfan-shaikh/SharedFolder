package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/agent"
	agentConfig "gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/agent-config"
	dynamicGroups "gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/dynamic-groups"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/sites"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/permission_temp"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/request_info"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/legacy"

	"github.com/urfave/negroni"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/messaging"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/access-control"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/app-loader"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/handlers/api"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/handlers/kafka"
	mc "gitlab.connectwisedev.com/platform/platform-tasking-service/src/handlers/memcache"
	sh "gitlab.connectwisedev.com/platform/platform-tasking-service/src/handlers/scheduler"
	k "gitlab.connectwisedev.com/platform/platform-tasking-service/src/infrastructure/kafka"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/infrastructure/scheduler"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/asset"
	automationEngine "gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/automation-engine"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/partner-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	u "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/user"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/freecache"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/memcached"
	as "gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository/asset"
	crepo "gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository/cassandra"
	krepo "gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository/kafka"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository/site"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository/task-counter-cassandra"
	te "gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository/task-execution"
	triggersRepo "gitlab.connectwisedev.com/platform/platform-tasking-service/src/repository/triggers"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/encryption"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/execution-results"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/execution-results-update"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/task-counters"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/task-definitions"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/tasks"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/templates"
	us "gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/scheduler"
	tasksUseCase "gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/tasks"
	triggerUseCase "gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/trigger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/trigger/trigger-definition"
	userUseCase "gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/user"
)

var (
	wg          = &sync.WaitGroup{}
	ctx, cancel = context.WithCancel(context.Background())
)

func main() {
	idCtx := transactionID.NewContext()
	defer func() { cancel(); wg.Wait() }()

	configFile := flag.String("config", "config.json", "Configuration file in JSON-format")
	flag.Parse()

	if len(*configFile) > 0 {
		config.ConfigFilePath = *configFile
	}

	// initialising global Application Service
	appLoader.LoadApplicationServices(false)
	conf := config.Config

	httpClient := &http.Client{
		Timeout: time.Duration(config.Config.HTTPClientTimeoutSec) * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout:     time.Duration(config.Config.HTTPClientTimeoutSec) * time.Second,
			MaxIdleConns:        2 * config.Config.HTTPClientMaxIdleConnPerHost,
			MaxIdleConnsPerHost: config.Config.HTTPClientMaxIdleConnPerHost,
			DisableKeepAlives:   false,
		},
	}

	accessControl.Load(httpClient)

	var (
		cacheDAO       = freecache.New()
		taskCounterDAO = taskCounterCassandra.New(config.Config.CassandraBatchSize)
		memcachedDAO   = memcached.MemCacheInstance
	)

	// User repository form asset service
	assetUserRepo := as.NewUser(httpClient, config.Config.AssetMsURL)
	// Site repository
	siteRepo := site.NewSite(httpClient, config.Config.SitesMsURL)
	// User usecase
	userUC := userUseCase.NewUser(models.UserSitesPersistenceInstance, assetUserRepo, siteRepo)

	// User middleware
	userMD := u.NewUser(logger.Log)

	userService := user.NewService()
	assetsService := asset.NewAssetsService(memcachedDAO, httpClient)
	permissionMD := permission_temp.NewPermissionMiddleware(logger.Log, userService, httpClient, cacheDAO)

	kafkaServiceForDG := kafka.NewDynamicGroups(config.Config)
	kafkaServiceForTasking := kafka.NewTasking(config.Config)
	dynamicGroupsClient := dynamicGroups.NewDynamicGroupsClient(kafkaServiceForDG, httpClient, models.UserSitesPersistenceInstance)
	sitesClient := sites.NewClient(httpClient)
	aeClient := automationEngine.New(logger.Log, httpClient, config.Config.AutomationEngineMSURL, config.Config.TaskingMsURL)

	atRepo := triggersRepo.NewTriggersRepo(cassandra.Session, memcachedDAO)

	agentConf := agentConfig.NewAgentConfClient(httpClient, config.Config.AgentConfigMsURL, logger.Log)
	encryptionService := encryption.NewService(config.Config.EncryptionKey)
	agentClient := agent.NewClient(httpClient, config.Config.AgentServiceURL, logger.Log)

	externalClients := integration.ExternalClients{
		Asset:            assetsService,
		Sites:            sitesClient,
		AgentConfig:      agentConf,
		DynamicGroups:    dynamicGroupsClient,
		AutomationEngine: aeClient,
		HTTP:             httpClient,
		AgentEncryption:  agentClient,
	}

	executionResultConf := messaging.Config{
		Address: strings.Split(conf.KafkaBrokers, ","),
		GroupID: conf.KafkaConsumerGroup,
		Topics:  []string{conf.ExecutionResultKafkaTopic},
	}
	executionResultKafka := k.NewExecutionResult(executionResultConf)

	profilesRepo := crepo.NewProfilesRepo(cassandra.Session)
	taskCounterService := taskCounters.New(taskCounterDAO)

	tr := crepo.NewTask(cassandra.Session)
	tir := crepo.NewTaskInstance(cassandra.Session)
	targetsRepo := crepo.NewTargets(cassandra.Session)
	schedulerRepo := crepo.NewScheduler(cassandra.Session)
	ser := crepo.NewScriptExecutionResults(cassandra.Session)
	migrationRepo := crepo.NewLegacyMigration(cassandra.Session)
	taskInstanceRepo := crepo.NewTaskInstance(cassandra.Session)
	taskExecHistory := crepo.NewTaskExecutionHistory(cassandra.Session)
	executionExpirationRepo := crepo.NewExecutionExpiration(cassandra.Session)
	tasksExecutionRepo := te.NewTaskExecution(httpClient, conf.TaskTypes, logger.Log)
	executionResultRepo := krepo.NewExecutionResult(conf.ExecutionResultKafkaTopic, executionResultKafka)

	repoDTO := repository.DatabaseRepositories{
		ExecutionExpiration: executionExpirationRepo,
		LegacyMigration:     migrationRepo,
		Scheduler:           schedulerRepo,
		ExecutionResults:    ser,
		Task:                tr,
		TaskInstance:        tir,
		Triggers:            atRepo,
		Targets:             targetsRepo,
		Profiles:            profilesRepo,
		ExecutionHistory:    taskExecHistory,
	}
	modelsDTO := models.DataBaseConnectors{
		TemplateCache:       models.TemplateCacheInstance,
		Task:                models.TaskPersistenceInstance,
		UserSites:           models.UserSitesPersistenceInstance,
		TaskSummary:         models.TaskSummaryPersistenceInstance,
		TaskInstance:        models.TaskInstancePersistenceInstance,
		TaskDefinition:      models.TaskDefinitionPersistenceInstance,
		ExecutionResult:     models.ExecutionResultPersistenceInstance,
		ExecResultView:      models.ExecutionResultViewPersistenceInstance,
		ExecutionExpiration: models.ExecutionExpirationPersistenceInstance,
	}

	triggerHandlerUC := triggerUseCase.New(cacheDAO, memcachedDAO, modelsDTO, repoDTO, logger.Log, externalClients)
	triggerDefinition := triggerDefinition.NewDefinitionService(atRepo)

	executionResultUpdateService := executionResultsUpdate.NewExecutionResultUpdateService(
		models.ExecutionResultPersistenceInstance,
		models.TaskPersistenceInstance,
		models.TaskInstancePersistenceInstance,
		taskExecHistory,
		assetsService,
	)
	tasksService := tasks.New(
		models.TaskDefinitionPersistenceInstance,
		models.TaskPersistenceInstance,
		models.TaskInstancePersistenceInstance,
		models.TemplateCacheInstance,
		models.ExecutionResultPersistenceInstance,
		models.TaskSummaryPersistenceInstance,
		models.ExecutionExpirationPersistenceInstance,
		models.UserSitesPersistenceInstance,
		userService,
		assetsService,
		sitesClient,
		dynamicGroupsClient,
		taskCounterDAO,
		cacheDAO,
		kafkaServiceForTasking,
		httpClient,
		userUC,
		triggerHandlerUC,
		triggerDefinition,
		targetsRepo,
		&executionResultUpdateService,
		encryptionService,
		agentClient,
	)
	taskResultsService := taskExecutionResults.NewTaskResultsService(
		models.ExecutionResultViewPersistenceInstance,
		models.ExecutionResultPersistenceInstance,
		userService,
		httpClient,
	)
	taskDefinitionService := taskDefinitions.NewTaskDefinitionService(
		models.TaskDefinitionPersistenceInstance,
		models.TemplateCacheInstance,
		httpClient,
		userService,
		encryptionService,
	)
	templateService := templates.NewTemplateService(
		models.TemplateCacheInstance,
		userService,
		httpClient,
	)

	tasksUC := tasksUseCase.NewTasks(tr, models.TaskPersistenceInstance, tir, ser, cacheDAO, triggerHandlerUC, logger.Log)
	legacyUC := legacy.NewMigrationUsecase(migrationRepo, models.TaskDefinitionPersistenceInstance, logger.Log)

	closestTasksHandler := api.NewClosestTasks(tasksUC, logger.Log)
	scheduledTasksHandler := api.NewScheduledTasksApi(tasksUC, logger.Log)
	tasksHistoryHandler := api.NewLastRunTasksApi(tasksUC, logger.Log)
	legacyHandler := api.NewLegacyMigration(legacyUC, logger.Log)

	handlers := make(map[string]func(http.ResponseWriter, *http.Request))

	handlers[api.TasksClosestHandler] = closestTasksHandler.ServeHTTP
	handlers[api.TasksScheduledHandler] = scheduledTasksHandler.ServeHTTP
	handlers[api.GetLegacyScriptByPartnerHandler] = legacyHandler.GetByPartner
	handlers[api.GetLegacyScriptByScriptHandler] = legacyHandler.GetByLegacyScript
	handlers[api.InsertScriptInfoHandler] = legacyHandler.InsertLegacyInfo
	handlers[api.TasksHistoryHandler] = tasksHistoryHandler.ServeHTTP

	handlers[api.JobInsertHandler] = legacyHandler.InsertJobInfo
	handlers[api.JobGetByPartnerHandler] = legacyHandler.GetJobInfoByPartner
	handlers[api.JobGetByJobIDHandler] = legacyHandler.GetByLegacyJob

	routerDTO := services.RouterDTO{
		tasksService,
		templateService,
		taskResultsService,
		executionResultUpdateService,
		taskDefinitionService,
		taskCounterService,
		userMD,
		permissionMD,
		handlers}
	router := services.NewRouter(routerDTO)

	oneTimeTaskUCBuilder := us.OneTimeTasksBuilder{
		Conf:                    conf,
		TaskExecutionRepo:       tasksExecutionRepo,
		TaskInstanceRepo:        taskInstanceRepo,
		ExecutionResultRepo:     executionResultRepo,
		TaskRepo:                models.TaskPersistenceInstance,
		ExecutionExpirationRepo: executionExpirationRepo,
		CacheRepo:               models.TemplateCacheInstance,
		Log:                     logger.Log,
		EncryptionService:       encryptionService,
		AgentEncryptionService:  agentClient,
	}
	// --------------------------------------scheduler--------------------------------
	// general trigger UC for scheduler that contains different activation implementations of triggers
	schedulerTriggerUC := us.NewTrigger(logger.Log, triggerHandlerUC)
	oneTimeTaskUC := oneTimeTaskUCBuilder.Build()
	recurrentTaskUC := us.New(
		conf,
		logger.Log,
		triggerHandlerUC,
		taskInstanceRepo,
		targetsRepo,
		dynamicGroupsClient,
		sitesClient,
		tasksExecutionRepo,
		executionResultRepo,
		models.TaskPersistenceInstance,
		executionExpirationRepo,
		models.TemplateCacheInstance,
		assetsService,
		encryptionService,
		agentClient,
		assetsService,
		httpClient,
	)

	schedulerUCs := map[tasking.Regularity]sh.SchedulerTypeUC{
		tasking.OneTime:   oneTimeTaskUC,
		tasking.Trigger:   schedulerTriggerUC,
		tasking.Recurrent: recurrentTaskUC,
	}

	kafkaEndpoint := kafka.NewEndpoints(models.TaskPersistenceInstance, models.TaskInstancePersistenceInstance, taskInstanceRepo, *recurrentTaskUC, logger.Log, config.Config, assetsService)
	kafkaEndpoint.Init()

	schedulerInstance := sh.New(
		taskCounterService,
		executionResultRepo,
		&executionResultUpdateService,
		schedulerRepo,
	)

	schedulerTasksUC := us.NewScheduler(schedulerRepo, models.TaskPersistenceInstance)
	memcacheLoader := mc.NewLoader(atRepo, logger.Log)
	schedulerLoadDTO := scheduler.LoadDTO{
		Ctx:       ctx,
		Log:       logger.Log,
		Conf:      conf,
		WG:        wg,
		DS:        appLoader.AppLoaderService.Scheduler,
		Service:   schedulerInstance,
		Scheduler: sh.NewScheduler(schedulerTasksUC, schedulerUCs, logger.Log),
		Loader:    memcacheLoader,
		Trigger:   triggerHandlerUC,
	}

	if err := scheduler.Load(schedulerLoadDTO); err != nil {
		logger.Log.ErrfCtx(idCtx, errorcode.ErrorApplication, "Error while loading scheduler %v", err)
		os.Exit(1)
	}
	// --------------------------------------scheduler--------------------------------

	// setting up web server middleware
	middlewareManager := negroni.New()
	middlewareManager.Use(negroni.NewRecovery())
	middlewareManager.Use(negroni.HandlerFunc(transactionID.Middleware))
	middlewareManager.Use(negroni.HandlerFunc(partnerID.Middleware))
	middlewareManager.Use(request_info.NewMiddlewareFromLogger("ServiceLogger"))
	middlewareManager.UseHandler(router)

	if err := http.ListenAndServe(appLoader.AppLoaderService.ConfigService.ListenURL, middlewareManager); err != nil {
		logger.Log.ErrfCtx(idCtx, errorcode.ErrorApplication,"Stop running application: %v", err)
		os.Exit(1)
	}
}
