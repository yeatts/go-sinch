package numbers

import (
	"testing"

	"github.com/thezmc/go-sinch/pkg/sinch"
)

func Test_Availability_Implementations(t *testing.T) {
	var _ sinch.Action[*AvailabilityRequest, *AvailabilityResponse] = new(AvailabilityAction)
	var _ sinch.APIRequest = &AvailabilityRequest{}
	var _ sinch.APIResponse = &AvailabilityResponse{}
}
