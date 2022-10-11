mocks:
	mkdir -p mock
	rm -rf mock/*
	mockgen -package mock github.com/thezmc/go-sinch/pkg/sinch APIRequest,APIResponse
	mockgen -package mock net/http Client
