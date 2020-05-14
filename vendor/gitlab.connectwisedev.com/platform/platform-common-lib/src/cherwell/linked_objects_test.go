package cherwell

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const SaveRelatedBusinessObjectPath = "/api/V1/saverelatedbusinessobject"

func newMockHandlerWithBodyCheck(method, path, resp string, statusCode int, expectedRequestBody string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			panic(fmt.Errorf("unexpected method in call to http mock: path = %s, method = %s", path, method))
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(fmt.Errorf("error while reading body in call to http mock: path = %s, body = %s ___ expected body = %s", path, string(body), expectedRequestBody))
		}
		if string(body) != expectedRequestBody {
			panic(fmt.Errorf("unexpected body in call to http mock: path = %s, body = %s", path, string(body)))
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))
	})
}

var (
	SaveRelatedTestcases = map[string]struct {
		rBo             RelatedBusinessObject
		expectedRBoInfo *RelatedBusinessObjectInfo
		err             error
		path, method    string
		resp            string
		statusCode      int
	}{
		"Test SaveRelatedBusinessObject success": {
			expectedRBoInfo: &RelatedBusinessObjectInfo{
				BusinessObjectInfo: BusinessObjectInfo{PublicID: "pub_id_1", RecordID: "rec_id_1"},
				RelatedInfo: RelatedInfo{
					ParentBusObID:       "par_id_1",
					ParentBusObRecID:    "par_rec_id_1",
					ParentBusObPublicID: "par_pub_id_1",
					RelationshipID:      "rel_id_1",
				},
			},
			rBo: RelatedBusinessObject{
				BusinessObject: BusinessObject{
					BusinessObjectInfo: BusinessObjectInfo{ID: "rec_id_1", PublicID: "pub_id_1", RecordID: "rec_id_1"},
					Fields:             []FieldTemplateItem{},
				},
				RelatedInfo: RelatedInfo{
					ParentBusObID:    "par_id_1",
					ParentBusObRecID: "par_rec_id_1",
					RelationshipID:   "rel_id_1",
				},
			},
			err:    nil,
			path:   SaveRelatedBusinessObjectPath,
			method: http.MethodPost,
			resp: `{
				"parentBusObId":"par_id_1",
				"parentBusObPublicId": "par_pub_id_1",
  				"parentBusObRecId": "par_rec_id_1",
  				"relationshipId": "rel_id_1",
				"busObPublicId": "pub_id_1",
				"busObRecId": "rec_id_1",
				"cacheKey": "string",
				"fieldValidationErrors": [],
				"notificationTriggers": [],
				"errorCode": "",
				"errorMessage": "",
				"hasError": false
			  }`,
			statusCode: http.StatusOK,
		},
		"Test SaveRelatedBusinessObject response with cherwell error": {
			rBo: RelatedBusinessObject{
				BusinessObject: BusinessObject{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_wrong", PublicID: "pub_id_wrong", RecordID: "rec_id_wrong"},
					Fields:             []FieldTemplateItem{},
				},
				RelatedInfo: RelatedInfo{
					ParentBusObID:    "par_id_1",
					ParentBusObRecID: "par_rec_id_1",
					RelationshipID:   "rel_id_1",
				},
			},
			expectedRBoInfo: &RelatedBusinessObjectInfo{},
			err: &CherwError{
				Code:    "RecordNotFound",
				Message: "Record not found",
			},
			path:   SaveRelatedBusinessObjectPath,
			method: http.MethodPost,
			resp: `{
					"busObId": null,
					"busObPublicId": null,
					"busObRecId": null,
					"fields": [],
					"links": [],
					"Code": "RecordNotFound",
					"Message": "Record not found",
					"hasError": true
				}`,
			statusCode: http.StatusInternalServerError,
		},
		"Test SaveRelatedBusinessObject response with client error": {
			expectedRBoInfo: nil,
			rBo: RelatedBusinessObject{
				BusinessObject: BusinessObject{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_wrong_json", PublicID: "pub_id_wrong_json", RecordID: "rec_id_wrong_json"},
					Fields:             []FieldTemplateItem{},
				},
			},
			err:        errors.New("non-JSON response received HTTP status: 500 Internal Server Error; parse error: unexpected end of JSON input; response: {"),
			path:       SaveRelatedBusinessObjectPath,
			method:     http.MethodPost,
			resp:       `{`,
			statusCode: http.StatusInternalServerError,
		},
		"Test SaveRelatedBusinessObject response with client error with non-escaped symbols": {
			expectedRBoInfo: nil,
			rBo: RelatedBusinessObject{
				BusinessObject: BusinessObject{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_wrong_json", PublicID: "pub_id_wrong_json", RecordID: "rec_id_wrong_json"},
					Fields:             []FieldTemplateItem{},
				},
			},
			err:        errors.New("non-JSON response received HTTP status: 500 Internal Server Error; parse error: invalid character '<' looking for beginning of object key string; response: {  <BODY>"),
			path:       SaveRelatedBusinessObjectPath,
			method:     http.MethodPost,
			resp:       "{\r\n<BODY>",
			statusCode: http.StatusInternalServerError,
		},
	}
	LinkageTestCases = map[string]struct {
		parentBusobID, parentRecID string
		err                        error
		path, method               string
		busobRecID, publicID       string
		relationshipID             string
		statusCode                 int
		expectedStatusCode         int
	}{
		"Successful link of Business Objects": {
			parentBusobID:      "pbo_1",
			parentRecID:        "p_rec_id1",
			busobRecID:         "bo_rec_id1",
			publicID:           "bo_1",
			relationshipID:     "relid_1",
			path:               "/api/V1/linkrelatedbusinessobject/parentbusobid/pbo_1/parentbusobrecid/p_rec_id1/relationshipid/relid_1/busobid/bo_1/busobrecid/bo_rec_id1",
			method:             http.MethodGet,
			statusCode:         http.StatusOK,
			expectedStatusCode: http.StatusOK,
		},
	}

	UnLinkageTestCases = map[string]struct {
		parentBusobID, parentRecID string
		err                        error
		path, method               string
		busobRecID, publicID       string
		relationshipID             string
		statusCode                 int
		expectedStatusCode         int
	}{
		"Successful unlink of Business Objects": {
			parentBusobID:      "pbo_1",
			parentRecID:        "p_rec_id1",
			relationshipID:     "relid_1",
			busobRecID:         "bo_rec_id1",
			publicID:           "bo_1",
			path:               "/api/V1/unlinkrelatedbusinessobject/parentbusobid/pbo_1/parentbusobrecid/p_rec_id1/relationshipid/relid_1/busobid/bo_1/busobrecid/bo_rec_id1",
			method:             http.MethodDelete,
			statusCode:         http.StatusOK,
			expectedStatusCode: http.StatusOK,
		},
	}
)

func TestSaveRelated(t *testing.T) {
	for name, tc := range SaveRelatedTestcases {
		t.Run(name, func(t *testing.T) {
			server, mux := newTestServer()
			mux.Handle(tc.path, newMockHandler(tc.method, tc.path, tc.resp, tc.statusCode))
			client, err := NewClient(Config{Host: server.URL}, getWebClient())
			if err != nil {
				assert.NoError(t, err, "Can not create client: %v", err)
			}

			actualBoInfo, err := client.SaveRelatedBusinessObject(tc.rBo)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
				server.CloseClientConnections()
				server.Close()
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRBoInfo, actualBoInfo, "Unexpected result on testcase '%s': ", name)
			server.CloseClientConnections()
			server.Close()
		})
	}
}

func TestLinkage(t *testing.T) {
	for _, tc := range LinkageTestCases {
		server, mux := newTestServer()
		mux.Handle(tc.path, newMockHandler(tc.method, tc.path, "", tc.statusCode))
		client, err := NewClient(Config{Host: server.URL}, getWebClient())
		if err != nil {
			assert.NoError(t, err, "Can not create client: %v", err)
		}
		linkage := LinkedObject{
			ParentBusObID:    tc.parentBusobID,
			ParentBusObRecID: tc.parentRecID,
			BusObID:          tc.publicID,
			BusObRecID:       tc.busobRecID,
			RelationshipID:   tc.relationshipID,
		}
		err = client.Link(&linkage)

		if tc.err != nil {
			assert.EqualError(t, err, tc.err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestUnLinkage(t *testing.T) {
	for _, tc := range UnLinkageTestCases {
		server, mux := newTestServer()
		mux.Handle(tc.path, newMockHandler(tc.method, tc.path, "", tc.statusCode))
		client, err := NewClient(Config{Host: server.URL}, getWebClient())
		if err != nil {
			assert.NoError(t, err, "Can not create client: %v", err)
		}
		linkage := LinkedObject{
			ParentBusObID:    tc.parentBusobID,
			ParentBusObRecID: tc.parentRecID,
			BusObID:          tc.publicID,
			BusObRecID:       tc.busobRecID,
			RelationshipID:   tc.relationshipID,
		}
		err = client.Unlink(&linkage)

		if tc.err != nil {
			assert.EqualError(t, err, tc.err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestGetRelatedBusinessObject(t *testing.T) {
	testcases := map[string]struct {
		request             RelatedBusinessObjectsRequest
		expectedRequestJSON string
		expected            *RelatedBusinessObjectsResponse
		err                 error
		path, method        string
		req, resp           string
		statusCode          int
	}{
		"GetRelatedObject success": {
			request: RelatedBusinessObjectsRequest{
				ParentBusObID:    "pbo_1",
				ParentBusObRecID: "p_rec_id1",
				RelationshipID:   "relid_1",
				Limit:            10,
				Offset:           1,
				Fields:           []string{"field_id_1", "field_id_2"},
				UseDefaultGrid:   true,
				CustomGridID:     "some_grid_id",
				IncludeAllFields: true,
				Sorting: []Sorting{
					{
						FieldID:       "field_id_1",
						SortDirection: SortAsc,
					},
					{
						FieldID:       "field_id_2",
						SortDirection: SortDesc,
					},
				},
				Filters: []Filter{
					{
						FieldID:  "field_id_1",
						Operator: OpEqual,
						Value:    "value_1",
					},
					{
						FieldID:  "field_id_2",
						Operator: OpLessThan,
						Value:    "1234",
					},
				},
			},
			expectedRequestJSON: `{"customGridId":"some_grid_id","parentBusObId":"pbo_1","parentBusObRecId":"p_rec_id1","relationshipId":"relid_1","filters":[{"fieldId":"field_id_1","operator":"eq","value":"value_1"},{"fieldId":"field_id_2","operator":"lt","value":"1234"}],"sorting":[{"fieldId":"field_id_1","sortDirection":1},{"fieldId":"field_id_2"}],"pageSize":10,"pageNumber":1,"useDefaultGrid":true,"allFields":true,"fieldsList":["field_id_1","field_id_2"]}`,
			expected: &RelatedBusinessObjectsResponse{
				TotalRows: 3,
				BusinessObjects: []BusinessObject{
					{
						BusinessObjectInfo: BusinessObjectInfo{ID: "id_1", PublicID: "pub_id_1", RecordID: "rec_id_1"},
						Fields: []FieldTemplateItem{
							{
								Dirty:       true,
								DisplayName: "display_name_1",
								FieldID:     "field_id_1",
								HTML:        `<div id='somediv1'></div>`,
								Name:        "name_1",
								Value:       "value_1",
							},
							{
								Dirty:       false,
								DisplayName: "display_name_2",
								FieldID:     "field_id_2",
								HTML:        `<span id='somediv2'></span>`,
								Name:        "name_2",
								Value:       "value_2",
							},
						},
					},
					{
						BusinessObjectInfo: BusinessObjectInfo{ID: "id_2", PublicID: "pub_id_2", RecordID: "rec_id_2"},
						Fields: []FieldTemplateItem{
							{
								Dirty:       false,
								DisplayName: "display_name_2",
								FieldID:     "field_id_2",
								HTML:        `<span id='somediv2'></span>`,
								Name:        "name_2",
								Value:       "value_2",
							},
						},
					},
					{
						BusinessObjectInfo: BusinessObjectInfo{ID: "id_3", PublicID: "pub_id_3", RecordID: "rec_id_3"},
						Fields: []FieldTemplateItem{
							{
								Dirty:       true,
								DisplayName: "display_name_1",
								FieldID:     "field_id_1",
								HTML:        `<div id='somediv1'></div>`,
								Name:        "name_1",
								Value:       "value_1",
							},
						},
					},
				},
			},
			err:    nil,
			path:   "/api/V1/getrelatedbusinessobject",
			method: http.MethodPost,
			resp: `{
				"errorCode": null,
				"errorMessage": null,
				"hasError": false,
				"links": [],
				"relatedBusinessObjects": [
					{
						"busObId": "id_1",
						"busObPublicId": "pub_id_1",
						"busObRecId": "rec_id_1",
						"fields": [
							{
								"dirty": true,
								"displayName": "display_name_1",
								"fieldId": "field_id_1",
								"html": "<div id='somediv1'></div>",
								"name": "name_1",
								"value": "value_1"
							},
							{
								"dirty": false,
								"displayName": "display_name_2",
								"fieldId": "field_id_2",
								"html": "<span id='somediv2'></span>",
								"name": "name_2",
								"value": "value_2"
							}
						],
						"errorCode": "",
						"errorMessage": "",
						"hasError": false
					},
					{
						"busObId": "id_2",
						"busObPublicId": "pub_id_2",
						"busObRecId": "rec_id_2",
						"fields": [
							{
								"dirty": false,
								"displayName": "display_name_2",
								"fieldId": "field_id_2",
								"html": "<span id='somediv2'></span>",
								"name": "name_2",
								"value": "value_2"
							}
						],
						"errorCode": "",
						"errorMessage": "",
						"hasError": false
					},
					{
						"busObId": "id_3",
						"busObPublicId": "pub_id_3",
						"busObRecId": "rec_id_3",
						"fields": [
							{
								"dirty": true,
								"displayName": "display_name_1",
								"fieldId": "field_id_1",
								"html": "<div id='somediv1'></div>",
								"name": "name_1",
								"value": "value_1"
							}
						],
						"errorCode": "",
						"errorMessage": "",
						"hasError": false
					}
				],
				"relationshipId": "944410200cf536ec8b2715426a88d6679a23d9bab1",
				"totalRecords": 3
			}`,
			statusCode: http.StatusOK,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			server, mux := newTestServer()
			mux.Handle(tc.path, newMockHandlerWithBodyCheck(tc.method, tc.path, tc.resp, tc.statusCode, tc.expectedRequestJSON))
			client, err := NewClient(Config{Host: server.URL}, getWebClient())
			if err != nil {
				assert.NoError(t, err, "Can not create client: %v", err)
			}
			result, err := client.GetRelatedBusinessObjects(&tc.request)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected.TotalRows, result.TotalRows)
			assert.ElementsMatch(t, tc.expected.BusinessObjects, result.BusinessObjects)

			server.CloseClientConnections()
			server.Close()
		})
	}
}
