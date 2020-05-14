package common

import (
	"strings"
	"time"

	"github.com/robfig/cron"
)

// CronParser custom cron Parser without Second field
var CronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

// ValidTime is used to validate optional `end` timestamp
func ValidTime(end, start time.Time) bool {
	return !end.IsZero() && start.Before(end) || end.IsZero()
}

// ValidFutureTime is used to validate t
func ValidFutureTime(t time.Time) bool {
	return time.Now().UTC().Truncate(time.Minute).Before(t.UTC())
}

// IsElementAlreadyExists is used to check string slice for duplicates
func IsElementAlreadyExists(element string, uniqueStringSlice []string) bool {
	for _, v := range uniqueStringSlice {
		if strings.ToLower(v) == strings.ToLower(element) {
			return true
		}
	}
	return false
}
