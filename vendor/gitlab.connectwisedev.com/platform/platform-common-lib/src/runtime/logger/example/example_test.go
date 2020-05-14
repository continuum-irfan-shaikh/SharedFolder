package main

import (
	"testing"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
)

func Test_main_function(t *testing.T) {
	t.Run("main function", func(t *testing.T) {
		logger.Create(logger.Config{Name: "Logger-1", MaxSize: 1, Destination: logger.DISCARD})
		logger.Create(logger.Config{Name: "Logger-2", MaxSize: 1, Destination: logger.DISCARD})
		main()
	})

}
