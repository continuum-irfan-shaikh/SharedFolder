package consumer

import (
	"testing"
)

func Test_workerPool_initialize(t *testing.T) {
	type fields struct {
		cfg Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "Pool size zero", wantErr: true, fields: fields{cfg: Config{SubscriberPerCore: 0}}},
		{name: "Pool size blank", wantErr: true, fields: fields{cfg: Config{}}},
		{name: "Pool size Negative", wantErr: true, fields: fields{cfg: Config{SubscriberPerCore: -10}}},
		{name: "Pool size Positive", wantErr: false, fields: fields{cfg: Config{SubscriberPerCore: 10}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &workerPool{
				cfg: tt.fields.cfg,
			}
			if err := w.initialize(); (err != nil) != tt.wantErr {
				t.Errorf("workerPool.initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
