package sinch

import "github.com/stretchr/testify/mock"

type MockAPIRequest struct {
	mock.Mock
}

func (m *MockAPIRequest) Validate() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAPIRequest) ExpectedStatusCode() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockAPIRequest) Method() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockAPIRequest) QueryString() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockAPIRequest) Body() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockAPIRequest) Path() string {
	args := m.Called()
	return args.String(0)
}

type MockAPIResponse struct {
	mock.Mock
}

func (m *MockAPIResponse) FromJSON([]byte) error {
	args := m.Called()
	return args.Error(0)
}

type MockAPI struct {
	mock.Mock
}

func (m *MockAPI) Do(client APIClient, req APIRequest, recv APIResponse) error {
	args := m.Called(client, req, recv)
	return args.Error(0)
}

type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) Validate() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAPIClient) URL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockAPIClient) Do(req APIRequest, resp APIResponse) error {
	args := m.Called(req, resp)
	return args.Error(0)
}

type MockValidatable struct {
	mock.Mock
}

func (m *MockValidatable) Validate() error {
	args := m.Called()
	return args.Error(0)
}

type MockAPIAction struct {
	mock.Mock
}

func (m *MockAPIAction) Request() *MockAPIRequest {
	args := m.Called()
	return args.Get(0).(*MockAPIRequest)
}

func (m *MockAPIAction) Response() *MockAPIResponse {
	args := m.Called()
	return args.Get(0).(*MockAPIResponse)
}
