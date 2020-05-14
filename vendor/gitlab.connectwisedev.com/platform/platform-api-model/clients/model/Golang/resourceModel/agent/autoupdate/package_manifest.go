package autoupdate

//PackageManifest is a struct defining the Package level manifest having operations list
type PackageManifest struct {
	ManifestVersion     string      `json:"manifest_version,omitempty"`
	Name                string      `json:"name,omitempty"`
	Type                string      `json:"type,omitempty"`
	Version             string      `json:"version,omitempty"`
	MinimumAgentVersion string      `json:"minimumAgentVersion,omitempty"`
	Operations          []Operation `json:"operations,omitempty"`
	UninstallOperations []Operation `json:"uninstallOperations,omitempty"`
	UnsupportedOS       []OS        `json:"unsupportedOS,omitempty"`
	Arch                []string    `json:"supportedArch,omitempty"`
	Backup              []string    `json:"backup,omitempty"`
	StallInstallation   bool        `json:"stallInstallation,omitempty"`
}

//Operation is a struct defining the operation structure
type Operation struct {
	Type                          string   `json:"type,omitempty"`
	BaseService                   bool     `json:"baseService,omitempty"`
	Action                        string   `json:"action,omitempty"`
	Name                          string   `json:"name,omitempty"`
	NewFileName                   string   `json:"newFileName,omitempty"`
	AssertSrcFileExist            bool     `json:"assertSrcFileExist,omitempty"`
	RestoreOnFailure              bool     `json:"restoreOnFailure,omitempty"`
	InstallationPath              string   `json:"installationPath,omitempty"`
	FileHash                      string   `json:"fileHash,omitempty"`
	URL                           string   `json:"url,omitempty"`
	Arguments                     []string `json:"arguments,omitempty"`
	ReadResult                    bool     `json:"readResult,omitempty"`
	Input                         string   `json:"input,omitempty"`
	FileHashType                  string   `json:"fileHashType,omitempty"`
	FileExecutionTimeoutinSeconds int      `json:"fileExecutionTimeoutinSeconds,omitempty"`
	IntegrityCheckIgnore          bool     `json:"integrityCheckIgnore,omitempty"`
	UnsupportedOS                 []OS     `json:"unsupportedOS,omitempty"`
	SupportedOnlyForOS            []OS     `json:"supportedOnlyForOS,omitempty"`
	JunoBrokerPath                string   `json:"junoBrokerPath,omitempty"`
}
