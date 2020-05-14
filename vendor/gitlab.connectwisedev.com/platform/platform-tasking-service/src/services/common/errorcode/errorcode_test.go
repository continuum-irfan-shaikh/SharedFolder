package errorcode

import (
	"testing"
)

func TestName(t *testing.T) {
	b := NewBadRequestErr("", "")
	c := b.Error()

	i := NewInternalServerErr("", "")
	c = i.Error()

	n := NewNotFoundErr("", "")
	c = n.Error()

	c = ErrorCantDecodeInputData
	c = ErrorTaskDefinitionExists
	c = ErrorCantGetTaskDefinitionTemplate
	c = ErrorCantGetTemplatesForExecutionMS
	c = ErrorWrongTemplateType
	c = ErrorCantGetTaskExecutionResultsForManagedEndpoint
	c = ErrorCantProcessTaskExecutionResults
	c = CodeUpdated
	c = CodeCreated
	c = ErrorCantGetPartnersForCounters
	c = ErrorCantCreateNewTask
	c = ErrorCantSaveTaskToDB
	c = ErrorCantUpdateTask
	c = ErrorCantUpdateTaskInstances
	c = ErrorCantGetListOfTasksByManagedEndpoint
	c = ErrorCantGetTaskByTaskID
	c = ErrorCantGetTriggerTypes
	c = ErrorCantProcessTriggerDefinitions
	c = ErrorTaskIsNotFoundByTaskID
	c = ErrorCantGetTaskByTaskInstanceIDs
	c = ErrorCantGetCountByPartnerID
	c = ErrorCantGetTaskInstanceCountByTaskID
	c = ErrorCantGetTaskExecutionResults
	c = ErrorCantGetTaskInstances
	c = ErrorCountVarHasBadFormat
	c = ErrorTaskIDHasBadFormat
	c = ErrorTaskInstanceIDHasBadFormat
	c = ErrorEndpointIDHasBadFormat
	c = ErrorTimeFrameHasBadFormat
	c = ErrorCantSaveTaskDefinitionToDB
	c = ErrorTaskDefinitionNotFound
	c = ErrorTaskDefinitionByPartnerNotFound
	c = ErrorCantDeleteTaskDefinition
	c = ErrorCantPrepareTaskForSendingOnExecution
	c = ErrorTaskDefinitionIDHasBadFormat
	c = ErrorUIDHeaderIsEmpty
	c = ErrorCantGetTasksSummaryData
	c = ErrorCantOpenTargets
	c = ErrorNoEndpointsForTargets
	c = ErrorCantExportTasksSummaryDataToXLSX
	c = ErrorCantGetTaskSummaryDetails
	c = ErrorAccessDenied
	c = ErrorCantDeleteTask
	c = ErrorCannotGetUserInfo
	c = ErrorCantDeleteTaskInstance
	c = ErrorCantDeleteExecutionResults
	c = ErrorCantGetUserSites
	c = ErrorCantGetUserEndPointsBySites
	c = ErrorCantValidateUserParameters
	c = ErrorCantGetUser
	c = ErrorCantGetActiveTriggers
	c = ErrorCantGetClosestTasks
	c = ErrorInternalServerError
	c = ErrorNotFound
	c = ErrorCantEcryptCredentials
	NewNotFoundErr(c, ErrorCantEcryptCredentials)
}
