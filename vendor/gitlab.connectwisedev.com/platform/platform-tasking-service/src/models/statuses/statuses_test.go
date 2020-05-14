package statuses

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	var stat TaskState
	stat.MarshalJSON()
	stat.Parse("")
	stat.Parse("Active")
	stat.MarshalJSON()
	stat.Parse("Inactive")
	stat.MarshalJSON()
	stat.Parse("Disabled")
	stat.MarshalJSON()
	stat.Parse("Wrong")
	stat.MarshalJSON()

	s := taskInstanceStatuses
	fmt.Println(s)
	ss := taskInstanceStatusText
	fmt.Println(ss)

	consta := TaskInstanceRunning
	consta = TaskInstanceSuccess
	consta = TaskInstanceFailed
	consta = TaskInstanceScheduled
	consta = TaskInstanceDisabled
	consta = TaskInstanceSomeFailures
	consta = TaskInstancePending
	consta = TaskInstanceStopped
	consta = TaskInstancePostponed
	consta = TaskInstanceCanceled
	fmt.Println(consta)

	constaText := TaskInstanceRunningText
	constaText = TaskInstanceSuccessText
	constaText = TaskInstanceFailedText
	constaText = TaskInstanceScheduledText
	constaText = TaskInstanceDisabledText
	constaText = TaskInstanceSomeFailuresText
	constaText = TaskInstancePendingText
	constaText = TaskInstanceStoppedText
	constaText = TaskInstancePostponedText
	constaText = TaskInstanceCanceledText
	fmt.Println(constaText)

	TaskInstanceStatusText(consta)
	taskInstanceStatusText = map[TaskInstanceStatus]string{}
	TaskInstanceStatusText(consta)

	TaskInstanceStatusFromText("Success")
	TaskInstanceStatusFromText("waea")
	CalculateForStartedTask(1, 1, 0)
	CalculateForStartedTask(1, 0, 1)
	CalculateForStartedTask(2, 1, 1)
	CalculateForStartedTask(1, 0, 0)

	o := OverallSuccess
	o = OverallFailed
	o = OverallPartialFailed
	o = OverallNew
	o = OverallScheduled
	o = OverallSuspended
	o = OverallRunning
	fmt.Println(o)
}

func TestTaskInstanceStatusFromText(t *testing.T) {
	t.Run("ok Upper case", func(t *testing.T) {
		status, err := TaskInstanceStatusFromText("SUCCESS")
		if err != nil {
			t.Fatalf("expected nil but got %v", err)
		}

		if status != TaskInstanceSuccess {
			t.Fatalf("expected nil but got %v", err)
		}
	})

	t.Run("wrong status", func(t *testing.T) {
		_, err := TaskInstanceStatusFromText("Not okay")
		if err == nil {
			t.Fatal("expected err but got nil")
		}
	})

	t.Run("ok success", func(t *testing.T) {
		status, err := TaskInstanceStatusFromText("Success")
		if err != nil {
			t.Fatalf("expected nil but got %v", err)
		}

		if status != TaskInstanceSuccess {
			t.Fatalf("expected nil but got %v", err)
		}
	})

	t.Run("ok failed", func(t *testing.T) {
		status, err := TaskInstanceStatusFromText("failed")
		if err != nil {
			t.Fatalf("expected nil but got %v", err)
		}

		if status != TaskInstanceFailed {
			t.Fatalf("expected nil but got %v", err)
		}
	})

	t.Run("empty", func(t *testing.T) {
		_, err := TaskInstanceStatusFromText("Fa")
		if err == nil {
			t.Fatal("expected err but got nil")
		}
	})
}