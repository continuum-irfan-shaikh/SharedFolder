package metric

import (
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/communication/udp"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/util"
)

//Config - Holds all the configuration for Metric object Publishing
type Config struct {
	// Communication - UDP Communication configuration
	// default - Communication : Default UDP Config
	Communication *udp.Config

	// Namespace - Namespace of a Metric collector for unique identification
	// Default - value is <HostName>
	Namespace string
}

// New - Default configuration object having default values
// values - Address: "localhost", PortNumber: "7000", Namespace : ""
var New = func() *Config {
	return &Config{
		Communication: udp.New(),
		Namespace:     "",
	}
}

// GetNamespace - Return service name space for metric
func (c *Config) GetNamespace() string {
	ns := c.Namespace
	if ns == "" {
		ns = util.Hostname(util.ProcessName())
	}
	return ns
}
