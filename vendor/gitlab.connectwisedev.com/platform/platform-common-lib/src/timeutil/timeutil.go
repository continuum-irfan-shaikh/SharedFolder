package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "01/02/2006"
	timeFormat = "15:04:05"
	dayFormat  = "Monday"
)

// CurrentTime represents message format for time api
type CurrentTime struct {
	Time string `json:"Time"`
	Date string `json:"Date"`
	Day  string `json:"Day"`
}

//ToLongYYYYMMDDHH returns the time in int format as YYYYMMDDHH
func ToLongYYYYMMDDHH(tm *time.Time) (int, error) {
	return strconv.Atoi(fmt.Sprintf("%d%02d%02d%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour()))
}

//ToLongYYYYMMDD returns the time in int format as YYYYMMDD
func ToLongYYYYMMDD(tm *time.Time) (int, error) {
	return strconv.Atoi(fmt.Sprintf("%d%02d%02d", tm.Year(), tm.Month(), tm.Day()))
}

//ToHourLong returns the time array in long format for the time range given in the input
func ToHourLong(fromTime, toTime time.Time) []int {
	var tmLong []int
	toTm, _ := ToLongYYYYMMDD(&toTime) //nolint
	for tm := fromTime; ; {
		tmInt, _ := ToLongYYYYMMDD(&tm)
		if tmInt > toTm {
			break
		}
		tmLong = append(tmLong, tmInt)
		tm = tm.Add(24 * time.Hour)

	}
	return tmLong
}

//GetCurrentTime @inputParameter : locationName eg Asia/Calcutta , @Output currentTime response having day,time,Date
func GetCurrentTime(locationName string) (*CurrentTime, error) {
	now := time.Now()
	location, err := time.LoadLocation(strings.TrimSpace(locationName))
	if err != nil {
		return nil, fmt.Errorf("GetCurrentTime: Failed to load Location %v", locationName)
	}

	specifiedZoneTime := now.In(location)

	return &CurrentTime{
		Time: specifiedZoneTime.Format(timeFormat),
		Date: specifiedZoneTime.Format(dateFormat),
		Day:  specifiedZoneTime.Format(dayFormat),
	}, nil
}
