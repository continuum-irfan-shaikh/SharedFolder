package gateway

// Statistics is struct for reporting statistics
type Statistics struct {
	DownloadFolderSizeInBytes int              `json:"downloadFolderSizeInBytes,omitempty"`
	HostedPackages            []PackageDetails `json:"hostedPackages,omitempty"`
	BlacklistedPackages       []PackageDetails `json:"blacklistedPackages,omitempty"`
}
