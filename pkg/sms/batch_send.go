package sms

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/thezmc/go-sinch/pkg/sinch"
	"golang.org/x/exp/slices"
)

type BatchSend struct {
	request  *BatchSendRequest
	response *BatchSendResponse
}

func (bs *BatchSend) Request() *BatchSendRequest {
	return bs.request
}

func (bs *BatchSend) Response() *BatchSendResponse {
	return bs.response
}

type BatchSendRequest struct {
	MessageBody             string                       `json:"body"`                                  // The message content
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

type BatchSendResponse struct {
	BatchSendRequest
	ID         string `json:"id"`          // Unique identifier for batch
	Canceled   bool   `json:"canceled"`    // Indicates if the batch has been canceled or not
	CreatedAt  string `json:"created_at"`  // Timestamp for when batch was created. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
	ModifiedAt string `json:"modified_at"` // Timestamp for when batch was last updated. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
}

// WithBody sets the message body for the request.
func (bsr *BatchSendRequest) WithMessageBody(body string) *BatchSendRequest {
	bsr.MessageBody = body
	return bsr
}

// WithDeliveryReport sets the delivery report option for the request.
func (bsr *BatchSendRequest) WithDeliveryReport(deliveryReport DeliveryReport) *BatchSendRequest {
	bsr.DeliveryReport = deliveryReport
	return bsr
}

// To sets the recipient(s) for the request.
func (bsr *BatchSendRequest) To(to ...string) *BatchSendRequest {
	bsr.ToNumbers = append(bsr.ToNumbers, to...)
	return bsr
}

// From sets the sending number for the request.
func (bsr *BatchSendRequest) From(from string) *BatchSendRequest {
	bsr.FromNumber = from
	return bsr
}

// WithParameters sets the provided parameters for the request.
func (bsr *BatchSendRequest) WithParameters(parameters map[string]map[string]string) *BatchSendRequest {
	for k, v := range parameters {
		bsr.Parameters[k] = v
	}
	return bsr
}

// WithParameter sets a single parameter for the request.
func (bsr *BatchSendRequest) WithParameter(parameterName string, valueMap map[string]string) *BatchSendRequest {
	bsr.Parameters[parameterName] = valueMap
	return bsr
}

// WithCampaignID sets the campaign ID for the request.
func (bsr *BatchSendRequest) WithCampaignID(campaignID string) *BatchSendRequest {
	bsr.CampaignID = campaignID
	return bsr
}

// SendingAt sets the date and time to deliver the batch.
func (bsr *BatchSendRequest) SendingAt(sendAt string) *BatchSendRequest {
	bsr.SendAt = sendAt
	return bsr
}

// ExpiringAt sets the date and time to stop attempting to deliver a batch if failures occur.
func (bsr *BatchSendRequest) ExpiringAt(expireAt string) *BatchSendRequest {
	bsr.ExpireAt = expireAt
	return bsr
}

// WithCallbackURL sets the callback URL for the request.
func (bsr *BatchSendRequest) WithCallbackURL(callbackURL string) *BatchSendRequest {
	bsr.CallbackURL = callbackURL
	return bsr
}

// WithClientReference sets the client reference for the request.
func (bsr *BatchSendRequest) WithClientReference(clientReference string) *BatchSendRequest {
	bsr.ClientReference = clientReference
	return bsr
}

// WithFeedbackEnabled enables feedback for the request. By default this is false.
func (bsr *BatchSendRequest) WithFeedbackEnabled() *BatchSendRequest {
	bsr.FeedbackEnabled = true
	return bsr
}

// WithFlashMessageEnabled enables a flash message for the request. By default this is false.
func (bsr *BatchSendRequest) WithFlashMessageEnabled() *BatchSendRequest {
	bsr.FlashMessage = true
	return bsr
}

// WithTruncateConcatEnabled enables truncating the concatenated message for the request. This means if the message is longer
// than can be delivered by a single SMS, only the first SMS will be sent. By default this is false.
func (bsr *BatchSendRequest) WithTruncateConcatEnabled() *BatchSendRequest {
	bsr.TruncateConcat = true
	return bsr
}

// WithMaxNumberOfMessageParts sets the maximum number of message parts for the request.
func (bsr *BatchSendRequest) WithMaxNumberOfMessageParts(maxNumberOfMessageParts int) *BatchSendRequest {
	bsr.MaxNumberOfMessageParts = maxNumberOfMessageParts
	return bsr
}

// WithTonOverride overrides the type of number for the request. By default this is determined automatically. Only
// use this option if you know what you're doing.
func (bsr *BatchSendRequest) WithTonOverride(fromTypeOfNumber int) *BatchSendRequest {
	bsr.FromTypeOfNumber = fromTypeOfNumber
	return bsr
}

// WithNPIOverride overrides the type of Number Plan Indicator for the request. By default this is determined
// automatically. Only use this option if you know what you're doing.
func (bsr *BatchSendRequest) WithNPIOverride(fromNumberPlanIndicator int) *BatchSendRequest {
	bsr.FromNumberPlanIndicator = fromNumberPlanIndicator
	return bsr
}

// validate makes sure all request parameters are set within the limits specified by the Sinch API documentation.
// This should be called before sending the request to keep performance to a maximum and API errors to a minimum.
//
// Ref: https://developers.sinch.com/docs/sms/api-reference/sms/tag/Batches/#tag/Batches/operation/SendSMS
func (bsr *BatchSendRequest) Validate() error {
	var errors sinch.Errors
	if len(bsr.ToNumbers) == 0 || len(bsr.ToNumbers) > 0 && slices.Contains(bsr.ToNumbers, "") || len(bsr.ToNumbers) > 1000 {
		errors = append(errors, InvalidToNumberError)
	}
	if bsr.FromNumber == "" {
		errors = append(errors, InvalidFromNumberError)
	}
	if bsr.FromTypeOfNumber < 0 || bsr.FromTypeOfNumber > 6 {
		errors = append(errors, InvalidTypeOfNumberError)
	}
	if bsr.FromNumberPlanIndicator < 0 || bsr.FromNumberPlanIndicator > 18 {
		errors = append(errors, InvalidNPIError)
	}
	if bsr.MessageBody == "" || len(bsr.MessageBody) > 2000 {
		errors = append(errors, InvalidBodyError)
	}
	if bsr.CallbackURL != "" && !strings.HasPrefix(bsr.CallbackURL, "http") || len(bsr.CallbackURL) > 2048 {
		errors = append(errors, InvalidCallbackURLError)
	}
	if bsr.SendAt != "" {
		if _, err := time.Parse(TimeFormat, bsr.SendAt); err != nil {
			errors = append(errors, InvalidSendAtError)
		}
	}
	if bsr.ExpireAt != "" {
		if _, err := time.Parse(TimeFormat, bsr.ExpireAt); err != nil {
			errors = append(errors, InvalidExpireAtError)
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (bsr *BatchSendRequest) ExpectedStatusCode() int {
	return http.StatusCreated
}

func (bsr *BatchSendRequest) Method() string {
	return http.MethodPost
}

func (bsr *BatchSendRequest) QueryString() (string, error) {
	return "", nil
}

func (bsr *BatchSendRequest) Body() ([]byte, error) {
	return json.Marshal(bsr)
}

func (bsr *BatchSendRequest) Path() string {
	return "/batches"
}

func (bsr *BatchSendResponse) FromJSON(data []byte) error {
	return json.Unmarshal(data, bsr)
}
