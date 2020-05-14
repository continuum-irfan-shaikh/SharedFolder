package models_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration/asset"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-cassandra"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	. "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/triggers"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/cassandra"
)

const scheduleStringIndex = 15

var fail = errors.New("fail")

func TestGetGlobalTaskState(t *testing.T) {
	tests := []struct {
		name      string
		tasks     []models.Task
		wantState TaskState
	}{
		{
			name: `Active`,
			tasks: []models.Task{
				{
					State: TaskStateActive,
				},
				{
					State: TaskStateInactive,
				},
				{
					State: TaskStateDisabled,
				},
			},
			wantState: TaskStateActive,
		},
		{
			name: `Inactive`,
			tasks: []models.Task{
				{
					State: TaskStateInactive,
				},
				{
					State: TaskStateDisabled,
				},
				{
					State: TaskStateInactive,
				},
			},
			wantState: TaskStateInactive,
		},
		{
			name: `Disabled`,
			tasks: []models.Task{
				{
					State: TaskStateDisabled,
				},
				{
					State: TaskStateDisabled,
				},
				{
					State: TaskStateDisabled,
				},
			},
			wantState: TaskStateDisabled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotState := models.GetGlobalTaskState(tt.tasks); !reflect.DeepEqual(gotState, tt.wantState) {
				t.Fatalf("%s: getGlobalTaskState() = %v, want %v", tt.name, gotState, tt.wantState)
			}
		})
	}
}

func Test_HasEndpointTypeOnly(t *testing.T) {
	RegisterTestingT(t)

	Expect(models.TargetsByType{}.HasEndpointTypeOnly()).To(BeFalse())
	Expect(models.TargetsByType{models.ManagedEndpoint: []string{"test"}}.HasEndpointTypeOnly()).To(BeTrue())
	Expect(models.TargetsByType{models.ManagedEndpoint: []string{"test"}, models.Site: []string{"test"}}.HasEndpointTypeOnly()).To(BeFalse())
	Expect(models.TargetsByType{models.ManagedEndpoint: []string{"test"}, models.Site: []string{"test"}, models.DynamicGroup: []string{"test"}}.HasEndpointTypeOnly()).To(BeFalse())
	Expect(models.TargetsByType{models.ManagedEndpoint: []string{"test"}, models.Site: []string{}, models.DynamicGroup: []string{}}.HasEndpointTypeOnly()).To(BeTrue())
	Expect(models.TargetsByType{models.ManagedEndpoint: []string{}, models.Site: []string{}, models.DynamicGroup: []string{}}.HasEndpointTypeOnly()).To(BeFalse())
}

func Test_Validate(t *testing.T) {
	RegisterTestingT(t)

	Expect(models.TargetsByType{}.Validate()).To(BeNil())
	Expect(models.TargetsByType{models.ManagedEndpoint: []string{"test"}}.Validate()).To(BeNil())
	Expect(models.TargetsByType{models.ManagedEndpoint: []string{"test"}, models.Site: []string{"test"}}.Validate()).To(BeNil())
	Expect(models.TargetsByType{models.ManagedEndpoint: []string{"test"}, models.Site: []string{"test"}, models.DynamicGroup: []string{"test"}}.Validate()).To(BeNil())
	Expect(models.TargetsByType{0: []string{"test"}, models.Site: []string{}, models.DynamicGroup: []string{}}.Validate()).NotTo(BeNil())
}

func Test_Contains(t *testing.T) {
	RegisterTestingT(t)

	Expect(models.TargetsByType{}.Contains(models.Site, "test")).To(BeFalse())
	Expect(models.TargetsByType{models.ManagedEndpoint: []string{"test"}}.Contains(models.ManagedEndpoint, "test")).To(BeTrue())
	Expect(models.TargetsByType{models.ManagedEndpoint: []string{"test0"}}.Contains(models.ManagedEndpoint, "test")).To(BeFalse())
}

func Test_UnmarshalJSON(t *testing.T) {
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		json    string
		wantTBT models.TargetsByType
		wantErr bool
	}{
		{
			name: "Success",
			json: `{"SITE": ["50120799"], "MANAGED_ENDPOINT": ["00000000-0000-0000-0000-0000"], "DYNAMIC_GROUP": ["test"]}`,
			wantTBT: models.TargetsByType{
				models.Site:            []string{"50120799"},
				models.ManagedEndpoint: []string{"00000000-0000-0000-0000-0000"},
				models.DynamicGroup:    []string{"test"},
			},
		},
		{
			name:    "Success empty",
			json:    `{}`,
			wantTBT: models.TargetsByType{},
		},
		{
			name:    "Fail empty",
			json:    ``,
			wantErr: true,
		},
		{
			name:    "Fail empty",
			json:    `{"test": ["50120799"]}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tbt := models.TargetsByType{}
		err := tbt.UnmarshalJSON([]byte(tt.json))

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("UnmarshalJSON() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			Expect(tbt).To(Equal(tt.wantTBT), fmt.Sprintf("UnmarshalJSON() name = %s, result = %v, wantResult %v", tt.name, tbt, tt.wantTBT))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("UnmarshalJSON() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_MarshalJSON(t *testing.T) {
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name     string
		tbt      models.TargetsByType
		wantJson string
		wantErr  bool
	}{
		{
			name: "Success site",
			tbt: models.TargetsByType{
				models.Site: []string{"50120799"},
			},
			wantJson: `{"SITE":["50120799"]}`,
		},
		{
			name: "Success endpoint",
			tbt: models.TargetsByType{
				models.ManagedEndpoint: []string{"00000000-0000-0000-0000-0000"},
			},
			wantJson: `{"MANAGED_ENDPOINT":["00000000-0000-0000-0000-0000"]}`,
		},
		{
			name: "Success group",
			tbt: models.TargetsByType{
				models.DynamicGroup: []string{"test"},
			},
			wantJson: `{"DYNAMIC_GROUP":["test"]}`,
		},
		{
			name:     "Success empty",
			tbt:      models.TargetsByType{},
			wantJson: `{}`,
		},
		{
			name: "Fail type",
			tbt: models.TargetsByType{
				999999999: []string{"test"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		json, err := tt.tbt.MarshalJSON()

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("UnmarshalJSON() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			Expect(string(json)).To(Equal(tt.wantJson), fmt.Sprintf("UnmarshalJSON() name = %s, result = %v, wantResult %v", tt.name, string(json), tt.wantJson))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("UnmarshalJSON() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_IsRecurringlyScheduled(t *testing.T) {
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name string
		task models.Task
		want bool
	}{
		{
			name: "True 0",
			task: models.Task{Schedule: tasking.Schedule{
				Regularity:   tasking.Recurrent,
				StartRunTime: time.Unix(0, 0),
				EndRunTime:   time.Unix(1, 0),
			}},
			want: true,
		},
		{
			name: "True 1",
			task: models.Task{Schedule: tasking.Schedule{
				Regularity:   tasking.Recurrent,
				StartRunTime: time.Unix(0, 0),
			}},
			want: true,
		},
		{
			name: "False",
			task: models.Task{Schedule: tasking.Schedule{
				Regularity: tasking.OneTime,
			}},
		},
	}
	for _, tt := range tests {
		result := tt.task.IsRecurringlyScheduled()

		Expect(result).To(Equal(tt.want), fmt.Sprintf("IsRecurringlyScheduled() name = %s, result = %v, wantResult %v", tt.name, result, tt.want))
	}
}

func Test_IsScheduled(t *testing.T) {
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name string
		task models.Task
		want bool
	}{
		{
			name: "True 0",
			task: models.Task{
				State: TaskStateActive,
				Schedule: tasking.Schedule{
					Regularity:   tasking.Recurrent,
					StartRunTime: time.Unix(0, 0),
				}},
			want: true,
		},
		{
			name: "True 1",
			task: models.Task{
				State: TaskStateActive,
				Schedule: tasking.Schedule{
					Regularity:   tasking.OneTime,
					StartRunTime: time.Unix(0, 0),
				}},
			want: true,
		},
		{
			name: "False",
			task: models.Task{
				State: TaskStateDisabled,
				Schedule: tasking.Schedule{
					Regularity:   tasking.Recurrent,
					StartRunTime: time.Unix(0, 0),
				}},
		},
	}
	for _, tt := range tests {
		result := tt.task.IsScheduled()

		Expect(result).To(Equal(tt.want), fmt.Sprintf("IsScheduled() name = %s, result = %v, wantResult %v", tt.name, result, tt.want))
	}
}

func Test_IsTaskAndTriggerNotActivated(t *testing.T) {
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name string
		task models.Task
		want bool
	}{
		{
			name: "True",
			task: models.Task{
				Schedule: tasking.Schedule{
					Regularity:   tasking.Recurrent,
					TriggerTypes: []string{"test"},
					StartRunTime: time.Unix(1, 0).UTC(),
				},
				OriginalNextRunTime: time.Unix(1, 0),
				RunTimeUTC:          time.Unix(1, 0).UTC(),
			},
			want: true,
		},
		{
			name: "False",
			task: models.Task{
				State: TaskStateDisabled,
				Schedule: tasking.Schedule{
					Regularity:   tasking.Recurrent,
					StartRunTime: time.Unix(0, 0),
				}},
		},
	}
	for _, tt := range tests {
		result := tt.task.IsTaskAndTriggerNotActivated()

		Expect(result).To(Equal(tt.want), fmt.Sprintf("IsTaskAndTriggerNotActivated() name = %s, result = %v, wantResult %v", tt.name, result, tt.want))
	}
}

func Test_IsActivatedTrigger(t *testing.T) {
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name string
		task models.Task
		want bool
	}{
		{
			name: "True 0",
			task: models.Task{
				RunTimeUTC: time.Now().Add(100 * time.Hour),
			},
			want: true,
		},
		{
			name: "True 1",
			task: models.Task{
				Schedule: tasking.Schedule{
					Regularity: tasking.Recurrent,
					EndRunTime: time.Unix(1, 0),
				},
				RunTimeUTC: time.Unix(1, 0),
			},
			want: true,
		},
		{
			name: "False",
			task: models.Task{
				Schedule: tasking.Schedule{
					EndRunTime: time.Unix(0, 0),
				}},
		},
	}
	for _, tt := range tests {
		result := tt.task.IsActivatedTrigger()

		Expect(result).To(Equal(tt.want), fmt.Sprintf("IsActivatedTrigger() name = %s, result = %v, wantResult %v", tt.name, result, tt.want))
	}
}

func Test_IsDynamicGroupBasedTrigger(t *testing.T) {
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name string
		task models.Task
		want bool
	}{
		{
			name: "True 0",
			task: models.Task{
				Schedule: tasking.Schedule{
					Regularity:   tasking.Trigger,
					TriggerTypes: []string{triggers.DynamicGroupEnterTrigger},
				},
			},
			want: true,
		},
		{
			name: "True 1",
			task: models.Task{
				Schedule: tasking.Schedule{
					Regularity:   tasking.Trigger,
					TriggerTypes: []string{triggers.DynamicGroupExitTrigger},
				},
			},
			want: true,
		},
		{
			name: "False 0",
			task: models.Task{
				Schedule: tasking.Schedule{
					Regularity: tasking.Trigger,
				},
			},
		},
		{
			name: "False 1",
			task: models.Task{
				Schedule: tasking.Schedule{
					Regularity:   tasking.Trigger,
					TriggerTypes: []string{triggers.FirstCheckInTrigger},
				},
			},
		},
		{
			name: "False 2",
			task: models.Task{
				Schedule: tasking.Schedule{
					Regularity: tasking.Recurrent,
				},
			},
		},
	}
	for _, tt := range tests {

		result := tt.task.IsDynamicGroupBasedTrigger()
		Expect(result).To(Equal(tt.want), fmt.Sprintf("IsDynamicGroupBasedTrigger() name = %s, result = %v, wantResult %v", tt.name, result, tt.want))
	}
}

func Test_IsRunAsUserApplied(t *testing.T) {
	RegisterTestingT(t)

	ttc := models.Task{Credentials: &agent.Credentials{}}
	tt := models.Task{}
	Expect(ttc.IsRunAsUserApplied()).To(BeTrue())
	Expect(tt.IsRunAsUserApplied()).To(BeFalse())
}

func Test_IsExpired(t *testing.T) {
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name string
		task models.Task
		want bool
	}{
		{
			name: "True 0",
			task: models.Task{
				Schedule: tasking.Schedule{
					Regularity: tasking.OneTime,
				},
			},
			want: true,
		},
		{
			name: "True 1",
			task: models.Task{
				Schedule: tasking.Schedule{
					Regularity: tasking.Trigger,
				},
			},
			want: true,
		},
		{
			name: "True 2",
			task: models.Task{},
			want: true,
		},
		{
			name: "False 0",
			task: models.Task{
				Schedule: tasking.Schedule{
					Regularity: tasking.Recurrent,
					EndRunTime: time.Now().Add(100 * time.Hour),
				},
			},
		},
	}
	for _, tt := range tests {
		result := tt.task.IsExpired()

		Expect(result).To(Equal(tt.want), fmt.Sprintf("IsExpired() name = %s, result = %v, wantResult %v", tt.name, result, tt.want))
	}
}

func Test_HasOriginalRunTime(t *testing.T) {
	RegisterTestingT(t)

	ttc := models.Task{OriginalNextRunTime: time.Now()}
	tt := models.Task{}
	Expect(ttc.HasOriginalRunTime()).To(BeTrue())
	Expect(tt.HasOriginalRunTime()).To(BeFalse())
}

func Test_HasPostponedTime(t *testing.T) {
	RegisterTestingT(t)

	ttc := models.Task{PostponedRunTime: time.Now()}
	tt := models.Task{}
	Expect(ttc.HasPostponedTime()).To(BeTrue())
	Expect(tt.HasPostponedTime()).To(BeFalse())
}

func Test_Enable(t *testing.T) {
	var ctrl *gomock.Controller
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		task    models.Task
		asset   func() integration.Asset
		want    bool
		wantErr bool
	}{
		{
			name: "success",
			task: models.Task{
				State:      TaskStateDisabled,
				RunTimeUTC: time.Now().UTC().Add(-100 * time.Hour),
				Schedule: tasking.Schedule{
					Regularity: tasking.Recurrent,
					Repeat: tasking.Repeat{
						Frequency: tasking.Daily,
						Every:     1,
					},
					EndRunTime: time.Now().Add(100 * time.Hour),
				},
			},
			asset: func() integration.Asset {
				a := mock.NewMockAsset(ctrl)
				a.EXPECT().GetLocationByEndpointID(gomock.Any(), gomock.Any(), gomock.Any()).Return(time.UTC, nil)
				return a
			},
			want: true,
		},
		{
			name: "false 0",
			task: models.Task{},
			asset: func() integration.Asset {
				a := mock.NewMockAsset(ctrl)
				return a
			},
			want: false,
		},
		{
			name: "false 0",
			task: models.Task{},
			asset: func() integration.Asset {
				a := mock.NewMockAsset(ctrl)
				return a
			},
			want: false,
		},
		{
			name: "expired",
			task: models.Task{
				State: TaskStateDisabled,
			},
			asset: func() integration.Asset {
				a := mock.NewMockAsset(ctrl)
				return a
			},
			wantErr: true,
		},
		{
			name: "success 1",
			task: models.Task{
				State:      TaskStateDisabled,
				RunTimeUTC: time.Now().UTC().Add(100 * time.Hour),
				Schedule: tasking.Schedule{
					Regularity: tasking.Recurrent,
					Repeat: tasking.Repeat{
						Frequency: tasking.Daily,
						Every:     1,
					},
				},
			},
			asset: func() integration.Asset {
				a := mock.NewMockAsset(ctrl)
				return a
			},
			want: true,
		},
		{
			name: "unable to calc next run time",
			task: models.Task{
				State:      TaskStateDisabled,
				RunTimeUTC: time.Now().UTC().Add(-100 * time.Hour),
				Schedule: tasking.Schedule{
					Regularity: tasking.Recurrent,
				},
			},
			asset: func() integration.Asset {
				a := mock.NewMockAsset(ctrl)
				a.EXPECT().GetLocationByEndpointID(gomock.Any(), gomock.Any(), gomock.Any()).Return(time.UTC, nil)
				return a
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		asset.ServiceInstance = tt.asset()

		result, err := tt.task.Enable(context.Background())
		ctrl.Finish()

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("Enable() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			Expect(result).To(Equal(tt.want), fmt.Sprintf("Enable() name = %s, result = %v, wantResult %v", tt.name, result, tt.want))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("Enable() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_Disable(t *testing.T) {
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		task    models.Task
		want    bool
		wantErr bool
	}{
		{
			name: "True 0",
			task: models.Task{
				State:      TaskStateActive,
				RunTimeUTC: time.Now().UTC().Add(100 * time.Hour),
				Schedule: tasking.Schedule{
					Regularity: tasking.OneTime,
				},
			},
			want: true,
		},
		{
			name: "False 0",
			task: models.Task{
				State: TaskStateInactive,
			},
			want: false,
		},
		{
			name: "expired",
			task: models.Task{
				State: TaskStateActive,
				Schedule: tasking.Schedule{
					Regularity: tasking.OneTime,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		result, err := tt.task.Disable()

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("Disable() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			Expect(result).To(Equal(tt.want), fmt.Sprintf("Disable() name = %s, result = %v, wantResult %v", tt.name, result, tt.want))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("Disable() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func TestTask_CopyWithRunTime(t *testing.T) {
	var task models.Task
	task.Schedule.Regularity = tasking.Trigger
	task.IsTrigger()

	task.Schedule.Regularity = tasking.Recurrent
	task.IsTrigger()

	task.CopyWithRunTime(gocql.TimeUUID())
}

func Test_GetExecutionResultTaskData(t *testing.T) {
	var ctrl *gomock.Controller
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		session func() cassandra.ISession
		want    models.ExecutionResultTaskData
		wantErr bool
	}{
		{
			name: "success",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				q := mocks.NewMockIQuery(ctrl)

				q.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(q).
					Times(1)

				return s
			},
			want: models.ExecutionResultTaskData{},
		},
		{
			name: "failed first",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				q := mocks.NewMockIQuery(ctrl)

				q.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fail).
					Times(1)

				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(q).
					Times(1)

				return s
			},
			wantErr: true,
		},
		{
			name: "first empty",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				q := mocks.NewMockIQuery(ctrl)

				q.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(gocql.ErrNotFound).
					Times(1)

				q.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(q).
					Times(2)

				return s
			},
			want: models.ExecutionResultTaskData{},
		},
		{
			name: "failed second",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				q := mocks.NewMockIQuery(ctrl)

				q.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(gocql.ErrNotFound).
					Times(1)

				q.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fail).
					Times(1)

				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(q).
					Times(2)
				return s
			},
			wantErr: true,
		},
		{
			name: "both empty",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				q := mocks.NewMockIQuery(ctrl)

				q.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(gocql.ErrNotFound).
					Times(2)

				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(q).
					Times(2)

				return s
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		cassandra.Session = tt.session()
		t := models.TaskRepoCassandra{}
		result, err := t.GetExecutionResultTaskData("test", gocql.UUID{}, gocql.UUID{})
		ctrl.Finish()

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("GetExecutionResultTaskData() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
			Expect(result).To(Equal(tt.want), fmt.Sprintf("GetExecutionResultTaskData() name = %s, result = %v, wantResult %v", tt.name, result, tt.want))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("GetExecutionResultTaskData() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_InsertOrUpdate(t *testing.T) {
	var ctrl *gomock.Controller
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		session func() cassandra.ISession
		wantErr bool
	}{
		{
			name: "success",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				b := mocks.NewMockIBatch(ctrl)
				b2 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b)
				s.EXPECT().NewBatch(gomock.Any()).Return(b2)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!select old tasks (selectTasks())

				b.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Size().Return(1)

				s.EXPECT().ExecuteBatch(b).Return(nil)
				s.EXPECT().ExecuteBatch(b2).Return(nil)

				//mutateMVTables()
				b3 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b3).Times(4)
				b3.EXPECT().Query(gomock.Any(), gomock.Any()).Times(4)
				b3.EXPECT().Query(gomock.Any(), gomock.Any()).Times(4)
				s.EXPECT().ExecuteBatch(b3).Return(nil).Times(4)
				//!mutateMVTables()

				s.EXPECT().NewBatch(gomock.Any()).Return(b)
				s.EXPECT().NewBatch(gomock.Any()).Return(b2)
				return s
			},
			wantErr: false,
		},
		{
			name: "select tasks error",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				b := mocks.NewMockIBatch(ctrl)
				b2 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b)
				s.EXPECT().NewBatch(gomock.Any()).Return(b2)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Return(fail).Times(1)
				//!selectTaskCommonData()
				return s
			},
			wantErr: true,
		},
		{
			name: "failed to exec batch",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				b := mocks.NewMockIBatch(ctrl)
				b2 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b)
				s.EXPECT().NewBatch(gomock.Any()).Return(b2)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!select old tasks (selectTasks())

				b.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Query(gomock.Any(), gomock.Any())

				s.EXPECT().ExecuteBatch(b).Return(fail)
				return s
			},
			wantErr: true,
		},
		{
			name: "failed to execute batch2",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				b := mocks.NewMockIBatch(ctrl)
				b2 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b)
				s.EXPECT().NewBatch(gomock.Any()).Return(b2)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!select old tasks (selectTasks())

				b.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Size().Return(1)

				s.EXPECT().ExecuteBatch(b).Return(nil)
				s.EXPECT().ExecuteBatch(b2).Return(fail)
				return s
			},
			wantErr: true,
		},
		{
			name: "fail",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				b := mocks.NewMockIBatch(ctrl)
				b2 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b)
				s.EXPECT().NewBatch(gomock.Any()).Return(b2)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!select old tasks (selectTasks())

				b.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Size().Return(1)

				s.EXPECT().ExecuteBatch(b).Return(nil)
				s.EXPECT().ExecuteBatch(b2).Return(nil)

				//mutateMVTables()
				b3 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b3).MinTimes(1)
				b3.EXPECT().Query(gomock.Any(), gomock.Any()).MinTimes(1)
				s.EXPECT().ExecuteBatch(b3).Return(fail).MinTimes(1)
				//!mutateMVTables()
				return s
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		cassandra.Session = tt.session()
		t := models.TaskRepoCassandra{}
		err := t.InsertOrUpdate(context.Background(), models.Task{})
		ctrl.Finish()

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("InsertOrUpdate() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("InsertOrUpdate() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_Delete(t *testing.T) {
	var ctrl *gomock.Controller
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		session func() cassandra.ISession
		wantErr bool
	}{
		{
			name: "success",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				q := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q).MinTimes(1)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q).MinTimes(1)
				q.EXPECT().Exec().MinTimes(1)
				return s
			},
			wantErr: false,
		},
		{
			name: "fail",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				q := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q).MinTimes(1)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q).MinTimes(1)
				q.EXPECT().Exec().Return(fail).MinTimes(1)
				return s
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		cassandra.Session = tt.session()
		t := models.TaskRepoCassandra{}
		err := t.Delete(context.Background(), []models.Task{{}})
		ctrl.Finish()

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("Delete() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("Delete() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_UpdateTask(t *testing.T) {
	var ctrl *gomock.Controller
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name         string
		tPersist     func() models.TaskPersistence
		tInstPersist func() models.TaskInstancePersistence
		input        func() interface{}
		wantErr      bool
	}{
		{
			name: "success SelectedManagedEndpointEnable",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDAndManagedEndpoints(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC: time.Now().UTC().Add(100 * time.Hour),
						State:      TaskStateDisabled,
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				tp.EXPECT().InsertOrUpdate(gomock.Any(), gomock.Any()).Return(nil)

				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				ti.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return([]models.TaskInstance{{}}, nil)
				ti.EXPECT().GetNearestInstanceAfter(gomock.Any(), gomock.Any()).Return(models.TaskInstance{
					Statuses: map[gocql.UUID]TaskInstanceStatus{}}, nil)
				ti.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
				return ti
			},
			input: func() interface{} {
				return &models.SelectedManagedEndpointEnable{ManagedEndpoints: map[string]bool{"00000000-0000-0000-0000-000000000000": true}}
			},
			wantErr: false,
		},
		{
			name: "success AllTargetsEnable",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC: time.Now().UTC().Add(100 * time.Hour),
						State:      TaskStateDisabled,
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				tp.EXPECT().InsertOrUpdate(gomock.Any(), gomock.Any()).Return(nil)

				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				ti.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return([]models.TaskInstance{{}}, nil)
				ti.EXPECT().GetNearestInstanceAfter(gomock.Any(), gomock.Any()).Return(models.TaskInstance{
					Statuses: map[gocql.UUID]TaskInstanceStatus{}}, nil)
				ti.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
				return ti
			},
			input: func() interface{} {
				return &models.AllTargetsEnable{Active: true}
			},
			wantErr: false,
		},
		{
			name: "success no nearest",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDAndManagedEndpoints(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC: time.Now().UTC().Add(100 * time.Hour),
						State:      TaskStateDisabled,
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				tp.EXPECT().InsertOrUpdate(gomock.Any(), gomock.Any()).Return(nil)

				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				ti.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return([]models.TaskInstance{{
					Statuses: map[gocql.UUID]TaskInstanceStatus{}}}, nil)
				ti.EXPECT().GetNearestInstanceAfter(gomock.Any(), gomock.Any()).Return(models.TaskInstance{}, gocql.ErrNotFound)
				ti.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
				return ti
			},
			input: func() interface{} {
				return &models.SelectedManagedEndpointEnable{ManagedEndpoints: map[string]bool{"00000000-0000-0000-0000-000000000000": true}}
			},
			wantErr: false,
		},
		{
			name: "success postponed",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDAndManagedEndpoints(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC:       time.Now().UTC().Add(100 * time.Hour),
						State:            TaskStateDisabled,
						PostponedRunTime: time.Now().UTC().Add(100 * time.Hour),
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				tp.EXPECT().InsertOrUpdate(gomock.Any(), gomock.Any()).Return(nil)

				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				ti.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return([]models.TaskInstance{{
					Statuses: map[gocql.UUID]TaskInstanceStatus{}}}, nil)
				ti.EXPECT().GetNearestInstanceAfter(gomock.Any(), gomock.Any()).Return(models.TaskInstance{}, gocql.ErrNotFound)
				ti.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
				return ti
			},
			input: func() interface{} {
				return &models.SelectedManagedEndpointEnable{ManagedEndpoints: map[string]bool{"00000000-0000-0000-0000-000000000000": true}}
			},
			wantErr: false,
		},
		{
			name: "success to be disabled",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDAndManagedEndpoints(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC:       time.Now().UTC().Add(100 * time.Hour),
						State:            TaskStateActive,
						PostponedRunTime: time.Now().UTC().Add(100 * time.Hour),
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				tp.EXPECT().InsertOrUpdate(gomock.Any(), gomock.Any()).Return(nil)

				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				ti.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return([]models.TaskInstance{{
					Statuses: map[gocql.UUID]TaskInstanceStatus{}}}, nil)
				ti.EXPECT().GetNearestInstanceAfter(gomock.Any(), gomock.Any()).Return(models.TaskInstance{}, gocql.ErrNotFound)
				ti.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
				return ti
			},
			input: func() interface{} {
				return &models.SelectedManagedEndpointEnable{ManagedEndpoints: map[string]bool{"00000000-0000-0000-0000-000000000000": false}}
			},
			wantErr: false,
		},
		{
			name: "success to be disabled 2",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC:       time.Now().UTC().Add(100 * time.Hour),
						State:            TaskStateActive,
						PostponedRunTime: time.Now().UTC().Add(100 * time.Hour),
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				tp.EXPECT().InsertOrUpdate(gomock.Any(), gomock.Any()).Return(nil)

				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				ti.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return([]models.TaskInstance{{
					Statuses: map[gocql.UUID]TaskInstanceStatus{}}}, nil)
				ti.EXPECT().GetNearestInstanceAfter(gomock.Any(), gomock.Any()).Return(models.TaskInstance{}, gocql.ErrNotFound)
				ti.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
				return ti
			},
			input: func() interface{} {
				return &models.AllTargetsEnable{Active: false}
			},
			wantErr: false,
		},
		{
			name: "failed to enable",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC:       time.Now().UTC().Add(-100 * time.Hour),
						State:            TaskStateDisabled,
						PostponedRunTime: time.Now().UTC().Add(100 * time.Hour),
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				return ti
			},
			input: func() interface{} {
				return &models.AllTargetsEnable{Active: true}
			},
			wantErr: true,
		},
		{
			name: "failed to disable 2",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDAndManagedEndpoints(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC:       time.Now().UTC().Add(-100 * time.Hour),
						State:            TaskStateDisabled,
						PostponedRunTime: time.Now().UTC().Add(100 * time.Hour),
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				return ti
			},
			input: func() interface{} {
				return &models.SelectedManagedEndpointEnable{ManagedEndpoints: map[string]bool{"00000000-0000-0000-0000-000000000000": true}}
			},
			wantErr: true,
		},
		{
			name: "failed - emty tasks",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDAndManagedEndpoints(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC:       time.Now().UTC().Add(-100 * time.Hour),
						State:            TaskStateActive,
						PostponedRunTime: time.Now().UTC().Add(100 * time.Hour),
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				return ti
			},
			input: func() interface{} {
				return &models.SelectedManagedEndpointEnable{ManagedEndpoints: map[string]bool{"00000000-0000-0000-0000-000000000000": true}}
			},
			wantErr: true,
		},
		{
			name: "fail - empty instances",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDAndManagedEndpoints(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC: time.Now().UTC().Add(100 * time.Hour),
						State:      TaskStateDisabled,
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				ti.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return([]models.TaskInstance{}, nil)
				return ti
			},
			input: func() interface{} {
				return &models.SelectedManagedEndpointEnable{ManagedEndpoints: map[string]bool{"00000000-0000-0000-0000-000000000000": true}}
			},
			wantErr: true,
		},
		{
			name: "failed invalid uuid",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				return ti
			},
			input: func() interface{} {
				return &models.SelectedManagedEndpointEnable{ManagedEndpoints: map[string]bool{"test": true}}
			},
			wantErr: true,
		},
		{
			name: "failed GetByIDAndManagedEndpoints",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDAndManagedEndpoints(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{}, fail)
				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				return ti
			},
			input: func() interface{} {
				return &models.SelectedManagedEndpointEnable{ManagedEndpoints: map[string]bool{"00000000-0000-0000-0000-000000000000": true}}
			},
			wantErr: true,
		},
		{
			name: "failed to GetByIDs",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{}, fail)
				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				return ti
			},
			input: func() interface{} {
				return &models.AllTargetsEnable{Active: true}
			},
			wantErr: true,
		},
		{
			name: "invalid input",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				return ti
			},
			input: func() interface{} {
				return fail
			},
			wantErr: true,
		},
		{
			name: "failed to get instances by id",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDAndManagedEndpoints(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC: time.Now().UTC().Add(100 * time.Hour),
						State:      TaskStateDisabled,
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				ti.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return([]models.TaskInstance{}, fail)
				return ti
			},
			input: func() interface{} {
				return &models.SelectedManagedEndpointEnable{ManagedEndpoints: map[string]bool{"00000000-0000-0000-0000-000000000000": true}}
			},
			wantErr: true,
		},
		{
			name: "failed to get nearest instance",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDAndManagedEndpoints(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC: time.Now().UTC().Add(100 * time.Hour),
						State:      TaskStateDisabled,
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				ti.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return([]models.TaskInstance{{}}, nil)
				ti.EXPECT().GetNearestInstanceAfter(gomock.Any(), gomock.Any()).Return(models.TaskInstance{}, fail)
				return ti
			},
			input: func() interface{} {
				return &models.SelectedManagedEndpointEnable{ManagedEndpoints: map[string]bool{"00000000-0000-0000-0000-000000000000": true}}
			},
			wantErr: true,
		},
		{
			name: "failed to insert instance",
			tPersist: func() models.TaskPersistence {
				tp := mock.NewMockTaskPersistence(ctrl)
				tp.EXPECT().GetByIDAndManagedEndpoints(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Task{{
						RunTimeUTC: time.Now().UTC().Add(100 * time.Hour),
						State:      TaskStateDisabled,
						Schedule: tasking.Schedule{
							Regularity: tasking.OneTime,
							Repeat: tasking.Repeat{
								RunTime: time.Now().Add(time.Hour),
							},
						},
					}}, nil)
				return tp
			},
			tInstPersist: func() models.TaskInstancePersistence {
				ti := mock.NewMockTaskInstancePersistence(ctrl)
				ti.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return([]models.TaskInstance{{}}, nil)
				ti.EXPECT().GetNearestInstanceAfter(gomock.Any(), gomock.Any()).Return(models.TaskInstance{
					Statuses: map[gocql.UUID]TaskInstanceStatus{}}, nil)
				ti.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(fail)
				return ti
			},
			input: func() interface{} {
				return &models.SelectedManagedEndpointEnable{ManagedEndpoints: map[string]bool{"00000000-0000-0000-0000-000000000000": true}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		models.TaskPersistenceInstance = tt.tPersist()
		models.TaskInstancePersistenceInstance = tt.tInstPersist()

		t := models.TaskRepoCassandra{}
		err := t.UpdateTask(context.Background(), tt.input(), "test", gocql.UUID{})
		ctrl.Finish()

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("UpdateTask() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("UpdateTask() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_UpdateSchedulerFields(t *testing.T) {
	var ctrl *gomock.Controller
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		session func() cassandra.ISession
		wantErr bool
	}{
		{
			name: "success",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				b := mocks.NewMockIBatch(ctrl)
				b2 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b)
				s.EXPECT().NewBatch(gomock.Any()).Return(b2)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!select old tasks (selectTasks())

				b.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Size().Return(1)

				s.EXPECT().ExecuteBatch(b).Return(nil)
				s.EXPECT().ExecuteBatch(b2).Return(nil)

				//mutateMVTables()
				b3 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b3).Times(4)
				b3.EXPECT().Query(gomock.Any(), gomock.Any()).Times(4)
				b3.EXPECT().Query(gomock.Any(), gomock.Any()).Times(4)
				s.EXPECT().ExecuteBatch(b3).Return(nil).Times(4)
				//!mutateMVTables()

				s.EXPECT().NewBatch(gomock.Any()).Return(b)
				s.EXPECT().NewBatch(gomock.Any()).Return(b2)
				return s
			},
			wantErr: false,
		},
		{
			name: "select tasks error",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				b := mocks.NewMockIBatch(ctrl)
				b2 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b)
				s.EXPECT().NewBatch(gomock.Any()).Return(b2)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Return(fail).Times(1)
				//!selectTaskCommonData()
				return s
			},
			wantErr: true,
		},
		{
			name: "failed to exec batch",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				b := mocks.NewMockIBatch(ctrl)
				b2 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b)
				s.EXPECT().NewBatch(gomock.Any()).Return(b2)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!select old tasks (selectTasks())

				b.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Query(gomock.Any(), gomock.Any())

				s.EXPECT().ExecuteBatch(b).Return(fail)
				return s
			},
			wantErr: true,
		},
		{
			name: "failed to execute batch2",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				b := mocks.NewMockIBatch(ctrl)
				b2 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b)
				s.EXPECT().NewBatch(gomock.Any()).Return(b2)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!select old tasks (selectTasks())

				b.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Size().Return(1)

				s.EXPECT().ExecuteBatch(b).Return(nil)
				s.EXPECT().ExecuteBatch(b2).Return(fail)
				return s
			},
			wantErr: true,
		},
		{
			name: "fail",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				b := mocks.NewMockIBatch(ctrl)
				b2 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b)
				s.EXPECT().NewBatch(gomock.Any()).Return(b2)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!select old tasks (selectTasks())

				b.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Query(gomock.Any(), gomock.Any())
				b2.EXPECT().Size().Return(1)

				s.EXPECT().ExecuteBatch(b).Return(nil)
				s.EXPECT().ExecuteBatch(b2).Return(nil)

				//mutateMVTables()
				b3 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b3).MinTimes(1)
				b3.EXPECT().Query(gomock.Any(), gomock.Any()).MinTimes(1)
				s.EXPECT().ExecuteBatch(b3).Return(fail).MinTimes(1)
				//!mutateMVTables()
				return s
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		cassandra.Session = tt.session()
		t := models.TaskRepoCassandra{}
		err := t.UpdateSchedulerFields(context.Background(), models.Task{})
		ctrl.Finish()

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("InsertOrUpdate() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("InsertOrUpdate() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_UpdateModifiedFieldsByMEs(t *testing.T) {
	var ctrl *gomock.Controller
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		session func() cassandra.ISession
		wantErr bool
	}{
		{
			name: "success",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!select old tasks (selectTasks())

				q2 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q2)
				q2.EXPECT().Exec().Return(nil)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q2)
				q2.EXPECT().Exec().Return(nil)

				//mutateMVTables()
				b3 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b3).Times(4)
				b3.EXPECT().Query(gomock.Any(), gomock.Any()).Times(4)
				b3.EXPECT().Query(gomock.Any(), gomock.Any()).Times(4)
				s.EXPECT().ExecuteBatch(b3).Return(nil).Times(4)
				//!mutateMVTables()

				return s
			},
			wantErr: false,
		},
		{
			name: "failed to select",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Return(fail).Times(1)
				//!selectTaskCommonData()
				return s
			},
			wantErr: true,
		},
		{
			name: "failed to update 1",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!select old tasks (selectTasks())

				q2 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q2)
				q2.EXPECT().Exec().Return(fail)
				return s
			},
			wantErr: true,
		},
		{
			name: "failed to update 2",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!select old tasks (selectTasks())

				q2 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q2)
				q2.EXPECT().Exec().Return(nil)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q2)
				q2.EXPECT().Exec().Return(fail)
				return s
			},
			wantErr: true,
		},
		{
			name: "failed to mutate mvs",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				//select old tasks (selectTasks())
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!select old tasks (selectTasks())

				q2 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q2)
				q2.EXPECT().Exec().Return(nil)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q2)
				q2.EXPECT().Exec().Return(nil)

				//mutateMVTables()
				b3 := mocks.NewMockIBatch(ctrl)
				s.EXPECT().NewBatch(gomock.Any()).Return(b3).AnyTimes()
				b3.EXPECT().Query(gomock.Any(), gomock.Any()).AnyTimes()
				b3.EXPECT().Query(gomock.Any(), gomock.Any()).AnyTimes()
				s.EXPECT().ExecuteBatch(b3).Return(fail).AnyTimes()
				//!mutateMVTables()

				return s
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		cassandra.Session = tt.session()
		t := models.TaskRepoCassandra{}
		err := t.UpdateModifiedFieldsByMEs(context.Background(), models.Task{}, gocql.UUID{})
		ctrl.Finish()

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("InsertOrUpdate() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("InsertOrUpdate() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_GetByIDs(t *testing.T) {
	var ctrl *gomock.Controller
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		cache   func() persistency.Cache
		session func() cassandra.ISession
		wantErr bool
	}{
		{
			name: "success",
			cache: func() persistency.Cache {
				c := mock.NewMockCache(ctrl)
				c.EXPECT().Get(gomock.Any()).Return(nil, fail)
				c.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return c
			},
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				//selectTasks()
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!selectTasks()

				return s
			},
			wantErr: false,
		},
		{
			name: "success w/o cache",
			cache: func() persistency.Cache {
				return nil
			},
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				//selectTasks()
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!selectTasks()

				return s
			},
			wantErr: false,
		},
		{
			name: "failed to select tasks",
			cache: func() persistency.Cache {
				c := mock.NewMockCache(ctrl)
				c.EXPECT().Get(gomock.Any()).Return(nil, fail)
				return c
			},
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				//selectTasks()
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Return(fail).Times(1)
				//!selectTaskCommonData()

				//!selectTasks()

				return s
			},
			wantErr: true,
		},
		{
			name: "failed to select tasks w/o cache",
			cache: func() persistency.Cache {
				return nil
			},
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				//selectTasks()
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Return(fail).Times(1)
				//!selectTaskCommonData()

				//!selectTasks()

				return s
			},
			wantErr: true,
		},
		{
			name: "success from cache",
			cache: func() persistency.Cache {
				c := mock.NewMockCache(ctrl)
				c.EXPECT().Get(gomock.Any()).Return([]byte("{}"), nil)
				return c
			},
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)
				return s
			},
			wantErr: false,
		},
		{
			name: "failed to set cache",
			cache: func() persistency.Cache {
				c := mock.NewMockCache(ctrl)
				c.EXPECT().Get(gomock.Any()).Return(nil, fail)
				c.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(fail)
				return c
			},
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				//selectTasks()
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!selectTasks()

				return s
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		config.Config.AssetCacheEnabled = true

		cassandra.Session = tt.session()
		t := models.TaskRepoCassandra{}
		_, err := t.GetByIDs(context.Background(), tt.cache(), "test", true, gocql.UUID{})
		ctrl.Finish()

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("InsertOrUpdate() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("InsertOrUpdate() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_GetByIDAndManagedEndpoints(t *testing.T) {
	var ctrl *gomock.Controller
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		session func() cassandra.ISession
		wantErr bool
	}{
		{
			name: "success",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				//selectTasks()
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)
				//!selectTaskCommonData()

				i.EXPECT().Scan(gomock.Any()).Return(false)
				//!selectTasks()

				return s
			},
		},
		{
			name: "failed",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				//selectTasks()
				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				//selectTaskCommonData()
				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Return(fail).Times(1)
				//!selectTaskCommonData()

				return s
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		cassandra.Session = tt.session()
		t := models.TaskRepoCassandra{}
		_, err := t.GetByIDAndManagedEndpoints(context.Background(), "test", gocql.UUID{}, gocql.UUID{})
		ctrl.Finish()

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("GetExecutionResultTaskData() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("GetExecutionResultTaskData() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}

func Test_GetByRunTimeRange(t *testing.T) {
	var ctrl = gomock.NewController(t)
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	s := mocks.NewMockISession(ctrl)

	//selectTasks()
	q := mocks.NewMockIQuery(ctrl)
	i := mocks.NewMockIIter(ctrl)
	s.EXPECT().Query(gomock.Any(), gomock.Any()).Return(q)
	q.EXPECT().Iter().Return(i)
	i.EXPECT().Scan(gomock.Any()).Return(true)

	//selectTaskCommonData()
	q1 := mocks.NewMockIQuery(ctrl)
	s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
	q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
		schedule := tasking.Schedule{
			Regularity: 1,
			Repeat: tasking.Repeat{
				RunTime: time.Now().Add(time.Hour),
			},
		}

		b, err := json.Marshal(schedule)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// set &scheduleString to pass validation
		*(p[scheduleStringIndex].(*string)) = string(b)
	}).Return(nil).Times(1)
	//!selectTaskCommonData()

	i.EXPECT().Scan(gomock.Any()).Return(false)
	//!selectTasks()

	cassandra.Session = s
	tr := models.TaskRepoCassandra{}
	tr.GetByRunTimeRange(context.Background(), []time.Time{time.Now()})
	ctrl.Finish()
}

func Test_GetByPartnerAndTime(t *testing.T) {
	var ctrl = gomock.NewController(t)
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	s := mocks.NewMockISession(ctrl)

	//selectTasks()
	q := mocks.NewMockIQuery(ctrl)
	i := mocks.NewMockIIter(ctrl)
	s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
	q.EXPECT().Iter().Return(i)
	i.EXPECT().Scan(gomock.Any()).Return(true)

	//selectTaskCommonData()
	q1 := mocks.NewMockIQuery(ctrl)
	s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
	q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
		schedule := tasking.Schedule{
			Regularity: 1,
			Repeat: tasking.Repeat{
				RunTime: time.Now().Add(time.Hour),
			},
		}

		b, err := json.Marshal(schedule)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// set &scheduleString to pass validation
		*(p[scheduleStringIndex].(*string)) = string(b)
	}).Return(nil).Times(1)
	//!selectTaskCommonData()

	i.EXPECT().Scan(gomock.Any()).Return(false)
	//!selectTasks()

	cassandra.Session = s
	tr := models.TaskRepoCassandra{}
	tr.GetByPartnerAndTime(context.Background(), "test", time.Now())
	ctrl.Finish()
}

func Test_GetByPartnerAndManagedEndpointID(t *testing.T) {
	var ctrl = gomock.NewController(t)
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	s := mocks.NewMockISession(ctrl)

	//selectTasks()
	q := mocks.NewMockIQuery(ctrl)
	i := mocks.NewMockIIter(ctrl)
	s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
	q.EXPECT().Iter().Return(i)
	i.EXPECT().Scan(gomock.Any()).Return(true)

	//selectTaskCommonData()
	q1 := mocks.NewMockIQuery(ctrl)
	s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
	q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
		schedule := tasking.Schedule{
			Regularity: 1,
			Repeat: tasking.Repeat{
				RunTime: time.Now().Add(time.Hour),
			},
		}

		b, err := json.Marshal(schedule)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// set &scheduleString to pass validation
		*(p[scheduleStringIndex].(*string)) = string(b)
	}).Return(nil).Times(1)
	//!selectTaskCommonData()

	i.EXPECT().Scan(gomock.Any()).Return(false)
	//!selectTasks()

	cassandra.Session = s
	tr := models.TaskRepoCassandra{}
	tr.GetByPartnerAndManagedEndpointID(context.Background(), "test", gocql.UUID{}, 1)
	ctrl.Finish()
}

func Test_GetByLastTaskInstanceIDs(t *testing.T) {
	var ctrl *gomock.Controller
	RegisterTestingT(t)
	logger.Load(config.Config.Log)

	tests := []struct {
		name    string
		session func() cassandra.ISession
		wantErr bool
	}{
		{
			name: "success",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)

				i.EXPECT().Scan(gomock.Any()).Return(false)
				i.EXPECT().Close().Return(nil)

				return s
			},
		},
		{
			name: "fail 1",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Return(fail).Times(1)

				return s
			},
			wantErr: true,
		},
		{
			name: "fail 2",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					*(p[scheduleStringIndex].(*string)) = "test"
				}).Return(nil).Times(1)

				i.EXPECT().Scan(gomock.Any()).Return(false)
				i.EXPECT().Close().Return(nil)
				return s
			},
		},
		{
			name: "fail 3",
			session: func() cassandra.ISession {
				s := mocks.NewMockISession(ctrl)

				q := mocks.NewMockIQuery(ctrl)
				i := mocks.NewMockIIter(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(q)
				q.EXPECT().Iter().Return(i)
				i.EXPECT().Scan(gomock.Any()).Return(true)

				q1 := mocks.NewMockIQuery(ctrl)
				s.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(q1)
				q1.EXPECT().Scan(gomock.Any()).Do(func(p ...interface{}) {
					schedule := tasking.Schedule{
						Regularity: 1,
						Repeat: tasking.Repeat{
							RunTime: time.Now().Add(time.Hour),
						},
					}

					b, err := json.Marshal(schedule)
					if err != nil {
						t.Fatalf(err.Error())
					}

					// set &scheduleString to pass validation
					*(p[scheduleStringIndex].(*string)) = string(b)
				}).Return(nil).Times(1)

				i.EXPECT().Scan(gomock.Any()).Return(false)
				i.EXPECT().Close().Return(fail)

				return s
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctrl = gomock.NewController(t)
		cassandra.Session = tt.session()
		t := models.TaskRepoCassandra{}
		_, err := t.GetByLastTaskInstanceIDs(context.Background(), "test", gocql.UUID{}, gocql.UUID{})
		ctrl.Finish()

		if !tt.wantErr {
			Expect(err).To(BeNil(), fmt.Sprintf("GetExecutionResultTaskData() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		} else {
			Expect(err).ToNot(BeNil(), fmt.Sprintf("GetExecutionResultTaskData() name = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr))
		}
	}
}
