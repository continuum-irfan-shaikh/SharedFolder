package cherwell

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/webClient"
)

var (
	removableSymbols  = []byte{'\n', '\r'}
	replacementSymbol = []byte{' '}
)

// Client contains information about Cherwell API response after authorization request
type Client struct {
	tokenResponse *tokenResponse
	client        webClient.HTTPClientService
	config        Config
	mutex         sync.RWMutex
	logger        Logger
}

// NewClient creates an instance of Client and obtains access token
func NewClient(conf Config, webClient webClient.HTTPClientService) (*Client, error) {
	if webClient == nil {
		return nil, errors.New("cherwell.NewClient: httpClient can not be nil")
	}

	client := &Client{config: conf, client: webClient, logger: NewNopLogger()}

	if err := client.getAccessToken(); err != nil {
		return nil, fmt.Errorf("cherwell.NewClient: authentication error: %s", err)
	}
	return client, nil
}

// NewClientWithLogger creates an instance of Client with custom logger and obtains access token
func NewClientWithLogger(conf Config, webClient webClient.HTTPClientService, logger Logger) (*Client, error) {
	if logger == nil {
		return nil, fmt.Errorf("cherwell.NewClientWithLogger: logger can not be nil")
	}

	client, err := NewClient(conf, webClient)
	if err != nil {
		return nil, fmt.Errorf("cherwell.NewClientWithLogger: %v", err)
	}

	client.logger = logger

	return client, nil
}

// getAccessToken is a func for authenticating using the Internal Mode. In this scenario, the User logs in to the
// REST API using CSM credentials. CSM returns a JSON response that includes information about the access token or error
// with status code.
func (c *Client) getAccessToken() error {
	var errRes errorResponse

	vals := url.Values{
		"grant_type": {passwordGrantType},
		"client_id":  {c.config.ClientID},
		"username":   {c.config.UserName},
		"password":   {c.config.Password},
		"auth_mode":  {c.config.AuthMode},
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	req, err := http.NewRequest(http.MethodPost, c.config.Host+tokenEndpoint, strings.NewReader(vals.Encode()))
	if err != nil {
		return &CherwError{
			Code:    AuthorizationError,
			Message: fmt.Sprintf("getAccessToken authentication failed while creating request: %s", err),
		}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		return &CherwError{
			Code:    AuthorizationError,
			Message: fmt.Sprintf("getAccessToken authentication request failed: %s", err),
		}
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode == http.StatusOK {
		return unmarshalRespBody(resp, &c.tokenResponse)
	}

	if err = unmarshalRespBody(resp, &errRes); err != nil {
		return &CherwError{
			Code:    AuthorizationError,
			Message: fmt.Sprintf("getAccessToken failed with error: %s , during deserialization response body: %s", err, resp.Body),
		}
	}

	return errRes
}

func formatCherwellResponse(data []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)
	for _, s := range removableSymbols {
		result = bytes.Replace(result, []byte{s}, replacementSymbol, -1)
	}
	return result
}

// unmarshalRespBody perform unmarshals Cherwell API response body into given structure
func unmarshalRespBody(resp *http.Response, out interface{}) (err error) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &CherwError{
			Code:    ReadResponseError,
			Message: err.Error(),
		}
	}

	// linkrelatedbusinessobject returns empty body and status code 200 instead of response described in swagger
	if (out == nil || len(data) == 0) && resp.StatusCode == http.StatusOK {
		return nil
	}

	err = json.Unmarshal(data, &out)
	if err != nil {
		return &CherwError{
			Code: UnmarshalError,
			Message: fmt.Sprintf(
				"non-JSON response received HTTP status: %s; parse error: %s; response: %s",
				resp.Status,
				err.Error(),
				formatCherwellResponse(data),
			),
		}
	}
	return nil
}

func (c *Client) performRequest(method, path string, reqEntity interface{}, respEntity interface{}) error {
	resp, err := c.getResponse(method, path, reqEntity)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint: errcheck
	if resp.StatusCode == http.StatusUnauthorized {
		if err = c.refreshAccessToken(); err != nil {
			return err
		}
		resp, err = c.getResponse(method, path, reqEntity)
		if err != nil {
			return err
		}
		defer resp.Body.Close() // nolint: errcheck
	}

	return unmarshalRespBody(resp, respEntity)
}

func (c *Client) refreshAccessToken() error {
	retry := 0
	for retry < retryCount {
		err := c.getAccessToken()
		if err != nil {
			retry++
			continue
		} else {
			return nil
		}
	}
	return &CherwError{
		Code:    AuthorizationError,
		Message: fmt.Sprintf("failed to refresh access token: retry count is exceeded %d", retryCount),
	}
}

func (c *Client) createRequest(method, path string, data io.Reader) (*http.Request, error) {
	requestPath := c.config.Host + path
	req, err := http.NewRequest(method, requestPath, data)
	if err != nil {
		return nil, &CherwError{
			Code:    CreateRequestError,
			Message: err.Error(),
		}
	}

	c.mutex.RLock()
	req.Header.Set("Authorization", "Bearer "+c.tokenResponse.AccessToken)
	c.mutex.RUnlock()
	return req, nil
}

func (c *Client) getResponse(method, path string, reqEntity interface{}) (*http.Response, error) {
	var (
		reqBody io.Reader
		data    []byte
		err     error
	)

	if reqEntity != nil {
		data, err = json.Marshal(&reqEntity)
		if err != nil {
			c.logFailedRequestPayload("", method, path, data, err)

			return nil, &CherwError{
				Code:    MarshalError,
				Message: fmt.Sprintf("getResponse: failed to marshal request: %s", err),
			}
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := c.createRequest(method, path, reqBody)
	if err != nil {
		c.logFailedRequestPayload("", method, path, data, err)

		return nil, &CherwError{
			Code:    CreateRequestError,
			Message: fmt.Sprintf("getResponse: failed to create request: %s", err),
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logFailedRequestPayload("", method, path, data, err)

		return nil, &CherwError{
			Code:    DoRequestError,
			Message: fmt.Sprintf("getResponse: failed to send request: %s", err),
		}
	}

	if failedHTTPCode(resp.StatusCode) {
		c.logFailedRequestPayload("", method, path, data, nil)
	}

	return resp, err
}

// logFailedRequestPayload logs the payload if exists
func (c *Client) logFailedRequestPayload(trID string, method, url string, data []byte, err error) {
	if len(data) != 0 {
		c.logger.Log(trID, "failed cherwell request payload: method=%s url=%s err=%v payload: %s",
			method, url, err, string(data))
	}
}

func failedHTTPCode(code int) bool {
	return code < http.StatusContinue || code >= http.StatusBadRequest
}
