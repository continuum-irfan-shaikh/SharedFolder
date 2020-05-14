package cherwell

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValuesLookup(t *testing.T) {

	type valuesLookupTestcases struct {
		resp          string
		expected      *LookupResponse
		inputID       string
		inputFID      string
		expectedError error
		statusCode    int
	}

	testcases := map[string]valuesLookupTestcases{
		"ValuesLookup OK": {
			resp: `{
				"values": [
				  "New"
				],
				"hasError": false
			  }`,
			expected: &LookupResponse{
				Values: []string{
					"New",
				},
			},
			inputID:    "123",
			inputFID:   "213213",
			statusCode: 200,
		},
		"ValuesLookup with bad response": {
			resp: `{
				"values": [
					"New"
				  ],
	  			}`,
			expected:      nil,
			expectedError: errors.New("BAD_REQUEST"),
			statusCode:    400,
		},
		"ValuesLookup with err in response": {
			resp: `{
				"values": [
				  "New"
				],
				"hasError": true
			  }`,
			expected:      nil,
			expectedError: errors.New("BAD_REQUEST"),
			statusCode:    400,
		},
	}

	for _, tc := range testcases {
		server, mux := newTestServer()
		mockHandler := newMockHandler(http.MethodPost, fieldValuesLookupEndpoint, tc.resp, tc.statusCode)

		mux.Handle(fieldValuesLookupEndpoint, mockHandler)
		client, err := NewClient(Config{Host: server.URL}, getWebClient())
		assert.NoError(t, err, "Can not create client: %v", err)

		resp, err := client.ValuesLookup(tc.inputID, tc.inputFID)
		if tc.expectedError == nil {
			assert.NoError(t, err)
		}
		assert.Equal(t, tc.expected, resp)

		server.CloseClientConnections()
		server.Close()
	}
}
