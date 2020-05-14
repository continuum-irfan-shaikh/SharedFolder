package wal

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/util"
)

// Config - Wal configuration
type Config struct {
	// Name - User defined name of wal, this is used to handle multiple wal files.
	// The default is <processname>.wal, if empty
	Name string

	// MaxSegments - Max number of sagements to be published in the Wal
	// The default is 100 Segments
	MaxSegments int

	// SegmentSizeInKB - A Segment size; this wil be used to calculate max file size for wal
	// The default is 2 KB
	SegmentSizeInKB int

	// MaxAgeDays - Maximum number of days to retain old log files based on the timestamp encoded in their filename.
	// The default is 30 Days to remove old log files based on age.
	MaxAgeDays int

	// MaxFiles - Maximum number of wal files to retain.
	// The default is to retain 5 wal files (though MaxAgeDays may still cause them to get deleted)
	MaxFiles int
}

// NewConfig - Return a new wal configuration object with default values
func NewConfig() *Config {
	return &Config{
		Name:            fmt.Sprintf("./wal/%s.wal", util.ProcessName()),
		MaxSegments:     100,
		SegmentSizeInKB: 2,
		MaxAgeDays:      10,
		MaxFiles:        20,
	}
}
