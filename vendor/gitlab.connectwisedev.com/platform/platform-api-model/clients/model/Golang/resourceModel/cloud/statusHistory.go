package cloud

import "time"

//ResourceStatusHistory return resource status history
type ResourceStatusHistory struct {
	Summary            string    `json:"summary"`
	AvailabilityStatus string    `json:"availabilitystatus"`
	OccurrenceTimeUTC  time.Time `json:"occurrencetimeutc"`
}
