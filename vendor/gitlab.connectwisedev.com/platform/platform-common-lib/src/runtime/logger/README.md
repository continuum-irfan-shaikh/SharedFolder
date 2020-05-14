<p align="center">
<img height=70px src="docs/images/continuum-logo.png">
<img height=70px src="docs/images/Go-Logo_Blue.png">
</p>

# Logger

This is a Standard logger implementation used by all the Go projects in the Continuum. So that we can implement different serach patern in the Graylog and generate alerts on any anomoly.

### Third-Party Libraties

- [lumberjack](https://gopkg.in/natefinch/lumberjack.v2) - **License** [MIT License](https://github.com/natefinch/lumberjack/blob/v2.0/LICENSE) - **Description** - Lumberjack is intended to be one part of a logging infrastructure. It is not an all-in-one solution, but instead is a pluggable component at the bottom of the logging stack that simply controls the files to which logs are written.
  Lumberjack assumes that only one process is writing to the output files. Using the same lumberjack configuration from multiple processes on the same machine will result in improper behavior.

### [Example](example/example.go)

**Import Statement**

```go
import	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
```

**Logger Instance**

```go
// Create Logger instance
log, err := logger.Create(logger.Config{Name: name, MaxSize: 1})

// Update logger instance
log, err := logger.Create(logger.Config{Name: name, MaxSize: 1})
```

**Writing Logs**

```go
log.Trace(transaction, "This is a TRACE Message")
log.Debug(transaction, "This is a DEBUG Message")
log.Info(transaction, "This is a INFO Message")
log.Warn(transaction, "This is a WARN Message")
log.Error(transaction, "ERROR-CODE", "This is a ERROR Message")
log.Fatal(transaction, "FATAL-CODE", "This is a FATAL Message")
```

**Helper functions**

```go
// Return an instance of internal io.writer used by this logger
log.GetWriter() io.WriteCloser

// Set a instance of internal io.writer used by this logger
log.SetWriter(writer io.Writer)

// Current log level of a Logger
log.LogLevel() LogLevel
```

**Configuration**

```go
// Config is a struct to hold logger configuration
type Config struct {
	// Name is a user defined name of logger, this is used to handle multiple log files.
	// The default is <processname>, if empty
	Name string `json:"name"`

	// FileName is the file to write logs to and backup log files will be retained in the same directory.
	// It uses <processname>-Name.log in the same directory where process binary available, if empty.
	FileName string `json:"filename"`

	// MaxSize is the maximum size in megabytes of the log file before it gets rotated.
	//It defaults to 20 megabytes.
	MaxSize int `json:"maxsize"`

	// MaxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename.
	// The default is 30 Days to remove old log files based on age.
	MaxAge int `json:"maxage"`

	// MaxBackups is the maximum number of old log files to retain.
	// The default is to retain 5 old log files (though MaxAge may still cause them to get deleted)
	MaxBackups int `json:"maxbackups"`

	// ServiceName is a user defined name of Service.
	// The default is <processname>, if empty
	ServiceName string `json:"servicename"`

	// Filler is a string used to fillup required logger attribute in case these are not available
	// The default is -
	Filler string `json:"filler"`

	// LogLevel is a allowed log level to write logs
	// The default is INFO
	LogLevel LogLevel `json:"logLevel"`

	// Destination is a location to write logs
	// The default is FILE
	Destination Destination `json:"destination"`

	// CallDepth is a depth for runtime to find a caller file name
	// The default value is 4 because current loger has 3 function layes on top of the runtime
	CallDepth int
}
```

## Example

```json
"Logging": {
	"LogLevel": "INFO",
	"Destination": "FILE",
	"MaxSize": 20,
	"MaxAge": 30,
	"MaxBackups": 5,
	"Filler": "",
	"ServiceName":"<processname>",
	"Name": "<processname>",
	"FileName": "<processname>-<Name>.log"
},
```

## FAQ:

**What is Transaction ID**

- This should be **business transaction id** to track complete business flow across services.
- We should have **transaction id from each of caller in the continuum eco-system**, in case API is the originator it should generate the new transaction id for subsequent usage
- Whenever you generate a log file, include the Transaction ID in the log message
- Transction should be an **UUID**
  - [Helper Functions](../../utils)
- Some of the examples for request originator are:
  - Agent in case of scheduler
  - Portal in case of user request
  - LRP in case of event handling
  - Job trigger point

**Where I can find Transaction ID**
- Read Transaction ID from your incoming requests, and if one is provided, send it on outgoing requests
  - `utils.GetTransactionIDFromRequest` retrieves transaction Id, if from the http request header, if this does not present it creates a new one
- If you don’t get a Transaction ID on incoming request, then generate one, and send it on outgoing requests
  - `utils.GetTransactionID` generates new transaction id

**ERROR/FATAL Code**
- Repository owners are free to define Error/Fatal codes as these will be used only for identifying logs in the Graylag


### Contribution

Any changes in this package should be communicated to Juno Team.
