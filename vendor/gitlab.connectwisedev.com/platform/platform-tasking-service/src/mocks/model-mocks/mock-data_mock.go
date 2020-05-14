package modelMocks

import (
	"encoding/json"
	"fmt"
	"time"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"github.com/gocql/gocql"
)

var (
	// ExistedManagedEndpointIDStr represents ManagedEndpointID for the first Result in the TaskExecutionResultsView list
	ExistedManagedEndpointIDStr = "58a1af2f-6579-4aec-b45d-5dfde879ef01"
	// ExistedManagedEndpointID represents ManagedEndpointID for the first Result in the TaskExecutionResultsView list
	ExistedManagedEndpointID = str2uuid(ExistedManagedEndpointIDStr)
	// ExistedTaskID represents TaskInstanceID for the first Result in TaskExecutionResultsView
	ExistedTaskID = str2uuid("58a1af2f-6579-4aec-b45d-5dfde879ef01")
	// ExistedTaskInstanceID represents TaskInstanceID for the first Result in TaskExecutionResultsView
	ExistedTaskInstanceID = str2uuid("58a1af2f-6579-4aec-b45d-000000000001")
	// TargetIDStr is a target uuid used for testing
	TargetIDStr = "11111111-1111-1111-1111-111111111111"
	//EmptyUUID represents nullable uuid
	EmptyUUID = str2uuid("00000000-0000-0000-0000-000000000000")
	// TargetID is a target uuid used for testing
	TargetID = str2uuid(TargetIDStr)
	// NotExistedManagedEndpointID is an empty (not existed) ManagedEndpointID
	NotExistedManagedEndpointID gocql.UUID
	// NotExistedTaskID is an empty (not existed) TaskID
	NotExistedTaskID gocql.UUID
	// NotExistedTaskInstanceID is an empty (not existed) TaskInstanceID
	NotExistedTaskInstanceID gocql.UUID

	// ExistedResult is a first Result in the TaskExecutionResultsView list
	ExistedResult models.ExecutionResult

	// NewManagedEndpointID represents ManagedEndpointID for a new Result
	NewManagedEndpointID = str2uuid("4e1216b3-4b66-4a99-9d3f-1a2b93763e88")
	// NewTaskInstanceID represents TaskInstanceID for a new Result
	NewTaskInstanceID = str2uuid("9af406bd-bab0-4936-8b72-ba738495d8d2")
	// NewCurrentTime represents a fixed time not used in predefined results
	NewCurrentTime = time.Now().UTC()

	// TaskExecutionResultsView stores a Mock data as a slice of predefined TaskExecutionResultsView
	TaskExecutionResultsView = []models.ExecutionResultView{}
)

func init() {
	executionResultsViewByte, err := json.Marshal(mockExecutionResultsView)
	if err != nil {
		panic(fmt.Sprintf("can not Marshal mockExecutionResultsView, err: %v\n", err))
	}
	populateTaskExecutionResultsView(&TaskExecutionResultsView, executionResultsViewByte)
	executionResultView := TaskExecutionResultsView[0]
	ExistedResult = models.ExecutionResult{
		ManagedEndpointID: executionResultView.ManagedEndpointID,
		TaskInstanceID:    executionResultView.ExecutionID,
		UpdatedAt:         executionResultView.LastRunTime,
		StdOut:            executionResultView.LastRunStdOut,
		StdErr:            executionResultView.LastRunStdErr,
		ExecutionStatus:   executionResultView.LastRunStatus,
	}
}

// GenerateTargetsMock is used to return array of task targets
func GenerateTargetsMock() map[string]bool {
	return map[string]bool{TargetIDStr: true}
}

// GenerateDisabledTargetsMock is used to return array of disabled task targets
func GenerateDisabledTargetsMock() map[string]bool {
	return map[string]bool{TargetIDStr: false}
}

func populateTaskExecutionResultsView(executionResultsView *[]models.ExecutionResultView, fixture []byte) {
	if err := json.Unmarshal(fixture, executionResultsView); err != nil {
		message := fmt.Sprintf("can not populate ExecutionResultViewRepoMock with data (err: %v)\n", err)
		panic(message)
	}
}

var (
	partnerID                = "50016364"
	mockExecutionResultsView = []models.ExecutionResultView{
		{
			PartnerID:         partnerID,
			ManagedEndpointID: ExistedManagedEndpointID,
			TaskID:            ExistedTaskID,
			ExecutionID:       str2uuid("58a1af2f-6579-4aec-b45d-000000000001"),
			OriginID:          str2uuid("00000000-0000-0000-1111-000000000000"),
			TaskName:          "Wake Up Device",
			Type:              config.ScriptTaskType,
			Description:       "This task wakes up the device.",
			Regularity:        apiModels.Recurrent,
			InitiatedBy:       "Andy",
			Status:            statuses.TaskStateInactive,
			LastRunTime:       someTime,
			LastRunStatus:     statuses.TaskInstanceSuccess,
			LastRunStdOut:     "Done",
			LastRunStdErr:     "",
		},
		{
			PartnerID:         partnerID,
			ManagedEndpointID: str2uuid("58a1af2f-6579-4aec-b45d-5dfde879ef99"),
			TaskID:            ExistedTaskID,
			ExecutionID:       str2uuid("58a1af2f-6579-4aec-b45d-000000000002"),
			OriginID:          str2uuid("00000000-0000-0000-2222-000000000000"),
			TaskName:          "Install Adobe Acrobat Reader",
			Type:              config.ScriptTaskType,
			Description:       "This task installs Adobe Acrobat Reader on the device",
			Regularity:        apiModels.Recurrent,
			InitiatedBy:       "Andy",
			Status:            statuses.TaskStateInactive,
			LastRunTime:       someTime,
			LastRunStatus:     statuses.TaskInstanceSuccess,
			LastRunStdOut:     "Adobe Acrobat Reader has been installed successfully",
			LastRunStdErr:     "",
		},
		{
			PartnerID:         partnerID,
			ManagedEndpointID: str2uuid("58a1af2f-6579-4aec-b45d-5dfde879ef99"),
			TaskID:            ExistedTaskID,
			ExecutionID:       str2uuid("58a1af2f-6579-4aec-b45d-000000000003"),
			OriginID:          str2uuid("1d33d253-2a9b-47bc-ba3c-44a9940d8063"),
			TaskName:          "Delete Temp Files",
			Type:              config.ScriptTaskType,
			Description:       "This task deleted the temporary files on the device.",
			Regularity:        apiModels.Recurrent,
			InitiatedBy:       "Andy",
			Status:            statuses.TaskStateInactive,
			LastRunTime:       someTime,
			LastRunStatus:     statuses.TaskInstanceSuccess,
			LastRunStdOut:     "507 files deleted successfully",
			LastRunStdErr:     "",
		},
		{
			PartnerID:         partnerID,
			ManagedEndpointID: str2uuid("58a1af2f-6579-4aec-b45d-5dfde879ef99"),
			TaskID:            ExistedTaskID,
			ExecutionID:       str2uuid("58a1af2f-6579-4aec-b45d-000000000004"),
			OriginID:          str2uuid("9ad73107-6fe3-462e-8e97-b1639516f3f7"),
			TaskName:          "Add Local User",
			Type:              config.ScriptTaskType,
			Description:       "This task adds a local user account to the device.",
			Regularity:        apiModels.Recurrent,
			InitiatedBy:       "Andy",
			Status:            statuses.TaskStateInactive,
			LastRunTime:       someTime,
			LastRunStatus:     statuses.TaskInstanceFailed,
			LastRunStdOut:     "",
			LastRunStdErr:     "Can not finish operation: account already exists.",
		},
		{
			PartnerID:         partnerID,
			ManagedEndpointID: str2uuid("74a86e19-2d9b-413a-a820-cf017e22026b"),
			TaskID:            str2uuid("58a1af2f-6579-4aec-b45d-5dfde879ef01"),
			ExecutionID:       str2uuid("74a86e19-2d9b-413a-a820-000000000001"),
			OriginID:          str2uuid("00000000-0000-0000-1111-000000000000"),
			TaskName:          "Wake Up Device",
			Type:              config.ScriptTaskType,
			Description:       "This task wakes up the device.",
			Regularity:        apiModels.Recurrent,
			InitiatedBy:       "Andy",
			Status:            statuses.TaskStateInactive,
			LastRunTime:       someTime,
			LastRunStatus:     statuses.TaskInstanceSuccess,
			LastRunStdOut:     "Done",
			LastRunStdErr:     "",
		},
		{
			PartnerID:         partnerID,
			ManagedEndpointID: str2uuid("74a86e19-2d9b-413a-a820-cf017e22026b"),
			TaskID:            str2uuid("58a1af2f-6579-4aec-b45d-5dfde879ef02"),
			ExecutionID:       str2uuid("74a86e19-2d9b-413a-a820-000000000002"),
			OriginID:          str2uuid("00000000-0000-0000-2222-000000000000"),
			TaskName:          "Install Adobe Acrobat Reader",
			Type:              config.ScriptTaskType,
			Description:       "This task installs Adobe Acrobat Reader on the device",
			Regularity:        apiModels.Recurrent,
			InitiatedBy:       "Andy",
			Status:            statuses.TaskStateInactive,
			LastRunTime:       someTime,
			LastRunStatus:     statuses.TaskInstanceSuccess,
			LastRunStdOut:     "Adobe Acrobat Reader has been installed successfully",
			LastRunStdErr:     "",
		},
	}
)
