package models_test

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
)

var someOriginID, _ = gocql.ParseUUID("049d82af-6b41-4eae-9489-34a77b29e864")
var startRunTime, _ = time.Parse(time.RFC3339, "2030-05-20T07:40:00.000Z")
var endRunTime = time.Now().Add(5 * time.Minute)
var betweenEndTime, _ = time.Parse(time.RFC3339, "2030-05-20T09:40:00.000Z")
var runTime, _ = time.Parse(time.RFC3339, "2030-05-20T08:40:00.000Z")

func init()  {
	logger.Load(config.Config.Log)
}

func TestTemplateCacheLocal_GetAllTemplates(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrlr *gomock.Controller

	type payload struct {
		hasNOC     bool
		cache      func()
	}

	type expected struct {
		templates []models.Template
		err       string
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "error with empty templates",
			payload: payload{
				hasNOC: false,
				cache: func() {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					cacheMock.EXPECT().Get([]byte(gocql.UUID{}.String())).Return(nil, errors.New("some_err")).Times(1)
					cacheMock.EXPECT().Get([]byte("partnerID")).Return([]byte("}}"), nil).Times(1)
					models.TemplatesCache = cacheMock
				},
			},
			expected: expected{
				err: "couldn't retrieve templates",
			},
		},
		{
			name: "Success",
			payload: payload{
				hasNOC: false,
				cache: func() {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					cacheMock.EXPECT().Get([]byte(gocql.UUID{}.String())).Return(nil, errors.New("some_err")).Times(1)
					cacheMock.EXPECT().Get([]byte("partnerID")).Return([]byte("[{\"isRequireNOCAccess\":false,"+
						"\"originId\":\"049d82af-6b41-4eae-9489-34a77b29e864\",\"partnerId\":\"partnerID\","+
						"\"name\":\"someName\"},"+
						"{\"isRequireNOCAccess\":true,"+
						"\"originId\":\"3f8f2eb5-57f8-455b-9007-8b00851b1bd1\",\"partnerId\":\"partnerID\","+
						"\"name\":\"templateName\"}]"), nil).Times(1)

					models.TemplatesCache = cacheMock
				},
			},
			expected: expected{
				templates: []models.Template{
					{
						PartnerID: "partnerID",
						OriginID:  someOriginID,
						Name:      "someName",
					},
				},
			},
		},
	}

	for _, test := range tc {
		mockCtrlr = gomock.NewController(t)
		test.payload.cache()

		t := models.TemplateCacheLocal{}
		templates, err := t.GetAllTemplates(context.Background(), "partnerID", test.payload.hasNOC)
		mockCtrlr.Finish()

		if err == nil {
			Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
			Ω(templates).To(ConsistOf(test.expected.templates), fmt.Sprintf(defaultMsg, test.name))
			continue
		}
		Ω(err.Error()).To(Equal(test.expected.err), fmt.Sprintf(defaultMsg, test.name))
		Ω(templates).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
	}
}

func TestTemplateCacheLocal_GetByType(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrlr *gomock.Controller

	type payload struct {
		taskType string
		hasNoc   bool
		cache    func()
	}

	type expected struct {
		templates []models.Template
		err       string
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "error while getting all templates",
			payload: payload{
				hasNoc: true,
				cache: func() {
					cacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					cacheMock.EXPECT().GetAllTemplates(context.Background(), "partnerID", true).
						Return(nil, errors.New("some_err")).Times(1)
					models.TemplateCacheInstance = cacheMock
				},
			},
			expected: expected{
				err: "No Task Definition Templates with Type:  for Partner: partnerID, err: some_err",
			},
		},
		{
			name: "success",
			payload: payload{
				hasNoc:   false,
				taskType: "someType",
				cache: func() {
					cacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					cacheMock.EXPECT().GetAllTemplates(context.Background(), "partnerID", false).
						Return([]models.Template{
							{
								IsRequireNOCAccess: true,
							},
							{
								Type: "someType",
							},
						}, nil).Times(1)
					models.TemplateCacheInstance = cacheMock
				},
			},
			expected: expected{
				templates: []models.Template{
					{
						Type: "someType",
					},
				},
			},
		},
	}

	for _, test := range tc {
		mockCtrlr = gomock.NewController(t)
		test.payload.cache()

		t := models.TemplateCacheLocal{}
		templates, err := t.GetByType(context.Background(), "partnerID", test.payload.taskType, test.payload.hasNoc)
		mockCtrlr.Finish()

		if err == nil {
			Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
			Ω(templates).To(ConsistOf(test.expected.templates), fmt.Sprintf(defaultMsg, test.name))
			continue
		}
		Ω(err.Error()).To(Equal(test.expected.err), fmt.Sprintf(defaultMsg, test.name))
		Ω(templates).To(BeEmpty(), fmt.Sprintf(defaultMsg, test.name))
	}
}

func TestTemplateCacheLocal_ExistsWithName(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrlr *gomock.Controller
	type payload struct {
		cache func()
	}

	type expected struct {
		exists bool
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "error",
			payload: payload{
				cache: func() {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					cacheMock.EXPECT().Get([]byte("TKS_TEMPLATE_NAME_scriptName")).Return(nil, errors.New("some_err")).Times(1)
					cacheMock.EXPECT().Get([]byte("TKS_PARTNER_partnerID_TEMPLATE_NAME_scriptName")).Return(nil, errors.New("some_err")).Times(1)
					models.TemplatesCache = cacheMock
				},
			},
			expected: expected{
				exists: false,
			},
		},
		{
			name: "Success_#1",
			payload: payload{
				cache: func() {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					cacheMock.EXPECT().Get([]byte("TKS_TEMPLATE_NAME_scriptName")).Return(nil, errors.New("some_err")).Times(1)
					cacheMock.EXPECT().Get([]byte("TKS_PARTNER_partnerID_TEMPLATE_NAME_scriptName")).Return([]byte("someTemplate"), nil).Times(1)
					models.TemplatesCache = cacheMock
				},
			},
			expected: expected{
				exists: true,
			},
		},
		{
			name: "Success_#2",
			payload: payload{
				cache: func() {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					cacheMock.EXPECT().Get([]byte("TKS_TEMPLATE_NAME_scriptName")).Return([]byte("someTemplate"), nil).Times(1)
					models.TemplatesCache = cacheMock
				},
			},
			expected: expected{
				exists: true,
			},
		},
	}

	for _, test := range tc {
		mockCtrlr = gomock.NewController(t)
		test.payload.cache()

		t := models.TemplateCacheLocal{}
		exist := t.ExistsWithName(context.Background(), "partnerID", "scriptName")
		mockCtrlr.Finish()

		Ω(exist).To(Equal(test.expected.exists), fmt.Sprintf(defaultMsg, test.name))
	}
}

func TestTemplateCacheLocal_GetAllTemplatesDetails(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrlr *gomock.Controller

	type payload struct {
		cache  func()
	}

	type expected struct {
		details []models.TemplateDetails
		err     string
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "error",
			payload: payload{
				cache: func() {
					cacheMock := mocks.NewMockCache(mockCtrlr)
					cacheMock.EXPECT().Get([]byte(gocql.UUID{}.String())).Return(nil, errors.New("some_err")).Times(1)
					cacheMock.EXPECT().Get([]byte("partnerID")).Return(nil, errors.New("some_err")).Times(1)
					models.TemplatesCache = cacheMock
				},
			},
			expected: expected{
				err: "couldn't retrieve templates",
			},
		},
	}

	for _, test := range tc {
		mockCtrlr = gomock.NewController(t)
		test.payload.cache()

		t := models.TemplateCacheLocal{}
		d, err := t.GetAllTemplatesDetails(context.Background(), "partnerID")
		mockCtrlr.Finish()

		if err == nil {
			Ω(err).To(BeNil(), fmt.Sprintf(defaultMsg, test.name))
			Ω(d).To(ConsistOf(test.expected.details), fmt.Sprintf(defaultMsg, test.name))
			continue
		}
		Ω(err.Error()).To(Equal(test.expected.err), fmt.Sprintf(defaultMsg, test.name))
		Ω(d).To(BeEmpty(), fmt.Sprintf(defaultMsg, test.name))
	}
}

func TestTemplateCacheLocal_CalculateExpectedExecutionTimeSec(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrlr *gomock.Controller

	type payload struct {
		task   models.Task
		cache  func()
	}

	type expected struct {
		expectedTime int
	}

	tc := []struct {
		name string
		payload
		expected
	}{
		{
			name: "unmarshal error",
			payload: payload{
				task: models.Task{
					Parameters: "}\"expectedExecutionTimeSec\":200}",
				},
				cache: func() {
					cacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					models.TemplateCacheInstance = cacheMock
				},
			},
			expected: expected{
				expectedTime: 300,
			},
		},
		{
			name: "Success_#1",
			payload: payload{
				task: models.Task{
					Parameters: "{\"expectedExecutionTimeSec\":200}",
				},
				cache: func() {
					cacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					models.TemplateCacheInstance = cacheMock
				},
			},
			expected: expected{
				expectedTime: 200,
			},
		},
		{
			name: "Success_#2",
			payload: payload{
				task: models.Task{
					Schedule: tasking.Schedule{
						Regularity: tasking.RunNow,
						EndRunTime: endRunTime,
					},
				},
				cache: func() {
					cacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					models.TemplateCacheInstance = cacheMock
				},
			},
			expected: expected{
				expectedTime: 300,
			},
		},
		{
			name: "Success_#3",
			payload: payload{
				task: models.Task{
					PartnerID: "partnerID",
					Type:      models.TaskTypeScript,
					OriginID:  someOriginID,
				},
				cache: func() {
					cacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					cacheMock.EXPECT().GetAllTemplatesDetails(context.Background(), "partnerID").
						Return(nil, errors.New("some_err")).Times(1)
					models.TemplateCacheInstance = cacheMock
				},
			},
			expected: expected{
				expectedTime: 86400,
			},
		},
		{
			name: "Success_#4",
			payload: payload{
				task: models.Task{
					PartnerID: "partnerID",
					Type:      models.TaskTypeScript,
					OriginID:  someOriginID,
					Schedule: tasking.Schedule{
						Regularity:     tasking.OneTime,
						BetweenEndTime: betweenEndTime,
						Repeat: tasking.Repeat{
							RunTime: runTime,
						},
					},
				},
				cache: func() {
					cacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					cacheMock.EXPECT().GetAllTemplatesDetails(context.Background(), "partnerID").
						Return(nil, errors.New("some_err")).Times(1)
					models.TemplateCacheInstance = cacheMock
				},
			},
			expected: expected{
				expectedTime: 90000,
			},
		},
		{
			name: "Success_#4",
			payload: payload{
				task: models.Task{
					PartnerID: "partnerID",
					Type:      models.TaskTypeScript,
					OriginID:  someOriginID,
					Schedule: tasking.Schedule{
						Regularity:     tasking.OneTime,
						BetweenEndTime: betweenEndTime,
						StartRunTime:   startRunTime,
					},
				},
				cache: func() {
					cacheMock := mocks.NewMockTemplateCache(mockCtrlr)
					cacheMock.EXPECT().GetAllTemplatesDetails(context.Background(), "partnerID").
						Return(nil, errors.New("some_err")).Times(1)
					models.TemplateCacheInstance = cacheMock
				},
			},
			expected: expected{
				expectedTime: 93600,
			},
		},
	}

	for _, test := range tc {
		fmt.Println(test.name)
		mockCtrlr = gomock.NewController(t)
		test.payload.cache()

		t := models.TemplateCacheLocal{}
		timeSec := t.CalculateExpectedExecutionTimeSec(context.Background(), test.payload.task)
		mockCtrlr.Finish()

		Ω(test.expected.expectedTime-timeSec).Should(And(BeNumerically("<=", 10), BeNumerically(">=", 0)), fmt.Sprintf(defaultMsg, test.name))
	}
}
