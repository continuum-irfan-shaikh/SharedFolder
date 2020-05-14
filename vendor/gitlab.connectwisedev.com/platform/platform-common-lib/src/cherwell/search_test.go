package cherwell

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {

	type findTestcases struct {
		resp                  string
		expected              *SearchResponse
		mockedID              string
		expectedError         error
		statusCode            int
		operator              string
		pageNumber            int
		pageSize              int
		expectedSearchRequest SearchRequest
	}

	testcases := map[string]findTestcases{
		"Find with all fields": {
			resp: `{
		"businessObjects": [
		 {
			"busObId": "6d",
			"busObPublicId": "19",
			"busObRecId": "94",
			"fields": [
			  {
				"dirty": false,
				"displayName": "Service Order Number",
				"fieldId": "1",
				"html": null,
				"name": "CartItemID",
				"value": ""
			  },
			  {
				"dirty": false,
				"displayName": "Description",
				"fieldId": "2",
				"html": null,
				"name": "Incident description",
				"value": ""
			  }
			],
			"links": [],
			"errorCode": null,
			"errorMessage": null,
			"hasError": false
		 },
		 {
			"busObId": "7d",
			"busObPublicId": "23",
			"busObRecId": "78",
			"fields": [
			  {
				"dirty": false,
				"displayName": "Service Order Number",
				"fieldId": "1",
				"html": null,
				"name": "CartItemID",
				"value": ""
			  },
			  {
				"dirty": false,
				"displayName": "Description",
				"fieldId": "2",
				"html": null,
				"name": "Incident description",
				"value": ""
  	         }
			],
			"links": [],
			"errorCode": null,
			"errorMessage": null,
			"hasError": false
		 }
		],
		"hasPrompts": false,
		"links": [],
		"prompts": [],
		"searchResultsFields": [],
		"simpleResults": null,
		"totalRows": 2,
		"errorCode": null,
		"errorMessage": null,
		"hasError": false
	  }`,
			expected: &SearchResponse{
				BusinessObjects: []BusinessObject{
					{
						BusinessObjectInfo: BusinessObjectInfo{
							ID:       "6d",
							RecordID: "94",
							PublicID: "19",
						},
						Fields: []FieldTemplateItem{
							{
								DisplayName: "Service Order Number",
								FieldID:     "1",
								Name:        "CartItemID",
							},
							{
								DisplayName: "Description",
								FieldID:     "2",
								Name:        "Incident description",
							},
						},
					},
					{
						BusinessObjectInfo: BusinessObjectInfo{
							ID:       "7d",
							RecordID: "78",
							PublicID: "23",
						},
						Fields: []FieldTemplateItem{
							{
								DisplayName: "Service Order Number",
								FieldID:     "1",
								Name:        "CartItemID",
							},
							{
								DisplayName: "Description",
								FieldID:     "2",
								Name:        "Incident description",
							},
						},
					},
				},
				TotalRows: 2,
			},
			mockedID:   "123",
			statusCode: 200,
			operator:   "eq",
			expectedSearchRequest: SearchRequest{
				ID:               "123",
				IncludeAllFields: true,
				PageNumber:       0,
				PageSize:         0,
				Fields:           nil,
				Filters: []Filter{
					{
						FieldID:  "1",
						Operator: "eq",
						Value:    "fieldvalue1",
					},
					{
						FieldID:  "2",
						Operator: "eq",
						Value:    "fieldvalue2",
					},
				},
			},
		},

		"Find with err in response": {
			resp: `{
		"businessObjects": null,
		"hasPrompts": false,
		"links": [],
		"prompts": [],
		"searchResultsFields": [],
		"simpleResults": null,
		"totalRows": 1,
		"errorCode": "400",
		"errorMessage": "Some Error",
		"hasError": true
	  }`,
			expected:      nil,
			mockedID:      "",
			expectedError: errors.New("Some Error"),
			statusCode:    400,
			operator:      "eq",
			expectedSearchRequest: SearchRequest{
				ID:               "",
				IncludeAllFields: true,
				PageNumber:       0,
				PageSize:         0,
				Fields:           nil,
				Filters: []Filter{
					{
						FieldID:  "1",
						Operator: "eq",
						Value:    "fieldvalue1",
					},
					{
						FieldID:  "2",
						Operator: "eq",
						Value:    "fieldvalue2",
					},
				},
			},
		},

		"Find with error in performRequest method": {
			resp:          `test error`,
			expected:      nil,
			mockedID:      "123",
			expectedError: errors.New("find: non-JSON response received HTTP status: 500 Internal Server Error; parse error: invalid character 'e' in literal true (expecting 'r'); response: test error"),
			statusCode:    500,
			operator:      "eq",
			expectedSearchRequest: SearchRequest{
				ID:               "123",
				IncludeAllFields: true,
				PageNumber:       0,
				PageSize:         0,
				Fields:           nil,
				Filters: []Filter{
					{
						FieldID:  "1",
						Operator: "eq",
						Value:    "fieldvalue1",
					},
					{
						FieldID:  "2",
						Operator: "eq",
						Value:    "fieldvalue2",
					},
				},
			},
		},

		"Find with error in resp.Responses": {
			resp: `{
		"businessObjects": [
		{
			"busObId": "6d",
			"busObPublicId": "19",
			"busObRecId": "94",
			"fields": [
			  {
				"dirty": false,
				"displayName": "Service Order Number",
				"fieldId": "1",
				"html": null,
				"name": "CartItemID",
				"value": ""
			  }
			],
			"links": [],
			"errorCode": "",
			"errorMessage": "some error",
			"hasError": true
		}
		],
		"hasPrompts": false,
		"links": [],
		"prompts": [],
		"searchResultsFields": [],
		"simpleResults": null,
		"totalRows": 1,
		"errorCode": null,
		"errorMessage": null,
		"hasError": false
	  }`,
			mockedID:      "123",
			statusCode:    400,
			expectedError: errors.New("some error"),
			operator:      "eq",
			expectedSearchRequest: SearchRequest{
				ID:               "123",
				IncludeAllFields: true,
				PageNumber:       0,
				PageSize:         0,
				Fields:           nil,
				Filters: []Filter{
					{
						FieldID:  "1",
						Operator: "eq",
						Value:    "fieldvalue1",
					},
					{
						FieldID:  "2",
						Operator: "eq",
						Value:    "fieldvalue2",
					},
				},
			},
		},

		"Find with all fields with page number '1'": {
			resp: `{
		"businessObjects": [
		 {
			"busObId": "6d",
			"busObPublicId": "19",
			"busObRecId": "94",
			"fields": [
			  {
				"dirty": false,
				"displayName": "Service Order Number",
				"fieldId": "1",
				"html": null,
				"name": "CartItemID",
				"value": ""
			  },
			  {
				"dirty": false,
				"displayName": "Description",
				"fieldId": "2",
				"html": null,
				"name": "Incident description",
				"value": ""
			  }
			],
			"links": [],
			"errorCode": null,
			"errorMessage": null,
			"hasError": false
		 }
		],
		"hasPrompts": false,
		"links": [],
		"prompts": [],
		"searchResultsFields": [],
		"simpleResults": null,
		"totalRows": 1,
		"errorCode": null,
		"errorMessage": null,
		"hasError": false
	  }`,
			expected: &SearchResponse{
				BusinessObjects: []BusinessObject{
					{
						BusinessObjectInfo: BusinessObjectInfo{
							ID:       "6d",
							RecordID: "94",
							PublicID: "19",
						},
						Fields: []FieldTemplateItem{
							{
								DisplayName: "Service Order Number",
								FieldID:     "1",
								Name:        "CartItemID",
							},
							{
								DisplayName: "Description",
								FieldID:     "2",
								Name:        "Incident description",
							},
						},
					},
				},
				TotalRows: 1,
			},
			mockedID:   "123",
			statusCode: 200,
			operator:   "eq",
			pageNumber: 1,
			pageSize:   10,
			expectedSearchRequest: SearchRequest{
				ID:               "123",
				IncludeAllFields: true,
				PageNumber:       1,
				PageSize:         10,
				Fields:           nil,
				Filters: []Filter{
					{
						FieldID:  "1",
						Operator: "eq",
						Value:    "fieldvalue1",
					},
					{
						FieldID:  "2",
						Operator: "eq",
						Value:    "fieldvalue2",
					},
				},
			},
		},

		"Find with all fields with page number 1 and page size 0": {
			resp: `{
		"businessObjects": [
		 {
			"busObId": "6d",
			"busObPublicId": "19",
			"busObRecId": "94",
			"fields": [
			  {
				"dirty": false,
				"displayName": "Service Order Number",
				"fieldId": "1",
				"html": null,
				"name": "CartItemID",
				"value": ""
			  },
			  {
				"dirty": false,
				"displayName": "Description",
				"fieldId": "2",
				"html": null,
				"name": "Incident description",
				"value": ""
			  }
			],
			"links": [],
			"errorCode": null,
			"errorMessage": null,
			"hasError": false
		 }
		],
		"hasPrompts": false,
		"links": [],
		"prompts": [],
		"searchResultsFields": [],
		"simpleResults": null,
		"totalRows": 1,
		"errorCode": null,
		"errorMessage": null,
		"hasError": false
	  }`,
			expected: &SearchResponse{
				BusinessObjects: []BusinessObject{
					{
						BusinessObjectInfo: BusinessObjectInfo{
							ID:       "6d",
							RecordID: "94",
							PublicID: "19",
						},
						Fields: []FieldTemplateItem{
							{
								DisplayName: "Service Order Number",
								FieldID:     "1",
								Name:        "CartItemID",
							},
							{
								DisplayName: "Description",
								FieldID:     "2",
								Name:        "Incident description",
							},
						},
					},
				},
				TotalRows: 1,
			},
			mockedID:   "123",
			statusCode: 200,
			operator:   "eq",
			pageNumber: 1,
			pageSize:   0,
			expectedSearchRequest: SearchRequest{
				ID:               "123",
				IncludeAllFields: true,
				PageNumber:       1,
				PageSize:         0,
				Fields:           nil,
				Filters: []Filter{
					{
						FieldID:  "1",
						Operator: "eq",
						Value:    "fieldvalue1",
					},
					{
						FieldID:  "2",
						Operator: "eq",
						Value:    "fieldvalue2",
					},
				},
			},
		},

		"Find with all fields with page number more then results count": {
			resp: `{
		"businessObjects": [],
		"hasPrompts": false,
		"links": [],
		"prompts": [],
		"searchResultsFields": [],
		"simpleResults": null,
		"totalRows": 1,
		"errorCode": null,
		"errorMessage": null,
		"hasError": false
	  }`,
			expected: &SearchResponse{
				BusinessObjects: nil,
				TotalRows:       1,
			},
			mockedID:   "123",
			statusCode: 200,
			operator:   "eq",
			pageNumber: 2,
			pageSize:   10,
			expectedSearchRequest: SearchRequest{
				ID:               "123",
				IncludeAllFields: true,
				PageNumber:       11,
				PageSize:         10,
				Fields:           nil,
				Filters: []Filter{
					{
						FieldID:  "1",
						Operator: "eq",
						Value:    "fieldvalue1",
					},
					{
						FieldID:  "2",
						Operator: "eq",
						Value:    "fieldvalue2",
					},
				},
			},
		},

		"Find with negative page size": {
			expected: nil,
			expectedError: &CherwError{
				Code:    BadRequestError,
				Message: "Page Number/Page Size is invalid. PageNumber: 5, PageSize: -7",
			},
			pageNumber: 5,
			pageSize:   -7,
		},

		"Find with negative page number": {
			expected: nil,
			expectedError: &CherwError{
				Code:    BadRequestError,
				Message: "Page Number/Page Size is invalid. PageNumber: -5, PageSize: 7",
			},
			pageNumber: -5,
			pageSize:   7,
		},
	}

	for _, tc := range testcases {
		server, mux := newTestServer()
		mockHandler := newSearchHandler(t, http.MethodPost, tc.resp, tc.statusCode, tc.expectedSearchRequest)

		req := NewSearchRequest(tc.mockedID)
		req.AddFilter("1", tc.operator, "fieldvalue1")
		req.AddFilter("2", tc.operator, "fieldvalue2")
		req.PageSize = tc.pageSize
		req.PageNumber = tc.pageNumber

		mux.Handle(searchEndpoint, mockHandler)
		client, err := NewClient(Config{Host: server.URL}, getWebClient())
		assert.NoError(t, err, "Can not create client: %v", err)

		resp, err := client.Find(*req)
		if tc.expectedError == nil {
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, resp)
		} else {
			assert.Equal(t, tc.expected, resp)
			assert.EqualError(t, err, tc.expectedError.Error())
		}
	}
}

func TestFindBoInfos(t *testing.T) {
	server, mux := newTestServer()

	type findTestcases struct {
		resp     string
		expected *SearchResponse
	}

	testcases := map[string]findTestcases{
		"Find with all fields": {
			resp: `{
		"businessObjects": [
		  {
			"busObId": "6d",
			"busObPublicId": "19",
			"busObRecId": "94",
			"links": [],
			"errorCode": null,
			"errorMessage": null,
			"hasError": false
		  },
		  {
			"busObId": "7d",
			"busObPublicId": "23",
			"busObRecId": "78",
			"links": [],
			"errorCode": null,
			"errorMessage": null,
			"hasError": false
		  }
		],
		"hasPrompts": false,
		"links": [],
		"prompts": [],
		"searchResultsFields": [],
		"simpleResults": null,
		"totalRows": 1,
		"errorCode": null,
		"errorMessage": null,
		"hasError": false
	  }`,
			expected: &SearchResponse{BusinessObjects: []BusinessObject{
				{
					BusinessObjectInfo: BusinessObjectInfo{
						ID:       "6d",
						RecordID: "94",
						PublicID: "19",
					},
				},
				{
					BusinessObjectInfo: BusinessObjectInfo{
						ID:       "7d",
						RecordID: "78",
						PublicID: "23",
					},
				},
			},
				TotalRows: 1,
			},
		},
	}

	for _, tc := range testcases {
		mockHandler := newMockHandler(http.MethodPost, searchEndpoint, tc.resp, http.StatusOK)

		req := NewSearchRequest("123")
		req.AddFilter("1", " eq", "fieldvalue1")
		req.AddFilter("2", "eq", "fieldvalue2")

		mux.Handle(searchEndpoint, mockHandler)
		client, err := NewClient(Config{Host: server.URL}, getWebClient())
		assert.NoError(t, err, "Can not create client: %v", err)

		bos, err := client.FindBoInfos(*req)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, bos)
	}
}

func TestSetSpecificFields(t *testing.T) {

	type findTestcase struct {
		input       SearchRequest
		expected    SearchRequest
		inputFields []string
	}

	tc := findTestcase{
		input: SearchRequest{
			IncludeAllFields: true,
		},
		expected: SearchRequest{
			Fields: []string{
				"1", "2", "3",
			},
			IncludeAllFields: false,
		},
		inputFields: []string{
			"1", "2", "3",
		},
	}

	tc.input.SetSpecificFields(tc.inputFields)
	assert.Equal(t, tc.expected, tc.input)
}

func TestAppendSpecificFields(t *testing.T) {

	type findTestcase struct {
		input        SearchRequest
		expected     SearchRequest
		appendFields []string
	}

	tc := findTestcase{
		input: SearchRequest{
			Fields: []string{
				"1",
			},
			IncludeAllFields: true,
		},
		expected: SearchRequest{
			Fields: []string{
				"1", "2", "3",
			},
			IncludeAllFields: false,
		},
		appendFields: []string{
			"2", "3",
		},
	}

	tc.input.AppendSpecificFields(tc.appendFields...)
	assert.Equal(t, tc.expected, tc.input)
}

func newSearchHandler(t *testing.T, method, resp string, statusCode int, expectedSR SearchRequest) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var searchRequest SearchRequest

		err = json.Unmarshal(b, &searchRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		assert.Equal(t, expectedSR, searchRequest)

		w.WriteHeader(statusCode)
		w.Write([]byte(resp))
	})
}
