package kafka

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/scheduler"
)

func TestMessageHandler(t *testing.T) {
	uuid, _ := gocql.ParseUUID("123e4567-e89b-12d3-a456-426655440000")
	testCases := []struct {
		name     string
		taskRepo func(ctl *gomock.Controller) models.TaskPersistence
		message  Envelope
	}{
		{
			name: "testCase 1 - no headers",
			message: Envelope{
				Message: ``,
			},
			taskRepo: func(ctl *gomock.Controller) models.TaskPersistence {
				return nil
			},
		},
		{
			name: "testCase 2 - not delete header",
			message: Envelope{
				Message: `BAD`,
			},
			taskRepo: func(ctl *gomock.Controller) models.TaskPersistence {
				return nil
			},

		},
		{
			name: "testCase 3 - can't unmarshal",
			message: Envelope{
				Message: `BAD`,
			},
			taskRepo: func(ctl *gomock.Controller) models.TaskPersistence {
				return nil
			},
		},
		{
			name: "testCase 4 - can't change task state",
			message: Envelope{
				Message: `{"endpointID":"123e4567-e89b-12d3-a456-426655440000"}`,
			},
			taskRepo: func(ctl *gomock.Controller) models.TaskPersistence {
				task := mocks.NewMockTaskPersistence(ctl)
				task.EXPECT().GetByPartnerAndManagedEndpointID(gomock.Any(), gomock.Any(), uuid, gomock.Any()).Return([]models.Task{}, fmt.Errorf("err"))
				return task
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			taskRepo := tc.taskRepo(ctl)
			logger := logger.Log

			endp := NewEndpoints(taskRepo, nil, nil, scheduler.RecurrentTaskProcessor{}, logger, config.Config, nil)
			endp.processMessage(context.TODO(), tc.message)
		})
	}
}

func TestChangeTaskState(t *testing.T) {
	uuid, _ := gocql.ParseUUID("123e4567-e89b-12d3-a456-426655440000")
	partnerIDStr := "1"
	logger.Load(config.Config.Log)
	testCases := []struct {
		name       string
		repos      func(ctl *gomock.Controller) (models.TaskPersistence, models.TaskInstancePersistence, InstanceEndpointsRepo, logger.Logger)
		partnerID  string
		endpointID gocql.UUID
		expectErr  bool
	}{
		{
			name:       "testCase 1 - can't change task state",
			partnerID:  partnerIDStr,
			endpointID: uuid,
			repos: func(ctl *gomock.Controller) (models.TaskPersistence, models.TaskInstancePersistence, InstanceEndpointsRepo, logger.Logger) {
				task := mocks.NewMockTaskPersistence(ctl)

				tasks := []models.Task{{
					ID:            uuid,
					TargetType:    models.ManagedEndpoint,
					TargetsByType: map[models.TargetType][]string{models.ManagedEndpoint: {uuid.String()}},
					State:         statuses.TaskStateActive,
					Schedule: tasking.Schedule{
						Regularity: tasking.OneTime,
					},
				}}
				task.EXPECT().GetByPartnerAndManagedEndpointID(gomock.Any(), gomock.Any(), uuid, gomock.Any()).Return(tasks, nil)
				var expectedTasks = make([]models.Task, 1)
				copy(expectedTasks, tasks)
				expectedTasks[0].State = statuses.TaskStateInactive
				task.EXPECT().InsertOrUpdate(gomock.Any(), expectedTasks).Return(fmt.Errorf("err"))
				tie := mocks.NewMockInstanceEndpointsRepo(ctl)

				return task, nil, tie, nil
			},
			expectErr: true,
		},
		{
			name:       "testCase 2 - no tasks To update",
			partnerID:  partnerIDStr,
			endpointID: uuid,
			repos: func(ctl *gomock.Controller) (models.TaskPersistence, models.TaskInstancePersistence, InstanceEndpointsRepo, logger.Logger) {
				task := mocks.NewMockTaskPersistence(ctl)
				task.EXPECT().GetByPartnerAndManagedEndpointID(gomock.Any(), gomock.Any(), uuid, gomock.Any()).
					Return(nil, nil)
				task.EXPECT().InsertOrUpdate(gomock.Any(), gomock.Any())
				tie := mocks.NewMockInstanceEndpointsRepo(ctl)

				return task, nil, tie, nil
			},
			expectErr: false,
		},
		{
			name:       "testCase 3 - can't get Instance for update",
			partnerID:  partnerIDStr,
			endpointID: uuid,
			repos: func(ctl *gomock.Controller) (models.TaskPersistence, models.TaskInstancePersistence, InstanceEndpointsRepo, logger.Logger) {
				task := mocks.NewMockTaskPersistence(ctl)
				ti := mocks.NewMockTaskInstancePersistence(ctl)
				tasks := []models.Task{{
					ID:                 uuid,
					LastTaskInstanceID: uuid,
					TargetType:         models.ManagedEndpoint,
					TargetsByType:      map[models.TargetType][]string{models.ManagedEndpoint: {uuid.String()}},
					State:              statuses.TaskStateActive,
					Schedule: tasking.Schedule{
						Regularity: tasking.OneTime,
					},
				}}
				task.EXPECT().GetByPartnerAndManagedEndpointID(gomock.Any(), gomock.Any(), uuid, gomock.Any()).Return(tasks, nil)
				var expectedTasks = make([]models.Task, 1)
				copy(expectedTasks, tasks)
				expectedTasks[0].State = statuses.TaskStateInactive
				task.EXPECT().InsertOrUpdate(gomock.Any(), expectedTasks)

				ti.EXPECT().GetByIDs(gomock.Any(), tasks[0].LastTaskInstanceID).
					Return(nil, fmt.Errorf("err"))
				tie := mocks.NewMockInstanceEndpointsRepo(ctl)

				return task, ti, tie, nil
			},
		},
		{
			name:       "testCase 4 - can't get Instance for update, zero",
			partnerID:  partnerIDStr,
			endpointID: uuid,
			repos: func(ctl *gomock.Controller) (models.TaskPersistence, models.TaskInstancePersistence, InstanceEndpointsRepo, logger.Logger) {
				task := mocks.NewMockTaskPersistence(ctl)
				ti := mocks.NewMockTaskInstancePersistence(ctl)
				tasks := []models.Task{{
					ID:                 uuid,
					LastTaskInstanceID: uuid,
					TargetType:         models.ManagedEndpoint,
					TargetsByType:      map[models.TargetType][]string{models.ManagedEndpoint: {uuid.String()}},
					State:              statuses.TaskStateActive,
					Schedule: tasking.Schedule{
						Regularity: tasking.OneTime,
					},
				}}
				task.EXPECT().GetByPartnerAndManagedEndpointID(gomock.Any(), gomock.Any(), uuid, gomock.Any()).Return(tasks, nil)
				var expectedTasks = make([]models.Task, 1)
				copy(expectedTasks, tasks)
				expectedTasks[0].State = statuses.TaskStateInactive
				task.EXPECT().InsertOrUpdate(gomock.Any(), expectedTasks)

				ti.EXPECT().GetByIDs(gomock.Any(), tasks[0].LastTaskInstanceID).
					Return(nil, nil)
				tie := mocks.NewMockInstanceEndpointsRepo(ctl)

				return task, ti, tie, nil
			},
		},
		{
			name:       "testCase 5 - can't get instance After",
			partnerID:  partnerIDStr,
			endpointID: uuid,
			repos: func(ctl *gomock.Controller) (models.TaskPersistence, models.TaskInstancePersistence, InstanceEndpointsRepo, logger.Logger) {
				task := mocks.NewMockTaskPersistence(ctl)
				ti := mocks.NewMockTaskInstancePersistence(ctl)
				tasks := []models.Task{{
					ID:                 uuid,
					LastTaskInstanceID: uuid,
					TargetType:         models.ManagedEndpoint,
					State:              statuses.TaskStateActive,
					TargetsByType:      map[models.TargetType][]string{models.ManagedEndpoint: {uuid.String()}},
					Schedule: tasking.Schedule{
						Regularity:   tasking.Recurrent,
						StartRunTime: time.Now().UTC(),
						Repeat: tasking.Repeat{
							Every:     1,
							RunTime:   time.Now().UTC().Add(time.Minute),
							Frequency: tasking.Hourly,
						},
					},
				}}
				task.EXPECT().GetByPartnerAndManagedEndpointID(gomock.Any(), gomock.Any(), uuid, gomock.Any()).Return(tasks, nil)
				var expectedTasks = make([]models.Task, 1)
				copy(expectedTasks, tasks)
				expectedTasks[0].State = statuses.TaskStateInactive
				task.EXPECT().InsertOrUpdate(gomock.Any(), expectedTasks)

				ti.EXPECT().GetByIDs(gomock.Any(), tasks[0].LastTaskInstanceID).
					Return([]models.TaskInstance{{ID: uuid}}, nil)
				ti.EXPECT().GetNearestInstanceAfter(tasks[0].ID, gomock.Any()).
					Return(models.TaskInstance{}, fmt.Errorf("err"))
				ti.EXPECT().Insert(gomock.Any(), models.TaskInstance{ID: uuid}).
					Return(fmt.Errorf("err"))

				tie := mocks.NewMockInstanceEndpointsRepo(ctl)
				tie.EXPECT().RemoveInactiveEndpoints(gomock.Any(), gomock.Any()).Return(nil)

				return task, ti, tie, nil
			},
		},
		{
			name:       "testCase 6 - Insert of instance returned error",
			partnerID:  partnerIDStr,
			endpointID: uuid,
			repos: func(ctl *gomock.Controller) (models.TaskPersistence, models.TaskInstancePersistence, InstanceEndpointsRepo, logger.Logger) {
				task := mocks.NewMockTaskPersistence(ctl)
				ti := mocks.NewMockTaskInstancePersistence(ctl)
				tasks := []models.Task{{
					ID:                 uuid,
					LastTaskInstanceID: uuid,
					ManagedEndpointID:  uuid,
					TargetType:         models.ManagedEndpoint,
					TargetsByType:      map[models.TargetType][]string{models.ManagedEndpoint: {uuid.String()}},
					State:              statuses.TaskStateActive,
					Schedule: tasking.Schedule{
						Regularity:   tasking.Recurrent,
						StartRunTime: time.Now().UTC(),
						Repeat: tasking.Repeat{
							Every:     1,
							RunTime:   time.Now().UTC().Add(time.Minute),
							Frequency: tasking.Hourly,
						},
					},
				}}
				task.EXPECT().GetByPartnerAndManagedEndpointID(gomock.Any(), gomock.Any(), uuid, gomock.Any()).Return(tasks, nil)
				var expectedTasks = make([]models.Task, 1)
				copy(expectedTasks, tasks)
				expectedTasks[0].State = statuses.TaskStateInactive
				task.EXPECT().InsertOrUpdate(gomock.Any(), expectedTasks)

				ti.EXPECT().GetByIDs(gomock.Any(), tasks[0].LastTaskInstanceID).
					Return([]models.TaskInstance{{ID: uuid}}, nil)
				ti.EXPECT().GetNearestInstanceAfter(tasks[0].ID, gomock.Any()).
					Return(models.TaskInstance{ID: uuid}, errors.New(""))
				ti.EXPECT().Insert(gomock.Any(), models.TaskInstance{ID: uuid}).
					Return(fmt.Errorf("err"))

				tie := mocks.NewMockInstanceEndpointsRepo(ctl)
				tie.EXPECT().RemoveInactiveEndpoints(gomock.Any(), gomock.Any()).Return(nil)
				return task, ti, tie, nil
			},
		},
		{
			name:       "testCase 6.5 - Removing of endpoints returned error",
			partnerID:  partnerIDStr,
			endpointID: uuid,
			repos: func(ctl *gomock.Controller) (models.TaskPersistence, models.TaskInstancePersistence, InstanceEndpointsRepo, logger.Logger) {
				task := mocks.NewMockTaskPersistence(ctl)
				ti := mocks.NewMockTaskInstancePersistence(ctl)
				tasks := []models.Task{{
					ID:                 uuid,
					LastTaskInstanceID: uuid,
					ManagedEndpointID:  uuid,
					TargetType:         models.ManagedEndpoint,
					TargetsByType:      map[models.TargetType][]string{models.ManagedEndpoint: {uuid.String()}},
					State:              statuses.TaskStateActive,
					Schedule: tasking.Schedule{
						Regularity:   tasking.Recurrent,
						StartRunTime: time.Now().UTC(),
						Repeat: tasking.Repeat{
							Every:     1,
							RunTime:   time.Now().UTC().Add(time.Minute),
							Frequency: tasking.Hourly,
						},
					},
				}}
				task.EXPECT().GetByPartnerAndManagedEndpointID(gomock.Any(), gomock.Any(), uuid, gomock.Any()).Return(tasks, nil)
				var expectedTasks = make([]models.Task, 1)
				copy(expectedTasks, tasks)
				expectedTasks[0].State = statuses.TaskStateInactive
				task.EXPECT().InsertOrUpdate(gomock.Any(), expectedTasks)

				ti.EXPECT().GetByIDs(gomock.Any(), tasks[0].LastTaskInstanceID).
					Return([]models.TaskInstance{{ID: uuid}}, nil)
				ti.EXPECT().GetNearestInstanceAfter(tasks[0].ID, gomock.Any()).
					Return(models.TaskInstance{ID: uuid}, errors.New(""))

				tie := mocks.NewMockInstanceEndpointsRepo(ctl)
				tie.EXPECT().RemoveInactiveEndpoints(gomock.Any(), gomock.Any()).Return(errors.New("err"))


				return task, ti, tie, nil
			},
		},
		{
			name:       "testCase 7 - ok",
			partnerID:  partnerIDStr,
			endpointID: uuid,
			repos: func(ctl *gomock.Controller) (models.TaskPersistence, models.TaskInstancePersistence, InstanceEndpointsRepo, logger.Logger) {
				task := mocks.NewMockTaskPersistence(ctl)
				ti := mocks.NewMockTaskInstancePersistence(ctl)
				tasks := []models.Task{{
					ID:                 uuid,
					LastTaskInstanceID: uuid,
					ManagedEndpointID:  uuid,
					TargetType:         models.ManagedEndpoint,
					TargetsByType:      map[models.TargetType][]string{models.ManagedEndpoint: {uuid.String()}},
					State:              statuses.TaskStateActive,
					Schedule: tasking.Schedule{
						Regularity:   tasking.Recurrent,
						StartRunTime: time.Now().UTC(),
						Repeat: tasking.Repeat{
							Every:     1,
							RunTime:   time.Now().UTC().Add(time.Minute),
							Frequency: tasking.Hourly,
						},
					},
				}}
				task.EXPECT().GetByPartnerAndManagedEndpointID(gomock.Any(), gomock.Any(), uuid, gomock.Any()).Return(tasks, nil)
				var expectedTasks = make([]models.Task, 1)
				copy(expectedTasks, tasks)
				expectedTasks[0].State = statuses.TaskStateInactive
				task.EXPECT().InsertOrUpdate(gomock.Any(), expectedTasks)

				uuidDefault, _ := gocql.ParseUUID(models.DefaultEndpointUID)
				expectedTasksAfterActivationWithDefaultUUID := expectedTasks[0]
				expectedTasksAfterActivationWithDefaultUUID.State = statuses.TaskStateActive
				expectedTasksAfterActivationWithDefaultUUID.ManagedEndpointID = uuidDefault
				task.EXPECT().InsertOrUpdate(gomock.Any(), expectedTasksAfterActivationWithDefaultUUID)

				ti.EXPECT().GetByIDs(gomock.Any(), tasks[0].LastTaskInstanceID).
					Return([]models.TaskInstance{{ID: uuid}}, nil)
				ti.EXPECT().GetNearestInstanceAfter(tasks[0].ID, gomock.Any()).
					Return(models.TaskInstance{ID: uuid}, errors.New(""))
				ti.EXPECT().
					Insert(gomock.Any(), models.TaskInstance{ID: uuid}).
					Return(nil)
				tie := mocks.NewMockInstanceEndpointsRepo(ctl)
				tie.EXPECT().RemoveInactiveEndpoints(gomock.Any(), gomock.Any()).Return(nil)

				return task, ti, tie, nil
			},
		},
		{
			name:       "testCase 8 - ok for one time",
			partnerID:  partnerIDStr,
			endpointID: uuid,
			repos: func(ctl *gomock.Controller) (models.TaskPersistence, models.TaskInstancePersistence, InstanceEndpointsRepo, logger.Logger) {
				task := mocks.NewMockTaskPersistence(ctl)
				ti := mocks.NewMockTaskInstancePersistence(ctl)
				tasks := []models.Task{{
					ID:                 uuid,
					LastTaskInstanceID: uuid,
					TargetsByType:      map[models.TargetType][]string{models.ManagedEndpoint: {uuid.String()}},
					ManagedEndpointID:  uuid,
					TargetType:         models.ManagedEndpoint,
					State:              statuses.TaskStateActive,
					Schedule: tasking.Schedule{
						Regularity:   tasking.OneTime,
						StartRunTime: time.Now().UTC(),
					},
				}}
				task.EXPECT().GetByPartnerAndManagedEndpointID(gomock.Any(), gomock.Any(), uuid, gomock.Any()).Return(tasks, nil)
				var expectedTasks = make([]models.Task, 1)
				copy(expectedTasks, tasks)
				expectedTasks[0].State = statuses.TaskStateInactive
				task.EXPECT().InsertOrUpdate(gomock.Any(), expectedTasks)

				uuidDefault, _ := gocql.ParseUUID(models.DefaultEndpointUID)
				expectedTasksAfterActivationWithDefaultUUID := expectedTasks[0]
				expectedTasksAfterActivationWithDefaultUUID.State = statuses.TaskStateActive
				expectedTasksAfterActivationWithDefaultUUID.ManagedEndpointID = uuidDefault
				task.EXPECT().InsertOrUpdate(gomock.Any(), expectedTasksAfterActivationWithDefaultUUID)

				ti.EXPECT().GetByIDs(gomock.Any(), tasks[0].LastTaskInstanceID).
					Return([]models.TaskInstance{{ID: uuid}}, nil)
				ti.EXPECT().Insert(gomock.Any(), models.TaskInstance{ID: uuid})
				tie := mocks.NewMockInstanceEndpointsRepo(ctl)
				tie.EXPECT().RemoveInactiveEndpoints(gomock.Any(), gomock.Any()).Return(nil)

				return task, ti, tie, nil
			},
		},
		{
			name:       "testCase 9 - ok for one time, but can't make another insert",
			partnerID:  partnerIDStr,
			endpointID: uuid,
			repos: func(ctl *gomock.Controller) (models.TaskPersistence, models.TaskInstancePersistence, InstanceEndpointsRepo, logger.Logger) {
				task := mocks.NewMockTaskPersistence(ctl)
				ti := mocks.NewMockTaskInstancePersistence(ctl)
				tasks := []models.Task{{
					ID:                 uuid,
					LastTaskInstanceID: uuid,
					ManagedEndpointID:  uuid,
					TargetType:         models.ManagedEndpoint,
					TargetsByType:      map[models.TargetType][]string{models.ManagedEndpoint: {uuid.String()}},
					State:              statuses.TaskStateActive,
					Schedule: tasking.Schedule{
						Regularity:   tasking.OneTime,
						StartRunTime: time.Now().UTC(),
					},
				}}
				task.EXPECT().GetByPartnerAndManagedEndpointID(gomock.Any(), gomock.Any(), uuid, gomock.Any()).Return(tasks, nil)
				var expectedTasks = make([]models.Task, 1)
				copy(expectedTasks, tasks)
				expectedTasks[0].State = statuses.TaskStateInactive
				task.EXPECT().InsertOrUpdate(gomock.Any(), expectedTasks)

				uuidDefault, _ := gocql.ParseUUID(models.DefaultEndpointUID)
				expectedTasksAfterActivationWithDefaultUUID := expectedTasks[0]
				expectedTasksAfterActivationWithDefaultUUID.State = statuses.TaskStateActive
				expectedTasksAfterActivationWithDefaultUUID.ManagedEndpointID = uuidDefault
				task.EXPECT().InsertOrUpdate(gomock.Any(), expectedTasksAfterActivationWithDefaultUUID).Return(fmt.Errorf("err"))

				ti.EXPECT().GetByIDs(gomock.Any(), tasks[0].LastTaskInstanceID).
					Return([]models.TaskInstance{{ID: uuid}}, nil)
				ti.EXPECT().Insert(gomock.Any(), models.TaskInstance{ID: uuid})
				tie := mocks.NewMockInstanceEndpointsRepo(ctl)
				tie.EXPECT().RemoveInactiveEndpoints(gomock.Any(), gomock.Any()).Return(nil)

				return task, ti, tie, nil
			},
		},
	}

	RegisterTestingT(t)
	for _, tc := range testCases {
		fmt.Println(tc.name)
		ctl := gomock.NewController(t)

		taskRepo, tiRepo, tieRepo, _ := tc.repos(ctl)
		assetMock := mocks.NewMockAsset(ctl)
		assetMock.EXPECT().GetLocationByEndpointID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		endp := NewEndpoints(taskRepo, tiRepo, tieRepo, scheduler.RecurrentTaskProcessor{}, logger.Log, config.Config, assetMock)
		endp.Init() // for coverage

		gotErr := endp.changeTaskState(context.Background(), tc.partnerID, tc.endpointID)
		if tc.expectErr {
			Ω(gotErr).ShouldNot(BeNil(), tc.name)
		} else {
			Ω(gotErr).Should(BeNil(), tc.name)
		}
		ctl.Finish()
	}
}
