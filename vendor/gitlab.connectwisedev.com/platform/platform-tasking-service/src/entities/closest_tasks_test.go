package entities

import (
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/statuses"
)

func TestPreparePendingStatuses(t *testing.T) {
	t.Run("has running - transformed", func(t *testing.T) {
		ti := TaskInstance{
			Statuses: map[string]statuses.TaskInstanceStatus{
				"1": statuses.TaskInstanceRunning,
				"2": statuses.TaskInstanceScheduled,
			},
		}

		ti.PreparePendingStatuses()
		if ti.Statuses["2"] != statuses.TaskInstancePending {
			t.Fatalf("expected pending status but got %v", ti.Statuses["2"])
		}
	})

	t.Run("has running - transformed", func(t *testing.T) {
		ti := TaskInstance{
			Statuses: map[string]statuses.TaskInstanceStatus{
				"1": statuses.TaskInstanceFailed,
				"2": statuses.TaskInstanceScheduled,
			},
		}

		ti.PreparePendingStatuses()
		if ti.Statuses["2"] != statuses.TaskInstancePending {
			t.Fatalf("expected pending status but got %v", ti.Statuses["2"])
		}
	})

	t.Run("has running - transformed", func(t *testing.T) {
		ti := TaskInstance{
			Statuses: map[string]statuses.TaskInstanceStatus{
				"1": statuses.TaskInstanceScheduled,
			},
		}

		ti.PreparePendingStatuses()
		if ti.Statuses["1"] != statuses.TaskInstanceScheduled {
			t.Fatalf("expected scheduled status but got %v", ti.Statuses["2"])
		}
	})
}

func TestTaskInstance(t *testing.T) {
	t.Run("calculate statuses", func(t *testing.T) {
		ti := TaskInstance{}
		_, err := ti.CalculateStatuses()
		if err == nil {
			t.Errorf("expected err")
		}

		// wrong status
		ti.Statuses = map[string]statuses.TaskInstanceStatus{
			"id": statuses.TaskInstanceStatus(999),
		}
		_, err = ti.CalculateStatuses()
		if err == nil {
			t.Errorf("expected err")
		}

		ti.Statuses = map[string]statuses.TaskInstanceStatus{
			"id": statuses.TaskInstancePending,
		}
		got, err := ti.CalculateStatuses()
		if err != nil {
			t.Errorf("expected err")
		}

		if got[statuses.TaskInstancePendingText] == 0 {
			t.Errorf("expected pending status != 0")
		}
	})

	testCases := []struct {
		statuses map[string]statuses.TaskInstanceStatus
		expected statuses.OverallStatus
	}{
		{
			expected: statuses.OverallNew,
		},
		{
			statuses: map[string]statuses.TaskInstanceStatus{
				"id": statuses.TaskInstanceRunning,
			},
			expected: statuses.OverallRunning,
		},
		{
			statuses: map[string]statuses.TaskInstanceStatus{
				"id": statuses.TaskInstancePending,
			},
			expected: statuses.OverallRunning,
		},
		{
			statuses: map[string]statuses.TaskInstanceStatus{
				"id": statuses.TaskInstanceScheduled,
			},
			expected: statuses.OverallNew,
		},
		{
			statuses: map[string]statuses.TaskInstanceStatus{
				"id": statuses.TaskInstanceFailed,
			},
			expected: statuses.OverallFailed,
		},
		{
			statuses: map[string]statuses.TaskInstanceStatus{
				"id": statuses.TaskInstanceSuccess,
			},
			expected: statuses.OverallSuccess,
		},
		{
			statuses: map[string]statuses.TaskInstanceStatus{
				"id":  statuses.TaskInstanceSuccess,
				"id2": statuses.TaskInstanceFailed,
			},
			expected: statuses.OverallPartialFailed,
		},
		{
			statuses: map[string]statuses.TaskInstanceStatus{
				"id": statuses.TaskInstancePostponed,
			},
			expected: statuses.OverallSuspended,
		},
		{
			statuses: map[string]statuses.TaskInstanceStatus{
				"id": statuses.TaskInstanceDisabled,
			},
			expected: statuses.OverallSuspended,
		},
		{
			statuses: map[string]statuses.TaskInstanceStatus{
				"id":  statuses.TaskInstanceSuccess,
				"id2": statuses.TaskInstanceSomeFailures,
			},
			expected: statuses.OverallPartialFailed,
		},
		{
			statuses: map[string]statuses.TaskInstanceStatus{
				"id":  statuses.TaskInstanceCanceled,
				"id2": statuses.TaskInstanceSuccess,
			},
			expected: statuses.OverallSuspended,
		},
		{
			statuses: map[string]statuses.TaskInstanceStatus{
				"id": statuses.TaskInstanceStatus(999),
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		ti := TaskInstance{}
		ti.Statuses = tc.statuses

		got := ti.CalculateOverallStatus()
		if got != tc.expected {
			t.Errorf("expected %v but got %v", tc.expected, got)
		}
	}

}
