package task_execution

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
)

// NewTaskExecution returns new task exec repo
func NewTaskExecution(client *http.Client, domains map[string]string, log logger.Logger) *TaskExecution {
	return &TaskExecution{
		client:  client,
		domains: domains,
		log:     log,
	}
}

// TaskExecution represents repo to execute task
type TaskExecution struct {
	client  *http.Client
	domains map[string]string
	log     logger.Logger
}

// ExecuteTasks executes a task with payload
func (t *TaskExecution) ExecuteTasks(ctx context.Context, payload apiModels.ExecutionPayload, partnerID, taskType string) error {
	domain, ok := t.domains[taskType]
	if !ok {
		return fmt.Errorf("there isn't domain for task type %s", taskType)
	}

	url := fmt.Sprintf("%s/partners/%s/executions", domain, partnerID)

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := common.HTTPRequestWithRetry(ctx, t.client, http.MethodPost, url, body)
	if err != nil {
		return err
	}
	defer t.close(ctx, resp.Body)
	return nil
}

func (t *TaskExecution) close(ctx context.Context, c io.Closer) {
	if err := c.Close(); err != nil {
		t.log.WarnfCtx(ctx, "TaskExecution: can't close body")
	}
}
