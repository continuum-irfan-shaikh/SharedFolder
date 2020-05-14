package memcache

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mockLoggerTasking"
	mockrepositories "gitlab.connectwisedev.com/platform/platform-tasking-service/src/mocks/mocks-repository"
)

func TestLoader_LoadTriggersToCache(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name   string
		fields func() (TriggerRepo, logger.Logger)
	}{
		{
			name: "success",
			fields: func() (TriggerRepo, logger.Logger) {
				r := mockrepositories.NewMockTriggerRepo(ctrl)
				r.EXPECT().LoadTriggersToCache().Return(nil)
				return r, nil
			},
		},
		{
			name: "err",
			fields: func() (TriggerRepo, logger.Logger) {
				r, l := mockrepositories.NewMockTriggerRepo(ctrl), mockLoggerTasking.NewMockLogger(ctrl)
				r.EXPECT().LoadTriggersToCache().Return(errors.New("fail"))
				return r, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := tt.fields()
			loader := NewLoader(r, logger.Log)
			loader.LoadTriggersToCache(context.Background())
		})
	}
}

func TestNewLoader(t *testing.T) {
	r, l := mockrepositories.NewMockTriggerRepo(nil), mockLoggerTasking.NewMockLogger(nil)
	expected := &Loader{r, logger.Log}
	actual := NewLoader(r, logger.Log)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("NewLoader() = %v, want %v", actual, expected)
	}
}
