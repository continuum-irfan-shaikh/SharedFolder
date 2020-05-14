package performance

import "time"

const (
	CreatedBy string = "/agent/plugin/performance"
)

//Plugin is the struct definition of /resources/performance/performancePlugin
type PerformancePlugin struct {
	TimeStampUTC             time.Time `json:"timeStampUTC"`
	Name                     string    `json:"name"`
	Type                     string    `json:"type"`
	PerformancePluginVersion string    `json:"performancePluginVersion"`
}
