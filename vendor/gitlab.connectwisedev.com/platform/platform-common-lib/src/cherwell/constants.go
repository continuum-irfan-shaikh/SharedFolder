package cherwell

// Group of constants related to BusinessObject swagger tag
// it provides uri patterns for formatting with required values
const (
	createUpdateBOEndpoint        = "/api/V1/savebusinessobject"
	createUpdateRelatedBOEndpoint = "/api/V1/saverelatedbusinessobject"
	createUpdateBOsEndpoint       = "/api/V1/savebusinessobjectbatch"
	deleteBOByPubIDEndpoint       = "/api/V1/deletebusinessobject/busobid/%v/publicid/%v"
	deleteBOByRecIDEndpoint       = "/api/V1/deletebusinessobject/busobid/%v/busobrecid/%v"
	deleteBOsEndpoint             = "/api/V1/deletebusinessobjectbatch"
	getBOByPubIDEndpoint          = "/api/V1/getbusinessobject/busobid/%v/publicid/%v"
	getBOByRecIDEndpoint          = "/api/V1/getbusinessobject/busobid/%v/busobrecid/%v"
	getBOsEndpoint                = "/api/V1/getbusinessobjectbatch"
	getRelatedObjectsEndpoint     = "/api/V1/getrelatedbusinessobject"
	retryCount                    = 3

	// searchEndpoint is a constant declared with uri for performing search requests
	searchEndpoint = "/api/V1/getsearchresults"

	// fieldValuesLookupEndpoint is a constant declared with uri for performing lookup requests
	fieldValuesLookupEndpoint = "/api/V1/fieldvalueslookup"

	// tokenEndpoint is a constant declared with uri for token obtaining
	tokenEndpoint = "/token"

	// passwordGrantType is a grant type value to be sent in token refresh request
	passwordGrantType = "password"

	attachmentUploadPath         = "/api/V1/uploadbusinessobjectattachment"
	attachmentGetPath            = "/api/V1/getbusinessobjectattachments/busobid/%s/busobrecid/%s/type/%s/attachmenttype/%s"
	attachmentUploadPathTemplate = "/filename/%s/busobid/%s/busobrecid/%s/offset/%d/totalsize/%s"
	attachmentDeletePath         = "/api/V1/removebusinessobjectattachment/attachmentid/%s/busobid/%s/busobrecid/%s"
	attachmentDownloadPath       = "/api/V1/getbusinessobjectattachment/attachmentid/%s/busobid/%s/busobrecid/%s"

	linkBOsPath   = "/api/V1/linkrelatedbusinessobject/parentbusobid/%s/parentbusobrecid/%s/relationshipid/%s/busobid/%s/busobrecid/%s"
	unlinkBOsPath = "/api/V1/unlinkrelatedbusinessobject/parentbusobid/%s/parentbusobrecid/%s/relationshipid/%s/busobid/%s/busobrecid/%s"
)
