package zookeeper

import "testing"

func TestJob(t *testing.T) {
	job := Job{
		Name:     "Name",
		Schedule: "Schedule",
		Task:     "Task",
	}

	if job.GetName() != job.Name {
		t.Errorf("expected name: %s, got: %s", job.Name, job.GetName())
	}
	if job.GetSchedule() != job.Schedule {
		t.Errorf("expected schedule: %s, got: %s", job.Schedule, job.GetSchedule())
	}
	if job.GetTask() != job.Task {
		t.Errorf("expected task: %s, got: %s", job.Task, job.GetTask())
	}
}
