package cherwell

import "net/http"

const boTemplateURL = "/api/V1/getbusinessobjecttemplate"

// BOTemplateQuery contains params for querying business object templates
type BOTemplateQuery struct {
	// ID is a business object ID
	ID string `json:"busObId"`

	// FieldNames is a list of fields to get by field names
	FieldNames []string `json:"fieldNames"`

	// FieldIDs is a list of fields to get by field IDs
	FieldIDs []string `json:"fieldIds"`

	// IncludeAll flag includes all business object fields
	IncludeAll bool `json:"includeAll"`

	// IncludeRequired flag includes required fields
	IncludeRequired bool `json:"includeRequired"`
}

// BOTemplate is business object template
type BOTemplate struct {
	ErrorData

	// ID is a business object ID
	ID string

	// Fields is a list of business object fields
	Fields []FieldTemplateItem `json:"fields"`
}

// NewBusinessObject creates a new business object from the template
func (t *BOTemplate) NewBusinessObject() BusinessObject {
	bo := BusinessObject{
		BusinessObjectInfo: BusinessObjectInfo{
			ID: t.ID,
		},
		Fields: make([]FieldTemplateItem, len(t.Fields)),
	}

	copy(bo.Fields, t.Fields)
	return bo
}

// FieldByName gets field from the template by name
func (t *BOTemplate) FieldByName(name string) (*FieldTemplateItem, bool) {
	for index, field := range t.Fields {
		if field.Name == name {
			return &t.Fields[index], true
		}
	}

	return nil, false
}

// GetBusinessObjectTemplate returns a template to create Business Objects.
//
// The template includes placeholders for field values.
//
// You can then send the template with these values to the Business Object Save operation.
func (c *Client) GetBusinessObjectTemplate(query BOTemplateQuery) (template *BOTemplate, err error) {
	template = new(BOTemplate)
	err = c.performRequest(http.MethodPost, boTemplateURL, query, template)
	if err != nil {
		return nil, err
	}

	if template.HasError {
		return nil, template.GetErrorObject()
	}

	template.ID = query.ID
	return template, nil
}
