package main

import (
	"fmt"

	"github.com/thezmc/go-sinch/pkg/sms"
)

func main() {
	client := sms.NewClient().US().WithAuthToken("authToken").WithPlanID("planID")

	req := client.NewSendRequest().
		To("+12345678901").
		From("+12345678901").
		WithBody("Hello World!")

	if _, err := req.Execute(); err != nil {
		fmt.Println(err)
	}
}
