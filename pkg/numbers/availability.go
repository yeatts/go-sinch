package numbers

import (
	"encoding/json"
	"net/http"

	"github.com/biter777/countries"
	"github.com/google/go-querystring/query"
	"github.com/thezmc/go-sinch/pkg/sinch"
	"go.uber.org/multierr"
)

type AvailabilityAction struct {
	request  *AvailabilityRequest
	response *AvailabilityResponse
}

func (aa *AvailabilityAction) IsNumbersAction() {}

func (aa *AvailabilityAction) Request() *AvailabilityRequest {
	return aa.request
}

func (aa *AvailabilityAction) Response() *AvailabilityResponse {
	return aa.response
}

type AvailabilityRequest struct {
	// Sequence of digits to search for. If you prefer or need certain digits in sequential order, you can enter the sequence of numbers here. For example, 2020.
	Pattern string `url:"numberPattern.pattern,omitempty"`
	// Search pattern to apply. The options are, START, CONTAIN, and END.
	//	START
	// Numbers that begin with the numberPattern.pattern entered. Often used to search for a specific area code. When using START, a plus sign (+) must be included and URL encoded, so %2B. For example, to search for area code 206 in the US, you would enter, %2b1206.
	// 	CONTAIN
	// The number pattern entered is contained somewhere in the number, the location being undefined.
	// 	END
	// The number ends with the number pattern entered.
	SearchPattern string   `url:"numberPattern.searchPattern,omitempty"`
	RegionCode    string   `url:"regionCode"`             // Region code to filter by. ISO 3166-1 alpha-2 country code of the phone number. Example: US, GB or SE.
	Type          string   `url:"type"`                   // Number type to filter by. Options include MOBILE, LOCAL or TOLL_FREE.
	Capabilities  []string `url:"capabilities,omitempty"` // Capabilities to filter by. Options include SMS or VOICE.
	Size          int      `url:"size,omitempty"`         // Number of available numbers to return.
}

type AvailabilityResponse struct {
	AvailableNumbers []AvailableNumber `json:"availableNumbers"`
}

type AvailableNumber struct {
	PhoneNumber string `json:"phoneNumber"` // The phone number in E.164 format with leading +. Example +12025550134.
	RegionCode  string `json:"regionCode"`  // ISO 3166-1 alpha-2 country code of the phone number. Example: US, UK or SE.
	// The number type.
	//	MOBILE
	// Numbers that belong to a specific range.
	//	LOCAL
	// Numbers that are assigned to a specific geographic region.
	//	TOLL_FREE
	// Numbers that are free of charge for the calling party but billed for all arriving calls.
	Type string `json:"type"`
	// The capabilities of the number.
	//	SMS
	// The number can receive/send SMS messages.
	//	VOICE
	// The number can receive voice calls.
	Capability                      []string `json:"capability"`
	SetupPrice                      Price    `json:"setupPrice"`
	MonthlyPrice                    Price    `json:"monthlyPrice"`
	PaymentIntervalMonths           int      `json:"paymentIntervalMonths"`
	SupportingDocumentationRequired bool     `json:"supportingDocumentationRequired"`
}

type Price struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

type SearchPattern int

const (
	SearchPatternStart SearchPattern = iota
	SearchPatternExact
	SearchPatternEnd
)

func (sp SearchPattern) String() string {
	return [...]string{"START", "EXACT", "END"}[sp]
}

type Type int

const (
	TypeLocal Type = iota
	TypeTollFree
	TypeMobile
)

func (t Type) String() string {
	return [...]string{"LOCAL", "TOLL_FREE", "MOBILE"}[t]
}

type Capability int

const (
	CapabilitySMS Capability = iota
	CapabilityVoice
)

func (c Capability) String() string {
	return [...]string{"SMS", "VOICE"}[c]
}

func (anr *AvailabilityRequest) WithPattern(pattern string) *AvailabilityRequest {
	anr.Pattern = pattern
	return anr
}

func (anr *AvailabilityRequest) WithSearchPattern(sp SearchPattern) *AvailabilityRequest {
	anr.SearchPattern = sp.String()
	return anr
}

func (anr *AvailabilityRequest) WithRegionCode(cc countries.CountryCode) *AvailabilityRequest {
	anr.RegionCode = cc.Alpha2()
	return anr
}

func (anr *AvailabilityRequest) WithType(t Type) *AvailabilityRequest {
	anr.Type = t.String()
	return anr
}

func (anr *AvailabilityRequest) WithCapability(c ...Capability) *AvailabilityRequest {
	for _, cap := range c {
		anr.Capabilities = append(anr.Capabilities, cap.String())
	}
	return anr
}

func (anr *AvailabilityRequest) WithSize(size int) *AvailabilityRequest {
	anr.Size = size
	return anr
}

func (ac *AvailabilityRequest) Validate() error {
	var errs sinch.Errors
	if ac.RegionCode == "" {
		errs = append(errs, RegionCodeRequiredError)
	}
	if ac.Type == "" {
		errs = append(errs, TypeRequiredError)
	}
	if len(errs) > 0 {
		return multierr.Combine(errs)
	}
	return nil
}

func (ac *AvailabilityRequest) ExpectedStatusCode() int {
	return http.StatusOK
}

func (ac *AvailabilityRequest) Method() string {
	return http.MethodGet
}

func (ac *AvailabilityRequest) Path() string {
	return "/availableNumbers"
}

func (ac *AvailabilityRequest) QueryString() (string, error) {
	v, err := query.Values(ac)
	if err != nil {
		return "", err
	}
	return "?" + v.Encode(), nil
}

func (ac *AvailabilityRequest) Body() ([]byte, error) {
	body, err := json.Marshal(ac)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (ar *AvailabilityResponse) FromJSON(bytes []byte) error {
	return json.Unmarshal(bytes, ar)
}
