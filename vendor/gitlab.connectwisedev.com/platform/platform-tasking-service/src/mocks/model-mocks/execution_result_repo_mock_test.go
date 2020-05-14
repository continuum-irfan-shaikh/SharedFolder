package modelMocks

import (
	"context"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"github.com/gocql/gocql"
	"reflect"
	"testing"
	"time"
)

func TestScriptExecutionResultRepoMockUpsirtErrInRepo(t *testing.T) {
	repo := NewExecutionResultRepoMock(true, true)
	err := repo.Upsert(context.TODO(), "", "", models.ExecutionResult{})
	if err == nil {
		t.Errorf("Upsirt :Error expected but none has been gotten from the Repository")
	}
}

func TestScriptExecutionResultRepoMockGetByTaskInstanceID(t *testing.T) {
	for _, tc := range scriptExecutionResultGetByTaskInstanceIDTestCases {
		testCase := tc
		t.Run(testCase.Name, func(t *testing.T) {
			repo := NewExecutionResultRepoMock(testCase.IsFilled, testCase.RepoError)
			actualResult, err := repo.GetByTargetAndTaskInstanceIDs(testCase.ManagedEndpointID, testCase.TaskInstanceID)

			if err == nil && (testCase.RepoError || testCase.WantError) {
				t.Fatalf("\nGetByTaskInstanceID [%s]:\nError expected but none has been gotten from the Repository", testCase.Name)
			}

			if testCase.ResultMustExist && !reflect.DeepEqual(actualResult[0], ExistedResult) {
				t.Fatalf("\nGetByTaskInstanceID [%s]:\nError getting the Result: want %v but got %v", testCase.Name, ExistedResult, actualResult)
			}

		})

	}
}

func TestScriptExecutionResultRepoMockUpsert(t *testing.T) {
	for _, tc := range executionResultUpsertTestCases {
		testCase := tc
		t.Run(testCase.Name, func(t *testing.T) {
			repo := NewExecutionResultRepoMock(testCase.IsFilled, false)
			err := repo.Upsert(context.TODO(), "", "", testCase.Result)
			if testCase.WantError && err == nil {
				t.Fatalf("\nUpsert [%s]:\nError expected but none has been gotten from the Repository", testCase.Name)
			}
			if !testCase.WantError && err != nil {
				t.Fatalf("\nUpsert [%s]:\nGot error upserting result %v (err: %v)", testCase.Name, testCase.Result, err)
			}
			expectedResult := testCase.Result
			actualResult, err := repo.GetByTargetAndTaskInstanceIDs(testCase.Result.ManagedEndpointID, testCase.Result.TaskInstanceID)
			if err != nil {
				t.Fatalf("\nUpsert [%s]:\nGot error upserting result %v (err: %v)", testCase.Name, expectedResult, err)
			}
			if !reflect.DeepEqual(expectedResult, actualResult[0]) {
				t.Fatalf("\nUpsert [%s]:\nResult is not upserted: want %v but got %v", testCase.Name, expectedResult, actualResult)
			}
		})
	}
}

var scriptExecutionResultGetByTaskInstanceIDTestCases = []struct {
	Name              string
	ManagedEndpointID gocql.UUID
	TaskInstanceID    gocql.UUID
	IsFilled          bool
	RepoError         bool
	ResultMustExist   bool
	WantError         bool
	Message           string
}{
	{
		Name:              "Get an Existed Result from the Filled Repository w/o errors",
		ManagedEndpointID: ExistedManagedEndpointID,
		TaskInstanceID:    ExistedTaskInstanceID,
		IsFilled:          true,
		RepoError:         false,
		ResultMustExist:   true,
		WantError:         false,
	},
	{
		Name:              "Get no Result from the Filled Repository w/ errors",
		ManagedEndpointID: ExistedManagedEndpointID,
		TaskInstanceID:    ExistedTaskInstanceID,
		IsFilled:          true,
		RepoError:         true,
		ResultMustExist:   false,
		WantError:         true,
	},
	{
		Name:              "Get no Result from the Empty Repository w/o errors",
		ManagedEndpointID: ExistedManagedEndpointID,
		TaskInstanceID:    ExistedTaskInstanceID,
		IsFilled:          false,
		RepoError:         false,
		ResultMustExist:   false,
		WantError:         true,
	},
	{
		Name:              "Get no Result from the Empty Repository w/ errors",
		ManagedEndpointID: ExistedManagedEndpointID,
		TaskInstanceID:    ExistedTaskInstanceID,
		IsFilled:          false,
		RepoError:         true,
		ResultMustExist:   false,
		WantError:         true,
	},

	{
		Name:              "Get no Result (NotExistedManagedEndpointID) from the Filled Repository w/o errors",
		ManagedEndpointID: NotExistedManagedEndpointID,
		TaskInstanceID:    ExistedTaskInstanceID,
		IsFilled:          true,
		RepoError:         false,
		ResultMustExist:   false,
		WantError:         true,
	},
	{
		Name:              "Get no Result (NotExistedTaskInstanceID) from the Filled Repository w/o errors",
		ManagedEndpointID: ExistedManagedEndpointID,
		TaskInstanceID:    NotExistedTaskInstanceID,
		IsFilled:          true,
		RepoError:         false,
		ResultMustExist:   false,
		WantError:         true,
	},
	{
		Name:              "Get no Result (both IDs are not existed) from the Filled Repository w/o errors",
		ManagedEndpointID: NotExistedManagedEndpointID,
		TaskInstanceID:    NotExistedTaskInstanceID,
		IsFilled:          true,
		RepoError:         false,
		ResultMustExist:   false,
		WantError:         true,
	},
}

var executionResultUpsertTestCases = []struct {
	Name      string
	IsFilled  bool
	WantError bool
	Result    models.ExecutionResult
}{
	{
		Name:      "Repo: filled, ManagedEndpointID: new,  TaskInstanceID: new",
		IsFilled:  true,
		WantError: false,
		Result: models.ExecutionResult{
			ManagedEndpointID: NewManagedEndpointID,
			TaskInstanceID:    NewTaskInstanceID,
			UpdatedAt:         NewCurrentTime,
			StdOut:            "Done",
			StdErr:            "",
			ExecutionStatus:   statuses.TaskInstanceSuccess,
		},
	},
	{
		Name:      "Repo: filled, ManagedEndpointID: existed,  TaskInstanceID: new",
		IsFilled:  true,
		WantError: false,
		Result: models.ExecutionResult{
			ManagedEndpointID: ExistedManagedEndpointID,
			TaskInstanceID:    NewTaskInstanceID,
			UpdatedAt:         NewCurrentTime,
			StdOut:            "Done",
			StdErr:            "",
			ExecutionStatus:   statuses.TaskInstanceSuccess,
		},
	},
	{
		Name:      "Repo: filled, ManagedEndpointID: new,  TaskInstanceID: existed",
		IsFilled:  true,
		WantError: false,
		Result: models.ExecutionResult{
			ManagedEndpointID: NewManagedEndpointID,
			TaskInstanceID:    ExistedTaskInstanceID,
			UpdatedAt:         NewCurrentTime,
			StdOut:            "Done",
			StdErr:            "",
			ExecutionStatus:   statuses.TaskInstanceSuccess,
		},
	},
	{
		Name:      "Repo: filled, ManagedEndpointID: existed,  TaskInstanceID: existed",
		IsFilled:  true,
		WantError: false,
		Result: models.ExecutionResult{
			ManagedEndpointID: NewManagedEndpointID,
			TaskInstanceID:    NewTaskInstanceID,
			UpdatedAt:         NewCurrentTime,
			StdOut:            "Done",
			StdErr:            "",
			ExecutionStatus:   statuses.TaskInstanceSuccess,
		},
	},

	{
		Name:      "Repo: empty, ManagedEndpointID: new,  TaskInstanceID: new",
		IsFilled:  false,
		WantError: false,
		Result: models.ExecutionResult{
			ManagedEndpointID: NewManagedEndpointID,
			TaskInstanceID:    NewTaskInstanceID,
			UpdatedAt:         NewCurrentTime,
			StdOut:            "Done",
			StdErr:            "",
			ExecutionStatus:   statuses.TaskInstanceSuccess,
		},
	},
	{
		Name:      "Repo: empty, ManagedEndpointID: existed,  TaskInstanceID: new",
		IsFilled:  false,
		WantError: false,
		Result: models.ExecutionResult{
			ManagedEndpointID: ExistedManagedEndpointID,
			TaskInstanceID:    NewTaskInstanceID,
			UpdatedAt:         NewCurrentTime,
			StdOut:            "Done",
			StdErr:            "",
			ExecutionStatus:   statuses.TaskInstanceSuccess,
		},
	},
	{
		Name:      "Repo: empty, ManagedEndpointID: new,  TaskInstanceID: existed",
		IsFilled:  false,
		WantError: false,
		Result: models.ExecutionResult{
			ManagedEndpointID: NewManagedEndpointID,
			TaskInstanceID:    ExistedTaskInstanceID,
			UpdatedAt:         NewCurrentTime,
			StdOut:            "Done",
			StdErr:            "",
			ExecutionStatus:   statuses.TaskInstanceSuccess,
		},
	},
	{
		Name:      "Repo: empty, ManagedEndpointID: existed,  TaskInstanceID: existed",
		IsFilled:  false,
		WantError: false,
		Result: models.ExecutionResult{
			ManagedEndpointID: NewManagedEndpointID,
			TaskInstanceID:    NewTaskInstanceID,
			UpdatedAt:         NewCurrentTime,
			StdOut:            "Done",
			StdErr:            "",
			ExecutionStatus:   statuses.TaskInstanceSuccess,
		},
	},
}

func TestGetByTargetsAndTaskInstanceIDs(t *testing.T) {
	getTime := func(s string) (t time.Time) {
		t, _ = time.Parse(time.RFC3339, s)
		return t
	}

	tests := []struct {
		name            string
		needErr         bool
		needFill        bool
		taskInstanceIDs []gocql.UUID
		want            []models.ExecutionResult
	}{
		{
			name:     "Good",
			needErr:  false,
			needFill: true,
			taskInstanceIDs: []gocql.UUID{
				str2uuid("58a1af2f-6579-4aec-b45d-000000000001"),
			},
			want: []models.ExecutionResult{
				{
					ManagedEndpointID: str2uuid("58a1af2f-6579-4aec-b45d-5dfde879ef01"),
					TaskInstanceID:    str2uuid("58a1af2f-6579-4aec-b45d-000000000001"),
					UpdatedAt:         someTime,
					ExecutionStatus:   statuses.TaskInstanceSuccess,
					StdOut:            "Done",
				},
			},
		},
		{
			name:     "Bad",
			needErr:  true,
			needFill: true,
			taskInstanceIDs: []gocql.UUID{
				str2uuid("58a1af2f-6579-4aec-b45d-000000000001"),
			},
			want: []models.ExecutionResult{
				{
					ManagedEndpointID: str2uuid("58a1af2f-6579-4aec-b45d-5dfde879ef01"),
					TaskInstanceID:    str2uuid("58a1af2f-6579-4aec-b45d-000000000001"),
					UpdatedAt:         getTime("2017-10-11T10:04:05Z"),
					ExecutionStatus:   statuses.TaskInstanceSuccess,
					StdOut:            "Done",
				},
			},
		},
		{
			name:     "Non-existed exec result",
			needErr:  false,
			needFill: false,
			taskInstanceIDs: []gocql.UUID{
				gocql.TimeUUID(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewExecutionResultRepoMock(test.needFill, test.needErr)
			results, err := repo.GetByTaskInstanceIDs(test.taskInstanceIDs)
			if (err != nil) != test.needErr {
				t.Fatalf("Got err=%v, but want err=%v", err, test.needErr)
			}
			if !test.needFill && !reflect.DeepEqual(results, []models.ExecutionResult{}) {
				t.Fatalf("Wanted nil slice of Script Execution Results, but got %v", results)
			}
			if !test.needErr {
				if test.needFill && !reflect.DeepEqual(results, test.want) {
					t.Fatalf("Got %v\nwant %v", results, test.want)
				}
			}
		})
	}
}
