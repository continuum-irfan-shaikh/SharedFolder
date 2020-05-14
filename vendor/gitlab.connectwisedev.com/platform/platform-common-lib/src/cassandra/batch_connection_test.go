// +build !integration

package cassandra

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/mock/gomock"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/cassandra/cql"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/cassandra/cql/mock"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/circuit"
)

func TestNewBatchDbConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	session := mock.NewMockSession(ctrl)

	type args struct {
		conf *DbConfig
	}
	tests := []struct {
		name    string
		args    args
		want    BatchDbConnector
		setup   func()
		wantErr bool
	}{
		{
			name: "Nil Config", wantErr: true, want: nil, args: args{conf: nil}, setup: func() {
				cql.CreateSession = func(cluster *gocql.ClusterConfig) (cql.Session, error) { return session, errors.New("Error") }
			},
		},
		{
			name: "Blank Config", wantErr: true, want: nil, args: args{conf: &DbConfig{}}, setup: func() {
				cql.CreateSession = func(cluster *gocql.ClusterConfig) (cql.Session, error) { return session, errors.New("Error") }
			},
		},
		{
			name: "Blank Host", wantErr: true, want: nil, args: args{conf: &DbConfig{Hosts: []string{}}}, setup: func() {
				cql.CreateSession = func(cluster *gocql.ClusterConfig) (cql.Session, error) { return session, errors.New("Error") }
			},
		},
		{
			name: "Blank Keyspace", wantErr: true, want: nil, args: args{conf: &DbConfig{Hosts: []string{"Sever"}, Keyspace: ""}}, setup: func() {
				cql.CreateSession = func(cluster *gocql.ClusterConfig) (cql.Session, error) { return session, errors.New("Error") }
			},
		},
		{
			name: "Session Creation Error", wantErr: true, want: nil,
			args: args{conf: &DbConfig{Hosts: []string{"Sever"}, Keyspace: "test"}}, setup: func() {
				cql.CreateSession = func(cluster *gocql.ClusterConfig) (cql.Session, error) { return session, errors.New("Error") }
			},
		},
		{
			name: "Session Creation Success", wantErr: false,
			want: &batchConnection{Connection: connection{conf: &DbConfig{Hosts: []string{"Server"}, Keyspace: "test",
				NumConn: 10, TimeoutMillisecond: time.Millisecond * 10, CircuitBreaker: &circuit.Config{
					Enabled: true, TimeoutInSecond: 1, MaxConcurrentRequests: 15000,
					ErrorPercentThreshold: 25, RequestVolumeThreshold: 500, SleepWindowInSecond: 10,
				}, CommandName: "Database-Command1", ValidErrors: []string{"Error1"}}, session: session},
				batch: &gocql.Batch{}},
			args: args{conf: &DbConfig{Hosts: []string{"Server"}, Keyspace: "test",
				NumConn: 10, TimeoutMillisecond: time.Millisecond * 10, CircuitBreaker: &circuit.Config{
					Enabled: true, TimeoutInSecond: 1, MaxConcurrentRequests: 15000,
					ErrorPercentThreshold: 25, RequestVolumeThreshold: 500, SleepWindowInSecond: 10,
				}, CommandName: "Database-Command1", ValidErrors: []string{"Error1"}}},
			setup: func() {
				session.EXPECT().NewBatch(gomock.Any()).Return(&gocql.Batch{})
				cql.CreateSession = func(cluster *gocql.ClusterConfig) (cql.Session, error) { return session, nil }
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := NewBatchDbConnection(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBatchDbConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBatchDbConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_batchConnection_BatchExecution(t *testing.T) {
	ctrl := gomock.NewController(t)
	session := mock.NewMockSession(ctrl)
	config := DbConfig{CommandName: "DB", CircuitBreaker: circuit.New()}

	type fields struct {
		Connection connection
		batch      *gocql.Batch
	}
	type args struct {
		query  string
		values [][]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func()
		wantErr bool
	}{
		{
			name: "Nil Session", wantErr: true,
			fields: fields{Connection: connection{conf: &config, session: nil}, batch: &gocql.Batch{}},
			setup:  func() {},
		},
		{
			name: "ExecuteBatch Error", wantErr: true,
			fields: fields{Connection: connection{conf: &config, session: session}, batch: &gocql.Batch{}},
			setup: func() {
				session.EXPECT().ExecuteBatch(gomock.Any()).Return(errors.New("Error"))
			},
		},
		{
			name: "ExecuteBatch Success", wantErr: false, args: args{values: [][]interface{}{[]interface{}{}}},
			fields: fields{Connection: connection{conf: &config, session: session}, batch: &gocql.Batch{}},
			setup: func() {
				session.EXPECT().ExecuteBatch(gomock.Any()).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			d := batchConnection{
				Connection: tt.fields.Connection,
				batch:      tt.fields.batch,
			}
			if err := d.BatchExecution(tt.args.query, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("batchConnection.BatchExecution() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_batchConnection_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	session := mock.NewMockSession(ctrl)
	type fields struct {
		Connection connection
		batch      *gocql.Batch
	}
	tests := []struct {
		name   string
		fields fields
		setup  func()
	}{
		{
			name:   "Nil Session",
			fields: fields{Connection: connection{conf: &DbConfig{}, session: nil}},
			setup:  func() {},
		},
		{
			name:   "Session",
			fields: fields{Connection: connection{conf: &DbConfig{}, session: session}},
			setup:  func() { session.EXPECT().Close() },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			d := batchConnection{
				Connection: tt.fields.Connection,
				batch:      tt.fields.batch,
			}
			d.Close()
		})
	}
}

func Test_batchConnection_Closed(t *testing.T) {
	ctrl := gomock.NewController(t)
	session := mock.NewMockSession(ctrl)
	type fields struct {
		Connection connection
		batch      *gocql.Batch
	}
	tests := []struct {
		name   string
		fields fields
		setup  func()
		want   bool
	}{
		{
			name: "Nil Session", want: true,
			fields: fields{Connection: connection{conf: &DbConfig{}, session: nil}},
			setup:  func() {},
		},
		{
			name: "Session Closed", want: true,
			fields: fields{Connection: connection{conf: &DbConfig{}, session: session}},
			setup:  func() { session.EXPECT().Closed().Return(true) },
		},
		{
			name: "Session Open", want: false,
			fields: fields{Connection: connection{conf: &DbConfig{}, session: session}},
			setup:  func() { session.EXPECT().Closed().Return(false) },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			d := batchConnection{
				Connection: tt.fields.Connection,
				batch:      tt.fields.batch,
			}
			if got := d.Closed(); got != tt.want {
				t.Errorf("connection.Closed() = %v, want %v", got, tt.want)
			}
		})
	}
}
