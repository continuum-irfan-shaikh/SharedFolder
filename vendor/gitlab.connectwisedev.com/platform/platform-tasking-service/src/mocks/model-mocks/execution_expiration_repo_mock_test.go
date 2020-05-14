package modelMocks

import (
	"context"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"github.com/gocql/gocql"
	"reflect"
	"testing"
	"time"
)

var executionExpirationRepoMock ExecutionExpirationRepoMock

func setUpRepoMock(fillMock bool) {
	executionExpirationRepoMock = NewExecutionExpirationRepoMock(fillMock)
}

func clearUpRepoMock() {
	executionExpirationRepoMock = ExecutionExpirationRepoMock{}
}

func TestNewExecutionExpirationRepoMock(t *testing.T) {
	test := NewExecutionExpirationRepoMock(true)
	if test.repo[CtrlExpirationTime.Unix()][0].TaskInstanceID != CtrlTaskInstanceID {
		t.Fatalf("Something went wrong")
	}
}

func TestInsertExecutionExpiration(t *testing.T) {
	setUpRepoMock(true)
	defer clearUpRepoMock()

	expTime := time.Now().Truncate(time.Minute)
	randomUUID, _ := gocql.RandomUUID()

	testCases := []struct {
		name        string
		exp         models.ExecutionExpiration
		isNeedError bool
	}{
		{
			name: "testCase1 - success",
			exp: models.ExecutionExpiration{
				ExpirationTimeUTC: expTime,
				PartnerID:         partnerID,
				TaskInstanceID:    randomUUID,
			},
			isNeedError: false,
		},
		{
			name: "testCase2 - error",
			exp: models.ExecutionExpiration{
				ExpirationTimeUTC: expTime,
				PartnerID:         partnerID,
				TaskInstanceID:    randomUUID,
			},
			isNeedError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := executionExpirationRepoMock.InsertExecutionExpiration(context.WithValue(context.TODO(), IsNeedError, tt.isNeedError), tt.exp)
			if err != nil && !tt.isNeedError {
				t.Fatalf("%s: expected no error, but got err: %v", tt.name, err)
			}

			if err != nil {
				return
			}

			ctrlExps := []models.ExecutionExpiration{
				{
					ExpirationTimeUTC: expTime,
					PartnerID:         partnerID,
					TaskInstanceID:    randomUUID,
				},
			}
			repoExps := executionExpirationRepoMock.GetMockRepo()

			if ok := reflect.DeepEqual(repoExps[expTime.Unix()], ctrlExps); !ok {
				t.Fatalf("%s: expected %v, got %v", tt.name, ctrlExps, repoExps[expTime.Unix()])
			}
		})
	}
}

func TestGetByExpirationTime(t *testing.T) {
	setUpRepoMock(true)
	defer clearUpRepoMock()

	testCases := []struct {
		name           string
		expirationTime time.Time
		isNeedError    bool
	}{
		{
			name:           "testCase1 - success",
			expirationTime: CtrlExpirationTime,
			isNeedError:    false,
		},
		{
			name:           "testCase2 - error",
			expirationTime: CtrlExpirationTime,
			isNeedError:    true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			exps, err := executionExpirationRepoMock.GetByExpirationTime(context.WithValue(context.TODO(), IsNeedError, tt.isNeedError), tt.expirationTime)
			if err != nil && !tt.isNeedError {
				t.Fatalf("%s: expected no error, but got err: %v", tt.name, err)
			}
			if err != nil {
				return
			}

			repoExps := executionExpirationRepoMock.GetMockRepo()
			ok := reflect.DeepEqual(exps, repoExps[tt.expirationTime.Unix()])
			if !ok {
				t.Fatalf("%s: expected %v, got %v", tt.name, repoExps[tt.expirationTime.Unix()], exps)
			}
		})
	}
}
