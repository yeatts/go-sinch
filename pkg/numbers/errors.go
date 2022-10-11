package numbers

import "github.com/thezmc/go-sinch/pkg/sinch"

const (
	ProjectIDRequiredError     = sinch.Error("project ID is required")
	RegionCodeRequiredError    = sinch.Error("region code is required")
	TypeRequiredError          = sinch.Error("type is required")
	PhoneNumberRequiredError   = sinch.Error("phone number is required")
	MissingConfigurationError  = sinch.Error("either smsConfiguration or voiceConfiguration or both must be set")
	ServicePlanIDRequiredError = sinch.Error("service plan ID is required")
	AppIDRequiredError         = sinch.Error("app ID is required")
)
