package wal

import (
	"fmt"
	"testing"
	"time"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/util"
)

func TestCreate(t *testing.T) {
	t.Run("Create Instance", func(t *testing.T) {
		got := Create(&Config{})
		_, ok := got.(*walImpl)
		if !ok {
			t.Errorf("Create() = %v, want walImpl", got)
		}
	})
}

func Test_walImpl_WriteRead(t *testing.T) {
	t.Run("Write Read test", func(t *testing.T) {
		w := Create(&Config{Name: fmt.Sprintf("./wal/%s.wal", util.ProcessName()),
			MaxSegments: 100, SegmentSizeInKB: 2, MaxAgeDays: 10, MaxFiles: 20})

		err := w.Write("Test-1", "Test-\n2")
		if err != nil {
			t.Errorf("walImpl.Write() Expected nil but got error = %v", err)
		}

		records, err := w.Read(1)
		size := len(records)
		if size != 0 {
			t.Errorf("walImpl.Read() Expected 0 but got = %v", size)
		}

		w.Flush()
		time.Sleep(time.Millisecond)

		records, err = w.Read(1)
		size = len(records)
		if err != nil || size != 1 {
			t.Errorf("walImpl.Read() Expected 1 but got = %v", size)
		}

		test := "Test-3"
		err = w.WriteObject(test, test)
		if err != nil {
			t.Errorf("walImpl.WriteObject() Expected nil but got error = %v", err)
		}

		w.Flush()
		records, err = w.Read(1)
		size = len(records)
		if err != nil || size != 1 {
			t.Errorf("walImpl.Read() Expected 1 but got = %v", size)
		}

		records, err = w.Read(2)
		size = len(records)
		if err != nil || size != 2 {
			t.Errorf("walImpl.Read() Expected 2 but got = %v", size)
		}

		for _, r := range records {
			if len(r.Messages) != 2 {
				t.Errorf("walImpl.Read() Expected 2 messages but got = %v", len(r.Messages))
			}

			err = r.Commit()
			if err != nil {
				t.Errorf("walImpl.Read() Expected nil but got = %v", err)
			}
		}
	})
}

func TestRecord_Commit(t *testing.T) {
	type fields struct {
		fileName string
		Messages []string
		Err      error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "blank file", fields: fields{fileName: ""}, wantErr: false},
		{name: "dot file", fields: fields{fileName: "."}, wantErr: false},
		{name: "actual file", fields: fields{fileName: "aaa"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				fileName: tt.fields.fileName,
				Messages: tt.fields.Messages,
				Err:      tt.fields.Err,
			}
			if err := r.Commit(); (err != nil) != tt.wantErr {
				t.Errorf("Record.Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
