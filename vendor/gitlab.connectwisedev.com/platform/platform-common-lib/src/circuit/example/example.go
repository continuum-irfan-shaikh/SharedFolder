package main

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/circuit"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/utils"
)

const commandName = "Example-Command"

func main() {
	logger.Create(logger.Config{}) // nolint
	transaction := utils.GetTransactionID()
	circuit.Logger = logger.Get
	circuit.Register(transaction, commandName, circuit.New(), nil)

	for index := 0; index < 100; index++ {
		err := circuit.Do(commandName, true, func() error {
			return fmt.Errorf("Error ==> %v", index)
		}, nil)

		if err != nil {
			logger.Get().Error(transaction, "Error", "%v ==> %v", err, index)
		}
	}
}
