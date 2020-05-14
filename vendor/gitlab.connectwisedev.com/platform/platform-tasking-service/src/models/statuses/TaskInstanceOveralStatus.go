package statuses

// OverallStatus represents overal status of an instance
type OverallStatus string

// represents all overall statuses
const (
	OverallSuccess       OverallStatus = "Success"
	OverallFailed        OverallStatus = "Failed"
	OverallPartialFailed OverallStatus = "Partial Failed"
	OverallNew           OverallStatus = "New"
	OverallScheduled     OverallStatus = "Scheduled"
	OverallSuspended     OverallStatus = "Suspended"
	OverallRunning       OverallStatus = "Running"
)
