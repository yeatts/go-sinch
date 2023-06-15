package numbers

import (
	"encoding/json"
	"net/http"

	"github.com/thezmc/go-sinch/pkg/sinch"
)

type Update struct {
	request  *UpdateRequest
	response *UpdateResponse
}

type UpdateRequest struct {
	PhoneNumber        string                     `url:"phoneNumber" json:"-"`
	DisplayName        string                     `url:"-" json:"displayName,omitempty"`
	SMSConfiguration   *RequestSMSConfiguration   `url:"-" json:"smsConfiguration,omitempty"`
	VoiceConfiguration *RequestVoiceConfiguration `url:"-" json:"voiceConfiguration,omitempty"`
}

type UpdateResponse struct {
	PhoneNumber           string                      `json:"phoneNumber"`
	ProjectId             string                      `json:"projectId"`
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

func (u *Update) IsNumbersAction() {}

func (u *Update) WithRequest(request *UpdateRequest) *Update {
	u.request = request
	return u
}

func (u *Update) WithResponse(response *UpdateResponse) *Update {
	u.response = response
	return u
}

func (u *Update) Request() *UpdateRequest {
	return u.request
}

func (u *Update) Response() *UpdateResponse {
	return u.response
}

func (ur *UpdateRequest) WithPhoneNumber(phoneNumber string) *UpdateRequest {
	ur.PhoneNumber = phoneNumber
	return ur
}

func (ur *UpdateRequest) WithDisplayName(displayName string) *UpdateRequest {
	ur.DisplayName = displayName
	return ur
}

func (ur *UpdateRequest) WithSMSConfiguration(servicePlanID string, campaignID string) *UpdateRequest {
	if ur.SMSConfiguration == nil {
		ur.SMSConfiguration = new(RequestSMSConfiguration)
	}
	ur.SMSConfiguration.ServicePlanID = servicePlanID
	ur.SMSConfiguration.CampaignID = campaignID
	return ur
}

func (ur *UpdateRequest) WithVoiceConfiguration(appID string) *UpdateRequest {
	if ur.VoiceConfiguration == nil {
		ur.VoiceConfiguration = new(RequestVoiceConfiguration)
	}
	ur.VoiceConfiguration.AppID = appID
	return ur
}

func (ur *UpdateRequest) Validate() error {
	var errors sinch.Errors
	if ur.PhoneNumber == "" {
		errors = append(errors, PhoneNumberRequiredError)
	}

	if ur.SMSConfiguration != nil {
		if ur.SMSConfiguration.ServicePlanID == "" {
			errors = append(errors, ServicePlanIDRequiredError)
		}
	}
	if ur.VoiceConfiguration != nil {
		if ur.VoiceConfiguration.AppID == "" {
			errors = append(errors, AppIDRequiredError)
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (ur *UpdateRequest) Body() ([]byte, error) {
	return json.Marshal(ur)
}

func (ur *UpdateRequest) Method() string {
	return http.MethodPatch
}

func (ur *UpdateRequest) Path() string {
	return "/activeNumbers/" + ur.PhoneNumber
}

func (ur *UpdateRequest) QueryString() (string, error) {
	return "", nil
}

func (ur *UpdateRequest) ExpectedStatusCode() int {
	return http.StatusOK
}

func (ur *UpdateResponse) FromJSON(bytes []byte) error {
	return json.Unmarshal(bytes, ur)
}
