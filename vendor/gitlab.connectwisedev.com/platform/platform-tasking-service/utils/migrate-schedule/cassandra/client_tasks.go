package cassandra

import (
	"encoding/json"
	"fmt"
	apiModels "gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/tasking"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TargetType is used for targets definition
type TargetType int

// These constants describe the type of entities the task was run on
const (
	_ TargetType = iota
	managedEndpoint
	dynamicGroup
)

const (
	selectAllTasksCqlQuery = `SELECT partner_id, id, managed_endpoint_id, external_task, schedule, regularity, run_time, start_run_time, end_run_time, location, TTL(regularity) FROM tasks`
	updateTaskCqlQuery     = `UPDATE tasks USING TTL ? SET schedule = ?, old_schedule = ? WHERE partner_id = ? AND external_task = ? AND managed_endpoint_id = ? AND id = ?`
)

var (
	hourlyRE  = regexp.MustCompile(`^\d+ \* \* \* \*\|@every (\d+)h$`)
	dailyRE   = regexp.MustCompile(`^(\d+) (\d+) \* \* \*\|@every (\d+)d$`)
	weeklyRE  = regexp.MustCompile(`^(\d+) (\d+) \* \* ([0-6,]+)\|@every (\d+)w$`)
	monthlyRE = regexp.MustCompile(`^(\d+) (\d+) ([\d,]+) \* \*\|@every (\d+)m$`)
)

// TaskFields is a set of Primary Keys of tasks table
type TaskFields struct {
	PartnerID         string
	ID                gocql.UUID
	ManagedEndpointID gocql.UUID
	ExternalTask      bool
	Schedule          string
	Regularity        apiModels.Regularity
	RunTime           time.Time
	StartRunTime      time.Time
	EndRunTime        time.Time
	Location          string
	TTL               int
}

func selectAllTasks() ([]TaskFields, error) {
	var (
		cassandraQuery = Session.Query(selectAllTasksCqlQuery)
		taskFields     TaskFields
		taskFieldsList []TaskFields
		iter           = cassandraQuery.Iter()
	)
	for iter.Scan(
		&taskFields.PartnerID,
		&taskFields.ID,
		&taskFields.ManagedEndpointID,
		&taskFields.ExternalTask,
		&taskFields.Schedule,
		&taskFields.Regularity,
		&taskFields.RunTime,
		&taskFields.StartRunTime,
		&taskFields.EndRunTime,
		&taskFields.Location,
		&taskFields.TTL,
	) {
		taskFieldsList = append(
			taskFieldsList,
			TaskFields{
				PartnerID:         taskFields.PartnerID,
				ID:                taskFields.ID,
				ManagedEndpointID: taskFields.ManagedEndpointID,
				ExternalTask:      taskFields.ExternalTask,
				Schedule:          taskFields.Schedule,
				Regularity:        taskFields.Regularity,
				RunTime:           taskFields.RunTime,
				StartRunTime:      taskFields.StartRunTime,
				EndRunTime:        taskFields.EndRunTime,
				Location:          taskFields.Location,
				TTL:               taskFields.TTL,
			},
		)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.Wrapf(err, "can't perform query '%s'", selectAllTasksCqlQuery)
	}

	return taskFieldsList, nil
}

func updateTasks(taskFieldsList []TaskFields) (resErr error) {
	for _, task := range taskFieldsList {
		location, err := time.LoadLocation(task.Location)
		if err != nil {
			location = time.UTC
		}

		newSchedule, errExec := ParseSchedule(task.Schedule, location)
		if errExec != nil {
			resErr = fmt.Errorf("\ntask:%v, err:%s. %s", task, errExec, resErr)
			continue
		}

		if newSchedule == nil {
			newSchedule = &apiModels.Schedule{}
			if task.Regularity == apiModels.Recurrent {
				newSchedule.EndRunTime = task.EndRunTime
				newSchedule.Repeat = apiModels.Repeat{}
				newSchedule.Repeat.RunTime = task.RunTime
			}
		}

		newSchedule.Regularity = task.Regularity
		newSchedule.StartRunTime = task.StartRunTime
		newSchedule.EndRunTime = task.EndRunTime
		newSchedule.Location = task.Location

		scheduleBytes, errExec := json.Marshal(*newSchedule)
		if errExec != nil {
			resErr = fmt.Errorf("\ntask:%v, err:%s. %s", task, errExec, resErr)
			continue
		}

		errExec = Session.Query(updateTaskCqlQuery, task.TTL, string(scheduleBytes), task.Schedule, task.PartnerID, task.ExternalTask, task.ManagedEndpointID, task.ID).Exec()
		if errExec != nil {
			resErr = fmt.Errorf("\ntask:%v, err:%s. %s", task, errExec, resErr)
		}
	}

	return resErr
}

// UpdateTasksTable sets external_task to False for all existing tasks
func UpdateTasksTable() error {
	taskFieldsList, err := selectAllTasks()
	if err != nil {
		return err
	}
	if err = updateTasks(taskFieldsList); err != nil {
		return errors.Wrap(err, "error while updating tasks")
	}
	return nil
}

//ParseSchedule parses string that represents cron schedule
func ParseSchedule(cron string, location *time.Location) (*apiModels.Schedule, error) {
	var (
		schedule      = &apiModels.Schedule{}
		scheduleSlice = strings.Split(cron, "|")
		err           error
	)
	schedule.Repeat = apiModels.Repeat{}

	switch len(scheduleSlice) {
	case 0, 1:
		return nil, nil
	case 2:
	case 3:
		if scheduleSlice[2] != "" {
			schedule.Repeat.Period, err = strconv.Atoi(scheduleSlice[2])
			if err != nil {
				return nil, err
			}
		}
		cron = strings.Join(scheduleSlice[:2], "|")
	default:
		return nil, fmt.Errorf("unknown format of cron string %s", cron)

	}

	var (
		// default cron values
		cronMinutes = 0
		cronHours   = 10
	)

	if monthlyRE.MatchString(cron) {
		//no needs to handle error here, this work already done by MatchString method
		cronParts := monthlyRE.FindStringSubmatch(cron)
		schedule.Repeat.Frequency = apiModels.Monthly
		cronMinutes, _ = strconv.Atoi(cronParts[1])
		cronHours, _ = strconv.Atoi(cronParts[2])
		schedule.Repeat.DaysOfMonth, err = splitToInts(cronParts[3], ",")
		schedule.Repeat.Every, _ = strconv.Atoi(cronParts[4])

	} else if weeklyRE.MatchString(cron) {
		cronParts := weeklyRE.FindStringSubmatch(cron)
		schedule.Repeat.Frequency = apiModels.Weekly
		cronMinutes, _ = strconv.Atoi(cronParts[1])
		cronHours, _ = strconv.Atoi(cronParts[2])
		schedule.Repeat.DaysOfWeek, err = splitToInts(cronParts[3], ",")
		schedule.Repeat.Every, _ = strconv.Atoi(cronParts[4])

	} else if dailyRE.MatchString(cron) {
		cronParts := dailyRE.FindStringSubmatch(cron)
		schedule.Repeat.Frequency = apiModels.Daily
		cronMinutes, _ = strconv.Atoi(cronParts[1])
		cronHours, _ = strconv.Atoi(cronParts[2])
		schedule.Repeat.Every, _ = strconv.Atoi(cronParts[3])

	} else if hourlyRE.MatchString(cron) {
		cronParts := dailyRE.FindStringSubmatch(cron)
		schedule.Repeat.Frequency = apiModels.Hourly
		cronMinutes, _ = strconv.Atoi(cronParts[1])
		schedule.Repeat.Every, _ = strconv.Atoi(cronParts[1])
	} else {
		return nil, fmt.Errorf("unknown format of cron string %s", cron)
	}

	schedule.Repeat.RunTime = calculateRepeatRunTime(location, cronHours, cronMinutes)
	schedule.Regularity = apiModels.Recurrent

	return schedule, err
}

func calculateRepeatRunTime(location *time.Location, hours, minutes int) time.Time {
	const (
		defaultYear     = 2009
		defaultMonth    = time.November
		defaultDay      = 10
		defaultSecs     = 0
		defaultNanosecs = 0
	)
	return time.Date(defaultYear, defaultMonth, defaultDay, hours, minutes, defaultSecs, defaultNanosecs, location)
}

func splitToInts(s, sep string) (ints []int, err error) {
	for _, v := range strings.Split(s, sep) {

		if val, err := strconv.Atoi(v); err == nil {
			ints = append(ints, val)
		} else {
			return nil, err
		}

	}
	return
}
