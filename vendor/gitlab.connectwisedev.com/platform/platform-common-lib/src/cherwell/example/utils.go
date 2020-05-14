package main

import (
	"net/http"
	"net/http/httptest"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/cherwell"
)

const tokenEndpoint = "/token"

var (
	getByRecordIDExampleCase = struct {
		busObID, recordID string
		expected          *cherwell.BusinessObject
		err               error
		path, method      string
		resp              string
		statusCode        int
	}{
		busObID:  "id_1",
		recordID: "pub_id_1",
		expected: &cherwell.BusinessObject{
			BusinessObjectInfo: cherwell.BusinessObjectInfo{ID: "id_1", PublicID: "pub_id_1", RecordID: "rec_id_1"},
			Fields: []cherwell.FieldTemplateItem{
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
		err:    nil,
		path:   "/api/V1/getbusinessobject/busobid/id_1/busobrecid/pub_id_1",
		method: http.MethodGet,
		resp: `{
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
				"links": [
				  {
					"name": "string",
					"url": "string"
				  }
				],
				"errorCode": "",
				"errorMessage": "",
				"hasError": false
			  }`,
		statusCode: http.StatusOK,
	}
	saveExampleCase = struct {
		bo             cherwell.BusinessObject
		expectedBoInfo *cherwell.BusinessObjectInfo
		err            error
		path, method   string
		resp           string
		statusCode     int
	}{
		expectedBoInfo: &cherwell.BusinessObjectInfo{PublicID: "pub_id_1", RecordID: "rec_id_1"},
		bo: cherwell.BusinessObject{
			BusinessObjectInfo: cherwell.BusinessObjectInfo{ID: "rec_id_1", PublicID: "pub_id_1", RecordID: "rec_id_1"},
			Fields:             []cherwell.FieldTemplateItem{},
		},
		err:    nil,
		path:   "/api/V1/savebusinessobject",
		method: http.MethodPost,
		resp: `{
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
	}
	deleteByRecordIDExampleCase = struct {
		busobID, recordID string
		expectedBoInfo    *cherwell.BusinessObjectInfo
		err               error
		path, method      string
		resp              string
		statusCode        int
	}{

		recordID:       "rec_id_1",
		busobID:        "id_1",
		expectedBoInfo: &cherwell.BusinessObjectInfo{ID: "id_1", PublicID: "pub_id_1", RecordID: "rec_id_1"},
		err:            nil,
		path:           "/api/V1/deletebusinessobject/busobid/id_1/busobrecid/rec_id_1",
		method:         http.MethodDelete,
		resp: `{
				"busObId": "id_1",
				"busObPublicId": "pub_id_1",
				"busObRecId": "rec_id_1",
				"errorCode": "",
				"errorMessage": "",
				"hasError": false
			  }`,
		statusCode: http.StatusOK,
	}
)

func handleBOInfo(info *cherwell.BusinessObjectInfo) {
	//actions to handle link to business object info
}

func handleBO(info *cherwell.BusinessObject) {
	//actions to handle link to business object
}

func handleError(err error) {
	//actions to handle link to business object
}

func newTestServer() (*httptest.Server, *http.ServeMux) {
	defaultTokenResponse := []byte(`{
		"access_token": "access_token",
		"token_type": "bearer",
		"expires_in": 14399,
		"refresh_token": "refresh_token",
		"as:client_id": "client_id",
		"username": "username",
		".issued": "Tue, 31 Jul 2018 14:46:46 GMT",
		".expires": "Tue, 31 Jul 2018 18:46:46 GMT"
	  }`)

	mux := http.NewServeMux()
	mux.HandleFunc(tokenEndpoint, func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(defaultTokenResponse)
		if err != nil {
			// handle error
			handleError(err)
		}
	})
	server := httptest.NewServer(mux)
	return server, mux
}

func newMockHandler(method, path, resp string, statusCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(statusCode)
		_, err := w.Write([]byte(resp))
		if err != nil {
			// handle error
			handleError(err)
		}
	})
}
func mockExampleHandlers(mux *http.ServeMux) {

	mux.Handle(getByRecordIDExampleCase.path, newMockHandler(getByRecordIDExampleCase.method, getByRecordIDExampleCase.path, getByRecordIDExampleCase.resp, getByRecordIDExampleCase.statusCode))
	mux.Handle(saveExampleCase.path, newMockHandler(saveExampleCase.method, saveExampleCase.path, saveExampleCase.resp, saveExampleCase.statusCode))
	mux.Handle(deleteByRecordIDExampleCase.path, newMockHandler(deleteByRecordIDExampleCase.method, deleteByRecordIDExampleCase.path, deleteByRecordIDExampleCase.resp, deleteByRecordIDExampleCase.statusCode))

}
