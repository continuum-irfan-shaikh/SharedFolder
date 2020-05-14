package logger

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	rLog "gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/config"
	mock "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mock-logger"
)

func TestLoadLogger(t *testing.T) {
	config.Config.Log.FileName = "logs_test.log"
	if _, err := os.Create(config.Config.Log.FileName); err != nil {
		t.Fatalf("can't create file: %s", config.Config.Log.FileName)
	}
	defer os.Remove(config.Config.Log.FileName)

	config.Config.Log.LogLevel = rLog.DEBUG
	err := Load(config.Config.Log)
	if err != nil {
		t.Fatalf("Load method return an error:%v", err)
	}

	if &CassandraLogger == nil {
		t.Fatal(`CassandraLogger should be initialized`)
	}
	if &ZKLogger == nil {
		t.Fatal("ZKLogger should be initialized")
	}

	CassandraLogger.Print("CassandraLogger")
	CassandraLogger.Printf("CassandraLogger - f")
	CassandraLogger.Println("CassandraLogger - ln")
	ZKLogger.LogInfo("Zookeeper info")
	ZKLogger.LogError("Zookeeper error")

	time.Sleep(1 * time.Second)

	b, _ := ioutil.ReadFile(config.Config.Log.FileName)
	strs := []string{
		"CassandraLogger",
		"CassandraLogger - f",
		"CassandraLogger - ln",
		"Zookeeper info",
		"Zookeeper error",
	}

	for _, str := range strs {
		if !strings.Contains(string(b), str) {
			t.Fatalf("config.Config.LogFile does not contain substr \"%s\"", str)
		}
	}

	err = Load(config.Config.Log)
	if err != nil {
		t.Fatalf("Second call of the Load method return an error:%v", err)
	}
}

func TestLogWrapper_InfofCtx(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	logMock := mock.NewMockLog(ctrl)
	format, err := "error: %v", "somme error"
	log := logWrapper{logMock}
	ctx := context.Background()
	transactionID := log.transactionID(ctx)

	logMock.EXPECT().Info(transactionID, format, []interface{}{err}).Times(1)
	log.InfofCtx(ctx, format, err)
}

func TestLogWrapper_WarnfCtx(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	logMock := mock.NewMockLog(ctrl)
	format, err := "error: %v", "somme error"
	log := logWrapper{logMock}
	ctx := context.Background()
	transactionID := log.transactionID(ctx)

	logMock.EXPECT().Warn(transactionID, format, []interface{}{err}).Times(1)
	log.WarnfCtx(ctx, format, err)
}

func TestLogWrapper_ErrfCtx(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	logMock := mock.NewMockLog(ctrl)
	format, err := "error: %v", "somme error"
	log := logWrapper{logMock}
	ctx := context.Background()
	transactionID := log.transactionID(ctx)

	logMock.EXPECT().Error(transactionID, "", format, []interface{}{err}).Times(1)
	log.ErrfCtx(ctx, "", format, err)
}

func TestLogWrapper_Debug(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	logMock := mock.NewMockLog(ctrl)
	err := "error: %v"
	log := logWrapper{logMock}

	logMock.EXPECT().Debug("", err).Times(1)
	log.Debug(err)
}

func TestLogWrapper_Debugln(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)
	logMock := mock.NewMockLog(ctrl)
	err := "error: %v"
	log := logWrapper{logMock}

	logMock.EXPECT().Debug("", fmt.Sprintln(err)).Times(1)
	log.Debugln(err)
}
