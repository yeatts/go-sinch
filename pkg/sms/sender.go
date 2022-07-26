package sms // import sinchsms "github.com/thezmc/go-sinch/sms"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/thezmc/go-sinch/pkg/interfaces"
	"github.com/thezmc/go-sinch/pkg/models"
	"golang.org/x/exp/slices"
)

const (
	sendResourceName = "batches"
	iso8601          = "2006-01-02T15:04:05.000Z"
)

type batchSender struct {
	client   *client
	request  *models.SendRequest
	response *models.SendResponse
	err      error
}

// NewBatchSender returns a new Sender using this SMS client.
func (c *client) NewBatchSender() interfaces.SMSBatchSender {
	return NewBatchSender(c)
}

// NewBatchSender returns a new Sender using the provided SMS client.
func NewBatchSender(c *client) interfaces.SMSBatchSender {
	return &batchSender{
		client:  c,
		request: &models.SendRequest{},
	}
}

// FromRequest allows you to create a Sender from an existing SendRequest
func (s *batchSender) FromRequest(req *models.SendRequest) interfaces.SMSBatchSender {
	s.request = req
	return s
}

// WithBody sets the message body for the request.
func (s *batchSender) WithBody(body string) interfaces.SMSBatchSender {
	s.request.Body = body
	return s
}

// WithDeliveryReport sets the delivery report option for the request.
func (s *batchSender) WithDeliveryReport(deliveryReport models.DeliveryReport) interfaces.SMSBatchSender {
	s.request.DeliveryReport = deliveryReport
	return s
}

// To sets the recipient(s) for the request.
func (s *batchSender) To(to ...string) interfaces.SMSBatchSender {
	s.request.ToNumbers = append(s.request.ToNumbers, to...)
	return s
}

// From sets the sending number for the request.
func (s *batchSender) From(from string) interfaces.SMSBatchSender {
	s.request.FromNumber = from
	return s
}

// WithParameters sets the provided parameters for the request.
func (s *batchSender) WithParameters(parameters map[string]map[string]string) interfaces.SMSBatchSender {
	for k, v := range parameters {
		s.request.Parameters[k] = v
	}
	return s
}

// WithParameter sets a single parameter for the request.
func (s *batchSender) WithParameter(parameterName string, valueMap map[string]string) interfaces.SMSBatchSender {
	s.request.Parameters[parameterName] = valueMap
	return s
}

// WithCampaignID sets the campaign ID for the request.
func (s *batchSender) WithCampaignID(campaignID string) interfaces.SMSBatchSender {
	s.request.CampaignID = campaignID
	return s
}

// SendingAt sets the date and time to deliver the batch.
func (s *batchSender) SendingAt(sendAt string) interfaces.SMSBatchSender {
	s.request.SendAt = sendAt
	return s
}

// ExpiringAt sets the date and time to stop attempting to deliver a batch if failures occur.
func (s *batchSender) ExpiringAt(expireAt string) interfaces.SMSBatchSender {
	s.request.ExpireAt = expireAt
	return s
}

// WithCallbackURL sets the callback URL for the request.
func (s *batchSender) WithCallbackURL(callbackURL string) interfaces.SMSBatchSender {
	s.request.CallbackURL = callbackURL
	return s
}

// WithClientReference sets the client reference for the request.
func (s *batchSender) WithClientReference(clientReference string) interfaces.SMSBatchSender {
	s.request.ClientReference = clientReference
	return s
}

// WithFeedbackEnabled enables feedback for the request. By default this is false.
func (s *batchSender) WithFeedbackEnabled() interfaces.SMSBatchSender {
	s.request.FeedbackEnabled = true
	return s
}

// WithFlashMessageEnabled enables a flash message for the request. By default this is false.
func (s *batchSender) WithFlashMessageEnabled() interfaces.SMSBatchSender {
	s.request.FlashMessage = true
	return s
}

// WithTruncateConcatEnabled enables truncating the concatenated message for the request. This means if the message is longer
// than can be delivered by a single SMS, only the first SMS will be sent. By default this is false.
func (s *batchSender) WithTruncateConcatEnabled() interfaces.SMSBatchSender {
	s.request.TruncateConcat = true
	return s
}

// WithMaxNumberOfMessageParts sets the maximum number of message parts for the request.
func (s *batchSender) WithMaxNumberOfMessageParts(maxNumberOfMessageParts int) interfaces.SMSBatchSender {
	s.request.MaxNumberOfMessageParts = maxNumberOfMessageParts
	return s
}

// WithTonOverride overrides the type of number for the request. By default this is determined automatically. Only
// use this option if you know what you're doing.
func (s *batchSender) WithTonOverride(fromTypeOfNumber int) interfaces.SMSBatchSender {
	s.request.FromTypeOfNumber = fromTypeOfNumber
	return s
}

// WithNPIOverride overrides the type of Number Plan Indicator for the request. By default this is determined
// automatically. Only use this option if you know what you're doing.
func (s *batchSender) WithNPIOverride(fromNumberPlanIndicator int) interfaces.SMSBatchSender {
	s.request.FromNumberPlanIndicator = fromNumberPlanIndicator
	return s
}

// validate makes sure all request parameters are set within the limits specified by the Sinch API documentation.
// This should be called before sending the request to keep performance to a maximum and API errors to a minimum.
//
// Ref: https://developers.sinch.com/docs/sms/api-reference/sms/tag/Batches/#tag/Batches/operation/SendSMS
func (s *batchSender) validate() error {
	if len(s.request.ToNumbers) == 0 || len(s.request.ToNumbers) > 0 && slices.Contains(s.request.ToNumbers, "") || len(s.request.ToNumbers) > 1000 {
		return InvalidToNumberError
	}
	if s.request.FromNumber == "" {
		return InvalidFromNumberError
	}
	if s.request.FromTypeOfNumber < 0 || s.request.FromTypeOfNumber > 6 {
		return InvalidTypeOfNumberError
	}
	if s.request.FromNumberPlanIndicator < 0 || s.request.FromNumberPlanIndicator > 18 {
		return InvalidNPIError
	}
	if s.request.Body == "" || len(s.request.Body) > 2000 {
		return InvalidBodyError
	}
	if s.request.CallbackURL != "" && !strings.HasPrefix(s.request.CallbackURL, "http") || len(s.request.CallbackURL) > 2048 {
		return InvalidCallbackURLError
	}
	if s.request.SendAt != "" {
		if _, err := time.Parse(iso8601, s.request.SendAt); err != nil {
			return InvalidSendAtError
		}
	}
	if s.request.ExpireAt != "" {
		if _, err := time.Parse(iso8601, s.request.ExpireAt); err != nil {
			return InvalidExpireAtError
		}
	}
	return nil
}

// toIOReader converts the SMSSendRequest to json and encodes it into a buffer. Returns an io.Reader with the encoded
// data.
func (s *batchSender) toIOReader() (io.Reader, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(s.request); err != nil {
		return nil, err
	}
	return buffer, nil
}

// toRequest converts the SMSSendRequest to an http.Request.
func (s *batchSender) toRequest() (*http.Request, error) {
	json, err := s.toIOReader()
	if err != nil {
		return nil, err
	}
	sendURL := fmt.Sprintf("%s/%s/%s", s.client.baseURL, s.client.planID, sendResourceName)
	req, err := http.NewRequest(http.MethodPost, sendURL, json)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Send sends the request to the Sinch API and returns the response.
func (s *batchSender) Send() interfaces.SMSBatchSender {
	if err := s.validate(); err != nil {
		s.err = err
		return s
	}
	req, err := s.toRequest()
	if err != nil {
		s.err = err
		return s
	}
	resp, err := s.client.Execute(req, sendResourceName)
	if err != nil {
		s.err = err
		return s
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		s.err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return s
	}
	var response *models.SendResponse
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		s.err = err
		return s
	}
	s.response = response
	return nil
}

// Request returns the SMSSendRequest used to send the SMS.
func (s *batchSender) Request() *models.SendRequest {
	return s.request
}

// Error returns the error that occurred during the request or nil if no error occurred.
func (s *batchSender) Error() error {
	return s.err
}

// Response returns the response from the Sinch API or nil if no response was received.
func (s *batchSender) Response() *models.SendResponse {
	return s.response
}
