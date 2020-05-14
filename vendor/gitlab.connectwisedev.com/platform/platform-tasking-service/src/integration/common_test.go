package integration

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

type errReadCloser int

func (errReadCloser) Read(p []byte) (int, error) {
	return 0, errors.New("can't read response body")
}

func (errReadCloser) Close() error {
	return nil
}

func TestOpenDynamicGroups(t *testing.T) {
	logger.Load(config.Config.Log)
	var (
		outputStruct = struct {
			test string
		}{}
		tests = []struct {
			name        string
			expectedErr bool
			body        []byte
			err         error
			url         string
			status      int
			token       string
		}{
			{
				name:        "error while performing request, token not empty",
				expectedErr: true,
				err:         errors.New("error while performing request"),
				url:         "url1",
				token:       "something",
			},
			{
				name:        "error while performing request",
				expectedErr: true,
				err:         errors.New("error while performing request"),
				url:         "url1",
			},
			{
				name:        "bad response body: StatusNotFound",
				expectedErr: true,
				body:        nil,
				status:      http.StatusNotFound,
				err:         nil,
				url:         "url6",
			},
			{
				name:        "bad response body: StatusInternalServerError",
				expectedErr: true,
				body:        nil,
				status:      http.StatusInternalServerError,
				err:         nil,
				url:         "url2",
			},
			{
				name:        "bad response body: ioutil.ReadAll error",
				expectedErr: true,
				body:        nil,
				status:      http.StatusOK,
				err:         nil,
				url:         "url3",
			},
			{
				name:        "invalid response body: Unmarshal error",
				expectedErr: true,
				body:        []byte("invalid response body"),
				status:      http.StatusOK,
				err:         nil,
				url:         "url4",
			},
			{
				name:        "good case",
				expectedErr: false,
				body:        []byte("{\"test\":\"goodCase\"}"),
				status:      http.StatusOK,
				err:         nil,
				url:         "url5",
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", tt.url,
				func(req *http.Request) (*http.Response, error) {
					var body io.ReadCloser
					if tt.body != nil {
						body = httpmock.NewRespBodyFromBytes(tt.body)
					} else {
						body = errReadCloser(0)
					}
					return &http.Response{
						Status:     strconv.Itoa(tt.status),
						StatusCode: tt.status,
						Body:       body,
						Header:     http.Header{},
					}, tt.err
				},
			)
			err := GetDataByURL(context.Background(), &outputStruct, http.DefaultClient, tt.url, tt.token, true)
			if (err != nil) != tt.expectedErr {
				t.Fatalf("GetManagedEndpointsFromTargets: test case '%v'. Error = %v, expectedErr %v", tt.name, err, tt.expectedErr)
				return
			}
			if err != nil {
				_ = err.Error()
			}
		})
	}
}
