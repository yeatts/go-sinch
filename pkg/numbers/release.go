package numbers

import (
	"encoding/json"
	"net/http"

	"github.com/thezmc/go-sinch/pkg/sinch"
)

type Release struct {
	request  *ReleaseRequest
	response *ReleaseResponse
}

type ReleaseRequest struct {
	PhoneNumber string `url:"phoneNumber" json:"-"`
}

type ReleaseResponse struct {
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
	CallbackUrl           string                      `json:"callbackUrl"`
}

func (r *Release) IsNumbersAction() {}

func (r *Release) WithRequest(request *ReleaseRequest) *Release {
	r.request = request
	return r
}

func (r *Release) WithResponse(response *ReleaseResponse) *Release {
	r.response = response
	return r
}

func (r *Release) Request() *ReleaseRequest {
	return r.request
}

func (r *Release) Response() *ReleaseResponse {
	return r.response
}

func (rr *ReleaseRequest) WithPhoneNumber(phoneNumber string) *ReleaseRequest {
	rr.PhoneNumber = phoneNumber
	return rr
}

func (rr *ReleaseRequest) Validate() error {
	var errors sinch.Errors
	if rr.PhoneNumber == "" {
		errors = append(errors, PhoneNumberRequiredError)
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (rr *ReleaseRequest) Body() ([]byte, error) {
	return json.Marshal(rr)
}

func (rr *ReleaseRequest) Method() string {
	return http.MethodPost
}

func (rr *ReleaseRequest) Path() string {
	return "/activeNumbers/" + rr.PhoneNumber + ":release"
}

func (rr *ReleaseRequest) QueryString() (string, error) {
	return "", nil
}

func (rr *ReleaseRequest) ExpectedStatusCode() int {
	return http.StatusOK
}

func (rr *ReleaseResponse) FromJSON(bytes []byte) error {
	return json.Unmarshal(bytes, rr)
}
