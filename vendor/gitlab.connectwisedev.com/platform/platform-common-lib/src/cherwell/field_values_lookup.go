package cherwell

import (
	"net/http"
)

// LookupRequest is a structure representing request for endpoint "/api/V1/fieldvalueslookup"
type LookupRequest struct {
	BusinessObject
	FieldID   string `json:"fieldId"`
	FieldName string `json:"fieldName"`
}

// LookupResponse is a structure representing response for endpoint "/api/V1/fieldvalueslookup"
type LookupResponse struct {
	Values []string `json:"values"`
	ErrorData
	HTTPStatusCode string `json:"httpStatusCode"`
}

// ValuesLookup get potentially valid values for Business Object fields.
func (c *Client) ValuesLookup(boid, fieldID string) (*LookupResponse, error) {
	var resp LookupResponse
	req := LookupRequest{
		BusinessObject: BusinessObject{
			BusinessObjectInfo: BusinessObjectInfo{ID: boid},
			// Fields should be initialized
			// in other way Cherwell will return an error
			Fields: []FieldTemplateItem{},
		},
		FieldID: fieldID,
	}
	err := c.performRequest(http.MethodPost, fieldValuesLookupEndpoint, &req, &resp)
	if err != nil {
		return nil, ErrorData{ErrorCode: GeneralFailureError, ErrorMessage: err.Error()}
	}

	if resp.HasError {
		return nil, resp.ErrorData
	}

	return &resp, nil
}
