package taskResultWebhook

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

// TaskResult is Task result, apples are green, oranges are orange
type TaskResult struct {
	ID            gocql.UUID `json:"id"`
	ResultMessage string     `json:"result_message"`
	StdOut        string     `json:"std_out"`
	StdErr        string     `json:"std_err"`
	Success       bool       `json:"success"`
	ResultWebhook string     `json:"result_webhook"`
}

// CallWebhooksFor making HTTP call by provided in ResultWebhook of task
func CallWebhooksFor(ctx context.Context, hookedTasks []TaskResult) {
	for _, task := range hookedTasks {
		respCode, err := request(ctx, "POST", task.ResultWebhook, task)
		if err != nil {
			logger.Log.ErrfCtx(
				ctx,
				errorcode.ErrorCantPerformRequest,
				"CallWebhooksFor: error during sending result webhook for Task %v. Err: %v", task.ID, err,
			)
			continue
		}

		if respCode != http.StatusOK {
			logger.Log.ErrfCtx(
				ctx,
				errorcode.ErrorCantPerformRequest,
				"CallWebhooksFor: service not accepted result webhook %v for Task %v. response code: %v", task.ResultWebhook, task.ID, respCode,
			)
		}
	}
}

func request(ctx context.Context, method, url string, body interface{}) (responseCode int, err error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return responseCode, err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return responseCode, err
	}
	request.Header.Set(common.ContentType, common.ApplicationJSON)
	request.Header.Set(transactionID.Key, transactionID.FromContext(ctx))

	response, err := (&http.Client{
		Timeout: time.Duration(config.Config.HTTPClientTimeoutSec) * time.Second,
	}).Do(request)
	if err != nil {
		return responseCode, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			logger.Log.WarnfCtx(ctx,"func request: error while closing body: %v", err)
		}
	}()
	return response.StatusCode, nil
}
