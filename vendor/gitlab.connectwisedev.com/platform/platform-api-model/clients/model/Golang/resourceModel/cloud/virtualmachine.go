package cloud

//VMCPUPerfJobs represents backup jobs
type VMCPUPerfJobs struct {
	ResourceInfo
	VMCPUPerfJobs VMCPUPerfJob `json:"vmcpuperfjobs"`
}

//VMCPUPerfJob represents single vm cpu performance job
type VMCPUPerfJob struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	TimeSeries TimeSeries `json:"azureIaaSVMErrorInfo"`
}

//TimeSeries represent performance Data
type TimeSeries struct {
	Data PerfData `json:"data"`
}

//PerfData represnt cpu performance details
type PerfData struct {
	TimeStamp string `json:"timeStamp"`
	Average   string `json:"average"`
}
