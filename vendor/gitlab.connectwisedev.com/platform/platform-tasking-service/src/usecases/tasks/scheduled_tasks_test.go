package tasks

import (
	"context"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	e "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/usecases/trigger"
)

const (
	partnerID      = "1"
	taskID1        = "id1"
	taskID2        = "id2"
	taskID1Trigger = "id1t"
	taskID2Trigger = "id2t"
	tiID1          = "123e4567-e89b-12d3-a456-426655440000"
	tiID1Trigger   = "123e4567-e89b-12d3-a456-426655440000"
	tiID2Trigger   = "223e4567-e89b-12d3-a456-426655440000"
	tiID3Trigger   = "323e4567-e89b-12d3-a456-426655440000"
	endpID         = "endp1"
	endpID2        = "endp2"
)

func init()  {
	logger.Load(config.Config.Log)
}

func TestTasks_GetScheduledTasks(t *testing.T) {
	RegisterTestingT(t)

	timeUUID := gocql.TimeUUID()
	ctx := context.WithValue(context.Background(), "s", "s")
	ctxWithUser := context.WithValue(ctx, config.UserKeyCTX, e.User{})
	goodCtx := context.WithValue(ctxWithUser, config.PartnerIDKeyCTX, partnerID)

	now := time.Now().UTC()
	tomorrow := now.Add(time.Hour * (24))

	testCases := []struct {
		name          string
		expectedError bool
		expectedData  []e.ScheduledTasks
		mock          func(ctl *gomock.Controller) (Repo, InstancesRepo, logger.Logger)
		ctx           context.Context
	}{
		{
			name:          "testCase 1 - cant get user from context",
			expectedError: true,
			mock: func(ctl *gomock.Controller) (repo Repo, ti InstancesRepo, log logger.Logger) {
				return nil, nil, nil
			},
			ctx: ctx,
		},
		{
			name:          "testCase 2 - cant get partner from context",
			expectedError: true,
			mock: func(ctl *gomock.Controller) (repo Repo, ti InstancesRepo, log logger.Logger) {
				return nil, nil, nil
			},
			ctx: ctxWithUser,
		},
		{
			name:          "testCase 3 - cant get partner from context",
			expectedError: true,
			mock: func(ctl *gomock.Controller) (Repo, InstancesRepo, logger.Logger) {
				tr := NewMockRepo(ctl)
				tr.EXPECT().GetScheduledTasks(partnerID).Return([]e.ScheduledTasks{}, errors.New("err"))
				return tr, nil, nil
			},
			ctx: goodCtx,
		},
		{
			name: "testCase 4 - cant get instences from DB, but still ok",
			mock: func(ctl *gomock.Controller) (Repo, InstancesRepo, logger.Logger) {
				tr := NewMockRepo(ctl)
				ti := NewMockInstancesRepo(ctl)

				err := errors.New("err")
				tr.EXPECT().GetScheduledTasks(partnerID).Return([]e.ScheduledTasks{{ID: taskID1, LastTaskInstanceID: tiID1, Regularity: tasking.OneTime}}, nil)
				ti.EXPECT().GetInstancesForScheduled([]string{tiID1}).Return([]e.TaskInstance{{TaskID: taskID1, ID: tiID1, Statuses: map[string]statuses.TaskInstanceStatus{"tiID1": statuses.TaskInstanceScheduled}}}, err)
				return tr, ti, nil
			},
			expectedData: []e.ScheduledTasks{
				{
					ID:                 taskID1,
					OverallStatus:      statuses.OverallNew,
					CanBeCanceled:      true,
					LastTaskInstanceID: tiID1,
					ExecutionInfo: e.ExecutionInfo{
						DeviceCount: 1,
					},
					Regularity: tasking.OneTime,
				},
			},
			ctx: goodCtx,
		},
		{
			name: "testCase 4.1 - cant get trigger instences from DB, but still ok",
			mock: func(ctl *gomock.Controller) (Repo, InstancesRepo, logger.Logger) {
				tr := NewMockRepo(ctl)
				ti := NewMockInstancesRepo(ctl)

				err := errors.New("err")
				tr.EXPECT().GetScheduledTasks(partnerID).Return([]e.ScheduledTasks{
					{ID: taskID1, LastTaskInstanceID: tiID1, Regularity: tasking.OneTime},
					{
						ID:                 taskID1Trigger,
						Regularity:         tasking.Trigger,
						LastTaskInstanceID: tiID1Trigger,
						TriggerTypes:       []string{"logOn"},
						TriggerFrames: []tasking.TriggerFrame{{
							TriggerType: "logOn",
						}},
					}}, nil)
				ti.EXPECT().GetInstancesForScheduled([]string{tiID1, tiID1Trigger}).Return([]e.TaskInstance{{TaskID: taskID1, ID: tiID1, Statuses: map[string]statuses.TaskInstanceStatus{"tiID1": statuses.TaskInstanceScheduled}}}, nil)
				ti.EXPECT().GetTopInstancesForScheduledByTaskIDs(gomock.Any()).Return([]e.TaskInstance{}, err)
				return tr, ti, nil
			},
			expectedData: []e.ScheduledTasks{
				{
					ID:                 taskID1,
					OverallStatus:      statuses.OverallNew,
					CanBeCanceled:      true,
					LastTaskInstanceID: tiID1,
					ExecutionInfo: e.ExecutionInfo{
						DeviceCount: 1,
					},
					Regularity: tasking.OneTime,
				},
				{
					ID:                 taskID1Trigger,
					OverallStatus:      statuses.OverallNew,
					Regularity:         tasking.Trigger,
					CanBeCanceled:      false,
					LastTaskInstanceID: tiID1Trigger,
				},
			},
			ctx: goodCtx,
		},
		{
			name: "testCase 5 - ok",
			mock: func(ctl *gomock.Controller) (Repo, InstancesRepo, logger.Logger) {
				tr := NewMockRepo(ctl)
				ti := NewMockInstancesRepo(ctl)

				tr.EXPECT().GetScheduledTasks(partnerID).Return([]e.ScheduledTasks{
					{ID: taskID1, LastTaskInstanceID: tiID1, CreatedAt: now, Regularity: tasking.OneTime},
					{ID: taskID2, CreatedAt: tomorrow}, // new task w\o ti
				}, nil)
				ti.EXPECT().GetInstancesForScheduled([]string{tiID1}).Return([]e.TaskInstance{{
					ID:       tiID1,
					TaskID:   taskID1,
					Statuses: map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceFailed},
				}}, nil)
				return tr, ti, nil
			},
			expectedData: []e.ScheduledTasks{
				{
					ID:            taskID2,
					CreatedAt:     tomorrow,
					CanBeCanceled: true,
					OverallStatus: statuses.OverallNew,
				},
				{
					ID:        taskID1,
					CreatedAt: now,
					ExecutionInfo: e.ExecutionInfo{
						FailedCount: 1,
						DeviceCount: 1,
					},
					OverallStatus:      statuses.OverallFailed,
					LastTaskInstanceID: tiID1,
					CanBeCanceled:      true,
					Regularity:         tasking.OneTime,
				},
			},
			ctx: goodCtx,
		},
		{
			name: "testCase 6 - ok with trigger tasks",
			mock: func(ctl *gomock.Controller) (Repo, InstancesRepo, logger.Logger) {
				tr := NewMockRepo(ctl)
				ti := NewMockInstancesRepo(ctl)

				tr.EXPECT().GetScheduledTasks(partnerID).Return([]e.ScheduledTasks{
					{ID: taskID1, LastTaskInstanceID: tiID1, CreatedAt: now, Regularity: tasking.OneTime},
					{ID: taskID1Trigger, LastTaskInstanceID: tiID1Trigger, CreatedAt: now,
						TriggerTypes: []string{"logOn"},
						Regularity:   tasking.Trigger,
						TriggerFrames: []tasking.TriggerFrame{{
							TriggerType: "logOn",
						}}},
					{ID: taskID2Trigger, CreatedAt: tomorrow, Regularity: tasking.Trigger,
						TriggerTypes: []string{"logOn"},
						TriggerFrames: []tasking.TriggerFrame{{
							TriggerType: "logOn",
						}}},
				}, nil)
				ti.EXPECT().GetInstancesForScheduled([]string{tiID1, tiID1Trigger}).Return([]e.TaskInstance{{
					ID:       tiID1,
					TaskID:   taskID1,
					Statuses: map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceFailed, endpID2: statuses.TaskInstancePending},
				}}, nil)
				ti.EXPECT().GetTopInstancesForScheduledByTaskIDs(gomock.Any()).Return([]e.TaskInstance{
					{
						ID:          tiID2Trigger,
						TaskID:      taskID1Trigger,
						LastRunTime: now,
						Statuses:    map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceFailed},
					}, {
						ID:       tiID1Trigger,
						TaskID:   taskID1Trigger,
						Statuses: map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceScheduled},
					},
				}, nil)
				return tr, ti, nil
			},
			expectedData: []e.ScheduledTasks{
				{
					ID:            taskID2Trigger,
					CreatedAt:     tomorrow,
					Regularity:    tasking.Trigger,
					OverallStatus: statuses.OverallNew,
					CanBeCanceled: false,
				},
				{
					ID:         taskID1,
					CreatedAt:  now,
					Regularity: tasking.OneTime,
					ExecutionInfo: e.ExecutionInfo{
						FailedCount: 1,
						DeviceCount: 2,
					},
					OverallStatus:      statuses.OverallRunning,
					LastTaskInstanceID: tiID1,
					CanBeCanceled:      false,
				},
				{
					ID:         taskID1Trigger,
					CreatedAt:  now,
					Regularity: tasking.Trigger,
					ExecutionInfo: e.ExecutionInfo{
						FailedCount: 1,
						DeviceCount: 1,
					},
					OverallStatus:      statuses.OverallFailed,
					LastTaskInstanceID: tiID2Trigger,
					LastRunTime:        now,
					CanBeCanceled:      false,
				},
			},
			ctx: goodCtx,
		},
		{
			name: "testCase 6.1 - ok with trigger tasks + postponed one",
			mock: func(ctl *gomock.Controller) (Repo, InstancesRepo, logger.Logger) {
				tr := NewMockRepo(ctl)
				ti := NewMockInstancesRepo(ctl)

				tr.EXPECT().GetScheduledTasks(partnerID).Return([]e.ScheduledTasks{
					{ID: taskID1, LastTaskInstanceID: timeUUID.String(),
						PostponedTime: tomorrow,
						CreatedAt:     now, Regularity: tasking.OneTime},
					{ID: taskID1Trigger, LastTaskInstanceID: tiID1Trigger, CreatedAt: now,
						TriggerTypes: []string{"logOn"},
						Regularity:   tasking.Trigger,
						TriggerFrames: []tasking.TriggerFrame{{
							TriggerType: "logOn",
						}}},
					{ID: taskID2Trigger, CreatedAt: tomorrow, Regularity: tasking.Trigger,
						TriggerTypes: []string{"logOn"},
						TriggerFrames: []tasking.TriggerFrame{{
							TriggerType: "logOn",
						}}},
				}, nil)
				ti.EXPECT().GetInstancesForScheduled([]string{timeUUID.String(), tiID1}).Return([]e.TaskInstance{{
					ID:       timeUUID.String(),
					TaskID:   taskID1,
					Statuses: map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceFailed, endpID2: statuses.TaskInstancePostponed},
				}}, nil).AnyTimes()
				ti.EXPECT().GetInstancesForScheduled([]string{tiID1, timeUUID.String()}).Return([]e.TaskInstance{{
					ID:       timeUUID.String(),
					TaskID:   taskID1,
					Statuses: map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceFailed, endpID2: statuses.TaskInstancePostponed},
				}}, nil).AnyTimes()
				ti.EXPECT().GetTopInstancesForScheduledByTaskIDs(gomock.Any()).Return([]e.TaskInstance{
					{
						ID:          tiID2Trigger,
						TaskID:      taskID1Trigger,
						LastRunTime: now,
						Statuses:    map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceFailed},
					}, {
						ID:       tiID1Trigger,
						TaskID:   taskID1Trigger,
						Statuses: map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceScheduled},
					},
				}, nil)
				return tr, ti, nil
			},
			expectedData: []e.ScheduledTasks{
				{
					ID:            taskID2Trigger,
					CreatedAt:     tomorrow,
					Regularity:    tasking.Trigger,
					OverallStatus: statuses.OverallNew,
					CanBeCanceled: false,
				},
				{
					ID:         taskID1,
					CreatedAt:  now,
					Regularity: tasking.OneTime,
					ExecutionInfo: e.ExecutionInfo{
						FailedCount: 1,
						DeviceCount: 2,
					},
					OverallStatus:      statuses.OverallPartialFailed,
					LastTaskInstanceID: timeUUID.String(),
					CanBeCanceled:      false,
				},
				{
					ID:         taskID1Trigger,
					CreatedAt:  now,
					Regularity: tasking.Trigger,
					ExecutionInfo: e.ExecutionInfo{
						FailedCount: 1,
						DeviceCount: 1,
					},
					OverallStatus:      statuses.OverallFailed,
					LastTaskInstanceID: tiID2Trigger,
					LastRunTime:        now,
					CanBeCanceled:      false,
				},
			},
			ctx: goodCtx,
		},
		{
			name: "testCase 7 - ok with recurrent+trigger tasks",
			mock: func(ctl *gomock.Controller) (Repo, InstancesRepo, logger.Logger) {
				tr := NewMockRepo(ctl)
				ti := NewMockInstancesRepo(ctl)

				tr.EXPECT().GetScheduledTasks(partnerID).Return([]e.ScheduledTasks{
					{ID: taskID1Trigger, LastTaskInstanceID: tiID1Trigger, CreatedAt: now,
						Regularity:   tasking.Recurrent,
						TriggerTypes: []string{"logOn"},
						TriggerFrames: []tasking.TriggerFrame{{
							TriggerType: "logOn",
						}}},
					{ID: taskID2Trigger, CreatedAt: tomorrow, Regularity: tasking.Recurrent,
						TriggerTypes: []string{"logOn"},
						TriggerFrames: []tasking.TriggerFrame{{
							TriggerType: "logOn",
						}}},
				}, nil)
				ti.EXPECT().GetInstancesForScheduled([]string{tiID1Trigger}).Return([]e.TaskInstance{
					{
						ID:       tiID1Trigger,
						TaskID:   taskID1Trigger,
						Statuses: map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceScheduled},
					},
				}, nil)
				ti.EXPECT().GetTopInstancesForScheduledByTaskIDs(gomock.Any()).Return([]e.TaskInstance{
					{
						ID:          tiID2Trigger,
						TaskID:      taskID1Trigger,
						LastRunTime: now,
						Statuses:    map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceFailed},
					}, {
						ID:       tiID1Trigger,
						TaskID:   taskID1Trigger,
						Statuses: map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceScheduled},
					}, {
						ID:       tiID1Trigger,
						TaskID:   taskID1Trigger,
						Statuses: map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceScheduled},
					},
				}, nil)
				return tr, ti, nil
			},
			expectedData: []e.ScheduledTasks{
				{
					ID:            taskID2Trigger,
					CreatedAt:     tomorrow,
					Regularity:    tasking.Recurrent,
					OverallStatus: statuses.OverallNew,
				},
				{
					ID:         taskID1Trigger,
					CreatedAt:  now,
					Regularity: tasking.Recurrent,
					ExecutionInfo: e.ExecutionInfo{
						FailedCount: 1,
						DeviceCount: 1,
					},
					OverallStatus:      statuses.OverallFailed,
					LastTaskInstanceID: tiID2Trigger,
					LastRunTime:        now,
					CanBeCanceled:      false,
				},
			},
			ctx: goodCtx,
		},
		{
			name: "testCase 8 - ok with recurrent+trigger tasks with recurrent last instance",
			mock: func(ctl *gomock.Controller) (Repo, InstancesRepo, logger.Logger) {
				tr := NewMockRepo(ctl)
				ti := NewMockInstancesRepo(ctl)

				tr.EXPECT().GetScheduledTasks(partnerID).Return([]e.ScheduledTasks{
					{ID: taskID1Trigger, LastTaskInstanceID: tiID1Trigger, CreatedAt: now, Regularity: tasking.Recurrent,
						TriggerTypes: []string{"logOn"},
						TriggerFrames: []tasking.TriggerFrame{{
							TriggerType: "logOn",
						}}},
				}, nil)
				ti.EXPECT().GetInstancesForScheduled([]string{tiID1Trigger}).Return([]e.TaskInstance{
					{
						ID:          tiID1Trigger,
						TaskID:      taskID1Trigger,
						LastRunTime: now,
						Statuses:    map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceFailed},
					},
				}, nil)
				ti.EXPECT().GetTopInstancesForScheduledByTaskIDs([]string{taskID1Trigger}).Return([]e.TaskInstance{
					{
						ID:          tiID2Trigger,
						TaskID:      taskID1Trigger,
						LastRunTime: now.Add(-2 * time.Hour),
						Statuses:    map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceFailed},
					},
					{
						ID:       tiID3Trigger,
						TaskID:   taskID1Trigger,
						Statuses: map[string]statuses.TaskInstanceStatus{endpID: statuses.TaskInstanceScheduled},
					},
				}, nil)
				return tr, ti, nil
			},
			expectedData: []e.ScheduledTasks{
				{
					ID:         taskID1Trigger,
					CreatedAt:  now,
					Regularity: tasking.Recurrent,
					ExecutionInfo: e.ExecutionInfo{
						FailedCount: 1,
						DeviceCount: 1,
					},
					OverallStatus:      statuses.OverallFailed,
					LastTaskInstanceID: tiID1Trigger,
					LastRunTime:        now,
					CanBeCanceled:      false,
				},
			},
			ctx: goodCtx,
		},
	}

	for _, tc := range testCases {
		ctl := gomock.NewController(t)

		tr, tiRepo, _ := tc.mock(ctl)
		uc := NewTasks(tr, nil, tiRepo, nil, nil, nil, logger.Log)

		gotData, gotErr := uc.GetScheduledTasks(tc.ctx)
		ctl.Finish()

		if tc.expectedError {
			Ω(gotErr).ShouldNot(BeNil(), tc.name)
		} else {
			Ω(gotErr).Should(BeNil(), tc.name)
			Ω(gotData).Should(ConsistOf(tc.expectedData), tc.name)
		}
	}
}

func TestDeleteScheduledTasks(t *testing.T) {
	RegisterTestingT(t)
	var mockCtrl *gomock.Controller

	now := time.Now().UTC()
	uuid := gocql.UUIDFromTime(now)
	ctx := context.WithValue(context.Background(), config.PartnerIDKeyCTX, partnerID)

	type fields struct {
		legacyRepo func() LegacyRepo
		usecase    func() trigger.Usecase
	}
	type args struct {
		ctx context.Context
		ids e.TaskIDs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful",
			fields: fields{
				legacyRepo: func() LegacyRepo {
					r := NewMockLegacyRepo(mockCtrl)

					r.EXPECT().GetByIDs(gomock.Any(), nil, partnerID, false, uuid).
						Return([]models.Task{{ID: uuid, State: statuses.TaskStateActive, RunTimeUTC: now.UTC()}}, nil).
						Times(1)

					r.EXPECT().UpdateSchedulerFields(gomock.Any(), models.Task{ID: uuid, State: statuses.TaskStateInactive, RunTimeUTC: time.Now().Truncate(time.Minute)}).
						Return(nil).
						Times(1)

					return r
				},
			},
			args: args{
				ctx: ctx,
				ids: e.TaskIDs{IDs: []string{uuid.String()}},
			},
		},
		{
			name: "failed to load partnerID from context",
			fields: fields{
				legacyRepo: func() LegacyRepo {
					r := NewMockLegacyRepo(mockCtrl)
					return r
				},
			},
			args: args{
				ctx: context.Background(),
				ids: e.TaskIDs{IDs: []string{uuid.String()}},
			},
			wantErr: true,
			errMsg:  "can't get parameter from context: partnerID",
		},
		{
			name: "failed to map uuids",
			fields: fields{
				legacyRepo: func() LegacyRepo {
					r := NewMockLegacyRepo(mockCtrl)
					return r
				},
			},
			args: args{
				ctx: ctx,
				ids: e.TaskIDs{IDs: []string{"111"}},
			},
			wantErr: true,
			errMsg:  "uuid is invalid: 111",
		},
		{
			name: "failed while getting tasks",
			fields: fields{
				legacyRepo: func() LegacyRepo {
					r := NewMockLegacyRepo(mockCtrl)

					r.EXPECT().GetByIDs(gomock.Any(), nil, partnerID, false, uuid).
						Return([]models.Task{}, errors.New("fail")).
						Times(1)

					return r
				},
			},
			args: args{
				ctx: ctx,
				ids: e.TaskIDs{IDs: []string{uuid.String()}},
			},
			wantErr: true,
			errMsg:  "fail",
		},
		{
			name: "returned zero tasks",
			fields: fields{
				legacyRepo: func() LegacyRepo {
					r := NewMockLegacyRepo(mockCtrl)
					r.EXPECT().GetByIDs(gomock.Any(), nil, partnerID, false, uuid).
						Return([]models.Task{}, nil)
					return r
				},
			},
			args: args{
				ctx: ctx,
				ids: e.TaskIDs{IDs: []string{uuid.String()}},
			},
			errMsg:  "got zero tasks",
			wantErr: true,
		},
		{
			name: "failed while updating tasks",
			fields: fields{
				legacyRepo: func() LegacyRepo {
					r := NewMockLegacyRepo(mockCtrl)

					r.EXPECT().GetByIDs(gomock.Any(), nil, partnerID, false, uuid).
						Return([]models.Task{{ID: uuid, State: statuses.TaskStateActive, RunTimeUTC: now}}, nil).
						Times(1)

					r.EXPECT().UpdateSchedulerFields(gomock.Any(), models.Task{ID: uuid, State: statuses.TaskStateInactive, RunTimeUTC: now.Truncate(time.Minute).Local()}).
						Return(errors.New("fail")).
						Times(1)

					return r
				},
			},
			args: args{
				ctx: ctx,
				ids: e.TaskIDs{IDs: []string{uuid.String()}},
			},
			wantErr: true,
			errMsg:  "fail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl = gomock.NewController(t)
			var tasks *Tasks
			if tt.fields.usecase != nil {
				tasks = &Tasks{
					tasksRepo:        nil,
					legacyRepo:       tt.fields.legacyRepo(),
					taskInstanceRepo: nil,
					execResultsRepo:  nil,
					tr:               tt.fields.usecase(),
					log:              nil,
				}
			} else {
				tasks = &Tasks{
					tasksRepo:        nil,
					legacyRepo:       tt.fields.legacyRepo(),
					taskInstanceRepo: nil,
					execResultsRepo:  nil,
					tr:               nil,
					log:              nil,
				}
			}

			err := tasks.DeleteScheduledTasks(tt.args.ctx, tt.args.ids)

			mockCtrl.Finish()

			if tt.wantErr {
				Ω(errors.Cause(err)).Should(MatchError(tt.errMsg), tt.name)
			} else {
				Ω(err).To(BeNil(), tt.name)
			}
		})
	}
}
