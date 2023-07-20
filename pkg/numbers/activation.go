package numbers

import (
	"encoding/json"
	"net/http"

	"github.com/thezmc/go-sinch/pkg/sinch"
)

type Activation struct {
	request  *ActivationRequest
	response *ActivationResponse
}

type ActivationRequest struct {
	PhoneNumber        string                     `url:"phoneNumber" json:"-"`
	SMSConfiguration   *RequestSMSConfiguration   `url:"-" json:"smsConfiguration,omitempty"`
	VoiceConfiguration *RequestVoiceConfiguration `url:"-" json:"voiceConfiguration,omitempty"`
}

type ActivationResponse struct {
	PhoneNumber           string                      `json:"phoneNumber"`
	ProjectID             string                      `json:"projectId"`
	DisplayName           string                      `json:"displayName"`
	RegionCode            string                      `json:"regionCode"`
	Type                  string                      `json:"type"`
	Capability            []string                    `json:"capability"`
	Money                 Price                       `json:"money"`
	PaymentIntervalMonths int                         `json:"paymentIntervalMonths"`
	NextChargeDate        string                      `json:"nextChargeDate"`
	ExpireAt              string                      `json:"expireAt"`
	SMSConfiguration      *ResponseSMSConfiguration   `json:"smsConfiguration"`
	VoiceConfiguration    *ResponseVoiceConfiguration `json:"voiceConfiguration,omitempty"`
}

func (a *Activation) IsNumbersAction() {}

func (a *Activation) WithRequest(request *ActivationRequest) *Activation {
	a.request = request
	return a
}

func (a *Activation) WithResponse(response *ActivationResponse) *Activation {
	a.response = response
	return a
}

func (a *Activation) Request() *ActivationRequest {
	return a.request
}

func (a *Activation) Response() *ActivationResponse {
	return a.response
}

func (ar *ActivationRequest) WithPhoneNumber(phoneNumber string) *ActivationRequest {
	ar.PhoneNumber = phoneNumber
	return ar
}

func (ar *ActivationRequest) WithSMSConfiguration(servicePlanID string) *ActivationRequest {
	if ar.SMSConfiguration == nil {
		ar.SMSConfiguration = new(RequestSMSConfiguration)
	}
	ar.SMSConfiguration.ServicePlanID = servicePlanID
	return ar
}

func (ar *ActivationRequest) WithVoiceConfiguration(appID string) *ActivationRequest {
	if ar.VoiceConfiguration == nil {
		ar.VoiceConfiguration = new(RequestVoiceConfiguration)
	}
	ar.VoiceConfiguration.AppID = appID
	return ar
}

func (ar *ActivationRequest) Validate() error {
	var errors sinch.Errors
	if ar.SMSConfiguration == nil && ar.VoiceConfiguration == nil {
		errors = append(errors, MissingConfigurationError)
	}
	if ar.PhoneNumber == "" {
		errors = append(errors, PhoneNumberRequiredError)
	}
	if ar.SMSConfiguration != nil {
		if ar.SMSConfiguration.ServicePlanID == "" {
			errors = append(errors, ServicePlanIDRequiredError)
		}
	}
	if ar.VoiceConfiguration != nil {
		if ar.VoiceConfiguration.AppID == "" {
			errors = append(errors, AppIDRequiredError)
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (ar *ActivationRequest) Body() ([]byte, error) {
	return json.Marshal(ar)
}

func (ar *ActivationRequest) Method() string {
	return http.MethodPost
}

func (ar *ActivationRequest) Path() string {
	return "/availableNumbers/" + ar.PhoneNumber + ":rent"
}

func (ar *ActivationRequest) QueryString() (string, error) {
	return "", nil
}

func (ar *ActivationRequest) ExpectedStatusCode() int {
	return http.StatusOK
}

func (ar *ActivationResponse) FromJSON(bytes []byte) error {
	return json.Unmarshal(bytes, ar)
}
