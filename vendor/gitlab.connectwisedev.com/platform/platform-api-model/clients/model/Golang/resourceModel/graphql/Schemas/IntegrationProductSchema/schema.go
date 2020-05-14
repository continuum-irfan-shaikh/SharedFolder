package IntegrationProductSchema

// IntegrationProduct integration product type
type IntegrationProduct struct {
	ProductID   string          `json:"productId"`
	ProductName string          `json:"productName"`
	Statuses    []ProductStatus `json:"statuses"`
}

// ProductStatus product status type
type ProductStatus struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Format      string `json:"format"`
}
