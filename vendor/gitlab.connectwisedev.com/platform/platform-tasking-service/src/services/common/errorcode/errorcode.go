package errorcode

import "fmt"

// These constants store error codes for internationalization
const (
	ErrorCantDecodeInputData                           = `error_cant_decode_input_data`
	ErrorTaskDefinitionExists                          = `error_task_template_exists`
	ErrorCantGetTaskDefinitionTemplate                 = `error_cant_get_task_definition_template`
	ErrorCantGetTemplatesForExecutionMS                = `error_cant_get_templates_for_execution_ms`
	ErrorWrongTemplateType                             = `error_wrong_template_type`
	ErrorCantGetTaskExecutionResultsForManagedEndpoint = `error_cant_get_task_execution_results_for_managed_endpoint`
	ErrorCantProcessTaskExecutionResults               = `error_cant_process_task_execution_results`
	CodeUpdated                                        = `updated`
	CodeCreated                                        = `created`
	ErrorCantGetPartnersForCounters                    = `error_cant_get_partners_ids`
	ErrorCantGetCounters                               = `error_cant_get_counters`
	ErrorCantDecrypt                                   = `error_cant_decrypt`
	ErrorCantPerformRequest                            = `error_cant_perform_request`
	ErrorCantCreateNewTask                             = `error_cant_create_new_task`
	ErrorCantInsertData                                = `error_cant_insert_data`
	ErrorCantMarshall                                  = `error_cant_marshall`
	ErrorCantSaveExecutionExpiration                   = `error_save_execution_expiration`
	ErrorUsecaseProcessing                             = `error_usecase_processing_error`
	ErrorCantExecuteTasks                              = `error_cant_execute_tasks`
	ErrorCantSaveTaskToDB                              = `error_cant_save_task_to_db`
	ErrorCantUpdateTask                                = `error_cant_update_task`
	ErrorCantUpdateTaskInstances                       = `error_cant_update_task_instances`
	ErrorCantGetListOfTasksByManagedEndpoint           = `error_cannot_get_list_of_tasks_by_managed_endpoint`
	ErrorCantGetTaskByTaskID                           = `error_cant_get_task_by_task_id`
	ErrorCantGetTriggerTypes                           = `error_cant_get_trigger_types`
	ErrorCantProcessTriggerDefinitions                 = `error_cant_process_trigger_definitions`
	ErrorTaskIsNotFoundByTaskID                        = `error_task_is_not_found_by_task_id`
	ErrorCantGetTaskByTaskInstanceIDs                  = `error_cant_get_task_by_task_instance_ids`
	ErrorCantGetCountByPartnerID                       = `error_cant_get_count_by_partner_id`
	ErrorCantGetTaskInstanceCountByTaskID              = `error_cant_get_count_by_task_id`
	ErrorCantGetTaskExecutionResults                   = `error_cant_get_task_execution_results`
	ErrorCantGetTaskInstances                          = `error_cant_get_task_instances`
	ErrorCountVarHasBadFormat                          = `error_count_var_has_bad_format`
	ErrorTaskIDHasBadFormat                            = `error_task_id_has_bad_format`
	ErrorTaskInstanceIDHasBadFormat                    = `error_task_instance_id_has_bad_format`
	ErrorEndpointIDHasBadFormat                        = `error_endpoint_id_has_bad_format`
	ErrorTimeFrameHasBadFormat                         = `error_time_frame_bad_format`
	ErrorCantSaveTaskDefinitionToDB                    = `error_cant_save_task_definition_to_db`
	ErrorTaskDefinitionNotFound                        = `error_task_definition_not_found`
	ErrorTaskDefinitionByPartnerNotFound               = `error_task_definition_by_partner_not_found`
	ErrorCantDeleteTaskDefinition                      = `error_cant_delete_task_definition`
	ErrorCantPrepareTaskForSendingOnExecution          = `error_cant_prepare_task_for_sending_on_execution`
	ErrorTaskDefinitionIDHasBadFormat                  = `error_task_definition_id_has_bad_format`
	ErrorUIDHeaderIsEmpty                              = `error_uid_header_is_empty`
	ErrorCantGetTasksSummaryData                       = `error_cant_get_tasks_summary_data`
	ErrorCantOpenTargets                               = `error_cant_open_targets`
	ErrorNoEndpointsForTargets                         = `error_no_endpoints_for_targets`
	ErrorCantExportTasksSummaryDataToXLSX              = `error_cant_export_tasks_summary_data_to_xlsx`
	ErrorCantGetTaskSummaryDetails                     = `error_cant_get_task_summary_details`
	ErrorAccessDenied                                  = `error_access_denied_for_partner`
	ErrorCantDeleteTask                                = `error_cant_delete_task`
	ErrorCannotGetUserInfo                             = `error_cannot_get_user_info`
	ErrorCantDeleteTaskInstance                        = `error_cant_delete_task_instance`
	ErrorCantDeleteExecutionResults                    = `error_cant_delete_execution_results`
	ErrorCantGetUserSites                              = `error_cant_get_user_sites`
	ErrorCantGetUserEndPointsBySites                   = `error_cant_get_user_end_points_by_sites`
	ErrorCantValidateUserParameters                    = `error_cant_validate_user_parameters`
	ErrorCantGetUser                                   = `error_cant_get_user`
	ErrorCantGetActiveTriggers                         = `error_cant_get_active_triggers`
	ErrorCantExecuteTrigger                            = `error_cant_execute_trigger`
	ErrorCantGetClosestTasks                           = `error_cannot_get_closest_tasks`
	ErrorInternalServerError                           = `error_internal_server_error`
	ErrorNotFound                                      = `error_not_found`
	ErrorCantEcryptCredentials                         = `error_cant_encrypt_credentials`
	ErrorCantRecalculate                               = `error_cant_recalculate_tasks`
	ErrorKafka                                         = `error_kafka`
	ErrorApplication                                   = `error_application_error`
	ErrorCache                                         = `error_cache`
	ErrorCantProcessData                               = `error_cant_process_data`
)

// NotFoundErr is a not found connection err
type NotFoundErr struct {
	ErrorCode  string
	LogMessage string
}

// NewNotFoundErr constructor
func NewNotFoundErr(errorCode string, logMessage string) NotFoundErr {
	return NotFoundErr{
		ErrorCode:  errorCode,
		LogMessage: logMessage,
	}
}

// Error - error interface implementation for TaskNotFoundError type
func (err NotFoundErr) Error() string {
	return fmt.Sprintf("not found. %s", err.LogMessage)
}

// InternalServerErr is an internalServer connection err
type InternalServerErr struct {
	ErrorCode  string
	LogMessage string
}

// NewInternalServerErr constructor
func NewInternalServerErr(errorCode string, logMessage string) InternalServerErr {
	return InternalServerErr{ErrorCode: errorCode, LogMessage: logMessage}
}

// Error - error interface implementation for TaskNotFoundError type
func (err InternalServerErr) Error() string {
	return fmt.Sprintf("server error. %s", err.LogMessage)
}

// BadRequestErr is an BadRequest connection err
type BadRequestErr struct {
	ErrorCode  string
	LogMessage string
}

// NewBadRequestErr constructor
func NewBadRequestErr(errorCode string, logMessage string) BadRequestErr {
	return BadRequestErr{
		ErrorCode:  errorCode,
		LogMessage: logMessage,
	}
}

// Error - error interface implementation for TaskNotFoundError type
func (err BadRequestErr) Error() string {
	return fmt.Sprintf("bad request. %s", err.LogMessage)
}
