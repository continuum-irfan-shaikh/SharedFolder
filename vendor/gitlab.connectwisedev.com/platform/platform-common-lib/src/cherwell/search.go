package cherwell

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	// A bucket of possible filtering operations to be used for check if object field value satisfies operation condition applied to value provided in filter.

	// OpEqual is operation for value equality check.
	OpEqual = "eq"
	// OpLessThan checks if field of searched object is less than field in filter. Cherwell perform simple string comparsion, so digits and dates comparsion will not eork properly in most of cases.
	OpLessThan = "lt"
	// OpGreaterThan checks if field of searched object is greater than field in filter. Cherwell perform simple string comparsion, so digits and dates comparsion will not eork properly in most of cases.
	OpGreaterThan = "gt"
	// OpContains checks if field value contains filter value.
	OpContains = "contains"
	// OpStartsWith checks if field value has suffix equal to filter value.
	OpStartsWith = "startswith"

	// SortDesc direction
	SortDesc = 0
	// SortAsc direction
	SortAsc = 1
)

// PromptValue is a structure representing prompt value for search request
type PromptValue struct {
	ID                       string      `json:"busObId,omitempty"`
	CollectionStoreEntireRow string      `json:"collectionStoreEntireRow,omitempty"`
	CollectionValueField     string      `json:"collectionValueField,omitempty"`
	FieldID                  string      `json:"fieldId,omitempty"`
	ListReturnFieldID        string      `json:"listReturnFieldId,omitempty"`
	PromptID                 string      `json:"promptId,omitempty"`
	Value                    interface{} `json:"value,omitempty"`
	ValueIsRecID             bool        `json:"valueIsRecId,omitempty"`
}

// Sorting is a structure representing sorting for search request
type Sorting struct {
	FieldID       string `json:"fieldId,omitempty"`
	SortDirection int    `json:"sortDirection,omitempty"`
}

// Filter is a structure representing rules for filtering by business object value
type Filter struct {
	FieldID  string `json:"fieldId"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// SearchRequest request is a structure representing request for endpoint "/api/V1/getsearchresults"
type SearchRequest struct {
	Association        string        `json:"association,omitempty"`
	ID                 string        `json:"busObId,omitempty"`
	CustomGridDefID    string        `json:"customGridDefId,omitempty"`
	DateTimeFormatting string        `json:"dateTimeFormatting,omitempty"`
	FieldID            string        `json:"fieldId,omitempty"`
	Scope              string        `json:"scope,omitempty"`
	ScopeOwner         string        `json:"scopeOwner,omitempty"`
	SearchID           string        `json:"searchId,omitempty"`
	SearchName         string        `json:"searchName,omitempty"`
	SearchText         string        `json:"searchText,omitempty"`
	IncludeAllFields   bool          `json:"includeAllFields"`
	IncludeSchema      bool          `json:"includeSchema,omitempty"`
	PageNumber         int           `json:"pageNumber,omitempty"`
	PageSize           int           `json:"pageSize,omitempty"`
	Sorting            []Sorting     `json:"sorting,omitempty"`
	PromptValues       []PromptValue `json:"promptValues,omitempty"`
	Fields             []string      `json:"fields,omitempty"`
	Filters            []Filter      `json:"filters,omitempty"`
}

// SearchResponse is a structure representing response for endpoint "/api/V1/getsearchresults"
type SearchResponse struct {
	BusinessObjects []BusinessObject
	TotalRows       int
	ErrorData
}

type searchResultsResponse struct {
	Responses []*businessObjectResponse `json:"businessObjects"`
	TotalRows int                       `json:"totalRows"`
	ErrorData
}

var filterOperators = map[string]struct{}{
	"eq":         {},
	"lt":         {},
	"gt":         {},
	"contains":   {},
	"startswith": {},
}

// NewSearchRequest creates new SearchRequest with ID with specific fields in response
func NewSearchRequest(id string) *SearchRequest {
	return &SearchRequest{ID: id, Filters: []Filter{}, IncludeAllFields: true}
}

// AddFilter adds new filter to SearchRequest
func (sc *SearchRequest) AddFilter(filedID, operator, value string) error {
	_, ok := filterOperators[operator]
	if !ok {
		return InvalidFilterOperator{Message: "invalid filter operator"}
	}
	f := Filter{
		FieldID:  filedID,
		Operator: operator,
		Value:    value,
	}
	sc.Filters = append(sc.Filters, f)

	return nil
}

// OrderBy adds ordering rule to search request
func (sc *SearchRequest) OrderBy(filedID string, direction int) error {
	if direction != SortAsc && direction != SortDesc {
		return errors.New("adding order by: invalid direction value '%s")
	}
	sc.Sorting = append(sc.Sorting, Sorting{FieldID: filedID, SortDirection: direction})
	return nil
}

// SetSpecificFields sets specific fields to be returned in BO response
func (sc *SearchRequest) SetSpecificFields(requiredFieldIDs []string) {
	if len(requiredFieldIDs) == 0 {
		return
	}
	sc.IncludeAllFields = false
	sc.Fields = requiredFieldIDs
}

// AppendSpecificFields appends specific fields to be returned in BO response
func (sc *SearchRequest) AppendSpecificFields(requiredFieldIDs ...string) {
	if len(requiredFieldIDs) == 0 {
		return
	}
	sc.IncludeAllFields = false
	sc.Fields = append(sc.Fields, requiredFieldIDs...)
}

// FindBoInfos sets SearchRequest for returning only BoInfos in response
func (c *Client) FindBoInfos(req SearchRequest) (*SearchResponse, error) {
	input := SearchRequest{ID: req.ID, Filters: req.Filters, IncludeAllFields: false, Fields: []string{""}}

	return c.Find(input)
}

// Find returns all business objects defined by SearchRequest
func (c *Client) Find(req SearchRequest) (*SearchResponse, error) {
	if req.PageNumber < 0 || req.PageSize < 0 {
		return nil, &CherwError{
			Code:    BadRequestError,
			Message: fmt.Sprintf("Page Number/Page Size is invalid. PageNumber: %d, PageSize: %d", req.PageNumber, req.PageSize),
		}
	}

	var resp searchResultsResponse

	// calculate page number if it is greater then 0,
	// PageNumber works as "offset", so we need to recalculate it on fly to align the logic with Cherwell Search API
	if req.PageNumber > 0 && req.PageSize > 0 {
		req.PageNumber = req.PageNumber*req.PageSize - req.PageSize + 1
	}

	err := c.performRequest(http.MethodPost, searchEndpoint, &req, &resp)
	if err != nil {
		return nil, fmt.Errorf("find: %s", err)
	}

	if resp.HasError {
		err = resp.GetErrorObject()
		return nil, err
	}

	errs := NewErrorSet()

	var bos []BusinessObject
	var searchResponse SearchResponse
	for _, b := range resp.Responses {
		if b.HasError {
			errs.Add(&CherwError{
				Code: b.ErrorCode, Message: b.ErrorMessage,
			})
		}
		bos = append(bos, BusinessObject{
			BusinessObjectInfo: BusinessObjectInfo{ID: b.ID, PublicID: b.PublicID, RecordID: b.RecordID},
			Fields:             b.Fields,
		})
	}

	if !errs.IsEmpty() {
		return nil, errs
	}

	searchResponse.BusinessObjects = bos
	searchResponse.TotalRows = resp.TotalRows
	searchResponse.ErrorData = resp.ErrorData

	return &searchResponse, nil
}
