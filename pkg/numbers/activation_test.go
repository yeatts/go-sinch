package numbers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thezmc/go-sinch/pkg/sinch"
)

func Test_Activation_Implementations(t *testing.T) {
	var _ sinch.Action[*ActivationRequest, *ActivationResponse] = new(Activation)
	var _ sinch.APIRequest = new(ActivationRequest)
	var _ sinch.APIResponse = new(ActivationResponse)
}

func Test_ActivationRequest_Validate(t *testing.T) {
	ar := new(ActivationRequest)

	tests := map[string]struct {
		configFn    func()
		expectedErr error
	}{
		"missing number": {
			configFn: func() {
				ar = new(ActivationRequest).WithPhoneNumber("")
			},
			expectedErr: PhoneNumberRequiredError,
		},
		"missing configuation": {
			configFn: func() {
				ar = new(ActivationRequest)
			},
			expectedErr: MissingConfigurationError,
		},
		"missing service plan": {
			configFn: func() {
				ar = new(ActivationRequest).WithPhoneNumber("1234567890").WithSMSConfiguration("")
			},
			expectedErr: ServicePlanIDRequiredError,
		},
		"missing app ID": {
			configFn: func() {
				ar = new(ActivationRequest).WithPhoneNumber("1234567890").WithVoiceConfiguration("")
			},
			expectedErr: AppIDRequiredError,
		},
		"no errors": {
			configFn: func() {
				ar = new(ActivationRequest).WithPhoneNumber("1234567890").WithSMSConfiguration("test")
			},
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.configFn()
			if test.expectedErr != nil {
				assert.ErrorContains(t, ar.Validate(), test.expectedErr.Error())
			} else {
				assert.NoError(t, ar.Validate())
			}
		})
	}
}
