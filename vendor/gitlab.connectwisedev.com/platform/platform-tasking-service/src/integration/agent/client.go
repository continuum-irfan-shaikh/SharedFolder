package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gocql/gocql"
	agentModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
)

type client struct {
	httpClient integration.HTTPClient
	agentELB   string
	log        logger.Logger
}

// NewClient constructs agent client
func NewClient(httpClient integration.HTTPClient, agentELB string, log logger.Logger) *client {
	return &client{
		httpClient: httpClient,
		agentELB:   agentELB,
		log:        log,
	}
}

type payload struct {
	EndpointID gocql.UUID `json:"endpointID"`
	Data       string     `json:"data"`
}

//Encrypt encrypts credentials by public key stored for particular endpoint
func (c *client) Encrypt(ctx context.Context, endpointID gocql.UUID, credentials agentModels.Credentials) (encrypted agentModels.Credentials, err error) {
	encrypted.UseCurrentUser = credentials.UseCurrentUser

	encrypted.Password, err = c.encrypt(ctx, payload{
		EndpointID: endpointID,
		Data:       credentials.Password,
	})

	if err != nil {
		return encrypted, fmt.Errorf("Encrypt: can't encrypt password. endpointID: %v. err :%s", endpointID, err.Error())
	}

	encrypted.Username, err = c.encrypt(ctx, payload{
		EndpointID: endpointID,
		Data:       credentials.Username,
	})

	if err != nil {
		return encrypted, fmt.Errorf("Encrypt: can't encrypt username. endpointID: %v. err :%s", endpointID, err.Error())
	}

	encrypted.Domain, err = c.encrypt(ctx, payload{
		EndpointID: endpointID,
		Data:       credentials.Domain,
	})

	if err != nil {
		return encrypted, fmt.Errorf("Encrypt: can't encrypt domain. endpointID: %v. err :%s", endpointID, err.Error())
	}

	return
}

func (c *client) encrypt(ctx context.Context, body payload) (result string, err error) {
	if len(body.Data) == 0 {
		return "", nil
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("can't marshal body %v .err: %s", body, err.Error())
	}

	url := c.agentELB + "/encrypt"

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("can't create post request %v.err: %s", request, err.Error())
	}
	request.Header.Set(common.ContentType, common.ApplicationJSON)
	request.Header.Set(transactionID.Key, transactionID.FromContext(ctx))

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("can't proceed with post request %v.err: %s", request, err.Error())
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			c.log.WarnfCtx(ctx,"func request: error while closing body: %v", err)
		}
	}()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("can't read responseBytes %v.err: %s", responseBytes, err.Error())
	}

	c.log.DebugfCtx(ctx, "Response from POST url %v , status code '%v',  payload '%s'", url, response.StatusCode, string(responseBytes))
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("StatusCode  %v recieved. response: %s", response.StatusCode, string(responseBytes))
	}

	var responsePayload payload
	err = json.Unmarshal(responseBytes, &responsePayload)
	if err != nil {
		return "", fmt.Errorf("can't unmarshal responseBytes %s.err: %s", string(responseBytes), err.Error())
	}

	return responsePayload.Data, nil
}
