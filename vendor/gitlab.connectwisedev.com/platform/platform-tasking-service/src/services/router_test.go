package services

import (
	"testing"

	taskExecutionResults "gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/execution-results"
	executionResultsUpdate "gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/execution-results-update"
	taskCounters "gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/task-counters"
	taskDefinitions "gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/task-definitions"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/tasks"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/templates"
)

func TestNewRouter(t *testing.T) {
	ts := tasks.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	templ := templates.NewTemplateService(nil, nil, nil)
	taskResults := taskExecutionResults.NewTaskResultsService(nil, nil, nil, nil)
	execRes := executionResultsUpdate.NewExecutionResultUpdateService(nil, nil, nil, nil, nil)
	tc := taskCounters.Service{}
	td := taskDefinitions.TaskDefinitionService{}
	dto := RouterDTO{
		ts,
		templ,
		taskResults,
		execRes,
		td,
		tc,
		nil,
		nil,
		nil,
	}
	r := NewRouter(dto)
	if r == nil {
		t.Fatal("must not be nil")
	}
}
