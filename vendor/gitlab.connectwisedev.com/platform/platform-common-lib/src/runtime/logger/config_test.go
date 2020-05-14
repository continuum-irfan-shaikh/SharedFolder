package logger

import (
	"path/filepath"
	"reflect"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/util"
)

func TestConfig_name(t *testing.T) {
	type fields struct {
		Name        string
		FileName    string
		MaxSize     int
		MaxAge      int
		MaxBackups  int
		ServiceName string
		Filler      string
		LogLevel    LogLevel
		Destination Destination
		CallDepth   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Default-Name", fields: fields{}, want: "logger.test"},
		{name: "Given Name", fields: fields{Name: "Test"}, want: "Test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Name:        tt.fields.Name,
				FileName:    tt.fields.FileName,
				MaxSize:     tt.fields.MaxSize,
				MaxAge:      tt.fields.MaxAge,
				MaxBackups:  tt.fields.MaxBackups,
				ServiceName: tt.fields.ServiceName,
				Filler:      tt.fields.Filler,
				LogLevel:    tt.fields.LogLevel,
				Destination: tt.fields.Destination,
				CallDepth:   tt.fields.CallDepth,
			}
			if got := c.name(); got != tt.want {
				t.Errorf("Config.name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_fileName(t *testing.T) {
	type fields struct {
		Name        string
		FileName    string
		MaxSize     int
		MaxAge      int
		MaxBackups  int
		ServiceName string
		Filler      string
		LogLevel    LogLevel
		Destination Destination
		CallDepth   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Default-Name", fields: fields{}, want: filepath.Join(util.InvocationPath(), "logger.test.log")},
		{name: "Default-Name", fields: fields{Name: "test"}, want: filepath.Join(util.InvocationPath(), "logger.test-test.log")},
		{name: "Given Name", fields: fields{FileName: "Test"}, want: "Test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Name:        tt.fields.Name,
				FileName:    tt.fields.FileName,
				MaxSize:     tt.fields.MaxSize,
				MaxAge:      tt.fields.MaxAge,
				MaxBackups:  tt.fields.MaxBackups,
				ServiceName: tt.fields.ServiceName,
				Filler:      tt.fields.Filler,
				LogLevel:    tt.fields.LogLevel,
				Destination: tt.fields.Destination,
				CallDepth:   tt.fields.CallDepth,
			}
			if got := c.fileName(); got != tt.want {
				t.Errorf("Config.fileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_maxSize(t *testing.T) {
	type fields struct {
		Name        string
		FileName    string
		MaxSize     int
		MaxAge      int
		MaxBackups  int
		ServiceName string
		Filler      string
		LogLevel    LogLevel
		Destination Destination
		CallDepth   int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{name: "Default-Value", fields: fields{}, want: 20},
		{name: "Default-Value-Nagative", fields: fields{MaxSize: -2}, want: 20},
		{name: "Given Value", fields: fields{FileName: "Test", MaxSize: 25}, want: 25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Name:        tt.fields.Name,
				FileName:    tt.fields.FileName,
				MaxSize:     tt.fields.MaxSize,
				MaxAge:      tt.fields.MaxAge,
				MaxBackups:  tt.fields.MaxBackups,
				ServiceName: tt.fields.ServiceName,
				Filler:      tt.fields.Filler,
				LogLevel:    tt.fields.LogLevel,
				Destination: tt.fields.Destination,
				CallDepth:   tt.fields.CallDepth,
			}
			if got := c.maxSize(); got != tt.want {
				t.Errorf("Config.maxSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_maxAge(t *testing.T) {
	type fields struct {
		Name        string
		FileName    string
		MaxSize     int
		MaxAge      int
		MaxBackups  int
		ServiceName string
		Filler      string
		LogLevel    LogLevel
		Destination Destination
		CallDepth   int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{name: "Default-Value", fields: fields{}, want: 30},
		{name: "Default-Value-Nagative", fields: fields{MaxAge: -2}, want: 30},
		{name: "Given Value", fields: fields{FileName: "Test", MaxAge: 25}, want: 25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Name:        tt.fields.Name,
				FileName:    tt.fields.FileName,
				MaxSize:     tt.fields.MaxSize,
				MaxAge:      tt.fields.MaxAge,
				MaxBackups:  tt.fields.MaxBackups,
				ServiceName: tt.fields.ServiceName,
				Filler:      tt.fields.Filler,
				LogLevel:    tt.fields.LogLevel,
				Destination: tt.fields.Destination,
				CallDepth:   tt.fields.CallDepth,
			}
			if got := c.maxAge(); got != tt.want {
				t.Errorf("Config.maxAge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_maxBackups(t *testing.T) {
	type fields struct {
		Name        string
		FileName    string
		MaxSize     int
		MaxAge      int
		MaxBackups  int
		ServiceName string
		Filler      string
		LogLevel    LogLevel
		Destination Destination
		CallDepth   int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{name: "Default-Value", fields: fields{}, want: 5},
		{name: "Default-Value-Nagative", fields: fields{MaxBackups: -2}, want: 5},
		{name: "Given Value", fields: fields{FileName: "Test", MaxBackups: 25}, want: 25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Name:        tt.fields.Name,
				FileName:    tt.fields.FileName,
				MaxSize:     tt.fields.MaxSize,
				MaxAge:      tt.fields.MaxAge,
				MaxBackups:  tt.fields.MaxBackups,
				ServiceName: tt.fields.ServiceName,
				Filler:      tt.fields.Filler,
				LogLevel:    tt.fields.LogLevel,
				Destination: tt.fields.Destination,
				CallDepth:   tt.fields.CallDepth,
			}
			if got := c.maxBackups(); got != tt.want {
				t.Errorf("Config.maxBackups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_serviceName(t *testing.T) {
	type fields struct {
		Name        string
		FileName    string
		MaxSize     int
		MaxAge      int
		MaxBackups  int
		ServiceName string
		Filler      string
		LogLevel    LogLevel
		Destination Destination
		CallDepth   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Default-Name", fields: fields{}, want: "logger.test"},
		{name: "Default-Name", fields: fields{Name: "test"}, want: "logger.test"},
		{name: "Given Name", fields: fields{ServiceName: "Test"}, want: "Test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Name:        tt.fields.Name,
				FileName:    tt.fields.FileName,
				MaxSize:     tt.fields.MaxSize,
				MaxAge:      tt.fields.MaxAge,
				MaxBackups:  tt.fields.MaxBackups,
				ServiceName: tt.fields.ServiceName,
				Filler:      tt.fields.Filler,
				LogLevel:    tt.fields.LogLevel,
				Destination: tt.fields.Destination,
				CallDepth:   tt.fields.CallDepth,
			}
			if got := c.serviceName(); got != tt.want {
				t.Errorf("Config.serviceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_filler(t *testing.T) {
	type fields struct {
		Name        string
		FileName    string
		MaxSize     int
		MaxAge      int
		MaxBackups  int
		ServiceName string
		Filler      string
		LogLevel    LogLevel
		Destination Destination
		CallDepth   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Default-Name", fields: fields{}, want: "-"},
		{name: "Given Name", fields: fields{Filler: "Test"}, want: "Test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Name:        tt.fields.Name,
				FileName:    tt.fields.FileName,
				MaxSize:     tt.fields.MaxSize,
				MaxAge:      tt.fields.MaxAge,
				MaxBackups:  tt.fields.MaxBackups,
				ServiceName: tt.fields.ServiceName,
				Filler:      tt.fields.Filler,
				LogLevel:    tt.fields.LogLevel,
				Destination: tt.fields.Destination,
				CallDepth:   tt.fields.CallDepth,
			}
			if got := c.filler(); got != tt.want {
				t.Errorf("Config.filler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_logLevel(t *testing.T) {
	type fields struct {
		Name        string
		FileName    string
		MaxSize     int
		MaxAge      int
		MaxBackups  int
		ServiceName string
		Filler      string
		LogLevel    LogLevel
		Destination Destination
		CallDepth   int
	}
	tests := []struct {
		name   string
		fields fields
		want   LogLevel
	}{
		{name: "Default-Name", fields: fields{}, want: INFO},
		{name: "Given Name - TRACE", fields: fields{LogLevel: TRACE}, want: TRACE},
		{name: "Given Name - DEBUG", fields: fields{LogLevel: DEBUG}, want: DEBUG},
		{name: "Given Name - INFO", fields: fields{LogLevel: INFO}, want: INFO},
		{name: "Given Name - WARN", fields: fields{LogLevel: WARN}, want: WARN},
		{name: "Given Name - ERROR", fields: fields{LogLevel: ERROR}, want: ERROR},
		{name: "Given Name - FATAL", fields: fields{LogLevel: FATAL}, want: FATAL},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Name:        tt.fields.Name,
				FileName:    tt.fields.FileName,
				MaxSize:     tt.fields.MaxSize,
				MaxAge:      tt.fields.MaxAge,
				MaxBackups:  tt.fields.MaxBackups,
				ServiceName: tt.fields.ServiceName,
				Filler:      tt.fields.Filler,
				LogLevel:    tt.fields.LogLevel,
				Destination: tt.fields.Destination,
				CallDepth:   tt.fields.CallDepth,
			}
			if got := c.logLevel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.logLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_destination(t *testing.T) {
	type fields struct {
		Name        string
		FileName    string
		MaxSize     int
		MaxAge      int
		MaxBackups  int
		ServiceName string
		Filler      string
		LogLevel    LogLevel
		Destination Destination
		CallDepth   int
	}
	tests := []struct {
		name   string
		fields fields
		want   Destination
	}{
		{name: "Default-Value", fields: fields{Destination: Destination{order: 5, value: "Unknown"}}, want: FILE},
		{name: "Default-Value-FILE", fields: fields{Destination: FILE}, want: FILE},
		{name: "Given-Value-STDERR", fields: fields{Destination: STDERR}, want: STDERR},
		{name: "Given-Value-DISCARD", fields: fields{Destination: DISCARD}, want: DISCARD},
		{name: "Given-Value-STDOUT", fields: fields{Destination: STDOUT}, want: STDOUT},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Name:        tt.fields.Name,
				FileName:    tt.fields.FileName,
				MaxSize:     tt.fields.MaxSize,
				MaxAge:      tt.fields.MaxAge,
				MaxBackups:  tt.fields.MaxBackups,
				ServiceName: tt.fields.ServiceName,
				Filler:      tt.fields.Filler,
				LogLevel:    tt.fields.LogLevel,
				Destination: tt.fields.Destination,
				CallDepth:   tt.fields.CallDepth,
			}
			if got := c.destination(); got != tt.want {
				t.Errorf("Config.destination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_calldepth(t *testing.T) {
	type fields struct {
		Name        string
		FileName    string
		MaxSize     int
		MaxAge      int
		MaxBackups  int
		ServiceName string
		Filler      string
		LogLevel    LogLevel
		Destination Destination
		CallDepth   int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{name: "Default-Value", fields: fields{}, want: 4},
		{name: "Default-Value-10", fields: fields{CallDepth: 10}, want: 10},
		{name: "Default-Value-20", fields: fields{CallDepth: 20}, want: 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Name:        tt.fields.Name,
				FileName:    tt.fields.FileName,
				MaxSize:     tt.fields.MaxSize,
				MaxAge:      tt.fields.MaxAge,
				MaxBackups:  tt.fields.MaxBackups,
				ServiceName: tt.fields.ServiceName,
				Filler:      tt.fields.Filler,
				LogLevel:    tt.fields.LogLevel,
				Destination: tt.fields.Destination,
				CallDepth:   tt.fields.CallDepth,
			}
			if got := c.calldepth(); got != tt.want {
				t.Errorf("Config.calldepth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogLevel_UnmarshalJSON(t *testing.T) {
	type fields struct {
		order int
		name  string
		value string
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    LogLevel
	}{
		{name: "Default-Name", fields: fields{}, args: args{[]byte("")}, wantErr: false, want: INFO},
		{name: "Given Name - OFF", fields: fields{}, args: args{[]byte("\"OFF\"")}, wantErr: false, want: OFF},
		{name: "Given Name - OFF", fields: fields{}, args: args{[]byte("\"off\"")}, wantErr: false, want: OFF},
		{name: "Given Name - TRACE", fields: fields{}, args: args{[]byte("\"TRACE\"")}, wantErr: false, want: TRACE},
		{name: "Given Name - TRACE", fields: fields{}, args: args{[]byte("\"trace\"")}, wantErr: false, want: TRACE},
		{name: "Given Name - DEBUG", fields: fields{}, args: args{[]byte("\"DEBUG\"")}, wantErr: false, want: DEBUG},
		{name: "Given Name - DEBUG", fields: fields{}, args: args{[]byte("\"debug\"")}, wantErr: false, want: DEBUG},
		{name: "Given Name - INFO", fields: fields{}, args: args{[]byte("\"INFO\"")}, wantErr: false, want: INFO},
		{name: "Given Name - INFO", fields: fields{}, args: args{[]byte("\"info\"")}, wantErr: false, want: INFO},
		{name: "Given Name - WARN", fields: fields{}, args: args{[]byte("\"WARN\"")}, wantErr: false, want: WARN},
		{name: "Given Name - WARN", fields: fields{}, args: args{[]byte("\"warn\"")}, wantErr: false, want: WARN},
		{name: "Given Name - ERROR", fields: fields{}, args: args{[]byte("\"ERROR\"")}, wantErr: false, want: ERROR},
		{name: "Given Name - ERROR", fields: fields{}, args: args{[]byte("\"error\"")}, wantErr: false, want: ERROR},
		{name: "Given Name - FATAL", fields: fields{}, args: args{[]byte("\"FATAL\"")}, wantErr: false, want: FATAL},
		{name: "Given Name - FATAL", fields: fields{}, args: args{[]byte("\"fatal\"")}, wantErr: false, want: FATAL},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LogLevel{
				order: tt.fields.order,
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if err := l.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("LogLevel.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *l != tt.want {
				t.Errorf("Config.calldepth() = %v, want %v", l, tt.want)
			}
		})
	}
}

func TestLogLevel_MarshalJSON(t *testing.T) {
	type fields struct {
		order int
		name  string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{name: "OFF", fields: fields{value: OFF.value}, wantErr: false, want: []byte("\"OFF\"")},
		{name: "TRACE", fields: fields{value: TRACE.value}, wantErr: false, want: []byte("\"TRACE\"")},
		{name: "DEBUG", fields: fields{value: DEBUG.value}, wantErr: false, want: []byte("\"DEBUG\"")},
		{name: "INFO", fields: fields{value: INFO.value}, wantErr: false, want: []byte("\"INFO\"")},
		{name: "WARN", fields: fields{value: WARN.value}, wantErr: false, want: []byte("\"WARN\"")},
		{name: "ERROR", fields: fields{value: ERROR.value}, wantErr: false, want: []byte("\"ERROR\"")},
		{name: "FATAL", fields: fields{value: FATAL.value}, wantErr: false, want: []byte("\"FATAL\"")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := LogLevel{
				order: tt.fields.order,
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			got, err := l.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("LogLevel.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogLevel.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDestination_UnmarshalJSON(t *testing.T) {
	type fields struct {
		order int
		value string
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    Destination
	}{
		{name: "Default-Name", fields: fields{}, args: args{[]byte("")}, wantErr: false, want: FILE},
		{name: "Given Name - FILE", fields: fields{}, args: args{[]byte("\"FILE\"")}, wantErr: false, want: FILE},
		{name: "Given Name - FILE", fields: fields{}, args: args{[]byte("\"file\"")}, wantErr: false, want: FILE},
		{name: "Given Name - STDERR", fields: fields{}, args: args{[]byte("\"STDERR\"")}, wantErr: false, want: STDERR},
		{name: "Given Name - STDERR", fields: fields{}, args: args{[]byte("\"stderr\"")}, wantErr: false, want: STDERR},
		{name: "Given Name - STDOUT", fields: fields{}, args: args{[]byte("\"STDOUT\"")}, wantErr: false, want: STDOUT},
		{name: "Given Name - STDOUT", fields: fields{}, args: args{[]byte("\"stdout\"")}, wantErr: false, want: STDOUT},
		{name: "Given Name - DISCARD", fields: fields{}, args: args{[]byte("\"DISCARD\"")}, wantErr: false, want: DISCARD},
		{name: "Given Name - DISCARD", fields: fields{}, args: args{[]byte("\"discard\"")}, wantErr: false, want: DISCARD},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Destination{
				order: tt.fields.order,
				value: tt.fields.value,
			}
			if err := d.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Destination.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *d != tt.want {
				t.Errorf("Config.calldepth() = %v, want %v", d, tt.want)
			}
		})
	}
}

func TestDestination_MarshalJSON(t *testing.T) {
	type fields struct {
		order int
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{name: "FILE", fields: fields{value: FILE.value}, wantErr: false, want: []byte("\"FILE\"")},
		{name: "STDERR", fields: fields{value: STDERR.value}, wantErr: false, want: []byte("\"STDERR\"")},
		{name: "STDOUT", fields: fields{value: STDOUT.value}, wantErr: false, want: []byte("\"STDOUT\"")},
		{name: "DISCARD", fields: fields{value: DISCARD.value}, wantErr: false, want: []byte("\"DISCARD\"")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Destination{
				order: tt.fields.order,
				value: tt.fields.value,
			}
			got, err := d.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Destination.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Destination.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
