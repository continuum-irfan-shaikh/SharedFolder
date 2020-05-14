package cassandra

import "testing"

func TestNew(t *testing.T) {
	NewTaskInstance(nil)
	NewTask(nil)
	NewExecutionExpiration(nil)
	NewLegacyMigration(nil)
	NewProfilesRepo(nil)
	NewScheduler(nil)
	NewScriptExecutionResults(nil)
	NewTaskExecutionHistory(nil)
	NewTargets(nil)
	NewExecutionExpiration(nil)
}
