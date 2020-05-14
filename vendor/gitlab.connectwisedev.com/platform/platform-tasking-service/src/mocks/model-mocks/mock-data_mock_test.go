package modelMocks

import (
	"reflect"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
)

func TestGenerateTargetsMock(t *testing.T) {
	tests := []struct {
		name     string
		want     map[string]bool
		dontWant map[string]bool
	}{
		{
			name:     "case1",
			want:     map[string]bool{TargetIDStr: true},
			dontWant: map[string]bool{TargetIDStr: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateTargetsMock(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GenerateTargetsMock() = %v, want %v", got, tt.want)
			}
			if got := GenerateTargetsMock(); reflect.DeepEqual(got, tt.dontWant) {
				t.Fatalf("GenerateTargetsMock() = %v must not be equal to %v", got, tt.dontWant)
			}
		})
	}
}

func TestGenerateDisabledTargetsMock(t *testing.T) {
	tests := []struct {
		name     string
		want     map[string]bool
		dontWant map[string]bool
	}{
		{
			name:     "case1",
			want:     map[string]bool{TargetIDStr: false},
			dontWant: map[string]bool{TargetIDStr: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateDisabledTargetsMock(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GenerateDisabledTargetsMock() = %v, want %v", got, tt.want)
			}
			if got := GenerateDisabledTargetsMock(); reflect.DeepEqual(got, tt.dontWant) {
				t.Fatalf("GenerateDisabledTargetsMock() = %v must not be equal to %v", got, tt.dontWant)
			}
		})
	}
}

func TestPopulateTaskExecutionResults(t *testing.T) {
	defer func(t *testing.T) {
		if err := recover(); err == nil {
			t.Errorf("Error expected when a bad fixture data is being parsed")
		}
	}(t)

	var (
		executionResultsView []models.ExecutionResultView
		fixture              = []byte(`bad json`)
	)
	populateTaskExecutionResultsView(&executionResultsView, fixture)
}
