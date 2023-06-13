package numbers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thezmc/go-sinch/pkg/sinch"
)

func Test_Update_Implementations(t *testing.T) {
	var _ sinch.Action[*UpdateRequest, *UpdateResponse] = new(Update)
	var _ sinch.APIRequest = new(UpdateRequest)
	var _ sinch.APIResponse = new(UpdateResponse)
}

func Test_UpdateRequest_Validate(t *testing.T) {
	ur := new(UpdateRequest)

	tests := map[string]struct {
		configFn    func()
		expectedErr error
	}{
		"missing number": {
			configFn: func() {
				ur = new(UpdateRequest).WithPhoneNumber("")
			},
			expectedErr: PhoneNumberRequiredError,
		},
		"missing configuation": {
			configFn: func() {
				ur = new(UpdateRequest)
			},
			expectedErr: MissingConfigurationError,
		},
		"missing service plan": {
			configFn: func() {
				ur = new(UpdateRequest).WithPhoneNumber("1234567890").WithSMSConfiguration(&RequestSMSConfiguration{ServicePlanID: ""})
			},
			expectedErr: ServicePlanIDRequiredError,
		},
		"with request": {
			configFn: func() {
				u := new(Update).WithRequest(new(UpdateRequest).WithPhoneNumber("1234567890"))
				ur = u.Request()
			},
		},
		"with response": {
			configFn: func() {
				u := new(Update).WithResponse(new(UpdateResponse))
				ur = u.Request()
			},
		},
		"with display name": {
			configFn: func() {
				ur = new(UpdateRequest).WithPhoneNumber("1234567890").WithDisplayName("test")
			},
			expectedErr: nil,
		},
		"with voice configuration": {
			configFn: func() {
				ur = new(UpdateRequest).WithPhoneNumber("1234567890").WithVoiceConfiguration(&RequestVoiceConfiguration{AppID: "test"})
			},
			expectedErr: nil,
		},
		"with sms configuration service plan id": {
			configFn: func() {
				ur = new(UpdateRequest).WithPhoneNumber("1234567890").WithSMSConfigurationServicePlanID("test")
			},
			expectedErr: nil,
		},
		"with sms configuration campaign id": {
			configFn: func() {
				ur = new(UpdateRequest).WithPhoneNumber("1234567890").WithSMSConfigurationCampaignID("test")
			},
			expectedErr: nil,
		},
		"with voice configuration app id": {
			configFn: func() {
				ur = new(UpdateRequest).WithPhoneNumber("1234567890").WithVoiceConfigurationAppID("test")
			},
			expectedErr: nil,
		},
		"get update request body": {
			configFn: func() {
				_, err := new(UpdateRequest).WithPhoneNumber("1234567890").WithDisplayName("test").WithSMSConfiguration(&RequestSMSConfiguration{ServicePlanID: "test"}).WithVoiceConfiguration(&RequestVoiceConfiguration{AppID: "test"}).Body()
				assert.NoError(t, err)
			},
		},
		"get method type": {
			configFn: func() {
				assert.Equal(t, http.MethodPatch, new(UpdateRequest).Method())
			},
		},
		"get path": {
			configFn: func() {
				assert.Equal(t, "/activeNumbers/1234567890", new(UpdateRequest).Path())
			},
		},
		"get expected status code": {
			configFn: func() {
				assert.Equal(t, http.StatusOK, new(UpdateRequest).ExpectedStatusCode())
			},
		},
		"get query string": {
			configFn: func() {
				_, err := new(UpdateRequest).QueryString()
				assert.NoError(t, err)
			},
		},
		"get response from json": {
			configFn: func() {
				err := new(UpdateResponse).FromJSON([]byte(`{"phoneNumber": "1234567890"}`))
				assert.NoError(t, err)
			},
		},
		"no errors": {
			configFn: func() {
				ur = new(UpdateRequest).WithPhoneNumber("1234567890").WithSMSConfiguration(&RequestSMSConfiguration{ServicePlanID: "test"})
			},
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.configFn()
			if test.expectedErr != nil {
				assert.ErrorContains(t, ur.Validate(), test.expectedErr.Error())
			} else {
				assert.NoError(t, ur.Validate())
			}
		})
	}
}
