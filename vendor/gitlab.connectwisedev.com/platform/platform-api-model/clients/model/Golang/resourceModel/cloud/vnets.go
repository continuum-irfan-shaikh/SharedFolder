package cloud

//Vnet represents Virtual Network structure
type Vnet struct {
	ResourceInfo
	Peerings []VnetPeering `json:"virtualNetworkPeerings"`
}

//VnetPeering represents single vnet peering
type VnetPeering struct {
	PeeringName  string `json:"peeringName"`
	PeeringState string `json:"peeringState"`
}
