package sms

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	NoAuthTokenError            = Error("an auth token is required")
	NoPlanIDError               = Error("a plan ID is required")
	InvalidToNumberError        = Error("at least one to_number is required and no more than 1000 are allowed")
	InvalidFromNumberError      = Error("a from_number is required")
	InvalidTypeOfNumberError    = Error("type_of_number must be an int in the range 0-6")
	InvalidNPIError             = Error("npi must be an int in the range 0-18")
	InvalidBodyError            = Error("body must be between 0 and 2000 characters long")
	InvalidCallbackURLError     = Error("callback_url must start with http and be between 0 and 2048 characters long")
	InvalidClientReferenceError = Error("client_reference must be between 0 and 255 characters long")
	InvalidSendAtError          = Error("send_at must be in ISO-8601 format")
	InvalidExpireAtError        = Error("expire_at must be in ISO-8601 format")
	InvalidMaxMessagePartsError = Error("max_number_of_message_parts must be greater than 0")
)
