package main

import (
	"fmt"

	"github.com/thezmc/go-sinch/pkg/api"
	"github.com/thezmc/go-sinch/pkg/sms"
)

func main() {
	apiClient := new(api.Client)
	smsClient := new(sms.Client).WithPlanID("YOUR_PLAN_ID").WithAuthToken("YOUR_AUTH_TOKEN").WithSinchAPI(apiClient)

	request := new(sms.BatchSendRequest).
		To("RECIPIENT_PHONE_NUMBER").
		From("SENDER_PHONE_NUMBER").
		WithMessageBody("YOUR_MESSAGE_BODY")

	response := new(sms.BatchSendResponse)
	if err := smsClient.Do(request, response); err != nil {
		panic(err)
	}

	fmt.Printf("Send Response: %+v", response)
}
