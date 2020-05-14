package asset

// AssetCardReader defines the structure of all cardreader related info
type AssetCardReader struct {
	VendorID  string `json:"vendorID,omitempty" cql:"vendor_id"`
	ProductID string `json:"productID,omitempty" cql:"product_id"`
	Revision  string `json:"revision,omitempty" cql:"revision"`
}
