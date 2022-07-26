package sms // import sinchsms "github.com/thezmc/go-sinch/sms"

// SendRequest represents the request body for sending SMS messages.
//
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

type DeliveryReport string

const (
	None         DeliveryReport = "none"
	Summary      DeliveryReport = "summary"
	Full         DeliveryReport = "full"
	PerRecipient DeliveryReport = "per_recipient"
)

type SendResponse struct {
	SendRequest
	ID         string `json:"id"`          // Unique identifier for batch
	Canceled   bool   `json:"canceled"`    // Indicates if the batch has been canceled or not
	CreatedAt  string `json:"created_at"`  // Timestamp for when batch was created. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
	ModifiedAt string `json:"modified_at"` // Timestamp for when batch was last updated. Formatted as ISO-8601: YYYY-MM-DDThh:mm:ss.SSSZ
}
