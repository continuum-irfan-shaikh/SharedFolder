package tasks

import (
	"context"
	"testing"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	e "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

func TestTasks_GetTasksHistory(t *testing.T) {
	RegisterTestingT(t)

	ctx := context.WithValue(context.Background(), "s", "s")
	ctxWithUser := context.WithValue(ctx, config.UserKeyCTX, e.User{})
	goodCtx := context.WithValue(ctxWithUser, config.PartnerIDKeyCTX, partnerID)

	now := time.Now().UTC()
	tomorrow := now.Add(time.Hour * (24))
	uuidID, _ := gocql.ParseUUID("11111111-2222-1111-1111-111111111111")
	taskIDs := []gocql.UUID{uuidID}

	instances := []e.TaskInstance{{
		ID:            "11111111-2222-1111-1111-111111111112",
		TaskID:        "11111111-2222-1111-1111-111111111111",
		StartedAt:     time.Time{},
		LastRunTime:   time.Time{},
		Statuses:      nil,
		TaskName:      "hi",
		PartnerID:     partnerID,
		FailureCount:  0,
		SuccessCount:  0,
		StatusesCount: nil,
	},
	}

	tasks := []models.Task{{
		ID:                 uuidID,
		Name:               "hi",
		CreatedAt:          time.Time{},
		PartnerID:          partnerID,
		OriginID:           gocql.UUID{},
		State:              0,
		RunTimeUTC:         time.Time{},
		ExternalTask:       false,
		LastTaskInstanceID: gocql.UUID{},
		ModifiedAt:         time.Time{},
	},
	}

	testCases := []struct {
		name          string
		expectedError bool
		expectedData  []models.Task
		mock          func(ctl *gomock.Controller) (Repo, InstancesRepo, LegacyRepo, logger.Logger)
		ctx           context.Context
		from          time.Time
		to            time.Time
	}{
		{
			name:          "testCase 1 - can't get user from context",
			expectedError: true,
			mock: func(ctl *gomock.Controller) (repo Repo, ti InstancesRepo, legacy LegacyRepo, log logger.Logger) {
				return nil, nil, nil, nil
			},
			ctx:  ctx,
			from: now,
			to:   tomorrow,
		},
		{
			name:          "testCase 2 - can't get partner from context",
			expectedError: true,
			mock: func(ctl *gomock.Controller) (repo Repo, ti InstancesRepo, legacy LegacyRepo, log logger.Logger) {
				return nil, nil, nil, nil
			},
			ctx: ctxWithUser,
		},
		{
			name:          "testCase 3 - can't get task instance",
			expectedError: true,
			mock: func(ctl *gomock.Controller) (Repo, InstancesRepo, LegacyRepo, logger.Logger) {
				tr := NewMockInstancesRepo(ctl)
				tr.EXPECT().GetByStartedAtAfter(partnerID, time.Time{}, time.Time{}).Return([]e.TaskInstance{}, errors.New("err"))
				return nil, tr, nil, nil
			},
			ctx: goodCtx,
		},
		{
			name:          "testCase 4 - can't get task",
			expectedError: true,
			mock: func(ctl *gomock.Controller) (Repo, InstancesRepo, LegacyRepo, logger.Logger) {
				ti := NewMockInstancesRepo(ctl)
				tr := NewMockLegacyRepo(ctl)
				ti.EXPECT().GetByStartedAtAfter(partnerID, time.Time{}, time.Time{}).Return(instances, nil)
				tr.EXPECT().GetByIDs(goodCtx, nil, partnerID, true, taskIDs).Return([]models.Task{}, errors.New("err"))
				return nil, ti, tr, nil
			},
			ctx: goodCtx,
		},
		{
			name:          "testCase 5 - can get task",
			expectedError: false,
			mock: func(ctl *gomock.Controller) (Repo, InstancesRepo, LegacyRepo, logger.Logger) {
				ti := NewMockInstancesRepo(ctl)
				tr := NewMockLegacyRepo(ctl)
				ti.EXPECT().GetByStartedAtAfter(partnerID, time.Time{}, time.Time{}).Return(instances, nil)
				tr.EXPECT().GetByIDs(goodCtx, nil, partnerID, true, taskIDs).Return(tasks, nil)
				return nil, ti, tr, nil
			},
			ctx:          goodCtx,
			expectedData: tasks,
		},
	}

	for _, tc := range testCases {
		ctl := gomock.NewController(t)

		tr, tiRepo, legacyRepo, log := tc.mock(ctl)
		uc := NewTasks(tr, legacyRepo, tiRepo, nil, nil, nil, log)

		_, gotErr := uc.GetTasksHistory(tc.ctx, tc.from, tc.to)
		ctl.Finish()

		if tc.expectedError {
			Ω(gotErr).ShouldNot(BeNil(), tc.name)
		} else {
			Ω(gotErr).Should(BeNil(), tc.name)
		}
	}
}
