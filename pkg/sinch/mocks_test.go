package sinch

import "testing"

func Test_Mock_Implementations(t *testing.T) {
	var _ APIAction[*MockAPIRequest, *MockAPIResponse] = new(MockAPIAction)
	var _ APIRequest = new(MockAPIRequest)
	var _ APIResponse = new(MockAPIResponse)
}
