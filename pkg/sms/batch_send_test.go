package sms

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thezmc/go-sinch/pkg/sinch"
)

func Test_BatchSend_Implementations(t *testing.T) {
	var _ sinch.Action[*BatchSendRequest, *BatchSendResponse] = new(BatchSend)
	var _ sinch.APIRequest = new(BatchSendRequest)
	var _ sinch.APIResponse = new(BatchSendResponse)
}

func Test_BatchSendRequest_Validate(t *testing.T) {
	var bsr *BatchSendRequest
	tests := map[string]struct {
		configFn    func()
		expectedErr error
	}{
		"missing to": {
			configFn: func() {
				bsr = new(BatchSendRequest)
			},
			expectedErr: InvalidToNumberError,
		},
		"missing from": {
			configFn: func() {
				bsr = new(BatchSendRequest)
			},
			expectedErr: InvalidFromNumberError,
		},
		"bad type of number from": {
			configFn: func() {
				bsr = new(BatchSendRequest)
				bsr.FromTypeOfNumber = 7
			},
			expectedErr: InvalidTypeOfNumberError,
		},
		"bad npi": {
			configFn: func() {
				bsr = new(BatchSendRequest)
				bsr.FromNumberPlanIndicator = 19
			},
			expectedErr: InvalidNPIError,
		},
		"bad callback url": {
			configFn: func() {
				bsr = new(BatchSendRequest)
				bsr.CallbackURL = "test"
			},
			expectedErr: InvalidCallbackURLError,
		},
		"bad send at time": {
			configFn: func() {
				bsr = new(BatchSendRequest)
				bsr.SendAt = "test"
			},
			expectedErr: InvalidSendAtError,
		},
		"bad expiry time": {
			configFn: func() {
				bsr = new(BatchSendRequest)
				bsr.ExpireAt = "test"
			},
			expectedErr: InvalidExpireAtError,
		},
		"no errors": {
			configFn: func() {
				bsr = new(BatchSendRequest).
					From("1234567890").
					To("1234567890").
					WithMessageBody("test")
			},
			expectedErr: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.configFn()
			if test.expectedErr != nil {
				assert.ErrorContains(t, bsr.Validate(), test.expectedErr.Error())
			} else {
				assert.NoError(t, bsr.Validate())
			}
		})
	}
}
