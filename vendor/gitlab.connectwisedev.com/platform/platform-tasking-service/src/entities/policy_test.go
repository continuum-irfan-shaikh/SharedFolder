package entities

import (
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models/triggers"
)

func TestPolicy_IsAlertingPolicy(t *testing.T) {
	p := Policy{ID: triggers.AlertTypePrefix}
	if !p.IsAlertingPolicy() {
		t.Fatalf("expected true but got  false")
	}
	
	p = Policy{ID: triggers.GenericTypePrefix}
	if p.IsAlertingPolicy() {
		t.Fatalf("expected false but got true")
	}
}
