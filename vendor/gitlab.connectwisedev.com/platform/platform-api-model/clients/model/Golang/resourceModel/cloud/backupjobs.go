package cloud

//BackupJobs represents backup jobs
type BackupJobs struct {
	ResourceInfo
	BackupJobs BackupJob `json:"vaultbackupjobs"`
}

//BackupJob represents single backup jobs
type BackupJob struct {
	JobType   string       `json:"jobType"`
	VMName    string       `json:"entityFriendlyName"`
	StartTime string       `json:"startTime"`
	EndTime   string       `json:"endTime"`
	Status    string       `json:"status"`
	ErrorInfo ErrorDetails `json:"azureIaaSVMErrorInfo"`
}

//ErrorDetails represent Error Info
type ErrorDetails struct {
	ErrorCode   int64  `json:"errorCode"`
	ErrorString string `json:"errorString"`
}
