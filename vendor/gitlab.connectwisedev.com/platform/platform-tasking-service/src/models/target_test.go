package models

import (
	"testing"
)

func TestTargetType_MarshalJSON(t *testing.T) {
	var (
		target = ManagedEndpoint
	)

	t.Run("positive", func(t *testing.T) {
		jsn, err := target.MarshalJSON()
		if err != nil {
			t.Fatalf(err.Error())
		}

		if string(jsn) != `"MANAGED_ENDPOINT"` {
			t.Fatalf("got unexpected value")
		}
	})

	t.Run("negative", func(t *testing.T) {
		var target TargetType = 10

		_, err := target.MarshalJSON()
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}

func TestTargetType_UnmarshalJSON(t *testing.T) {
	var (
		target = ManagedEndpoint
	)

	t.Run("positive", func(t *testing.T) {
		s := `"MANAGED_ENDPOINT"`
		err := target.UnmarshalJSON([]byte(s))
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("negative_1", func(t *testing.T) {
		s := `invalid`
		err := target.UnmarshalJSON([]byte(s))
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})

	t.Run("negative_2", func(t *testing.T) {
		s := `"INVALID"`
		err := target.UnmarshalJSON([]byte(s))
		if err == nil {
			t.Fatalf("error cannot be nil")
		}
	})
}
