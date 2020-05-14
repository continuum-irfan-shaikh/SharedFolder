package logger

//go:generate mockgen -destination=../mocks/mock-logger/logger_mock.go -package=mocks gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger Log
//go:generate mockgen -destination=../mocks/mockLoggerTasking/loggerTasking_mock.go -package=mockLoggerTasking -source=./logger.go

import (
	"context"
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/middlewares/transaction-id"
)

type Logger interface {
	InfofCtx(ctx context.Context, format string, v ...interface{})
	WarnfCtx(ctx context.Context, format string, v ...interface{})
	ErrfCtx(ctx context.Context, errorcode, format string, v ...interface{})
	DebugfCtx(ctx context.Context, format string, v ...interface{})

	Debug(v ...interface{}) // used by Cassandra logger
}

type logWrapper struct {
	log logger.Log
}

// InfofCtx - logging (level info) with formation, wrap message with transactionID
func (l *logWrapper) InfofCtx(ctx context.Context, format string, v ...interface{}) {
	l.log.Info(l.transactionID(ctx), format, v...)
}

// WarnfCtx - logging (level warning) with formation, wrap message with transactionID
func (l *logWrapper) WarnfCtx(ctx context.Context, format string, v ...interface{}) {
	l.log.Warn(l.transactionID(ctx), format, v...)
}

// ErrfCtx - logging (level error) with formation, wrap message with transactionID
func (l *logWrapper) ErrfCtx(ctx context.Context, errorcode, format string, v ...interface{}) {
	l.log.Error(l.transactionID(ctx), errorcode, format, v...)
}

// Debug - logging (level debug)
func (l *logWrapper) Debug(v ...interface{}) {
	l.log.Debug("", fmt.Sprint(v...))
}

// Debugln - logging (level debug) with new line
func (l *logWrapper) Debugln(v ...interface{}) {
	l.log.Debug("", fmt.Sprintln(v...))
}

// DebugfCtx - logging (level debug) with formation transactionID
func (l *logWrapper) DebugfCtx(ctx context.Context, format string, v ...interface{}) {
	l.log.Debug(l.transactionID(ctx), format, v...)
}

func (l *logWrapper) transactionID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	id, ok := ctx.Value(transactionID.Key).(string)
	if !ok {
		return ""
	}

	return id
}

//ZooKeeperLogger struct implements ZooKeeper `logger` interface
type ZooKeeperLogger struct{}

// LogInfo writes the message to the log file with INFO LogLevel
func (ZooKeeperLogger) LogInfo(format string, v ...interface{}) {
	Log.InfofCtx(context.Background(), format, v...)
}

// LogError writes the message to the log file with INFO LogLevel
func (ZooKeeperLogger) LogError(format string, v ...interface{}) {
	Log.ErrfCtx(context.Background(), "", format, v...)
}

//CassandraLoggerT struct for cassandra logger
type CassandraLoggerT struct {
	logger.Log
}

//Print for cassandra. Use DEBUG to use logger only for debug
func (c CassandraLoggerT) Print(v ...interface{}) {
	Log.Debug(v)
}

//Printf for cassandra. Use DEBUG to use logger only for debug
func (c CassandraLoggerT) Printf(format string, v ...interface{}) {
	Log.DebugfCtx(context.Background(), format, v)
}

//Println for cassandra. Use DEBUG to use logger only for debug
func (c CassandraLoggerT) Println(v ...interface{}) {
	Log.DebugfCtx(context.Background(), "", v)
}

var (
	// Log used for logging through app in DI manner
	Log Logger
	// CassandraLogger is impl of logger
	CassandraLogger CassandraLoggerT
	// ZKLogger is impl of logger
	ZKLogger ZooKeeperLogger
)

func Load(cfg logger.Config) error {
	l, err := logger.Update(cfg)
	if err != nil {
		return err
	}
	CassandraLogger = CassandraLoggerT{l}
	Log = &logWrapper{log: l}

	return nil
}
