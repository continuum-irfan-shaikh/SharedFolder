package validator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/ContinuumLLC/godep-govalidator"
	"github.com/gocql/gocql"
	"github.com/xeipuuv/gojsonschema"
	agentModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/agent"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/triggers"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common"
)

var (
	parametersPattern = regexp.MustCompile(`"required"\s*:\s*\[[\s",a-zA-Z]+]`)
)

const (
	emptyRepeatString = "{0 0001-01-01 00:00:00 +0000 UTC 0 [] <nil> [] [] 0}"
	firstDayOfMonth   = 1
	lastDayOfMonth    = 31
)

// logical XNOR
func xnor(x1, x2 bool) bool {
	return (!x1 && !x2) || (x1 && x2)
}

// validatorFunction returns the function for custom validators based on the `empty` parameter
// empty == true gives us a function for unsettableByUsers fields validator
// empty == false gives us a function for requiredForUsers fields validator
func validatorFunction(empty bool) func(i interface{}, context interface{}) bool {
	return func(i interface{}, context interface{}) bool {
		switch context.(type) {
		case models.Task, models.TaskDefinitionDetails, models.TaskDefinition, models.SelectedManagedEndpointEnable, models.Target, apiModels.Schedule, apiModels.Repeat:
		default:
			return false
		}
		switch v := i.(type) {
		case gocql.UUID:
			var emptyUUID gocql.UUID
			return xnor(empty, v == emptyUUID)
		case time.Time:
			return xnor(empty, v.IsZero()) && common.ValidTime(v, time.Now().UTC())
		case string, []string, map[string]bool, []models.ManagedEndpointDetailed:
			return xnor(empty, reflect.ValueOf(v).Len() == 0)
		case int, statuses.TaskState, apiModels.Regularity, models.TargetType:
			return xnor(empty, reflect.ValueOf(v).Int() == 0) // Only for int based types
		}
		return false
	}
}

func validatorType(i interface{}, ctx interface{}) bool {
	switch ctx.(type) {
	case models.Task, models.TaskDefinition, models.TaskDefinitionDetails:
	default:
		return false
	}
	field, ok := i.(string)
	if !ok {
		return false
	}
	_, ok = config.Config.TaskTypes[field]
	return ok
}

func validPositive(i interface{}, _ interface{}) bool {
	number, ok := i.(int)
	if !ok {
		return false
	}
	return number >= 0
}

func validatorCreds(i interface{}, ctx interface{}) bool {
	switch ctx.(type) {
	case models.Task, models.TaskDefinitionDetails:
	default:
		return false
	}

	field, ok := i.(*agentModels.Credentials)
	if !ok {
		return false
	}

	if field == nil {
		return true
	}

	if field.UseCurrentUser && field.Username == "" && field.Password == "" && field.Domain == "" {
		return true
	}

	if !field.UseCurrentUser && field.Username != "" {
		return true
	}

	return false
}

func validatorOneTime(i interface{}, ctx interface{}) bool {
	if structInstance, ok := ctx.(models.Task); ok {
		switch field := i.(type) {
		case string: // TaskRequest.Trigger
			if structInstance.Schedule.Regularity == apiModels.OneTime && structInstance.Schedule.StartRunTime.IsZero() {
				return len(field) != 0
			}
			return len(field) == 0
		}
	}
	return false
}

func validResourceType(i interface{}, ctx interface{}) bool {
	var ok bool
	if _, ok = ctx.(models.Task); !ok {
		return false
	}

	r, ok := i.(integration.ResourceType)
	if !ok {
		return false
	}

	if r.IsAllResources() {
		return true
	}

	switch r {
	case integration.Desktop, integration.Server:
		return true
	default:
		return false
	}
}

func validTargets(i interface{}, ctx interface{}) bool {
	var ok bool
	var task models.Task
	if task, ok = ctx.(models.Task); !ok {
		return false
	}

	t, ok := i.(models.Target)
	if !ok {
		return false
	}

	if len(task.TargetsByType) > 0 && len(t.IDs) > 0 {
		return false
	}

	if len(t.IDs) == 0 && len(task.TargetsByType) == 0 {
		return false
	}

	for _, ids := range task.TargetsByType {
		if !consistsOfUniqueNotEmptyElements(ids) {
			return false
		}
	}

	_, hasSites := task.TargetsByType[models.Site]
	_, hasDynamicSites := task.TargetsByType[models.DynamicSite]
	if hasSites && hasDynamicSites {
		return false
	}

	return true
}

func validatorDynamicGroup(i interface{}, ctx interface{}) bool {
	if structInstance, ok := ctx.(models.Task); ok {
		_, ok := i.(apiModels.Schedule)
		if !ok {
			return false
		}

		if structInstance.IsDynamicGroupBasedTrigger() && len(structInstance.TargetsByType) > 0 && len(structInstance.TargetsByType[models.DynamicGroup]) == 0 {
			return false
		}

		if structInstance.IsDynamicGroupBasedTrigger() && len(structInstance.TargetsByType) > 0 && len(structInstance.TargetsByType[models.DynamicGroup]) != 0 {
			return true
		}

		if structInstance.IsDynamicGroupBasedTrigger() && structInstance.Targets.Type != models.DynamicGroup {
			return false
		}
		return true
	}
	return false
}

func validatorLocation(i interface{}, ctx interface{}) bool {
	structInstance, ok := ctx.(apiModels.Schedule)
	if !ok {
		return false
	}

	field, ok := i.(string)
	if !ok {
		return false
	}

	if structInstance.Regularity == apiModels.RunNow {
		return len(field) == 0
	}

	if _, err := time.LoadLocation(field); err != nil {
		return false
	}

	return true
}

func validatorTargets(i interface{}, ctx interface{}) bool {
	structInstance, ok := ctx.(models.Target)
	if !ok {
		return false
	}

	field, ok := i.([]string)
	if !ok {
		return false
	}

	// MANAGED_ENDPOINT ids and DYNAMIC_GROUP ids should be in the uuid format
	if structInstance.Type != models.Site &&  structInstance.Type != models.DynamicSite {
		_, err := common.ConvertStringsToUUIDs(field)
		if err != nil {
			return false
		}
	}

	return consistsOfUniqueNotEmptyElements(field)
}

func validatorRequiredRecurrent(i interface{}, ctx interface{}) bool {
	if structInstance, ok := ctx.(apiModels.Schedule); ok {
		switch field := i.(type) {
		case apiModels.Repeat:
			if structInstance.Regularity != apiModels.Recurrent {
				return fmt.Sprint(field) == emptyRepeatString
			}

			var isValid = field.Every > 0 && field.Frequency > 0 && field.Frequency <= apiModels.Monthly

			if structInstance.Repeat.Frequency == apiModels.Hourly {
				return isValid
			}

			return isValid && !field.RunTime.IsZero()
		default:
			if structInstance.Regularity == apiModels.Trigger {
				return true
			}
		}
	}

	return false
}

func validatorRequiredRecurrendAndOneTime(i interface{}, ctx interface{}) bool {
	if structInstance, ok := ctx.(apiModels.Schedule); ok {
		field, ok := i.(time.Time) // Schedule.StartRunTime
		if !ok {
			return false
		}

		if structInstance.Location != "" && structInstance.Regularity != apiModels.RunNow {
			location, err := time.LoadLocation(structInstance.Location)
			if err != nil {
				return false
			}
			fieldWithLoc := common.AddLocationToTime(field, location)
			return common.ValidTime(fieldWithLoc, time.Now().In(location))
		}

		if structInstance.Regularity == apiModels.RunNow {
			return field.IsZero()
		}
		return !field.IsZero()
	}
	return false
}

func validatorRequiredMonthly(i interface{}, ctx interface{}) bool {
	structInstance, ok := ctx.(apiModels.Repeat)
	if !ok {
		return false
	}

	field, ok := i.([]int)
	if !ok {
		return false
	}

	if structInstance.Frequency != apiModels.Monthly {
		return len(field) == 0
	}
	return len(field) > 0
}

func validOptionalBetween(i interface{}, ctx interface{}) bool {
	schedule, ok := ctx.(apiModels.Schedule)
	if !ok {
		return false
	}

	betweenEndTIme, ok := i.(time.Time)
	if !ok {
		return false
	}

	if betweenEndTIme.IsZero() {
		return true
	}

	if schedule.Regularity == apiModels.RunNow {
		return false
	}

	if schedule.Repeat.Frequency == apiModels.Hourly {
		return false
	}

	if schedule.Regularity == apiModels.OneTime {
		if !schedule.StartRunTime.Before(betweenEndTIme) {
			return false
		}

		return schedule.StartRunTime.Add(time.Hour * 24).After(betweenEndTIme)
	}

	if !schedule.Repeat.RunTime.Before(betweenEndTIme) {
		return false
	}

	return schedule.Repeat.RunTime.Add(time.Hour * 24).After(betweenEndTIme)
}

func validatorOptionalMonthly(i interface{}, ctx interface{}) bool {
	structInstance, ok := ctx.(apiModels.Repeat)
	if !ok {
		return false
	}

	if structInstance.WeekDay != nil && len(structInstance.WeekDays) > 0 {
		return false
	}

	if structInstance.WeekDay != nil {
		structInstance.WeekDays = append(structInstance.WeekDays, *structInstance.WeekDay)
	}

	if structInstance.Frequency != apiModels.Monthly {
		return len(structInstance.DaysOfMonth) == 0 && len(structInstance.WeekDays) == 0
	}

	if (len(structInstance.DaysOfMonth) == 0) == (len(structInstance.WeekDays) == 0) {
		return false
	}

	switch field := i.(type) {
	case []int:
		return containsValidMonthDays(field)
	case []apiModels.WeekDay:
		return containsValidWeekDay(field)
	default:
		return false
	}
}

func containsValidMonthDays(monthDays []int) bool {
	for _, monthDay := range monthDays {
		if monthDay < firstDayOfMonth || monthDay > lastDayOfMonth {
			return false
		}
	}

	return true
}

func containsValidWeekDay(weekDays []apiModels.WeekDay) bool {
	if len(weekDays) > config.MaxWeekDaysInSchedule {
		return false
	}

	for _, wd := range weekDays {
		if wd.Day < time.Sunday || wd.Day > time.Saturday || wd.Index < apiModels.First || wd.Index > apiModels.Last {
			return false
		}
	}

	return true
}

func validatorRequiredWeekly(i interface{}, ctx interface{}) bool {
	structInstance, ok := ctx.(apiModels.Repeat)
	if !ok {
		return false
	}

	field, ok := i.([]int)
	if !ok {
		return false
	}

	if structInstance.Frequency != apiModels.Weekly {
		return len(field) == 0
	}

	if len(field) < 1 {
		return false
	}

	for _, day := range field {
		if day < 0 || day > 6 {
			return false
		}
	}
	return true
}

func validatorOptionalRecurrent(i interface{}, ctx interface{}) bool {
	structInstance, ok := ctx.(apiModels.Schedule)
	if !ok {
		return false
	}

	field, ok := i.(time.Time) // TaskRequest.EndRunTime
	if !ok {
		return false
	}

	switch structInstance.Regularity {
	case apiModels.Recurrent:
		_, err := common.GetNextRunTimeByCron(structInstance.StartRunTime.Add((-1)*time.Minute), structInstance.Cron())
		if err != nil { // invalid cron TaskRequest.Schedule
			return common.ValidTime(field, structInstance.StartRunTime)
		}
		return common.ValidTime(field, structInstance.StartRunTime)

	case apiModels.Trigger, apiModels.RunNow:
		if field.IsZero() {
			return true
		}
		return common.ValidTime(field, structInstance.StartRunTime)
	}
	return field.IsZero()
}

func validatorOptionalTriggerTypes(i interface{}, ctx interface{}) bool {
	var task models.Task
	var ok bool
	if task, ok = ctx.(models.Task); !ok {
		return false
	}

	schedule, ok := i.(apiModels.Schedule)
	if !ok {
		return false
	}

	for _, trigger := range schedule.TriggerTypes {
		if trigger != triggers.FirstCheckInTrigger {
			continue
		}

		if len(task.TargetsByType) > 0 &&
			(len(task.TargetsByType[models.Site]) == 0 && len(task.TargetsByType[models.DynamicSite]) == 0) {
			return false
		}

		if (task.Targets.Type != models.Site && task.Targets.Type != models.DynamicSite) && len(task.TargetsByType) == 0 {
			return false
		}
	}
	return true
}

func validatorCategories(i interface{}, ctx interface{}) bool {
	if _, ok := ctx.(models.TaskDefinition); !ok {
		return false
	}

	field, ok := i.([]string)
	if !ok {
		return false
	}

	// if it's empty - it's ok. it'll be set as custom
	if len(field) == 0 {
		return true
	}

	return consistsOfUniqueNotEmptyElements(field)
}

func consistsOfUniqueNotEmptyElements(stringSlice []string) bool {
	for i, v := range stringSlice {
		if len(strings.TrimSpace(v)) == 0 || common.IsElementAlreadyExists(v, stringSlice[i+1:]) {
			return false
		}
	}
	return true
}

func validatorRequiredForTriggers(i interface{}, ctx interface{}) bool {
	structInstance, ok := ctx.(apiModels.Schedule)
	if !ok {
		return false
	}

	field, ok := i.([]apiModels.TriggerFrame)
	if !ok {
		return false
	}

	switch structInstance.Regularity {
	case apiModels.OneTime, apiModels.RunNow:
		return len(field) == 0
	case apiModels.Recurrent:
		if len(field) == 0 {
			return true
		}
	case apiModels.Trigger:
		if len(field) == 0 {
			return false
		}
	default:
		return false
	}

	return validateTriggerTypes(field)
}

func validateTriggerTypes(field []apiModels.TriggerFrame) bool {
	for _, triggerFrame := range field {
		if triggerFrame.TriggerType == "" {
			return false
		}

		if !triggerFrame.StartTimeFrame.IsZero() && !triggerFrame.EndTimeFrame.IsZero() {
			endTime := setTimeToOneDay(triggerFrame.StartTimeFrame, triggerFrame.EndTimeFrame)
			if triggerFrame.StartTimeFrame.After(endTime) {
				return false
			}
			continue
		}

		if triggerFrame.StartTimeFrame.IsZero() && triggerFrame.EndTimeFrame.IsZero() {
			continue
		}
		return false
	}
	return true
}

func validateRecurrentDGTriggerTarget(i interface{}, ctx interface{}) bool {
	if structInstance, ok := ctx.(models.Task); ok {
		field, ok := i.(apiModels.Schedule)
		if !ok {
			return false
		}

		// create task with recurrence + DG based trigger only against DG targetType
		if containsDGTrigger(field.TriggerFrames) && field.Regularity == apiModels.Recurrent &&
			(len(structInstance.TargetsByType) > 1 || len(structInstance.TargetsByType[models.DynamicGroup]) == 0) {
			return false
		}
		return true
	}
	return false
}

func validatorTriggerTypeOptional(i interface{}, ctx interface{}) bool {
	schedule, ok := ctx.(apiModels.Schedule)
	if !ok {
		return false
	}

	triggerTypes, ok := i.([]string)
	if !ok {
		return false
	}

	if len(triggerTypes) == 0 && schedule.Regularity != apiModels.Trigger {
		return true
	}

	if len(triggerTypes) == 0 && schedule.Regularity == apiModels.Trigger {
		return false
	}

	if len(triggerTypes) != len(schedule.TriggerFrames) {
		return false
	}
	return validTypesFrame(triggerTypes, schedule)
}

func validTypesFrame(triggerTypes []string, schedule apiModels.Schedule) bool {
	triggersMap := make(map[string]struct{})
	for _, triggerType := range triggerTypes {
		// checking for uniqueness
		if _, ok := triggersMap[triggerType]; ok {
			return false
		}
		triggersMap[triggerType] = struct{}{}

		// checking for mandatory duplication in triggerFrames slice
		existsInTriggerFrames := false
		for i := range schedule.TriggerFrames {
			if schedule.TriggerFrames[i].TriggerType == triggerType {
				existsInTriggerFrames = true
				break
			}
		}

		if !existsInTriggerFrames {
			return false
		}
	}
	return true
}

func containsDGTrigger(frames []apiModels.TriggerFrame) bool {
	for _, v := range frames {
		if v.TriggerType == triggers.DynamicGroupEnterTrigger || v.TriggerType == triggers.DynamicGroupExitTrigger {
			return true
		}
	}
	return false
}

func setTimeToOneDay(dayToSet, timeToSet time.Time) time.Time {
	if dayToSet.Year() != timeToSet.Year() || dayToSet.Month() != timeToSet.Month() || dayToSet.Day() != timeToSet.Day() {
		return time.Date(dayToSet.Year(), dayToSet.Month(), dayToSet.Day(), timeToSet.Hour(), timeToSet.Minute(), timeToSet.Second(), timeToSet.Nanosecond(), timeToSet.Location())
	}
	return timeToSet
}

// SetupCustomValidators is used to set up the validators of the request
func SetupCustomValidators() {
	govalidator.CustomTypeTagMap.Set("unsettableByUsers", govalidator.CustomTypeValidator(validatorFunction(true)))
	govalidator.CustomTypeTagMap.Set("requiredForUsers", govalidator.CustomTypeValidator(validatorFunction(false)))
	govalidator.CustomTypeTagMap.Set("requiredOnlyForOneTime", govalidator.CustomTypeValidator(validatorOneTime))
	govalidator.CustomTypeTagMap.Set("validResourceType", govalidator.CustomTypeValidator(validResourceType))
	govalidator.CustomTypeTagMap.Set("requiredOnlyForRecurrent", govalidator.CustomTypeValidator(validatorRequiredRecurrent))
	govalidator.CustomTypeTagMap.Set("optionalOnlyForRecurrent", govalidator.CustomTypeValidator(validatorOptionalRecurrent))
	govalidator.CustomTypeTagMap.Set("validCategories", govalidator.CustomTypeValidator(validatorCategories))
	govalidator.CustomTypeTagMap.Set("validType", govalidator.CustomTypeValidator(validatorType))
	govalidator.CustomTypeTagMap.Set("validLocation", govalidator.CustomTypeValidator(validatorLocation))
	govalidator.CustomTypeTagMap.Set("requiredUniqueTargetIDs", govalidator.CustomTypeValidator(validatorTargets))
	govalidator.CustomTypeTagMap.Set("validatorDynamicGroup", govalidator.CustomTypeValidator(validatorDynamicGroup))
	govalidator.CustomTypeTagMap.Set("optionalTriggerTypes", govalidator.CustomTypeValidator(validatorOptionalTriggerTypes))
	govalidator.CustomTypeTagMap.Set("requiredOnlyForMonthly", govalidator.CustomTypeValidator(validatorRequiredMonthly))
	govalidator.CustomTypeTagMap.Set("optionalOnlyForMonthly", govalidator.CustomTypeValidator(validatorOptionalMonthly))
	govalidator.CustomTypeTagMap.Set("requiredOnlyForWeekly", govalidator.CustomTypeValidator(validatorRequiredWeekly))
	govalidator.CustomTypeTagMap.Set("requiredForRecurrendAndOneTime", govalidator.CustomTypeValidator(validatorRequiredRecurrendAndOneTime))
	govalidator.CustomTypeTagMap.Set("requiredForTriggers", govalidator.CustomTypeValidator(validatorRequiredForTriggers))
	govalidator.CustomTypeTagMap.Set("recurrentDGTriggerTarget", govalidator.CustomTypeValidator(validateRecurrentDGTriggerTarget))
	govalidator.CustomTypeTagMap.Set("optionalValidTriggerTypes", govalidator.CustomTypeValidator(validatorTriggerTypeOptional))
	govalidator.CustomTypeTagMap.Set("validCreds", govalidator.CustomTypeValidator(validatorCreds))
	govalidator.CustomTypeTagMap.Set("validTargets", govalidator.CustomTypeValidator(validTargets))
	govalidator.CustomTypeTagMap.Set("validPositive", govalidator.CustomTypeValidator(validPositive))
	govalidator.CustomTypeTagMap.Set("optionalBetween", govalidator.CustomTypeValidator(validOptionalBetween))
}

// ExtractStructFromRequest Unmarshal Request body to inputStructPtr struct
func ExtractStructFromRequest(r *http.Request, inputStructPtr interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("can't read input data: %v, err: %v", string(b), err)
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			logger.Log.WarnfCtx(r.Context(), "cannot close request.Body err: %s", err)
		}
	}()

	logger.Log.DebugfCtx(r.Context(), "ExtractStructFromRequest: payload received %s for request from URL %s", string(b), r.URL.String())

	// decode input data
	if err = json.Unmarshal(b, inputStructPtr); err != nil {
		return fmt.Errorf("can't decode input data: %v, err: %v", string(b), err)
	}

	switch inputStructPtr.(type) {
	case *models.AllTargetsEnable:
		err = validateAllTargetsEnable(string(b))
	default:
		err = ValidateByCustomValidators(inputStructPtr)
	}

	if err != nil {
		return fmt.Errorf("can't validate input data: %v, err: %v", string(b), err)
	}
	return nil
}

func ValidateByCustomValidators(s interface{}) error {
	SetupCustomValidators()
	_, err := govalidator.ValidateStruct(s)
	return err
}

func validateAllTargetsEnable(inputJSON string) error {
	if !strings.Contains(inputJSON, `"active":`) {
		return fmt.Errorf("can't decode input data: %s. \"active\" field is required", inputJSON)
	}
	return nil
}

// ValidateParametersField is used to validate Parameters field
func ValidateParametersField(jsonSchema string, parameters string, strict bool) error {

	if len(jsonSchema) == 0 {
		if len(parameters) != 0 {
			return fmt.Errorf("parameters for non-parameterized script should be empty, but got: %s", parameters)
		}
		return nil
	}

	// to prevent error while validating empty json
	if len(parameters) == 0 {
		parameters = "{}"
	}

	jsonSchema = fmt.Sprint(jsonSchema)
	if !strict {
		jsonSchema = parametersPattern.ReplaceAllString(jsonSchema, `"required":[]`)
	}

	// use fmt.Sprint for de-escaping input strings before validation
	schemaLoader := gojsonschema.NewStringLoader(jsonSchema)
	inputLoader := gojsonschema.NewStringLoader(fmt.Sprint(parameters))

	result, err := gojsonschema.Validate(schemaLoader, inputLoader)
	if err != nil {
		return fmt.Errorf("parameters validation error: %v", err)
	}
	if !result.Valid() {
		return fmt.Errorf("parameters are not valid, err: %v", result.Errors())
	}
	return nil
}
