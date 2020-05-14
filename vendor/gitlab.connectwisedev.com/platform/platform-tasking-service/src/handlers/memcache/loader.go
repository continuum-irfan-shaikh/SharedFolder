package memcache

import (
	"context"

	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/logger"
	"gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/common/errorcode"
)

//go:generate mockgen -destination=../../mocks/mocks-repository/trigger_mock.go -package=mockrepositories -source=./loader.go

//TriggerRepo - represents repository responsible for loading data about active triggers into cache
type TriggerRepo interface {
	LoadTriggersToCache() error
}

//Loader - represents loader handler
type Loader struct {
	triggerRepo TriggerRepo
	log         logger.Logger
}

//NewLoader - returns instance of loader handler
func NewLoader(triggerRepo TriggerRepo, log logger.Logger) *Loader {
	return &Loader{triggerRepo: triggerRepo, log: log}
}

//LoadTriggersToCache - runs process of loading trigger data into cache
func (l *Loader) LoadTriggersToCache(ctx context.Context) {
	if err := l.triggerRepo.LoadTriggersToCache(); err != nil {
		l.log.ErrfCtx(ctx,errorcode.ErrorCantProcessData, "can't load active triggers to cache: %s", err)
	}
}
