package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/util"
)

func TestCreate(t *testing.T) {
	t.Run("Create Instance", func(t *testing.T) {
		_, err := Create(Config{})
		if err != nil {
			t.Errorf("Create() error = %v, wantErr %v", err, nil)
			return
		}
		_, err = Create(Config{})
		if err == nil {
			t.Errorf("Create() error = %v, wantErr %v", nil, "LoggerAlreadyInitialized")
			return
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update Existing Or Create", func(t *testing.T) {
		got, err := Update(Config{Name: "Test"})
		if err != nil {
			t.Errorf("Update()  error = %v, wantErr %v", err, nil)
			return
		}

		u, err := Update(Config{Name: "Test", Destination: STDOUT})
		if err != nil {
			t.Errorf("Update()  error = %v, wantErr %v", err, nil)
			return
		}

		if !reflect.DeepEqual(got, u) {
			t.Errorf("Update() = %v, want %v", u, got)
		}

		u, err = Update(Config{Name: "Test", Destination: STDERR})
		if err != nil {
			t.Errorf("Update()  error = %v, wantErr %v", err, nil)
			return
		}

		if !reflect.DeepEqual(got, u) {
			t.Errorf("Update() = %v, want %v", u, got)
		}

		u, err = Update(Config{Name: "Test", Destination: DISCARD})
		if err != nil {
			t.Errorf("Update()  error = %v, wantErr %v", err, nil)
			return
		}

		if !reflect.DeepEqual(got, u) {
			t.Errorf("Update() = %v, want %v", u, got)
		}
	})
}

func TestGetConfig(t *testing.T) {
	t.Run("Get Default Config", func(t *testing.T) {
		if got := GetConfig("NOT-AVAILABLE"); !reflect.DeepEqual(got, Config{}) {
			t.Errorf("GetConfig() = %v, want %v", got, Config{})
		}
	})

	t.Run("Get Config", func(t *testing.T) {
		Update(Config{Name: "Test", Destination: DISCARD, FileName: "Test"})
		if got := GetConfig("Test"); !reflect.DeepEqual(got, Config{Name: "Test", Destination: DISCARD, FileName: "Test"}) {
			t.Errorf("GetConfig() = %v, want %v", got, Config{Name: "Test", Destination: DISCARD, FileName: "Test"})
		}

		Update(Config{Destination: DISCARD, FileName: "Test1"})
		if got := GetConfig(""); !reflect.DeepEqual(got, Config{Name: "logger.test", Destination: DISCARD, FileName: "Test1"}) {
			t.Errorf("GetConfig() = %v, want %v", got, Config{Name: "logger.test", Destination: DISCARD, FileName: "Test1"})
		}
	})
}

func TestGetViaName(t *testing.T) {
	l, _ := Update(Config{Name: "Test", Destination: DISCARD, FileName: "Test"})
	t.Run("Get Logger Instance", func(t *testing.T) {
		if got := GetViaName("Test"); !reflect.DeepEqual(got, l) {
			t.Errorf("GetViaName() = %v, want %v", got, l)
		}
	})

	l, _ = Update(Config{Destination: DISCARD, FileName: "Test"})
	t.Run("Get Logger Instance", func(t *testing.T) {
		if got := GetViaName(""); !reflect.DeepEqual(got, l) {
			t.Errorf("GetViaName() = %v, want %v", got, l)
		}
	})

	t.Run("Get Logger Instance", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in Get Logger Instance", r)
			}
		}()
		if got := GetViaName("Test-2"); got == nil {
			t.Errorf("GetViaName() = %v, want %v", nil, "Instance")
		}
	})

	t.Run("Get Process Name Logger Instance", func(t *testing.T) {
		delete(nameToLogger, util.ProcessName())
		name := util.ProcessName()
		if got := GetViaName(name); got == nil {
			t.Errorf("GetViaName() = %v, want %v", nil, "Instance")
		}
	})
}

func TestGet(t *testing.T) {
	l, _ := Update(Config{Destination: DISCARD})
	t.Run("Get Logger Instance", func(t *testing.T) {
		if got := Get(); !reflect.DeepEqual(got, l) {
			t.Errorf("GetViaName() = %v, want %v", got, l)
		}
	})
}

func Test_loggerImpl_formatHeader(t *testing.T) {
	t.Run("formatHeader", func(t *testing.T) {
		date, _ := time.Parse("2006/01/02 15:04:05.999999999", "2019/04/22 11:36:11.109121")
		want := []byte("2019/04/22 11:36:11.109121 hostName logger.test transactionID logger.go")
		l := &loggerImpl{
			writer:   os.Stderr,
			config:   &Config{},
			hostName: "hostName",
		}
		if got := l.formatHeader(date, 1, "transactionID", INFO); !strings.Contains(string(got), string(want)) {
			t.Errorf("loggerImpl.formatHeader() = %v, want %v", string(got), string(want))
		}
	})
}

func TestDiscardLogger(t *testing.T) {
	tests := []struct {
		name string
		want Log
	}{
		{name: "Discard", want: GetViaName(discardLoggerName)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiscardLogger(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiscardLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggerImpl_LogMessages(t *testing.T) {
	validate := func(got string, want []string, dontWant []string, t *testing.T) {
		success := true
		for _, w := range want {
			if !strings.Contains(got, w) {
				success = false
			}
		}

		if !success {
			t.Errorf("loggerImpl_LogLevel() = %v, want %v", got, want)
			return
		}

		success = false
		for _, w := range dontWant {
			if strings.Contains(got, w) {
				success = true
			}
		}

		if success {
			t.Errorf("loggerImpl_LogLevel() = %v, dontWant %v", got, dontWant)
		}
	}

	t.Run("Log level TRACE", func(t *testing.T) {
		b := &bytes.Buffer{}
		l := &loggerImpl{writer: nopCloser{b}, config: &Config{LogLevel: TRACE}, hostName: "hostName"}
		l.Trace("transactionID", "Test Message with %v", "test arg")
		l.Debug("transactionID", "Test Message with %v", "test arg")
		l.Info("transactionID", "Test Message with %v", "test arg")
		l.Warn("transactionID", "Test Message with %v", "test arg")
		l.Error("transactionID", "Error Code", "Test Message with %v", "test arg")
		l.Fatal("transactionID", "Fatal Code", "Test Message with %v", "test arg")

		want := []string{"TRACE  Test Message with test arg", "DEBUG  Test Message with test arg",
			"INFO   Test Message with test arg", "WARN   Test Message with test arg",
			"ERROR  Error Code Test Message with test arg", "FATAL  Fatal Code Test Message with test arg",
		}
		validate(b.String(), want, []string{}, t)
	})

	t.Run("Log level DEBUG", func(t *testing.T) {
		b := &bytes.Buffer{}
		l := &loggerImpl{writer: nopCloser{b}, config: &Config{LogLevel: DEBUG}, hostName: "hostName"}
		l.Trace("transactionID", "Test Message with %v", "test arg")
		l.Debug("transactionID", "Test Message with %v", "test arg")
		l.Info("transactionID", "Test Message with %v", "test arg")
		l.Warn("transactionID", "Test Message with %v", "test arg")
		l.Error("transactionID", "Error Code", "Test Message with %v", "test arg")
		l.Fatal("transactionID", "Fatal Code", "Test Message with %v", "test arg")

		want := []string{"DEBUG  Test Message with test arg", "INFO   Test Message with test arg",
			"WARN   Test Message with test arg", "ERROR  Error Code Test Message with test arg",
			"FATAL  Fatal Code Test Message with test arg",
		}
		dontWant := []string{"TRACE  Test Message with test arg"}
		validate(b.String(), want, dontWant, t)
	})

	t.Run("Log level INFO", func(t *testing.T) {
		b := &bytes.Buffer{}
		l := &loggerImpl{writer: nopCloser{b}, config: &Config{LogLevel: INFO}, hostName: "hostName"}
		l.Trace("transactionID", "Test Message with %v", "test arg")
		l.Debug("transactionID", "Test Message with %v", "test arg")
		l.Info("transactionID", "Test Message with %v", "test arg")
		l.Warn("transactionID", "Test Message with %v", "test arg")
		l.Error("transactionID", "Error Code", "Test Message with %v", "test arg")
		l.Fatal("transactionID", "Fatal Code", "Test Message with %v", "test arg")

		want := []string{"INFO   Test Message with test arg",
			"WARN   Test Message with test arg", "ERROR  Error Code Test Message with test arg",
			"FATAL  Fatal Code Test Message with test arg",
		}
		dontWant := []string{"TRACE  Test Message with test arg", "DEBUG  Test Message with test arg"}
		validate(b.String(), want, dontWant, t)
	})

	t.Run("Log level WARN", func(t *testing.T) {
		b := &bytes.Buffer{}
		l := &loggerImpl{writer: nopCloser{b}, config: &Config{LogLevel: WARN}, hostName: "hostName"}
		l.Trace("transactionID", "Test Message with %v", "test arg")
		l.Debug("transactionID", "Test Message with %v", "test arg")
		l.Info("transactionID", "Test Message with %v", "test arg")
		l.Warn("transactionID", "Test Message with %v", "test arg")
		l.Error("transactionID", "Error Code", "Test Message with %v", "test arg")
		l.Fatal("transactionID", "Fatal Code", "Test Message with %v", "test arg")

		want := []string{"WARN   Test Message with test arg", "ERROR  Error Code Test Message with test arg",
			"FATAL  Fatal Code Test Message with test arg",
		}
		dontWant := []string{"TRACE  Test Message with test arg", "DEBUG  Test Message with test arg",
			"INFO   Test Message with test arg"}
		validate(b.String(), want, dontWant, t)
	})

	t.Run("Log level Error", func(t *testing.T) {
		b := &bytes.Buffer{}
		l := &loggerImpl{writer: nopCloser{b}, config: &Config{LogLevel: ERROR}, hostName: "hostName"}
		l.Trace("transactionID", "Test Message with %v", "test arg")
		l.Debug("transactionID", "Test Message with %v", "test arg")
		l.Info("transactionID", "Test Message with %v", "test arg")
		l.Warn("transactionID", "Test Message with %v", "test arg")
		l.Error("transactionID", "Error Code", "Test Message with %v", "test arg")
		l.Fatal("transactionID", "Fatal Code", "Test Message with %v", "test arg")

		want := []string{"ERROR  Error Code Test Message with test arg", "FATAL  Fatal Code Test Message with test arg"}
		dontWant := []string{"TRACE  Test Message with test arg", "DEBUG  Test Message with test arg",
			"INFO   Test Message with test arg", "WARN   Test Message with test arg"}
		validate(b.String(), want, dontWant, t)
	})

	t.Run("Log level FATAL", func(t *testing.T) {
		b := &bytes.Buffer{}
		l := &loggerImpl{writer: nopCloser{b}, config: &Config{LogLevel: FATAL}, hostName: "hostName"}
		l.Trace("transactionID", "Test Message with %v", "test arg")
		l.Debug("transactionID", "Test Message with %v", "test arg")
		l.Info("transactionID", "Test Message with %v", "test arg")
		l.Warn("transactionID", "Test Message with %v", "test arg")
		l.Error("transactionID", "Error Code", "Test Message with %v", "test arg")
		l.Fatal("transactionID", "Fatal Code", "Test Message with %v", "test arg")

		want := []string{"FATAL  Fatal Code Test Message with test arg"}
		dontWant := []string{"TRACE  Test Message with test arg", "DEBUG  Test Message with test arg",
			"INFO   Test Message with test arg", "WARN   Test Message with test arg", "ERROR  Error Code Test Message with test arg"}
		validate(b.String(), want, dontWant, t)
	})
}

func Test_loggerImpl_SetWriter(t *testing.T) {
	type fields struct {
		writer   io.WriteCloser
		config   *Config
		hostName string
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
	}{
		{name: "default"},
		{name: "buffer writer", fields: fields{writer: nopCloser{&bytes.Buffer{}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &loggerImpl{
				writer:   tt.fields.writer,
				config:   tt.fields.config,
				hostName: tt.fields.hostName,
			}
			writer := &bytes.Buffer{}
			l.SetWriter(writer)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("loggerImpl.SetWriter() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func Test_loggerImpl_LogLevel(t *testing.T) {
	type fields struct {
		writer   io.WriteCloser
		config   *Config
		hostName string
	}
	tests := []struct {
		name   string
		fields fields
		want   LogLevel
	}{
		{name: "Default", want: INFO, fields: fields{config: &Config{}}},
		{name: "TRACE", want: TRACE, fields: fields{config: &Config{LogLevel: TRACE}}},
		{name: "DEBUG", want: DEBUG, fields: fields{config: &Config{LogLevel: DEBUG}}},
		{name: "INFO", want: INFO, fields: fields{config: &Config{LogLevel: INFO}}},
		{name: "WARN", want: WARN, fields: fields{config: &Config{LogLevel: WARN}}},
		{name: "ERROR", want: ERROR, fields: fields{config: &Config{LogLevel: ERROR}}},
		{name: "FATAL", want: FATAL, fields: fields{config: &Config{LogLevel: FATAL}}},
		{name: "OFF", want: OFF, fields: fields{config: &Config{LogLevel: OFF}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &loggerImpl{
				writer:   tt.fields.writer,
				config:   tt.fields.config,
				hostName: tt.fields.hostName,
			}
			if got := l.LogLevel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loggerImpl.LogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggerImpl_GetWriter(t *testing.T) {
	type fields struct {
		writer   io.WriteCloser
		config   *Config
		hostName string
	}
	tests := []struct {
		name   string
		fields fields
		want   io.WriteCloser
	}{
		{name: "buffer writer", want: nopCloser{&bytes.Buffer{}}, fields: fields{writer: nopCloser{&bytes.Buffer{}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &loggerImpl{
				writer:   tt.fields.writer,
				config:   tt.fields.config,
				hostName: tt.fields.hostName,
			}
			if got := l.GetWriter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loggerImpl.GetWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}
