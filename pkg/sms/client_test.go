package sms

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thezmc/go-sinch/pkg/sinch"
)

func Test_Client_Implementations(t *testing.T) {
	var _ sinch.APIClient = new(Client)
}

func Test_Client_Validate(t *testing.T) {
	var c *Client
	tests := map[string]struct {
		configFn    func()
		expectedErr error
	}{
		"missing plan id": {
			configFn: func() {
				c = new(Client)
			},
			expectedErr: NoPlanIDError,
		},
		"no error": {
			configFn: func() {
				c = new(Client).WithPlanID("test")
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
