package main

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
)

func main() {
	createLogger([]string{"Logger-1", "Logger-2"})
	printMessage("Logger-1", "test")
	printMessageWithParam("Logger-2", "test")

	updateLogLevel("Logger-1", "test", logger.DEBUG)
	printMessage("Logger-1", "test")
	updateLogLevel("Logger-2", "test", logger.DEBUG)
	printMessageWithParam("Logger-2", "test")

	updateLogLevel("Logger-1", "test", logger.TRACE)
	printMessage("Logger-1", "test")
	updateLogLevel("Logger-2", "test", logger.TRACE)
	printMessageWithParam("Logger-2", "test")

	updateLogLevel("Logger-1", "test", logger.WARN)
	printMessage("Logger-2", "test")
	updateLogLevel("Logger-2", "test1", logger.WARN)
	printMessageWithParam("Logger-1", "test")

	updateLogLevel("Logger-1", "test", logger.ERROR)
	printMessage("Logger-1", "test")
	updateLogLevel("Logger-2", "test", logger.ERROR)
	printMessageWithParam("Logger-2", "test")

	updateLogLevel("Logger-1", "test", logger.FATAL)
	printMessage("Logger-1", "test")
	updateLogLevel("Logger-2", "test", logger.FATAL)
	printMessageWithParam("Logger-2", "test")

	updateLogLevel("Logger-1", "test", logger.OFF)
	printMessage("Logger-1", "test")
	updateLogLevel("Logger-2", "test", logger.OFF)
	printMessageWithParam("Logger-2", "test")
}

// createLogger is a function to explain how to create a logger instance
func createLogger(names []string) {
	for _, name := range names {
		_, err := logger.Create(logger.Config{Name: name, MaxSize: 1})
		if err != nil {
			fmt.Println(err)
		}
	}
}

//printMessage is a function to showcase all message type prinitng
func printMessage(loggerName string, transaction string) {
	log := logger.GetViaName(loggerName)
	log.Trace(transaction, "This is a TRACE Message")
	log.Debug(transaction, "This is a DEBUG Message")
	log.Info(transaction, "This is a INFO Message")
	log.Warn(transaction, "This is a WARN Message")
	log.Error(transaction, "ERROR-CODE", "This is a ERROR Message")
	log.Fatal(transaction, "FATAL-CODE", "This is a FATAL Message")
}

//printMessageWithParam is a function to showcase all message type prinitng
func printMessageWithParam(loggerName string, transaction string) {
	log := logger.GetViaName(loggerName)
	log.Trace(transaction, "This is a %s Message", "TRACE")
	log.Debug(transaction, "This is a %s Message", "DEBUG")
	log.Info(transaction, "This is a %s Message", "INFO")
	log.Warn(transaction, "This is a %s Message", "WARN")
	log.Error(transaction, "ERROR-CODE", "This is a %s Message", "ERROR")
	log.Fatal(transaction, "FATAL-CODE", "This is a %s Message", "FATAL")
}

// updateLogLevel is a function to showcase how to update log level
func updateLogLevel(loggerName string, transaction string, loglevel logger.LogLevel) {
	log, _ := logger.Update(logger.Config{Name: loggerName, MaxSize: 1, LogLevel: loglevel}) //nolint
	log.Info(transaction, "-----------------------------------------------------")
	log.Info(transaction, "Update Loglevel to %v", loglevel)
	log.Info(transaction, "-----------------------------------------------------")
}
