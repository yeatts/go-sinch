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
	"fmt"

	sinchsms "github.com/theZMC/go-sinch/pkg/sms"
)

func main() {
	sms := sinchsms.New().US().WithAuthToken("authToken").WithPlanID("planID")

	req := sms.NewSendRequest().
		To("+12345678901").
		From("+12345678901").
		WithBody("Hello World!")

	if _, err := req.Execute(); err != nil {
		fmt.Println(err)
	}
}
```