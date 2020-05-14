package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	infrastructureConsumer "gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/messaging/consumer"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/messaging"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/kafka/consumer"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/scheduler"
)

//go:generate mockgen -destination=../../mocks/mocks-gomock/taskInstanceEndpointRepo_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-tasking-service/src/handlers/kafka InstanceEndpointsRepo

type InstanceEndpointsRepo interface {
	RemoveInactiveEndpoints(ti models.TaskInstance, endpoints ...gocql.UUID) (err error)
}

// EndpointForDeletion struct contains info about endpoint where agent was uninstalled
type EndpointForDeletion struct {
	EndpointID gocql.UUID `json:"endpointID"`
	PartnerID  string     `json:"partnerID"`
}

// Endpoints represent impl of endpoints kafka handler
type Endpoints struct {
	taskRepo           models.TaskPersistence
	tiRepo             models.TaskInstancePersistence
	tiEndpointsRepo    InstanceEndpointsRepo
	recurrentScheduler scheduler.RecurrentTaskProcessor
	logger             logger.Logger
	config             config.Configuration
	asset              integration.Asset
}

// NewEndpoints returns new Endpoints impl
func NewEndpoints(t models.TaskPersistence, ti models.TaskInstancePersistence, tiEndpointsRepo InstanceEndpointsRepo, rc scheduler.RecurrentTaskProcessor, logger logger.Logger, cfg config.Configuration, asset integration.Asset) *Endpoints {
	return &Endpoints{
		taskRepo:           t,
		tiRepo:             ti,
		tiEndpointsRepo:    tiEndpointsRepo,
		recurrentScheduler: rc,
		logger:             logger,
		config:             cfg,
		asset:              asset,
	}
}

// Init initializes/executes all components related to processing Kafka messages from managed-endpoint-change topic
func (e *Endpoints) Init() {
	endpointsConsumer := consumer.NewConsumer(e.logger, e.config.ManagedEndpointChangeTopic,
		consumer.GetGoroutineLimits(e.config, e.config.ManagedEndpointChangeTopic),
		e.config.Kafka.ConsumerGroup, e.config.Kafka.Brokers, e.Handle)
	go endpointsConsumer.Consume()
}

// Envelope represents envelope kafka message
type Envelope struct {
	Type    string
	Topic   string
	Message string // should be removed
	Header  messaging.Header
	Context json.RawMessage
}

// Handle handle ManagedEndpointChangeTopic kafka topic
func (e *Endpoints) Handle(msg infrastructureConsumer.Message) (err error) {
	ctx := transactionID.NewContext()
	message := Envelope{}
	if err = json.Unmarshal(msg.Message, &message); err != nil {
		err = fmt.Errorf("Unable to umarshall Kafka message '%v', error: %v\n", msg, err)
		e.logger.ErrfCtx(ctx, errorcode.ErrorKafka, err.Error())
		return
	}

	messageHeader := message.Header.Get(messaging.MessageType)
	if len(messageHeader) == 0 {
		err = fmt.Errorf("Unable to parse Kafka message '%v', error: no headers\n", msg)
		e.logger.ErrfCtx(ctx, errorcode.ErrorKafka, err.Error())
		return
	}

	e.logger.DebugfCtx(ctx,"Processing %v managed_endpoint message type. (%v)", messageHeader[0], messageHeader)

	messageType := strings.ToUpper(messageHeader[0])
	switch messageType {
	case "DELETE":
		return e.processMessage(ctx, message)
	default:
		return nil
	}
}

func (e *Endpoints) processMessage(ctx context.Context, msg Envelope) (err error) {
	e.logger.DebugfCtx(ctx, "Processing DELETE msg %v", msg)

	var messageObj EndpointForDeletion
	if err = json.Unmarshal([]byte(msg.Message), &messageObj); err != nil {
		err = fmt.Errorf("Unable to parse Kafka message '%v', error: %v\n", string(msg.Message), err)
		e.logger.ErrfCtx(ctx, errorcode.ErrorCantDecodeInputData, err.Error())
		return
	}

	if err = e.changeTaskState(ctx, messageObj.PartnerID, messageObj.EndpointID); err != nil {
		msg := "Unable to stop scheduled tasks for managed endpoint with ID %v for partnerID %v, error: %v\n"
		err = fmt.Errorf(msg, messageObj.EndpointID, messageObj.PartnerID, err)
		e.logger.ErrfCtx(ctx, errorcode.ErrorCantProcessData, err.Error())
		return
	}

	return nil
}

func (e *Endpoints) removeEndpoint(endpointToRemove string, endpoints []string) (res []string) {
	for _, e := range endpoints {
		if e == endpointToRemove {
			continue
		}

		res = append(res, e)
	}
	return res
}

func (e *Endpoints) changeTaskState(ctx context.Context, partnerID string, endpointID gocql.UUID) (err error) {
	tasks, err := e.taskRepo.GetByPartnerAndManagedEndpointID(ctx, partnerID, endpointID, common.UnlimitedCount)
	if err != nil {
		return err
	}

	var tasksForUpdate []models.Task
	for _, task := range tasks {
		if task.IsScheduled() {
			task.State = statuses.TaskStateInactive
			if task.TargetsByType == nil {
				task.TargetsByType = make(models.TargetsByType)
			}
			if task.TargetType == models.ManagedEndpoint {
				task.TargetsByType[models.ManagedEndpoint] = e.removeEndpoint(task.ManagedEndpointID.String(), task.TargetsByType[models.ManagedEndpoint])
			}
			tasksForUpdate = append(tasksForUpdate, task)
		}
	}

	if err = e.taskRepo.InsertOrUpdate(ctx, tasksForUpdate...); err != nil {
		e.logger.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "changeTaskState: InsertOrUpdate %v", err)
		return
	}

	if len(tasksForUpdate) > 0 {
		e.processInstances(ctx, tasks)
	}
	return
}

func (e *Endpoints) processInstances(ctx context.Context, tasks []models.Task) {
	for _, task := range tasks {
		if task.IsScheduled() {
			isLastEndpoint := e.processInstance(ctx, task)
			if isLastEndpoint {
				task.State = statuses.TaskStateActive
				defaultUUID, err := gocql.ParseUUID(models.DefaultEndpointUID)
				if err != nil {
					continue
				}

				task.ManagedEndpointID = defaultUUID
				if err = e.taskRepo.InsertOrUpdate(context.Background(), task); err != nil {
					e.logger.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "processInstances: InsertOrUpdate %v", err)
					return
				}
			}
		}
	}
}

func (e *Endpoints) processInstance(ctx context.Context, task models.Task) (isLastEndpoint bool) {
	ti, err := e.getInstanceToUpdate(ctx, task)
	if err != nil {
		e.logger.ErrfCtx(ctx, errorcode.ErrorCantProcessData,"getInstanceToUpdate: %v", err)
		return
	}

	if err := e.tiEndpointsRepo.RemoveInactiveEndpoints(ti, task.ManagedEndpointID); err != nil {
		e.logger.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "processInstance: unable to remove endpoint %v", err)
		return
	}
	delete(ti.Statuses, task.ManagedEndpointID)

	if err = e.tiRepo.Insert(context.Background(), ti); err != nil {
		e.logger.ErrfCtx(ctx, errorcode.ErrorCantInsertData, "processInstance: Insert%v", err)
		return
	}

	if len(ti.Statuses) == 0 {
		isLastEndpoint = true
	}
	return
}

func (e *Endpoints) getInstanceToUpdate(ctx context.Context, task models.Task) (ti models.TaskInstance, err error) {
	instances, err := e.tiRepo.GetByIDs(ctx, task.LastTaskInstanceID)
	if err != nil {
		return ti, fmt.Errorf("getInstanceToUpdate: can't get instance by LastTaskInstanceID, cause: %s", err.Error())
	}

	if len(instances) == 0 {
		return ti, fmt.Errorf("instances not found")
	}

	if task.Schedule.Regularity != tasking.Recurrent {
		return instances[0], nil
	}

	inst := instances[0]
	startedAt := inst.StartedAt
	for {
		gotTI, err := e.tiRepo.GetNearestInstanceAfter(task.ID, startedAt)
		if err != nil {
			return inst, nil
		}

		if gotTI.TriggeredBy == "" {
			inst = gotTI
		}
		startedAt = gotTI.StartedAt
	}
}
