tidy: 
	go mod tidy
	go mod vendor
test:
	go test ./... -count=1
	# staticcheck -checks=all ./...
test-verbose:
	go test ./... -count=1 -v
	# staticcheck -checks=all ./...
mocks:
	mockgen -source=business/usecase/user/interface.go -destination=business/usecase/user/mock/mock.go -package=mock