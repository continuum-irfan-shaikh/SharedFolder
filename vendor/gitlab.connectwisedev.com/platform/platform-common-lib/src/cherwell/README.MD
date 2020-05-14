# Cherwell Service Manager Client

Go client for Cherwell Service Management (CSM) 

## NOTE:
- This should not be part of common-lib. It should be added in the glide as the dependency. 
- Going forward, the plan is to façade that can be used by relevant microservice.
- Cherwell apis will not be used by any microservice except ITSM integration one.


## Documentation
For details on functionality, see [Cherwell API Golang library](https://continuum.atlassian.net/wiki/spaces/C2E/pages/781945314/Cherwell+API+Golang+library)

### Setup
How to create a client:
```go
import "gitlab.connectwisedev.com/platform/platform-common-lib/src/cherwell"
        
var conf = cherwell.Config{
                        Host    : "CSM_HOST",
                        ClientId: "CSM_CLIENT_ID",
                        UserName: "CSM_USERNAME",
                        PassWord: "CSM_PASSWORD"
                        }
cl, err := cherwell.NewClient(conf, &http.Client{})
 	if err != nil {
 		// handle error
 	}
```
#### Client purpose
Cherwell client provides calls for interacting with objects in Cherwell.
It also contains structure representations of:
- Cherwell Business Object 
- Attachments 
- Linked Objects 
- Request/Response  for Save/Get/Update/Delete/Search operations 
- Request/Response for Attachments
- Request/Response for Related Business Objects

#### Entities representation
- Cherwell Business Object  
    BusinessObject is a structure which represents all Business Object in Cherwell API
    ```go
    type BusinessObject struct {
	    BusinessObjectInfo
	    Fields []FieldTemplateItem 
    }
    ```
    
    BusinessObjectInfo represents request to Cherwell API for reading and deleting of Business Object
    ```go
    type BusinessObjectInfo struct {
        ID       string 
        PublicID string 
        RecordID string 
        ErrorData
    }
    ```
    FieldTemplateItem represents structure of the field for Business Object
    ```go
    type FieldTemplateItem struct {
        Dirty       bool  
        DisplayName string 
        FieldID     string 
        HTML        string 
        Name        string 
        Value       string 
    }
    ```
- Attachments  
    Attachment represents Cherwell BO attachment
    ```go
    type Attachment struct {
        Owner        BusinessObjectInfo
        FileName     string
        Description  string
        Data         io.Reader
        Offset       int
        AttachmentID string
    }
    ```
    AttachmentSummary contains information about attachment
    ```go
    type AttachmentSummary struct {
	    BusinessObjectInfo
	    FileID         string 
	    FileName       string 
	    FileType       string 
	    AttachmentID   string 
	    AttachmentType int    
	    Comment        string 
	    CreatedAt      string 
	    DisplayText    string 
	    Owner          string 
	    Scope          int    
	    ScopeOwner     string 
	    Type           int    
	    Links          []struct {
	        Name string 
	        URL  string 
	    } 
    }
    ```
    
    AttachedFile is downloaded attachment file
    ```go
    type AttachedFile struct {
        FileName    string
        ContentType string
        SizeBytes   string
        Data        io.ReadCloser
    }
    ```
- Linked Object  
    ```go
    type LinkedObject struct {
        ParentBusObID    string
        ParentBusObRecID string
        RelationshipID   string
        BusObID          string
        BusObRecID       string
    }   
    ```
- Related Business Objects Request/Response  
    RelatedBusinessObjectsRequest represents request to Cherwell API for retrieving business objects related to parent
    ```go
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
    ```
    RelatedBusinessObjectsResponse represents response from Cherwell API containing business objects related to parent
    ```go
    type RelatedBusinessObjectsResponse struct {
        Responses    []businessObjectResponse 
        TotalRecords int                     
    }
    ```
- Search    
    SearchRequest request is a structure representing request for endpoint "/api/V1/getsearchresults"
    ```go 
    type SearchRequest struct {
        Association        string       
        ID                 string        
        CustomGridDefID    string         
        DateTimeFormatting string        
        FieldID            string        
        Scope              string       
        ScopeOwner         string         
        SearchID           string        
        SearchName         string        
        SearchText         string       
        IncludeAllFields   bool          
        IncludeSchema      bool         
        PageNumber         int            
        PageSize           int           
        Sorting            []*Sorting   
        PromptValues       []*PromptValue 
        Fields             []string      
        Filters            []*Filter      
    }
    ```
    SearchResponse is a structure representing response for endpoint "/api/V1/getsearchresults"  
    ```go
    type SearchResponse struct {
	    BusinessObjects []BusinessObject
	    TotalRows       int
	    ErrorData
    }
    ```

    Sorting is a structure representing sorting for search request
    ```go
    type Sorting struct {
        FieldID       string 
	    SortDirection int   
    }
    ```
    PromptValue is a structure representing prompt value for search request
    ```go
    type PromptValue struct {
	    ID                       string     
	    CollectionStoreEntireRow string      
	    CollectionValueField     string      
	    FieldID                  string     
	    ListReturnFieldID        string     
	    PromptID                 string      
	    Value                    interface{}
	    ValueIsRecID             bool        
    }   
    ```
    Filter is a structure representing rules for filtering by business object value
    ```go
    type Filter struct {
        FieldID  string 
        Operator string
        Value    string 
    }
    ```
- Errors  
    ErrorData holds common error info in response
    ```go
    type ErrorData struct {
	    ErrorCode    string 
	    ErrorMessage string 
	    HasError     bool  
    }
    ```

    Errors is a set of errors
    ```go
    type Errors struct {
	    Errors []error
    }
    ```

    BusObNotValid recognizes cherwell invalid business object error
    ```go
    type BusObNotValid struct {
	    Message string
    }
    ```

    RecordNotFound recognizes cherwell record not found error
    ```go
    type RecordNotFound struct {
	    Message string
    }
    ```

    GeneralFailure recognizes cherwell general failure error
    ```go
    type GeneralFailure struct {
	    Message string
    }
    ```

    InvalidFilterOperator recognizes cherwell invalid filter operator error
    ```go
    type InvalidFilterOperator struct {
	    Message string
    }
    ```

    FieldValidationError represents errors (about Business Object's fields) which occurs after sending invalid request to Cherwell API filter operator error
    ```go
    type FieldValidationError struct {
	    Error     string 
	    ErrorCode string 
	    FieldID   string 
    }
    ```

#### Reading Business Objects from CSM
To get information about Business Objects from CSM use following methods:

- `func (c *Client) GetByPublicID(busobID, publicID string) (*BusinessObject, error)` - 
method that returns a Business Object record that includes a list of fields and their record IDs, names, and set values.
Specify the Business Object ID and Business Object public ID.
- `func (c *Client) GetByRecordID(busobID, busobRecID string) (*BusinessObject, error)` - 
method that returns a Business Object record that includes a list of fields and their record IDs, 
names, and set values. Specify the Business Object ID and Business Object record ID.
- `func (c *Client) GetBatch(items []BusinessObjectInfo) ([]BusinessObject, error)` - 
method that returns a batch of Business Object records that includes a list of field record IDs, display names, 
and values for each record. Specify an array of Business Object IDs, record IDs, or public IDs ([]BusinessObjectInfo). 
Use a flag to stop on error or continue on error.

#### Creating Business Objects in CSM

To create Business Objects in CSM use following methods:

- `func (c *Client) Save(bo BusinessObject) (string, string, error)` - to create Business Object, 
leave record ID and public ID empty. 
Upon creating a cache key is returned to use for subsequent requests. If the object is not found in the cache 
with said cache key, specify record ID or public ID to save and return a new cache key. Set persist = true, to actually 
save the Business Object to disk, persist = false will just cache it.
- `func (c *Client) SaveBatch(bos []BusinessObject) ([]string, []string, error)` - to create an array of Business 
Objects in a batch, leave record ID and public ID empty.

#### Updating Business Objects in CSM

To update Business Objects in CSM use the same methods like for creating:

- `func (c *Client) Save(bo BusinessObject) (string, string, error)` - to update Business Object specify
Business Object ID and record ID or public ID.
- `func (c *Client) SaveBatch(bos []BusinessObject) ([]string, []string, error)` - to update Business Objects 
specify Business Object IDs and record IDs or public IDs specify record ID or public ID.

#### Deleting Business Objects from CSM

To delete Business Objects from CSM use following methods:

- `func (c *Client) DeleteByPublicID(busobID, publicID string) (string, string, error)` - 
method that deletes a Business Object record by public ID. Specify the Business Object ID and Business Object public ID.
- `func (c *Client) DeleteByRecordID(busobID, busobRecID string) (string, string, error)` - 
method that deletes a Business Object record by Record ID. Specify the Business Object ID and Business Object record ID.
- `func (c *Client) DeleteBatch(itemsToDelete []BusinessObjectInfo) ([]string, []string, error)` - 
method that deletes a batch of Business Object records. Specify an array of Business Object IDs, record IDs
or public IDs. Use a flag to stop on error or continue on error.

#### Search Business Objects 

To find Business Objects use following methods:
- `func NewSearchRequest(id string) *SearchRequest` -
function that creates new SearchRequest with ID with specific fields in response
- `func (sc *SearchRequest) AddFilter(filedID, operator, value string) error` -
method that adds new filter to SearchRequest
- `func (sc *SearchRequest) SetSpecificFields(requiredFieldIDs []string)` -
method that sets specific fields to be returned in Business Object response
- `func (sc *SearchRequest) AppendSpecificFields(requiredFieldIDs ...string)` -
method that appends specific fields to be returned in Business Object response
- `func (c *Client) FindBoInfos(req SearchRequest) (*SearchResponse, error)` -
method that sets SearchRequest for returning only BoInfos in response
- `func (c *Client) Find(req SearchRequest) (*SearchResponse, error)` -
method that returns all business objects defined by SearchRequest

#### Link/Unlink Business Object
To link and unlink Business Objects use following methods:
- `func (c *Client) Link(l *LinkedObject) error` -
method that links related Business Objects
- `func (c *Client) Unlink(l *LinkedObject) error` -
method that unlinks related Business Objects
- `func (c *Client) GetRelatedBusinessObjects(req *RelatedBusinessObjectsRequest) ([]BusinessObject, error)` -
method that retrieves a bucket of business objects related to parent

#### Attachments
- `func NewAttachment(fileName string, data io.Reader, owner *BusinessObjectInfo) *Attachment` -
function that creates a new attachment object
- `func (c *Client) UploadAttachment(a *Attachment) (fileID string, err error)` -
method that uploads a file as attachment to business object
- `func (c *Client) GetObjectAttachments(boid, recordID string, recordType RecordAttachmentType, attachmentType AttachmentType) (resp *AttachmentResponse, err error)` -
method that gets attachments for specific business object
- `func (c *Client) DeleteAttachment(attachmentID string, owner *BusinessObjectInfo) error` -
method that removes an attachment
- `func (c *Client) AttachmentByID(attachmentID string, owner *BusinessObjectInfo) (*AttachedFile, error)` - 
method that gets attachment contents by ID  

#### Errors
- `func (r *ErrorData) GetErrorObject() error` -
method that gets error object based on error response
- `func NewErrorSet(es ...error) *Errors` -
function that creates a new Errors instance
- `func (es *Errors) Add(e ...error)` -
method that appends an error or set of errors to existing set
- `func (es *Errors) IsEmpty() bool` -
method that checks if set of errors is not empty

### Examples
- Create Business Object:  
bo - variable which stores information about BusinessObject
```go
    busObPublicID, busObRecID, err := cl.Save(bo)
    if err != nil { 
            // handle error
    }	
```
- Create Business Objects in batch:  
bos - variable which stores information about BusinessObjects
```go

    publicIds, recIds, err := cl.SaveBatch(bos)
    if err != nil {
        // handle error
    }
```
- Read Business Object by Public ID:
```go
    bo, err := cl.GetByPublicID(busobID, publicID)
    if err != nil {
		// handle error
    }
```
- Read Business Object by Record ID:
```go
    bo, err := cl.GetByRecordID(busobID, recID)
    if err != nil {
		// handle error
    }
```

- Read Business Objects in batch:
```go
    var boInfo []BusinessObjectInfo
    // set up particular boInfo
    bos, err := cl.GetBatch(boInfo)
    if err != nil {
		// handle error
    }
```

- Delete Business Object by Public ID:
```go
    bo, err := cl.DeleteByPublicID(busobID, publicID)
    if err != nil {
		// handle error
    }
```

- Delete Business Object by Record ID:
```go
    bo, err := cl.DeleteByRecordID(busobID, recID)
    if err != nil {
		// handle error
    }
```

- Delete Business Objects in batch:
```go
    var itemsToDelete []BusinessObjectInfo
    pubIDs, recIDs, err := cl.DeleteBatch(itemsToDelete)
    if err != nil {
		// handle error
    }
```
- Update Business Object (the same method as for saving of Business Object):  
bo - variable which stores information about BusinessObject with fields to update
```go
    busObPublicID, busObRecID, err := cl.Save(bo)
    if err != nil { 
            // handle error
    }	
```

- Update Business Objects in batch (the same method as for saving of Business Objects):  
bos - variable which stores information about BusinessObjects with fields to update
```go
    publicIds, recIds, err := cl.SaveBatch(bos)
    if err != nil {
        // handle error
    }
```

- Search Business Objects
busObID - ID of Business Object in Cherwell (Incident, Contact, etc.)  
fieldID - ID of Business Object field (title, description, etc.)  
operator - search operator (eq, lt, gt, contains, startswith)  
value - Business Object field value  
```go
    request = cl.NewSearchRequest(busObID)
    err = request.AddFilter(fieldID, operator, value)
    if err != nil {
        // handle error
    }
    response, err := cl.Find(*request)
    if err != nil {
        // handle error
    }
```
- Link Business Objects  
linkage - LinkedObject
```go
    err = cl.Link(&linkage)
    if err != nil {
        // handle error
    }
```
- Unlink Business Objects  
linkage - LinkedObject
```go
    err = cl.Unlink(&linkage)
    if err != nil {
        // handle error
    }
```

- Get Related Business Objects  
request - RelatedBusinessObjectsRequest
result - recieved related Business Objects
```go
    result, err = cl.GetRelatedBusinessObjects(&request)
    if err != nil {
        // handle error
    }
```

- Create new Attachment  
fileName - name of the file  
owner - Business Object Info of the file owner  
```go
    attachment := NewAttachment(fileName, file, owner)
```
- Upload Attachment  
fileID - ID of attached file
```go
    fileID, err := cl.UploadAttachment(attachment)
    if err != nil {
        // handle error
    }
```
- Get object Attachments  
ownerID - owner business object ID  
ownerRecordID - owner record ID  
fileRecord - record attachment type (File, History, etc.)  
importedAttachment - attachment type (Imported, Linked, etc.)  
response - AttachmentResponse  
```go
    response, err := cl.GetObjectAttachments(
        ownerID,
        ownerRecordID,
        fileRecord,
        importedAttachment,
    )
    if err != nil {
        // handle error
    }
```

- Delete Attachment  
atID - ID of attachment  
owner - Business Object Info of the attachment owner  
```go
    err = cl.DeleteAttachment(atID, &owner)
    if err != nil {
        // handle error
    }
```

- Get Attachment by ID  
atID - ID of attachment  
owner - Business Object Info of the attachment owner  
```go
    result, err := cl.AttachmentByID(atID, &owner)
    if err != nil {
        // handle error
    }
```

- Get Error object    
```go
    err = resp.GetErrorObject()
    return err
```

- Create new ErrorSet 
```go
    errs := NewErrorSet()
```

- Add error to ErrorSet    
```go
    errs.Add(&CherwellError)
```

- Check if ErrorSet is empty    
```go
    if !errs.IsEmpty() {
		return nil, errs
    }
```
