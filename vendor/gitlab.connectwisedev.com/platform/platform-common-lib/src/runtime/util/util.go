package util

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var (
	processName    string
	invocationPath string
	hostName       string
)

// ProcessName is a function to return process name for a binary
func ProcessName() string {
	if processName == "" {
		processName = filepath.Base(os.Args[0])
	}
	return processName
}

// InvocationPath is a function to return Invocation path for a binary
func InvocationPath() string {
	if invocationPath == "" {
		invocationPath = filepath.Dir(os.Args[0])
	}
	return invocationPath
}

// Hostname is a function to return Hostname for a machine; in case of error it sends default value
func Hostname(defaultValue string) string {
	if hostName == "" {
		h, err := os.Hostname()
		if err != nil {
			return defaultValue
		}
		hostName = h
	}
	return hostName
}

// LocalIPAddress returns the non loopback local IP of the host
func LocalIPAddress() []string {
	localIPAddress := make([]string, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return localIPAddress
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIPAddress = append(localIPAddress, ipnet.IP.String())
			}
		}
	}
	return localIPAddress
}

// NotifyStopSignal - Function to execute callback on reciving a quit signal
func NotifyStopSignal(stop <-chan bool, callbacks ...func()) error {
	if stop == nil {
		return errors.New("stop: Notify using nil channel")
	}
	if callbacks == nil {
		return errors.New("NotifyStopSignal | Received nil callback")
	}

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		for _, callback := range callbacks {
			callback()
		}
	case <-stop:
	}
	return nil
}
