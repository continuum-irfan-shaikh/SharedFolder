package cherwell

// Logger is an interface for cherwell client logging
type Logger interface {
	Log(trID string, message string, v ...interface{})
}

// NopLogger is a logger that does not log
type NopLogger struct{}

// NewNopLogger is a constructor for NopLogger
func NewNopLogger() Logger {
	return &NopLogger{}
}

// Log implementation of Logger interface for NopLogger
func (nl *NopLogger) Log(trID string, message string, v ...interface{}) {}
