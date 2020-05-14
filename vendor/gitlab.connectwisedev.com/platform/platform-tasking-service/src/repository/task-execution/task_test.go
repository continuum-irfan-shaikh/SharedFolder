package task_execution

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestNew(t *testing.T) {
	ex := NewTaskExecution(http.DefaultClient, map[string]string{"script": "script"}, nil)

	httpmock.Activate()
	resp := httpmock.NewBytesResponder(http.StatusOK, []byte(`{"siteDetailList":[{"siteId":1}]}`))
	url := fmt.Sprintf("%s/partners/%s/executions", "script", "1")

	httpmock.RegisterResponder(http.MethodPost, url, resp)
	config.Config.RetryStrategy.MaxNumberOfRetries = 1
	logger.Load(config.Config.Log)
	ex.ExecuteTasks(context.TODO(), apiModels.ExecutionPayload{}, "1", "script")
	ex.ExecuteTasks(context.TODO(), apiModels.ExecutionPayload{}, "1", "123")

}
