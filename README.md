# go-sinch
`go-sinch` is an API client for the [Sinch APIs](https://developers.sinch.com/) written in Go(Lang).

## Quickstart
Download the package
```shell
go get github.com/theZMC/go-sinch
```

### SMS
Send a message using the SMS client
```go
package main

import (
	"log"

	sinchsms "github.com/thezmc/go-sinch/pkg/sms"
)

func main() {
	client := sinchsms.NewClient().
		WithAuthToken("authToken").
		WithPlanID("planID")

	req := client.NewBatchSender().
		To("+12345678901").
		From("+12345678901").
		WithBody("Hello World!")

	if err := req.Send().Error(); err != nil {
		log.Fatalf("Error sending SMS message: %v", err)
	}
}
```