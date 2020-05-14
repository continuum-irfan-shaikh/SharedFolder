package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	transactionID "gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/grsmv/goweek"
	"github.com/pkg/errors"
	"github.com/robfig/cron"

	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	commonlibRest "gitlab.connectwisedev.com/platform/platform-common-lib/src/web/rest"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/translation"
)

const (
	// StartIndexDefault defines a default value for Start Index
	StartIndexDefault = 0
	// NumberOfRowsDefault defines a default value for Number of Rows
	NumberOfRowsDefault = 100
	// NumberOfRowsMaximum defines a default value for Max Number of Rows
	NumberOfRowsMaximum  = 1000
	acceptLanguageHeader = "Accept-Language"
	// InitiatedByHeader defines HTTP header with user account UUID.
	InitiatedByHeader = "uid"
	monthInYear       = 12
	daysInWeek        = 7
	firstDay          = 1
	minDaysInMonth    = 28
	countKey          = "count"
	// UnlimitedCount is used when we don't need to LIMIT result rows in cql query
	UnlimitedCount = 0
	// CassandraTimeFormat defines format for Cassandra timestamps
	CassandraTimeFormat = `2006-01-02 15:04:05-0700`

	// ContentType is a key of request header
	ContentType = "Content-Type"
	// ApplicationJSON is a value of Content-Type request header
	ApplicationJSON = "application/json"

	// NextRunTimeExceedsEndRunTime  describe case when end run time is expired
	NextRunTimeExceedsEndRunTime = "next run time exceeds scheduled period"

	taskWillNeverRunError = "task will never run"
)

var (
	major          = "1"
	minor          = "0"
	serviceVersion = "v" + major + "." + minor
	// hourlyPattern  = regexp.MustCompile(`^@every\s*[1-9][0-9]*[h]$`)
)

// Message contains the message (for example about error)
type Message struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode,omitempty"`
}

// RunTimeInvalidError is a custom error
type RunTimeInvalidError struct {
	ErrMsg       string
	StartRunTime time.Time
	EndRunTime   time.Time
	Schedule     apiModels.Schedule
}

func (e *RunTimeInvalidError) Error() string {
	return fmt.Sprintf("%s. StartRunTime[%v], EndRunTime[%v], Schedule[%v]", e.ErrMsg, e.StartRunTime, e.EndRunTime, e.Schedule)
}

// RenderJSON is used for rendering JSON response body with appropriate headers
func RenderJSON(w http.ResponseWriter, response interface{}) {
	if fmt.Sprint(response) == "[]" {
		response = make([]interface{}, 0)
	}
	data, err := json.Marshal(response)
	if err != nil {
		logger.Log.WarnfCtx(context.TODO(), "RenderJSON: error while marshaling response: %v", err)
		render(w, http.StatusInternalServerError, data)
	} else {
		render(w, http.StatusOK, data)
	}
}

// ParseUUID makes parseUUID
func ParseUUID(s string) (gocql.UUID, error) {
	return gocql.ParseUUID(s)
}

// RenderJSONFromBytes is used for rendering JSON response body from raw bytes with appropriate headers
func RenderJSONFromBytes(w http.ResponseWriter, response []byte) {
	render(w, http.StatusOK, response)
}

// RenderJSONCreated is used for rendering JSON response body when new resource has been created
func RenderJSONCreated(w http.ResponseWriter, response interface{}) {
	if data, err := json.Marshal(response); err != nil {
		logger.Log.ErrfCtx(context.TODO(), errorcode.ErrorCantProcessData, "RenderJSONCreated: error while marshaling response: %v", err)
		render(w, http.StatusInternalServerError, data)
	} else {
		render(w, http.StatusCreated, data)
	}
}

func render(w http.ResponseWriter, code int, response []byte) {
	w.Header().Set(ContentType, ApplicationJSON+"; charset=utf-8")
	w.WriteHeader(code)
	if _, err := w.Write(response); err != nil {
		logger.Log.ErrfCtx(context.TODO(), errorcode.ErrorCantProcessData, "render: error while writing response: %v", err)
	}
}

// SendNoContent sends to the client an empty response with the 204 (NOCONTENT) status
func SendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// SendStatusCodeWithMessage writes a defined string as an error message
// with appropriate headers to the HTTP response
func SendStatusCodeWithMessage(w http.ResponseWriter, r *http.Request, code int, errorCode string, args ...interface{}) {
	translator, err := translation.New(GetLanguage(r))
	if err != nil {
		data, err := json.Marshal(Message{fmt.Sprintf(translation.ErrTranslator), translation.ErrorCodeTranslator})
		if err != nil {
			logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantMarshall, "SendStatusCodeWithMessage: error while marshaling message: %v", err)
		}
		render(w, 500, data)
		return
	}

	msg := Message{Message: fmt.Sprintf(translator.Translate(errorCode), args...)}
	if code != http.StatusOK && code != http.StatusCreated {
		msg.ErrorCode = errorCode
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Log.ErrfCtx(r.Context(), errorcode.ErrorCantMarshall, "SendStatusCodeWithMessage: error while marshaling message: %v", err)
	}
	render(w, code, data)
}

// SendBadRequest sends Bad Request Status and logs an error if it exists
func SendBadRequest(w http.ResponseWriter, r *http.Request, message string, args ...interface{}) {
	SendStatusCodeWithMessage(w, r, http.StatusBadRequest, message, args...)
}

// SendForbidden sends Forbidden Status and logs an error if exists
func SendForbidden(w http.ResponseWriter, r *http.Request, message string, args ...interface{}) {
	SendStatusCodeWithMessage(w, r, http.StatusForbidden, message, args...)
}

// SendNotFound sends Not Fount Status and logs an error if it exists
func SendNotFound(w http.ResponseWriter, r *http.Request, message string, args ...interface{}) {
	SendStatusCodeWithMessage(w, r, http.StatusNotFound, message, args...)
}

// SendInternalServerError sends Internal Server Error Status and logs an error if it exists
func SendInternalServerError(w http.ResponseWriter, r *http.Request, message string, args ...interface{}) {
	SendStatusCodeWithMessage(w, r, http.StatusInternalServerError, message, args...)
}

// SendStatusOkWithMessage sends Status Ok and message
func SendStatusOkWithMessage(w http.ResponseWriter, r *http.Request, message string, args ...interface{}) {
	SendStatusCodeWithMessage(w, r, http.StatusOK, message, args...)
}

// SendCreated sends Created Status with message
func SendCreated(w http.ResponseWriter, r *http.Request, message string, args ...interface{}) {
	SendStatusCodeWithMessage(w, r, http.StatusCreated, message, args...)
}

// GetGeneralInfo used for getting general data about the Tasking service
func GetGeneralInfo() commonlibRest.GeneralInfo {
	return commonlibRest.GeneralInfo{
		TimeStampUTC:    time.Now().UTC(),
		ServiceName:     config.Config.Version.ServiceName,
		ServiceProvider: config.Config.Version.ServiceProvider,
		ServiceVersion:  serviceVersion,
		Name:            config.Config.Version.SolutionName,
	}
}

// FetchStartIndexAndNumberOfRows extracts from the request vars Start Index & # of Rows
func FetchStartIndexAndNumberOfRows(vals url.Values) (startIndex int, numberOfRows int) {

	startIndexStr, ok := vals["startIndex"] // Note type, not ID. ID wasn't specified anywhere.

	startIndex = StartIndexDefault
	var err error

	if ok {
		startIndex, err = strconv.Atoi(startIndexStr[0])
		if err != nil || startIndex <= 0 {
			startIndex = StartIndexDefault
		}
	}

	numberOfRows = NumberOfRowsDefault
	numberOfRowsStr, ok := vals["numberOfRows"]

	if ok {
		numberOfRows, err = strconv.Atoi(numberOfRowsStr[0])
		if err != nil || numberOfRows <= 0 {
			numberOfRows = NumberOfRowsDefault
		}
		if numberOfRows > NumberOfRowsMaximum {
			numberOfRows = NumberOfRowsMaximum
		}
	}

	return
}

// GetLanguage gets the language from Header
func GetLanguage(r *http.Request) string {
	return r.Header.Get(acceptLanguageHeader)
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

// GenerateTargetIDs return random some targets id
func GenerateTargetIDs() []string {
	var targets []string
	for amount := 0; amount < randInt(1, 10); amount++ {
		uuid, _ := gocql.RandomUUID()
		targets = append(targets, uuid.String())
	}
	return targets
}

// GetQuestionMarkString is used to build cassandra SELECT query with IN
func GetQuestionMarkString(count int) string {
	if count <= 0 {
		return ""
	}

	return strings.Repeat("?, ", count-1) + "?"
}

// ExtractUUID is used to extract uuidKey from the request and SendBadRequest if the type of the value is not gocql.UUID
func ExtractUUID(methodName string, responseWriter http.ResponseWriter, request *http.Request, uuidKey string) (uuid gocql.UUID, err error) {
	uuid, err = gocql.ParseUUID(mux.Vars(request)[uuidKey])
	if err != nil {
		logger.Log.ErrfCtx(request.Context(), errorcode.ErrorCantDecodeInputData, "%s: %s has bad format. Err=%s", methodName, uuidKey, err)
		errorCodes := map[string]string{
			"managedEndpointID": errorcode.ErrorEndpointIDHasBadFormat,
			"taskID":            errorcode.ErrorTaskIDHasBadFormat,
			"taskInstanceID":    errorcode.ErrorTaskInstanceIDHasBadFormat,
			"definitionID":      errorcode.ErrorTaskDefinitionIDHasBadFormat,
		}

		errorCode, ok := errorCodes[uuidKey]
		if !ok {
			errorCode = errorcode.ErrorCantDecodeInputData
		}
		SendBadRequest(responseWriter, request, errorCode)
	}
	return
}

func ExtractUUIDs(request *http.Request, uuidKey string) (uuids []gocql.UUID, err error) {
	for _, uuidStr := range strings.Split(mux.Vars(request)[uuidKey], ",") {
		uuid, err := gocql.ParseUUID(uuidStr)
		if err != nil {
			return nil, fmt.Errorf("ExtractUUIDs: uuid has bad format: %s", uuidStr)
		}
		uuids = append(uuids, uuid)
	}

	if len(uuids) == 0 {
		return nil, fmt.Errorf("ExtractUUIDs: there was no uuids found under key : %v", mux.Vars(request)[uuidKey])
	}

	return uuids, nil

}

// ExtractOptionalCount is used to extract optional countKey from the request and SendBadRequest in case of bad input data
func ExtractOptionalCount(methodName string, responseWriter http.ResponseWriter, request *http.Request) (count int, isSpecified bool, err error) {
	countStr := request.URL.Query().Get(countKey)
	if len(countStr) == 0 {
		return 0, false, nil
	}
	count, err = strconv.Atoi(countStr)
	if err != nil || count < 0 {
		logger.Log.WarnfCtx(request.Context(), "%s: variable '%s' should be non-negative integer. Err=%v", methodName, countKey, err)
		SendBadRequest(responseWriter, request, errorcode.ErrorCountVarHasBadFormat)
		return 0, true, fmt.Errorf("ExtractOptionalCount: invalid %s=%d. Err: %s", countKey, count, err)
	}
	return count, true, nil
}

// HTTPRequestWithRetry wrapper on HTTP call with retry logic.
func HTTPRequestWithRetry(ctx context.Context, httpClient *http.Client, httpMethod string, executionURL string, body []byte) (response *http.Response, err error) {
	for retriesCnt := 1; retriesCnt <= config.Config.RetryStrategy.MaxNumberOfRetries; retriesCnt++ {
		response, err = performRequest(ctx, httpClient, httpMethod, executionURL, body)
		if err == nil &&
			(response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated) {
			logger.Log.DebugfCtx(ctx, "Successfully completed URL request %s. Response code: %v",
				executionURL, response.StatusCode)
			return
		}

		var statusCode int
		if response != nil {
			statusCode = response.StatusCode
		}

		if retriesCnt == config.Config.RetryStrategy.MaxNumberOfRetries {
			errMsg := fmt.Sprintf("An error occurred during sending execution payload message to URL: %s. status code %v "+
				"Origin service is out of service. Limit of retry attempts was exhausted. Err: %v", executionURL, statusCode, err)
			err = fmt.Errorf(errMsg)
			logger.Log.DebugfCtx(ctx, errMsg)
			return
		}
		exponentialSleep(retriesCnt)
	}
	return
}

// CloseRespBody will read and close HTTP response body.
func CloseRespBody(response *http.Response) (err error) {
	if response != nil && response.Body != nil {
		if _, err = io.Copy(ioutil.Discard, response.Body); err != nil {
			return err
		}
		return response.Body.Close()
	}
	return nil
}

func exponentialSleep(retriesCnt int) {
	time.Sleep(time.Duration(retriesCnt*config.Config.RetryStrategy.RetrySleepIntervalSec) * time.Second)
}

func performRequest(ctx context.Context, httpClient *http.Client, httpMethod string, executionURL string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(httpMethod, executionURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set(transactionID.Key, transactionID.FromContext(ctx))
	req.Header.Set(ContentType, ApplicationJSON)

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetTimeInLocation is used to parse time in specified location
func GetTimeInLocation(datetime time.Time, location *time.Location) time.Time {
	return time.Date(
		datetime.Year(),
		datetime.Month(),
		datetime.Day(),
		datetime.Hour(),
		datetime.Minute(),
		datetime.Second(),
		datetime.Nanosecond(),
		location,
	)
}

func hourlyNextRunTime(lastRun time.Time, interval int) time.Time {
	return lastRun.Add(time.Hour * time.Duration(interval))
}

func dailyNextRunTime(lastRunTime time.Time, cronSchedule cron.Schedule, schedule apiModels.Schedule) (time.Time, apiModels.Schedule, error) {
	d := cronSchedule.Next(lastRunTime).YearDay()
	if d == schedule.Repeat.Period {
		return cronSchedule.Next(lastRunTime).Truncate(time.Minute), schedule, nil
	}

	newPeriod := schedule.Repeat.Period + schedule.Repeat.Every
	daysInYear := time.Date(lastRunTime.Year(), time.December, 31, 0, 0, 0, 0, lastRunTime.Location()).YearDay()
	for newPeriod > daysInYear {
		newPeriod -= daysInYear
	}
	// Add interval
	lastRunTime = lastRunTime.AddDate(0, 0, schedule.Repeat.Every)
	// Set hour and min to zero
	lastRunTime = time.Date(lastRunTime.Year(), lastRunTime.Month(), lastRunTime.Day(), 0, 0, 0, 0, lastRunTime.Location()).Add(-time.Minute)
	schedule.Repeat.Period = newPeriod
	return dailyNextRunTime(lastRunTime, cronSchedule, schedule)
}

func weeklyNextRunTime(lastRunTime time.Time, cronSchedule cron.Schedule, schedule apiModels.Schedule) (time.Time, apiModels.Schedule, error) {
	_, w := cronSchedule.Next(lastRunTime).ISOWeek()
	if w == schedule.Repeat.Period {
		return cronSchedule.Next(lastRunTime).Truncate(time.Minute), schedule, nil
	}

	newPeriod := schedule.Repeat.Period + schedule.Repeat.Every
	yearsCount := 0
	// last day of the year which can't belong to week 1 for sure based on ISOWeek documentation
	lastDayYear := time.Date(lastRunTime.Year(), 12, 28, 0, 0, 0, 0, lastRunTime.Location())
	_, weeksInYear := lastDayYear.ISOWeek()
	for newPeriod >= weeksInYear {
		newPeriod -= weeksInYear
		lastDayYear = lastDayYear.AddDate(1, 0, 0)
		_, weeksInYear = lastDayYear.ISOWeek()

		yearsCount++
	}
	if newPeriod == 0 {
		newPeriod++
	}
	schedule.Repeat.Period = newPeriod
	newWeek, err := goweek.NewWeek(lastRunTime.Year()+yearsCount, newPeriod)
	if err != nil {
		return weeklyNextRunTime(lastRunTime, cronSchedule, schedule)
	}

	lastRunTime = GetTimeInLocation(newWeek.Days[0], lastRunTime.Location()).Add(-time.Minute)
	return weeklyNextRunTime(lastRunTime, cronSchedule, schedule)
}

func monthlyNextRunTime(lastRunTime time.Time, cronSchedule cron.Schedule, schedule apiModels.Schedule) (time.Time, apiModels.Schedule, error) {
	currentYear, currentMonth, currentDay := lastRunTime.Date()
	_, nextMonth, nextDay := cronSchedule.Next(lastRunTime).Date()

	var maxDaysOfMonth uint = 31
	var minDaysOfMonth uint = 28
	specSchedule, ok := cronSchedule.(*cron.SpecSchedule)
	if !ok {
		return lastRunTime, schedule, fmt.Errorf("wrong scheduler type")
	}

	_, _, lastDayOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, lastRunTime.Location()).AddDate(0, 1, -1).Date()

	initialSpecSchedule := *specSchedule
	for (currentDay < nextDay && currentMonth != nextMonth && currentDay != lastDayOfMonth) ||
		(currentMonth != time.December && math.Abs(float64(nextMonth-currentMonth)) > 1 && maxDaysOfMonth >= minDaysOfMonth) {
		var dayOfMonthBits uint64
		maxDaysOfMonth--
		dayOfMonthBits = ^(math.MaxUint64 << (maxDaysOfMonth + 1)) & (math.MaxUint64 << maxDaysOfMonth)
		specSchedule.Dom = dayOfMonthBits
		cronSchedule = specSchedule
		_, nextMonth, nextDay = cronSchedule.Next(lastRunTime).Date()

	}

	if int(nextMonth) == schedule.Repeat.Period {
		return cronSchedule.Next(lastRunTime).Truncate(time.Minute), schedule, nil
	}

	cronSchedule = &initialSpecSchedule
	newPeriod := schedule.Repeat.Period + schedule.Repeat.Every

	for newPeriod > monthInYear {
		newPeriod -= monthInYear
	}
	lastRunTime = lastRunTime.AddDate(0, schedule.Repeat.Every, 0)
	lastRunTime = time.Date(lastRunTime.Year(), time.Month(newPeriod), 1, 0, 0, 0, 0, lastRunTime.Location()).Add(-time.Minute)
	schedule.Repeat.Period = newPeriod
	return monthlyNextRunTime(lastRunTime, cronSchedule, schedule)
}

// GetNextRunTime returns next run time based on given time, cron string and interval
func GetNextRunTime(lastRunTime time.Time, schedule apiModels.Schedule) (time.Time, apiModels.Schedule, error) {
	cronStr := schedule.Cron()

	// parse cron string
	cronSchedule, err := CronParser.Parse(cronStr)
	if err != nil {
		return lastRunTime, schedule, err
	}

	switch schedule.Repeat.Frequency {
	case apiModels.Hourly:
		return hourlyNextRunTime(lastRunTime, schedule.Repeat.Every), schedule, err
	case apiModels.Daily:
		return dailyNextRunTime(lastRunTime, cronSchedule, schedule)
	case apiModels.Weekly:
		return weeklyNextRunTime(lastRunTime, cronSchedule, schedule)
	case apiModels.Monthly:
		return monthlyNextRunTime(lastRunTime, cronSchedule, schedule)
	default:
		return lastRunTime, schedule, fmt.Errorf("wrong scheduling interval %d", schedule.Repeat.Frequency)
	}
}

// CalcFirstNextRunTime - returns first run time for task in UTC
func CalcFirstNextRunTime(currentTime time.Time, s apiModels.Schedule) (time.Time, error) {
	runTime, err := calcFirstNextRunTime(s)
	if err != nil {
		return time.Time{}, err
	}

	runTime = runTime.UTC()
	for !runTime.Truncate(time.Minute).UTC().After(currentTime.Truncate(time.Minute).UTC()) {
		runTime, err = CalcNextRunTime(runTime, s, *s.StartRunTime.Location())
		if err != nil {
			return time.Time{}, err
		}
	}
	return runTime, nil
}

func calcFirstNextRunTime(s apiModels.Schedule) (time.Time, error) {
	switch s.Repeat.Frequency {
	case apiModels.Hourly:
		nowInLoc := time.Now().In(s.StartRunTime.Location()).Truncate(time.Minute)
		execHour := s.StartRunTime
		if nowInLoc.Before(execHour) {
			return checkEndRunTime(s.EndRunTime, execHour)
		}

		execHour = execHour.Add(1 * time.Hour)
		return checkEndRunTime(s.EndRunTime, execHour)

	case apiModels.Daily:
		startDay := s.StartRunTime
		startTime := s.Repeat.RunTime
		execDay := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), startTime.Hour(), startTime.Minute(), 0, 0, startDay.Location())
		if !execDay.Before(startDay) {
			return checkEndRunTime(s.EndRunTime, execDay)
		}

		execDay = execDay.AddDate(0, 0, 1)
		return checkEndRunTime(s.EndRunTime, execDay)
	case apiModels.Weekly:
		daysOfWeek := s.Repeat.DaysOfWeek
		if !sort.IntsAreSorted(daysOfWeek) {
			sort.Ints(daysOfWeek)
		}

		startDay := getNextWeekDayFrom(s.StartRunTime, daysOfWeek)
		startTime := s.Repeat.RunTime
		execDay := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), startTime.Hour(), startTime.Minute(), 0, 0, startDay.Location())
		if !execDay.Before(s.StartRunTime) {
			return checkEndRunTime(s.EndRunTime, execDay)
		}

		startDay = getNextWeekDayFrom(s.StartRunTime.AddDate(0, 0, 1), daysOfWeek)
		execDay = time.Date(startDay.Year(), startDay.Month(), startDay.Day(), startTime.Hour(), startTime.Minute(), 0, 0, startDay.Location())
		return checkEndRunTime(s.EndRunTime, execDay)

	case apiModels.Monthly:
		execDay, err := findExecMonthDay(s.StartRunTime, s)
		if err != nil {
			return time.Time{}, err
		}
		if !execDay.IsZero() {
			return execDay, nil
		}

		nextMonth := getNextMonth(s.StartRunTime)
		execDay, err = findExecMonthDay(nextMonth, s)
		if err != nil {
			return time.Time{}, err
		}
		return execDay, nil

	default:
		return time.Time{}, errors.New("such frequency is not allowed")
	}
}

func getNextWeekDayFrom(current time.Time, weekDays []int) time.Time {
	currentWeekDay := int(current.Weekday())
	for _, day := range weekDays {
		if day >= currentWeekDay {
			dayDiff := day - currentWeekDay
			return current.AddDate(0, 0, dayDiff)
		}
	}
	dayDiff := weekDays[0] - currentWeekDay
	return current.AddDate(0, 0, dayDiff+daysInWeek)
}

func findExecMonthDay(month time.Time, s apiModels.Schedule) (time.Time, error) {
	scheduledDays := getScheduledDays(month, s.Repeat)
	for _, scheduledDay := range scheduledDays {
		date := time.Date(month.Year(), month.Month(), scheduledDay, s.Repeat.RunTime.Hour(), s.Repeat.RunTime.Minute(), 0, 0, month.Location())

		// need to check if produced date is still current month
		// if not should create date for last day of current month
		if date.Month() != month.Month() {
			date = lastDayInMonth(month)
			date = time.Date(date.Year(), date.Month(), date.Day(), s.Repeat.RunTime.Hour(), s.Repeat.RunTime.Minute(), 0, 0, month.Location())
		}
		if date.Before(s.StartRunTime) {
			continue
		}
		return checkEndRunTime(s.EndRunTime, date)
	}
	return time.Time{}, nil
}

func getNextMonth(current time.Time) time.Time {
	firstDayOfCurrentMonth := time.Date(current.Year(), current.Month(), firstDay, 0, 0, 0, 0, current.Location())
	return firstDayOfCurrentMonth.AddDate(0, 1, 0)
}

func checkEndRunTime(endRunTime, nextRunTime time.Time) (time.Time, error) {
	execHour, err := checkScheduledRange(endRunTime, nextRunTime)
	if err != nil {
		return time.Time{}, errors.New(taskWillNeverRunError)
	}
	return execHour, nil
}

// CalcNextRunTime returns next run time based on given time and schedule
func CalcNextRunTime(lastRunTime time.Time, s apiModels.Schedule, actualLocation time.Location) (time.Time, error) {
	if s.Regularity == apiModels.RunNow {
		return time.Time{}, errors.New("schedule can't be used for run now tasks")
	}
	if s.Repeat.Every < 1 {
		return time.Time{}, errors.New("field every for recurrent task should be greater than 0")
	}
	// need to set actual location for EndRunTime to get right nextRunTime in local machine time
	if s.Location == "" && !s.EndRunTime.IsZero() {
		s.EndRunTime = s.EndRunTime.In(&actualLocation)
	}

	switch s.Repeat.Frequency {
	case apiModels.Hourly:
		return handleHourlyFrequency(lastRunTime, s)
	case apiModels.Daily:
		return handleDailyFrequency(lastRunTime, s)
	case apiModels.Weekly:
		return handleWeeklyFrequency(lastRunTime, s)
	case apiModels.Monthly:
		return handleMonthlyFrequency(lastRunTime, s)
	default:
		return time.Time{}, errors.Errorf("wrong scheduling interval %d", s.Repeat.Frequency)
	}
}

func handleHourlyFrequency(lastRunTime time.Time, s apiModels.Schedule) (time.Time, error) {
	timeOffset := time.Duration(s.Repeat.Every) * time.Hour
	nextRunTime := lastRunTime.Add(timeOffset)
	return checkScheduledRange(s.EndRunTime, nextRunTime)
}

func handleDailyFrequency(lastRunTime time.Time, s apiModels.Schedule) (time.Time, error) {
	return getNextRunDay(lastRunTime, s.EndRunTime, s.Repeat.Every)
}

func handleWeeklyFrequency(lastRunTime time.Time, s apiModels.Schedule) (time.Time, error) {
	if len(s.Repeat.DaysOfWeek) == 0 {
		return time.Time{}, errors.New("Days of week can't be empty")
	}

	daysOfWeek := s.Repeat.DaysOfWeek
	daysOfWeek = checkAndRemoveDuplicates(daysOfWeek)

	if !sort.IntsAreSorted(daysOfWeek) {
		sort.Ints(daysOfWeek)
	}

	prevWeekDay := int(lastRunTime.Weekday())
	nextWeekDay := getNextScheduledDay(prevWeekDay, daysOfWeek)

	if nextWeekDay > prevWeekDay {
		daysOffset := nextWeekDay - prevWeekDay
		return getNextRunDay(lastRunTime, s.EndRunTime, daysOffset)
	}

	dayDiff := daysOfWeek[0] - daysOfWeek[len(daysOfWeek)-1]
	firstScheduledDayOfWeek := lastRunTime.AddDate(0, 0, dayDiff)

	return getNextRunDay(firstScheduledDayOfWeek, s.EndRunTime, s.Repeat.Every*daysInWeek)
}

func handleMonthlyFrequency(lastRunTime time.Time, s apiModels.Schedule) (time.Time, error) {
	if len(s.Repeat.DaysOfMonth) == 0 && len(s.Repeat.WeekDays) == 0 {
		return time.Time{}, errors.New("there should be month days for execution")
	}

	daysOfCurrentMonth := getScheduledDays(lastRunTime, s.Repeat)

	prevMonthDay := lastRunTime.Day()
	nextMonthDay := getNextScheduledDay(prevMonthDay, daysOfCurrentMonth)

	if nextMonthDay > prevMonthDay {
		// in order to avoid possibly redundant steps to find last day of month
		// check if nextMonthDay is less or equal 28 - minimum days in each months
		if nextMonthDay <= minDaysInMonth {
			return getNextRunDay(lastRunTime, s.EndRunTime, nextMonthDay-prevMonthDay)
		}

		// if first check fails than need to find last day of months
		lastDayOfMonth := lastDayInMonth(lastRunTime).Day()
		if nextMonthDay <= lastDayOfMonth {
			return getNextRunDay(lastRunTime, s.EndRunTime, nextMonthDay-prevMonthDay)
		}

		// need to get next run time for last day of month if lastRunTime is not such day
		daysOffset := lastDayOfMonth - prevMonthDay
		if daysOffset > 0 {
			return getNextRunDay(lastRunTime, s.EndRunTime, daysOffset)
		}
	}

	// in order to safely get next run month make sense to get first day of current month and add to this date month's offset
	firstDayOfMonth := getFirstDayOfMonth(lastRunTime)
	nextMonth := firstDayOfMonth.AddDate(0, s.Repeat.Every, 0)

	lastDayOfNextMonth := lastDayInMonth(nextMonth).Day()
	daysOfNextMonth := getScheduledDays(nextMonth, s.Repeat)
	firstScheduledDay := daysOfNextMonth[0]
	if firstScheduledDay <= lastDayOfNextMonth {
		daysOffset := firstScheduledDay - firstDay
		return getNextRunDay(nextMonth, s.EndRunTime, daysOffset)
	}
	daysOffset := lastDayOfNextMonth - firstDay
	return getNextRunDay(nextMonth, s.EndRunTime, daysOffset)
}

// getScheduledDays - returns updated slice of scheduled days including special day which date changes from month to month
func getScheduledDays(lastRunTime time.Time, repeat apiModels.Repeat) []int {
	var scheduledDays []int
	if len(repeat.DaysOfMonth) != 0 {
		scheduledDays = make([]int, len(repeat.DaysOfMonth))
		copy(scheduledDays, repeat.DaysOfMonth)
	} else {
		scheduledDays = make([]int, 0, config.MaxWeekDaysInSchedule)
		for _, wd := range repeat.WeekDays {
			wdDate := getWeekDayDate(lastRunTime, wd)
			scheduledDays = append(scheduledDays, wdDate)
		}
	}

	scheduledDays = checkAndRemoveDuplicates(scheduledDays)

	if !sort.IntsAreSorted(scheduledDays) {
		sort.Ints(scheduledDays)
	}
	return scheduledDays
}

func getNextScheduledDay(previous int, days []int) int {
	for _, day := range days {
		if day > previous {
			return day
		}
	}
	return days[0]
}

func checkAndRemoveDuplicates(days []int) []int {
	daysBuff := make(map[int]struct{})
	for _, day := range days {
		daysBuff[day] = struct{}{}
	}

	checked := make([]int, 0)
	for day := range daysBuff {
		checked = append(checked, day)
	}
	return checked
}

func lastDayInMonth(date time.Time) time.Time {
	firstDayOfMonth := time.Date(date.Year(), date.Month(), firstDay, 0, 0, 0, 0, date.Location())
	return firstDayOfMonth.AddDate(0, 1, -1)
}

func isSameMonth(date time.Time, month time.Month) bool {
	return date.Month() == month
}

func getNextRunDay(lastRunTime time.Time, endRunTime time.Time, daysOffset int) (time.Time, error) {
	nextRunTime := lastRunTime.AddDate(0, 0, daysOffset)
	return checkScheduledRange(endRunTime, nextRunTime)
}

func checkScheduledRange(endRunTime, nextRunTime time.Time) (time.Time, error) {
	if endRunTime.IsZero() || nextRunTime.Before(endRunTime) {
		return nextRunTime, nil
	}
	return time.Time{}, errors.New(NextRunTimeExceedsEndRunTime)
}

// getWeekDayDate - returns date of special month's day like "last friday of month" or "second wednesday of month" etc.
func getWeekDayDate(runTime time.Time, weekDay apiModels.WeekDay) int {
	// need to find date of WeekDay with some index from the beginning of month
	firstDayOfMonth := getFirstDayOfMonth(runTime)
	dayOfMonth := firstDayOfMonth
	var weekDayDate int
	for i := 0; i < daysInWeek; i++ {
		// check if current weekday equal needed weekday
		// if yes - first weekday of this month found we can break search loop
		if dayOfMonth.Weekday() == weekDay.Day {
			weekDayDate = dayOfMonth.Day()
			break
		}
		// if upper comparison failed go to the next day and do it again
		dayOfMonth = dayOfMonth.AddDate(0, 0, 1)
	}

	// when needed weekday found need to find date of such weekday that matches index by adding on each iteration step
	// daysInWeek = 7
	for i := 0; i < int(weekDay.Index); i++ {
		weekDayDate += daysInWeek
	}

	// it can be 4 or 5 particular weekdays in some month
	// for one month index Last could be 4 for other 5.
	// if adding a date on last iteration moves date to the next month then Last = 4 and
	// need to move for 7 days later
	if weekDay.Index == apiModels.Last {
		actualDay := firstDayOfMonth.AddDate(0, 0, weekDayDate-firstDay)
		if !isSameMonth(actualDay, firstDayOfMonth.Month()) {
			weekDayDate -= daysInWeek
		}
	}

	return weekDayDate
}

func getFirstDayOfMonth(current time.Time) time.Time {
	daysFromBeginning := firstDay - current.Day()
	return current.AddDate(0, 0, daysFromBeginning)
}

// GetNextRunTimeByCron returns next run time based on given time and cron string
func GetNextRunTimeByCron(lastRunTime time.Time, scheduleStr string) (time.Time, error) {
	schedule, err := CronParser.Parse(scheduleStr)
	if err == nil {
		return schedule.Next(lastRunTime).Truncate(time.Minute), err
	}
	return lastRunTime, err
}

// AddLocationToTime adds location to time
func AddLocationToTime(t time.Time, location *time.Location) time.Time {
	if location != nil {
		return time.Date(
			t.Year(),
			t.Month(),
			t.Day(),
			t.Hour(),
			t.Minute(),
			t.Second(),
			t.Nanosecond(),
			location,
		)
	}

	return t
}

// ConvertUUIDsToInterfaces converts a slice of gocql.UUIDs to slice of interfaces
func ConvertUUIDsToInterfaces(in []gocql.UUID) []interface{} {
	var out = make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

// ConvertStringsToUUIDs converts a slice of string to a slice of gocql.UUIDs
func ConvertStringsToUUIDs(in []string) ([]gocql.UUID, error) {
	var out = make([]gocql.UUID, len(in))
	for i, v := range in {
		uuid, err := gocql.ParseUUID(v)
		if err != nil {
			return nil, err
		}
		out[i] = uuid
	}
	return out, nil
}

// UUIDSliceContainsElement is used to check uuid slice for duplicates
func UUIDSliceContainsElement(uniqueUUIDSlice []gocql.UUID, element gocql.UUID) bool {
	for _, v := range uniqueUUIDSlice {
		if v == element {
			return true
		}
	}
	return false
}

// UserFromCtx returns user from context
func UserFromCtx(ctx context.Context) (us entities.User, err error) {
	var (
		user interface{}
		ok   bool
	)

	if user = ctx.Value(config.UserKeyCTX); user == nil {
		return entities.User{}, fmt.Errorf("cannot get user from context")
	}

	us, ok = user.(entities.User)
	if !ok {
		return entities.User{}, fmt.Errorf("cannot assert interface to user model")
	}
	return
}

// UsersEndpointsFromCtx  returns user endpoints from context
func UsersEndpointsFromCtx(ctx context.Context) (endpointsSlice []string, err error) {
	var (
		endpoints interface{}
		ok        bool
	)

	if endpoints = ctx.Value(config.UserEndPointsKeyCTX); endpoints == nil {
		return nil, fmt.Errorf("cannot get endpoints from context")
	}
	endpointsSlice, ok = endpoints.([]string)
	if !ok {
		return nil, fmt.Errorf("cannot assert interface to slice")
	}
	return
}

// SliceToMap converts slice to map
func SliceToMap(slice []string) (mapString map[string]struct{}) {
	mapString = make(map[string]struct{})
	for _, value := range slice {
		mapString[value] = struct{}{}
	}
	return
}
