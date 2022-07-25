package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

const (
	SendResourceName = "batches"
	ISO8601          = "2006-01-02T15:04:05.000Z"
)

// SMSSendRequest represents the request body for sending SMS messages.
//
// Ref: https://developers.sinch.com/docs/sms/api-reference/sms/tag/Batches/#tag/Batches/operation/SendSMS
type SMSSendRequest struct {
	Client                  *Client
	Body                    string                       `json:"body"`                                  // The message content
	DeliveryReport          DeliveryReport               `json:"delivery_report"`                       // Request delivery report callback. Note that delivery reports can be fetched from the API regardless of this setting.
	ToNumbers               []string                     `json:"to"`                                    // List of Phone numbers and group IDs that will receive the batch.
	FromNumber              string                       `json:"from,omitempty"`                        // Sender number. Must be valid phone number, short code or alphanumeric. Required if Automatic Default Originator not configured.
	Parameters              map[string]map[string]string `json:"parameters,omitempty"`                  // Contains the parameters that will be used for customizing the message for each recipient. Ref: https://developers.sinch.com/docs/sms/resources/message-info/message-parameterization/
	CampaignID              string                       `json:"campaign_id,omitempty"`                 // The campaign and service IDs this message belongs to. You generate your own campaign and service IDs. US only. Ref: https://dashboard.sinch.com/sms/us-campaigns
	SendAt                  string                       `json:"send_at,omitempty"`                     // If set in the future, the message will be delayed until send_at occurs. Must be before expire_at. If set in the past, messages will be sent immediately. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
	ExpireAt                string                       `json:"expire_at,omitempty"`                   // If set, the system will stop trying to deliver the message at this point. Must be after send_at. Default and max is 3 days after send_at. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
	CallbackURL             string                       `json:"callback_url,omitempty"`                // Override the default callback URL for this batch. Must be valid URL.
	ClientReference         string                       `json:"client_reference,omitempty"`            // The client identifier of a batch message. If set, the identifier will be added in the delivery report/callback of this batch
	FeedbackEnabled         bool                         `json:"feedback_enabled,omitempty"`            // If set to true, then feedback is expected after successful delivery. Ref: https://developers.sinch.com/docs/sms/api-reference/sms/tag/Batches/#tag/Batches/operation/deliveryFeedback
	FlashMessage            bool                         `json:"flash_message,omitempty"`               // Shows message on screen without user interaction while not saving the message to the inbox.
	TruncateConcat          bool                         `json:"truncate_concat,omitempty"`             // If set to true the message will be shortened when exceeding one part.
	MaxNumberOfMessageParts int                          `json:"max_number_of_message_parts,omitempty"` // Message will be dispatched only if it is not split to more parts than Max Number of Message Parts
	FromTypeOfNumber        int                          `json:"from_ton,omitempty"`                    // The type of number for the sender number. Use to override the automatic detection.
	FromNumberPlanIndicator int                          `json:"from_npi,omitempty"`                    // Number Plan Indicator for the sender number. Use to override the automatic detection.
}

type DeliveryReport string

const (
	None         DeliveryReport = "none"
	Summary      DeliveryReport = "summary"
	Full         DeliveryReport = "full"
	PerRecipient DeliveryReport = "per_recipient"
)

type SMSSendResponse struct {
	SMSSendRequest
	ID         string `json:"id"`          // Unique identifier for batch
	Canceled   bool   `json:"canceled"`    // Indicates if the batch has been canceled or not
	CreatedAt  string `json:"created_at"`  // Timestamp for when batch was created. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
	ModifiedAt string `json:"modified_at"` // Timestamp for when batch was last updated. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
}

type SMSListBatchesRequest struct {
	Page            int    `url:"page,omitempty"`             // The page number starting from 0.
	PageSize        int    `url:"page_size,omitempty"`        // integer <= 100. Default: 30. Determines the size of a page. Example: page_size=50
	From            string `url:"from,omitempty"`             // Only list messages sent from this sender number. Multiple originating numbers can be comma separated. Must be phone numbers or short code.
	StartDate       string `url:"start_date,omitempty"`       // Only list messages received at or after this date/time. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ.
	EndDate         string `url:"end_date,omitempty"`         // Only list messages received before this date/time. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ.
	ClientReference string `url:"client_reference,omitempty"` // Client reference to include
}

// NewSendRequest returns a new SMSSendRequest using this SMS client.
func (c *Client) NewSendRequest() *SMSSendRequest {
	return NewSendRequest(c)
}

// NewSendRequest returns a new SMSSendRequest using the provided SMS client.
func NewSendRequest(client *Client) *SMSSendRequest {
	return &SMSSendRequest{
		Client: client,
	}
}

// WithBody sets the message body for the request.
func (r *SMSSendRequest) WithBody(body string) *SMSSendRequest {
	r.Body = body
	return r
}

// WithDeliveryReport sets the delivery report option for the request.
func (r *SMSSendRequest) WithDeliveryReport(deliveryReport DeliveryReport) *SMSSendRequest {
	r.DeliveryReport = deliveryReport
	return r
}

// To sets the recipient(s) for the request.
func (r *SMSSendRequest) To(to ...string) *SMSSendRequest {
	r.ToNumbers = append(r.ToNumbers, to...)
	return r
}

// From sets the sending number for the request.
func (r *SMSSendRequest) From(from string) *SMSSendRequest {
	r.FromNumber = from
	return r
}

// WithParameters sets the provided parameters for the request.
func (r *SMSSendRequest) WithParameters(parameters map[string]map[string]string) *SMSSendRequest {
	for k, v := range parameters {
		r.Parameters[k] = v
	}
	return r
}

// WithParameter sets a single parameter for the request.
func (r *SMSSendRequest) WithParameter(parameterName string, valueMap map[string]string) *SMSSendRequest {
	r.Parameters[parameterName] = valueMap
	return r
}

// WithCampaignID sets the campaign ID for the request.
func (r *SMSSendRequest) WithCampaignID(campaignID string) *SMSSendRequest {
	r.CampaignID = campaignID
	return r
}

// SendingAt sets the date and time to deliver the batch.
func (r *SMSSendRequest) SendingAt(sendAt string) *SMSSendRequest {
	r.SendAt = sendAt
	return r
}

// ExpiringAt sets the date and time to stop attempting to deliver a batch if failures occur.
func (r *SMSSendRequest) ExpiringAt(expireAt string) *SMSSendRequest {
	r.ExpireAt = expireAt
	return r
}

// WithCallbackURL sets the callback URL for the request.
func (r *SMSSendRequest) WithCallbackURL(callbackURL string) *SMSSendRequest {
	r.CallbackURL = callbackURL
	return r
}

// WithClientReference sets the client reference for the request.
func (r *SMSSendRequest) WithClientReference(clientReference string) *SMSSendRequest {
	r.ClientReference = clientReference
	return r
}

// WithFeedbackEnabled enables feedback for the request. By default this is false.
func (r *SMSSendRequest) WithFeedbackEnabled() *SMSSendRequest {
	r.FeedbackEnabled = true
	return r
}

// WithFlashMessageEnabled enables a flash message for the request. By default this is false.
func (r *SMSSendRequest) WithFlashMessageEnabled() *SMSSendRequest {
	r.FlashMessage = true
	return r
}

// WithTruncateConcatEnabled enables truncating the concatenated message for the request. This means if the message is longer
// than can be delivered by a single SMS, only the first SMS will be sent. By default this is false.
func (r *SMSSendRequest) WithTruncateConcatEnabled() *SMSSendRequest {
	r.TruncateConcat = true
	return r
}

// WithMaxNumberOfMessageParts sets the maximum number of message parts for the request.
func (r *SMSSendRequest) WithMaxNumberOfMessageParts(maxNumberOfMessageParts int) *SMSSendRequest {
	r.MaxNumberOfMessageParts = maxNumberOfMessageParts
	return r
}

// WithTonOverride overrides the type of number for the request. By default this is determined automatically. Only
// use this option if you know what you're doing.
func (r *SMSSendRequest) WithTonOverride(fromTypeOfNumber int) *SMSSendRequest {
	r.FromTypeOfNumber = fromTypeOfNumber
	return r
}

// WithNPIOverride overrides the type of Number Plan Indicator for the request. By default this is determined
// automatically. Only use this option if you know what you're doing.
func (r *SMSSendRequest) WithNPIOverride(fromNumberPlanIndicator int) *SMSSendRequest {
	r.FromNumberPlanIndicator = fromNumberPlanIndicator
	return r
}

// Validate makes sure all request parameters are set within the limits specified by the Sinch API documentation.
// this should be called before sending the request to keep performance to a maximum and API errors to a minimum.
//
// Ref: https://developers.sinch.com/docs/sms/api-reference/sms/tag/Batches/#tag/Batches/operation/SendSMS
func (r *SMSSendRequest) Validate() error {
	if len(r.ToNumbers) == 0 || len(r.ToNumbers) > 0 && slices.Contains(r.ToNumbers, "") || len(r.ToNumbers) > 1000 {
		return InvalidToNumberError
	}
	if r.FromNumber == "" {
		return InvalidFromNumberError
	}
	if r.FromTypeOfNumber < 0 || r.FromTypeOfNumber > 6 {
		return InvalidTypeOfNumberError
	}
	if r.FromNumberPlanIndicator < 0 || r.FromNumberPlanIndicator > 18 {
		return InvalidNPIError
	}
	if r.Body == "" || len(r.Body) > 2000 {
		return InvalidBodyError
	}
	if r.CallbackURL != "" && !strings.HasPrefix(r.CallbackURL, "http") || len(r.CallbackURL) > 2048 {
		return InvalidCallbackURLError
	}
	if r.SendAt != "" {
		if _, err := time.Parse(ISO8601, r.SendAt); err != nil {
			return InvalidSendAtError
		}
	}
	if r.ExpireAt != "" {
		if _, err := time.Parse(ISO8601, r.ExpireAt); err != nil {
			return InvalidExpireAtError
		}
	}
	return nil
}

// toIOReader converts the SMSSendRequest to json and encodes it to an io.Reader.
func (r *SMSSendRequest) toIOReader() (io.Reader, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(r); err != nil {
		return nil, err
	}
	return buffer, nil
}

// toRequest converts the SMSSendRequest to an http.Request.
func (r *SMSSendRequest) toRequest() (*http.Request, error) {
	json, err := r.toIOReader()
	if err != nil {
		return nil, err
	}
	sendURL := fmt.Sprintf("%s/%s/%s", r.Client.baseURL, r.Client.planID, SendResourceName)
	req, err := http.NewRequest(http.MethodPost, sendURL, json)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Execute sends the request to the Sinch API and returns the response.
func (r *SMSSendRequest) Execute() (*SMSSendResponse, error) {
	req, err := r.toRequest()

	if err != nil {
		return nil, err
	}
	resp, err := r.Client.Execute(req, SendResourceName)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response SMSSendResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
