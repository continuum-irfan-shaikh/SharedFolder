package wal

import (
	"fmt"
	"reflect"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/util"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Default Configuration", want: &Config{Name: fmt.Sprintf("./wal/%s.wal", util.ProcessName()),
				MaxSegments: 100, SegmentSizeInKB: 2, MaxAgeDays: 10, MaxFiles: 20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
