<p align="center">
<img height=70px src="docs/images/continuum-logo.png">
<img height=70px src="docs/images/Go-Logo_Blue.png">
</p>

# WAL

This is a Standard wal implementation used by all the Go projects in the Continuum.

### Third-Party Libraties

- [lumberjack](https://gopkg.in/natefinch/lumberjack.v2) 
  - **License** [MIT License](https://github.com/natefinch/lumberjack/blob/v2.0/LICENSE) 
  - **Description** - Lumberjack is intended to be one part of a logging infrastructure. It is not an all-in-one solution, but instead is a pluggable component at the bottom of the logging stack that simply controls the files to which logs are written.
  Lumberjack assumes that only one process is writing to the output files. Using the same lumberjack configuration from multiple processes on the same machine will result in improper behavior.

### [Example](example/example.go)

**Import Statement**

```go
import	"gitlab.connectwisedev.com/platform/platform-common-lib/src/wal"
```

**Configuration**
```go
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
```
**Wal Instance**
```go
// Create - returns instance of a wal service
wal.Create(conf *Config) WAL
```

**Wal Interface**
```go
// WAL - Interface to perform action on wal file
type WAL interface {
	// Write - write messages in the wal file
	Write(...string) error

	// WriteObject - serialize object and write this into wal file
	WriteObject(...interface{}) error

	//Read - x wal files, and convert all the enrties in to records
	Read(batchSize int) ([]Record, error)

	// Flush - Commit / Lock / Rotate current wal file, and
	// create a new one for wal messages
	Flush() error
}
```


### Contribution

Any changes in this package should be communicated to Juno Team.
