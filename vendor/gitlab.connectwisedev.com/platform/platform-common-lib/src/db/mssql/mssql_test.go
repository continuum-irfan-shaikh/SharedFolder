package mssql

import (
	"testing"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/db"
)

func TestGetConnectionString(t *testing.T) {
	m := mssql{}
	t.Run("Error missing config", func(t *testing.T) {
		_, err := m.GetConnectionString(db.Config{})
		if err == nil {
			t.Errorf("Expecting error but found nil")
		}
	})

	t.Run("Success", func(t *testing.T) {
		conStr, err := m.GetConnectionString(db.Config{DbName: "NOCBO",
			Server:     "10.2.27.41",
			Password:   "its",
			UserID:     "its",
			CacheLimit: 200})
		if err != nil {
			t.Errorf("Expecting nil but found err := %v", err)
		}

		if conStr == "" {
			t.Errorf("Expecting connection string found empty")
		}
	})
}
