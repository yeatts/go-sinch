package sms

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thezmc/go-sinch/pkg/sinch"
)

var FakeError = sinch.Error("test error")

func Test_Client_Implementations(t *testing.T) {
	var _ sinch.APIClient = new(Client)
}

func Test_Client_Validate(t *testing.T) {
	var c *Client
	testAuthToken := "testToken"
	testPlanID := "testPlanID"
	tests := map[string]struct {
		configFn    func()
		expectedErr error
	}{
		"missing plan id": {
			configFn: func() {
				c = new(Client).WithAuthToken(testAuthToken)
			},
			expectedErr: NoPlanIDError,
		},
		"missing auth token": {
			configFn: func() {
				c = new(Client).WithPlanID(testPlanID)
			},
			expectedErr: NoAuthTokenError,
		},
		"no error": {
			configFn: func() {
				c = new(Client).WithPlanID(testPlanID).WithAuthToken(testAuthToken)
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
