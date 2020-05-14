package cloud

//Hierarchy represents an hierarchy of the cloud vendors product suite
type Hierarchy struct {
	Title  string           `json:"title"`
	Level  int              `json:"level"`
	Values []HierarchyValue `json:"values"`
}

//HierarchyValue represents a single value in the list of values for a hierarchy
type HierarchyValue struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

//HierarchySummary is used to display the hierarchy information along with other information
type HierarchySummary struct {
	Title string                  `json:"title"`
	Level int                     `json:"level"`
	Value []HierarchySummaryValue `json:"values"`
}

//HierarchySummaryValue is used to display the hierarchy summary value
type HierarchySummaryValue struct {
	VendorID   string          `json:"vendorid"`
	Name       string          `json:"name"`
	State      string          `json:"state"`
	Properties ValueProperties `json:"properties"`
}

//ValueProperties additional properties of a hierarchy
type ValueProperties struct {
	LocationPlacementID string `json:"locationplacementid"` //Azure Policy value
	QuotaID             string `json:"quotaid"`             //Azure Policy value
	SpendingLimit       string `json:"spendinglimit"`       //Azure Policy value
}

//MapHierarchy this is for posting the Hierarchy information
type MapHierarchy struct {
	Hierarchies []string `json:"Hierarchies"`
}
