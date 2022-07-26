package sms // import sinchsms "github.com/thezmc/go-sinch/sms"

import "net/http"

type SMSClient interface {
	US() SMSClient
	EU() SMSClient
	AU() SMSClient
	BR() SMSClient
	CA() SMSClient

	WithAuthToken(authToken string) SMSClient
	WithPlanID(planID string) SMSClient
	WithCustomBaseURL(baseURL string) SMSClient
	WithCustomHTTPClient(httpClient *http.Client) SMSClient
	Execute(req *http.Request, resourceName string) (*http.Response, error)

	NewBatchSender() SMSBatchSender
}

type SMSBatchSender interface {
	WithBody(body string) SMSBatchSender
	To(to ...string) SMSBatchSender
	From(from string) SMSBatchSender

	WithDeliveryReport(deliveryReport DeliveryReport) SMSBatchSender
	WithParameters(parameters map[string]map[string]string) SMSBatchSender
	WithParameter(parameterName string, valueMap map[string]string) SMSBatchSender
	WithCampaignID(campaignID string) SMSBatchSender
	SendingAt(sendAt string) SMSBatchSender
	ExpiringAt(expireAt string) SMSBatchSender
	WithCallbackURL(callbackURL string) SMSBatchSender
	WithClientReference(clientReference string) SMSBatchSender
	WithFeedbackEnabled() SMSBatchSender
	WithFlashMessageEnabled() SMSBatchSender
	WithTruncateConcatEnabled() SMSBatchSender
	WithMaxNumberOfMessageParts(maxNumberOfMessageParts int) SMSBatchSender
	WithTonOverride(fromTypeOfNumber int) SMSBatchSender
	WithNPIOverride(fromNumberPlanIndicator int) SMSBatchSender

	Send() SMSBatchSender
	Error() error
}
