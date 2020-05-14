package PatchingStatusSchema

import "sort"

// StateFailed state "Failed" from patching service
const StateFailed = "Failed"

// EndpointUnassigned status "Unassigned" for patching service
const EndpointUnassigned = "Unassigned"

// EndpointNoDataAvailable status "No Data Available" for patching service
const EndpointNoDataAvailable = "No Data Available"

// EndpointStatusItem endpoint status item
type EndpointStatusItem struct {
	EndpointID   string `json:"endpointID"`
	Status       string `json:"status"`
	State        string `json:"state"`
	LastUpToDate string `json:"lastUpToDate"`
	ResourceType string `json:"resourceType"`
}

// PatchingStatusData patching status response
type PatchingStatusData struct {
	Statuses   []EndpointStatusItem `json:"statuses"`
	TotalCount int                  `json:"totalCount"`
}

// By is function type of less function
type By func(p1, p2 *EndpointStatusItem) bool

// Sort sort
func (by By) Sort(patchStatus []EndpointStatusItem) {
	ps := &Sorter{patchStatus: patchStatus, by: by}
	sort.Sort(ps)
}

// Sorter combines By function and data to sort
type Sorter struct {
	patchStatus []EndpointStatusItem
	by          func(p1, p2 *EndpointStatusItem) bool
}

// Len Len gives length of data to sort
func (s *Sorter) Len() int {
	return len(s.patchStatus)
}

// Swap Swap is function to swap elements
func (s *Sorter) Swap(i, j int) {
	s.patchStatus[i], s.patchStatus[j] = s.patchStatus[j], s.patchStatus[i]
}

// Less will be called by calling By closure in the sorter
func (s *Sorter) Less(i, j int) bool {
	return s.by(&s.patchStatus[i], &s.patchStatus[j])
}
