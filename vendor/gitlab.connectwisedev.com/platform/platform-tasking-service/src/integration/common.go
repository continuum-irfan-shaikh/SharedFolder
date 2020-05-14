package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
)

const (
	SuperUserRole = "SuperUser"
	RoleHeader    = "role"
)

//go:generate mockgen -destination=../mocks/mocks-integration/http_client_mock.go -package=common -source=./common.go

//HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
}

// NotFound is used to return an error in case of StatusNotFound response
type NotFound struct {
	url string
}

// Error - error interface implementation for NotFound type
func (err NotFound) Error() string {
	return fmt.Sprintf("Nothing found by URL %s", err.url)
}

// GetDataByURL makes Get request by the provided URL and performs Unmarshal response body to outputStructPtr struct
func GetDataByURL(ctx context.Context, outputStructPtr interface{}, httpClient HTTPClient, url, token string, isSuperUser bool) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if len(token) != 0 {
		req.Header.Set("iPlanetDirectoryPro", token)
	}

	if isSuperUser {
		req.Header.Set(RoleHeader, SuperUserRole)
	}

	req.Header.Set(transactionID.Key, transactionID.FromContext(ctx))

	logger.Log.DebugfCtx(ctx, "GetDataByURL: Performing GET by url %v", url)
	response, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			logger.Log.WarnfCtx(ctx, "GetDataByURL: error while closing body: %v", err)
		}
	}()
	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return NotFound{url: url}
	default:
		return errors.Errorf("got %s Status", http.StatusText(response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	logger.Log.DebugfCtx(ctx, "GetDataByURL: Response from %v : '%v'", url, body)

	return json.Unmarshal(body, outputStructPtr)
}
