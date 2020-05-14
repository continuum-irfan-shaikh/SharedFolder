package consumer

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/utils"

	"github.com/lkumarjain/work"
)

type workerPool struct {
	pool *work.Pool
	cfg  Config
}

func (w *workerPool) initialize() error {
	var err error
	numCPUs := runtime.NumCPU()
	Logger().Info(utils.GetTransactionID(), "Worker pool size : %d", numCPUs*w.cfg.SubscriberPerCore)
	w.pool, err = work.New(numCPUs*w.cfg.SubscriberPerCore, time.Second, nil)
	return err
}

func (w *workerPool) addJob(message Message, commitStrategy commitStrategy) {
	w.pool.Run(job{
		cfg:            w.cfg,
		message:        message,
		commitStrategy: commitStrategy,
	})
}

func (w *workerPool) Shutdown() {
	w.pool.Shutdown()
}

type job struct {
	cfg            Config
	message        Message
	commitStrategy commitStrategy
}

func (j job) Work(id int) {
	transactionID := j.message.GetTransactionID()
	j.commitStrategy.beforeHandler(transactionID, j.message.Topic, j.message.Partition, j.message.Offset)

	ctx, cancel := context.WithTimeout(context.Background(), j.cfg.Timeout)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- invokeMessageHandler(j.message, j.cfg)
		j.commitStrategy.afterHandler(transactionID, j.message.Topic, j.message.Partition, j.message.Offset)
	}()

	select {
	case <-ctx.Done():
		Logger().Debug(transactionID, "Topic: %s at %d/%d ==> %v\n", j.message.Topic, j.message.Partition, j.message.Offset, ctx.Err())
		j.commitStrategy.afterHandler(transactionID, j.message.Topic, j.message.Partition, j.message.Offset)
		<-done
	case err := <-done:
		if err != nil && j.cfg.ErrorHandler != nil {
			j.cfg.ErrorHandler(err, &j.message)
		}
	}
}

func invokeMessageHandler(message Message, cfg Config) error {
	defer func() {
		if r := recover(); r != nil {
			invokeErrorHandler(
				fmt.Errorf("invokeMessageHandler.Panic: While processing %v, trace : %s", r, string(debug.Stack())),
				&message,
				cfg,
			)
		}
	}()
	return cfg.MessageHandler(message)
}
