package webClient

import (
	"errors"
	"net/http"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/circuit"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/communication/http/client"
)

const (
	proxyURL string = "%s://%s:%d"
)

var err5xx = errors.New("Server returned 5xx status code")

//httpClientServiceImpl implements HTTPCommandService
type httpClientServiceImpl struct {
	config     ClientConfig
	httpClient *http.Client
}

// Create create client
func (hc *httpClientServiceImpl) Create() {
	if hc.httpClient != nil {
		return
	}

	//  Here false value present that we want to use proxy configuration, if it present
	hc.httpClient = client.Basic(clientConfig(hc.config), false)
}

//Do sends Request to the Server
func (hc *httpClientServiceImpl) Do(request *http.Request) (*http.Response, error) {
	hc.Create()

	var (
		response  *http.Response
		err       error
		cbEnabled bool
	)

	commandName := request.URL.Host
	if enabled, ok := circuitBreaker[commandName]; ok {
		cbEnabled = enabled
	}

	err = circuit.Do(commandName, cbEnabled, func() error {
		response, err = hc.httpClient.Do(request)
		if err != nil {
			err = checkOffline(err)
			return err
		}

		if response.StatusCode >= http.StatusInternalServerError {
			return err5xx
		}

		return nil
	}, nil)

	if err == err5xx {
		return response, nil
	}

	return response, err
}

// SetCheckRedirect set CheckRedirect to client
func (hc *httpClientServiceImpl) SetCheckRedirect(cr func(req *http.Request, via []*http.Request) error) {
	hc.httpClient.CheckRedirect = cr
}

func clientConfig(config ClientConfig) *client.Config {
	cfg := &client.Config{
		TimeoutMinute:               config.TimeoutMinute,
		MaxIdleConns:                config.MaxIdleConns,
		MaxIdleConnsPerHost:         config.MaxIdleConnsPerHost,
		MaxConnsPerHost:             0,
		IdleConnTimeoutMinute:       config.IdleConnTimeoutMinute,
		DialTimeoutSecond:           config.DialTimeoutSecond,
		DialKeepAliveSecond:         config.DialKeepAliveSecond,
		TLSHandshakeTimeoutSecond:   config.TLSHandshakeTimeoutSecond,
		ExpectContinueTimeoutSecond: config.ExpectContinueTimeoutSecond,
		UseIEProxy:                  config.UseIEProxy,
	}
	cfg.Proxy = client.Proxy{
		Address:  config.ProxySetting.IP,
		Port:     config.ProxySetting.Port,
		UserName: config.ProxySetting.UserName,
		Password: config.ProxySetting.Password,
		Protocol: config.ProxySetting.Protocol,
	}
	return cfg
}
