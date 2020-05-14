package webClient

import (
	"net/url"

	"github.com/pkg/errors"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/circuit"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
)

//circuitBreaker hold circuit breaker command name and its enabled status
var circuitBreaker = make(map[string]bool)

//CircuitBreakerConfig hold circuit breaker config for a host, service connects to
type CircuitBreakerConfig struct {
	// CircuitBreaker - Configuration for the Circuit breaker
	// Default config - circuitBreaker()
	CircuitBreaker *circuit.Config
	//BaseURL of the host for which circuit breaker needs to be configured
	BaseURL string
	//StateChangeCallback - callback to be used on circuit breaker state change
	//Default - logging the circuit breaker state change
	StateChangeCallback func(transaction string, commandName string, state string)
}

// RegisterCircuitBreaker - register circuit breaker for all hosts, service connects to.
//  This function is not safe for concurrent use.
func RegisterCircuitBreaker(cb []CircuitBreakerConfig) error {

	circuit.Logger = logger.Get

	for _, cbCfg := range cb {
		if cbCfg.BaseURL != "" {

			u, err := url.Parse(cbCfg.BaseURL)
			if err != nil {
				return errors.Wrapf(err, "Failed to parse url: %s", cbCfg.BaseURL)
			}

			hostname := u.Host
			circuit.Register("", hostname, cbCfg.circuitBreaker(), cbCfg.StateChangeCallback)
			circuitBreaker[hostname] = cbCfg.CircuitBreaker.Enabled
		} else {
			return errors.New("Missing BaseURL in circuit breaker config")
		}
	}
	return nil
}

//circuitBreaker- set default config for circuit breaker
func (c *CircuitBreakerConfig) circuitBreaker() *circuit.Config {
	if c.CircuitBreaker == nil {
		c.CircuitBreaker = &circuit.Config{
			Enabled: true, TimeoutInSecond: 3, MaxConcurrentRequests: 15000,
			ErrorPercentThreshold: 25, RequestVolumeThreshold: 500, SleepWindowInSecond: 10,
		}
	}
	return c.CircuitBreaker
}
