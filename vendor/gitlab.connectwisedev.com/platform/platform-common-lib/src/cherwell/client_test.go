package cherwell

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/webClient"
)

func catchPanics(t *testing.T) func() {
	return func() {
		if r := recover(); r != nil {
			t.Fatalf("panic occured:\n%v", r)
		}
	}
}

func newMockHandler(method, path, resp string, statusCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))
	})
}

func newTestServer() (*httptest.Server, *http.ServeMux) {
	defaultTokenResponse := []byte(`{
		"access_token": "access_token",
		"token_type": "bearer",
		"expires_in": 14399,
		"refresh_token": "refresh_token",
		"as:client_id": "client_id",
		"username": "username",
		".issued": "Tue, 31 Jul 2018 14:46:46 GMT",
		".expires": "Tue, 31 Jul 2018 18:46:46 GMT"
	  }`)

	mux := http.NewServeMux()
	mux.HandleFunc(tokenEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Write(defaultTokenResponse)
	})
	server := httptest.NewServer(mux)
	return server, mux
}

func newFailedTestServer() (*httptest.Server, *http.ServeMux) {
	defaultTokenResponse := []byte(`{
		"error": "access_token",		
		"error_description": "Tue, 31 Jul 2018 18:46:46 GMT"
	  }`)

	mux := http.NewServeMux()
	mux.HandleFunc(tokenEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Write(defaultTokenResponse)
	})
	server := httptest.NewServer(mux)
	return server, mux
}

func TestNewClientSuccess(t *testing.T) {
	server, _ := newTestServer()
	conf := Config{
		Host: server.URL,
	}
	_, err := NewClient(conf, getWebClient())
	assert.NoError(t, err)
}

func TestNewClientFailNilHTTPClient(t *testing.T) {
	conf := Config{
		Host: "empty/URL",
	}
	_, err := NewClient(conf, nil)
	assert.EqualError(t, err, "cherwell.NewClient: httpClient can not be nil")
}

func TestNewClientFailObtainAccesToken(t *testing.T) {
	conf := Config{
		Host: "http://Invalid hosthname",
	}
	_, err := NewClient(conf, getWebClient())
	assert.EqualError(t, err, "cherwell.NewClient: authentication error: getAccessToken authentication failed while creating request: parse http://Invalid hosthname/token: invalid character \" \" in host name")
}

func TestObtainAccesTokenSuccess(t *testing.T) {
	server, _ := newTestServer()
	conf := Config{
		Host: server.URL,
	}
	client := &Client{config: conf, client: getWebClient()}
	assert.NoError(t, client.getAccessToken())
}

func TestObtainAccesTokenFailOnRequest(t *testing.T) {
	conf := Config{
		Host: "http://Invalid hosthname",
	}
	client := &Client{config: conf, client: getWebClient()}
	err := client.getAccessToken()
	assert.EqualError(t, err, "getAccessToken authentication failed while creating request: parse http://Invalid hosthname/token: invalid character \" \" in host name")
}

func TestObtainAccesTokenFailOnBadResponseCode(t *testing.T) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(errorResponse{Err: "Invalid request", ErrorDescription: "password is empty"})
		assert.NoError(t, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))

	conf := Config{
		Host: server.URL,
	}
	client := &Client{config: conf, client: getWebClient()}
	err := client.getAccessToken()
	assert.EqualError(t, err, "Error: Invalid request\nDescription: password is empty\n")
}

func TestClientRefreshAccessToken(t *testing.T) {
	server, _ := newTestServer()
	failedServer, _ := newFailedTestServer()
	conf := Config{
		Host: server.URL,
	}
	failedConf := Config{
		Host: failedServer.URL,
	}
	client := &Client{config: conf, client: getWebClient()}

	tests := []struct {
		name          string
		tokenResponse *tokenResponse
		client        webClient.HTTPClientService
		config        Config
		wantErr       bool
	}{
		{name: "Success token refresh",
			tokenResponse: nil,
			client:        client.client,
			config:        conf,
			wantErr:       false,
		},
		{name: "Failed token refresh",
			tokenResponse: nil,
			client:        client.client,
			config:        failedConf,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				tokenResponse: tt.tokenResponse,
				client:        tt.client,
				config:        tt.config,
			}
			if err := c.refreshAccessToken(); (err != nil) != tt.wantErr {
				t.Errorf("Client.refreshAccessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRaceConditionOnRefreshToken(t *testing.T) {
	server, mux := newTestServer()
	conf := Config{
		Host: server.URL,
	}
	client := &Client{config: conf, client: getWebClient()}
	client.getAccessToken()

	resp := `{
          "busObPublicId": "pub_id_1",
          "busObRecId": "rec_id_1",
          "cacheKey": "string",
          "fieldValidationErrors": [],
          "notificationTriggers": [],
          "errorCode": "",
          "errorMessage": "",
          "hasError": false
    }`

	test := struct {
		name   string
		method string
		path   string
	}{
		name:   "Success token refresh",
		method: http.MethodGet,
		path:   fmt.Sprintf(getBOByRecIDEndpoint, "12", "123"),
	}

	mux.Handle(test.path, newMockHandler(test.method, test.path, resp, 401))

	t.Run(test.name, func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, method string, path string) {
				for i := 0; i < 500; i++ {
					respEntity := new(businessObjectResponse)
					client.performRequest(test.method, test.path, nil, respEntity)
				}
				wg.Done()
			}(&wg, test.method, test.path)
		}
		wg.Wait()
	})
}

func TestFormatCherwellResponse(t *testing.T) {
	assert.Equal(t, []byte("  abc abc def  hijc  df  "), formatCherwellResponse([]byte("\n\rabc\nabc\rdef\r\rhijc\n\ndf\n\r")))
}

func TestFailedHTTPCode(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test Case 1: HTTP Code 418 - failed request",
			args: args{code: http.StatusTeapot},
			want: true,
		},
		{
			name: "Test Case 2: HTTP Code 399 - successful request",
			args: args{code: http.StatusPermanentRedirect},
			want: false,
		},
		{
			name: "Test Case 3: HTTP Code 200 - successful request",
			args: args{code: http.StatusOK},
			want: false,
		},
		{
			name: "Test Case 4: HTTP Code 777 - failed request",
			args: args{code: 777},
			want: true,
		},
		{
			name: "Test Case 5: HTTP Code 0 - failed request",
			args: args{code: 0},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := failedHTTPCode(tt.args.code); got != tt.want {
				t.Errorf("failedHTTPCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClientWithLogger(t *testing.T) {
	server, _ := newTestServer()
	type args struct {
		conf   Config
		client webClient.HTTPClientService
		logger Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{
			name: "Test Case 1: successful",
			args: args{
				conf:   Config{Host: server.URL},
				client: getWebClient(),
				logger: NewNopLogger(),
			},
			wantErr: false,
		},
		{
			name: "Test Case 2: failed because of nil http client",
			args: args{
				conf:   Config{Host: server.URL},
				client: nil,
				logger: NewNopLogger(),
			},
			wantErr: true,
			err:     "cherwell.NewClientWithLogger: cherwell.NewClient: httpClient can not be nil",
		},
		{
			name: "Test Case 3: failed because of nil logger",
			args: args{
				conf:   Config{Host: server.URL},
				client: getWebClient(),
				logger: nil,
			},
			wantErr: true,
			err:     "cherwell.NewClientWithLogger: logger can not be nil",
		},
		{
			name: "Test Case 4: failed because of bad auth response",
			args: args{
				conf:   Config{Host: "testfail/com"},
				client: getWebClient(),
				logger: NewNopLogger(),
			},
			wantErr: true,
			err:     "cherwell.NewClientWithLogger: cherwell.NewClient: authentication error: getAccessToken authentication request failed: Post testfail/com/token: unsupported protocol scheme \"\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClientWithLogger(tt.args.conf, tt.args.client, tt.args.logger)
			if err != nil {
				assert.EqualError(t, err, tt.err)
				assert.Empty(t, got)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}
