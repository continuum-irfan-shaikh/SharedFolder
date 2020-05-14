package main

import (
	"errors"
	"testing"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/wal"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/wal/mock"

	"github.com/golang/mock/gomock"
)

func Test_process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockWAL(ctrl)

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Write Error", wantErr: true, setup: func() {
				wal.Create = func(*wal.Config) wal.WAL { return m }
				m.EXPECT().Write("Test-1").Return(errors.New("Error"))
			},
		},
		{
			name: "Write Object Error", wantErr: true, setup: func() {
				wal.Create = func(*wal.Config) wal.WAL { return m }
				m.EXPECT().Write("Test-1").Return(nil)
				m.EXPECT().WriteObject("Test-2").Return(errors.New("Error"))
			},
		},
		{
			name: "Read Error", wantErr: true, setup: func() {
				wal.Create = func(*wal.Config) wal.WAL { return m }
				m.EXPECT().Write("Test-1").Return(nil)
				m.EXPECT().WriteObject("Test-2").Return(nil)
				m.EXPECT().Flush()
				m.EXPECT().Read(5).Return(nil, errors.New("Error"))
			},
		},
		{
			name: "Success", wantErr: false, setup: func() {
				wal.Create = func(*wal.Config) wal.WAL { return m }
				m.EXPECT().Write("Test-1").Return(nil)
				m.EXPECT().WriteObject("Test-2").Return(nil)
				m.EXPECT().Flush()
				m.EXPECT().Read(5).Return([]wal.Record{wal.Record{}}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if err := process(); (err != nil) != tt.wantErr {
				t.Errorf("process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
