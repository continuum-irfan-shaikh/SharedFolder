package autoupdate

const (
	//InstallOperation ...
	InstallOperation = "install"
	//UninstallOperation ...
	UninstallOperation = "uninstall"
	//UpdateOperation ...
	UpdateOperation = "update"
)

const (
	//ContinuumOSName : Operating system name header
	ContinuumOSName = "ITSPlatform-OS-Name"
	//ContinuumOSType : Operating system type header
	ContinuumOSType = "ITSPlatform-OS-Type"
	//ContinuumOSVersion : Operating System Version header
	ContinuumOSVersion = "ITSPlatform-OS-Version"
	//ContinuumArchitecture : system architecture header
	ContinuumArchitecture = "ITSPlatform-Architecture"
)

//GlobalManifest is a struct defining the Global manifest having package list
type GlobalManifest struct {
	ProductVersion string    `json:"productVersion,omitempty"`
	Version        string    `json:"version,omitempty"`
	SupportedOS    []OS      `json:"supportedOS,omitempty"`
	SupportedArch  []string  `json:"supportedArch,omitempty"`
	Packages       []Package `json:"packages,omitempty"`
}

//OS is a struct defining the Operating System Details
type OS struct {
	Name    string `json:"name,omitempty" cql:"name"`
	Type    string `json:"type,omitempty" cql:"type"`
	Version string `json:"version,omitempty" cql:"version"`
}

//Package is a struct defining the Package Details
type Package struct {
	Name            string     `json:"name,omitempty" cql:"name"`
	Type            string     `json:"type,omitempty" cql:"type"`
	OperatingSystem string     `json:"operating_system,omitempty" cql:"operating_system"`
	Version         string     `json:"version,omitempty" cql:"version"`
	SourceURL       string     `json:"sourceUrl,omitempty" cql:"source_url"` // exists for backward compatibility can be removed after all agents start supporting URLs
	URLs            []URL      `json:"urls,omitempty" cql:"urls"`            // supporting multiple download paths, each with multiple protocols/schemes if applicable
	Operation       string     `json:"operation,omitempty"`
	Arch            string     `json:"arch,omitempty" cql:"arch"`
	Checksum        []Checksum `json:"checksum,omitempty" cql:"checksum"`
	Cacheable       bool       `json:"cacheable,omitempty" cql:"cacheable"`
}

// URL is a struct holding supported schemes and base url for a specific resource path
type URL struct {
	Schemes      []string `json:"schemes,omitempty" cql:"schemes"`            // scheme/protocol i.e. "http", "https", "ftp", etc
	BaseURL      string   `json:"baseUrl,omitempty" cql:"base_url"`           // domain name i.e. "domain.com"
	ResourcePath string   `json:"resourcePath,omitempty" cql:"resource_path"` // resource path suffixed with a '/' i.e. "/location/resource.ext"
}

// Checksum is a struct defining the Checksum Details
type Checksum struct {
	Type  string `json:"type,omitempty" cql:"type"`
	Value string `json:"value,omitempty" cql:"value"`
}
