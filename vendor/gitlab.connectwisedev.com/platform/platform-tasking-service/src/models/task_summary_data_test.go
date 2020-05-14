package models_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestTaskSummaryRepoCassandra_GetTasksSummaryData(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrlr *gomock.Controller

	type payload struct {
		isNoc                   bool
		taskIDs                 []gocql.UUID
		cache                   func() persistency.Cache
		taskPersistence         func()
		taskInstancePersistence func()
		taskSummaryPersistence  func()
	}

	type expected struct {
		data []models.TaskSummaryData
		err  error
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "error with get by partner",
			payload: payload{
				taskIDs: make([]gocql.UUID, 0),
				cache: func() persistency.Cache {
					return nil
				},
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByPartner(context.Background(), "partnerID").
						Return(nil, errors.New("some_err")).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				taskInstancePersistence: func() {},
				taskSummaryPersistence:  func() {},
			},
			expected: expected{
				err: errors.New("some_err"),
			},
		},
		{
			name: "error with get by id's",
			payload: payload{
				taskIDs: []gocql.UUID{taskId},
				cache: func() persistency.Cache {
					return nil
				},
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByIDs(context.Background(), nil, "partnerID", false, []gocql.UUID{taskId}).
						Return(nil, errors.New("some_err")).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				taskInstancePersistence: func() {},
				taskSummaryPersistence:  func() {},
			},
			expected: expected{
				err: errors.New("some_err"),
			},
		},
		{
			name: "error with get task instance by id's",
			payload: payload{
				taskIDs: []gocql.UUID{taskId},
				cache: func() persistency.Cache {
					return nil
				},
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByIDs(context.Background(), nil, "partnerID", false, []gocql.UUID{taskId}).
						Return([]models.Task{
							{
								IsRequireNOCAccess: true,
							},
							{
								IsRequireNOCAccess: false,
							},
							{
								IsRequireNOCAccess: false,
								LastTaskInstanceID: lastTaskInstID1,
							},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByIDs(context.Background(), lastTaskInstID1).
						Return(nil, errors.New("some_err")).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				taskSummaryPersistence: func() {},
			},
			expected: expected{
				err: errors.New("some_err"),
			},
		},
		{
			name: "error with get status counts by id",
			payload: payload{
				taskIDs: []gocql.UUID{taskId},
				cache: func() persistency.Cache {
					return nil
				},
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByIDs(context.Background(), nil, "partnerID", false, []gocql.UUID{taskId}).
						Return([]models.Task{
							{
								IsRequireNOCAccess: true,
							},
							{
								IsRequireNOCAccess: false,
							},
							{
								IsRequireNOCAccess: false,
								LastTaskInstanceID: lastTaskInstID1,
							},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByIDs(context.Background(), lastTaskInstID1).
						Return([]models.TaskInstance{
							{
								ID: taskInstID1,
							},
						}, nil).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				taskSummaryPersistence: func() {
					taskSummaryPersistenceMock := mocks.NewMockTaskSummaryPersistence(mockCtrlr)
					taskSummaryPersistenceMock.EXPECT().GetStatusCountsByIDs(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil, errors.New("some_err")).Times(1)
					models.TaskSummaryPersistenceInstance = taskSummaryPersistenceMock
				},
			},
			expected: expected{
				err: errors.New("some_err"),
			},
		},
		{
			name: "Success_#1",
			payload: payload{
				taskIDs: []gocql.UUID{taskId},
				cache: func() persistency.Cache {
					return nil
				},
				taskPersistence: func() {
					taskPersistenceMock := mocks.NewMockTaskPersistence(mockCtrlr)
					taskPersistenceMock.EXPECT().GetByIDs(context.Background(), nil, "partnerID", false, []gocql.UUID{taskId}).
						Return([]models.Task{
							{
								IsRequireNOCAccess: true,
							},
							{
								IsRequireNOCAccess: false,
							},
							{
								IsRequireNOCAccess: false,
								LastTaskInstanceID: lastTaskInstID1,
							},
						}, nil).Times(1)
					models.TaskPersistenceInstance = taskPersistenceMock
				},
				taskInstancePersistence: func() {
					taskInstancePersistenceMock := mocks.NewMockTaskInstancePersistence(mockCtrlr)
					taskInstancePersistenceMock.EXPECT().GetByIDs(context.Background(), lastTaskInstID1).
						Return([]models.TaskInstance{
							{
								ID: taskInstID1,
							},
						}, nil).Times(1)
					models.TaskInstancePersistenceInstance = taskInstancePersistenceMock
				},
				taskSummaryPersistence: func() {
					taskSummaryPersistenceMock := mocks.NewMockTaskSummaryPersistence(mockCtrlr)
					taskSummaryPersistenceMock.EXPECT().GetStatusCountsByIDs(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(map[gocql.UUID]models.TaskInstanceStatusCount{taskInstID1: {TaskInstanceID: taskInstID1, SuccessCount: 1}}, nil).Times(1)
					models.TaskSummaryPersistenceInstance = taskSummaryPersistenceMock
				},
			},
			expected: expected{
				data: []models.TaskSummaryData{
					{
						RunOn: models.TargetData{
							Count: 2,
						},
						Status: 3,
					},
				},
			},
		},
	}

	for _, test := range tc {
		mockCtrlr = gomock.NewController(t)
		test.payload.cache()
		test.payload.taskPersistence()
		test.payload.taskInstancePersistence()
		test.payload.taskSummaryPersistence()

		s := models.TaskSummaryRepoCassandra{}
		d, err := s.GetTasksSummaryData(context.Background(), test.payload.isNoc, test.payload.cache(), "partnerID", test.payload.taskIDs...)
		mockCtrlr.Finish()

		if test.expected.err == nil {
			Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
			Ω(d).To(ConsistOf(test.expected.data), fmt.Sprintf(defaultMsg, test.name))
			continue
		}

		Ω(err.Error()).To(Equal(test.expected.err.Error()), fmt.Sprintf(defaultMsg, test.name))
		Ω(d).To(BeEmpty(), fmt.Sprintf(defaultMsg, test.name))
	}
}
