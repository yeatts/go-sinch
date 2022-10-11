package numbers

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/thezmc/go-sinch/pkg/api"
	"github.com/thezmc/go-sinch/pkg/sinch"
)

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

func Test_Validate(t *testing.T) {
	client := new(Client)
	if err := client.Validate(); err != ProjectIDRequiredError {
		t.Errorf("expected ProjectIDRequiredError, got %v", err)
	}
}

func Test_URL(t *testing.T) {
	client := new(Client).WithProjectID("foo").WithSinchAPI(new(api.Client).WithBaseURL(BaseURLv1))
	if client.URL() != "https://numbers.api.sinch.com/v1/projects/foo" {
		t.Errorf("expected URL to be https://numbers.api.sinch.com/v1/projects/foo, got %s", client.URL())
	}
}
