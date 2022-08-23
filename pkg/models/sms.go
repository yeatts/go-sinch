package models // import sinchsms "github.com/thezmc/go-sinch/sms"

import "encoding/json" // SendRequest represents the request body for sending SMS messages.
// Ref: https://developers.sinch.com/docs/sms/api-reference/sms/tag/Batches/#tag/Batches/operation/SendSMS
type SendRequest struct {
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

type DeliveryReport int

const (
	None DeliveryReport = iota
	Summary
	Full
	PerRecipient
)

func (dr DeliveryReport) String() string {
	return dr.toString()
}

func (dr DeliveryReport) toString() string {
	switch dr {
	case None:
		return "none"
	case Summary:
		return "summary"
	case Full:
		return "full"
	case PerRecipient:
		return "per_recipient"
	}
	return "none"
}

func toDR(s string) DeliveryReport {
	switch s {
	case "none":
		return None
	case "summary":
		return Summary
	case "full":
		return Full
	case "per_recipient":
		return PerRecipient
	}
	return None
}

func (dr DeliveryReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(dr.String())
}

func (dr *DeliveryReport) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*dr = toDR(s)
	return nil
}

type Type int

const (
	Text Type = iota
	Binary
)

func (t Type) String() string {
	switch t {
	case Text:
		return "mt_text"
	case Binary:
		return "mt_binary"
	}
	return "mt_text"
}

type SendResponse struct {
	SendRequest
	ID         string `json:"id"`          // Unique identifier for batch
	Canceled   bool   `json:"canceled"`    // Indicates if the batch has been canceled or not
	CreatedAt  string `json:"created_at"`  // Timestamp for when batch was created. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
	ModifiedAt string `json:"modified_at"` // Timestamp for when batch was last updated. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
}

type ListRequest struct {
	Page            int      `url:"page"`             // The page number of the results to be returned. Default is 0.
	PageSize        int      `url:"page_size"`        // The number of results to be returned per page. Default is 30.
	From            []string `url:"from"`             // List of phone numbers from where the batch was sent.
	StartDate       string   `url:"start_date"`       // Only list batches sent after this date. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
	EndDate         string   `url:"end_date"`         // Only list batches sent before this date. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
	ClientReference string   `url:"client_reference"` // Only list batches with this client reference.
}

type ListResponse struct {
	Page     int     `json:"page"`      // The page number of the results being returned.
	Count    int     `json:"count"`     // The number of results in the full data set.
	PageSize int     `json:"page_size"` // The number of results per page.
	Batches  []Batch `json:"batches"`   // The list of batches.
}

type Batch struct {
	To                      []string                     `json:"to"`                                    // List of phone numbers to which the batch was sent.
	Body                    string                       `json:"body"`                                  // The message content.
	From                    string                       `json:"from,omitempty"`                        // Sender phone number. Could also be a Sender ID.
	Type                    Type                         `json:"type,omitempty"`                        // The type of number for the sender number.
	UDH                     string                       `json:"udh,omitempty"`                         // The UDH header for a binary message. Max 140 bytes with the body.
	DeliveryReport          DeliveryReport               `json:"delivery_report,omitempty"`             // The delivery report type.
	SendAt                  string                       `json:"send_at,omitempty"`                     // The timestamp for when the batch was/will be sent. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
	ExpireAt                string                       `json:"expire_at,omitempty"`                   // The timestamp for when the batch will expire. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
	CallbackURL             string                       `json:"callback_url,omitempty"`                // The callback URL for the batch.
	FlashMessage            bool                         `json:"flash_message,omitempty"`               // Indicates if the message is a flash message.
	Parameters              map[string]map[string]string `json:"parameters,omitempty"`                  // The parameters for the batch.
	ClientReference         string                       `json:"client_reference,omitempty"`            // The client identifier of a batch message.
	MaxNumberOfMessageParts int                          `json:"max_number_of_message_parts,omitempty"` // The maximum number of message parts for the batch.
}
