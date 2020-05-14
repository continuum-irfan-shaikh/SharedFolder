package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/util"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var nameToLogger = make(map[string]*loggerImpl)

const discardLoggerName = "Discard-Logger"

// Log is a interface to hold instace of loggerImpl
type Log interface {
	Trace(transactionID string, message string, v ...interface{})
	Debug(transactionID string, message string, v ...interface{})
	Info(transactionID string, message string, v ...interface{})
	Warn(transactionID string, message string, v ...interface{})
	Error(transactionID string, errorCode string, message string, v ...interface{})
	Fatal(transactionID string, fatalCode string, message string, v ...interface{})
	GetWriter() io.WriteCloser
	SetWriter(writer io.Writer)
	LogLevel() LogLevel

	// TO BE REMOVED AFTER MIGRATION TO NEW LOGGER COMPLETES

	// TraceWithLevel log trace message with calldepth
	//
	// Deprecated: TraceWithLevel should not be used except for compatibility with legacy systems.
	// Instead use trace
	TraceWithLevel(transactionID string, calldepth int, message string, v ...interface{})
	// DebugWithLevel log debug message with calldepth
	//
	// Deprecated: DebugWithLevel should not be used except for compatibility with legacy systems.
	// Instead use debug
	DebugWithLevel(transactionID string, calldepth int, message string, v ...interface{})
	// InfoWithLevel log info message with calldepth
	//
	// Deprecated: InfoWithLevel should not be used except for compatibility with legacy systems.
	// Instead use info
	InfoWithLevel(transactionID string, calldepth int, message string, v ...interface{})
	// WarnWithLevel log warn message with calldepth
	//
	// Deprecated: WarnWithLevel should not be used except for compatibility with legacy systems.
	// Instead use warn
	WarnWithLevel(transactionID string, calldepth int, message string, v ...interface{})
	// ErrorWithLevel log trace message with calldepth
	//
	// Deprecated: ErrorWithLevel should not be used except for compatibility with legacy systems.
	// Instead use error
	ErrorWithLevel(transactionID string, calldepth int, errorCode string, message string, v ...interface{})
	// FatalWithLevel log trace message with calldepth
	//
	// Deprecated: FatalWithLevel should not be used except for compatibility with legacy systems.
	// Instead use fetal
	FatalWithLevel(transactionID string, calldepth int, fatalCode string, message string, v ...interface{})
}

type loggerImpl struct {
	writer   io.WriteCloser
	config   *Config
	hostName string

	mutex sync.Mutex // ensures atomic writes; protects the following fields
}

func init() {
	Create(Config{Name: discardLoggerName, Destination: DISCARD}) // nolint
}

// Create is a function to create an instance of a loggerImpl
func Create(config Config) (Log, error) {
	instance, ok := nameToLogger[config.name()]

	if ok {
		return instance, fmt.Errorf("LoggerAlreadyInitialized for name :: %s and File :: %s", config.name(), config.fileName())
	}

	l := &loggerImpl{config: &config, hostName: util.Hostname(config.filler())}
	l.setOutput()
	nameToLogger[config.name()] = l
	return l, nil
}

// Update is a function to Update an instance of a loggerImpl
// This creates a new logger instance if logger for given name does not exist
func Update(config Config) (Log, error) {
	instance, ok := nameToLogger[config.name()]

	if !ok {
		return Create(config)
	}

	instance.config = &config
	instance.setOutput()
	return instance, nil
}

// GetConfig is a function to return an instance of a configuration used for soecified loggerImpl
func GetConfig(name string) Config {
	if name == "" {
		name = util.ProcessName()
	}

	instance, ok := nameToLogger[name]
	if ok {
		return *instance.config
	}
	return Config{}
}

// GetViaName is a function to return a logger instance for given Name
// This return a default logger instance having FILE as a writer for Process name
func GetViaName(name string) Log {
	if name == "" {
		name = util.ProcessName()
	}

	instance, ok := nameToLogger[name]

	if ok {
		return instance
	}

	if name == util.ProcessName() {
		instance, _ := Create(Config{}) //nolint
		return instance
	}

	panic(fmt.Errorf("logger with Name: %s does not exist", name))
}

// Get is a function to return a logger instance
// Default name will be used as a process name
func Get() Log {
	return GetViaName(util.ProcessName())
}

// DiscardLogger is a function to return a logger instance having ioutil.Discard as a writer
func DiscardLogger() Log {
	return GetViaName(discardLoggerName)
}

func (l *loggerImpl) Trace(transactionID string, message string, v ...interface{}) {
	if l.config.logLevel().order >= TRACE.order {
		l.output(l.config.calldepth(), transactionID, TRACE, fmt.Sprintf(message, v...))
	}
}

func (l *loggerImpl) Debug(transactionID string, message string, v ...interface{}) {
	if l.config.logLevel().order >= DEBUG.order {
		l.output(l.config.calldepth(), transactionID, DEBUG, fmt.Sprintf(message, v...))
	}
}

func (l *loggerImpl) Info(transactionID string, message string, v ...interface{}) {
	if l.config.logLevel().order >= INFO.order {
		l.output(l.config.calldepth(), transactionID, INFO, fmt.Sprintf(message, v...))
	}
}

func (l *loggerImpl) Warn(transactionID string, message string, v ...interface{}) {
	if l.config.logLevel().order >= WARN.order {
		l.output(l.config.calldepth(), transactionID, WARN, fmt.Sprintf(message, v...))
	}
}

func (l *loggerImpl) Error(transactionID string, errorCode string, message string, v ...interface{}) {
	if l.config.logLevel().order >= ERROR.order {
		l.output(l.config.calldepth(), transactionID, ERROR, fmt.Sprintf(errorCode+" "+message, v...))
	}
}

func (l *loggerImpl) Fatal(transactionID string, fatalCode string, message string, v ...interface{}) {
	if l.config.logLevel().order >= FATAL.order {
		l.output(l.config.calldepth(), transactionID, FATAL, fmt.Sprintf(fatalCode+" "+message, v...))
	}
}

// GetWriter is a function to return an instance of internal io.writer used by this logger
func (l *loggerImpl) GetWriter() io.WriteCloser {
	return l.writer
}

// SetWriter is a function to set a instance of internal io.writer used by this logger
func (l *loggerImpl) SetWriter(writer io.Writer) {
	//Adding locks as we only update output while creating or updating logger instance
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.writer != nil {
		l.writer.Close() //nolint
	}

	l.writer = &nopCloser{writer}
}

// LogLevel - Current log level of a Logger
func (l *loggerImpl) LogLevel() LogLevel {
	return l.config.logLevel()
}

// Helper functions used for logging
func (l *loggerImpl) setOutput() {
	//Adding locks as we only update output while creating or updating logger instance
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.writer != nil {
		l.writer.Close() //nolint
	}

	switch l.config.destination() {
	case STDOUT:
		l.writer = nopCloser{os.Stdout}
	case STDERR:
		l.writer = nopCloser{os.Stderr}
	case FILE:
		l.writer = &lumberjack.Logger{
			Filename:   l.config.fileName(),
			MaxSize:    l.config.maxSize(), // megabytes
			MaxBackups: l.config.maxBackups(),
			MaxAge:     l.config.maxAge(), //days
			Compress:   true,              // disabled by default
			LocalTime:  false,             // UTC by default
		}
	default:
		l.writer = nopCloser{ioutil.Discard}
	}
}

func (l *loggerImpl) output(calldepth int, transactionID string, level LogLevel, message string) {
	now := time.Now().UTC() // get this early and always use UTC for logging
	buf := l.formatHeader(now, calldepth, transactionID, level)
	buf = append(buf, message...)

	if len(message) == 0 || message[len(message)-1] != '\n' {
		buf = append(buf, '\n')
	}

	l.writer.Write(buf) //nolint Ignoring this error as we dont want to handle
}

func (l *loggerImpl) formatHeader(t time.Time, calldepth int, transactionID string, level LogLevel) []byte {
	buf := make([]byte, 0)
	buf = append(buf, l.formatTime(t)...)
	buf = append(buf, ' ')
	buf = append(buf, l.hostName...)
	buf = append(buf, ' ')
	buf = append(buf, l.config.serviceName()...)
	buf = append(buf, ' ')
	buf = append(buf, transactionID...)
	buf = append(buf, ' ')
	buf = append(buf, l.formatFileName(calldepth)...)
	buf = append(buf, ' ')
	buf = append(buf, level.name...)
	buf = append(buf, ' ')
	return buf
}

func (l *loggerImpl) formatTime(t time.Time) []byte {
	buf := make([]byte, 0)

	// Add year, month, day
	year, month, day := t.Date()
	l.itoa(&buf, year, 4)
	buf = append(buf, '/')
	l.itoa(&buf, int(month), 2)
	buf = append(buf, '/')
	l.itoa(&buf, day, 2)
	buf = append(buf, ' ')

	// Add hour, min, sec
	hour, min, sec := t.Clock()
	l.itoa(&buf, hour, 2)
	buf = append(buf, ':')
	l.itoa(&buf, min, 2)
	buf = append(buf, ':')
	l.itoa(&buf, sec, 2)

	// Add microseconds
	buf = append(buf, '.')
	l.itoa(&buf, t.Nanosecond()/1e3, 6)

	return buf
}

func (l *loggerImpl) formatFileName(calldepth int) []byte {
	buf := make([]byte, 0)
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "????"
		line = 0
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	buf = append(buf, short...)
	buf = append(buf, ':')
	l.itoa(&buf, line, -1)
	buf = append(buf, ':')
	return buf
}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
func (l *loggerImpl) itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

// TO BE REMOVED AFTER MIGRATION TO NEW LOGGER COMPLETES
func (l *loggerImpl) TraceWithLevel(transactionID string, calldepth int, message string, v ...interface{}) {
	if l.config.logLevel().order >= TRACE.order {
		l.output(l.config.calldepth()+calldepth, transactionID, TRACE, fmt.Sprintf(message, v...))
	}
}

func (l *loggerImpl) DebugWithLevel(transactionID string, calldepth int, message string, v ...interface{}) {
	if l.config.logLevel().order >= DEBUG.order {
		l.output(l.config.calldepth()+calldepth, transactionID, DEBUG, fmt.Sprintf(message, v...))
	}
}

func (l *loggerImpl) InfoWithLevel(transactionID string, calldepth int, message string, v ...interface{}) {
	if l.config.logLevel().order >= INFO.order {
		l.output(l.config.calldepth()+calldepth, transactionID, INFO, fmt.Sprintf(message, v...))
	}
}

func (l *loggerImpl) WarnWithLevel(transactionID string, calldepth int, message string, v ...interface{}) {
	if l.config.logLevel().order >= WARN.order {
		l.output(l.config.calldepth()+calldepth, transactionID, WARN, fmt.Sprintf(message, v...))
	}
}

func (l *loggerImpl) ErrorWithLevel(transactionID string, calldepth int, errorCode string, message string, v ...interface{}) {
	if l.config.logLevel().order >= ERROR.order {
		l.output(l.config.calldepth()+calldepth, transactionID, ERROR, fmt.Sprintf(errorCode+" "+message, v...))
	}
}

func (l *loggerImpl) FatalWithLevel(transactionID string, calldepth int, fatalCode string, message string, v ...interface{}) {
	if l.config.logLevel().order >= FATAL.order {
		l.output(l.config.calldepth()+calldepth, transactionID, FATAL, fmt.Sprintf(fatalCode+" "+message, v...))
	}
}
