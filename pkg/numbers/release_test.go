package numbers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thezmc/go-sinch/pkg/sinch"
)

func Test_Release_Implementations(t *testing.T) {
	var _ sinch.Action[*ReleaseRequest, *ReleaseResponse] = new(Release)
	var _ sinch.APIRequest = new(ReleaseRequest)
	var _ sinch.APIResponse = new(ReleaseResponse)
}

func Test_ReleaseRequest_Validate(t *testing.T) {
	rr := new(ReleaseRequest)

	tests := map[string]struct {
		configFn    func()
		expectedErr error
	}{
		"missing number": {
			configFn: func() {
				rr = new(ReleaseRequest).WithPhoneNumber("")
			},
			expectedErr: PhoneNumberRequiredError,
		},
		"no errors": {
			configFn: func() {
				rr = new(ReleaseRequest).WithPhoneNumber("1234567890")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.configFn()
			if test.expectedErr != nil {
				assert.ErrorContains(t, rr.Validate(), test.expectedErr.Error())
			} else {
				assert.NoError(t, rr.Validate())
			}
		})
	}
}
