package tasking

import (
	"fmt"
	"strings"
	"time"
)

const (
	hourlyTmpl  = `%d * * * *|@every %dh`
	dailyTmpl   = `%d %d * * *|@every %dd`
	weeklyTmpl  = `%d %d * * %s|@every %dw`
	monthlyTmpl = `%d %d %s * *|@every %dm`
)

type Index int

const (
	First Index = iota
	Second
	Third
	Fourth
	Last
)

type WeekDay struct {
	Day   time.Weekday `json:"day"`
	Index Index        `json:"index"`
}

// Schedule of tasks
//
// Attributes:
//  - Regularity
// 	- TriggerTypes: represents types of triggers for current task
//  - StartRunTime: Date when cron must be activated
//  - Location: Time zone name.  Ex America/Chicago
//  - EndRunTime: Date when cron must be stopped
//  - Repeat
type Schedule struct {
	Regularity     Regularity     `json:"regularity"                   valid:"requiredForUsers"`
	TriggerTypes   []string       `json:"triggerTypes,omitempty"       valid:"optionalValidTriggerTypes"`
	TriggerFrames  []TriggerFrame `json:"triggerFrames,omitempty"      valid:"requiredForTriggers"`
	StartRunTime   time.Time      `json:"startDate"                    valid:"requiredForRecurrendAndOneTime" `
	Location       string         `json:"timeZone"                     valid:"validLocation"`
	EndRunTime     time.Time      `json:"endDate"                      valid:"optionalOnlyForRecurrent"`
	Repeat         Repeat         `json:"repeat"                       valid:"requiredOnlyForRecurrent"`
	BetweenEndTime time.Time      `json:"betweenEndTime,omitempty"     valid:"optionalBetween"`
}

// Repeat indicates how many times and when exactly task will be executed.
//
// - Attributes:
// - Every: fixed intervals of running
// - RunTime: time of execution
// - Frequency: basis for recurrent execution (Monthly, Weekly etc.)
// - DaysOfMonth: 0-31
// - DaysOfWeek  0-6
type Repeat struct {
	Every       int       `json:"every"`
	RunTime     time.Time `json:"runTime"`
	Frequency   Frequency `json:"frequency"`
	DaysOfMonth []int     `json:"daysOfMonth,omitempty"    valid:"optionalOnlyForMonthly"`
	WeekDay     *WeekDay  `json:"weekDays,omitempty"       valid:"optionalOnlyForMonthly"`
	WeekDays    []WeekDay `json:"weekDayList,omitempty"`
	DaysOfWeek  []int     `json:"daysOfWeek"               valid:"requiredOnlyForWeekly"`
	Period      int       `json:"period"                   valid:"unsettableByUsers"`
}

// TriggerSchedule - represents info about trigger schedule
type TriggerFrame struct {
	TriggerType    string    `json:"triggerType"`
	StartTimeFrame time.Time `json:"startTimeFrame,omitempty"`
	EndTimeFrame   time.Time `json:"endTimeFrame,omitempty"`
}

// String converts Schedule type object to string of cron format
func (s *Schedule) String() string {
	if s.Regularity != Recurrent {
		return ""
	}

	switch s.Repeat.Frequency {
	case Hourly:
		return fmt.Sprintf(hourlyTmpl, s.StartRunTime.Minute(), s.Repeat.Every)
	case Daily:
		return fmt.Sprintf(dailyTmpl, s.Repeat.RunTime.Minute(), s.Repeat.RunTime.Hour(), s.Repeat.Every)
	case Weekly:
		return fmt.Sprintf(weeklyTmpl, s.Repeat.RunTime.Minute(), s.Repeat.RunTime.Hour(), joinInts(s.Repeat.DaysOfWeek, ","), s.Repeat.Every)
	case Monthly:
		return fmt.Sprintf(monthlyTmpl, s.Repeat.RunTime.Minute(), s.Repeat.RunTime.Hour(), joinInts(s.Repeat.DaysOfMonth, ","), s.Repeat.Every)
	default:
		return ""
	}
}

// Cron converts Schedule type object to cron format string
func (s *Schedule) Cron() string {
	return strings.Split(s.String(), "|")[0]
}

func joinInts(ints []int, delim string) string {
	result := strings.Join(strings.Split(strings.Trim(fmt.Sprint(ints), "[]"), " "), delim)
	if result == "" {
		return "*"
	}

	return result
}
