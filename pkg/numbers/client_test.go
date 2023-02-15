package numbers

import (
	"net/http"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thezmc/go-sinch/pkg/api"
	"github.com/thezmc/go-sinch/pkg/sinch"
)

var FakeError = sinch.Error("test error")
var OkHandler = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type MockNumbersAction struct {
	mock.Mock
}

func (m *MockNumbersAction) IsNumbersAction() {}

func (m *MockNumbersAction) Request() *sinch.MockAPIRequest {
	args := m.Called()
	return lo.ToPtr(args.Get(0).(sinch.MockAPIRequest)) //nolint no other way to do this
}

func (m *MockNumbersAction) Response() *sinch.MockAPIResponse {
	args := m.Called()
	return lo.ToPtr(args.Get(0).(sinch.MockAPIResponse)) //nolint no other way to do this
}

func Test_Mock_Implementations(t *testing.T) {
	var _ sinch.Action[*sinch.MockAPIRequest, *sinch.MockAPIResponse] = new(MockNumbersAction)
	var _ sinch.APIRequest = new(sinch.MockAPIRequest)
	var _ sinch.APIResponse = new(sinch.MockAPIResponse)
}

func Test_WithProjectID(t *testing.T) {
	client := new(Client).WithProjectID("foo")
	if client.ProjectID != "foo" {
		t.Errorf("expected ProjectID to be foo, got %s", client.ProjectID)
	}
}

func Test_WithKeyID(t *testing.T) {
	client := new(Client).WithKeyID("bar")
	if client.KeyID != "bar" {
		t.Errorf("expected KeyID to be bar, got " + client.KeyID)
	}
}

func Test_WithKeySecret(t *testing.T) {
	client := new(Client).WithKeySecret("baz")
	if client.KeySecret != "baz" {
		t.Errorf("expected KeySecret to be baz, got " + client.KeyID)
	}
}

func Test_Client_Validate(t *testing.T) {
	var c *Client
	testProjectID := "testProjectID"
	testKeyID := "testKeyID"
	testKeySecret := "testKeySecret"
	tests := map[string]struct {
		configFn    func()
		expectedErr error
	}{
		"missing project id": {
			configFn: func() {
				c = new(Client).WithKeyID(testKeyID).WithKeySecret(testKeySecret)
			},
			expectedErr: ProjectIDRequiredError,
		},
		"missing key ID": {
			configFn: func() {
				c = new(Client).WithProjectID(testProjectID).WithKeySecret(testKeySecret)
			},
			expectedErr: KeyIDRequiredError,
		},
		"missing key secret": {
			configFn: func() {
				c = new(Client).WithProjectID(testProjectID).WithKeyID(testKeyID)
			},
			expectedErr: KeySecretRequiredError,
		},
		"no error": {
			configFn: func() {
				c = new(Client).WithProjectID(testProjectID).WithKeyID(testKeyID).WithKeySecret(testKeySecret)
			},
			expectedErr: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.configFn()
			if test.expectedErr != nil {
				assert.ErrorContains(t, c.Validate(), test.expectedErr.Error())
			} else {
				assert.NoError(t, c.Validate())
			}
		})
	}
}

func Test_URL(t *testing.T) {
	client := new(Client).WithProjectID("foo").WithSinchAPI(new(api.Client).WithBaseURL(BaseURLv1))
	if client.URL() != "https://numbers.api.sinch.com/v1/projects/foo" {
		t.Errorf("expected URL to be https://numbers.api.sinch.com/v1/projects/foo, got %s", client.URL())
	}
}
