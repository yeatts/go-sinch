package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/thezmc/go-sinch/pkg/sinch"
)

var FakeError = sinch.Error("test error")
var OkHandler = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Test_Do(t *testing.T) {
	var mockRequest *sinch.MockAPIRequest
	var mockResponse *sinch.MockAPIResponse
	var mockClient *sinch.MockAPIClient
	var mockHTTPSrv *httptest.Server

	client := new(Client).WithBaseURL("https://notareal.domain")

	tests := map[string]struct {
		configFn func()
		wantErr  bool
	}{
		"validation error": {
			configFn: func() {
				mockClient.On("Validate").Return(FakeError)
				mockRequest.On("Validate").Return(FakeError)
			},
			wantErr: true,
		},
		"query string error": {
			configFn: func() {
				mockClient.On("Validate").Return(nil)
				mockRequest.On("Validate").Return(nil)
				mockRequest.On("QueryString").Return("", FakeError)
			},
			wantErr: true,
		},
		"body error": {
			configFn: func() {
				mockClient.On("Validate").Return(nil)
				mockRequest.On("Validate").Return(nil)
				mockRequest.On("QueryString").Return("", nil)
				mockRequest.On("Body").Return([]byte{}, FakeError)
			},
			wantErr: true,
		},
		"new request error": {
			configFn: func() {
				mockClient.On("Validate").Return(nil)
				mockRequest.On("Validate").Return(nil)
				mockRequest.On("QueryString").Return("", nil)
				mockRequest.On("Body").Return([]byte{}, nil)
				mockRequest.On("Method").Return("NOT A REAL METHOD")
				mockRequest.On("Path").Return("/not/a/real/path")
				mockClient.On("URL").Return("https://notareal.domain")
			},
			wantErr: true,
		},
		"http client do error": {
			configFn: func() {
				mockHTTPSrv = httptest.NewServer(http.RedirectHandler("https://notareal.domain", http.StatusMovedPermanently))
				mockClient.On("Validate").Return(nil)
				mockRequest.On("Validate").Return(nil)
				mockRequest.On("QueryString").Return("", nil)
				mockRequest.On("Body").Return([]byte{}, nil)
				mockRequest.On("Method").Return("GET")
				mockRequest.On("Path").Return("/not/a/real/path")
				mockClient.On("URL").Return(mockHTTPSrv.URL)
				mockClient.On("Authenticate", mock.Anything).Return((*http.Request)(nil), nil)
				httpClient := mockHTTPSrv.Client()
				httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return FakeError
				}
				client.WithHTTPClient(mockHTTPSrv.Client())
			},
			wantErr: true,
		},
		"unexpected status code": {
			configFn: func() {
				mockHTTPSrv = httptest.NewServer(http.HandlerFunc(OkHandler))
				mockClient.On("Validate").Return(nil)
				mockRequest.On("Validate").Return(nil)
				mockRequest.On("QueryString").Return("", nil)
				mockRequest.On("Body").Return([]byte{}, nil)
				mockRequest.On("Method").Return("GET")
				mockRequest.On("Path").Return("/not/a/real/path")
				mockClient.On("Authenticate", mock.Anything).Return((*http.Request)(nil), nil)
				mockRequest.On("ExpectedStatusCode").Return(http.StatusTeapot)
				mockClient.On("URL").Return(mockHTTPSrv.URL)
				client.WithHTTPClient(mockHTTPSrv.Client())
			},
			wantErr: true,
		},
		"unmarshal error": {
			configFn: func() {
				mockHTTPSrv = httptest.NewServer(http.HandlerFunc(OkHandler))
				mockClient.On("Validate").Return(nil)
				mockRequest.On("Validate").Return(nil)
				mockRequest.On("QueryString").Return("", nil)
				mockRequest.On("Body").Return([]byte{}, nil)
				mockRequest.On("Method").Return("GET")
				mockRequest.On("Path").Return("/not/a/real/path")
				mockClient.On("Authenticate", mock.Anything).Return((*http.Request)(nil), nil)
				mockRequest.On("ExpectedStatusCode").Return(http.StatusOK)
				mockResponse.On("FromJSON").Return(FakeError)
				mockClient.On("URL").Return(mockHTTPSrv.URL)
				client.WithHTTPClient(mockHTTPSrv.Client())
			},
			wantErr: true,
		},
		"success": {
			configFn: func() {
				mockHTTPSrv = httptest.NewServer(http.HandlerFunc(OkHandler))
				mockClient.On("Validate").Return(nil)
				mockRequest.On("Validate").Return(nil)
				mockRequest.On("QueryString").Return("", nil)
				mockRequest.On("Body").Return([]byte{}, nil)
				mockRequest.On("Method").Return("GET")
				mockRequest.On("Path").Return("/not/a/real/path")
				mockClient.On("Authenticate", mock.Anything).Return((*http.Request)(nil), nil)
				mockRequest.On("ExpectedStatusCode").Return(http.StatusOK)
				mockResponse.On("FromJSON").Return(nil)
				mockClient.On("URL").Return(mockHTTPSrv.URL)
				client.WithHTTPClient(mockHTTPSrv.Client())
			},
			wantErr: false,
		},
		"success with body greater than 0": {
			configFn: func() {
				mockHTTPSrv = httptest.NewServer(http.HandlerFunc(OkHandler))
				mockClient.On("Validate").Return(nil)
				mockRequest.On("Validate").Return(nil)
				mockRequest.On("QueryString").Return("", nil)
				mockRequest.On("Body").Return([]byte("hello"), nil)
				mockRequest.On("Method").Return("GET")
				mockRequest.On("Path").Return("/not/a/real/path")
				mockClient.On("Authenticate", mock.Anything).Return((*http.Request)(nil), nil)
				mockRequest.On("ExpectedStatusCode").Return(http.StatusOK)
				mockResponse.On("FromJSON").Return(nil)
				mockClient.On("URL").Return(mockHTTPSrv.URL)
				client.WithHTTPClient(mockHTTPSrv.Client())
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockRequest = new(sinch.MockAPIRequest)
			mockResponse = new(sinch.MockAPIResponse)
			mockClient = new(sinch.MockAPIClient)
			test.configFn()
			if err := client.Do(mockClient, mockRequest, mockResponse); (err != nil) != test.wantErr {
				t.Errorf("Client.Do() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func Test_Validate(t *testing.T) {
	client := new(Client)

	tests := map[string]struct {
		configFn    func()
		expectedErr error
	}{
		"no base url": {
			configFn: func() {
				client.WithAuthToken("test").WithBaseURL("")
			},
			expectedErr: NoBaseURLError,
		},
		"no http client": {
			configFn: func() {
				client.WithAuthToken("test").WithBaseURL("test").WithHTTPClient(nil)
			},
			expectedErr: NilHTTPClientError,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.configFn()
			err := client.Validate()
			if err != test.expectedErr {
				t.Errorf("Client.Validate() error = %v, expectedErr %v", err, test.expectedErr)
			}
		})
	}
}
