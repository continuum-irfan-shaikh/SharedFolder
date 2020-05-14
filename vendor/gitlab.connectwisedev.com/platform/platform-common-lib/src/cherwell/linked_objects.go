package cherwell

import (
	"fmt"
	"net/http"
	"strings"
)

// List of p[arameters for link calls
const (
	// Cherwell swagger documentation says that allfields is a:
	// "Flag to include all related Business Object fields. Default is true if not supplied. If true, then UseDefaultGrid is not used."
	// But call to getrelatedbusinessobjects without allfields=false and usedefaultgrid=false
	// in different cases can return one or zero records when there are more than one.
	allFieldsParamVal   = false
	defaultGrigParamVal = false
)

// RelatedInfo represents Business object and its relation to parent business object
type RelatedInfo struct {
	ParentBusObID       string `json:"parentBusObId"`
	ParentBusObRecID    string `json:"parentBusObRecId"`
	RelationshipID      string `json:"relationshipId"`
	ParentBusObPublicID string `json:"parentBusObPublicId"`
}

// relatedBusinessObjectsResponse represents response from Cherwell API containing business objects related to parent
type relatedBusinessObjectsResponse struct {
	Responses    []businessObjectResponse `json:"relatedBusinessObjects"`
	TotalRecords int                      `json:"totalRecords"`
	ErrorData
}

// saveUpdateRelatedResponse represents structure of response which Cherwell API sent after request from client to save or update
// Business Object with relations
type saveUpdateRelatedResponse struct {
	saveUpdateResponse
	RelatedInfo
	CherwError
}

// LinkedObject is for linking/unlinking actions
type LinkedObject struct {
	ParentBusObID    string
	ParentBusObRecID string
	RelationshipID   string
	BusObID          string
	BusObRecID       string
}

// RelatedBusinessObjectsRequest represents request to Cherwell API for retrieving business objects related to parent
type RelatedBusinessObjectsRequest struct {
	CustomGridID     string    `json:"customGridId,omitempty"`
	ParentBusObID    string    `json:"parentBusObId"`
	ParentBusObRecID string    `json:"parentBusObRecId"`
	RelationshipID   string    `json:"relationshipId"`
	Filters          []Filter  `json:"filters,omitempty"`
	Sorting          []Sorting `json:"sorting,omitempty"`
	Limit            int       `json:"pageSize,omitempty"`
	Offset           int       `json:"pageNumber,omitempty"`
	UseDefaultGrid   bool      `json:"useDefaultGrid"`
	IncludeAllFields bool      `json:"allFields"`
	Fields           []string  `json:"fieldsList"`
}

// RelatedBusinessObjectsResponse is a structure representing response on retrieving business objects related to parent
type RelatedBusinessObjectsResponse struct {
	BusinessObjects []BusinessObject
	TotalRows       int
}

// SaveRelatedBusinessObject creates or updates an existing Business Object related to parent.
// API endpoint: POST /api/V1/saverelatedbusinessobject
func (c *Client) SaveRelatedBusinessObject(bo RelatedBusinessObject) (*RelatedBusinessObjectInfo, error) {
	saveResp := new(saveUpdateRelatedResponse)
	request := saveUpdateRelatedRequest{
		RelatedBusinessObject: bo,
		CacheScope:            cacheScopeSession,
		Persist:               true,
	}

	err := c.performRequest(http.MethodPost, createUpdateRelatedBOEndpoint, &request, saveResp)
	if err != nil {
		return nil, err
	}
	if saveResp.CherwError.Message != "" {
		return nil, processSaveRelatedError(&saveResp.CherwError)
	}

	rBoInfo := &RelatedBusinessObjectInfo{
		BusinessObjectInfo: BusinessObjectInfo{
			PublicID: saveResp.PublicID,
			RecordID: saveResp.RecordID,
		},
		RelatedInfo: saveResp.RelatedInfo,
	}
	return rBoInfo, nil
}

func processSaveRelatedError(err *CherwError) error {
	switch {
	case strings.Contains(err.Message, RecordNotFoundError):
		return ErrorData{ErrorCode: RecordNotFoundError, ErrorMessage: err.Message}
	case strings.Contains(err.Message, GeneralFailureError):
		return ErrorData{ErrorCode: GeneralFailureError, ErrorMessage: err.Message}
	}
	return ErrorData{ErrorCode: UndefinedError, ErrorMessage: err.Message}
}

// Link links related Business Objects
func (c *Client) Link(l *LinkedObject) error {
	// no content in response body
	reqPath := fmt.Sprintf(linkBOsPath, l.ParentBusObID, l.ParentBusObRecID, l.RelationshipID, l.BusObID, l.BusObRecID)
	var errorResponse *CherwError
	if err := c.performRequest(http.MethodGet, reqPath, nil, &errorResponse); err != nil {
		return err
	}

	if errorResponse != nil {
		return processLinkedObjectError(errorResponse)
	}
	return nil
}

// Unlink unlinks related Business Objects
func (c *Client) Unlink(l *LinkedObject) error {
	// no content in response body
	reqPath := fmt.Sprintf(unlinkBOsPath, l.ParentBusObID, l.ParentBusObRecID, l.RelationshipID, l.BusObID, l.BusObRecID)
	// next lines are a kind of a hack to check error response appearance
	// we provide pointer on pointer value of cherwell error to check if error was occurred, if it was deadressed pointer to error will be not nil
	var errorResponse *CherwError
	if err := c.performRequest(http.MethodDelete, reqPath, nil, &errorResponse); err != nil {
		return err
	}

	if errorResponse != nil {
		return processLinkedObjectError(errorResponse)
	}
	return nil
}

func processLinkedObjectError(err *CherwError) error {
	switch {
	case strings.Contains(err.Message, RecordNotFoundError):
		return ErrorData{ErrorCode: RecordNotFoundError, ErrorMessage: err.Message}
	case strings.Contains(err.Message, GeneralFailureError):
		return ErrorData{ErrorCode: GeneralFailureError, ErrorMessage: err.Message}
	case strings.Contains(err.Message, DuplicateEntryError):
		return ErrorData{ErrorCode: DuplicateLinkError, ErrorMessage: err.Message}
	}

	return &CherwError{Code: UndefinedError, Message: err.Message}
}

// GetRelatedBusinessObjects retrieves a bucket of business objects related to parent
func (c *Client) GetRelatedBusinessObjects(req *RelatedBusinessObjectsRequest) (*RelatedBusinessObjectsResponse, error) {
	bosResp := new(relatedBusinessObjectsResponse)

	if err := c.performRequest(http.MethodPost, getRelatedObjectsEndpoint, req, bosResp); err != nil {
		return nil, err
	}

	response := &RelatedBusinessObjectsResponse{
		BusinessObjects: make([]BusinessObject, 0, len(bosResp.Responses)),
		TotalRows:       bosResp.TotalRecords,
	}

	for _, b := range bosResp.Responses {
		if b.HasError {
			return nil, &CherwError{Code: b.ErrorCode, Message: b.ErrorMessage}
		}
		response.BusinessObjects = append(response.BusinessObjects, BusinessObject{
			BusinessObjectInfo: BusinessObjectInfo{
				ID:       b.ID,
				PublicID: b.PublicID,
				RecordID: b.RecordID,
			},
			Fields: b.Fields,
		})
	}

	if bosResp.HasError {
		return nil, &CherwError{Code: bosResp.ErrorCode, Message: bosResp.ErrorMessage}
	}

	return response, nil
}
