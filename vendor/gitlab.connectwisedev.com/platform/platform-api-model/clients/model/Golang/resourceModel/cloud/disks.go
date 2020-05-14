package cloud

//Disks represents status of Azure Managed Disks
type Disks struct {
	Disks []Disk
}

//Disk represents status of Disk
type Disk struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ResourceType  string `json:"type"`
	Location      string `json:"location"`
	State         string `json:"diskState"`
	DiskSizeBytes int64  `json:"diskSizeBytes"`
}
